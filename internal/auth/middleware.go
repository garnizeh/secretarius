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
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := bearerToken[1]
		claims, err := a.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		if claims.TokenType != models.TokenAccess {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token type for this endpoint",
			})
			c.Abort()
			return
		}

		// Set user ID in context for use in handlers
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func (a *AuthService) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			c.Next()
			return
		}

		token := bearerToken[1]
		claims, err := a.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		if claims.TokenType == models.TokenAccess {
			c.Set("user_id", claims.UserID)
		}

		c.Next()
	}
}
