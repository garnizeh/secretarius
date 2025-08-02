package ai

import (
	"context"
	"testing"

	"github.com/garnizeh/englog/internal/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewOllamaService tests the creation of OllamaService instances with langchaingo
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
				assert.NotNil(t, service.logger)
				assert.NotNil(t, service.llm)
				assert.Equal(t, "qwen2.5-coder:7b", service.modelName)
			}
		})
	}
}

// TestInsightRequestPrompt tests the prompt generation with various contexts
func TestInsightRequestPrompt(t *testing.T) {
	ctx := context.Background()
	logger := logging.NewTestLogger()

	tests := []struct {
		name          string
		request       *InsightRequest
		expectEmpty   bool
		shouldContain []string
	}{
		{
			name: "Basic productivity insight",
			request: &InsightRequest{
				UserID:      "user123",
				EntryIDs:    []string{"entry1", "entry2"},
				InsightType: "productivity",
				Context:     "Simple string context",
			},
			expectEmpty: false,
			shouldContain: []string{
				"User ID: user123",
				"Insight Type: productivity",
				"Number of Log Entries: 2",
				"productivity",
				"efficiency patterns",
			},
		},
		{
			name: "Empty user ID",
			request: &InsightRequest{
				UserID:      "",
				EntryIDs:    []string{"entry1"},
				InsightType: "productivity",
			},
			expectEmpty: false,
			shouldContain: []string{
				"User ID: ",
				"Number of Log Entries: 1",
			},
		},
		{
			name: "Structured context",
			request: &InsightRequest{
				UserID:      "user123",
				EntryIDs:    []string{"entry1"},
				InsightType: "time_management",
				Context: map[string]any{
					"date_range": map[string]any{
						"start": "2025-08-01",
						"end":   "2025-08-07",
					},
					"focus_areas": []string{"meetings", "development"},
				},
			},
			expectEmpty: false,
			shouldContain: []string{
				"time_management",
				"Structured Context",
				"date_range",
				"focus_areas",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := tt.request.Prompt(ctx, logger)

			if tt.expectEmpty {
				assert.Empty(t, prompt)
			} else {
				assert.NotEmpty(t, prompt)
				for _, expected := range tt.shouldContain {
					assert.Contains(t, prompt, expected)
				}
			}
		})
	}
}

// TestValidateInsightRequest tests the validation logic
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
				UserID:      "user123",
				EntryIDs:    []string{"entry1", "entry2"},
				InsightType: "productivity",
				Context:     "test context",
			},
			expectError: false,
		},
		{
			name: "Empty user ID",
			request: &InsightRequest{
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

// TestJoinStrings tests the helper function
func TestJoinStrings(t *testing.T) {
	tests := []struct {
		name      string
		strs      []string
		separator string
		expected  string
	}{
		{
			name:      "Empty slice",
			strs:      []string{},
			separator: ", ",
			expected:  "",
		},
		{
			name:      "Single element",
			strs:      []string{"one"},
			separator: ", ",
			expected:  "one",
		},
		{
			name:      "Multiple elements",
			strs:      []string{"one", "two", "three"},
			separator: ", ",
			expected:  "one, two, three",
		},
		{
			name:      "Different separator",
			strs:      []string{"a", "b", "c"},
			separator: " | ",
			expected:  "a | b | c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := joinStrings(tt.strs, tt.separator)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Note: Integration tests that require actual Ollama service would be in a separate file
// or marked with build tags for integration testing only.
