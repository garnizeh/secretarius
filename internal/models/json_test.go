package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserJSON(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	user := User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Timezone:  "UTC",
		Preferences: map[string]any{
			"theme": "dark",
			"lang":  "en",
		},
		LastLoginAt: &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(user)
	require.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaled User
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Verify data integrity
	assert.Equal(t, user.ID, unmarshaled.ID)
	assert.Equal(t, user.Email, unmarshaled.Email)
	assert.Equal(t, user.FirstName, unmarshaled.FirstName)
	assert.Equal(t, user.LastName, unmarshaled.LastName)
	assert.Equal(t, user.Timezone, unmarshaled.Timezone)
	assert.Equal(t, user.Preferences, unmarshaled.Preferences)

	// Verify password hash is not in JSON
	assert.NotContains(t, string(jsonData), "password_hash")
}

func TestLogEntryJSON(t *testing.T) {
	entryID := uuid.New()
	userID := uuid.New()
	projectID := uuid.New()
	startTime := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
	endTime := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)

	logEntry := LogEntry{
		ID:              entryID,
		UserID:          userID,
		Title:           "Development Work",
		Type:            ActivityDevelopment,
		ProjectID:       &projectID,
		StartTime:       startTime,
		EndTime:         endTime,
		DurationMinutes: 120,
		ValueRating:     ValueHigh,
		ImpactLevel:     ImpactTeam,
		Tags:            []string{"backend", "api"},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(logEntry)
	require.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaled LogEntry
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Verify data integrity
	assert.Equal(t, logEntry.ID, unmarshaled.ID)
	assert.Equal(t, logEntry.UserID, unmarshaled.UserID)
	assert.Equal(t, logEntry.Title, unmarshaled.Title)
	assert.Equal(t, logEntry.Type, unmarshaled.Type)
	assert.Equal(t, logEntry.ProjectID, unmarshaled.ProjectID)
	assert.Equal(t, logEntry.ValueRating, unmarshaled.ValueRating)
	assert.Equal(t, logEntry.ImpactLevel, unmarshaled.ImpactLevel)
	assert.Equal(t, logEntry.Tags, unmarshaled.Tags)
}

func TestProjectJSON(t *testing.T) {
	projectID := uuid.New()
	createdBy := uuid.New()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	project := Project{
		ID:          projectID,
		Name:        "Test Project",
		Description: &[]string{"A test project"}[0],
		Color:       "#FF5733",
		Status:      ProjectActive,
		StartDate:   &startDate,
		CreatedBy:   createdBy,
		IsDefault:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(project)
	require.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaled Project
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Verify data integrity
	assert.Equal(t, project.ID, unmarshaled.ID)
	assert.Equal(t, project.Name, unmarshaled.Name)
	assert.Equal(t, project.Description, unmarshaled.Description)
	assert.Equal(t, project.Color, unmarshaled.Color)
	assert.Equal(t, project.Status, unmarshaled.Status)
	assert.Equal(t, project.StartDate, unmarshaled.StartDate)
	assert.Equal(t, project.CreatedBy, unmarshaled.CreatedBy)
	assert.Equal(t, project.IsDefault, unmarshaled.IsDefault)
}

func TestGeneratedInsightJSON(t *testing.T) {
	insightID := uuid.New()
	userID := uuid.New()
	periodStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, 1, 7, 23, 59, 59, 0, time.UTC)
	qualityScore := 0.85

	insight := GeneratedInsight{
		ID:           insightID,
		UserID:       userID,
		ReportType:   ReportWeeklySummary,
		PeriodStart:  periodStart,
		PeriodEnd:    periodEnd,
		Title:        "Weekly Summary",
		Content:      "This week you worked on...",
		Summary:      &[]string{"Great productivity this week!"}[0],
		Metadata:     map[string]any{"total_hours": 40.0, "projects": 3.0},
		QualityScore: &qualityScore,
		Status:       InsightActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(insight)
	require.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaled GeneratedInsight
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Verify data integrity
	assert.Equal(t, insight.ID, unmarshaled.ID)
	assert.Equal(t, insight.UserID, unmarshaled.UserID)
	assert.Equal(t, insight.ReportType, unmarshaled.ReportType)
	assert.Equal(t, insight.PeriodStart, unmarshaled.PeriodStart)
	assert.Equal(t, insight.PeriodEnd, unmarshaled.PeriodEnd)
	assert.Equal(t, insight.Title, unmarshaled.Title)
	assert.Equal(t, insight.Content, unmarshaled.Content)
	assert.Equal(t, insight.Summary, unmarshaled.Summary)
	assert.Equal(t, insight.Metadata, unmarshaled.Metadata)
	assert.Equal(t, insight.QualityScore, unmarshaled.QualityScore)
	assert.Equal(t, insight.Status, unmarshaled.Status)
}
