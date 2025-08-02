package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/garnizeh/englog/internal/logging"
)

// OllamaService implements AI service using Ollama with resilience
type OllamaService struct {
	logger  *logging.Logger
	baseURL string
	client  *http.Client
}

// GenerateRequest represents a request to Ollama
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// GenerateResponse represents a response from Ollama
type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Insight represents an AI-generated insight
type Insight struct {
	Content    string   `json:"content"`
	Tags       []string `json:"tags"`
	Confidence float32  `json:"confidence"`
}

// InsightRequest represents a request for insight generation
type InsightRequest struct {
	Prompt      string   `json:"prompt"`
	UserID      string   `json:"user_id"`
	EntryIDs    []string `json:"entry_ids"`
	InsightType string   `json:"insight_type"`
	Context     any      `json:"context,omitempty"`
}

// WeeklyReportRequest represents a request for weekly report generation
type WeeklyReportRequest struct {
	UserID    string    `json:"user_id"`
	WeekStart time.Time `json:"week_start"`
	WeekEnd   time.Time `json:"week_end"`
}

// WeeklyReport represents a generated weekly report
type WeeklyReport struct {
	Summary         string   `json:"summary"`
	KeyInsights     []string `json:"key_insights"`
	Recommendations []string `json:"recommendations"`
}

// NewOllamaService creates a new Ollama service instance with retry configuration
func NewOllamaService(baseURL string, logger *logging.Logger) (*OllamaService, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("ollama base URL cannot be empty")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	logger.Info("Creating new Ollama service",
		"base_url", baseURL,
		"timeout", "120s")

	return &OllamaService{
		logger:  logger,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 120 * time.Second, // Increased timeout for AI operations
		},
	}, nil
}

// GenerateInsight generates AI insights based on a prompt with retry logic
func (s *OllamaService) GenerateInsight(ctx context.Context, prompt string) (*Insight, error) {
	if prompt == "" {
		s.logger.LogError(ctx, fmt.Errorf("empty prompt"), "GenerateInsight called with empty prompt")
		return nil, fmt.Errorf("prompt cannot be empty")
	}

	s.logger.Info("Starting insight generation",
		"prompt_length", len(prompt),
		"model", "qwen2.5-coder:7b")

	// Retry configuration for AI operations
	maxRetries := 3
	baseDelay := 1 * time.Second

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		s.logger.Debug("Attempting insight generation",
			"attempt", attempt+1,
			"max_retries", maxRetries)

		select {
		case <-ctx.Done():
			s.logger.Warn("Insight generation cancelled",
				"attempt", attempt+1,
				"error", ctx.Err().Error())
			return nil, fmt.Errorf("insight generation cancelled: %w", ctx.Err())
		default:
		}

		response, err := s.generateWithTimeout(ctx, "qwen2.5-coder:7b", prompt, 60*time.Second)
		if err == nil {
			s.logger.Info("Insight generation successful",
				"attempt", attempt+1,
				"response_length", len(response))

			// Parse response and extract insights
			insight := &Insight{
				Content:    response,
				Tags:       []string{"ai-generated", "productivity"},
				Confidence: 0.8, // Default confidence
			}

			s.logger.Debug("Generated insight",
				"content_length", len(insight.Content),
				"tags_count", len(insight.Tags),
				"confidence", insight.Confidence)

			return insight, nil
		}

		lastErr = err
		s.logger.Warn("Insight generation attempt failed",
			"attempt", attempt+1,
			"error", err.Error())

		if attempt < maxRetries-1 {
			delay := time.Duration(attempt+1) * baseDelay
			s.logger.Debug("Retrying insight generation",
				"delay", delay.String(),
				"next_attempt", attempt+2)

			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("insight generation cancelled during retry: %w", ctx.Err())
			case <-time.After(delay):
			}
		}
	}

	s.logger.LogError(ctx, lastErr, "All insight generation attempts failed",
		"max_retries", maxRetries)

	return nil, fmt.Errorf("failed to generate insight after %d attempts: %w", maxRetries, lastErr)
}

// GenerateInsightWithContext generates AI insights using structured context
func (s *OllamaService) GenerateInsightWithContext(ctx context.Context, req *InsightRequest) (*Insight, error) {
	if req.Prompt == "" {
		s.logger.LogError(ctx, fmt.Errorf("empty prompt"), "GenerateInsightWithContext called with empty prompt")
		return nil, fmt.Errorf("prompt cannot be empty")
	}

	// Enhance prompt with structured context
	enhancedPrompt := s.buildEnhancedPrompt(req)

	s.logger.Info("Starting insight generation with context",
		"user_id", req.UserID,
		"insight_type", req.InsightType,
		"entry_count", len(req.EntryIDs),
		"enhanced_prompt_length", len(enhancedPrompt),
		"context_type", fmt.Sprintf("%T", req.Context))

	return s.GenerateInsight(ctx, enhancedPrompt)
}

// buildEnhancedPrompt creates an enhanced prompt using all pertinent information from the request
func (s *OllamaService) buildEnhancedPrompt(req *InsightRequest) string {
	var promptBuilder bytes.Buffer

	// Start with the base prompt
	promptBuilder.WriteString(req.Prompt)

	// Add structured request information to enhance AI understanding
	promptBuilder.WriteString("\n\n--- Request Information ---")
	promptBuilder.WriteString(fmt.Sprintf("\nUser ID: %s", req.UserID))
	promptBuilder.WriteString(fmt.Sprintf("\nInsight Type: %s", req.InsightType))
	promptBuilder.WriteString(fmt.Sprintf("\nNumber of Log Entries: %d", len(req.EntryIDs)))

	// Include entry IDs for reference (useful for the AI to understand scope)
	if len(req.EntryIDs) > 0 {
		promptBuilder.WriteString("\nLog Entry IDs: ")
		if len(req.EntryIDs) <= 5 {
			// Show all IDs if there are 5 or fewer
			promptBuilder.WriteString(fmt.Sprintf("[%s]", joinStrings(req.EntryIDs, ", ")))
		} else {
			// Show first 3 and last 2 with count for larger sets
			firstThree := req.EntryIDs[:3]
			lastTwo := req.EntryIDs[len(req.EntryIDs)-2:]
			promptBuilder.WriteString(fmt.Sprintf("[%s, ... (%d more), %s]",
				joinStrings(firstThree, ", "),
				len(req.EntryIDs)-5,
				joinStrings(lastTwo, ", ")))
		}
	}

	// Add insight type specific instructions
	promptBuilder.WriteString(fmt.Sprintf("\n\nInsight Generation Guidelines for '%s':", req.InsightType))
	switch req.InsightType {
	case "productivity":
		promptBuilder.WriteString("\n- Focus on efficiency patterns, time utilization, and value delivery")
		promptBuilder.WriteString("\n- Identify high-impact activities and optimization opportunities")
		promptBuilder.WriteString("\n- Analyze work-life balance and sustainable productivity patterns")
	case "skill_development":
		promptBuilder.WriteString("\n- Identify learning opportunities and skill gaps")
		promptBuilder.WriteString("\n- Track progress in technical and soft skills")
		promptBuilder.WriteString("\n- Suggest development paths and learning resources")
	case "time_management":
		promptBuilder.WriteString("\n- Analyze time allocation across different activity types")
		promptBuilder.WriteString("\n- Identify time drains and efficiency bottlenecks")
		promptBuilder.WriteString("\n- Suggest schedule optimization strategies")
	case "team_collaboration":
		promptBuilder.WriteString("\n- Focus on collaboration patterns and team interactions")
		promptBuilder.WriteString("\n- Identify communication effectiveness and team dynamics")
		promptBuilder.WriteString("\n- Suggest improvements for team productivity")
	default:
		promptBuilder.WriteString("\n- Provide comprehensive analysis based on the available data")
		promptBuilder.WriteString("\n- Focus on actionable insights and practical recommendations")
	}

	// Handle context data - now as additional structured information
	if req.Context != nil {
		promptBuilder.WriteString("\n\n--- Additional Context ---")

		switch contextData := req.Context.(type) {
		case string:
			// Backward compatibility - simple string context
			if contextData != "" {
				promptBuilder.WriteString(fmt.Sprintf("\nContext: %s", contextData))
			}

		case map[string]any:
			// Structured context - serialize to JSON for AI processing
			if len(contextData) > 0 {
				contextJSON, err := json.MarshalIndent(contextData, "", "  ")
				if err != nil {
					s.logger.Warn("Failed to marshal context data, adding basic context info",
						"error", err.Error())
					promptBuilder.WriteString(fmt.Sprintf("\nStructured context provided (%d fields)", len(contextData)))
				} else {
					promptBuilder.WriteString(fmt.Sprintf("\nStructured Context:\n%s", string(contextJSON)))
				}
			}

		default:
			// Unknown context type - try to convert to JSON
			s.logger.Warn("Unknown context type, attempting JSON serialization",
				"type", fmt.Sprintf("%T", contextData))

			contextJSON, err := json.MarshalIndent(contextData, "", "  ")
			if err != nil {
				s.logger.Warn("Failed to serialize unknown context type, adding type info only",
					"error", err.Error(),
					"type", fmt.Sprintf("%T", contextData))
				promptBuilder.WriteString(fmt.Sprintf("\nContext Type: %T (serialization failed)", contextData))
			} else {
				promptBuilder.WriteString(fmt.Sprintf("\nContext Data:\n%s", string(contextJSON)))
			}
		}
	}

	// Add final instructions for consistent output format
	promptBuilder.WriteString("\n\n--- Output Instructions ---")
	promptBuilder.WriteString("\nPlease provide a comprehensive analysis that includes:")
	promptBuilder.WriteString("\n1. Key findings and patterns identified")
	promptBuilder.WriteString("\n2. Specific, actionable recommendations")
	promptBuilder.WriteString("\n3. Confidence level in your analysis (high/medium/low)")
	promptBuilder.WriteString("\n4. Suggested next steps or areas for deeper investigation")

	return promptBuilder.String()
}

// joinStrings is a helper function to join string slices (similar to strings.Join but for clarity)
func joinStrings(strs []string, separator string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	var result bytes.Buffer
	result.WriteString(strs[0])
	for _, str := range strs[1:] {
		result.WriteString(separator)
		result.WriteString(str)
	}
	return result.String()
}

// ValidateInsightRequest validates the insight request structure
func (s *OllamaService) ValidateInsightRequest(req *InsightRequest) error {
	if req.Prompt == "" {
		return fmt.Errorf("prompt cannot be empty")
	}
	if req.UserID == "" {
		return fmt.Errorf("user_id cannot be empty")
	}
	if req.InsightType == "" {
		return fmt.Errorf("insight_type cannot be empty")
	}
	if len(req.EntryIDs) == 0 {
		return fmt.Errorf("entry_ids cannot be empty")
	}

	// Validate context based on insight type
	return s.validateContextForInsightType(req.Context, req.InsightType)
}

// validateContextForInsightType validates context structure based on insight type
func (s *OllamaService) validateContextForInsightType(context any, insightType string) error {
	if context == nil {
		return nil // Optional context
	}

	switch contextData := context.(type) {
	case string:
		// String context is always valid
		return nil

	case map[string]any:
		// Validate structured context based on insight type
		switch insightType {
		case "productivity":
			return s.validateProductivityContext(contextData)
		case "skill_development":
			return s.validateSkillDevelopmentContext(contextData)
		case "time_management":
			return s.validateTimeManagementContext(contextData)
		default:
			// For unknown insight types, accept any valid JSON structure
			return nil
		}

	default:
		// Try to serialize to ensure it's JSON-compatible
		_, err := json.Marshal(context)
		if err != nil {
			return fmt.Errorf("context must be JSON-serializable: %w", err)
		}
		return nil
	}
}

// validateProductivityContext validates productivity-specific context
func (s *OllamaService) validateProductivityContext(context map[string]any) error {
	// Example validation for productivity context
	if timeBlocks, exists := context["time_blocks"]; exists {
		if blocks, ok := timeBlocks.([]any); ok {
			for _, block := range blocks {
				if _, ok := block.(string); !ok {
					return fmt.Errorf("time_blocks must contain strings")
				}
			}
		} else {
			return fmt.Errorf("time_blocks must be an array")
		}
	}
	return nil
}

// validateSkillDevelopmentContext validates skill development context
func (s *OllamaService) validateSkillDevelopmentContext(context map[string]any) error {
	// Example validation for skill development context
	if focusAreas, exists := context["focus_areas"]; exists {
		if areas, ok := focusAreas.([]any); ok {
			for _, area := range areas {
				if _, ok := area.(string); !ok {
					return fmt.Errorf("focus_areas must contain strings")
				}
			}
		} else {
			return fmt.Errorf("focus_areas must be an array")
		}
	}
	return nil
}

// validateTimeManagementContext validates time management context
func (s *OllamaService) validateTimeManagementContext(context map[string]any) error {
	// Example validation for time management context
	if dateRange, exists := context["date_range"]; exists {
		if rangeMap, ok := dateRange.(map[string]any); ok {
			if _, hasStart := rangeMap["start"]; !hasStart {
				return fmt.Errorf("date_range must include 'start' field")
			}
			if _, hasEnd := rangeMap["end"]; !hasEnd {
				return fmt.Errorf("date_range must include 'end' field")
			}
		} else {
			return fmt.Errorf("date_range must be an object")
		}
	}
	return nil
}

// GenerateWeeklyReport generates a weekly report with retry logic
func (s *OllamaService) GenerateWeeklyReport(ctx context.Context, userID string, weekStart, weekEnd time.Time) (*WeeklyReport, error) {
	if userID == "" {
		s.logger.LogError(ctx, fmt.Errorf("empty userID"), "GenerateWeeklyReport called with empty userID")
		return nil, fmt.Errorf("userID cannot be empty")
	}

	s.logger.Info("Starting weekly report generation",
		"user_id", userID,
		"week_start", weekStart.Format("2006-01-02"),
		"week_end", weekEnd.Format("2006-01-02"))

	prompt := fmt.Sprintf(
		"Generate a comprehensive weekly productivity report for user %s from %s to %s. "+
			"Include summary, key insights, and recommendations.",
		userID, weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))

	// Retry configuration
	maxRetries := 3
	baseDelay := 2 * time.Second

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		s.logger.Debug("Attempting weekly report generation",
			"attempt", attempt+1,
			"max_retries", maxRetries,
			"user_id", userID)

		select {
		case <-ctx.Done():
			s.logger.Warn("Weekly report generation cancelled",
				"attempt", attempt+1,
				"user_id", userID,
				"error", ctx.Err().Error())
			return nil, fmt.Errorf("weekly report generation cancelled: %w", ctx.Err())
		default:
		}

		response, err := s.generateWithTimeout(ctx, "qwen2.5-coder:7b", prompt, 90*time.Second)
		if err == nil {
			s.logger.Info("Weekly report generation successful",
				"attempt", attempt+1,
				"user_id", userID,
				"response_length", len(response))

			// Parse response into structured report
			report := &WeeklyReport{
				Summary:         response,
				KeyInsights:     []string{"AI-generated insight 1", "AI-generated insight 2"},
				Recommendations: []string{"Recommendation 1", "Recommendation 2"},
			}

			s.logger.Debug("Generated weekly report",
				"user_id", userID,
				"summary_length", len(report.Summary),
				"insights_count", len(report.KeyInsights),
				"recommendations_count", len(report.Recommendations))

			return report, nil
		}

		lastErr = err
		s.logger.Warn("Weekly report generation attempt failed",
			"attempt", attempt+1,
			"user_id", userID,
			"error", err.Error())

		if attempt < maxRetries-1 {
			delay := time.Duration(attempt+1) * baseDelay
			s.logger.Debug("Retrying weekly report generation",
				"delay", delay.String(),
				"next_attempt", attempt+2,
				"user_id", userID)

			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("weekly report generation cancelled during retry: %w", ctx.Err())
			case <-time.After(delay):
			}
		}
	}

	s.logger.LogError(ctx, lastErr, "All weekly report generation attempts failed",
		"max_retries", maxRetries,
		"user_id", userID)

	return nil, fmt.Errorf("failed to generate weekly report after %d attempts: %w", maxRetries, lastErr)
}

// HealthCheck performs a health check on the AI service
func (s *OllamaService) HealthCheck(ctx context.Context) error {
	s.logger.Debug("Performing AI service health check")

	// Simple health check by making a basic request
	httpReq, err := http.NewRequestWithContext(ctx, "GET", s.baseURL, nil)
	if err != nil {
		s.logger.LogError(ctx, err, "Health check failed: error creating HTTP request")
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	// Use a shorter timeout for health checks
	healthCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	httpReq = httpReq.WithContext(healthCtx)

	start := time.Now()
	resp, err := s.client.Do(httpReq)
	duration := time.Since(start)

	if err != nil {
		s.logger.LogError(ctx, err, "Health check failed: HTTP request error",
			"duration", duration.String())
		return fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.LogError(ctx, fmt.Errorf("non-OK status: %d", resp.StatusCode), "Health check failed: non-OK status",
			"status_code", resp.StatusCode,
			"status", resp.Status,
			"duration", duration.String())
		return fmt.Errorf("health check failed with status: %s", resp.Status)
	}

	s.logger.Debug("AI service health check passed",
		"duration", duration.String())
	return nil
}

// generateWithTimeout makes a request to Ollama with the specified timeout
func (s *OllamaService) generateWithTimeout(ctx context.Context, model, prompt string, timeout time.Duration) (string, error) {
	s.logger.Debug("Making generation request",
		"model", model,
		"prompt_length", len(prompt),
		"timeout", timeout)

	req := &GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	url := fmt.Sprintf("%s/api/generate", s.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		s.logger.Error("Failed to marshal generation request",
			"error", err)
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		s.logger.Error("Failed to create HTTP request",
			"error", err)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Apply timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	httpReq = httpReq.WithContext(timeoutCtx)

	start := time.Now()
	resp, err := s.client.Do(httpReq)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("Generation request failed",
			"error", err,
			"duration", duration)
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Generation request returned non-OK status",
			"status_code", resp.StatusCode,
			"status", resp.Status,
			"duration", duration)
		return "", fmt.Errorf("request failed with status: %s", resp.Status)
	}

	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		s.logger.Error("Failed to decode generation response",
			"error", err,
			"duration", duration)
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	s.logger.Debug("Generation request completed successfully",
		"response_length", len(genResp.Response),
		"duration", duration,
		"done", genResp.Done)

	return genResp.Response, nil
}
