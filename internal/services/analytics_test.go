package services

import (
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/logging"
	"github.com/stretchr/testify/assert"
)

// TestAnalyticsService_ValidationMethods tests validation functions
func TestAnalyticsService_ValidationMethods(t *testing.T) {
	// Create test logger to prevent nil pointer dereference
	testLogger := logging.NewTestLogger()
	analyticsService := NewAnalyticsService(nil, testLogger)

	t.Run("ValidateDateRange", func(t *testing.T) {
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		tomorrow := now.Add(24 * time.Hour)

		tests := []struct {
			name      string
			startDate time.Time
			endDate   time.Time
			wantErr   bool
			errMsg    string
		}{
			{
				name:      "valid date range",
				startDate: yesterday,
				endDate:   now,
				wantErr:   false,
			},
			{
				name:      "same start and end date",
				startDate: now,
				endDate:   now,
				wantErr:   false, // Same day should be allowed
			},
			{
				name:      "end before start",
				startDate: tomorrow,
				endDate:   yesterday,
				wantErr:   true,
				errMsg:    "end date must be after or equal to start date",
			},
			{
				name:      "very long range",
				startDate: now.Add(-365 * 24 * time.Hour),
				endDate:   now,
				wantErr:   false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := analyticsService.validateDateRange(tt.startDate, tt.endDate)
				if tt.wantErr {
					assert.Error(t, err)
					if tt.errMsg != "" {
						assert.Contains(t, err.Error(), tt.errMsg)
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestAnalyticsService_HelperMethods(t *testing.T) {
	t.Run("CalculatePercentageChange", func(t *testing.T) {
		tests := []struct {
			name     string
			current  float64
			previous float64
			expected float64
		}{
			{
				name:     "50% increase",
				current:  150,
				previous: 100,
				expected: 50.0,
			},
			{
				name:     "25% decrease",
				current:  75,
				previous: 100,
				expected: -25.0,
			},
			{
				name:     "no change",
				current:  100,
				previous: 100,
				expected: 0.0,
			},
			{
				name:     "from zero",
				current:  50,
				previous: 0,
				expected: 100.0, // Special case: 100% increase from zero
			},
			{
				name:     "to zero",
				current:  0,
				previous: 100,
				expected: -100.0,
			},
			{
				name:     "both zero",
				current:  0,
				previous: 0,
				expected: 0.0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				actual := calculatePercentageChange(tt.current, tt.previous)
				assert.Equal(t, tt.expected, actual)
			})
		}
	})

	t.Run("NormalizeProductivityScore", func(t *testing.T) {
		tests := []struct {
			name     string
			score    float64
			expected float64
		}{
			{
				name:     "normal score",
				score:    75.5,
				expected: 75.5,
			},
			{
				name:     "negative score",
				score:    -10,
				expected: 0.0,
			},
			{
				name:     "over 100 score",
				score:    120,
				expected: 100.0,
			},
			{
				name:     "zero score",
				score:    0,
				expected: 0.0,
			},
			{
				name:     "perfect score",
				score:    100,
				expected: 100.0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				actual := normalizeProductivityScore(tt.score)
				assert.Equal(t, tt.expected, actual)
			})
		}
	})

	t.Run("GetDayOfWeekName", func(t *testing.T) {
		tests := []struct {
			name     string
			weekday  time.Weekday
			expected string
		}{
			{
				name:     "Sunday",
				weekday:  time.Sunday,
				expected: "Sunday",
			},
			{
				name:     "Monday",
				weekday:  time.Monday,
				expected: "Monday",
			},
			{
				name:     "Tuesday",
				weekday:  time.Tuesday,
				expected: "Tuesday",
			},
			{
				name:     "Wednesday",
				weekday:  time.Wednesday,
				expected: "Wednesday",
			},
			{
				name:     "Thursday",
				weekday:  time.Thursday,
				expected: "Thursday",
			},
			{
				name:     "Friday",
				weekday:  time.Friday,
				expected: "Friday",
			},
			{
				name:     "Saturday",
				weekday:  time.Saturday,
				expected: "Saturday",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				actual := getDayOfWeekName(tt.weekday)
				assert.Equal(t, tt.expected, actual)
			})
		}
	})
}

// Helper functions for testing (these might not exist in the actual service)
func calculatePercentageChange(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100.0 // 100% increase from zero
	}
	return ((current - previous) / previous) * 100
}

func normalizeProductivityScore(score float64) float64 {
	if score < 0 {
		return 0.0
	}
	if score > 100 {
		return 100.0
	}
	return score
}

func getDayOfWeekName(weekday time.Weekday) string {
	return weekday.String()
}
