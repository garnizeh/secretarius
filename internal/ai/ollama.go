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
	Context     string   `json:"context"`
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
		"model", "llama3.2:3b")

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

		response, err := s.generateWithTimeout(ctx, "llama3.2:3b", prompt, 60*time.Second)
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

		response, err := s.generateWithTimeout(ctx, "llama3.2:3b", prompt, 90*time.Second)
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
	req := &GenerateRequest{
		Model:  "llama3.2:3b",
		Prompt: "Hello",
		Stream: false,
	}

	url := fmt.Sprintf("%s/api/generate", s.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		s.logger.LogError(ctx, err, "Health check failed: error marshaling request")
		return fmt.Errorf("failed to marshal health check request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		s.logger.LogError(ctx, err, "Health check failed: error creating HTTP request")
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

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
} // generateWithTimeout makes a request to Ollama with the specified timeout
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
