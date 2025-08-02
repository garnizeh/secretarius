package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test setup helper
func setupTestGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		message      string
		details      []string
		expectedBody APIResponse
	}{
		{
			name:       "error without details",
			statusCode: http.StatusBadRequest,
			message:    "Invalid input",
			details:    nil,
			expectedBody: APIResponse{
				Success: false,
				Error:   "Invalid input",
			},
		},
		{
			name:       "error with details",
			statusCode: http.StatusInternalServerError,
			message:    "Database error",
			details:    []string{"Connection timeout"},
			expectedBody: APIResponse{
				Success: false,
				Error:   "Database error",
				Message: "Connection timeout",
			},
		},
		{
			name:       "error with multiple details (only first is used)",
			statusCode: http.StatusUnauthorized,
			message:    "Authentication failed",
			details:    []string{"Invalid token", "Expired session"},
			expectedBody: APIResponse{
				Success: false,
				Error:   "Authentication failed",
				Message: "Invalid token",
			},
		},
		{
			name:       "not found error",
			statusCode: http.StatusNotFound,
			message:    "Resource not found",
			details:    []string{},
			expectedBody: APIResponse{
				Success: false,
				Error:   "Resource not found",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, w := setupTestGinContext()

			RespondWithError(c, tc.statusCode, tc.message, tc.details...)

			assert.Equal(t, tc.statusCode, w.Code)
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedBody, response)
		})
	}
}

func TestRespondWithSuccess(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		data         any
		message      []string
		expectedBody APIResponse
	}{
		{
			name:       "success without message",
			statusCode: http.StatusOK,
			data:       map[string]string{"id": "123", "name": "test"},
			message:    nil,
			expectedBody: APIResponse{
				Success: true,
				Data:    map[string]any{"id": "123", "name": "test"},
			},
		},
		{
			name:       "success with message",
			statusCode: http.StatusCreated,
			data:       map[string]string{"id": "456"},
			message:    []string{"Resource created successfully"},
			expectedBody: APIResponse{
				Success: true,
				Data:    map[string]any{"id": "456"},
				Message: "Resource created successfully",
			},
		},
		{
			name:       "success with multiple messages (only first is used)",
			statusCode: http.StatusAccepted,
			data:       []string{"item1", "item2"},
			message:    []string{"Operation accepted", "Processing in background"},
			expectedBody: APIResponse{
				Success: true,
				Data:    []any{"item1", "item2"},
				Message: "Operation accepted",
			},
		},
		{
			name:       "success with nil data",
			statusCode: http.StatusOK,
			data:       nil,
			message:    []string{"Deleted successfully"},
			expectedBody: APIResponse{
				Success: true,
				Data:    nil,
				Message: "Deleted successfully",
			},
		},
		{
			name:       "success with complex data structure",
			statusCode: http.StatusOK,
			data: struct {
				ID   int      `json:"id"`
				Name string   `json:"name"`
				Tags []string `json:"tags"`
			}{
				ID:   1,
				Name: "Test Item",
				Tags: []string{"tag1", "tag2"},
			},
			message: nil,
			expectedBody: APIResponse{
				Success: true,
				Data: map[string]any{
					"id":   float64(1), // JSON unmarshaling converts int to float64
					"name": "Test Item",
					"tags": []any{"tag1", "tag2"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, w := setupTestGinContext()

			RespondWithSuccess(c, tc.statusCode, tc.data, tc.message...)

			assert.Equal(t, tc.statusCode, w.Code)
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedBody, response)
		})
	}
}

func TestRespondWithPagination(t *testing.T) {
	tests := []struct {
		name       string
		data       any
		pagination *PaginationResponse
	}{
		{
			name: "paginated list with data",
			data: []map[string]string{
				{"id": "1", "name": "Item 1"},
				{"id": "2", "name": "Item 2"},
			},
			pagination: &PaginationResponse{
				Page:       1,
				Limit:      10,
				Total:      25,
				TotalPages: 3,
				HasNext:    true,
				HasPrev:    false,
			},
		},
		{
			name: "empty paginated result",
			data: []string{},
			pagination: &PaginationResponse{
				Page:       1,
				Limit:      10,
				Total:      0,
				TotalPages: 0,
				HasNext:    false,
				HasPrev:    false,
			},
		},
		{
			name: "last page pagination",
			data: []string{"last-item"},
			pagination: &PaginationResponse{
				Page:       3,
				Limit:      10,
				Total:      21,
				TotalPages: 3,
				HasNext:    false,
				HasPrev:    true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, w := setupTestGinContext()

			RespondWithPagination(c, tc.data, tc.pagination)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Check that it's successful
			assert.True(t, response.Success)
			assert.Empty(t, response.Error)
			assert.Empty(t, response.Message)

			// Check that the data structure is correct
			dataWithPagination, ok := response.Data.(map[string]any)
			require.True(t, ok, "Data should be a map with data and pagination keys")

			// Check data field
			actualData, exists := dataWithPagination["data"]
			assert.True(t, exists, "data field should exist")

			// Check pagination field
			actualPagination, exists := dataWithPagination["pagination"]
			assert.True(t, exists, "pagination field should exist")

			// Verify pagination structure
			paginationMap, ok := actualPagination.(map[string]any)
			require.True(t, ok, "pagination should be a map")

			assert.Equal(t, float64(tc.pagination.Page), paginationMap["page"])
			assert.Equal(t, float64(tc.pagination.Limit), paginationMap["limit"])
			assert.Equal(t, float64(tc.pagination.Total), paginationMap["total"])
			assert.Equal(t, float64(tc.pagination.TotalPages), paginationMap["total_pages"])
			assert.Equal(t, tc.pagination.HasNext, paginationMap["has_next"])
			assert.Equal(t, tc.pagination.HasPrev, paginationMap["has_prev"])

			// For data comparison, we need to handle the JSON marshaling differences
			expectedDataBytes, _ := json.Marshal(tc.data)
			actualDataBytes, _ := json.Marshal(actualData)
			assert.JSONEq(t, string(expectedDataBytes), string(actualDataBytes))
		})
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		expectedUserID string
		expectedExists bool
	}{
		{
			name: "valid user ID exists",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", "user-123")
			},
			expectedUserID: "user-123",
			expectedExists: true,
		},
		{
			name: "user ID not set",
			setupContext: func(c *gin.Context) {
				// Don't set anything
			},
			expectedUserID: "",
			expectedExists: false,
		},
		{
			name: "user ID is not a string",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", 12345) // Set as int instead of string
			},
			expectedUserID: "",
			expectedExists: false,
		},
		{
			name: "user ID is empty string",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", "")
			},
			expectedUserID: "",
			expectedExists: true,
		},
		{
			name: "user ID is nil",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", nil)
			},
			expectedUserID: "",
			expectedExists: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := setupTestGinContext()
			tc.setupContext(c)

			userID, exists := GetUserIDFromContext(c)

			assert.Equal(t, tc.expectedUserID, userID)
			assert.Equal(t, tc.expectedExists, exists)
		})
	}
}

func TestRequireUserID(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		expectAbort    bool
		expectedStatus int
		expectedBody   *APIResponse
	}{
		{
			name: "valid user ID - allows request to continue",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", "user-123")
			},
			expectAbort:    false,
			expectedStatus: 0, // No response written by middleware
			expectedBody:   nil,
		},
		{
			name: "no user ID - aborts request",
			setupContext: func(c *gin.Context) {
				// Don't set user_id
			},
			expectAbort:    true,
			expectedStatus: http.StatusUnauthorized,
			expectedBody: &APIResponse{
				Success: false,
				Error:   "Unauthorized",
			},
		},
		{
			name: "invalid user ID type - aborts request",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", 12345) // Set as int instead of string
			},
			expectAbort:    true,
			expectedStatus: http.StatusUnauthorized,
			expectedBody: &APIResponse{
				Success: false,
				Error:   "Unauthorized",
			},
		},
		{
			name: "empty user ID - allows request to continue",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", "") // Empty string is still a valid string
			},
			expectAbort:    false,
			expectedStatus: 0, // No response written by middleware
			expectedBody:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, w := setupTestGinContext()
			tc.setupContext(c)

			// Execute the middleware
			middleware := RequireUserID()
			middleware(c)

			// Check if the request was aborted
			assert.Equal(t, tc.expectAbort, c.IsAborted())

			if tc.expectAbort {
				assert.Equal(t, tc.expectedStatus, w.Code)

				var response APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, *tc.expectedBody, response)
			} else {
				// For successful cases, the middleware should not write any response
				// The status code should still be 0 because the middleware calls c.Next()
				// but there's no actual handler to set the status
				assert.Equal(t, 200, w.Code) // Gin sets 200 by default when headers are written
			}
		})
	}
} // Test edge cases and integration scenarios
func TestResponseStructures(t *testing.T) {
	t.Run("APIResponse JSON serialization", func(t *testing.T) {
		response := APIResponse{
			Success: true,
			Data:    map[string]string{"key": "value"},
			Message: "Success message",
		}

		data, err := json.Marshal(response)
		require.NoError(t, err)

		var unmarshaled APIResponse
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, response.Success, unmarshaled.Success)
		assert.Equal(t, response.Message, unmarshaled.Message)
		// Data will be map[string]any after unmarshaling
		expectedData := map[string]any{"key": "value"}
		assert.Equal(t, expectedData, unmarshaled.Data)
	})

	t.Run("PaginationResponse JSON serialization", func(t *testing.T) {
		pagination := PaginationResponse{
			Page:       2,
			Limit:      20,
			Total:      100,
			TotalPages: 5,
			HasNext:    true,
			HasPrev:    true,
		}

		data, err := json.Marshal(pagination)
		require.NoError(t, err)

		var unmarshaled PaginationResponse
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, pagination, unmarshaled)
	})

	t.Run("DataWithPagination JSON serialization", func(t *testing.T) {
		dataWithPagination := DataWithPagination{
			Data: []string{"item1", "item2"},
			Pagination: &PaginationResponse{
				Page:       1,
				Limit:      10,
				Total:      2,
				TotalPages: 1,
				HasNext:    false,
				HasPrev:    false,
			},
		}

		data, err := json.Marshal(dataWithPagination)
		require.NoError(t, err)

		var unmarshaled DataWithPagination
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		// Data will be []any after unmarshaling
		expectedData := []any{"item1", "item2"}
		assert.Equal(t, expectedData, unmarshaled.Data)
		assert.Equal(t, dataWithPagination.Pagination, unmarshaled.Pagination)
	})
}

// Benchmark tests for performance validation
func BenchmarkRespondWithError(b *testing.B) {
	gin.SetMode(gin.TestMode)

	b.ResetTimer()
	for range b.N {
		c, _ := setupTestGinContext()
		RespondWithError(c, http.StatusBadRequest, "Test error", "Test details")
	}
}

func BenchmarkRespondWithSuccess(b *testing.B) {
	gin.SetMode(gin.TestMode)
	data := map[string]string{"id": "123", "name": "test"}

	b.ResetTimer()
	for range b.N {
		c, _ := setupTestGinContext()
		RespondWithSuccess(c, http.StatusOK, data, "Test message")
	}
}

func BenchmarkRespondWithPagination(b *testing.B) {
	gin.SetMode(gin.TestMode)
	data := []string{"item1", "item2", "item3"}
	pagination := &PaginationResponse{
		Page:       1,
		Limit:      10,
		Total:      3,
		TotalPages: 1,
		HasNext:    false,
		HasPrev:    false,
	}

	b.ResetTimer()
	for range b.N {
		c, _ := setupTestGinContext()
		RespondWithPagination(c, data, pagination)
	}
}

func BenchmarkGetUserIDFromContext(b *testing.B) {
	gin.SetMode(gin.TestMode)
	c, _ := setupTestGinContext()
	c.Set("user_id", "user-123")

	b.ResetTimer()
	for range b.N {
		GetUserIDFromContext(c)
	}
}
