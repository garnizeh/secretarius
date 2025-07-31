package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents the core user entity in the system
type User struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	Email        string         `json:"email" db:"email" validate:"required,email"`
	PasswordHash string         `json:"-" db:"password_hash"`
	FirstName    string         `json:"first_name" db:"first_name" validate:"required,max=100"`
	LastName     string         `json:"last_name" db:"last_name" validate:"required,max=100"`
	Timezone     string         `json:"timezone" db:"timezone" validate:"required"`
	Preferences  map[string]any `json:"preferences" db:"preferences"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
}

// UserRegistration represents the data required for user registration
// @Description User registration request payload
type UserRegistration struct {
	Email     string `json:"email" validate:"required,email,max=255" example:"user@example.com" description:"Valid email address"`
	Password  string `json:"password" validate:"required,min=8,max=100" example:"securePassword123" description:"Password (minimum 8 characters)"`
	FirstName string `json:"first_name" validate:"required,max=100" example:"John" description:"User's first name"`
	LastName  string `json:"last_name" validate:"required,max=100" example:"Doe" description:"User's last name"`
	Timezone  string `json:"timezone" validate:"required,timezone" example:"UTC" description:"User's timezone (IANA timezone)"`
}

// UserLogin represents the data required for user login
// @Description User login request payload
type UserLogin struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com" description:"User's email address"`
	Password string `json:"password" validate:"required" example:"securePassword123" description:"User's password"`
}

// UserProfile represents the public user profile information
// @Description User profile information (public data)
type UserProfile struct {
	ID          uuid.UUID      `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" description:"Unique user identifier"`
	Email       string         `json:"email" example:"user@example.com" description:"User's email address"`
	FirstName   string         `json:"first_name" example:"John" description:"User's first name"`
	LastName    string         `json:"last_name" example:"Doe" description:"User's last name"`
	Timezone    string         `json:"timezone" example:"UTC" description:"User's timezone"`
	Preferences map[string]any `json:"preferences" description:"User preferences (JSON object)"`
	LastLoginAt *time.Time     `json:"last_login_at,omitempty" example:"2024-01-15T14:30:00Z" description:"Last login timestamp"`
	CreatedAt   time.Time      `json:"created_at" example:"2024-01-01T00:00:00Z" description:"Account creation timestamp"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2024-01-15T12:00:00Z" description:"Last profile update timestamp"`
}

// UserProfileRequest represents the data required for updating user profile
type UserProfileRequest struct {
	FirstName   string         `json:"first_name" validate:"required,max=100"`
	LastName    string         `json:"last_name" validate:"required,max=100"`
	Timezone    string         `json:"timezone" validate:"required,timezone"`
	Preferences map[string]any `json:"preferences"`
}

// UserPasswordChangeRequest represents the data required for changing user password
type UserPasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=72"` // NewPassword does not accept passwords longer than 72 bytes
}
