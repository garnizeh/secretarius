package middleware_test

import (
	"bytes"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestLogger(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create logger that writes to buffer for testing
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(&buf, opts)
	testLogger := &logging.Logger{Logger: slog.New(handler)}

	// Create router with middleware
	router := gin.New()
	router.Use(middleware.RequestLogger(testLogger))

	// Add test endpoint
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	// Test request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Header().Get("X-Trace-ID"), "-") // UUID format

	// Check if structured log was written
	logOutput := buf.String()
	assert.Contains(t, logOutput, "HTTP Request")
	assert.Contains(t, logOutput, "test-agent")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "GET")
}
