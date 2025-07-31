package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garnizeh/englog/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAnalyticsHandler_GetProductivityMetrics tests the GetProductivityMetrics endpoint
func TestAnalyticsHandler_GetProductivityMetrics(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedData   func(t *testing.T, response responseData[*models.ProductivityMetrics])
	}{
		{
			name:           "successful productivity metrics retrieval",
			path:           "/v1/analytics/productivity",
			expectedStatus: http.StatusOK,
			expectedData: func(t *testing.T, response responseData[*models.ProductivityMetrics]) {
				require.NotNil(t, response.Data)
				assert.GreaterOrEqual(t, response.Data.TotalActivities, 0)
				assert.GreaterOrEqual(t, response.Data.TotalMinutes, 0)
				assert.GreaterOrEqual(t, response.Data.ProjectsWorked, 0)
				assert.NotNil(t, response.Data.ActivityBreakdown)
				assert.NotNil(t, response.Data.ValueDistribution)
			},
		},
		{
			name:           "productivity metrics with date range",
			path:           "/v1/analytics/productivity?start_date=2024-01-01&end_date=2024-12-31",
			expectedStatus: http.StatusOK,
			expectedData: func(t *testing.T, response responseData[*models.ProductivityMetrics]) {
				require.NotNil(t, response.Data)
				assert.GreaterOrEqual(t, response.Data.TotalActivities, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			router, userService, _, _, _ := setupTestRouterWithServices(t)
			user := createTestUser(t, userService)
			token := loginUser(t, router, user.Email, "password123")

			// Make request
			req, _ := http.NewRequest("GET", tt.path, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK && tt.expectedData != nil {
				var response responseData[*models.ProductivityMetrics]
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				tt.expectedData(t, response)
			}
		})
	}
}

// TestAnalyticsHandler_GetActivitySummary tests the GetActivitySummary endpoint
func TestAnalyticsHandler_GetActivitySummary(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedData   func(t *testing.T, response responseData[[]*models.ActivitySummary])
	}{
		{
			name:           "successful activity summary retrieval",
			path:           "/v1/analytics/summary",
			expectedStatus: http.StatusOK,
			expectedData: func(t *testing.T, response responseData[[]*models.ActivitySummary]) {
				require.NotNil(t, response.Data)
				// Data can be empty for new users, so just check it's not nil
				assert.IsType(t, []*models.ActivitySummary{}, response.Data)
			},
		},
		{
			name:           "activity summary with date range",
			path:           "/v1/analytics/summary?start_date=2024-01-01&end_date=2024-12-31",
			expectedStatus: http.StatusOK,
			expectedData: func(t *testing.T, response responseData[[]*models.ActivitySummary]) {
				require.NotNil(t, response.Data)
				assert.IsType(t, []*models.ActivitySummary{}, response.Data)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			router, userService, _, _, _ := setupTestRouterWithServices(t)
			user := createTestUser(t, userService)
			token := loginUser(t, router, user.Email, "password123")

			// Make request
			req, _ := http.NewRequest("GET", tt.path, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK && tt.expectedData != nil {
				var response responseData[[]*models.ActivitySummary]
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				tt.expectedData(t, response)
			}
		})
	}
}

// TestAnalyticsHandler_ErrorHandling tests error scenarios
func TestAnalyticsHandler_ErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		token          string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "unauthorized productivity metrics request",
			path:           "/v1/analytics/productivity",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Unauthorized",
		},
		{
			name:           "unauthorized activity summary request",
			path:           "/v1/analytics/summary",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Unauthorized",
		},
		{
			name:           "invalid date format in productivity metrics",
			path:           "/v1/analytics/productivity?start_date=invalid-date",
			token:          "valid-token",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid date range",
		},
		{
			name:           "invalid date format in activity summary",
			path:           "/v1/analytics/summary?end_date=invalid-date",
			token:          "valid-token",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid date range",
		},
		{
			name:           "end date before start date in productivity metrics",
			path:           "/v1/analytics/productivity?start_date=2024-12-31&end_date=2024-01-01",
			token:          "valid-token",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid date range",
		},
		{
			name:           "end date before start date in activity summary",
			path:           "/v1/analytics/summary?start_date=2024-12-31&end_date=2024-01-01",
			token:          "valid-token",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid date range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			router, userService, _, _, _ := setupTestRouterWithServices(t)

			// Setup user and token for valid token tests
			if tt.token == "valid-token" {
				user := createTestUser(t, userService)
				tt.token = loginUser(t, router, user.Email, "password123")
			}

			// Make request
			req, _ := http.NewRequest("GET", tt.path, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Contains(t, response["error"], tt.expectedError)
		})
	}
}
