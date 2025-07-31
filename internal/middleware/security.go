package middleware

import (
	"strings"

	"github.com/garnizeh/englog/internal/config"
	"github.com/gin-gonic/gin"
)

// SecurityHeaders creates a middleware for setting security headers
func SecurityHeaders(cfg config.SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := cfg.SecurityHeaders

		// Set security headers
		c.Header("X-Content-Type-Options", headers.ContentTypeOptions)
		c.Header("X-Frame-Options", headers.FrameOptions)
		c.Header("X-XSS-Protection", headers.XSSProtection)
		c.Header("Strict-Transport-Security", headers.StrictTransportSecurity)
		c.Header("Content-Security-Policy", headers.ContentSecurityPolicy)
		c.Header("Referrer-Policy", headers.ReferrerPolicy)

		c.Next()
	}
}

// CORS creates a middleware for handling CORS with environment-aware configuration
func CORS(cfg config.SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		if isOriginAllowed(origin, cfg.CORSAllowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// Set CORS headers
		c.Header("Access-Control-Allow-Methods", strings.Join(cfg.CORSAllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(cfg.CORSAllowedHeaders, ", "))

		if cfg.CORSAllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Expose trace ID header for client debugging
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Trace-ID")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// isOriginAllowed checks if the origin is in the allowed list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}
