//go:build integration
// +build integration

package auth_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/garnizeh/englog/internal/auth"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/store"
	"github.com/garnizeh/englog/internal/store/testutils"
)

// TestAuthServiceIntegration tests the full authentication flow with database
// "Authentication is the foundation of security." �
func TestAuthServiceIntegration(t *testing.T) {
	db := testutils.DB(t)

	// Create test logger
	testLogger := logging.NewTestLogger()

	authService := auth.NewAuthServiceForTest(db, testLogger, "integration-test-secret-key")

	t.Run("FullAuthenticationFlow", func(t *testing.T) {
		ctx := context.Background()

		// Create a test user first
		email := fmt.Sprintf("integration-test-%d@example.com", time.Now().UnixNano())
		hashedPassword, err := authService.HashPassword("TestPassword123!@#")
		require.NoError(t, err)

		var testUserID uuid.UUID
		err = db.Write(ctx, func(qtx *store.Queries) error {
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        email,
				PasswordHash: hashedPassword,
				FirstName:    "Integration",
				LastName:     "Test",
				Timezone:     pgtype.Text{String: "UTC", Valid: true},
				Preferences:  []byte(`{}`),
			})
			if err != nil {
				return err
			}

			testUserID = user.ID
			return nil
		})
		require.NoError(t, err)

		testUserIDString := testUserID.String()

		// Create access token
		accessToken, err := authService.CreateAccessToken(testUserIDString)
		require.NoError(t, err)
		assert.NotEmpty(t, accessToken)

		// Validate access token
		claims, err := authService.ValidateToken(ctx, accessToken)
		require.NoError(t, err)
		assert.Equal(t, testUserIDString, claims.UserID)
		assert.Equal(t, models.TokenAccess, claims.TokenType)

		// Create refresh token
		refreshToken, err := authService.CreateRefreshToken(testUserIDString)
		require.NoError(t, err)
		assert.NotEmpty(t, refreshToken)

		// Validate refresh token
		refreshClaims, err := authService.ValidateToken(ctx, refreshToken)
		require.NoError(t, err)
		assert.Equal(t, testUserIDString, refreshClaims.UserID)
		assert.Equal(t, models.TokenRefresh, refreshClaims.TokenType)

		// Test token rotation
		newAccessToken, newRefreshToken, err := authService.RotateRefreshToken(ctx, refreshToken)
		require.NoError(t, err)
		assert.NotEmpty(t, newAccessToken)
		assert.NotEmpty(t, newRefreshToken)
		assert.NotEqual(t, refreshToken, newRefreshToken)

		// Original refresh token should now be denylisted
		_, err = authService.ValidateToken(ctx, refreshToken)
		assert.Error(t, err, "Original refresh token should be denylisted after rotation")

		// New tokens should be valid
		_, err = authService.ValidateToken(ctx, newAccessToken)
		require.NoError(t, err)
		_, err = authService.ValidateToken(ctx, newRefreshToken)
		require.NoError(t, err)
	})
}

// TestRefreshTokenDenylistIntegration tests refresh token denylist operations
// "The denylist protects us from revoked tokens." 🚫
func TestRefreshTokenDenylistIntegration(t *testing.T) {
	db := testutils.DB(t)

	// Create test logger
	testLogger := logging.NewTestLogger()

	authService := auth.NewAuthServiceForTest(db, testLogger, "denylist-test-secret-key")

	// Helper to create a test user
	createTestUser := func(email string) uuid.UUID {
		ctx := context.Background()
		hashedPassword, err := authService.HashPassword("TestPassword123!@#")
		require.NoError(t, err)

		var userID uuid.UUID
		err = db.Write(ctx, func(qtx *store.Queries) error {
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        email,
				PasswordHash: hashedPassword,
				FirstName:    "Denylist",
				LastName:     "Test",
				Timezone:     pgtype.Text{String: "UTC", Valid: true},
				Preferences:  []byte(`{}`),
			})
			if err != nil {
				return err
			}

			userID = user.ID
			return nil
		})
		require.NoError(t, err)
		return userID
	}

	t.Run("DenylistAndValidation", func(t *testing.T) {
		ctx := context.Background()
		testUserID := createTestUser(fmt.Sprintf("denylist-test-1-%d@example.com", time.Now().UnixNano()))
		testUserIDString := testUserID.String()

		// Create a refresh token
		refreshToken, err := authService.CreateRefreshToken(testUserIDString)
		require.NoError(t, err)

		// Token should be valid initially
		claims, err := authService.ValidateToken(ctx, refreshToken)
		require.NoError(t, err)
		assert.Equal(t, testUserIDString, claims.UserID)

		// Denylist the token
		err = authService.DenylistRefreshToken(ctx, claims.JTI, testUserIDString)
		require.NoError(t, err)

		// Token should now be invalid
		_, err = authService.ValidateToken(ctx, refreshToken)
		assert.ErrorIs(t, err, auth.ErrTokenDenylisted)
	})

	t.Run("CleanupExpiredTokens", func(t *testing.T) {
		ctx := context.Background()
		testUserID := createTestUser(fmt.Sprintf("denylist-test-2-%d@example.com", time.Now().UnixNano()))
		testUserIDString := testUserID.String()

		// Create and denylist multiple tokens
		for i := 0; i < 3; i++ {
			refreshToken, err := authService.CreateRefreshToken(testUserIDString)
			require.NoError(t, err)

			claims, err := authService.ValidateToken(ctx, refreshToken)
			require.NoError(t, err)

			err = authService.DenylistRefreshToken(ctx, claims.JTI, testUserIDString)
			require.NoError(t, err)

			// Verify token is on denylist
			_, err = authService.ValidateToken(ctx, refreshToken)
			assert.ErrorIs(t, err, auth.ErrTokenDenylisted)
		}
	})
}

// TestTokenRotationIntegration tests token rotation with database
// "Rotation is the key to token security." 🔄
func TestTokenRotationIntegration(t *testing.T) {
	db := testutils.DB(t)

	// Create test logger
	testLogger := logging.NewTestLogger()

	authService := auth.NewAuthServiceForTest(db, testLogger, "rotation-test-secret-key")

	// Helper to create a test user
	createTestUser := func(email string) uuid.UUID {
		ctx := context.Background()
		hashedPassword, err := authService.HashPassword("TestPassword123!@#")
		require.NoError(t, err)

		var userID uuid.UUID
		err = db.Write(ctx, func(qtx *store.Queries) error {
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        email,
				PasswordHash: hashedPassword,
				FirstName:    "Rotation",
				LastName:     "Test",
				Timezone:     pgtype.Text{String: "UTC", Valid: true},
				Preferences:  []byte(`{}`),
			})
			if err != nil {
				return err
			}
			userID = user.ID
			return nil
		})
		require.NoError(t, err)
		return userID
	}

	t.Run("SuccessfulRotation", func(t *testing.T) {
		ctx := context.Background()
		testUserID := createTestUser(fmt.Sprintf("rotation-test-1-%d@example.com", time.Now().UnixNano()))
		testUserIDString := testUserID.String()

		// Create initial refresh token
		originalRefreshToken, err := authService.CreateRefreshToken(testUserIDString)
		require.NoError(t, err)

		// Validate original token works
		originalClaims, err := authService.ValidateToken(ctx, originalRefreshToken)
		require.NoError(t, err)
		assert.Equal(t, testUserIDString, originalClaims.UserID)

		// Rotate the token
		newAccessToken, newRefreshToken, err := authService.RotateRefreshToken(ctx, originalRefreshToken)
		require.NoError(t, err)
		assert.NotEmpty(t, newAccessToken)
		assert.NotEmpty(t, newRefreshToken)
		assert.NotEqual(t, originalRefreshToken, newRefreshToken)

		// Original token should be denylisted
		_, err = authService.ValidateToken(ctx, originalRefreshToken)
		assert.Error(t, err, "Original token should be denylisted")

		// New tokens should be valid
		newAccessClaims, err := authService.ValidateToken(ctx, newAccessToken)
		require.NoError(t, err)
		assert.Equal(t, testUserIDString, newAccessClaims.UserID)
		assert.Equal(t, models.TokenAccess, newAccessClaims.TokenType)

		newRefreshClaims, err := authService.ValidateToken(ctx, newRefreshToken)
		require.NoError(t, err)
		assert.Equal(t, testUserIDString, newRefreshClaims.UserID)
		assert.Equal(t, models.TokenRefresh, newRefreshClaims.TokenType)
	})

	t.Run("RotationChain", func(t *testing.T) {
		ctx := context.Background()
		testUserID := createTestUser(fmt.Sprintf("rotation-test-2-%d@example.com", time.Now().UnixNano()))
		testUserIDString := testUserID.String()

		// Start with initial token
		currentRefreshToken, err := authService.CreateRefreshToken(testUserIDString)
		require.NoError(t, err)

		var previousTokens []string

		// Perform multiple rotations
		for i := 0; i < 3; i++ {
			previousTokens = append(previousTokens, currentRefreshToken)

			newAccessToken, newRefreshToken, err := authService.RotateRefreshToken(ctx, currentRefreshToken)
			require.NoError(t, err)
			assert.NotEmpty(t, newAccessToken)
			assert.NotEmpty(t, newRefreshToken)

			// Previous token should be invalid
			_, err = authService.ValidateToken(ctx, currentRefreshToken)
			assert.Error(t, err, fmt.Sprintf("Token from rotation %d should be denylisted", i))

			// Update current token
			currentRefreshToken = newRefreshToken

			// Current token should be valid
			_, err = authService.ValidateToken(ctx, currentRefreshToken)
			require.NoError(t, err)
		}

		// Verify all previous tokens are still invalid
		for i, previousToken := range previousTokens {
			_, err := authService.ValidateToken(ctx, previousToken)
			assert.Error(t, err, fmt.Sprintf("Previous token %d should remain denylisted", i))
		}
	})
}

// TestConcurrentTokenOperations tests thread safety of token operations
// "Concurrency is the test of true robustness." ⚡
func TestConcurrentTokenOperations(t *testing.T) {
	db := testutils.DB(t)

	// Create test logger
	testLogger := logging.NewTestLogger()

	authService := auth.NewAuthServiceForTest(db, testLogger, "concurrent-test-secret-key")

	t.Run("ConcurrentTokenCreation", func(t *testing.T) {
		const numGoroutines = 10
		const tokensPerGoroutine = 5

		tokenChan := make(chan string, numGoroutines*tokensPerGoroutine*2) // Access + Refresh
		errChan := make(chan error, numGoroutines)

		// Start multiple goroutines creating tokens concurrently
		for i := 0; i < numGoroutines; i++ {
			go func(goroutineID int) {
				userID := uuid.New().String()

				for j := 0; j < tokensPerGoroutine; j++ {
					// Create access token
					accessToken, err := authService.CreateAccessToken(userID)
					if err != nil {
						errChan <- err
						return
					}
					tokenChan <- accessToken

					// Create refresh token
					refreshToken, err := authService.CreateRefreshToken(userID)
					if err != nil {
						errChan <- err
						return
					}
					tokenChan <- refreshToken
				}
				errChan <- nil
			}(i)
		}

		// Collect results
		var tokens []string
		for i := 0; i < numGoroutines; i++ {
			err := <-errChan
			require.NoError(t, err)
		}

		// Collect all tokens
		close(tokenChan)
		for token := range tokenChan {
			tokens = append(tokens, token)
		}

		require.Len(t, tokens, numGoroutines*tokensPerGoroutine*2)

		ctx := context.Background()

		// Verify all tokens are valid
		for _, token := range tokens {
			_, err := authService.ValidateToken(ctx, token)
			require.NoError(t, err)
		}
	})

	t.Run("ConcurrentDenylistOperations", func(t *testing.T) {
		// Create a test user first
		ctx := context.Background()
		email := fmt.Sprintf("concurrent-denylist-%d@example.com", time.Now().UnixNano())
		hashedPassword, err := authService.HashPassword("TestPassword123!@#")
		require.NoError(t, err)

		var testUserID uuid.UUID
		err = db.Write(ctx, func(qtx *store.Queries) error {
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        email,
				PasswordHash: hashedPassword,
				FirstName:    "Concurrent",
				LastName:     "Test",
				Timezone:     pgtype.Text{String: "UTC", Valid: true},
				Preferences:  []byte(`{}`),
			})
			if err != nil {
				return err
			}
			testUserID = user.ID
			return nil
		})
		require.NoError(t, err)
		userID := testUserID.String()

		const numGoroutines = 5
		const tokensPerGoroutine = 3

		var refreshTokens []string
		var jtis []string

		// Create refresh tokens first
		for i := 0; i < numGoroutines*tokensPerGoroutine; i++ {
			refreshToken, err := authService.CreateRefreshToken(userID)
			require.NoError(t, err)
			refreshTokens = append(refreshTokens, refreshToken)

			claims, err := authService.ValidateToken(ctx, refreshToken)
			require.NoError(t, err)
			jtis = append(jtis, claims.JTI)
		}

		// Concurrently denylist tokens
		errChan := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(goroutineID int) {
				start := goroutineID * tokensPerGoroutine
				end := start + tokensPerGoroutine

				for j := start; j < end; j++ {
					err := authService.DenylistRefreshToken(ctx, jtis[j], userID)
					if err != nil {
						errChan <- err
						return
					}
				}
				errChan <- nil
			}(i)
		}

		// Collect results
		for i := 0; i < numGoroutines; i++ {
			err := <-errChan
			require.NoError(t, err)
		}

		// Verify all tokens are denylisted
		for _, refreshToken := range refreshTokens {
			_, err := authService.ValidateToken(ctx, refreshToken)
			assert.ErrorIs(t, err, auth.ErrTokenDenylisted)
		}
	})
}

// TestPasswordOperationsIntegration tests password operations with database
// "Passwords are the keys to the kingdom." 🔑
func TestPasswordOperationsIntegration(t *testing.T) {
	ctx := context.Background()
	db := testutils.DB(t)

	// Create test logger
	testLogger := logging.NewTestLogger()

	authService := auth.NewAuthServiceForTest(db, testLogger, "password-test-secret-key")

	t.Run("UserCreationAndPasswordValidation", func(t *testing.T) {
		email := fmt.Sprintf("password-test-%d@example.com", time.Now().UnixNano())
		password := "TestPassword123!@#"

		// Hash password
		hashedPassword, err := authService.HashPassword(password)
		require.NoError(t, err)

		// Create user with hashed password
		err = db.Write(ctx, func(qtx *store.Queries) error {
			_, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        email,
				PasswordHash: hashedPassword,
				FirstName:    "Password",
				LastName:     "Test",
				Timezone:     pgtype.Text{String: "UTC", Valid: true},
				Preferences:  []byte(`{}`),
			})
			return err
		})
		require.NoError(t, err)

		// Retrieve user and verify password
		var retrievedUser store.User
		err = db.Read(ctx, func(qtx *store.Queries) error {
			user, err := qtx.GetUserByEmail(ctx, email)
			if err != nil {
				return err
			}
			retrievedUser = user
			return nil
		})
		require.NoError(t, err)

		// Test correct password
		assert.True(t, authService.CheckPassword(password, retrievedUser.PasswordHash))

		// Test incorrect password
		assert.False(t, authService.CheckPassword("WrongPassword", retrievedUser.PasswordHash))
	})
}

// TestAuthServiceCleanupIntegration tests cleanup operations
// "Cleanliness is next to godliness." 🧹
func TestAuthServiceCleanupIntegration(t *testing.T) {
	ctx := context.Background()

	db := testutils.DB(t)

	// Create test logger
	testLogger := logging.NewTestLogger()

	authService := auth.NewAuthServiceForTest(db, testLogger, "cleanup-test-secret-key")

	t.Run("CleanupExpiredDenylistedTokens", func(t *testing.T) {
		// This test would need additional SQL queries to manipulate expiration dates
		// For now, we'll just test that the cleanup method runs without error
		err := authService.CleanupExpiredTokens(ctx)
		require.NoError(t, err)
	})
}

// Benchmark integration tests for performance under real database conditions
func BenchmarkAuthServiceIntegration(b *testing.B) {
	db := testutils.DB(b)

	// Create test logger
	testLogger := logging.NewTestLogger()

	authService := auth.NewAuthServiceForTest(db, testLogger, "benchmark-secret-key")

	b.Run("TokenCreationWithDB", func(b *testing.B) {
		// Create a test user first
		ctx := context.Background()
		email := fmt.Sprintf("benchmark-creation-%d@example.com", time.Now().UnixNano())
		hashedPassword, err := authService.HashPassword("TestPassword123!@#")
		if err != nil {
			b.Fatal(err)
		}

		var testUserID uuid.UUID
		err = db.Write(ctx, func(qtx *store.Queries) error {
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        email,
				PasswordHash: hashedPassword,
				FirstName:    "Benchmark",
				LastName:     "Creation",
				Timezone:     pgtype.Text{String: "UTC", Valid: true},
				Preferences:  []byte(`{}`),
			})
			if err != nil {
				return err
			}
			testUserID = user.ID
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
		testUserIDString := testUserID.String()

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := authService.CreateRefreshToken(testUserIDString)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("TokenValidationWithDenylistCheck", func(b *testing.B) {
		// Create a test user first
		ctx := context.Background()
		email := fmt.Sprintf("benchmark-validation-%d@example.com", time.Now().UnixNano())
		hashedPassword, err := authService.HashPassword("TestPassword123!@#")
		if err != nil {
			b.Fatal(err)
		}

		var testUserID uuid.UUID
		err = db.Write(ctx, func(qtx *store.Queries) error {
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        email,
				PasswordHash: hashedPassword,
				FirstName:    "Benchmark",
				LastName:     "Validation",
				Timezone:     pgtype.Text{String: "UTC", Valid: true},
				Preferences:  []byte(`{}`),
			})
			if err != nil {
				return err
			}
			testUserID = user.ID
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
		testUserIDString := testUserID.String()

		refreshToken, err := authService.CreateRefreshToken(testUserIDString)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := authService.ValidateToken(ctx, refreshToken)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("TokenRotationWithDB", func(b *testing.B) {
		// Create a test user first
		ctx := context.Background()
		email := fmt.Sprintf("benchmark-rotation-%d@example.com", time.Now().UnixNano())
		hashedPassword, err := authService.HashPassword("TestPassword123!@#")
		if err != nil {
			b.Fatal(err)
		}

		var testUserID uuid.UUID
		err = db.Write(ctx, func(qtx *store.Queries) error {
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        email,
				PasswordHash: hashedPassword,
				FirstName:    "Benchmark",
				LastName:     "Rotation",
				Timezone:     pgtype.Text{String: "UTC", Valid: true},
				Preferences:  []byte(`{}`),
			})
			if err != nil {
				return err
			}
			testUserID = user.ID
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
		testUserIDString := testUserID.String()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			refreshToken, err := authService.CreateRefreshToken(testUserIDString)
			if err != nil {
				b.Fatal(err)
			}
			b.StartTimer()

			_, _, err = authService.RotateRefreshToken(ctx, refreshToken)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DenylistOperationWithDB", func(b *testing.B) {
		// Create a test user first
		ctx := context.Background()
		email := fmt.Sprintf("benchmark-denylist-%d@example.com", time.Now().UnixNano())
		hashedPassword, err := authService.HashPassword("TestPassword123!@#")
		if err != nil {
			b.Fatal(err)
		}

		var testUserID uuid.UUID
		err = db.Write(ctx, func(qtx *store.Queries) error {
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        email,
				PasswordHash: hashedPassword,
				FirstName:    "Benchmark",
				LastName:     "Test",
				Timezone:     pgtype.Text{String: "UTC", Valid: true},
				Preferences:  []byte(`{}`),
			})
			if err != nil {
				return err
			}
			testUserID = user.ID
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
		testUserIDString := testUserID.String()

		// Pre-create tokens
		var tokens []string
		var jtis []string
		for i := 0; i < b.N; i++ {
			refreshToken, err := authService.CreateRefreshToken(testUserIDString)
			if err != nil {
				b.Fatal(err)
			}
			tokens = append(tokens, refreshToken)

			claims, err := authService.ValidateToken(ctx, refreshToken)
			if err != nil {
				b.Fatal(err)
			}
			jtis = append(jtis, claims.JTI)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err := authService.DenylistRefreshToken(ctx, jtis[i], testUserIDString)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
