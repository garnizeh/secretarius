package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/services"
	"github.com/gin-gonic/gin"
)

// LogEntryHandler handles HTTP requests for log entries
type LogEntryHandler struct {
	logEntryService *services.LogEntryService
}

// NewLogEntryHandler creates a new LogEntryHandler instance
func NewLogEntryHandler(logEntryService *services.LogEntryService) *LogEntryHandler {
	return &LogEntryHandler{
		logEntryService: logEntryService,
	}
}

// CreateLogEntry handles POST /v1/logs
func (h *LogEntryHandler) CreateLogEntry(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.LogEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	logEntry, err := h.logEntryService.CreateLogEntry(c.Request.Context(), userID, &req)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Failed to create log entry", err.Error())
		return
	}

	RespondWithSuccess(c, http.StatusCreated, logEntry)
}

// GetLogEntry handles GET /v1/logs/:id
func (h *LogEntryHandler) GetLogEntry(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		RespondWithError(c, 401, "Unauthorized")
		return
	}

	logEntryID := c.Param("id")
	if logEntryID == "" {
		RespondWithError(c, 400, "Log entry ID is required")
		return
	}

	logEntry, err := h.logEntryService.GetLogEntry(c.Request.Context(), userID.(string), logEntryID)
	if err != nil {
		RespondWithError(c, 404, "Log entry not found")
		return
	}

	RespondWithSuccess(c, 200, logEntry, "Log entry retrieved successfully")
}

// GetLogEntries handles GET /v1/logs
func (h *LogEntryHandler) GetLogEntries(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Parse query parameters
	filters, err := h.parseLogEntryFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	logEntries, err := h.logEntryService.GetLogEntries(c.Request.Context(), userID.(string), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get log entries",
			"details": err.Error(),
		})
		return
	}

	// Apply pagination
	page, limit := h.parsePagination(c)
	paginatedEntries, pagination := h.paginate(logEntries, page, limit)

	c.JSON(http.StatusOK, gin.H{
		"data":       paginatedEntries,
		"pagination": pagination,
		"total":      len(logEntries),
	})
}

// UpdateLogEntry handles PUT /v1/logs/:id
func (h *LogEntryHandler) UpdateLogEntry(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		RespondWithError(c, 401, "Unauthorized")
		return
	}

	logEntryID := c.Param("id")
	if logEntryID == "" {
		RespondWithError(c, 400, "Log entry ID is required")
		return
	}

	var req models.LogEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, 400, "Invalid request format")
		return
	}

	logEntry, err := h.logEntryService.UpdateLogEntry(c.Request.Context(), userID.(string), logEntryID, &req)
	if err != nil {
		RespondWithError(c, 400, "Failed to update log entry")
		return
	}

	RespondWithSuccess(c, 200, logEntry, "Log entry updated successfully")
}

// DeleteLogEntry handles DELETE /v1/logs/:id
func (h *LogEntryHandler) DeleteLogEntry(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		RespondWithError(c, 401, "Unauthorized")
		return
	}

	logEntryID := c.Param("id")
	if logEntryID == "" {
		RespondWithError(c, 400, "Log entry ID is required")
		return
	}

	err := h.logEntryService.DeleteLogEntry(c.Request.Context(), userID.(string), logEntryID)
	if err != nil {
		RespondWithError(c, 500, "Failed to delete log entry")
		return
	}

	RespondWithSuccess(c, 200, nil, "Log entry deleted successfully")
}

// BulkCreateLogEntries handles POST /v1/logs/bulk
func (h *LogEntryHandler) BulkCreateLogEntries(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var req struct {
		Entries []models.LogEntryRequest `json:"entries" binding:"required,min=1,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	results := make([]any, len(req.Entries))
	successCount := 0
	errors := make([]string, 0)

	for i, entryReq := range req.Entries {
		logEntry, err := h.logEntryService.CreateLogEntry(c.Request.Context(), userID.(string), &entryReq)
		if err != nil {
			results[i] = gin.H{
				"error": err.Error(),
				"index": i,
			}
			errors = append(errors, err.Error())
		} else {
			results[i] = logEntry
			successCount++
		}
	}

	statusCode := http.StatusCreated
	if successCount == 0 {
		statusCode = http.StatusBadRequest
	} else if len(errors) > 0 {
		statusCode = http.StatusMultiStatus
	}

	c.JSON(statusCode, gin.H{
		"data": results,
		"summary": gin.H{
			"total":   len(req.Entries),
			"success": successCount,
			"errors":  len(errors),
		},
	})
}

// parseLogEntryFilters parses query parameters into LogEntryFilters
func (h *LogEntryHandler) parseLogEntryFilters(c *gin.Context) (*services.LogEntryFilters, error) {
	filters := &services.LogEntryFilters{}

	// Parse date range
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if err := models.ValidateDateFormat(startDateStr); err != nil {
			return nil, fmt.Errorf("invalid start_date format: %v", err)
		}
		startDate, _ := time.Parse("2006-01-02", startDateStr)
		filters.StartDate = startDate
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if err := models.ValidateDateFormat(endDateStr); err != nil {
			return nil, fmt.Errorf("invalid end_date format: %v", err)
		}
		endDate, _ := time.Parse("2006-01-02", endDateStr)
		filters.EndDate = endDate
	}

	// Parse project ID
	if projectID := c.Query("project_id"); projectID != "" {
		filters.ProjectID = &projectID
	}

	// Parse activity type
	if activityType := c.Query("type"); activityType != "" {
		if !models.ValidateActivityType(activityType) {
			return nil, fmt.Errorf("invalid activity type: %s", activityType)
		}
		aType := models.ActivityType(activityType)
		filters.Type = &aType
	}

	// Parse value rating
	if valueRating := c.Query("value_rating"); valueRating != "" {
		if !models.ValidateValueRating(valueRating) {
			return nil, fmt.Errorf("invalid value rating: %s", valueRating)
		}
		vRating := models.ValueRating(valueRating)
		filters.ValueRating = &vRating
	}

	// Parse impact level
	if impactLevel := c.Query("impact_level"); impactLevel != "" {
		if !models.ValidateImpactLevel(impactLevel) {
			return nil, fmt.Errorf("invalid impact level: %s", impactLevel)
		}
		iLevel := models.ImpactLevel(impactLevel)
		filters.ImpactLevel = &iLevel
	}

	// Parse tags
	if tagsStr := c.Query("tags"); tagsStr != "" {
		filters.Tags = strings.Split(tagsStr, ",")
	}

	return filters, nil
}

// parsePagination parses pagination parameters from query string
func (h *LogEntryHandler) parsePagination(c *gin.Context) (int, int) {
	page := 1
	limit := 50 // Default limit

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	return page, limit
}

// paginate applies pagination to log entries slice
func (h *LogEntryHandler) paginate(entries []*models.LogEntry, page, limit int) ([]*models.LogEntry, map[string]any) {
	total := len(entries)
	totalPages := (total + limit - 1) / limit

	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedEntries := entries[start:end]

	pagination := map[string]any{
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
		"has_next":    page < totalPages,
		"has_prev":    page > 1,
	}

	return paginatedEntries, pagination
}
