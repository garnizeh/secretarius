package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActivityType_IsValid(t *testing.T) {
	tests := []struct {
		name         string
		activityType ActivityType
		want         bool
	}{
		{"valid development", ActivityDevelopment, true},
		{"valid meeting", ActivityMeeting, true},
		{"valid code review", ActivityCodeReview, true},
		{"valid debugging", ActivityDebugging, true},
		{"valid documentation", ActivityDocumentation, true},
		{"valid testing", ActivityTesting, true},
		{"valid deployment", ActivityDeployment, true},
		{"valid research", ActivityResearch, true},
		{"valid planning", ActivityPlanning, true},
		{"valid learning", ActivityLearning, true},
		{"valid maintenance", ActivityMaintenance, true},
		{"valid support", ActivitySupport, true},
		{"valid other", ActivityOther, true},
		{"invalid type", ActivityType("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.activityType.IsValid())
		})
	}
}

func TestValueRating_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		rating ValueRating
		want   bool
	}{
		{"valid low", ValueLow, true},
		{"valid medium", ValueMedium, true},
		{"valid high", ValueHigh, true},
		{"valid critical", ValueCritical, true},
		{"invalid rating", ValueRating("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.rating.IsValid())
		})
	}
}

func TestImpactLevel_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		impact ImpactLevel
		want   bool
	}{
		{"valid personal", ImpactPersonal, true},
		{"valid team", ImpactTeam, true},
		{"valid department", ImpactDepartment, true},
		{"valid company", ImpactCompany, true},
		{"invalid impact", ImpactLevel("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.impact.IsValid())
		})
	}
}

func TestLogEntry_CalculateDuration(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
	endTime := time.Date(2024, 1, 1, 11, 30, 0, 0, time.UTC)

	logEntry := &LogEntry{
		StartTime: startTime,
		EndTime:   endTime,
	}

	logEntry.CalculateDuration()

	expected := 150 // 2.5 hours * 60 minutes
	assert.Equal(t, expected, logEntry.DurationMinutes)
}

func TestLogEntry_Validate(t *testing.T) {
	tests := []struct {
		name     string
		logEntry LogEntry
		wantErr  bool
	}{
		{
			name: "valid log entry",
			logEntry: LogEntry{
				StartTime:   time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
				Type:        ActivityDevelopment,
				ValueRating: ValueHigh,
				ImpactLevel: ImpactTeam,
			},
			wantErr: false,
		},
		{
			name: "invalid time range",
			logEntry: LogEntry{
				StartTime:   time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				Type:        ActivityDevelopment,
				ValueRating: ValueHigh,
				ImpactLevel: ImpactTeam,
			},
			wantErr: true,
		},
		{
			name: "invalid activity type",
			logEntry: LogEntry{
				StartTime:   time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
				Type:        ActivityType("invalid"),
				ValueRating: ValueHigh,
				ImpactLevel: ImpactTeam,
			},
			wantErr: true,
		},
		{
			name: "invalid value rating",
			logEntry: LogEntry{
				StartTime:   time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
				Type:        ActivityDevelopment,
				ValueRating: ValueRating("invalid"),
				ImpactLevel: ImpactTeam,
			},
			wantErr: true,
		},
		{
			name: "invalid impact level",
			logEntry: LogEntry{
				StartTime:   time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
				Type:        ActivityDevelopment,
				ValueRating: ValueHigh,
				ImpactLevel: ImpactLevel("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.logEntry.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
