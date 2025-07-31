package middleware

import (
	"context"
	"time"

	"github.com/garnizeh/englog/internal/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Define a custom type for context keys to avoid collisions
type contextKey string

const traceIDKey contextKey = "trace_id"

// RequestLogger creates a structured logging middleware with request context
func RequestLogger(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceID := uuid.New().String()

		// Add trace ID to context for request tracing
		ctx := context.WithValue(c.Request.Context(), traceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)

		// Add trace ID to response headers for client correlation
		c.Header("X-Trace-ID", traceID)

		// Process request
		c.Next()

		// Log request details with structured logging
		duration := time.Since(start)

		// Get user ID from context if available (set by auth middleware)
		var userID string
		if uid, exists := c.Get("user_id"); exists {
			userID = uid.(string)
		}

		// Log with structured data
		logger.WithContext(ctx).Info("HTTP Request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status_code", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"user_id", userID,
			"request_id", traceID,
			"response_size", c.Writer.Size(),
		)
	}
}
