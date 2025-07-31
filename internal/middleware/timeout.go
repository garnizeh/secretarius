package middleware

import (
	"context"
	"net/http"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/gin-gonic/gin"
)

// RequestTimeout creates a middleware that enforces a timeout on individual requests
func RequestTimeout(cfg config.ServerConfig, logger *logging.Logger) gin.HandlerFunc {
	timeout := cfg.RequestTimeout

	return func(c *gin.Context) {
		// Skip timeout for health check endpoints
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/ready" {
			c.Next()
			return
		}

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace the request context with the timeout context
		c.Request = c.Request.WithContext(ctx)

		// Channel to capture if the request completed
		done := make(chan struct{})

		go func() {
			defer close(done)
			c.Next()
		}()

		select {
		case <-done:
			// Request completed successfully
			return
		case <-ctx.Done():
			// Request timed out
			logger.LogError(ctx, ctx.Err(), "Request timeout exceeded",
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"timeout", timeout,
				"client_ip", c.ClientIP(),
			)

			// Check if response was already written
			if !c.Writer.Written() {
				c.JSON(http.StatusRequestTimeout, gin.H{
					"error":   "Request timeout",
					"message": "The request took too long to process",
					"timeout": timeout.String(),
				})
			}
			c.Abort()
			return
		}
	}
}
