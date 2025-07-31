package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status ProjectStatus
		want   bool
	}{
		{"valid active", ProjectActive, true},
		{"valid completed", ProjectCompleted, true},
		{"valid on hold", ProjectOnHold, true},
		{"valid cancelled", ProjectCancelled, true},
		{"invalid status", ProjectStatus("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.IsValid())
		})
	}
}

func TestProject_Validate(t *testing.T) {
	tests := []struct {
		name    string
		project Project
		wantErr bool
	}{
		{
			name: "valid project",
			project: Project{
				Status: ProjectActive,
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			project: Project{
				Status: ProjectStatus("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.project.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
