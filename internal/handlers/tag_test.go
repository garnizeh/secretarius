package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTagHandler_CreateTag tests the CreateTag functionality
func TestTagHandler_CreateTag(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful tag creation", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Prepare tag creation request
		description := "Test tag for integration testing"
		tagRequest := models.TagRequest{
			Name:        fmt.Sprintf("test-tag-%d", time.Now().Unix()),
			Description: &description,
			Color:       "#FF5722",
		}

		body, err := json.Marshal(tagRequest)
		require.NoError(t, err)

		// Create tag
		req, _ := http.NewRequest("POST", "/v1/tags", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Debug: print response body if test fails
		if w.Code != http.StatusCreated {
			t.Logf("Response status: %d, body: %s", w.Code, w.Body.String())
		}

		assert.Equal(t, http.StatusCreated, w.Code)

		var response responseData[*models.Tag]
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		assert.Equal(t, tagRequest.Name, response.Data.Name)
		assert.Equal(t, *tagRequest.Description, *response.Data.Description)
		assert.Equal(t, tagRequest.Color, response.Data.Color)
		assert.NotEmpty(t, response.Data.ID)
		assert.False(t, response.Data.CreatedAt.IsZero())
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Send invalid JSON
		req, _ := http.NewRequest("POST", "/v1/tags", bytes.NewBuffer([]byte("invalid json")))
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
}

// TestTagHandler_GetTag tests the GetTag functionality
func TestTagHandler_GetTag(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful tag retrieval", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, tagService := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create a tag first
		tag := createTestTag(t, tagService)

		// Get the tag
		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/tags/%s", tag.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[*models.Tag]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		assert.Equal(t, tag.ID, response.Data.ID)
		assert.Equal(t, tag.Name, response.Data.Name)
		assert.Equal(t, tag.Description, response.Data.Description)
		assert.Equal(t, tag.Color, response.Data.Color)
	})

	t.Run("tag not found", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Try to get non-existent tag
		req, _ := http.NewRequest("GET", "/v1/tags/00000000-0000-0000-0000-000000000000", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var errorResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Tag not found")
	})
}

// TestTagHandler_GetTags tests the GetTags functionality
func TestTagHandler_GetTags(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful tags retrieval", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, tagService := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create multiple tags
		tag1 := createTestTag(t, tagService)
		tag2 := createTestTag(t, tagService)

		// Get all tags
		req, _ := http.NewRequest("GET", "/v1/tags", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[[]*models.Tag]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		// Should contain at least our created tags
		assert.GreaterOrEqual(t, len(response.Data), 2)

		// Find our tags in the response
		foundTag1, foundTag2 := false, false
		for _, tag := range response.Data {
			if tag.ID == tag1.ID {
				foundTag1 = true
			}
			if tag.ID == tag2.ID {
				foundTag2 = true
			}
		}
		assert.True(t, foundTag1, "First created tag should be in response")
		assert.True(t, foundTag2, "Second created tag should be in response")
	})
}

// TestTagHandler_GetPopularTags tests the GetPopularTags functionality
func TestTagHandler_GetPopularTags(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful popular tags retrieval", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Get popular tags
		req, _ := http.NewRequest("GET", "/v1/tags/popular", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[[]*models.Tag]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		// Response should be valid (can be empty for new user)
		assert.GreaterOrEqual(t, len(response.Data), 0)
	})

	t.Run("popular tags with limit parameter", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Get popular tags with limit
		req, _ := http.NewRequest("GET", "/v1/tags/popular?limit=5", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[[]*models.Tag]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		// Should respect limit (can be fewer than 5 if not enough tags)
		assert.LessOrEqual(t, len(response.Data), 5)
	})
}

// TestTagHandler_GetRecentlyUsedTags tests the GetRecentlyUsedTags functionality
func TestTagHandler_GetRecentlyUsedTags(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful recently used tags retrieval", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Get recently used tags
		req, _ := http.NewRequest("GET", "/v1/tags/recent", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Debug: print response body if test fails
		if w.Code != http.StatusOK {
			t.Logf("Response status: %d, body: %s", w.Code, w.Body.String())
		}

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[[]*models.Tag]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		// Response should be valid (can be empty for new user)
		assert.GreaterOrEqual(t, len(response.Data), 0)
	})
}

// TestTagHandler_SearchTags tests the SearchTags functionality
func TestTagHandler_SearchTags(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful tag search", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, tagService := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create a tag with searchable name
		description := "Backend development tasks"
		_, err := tagService.CreateTag(context.Background(), &models.TagRequest{
			Name:        "backend-development",
			Description: &description,
			Color:       "#3B82F6",
		})
		require.NoError(t, err)

		// Search for tags
		req, _ := http.NewRequest("GET", "/v1/tags/search?q=backend", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[[]*models.Tag]
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		// Should find our tag
		assert.GreaterOrEqual(t, len(response.Data), 1)

		// Verify the found tag contains our search term
		found := false
		for _, tag := range response.Data {
			if tag.Name == "backend-development" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find the tag with 'backend' in the name")
	})

	t.Run("search with missing query parameter", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Search without query parameter
		req, _ := http.NewRequest("GET", "/v1/tags/search", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, w.Body.String(), "Search query parameter 'q' is required")
	})
}

// TestTagHandler_UpdateTag tests the UpdateTag functionality
func TestTagHandler_UpdateTag(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful tag update", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, tagService := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create a tag first
		tag := createTestTag(t, tagService)

		// Prepare update request
		updateDescription := "Updated backend development tasks"
		updateRequest := models.TagRequest{
			Name:        "updated-backend",
			Description: &updateDescription,
			Color:       "#EF4444",
		}

		body, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		// Update the tag
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/tags/%s", tag.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[*models.Tag]
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		assert.Equal(t, tag.ID, response.Data.ID)
		assert.Equal(t, updateRequest.Name, response.Data.Name)
		assert.Equal(t, *updateRequest.Description, *response.Data.Description)
		assert.Equal(t, updateRequest.Color, response.Data.Color)
	})

	t.Run("update non-existent tag", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Prepare update request
		updateDescription := "This tag does not exist"
		updateRequest := models.TagRequest{
			Name:        "non-existent-tag",
			Description: &updateDescription,
			Color:       "#EF4444",
		}

		body, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		// Try to update non-existent tag
		req, _ := http.NewRequest("PUT", "/v1/tags/00000000-0000-0000-0000-000000000000", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResponse gin.H
		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Failed to update tag")
	})
}

// TestTagHandler_DeleteTag tests the DeleteTag functionality
func TestTagHandler_DeleteTag(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful tag deletion", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, tagService := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create a tag first
		tag := createTestTag(t, tagService)

		// Delete the tag
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/tags/%s", tag.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["message"], "Tag deleted successfully")

		// Verify tag is deleted by trying to get it
		req, _ = http.NewRequest("GET", fmt.Sprintf("/v1/tags/%s", tag.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("delete non-existent tag", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Try to delete non-existent tag
		req, _ := http.NewRequest("DELETE", "/v1/tags/00000000-0000-0000-0000-000000000000", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var errorResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Failed to delete tag")
	})
}

// TestTagHandler_GetUserTagUsage tests the GetUserTagUsage functionality
func TestTagHandler_GetUserTagUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful user tag usage retrieval", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Get user tag usage
		req, _ := http.NewRequest("GET", "/v1/tags/usage", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[[]*models.TagUsage]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotNil(t, response.Data)

		// Response should be valid (can be empty for new user)
		assert.GreaterOrEqual(t, len(response.Data), 0)
	})
}

// TestTagHandler_ErrorHandling tests error handling scenarios
func TestTagHandler_ErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("unauthorized create tag request", func(t *testing.T) {
		// Setup test environment
		router, _, _, _, _ := RouterWithServices(t)

		// Prepare tag creation request
		description := "This should fail"
		tagRequest := models.TagRequest{
			Name:        "unauthorized-tag",
			Description: &description,
			Color:       "#3B82F6",
		}

		body, err := json.Marshal(tagRequest)
		require.NoError(t, err)

		// Try to create tag without authentication
		req, _ := http.NewRequest("POST", "/v1/tags", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("unauthorized get tag request", func(t *testing.T) {
		// Setup test environment
		router, _, _, _, _ := RouterWithServices(t)

		// Try to get tag without authentication
		req, _ := http.NewRequest("GET", "/v1/tags/some-id", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid tag ID format", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Try to get tag with invalid ID format
		req, _ := http.NewRequest("GET", "/v1/tags/invalid-uuid", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var errorResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Tag not found")
	})
}
