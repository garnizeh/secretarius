package handlers

import (
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/auth"
	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/services"
	"github.com/garnizeh/englog/internal/testutils"
	"github.com/gin-gonic/gin"
)

// RouterWithServices creates a complete test router with all services returned
func RouterWithServices(t *testing.T) (*gin.Engine, *services.UserService, *services.ProjectService, *services.LogEntryService, *services.TagService) {
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
