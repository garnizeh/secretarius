package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLogEntryHandler_CreateLogEntry tests the CreateLogEntry functionality
func TestLogEntryHandler_CreateLogEntry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful log entry creation", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create project for the log entry
		project := createTestProject(t, projectService, user.ID.String())

		// Prepare log entry creation request
		description := "Test log entry for integration testing"
		now := time.Now()
		logRequest := models.LogEntryRequest{
			Title:       fmt.Sprintf("test-entry-%d", time.Now().Unix()),
			Description: &description,
			Type:        models.ActivityDevelopment,
			ProjectID:   &project.ID,
			StartTime:   now.Add(-2 * time.Hour),
			EndTime:     now.Add(-1 * time.Hour),
			ValueRating: models.ValueHigh,
			ImpactLevel: models.ImpactTeam,
			Tags:        []string{"testing", "development"},
		}

		body, err := json.Marshal(logRequest)
		require.NoError(t, err)

		// Create log entry
		req, _ := http.NewRequest("POST", "/v1/logs", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Debug: print response body if test fails
		if w.Code != http.StatusCreated {
			t.Logf("Response status: %d, body: %s", w.Code, w.Body.String())
		}

		assert.Equal(t, http.StatusCreated, w.Code)

		var response responseData[*models.LogEntry]
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		assert.Equal(t, logRequest.Title, response.Data.Title)
		assert.Equal(t, *logRequest.Description, *response.Data.Description)
		assert.Equal(t, logRequest.Type, response.Data.Type)
		assert.Equal(t, logRequest.ProjectID, response.Data.ProjectID)
		assert.Equal(t, logRequest.ValueRating, response.Data.ValueRating)
		assert.Equal(t, logRequest.ImpactLevel, response.Data.ImpactLevel)
		assert.Equal(t, logRequest.Tags, response.Data.Tags)
		assert.NotEmpty(t, response.Data.ID)
		assert.False(t, response.Data.CreatedAt.IsZero())
	})

	t.Run("unauthorized access", func(t *testing.T) {
		// Setup test environment
		router, _, _, _, _ := RouterWithServices(t)

		// Prepare log entry creation request
		logRequest := models.LogEntryRequest{
			Title:       "Test Entry",
			Type:        models.ActivityDevelopment,
			StartTime:   time.Now().Add(-time.Hour),
			EndTime:     time.Now(),
			ValueRating: models.ValueMedium,
			ImpactLevel: models.ImpactTeam,
		}

		body, err := json.Marshal(logRequest)
		require.NoError(t, err)

		// Create log entry without authorization
		req, _ := http.NewRequest("POST", "/v1/logs", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var errorResponse gin.H
		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Unauthorized")
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Send invalid JSON
		req, _ := http.NewRequest("POST", "/v1/logs", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Invalid request format")
	})

	t.Run("missing required fields", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create request with missing required fields
		logRequest := models.LogEntryRequest{
			// Missing Title, Type, StartTime, EndTime, etc.
		}

		body, err := json.Marshal(logRequest)
		require.NoError(t, err)

		req, _ := http.NewRequest("POST", "/v1/logs", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestLogEntryHandler_GetLogEntry tests the GetLogEntry functionality
func TestLogEntryHandler_GetLogEntry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful log entry retrieval", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, logEntryService, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create project and log entry
		project := createTestProject(t, projectService, user.ID.String())
		logEntry := createTestLogEntry(t, logEntryService, user.ID.String(), &project.ID)

		// Get the log entry
		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/logs/%s", logEntry.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[*models.LogEntry]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		assert.Equal(t, logEntry.ID, response.Data.ID)
		assert.Equal(t, logEntry.Title, response.Data.Title)
		assert.Equal(t, logEntry.Type, response.Data.Type)
		assert.Equal(t, logEntry.ProjectID, response.Data.ProjectID)
	})

	t.Run("unauthorized access", func(t *testing.T) {
		// Setup test environment
		router, _, _, _, _ := RouterWithServices(t)

		// Try to get log entry without authorization
		req, _ := http.NewRequest("GET", "/v1/logs/test-id", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("log entry not found", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Try to get non-existent log entry
		req, _ := http.NewRequest("GET", "/v1/logs/00000000-0000-0000-0000-000000000000", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var errorResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Log entry not found")
	})

	t.Run("missing log entry ID", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Try to get log entry with empty ID
		req, _ := http.NewRequest("GET", "/v1/logs/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 404 because the route doesn't match
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestLogEntryHandler_GetLogEntries tests the GetLogEntries functionality
func TestLogEntryHandler_GetLogEntries(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful log entries retrieval", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, logEntryService, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create project and multiple log entries
		project := createTestProject(t, projectService, user.ID.String())
		logEntry1 := createTestLogEntry(t, logEntryService, user.ID.String(), &project.ID)
		logEntry2 := createTestLogEntry(t, logEntryService, user.ID.String(), &project.ID)

		// Get all log entries
		req, _ := http.NewRequest("GET", "/v1/logs", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data       []*models.LogEntry `json:"data"`
			Pagination map[string]any     `json:"pagination"`
			Total      int                `json:"total"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		// Should contain at least our created log entries
		assert.GreaterOrEqual(t, len(response.Data), 2)
		assert.GreaterOrEqual(t, response.Total, 2)

		// Find our log entries in the response
		foundEntry1, foundEntry2 := false, false
		for _, entry := range response.Data {
			if entry.ID == logEntry1.ID {
				foundEntry1 = true
			}
			if entry.ID == logEntry2.ID {
				foundEntry2 = true
			}
		}
		assert.True(t, foundEntry1, "First created log entry should be in response")
		assert.True(t, foundEntry2, "Second created log entry should be in response")

		// Verify pagination structure
		assert.NotNil(t, response.Pagination)
		assert.Contains(t, response.Pagination, "page")
		assert.Contains(t, response.Pagination, "limit")
		assert.Contains(t, response.Pagination, "total_pages")
	})

	t.Run("unauthorized access", func(t *testing.T) {
		// Setup test environment
		router, _, _, _, _ := RouterWithServices(t)

		// Try to get log entries without authorization
		req, _ := http.NewRequest("GET", "/v1/logs", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var errorResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Unauthorized")
	})

	t.Run("with filters", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, logEntryService, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create project and log entry
		project := createTestProject(t, projectService, user.ID.String())
		_ = createTestLogEntry(t, logEntryService, user.ID.String(), &project.ID)

		// Get log entries with filters
		url := fmt.Sprintf("/v1/logs?project_id=%s&type=development&value_rating=medium", project.ID)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data       []*models.LogEntry `json:"data"`
			Pagination map[string]any     `json:"pagination"`
			Total      int                `json:"total"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
	})

	t.Run("with invalid filters", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Get log entries with invalid filters
		req, _ := http.NewRequest("GET", "/v1/logs?type=invalid_type", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Invalid query parameters")
	})

	t.Run("with pagination", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Get log entries with pagination
		req, _ := http.NewRequest("GET", "/v1/logs?page=1&limit=10", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data       []*models.LogEntry `json:"data"`
			Pagination map[string]any     `json:"pagination"`
			Total      int                `json:"total"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify pagination values
		assert.Equal(t, float64(1), response.Pagination["page"])
		assert.Equal(t, float64(10), response.Pagination["limit"])
	})
}

// TestLogEntryHandler_UpdateLogEntry tests the UpdateLogEntry functionality
func TestLogEntryHandler_UpdateLogEntry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful log entry update", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, logEntryService, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create project and log entry
		project := createTestProject(t, projectService, user.ID.String())
		logEntry := createTestLogEntry(t, logEntryService, user.ID.String(), &project.ID)

		// Prepare update request
		updatedDescription := "Updated log entry description"
		updateRequest := models.LogEntryRequest{
			Title:       "Updated Test Entry",
			Description: &updatedDescription,
			Type:        models.ActivityMeeting,
			ProjectID:   &project.ID,
			StartTime:   logEntry.StartTime,
			EndTime:     logEntry.EndTime,
			ValueRating: models.ValueLow,
			ImpactLevel: models.ImpactPersonal,
			Tags:        []string{"updated", "testing"},
		}

		body, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		// Update log entry
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/logs/%s", logEntry.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[*models.LogEntry]
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		assert.Equal(t, updateRequest.Title, response.Data.Title)
		assert.Equal(t, *updateRequest.Description, *response.Data.Description)
		assert.Equal(t, updateRequest.Type, response.Data.Type)
		assert.Equal(t, updateRequest.ValueRating, response.Data.ValueRating)
		assert.Equal(t, updateRequest.ImpactLevel, response.Data.ImpactLevel)
		assert.Equal(t, updateRequest.Tags, response.Data.Tags)
	})

	t.Run("unauthorized access", func(t *testing.T) {
		// Setup test environment
		router, _, _, _, _ := RouterWithServices(t)

		// Prepare update request
		updateRequest := models.LogEntryRequest{
			Title: "Updated Test Entry",
			Type:  models.ActivityMeeting,
		}

		body, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		// Try to update without authorization
		req, _ := http.NewRequest("PUT", "/v1/logs/test-id", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("log entry not found", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Prepare update request
		updateRequest := models.LogEntryRequest{
			Title: "Updated Test Entry",
			Type:  models.ActivityMeeting,
		}

		body, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		// Try to update non-existent log entry
		req, _ := http.NewRequest("PUT", "/v1/logs/00000000-0000-0000-0000-000000000000", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResponse gin.H
		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Failed to update log entry")
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, logEntryService, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create project and log entry
		project := createTestProject(t, projectService, user.ID.String())
		logEntry := createTestLogEntry(t, logEntryService, user.ID.String(), &project.ID)

		// Send invalid JSON
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/logs/%s", logEntry.ID), bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Invalid request format")
	})

	t.Run("missing log entry ID", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Prepare update request
		updateRequest := models.LogEntryRequest{
			Title: "Updated Test Entry",
			Type:  models.ActivityMeeting,
		}

		body, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		// Try to update with empty ID
		req, _ := http.NewRequest("PUT", "/v1/logs/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 404 because the route doesn't match
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestLogEntryHandler_DeleteLogEntry tests the DeleteLogEntry functionality
func TestLogEntryHandler_DeleteLogEntry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful log entry deletion", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, logEntryService, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create project and log entry
		project := createTestProject(t, projectService, user.ID.String())
		logEntry := createTestLogEntry(t, logEntryService, user.ID.String(), &project.ID)

		// Delete log entry
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/logs/%s", logEntry.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["message"], "Log entry deleted successfully")
	})

	t.Run("unauthorized access", func(t *testing.T) {
		// Setup test environment
		router, _, _, _, _ := RouterWithServices(t)

		// Try to delete without authorization
		req, _ := http.NewRequest("DELETE", "/v1/logs/test-id", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("log entry not found", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Try to delete non-existent log entry
		req, _ := http.NewRequest("DELETE", "/v1/logs/00000000-0000-0000-0000-000000000000", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var errorResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Failed to delete log entry")
	})

	t.Run("missing log entry ID", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Try to delete with empty ID
		req, _ := http.NewRequest("DELETE", "/v1/logs/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 404 because the route doesn't match
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestLogEntryHandler_BulkCreateLogEntries tests the BulkCreateLogEntries functionality
func TestLogEntryHandler_BulkCreateLogEntries(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful bulk creation", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create project for the log entries
		project := createTestProject(t, projectService, user.ID.String())

		// Prepare bulk creation request
		now := time.Now()
		bulkRequest := struct {
			Entries []models.LogEntryRequest `json:"entries"`
		}{
			Entries: []models.LogEntryRequest{
				{
					Title:       "Bulk Entry 1",
					Type:        models.ActivityDevelopment,
					ProjectID:   &project.ID,
					StartTime:   now.Add(-3 * time.Hour),
					EndTime:     now.Add(-2 * time.Hour),
					ValueRating: models.ValueHigh,
					ImpactLevel: models.ImpactTeam,
				},
				{
					Title:       "Bulk Entry 2",
					Type:        models.ActivityMeeting,
					ProjectID:   &project.ID,
					StartTime:   now.Add(-2 * time.Hour),
					EndTime:     now.Add(-1 * time.Hour),
					ValueRating: models.ValueMedium,
					ImpactLevel: models.ImpactPersonal,
				},
			},
		}

		body, err := json.Marshal(bulkRequest)
		require.NoError(t, err)

		// Create log entries in bulk
		req, _ := http.NewRequest("POST", "/v1/logs/bulk", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response struct {
			Data    []any `json:"data"`
			Summary struct {
				Total   int `json:"total"`
				Success int `json:"success"`
				Errors  int `json:"errors"`
			} `json:"summary"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 2, response.Summary.Total)
		assert.Equal(t, 2, response.Summary.Success)
		assert.Equal(t, 0, response.Summary.Errors)
		assert.Len(t, response.Data, 2)
	})

	t.Run("unauthorized access", func(t *testing.T) {
		// Setup test environment
		router, _, _, _, _ := RouterWithServices(t)

		// Prepare bulk creation request
		bulkRequest := struct {
			Entries []models.LogEntryRequest `json:"entries"`
		}{
			Entries: []models.LogEntryRequest{
				{
					Title: "Bulk Entry 1",
					Type:  models.ActivityDevelopment,
				},
			},
		}

		body, err := json.Marshal(bulkRequest)
		require.NoError(t, err)

		// Try to create without authorization
		req, _ := http.NewRequest("POST", "/v1/logs/bulk", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var errorResponse gin.H
		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Unauthorized")
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Send invalid JSON
		req, _ := http.NewRequest("POST", "/v1/logs/bulk", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Invalid request format")
	})

	t.Run("empty entries array", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Prepare request with empty entries
		bulkRequest := struct {
			Entries []models.LogEntryRequest `json:"entries"`
		}{
			Entries: []models.LogEntryRequest{},
		}

		body, err := json.Marshal(bulkRequest)
		require.NoError(t, err)

		req, _ := http.NewRequest("POST", "/v1/logs/bulk", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResponse gin.H
		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Invalid request format")
	})

	t.Run("partial success", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create project for valid entries
		project := createTestProject(t, projectService, user.ID.String())

		// Prepare bulk request with one valid and one invalid entry
		now := time.Now()
		bulkRequest := struct {
			Entries []models.LogEntryRequest `json:"entries"`
		}{
			Entries: []models.LogEntryRequest{
				{
					Title:       "Valid Entry",
					Type:        models.ActivityDevelopment,
					ProjectID:   &project.ID,
					StartTime:   now.Add(-2 * time.Hour),
					EndTime:     now.Add(-1 * time.Hour),
					ValueRating: models.ValueHigh,
					ImpactLevel: models.ImpactTeam,
				},
				{
					// Invalid entry - missing required fields
					Title: "Invalid Entry",
				},
			},
		}

		body, err := json.Marshal(bulkRequest)
		require.NoError(t, err)

		req, _ := http.NewRequest("POST", "/v1/logs/bulk", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMultiStatus, w.Code)

		var response struct {
			Data    []any `json:"data"`
			Summary struct {
				Total   int `json:"total"`
				Success int `json:"success"`
				Errors  int `json:"errors"`
			} `json:"summary"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 2, response.Summary.Total)
		assert.Equal(t, 1, response.Summary.Success)
		assert.Equal(t, 1, response.Summary.Errors)
		assert.Len(t, response.Data, 2)
	})
}

// TestParseLogEntryFilters tests the parseLogEntryFilters functionality
func TestParseLogEntryFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &LogEntryHandler{}

	t.Run("valid filters", func(t *testing.T) {
		// Create test context with query parameters
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/logs?start_date=2024-01-01&end_date=2024-01-31&type=development&value_rating=high&impact_level=team&tags=test,api&project_id=test-project", nil)

		// Execute
		filters, err := handler.parseLogEntryFilters(c)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, filters)

		assert.Equal(t, "2024-01-01", filters.StartDate.Format("2006-01-02"))
		assert.Equal(t, "2024-01-31", filters.EndDate.Format("2006-01-02"))
		assert.Equal(t, models.ActivityDevelopment, *filters.Type)
		assert.Equal(t, models.ValueHigh, *filters.ValueRating)
		assert.Equal(t, models.ImpactTeam, *filters.ImpactLevel)
		assert.Equal(t, []string{"test", "api"}, filters.Tags)
		assert.Equal(t, "test-project", *filters.ProjectID)
	})

	t.Run("no filters", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/logs", nil)

		filters, err := handler.parseLogEntryFilters(c)

		require.NoError(t, err)
		require.NotNil(t, filters)
		assert.Nil(t, filters.Type)
		assert.Nil(t, filters.ValueRating)
		assert.Nil(t, filters.ImpactLevel)
		assert.Nil(t, filters.ProjectID)
		assert.Nil(t, filters.Tags)
	})

	t.Run("invalid date format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/logs?start_date=invalid-date", nil)

		filters, err := handler.parseLogEntryFilters(c)

		assert.Error(t, err)
		assert.Nil(t, filters)
		assert.Contains(t, err.Error(), "invalid start_date format")
	})

	t.Run("invalid activity type", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/logs?type=invalid_type", nil)

		filters, err := handler.parseLogEntryFilters(c)

		assert.Error(t, err)
		assert.Nil(t, filters)
		assert.Contains(t, err.Error(), "invalid activity type")
	})

	t.Run("invalid value rating", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/logs?value_rating=invalid_rating", nil)

		filters, err := handler.parseLogEntryFilters(c)

		assert.Error(t, err)
		assert.Nil(t, filters)
		assert.Contains(t, err.Error(), "invalid value rating")
	})

	t.Run("invalid impact level", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/logs?impact_level=invalid_level", nil)

		filters, err := handler.parseLogEntryFilters(c)

		assert.Error(t, err)
		assert.Nil(t, filters)
		assert.Contains(t, err.Error(), "invalid impact level")
	})
}

// TestParsePagination tests the parsePagination functionality
func TestParsePagination(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &LogEntryHandler{}

	tests := []struct {
		name          string
		url           string
		expectedPage  int
		expectedLimit int
	}{
		{"default values", "/logs", 1, 50},
		{"custom values", "/logs?page=2&limit=20", 2, 20},
		{"invalid page", "/logs?page=0", 1, 50},
		{"negative page", "/logs?page=-1", 1, 50},
		{"limit too high", "/logs?page=1&limit=200", 1, 50},
		{"negative limit", "/logs?page=1&limit=-10", 1, 50},
		{"invalid page string", "/logs?page=abc", 1, 50},
		{"invalid limit string", "/logs?limit=xyz", 1, 50},
		{"max valid limit", "/logs?page=1&limit=100", 1, 100},
		{"boundary page", "/logs?page=1&limit=1", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", tt.url, nil)

			page, limit := handler.parsePagination(c)

			assert.Equal(t, tt.expectedPage, page, "page should match expected value")
			assert.Equal(t, tt.expectedLimit, limit, "limit should match expected value")
		})
	}
}

// TestPaginate tests the paginate functionality
func TestPaginate(t *testing.T) {
	handler := &LogEntryHandler{}

	// Create test data
	entries := make([]*models.LogEntry, 15)
	for i := 0; i < 15; i++ {
		entries[i] = &models.LogEntry{
			ID:    uuid.New(),
			Title: fmt.Sprintf("Entry %d", i+1),
		}
	}

	t.Run("normal pagination", func(t *testing.T) {
		// Test pagination
		paginatedEntries, pagination := handler.paginate(entries, 2, 5)

		// Assert results
		assert.Len(t, paginatedEntries, 5, "should return 5 entries")
		assert.Equal(t, 2, pagination["page"], "page should be 2")
		assert.Equal(t, 5, pagination["limit"], "limit should be 5")
		assert.Equal(t, 3, pagination["total_pages"], "should have 3 total pages")
		assert.Equal(t, true, pagination["has_next"], "should have next page")
		assert.Equal(t, true, pagination["has_prev"], "should have previous page")

		// Check that we got the correct entries (entries 5-9, 0-indexed)
		assert.Equal(t, "Entry 6", paginatedEntries[0].Title)
		assert.Equal(t, "Entry 10", paginatedEntries[4].Title)
	})

	t.Run("first page", func(t *testing.T) {
		paginatedEntries, pagination := handler.paginate(entries, 1, 5)

		assert.Len(t, paginatedEntries, 5)
		assert.Equal(t, 1, pagination["page"])
		assert.Equal(t, false, pagination["has_prev"])
		assert.Equal(t, true, pagination["has_next"])
		assert.Equal(t, "Entry 1", paginatedEntries[0].Title)
	})

	t.Run("last page", func(t *testing.T) {
		paginatedEntries, pagination := handler.paginate(entries, 3, 5)

		assert.Len(t, paginatedEntries, 5)
		assert.Equal(t, 3, pagination["page"])
		assert.Equal(t, true, pagination["has_prev"])
		assert.Equal(t, false, pagination["has_next"])
		assert.Equal(t, "Entry 11", paginatedEntries[0].Title)
		assert.Equal(t, "Entry 15", paginatedEntries[4].Title)
	})

	t.Run("page beyond available data", func(t *testing.T) {
		paginatedEntries, pagination := handler.paginate(entries, 10, 5)

		assert.Len(t, paginatedEntries, 0)
		assert.Equal(t, 10, pagination["page"])
		assert.Equal(t, true, pagination["has_prev"])
		assert.Equal(t, false, pagination["has_next"])
	})

	t.Run("partial last page", func(t *testing.T) {
		paginatedEntries, pagination := handler.paginate(entries, 2, 10)

		assert.Len(t, paginatedEntries, 5) // Only 5 entries remaining
		assert.Equal(t, 2, pagination["page"])
		assert.Equal(t, 2, pagination["total_pages"])
		assert.Equal(t, true, pagination["has_prev"])
		assert.Equal(t, false, pagination["has_next"])
	})

	t.Run("empty entries", func(t *testing.T) {
		emptyEntries := []*models.LogEntry{}
		paginatedEntries, pagination := handler.paginate(emptyEntries, 1, 5)

		assert.Len(t, paginatedEntries, 0)
		assert.Equal(t, 1, pagination["page"])
		assert.Equal(t, 0, pagination["total_pages"])
		assert.Equal(t, false, pagination["has_prev"])
		assert.Equal(t, false, pagination["has_next"])
	})
}
