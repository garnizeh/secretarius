package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garnizeh/englog/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidationMiddleware_ValidateUUIDParam(t *testing.T) {
	tests := []struct {
		name          string
		paramName     string
		paramValue    string
		expectAbort   bool
		expectedCode  int
		expectedError string
	}{
		{
			name:         "valid UUID parameter",
			paramName:    "id",
			paramValue:   "550e8400-e29b-41d4-a716-446655440000",
			expectAbort:  false,
			expectedCode: 0,
		},
		{
			name:          "missing parameter",
			paramName:     "id",
			paramValue:    "",
			expectAbort:   true,
			expectedCode:  400,
			expectedError: "id is required",
		},
		{
			name:         "valid UUID parameter (different param name)",
			paramName:    "project_id",
			paramValue:   "550e8400-e29b-41d4-a716-446655440000",
			expectAbort:  false,
			expectedCode: 0,
		},
		{
			name:          "invalid UUID format",
			paramName:     "id",
			paramValue:    "invalid-uuid",
			expectAbort:   true,
			expectedCode:  400,
			expectedError: "Invalid id format",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			_, router := gin.CreateTestContext(w)

			// Create validation middleware
			validator := middleware.NewValidationMiddleware()
			middleware := validator.ValidateUUIDParam(tc.paramName)

			// Set up route with middleware
			router.Use(middleware)
			router.GET("/test/:"+tc.paramName, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"success": true})
			})

			// Create request with parameter
			var url string
			if tc.paramValue != "" {
				url = "/test/" + tc.paramValue
			} else {
				url = "/test/"
			}

			req := httptest.NewRequest("GET", url, nil)
			router.ServeHTTP(w, req)

			if tc.expectAbort {
				assert.Equal(t, tc.expectedCode, w.Code)

				// Check error message
				if tc.expectedError != "" {
					assert.Contains(t, w.Body.String(), tc.expectedError)
				}
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
			}
		})
	}
}

func TestNewValidationMiddleware(t *testing.T) {
	mid := middleware.NewValidationMiddleware()
	assert.NotNil(t, mid)
	assert.IsType(t, &middleware.ValidationMiddleware{}, mid)
}
