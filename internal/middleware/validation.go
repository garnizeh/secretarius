package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ValidationMiddleware provides request validation utilities
type ValidationMiddleware struct{}

// NewValidationMiddleware creates a new validation middleware instance
func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{}
}

// ValidateUUIDParam validates that a URL parameter is a valid UUID
func (v *ValidationMiddleware) ValidateUUIDParam(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		paramValue := c.Param(paramName)
		if paramValue == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": paramName + " is required",
			})
			c.Abort()
			return
		}

		// Validate UUID format
		if _, err := uuid.Parse(paramValue); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid " + paramName + " format",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateUUID is a standalone function for UUID validation
func ValidateUUID(paramName string) gin.HandlerFunc {
	validator := NewValidationMiddleware()
	return validator.ValidateUUIDParam(paramName)
}
