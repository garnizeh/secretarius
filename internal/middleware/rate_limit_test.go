package middleware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test logger
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(&buf, opts)
	testLogger := &logging.Logger{Logger: slog.New(handler)}

	t.Run("rate limiting disabled", func(t *testing.T) {
		cfg := config.RateLimitConfig{
			Enabled:           false,
			RequestsPerMinute: 5,
			RedisEnabled:      false,
		}

		rateLimiter := middleware.NewRateLimiter(nil, cfg, testLogger)
		middlewareFunc := rateLimiter.Middleware()

		router := gin.New()
		router.Use(middlewareFunc)
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "test"})
		})

		// Make multiple requests - should all pass since rate limiting is disabled
		for i := 0; i < 10; i++ {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("rate limiting enabled but redis disabled", func(t *testing.T) {
		cfg := config.RateLimitConfig{
			Enabled:           true,
			RequestsPerMinute: 5,
			RedisEnabled:      false, // Falls back to no-op
		}

		rateLimiter := middleware.NewRateLimiter(nil, cfg, testLogger)
		middlewareFunc := rateLimiter.Middleware()

		router := gin.New()
		router.Use(middlewareFunc)
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "test"})
		})

		// Should pass since Redis is disabled (fallback to no-op)
		for i := 0; i < 10; i++ {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})
}
