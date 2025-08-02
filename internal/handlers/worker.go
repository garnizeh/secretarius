package handlers

import (
	"net/http"
	"time"

	"github.com/garnizeh/englog/internal/grpc"
	"github.com/gin-gonic/gin"
)

// WorkerHandlers provides HTTP endpoints for worker management
type WorkerHandlers struct {
	grpcManager *grpc.Manager
}

// NewWorkerHandlers creates a new WorkerHandlers instance
func NewWorkerHandlers(grpcManager *grpc.Manager) *WorkerHandlers {
	return &WorkerHandlers{
		grpcManager: grpcManager,
	}
}

// GetActiveWorkers returns information about active workers
func (h *WorkerHandlers) GetActiveWorkers(c *gin.Context) {
	ctx := c.Request.Context()
	workers := h.grpcManager.GetActiveWorkers(ctx)

	response := make([]gin.H, 0, len(workers))
	for _, worker := range workers {
		response = append(response, gin.H{
			"id":             worker.ID,
			"name":           worker.Name,
			"capabilities":   worker.Capabilities,
			"version":        worker.Version,
			"status":         worker.Status.String(),
			"last_heartbeat": worker.LastHeartbeat,
			"stats":          worker.Stats,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"workers": response,
		"count":   len(response),
	})
}

// RequestInsightGeneration queues an insight generation task
func (h *WorkerHandlers) RequestInsightGeneration(c *gin.Context) {
	var req struct {
		UserID      string   `json:"user_id" binding:"required"`
		EntryIDs    []string `json:"entry_ids" binding:"required"`
		InsightType string   `json:"insight_type" binding:"required"`
		Context     any      `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	taskID, err := h.grpcManager.QueueInsightGenerationTask(
		ctx,
		req.UserID,
		req.EntryIDs,
		req.InsightType,
		req.Context,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"task_id": taskID,
		"message": "Insight generation task queued successfully",
	})
}

// RequestWeeklyReport queues a weekly report generation task
func (h *WorkerHandlers) RequestWeeklyReport(c *gin.Context) {
	var req struct {
		UserID    string `json:"user_id" binding:"required"`
		WeekStart string `json:"week_start" binding:"required"`
		WeekEnd   string `json:"week_end" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	weekStart, err := time.Parse("2006-01-02", req.WeekStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week_start format. Use YYYY-MM-DD"})
		return
	}

	weekEnd, err := time.Parse("2006-01-02", req.WeekEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week_end format. Use YYYY-MM-DD"})
		return
	}

	ctx := c.Request.Context()
	taskID, err := h.grpcManager.QueueWeeklyReportTask(ctx, req.UserID, weekStart, weekEnd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"task_id": taskID,
		"message": "Weekly report generation task queued successfully",
	})
}

// GetTaskResult returns the result of a completed task
func (h *WorkerHandlers) GetTaskResult(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task_id is required"})
		return
	}

	ctx := c.Request.Context()
	result, exists := h.grpcManager.GetTaskResult(ctx, taskID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or not completed yet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task_id":      result.TaskID,
		"worker_id":    result.WorkerID,
		"status":       result.Status.String(),
		"result":       result.Result,
		"error":        result.ErrorMsg,
		"started_at":   result.StartedAt,
		"completed_at": result.CompletedAt,
	})
}

// HealthCheck provides health status of the worker system
func (h *WorkerHandlers) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()
	workers := h.grpcManager.GetActiveWorkers(ctx)

	healthStatus := "healthy"
	if len(workers) == 0 {
		healthStatus = "warning"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         healthStatus,
		"active_workers": len(workers),
		"timestamp":      time.Now(),
		"workers":        workers,
	})
}

// SetupWorkerRoutes adds worker-related routes to the router
func SetupWorkerRoutes(router *gin.RouterGroup, grpcManager *grpc.Manager) {
	workerHandlers := NewWorkerHandlers(grpcManager)

	// Worker management routes
	workers := router.Group("/workers")
	{
		workers.GET("", workerHandlers.GetActiveWorkers)
		workers.GET("/health", workerHandlers.HealthCheck)
	}

	// Task management routes
	tasks := router.Group("/tasks")
	{
		tasks.POST("/insights", workerHandlers.RequestInsightGeneration)
		tasks.POST("/reports", workerHandlers.RequestWeeklyReport)
		tasks.GET("/:task_id/result", workerHandlers.GetTaskResult)
	}
}
