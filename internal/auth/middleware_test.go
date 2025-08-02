package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/garnizeh/englog/internal/auth"
	"github.com/garnizeh/englog/internal/logging"
)

// TestRequireAuthMiddleware tests the authentication middleware
// "Middleware is the guardian of the gate." 🚪
func TestRequireAuthMiddleware(t *testing.T) {
	// Create test logger
	testLogger := logging.NewTestLogger()

	gin.SetMode(gin.TestMode)
	authService := auth.NewAuthService(nil, testLogger, "test-secret-key")

	testCases := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectUserID   bool
	}{
		{
			name:           "valid bearer token",
			authHeader:     "Bearer " + createValidToken(t, authService),
			expectedStatus: http.StatusOK,
			expectUserID:   true,
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectUserID:   false,
		},
		{
			name:           "invalid bearer format",
			authHeader:     "InvalidFormat token",
			expectedStatus: http.StatusUnauthorized,
			expectUserID:   false,
		},
		{
			name:           "bearer without token",
			authHeader:     "Bearer",
			expectedStatus: http.StatusUnauthorized,
			expectUserID:   false,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid.jwt.token",
			expectedStatus: http.StatusUnauthorized,
			expectUserID:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			c.Request = req

			// Create test handler that checks if user_id is set
			handlerCalled := false
			testHandler := func(c *gin.Context) {
				handlerCalled = true
				if tc.expectUserID {
					userID, exists := c.Get("user_id")
					assert.True(t, exists, "user_id should be set in context")
					assert.NotEmpty(t, userID, "user_id should not be empty")
				}
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}

			// Apply middleware
			middleware := authService.RequireAuth()
			middleware(c)

			// Check if request was aborted
			if !c.IsAborted() {
				testHandler(c)
			}

			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectedStatus == http.StatusOK {
				assert.True(t, handlerCalled, "Handler should have been called")
			} else {
				assert.True(t, c.IsAborted(), "Request should have been aborted")
			}
		})
	}
}

// TestOptionalAuthMiddleware tests the optional authentication middleware
// "Optional security is still security." 🔓
func TestOptionalAuthMiddleware(t *testing.T) {
	// Create test logger
	testLogger := logging.NewTestLogger()

	gin.SetMode(gin.TestMode)
	authService := auth.NewAuthService(nil, testLogger, "test-secret-key")

	testCases := []struct {
		name         string
		authHeader   string
		expectUserID bool
	}{
		{
			name:         "valid bearer token",
			authHeader:   "Bearer " + createValidToken(t, authService),
			expectUserID: true,
		},
		{
			name:         "missing authorization header",
			authHeader:   "",
			expectUserID: false,
		},
		{
			name:         "invalid bearer format",
			authHeader:   "InvalidFormat token",
			expectUserID: false,
		},
		{
			name:         "invalid token",
			authHeader:   "Bearer invalid.jwt.token",
			expectUserID: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			req, _ := http.NewRequest(http.MethodGet, "/optional", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			c.Request = req

			// Create test handler
			testHandler := func(c *gin.Context) {
				userID, exists := c.Get("user_id")
				if tc.expectUserID {
					assert.True(t, exists, "user_id should be set in context")
					assert.NotEmpty(t, userID, "user_id should not be empty")
				} else {
					assert.False(t, exists, "user_id should not be set in context")
				}
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}

			// Apply middleware
			middleware := authService.OptionalAuth()
			middleware(c)

			// Handler should always be called for optional auth
			assert.False(t, c.IsAborted(), "Request should not be aborted for optional auth")
			testHandler(c)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

// TestRefreshTokenWithMiddleware tests refresh token validation through middleware
// "Refresh is the cycle of security." 🔄
// NOTE: This test is commented out because it requires database integration
// For full testing, use integration tests with a real database connection
/*
func TestRefreshTokenWithMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authService := auth.NewAuthService(nil, testLogger, "test-secret-key")

	userID := uuid.New().String()

	// Create a refresh token
	refreshToken, err := authService.CreateRefreshToken(userID)
	require.NoError(t, err)

	// Create test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request with refresh token (should be rejected for access-only endpoints)
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+refreshToken)
	c.Request = req

	// Apply middleware
	middleware := authService.RequireAuth()
	middleware(c)

	// Should be rejected because refresh tokens can't be used for access
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())

	// Check error message
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Unauthorized", response["error"])
}
*/

// TestCaseInsensitiveBearerToken tests case insensitive bearer token parsing
// "Case sensitivity is the enemy of usability." 📝
func TestCaseInsensitiveBearerToken(t *testing.T) {
	// Create test logger
	testLogger := logging.NewTestLogger()

	gin.SetMode(gin.TestMode)
	authService := auth.NewAuthService(nil, testLogger, "test-secret-key")

	token := createValidToken(t, authService)

	testCases := []struct {
		name       string
		authHeader string
		shouldWork bool
	}{
		{
			name:       "lowercase bearer",
			authHeader: "bearer " + token,
			shouldWork: true,
		},
		{
			name:       "uppercase bearer",
			authHeader: "BEARER " + token,
			shouldWork: true,
		},
		{
			name:       "mixed case bearer",
			authHeader: "Bearer " + token,
			shouldWork: true,
		},
		{
			name:       "weird case bearer",
			authHeader: "BeArEr " + token,
			shouldWork: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tc.authHeader)
			c.Request = req

			// Apply middleware
			middleware := authService.RequireAuth()
			middleware(c)

			if tc.shouldWork {
				assert.False(t, c.IsAborted(), "Request should not be aborted for valid bearer token")
				// If not aborted, verify user_id is set
				if !c.IsAborted() {
					userID, exists := c.Get("user_id")
					assert.True(t, exists, "user_id should be set")
					assert.NotEmpty(t, userID, "user_id should not be empty")
				}
			} else {
				assert.True(t, c.IsAborted(), "Request should be aborted for invalid format")
				assert.Equal(t, http.StatusUnauthorized, w.Code)
			}
		})
	}
}

// TestMultipleSpacesInAuthHeader tests handling of multiple spaces in auth header
// "Spaces should not break security." 🌌
func TestMultipleSpacesInAuthHeader(t *testing.T) {
	// Create test logger
	testLogger := logging.NewTestLogger()

	gin.SetMode(gin.TestMode)
	authService := auth.NewAuthService(nil, testLogger, "test-secret-key")

	token := createValidToken(t, authService)

	testCases := []struct {
		name       string
		authHeader string
		shouldWork bool
	}{
		{
			name:       "single space",
			authHeader: "Bearer " + token,
			shouldWork: true,
		},
		{
			name:       "multiple spaces",
			authHeader: "Bearer   " + token,
			shouldWork: false, // Current implementation splits on single space
		},
		{
			name:       "tab character",
			authHeader: "Bearer\t" + token,
			shouldWork: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tc.authHeader)
			c.Request = req

			// Apply middleware
			middleware := authService.RequireAuth()
			middleware(c)

			if tc.shouldWork {
				assert.False(t, c.IsAborted(), "Request should not be aborted")
			} else {
				assert.True(t, c.IsAborted(), "Request should be aborted")
				assert.Equal(t, http.StatusUnauthorized, w.Code)
			}
		})
	}
}

// Helper function to create a valid access token for testing
func createValidToken(t *testing.T, authService *auth.AuthService) string {
	userID := uuid.New().String()
	token, err := authService.CreateAccessToken(context.Background(), userID)
	require.NoError(t, err)
	return token
}

// Benchmark tests for middleware performance
// "Speed is security in action." ⚡

func BenchmarkRequireAuthMiddleware(b *testing.B) {
	// Create test logger
	testLogger := logging.NewTestLogger()

	gin.SetMode(gin.TestMode)
	authService := auth.NewAuthService(nil, testLogger, "test-secret-key-for-benchmarking")

	userID := uuid.New().String()
	token, err := authService.CreateAccessToken(context.Background(), userID)
	if err != nil {
		b.Fatal(err)
	}

	middleware := authService.RequireAuth()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		c.Request = req

		middleware(c)
	}
}

func BenchmarkOptionalAuthMiddleware(b *testing.B) {
	// Create test logger
	testLogger := logging.NewTestLogger()

	gin.SetMode(gin.TestMode)
	authService := auth.NewAuthService(nil, testLogger, "test-secret-key-for-benchmarking")

	userID := uuid.New().String()
	token, err := authService.CreateAccessToken(context.Background(), userID)
	if err != nil {
		b.Fatal(err)
	}

	middleware := authService.OptionalAuth()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/optional", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		c.Request = req

		middleware(c)
	}
}
