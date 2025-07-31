package handlers

import (
	"net/http"
	"strconv"

	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/services"
	"github.com/gin-gonic/gin"
)

// TagHandler handles HTTP requests for tags
type TagHandler struct {
	tagService *services.TagService
}

// NewTagHandler creates a new TagHandler instance
func NewTagHandler(tagService *services.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// CreateTag handles POST /v1/tags
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req models.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	tag, err := h.tagService.CreateTag(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create tag",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": tag,
	})
}

// GetTag handles GET /v1/tags/:id
func (h *TagHandler) GetTag(c *gin.Context) {
	tagID := c.Param("id")

	tag, err := h.tagService.GetTag(c.Request.Context(), tagID)
	if err != nil {
		// Always return 404 "Tag not found" for any error (invalid format or not found)
		RespondWithError(c, 404, "Tag not found")
		return
	}

	RespondWithSuccess(c, 200, tag, "Tag retrieved successfully")
} // GetTags handles GET /v1/tags
func (h *TagHandler) GetTags(c *gin.Context) {
	tags, err := h.tagService.GetAllTags(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get tags",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tags,
		"total": len(tags),
	})
}

// GetPopularTags handles GET /v1/tags/popular
func (h *TagHandler) GetPopularTags(c *gin.Context) {
	limit := int32(10) // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = int32(l)
		}
	}

	tags, err := h.tagService.GetPopularTags(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get popular tags",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tags,
		"total": len(tags),
	})
}

// GetRecentlyUsedTags handles GET /v1/tags/recent
func (h *TagHandler) GetRecentlyUsedTags(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	limit := int32(10) // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = int32(l)
		}
	}

	tags, err := h.tagService.GetRecentlyUsedTags(c.Request.Context(), userID.(string), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get recently used tags",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tags,
		"total": len(tags),
	})
}

// SearchTags handles GET /v1/tags/search
func (h *TagHandler) SearchTags(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query parameter 'q' is required",
		})
		return
	}

	limit := int32(20) // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = int32(l)
		}
	}

	tags, err := h.tagService.SearchTags(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search tags",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tags,
		"total": len(tags),
		"query": query,
	})
}

// UpdateTag handles PUT /v1/tags/:id
func (h *TagHandler) UpdateTag(c *gin.Context) {
	tagID := c.Param("id")

	var req models.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, 400, "Invalid request format")
		return
	}

	tag, err := h.tagService.UpdateTag(c.Request.Context(), tagID, &req)
	if err != nil {
		RespondWithError(c, 400, "Failed to update tag")
		return
	}

	RespondWithSuccess(c, 200, tag, "Tag updated successfully")
}

// DeleteTag handles DELETE /v1/tags/:id
func (h *TagHandler) DeleteTag(c *gin.Context) {
	tagID := c.Param("id")

	err := h.tagService.DeleteTag(c.Request.Context(), tagID)
	if err != nil {
		// Always return 500 "Failed to delete tag" for any error (invalid format or not found)
		RespondWithError(c, 500, "Failed to delete tag")
		return
	}

	RespondWithSuccess(c, 200, nil, "Tag deleted successfully")
}

// GetUserTagUsage handles GET /v1/tags/usage
func (h *TagHandler) GetUserTagUsage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	usage, err := h.tagService.GetUserTagUsage(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get tag usage",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  usage,
		"total": len(usage),
	})
}
