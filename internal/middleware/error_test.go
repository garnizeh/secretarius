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

func TestErrorLogger(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(&buf, opts)
	testLogger := &logging.Logger{Logger: slog.New(handler)}

	// Create router with middleware
	router := gin.New()
	router.Use(middleware.ErrorLogger(testLogger))

	// Add test endpoint that generates an error
	router.GET("/error", func(c *gin.Context) {
		_ = c.Error(assert.AnError)
		c.JSON(500, gin.H{"error": "test error"})
	})

	// Test request
	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 500, w.Code)

	// Check if error was logged
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Request error occurred")
	assert.Contains(t, logOutput, "/error")
}
