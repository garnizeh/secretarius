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

// Helper function to parse time strings for tests
func parseTime(t *testing.T, timeStr string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		t.Fatalf("Failed to parse time %s: %v", timeStr, err)
	}
	return parsed
}
