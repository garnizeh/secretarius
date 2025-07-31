package services

import (
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestLogEntryService_ValidationMethods tests validation functions
func TestLogEntryService_ValidationMethods(t *testing.T) {
	// Criar logger de teste adequado
	logger := logging.NewTestLogger()

	// Criar service com dependências adequadas (usando nil para DB pois só testamos validação)
	logEntryService := NewLogEntryService(nil, logger)

	t.Run("validateLogEntryRequest", func(t *testing.T) {
		now := time.Now()
		later := now.Add(1 * time.Hour)
		validProjectID := uuid.New()

		tests := []struct {
			name    string
			req     *models.LogEntryRequest
			wantErr bool
			errMsg  string
		}{
			{
				name: "valid request minimal",
				req: &models.LogEntryRequest{
					Title:       "Test Entry",
					Type:        models.ActivityDevelopment,
					StartTime:   now,
					EndTime:     later,
					ValueRating: models.ValueMedium,
					ImpactLevel: models.ImpactPersonal,
				},
				wantErr: false,
			},
			{
				name: "valid request complete",
				req: &models.LogEntryRequest{
					ProjectID:   &validProjectID,
					Title:       "Complete Test Entry",
					Description: stringPtr("A complete test entry"),
					Type:        models.ActivityDevelopment,
					StartTime:   now,
					EndTime:     later,
					ValueRating: models.ValueHigh,
					ImpactLevel: models.ImpactCompany,
					Tags:        []string{"test", "work"},
				},
				wantErr: false,
			},
			{
				name: "missing title",
				req: &models.LogEntryRequest{
					Title:       "",
					Type:        models.ActivityDevelopment,
					StartTime:   now,
					EndTime:     later,
					ValueRating: models.ValueMedium,
					ImpactLevel: models.ImpactPersonal,
				},
				wantErr: true,
				errMsg:  "title is required",
			},
			{
				name: "title too long",
				req: &models.LogEntryRequest{
					Title:       stringRepeat("a", 201),
					Type:        models.ActivityDevelopment,
					StartTime:   now,
					EndTime:     later,
					ValueRating: models.ValueMedium,
					ImpactLevel: models.ImpactPersonal,
				},
				wantErr: true,
				errMsg:  "title must be at most 200 characters",
			},
			{
				name: "invalid activity type",
				req: &models.LogEntryRequest{
					Title:       "Test Entry",
					Type:        "invalid-type",
					StartTime:   now,
					EndTime:     later,
					ValueRating: models.ValueMedium,
					ImpactLevel: models.ImpactPersonal,
				},
				wantErr: true,
				errMsg:  "invalid activity type",
			},
			{
				name: "invalid value rating",
				req: &models.LogEntryRequest{
					Title:       "Test Entry",
					Type:        models.ActivityDevelopment,
					StartTime:   now,
					EndTime:     later,
					ValueRating: "invalid-rating",
					ImpactLevel: models.ImpactPersonal,
				},
				wantErr: true,
				errMsg:  "invalid value rating",
			},
			{
				name: "invalid impact level",
				req: &models.LogEntryRequest{
					Title:       "Test Entry",
					Type:        models.ActivityDevelopment,
					StartTime:   now,
					EndTime:     later,
					ValueRating: models.ValueMedium,
					ImpactLevel: "invalid-impact",
				},
				wantErr: true,
				errMsg:  "invalid impact level",
			},
			{
				name: "end time before start time",
				req: &models.LogEntryRequest{
					Title:       "Test Entry",
					Type:        models.ActivityDevelopment,
					StartTime:   later,
					EndTime:     now,
					ValueRating: models.ValueMedium,
					ImpactLevel: models.ImpactPersonal,
				},
				wantErr: true,
				errMsg:  "end time must be after start time",
			},
			{
				name: "same start and end time",
				req: &models.LogEntryRequest{
					Title:       "Test Entry",
					Type:        models.ActivityDevelopment,
					StartTime:   now,
					EndTime:     now,
					ValueRating: models.ValueMedium,
					ImpactLevel: models.ImpactPersonal,
				},
				wantErr: true,
				errMsg:  "end time must be after start time",
			},
			{
				name: "description too long",
				req: &models.LogEntryRequest{
					Title:       "Test Entry",
					Description: stringPtr(stringRepeat("a", 1001)),
					Type:        models.ActivityDevelopment,
					StartTime:   now,
					EndTime:     later,
					ValueRating: models.ValueMedium,
					ImpactLevel: models.ImpactPersonal,
				},
				wantErr: true,
				errMsg:  "description must be at most 1000 characters",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := logEntryService.validateLogEntryRequest(tt.req)
				if tt.wantErr {
					assert.Error(t, err)
					if err != nil && tt.errMsg != "" {
						assert.Contains(t, err.Error(), tt.errMsg)
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestLogEntryService_ActivityTypeValidation(t *testing.T) {
	tests := []struct {
		name    string
		actType models.ActivityType
		isValid bool
	}{
		{
			name:    "development activity",
			actType: models.ActivityDevelopment,
			isValid: true,
		},
		{
			name:    "meeting activity",
			actType: models.ActivityMeeting,
			isValid: true,
		},
		{
			name:    "learning activity",
			actType: models.ActivityLearning,
			isValid: true,
		},
		{
			name:    "debugging activity",
			actType: models.ActivityDebugging,
			isValid: true,
		},
		{
			name:    "other activity",
			actType: models.ActivityOther,
			isValid: true,
		},
		{
			name:    "invalid activity",
			actType: "invalid",
			isValid: false,
		},
		{
			name:    "empty activity",
			actType: "",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isValid, tt.actType.IsValid())
		})
	}
}

func TestLogEntryService_ValueRatingValidation(t *testing.T) {
	tests := []struct {
		name    string
		rating  models.ValueRating
		isValid bool
	}{
		{
			name:    "low value",
			rating:  models.ValueLow,
			isValid: true,
		},
		{
			name:    "medium value",
			rating:  models.ValueMedium,
			isValid: true,
		},
		{
			name:    "high value",
			rating:  models.ValueHigh,
			isValid: true,
		},
		{
			name:    "critical value",
			rating:  models.ValueCritical,
			isValid: true,
		},
		{
			name:    "invalid value",
			rating:  "invalid",
			isValid: false,
		},
		{
			name:    "empty value",
			rating:  "",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isValid, tt.rating.IsValid())
		})
	}
}

func TestLogEntryService_ImpactLevelValidation(t *testing.T) {
	tests := []struct {
		name    string
		impact  models.ImpactLevel
		isValid bool
	}{
		{
			name:    "personal impact",
			impact:  models.ImpactPersonal,
			isValid: true,
		},
		{
			name:    "team impact",
			impact:  models.ImpactTeam,
			isValid: true,
		},
		{
			name:    "department impact",
			impact:  models.ImpactDepartment,
			isValid: true,
		},
		{
			name:    "company impact",
			impact:  models.ImpactCompany,
			isValid: true,
		},
		{
			name:    "invalid impact",
			impact:  "invalid",
			isValid: false,
		},
		{
			name:    "empty impact",
			impact:  "",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isValid, tt.impact.IsValid())
		})
	}
}
