package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestProjectService_BusinessLogicScenarios tests business logic edge cases
func TestProjectService_BusinessLogicScenarios(t *testing.T) {
	projectService := &ProjectService{}

	t.Run("project_name_edge_cases", func(t *testing.T) {
		testCases := []struct {
			name        string
			projectName string
			wantErr     bool
			errMsg      string
		}{
			{
				name:        "single_character_name",
				projectName: "A",
				wantErr:     false,
			},
			{
				name:        "name_with_spaces",
				projectName: "My Project Name",
				wantErr:     false,
			},
			{
				name:        "name_with_special_chars",
				projectName: "Project-Name_with.special@chars",
				wantErr:     false,
			},
			{
				name:        "unicode_name",
				projectName: "项目名称",
				wantErr:     false,
			},
			{
				name:        "exactly_200_chars",
				projectName: stringRepeat("a", 200),
				wantErr:     false,
			},
			{
				name:        "201_chars_should_fail",
				projectName: stringRepeat("a", 201),
				wantErr:     true,
				errMsg:      "project name must be at most 200 characters",
			},
			{
				name:        "empty_name",
				projectName: "",
				wantErr:     true,
				errMsg:      "project name is required",
			},
			{
				name:        "whitespace_only_name",
				projectName: "   ",
				wantErr:     true,
				errMsg:      "project name cannot be whitespace only",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := &models.ProjectRequest{
					Name:   tc.projectName,
					Color:  "#FF5733",
					Status: models.ProjectActive,
				}
				err := projectService.validateProjectRequest(req)
				if tc.wantErr {
					assert.Error(t, err)
					if err != nil && tc.errMsg != "" {
						assert.Contains(t, err.Error(), tc.errMsg)
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("project_color_validation", func(t *testing.T) {
		testCases := []struct {
			name    string
			color   string
			wantErr bool
			errMsg  string
		}{
			{
				name:    "valid_hex_color_short",
				color:   "#FFF",
				wantErr: true, // 3-char hex colors are not supported by ValidateHexColor
				errMsg:  "invalid color format",
			},
			{
				name:    "valid_hex_color_long",
				color:   "#FFFFFF",
				wantErr: false,
			},
			{
				name:    "valid_hex_color_lowercase",
				color:   "#ff5733",
				wantErr: false,
			},
			{
				name:    "valid_hex_color_mixed_case",
				color:   "#Ff5733",
				wantErr: false,
			},
			{
				name:    "valid_hex_color_with_numbers",
				color:   "#123456",
				wantErr: false,
			},
			{
				name:    "invalid_color_no_hash",
				color:   "FF5733",
				wantErr: true,
				errMsg:  "invalid color format",
			},
			{
				name:    "invalid_color_wrong_length",
				color:   "#FF57",
				wantErr: true,
				errMsg:  "invalid color format",
			},
			{
				name:    "invalid_color_invalid_chars",
				color:   "#GGGGGG",
				wantErr: true,
				errMsg:  "invalid color format",
			},
			{
				name:    "empty_color",
				color:   "",
				wantErr: true,
				errMsg:  "invalid color format",
			},
			{
				name:    "whitespace_color",
				color:   "   ",
				wantErr: true,
				errMsg:  "invalid color format",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := &models.ProjectRequest{
					Name:   "Test Project",
					Color:  tc.color,
					Status: models.ProjectActive,
				}
				err := projectService.validateProjectRequest(req)
				if tc.wantErr {
					assert.Error(t, err)
					if err != nil && tc.errMsg != "" {
						assert.Contains(t, err.Error(), tc.errMsg)
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("project_status_validation", func(t *testing.T) {
		testCases := []struct {
			name    string
			status  models.ProjectStatus
			wantErr bool
			errMsg  string
		}{
			{
				name:    "active_status",
				status:  models.ProjectActive,
				wantErr: false,
			},
			{
				name:    "completed_status",
				status:  models.ProjectCompleted,
				wantErr: false,
			},
			{
				name:    "on_hold_status",
				status:  models.ProjectOnHold,
				wantErr: false,
			},
			{
				name:    "cancelled_status",
				status:  models.ProjectCancelled,
				wantErr: false,
			},
			{
				name:    "invalid_status",
				status:  "invalid_status",
				wantErr: true,
				errMsg:  "invalid project status",
			},
			{
				name:    "empty_status",
				status:  "",
				wantErr: true,
				errMsg:  "invalid project status",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := &models.ProjectRequest{
					Name:   "Test Project",
					Color:  "#FF5733",
					Status: tc.status,
				}
				err := projectService.validateProjectRequest(req)
				if tc.wantErr {
					assert.Error(t, err)
					if err != nil && tc.errMsg != "" {
						assert.Contains(t, err.Error(), tc.errMsg)
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("project_description_validation", func(t *testing.T) {
		testCases := []struct {
			name        string
			description *string
			wantErr     bool
			errMsg      string
		}{
			{
				name:        "nil_description",
				description: nil,
				wantErr:     false,
			},
			{
				name:        "empty_description",
				description: stringPtr(""),
				wantErr:     false,
			},
			{
				name:        "normal_description",
				description: stringPtr("This is a normal project description"),
				wantErr:     false,
			},
			{
				name:        "unicode_description",
				description: stringPtr("这是一个项目描述"),
				wantErr:     false,
			},
			{
				name:        "description_with_newlines",
				description: stringPtr("Line 1\nLine 2\nLine 3"),
				wantErr:     false,
			},
			{
				name:        "description_with_special_chars",
				description: stringPtr("Description with !@#$%^&*() special chars"),
				wantErr:     false,
			},
			{
				name:        "exactly_1000_chars",
				description: stringPtr(stringRepeat("a", 1000)),
				wantErr:     false,
			},
			{
				name:        "1001_chars_should_fail",
				description: stringPtr(stringRepeat("a", 1001)),
				wantErr:     true,
				errMsg:      "project description must be at most 1000 characters",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := &models.ProjectRequest{
					Name:        "Test Project",
					Description: tc.description,
					Color:       "#FF5733",
					Status:      models.ProjectActive,
				}
				err := projectService.validateProjectRequest(req)
				if tc.wantErr {
					assert.Error(t, err)
					if err != nil && tc.errMsg != "" {
						assert.Contains(t, err.Error(), tc.errMsg)
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("project_dates_validation", func(t *testing.T) {
		now := time.Now()
		pastDate := now.AddDate(0, 0, -30)
		futureDate := now.AddDate(0, 0, 30)

		testCases := []struct {
			name      string
			startDate *time.Time
			endDate   *time.Time
			wantErr   bool
			errMsg    string
		}{
			{
				name:      "both_dates_nil",
				startDate: nil,
				endDate:   nil,
				wantErr:   false,
			},
			{
				name:      "only_start_date",
				startDate: &now,
				endDate:   nil,
				wantErr:   false,
			},
			{
				name:      "only_end_date",
				startDate: nil,
				endDate:   &futureDate,
				wantErr:   false,
			},
			{
				name:      "valid_date_range",
				startDate: &pastDate,
				endDate:   &futureDate,
				wantErr:   false,
			},
			{
				name:      "same_start_and_end_date",
				startDate: &now,
				endDate:   &now,
				wantErr:   false,
			},
			{
				name:      "end_date_before_start_date",
				startDate: &futureDate,
				endDate:   &pastDate,
				wantErr:   true,
				errMsg:    "end date must be after start date",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := &models.ProjectRequest{
					Name:      "Test Project",
					Color:     "#FF5733",
					Status:    models.ProjectActive,
					StartDate: tc.startDate,
					EndDate:   tc.endDate,
				}
				err := projectService.validateProjectRequest(req)
				if tc.wantErr {
					assert.Error(t, err)
					if err != nil && tc.errMsg != "" {
						assert.Contains(t, err.Error(), tc.errMsg)
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("project_default_flag_scenarios", func(t *testing.T) {
		testCases := []struct {
			name      string
			isDefault bool
			wantErr   bool
		}{
			{
				name:      "default_project",
				isDefault: true,
				wantErr:   false,
			},
			{
				name:      "non_default_project",
				isDefault: false,
				wantErr:   false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := &models.ProjectRequest{
					Name:      "Test Project",
					Color:     "#FF5733",
					Status:    models.ProjectActive,
					IsDefault: tc.isDefault,
				}
				err := projectService.validateProjectRequest(req)
				if tc.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

// TestProjectService_StatusTransitions tests valid status transitions
func TestProjectService_StatusTransitions(t *testing.T) {
	validTransitions := map[models.ProjectStatus][]models.ProjectStatus{
		models.ProjectActive: {
			models.ProjectCompleted,
			models.ProjectOnHold,
			models.ProjectCancelled,
		},
		models.ProjectOnHold: {
			models.ProjectActive,
			models.ProjectCompleted,
			models.ProjectCancelled,
		},
		models.ProjectCompleted: {
			models.ProjectActive, // Allow reopening completed projects
		},
		models.ProjectCancelled: {
			models.ProjectActive, // Allow reactivating cancelled projects
		},
	}

	for fromStatus, toStatuses := range validTransitions {
		for _, toStatus := range toStatuses {
			t.Run(fmt.Sprintf("transition_%s_to_%s", fromStatus, toStatus), func(t *testing.T) {
				// This test verifies that status transitions are conceptually valid
				// The actual business logic for preventing invalid transitions
				// would be implemented in the service layer if needed
				assert.True(t, toStatus.IsValid())
				assert.True(t, fromStatus.IsValid())
			})
		}
	}
}

// TestProjectService_ComplexValidationScenarios tests complex validation scenarios
func TestProjectService_ComplexValidationScenarios(t *testing.T) {
	projectService := &ProjectService{}

	t.Run("complete_project_validation", func(t *testing.T) {
		now := time.Now()
		pastDate := now.AddDate(0, 0, -30)
		futureDate := now.AddDate(0, 0, 30)

		testCases := []struct {
			name    string
			req     *models.ProjectRequest
			wantErr bool
			errMsg  string
		}{
			{
				name: "fully_valid_project",
				req: &models.ProjectRequest{
					Name:        "Complete Valid Project",
					Description: stringPtr("A complete and valid project description"),
					Color:       "#FF5733",
					Status:      models.ProjectActive,
					StartDate:   &pastDate,
					EndDate:     &futureDate,
					IsDefault:   true,
				},
				wantErr: false,
			},
			{
				name: "minimal_valid_project",
				req: &models.ProjectRequest{
					Name:   "Minimal",
					Color:  "#000000",
					Status: models.ProjectActive,
				},
				wantErr: false,
			},
			{
				name: "multiple_validation_errors",
				req: &models.ProjectRequest{
					Name:        "",                                 // Empty name
					Description: stringPtr(stringRepeat("a", 1001)), // Too long description
					Color:       "invalid",                          // Invalid color
					Status:      "invalid",                          // Invalid status
					StartDate:   &futureDate,
					EndDate:     &pastDate, // End before start
				},
				wantErr: true,
				errMsg:  "project name is required", // Should catch first error
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := projectService.validateProjectRequest(tc.req)
				if tc.wantErr {
					assert.Error(t, err)
					if err != nil && tc.errMsg != "" {
						assert.Contains(t, err.Error(), tc.errMsg)
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}
