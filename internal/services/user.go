package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/garnizeh/englog/internal/database"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles all business logic for user management
type UserService struct {
	db     *database.DB
	logger *logging.Logger
}

// NewUserService creates a new UserService instance
func NewUserService(db *database.DB, logger *logging.Logger) *UserService {
	return &UserService{
		db:     db,
		logger: logger.WithComponent("user_service"),
	}
}

// GetUserProfile retrieves a user's profile information
func (s *UserService) GetUserProfile(ctx context.Context, userID string) (*models.UserProfile, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in GetUserProfile", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Getting user profile", "user_id", userID)

	var profile *models.UserProfile

	// Use read transaction to get user profile
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcUser, err := qtx.GetUserByID(ctx, userUUID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		// Convert to UserProfile model
		profile = s.sqlcToUserProfile(sqlcUser)
		return nil
	}); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			s.logger.Warn("User not found", "user_id", userID)
		} else {
			s.logger.LogError(ctx, err, "Failed to get user profile", "user_id", userID)
		}
		return nil, err
	}

	s.logger.Info("User profile retrieved successfully", "user_id", userID, "email", profile.Email)
	return profile, nil
}

// UpdateUserProfile updates a user's profile information
func (s *UserService) UpdateUserProfile(ctx context.Context, userID string, req *models.UserProfileRequest) (*models.UserProfile, error) {
	if err := s.validateProfileRequest(req); err != nil {
		s.logger.Warn("Invalid user profile update request", "user_id", userID, "error", err.Error())
		return nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Updating user profile", "user_id", userID, "first_name", req.FirstName, "last_name", req.LastName, "timezone", req.Timezone)

	var profile *models.UserProfile

	// Start write transaction to update user profile
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		// Convert preferences to JSON bytes
		preferencesBytes, err := s.preferencesToBytes(req.Preferences)
		if err != nil {
			return fmt.Errorf("failed to marshal preferences: %w", err)
		}

		sqlcUser, err := qtx.UpdateUserProfile(ctx, store.UpdateUserProfileParams{
			ID:          userUUID,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Timezone:    stringToPgTextRequired(req.Timezone),
			Preferences: preferencesBytes,
		})
		if err != nil {
			return fmt.Errorf("failed to update user profile: %w", err)
		}

		// Convert to UserProfile model
		profile = s.sqlcToUserProfile(sqlcUser)
		return nil
	}); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			s.logger.Warn("User not found for profile update", "user_id", userID)
		} else {
			s.logger.LogError(ctx, err, "Failed to update user profile", "user_id", userID)
		}
		return nil, err
	}

	s.logger.Info("User profile updated successfully", "user_id", userID, "email", profile.Email)
	return profile, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, userID string, req *models.UserPasswordChangeRequest) error {
	if err := s.validatePasswordChangeRequest(req); err != nil {
		s.logger.Warn("Invalid password change request", "user_id", userID, "error", err.Error())
		return err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format", "user_id", userID)
		return fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Changing user password", "user_id", userID)

	// Start write transaction to change password
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		// First verify current password
		currentUser, err := qtx.GetUserByID(ctx, userUUID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		// Check current password
		if err := bcrypt.CompareHashAndPassword([]byte(currentUser.PasswordHash), []byte(req.CurrentPassword)); err != nil {
			return fmt.Errorf("current password is incorrect")
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash new password: %w", err)
		}

		// Update password
		err = qtx.UpdateUserPassword(ctx, store.UpdateUserPasswordParams{
			ID:           userUUID,
			PasswordHash: string(hashedPassword),
		})
		if err != nil {
			return fmt.Errorf("failed to update password: %w", err)
		}

		return nil
	}); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			s.logger.Warn("User not found for password change", "user_id", userID)
		} else if strings.Contains(err.Error(), "current password is incorrect") {
			s.logger.Warn("Incorrect current password provided", "user_id", userID)
		} else {
			s.logger.LogError(ctx, err, "Failed to change password", "user_id", userID)
		}
		return err
	}

	s.logger.Info("User password changed successfully", "user_id", userID)
	return nil
}

// UpdateLastLogin updates the user's last login timestamp
func (s *UserService) UpdateLastLogin(ctx context.Context, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format", "user_id", userID)
		return fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Updating last login timestamp", "user_id", userID)

	// Start write transaction to update last login
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		err := qtx.UpdateUserLastLogin(ctx, userUUID)
		if err != nil {
			return fmt.Errorf("failed to update last login: %w", err)
		}
		return nil
	}); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			s.logger.Warn("User not found for last login update", "user_id", userID)
		} else {
			s.logger.LogError(ctx, err, "Failed to update last login", "user_id", userID)
		}
		return err
	}

	s.logger.Info("Last login timestamp updated successfully", "user_id", userID)
	return nil
}

// GetUserByEmail retrieves a user by email (for authentication)
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	s.logger.Info("Getting user by email", "email", email)

	var user *models.User

	// Use read transaction to get user by email
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcUser, err := qtx.GetUserByEmail(ctx, email)
		if err != nil {
			return fmt.Errorf("failed to get user by email: %w", err)
		}

		// Convert to User model
		user = s.sqlcToUser(sqlcUser)
		return nil
	}); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			s.logger.Info("User not found by email", "email", email)
		} else {
			s.logger.LogError(ctx, err, "Failed to get user by email", "email", email)
		}
		return nil, err
	}

	s.logger.Info("User retrieved successfully by email", "email", email, "user_id", user.ID)
	return user, nil
}

// CreateUser creates a new user account
func (s *UserService) CreateUser(ctx context.Context, req *models.UserRegistration) (*models.User, error) {
	if err := s.validateRegistrationRequest(req); err != nil {
		s.logger.Warn("Invalid user registration request", "email", req.Email, "error", err.Error())
		return nil, err
	}

	s.logger.Info("Creating new user account", "email", req.Email, "first_name", req.FirstName, "last_name", req.LastName)

	var user *models.User

	// Start write transaction to create user
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// Default preferences
		defaultPreferences := make(map[string]any)
		preferencesBytes, err := s.preferencesToBytes(defaultPreferences)
		if err != nil {
			return fmt.Errorf("failed to marshal default preferences: %w", err)
		}

		sqlcUser, err := qtx.CreateUser(ctx, store.CreateUserParams{
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			Timezone:     stringToPgTextRequired(req.Timezone),
			Preferences:  preferencesBytes,
		})
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// Convert to User model
		user = s.sqlcToUser(sqlcUser)
		return nil
	}); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			s.logger.Warn("User registration failed - email already exists", "email", req.Email)
		} else {
			s.logger.LogError(ctx, err, "Failed to create user", "email", req.Email)
		}
		return nil, err
	}

	s.logger.Info("User account created successfully", "email", req.Email, "user_id", user.ID)
	return user, nil
}

// DeleteUser deletes a user account
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format", "user_id", userID)
		return fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Deleting user account", "user_id", userID)

	// First verify the user exists
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		_, err := qtx.GetUserByID(ctx, userUUID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}
		return nil
	}); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			s.logger.Warn("User not found for deletion", "user_id", userID)
		} else {
			s.logger.LogError(ctx, err, "Failed to verify user for deletion", "user_id", userID)
		}
		return fmt.Errorf("failed to verify user: %w", err)
	}

	// Now delete the user
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		err := qtx.DeleteUser(ctx, userUUID)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to delete user", "user_id", userID)
		return err
	}

	s.logger.Info("User account deleted successfully", "user_id", userID)
	return nil
}

// GetUserCount returns the total number of users in the system
func (s *UserService) GetUserCount(ctx context.Context) (int64, error) {
	s.logger.Info("Getting user count")

	var count int64

	// Use read transaction to get user count
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		var err error
		count, err = qtx.GetUserCount(ctx)
		if err != nil {
			return fmt.Errorf("failed to get user count: %w", err)
		}
		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get user count")
		return 0, err
	}

	s.logger.Info("User count retrieved successfully", "user_count", count)
	return count, nil
}

// GetRecentUsers returns recently registered users
func (s *UserService) GetRecentUsers(ctx context.Context, limit int32) ([]*models.UserProfile, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	s.logger.Info("Getting recent users", "limit", limit)

	var profiles []*models.UserProfile

	// Use read transaction to get recent users
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcUsers, err := qtx.GetRecentUsers(ctx, limit)
		if err != nil {
			return fmt.Errorf("failed to get recent users: %w", err)
		}

		profiles = make([]*models.UserProfile, len(sqlcUsers))
		for i, sqlcUser := range sqlcUsers {
			profiles[i] = &models.UserProfile{
				ID:        sqlcUser.ID,
				Email:     sqlcUser.Email,
				FirstName: sqlcUser.FirstName,
				LastName:  sqlcUser.LastName,
				CreatedAt: pgTimestamptzToTime(sqlcUser.CreatedAt),
			}
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get recent users", "limit", limit)
		return nil, err
	}

	s.logger.Info("Recent users retrieved successfully", "limit", limit, "returned_count", len(profiles))
	return profiles, nil
}

// validateProfileRequest validates a user profile update request
func (s *UserService) validateProfileRequest(req *models.UserProfileRequest) error {
	if req.FirstName == "" {
		return fmt.Errorf("first name is required")
	}

	// Check for whitespace-only first name
	if strings.TrimSpace(req.FirstName) == "" {
		return fmt.Errorf("first name cannot be whitespace only")
	}

	if req.LastName == "" {
		return fmt.Errorf("last name is required")
	}

	// Check for whitespace-only last name
	if strings.TrimSpace(req.LastName) == "" {
		return fmt.Errorf("last name cannot be whitespace only")
	}

	if req.Timezone == "" {
		return fmt.Errorf("timezone is required")
	}

	// Validate timezone format
	if err := models.ValidateTimezone(req.Timezone); err != nil {
		return fmt.Errorf("invalid timezone: %w", err)
	}

	return nil
}

// validatePasswordChangeRequest validates a password change request
func (s *UserService) validatePasswordChangeRequest(req *models.UserPasswordChangeRequest) error {
	if req.CurrentPassword == "" {
		return fmt.Errorf("current password is required")
	}

	if req.NewPassword == "" {
		return fmt.Errorf("new password is required")
	}

	if len(req.NewPassword) < 8 {
		return fmt.Errorf("new password must be at least 8 characters long")
	}

	if len(req.NewPassword) > 100 {
		return fmt.Errorf("new password must be at most 100 characters long")
	}

	return nil
}

// validateRegistrationRequest validates a user registration request
func (s *UserService) validateRegistrationRequest(req *models.UserRegistration) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}

	// Check for whitespace-only email
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email cannot be whitespace only")
	}

	if req.Password == "" {
		return fmt.Errorf("password is required")
	}

	if len(req.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if len(req.Password) > 100 {
		return fmt.Errorf("password must be at most 100 characters long")
	}

	if req.FirstName == "" {
		return fmt.Errorf("first name is required")
	}

	// Check for whitespace-only first name
	if strings.TrimSpace(req.FirstName) == "" {
		return fmt.Errorf("first name cannot be whitespace only")
	}

	if req.LastName == "" {
		return fmt.Errorf("last name is required")
	}

	// Check for whitespace-only last name
	if strings.TrimSpace(req.LastName) == "" {
		return fmt.Errorf("last name cannot be whitespace only")
	}

	if req.Timezone == "" {
		return fmt.Errorf("timezone is required")
	}

	// Validate timezone format
	if err := models.ValidateTimezone(req.Timezone); err != nil {
		return fmt.Errorf("invalid timezone: %w", err)
	}

	return nil
}

// sqlcToUser converts a store.User to models.User
func (s *UserService) sqlcToUser(sqlcUser store.User) *models.User {
	preferences, _ := s.bytesToPreferences(sqlcUser.Preferences)

	return &models.User{
		ID:           sqlcUser.ID,
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		FirstName:    sqlcUser.FirstName,
		LastName:     sqlcUser.LastName,
		Timezone:     pgTextToStringRequired(sqlcUser.Timezone),
		Preferences:  preferences,
		LastLoginAt:  s.pgTimestamptzToTimePtr(sqlcUser.LastLoginAt),
		CreatedAt:    pgTimestamptzToTime(sqlcUser.CreatedAt),
		UpdatedAt:    pgTimestamptzToTime(sqlcUser.UpdatedAt),
	}
}

// sqlcToUserProfile converts a store.User to models.UserProfile
func (s *UserService) sqlcToUserProfile(sqlcUser store.User) *models.UserProfile {
	preferences, _ := s.bytesToPreferences(sqlcUser.Preferences)

	return &models.UserProfile{
		ID:          sqlcUser.ID,
		Email:       sqlcUser.Email,
		FirstName:   sqlcUser.FirstName,
		LastName:    sqlcUser.LastName,
		Timezone:    pgTextToStringRequired(sqlcUser.Timezone),
		Preferences: preferences,
		LastLoginAt: s.pgTimestamptzToTimePtr(sqlcUser.LastLoginAt),
		CreatedAt:   pgTimestamptzToTime(sqlcUser.CreatedAt),
		UpdatedAt:   pgTimestamptzToTime(sqlcUser.UpdatedAt),
	}
}

// preferencesToBytes converts map[string]any to JSON bytes
func (s *UserService) preferencesToBytes(preferences map[string]any) ([]byte, error) {
	if preferences == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(preferences)
}

// bytesToPreferences converts JSON bytes to map[string]any
func (s *UserService) bytesToPreferences(data []byte) (map[string]any, error) {
	if len(data) == 0 {
		return make(map[string]any), nil
	}

	var preferences map[string]any
	err := json.Unmarshal(data, &preferences)
	if err != nil {
		return make(map[string]any), err
	}

	return preferences, nil
}

// pgTimestamptzToTimePtr converts pgtype.Timestamptz to *time.Time
func (s *UserService) pgTimestamptzToTimePtr(pgTime pgtype.Timestamptz) *time.Time {
	if !pgTime.Valid {
		return nil
	}
	return &pgTime.Time
}
