package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAnalyticsHandler_GetProductivityMetrics_Comprehensive tests comprehensive scenarios for productivity metrics
func TestAnalyticsHandler_GetProductivityMetrics_Comprehensive(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		setupData      func(router *gin.Engine, userService *services.UserService, projectService *services.ProjectService, logEntryService *services.LogEntryService) string
		expectedStatus int
		expectedData   func(t *testing.T, response responseData[*models.ProductivityMetrics])
	}{
		{
			name: "productivity metrics with activity data",
			path: "/v1/analytics/productivity",
			setupData: func(router *gin.Engine, userService *services.UserService, projectService *services.ProjectService, logEntryService *services.LogEntryService) string {
				user := createTestUser(t, userService)
				token := loginUser(t, router, user.Email, "password123")

				// Create a project
				project := createTestProject(t, projectService, user.ID.String())

				// Create some log entries to generate analytics data
				createTestLogEntry(t, logEntryService, user.ID.String(), &project.ID)

				return token
			},
			expectedStatus: http.StatusOK,
			expectedData: func(t *testing.T, response responseData[*models.ProductivityMetrics]) {
				require.NotNil(t, response.Data)
				assert.GreaterOrEqual(t, response.Data.TotalActivities, 1)
				assert.GreaterOrEqual(t, response.Data.TotalMinutes, 50) // Should be around 60 minutes
				assert.GreaterOrEqual(t, response.Data.ProjectsWorked, 1)
				assert.NotNil(t, response.Data.ActivityBreakdown)
				assert.NotNil(t, response.Data.ValueDistribution)
			},
		},
		{
			name: "productivity metrics with empty date range",
			path: "/v1/analytics/productivity?start_date=2023-01-01&end_date=2023-01-31",
			setupData: func(router *gin.Engine, userService *services.UserService, projectService *services.ProjectService, logEntryService *services.LogEntryService) string {
				user := createTestUser(t, userService)
				return loginUser(t, router, user.Email, "password123")
			},
			expectedStatus: http.StatusOK,
			expectedData: func(t *testing.T, response responseData[*models.ProductivityMetrics]) {
				require.NotNil(t, response.Data)
				assert.Equal(t, 0, response.Data.TotalActivities)
				assert.Equal(t, 0, response.Data.TotalMinutes)
				assert.Equal(t, 0, response.Data.ProjectsWorked)
			},
		},
		{
			name: "productivity metrics with future date range",
			path: "/v1/analytics/productivity?start_date=2026-01-01&end_date=2026-12-31",
			setupData: func(router *gin.Engine, userService *services.UserService, projectService *services.ProjectService, logEntryService *services.LogEntryService) string {
				user := createTestUser(t, userService)
				return loginUser(t, router, user.Email, "password123")
			},
			expectedStatus: http.StatusOK,
			expectedData: func(t *testing.T, response responseData[*models.ProductivityMetrics]) {
				require.NotNil(t, response.Data)
				assert.Equal(t, 0, response.Data.TotalActivities)
				assert.Equal(t, 0, response.Data.TotalMinutes)
				assert.Equal(t, 0, response.Data.ProjectsWorked)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			router, userService, projectService, logEntryService, _ := setupTestRouterWithServices(t)
			token := tt.setupData(router, userService, projectService, logEntryService)

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

// TestAnalyticsHandler_GetActivitySummary_Comprehensive tests comprehensive scenarios for activity summary
func TestAnalyticsHandler_GetActivitySummary_Comprehensive(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		setupData      func(router *gin.Engine, userService *services.UserService, projectService *services.ProjectService, logEntryService *services.LogEntryService) string
		expectedStatus int
		expectedData   func(t *testing.T, response responseData[[]*models.ActivitySummary])
	}{
		{
			name: "activity summary with multiple activities",
			path: "/v1/analytics/summary",
			setupData: func(router *gin.Engine, userService *services.UserService, projectService *services.ProjectService, logEntryService *services.LogEntryService) string {
				user := createTestUser(t, userService)
				token := loginUser(t, router, user.Email, "password123")

				// Create projects
				project1 := createTestProject(t, projectService, user.ID.String())
				project2 := createTestProject(t, projectService, user.ID.String())

				// Create log entries for different activities
				createTestLogEntry(t, logEntryService, user.ID.String(), &project1.ID)
				createTestLogEntry(t, logEntryService, user.ID.String(), &project2.ID)

				return token
			},
			expectedStatus: http.StatusOK,
			expectedData: func(t *testing.T, response responseData[[]*models.ActivitySummary]) {
				require.NotNil(t, response.Data)
				// Data should contain at least one summary entry if activities were created
				if len(response.Data) > 0 {
					summary := response.Data[0]
					assert.GreaterOrEqual(t, summary.TotalMinutes, 50)
					assert.GreaterOrEqual(t, summary.Count, 1)
				}
			},
		},
		{
			name: "activity summary with date range filter",
			path: "/v1/analytics/summary?start_date=2024-01-01&end_date=2024-01-31",
			setupData: func(router *gin.Engine, userService *services.UserService, projectService *services.ProjectService, logEntryService *services.LogEntryService) string {
				user := createTestUser(t, userService)
				return loginUser(t, router, user.Email, "password123")
			},
			expectedStatus: http.StatusOK,
			expectedData: func(t *testing.T, response responseData[[]*models.ActivitySummary]) {
				require.NotNil(t, response.Data)
				// Should be empty for this specific date range in 2024
				assert.Equal(t, 0, len(response.Data))
			},
		},
		{
			name: "activity summary with no activities",
			path: "/v1/analytics/summary",
			setupData: func(router *gin.Engine, userService *services.UserService, projectService *services.ProjectService, logEntryService *services.LogEntryService) string {
				user := createTestUser(t, userService)
				return loginUser(t, router, user.Email, "password123")
			},
			expectedStatus: http.StatusOK,
			expectedData: func(t *testing.T, response responseData[[]*models.ActivitySummary]) {
				require.NotNil(t, response.Data)
				assert.Equal(t, 0, len(response.Data))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			router, userService, projectService, logEntryService, _ := setupTestRouterWithServices(t)
			token := tt.setupData(router, userService, projectService, logEntryService)

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

// TestAnalyticsHandler_DateParsing_Comprehensive tests comprehensive date parsing scenarios
func TestAnalyticsHandler_DateParsing_Comprehensive(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "various invalid date formats in productivity metrics",
			path:           "/v1/analytics/productivity?start_date=2024/01/01",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid date range",
		},
		{
			name:           "missing day in date format",
			path:           "/v1/analytics/productivity?start_date=2024-01",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid date range",
		},
		{
			name:           "invalid month in date",
			path:           "/v1/analytics/summary?start_date=2024-13-01",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid date range",
		},
		{
			name:           "invalid day in date",
			path:           "/v1/analytics/summary?end_date=2024-02-30",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid date range",
		},
		{
			name:           "malformed date with text",
			path:           "/v1/analytics/productivity?start_date=today",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid date range",
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

			var response map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Contains(t, response["error"], tt.expectedError)
		})
	}
}
