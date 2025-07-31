package services

import (
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// createTestLoggerForValidation creates a minimal logger for validation tests
func createTestLoggerForValidation() *logging.Logger {
	return logging.NewTestLogger()
}

// TestLogEntryService_ComprehensiveValidation tests comprehensive business logic validation
func TestLogEntryService_ComprehensiveValidation(t *testing.T) {
	// Criar logger de teste adequado
	logger := createTestLoggerForValidation()

	// Criar service com dependências adequadas (usando nil para DB pois só testamos validação)
	logEntryService := NewLogEntryService(nil, logger)

	t.Run("TitleValidationEdgeCases", func(t *testing.T) {
		now := time.Now()
		later := now.Add(1 * time.Hour)

		tests := []struct {
			name    string
			title   string
			wantErr bool
			errMsg  string
		}{
			{
				name:    "empty title",
				title:   "",
				wantErr: true,
				errMsg:  "title is required",
			},
			{
				name:    "whitespace only title",
				title:   "   ",
				wantErr: true,
				errMsg:  "title cannot be whitespace only",
			},
			{
				name:    "single character title",
				title:   "A",
				wantErr: false,
			},
			{
				name:    "maximum length title",
				title:   stringRepeat("A", 200),
				wantErr: false,
			},
			{
				name:    "title too long",
				title:   stringRepeat("A", 201),
				wantErr: true,
				errMsg:  "title must be at most 200 characters",
			},
			{
				name:    "title with special characters",
				title:   "Test Entry @#$%^&*()_+-=[]{}|;':\",./<>?",
				wantErr: false,
			},
			{
				name:    "title with unicode",
				title:   "Test Entry 测试 🚀",
				wantErr: false,
			},
			{
				name:    "title with newlines",
				title:   "Test\nEntry\nWith\nNewlines",
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := &models.LogEntryRequest{
					Title:       tt.title,
					Type:        models.ActivityDevelopment,
					StartTime:   now,
					EndTime:     later,
					ValueRating: models.ValueMedium,
					ImpactLevel: models.ImpactPersonal,
				}

				err := logEntryService.validateLogEntryRequest(req)
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

	t.Run("DescriptionValidationEdgeCases", func(t *testing.T) {
		now := time.Now()
		later := now.Add(1 * time.Hour)

		tests := []struct {
			name        string
			description *string
			wantErr     bool
			errMsg      string
		}{
			{
				name:        "nil description",
				description: nil,
				wantErr:     false,
			},
			{
				name:        "empty description",
				description: stringPtr(""),
				wantErr:     false,
			},
			{
				name:        "whitespace only description",
				description: stringPtr("   "),
				wantErr:     false,
			},
			{
				name:        "maximum length description",
				description: stringPtr(stringRepeat("A", 1000)),
				wantErr:     false,
			},
			{
				name:        "description too long",
				description: stringPtr(stringRepeat("A", 1001)),
				wantErr:     true,
				errMsg:      "description must be at most 1000 characters",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := &models.LogEntryRequest{
					Title:       "Valid Title",
					Description: tt.description,
					Type:        models.ActivityDevelopment,
					StartTime:   now,
					EndTime:     later,
					ValueRating: models.ValueMedium,
					ImpactLevel: models.ImpactPersonal,
				}

				err := logEntryService.validateLogEntryRequest(req)
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

	t.Run("TimeValidationEdgeCases", func(t *testing.T) {
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		tests := []struct {
			name      string
			startTime time.Time
			endTime   time.Time
			wantErr   bool
			errMsg    string
		}{
			{
				name:      "valid time range",
				startTime: baseTime,
				endTime:   baseTime.Add(1 * time.Hour),
				wantErr:   false,
			},
			{
				name:      "same start and end time",
				startTime: baseTime,
				endTime:   baseTime,
				wantErr:   true,
				errMsg:    "end time must be after start time",
			},
			{
				name:      "end before start",
				startTime: baseTime,
				endTime:   baseTime.Add(-1 * time.Hour),
				wantErr:   true,
				errMsg:    "end time must be after start time",
			},
			{
				name:      "1 minute duration",
				startTime: baseTime,
				endTime:   baseTime.Add(1 * time.Minute),
				wantErr:   false,
			},
			{
				name:      "1 second duration",
				startTime: baseTime,
				endTime:   baseTime.Add(1 * time.Second),
				wantErr:   false,
			},
			{
				name:      "24 hour duration",
				startTime: baseTime,
				endTime:   baseTime.Add(24 * time.Hour),
				wantErr:   false,
			},
			{
				name:      "very long duration",
				startTime: baseTime,
				endTime:   baseTime.Add(365 * 24 * time.Hour),
				wantErr:   false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := &models.LogEntryRequest{
					Title:       "Valid Title",
					Type:        models.ActivityDevelopment,
					StartTime:   tt.startTime,
					EndTime:     tt.endTime,
					ValueRating: models.ValueMedium,
					ImpactLevel: models.ImpactPersonal,
				}

				err := logEntryService.validateLogEntryRequest(req)
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

	t.Run("ActivityTypeValidationEdgeCases", func(t *testing.T) {
		now := time.Now()
		later := now.Add(1 * time.Hour)

		tests := []struct {
			name    string
			actType models.ActivityType
			wantErr bool
			errMsg  string
		}{
			{
				name:    "development type",
				actType: models.ActivityDevelopment,
				wantErr: false,
			},
			{
				name:    "meeting type",
				actType: models.ActivityMeeting,
				wantErr: false,
			},
			{
				name:    "code review type",
				actType: models.ActivityCodeReview,
				wantErr: false,
			},
			{
				name:    "debugging type",
				actType: models.ActivityDebugging,
				wantErr: false,
			},
			{
				name:    "documentation type",
				actType: models.ActivityDocumentation,
				wantErr: false,
			},
			{
				name:    "testing type",
				actType: models.ActivityTesting,
				wantErr: false,
			},
			{
				name:    "deployment type",
				actType: models.ActivityDeployment,
				wantErr: false,
			},
			{
				name:    "research type",
				actType: models.ActivityResearch,
				wantErr: false,
			},
			{
				name:    "planning type",
				actType: models.ActivityPlanning,
				wantErr: false,
			},
			{
				name:    "learning type",
				actType: models.ActivityLearning,
				wantErr: false,
			},
			{
				name:    "maintenance type",
				actType: models.ActivityMaintenance,
				wantErr: false,
			},
			{
				name:    "support type",
				actType: models.ActivitySupport,
				wantErr: false,
			},
			{
				name:    "other type",
				actType: models.ActivityOther,
				wantErr: false,
			},
			{
				name:    "invalid type",
				actType: "invalid-activity",
				wantErr: true,
				errMsg:  "invalid activity type",
			},
			{
				name:    "empty type",
				actType: "",
				wantErr: true,
				errMsg:  "invalid activity type",
			},
			{
				name:    "case sensitive type",
				actType: "DEVELOPMENT",
				wantErr: true,
				errMsg:  "invalid activity type",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := &models.LogEntryRequest{
					Title:       "Valid Title",
					Type:        tt.actType,
					StartTime:   now,
					EndTime:     later,
					ValueRating: models.ValueMedium,
					ImpactLevel: models.ImpactPersonal,
				}

				err := logEntryService.validateLogEntryRequest(req)
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

	t.Run("ValueRatingValidationEdgeCases", func(t *testing.T) {
		now := time.Now()
		later := now.Add(1 * time.Hour)

		tests := []struct {
			name    string
			rating  models.ValueRating
			wantErr bool
			errMsg  string
		}{
			{
				name:    "low value",
				rating:  models.ValueLow,
				wantErr: false,
			},
			{
				name:    "medium value",
				rating:  models.ValueMedium,
				wantErr: false,
			},
			{
				name:    "high value",
				rating:  models.ValueHigh,
				wantErr: false,
			},
			{
				name:    "critical value",
				rating:  models.ValueCritical,
				wantErr: false,
			},
			{
				name:    "invalid value",
				rating:  "invalid-value",
				wantErr: true,
				errMsg:  "invalid value rating",
			},
			{
				name:    "empty value",
				rating:  "",
				wantErr: true,
				errMsg:  "invalid value rating",
			},
			{
				name:    "case sensitive value",
				rating:  "LOW",
				wantErr: true,
				errMsg:  "invalid value rating",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := &models.LogEntryRequest{
					Title:       "Valid Title",
					Type:        models.ActivityDevelopment,
					StartTime:   now,
					EndTime:     later,
					ValueRating: tt.rating,
					ImpactLevel: models.ImpactPersonal,
				}

				err := logEntryService.validateLogEntryRequest(req)
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

	t.Run("ImpactLevelValidationEdgeCases", func(t *testing.T) {
		now := time.Now()
		later := now.Add(1 * time.Hour)

		tests := []struct {
			name    string
			impact  models.ImpactLevel
			wantErr bool
			errMsg  string
		}{
			{
				name:    "personal impact",
				impact:  models.ImpactPersonal,
				wantErr: false,
			},
			{
				name:    "team impact",
				impact:  models.ImpactTeam,
				wantErr: false,
			},
			{
				name:    "department impact",
				impact:  models.ImpactDepartment,
				wantErr: false,
			},
			{
				name:    "company impact",
				impact:  models.ImpactCompany,
				wantErr: false,
			},
			{
				name:    "invalid impact",
				impact:  "invalid-impact",
				wantErr: true,
				errMsg:  "invalid impact level",
			},
			{
				name:    "empty impact",
				impact:  "",
				wantErr: true,
				errMsg:  "invalid impact level",
			},
			{
				name:    "case sensitive impact",
				impact:  "PERSONAL",
				wantErr: true,
				errMsg:  "invalid impact level",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := &models.LogEntryRequest{
					Title:       "Valid Title",
					Type:        models.ActivityDevelopment,
					StartTime:   now,
					EndTime:     later,
					ValueRating: models.ValueMedium,
					ImpactLevel: tt.impact,
				}

				err := logEntryService.validateLogEntryRequest(req)
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

	t.Run("ProjectIDValidationEdgeCases", func(t *testing.T) {
		now := time.Now()
		later := now.Add(1 * time.Hour)

		tests := []struct {
			name      string
			projectID *uuid.UUID
			wantErr   bool
		}{
			{
				name:      "nil project ID",
				projectID: nil,
				wantErr:   false,
			},
			{
				name:      "valid project ID",
				projectID: func() *uuid.UUID { id := uuid.New(); return &id }(),
				wantErr:   false,
			},
			{
				name:      "zero UUID",
				projectID: func() *uuid.UUID { id := uuid.Nil; return &id }(),
				wantErr:   false, // Allow zero UUID as it might be valid in some contexts
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := &models.LogEntryRequest{
					Title:       "Valid Title",
					Type:        models.ActivityDevelopment,
					ProjectID:   tt.projectID,
					StartTime:   now,
					EndTime:     later,
					ValueRating: models.ValueMedium,
					ImpactLevel: models.ImpactPersonal,
				}

				err := logEntryService.validateLogEntryRequest(req)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("TagsValidationEdgeCases", func(t *testing.T) {
		now := time.Now()
		later := now.Add(1 * time.Hour)

		tests := []struct {
			name    string
			tags    []string
			wantErr bool
		}{
			{
				name:    "nil tags",
				tags:    nil,
				wantErr: false,
			},
			{
				name:    "empty tags",
				tags:    []string{},
				wantErr: false,
			},
			{
				name:    "single tag",
				tags:    []string{"development"},
				wantErr: false,
			},
			{
				name:    "multiple tags",
				tags:    []string{"development", "feature", "frontend"},
				wantErr: false,
			},
			{
				name:    "tags with special characters",
				tags:    []string{"tag-1", "tag_2", "tag@3"},
				wantErr: false,
			},
			{
				name:    "tags with unicode",
				tags:    []string{"开发", "🚀", "测试"},
				wantErr: false,
			},
			{
				name:    "empty tag in list",
				tags:    []string{"development", "", "feature"},
				wantErr: false, // Allow empty tags for now
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := &models.LogEntryRequest{
					Title:       "Valid Title",
					Type:        models.ActivityDevelopment,
					StartTime:   now,
					EndTime:     later,
					ValueRating: models.ValueMedium,
					ImpactLevel: models.ImpactPersonal,
					Tags:        tt.tags,
				}

				err := logEntryService.validateLogEntryRequest(req)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestLogEntryService_DurationCalculation(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		startTime        time.Time
		endTime          time.Time
		expectedDuration int // in minutes
	}{
		{
			name:             "1 minute",
			startTime:        baseTime,
			endTime:          baseTime.Add(1 * time.Minute),
			expectedDuration: 1,
		},
		{
			name:             "30 minutes",
			startTime:        baseTime,
			endTime:          baseTime.Add(30 * time.Minute),
			expectedDuration: 30,
		},
		{
			name:             "1 hour",
			startTime:        baseTime,
			endTime:          baseTime.Add(1 * time.Hour),
			expectedDuration: 60,
		},
		{
			name:             "1.5 hours",
			startTime:        baseTime,
			endTime:          baseTime.Add(90 * time.Minute),
			expectedDuration: 90,
		},
		{
			name:             "8 hours",
			startTime:        baseTime,
			endTime:          baseTime.Add(8 * time.Hour),
			expectedDuration: 480,
		},
		{
			name:             "partial minute (30 seconds)",
			startTime:        baseTime,
			endTime:          baseTime.Add(30 * time.Second),
			expectedDuration: 0, // Should round down
		},
		{
			name:             "partial minute (90 seconds)",
			startTime:        baseTime,
			endTime:          baseTime.Add(90 * time.Second),
			expectedDuration: 1, // Should round down to 1 minute
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := tt.endTime.Sub(tt.startTime)
			actualMinutes := int(duration.Minutes())
			assert.Equal(t, tt.expectedDuration, actualMinutes)
		})
	}
}

func TestLogEntryService_BoundaryConditions(t *testing.T) {
	logEntryService := &LogEntryService{}

	t.Run("MaxFieldLengthCombinations", func(t *testing.T) {
		now := time.Now()
		later := now.Add(1 * time.Hour)

		// Test request with all maximum field lengths
		req := &models.LogEntryRequest{
			Title:       stringRepeat("A", 200),
			Description: stringPtr(stringRepeat("B", 1000)),
			Type:        models.ActivityDevelopment,
			StartTime:   now,
			EndTime:     later,
			ValueRating: models.ValueMedium,
			ImpactLevel: models.ImpactPersonal,
			Tags:        []string{stringRepeat("C", 50), stringRepeat("D", 50)},
		}

		err := logEntryService.validateLogEntryRequest(req)
		assert.NoError(t, err)
	})

	t.Run("EmptyOptionalFieldsCombinations", func(t *testing.T) {
		now := time.Now()
		later := now.Add(1 * time.Hour)

		// Test request with minimal required fields only
		req := &models.LogEntryRequest{
			Title:       "Minimal Entry",
			Type:        models.ActivityDevelopment,
			StartTime:   now,
			EndTime:     later,
			ValueRating: models.ValueMedium,
			ImpactLevel: models.ImpactPersonal,
		}

		err := logEntryService.validateLogEntryRequest(req)
		assert.NoError(t, err)
	})
}
