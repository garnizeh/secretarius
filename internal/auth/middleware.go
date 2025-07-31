package auth

import (
	"net/http"
	"strings"

	"github.com/garnizeh/englog/internal/models"
	"github.com/gin-gonic/gin"
)

func (a *AuthService) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			a.logger.Warn("Authentication required - missing authorization header", "path", c.Request.URL.Path, "ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			a.logger.Warn("Authentication failed - invalid authorization header format", "path", c.Request.URL.Path, "ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		token := bearerToken[1]
		claims, err := a.ValidateToken(c.Request.Context(), token)
		if err != nil {
			a.logger.Warn("Authentication failed - invalid token", "error", err.Error(), "path", c.Request.URL.Path, "ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		if claims.TokenType != models.TokenAccess {
			a.logger.Warn("Authentication failed - invalid token type", "token_type", claims.TokenType, "user_id", claims.UserID, "path", c.Request.URL.Path, "ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		a.logger.Debug("Authentication successful", "user_id", claims.UserID, "path", c.Request.URL.Path, "ip", c.ClientIP())

		// Set user ID in context for use in handlers
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func (a *AuthService) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			a.logger.Debug("Optional auth - no authorization header provided", "path", c.Request.URL.Path, "ip", c.ClientIP())
			c.Next()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			a.logger.Debug("Optional auth - invalid authorization header format, continuing without auth", "path", c.Request.URL.Path, "ip", c.ClientIP())
			c.Next()
			return
		}

		token := bearerToken[1]
		claims, err := a.ValidateToken(c.Request.Context(), token)
		if err != nil {
			a.logger.Debug("Optional auth - invalid token, continuing without auth", "error", err.Error(), "path", c.Request.URL.Path, "ip", c.ClientIP())
			c.Next()
			return
		}

		if claims.TokenType == models.TokenAccess {
			a.logger.Debug("Optional auth - valid token found", "user_id", claims.UserID, "path", c.Request.URL.Path, "ip", c.ClientIP())
			c.Set("user_id", claims.UserID)
		} else {
			a.logger.Debug("Optional auth - non-access token type, continuing without auth", "token_type", claims.TokenType, "path", c.Request.URL.Path, "ip", c.ClientIP())
		}

		c.Next()
	}
}
