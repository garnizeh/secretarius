package middleware

import (
	"github.com/garnizeh/englog/internal/logging"
	"github.com/gin-gonic/gin"
)

// RecoveryLogger creates a middleware for panic recovery with structured logging
func RecoveryLogger(logger *logging.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		// Get user ID from context if available
		var userID string
		if uid, exists := c.Get("user_id"); exists {
			userID = uid.(string)
		}

		logger.WithContext(c.Request.Context()).Error("Panic recovered",
			"panic", recovered,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"client_ip", c.ClientIP(),
			"user_id", userID,
		)

		c.AbortWithStatus(500)
	})
}
