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

// TestUserHandler_GetProfile tests the GetProfile functionality
func TestUserHandler_GetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful profile retrieval", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Get user profile
		req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[*models.User]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		assert.Equal(t, user.Email, response.Data.Email)
		assert.Equal(t, user.FirstName, response.Data.FirstName)
		assert.Equal(t, user.LastName, response.Data.LastName)
		assert.NotEmpty(t, response.Data.ID)
		assert.False(t, response.Data.CreatedAt.IsZero())
		// Password hash should not be returned
		assert.Empty(t, response.Data.PasswordHash)
	})

	t.Run("unauthorized access", func(t *testing.T) {
		// Setup test environment
		router, _, _, _, _ := setupTestRouterWithServices(t)

		// Try to get profile without authentication
		req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// TestUserHandler_UpdateProfile tests the UpdateProfile functionality
func TestUserHandler_UpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful profile update", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Prepare profile update request
		updateRequest := models.UserProfileRequest{
			FirstName: "UpdatedFirstName",
			LastName:  "UpdatedLastName",
			Timezone:  "UTC",
		}

		body, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		// Update profile
		req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[*models.User]
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		assert.Equal(t, updateRequest.FirstName, response.Data.FirstName)
		assert.Equal(t, updateRequest.LastName, response.Data.LastName)
		assert.Equal(t, user.Email, response.Data.Email) // Email should remain unchanged
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Send invalid JSON
		req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request format")
	})

	t.Run("unauthorized update", func(t *testing.T) {
		// Setup test environment
		router, _, _, _, _ := setupTestRouterWithServices(t)

		updateRequest := models.UserProfileRequest{
			FirstName: "UpdatedFirstName",
			LastName:  "UpdatedLastName",
			Timezone:  "UTC",
		}

		body, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		// Try to update without authentication
		req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// TestUserHandler_ChangePassword tests the ChangePassword functionality
func TestUserHandler_ChangePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful password change", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Prepare password change request
		passwordRequest := models.UserPasswordChangeRequest{
			CurrentPassword: "password123",
			NewPassword:     "newpassword456",
		}

		body, err := json.Marshal(passwordRequest)
		require.NoError(t, err)

		// Change password
		req, _ := http.NewRequest("POST", "/v1/users/change-password", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response gin.H
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["message"], "Password changed successfully")

		// Verify old password no longer works
		loginRequest := models.UserLogin{
			Email:    user.Email,
			Password: "password123",
		}
		loginBody, _ := json.Marshal(loginRequest)
		req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(loginBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Verify new password works
		loginRequest.Password = "newpassword456"
		loginBody, _ = json.Marshal(loginRequest)
		req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(loginBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid current password", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Prepare password change request with wrong current password
		passwordRequest := models.UserPasswordChangeRequest{
			CurrentPassword: "wrongpassword",
			NewPassword:     "newpassword456",
		}

		body, err := json.Marshal(passwordRequest)
		require.NoError(t, err)

		// Try to change password
		req, _ := http.NewRequest("POST", "/v1/users/change-password", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "current password is incorrect")
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Send invalid JSON
		req, _ := http.NewRequest("POST", "/v1/users/change-password", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request format")
	})

	t.Run("unauthorized password change", func(t *testing.T) {
		// Setup test environment
		router, _, _, _, _ := setupTestRouterWithServices(t)

		passwordRequest := models.UserPasswordChangeRequest{
			CurrentPassword: "password123",
			NewPassword:     "newpassword456",
		}

		body, err := json.Marshal(passwordRequest)
		require.NoError(t, err)

		// Try to change password without authentication
		req, _ := http.NewRequest("POST", "/v1/users/change-password", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// TestUserHandler_DeleteAccount tests the DeleteAccount functionality
func TestUserHandler_DeleteAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful account deletion", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Delete account
		req, _ := http.NewRequest("DELETE", "/v1/users/account", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["message"], "User account deleted successfully")

		// Verify user can no longer login
		loginRequest := models.UserLogin{
			Email:    user.Email,
			Password: "password123",
		}
		loginBody, _ := json.Marshal(loginRequest)
		req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(loginBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("unauthorized account deletion", func(t *testing.T) {
		// Setup test environment
		router, _, _, _, _ := setupTestRouterWithServices(t)

		// Try to delete account without authentication
		req, _ := http.NewRequest("DELETE", "/v1/users/account", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// TestUserHandler_ErrorHandling tests various error scenarios
func TestUserHandler_ErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("profile retrieval with invalid token", func(t *testing.T) {
		// Setup test environment
		router, _, _, _, _ := setupTestRouterWithServices(t)

		// Try to get profile with invalid token
		req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("profile update with empty body", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Try to update with empty body
		req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should still work as empty fields are allowed
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("password change with missing fields", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := setupTestRouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Try to change password with missing fields
		passwordRequest := map[string]string{
			"current_password": "password123",
			// missing new_password
		}

		body, err := json.Marshal(passwordRequest)
		require.NoError(t, err)

		req, _ := http.NewRequest("POST", "/v1/users/change-password", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
