package services

import (
	"testing"

	"github.com/garnizeh/englog/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestProjectService_ValidationMethods tests validation functions
func TestProjectService_ValidationMethods(t *testing.T) {
	projectService := &ProjectService{}

	t.Run("validateProjectRequest", func(t *testing.T) {
		tests := []struct {
			name    string
			req     *models.ProjectRequest
			wantErr bool
			errMsg  string
		}{
			{
				name: "valid request",
				req: &models.ProjectRequest{
					Name:   "Test Project",
					Color:  "#FF5733",
					Status: models.ProjectActive,
				},
				wantErr: false,
			},
			{
				name: "valid request with description",
				req: &models.ProjectRequest{
					Name:        "Test Project",
					Description: stringPtr("A test project"),
					Color:       "#FF5733",
					Status:      models.ProjectActive,
				},
				wantErr: false,
			},
			{
				name: "valid request with all fields",
				req: &models.ProjectRequest{
					Name:        "Complete Project",
					Description: stringPtr("A complete project"),
					Color:       "#00FF00",
					Status:      models.ProjectCompleted,
					IsDefault:   true,
				},
				wantErr: false,
			},
			{
				name: "missing name",
				req: &models.ProjectRequest{
					Name:   "",
					Color:  "#FF5733",
					Status: models.ProjectActive,
				},
				wantErr: true,
				errMsg:  "project name is required",
			},
			{
				name: "name too long",
				req: &models.ProjectRequest{
					Name:   stringRepeat("a", 201),
					Color:  "#FF5733",
					Status: models.ProjectActive,
				},
				wantErr: true,
				errMsg:  "project name must be at most 200 characters",
			},
			{
				name: "missing color",
				req: &models.ProjectRequest{
					Name:   "Test Project",
					Color:  "",
					Status: models.ProjectActive,
				},
				wantErr: true,
				errMsg:  "invalid color format",
			},
			{
				name: "invalid color format",
				req: &models.ProjectRequest{
					Name:   "Test Project",
					Color:  "invalid-color",
					Status: models.ProjectActive,
				},
				wantErr: true,
				errMsg:  "invalid color format",
			},
			{
				name: "invalid status",
				req: &models.ProjectRequest{
					Name:   "Test Project",
					Color:  "#FF5733",
					Status: "invalid-status",
				},
				wantErr: true,
				errMsg:  "invalid project status",
			},
			{
				name: "description too long",
				req: &models.ProjectRequest{
					Name:        "Test Project",
					Description: stringPtr(stringRepeat("a", 1001)),
					Color:       "#FF5733",
					Status:      models.ProjectActive,
				},
				wantErr: true,
				errMsg:  "project description must be at most 1000 characters",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := projectService.validateProjectRequest(tt.req)
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

func TestProjectService_ProjectStatusValidation(t *testing.T) {
	tests := []struct {
		name    string
		status  models.ProjectStatus
		isValid bool
	}{
		{
			name:    "active status",
			status:  models.ProjectActive,
			isValid: true,
		},
		{
			name:    "completed status",
			status:  models.ProjectCompleted,
			isValid: true,
		},
		{
			name:    "on_hold status",
			status:  models.ProjectOnHold,
			isValid: true,
		},
		{
			name:    "cancelled status",
			status:  models.ProjectCancelled,
			isValid: true,
		},
		{
			name:    "invalid status",
			status:  "invalid",
			isValid: false,
		},
		{
			name:    "empty status",
			status:  "",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isValid, tt.status.IsValid())
		})
	}
}
