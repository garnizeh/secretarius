package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateTimezone(t *testing.T) {
	tests := []struct {
		name     string
		timezone string
		wantErr  bool
	}{
		{"valid UTC", "UTC", false},
		{"valid America/New_York", "America/New_York", false},
		{"valid Europe/London", "Europe/London", false},
		{"valid Asia/Tokyo", "Asia/Tokyo", false},
		{"invalid timezone", "Invalid/Timezone", true},
		{"empty timezone defaults to UTC", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTimezone(tt.timezone)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidTimezone, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateHexColor(t *testing.T) {
	tests := []struct {
		name    string
		color   string
		wantErr bool
	}{
		{"valid color lowercase", "#ff0000", false},
		{"valid color uppercase", "#FF0000", false},
		{"valid color mixed case", "#FfA500", false},
		{"invalid without hash", "ff0000", true},
		{"invalid short format", "#fff", true},
		{"invalid long format", "#ff00000", true},
		{"invalid characters", "#gggggg", true},
		{"empty color", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHexColor(tt.color)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidColor, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTimeRange(t *testing.T) {
	tests := []struct {
		name    string
		start   string
		end     string
		wantErr bool
	}{
		{"valid range", "2024-01-01T09:00:00Z", "2024-01-01T17:00:00Z", false},
		{"same time", "2024-01-01T09:00:00Z", "2024-01-01T09:00:00Z", false},
		{"invalid range", "2024-01-01T17:00:00Z", "2024-01-01T09:00:00Z", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := parseTime(t, tt.start)
			end := parseTime(t, tt.end)

			err := ValidateTimeRange(start, end)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidTimeRange, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateActivityType(t *testing.T) {
	tests := []struct {
		name         string
		activityType string
		expected     bool
	}{
		{
			name:         "empty activity type (optional)",
			activityType: "",
			expected:     true,
		},
		{
			name:         "valid activity type - development",
			activityType: "development",
			expected:     true,
		},
		{
			name:         "valid activity type - meeting",
			activityType: "meeting",
			expected:     true,
		},
		{
			name:         "valid activity type - code_review",
			activityType: "code_review",
			expected:     true,
		},
		{
			name:         "valid activity type - debugging",
			activityType: "debugging",
			expected:     true,
		},
		{
			name:         "valid activity type - research",
			activityType: "research",
			expected:     true,
		},
		{
			name:         "valid activity type - testing",
			activityType: "testing",
			expected:     true,
		},
		{
			name:         "valid activity type - documentation",
			activityType: "documentation",
			expected:     true,
		},
		{
			name:         "valid activity type - deployment",
			activityType: "deployment",
			expected:     true,
		},
		{
			name:         "valid activity type - learning",
			activityType: "learning",
			expected:     true,
		},
		{
			name:         "valid activity type - planning",
			activityType: "planning",
			expected:     true,
		},
		{
			name:         "valid activity type - maintenance",
			activityType: "maintenance",
			expected:     true,
		},
		{
			name:         "valid activity type - support",
			activityType: "support",
			expected:     true,
		},
		{
			name:         "valid activity type - other",
			activityType: "other",
			expected:     true,
		},
		{
			name:         "invalid activity type",
			activityType: "invalid_type",
			expected:     false,
		},
		{
			name:         "case sensitive validation",
			activityType: "DEVELOPMENT",
			expected:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateActivityType(tc.activityType)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidateValueRating(t *testing.T) {
	tests := []struct {
		name        string
		valueRating string
		expected    bool
	}{
		{
			name:        "empty value rating (optional)",
			valueRating: "",
			expected:    true,
		},
		{
			name:        "valid value rating - low",
			valueRating: "low",
			expected:    true,
		},
		{
			name:        "valid value rating - medium",
			valueRating: "medium",
			expected:    true,
		},
		{
			name:        "valid value rating - high",
			valueRating: "high",
			expected:    true,
		},
		{
			name:        "valid value rating - critical",
			valueRating: "critical",
			expected:    true,
		},
		{
			name:        "invalid value rating",
			valueRating: "invalid_rating",
			expected:    false,
		},
		{
			name:        "case sensitive validation",
			valueRating: "HIGH",
			expected:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateValueRating(tc.valueRating)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidateImpactLevel(t *testing.T) {
	tests := []struct {
		name        string
		impactLevel string
		expected    bool
	}{
		{
			name:        "empty impact level (optional)",
			impactLevel: "",
			expected:    true,
		},
		{
			name:        "valid impact level - personal",
			impactLevel: "personal",
			expected:    true,
		},
		{
			name:        "valid impact level - team",
			impactLevel: "team",
			expected:    true,
		},
		{
			name:        "valid impact level - department",
			impactLevel: "department",
			expected:    true,
		},
		{
			name:        "valid impact level - company",
			impactLevel: "company",
			expected:    true,
		},
		{
			name:        "invalid impact level",
			impactLevel: "invalid_level",
			expected:    false,
		},
		{
			name:        "case sensitive validation",
			impactLevel: "PERSONAL",
			expected:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateImpactLevel(tc.impactLevel)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidateDateFormat(t *testing.T) {
	tests := []struct {
		name        string
		dateStr     string
		expectError bool
	}{
		{
			name:        "empty date (optional)",
			dateStr:     "",
			expectError: false,
		},
		{
			name:        "valid date format",
			dateStr:     "2023-12-25",
			expectError: false,
		},
		{
			name:        "valid date format - different date",
			dateStr:     "2025-01-01",
			expectError: false,
		},
		{
			name:        "invalid date format - wrong separator",
			dateStr:     "2023/12/25",
			expectError: true,
		},
		{
			name:        "invalid date format - missing day",
			dateStr:     "2023-12",
			expectError: true,
		},
		{
			name:        "invalid date format - wrong order",
			dateStr:     "25-12-2023",
			expectError: true,
		},
		{
			name:        "invalid date format - with time",
			dateStr:     "2023-12-25 10:30:00",
			expectError: true,
		},
		{
			name:        "invalid date - impossible date",
			dateStr:     "2023-02-30",
			expectError: true,
		},
		{
			name:        "invalid date - wrong month",
			dateStr:     "2023-13-01",
			expectError: true,
		},
		{
			name:        "invalid date format - text",
			dateStr:     "not-a-date",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateDateFormat(tc.dateStr)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to parse time strings for tests
func parseTime(t *testing.T, timeStr string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		t.Fatalf("Failed to parse time %s: %v", timeStr, err)
	}
	return parsed
}

// Performance/benchmark tests
func BenchmarkValidateActivityType(b *testing.B) {
	for range b.N {
		ValidateActivityType("coding")
	}
}

func BenchmarkValidateValueRating(b *testing.B) {
	for range b.N {
		ValidateValueRating("high")
	}
}

func BenchmarkValidateImpactLevel(b *testing.B) {
	for range b.N {
		ValidateImpactLevel("medium")
	}
}

func BenchmarkValidateDateFormat(b *testing.B) {
	for range b.N {
		_ = ValidateDateFormat("2023-12-25")
	}
}

// Edge case tests
func TestMiddlewareEdgeCases(t *testing.T) {
	t.Run("date validation edge cases", func(t *testing.T) {
		edgeCases := []struct {
			date        string
			expectError bool
			description string
		}{
			{"2023-02-29", true, "non-leap year February 29th"},
			{"2024-02-29", false, "leap year February 29th"},
			{"2023-04-31", true, "April 31st (April has 30 days)"},
			{"2023-12-31", false, "valid December 31st"},
			{"0000-01-01", false, "year zero"},
			{"9999-12-31", false, "year 9999"},
		}

		for _, tc := range edgeCases {
			t.Run(tc.description, func(t *testing.T) {
				err := ValidateDateFormat(tc.date)
				if tc.expectError {
					assert.Error(t, err, "Expected error for date: %s", tc.date)
				} else {
					assert.NoError(t, err, "Expected no error for date: %s", tc.date)
				}
			})
		}
	})
}
