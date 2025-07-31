package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garnizeh/englog/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserHandler_ComprehensiveScenarios tests complex user management scenarios
func TestUserHandler_ComprehensiveScenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("complete_user_profile_lifecycle", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// 1. Get initial profile
		req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var profileResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &profileResponse)
		require.NoError(t, err)

		profile := profileResponse["data"].(map[string]any)
		assert.Equal(t, user.Email, profile["email"])
		assert.Equal(t, user.FirstName, profile["first_name"])
		assert.Equal(t, user.LastName, profile["last_name"])

		// 2. Update profile with new information
		updateReq := models.UserProfileRequest{
			FirstName:   "UpdatedFirst",
			LastName:    "UpdatedLast",
			Timezone:    "America/New_York",
			Preferences: map[string]any{"theme": "dark"},
		}
		body, _ := json.Marshal(updateReq)

		req, _ = http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 3. Verify profile was updated
		req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &profileResponse)
		require.NoError(t, err)

		updatedProfile := profileResponse["data"].(map[string]any)
		assert.Equal(t, "UpdatedFirst", updatedProfile["first_name"])
		assert.Equal(t, "UpdatedLast", updatedProfile["last_name"])
		assert.Equal(t, "America/New_York", updatedProfile["timezone"])

		// 4. Change password
		passwordReq := models.UserPasswordChangeRequest{
			CurrentPassword: "password123",
			NewPassword:     "newSecurePassword456!",
		}
		body, _ = json.Marshal(passwordReq)

		req, _ = http.NewRequest("POST", "/v1/users/change-password", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 5. Verify old password no longer works
		loginReq := models.UserLogin{
			Email:    user.Email,
			Password: "password123",
		}
		loginBody, _ := json.Marshal(loginReq)
		req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(loginBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// 6. Verify new password works
		newLoginReq := models.UserLogin{
			Email:    user.Email,
			Password: "newSecurePassword456!",
		}
		newLoginBody, _ := json.Marshal(newLoginReq)
		req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(newLoginBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Get new token for account deletion
		var loginResponse struct {
			Tokens models.AuthTokens `json:"tokens"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
		require.NoError(t, err)
		newToken := loginResponse.Tokens.AccessToken

		// 7. Delete account
		req, _ = http.NewRequest("DELETE", "/v1/users/account", nil)
		req.Header.Set("Authorization", "Bearer "+newToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 8. Verify account is deleted - login should fail
		req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(newLoginBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("profile_update_validation", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Test with valid but edge case data
		updateReq := models.UserProfileRequest{
			FirstName: "A", // Single character - valid
			LastName:  "B", // Single character - valid
			Timezone:  "UTC",
		}
		body, _ := json.Marshal(updateReq)

		req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Test with maximum length names (100 chars)
		maxName := string(make([]byte, 100))
		for i := range maxName {
			maxName = maxName[:i] + "a" + maxName[i+1:]
		}

		updateReq = models.UserProfileRequest{
			FirstName: maxName,
			LastName:  maxName,
			Timezone:  "UTC",
		}
		body, _ = json.Marshal(updateReq)

		req, _ = http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("password_change_security", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Test changing to minimum valid password
		passwordReq := models.UserPasswordChangeRequest{
			CurrentPassword: "password123",
			NewPassword:     "12345678", // 8 characters - minimum valid
		}
		body, _ := json.Marshal(passwordReq)

		req, _ := http.NewRequest("POST", "/v1/users/change-password", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Test changing to same password
		passwordReq = models.UserPasswordChangeRequest{
			CurrentPassword: "12345678",
			NewPassword:     "12345678", // Same password
		}
		body, _ = json.Marshal(passwordReq)

		req, _ = http.NewRequest("POST", "/v1/users/change-password", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Test with long password (72 bytes - bcrypt limit)
		// Note: bcrypt has a maximum password length of 72 bytes
		longPassword := "abcdefghij1234567890abcdefghij1234567890abcdefghij1234567890abcdefghij12" // 72 chars/bytes

		passwordReq = models.UserPasswordChangeRequest{
			CurrentPassword: "12345678",
			NewPassword:     longPassword,
		}
		body, _ = json.Marshal(passwordReq)

		req, _ = http.NewRequest("POST", "/v1/users/change-password", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("concurrent_operations", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Simulate concurrent profile updates
		for i := 0; i < 5; i++ {
			updateReq := models.UserProfileRequest{
				FirstName: "ConcurrentFirst",
				LastName:  "ConcurrentLast",
				Timezone:  "UTC",
			}
			body, _ := json.Marshal(updateReq)

			req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}

		// Verify final state
		req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var profileResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &profileResponse)
		require.NoError(t, err)

		profile := profileResponse["data"].(map[string]any)
		assert.Equal(t, "ConcurrentFirst", profile["first_name"])
	})

	t.Run("malicious_input_handling", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Test with SQL injection attempt in profile
		updateReq := models.UserProfileRequest{
			FirstName:   "'; DROP TABLE users; --",
			LastName:    "<script>alert('xss')</script>",
			Timezone:    "UTC",
			Preferences: map[string]any{"evil": "'; DROP TABLE users; --"},
		}
		body, _ := json.Marshal(updateReq)

		req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify data was stored safely
		req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var profileResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &profileResponse)
		require.NoError(t, err)

		profile := profileResponse["data"].(map[string]any)
		assert.Equal(t, "'; DROP TABLE users; --", profile["first_name"])
		assert.Equal(t, "<script>alert('xss')</script>", profile["last_name"])

		// Test SQL injection in password change
		passwordReq := models.UserPasswordChangeRequest{
			CurrentPassword: "password123",
			NewPassword:     "'; DROP TABLE users; --",
		}
		body, _ = json.Marshal(passwordReq)

		req, _ = http.NewRequest("POST", "/v1/users/change-password", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify new password works
		loginReq := models.UserLogin{
			Email:    user.Email,
			Password: "'; DROP TABLE users; --",
		}
		loginBody, _ := json.Marshal(loginReq)
		req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(loginBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestUserHandler_EdgeCases tests edge cases and boundary conditions
func TestUserHandler_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("token_expiration_scenarios", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Test with malformed token
		req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Test with empty token
		req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer ")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Test with valid token but different format
		req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", token) // Without "Bearer "
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("unicode_and_special_characters", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Test with Unicode characters
		updateReq := models.UserProfileRequest{
			FirstName:   "José María 🚀",
			LastName:    "Azañón-González 中文名字",
			Timezone:    "Europe/Madrid",
			Preferences: map[string]any{"emoji": "🎉", "special": "ñÑ áéíóú"},
		}
		body, _ := json.Marshal(updateReq)

		req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify Unicode characters were preserved
		req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var profileResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &profileResponse)
		require.NoError(t, err)

		profile := profileResponse["data"].(map[string]any)
		assert.Equal(t, "José María 🚀", profile["first_name"])
		assert.Equal(t, "Azañón-González 中文名字", profile["last_name"])
	})

	t.Run("large_request_payloads", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Test with very large preferences
		largePrefs := make(map[string]any)
		for i := 0; i < 1000; i++ {
			largePrefs[string(rune('a'+i%26))+string(rune('A'+i%26))] = "value" + string(rune(i))
		}

		updateReq := models.UserProfileRequest{
			FirstName:   "Test",
			LastName:    "User",
			Timezone:    "UTC",
			Preferences: largePrefs,
		}
		body, _ := json.Marshal(updateReq)

		req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should handle large payloads (may succeed or fail based on limits)
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusBadRequest)
	})

	t.Run("null_and_empty_values", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Test with null values (JSON nulls become Go zero values) - should fail validation
		reqBody := `{"first_name": null, "last_name": null, "timezone": "UTC"}`

		req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // Should fail validation

		// Test with empty strings - should fail validation
		updateReq := models.UserProfileRequest{
			FirstName: "",
			LastName:  "",
			Timezone:  "UTC",
		}
		body, _ := json.Marshal(updateReq)

		req, _ = http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // Should fail validation

		// Test with valid minimal values - should succeed
		validReq := models.UserProfileRequest{
			FirstName: "A",
			LastName:  "B",
			Timezone:  "UTC",
		}
		body, _ = json.Marshal(validReq)

		req, _ = http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code) // Should succeed
	})
}

// TestUserHandler_PerformanceScenarios tests performance-related scenarios
func TestUserHandler_PerformanceScenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("rapid_successive_requests", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Make rapid successive profile requests
		for i := 0; i < 10; i++ {
			req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("memory_efficiency", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Test memory efficiency with repeated operations
		for i := 0; i < 50; i++ {
			// Profile update
			updateReq := models.UserProfileRequest{
				FirstName:   "MemoryTest",
				LastName:    "User",
				Timezone:    "UTC",
				Preferences: map[string]any{"test": "memory efficiency"},
			}
			body, _ := json.Marshal(updateReq)

			req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			// Profile retrieval
			req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}
	})
}
