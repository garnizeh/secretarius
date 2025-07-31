package middleware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Integration tests combining multiple middlewares
func TestMiddlewareIntegration(t *testing.T) {
	t.Run("multiple middlewares working together", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		_, router := gin.CreateTestContext(w)

		validator := middleware.NewValidationMiddleware()

		var buf bytes.Buffer
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		handler := slog.NewJSONHandler(&buf, opts)
		testLogger := &logging.Logger{Logger: slog.New(handler)}

		securityCfg := config.SecurityConfig{
			CORSAllowedOrigins:   []string{"https://example.com", "https://app.example.com"},
			CORSAllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			CORSAllowedHeaders:   []string{"Content-Type", "Authorization"},
			CORSAllowCredentials: true,
		}

		rateLimitCfg := config.RateLimitConfig{
			Enabled:           true,
			RequestsPerMinute: 10,
			RedisEnabled:      false, // Use no-op for testing
		}

		// Create rate limiter
		rateLimiter := middleware.NewRateLimiter(nil, rateLimitCfg, testLogger)

		// Stack multiple middlewares
		router.Use(middleware.RequestLogger(testLogger))
		router.Use(middleware.CORS(securityCfg))
		router.Use(rateLimiter.Middleware())

		router.GET("/test/:id",
			validator.ValidateUUIDParam("id"),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "success",
					"id":      c.Param("id"),
				})
			})

		req := httptest.NewRequest("GET", "/test/550e8400-e29b-41d4-a716-446655440000", nil)
		req.Header.Set("Origin", "https://example.com")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Body.String(), "success")
		assert.Contains(t, w.Body.String(), "550e8400-e29b-41d4-a716-446655440000")
	})

	t.Run("validation middleware failure in middleware stack", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		_, router := gin.CreateTestContext(w)

		validator := middleware.NewValidationMiddleware()

		var buf bytes.Buffer
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		handler := slog.NewJSONHandler(&buf, opts)
		testLogger := &logging.Logger{Logger: slog.New(handler)}

		securityCfg := config.SecurityConfig{
			CORSAllowedOrigins:   []string{"https://example.com", "https://app.example.com"},
			CORSAllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			CORSAllowedHeaders:   []string{"Content-Type", "Authorization"},
			CORSAllowCredentials: true,
		}

		rateLimitCfg := config.RateLimitConfig{
			Enabled:           true,
			RequestsPerMinute: 10,
			RedisEnabled:      false, // Use no-op for testing
		}

		// Create rate limiter
		rateLimiter := middleware.NewRateLimiter(nil, rateLimitCfg, testLogger)

		// Stack multiple middlewares
		router.Use(middleware.RequestLogger(testLogger))
		router.Use(middleware.CORS(securityCfg))
		router.Use(rateLimiter.Middleware())
		router.Use(validator.ValidateUUIDParam("id"))

		router.GET("/test/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Request with empty ID should fail validation
		req := httptest.NewRequest("GET", "/test/", nil)
		req.Header.Set("Origin", "https://example.com")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Body.String(), "id is required")
	})
}

// Edge case tests
func TestMiddlewareEdgeCases(t *testing.T) {
	t.Run("ValidationMiddleware with special characters in param name", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		_, router := gin.CreateTestContext(w)

		validator := middleware.NewValidationMiddleware()

		router.Use(validator.ValidateUUIDParam("user-id"))
		router.GET("/test/:user-id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		// Test with missing parameter
		req := httptest.NewRequest("GET", "/test/", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "user-id is required")
	})

	t.Run("CORS middleware with complex requests", func(t *testing.T) {
		cfg := config.SecurityConfig{
			CORSAllowedOrigins:   []string{"https://example.com", "https://app.example.com"},
			CORSAllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			CORSAllowedHeaders:   []string{"Content-Type", "Authorization"},
			CORSAllowCredentials: true,
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		_, router := gin.CreateTestContext(w)

		router.Use(middleware.CORS(cfg))
		router.POST("/api/data", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{"created": true})
		})

		req := httptest.NewRequest("POST", "/api/data", strings.NewReader(`{"test": "data"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer token123")
		req.Header.Set("Origin", "https://example.com")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	})
}
