package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles HTTP requests for health checks
type HealthHandler struct {
	startTime time.Time
}

// NewHealthHandler creates a new HealthHandler instance
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// HealthCheck handles GET /health
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	uptime := time.Since(h.startTime)

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    uptime.String(),
		"version":   "1.0.0",
	})
}

// ReadinessCheck handles GET /ready
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
