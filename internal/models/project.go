package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ProjectStatus represents the status of a project
type ProjectStatus string

const (
	ProjectActive    ProjectStatus = "active"
	ProjectCompleted ProjectStatus = "completed"
	ProjectOnHold    ProjectStatus = "on_hold"
	ProjectCancelled ProjectStatus = "cancelled"
)

// IsValid checks if the ProjectStatus is valid
func (p ProjectStatus) IsValid() bool {
	switch p {
	case ProjectActive, ProjectCompleted, ProjectOnHold, ProjectCancelled:
		return true
	}
	return false
}

// Project represents a project in the system
type Project struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	Name        string        `json:"name" db:"name" validate:"required,max=200"`
	Description *string       `json:"description,omitempty" db:"description"`
	Color       string        `json:"color" db:"color" validate:"required,hexcolor"`
	Status      ProjectStatus `json:"status" db:"status" validate:"required"`
	StartDate   *time.Time    `json:"start_date,omitempty" db:"start_date"`
	EndDate     *time.Time    `json:"end_date,omitempty" db:"end_date"`
	CreatedBy   uuid.UUID     `json:"created_by" db:"created_by"`
	IsDefault   bool          `json:"is_default" db:"is_default"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
}

// ProjectRequest represents the data required to create or update a project
type ProjectRequest struct {
	Name        string        `json:"name" validate:"required,max=200"`
	Description *string       `json:"description,omitempty" validate:"omitempty,max=1000"`
	Color       string        `json:"color" validate:"required,hexcolor"`
	Status      ProjectStatus `json:"status" validate:"required"`
	StartDate   *time.Time    `json:"start_date,omitempty"`
	EndDate     *time.Time    `json:"end_date,omitempty"`
	IsDefault   bool          `json:"is_default"`
}

// Validate validates the project data
func (p *Project) Validate() error {
	if !p.Status.IsValid() {
		return errors.New("invalid project status")
	}
	return nil
}
