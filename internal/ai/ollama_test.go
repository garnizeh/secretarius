package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewOllamaService tests the creation of OllamaService instances
func TestNewOllamaService(t *testing.T) {
	ctx := context.Background()

	logger := logging.NewTestLogger()

	tests := []struct {
		name        string
		baseURL     string
		logger      *logging.Logger
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid parameters",
			baseURL:     "http://localhost:11434",
			logger:      logger,
			expectError: false,
		},
		{
			name:        "Empty base URL",
			baseURL:     "",
			logger:      logger,
			expectError: true,
			errorMsg:    "ollama base URL cannot be empty",
		},
		{
			name:        "Nil logger",
			baseURL:     "http://localhost:11434",
			logger:      nil,
			expectError: true,
			errorMsg:    "logger cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewOllamaService(ctx, tt.baseURL, tt.logger)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
				assert.Equal(t, tt.baseURL, service.baseURL)
				assert.Equal(t, tt.logger, service.logger)
				assert.NotNil(t, service.client)
				assert.Equal(t, 120*time.Second, service.client.Timeout)
			}
		})
	}
}

// TestGenerateInsight tests the basic insight generation functionality
func TestGenerateInsight(t *testing.T) {
	ctx := context.Background()

	logger := logging.NewTestLogger()

	tests := []struct {
		name           string
		prompt         string
		serverResponse string
		serverStatus   int
		expectError    bool
		errorMsg       string
		expectedTags   []string
		expectedConf   float32
	}{
		{
			name:           "Successful generation",
			prompt:         "Test prompt",
			serverResponse: `{"response": "Test insight response", "done": true}`,
			serverStatus:   http.StatusOK,
			expectError:    false,
			expectedTags:   []string{"ai-generated", "productivity"},
			expectedConf:   0.8,
		},
		{
			name:        "Empty prompt",
			prompt:      "",
			expectError: true,
			errorMsg:    "prompt cannot be empty",
		},
		{
			name:           "Server error",
			prompt:         "Test prompt",
			serverResponse: `{"error": "Server error"}`,
			serverStatus:   http.StatusInternalServerError,
			expectError:    true,
			errorMsg:       "request failed with status",
		},
		{
			name:           "Invalid JSON response",
			prompt:         "Test prompt",
			serverResponse: `invalid json`,
			serverStatus:   http.StatusOK,
			expectError:    true,
			errorMsg:       "failed to decode response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "/api/generate", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				// Verify request body
				var req GenerateRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				if err == nil {
					assert.Equal(t, "qwen2.5-coder:7b", req.Model)
					assert.Equal(t, tt.prompt, req.Prompt)
					assert.False(t, req.Stream)
				}

				w.WriteHeader(tt.serverStatus)
				_, _ = w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			service, err := NewOllamaService(ctx, server.URL, logger)
			require.NoError(t, err)

			ctx := context.Background()
			insight, err := service.GenerateInsight(ctx, tt.prompt)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, insight)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, insight)
				assert.Equal(t, "Test insight response", insight.Content)
				assert.Equal(t, tt.expectedTags, insight.Tags)
				assert.Equal(t, tt.expectedConf, insight.Confidence)
			}
		})
	}
}

// TestGenerateInsightWithContext tests the context-aware insight generation
func TestGenerateInsightWithContext(t *testing.T) {
	ctx := context.Background()

	logger := logging.NewTestLogger()

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response": "Context-aware insight response", "done": true}`))
	}))
	defer server.Close()

	service, err := NewOllamaService(ctx, server.URL, logger)
	require.NoError(t, err)

	tests := []struct {
		name        string
		request     *InsightRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid request with string context",
			request: &InsightRequest{
				Prompt:      "Test prompt",
				UserID:      "user123",
				EntryIDs:    []string{"entry1", "entry2"},
				InsightType: "productivity",
				Context:     "Simple string context",
			},
			expectError: false,
		},
		{
			name: "Valid request with structured context",
			request: &InsightRequest{
				Prompt:      "Test prompt",
				UserID:      "user123",
				EntryIDs:    []string{"entry1", "entry2"},
				InsightType: "productivity",
				Context: map[string]any{
					"time_blocks": []string{"morning", "afternoon"},
					"focus_level": 8,
				},
			},
			expectError: false,
		},
		{
			name: "Valid request with nil context",
			request: &InsightRequest{
				Prompt:      "Test prompt",
				UserID:      "user123",
				EntryIDs:    []string{"entry1", "entry2"},
				InsightType: "productivity",
				Context:     nil,
			},
			expectError: false,
		},
		{
			name: "Empty prompt",
			request: &InsightRequest{
				Prompt:      "",
				UserID:      "user123",
				EntryIDs:    []string{"entry1"},
				InsightType: "productivity",
				Context:     "context",
			},
			expectError: true,
			errorMsg:    "prompt cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			insight, err := service.GenerateInsightWithContext(ctx, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, insight)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, insight)
				assert.Equal(t, "Context-aware insight response", insight.Content)
			}
		})
	}
}

// TestBuildEnhancedPrompt tests the prompt enhancement functionality with complete request information
func TestBuildEnhancedPrompt(t *testing.T) {
	ctx := context.Background()
	logger := logging.NewTestLogger()
	service, err := NewOllamaService(ctx, "http://localhost:11434", logger)
	require.NoError(t, err)

	tests := []struct {
		name     string
		request  *InsightRequest
		contains []string // What the result should contain
	}{
		{
			name: "Complete request with string context",
			request: &InsightRequest{
				Prompt:      "Analyze my productivity",
				UserID:      "user-123",
				EntryIDs:    []string{"entry-1", "entry-2"},
				InsightType: "productivity",
				Context:     "String context data",
			},
			contains: []string{
				"Analyze my productivity",
				"User ID: user-123",
				"Insight Type: productivity",
				"Number of Log Entries: 2",
				"Log Entry IDs: [entry-1, entry-2]",
				"Focus on efficiency patterns",
				"Context: String context data",
				"Key findings and patterns identified",
			},
		},
		{
			name: "Request with structured context",
			request: &InsightRequest{
				Prompt:      "Base prompt",
				UserID:      "user-456",
				EntryIDs:    []string{"entry-1"},
				InsightType: "skill_development",
				Context: map[string]any{
					"key1": "value1",
					"key2": 42,
				},
			},
			contains: []string{
				"Base prompt",
				"User ID: user-456",
				"Insight Type: skill_development",
				"Number of Log Entries: 1",
				"Identify learning opportunities",
				"Structured Context:",
				"\"key1\": \"value1\"",
				"\"key2\": 42",
			},
		},
		{
			name: "Request with many entry IDs (truncated display)",
			request: &InsightRequest{
				Prompt:      "Weekly analysis",
				UserID:      "user-789",
				EntryIDs:    []string{"e1", "e2", "e3", "e4", "e5", "e6", "e7", "e8"},
				InsightType: "time_management",
				Context:     nil,
			},
			contains: []string{
				"Weekly analysis",
				"User ID: user-789",
				"Number of Log Entries: 8",
				"[e1, e2, e3, ... (3 more), e7, e8]",
				"time allocation across different",
				"Output Instructions",
			},
		},
		{
			name: "Request with minimal data",
			request: &InsightRequest{
				Prompt:      "Simple prompt",
				UserID:      "",
				EntryIDs:    []string{},
				InsightType: "",
				Context:     nil,
			},
			contains: []string{
				"Simple prompt",
				"User ID: ",
				"Number of Log Entries: 0",
				"Provide comprehensive analysis",
				"Output Instructions",
			},
		},
		{
			name: "Request with team collaboration type",
			request: &InsightRequest{
				Prompt:      "Team work analysis",
				UserID:      "team-lead",
				EntryIDs:    []string{"meeting-1", "standup-2"},
				InsightType: "team_collaboration",
				Context:     "Weekly team retrospective",
			},
			contains: []string{
				"Team work analysis",
				"Insight Type: team_collaboration",
				"Focus on collaboration patterns",
				"team interactions",
				"Context: Weekly team retrospective",
			},
		},
		{
			name: "Request with custom struct context",
			request: &InsightRequest{
				Prompt:      "Custom analysis",
				UserID:      "user-custom",
				EntryIDs:    []string{"custom-1"},
				InsightType: "productivity",
				Context: struct {
					Name  string `json:"name"`
					Value int    `json:"value"`
				}{
					Name:  "test",
					Value: 123,
				},
			},
			contains: []string{
				"Custom analysis",
				"Context Data:",
				"\"name\": \"test\"",
				"\"value\": 123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.buildEnhancedPrompt(ctx, tt.request)

			// Check that all expected content is present
			for _, expected := range tt.contains {
				assert.Contains(t, result, expected,
					"Expected '%s' to be found in result:\n%s", expected, result)
			}

			// Verify basic structure
			assert.Contains(t, result, tt.request.Prompt, "Should contain original prompt")
			assert.Contains(t, result, "--- Request Information ---", "Should contain request info section")
			assert.Contains(t, result, "--- Output Instructions ---", "Should contain output instructions")
		})
	}
}

// TestBuildEnhancedPromptExample demonstrates how the enhanced prompt looks
func TestBuildEnhancedPromptExample(t *testing.T) {
	ctx := context.Background()
	logger := logging.NewTestLogger()
	service, err := NewOllamaService(ctx, "http://localhost:11434", logger)
	require.NoError(t, err)

	// Create a comprehensive example request
	req := &InsightRequest{
		Prompt:      "Please analyze my productivity patterns and provide actionable insights for improvement.",
		UserID:      "user-12345",
		EntryIDs:    []string{"entry-001", "entry-002", "entry-003", "entry-004", "entry-005", "entry-006"},
		InsightType: "productivity",
		Context: map[string]any{
			"time_blocks": []string{"morning", "afternoon", "evening"},
			"focus_areas": []string{"development", "meetings", "documentation"},
			"date_range": map[string]string{
				"start": "2025-07-01",
				"end":   "2025-07-31",
			},
			"performance_metrics": map[string]float64{
				"avg_daily_hours":    8.5,
				"productivity_score": 0.85,
			},
		},
	}

	result := service.buildEnhancedPrompt(ctx, req)

	// Log the result for demonstration purposes
	t.Logf("=== Enhanced Prompt Example ===\n%s\n=== End of Example ===", result)

	// Verify it contains key sections
	assert.Contains(t, result, "Please analyze my productivity patterns")
	assert.Contains(t, result, "User ID: user-12345")
	assert.Contains(t, result, "Insight Type: productivity")
	assert.Contains(t, result, "Number of Log Entries: 6")
	assert.Contains(t, result, "Focus on efficiency patterns")
	assert.Contains(t, result, "time_blocks")
	assert.Contains(t, result, "performance_metrics")
} // TestValidateInsightRequest tests the insight request validation
func TestValidateInsightRequest(t *testing.T) {
	ctx := context.Background()
	logger := logging.NewTestLogger()
	service, err := NewOllamaService(ctx, "http://localhost:11434", logger)
	require.NoError(t, err)

	tests := []struct {
		name        string
		request     *InsightRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid request",
			request: &InsightRequest{
				Prompt:      "Test prompt",
				UserID:      "user123",
				EntryIDs:    []string{"entry1", "entry2"},
				InsightType: "productivity",
				Context:     "test context",
			},
			expectError: false,
		},
		{
			name: "Empty prompt",
			request: &InsightRequest{
				Prompt:      "",
				UserID:      "user123",
				EntryIDs:    []string{"entry1"},
				InsightType: "productivity",
			},
			expectError: true,
			errorMsg:    "prompt cannot be empty",
		},
		{
			name: "Empty user ID",
			request: &InsightRequest{
				Prompt:      "Test prompt",
				UserID:      "",
				EntryIDs:    []string{"entry1"},
				InsightType: "productivity",
			},
			expectError: true,
			errorMsg:    "user_id cannot be empty",
		},
		{
			name: "Empty insight type",
			request: &InsightRequest{
				Prompt:      "Test prompt",
				UserID:      "user123",
				EntryIDs:    []string{"entry1"},
				InsightType: "",
			},
			expectError: true,
			errorMsg:    "insight_type cannot be empty",
		},
		{
			name: "Empty entry IDs",
			request: &InsightRequest{
				Prompt:      "Test prompt",
				UserID:      "user123",
				EntryIDs:    []string{},
				InsightType: "productivity",
			},
			expectError: true,
			errorMsg:    "entry_ids cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateInsightRequest(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidateContextForInsightType tests context validation for different insight types
func TestValidateContextForInsightType(t *testing.T) {
	ctx := context.Background()
	logger := logging.NewTestLogger()
	service, err := NewOllamaService(ctx, "http://localhost:11434", logger)
	require.NoError(t, err)

	tests := []struct {
		name        string
		context     any
		insightType string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Nil context",
			context:     nil,
			insightType: "productivity",
			expectError: false,
		},
		{
			name:        "String context",
			context:     "test context",
			insightType: "productivity",
			expectError: false,
		},
		{
			name: "Valid productivity context",
			context: map[string]any{
				"time_blocks": []any{"morning", "afternoon"},
			},
			insightType: "productivity",
			expectError: false,
		},
		{
			name: "Invalid productivity context - wrong type",
			context: map[string]any{
				"time_blocks": []any{123, 456},
			},
			insightType: "productivity",
			expectError: true,
			errorMsg:    "time_blocks must contain strings",
		},
		{
			name: "Invalid productivity context - not array",
			context: map[string]any{
				"time_blocks": "not an array",
			},
			insightType: "productivity",
			expectError: true,
			errorMsg:    "time_blocks must be an array",
		},
		{
			name: "Valid skill development context",
			context: map[string]any{
				"focus_areas": []any{"golang", "testing"},
			},
			insightType: "skill_development",
			expectError: false,
		},
		{
			name: "Invalid skill development context",
			context: map[string]any{
				"focus_areas": []any{123, 456},
			},
			insightType: "skill_development",
			expectError: true,
			errorMsg:    "focus_areas must contain strings",
		},
		{
			name: "Valid time management context",
			context: map[string]any{
				"date_range": map[string]any{
					"start": "2025-01-01",
					"end":   "2025-01-31",
				},
			},
			insightType: "time_management",
			expectError: false,
		},
		{
			name: "Invalid time management context - missing start",
			context: map[string]any{
				"date_range": map[string]any{
					"end": "2025-01-31",
				},
			},
			insightType: "time_management",
			expectError: true,
			errorMsg:    "date_range must include 'start' field",
		},
		{
			name: "Invalid time management context - missing end",
			context: map[string]any{
				"date_range": map[string]any{
					"start": "2025-01-01",
				},
			},
			insightType: "time_management",
			expectError: true,
			errorMsg:    "date_range must include 'end' field",
		},
		{
			name: "Unknown insight type with valid JSON",
			context: map[string]any{
				"custom_field": "custom_value",
			},
			insightType: "unknown_type",
			expectError: false,
		},
		{
			name:        "Non-JSON serializable context",
			context:     make(chan int), // channels are not JSON serializable
			insightType: "productivity",
			expectError: true,
			errorMsg:    "context must be JSON-serializable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateContextForInsightType(tt.context, tt.insightType)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGenerateWeeklyReport tests the weekly report generation functionality
func TestGenerateWeeklyReport(t *testing.T) {
	ctx := context.Background()

	logger := logging.NewTestLogger()

	tests := []struct {
		name           string
		userID         string
		weekStart      time.Time
		weekEnd        time.Time
		serverResponse string
		serverStatus   int
		expectError    bool
		errorMsg       string
	}{
		{
			name:           "Successful generation",
			userID:         "user123",
			weekStart:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			weekEnd:        time.Date(2025, 1, 7, 23, 59, 59, 0, time.UTC),
			serverResponse: `{"response": "Weekly report content", "done": true}`,
			serverStatus:   http.StatusOK,
			expectError:    false,
		},
		{
			name:        "Empty user ID",
			userID:      "",
			weekStart:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			weekEnd:     time.Date(2025, 1, 7, 23, 59, 59, 0, time.UTC),
			expectError: true,
			errorMsg:    "userID cannot be empty",
		},
		{
			name:           "Server error",
			userID:         "user123",
			weekStart:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			weekEnd:        time.Date(2025, 1, 7, 23, 59, 59, 0, time.UTC),
			serverResponse: `{"error": "Server error"}`,
			serverStatus:   http.StatusInternalServerError,
			expectError:    true,
			errorMsg:       "request failed with status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverStatus)
				_, _ = w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			service, err := NewOllamaService(ctx, server.URL, logger)
			require.NoError(t, err)

			ctx := context.Background()
			report, err := service.GenerateWeeklyReport(ctx, tt.userID, tt.weekStart, tt.weekEnd)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, report)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, report)
				assert.Equal(t, "Weekly report content", report.Summary)
				assert.Len(t, report.KeyInsights, 2)
				assert.Len(t, report.Recommendations, 2)
			}
		})
	}
}

// TestHealthCheck tests the health check functionality
func TestHealthCheck(t *testing.T) {
	ctx := context.Background()

	logger := logging.NewTestLogger()

	tests := []struct {
		name           string
		serverResponse string
		serverStatus   int
		expectError    bool
		errorMsg       string
	}{
		{
			name:           "Healthy service",
			serverResponse: `{"response": "Hello response", "done": true}`,
			serverStatus:   http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Unhealthy service",
			serverResponse: `{"error": "Service unavailable"}`,
			serverStatus:   http.StatusServiceUnavailable,
			expectError:    true,
			errorMsg:       "health check failed with status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "/api/generate", r.URL.Path)

				// Verify it's a health check request
				var req GenerateRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				if err == nil {
					assert.Equal(t, "qwen2.5-coder:7b", req.Model)
					assert.Equal(t, "Hello", req.Prompt)
					assert.False(t, req.Stream)
				}

				w.WriteHeader(tt.serverStatus)
				_, _ = w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			service, err := NewOllamaService(ctx, server.URL, logger)
			require.NoError(t, err)

			ctx := context.Background()
			err = service.HealthCheck(ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGenerateWithTimeout tests the timeout functionality
func TestGenerateWithTimeout(t *testing.T) {
	ctx := context.Background()

	logger := logging.NewTestLogger()

	tests := []struct {
		name           string
		timeout        time.Duration
		serverDelay    time.Duration
		serverResponse string
		serverStatus   int
		expectError    bool
		errorMsg       string
	}{
		{
			name:           "Successful generation within timeout",
			timeout:        2 * time.Second,
			serverDelay:    100 * time.Millisecond,
			serverResponse: `{"response": "Test response", "done": true}`,
			serverStatus:   http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Timeout exceeded",
			timeout:        100 * time.Millisecond,
			serverDelay:    200 * time.Millisecond,
			serverResponse: `{"response": "Should not reach here", "done": true}`,
			serverStatus:   http.StatusOK,
			expectError:    true,
			errorMsg:       "request failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server with delay
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tt.serverDelay)
				w.WriteHeader(tt.serverStatus)
				_, _ = w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			service, err := NewOllamaService(ctx, server.URL, logger)
			require.NoError(t, err)

			ctx := context.Background()
			response, err := service.generateWithTimeout(ctx, "qwen2.5-coder:7b", "test prompt", tt.timeout)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Empty(t, response)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "Test response", response)
			}
		})
	}
}

// TestContextCancellation tests context cancellation scenarios
func TestContextCancellation(t *testing.T) {
	ctx := context.Background()

	logger := logging.NewTestLogger()

	// Create mock server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Long delay to allow cancellation
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response": "Should not reach here", "done": true}`))
	}))
	defer server.Close()

	service, err := NewOllamaService(ctx, server.URL, logger)
	require.NoError(t, err)

	t.Run("GenerateInsight context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		insight, err := service.GenerateInsight(ctx, "test prompt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insight generation cancelled")
		assert.Nil(t, insight)
	})

	t.Run("GenerateWeeklyReport context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		weekStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		weekEnd := time.Date(2025, 1, 7, 23, 59, 59, 0, time.UTC)

		report, err := service.GenerateWeeklyReport(ctx, "user123", weekStart, weekEnd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "weekly report generation cancelled")
		assert.Nil(t, report)
	})

	t.Run("HealthCheck context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		err := service.HealthCheck(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "health check request failed")
	})
}

// TestRetryMechanism tests the retry functionality for failed requests
func TestRetryMechanism(t *testing.T) {
	ctx := context.Background()
	logger := logging.NewTestLogger()
	attemptCount := 0

	// Create mock server that fails twice then succeeds
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Temporary error"}`))
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response": "Success after retry", "done": true}`))
		}
	}))
	defer server.Close()

	service, err := NewOllamaService(ctx, server.URL, logger)
	require.NoError(t, err)

	insight, err := service.GenerateInsight(ctx, "test prompt")

	assert.NoError(t, err)
	assert.NotNil(t, insight)
	assert.Equal(t, "Success after retry", insight.Content)
	assert.Equal(t, 3, attemptCount) // Should have made 3 attempts
}

// TestConcurrentRequests tests concurrent access to the service
func TestConcurrentRequests(t *testing.T) {
	ctx := context.Background()

	logger := logging.NewTestLogger()

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response": "Concurrent response", "done": true}`))
	}))
	defer server.Close()

	service, err := NewOllamaService(ctx, server.URL, logger)
	require.NoError(t, err)

	// Run multiple concurrent requests
	const numRequests = 10
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(id int) {
			ctx := context.Background()
			prompt := fmt.Sprintf("Concurrent test prompt %d", id)
			insight, err := service.GenerateInsight(ctx, prompt)
			if err != nil {
				results <- err
				return
			}
			if insight == nil || insight.Content != "Concurrent response" {
				results <- fmt.Errorf("unexpected response for request %d", id)
				return
			}
			results <- nil
		}(i)
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		select {
		case err := <-results:
			assert.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Fatal("Test timed out waiting for concurrent requests")
		}
	}
}

// TestJSONSerialization tests JSON marshaling and unmarshaling of structs
func TestJSONSerialization(t *testing.T) {
	tests := []struct {
		name   string
		object any
	}{
		{
			name: "GenerateRequest",
			object: &GenerateRequest{
				Model:  "qwen2.5-coder:7b",
				Prompt: "test prompt",
				Stream: false,
			},
		},
		{
			name: "GenerateResponse",
			object: &GenerateResponse{
				Response: "test response",
				Done:     true,
			},
		},
		{
			name: "Insight",
			object: &Insight{
				Content:    "test content",
				Tags:       []string{"tag1", "tag2"},
				Confidence: 0.85,
			},
		},
		{
			name: "InsightRequest",
			object: &InsightRequest{
				Prompt:      "test prompt",
				UserID:      "user123",
				EntryIDs:    []string{"entry1", "entry2"},
				InsightType: "productivity",
				Context:     map[string]any{"key": "value"},
			},
		},
		{
			name: "WeeklyReportRequest",
			object: &WeeklyReportRequest{
				UserID:    "user123",
				WeekStart: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				WeekEnd:   time.Date(2025, 1, 7, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name: "WeeklyReport",
			object: &WeeklyReport{
				Summary:         "test summary",
				KeyInsights:     []string{"insight1", "insight2"},
				Recommendations: []string{"rec1", "rec2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(tt.object)
			assert.NoError(t, err)
			assert.NotEmpty(t, jsonData)

			// Verify it's valid JSON
			var temp any
			err = json.Unmarshal(jsonData, &temp)
			assert.NoError(t, err)
		})
	}
}

// BenchmarkGenerateInsight benchmarks the insight generation performance
func BenchmarkGenerateInsight(b *testing.B) {
	ctx := context.Background()

	logger := logging.NewTestLogger()

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response": "Benchmark response", "done": true}`))
	}))
	defer server.Close()

	service, err := NewOllamaService(ctx, server.URL, logger)
	require.NoError(b, err)

	prompt := "Benchmark test prompt for performance measurement"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateInsight(ctx, prompt)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBuildEnhancedPrompt benchmarks the prompt enhancement performance
func BenchmarkBuildEnhancedPrompt(b *testing.B) {
	ctx := context.Background()
	logger := logging.NewTestLogger()
	service, err := NewOllamaService(ctx, "http://localhost:11434", logger)
	require.NoError(b, err)

	request := &InsightRequest{
		Prompt: "Base prompt for benchmark testing with various context types",
		Context: map[string]any{
			"key1":     "value1",
			"key2":     42,
			"key3":     []string{"item1", "item2", "item3"},
			"metadata": map[string]any{"nested": "data"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.buildEnhancedPrompt(ctx, request)
	}
}

// TestEdgeCases tests various edge cases and error conditions
func TestEdgeCases(t *testing.T) {
	ctx := context.Background()

	logger := logging.NewTestLogger()

	t.Run("Very long prompt", func(t *testing.T) {
		// Create a very long prompt
		longPrompt := strings.Repeat("This is a very long prompt. ", 1000)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response": "Long prompt response", "done": true}`))
		}))
		defer server.Close()

		service, err := NewOllamaService(ctx, server.URL, logger)
		require.NoError(t, err)

		ctx := context.Background()
		insight, err := service.GenerateInsight(ctx, longPrompt)
		assert.NoError(t, err)
		assert.NotNil(t, insight)
	})

	t.Run("Complex nested context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response": "Complex context response", "done": true}`))
		}))
		defer server.Close()

		service, err := NewOllamaService(ctx, server.URL, logger)
		require.NoError(t, err)

		complexContext := map[string]any{
			"level1": map[string]any{
				"level2": map[string]any{
					"level3": []any{
						map[string]any{"deep": "data"},
						[]string{"nested", "array"},
					},
				},
			},
		}

		request := &InsightRequest{
			Prompt:      "Test prompt",
			UserID:      "user123",
			EntryIDs:    []string{"entry1"},
			InsightType: "productivity",
			Context:     complexContext,
		}

		ctx := context.Background()
		insight, err := service.GenerateInsightWithContext(ctx, request)
		assert.NoError(t, err)
		assert.NotNil(t, insight)
	})

	t.Run("Malformed server response", func(t *testing.T) {
		ctx := context.Background()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response": "incomplete`)) // Malformed JSON
		}))
		defer server.Close()

		service, err := NewOllamaService(ctx, server.URL, logger)
		require.NoError(t, err)

		insight, err := service.GenerateInsight(ctx, "test prompt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode response")
		assert.Nil(t, insight)
	})
}
