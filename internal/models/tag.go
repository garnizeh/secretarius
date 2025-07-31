package models

import (
	"time"

	"github.com/google/uuid"
)

// Tag represents a tag that can be applied to log entries
type Tag struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,max=100"`
	Color       string    `json:"color" db:"color" validate:"required,hexcolor"`
	Description *string   `json:"description,omitempty" db:"description"`
	UsageCount  int       `json:"usage_count" db:"usage_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// TagRequest represents the data required to create or update a tag
type TagRequest struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Color       string  `json:"color" validate:"required,hexcolor"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
}

// TagUsage represents tag usage statistics for a user
type TagUsage struct {
	Tag        *Tag `json:"tag"`
	UsageCount int  `json:"usage_count"`
}

// LogEntryTag represents the many-to-many relationship between log entries and tags
type LogEntryTag struct {
	LogEntryID uuid.UUID `json:"log_entry_id" db:"log_entry_id"`
	TagID      uuid.UUID `json:"tag_id" db:"tag_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
