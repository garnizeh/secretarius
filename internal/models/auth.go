package models

import (
	"time"

	"github.com/google/uuid"
)

// TokenType represents the type of authentication token
type TokenType string

const (
	TokenAccess  TokenType = "access"
	TokenRefresh TokenType = "refresh"
)

// AuthTokens represents the authentication tokens returned after login
// @Description Authentication tokens for API access
type AuthTokens struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." description:"JWT access token for API authentication"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." description:"JWT refresh token for token renewal"`
	ExpiresIn    int    `json:"expires_in" example:"3600" description:"Access token expiration time in seconds"`
	TokenType    string `json:"token_type" example:"Bearer" description:"Token type, always 'Bearer'"`
}

// RefreshTokenDenylist represents a denylisted refresh token
type RefreshTokenDenylist struct {
	JTI          string    `json:"jti" db:"jti"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	DenylistedAt time.Time `json:"denylisted_at" db:"denylisted_at"`
	Reason       string    `json:"reason" db:"reason"`
}

// UserSession represents an active user session
type UserSession struct {
	ID               uuid.UUID `json:"id" db:"id"`
	UserID           uuid.UUID `json:"user_id" db:"user_id"`
	SessionTokenHash string    `json:"-" db:"session_token_hash"` // Hidden from JSON for security
	RefreshTokenHash string    `json:"-" db:"refresh_token_hash"` // Hidden from JSON for security
	ExpiresAt        time.Time `json:"expires_at" db:"expires_at"`
	LastActivity     time.Time `json:"last_activity" db:"last_activity"`
	IPAddress        *string   `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent        *string   `json:"user_agent,omitempty" db:"user_agent"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// ScheduledDeletion represents a scheduled user data deletion request
type ScheduledDeletion struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	UserID       uuid.UUID      `json:"user_id" db:"user_id"`
	ScheduledAt  time.Time      `json:"scheduled_at" db:"scheduled_at"`
	DeletionType string         `json:"deletion_type" db:"deletion_type"`
	Status       string         `json:"status" db:"status"`
	CompletedAt  *time.Time     `json:"completed_at,omitempty" db:"completed_at"`
	Metadata     map[string]any `json:"metadata" db:"metadata"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
}
