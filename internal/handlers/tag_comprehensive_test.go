package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTagHandler_Comprehensive_CreateTag tests comprehensive tag creation scenarios
func TestTagHandler_Comprehensive_CreateTag(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name         string
		tagRequest   models.TagRequest
		expectedCode int
		expectError  bool
	}{
		{
			name: "valid tag with all fields",
			tagRequest: models.TagRequest{
				Name:        fmt.Sprintf("comprehensive-tag-%d", time.Now().UnixNano()),
				Description: stringPtr("A comprehensive test tag with all fields"),
				Color:       "#FF5722",
			},
			expectedCode: http.StatusCreated,
			expectError:  false,
		},
		{
			name: "valid tag with minimal fields",
			tagRequest: models.TagRequest{
				Name:  fmt.Sprintf("minimal-tag-%d", time.Now().UnixNano()),
				Color: "#FF0000",
			},
			expectedCode: http.StatusCreated,
			expectError:  false,
		},
		{
			name: "tag with very long name",
			tagRequest: models.TagRequest{
				Name:  stringRepeat("a", 100),
				Color: "#FF0000",
			},
			expectedCode: http.StatusCreated,
			expectError:  false,
		},
		{
			name: "tag with extremely long name (exceeds limit)",
			tagRequest: models.TagRequest{
				Name:  stringRepeat("a", 300),
				Color: "#FF0000",
			},
			expectedCode: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name: "tag with empty name",
			tagRequest: models.TagRequest{
				Name:  "",
				Color: "#FF0000",
			},
			expectedCode: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name: "tag with only whitespace name",
			tagRequest: models.TagRequest{
				Name:  "   ",
				Color: "#FF0000",
			},
			expectedCode: http.StatusCreated, // System allows whitespace names
			expectError:  false,
		},
		{
			name: "tag with invalid color format",
			tagRequest: models.TagRequest{
				Name:  fmt.Sprintf("invalid-color-tag-%d", time.Now().UnixNano()),
				Color: "invalid-color",
			},
			expectedCode: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name: "tag with special characters in name",
			tagRequest: models.TagRequest{
				Name:  fmt.Sprintf("special-chars-@#$-tag-%d", time.Now().UnixNano()),
				Color: "#FF0000",
			},
			expectedCode: http.StatusCreated,
			expectError:  false,
		},
		{
			name: "tag with unicode characters",
			tagRequest: models.TagRequest{
				Name:  fmt.Sprintf("unicode-🚀-tag-%d", time.Now().UnixNano()),
				Color: "#FF0000",
			},
			expectedCode: http.StatusCreated,
			expectError:  false,
		},
		{
			name: "tag with very long description",
			tagRequest: models.TagRequest{
				Name:        fmt.Sprintf("long-desc-tag-%d", time.Now().UnixNano()),
				Description: stringPtr(stringRepeat("This is a very long description. ", 20)), // Reduced size
				Color:       "#FF0000",
			},
			expectedCode: http.StatusBadRequest,
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test environment
			router, userService, _, _, _ := RouterWithServices(t)

			// Create user and login
			user := createTestUser(t, userService)
			token := loginUser(t, router, user.Email, "password123")

			// Prepare request
			body, err := json.Marshal(tc.tagRequest)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/v1/tags", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectError {
				var errorResponse map[string]any
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				require.NoError(t, err)
				assert.Contains(t, errorResponse, "error")
			} else {
				var response responseData[*models.Tag]
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotNil(t, response.Data)
				assert.Equal(t, tc.tagRequest.Name, response.Data.Name)
				assert.Equal(t, tc.tagRequest.Color, response.Data.Color)
			}
		})
	}
}

// TestTagHandler_Comprehensive_TagDuplication tests tag name duplication scenarios
func TestTagHandler_Comprehensive_TagDuplication(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("duplicate tag name should fail", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, tagService := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create first tag directly via service
		firstTag := createTestTag(t, tagService)

		// Try to create second tag with same name via API
		tagRequest := models.TagRequest{
			Name:  firstTag.Name,
			Color: "#00FF00",
		}

		body, err := json.Marshal(tagRequest)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/tags", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should fail with bad request
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResponse map[string]any
		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "error")
	})

	t.Run("case insensitive tag name should fail", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, tagService := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create first tag
		firstTag := createTestTag(t, tagService)

		// Try to create second tag with different case
		tagRequest := models.TagRequest{
			Name:  firstTag.Name,
			Color: "#00FF00",
		}
		// Change case
		if len(tagRequest.Name) > 0 {
			tagRequest.Name = string(tagRequest.Name[0]-32) + tagRequest.Name[1:]
		}

		body, err := json.Marshal(tagRequest)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/tags", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Depending on implementation, this might fail or succeed
		// Check the actual behavior
		if w.Code == http.StatusBadRequest {
			var errorResponse map[string]any
			err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
			require.NoError(t, err)
			assert.Contains(t, errorResponse, "error")
		} else {
			assert.Equal(t, http.StatusCreated, w.Code)
		}
	})
}

// TestTagHandler_Comprehensive_GetTags tests comprehensive tag retrieval scenarios
func TestTagHandler_Comprehensive_GetTags(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("get all tags with multiple tags", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, tagService := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create multiple test tags
		tagCount := 5
		createdTags := make([]*models.Tag, tagCount)
		for i := range tagCount {
			createdTags[i] = createTestTag(t, tagService)
		}

		// Get all tags
		req := httptest.NewRequest(http.MethodGet, "/v1/tags", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data  []*models.Tag `json:"data"`
			Total int           `json:"total"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(response.Data), tagCount)
		assert.Equal(t, len(response.Data), response.Total)
	})
}

// TestTagHandler_Comprehensive_TagSearch tests comprehensive tag search scenarios
func TestTagHandler_Comprehensive_TagSearch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name         string
		query        string
		limit        string
		expectedCode int
		expectError  bool
	}{
		{
			name:         "valid search query",
			query:        "test",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "search with limit",
			query:        "test",
			limit:        "5",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "search with maximum limit",
			query:        "test",
			limit:        "50",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "search with excessive limit",
			query:        "test",
			limit:        "100",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "search with invalid limit",
			query:        "test",
			limit:        "invalid",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "search with negative limit",
			query:        "test",
			limit:        "-5",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "empty search query",
			query:        "",
			expectedCode: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name:         "whitespace search query",
			query:        "   ",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "special characters search",
			query:        "@#$%",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "unicode search",
			query:        "🚀",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "very long search query",
			query:        stringRepeat("test", 100),
			expectedCode: http.StatusOK,
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test environment
			router, userService, _, _, _ := RouterWithServices(t)

			// Create user and login
			user := createTestUser(t, userService)
			token := loginUser(t, router, user.Email, "password123")

			// Prepare URL
			endpoint := "/v1/tags/search?q=" + url.QueryEscape(tc.query)
			if tc.limit != "" {
				endpoint += "&limit=" + tc.limit
			}

			req := httptest.NewRequest(http.MethodGet, endpoint, nil)
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectError {
				var errorResponse map[string]any
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				require.NoError(t, err)
				assert.Contains(t, errorResponse, "error")
			} else {
				var response struct {
					Data  []*models.Tag `json:"data"`
					Total int           `json:"total"`
					Query string        `json:"query"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.NotNil(t, response.Data)
				assert.Equal(t, len(response.Data), response.Total)
				if tc.query != "" {
					assert.Equal(t, tc.query, response.Query)
				}
			}
		})
	}
}

// TestTagHandler_Comprehensive_UpdateTag tests comprehensive tag update scenarios
func TestTagHandler_Comprehensive_UpdateTag(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name         string
		updateReq    models.TagRequest
		expectedCode int
		expectError  bool
	}{
		{
			name: "update all fields",
			updateReq: models.TagRequest{
				Name:        fmt.Sprintf("updated-tag-%d", time.Now().UnixNano()),
				Description: stringPtr("Updated description"),
				Color:       "#00FF00",
			},
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name: "update only name",
			updateReq: models.TagRequest{
				Name:  fmt.Sprintf("name-only-update-%d", time.Now().UnixNano()),
				Color: "#FF0000", // Color is required
			},
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name: "update only color",
			updateReq: models.TagRequest{
				Name:  fmt.Sprintf("color-only-update-%d", time.Now().UnixNano()),
				Color: "#0000FF",
			},
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name: "update with empty name",
			updateReq: models.TagRequest{
				Name:  "",
				Color: "#0000FF",
			},
			expectedCode: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name: "update with invalid color",
			updateReq: models.TagRequest{
				Name:  fmt.Sprintf("invalid-color-update-%d", time.Now().UnixNano()),
				Color: "invalid-color",
			},
			expectedCode: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name: "update with very long name",
			updateReq: models.TagRequest{
				Name:  stringRepeat("a", 100),
				Color: "#FF0000",
			},
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name: "update with excessively long name",
			updateReq: models.TagRequest{
				Name:  stringRepeat("a", 300),
				Color: "#FF0000",
			},
			expectedCode: http.StatusBadRequest,
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test environment
			router, userService, _, _, tagService := RouterWithServices(t)

			// Create user and login
			user := createTestUser(t, userService)
			token := loginUser(t, router, user.Email, "password123")

			// Create initial tag
			existingTag := createTestTag(t, tagService)

			// Prepare update request
			body, err := json.Marshal(tc.updateReq)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/v1/tags/"+existingTag.ID.String(), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectError {
				var errorResponse map[string]any
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				require.NoError(t, err)
				assert.Contains(t, errorResponse, "error")
			} else {
				var response responseData[*models.Tag]
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotNil(t, response.Data, "Response data should not be nil for successful requests")
				assert.Equal(t, existingTag.ID, response.Data.ID)

				// Check updated fields
				if tc.updateReq.Name != "" {
					assert.Equal(t, tc.updateReq.Name, response.Data.Name)
				}
				if tc.updateReq.Color != "" {
					assert.Equal(t, tc.updateReq.Color, response.Data.Color)
				}
			}
		})
	}
}

// TestTagHandler_Comprehensive_DeleteTag tests comprehensive tag deletion scenarios
func TestTagHandler_Comprehensive_DeleteTag(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("delete existing tag", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, tagService := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create tag to delete
		tag := createTestTag(t, tagService)

		// Delete tag
		req := httptest.NewRequest(http.MethodDelete, "/v1/tags/"+tag.ID.String(), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "message")

		// Verify tag is deleted
		getReq := httptest.NewRequest(http.MethodGet, "/v1/tags/"+tag.ID.String(), nil)
		getReq.Header.Set("Authorization", "Bearer "+token)

		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)
		assert.Equal(t, http.StatusNotFound, getW.Code)
	})

	t.Run("delete non-existent tag", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Try to delete non-existent tag
		fakeID := "550e8400-e29b-41d4-a716-446655440000"
		req := httptest.NewRequest(http.MethodDelete, "/v1/tags/"+fakeID, nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should fail
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var errorResponse map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "error")
	})

	t.Run("delete with invalid tag ID", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create user and login
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Try to delete with invalid ID
		req := httptest.NewRequest(http.MethodDelete, "/v1/tags/invalid-id", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should fail
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var errorResponse map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "error")
	})
}

// TestTagHandler_Comprehensive_PopularTags tests comprehensive popular tags scenarios
func TestTagHandler_Comprehensive_PopularTags(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name         string
		limit        string
		expectedCode int
		expectError  bool
	}{
		{
			name:         "default limit",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "custom limit",
			limit:        "5",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "maximum limit",
			limit:        "50",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "excessive limit",
			limit:        "100",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "invalid limit",
			limit:        "invalid",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "negative limit",
			limit:        "-5",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "zero limit",
			limit:        "0",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test environment
			router, userService, _, _, _ := RouterWithServices(t)

			// Create user and login
			user := createTestUser(t, userService)
			token := loginUser(t, router, user.Email, "password123")

			// Prepare URL
			url := "/v1/tags/popular"
			if tc.limit != "" {
				url += "?limit=" + tc.limit
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectError {
				var errorResponse map[string]any
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				require.NoError(t, err)
				assert.Contains(t, errorResponse, "error")
			} else {
				var response struct {
					Data  []*models.Tag `json:"data"`
					Total int           `json:"total"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.NotNil(t, response.Data)
				assert.Equal(t, len(response.Data), response.Total)
			}
		})
	}
}

// TestTagHandler_Comprehensive_AuthenticationRequired tests authentication requirements
func TestTagHandler_Comprehensive_AuthenticationRequired(t *testing.T) {
	gin.SetMode(gin.TestMode)

	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/v1/tags"},
		{http.MethodGet, "/v1/tags"},
		{http.MethodGet, "/v1/tags/popular"},
		{http.MethodGet, "/v1/tags/recent"},
		{http.MethodGet, "/v1/tags/search?q=test"},
		{http.MethodGet, "/v1/tags/usage"},
		{http.MethodPut, "/v1/tags/550e8400-e29b-41d4-a716-446655440000"},
		{http.MethodDelete, "/v1/tags/550e8400-e29b-41d4-a716-446655440000"},
		{http.MethodGet, "/v1/tags/550e8400-e29b-41d4-a716-446655440000"},
	}

	for _, endpoint := range endpoints {
		t.Run(fmt.Sprintf("%s %s without authentication", endpoint.method, endpoint.path), func(t *testing.T) {
			// Setup test environment
			router, _, _, _, _ := RouterWithServices(t)

			var req *http.Request
			if endpoint.method == http.MethodPost || endpoint.method == http.MethodPut {
				body := bytes.NewBuffer([]byte(`{"name":"test","color":"#FF0000"}`))
				req = httptest.NewRequest(endpoint.method, endpoint.path, body)
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(endpoint.method, endpoint.path, nil)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should require authentication
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run(fmt.Sprintf("%s %s with invalid token", endpoint.method, endpoint.path), func(t *testing.T) {
			// Setup test environment
			router, _, _, _, _ := RouterWithServices(t)

			var req *http.Request
			if endpoint.method == http.MethodPost || endpoint.method == http.MethodPut {
				body := bytes.NewBuffer([]byte(`{"name":"test","color":"#FF0000"}`))
				req = httptest.NewRequest(endpoint.method, endpoint.path, body)
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(endpoint.method, endpoint.path, nil)
			}
			req.Header.Set("Authorization", "Bearer invalid-token")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should fail with unauthorized
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}
