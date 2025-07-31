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

	"github.com/garnizeh/englog/internal/auth"
	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/services"
	"github.com/garnizeh/englog/internal/store/testutils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// stringRepeat creates a string by repeating the input string a specified number of times
func stringRepeat(s string, count int) string {
	result := ""
	for range count {
		result += s
	}
	return result
}

// setupTestRouterWithServices creates a complete test router with all services returned
func setupTestRouterWithServices(t *testing.T) (*gin.Engine, *services.UserService, *services.ProjectService, *services.LogEntryService, *services.TagService) {
	t.Helper()

	db := testutils.DB(t)
	testLogger := logging.NewTestLogger()
	authService := auth.NewAuthService(db, testLogger, "test-secret")
	projectService := services.NewProjectService(db, testLogger)
	userService := services.NewUserService(db, testLogger)
	logEntryService := services.NewLogEntryService(db, testLogger)
	analyticsService := services.NewAnalyticsService(db, testLogger)
	tagService := services.NewTagService(db, testLogger)

	// Create test configuration
	cfg := &config.Config{
		Server: config.ServerConfig{
			RequestTimeout: 30 * time.Second, // Set reasonable timeout for tests
		},
		Security: config.SecurityConfig{
			CORSAllowedOrigins:   []string{"*"},
			CORSAllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			CORSAllowedHeaders:   []string{"Content-Type", "Authorization"},
			CORSAllowCredentials: true,
			SecurityHeaders: config.SecurityHeadersConfig{
				ContentTypeOptions:      "nosniff",
				FrameOptions:            "DENY",
				XSSProtection:           "1; mode=block",
				StrictTransportSecurity: "max-age=31536000",
				ContentSecurityPolicy:   "default-src 'self'",
				ReferrerPolicy:          "strict-origin",
			},
		},
	}

	router := SetupRoutes(
		cfg,
		testLogger,
		nil, // No Redis client in tests
		authService,
		logEntryService,
		projectService,
		analyticsService,
		tagService,
		userService,
		nil, // No gRPC manager in tests
	)

	return router, userService, projectService, logEntryService, tagService
}

// responseData is a generic structure for API responses with data field
type responseData[T any] struct {
	Data T `json:"data"`
}

// createTestUser creates a test user and returns it
func createTestUser(t *testing.T, userService *services.UserService) *models.User {
	t.Helper()

	ctx := context.Background()
	email := fmt.Sprintf("test-%d@example.com", time.Now().UnixNano())

	user, err := userService.CreateUser(ctx, &models.UserRegistration{
		Email:     email,
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Timezone:  "UTC",
	})
	require.NoError(t, err)

	// Cleanup on test end
	t.Cleanup(func() {
		_ = userService.DeleteUser(ctx, user.ID.String())
	})

	return user
}

// loginUser logs in a user and returns the access token
func loginUser(t *testing.T, router *gin.Engine, email, password string) string {
	t.Helper()

	loginReq := map[string]string{
		"email":    email,
		"password": password,
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Tokens models.AuthTokens `json:"tokens"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Tokens.AccessToken
}

// createTestProject creates a test project and returns it
func createTestProject(t *testing.T, projectService *services.ProjectService, userID string) *models.Project {
	t.Helper()

	ctx := context.Background()
	projectReq := &models.ProjectRequest{
		Name:        fmt.Sprintf("Test Project %d", time.Now().UnixNano()),
		Description: stringPtr("A test project"),
		Color:       "#FF5733",
		Status:      models.ProjectActive,
		IsDefault:   false,
	}

	project, err := projectService.CreateProject(ctx, userID, projectReq)
	require.NoError(t, err)

	// Cleanup on test end
	t.Cleanup(func() {
		_ = projectService.DeleteProject(ctx, userID, project.ID.String())
	})

	return project
}

// createTestLogEntry creates a test log entry and returns it
func createTestLogEntry(t *testing.T, logEntryService *services.LogEntryService, userID string, projectID *uuid.UUID) *models.LogEntry {
	t.Helper()

	ctx := context.Background()
	now := time.Now()
	logReq := &models.LogEntryRequest{
		Title:       fmt.Sprintf("Test log entry %d", time.Now().UnixNano()),
		Description: stringPtr("Test log entry description"),
		Type:        models.ActivityDevelopment,
		ProjectID:   projectID,
		StartTime:   now.Add(-1 * time.Hour),
		EndTime:     now,
		ValueRating: models.ValueMedium,
		ImpactLevel: models.ImpactTeam,
		Tags:        []string{"test", "automation"},
	}

	logEntry, err := logEntryService.CreateLogEntry(ctx, userID, logReq)
	require.NoError(t, err)

	// Cleanup on test end
	t.Cleanup(func() {
		_ = logEntryService.DeleteLogEntry(ctx, userID, logEntry.ID.String())
	})

	return logEntry
}

// createTestTag creates a test tag and returns it
func createTestTag(t *testing.T, tagService *services.TagService) *models.Tag {
	t.Helper()

	ctx := context.Background()
	tagReq := &models.TagRequest{
		Name:        fmt.Sprintf("test-tag-%d", time.Now().UnixNano()),
		Color:       "#FF5733",
		Description: stringPtr("A test tag"),
	}

	tag, err := tagService.CreateTag(ctx, tagReq)
	require.NoError(t, err)

	// Cleanup on test end
	t.Cleanup(func() {
		_ = tagService.DeleteTag(ctx, tag.ID.String())
	})

	return tag
}
