package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standardized API response format
type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// DataWithPagination represents paginated data response
type DataWithPagination struct {
	Data       any                 `json:"data"`
	Pagination *PaginationResponse `json:"pagination"`
}

// RespondWithError sends an error response
func RespondWithError(c *gin.Context, statusCode int, message string, details ...string) {
	response := APIResponse{
		Success: false,
		Error:   message,
	}

	if len(details) > 0 {
		response.Message = details[0]
	}

	c.JSON(statusCode, response)
}

// RespondWithSuccess sends a success response
func RespondWithSuccess(c *gin.Context, statusCode int, data any, message ...string) {
	response := APIResponse{
		Success: true,
		Data:    data,
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	c.JSON(statusCode, response)
}

// RespondWithPagination sends a paginated response
func RespondWithPagination(c *gin.Context, data any, pagination *PaginationResponse) {
	response := APIResponse{
		Success: true,
		Data: DataWithPagination{
			Data:       data,
			Pagination: pagination,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetUserIDFromContext extracts user ID from gin context
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	userIDStr, ok := userID.(string)
	return userIDStr, ok
}

// TODO(rodrigo): plug it into the handlers (now is manually done in each handler)
// RequireUserID middleware that ensures user ID exists in context
func RequireUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := GetUserIDFromContext(c); !exists {
			RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
