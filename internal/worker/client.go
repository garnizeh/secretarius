package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"time"

	"github.com/garnizeh/englog/internal/ai"
	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/logging"
	workerpb "github.com/garnizeh/englog/proto/worker"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Client represents the worker gRPC client with enhanced error handling
type Client struct {
	logger            *logging.Logger
	connectionManager *ConnectionManager
	aiService         *ai.OllamaService
	config            *config.Config
	workerID          string
	sessionToken      string
	stats             *WorkerStats
	taskManager       *TaskManager

	// Error handling and resilience
	retryConfig           *RetryConfig
	registrationBreaker   *CircuitBreaker
	taskProcessingBreaker *CircuitBreaker

	// Connection status
	connected      bool
	connectionMu   sync.RWMutex
	registrationMu sync.Mutex

	// Task processing
	taskSemaphore chan struct{}

	// Graceful shutdown
	ctx             context.Context
	cancel          context.CancelFunc
	shutdownTimeout time.Duration
}

// TaskManager manages active tasks
type TaskManager struct {
	activeTasks map[string]*ActiveTask
	mutex       sync.RWMutex
}

// ActiveTask represents a task being processed
type ActiveTask struct {
	ID        string
	Type      workerpb.TaskType
	Payload   string
	StartedAt time.Time
	Progress  int32
}

// WorkerStats tracks worker statistics
type WorkerStats struct {
	mutex          sync.RWMutex
	ActiveTasks    int32
	CompletedTasks int32
	FailedTasks    int32
	StartTime      time.Time
}

// NewClient creates a new worker client with enhanced error handling
func NewClient(logger *logging.Logger, connectionManager *ConnectionManager, aiService *ai.OllamaService, cfg *config.Config) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	maxTasks := cfg.Worker.MaxConcurrentTasks
	if maxTasks <= 0 {
		maxTasks = 5 // Default fallback
	}

	logger.Info("Creating new worker client",
		"component", "worker_client",
		"worker_id", cfg.Worker.WorkerID,
		"worker_name", cfg.Worker.WorkerName,
		"max_concurrent_tasks", maxTasks,
		"environment", cfg.Environment)

	return &Client{
		logger:            logger,
		connectionManager: connectionManager,
		aiService:         aiService,
		config:            cfg,
		workerID:          cfg.Worker.WorkerID,
		stats: &WorkerStats{
			StartTime: time.Now(),
		},
		taskManager: &TaskManager{
			activeTasks: make(map[string]*ActiveTask),
		},
		retryConfig:           DefaultRetryConfig(),
		registrationBreaker:   NewCircuitBreaker("worker_registration", 3, 2, 60*time.Second),
		taskProcessingBreaker: NewCircuitBreaker("task_processing", 5, 3, 30*time.Second),
		taskSemaphore:         make(chan struct{}, maxTasks),
		ctx:                   ctx,
		cancel:                cancel,
		shutdownTimeout:       30 * time.Second,
	}
}

// Start initializes the worker client and begins task processing with resilience
func (c *Client) Start(ctx context.Context) error {
	c.logger.Info("Starting worker client with enhanced error handling",
		"max_concurrent_tasks", cap(c.taskSemaphore),
		"shutdown_timeout", c.shutdownTimeout)

	// Connect to the server
	if err := c.connectionManager.Connect(ctx); err != nil {
		c.logger.Error("Failed to establish initial connection",
			"error", err,
			"target", c.connectionManager.target)
		return fmt.Errorf("failed to establish initial connection: %w", err)
	}

	c.logger.Info("Initial gRPC connection established successfully")

	// Register with API server using retry logic
	if err := c.registerWorkerWithRetry(ctx); err != nil {
		c.logger.Error("Failed to register worker after retries",
			"error", err)
		return fmt.Errorf("failed to register worker: %w", err)
	}

	c.logger.Info("Worker registration completed successfully")

	// Start background processes
	var wg sync.WaitGroup

	// Connection health monitoring routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.logger.Debug("Starting connection health monitoring routine")
		c.connectionHealthRoutine(ctx)
		c.logger.Debug("Connection health monitoring routine stopped")
	}()

	// Heartbeat routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.logger.Debug("Starting heartbeat routine")
		c.heartbeatRoutine(ctx)
		c.logger.Debug("Heartbeat routine stopped")
	}()

	// Task streaming routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.logger.Debug("Starting task streaming routine")
		c.taskStreamingRoutine(ctx)
		c.logger.Debug("Task streaming routine stopped")
	}()

	// Statistics reporting routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.logger.Debug("Starting statistics reporting routine")
		c.statisticsRoutine(ctx)
		c.logger.Debug("Statistics reporting routine stopped")
	}()

	c.logger.Info("All background routines started successfully")

	// Wait for context cancellation
	<-ctx.Done()

	c.logger.Info("Worker client shutting down...",
		"shutdown_timeout", c.shutdownTimeout,
		"reason", ctx.Err())
	c.cancel()

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), c.shutdownTimeout)
	defer shutdownCancel()

	// Wait for goroutines to finish or timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		c.logger.Info("Worker client stopped gracefully")
	case <-shutdownCtx.Done():
		c.logger.Warn("Worker client shutdown timeout reached")
	}

	// Close connection manager
	if err := c.connectionManager.Close(); err != nil {
		c.logger.Error("Failed to close connection manager", "error", err)
	} else {
		c.logger.Info("Connection manager closed successfully")
	}

	return nil
}

func (c *Client) registerWorkerWithRetry(ctx context.Context) error {
	c.registrationMu.Lock()
	defer c.registrationMu.Unlock()

	c.logger.Info("Starting worker registration process")

	return c.registrationBreaker.Execute(ctx, func() error {
		return RetryOperation(ctx, "worker_registration", c.retryConfig, func() error {
			return c.doRegisterWorker(ctx)
		})
	})
}

func (c *Client) doRegisterWorker(ctx context.Context) error {
	req := &workerpb.RegisterWorkerRequest{
		WorkerId:   c.workerID,
		WorkerName: c.config.Worker.WorkerName,
		Capabilities: []workerpb.WorkerCapability{
			workerpb.WorkerCapability_CAPABILITY_AI_INSIGHTS,
			workerpb.WorkerCapability_CAPABILITY_WEEKLY_REPORTS,
		},
		Version: c.config.Worker.Version,
		Metadata: map[string]string{
			"ai_model":    "llama3.2:3b",
			"max_tasks":   fmt.Sprintf("%d", c.config.Worker.MaxConcurrentTasks),
			"environment": c.config.Environment,
		},
	}

	c.logger.Debug("Sending worker registration request",
		"capabilities", req.Capabilities,
		"version", req.Version,
		"metadata", req.Metadata)

	return c.connectionManager.ExecuteWithRetry(ctx, "register_worker", func(client workerpb.APIWorkerServiceClient) error {
		resp, err := client.RegisterWorker(ctx, req)
		if err != nil {
			c.logger.Warn("Worker registration request failed", "error", err)
			return fmt.Errorf("registration request failed: %w", err)
		}

		if !resp.RegistrationSuccessful {
			c.logger.Error("Worker registration rejected by server",
				"message", resp.Message)
			return fmt.Errorf("worker registration failed: %s", resp.Message)
		}

		c.sessionToken = resp.SessionToken
		c.setConnected(true)

		c.logger.Info("Worker registered successfully",
			"session_token_length", len(c.sessionToken),
			"heartbeat_interval", resp.HeartbeatIntervalSeconds)
		return nil
	})
}

func (c *Client) connectionHealthRoutine(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	c.logger.Debug("Connection health monitoring started",
		"check_interval", "10s")

	for {
		select {
		case <-ctx.Done():
			c.logger.Debug("Connection health monitoring stopping due to context cancellation")
			return
		case <-ticker.C:
			if !c.connectionManager.IsConnected() {
				c.logger.Warn("Connection manager reports disconnected state",
					"connection_state", c.connectionManager.GetConnectionState())
				c.setConnected(false)

				// Attempt reconnection
				c.logger.Info("Attempting to reconnect to API server")
				if err := c.connectionManager.Reconnect(ctx); err != nil {
					c.logger.Error("Failed to reconnect",
						"error", err,
						"will_retry_next_cycle", true)
					continue
				}

				c.logger.Info("Reconnection successful, re-registering worker")
				// Re-register after successful reconnection
				if err := c.registerWorkerWithRetry(ctx); err != nil {
					c.logger.Error("Failed to re-register after reconnection",
						"error", err)
				} else {
					c.logger.Info("Worker re-registration completed successfully")
				}
			} else {
				c.logger.Debug("Connection health check passed")
			}
		}
	}
}

func (c *Client) heartbeatRoutine(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second) // Default heartbeat interval
	defer ticker.Stop()

	c.logger.Debug("Heartbeat routine started",
		"interval", "30s")

	for {
		select {
		case <-ctx.Done():
			c.logger.Debug("Heartbeat routine stopping due to context cancellation")
			return
		case <-ticker.C:
			c.logger.Debug("Sending heartbeat to API server")
			if err := c.sendHeartbeatWithRetry(ctx); err != nil {
				c.logger.Error("Heartbeat failed after retries",
					"error", err,
					"will_mark_disconnected", true)
				c.setConnected(false)
			} else {
				c.logger.Debug("Heartbeat sent successfully")
			}
		}
	}
}

func (c *Client) sendHeartbeatWithRetry(ctx context.Context) error {
	return RetryOperation(ctx, "heartbeat", c.retryConfig, func() error {
		return c.doSendHeartbeat(ctx)
	})
}

func (c *Client) doSendHeartbeat(ctx context.Context) error {
	c.stats.mutex.RLock()
	stats := &workerpb.WorkerStats{
		ActiveTasks:          c.stats.ActiveTasks,
		CompletedTasks:       c.stats.CompletedTasks,
		FailedTasks:          c.stats.FailedTasks,
		CpuUsage:             getCPUUsage(),
		MemoryUsage:          getMemoryUsage(),
		Uptime:               timestamppb.New(c.stats.StartTime),
		GrpcConnectionStatus: c.getConnectionStatus(),
	}
	c.stats.mutex.RUnlock()

	req := &workerpb.WorkerHeartbeatRequest{
		WorkerId:     c.workerID,
		SessionToken: c.sessionToken,
		Status:       c.getWorkerStatus(),
		Stats:        stats,
	}

	c.logger.Debug("Sending heartbeat with worker statistics",
		"active_tasks", stats.ActiveTasks,
		"completed_tasks", stats.CompletedTasks,
		"failed_tasks", stats.FailedTasks,
		"memory_usage_mb", stats.MemoryUsage,
		"connection_status", stats.GrpcConnectionStatus,
		"worker_status", req.Status)

	return c.connectionManager.ExecuteWithRetry(ctx, "heartbeat", func(client workerpb.APIWorkerServiceClient) error {
		resp, err := client.WorkerHeartbeat(ctx, req)
		if err != nil {
			c.logger.Warn("Heartbeat request failed", "error", err)
			return fmt.Errorf("heartbeat request failed: %w", err)
		}

		if !resp.ConnectionHealthy {
			c.logger.Warn("Server reports connection unhealthy",
				"message", resp.Message,
				"server_time", resp.ServerTime)
		} else {
			c.logger.Debug("Heartbeat acknowledged by server",
				"server_time", resp.ServerTime)
		}

		return nil
	})
}

func (c *Client) taskStreamingRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := c.streamTasksWithRetry(ctx); err != nil {
				slog.Error("Task streaming failed after retries", "error", err)
				c.setConnected(false)

				// Exponential backoff before retry
				backoffDelay := c.retryConfig.CalculateDelay(1)
				slog.Info("Backing off before reconnecting task stream", "delay", backoffDelay)

				select {
				case <-ctx.Done():
					return
				case <-time.After(backoffDelay):
					continue
				}
			}
		}
	}
}

func (c *Client) streamTasksWithRetry(ctx context.Context) error {
	return RetryOperation(ctx, "task_streaming", c.retryConfig, func() error {
		return c.doStreamTasks(ctx)
	})
}

func (c *Client) doStreamTasks(ctx context.Context) error {
	req := &workerpb.StreamTasksRequest{
		WorkerId:     c.workerID,
		SessionToken: c.sessionToken,
		Capabilities: []workerpb.WorkerCapability{
			workerpb.WorkerCapability_CAPABILITY_AI_INSIGHTS,
			workerpb.WorkerCapability_CAPABILITY_WEEKLY_REPORTS,
		},
	}

	var stream workerpb.APIWorkerService_StreamTasksClient

	err := c.connectionManager.ExecuteWithRetry(ctx, "create_task_stream", func(client workerpb.APIWorkerServiceClient) error {
		var err error
		stream, err = client.StreamTasks(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to create task stream: %w", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	c.setConnected(true)
	slog.Info("Task streaming started successfully")

	for {
		task, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("task stream receive error: %w", err)
		}

		// Process task with semaphore to limit concurrency
		select {
		case c.taskSemaphore <- struct{}{}:
			go func() {
				defer func() { <-c.taskSemaphore }()
				c.processTaskWithErrorHandling(ctx, task)
			}()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *Client) processTaskWithErrorHandling(ctx context.Context, task *workerpb.TaskRequest) {
	taskCtx, taskCancel := context.WithTimeout(ctx, 5*time.Minute) // Task timeout
	defer taskCancel()

	err := c.taskProcessingBreaker.Execute(taskCtx, func() error {
		return c.processTask(taskCtx, task)
	})

	if err != nil {
		slog.Error("Task processing failed through circuit breaker",
			"task_id", task.TaskId,
			"error", err)

		// Report failure
		c.reportTaskResultWithRetry(ctx, task.TaskId, workerpb.TaskStatus_TASK_STATUS_FAILED, "", err)
	}
}

func (c *Client) processTask(ctx context.Context, task *workerpb.TaskRequest) error {
	slog.Info("Processing task",
		"task_id", task.TaskId,
		"task_type", task.TaskType,
		"deadline", task.Deadline)

	// Add to active tasks
	activeTask := &ActiveTask{
		ID:        task.TaskId,
		Type:      task.TaskType,
		Payload:   task.Payload,
		StartedAt: time.Now(),
		Progress:  0,
	}

	c.taskManager.mutex.Lock()
	c.taskManager.activeTasks[task.TaskId] = activeTask
	c.taskManager.mutex.Unlock()

	// Update stats
	c.stats.mutex.Lock()
	c.stats.ActiveTasks++
	c.stats.mutex.Unlock()

	// Process with timeout handling
	startTime := time.Now()

	// Process based on task type with error handling
	var result string
	var processErr error

	switch task.TaskType {
	case workerpb.TaskType_TASK_TYPE_INSIGHT_GENERATION:
		result, processErr = c.processInsightTaskWithRetry(ctx, task)
	case workerpb.TaskType_TASK_TYPE_WEEKLY_REPORT:
		result, processErr = c.processWeeklyReportTaskWithRetry(ctx, task)
	default:
		processErr = fmt.Errorf("unsupported task type: %s", task.TaskType)
	}

	// Determine final status
	status := workerpb.TaskStatus_TASK_STATUS_COMPLETED
	if processErr != nil {
		status = workerpb.TaskStatus_TASK_STATUS_FAILED
		slog.Error("Task processing failed",
			"task_id", task.TaskId,
			"duration", time.Since(startTime),
			"error", processErr)
	} else {
		slog.Info("Task completed successfully",
			"task_id", task.TaskId,
			"duration", time.Since(startTime))
	}

	// Report result with retry
	c.reportTaskResultWithRetry(ctx, task.TaskId, status, result, processErr)

	// Clean up
	c.taskManager.mutex.Lock()
	delete(c.taskManager.activeTasks, task.TaskId)
	c.taskManager.mutex.Unlock()

	c.stats.mutex.Lock()
	c.stats.ActiveTasks--
	if processErr != nil {
		c.stats.FailedTasks++
	} else {
		c.stats.CompletedTasks++
	}
	c.stats.mutex.Unlock()

	return processErr
}

func (c *Client) processInsightTaskWithRetry(ctx context.Context, task *workerpb.TaskRequest) (string, error) {
	var result string

	err := RetryOperation(ctx, "insight_generation", c.retryConfig, func() error {
		var err error
		result, err = c.processInsightTask(ctx, task)
		return err
	})

	return result, err
}

func (c *Client) processWeeklyReportTaskWithRetry(ctx context.Context, task *workerpb.TaskRequest) (string, error) {
	var result string

	err := RetryOperation(ctx, "weekly_report_generation", c.retryConfig, func() error {
		var err error
		result, err = c.processWeeklyReportTask(ctx, task)
		return err
	})

	return result, err
}

func (c *Client) processInsightTask(ctx context.Context, task *workerpb.TaskRequest) (string, error) {
	var insightReq ai.InsightRequest
	if err := json.Unmarshal([]byte(task.Payload), &insightReq); err != nil {
		return "", fmt.Errorf("failed to unmarshal insight request: %w", err)
	}

	// Validate the insight request
	if err := c.aiService.ValidateInsightRequest(&insightReq); err != nil {
		return "", fmt.Errorf("invalid insight request: %w", err)
	}

	// Update progress
	c.updateTaskProgress(ctx, task.TaskId, 25, "Starting AI insight generation")

	// Use the enhanced context-aware insight generation
	insight, err := c.aiService.GenerateInsightWithContext(ctx, &insightReq)
	if err != nil {
		return "", fmt.Errorf("AI insight generation failed: %w", err)
	}

	// Update progress
	c.updateTaskProgress(ctx, task.TaskId, 100, "Insight generation completed")

	result, err := json.Marshal(insight)
	if err != nil {
		return "", fmt.Errorf("failed to marshal insight result: %w", err)
	}

	return string(result), nil
}

func (c *Client) processWeeklyReportTask(ctx context.Context, task *workerpb.TaskRequest) (string, error) {
	var reportReq ai.WeeklyReportRequest
	if err := json.Unmarshal([]byte(task.Payload), &reportReq); err != nil {
		return "", fmt.Errorf("failed to unmarshal weekly report request: %w", err)
	}

	// Update progress
	c.updateTaskProgress(ctx, task.TaskId, 25, "Starting weekly report generation")

	report, err := c.aiService.GenerateWeeklyReport(ctx, reportReq.UserID, reportReq.WeekStart, reportReq.WeekEnd)
	if err != nil {
		return "", fmt.Errorf("weekly report generation failed: %w", err)
	}

	// Update progress
	c.updateTaskProgress(ctx, task.TaskId, 100, "Weekly report completed")

	result, err := json.Marshal(report)
	if err != nil {
		return "", fmt.Errorf("failed to marshal report result: %w", err)
	}

	return string(result), nil
}

func (c *Client) updateTaskProgress(ctx context.Context, taskID string, progress int32, message string) {
	err := RetryOperation(ctx, "update_task_progress", c.retryConfig, func() error {
		return c.doUpdateTaskProgress(ctx, taskID, progress, message)
	})

	if err != nil {
		slog.Error("Failed to update task progress after retries",
			"task_id", taskID,
			"error", err)
	}
}

func (c *Client) doUpdateTaskProgress(ctx context.Context, taskID string, progress int32, message string) error {
	req := &workerpb.TaskProgressRequest{
		TaskId:          taskID,
		WorkerId:        c.workerID,
		ProgressPercent: progress,
		StatusMessage:   message,
		UpdatedAt:       timestamppb.Now(),
	}

	return c.connectionManager.ExecuteWithRetry(ctx, "update_progress", func(client workerpb.APIWorkerServiceClient) error {
		_, err := client.UpdateTaskProgress(ctx, req)
		if err != nil {
			return fmt.Errorf("update task progress failed: %w", err)
		}
		return nil
	})
}

func (c *Client) reportTaskResultWithRetry(ctx context.Context, taskID string, status workerpb.TaskStatus, result string, taskErr error) {
	err := RetryOperation(ctx, "report_task_result", c.retryConfig, func() error {
		return c.doReportTaskResult(ctx, taskID, status, result, taskErr)
	})

	if err != nil {
		slog.Error("Failed to report task result after retries",
			"task_id", taskID,
			"status", status,
			"error", err)
	}
}

func (c *Client) doReportTaskResult(ctx context.Context, taskID string, status workerpb.TaskStatus, result string, taskErr error) error {
	errorMessage := ""
	if taskErr != nil {
		errorMessage = taskErr.Error()
	}

	req := &workerpb.TaskResultRequest{
		TaskId:       taskID,
		WorkerId:     c.workerID,
		Status:       status,
		Result:       result,
		ErrorMessage: errorMessage,
		StartedAt:    timestamppb.Now(), // Should use actual start time
		CompletedAt:  timestamppb.Now(),
	}

	return c.connectionManager.ExecuteWithRetry(ctx, "report_result", func(client workerpb.APIWorkerServiceClient) error {
		_, err := client.ReportTaskResult(ctx, req)
		if err != nil {
			return fmt.Errorf("report task result failed: %w", err)
		}
		return nil
	})
}

func (c *Client) statisticsRoutine(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute) // Report stats every 5 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.logStatistics()
		}
	}
}

func (c *Client) logStatistics() {
	c.stats.mutex.RLock()
	stats := map[string]any{
		"worker_id":       c.workerID,
		"active_tasks":    c.stats.ActiveTasks,
		"completed_tasks": c.stats.CompletedTasks,
		"failed_tasks":    c.stats.FailedTasks,
		"uptime":          time.Since(c.stats.StartTime),
		"memory_usage_mb": getMemoryUsage(),
	}
	c.stats.mutex.RUnlock()

	// Add connection stats
	connectionStats := c.connectionManager.GetStats()
	for k, v := range connectionStats {
		stats["connection_"+k] = v
	}

	// Add circuit breaker stats
	registrationStats := c.registrationBreaker.GetStats()
	for k, v := range registrationStats {
		stats["registration_breaker_"+k] = v
	}

	taskProcessingStats := c.taskProcessingBreaker.GetStats()
	for k, v := range taskProcessingStats {
		stats["task_processing_breaker_"+k] = v
	}

	slog.Info("Worker statistics",
		"stats", stats)
}

func (c *Client) setConnected(connected bool) {
	c.connectionMu.Lock()
	defer c.connectionMu.Unlock()
	c.connected = connected
}

func (c *Client) isConnected() bool {
	c.connectionMu.RLock()
	defer c.connectionMu.RUnlock()
	return c.connected
}

func (c *Client) getConnectionStatus() string {
	if c.isConnected() {
		return "connected"
	}
	return "disconnected"
}

func (c *Client) getWorkerStatus() workerpb.WorkerStatus {
	c.stats.mutex.RLock()
	defer c.stats.mutex.RUnlock()

	if !c.isConnected() {
		return workerpb.WorkerStatus_WORKER_STATUS_ERROR
	}

	if c.stats.ActiveTasks > 0 {
		return workerpb.WorkerStatus_WORKER_STATUS_BUSY
	}

	return workerpb.WorkerStatus_WORKER_STATUS_IDLE
}

// Health check methods
func (c *Client) IsHealthy() bool {
	// Check gRPC connection and AI service
	return c.isConnected() && c.aiService.HealthCheck(context.Background()) == nil
}

func (c *Client) IsReady() bool {
	// Worker is ready if connected and not overloaded
	c.stats.mutex.RLock()
	defer c.stats.mutex.RUnlock()

	maxTasks := int32(c.config.Worker.MaxConcurrentTasks)
	return c.isConnected() && c.stats.ActiveTasks < maxTasks
}

func getCPUUsage() float32 {
	// Simplified CPU usage calculation
	// In production, use proper CPU monitoring
	return 0.0
}

func getMemoryUsage() float32 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float32(m.Alloc) / (1024 * 1024) // MB
}
