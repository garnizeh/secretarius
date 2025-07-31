package middleware

import (
	"github.com/garnizeh/englog/internal/logging"
	"github.com/gin-gonic/gin"
)

// ErrorLogger creates a middleware for structured error logging
func ErrorLogger(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log any errors that occurred during request processing
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				// Get user ID from context if available
				var userID string
				if uid, exists := c.Get("user_id"); exists {
					userID = uid.(string)
				}

				logger.LogError(c.Request.Context(), err.Err, "Request error occurred",
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"status_code", c.Writer.Status(),
					"client_ip", c.ClientIP(),
					"user_id", userID,
					"error_type", err.Type,
					"error_meta", err.Meta,
				)
			}
		}
	}
}
