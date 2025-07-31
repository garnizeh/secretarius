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

func TestRecoveryLogger(t *testing.T) {
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
	router.Use(middleware.RecoveryLogger(testLogger))

	// Add test endpoint that panics
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	// Test request
	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 500, w.Code)

	// Check if panic was logged
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Panic recovered")
	assert.Contains(t, logOutput, "test panic")
}
