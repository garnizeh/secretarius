package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/services"
	"github.com/gin-gonic/gin"
)

// AnalyticsHandler handles HTTP requests for analytics
type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

// NewAnalyticsHandler creates a new AnalyticsHandler instance
func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetProductivityMetrics handles GET /v1/analytics/productivity
func (h *AnalyticsHandler) GetProductivityMetrics(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Parse date range parameters
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid date range",
			"details": err.Error(),
		})
		return
	}

	metrics, err := h.analyticsService.GetProductivityMetrics(c.Request.Context(), userID.(string), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get productivity metrics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": metrics,
		"period": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
	})
}

// GetActivitySummary handles GET /v1/analytics/summary
func (h *AnalyticsHandler) GetActivitySummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Parse date range parameters
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid date range",
			"details": err.Error(),
		})
		return
	}

	summary, err := h.analyticsService.GetActivitySummary(c.Request.Context(), userID.(string), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get activity summary",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": summary,
		"period": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
	})
}

// parseDateRange parses date range parameters from query string
func (h *AnalyticsHandler) parseDateRange(c *gin.Context) (time.Time, time.Time, error) {
	// Default to last 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if err := models.ValidateDateFormat(startDateStr); err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid start_date format: %v", err)
		}
		parsed, _ := time.Parse("2006-01-02", startDateStr)
		startDate = parsed
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if err := models.ValidateDateFormat(endDateStr); err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end_date format: %v", err)
		}
		parsed, _ := time.Parse("2006-01-02", endDateStr)
		endDate = parsed
	}

	if endDate.Before(startDate) {
		return time.Time{}, time.Time{}, fmt.Errorf("end_date must be after start_date")
	}

	return startDate, endDate, nil
}
