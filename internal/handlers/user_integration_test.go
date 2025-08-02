//go:build integration
// +build integration

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserHandler_Integration_FullWorkflow tests end-to-end user workflows
func TestUserHandler_Integration_FullWorkflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("complete_user_lifecycle_workflow", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, logEntryService, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// 1. Get initial profile
		req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var initialResponse responseData[*models.User]
		err := json.Unmarshal(w.Body.Bytes(), &initialResponse)
		require.NoError(t, err)
		require.NotNil(t, initialResponse.Data)
		assert.Equal(t, user.Email, initialResponse.Data.Email)
		assert.Equal(t, user.FirstName, initialResponse.Data.FirstName)

		// 2. Create some user activity (projects and log entries)
		project1 := createTestProject(t, projectService, user.ID.String())
		project2 := createTestProject(t, projectService, user.ID.String())

		_ = createTestLogEntry(t, logEntryService, user.ID.String(), &project1.ID)
		_ = createTestLogEntry(t, logEntryService, user.ID.String(), &project2.ID)
		_ = createTestLogEntry(t, logEntryService, user.ID.String(), &project1.ID)

		// 3. Update user profile
		updateReq := models.UserProfileRequest{
			FirstName:   "UpdatedFirst",
			LastName:    "UpdatedLast",
			Timezone:    "America/New_York",
			Preferences: map[string]any{"theme": "dark", "notifications": true},
		}
		body, _ := json.Marshal(updateReq)

		req, _ = http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 4. Verify profile was updated
		req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var updatedResponse responseData[*models.User]
		err = json.Unmarshal(w.Body.Bytes(), &updatedResponse)
		require.NoError(t, err)
		require.NotNil(t, updatedResponse.Data)
		assert.Equal(t, "UpdatedFirst", updatedResponse.Data.FirstName)
		assert.Equal(t, "UpdatedLast", updatedResponse.Data.LastName)
		assert.Equal(t, "America/New_York", updatedResponse.Data.Timezone)
		assert.Equal(t, "dark", updatedResponse.Data.Preferences["theme"])

		// 5. Change password
		passwordReq := models.UserPasswordChangeRequest{
			CurrentPassword: "password123",
			NewPassword:     "newSecurePassword456",
		}
		body, _ = json.Marshal(passwordReq)

		req, _ = http.NewRequest("POST", "/v1/users/change-password", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 6. Verify old password no longer works
		loginReq := models.UserLogin{Email: user.Email, Password: "password123"}
		body, _ = json.Marshal(loginReq)

		req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// 7. Verify new password works
		newToken := loginUser(t, router, user.Email, "newSecurePassword456")

		// 8. Use new token to access profile
		req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+newToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var finalResponse responseData[*models.User]
		err = json.Unmarshal(w.Body.Bytes(), &finalResponse)
		require.NoError(t, err)
		require.NotNil(t, finalResponse.Data)
		assert.Equal(t, "UpdatedFirst", finalResponse.Data.FirstName)
		assert.Equal(t, "UpdatedLast", finalResponse.Data.LastName)

		// 9. Final step: Delete account
		req, _ = http.NewRequest("DELETE", "/v1/users/account", nil)
		req.Header.Set("Authorization", "Bearer "+newToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 10. Verify account is deleted - login should fail
		loginReq = models.UserLogin{Email: user.Email, Password: "newSecurePassword456"}
		body, _ = json.Marshal(loginReq)

		req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// 11. Verify access with deleted account token fails
		req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+newToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Accept either 401 (unauthorized) or 500 (server error) since user is deleted
		assert.True(t, w.Code == http.StatusUnauthorized || w.Code == http.StatusInternalServerError)
	})
}

// TestUserHandler_Integration_CrossSystemBehavior tests user interactions with other system components
func TestUserHandler_Integration_CrossSystemBehavior(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("user_profile_updates_affect_system_data", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, logEntryService, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create projects and log entries before profile change
		project := createTestProject(t, projectService, user.ID.String())
		logEntry := createTestLogEntry(t, logEntryService, user.ID.String(), &project.ID)

		// Verify initial user data in log entry
		req, _ := http.NewRequest("GET", "/v1/logs/"+logEntry.ID.String(), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Update user profile
		updateReq := models.UserProfileRequest{
			FirstName:   "NewFirstName",
			LastName:    "NewLastName",
			Timezone:    "Europe/London",
			Preferences: map[string]any{"language": "en"},
		}
		body, _ := json.Marshal(updateReq)

		req, _ = http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify user profile is updated
		req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var profileResponse responseData[*models.User]
		err := json.Unmarshal(w.Body.Bytes(), &profileResponse)
		require.NoError(t, err)
		require.NotNil(t, profileResponse.Data)
		assert.Equal(t, "NewFirstName", profileResponse.Data.FirstName)
		assert.Equal(t, "NewLastName", profileResponse.Data.LastName)
		assert.Equal(t, "Europe/London", profileResponse.Data.Timezone)
	})

	t.Run("user_account_deletion_affects_related_data", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, logEntryService, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create user's projects and log entries
		project1 := createTestProject(t, projectService, user.ID.String())
		project2 := createTestProject(t, projectService, user.ID.String())

		_ = createTestLogEntry(t, logEntryService, user.ID.String(), &project1.ID)
		_ = createTestLogEntry(t, logEntryService, user.ID.String(), &project2.ID)

		// Verify data exists before deletion
		req, _ := http.NewRequest("GET", "/v1/projects", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		req, _ = http.NewRequest("GET", "/v1/logs", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Delete user account
		req, _ = http.NewRequest("DELETE", "/v1/users/account", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify user can't login anymore
		loginReq := models.UserLogin{Email: user.Email, Password: "password123"}
		body, _ := json.Marshal(loginReq)

		req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Verify old token no longer works
		req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Accept either 401 (unauthorized) or 500 (server error) since user is deleted
		assert.True(t, w.Code == http.StatusUnauthorized || w.Code == http.StatusInternalServerError)
	})

	t.Run("user_timezone_changes_affect_timestamps", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, logEntryService, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create project and log entry
		project := createTestProject(t, projectService, user.ID.String())
		logEntry := createTestLogEntry(t, logEntryService, user.ID.String(), &project.ID)

		// Update user timezone
		updateReq := models.UserProfileRequest{
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Timezone:    "Asia/Tokyo",
			Preferences: map[string]any{"timezone_format": "24h"},
		}
		body, _ := json.Marshal(updateReq)

		req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify timezone was updated
		req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var timezoneResponse responseData[*models.User]
		err := json.Unmarshal(w.Body.Bytes(), &timezoneResponse)
		require.NoError(t, err)
		require.NotNil(t, timezoneResponse.Data)
		assert.Equal(t, "Asia/Tokyo", timezoneResponse.Data.Timezone)

		// Verify log entry is still accessible with new timezone preference
		req, _ = http.NewRequest("GET", "/v1/logs/"+logEntry.ID.String(), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestUserHandler_Integration_SecurityBehavior tests security-related integrations
func TestUserHandler_Integration_SecurityBehavior(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("token_behavior_after_password_change", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		oldToken := loginUser(t, router, user.Email, "password123")

		// Verify old token works
		req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+oldToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Change password using old token
		passwordReq := models.UserPasswordChangeRequest{
			CurrentPassword: "password123",
			NewPassword:     "superSecureNewPassword789",
		}
		body, _ := json.Marshal(passwordReq)

		req, _ = http.NewRequest("POST", "/v1/users/change-password", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+oldToken)
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Login with new password to get new token
		newToken := loginUser(t, router, user.Email, "superSecureNewPassword789")

		// Verify new token works
		req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+newToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Test if old token still works (implementation dependent)
		// Some implementations invalidate old tokens on password change, others don't
		req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+oldToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		// Accept either 200 (old token still valid) or 401 (old token invalidated)
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusUnauthorized,
			"Expected status 200 or 401, got %d", w.Code)
	})

	t.Run("concurrent_user_operations", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Perform multiple concurrent profile updates
		results := make(chan int, 3)

		// First update
		go func() {
			updateReq := models.UserProfileRequest{
				FirstName:   "Concurrent1",
				LastName:    "Test1",
				Timezone:    "UTC",
				Preferences: map[string]any{"test": 1},
			}
			body, _ := json.Marshal(updateReq)

			req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			results <- w.Code
		}()

		// Second update
		go func() {
			updateReq := models.UserProfileRequest{
				FirstName:   "Concurrent2",
				LastName:    "Test2",
				Timezone:    "America/Los_Angeles",
				Preferences: map[string]any{"test": 2},
			}
			body, _ := json.Marshal(updateReq)

			req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			results <- w.Code
		}()

		// Third update
		go func() {
			updateReq := models.UserProfileRequest{
				FirstName:   "Concurrent3",
				LastName:    "Test3",
				Timezone:    "Europe/Berlin",
				Preferences: map[string]any{"test": 3},
			}
			body, _ := json.Marshal(updateReq)

			req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			results <- w.Code
		}()

		// Wait for all updates to complete
		successCount := 0
		for i := 0; i < 3; i++ {
			code := <-results
			// Accept 200 (success) or 400 (conflict/validation error in concurrent scenario)
			if code == http.StatusOK || code == http.StatusBadRequest {
				successCount++
			}
		}

		// Most updates should succeed (allow some conflicts in concurrent scenario)
		assert.GreaterOrEqual(t, successCount, 1, "At least one concurrent update should succeed")

		// Verify final profile state is consistent
		req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var concurrentResponse responseData[*models.User]
		err := json.Unmarshal(w.Body.Bytes(), &concurrentResponse)
		require.NoError(t, err)
		require.NotNil(t, concurrentResponse.Data)

		// Profile should have data from one of the updates
		assert.True(t,
			(concurrentResponse.Data.FirstName == "Concurrent1" && concurrentResponse.Data.LastName == "Test1") ||
				(concurrentResponse.Data.FirstName == "Concurrent2" && concurrentResponse.Data.LastName == "Test2") ||
				(concurrentResponse.Data.FirstName == "Concurrent3" && concurrentResponse.Data.LastName == "Test3"))
	})

	t.Run("password_change_with_bcrypt_limits", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Test with maximum bcrypt password length (72 bytes)
		// Note: bcrypt has a 72-byte password limit, so we test with exactly 72 bytes
		maxPassword := "abcdefghij1234567890abcdefghij1234567890abcdefghij1234567890abcdefghij12" // exactly 72 bytes

		passwordReq := models.UserPasswordChangeRequest{
			CurrentPassword: "password123",
			NewPassword:     maxPassword,
		}
		body, _ := json.Marshal(passwordReq)

		req, _ := http.NewRequest("POST", "/v1/users/change-password", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify login with new password works
		loginReq := models.UserLogin{Email: user.Email, Password: maxPassword}
		body, _ = json.Marshal(loginReq)

		req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Test with password that exceeds bcrypt limit (73+ bytes) - should fail validation
		tooLongPassword := maxPassword + "x" // 73 bytes

		passwordReq = models.UserPasswordChangeRequest{
			CurrentPassword: maxPassword,
			NewPassword:     tooLongPassword,
		}
		body, _ = json.Marshal(passwordReq)

		req, _ = http.NewRequest("POST", "/v1/users/change-password", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should fail due to password length validation
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestUserHandler_Integration_PerformanceAndReliability tests performance and reliability aspects
func TestUserHandler_Integration_PerformanceAndReliability(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("bulk_user_operations_performance", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Measure performance of rapid profile reads
		start := time.Now()
		successCount := 0

		for i := 0; i < 20; i++ {
			req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code == http.StatusOK {
				successCount++
			}
		}

		duration := time.Since(start)

		// All requests should succeed
		assert.Equal(t, 20, successCount)

		// Should complete within reasonable time
		assert.True(t, duration < 10*time.Second, "20 profile reads took %v, expected < 10s", duration)

		// Calculate average response time
		avgTime := duration / 20
		t.Logf("Average profile read time: %v", avgTime)
	})

	t.Run("profile_update_resilience", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Test updating profile multiple times rapidly
		updateCount := 0
		successCount := 0

		for i := 0; i < 5; i++ {
			updateReq := models.UserProfileRequest{
				FirstName:   "UpdatedFirst" + string(rune('A'+i)),
				LastName:    "UpdatedLast" + string(rune('A'+i)),
				Timezone:    "UTC",
				Preferences: map[string]any{"iteration": i},
			}
			body, _ := json.Marshal(updateReq)

			req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			updateCount++
			if w.Code == http.StatusOK {
				successCount++
			}

			// Small delay to avoid overwhelming the system
			time.Sleep(10 * time.Millisecond)
		}

		// All updates should succeed
		assert.Equal(t, updateCount, successCount)

		// Verify final state is consistent
		req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resilienceResponse responseData[*models.User]
		err := json.Unmarshal(w.Body.Bytes(), &resilienceResponse)
		require.NoError(t, err)
		require.NotNil(t, resilienceResponse.Data)

		// Should have data from the last update
		assert.Equal(t, "UpdatedFirstE", resilienceResponse.Data.FirstName)
		assert.Equal(t, "UpdatedLastE", resilienceResponse.Data.LastName)
		assert.Equal(t, float64(4), resilienceResponse.Data.Preferences["iteration"])
	})

	t.Run("user_data_consistency_under_load", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, logEntryService, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create baseline data
		project := createTestProject(t, projectService, user.ID.String())
		logEntry := createTestLogEntry(t, logEntryService, user.ID.String(), &project.ID)

		// Perform multiple concurrent operations
		results := make(chan bool, 6)

		// Concurrent profile reads
		go func() {
			req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			results <- w.Code == http.StatusOK
		}()

		// Concurrent profile update
		go func() {
			updateReq := models.UserProfileRequest{
				FirstName:   "ConcurrentUpdate",
				LastName:    "LoadTest",
				Timezone:    "UTC",
				Preferences: map[string]any{"load_test": true},
			}
			body, _ := json.Marshal(updateReq)

			req, _ := http.NewRequest("PUT", "/v1/users/profile", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			results <- w.Code == http.StatusOK
		}()

		// Concurrent project access
		go func() {
			req, _ := http.NewRequest("GET", "/v1/projects/"+project.ID.String(), nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			results <- w.Code == http.StatusOK
		}()

		// Concurrent log entry access
		go func() {
			req, _ := http.NewRequest("GET", "/v1/logs/"+logEntry.ID.String(), nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			results <- w.Code == http.StatusOK
		}()

		// Additional profile reads
		go func() {
			req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			results <- w.Code == http.StatusOK
		}()

		go func() {
			req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			results <- w.Code == http.StatusOK
		}()

		// Wait for all operations to complete
		successCount := 0
		for i := 0; i < 6; i++ {
			if <-results {
				successCount++
			}
		}

		// All operations should succeed
		assert.Equal(t, 6, successCount)

		// Verify final system state is consistent
		req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var loadResponse responseData[*models.User]
		err := json.Unmarshal(w.Body.Bytes(), &loadResponse)
		require.NoError(t, err)
		require.NotNil(t, loadResponse.Data)

		// Profile should be in a consistent state
		assert.NotEmpty(t, loadResponse.Data.FirstName)
		assert.NotEmpty(t, loadResponse.Data.LastName)
		assert.Equal(t, user.Email, loadResponse.Data.Email)
	})
}
