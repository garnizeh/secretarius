package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/garnizeh/englog/internal/logging"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// OllamaService implements AI service using Ollama with langchaingo
type OllamaService struct {
	logger    *logging.Logger
	baseURL   string
	modelName string
	llm       llms.Model
}

// Insight represents an AI-generated insight
type Insight struct {
	Content    string   `json:"content"`
	Tags       []string `json:"tags"`
	Confidence float32  `json:"confidence"`
}

// InsightRequest represents a request for insight generation
type InsightRequest struct {
	UserID      string   `json:"user_id"`
	EntryIDs    []string `json:"entry_ids"`
	InsightType string   `json:"insight_type"`
	Context     any      `json:"context,omitempty"`
}

// Prompt generates the AI prompt based on the request
func (r *InsightRequest) Prompt(ctx context.Context, logger *logging.Logger) string {
	var promptBuilder bytes.Buffer

	// Add structured request information to enhance AI understanding
	promptBuilder.WriteString("\n\n--- Request Information ---")
	promptBuilder.WriteString(fmt.Sprintf("\nUser ID: %s", r.UserID))
	promptBuilder.WriteString(fmt.Sprintf("\nInsight Type: %s", r.InsightType))
	promptBuilder.WriteString(fmt.Sprintf("\nNumber of Log Entries: %d", len(r.EntryIDs)))

	// Include entry IDs for reference (useful for the AI to understand scope)
	if len(r.EntryIDs) > 0 {
		promptBuilder.WriteString("\nLog Entry IDs: ")
		if len(r.EntryIDs) <= 5 {
			// Show all IDs if there are 5 or fewer
			promptBuilder.WriteString(fmt.Sprintf("[%s]", joinStrings(r.EntryIDs, ", ")))
		} else {
			// Show first 3 and last 2 with count for larger sets
			firstThree := r.EntryIDs[:3]
			lastTwo := r.EntryIDs[len(r.EntryIDs)-2:]
			promptBuilder.WriteString(fmt.Sprintf("[%s, ... (%d more), %s]",
				joinStrings(firstThree, ", "),
				len(r.EntryIDs)-5,
				joinStrings(lastTwo, ", ")))
		}
	}

	// Add insight type specific instructions
	promptBuilder.WriteString(fmt.Sprintf("\n\nInsight Generation Guidelines for '%s':", r.InsightType))
	switch r.InsightType {
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
	if r.Context != nil {
		promptBuilder.WriteString("\n\n--- Additional Context ---")

		switch contextData := r.Context.(type) {
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
					logger.LogWarn(ctx, "Failed to marshal context data, adding basic context info",
						logging.OperationField, "insight_request_prompt",
						logging.ErrorField, err)
					promptBuilder.WriteString(fmt.Sprintf("\nStructured context provided (%d fields)", len(contextData)))
				} else {
					promptBuilder.WriteString(fmt.Sprintf("\nStructured Context:\n%s", string(contextJSON)))
				}
			}

		default:
			// Unknown context type - try to convert to JSON
			logger.LogWarn(ctx, "Unknown context type, attempting JSON serialization",
				logging.OperationField, "insight_request_prompt",
				"type", fmt.Sprintf("%T", contextData))

			contextJSON, err := json.MarshalIndent(contextData, "", "  ")
			if err != nil {
				logger.LogWarn(ctx, "Failed to serialize unknown context type, adding type info only",
					logging.OperationField, "insight_request_prompt",
					logging.ErrorField, err,
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

// NewOllamaService creates a new Ollama service instance using langchaingo
func NewOllamaService(ctx context.Context, baseURL string, logger *logging.Logger) (*OllamaService, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("ollama base URL cannot be empty")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	serviceLogger := logger.WithServiceAndComponent("worker", "ollama_service")

	// Default model - this should be configurable in the future
	modelName := "qwen2.5-coder:7b"

	serviceLogger.LogInfo(ctx, "Creating new Ollama service with langchaingo",
		logging.OperationField, "new_ollama_service",
		"base_url", baseURL,
		"model", modelName)

	// Create langchaingo Ollama LLM instance
	llm, err := ollama.New(
		ollama.WithServerURL(baseURL),
		ollama.WithModel(modelName),
	)
	if err != nil {
		serviceLogger.LogError(ctx, err, "Failed to create langchaingo Ollama LLM",
			logging.OperationField, "new_ollama_service")
		return nil, fmt.Errorf("failed to create Ollama LLM: %w", err)
	}

	serviceLogger.LogInfo(ctx, "Successfully created langchaingo Ollama LLM",
		logging.OperationField, "new_ollama_service",
		"model", modelName)

	return &OllamaService{
		logger:    serviceLogger,
		baseURL:   baseURL,
		modelName: modelName,
		llm:       llm,
	}, nil
}

// GenerateInsight generates AI insights using langchaingo
func (s *OllamaService) GenerateInsight(ctx context.Context, req *InsightRequest) (*Insight, error) {
	prompt := req.Prompt(ctx, s.logger)
	if prompt == "" {
		s.logger.LogError(ctx, fmt.Errorf("empty prompt"), "GenerateInsight req did not generated prompt")
		return nil, fmt.Errorf("invalid request: did not generated prompt")
	}

	s.logger.LogInfo(ctx, "Starting insight generation with langchaingo",
		logging.OperationField, "generate_insight",
		logging.UserIDField, req.UserID,
		"insight_type", req.InsightType,
		"entry_count", len(req.EntryIDs),
		"enhanced_prompt_length", len(prompt),
		"context_type", fmt.Sprintf("%T", req.Context),
		"model", s.modelName)

	// Retry configuration for AI operations
	maxRetries := 3
	baseDelay := 1 * time.Second

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		s.logger.LogDebug(ctx, "Attempting insight generation",
			logging.OperationField, "generate_insight",
			"attempt", attempt+1,
			"max_retries", maxRetries)

		select {
		case <-ctx.Done():
			s.logger.LogWarn(ctx, "Insight generation cancelled",
				logging.OperationField, "generate_insight",
				"attempt", attempt+1,
				logging.ErrorField, ctx.Err())
			return nil, fmt.Errorf("insight generation cancelled: %w", ctx.Err())
		default:
		}

		// Create a timeout context for this attempt
		timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		// Use langchaingo GenerateFromSinglePrompt method
		response, err := llms.GenerateFromSinglePrompt(timeoutCtx, s.llm, prompt)
		if err == nil {
			s.logger.LogInfo(ctx, "Insight generation successful",
				logging.OperationField, "generate_insight",
				"attempt", attempt+1,
				"response_length", len(response))

			// Parse response and extract insights
			insight := &Insight{
				Content:    response,
				Tags:       []string{"ai-generated", "productivity"},
				Confidence: 0.8, // Default confidence
			}

			s.logger.LogDebug(ctx, "Generated insight",
				logging.OperationField, "generate_insight",
				"content_length", len(insight.Content),
				"tags_count", len(insight.Tags),
				"confidence", insight.Confidence)

			return insight, nil
		}

		lastErr = err
		s.logger.LogWarn(ctx, "Insight generation attempt failed",
			logging.OperationField, "generate_insight",
			"attempt", attempt+1,
			logging.ErrorField, err)

		if attempt < maxRetries-1 {
			delay := time.Duration(attempt+1) * baseDelay
			s.logger.LogDebug(ctx, "Retrying insight generation",
				logging.OperationField, "generate_insight",
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

// GenerateWeeklyReport generates a weekly report using langchaingo
func (s *OllamaService) GenerateWeeklyReport(ctx context.Context, userID string, weekStart, weekEnd time.Time) (*WeeklyReport, error) {
	if userID == "" {
		s.logger.LogError(ctx, fmt.Errorf("empty userID"), "GenerateWeeklyReport called with empty userID")
		return nil, fmt.Errorf("userID cannot be empty")
	}

	s.logger.LogInfo(ctx, "Starting weekly report generation with langchaingo",
		logging.OperationField, "generate_weekly_report",
		logging.UserIDField, userID,
		"week_start", weekStart.Format("2006-01-02"),
		"week_end", weekEnd.Format("2006-01-02"),
		"model", s.modelName)

	prompt := fmt.Sprintf(
		"Generate a comprehensive weekly productivity report for user %s from %s to %s. "+
			"Include summary, key insights, and recommendations.",
		userID, weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))

	// Retry configuration
	maxRetries := 3
	baseDelay := 2 * time.Second

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		s.logger.LogDebug(ctx, "Attempting weekly report generation",
			logging.OperationField, "generate_weekly_report",
			"attempt", attempt+1,
			"max_retries", maxRetries,
			logging.UserIDField, userID)

		select {
		case <-ctx.Done():
			s.logger.LogWarn(ctx, "Weekly report generation cancelled",
				logging.OperationField, "generate_weekly_report",
				"attempt", attempt+1,
				logging.UserIDField, userID,
				logging.ErrorField, ctx.Err())
			return nil, fmt.Errorf("weekly report generation cancelled: %w", ctx.Err())
		default:
		}

		// Create a timeout context for this attempt
		timeoutCtx, cancel := context.WithTimeout(ctx, 90*time.Second)
		defer cancel()

		// Use langchaingo GenerateFromSinglePrompt method
		response, err := llms.GenerateFromSinglePrompt(timeoutCtx, s.llm, prompt)
		if err == nil {
			s.logger.LogInfo(ctx, "Weekly report generation successful",
				logging.OperationField, "generate_weekly_report",
				"attempt", attempt+1,
				logging.UserIDField, userID,
				"response_length", len(response))

			// Parse response into structured report
			report := &WeeklyReport{
				Summary:         response,
				KeyInsights:     []string{"AI-generated insight 1", "AI-generated insight 2"},
				Recommendations: []string{"Recommendation 1", "Recommendation 2"},
			}

			s.logger.LogDebug(ctx, "Generated weekly report",
				logging.OperationField, "generate_weekly_report",
				logging.UserIDField, userID,
				"summary_length", len(report.Summary),
				"insights_count", len(report.KeyInsights),
				"recommendations_count", len(report.Recommendations))

			return report, nil
		}

		lastErr = err
		s.logger.LogWarn(ctx, "Weekly report generation attempt failed",
			logging.OperationField, "generate_weekly_report",
			"attempt", attempt+1,
			logging.UserIDField, userID,
			logging.ErrorField, err)

		if attempt < maxRetries-1 {
			delay := time.Duration(attempt+1) * baseDelay
			s.logger.LogDebug(ctx, "Retrying weekly report generation",
				logging.OperationField, "generate_weekly_report",
				"delay", delay.String(),
				"next_attempt", attempt+2,
				logging.UserIDField, userID)

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

// HealthCheck performs a health check on the AI service using a simple prompt
func (s *OllamaService) HealthCheck(ctx context.Context) error {
	s.logger.LogDebug(ctx, "Performing AI service health check with langchaingo",
		logging.OperationField, "health_check",
		"model", s.modelName)

	// Create a shorter timeout for health checks
	healthCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Simple health check with a basic prompt
	testPrompt := "Respond with 'OK' to confirm you are working."

	start := time.Now()
	response, err := llms.GenerateFromSinglePrompt(healthCtx, s.llm, testPrompt)
	duration := time.Since(start)

	if err != nil {
		s.logger.LogError(ctx, err, "Health check failed: LLM call error",
			"duration", duration.String(),
			"model", s.modelName)
		return fmt.Errorf("health check failed: %w", err)
	}

	// Check if we got any response
	if len(response) == 0 {
		s.logger.LogError(ctx, fmt.Errorf("empty response"), "Health check failed: empty response",
			"duration", duration.String(),
			"model", s.modelName)
		return fmt.Errorf("health check failed: empty response from LLM")
	}

	s.logger.LogDebug(ctx, "AI service health check passed",
		logging.OperationField, "health_check",
		"duration", duration.String(),
		"response_length", len(response),
		"model", s.modelName)

	return nil
}
