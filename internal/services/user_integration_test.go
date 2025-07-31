//go:build integration
// +build integration

package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/services"
	"github.com/garnizeh/englog/internal/store/testutils"
)

// TestUserServiceIntegration tests the full user management flow with database
func TestUserServiceIntegration(t *testing.T) {
	db := testutils.DB(t)

	// Create test logger
	testLogger := logging.NewTestLogger()

	userService := services.NewUserService(db, testLogger)

	t.Run("FullUserLifecycle", func(t *testing.T) {
		ctx := context.Background()

		// Create a new user
		regReq := &models.UserRegistration{
			Email:     fmt.Sprintf("integration-user-%d@example.com", time.Now().UnixNano()),
			Password:  "securepassword123",
			FirstName: "Integration",
			LastName:  "User",
			Timezone:  "UTC",
		}

		// Test user creation
		createdUser, err := userService.CreateUser(ctx, regReq)
		require.NoError(t, err)
		require.NotNil(t, createdUser)
		assert.Equal(t, regReq.Email, createdUser.Email)
		assert.Equal(t, regReq.FirstName, createdUser.FirstName)
		assert.Equal(t, regReq.LastName, createdUser.LastName)
		assert.Equal(t, regReq.Timezone, createdUser.Timezone)
		assert.NotEqual(t, uuid.Nil, createdUser.ID)

		// Test user profile retrieval
		profile, err := userService.GetUserProfile(ctx, createdUser.ID.String())
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, profile.ID)
		assert.Equal(t, createdUser.Email, profile.Email)
		assert.Equal(t, createdUser.FirstName, profile.FirstName)

		// Test user profile update
		updateReq := &models.UserProfileRequest{
			FirstName: "Updated",
			LastName:  "Name",
			Timezone:  "America/New_York",
			Preferences: map[string]any{
				"theme":    "dark",
				"language": "en",
			},
		}

		updatedProfile, err := userService.UpdateUserProfile(ctx, createdUser.ID.String(), updateReq)
		require.NoError(t, err)
		assert.Equal(t, updateReq.FirstName, updatedProfile.FirstName)
		assert.Equal(t, updateReq.LastName, updatedProfile.LastName)
		assert.Equal(t, updateReq.Timezone, updatedProfile.Timezone)
		assert.Equal(t, "dark", updatedProfile.Preferences["theme"])

		// Test password change
		passwordChangeReq := &models.UserPasswordChangeRequest{
			CurrentPassword: "securepassword123",
			NewPassword:     "newsecurepassword456",
		}

		err = userService.ChangePassword(ctx, createdUser.ID.String(), passwordChangeReq)
		require.NoError(t, err)

		// Test GetUserByEmail
		userByEmail, err := userService.GetUserByEmail(ctx, regReq.Email)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, userByEmail.ID)

		// Test UpdateLastLogin
		err = userService.UpdateLastLogin(ctx, createdUser.ID.String())
		require.NoError(t, err)

		// Test user deletion
		err = userService.DeleteUser(ctx, createdUser.ID.String())
		require.NoError(t, err)

		// Verify user is deleted
		_, err = userService.GetUserProfile(ctx, createdUser.ID.String())
		assert.Error(t, err)
	})

	t.Run("UserPreferencesManagement", func(t *testing.T) {
		ctx := context.Background()

		// Create test user
		regReq := &models.UserRegistration{
			Email:     fmt.Sprintf("pref-user-%d@example.com", time.Now().UnixNano()),
			Password:  "testpassword123",
			FirstName: "Preferences",
			LastName:  "User",
			Timezone:  "UTC",
		}

		user, err := userService.CreateUser(ctx, regReq)
		require.NoError(t, err)

		defer func() {
			_ = userService.DeleteUser(ctx, user.ID.String())
		}()

		// Test complex preferences update
		complexPreferences := map[string]any{
			"theme":         "dark",
			"language":      "en",
			"notifications": true,
			"timezone":      "America/New_York",
			"dashboard": map[string]any{
				"layout":      "grid",
				"widgets":     []string{"calendar", "tasks", "analytics"},
				"autoRefresh": true,
			},
			"privacy": map[string]any{
				"showEmail":   false,
				"showProfile": true,
			},
		}

		updateReq := &models.UserProfileRequest{
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Timezone:    user.Timezone,
			Preferences: complexPreferences,
		}

		updatedProfile, err := userService.UpdateUserProfile(ctx, user.ID.String(), updateReq)
		require.NoError(t, err)

		// Verify complex preferences are stored correctly
		assert.Equal(t, "dark", updatedProfile.Preferences["theme"])
		assert.Equal(t, true, updatedProfile.Preferences["notifications"])

		dashboard, ok := updatedProfile.Preferences["dashboard"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "grid", dashboard["layout"])
	})

	t.Run("UserCountAndRecentUsers", func(t *testing.T) {
		ctx := context.Background()

		// Get initial count
		initialCount, err := userService.GetUserCount(ctx)
		require.NoError(t, err)

		// Create multiple users
		var createdUsers []*models.User
		for i := 0; i < 5; i++ {
			regReq := &models.UserRegistration{
				Email:     fmt.Sprintf("recent-user-%d-%d@example.com", i, time.Now().UnixNano()),
				Password:  "testpassword123",
				FirstName: fmt.Sprintf("Recent%d", i),
				LastName:  "User",
				Timezone:  "UTC",
			}

			user, err := userService.CreateUser(ctx, regReq)
			require.NoError(t, err)
			createdUsers = append(createdUsers, user)
		}

		// Clean up
		defer func() {
			for _, user := range createdUsers {
				_ = userService.DeleteUser(ctx, user.ID.String())
			}
		}()

		// Check updated count
		newCount, err := userService.GetUserCount(ctx)
		require.NoError(t, err)
		assert.Equal(t, initialCount+5, newCount)

		// Test GetRecentUsers
		recentUsers, err := userService.GetRecentUsers(ctx, 3)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(recentUsers), 3)

		// Verify the most recent users are returned (reverse chronological order)
		for i := 0; i < 3 && i < len(recentUsers); i++ {
			found := false
			for _, createdUser := range createdUsers {
				if recentUsers[i].ID == createdUser.ID {
					found = true
					break
				}
			}
			assert.True(t, found, "Recent user should be one of the created users")
		}
	})

	t.Run("ConcurrentUserOperations", func(t *testing.T) {
		ctx := context.Background()

		// Test concurrent user creation
		const numGoroutines = 5
		results := make(chan *models.User, numGoroutines)
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				regReq := &models.UserRegistration{
					Email:     fmt.Sprintf("concurrent-user-%d-%d@example.com", index, time.Now().UnixNano()),
					Password:  "testpassword123",
					FirstName: fmt.Sprintf("Concurrent%d", index),
					LastName:  "User",
					Timezone:  "UTC",
				}

				user, err := userService.CreateUser(ctx, regReq)
				if err != nil {
					errors <- err
					return
				}
				results <- user
			}(i)
		}

		// Collect results
		var createdUsers []*models.User
		for i := 0; i < numGoroutines; i++ {
			select {
			case user := <-results:
				createdUsers = append(createdUsers, user)
			case err := <-errors:
				t.Errorf("Concurrent operation failed: %v", err)
			case <-time.After(10 * time.Second):
				t.Error("Timeout waiting for concurrent operations")
			}
		}

		// Clean up created users
		defer func() {
			for _, user := range createdUsers {
				_ = userService.DeleteUser(ctx, user.ID.String())
			}
		}()

		assert.Equal(t, numGoroutines, len(createdUsers))

		// Verify all users are unique
		userEmails := make(map[string]bool)
		for _, user := range createdUsers {
			assert.False(t, userEmails[user.Email], "User email should be unique: %s", user.Email)
			userEmails[user.Email] = true
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		ctx := context.Background()

		// Test duplicate email registration
		regReq := &models.UserRegistration{
			Email:     fmt.Sprintf("duplicate-user-%d@example.com", time.Now().UnixNano()),
			Password:  "testpassword123",
			FirstName: "Original",
			LastName:  "User",
			Timezone:  "UTC",
		}

		user1, err := userService.CreateUser(ctx, regReq)
		require.NoError(t, err)

		defer func() {
			_ = userService.DeleteUser(ctx, user1.ID.String())
		}()

		// Try to create duplicate
		_, err = userService.CreateUser(ctx, regReq)
		assert.Error(t, err)

		// Test getting non-existent user
		_, err = userService.GetUserProfile(ctx, uuid.New().String())
		assert.Error(t, err)

		// Test updating non-existent user
		updateReq := &models.UserProfileRequest{
			FirstName: "Updated",
			LastName:  "Name",
			Timezone:  "UTC",
		}
		_, err = userService.UpdateUserProfile(ctx, uuid.New().String(), updateReq)
		assert.Error(t, err)

		// Test deleting non-existent user
		err = userService.DeleteUser(ctx, uuid.New().String())
		assert.Error(t, err)

		// Test invalid UUID formats
		invalidUUIDs := []string{"", "invalid", "not-a-uuid", "12345"}
		for _, invalidUUID := range invalidUUIDs {
			_, err = userService.GetUserProfile(ctx, invalidUUID)
			assert.Error(t, err)
		}

		// Test password change with wrong current password
		wrongPasswordReq := &models.UserPasswordChangeRequest{
			CurrentPassword: "wrongpassword",
			NewPassword:     "newsecurepassword456",
		}

		err = userService.ChangePassword(ctx, user1.ID.String(), wrongPasswordReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "current password is incorrect")
	})

	t.Run("EdgeCases", func(t *testing.T) {
		ctx := context.Background()

		// Test with minimal user data
		minimalReg := &models.UserRegistration{
			Email:     fmt.Sprintf("minimal-%d@example.com", time.Now().UnixNano()),
			Password:  "password123",
			FirstName: "A",
			LastName:  "B",
			Timezone:  "UTC",
		}

		user, err := userService.CreateUser(ctx, minimalReg)
		require.NoError(t, err)

		defer func() {
			_ = userService.DeleteUser(ctx, user.ID.String())
		}()

		// Test profile update with nil preferences
		updateReq := &models.UserProfileRequest{
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Timezone:    user.Timezone,
			Preferences: nil,
		}

		updatedProfile, err := userService.UpdateUserProfile(ctx, user.ID.String(), updateReq)
		require.NoError(t, err)
		assert.Empty(t, updatedProfile.Preferences)

		// Test GetUserByEmail with non-existent email
		_, err = userService.GetUserByEmail(ctx, "nonexistent@example.com")
		assert.Error(t, err)
	})
}

// BenchmarkUserService benchmarks user service operations
func BenchmarkUserService(b *testing.B) {
	db := testutils.DB(b)

	// Create test logger
	testLogger := logging.NewTestLogger()

	userService := services.NewUserService(db, testLogger)
	ctx := context.Background()

	b.Run("CreateUser", func(b *testing.B) {
		var createdIDs []string

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req := &models.UserRegistration{
				Email:     fmt.Sprintf("bench-user-%d@example.com", i),
				Password:  "password123",
				FirstName: "Bench",
				LastName:  "User",
				Timezone:  "UTC",
			}

			user, err := userService.CreateUser(ctx, req)
			if err != nil {
				b.Fatalf("Failed to create user: %v", err)
			}
			createdIDs = append(createdIDs, user.ID.String())
		}
		b.StopTimer()

		// Cleanup
		for _, id := range createdIDs {
			_ = userService.DeleteUser(ctx, id)
		}
	})

	b.Run("GetUserProfile", func(b *testing.B) {
		// Setup: create a user to fetch
		req := &models.UserRegistration{
			Email:     "bench-get-user@example.com",
			Password:  "password123",
			FirstName: "Bench",
			LastName:  "Get",
			Timezone:  "UTC",
		}

		user, err := userService.CreateUser(ctx, req)
		if err != nil {
			b.Fatalf("Failed to create test user: %v", err)
		}

		defer func() {
			_ = userService.DeleteUser(ctx, user.ID.String())
		}()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := userService.GetUserProfile(ctx, user.ID.String())
			if err != nil {
				b.Fatalf("Failed to get user profile: %v", err)
			}
		}
	})

	b.Run("UpdateUserProfile", func(b *testing.B) {
		// Setup: create a user to update
		regReq := &models.UserRegistration{
			Email:     "bench-update-user@example.com",
			Password:  "password123",
			FirstName: "Bench",
			LastName:  "Update",
			Timezone:  "UTC",
		}

		user, err := userService.CreateUser(ctx, regReq)
		if err != nil {
			b.Fatalf("Failed to create test user: %v", err)
		}

		defer func() {
			_ = userService.DeleteUser(ctx, user.ID.String())
		}()

		updateReq := &models.UserProfileRequest{
			FirstName: "Updated",
			LastName:  "Name",
			Timezone:  "America/New_York",
			Preferences: map[string]any{
				"theme": "dark",
			},
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := userService.UpdateUserProfile(ctx, user.ID.String(), updateReq)
			if err != nil {
				b.Fatalf("Failed to update user profile: %v", err)
			}
		}
	})
}
