package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents the core user entity in the system
type User struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	Email        string                 `json:"email" db:"email" validate:"required,email"`
	PasswordHash string                 `json:"-" db:"password_hash"`
	FirstName    string                 `json:"first_name" db:"first_name" validate:"required,max=100"`
	LastName     string                 `json:"last_name" db:"last_name" validate:"required,max=100"`
	Timezone     string                 `json:"timezone" db:"timezone" validate:"required"`
	Preferences  map[string]interface{} `json:"preferences" db:"preferences"`
	LastLoginAt  *time.Time             `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
}

// UserRegistration represents the data required for user registration
type UserRegistration struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"required,max=100"`
	Timezone  string `json:"timezone" validate:"required,timezone"`
}

// UserLogin represents the data required for user login
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserProfile represents the public user profile information
type UserProfile struct {
	ID          uuid.UUID              `json:"id"`
	Email       string                 `json:"email"`
	FirstName   string                 `json:"first_name"`
	LastName    string                 `json:"last_name"`
	Timezone    string                 `json:"timezone"`
	Preferences map[string]interface{} `json:"preferences"`
	LastLoginAt *time.Time             `json:"last_login_at,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}
