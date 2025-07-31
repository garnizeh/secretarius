package handlers_test

import (
	"testing"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/handlers"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestSetupRoutesWithRateLimit(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Security: config.SecurityConfig{
			CORSAllowedOrigins: []string{"*"},
			CORSAllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			CORSAllowedHeaders: []string{"Content-Type", "Authorization"},
		},
		RateLimit: config.RateLimitConfig{
			Enabled:           true,
			RequestsPerMinute: 60,
			RedisEnabled:      false, // Use fallback for test
		},
	}

	// Create test logger
	logger := logging.NewTestLogger()

	// Mock services (nil for test)
	router := handlers.SetupRoutes(
		cfg,
		logger,
		nil, // Redis client (nil for fallback)
		nil, // authService
		nil, // logEntryService
		nil, // projectService
		nil, // analyticsService
		nil, // tagService
		nil, // userService
		nil, // grpcManager
	)

	assert.NotNil(t, router)

	// Verify that basic routes are configured
	routes := router.Routes()
	assert.NotEmpty(t, routes)

	// Check for health endpoints
	var healthRouteFound bool
	for _, route := range routes {
		if route.Path == "/health" && route.Method == "GET" {
			healthRouteFound = true
			break
		}
	}
	assert.True(t, healthRouteFound, "Health route should be configured")
}
