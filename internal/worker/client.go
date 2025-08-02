package worker

import (
	"context"
	"encoding/json"
	"fmt"
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

	// Configure logger with worker client component
	clientLogger := logger.WithServiceAndComponent("worker", "client")

	clientLogger.LogInfo(ctx, "Creating new worker client",
		logging.OperationField, "create_worker_client",
		"worker_id", cfg.Worker.WorkerID,
		"worker_name", cfg.Worker.WorkerName,
		"max_concurrent_tasks", maxTasks,
		"environment", cfg.Environment)

	return &Client{
		logger:            clientLogger,
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
		registrationBreaker:   NewCircuitBreaker(clientLogger, "worker_registration", 3, 2, 60*time.Second),
		taskProcessingBreaker: NewCircuitBreaker(clientLogger, "task_processing", 5, 3, 30*time.Second),
		taskSemaphore:         make(chan struct{}, maxTasks),
		ctx:                   ctx,
		cancel:                cancel,
		shutdownTimeout:       30 * time.Second,
	}
}

// Start initializes the worker client and begins task processing with resilience
func (c *Client) Start(ctx context.Context) error {
	c.logger.LogInfo(ctx, "Starting worker client with enhanced error handling",
		logging.OperationField, "start_worker_client",
		"max_concurrent_tasks", cap(c.taskSemaphore),
		"shutdown_timeout", c.shutdownTimeout)

	// Connect to the server
	if err := c.connectionManager.Connect(ctx); err != nil {
		c.logger.LogError(ctx, err, "Failed to establish initial connection",
			logging.OperationField, "initial_connection",
			"target", c.connectionManager.target)
		return fmt.Errorf("failed to establish initial connection: %w", err)
	}

	c.logger.LogInfo(ctx, "Initial gRPC connection established successfully",
		logging.OperationField, "initial_connection")

	// Register with API server using retry logic
	if err := c.registerWorkerWithRetry(ctx); err != nil {
		c.logger.LogError(ctx, err, "Failed to register worker after retries",
			logging.OperationField, "worker_registration")
		return fmt.Errorf("failed to register worker: %w", err)
	}

	c.logger.LogInfo(ctx, "Worker registration completed successfully",
		logging.OperationField, "worker_registration")

	// Start background processes
	var wg sync.WaitGroup

	// Connection health monitoring routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.logger.LogDebug(ctx, "Starting connection health monitoring routine",
			logging.OperationField, "start_health_monitoring")
		c.connectionHealthRoutine(ctx)
		c.logger.LogDebug(ctx, "Connection health monitoring routine stopped",
			logging.OperationField, "stop_health_monitoring")
	}()

	// Heartbeat routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.logger.LogDebug(ctx, "Starting heartbeat routine",
			logging.OperationField, "start_heartbeat")
		c.heartbeatRoutine(ctx)
		c.logger.LogDebug(ctx, "Heartbeat routine stopped",
			logging.OperationField, "stop_heartbeat")
	}()

	// Task streaming routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.logger.LogDebug(ctx, "Starting task streaming routine",
			logging.OperationField, "start_task_streaming")
		c.taskStreamingRoutine(ctx)
		c.logger.LogDebug(ctx, "Task streaming routine stopped",
			logging.OperationField, "stop_task_streaming")
	}()

	// Statistics reporting routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.logger.LogDebug(ctx, "Starting statistics reporting routine",
			logging.OperationField, "start_statistics")
		c.statisticsRoutine(ctx)
		c.logger.LogDebug(ctx, "Statistics reporting routine stopped",
			logging.OperationField, "stop_statistics")
	}()

	c.logger.LogInfo(ctx, "All background routines started successfully",
		logging.OperationField, "start_background_routines")

	// Wait for context cancellation
	<-ctx.Done()

	c.logger.LogInfo(ctx, "Worker client shutting down...",
		logging.OperationField, "shutdown_start",
		"shutdown_timeout", c.shutdownTimeout,
		logging.ReasonField, ctx.Err())
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
		c.logger.LogInfo(ctx, "Worker client stopped gracefully",
			logging.OperationField, "shutdown_complete")
	case <-shutdownCtx.Done():
		c.logger.LogWarn(ctx, "Worker client shutdown timeout reached",
			logging.OperationField, "shutdown_timeout")
	}

	// Close connection manager
	if err := c.connectionManager.Close(); err != nil {
		c.logger.LogError(ctx, err, "Failed to close connection manager",
			logging.OperationField, "close_connection_manager")
	} else {
		c.logger.LogInfo(ctx, "Connection manager closed successfully",
			logging.OperationField, "close_connection_manager")
	}

	return nil
}

func (c *Client) registerWorkerWithRetry(ctx context.Context) error {
	c.registrationMu.Lock()
	defer c.registrationMu.Unlock()

	c.logger.LogInfo(ctx, "Starting worker registration process",
		logging.OperationField, "register_worker_with_retry")

	return c.registrationBreaker.Execute(ctx, func() error {
		return RetryOperation(ctx, c.logger, "worker_registration", c.retryConfig, func() error {
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
			"ai_model":    "qwen2.5-coder:7b",
			"max_tasks":   fmt.Sprintf("%d", c.config.Worker.MaxConcurrentTasks),
			"environment": c.config.Environment,
		},
	}

	c.logger.LogDebug(ctx, "Sending worker registration request",
		logging.OperationField, "do_register_worker",
		"capabilities", req.Capabilities,
		"version", req.Version,
		"metadata", req.Metadata)

	return c.connectionManager.ExecuteWithRetry(ctx, "register_worker", func(client workerpb.APIWorkerServiceClient) error {
		resp, err := client.RegisterWorker(ctx, req)
		if err != nil {
			c.logger.LogWarn(ctx, "Worker registration request failed",
				logging.OperationField, "do_register_worker",
				logging.ErrorField, err)
			return fmt.Errorf("registration request failed: %w", err)
		}

		if !resp.RegistrationSuccessful {
			c.logger.LogError(ctx, fmt.Errorf("worker registration rejected: %s", resp.Message),
				"Worker registration rejected by server",
				logging.OperationField, "do_register_worker",
				"message", resp.Message)
			return fmt.Errorf("worker registration failed: %s", resp.Message)
		}

		c.sessionToken = resp.SessionToken
		c.setConnected(true)

		c.logger.LogInfo(ctx, "Worker registered successfully",
			logging.OperationField, "do_register_worker",
			"session_token_length", len(c.sessionToken),
			"heartbeat_interval", resp.HeartbeatIntervalSeconds)
		return nil
	})
}

func (c *Client) connectionHealthRoutine(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	c.logger.LogDebug(ctx, "Connection health monitoring started",
		logging.OperationField, "connection_health_routine",
		"check_interval", "10s")

	for {
		select {
		case <-ctx.Done():
			c.logger.LogDebug(ctx, "Connection health monitoring stopping due to context cancellation",
				logging.OperationField, "connection_health_routine")
			return
		case <-ticker.C:
			if !c.connectionManager.IsConnected() {
				c.logger.LogWarn(ctx, "Connection manager reports disconnected state",
					logging.OperationField, "connection_health_routine",
					"connection_state", c.connectionManager.GetConnectionState())
				c.setConnected(false)

				// Attempt reconnection
				c.logger.LogInfo(ctx, "Attempting to reconnect to API server",
					logging.OperationField, "connection_health_routine")
				if err := c.connectionManager.Reconnect(ctx); err != nil {
					c.logger.LogError(ctx, err, "Failed to reconnect",
						logging.OperationField, "connection_health_routine",
						"will_retry_next_cycle", true)
					continue
				}

				c.logger.LogInfo(ctx, "Reconnection successful, re-registering worker",
					logging.OperationField, "connection_health_routine")
				// Re-register after successful reconnection
				if err := c.registerWorkerWithRetry(ctx); err != nil {
					c.logger.LogError(ctx, err, "Failed to re-register after reconnection",
						logging.OperationField, "connection_health_routine")
				} else {
					c.logger.LogInfo(ctx, "Worker re-registration completed successfully",
						logging.OperationField, "connection_health_routine")
				}
			} else {
				c.logger.LogDebug(ctx, "Connection health check passed",
					logging.OperationField, "connection_health_routine")
			}
		}
	}
}

func (c *Client) heartbeatRoutine(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second) // Default heartbeat interval
	defer ticker.Stop()

	c.logger.LogDebug(ctx, "Heartbeat routine started",
		logging.OperationField, "heartbeat_routine",
		"interval", "30s")

	for {
		select {
		case <-ctx.Done():
			c.logger.LogDebug(ctx, "Heartbeat routine stopping due to context cancellation",
				logging.OperationField, "heartbeat_routine")
			return
		case <-ticker.C:
			c.logger.LogDebug(ctx, "Sending heartbeat to API server",
				logging.OperationField, "heartbeat_routine")
			if err := c.sendHeartbeatWithRetry(ctx); err != nil {
				c.logger.LogError(ctx, err, "Heartbeat failed after retries",
					logging.OperationField, "heartbeat_routine",
					"will_mark_disconnected", true)
				c.setConnected(false)
			} else {
				c.logger.LogDebug(ctx, "Heartbeat sent successfully",
					logging.OperationField, "heartbeat_routine")
			}
		}
	}
}

func (c *Client) sendHeartbeatWithRetry(ctx context.Context) error {
	return RetryOperation(ctx, c.logger, "heartbeat", c.retryConfig, func() error {
		return c.doSendHeartbeat(ctx)
	})
}

func (c *Client) doSendHeartbeat(ctx context.Context) error {
	c.stats.mutex.RLock()

	// Check service health statuses
	services := make(map[string]string)

	// Check Ollama health
	ollamaCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := c.aiService.HealthCheck(ollamaCtx); err != nil {
		services["ollama"] = "unhealthy"
		c.logger.LogWarn(ctx, "Ollama health check failed",
			logging.OperationField, "do_send_heartbeat",
			logging.ErrorField, err)
	} else {
		services["ollama"] = "healthy"
	}

	// Check gRPC connection
	if c.isConnected() {
		services["grpc"] = "healthy"
	} else {
		services["grpc"] = "unhealthy"
	}

	stats := &workerpb.WorkerStats{
		ActiveTasks:          c.stats.ActiveTasks,
		CompletedTasks:       c.stats.CompletedTasks,
		FailedTasks:          c.stats.FailedTasks,
		CpuUsage:             getCPUUsage(),
		MemoryUsage:          getMemoryUsage(),
		Uptime:               timestamppb.New(c.stats.StartTime),
		GrpcConnectionStatus: c.getConnectionStatus(),
		Services:             services,
	}
	c.stats.mutex.RUnlock()

	req := &workerpb.WorkerHeartbeatRequest{
		WorkerId:     c.workerID,
		SessionToken: c.sessionToken,
		Status:       c.getWorkerStatus(),
		Stats:        stats,
	}

	c.logger.LogDebug(ctx, "Sending heartbeat with worker statistics",
		logging.OperationField, "do_send_heartbeat",
		"active_tasks", stats.ActiveTasks,
		"completed_tasks", stats.CompletedTasks,
		"failed_tasks", stats.FailedTasks,
		"memory_usage_mb", stats.MemoryUsage,
		"connection_status", stats.GrpcConnectionStatus,
		"ollama_status", services["ollama"],
		"grpc_status", services["grpc"],
		"worker_status", req.Status)

	return c.connectionManager.ExecuteWithRetry(ctx, "heartbeat", func(client workerpb.APIWorkerServiceClient) error {
		resp, err := client.WorkerHeartbeat(ctx, req)
		if err != nil {
			c.logger.LogWarn(ctx, "Heartbeat request failed",
				logging.OperationField, "do_send_heartbeat",
				logging.ErrorField, err)
			return fmt.Errorf("heartbeat request failed: %w", err)
		}

		if !resp.ConnectionHealthy {
			c.logger.LogWarn(ctx, "Server reports connection unhealthy",
				logging.OperationField, "do_send_heartbeat",
				"message", resp.Message,
				"server_time", resp.ServerTime)
		} else {
			c.logger.LogDebug(ctx, "Heartbeat acknowledged by server",
				logging.OperationField, "do_send_heartbeat",
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
				c.logger.LogError(ctx, err, "Task streaming failed after retries",
					logging.OperationField, "task_streaming_routine")
				c.setConnected(false)

				// Exponential backoff before retry
				backoffDelay := c.retryConfig.CalculateDelay(1)
				c.logger.LogInfo(ctx, "Backing off before reconnecting task stream",
					logging.OperationField, "task_streaming_routine",
					"delay", backoffDelay)

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
	return RetryOperation(ctx, c.logger, "task_streaming", c.retryConfig, func() error {
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
	c.logger.LogInfo(ctx, "Task streaming started successfully",
		logging.OperationField, "stream_tasks")

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
		c.logger.LogError(ctx, err, "Task processing failed through circuit breaker",
			logging.OperationField, "handle_task",
			"task_id", task.TaskId)

		// Report failure
		c.reportTaskResultWithRetry(ctx, task.TaskId, workerpb.TaskStatus_TASK_STATUS_FAILED, "", err)
	}
}

func (c *Client) processTask(ctx context.Context, task *workerpb.TaskRequest) error {
	c.logger.LogInfo(ctx, "Processing task",
		logging.OperationField, "process_task",
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
		c.logger.LogError(ctx, processErr, "Task processing failed",
			logging.OperationField, "process_task",
			"task_id", task.TaskId,
			"duration", time.Since(startTime))
	} else {
		c.logger.LogInfo(ctx, "Task completed successfully",
			logging.OperationField, "process_task",
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

	err := RetryOperation(ctx, c.logger, "insight_generation", c.retryConfig, func() error {
		var err error
		result, err = c.processInsightTask(ctx, task)
		return err
	})

	return result, err
}

func (c *Client) processWeeklyReportTaskWithRetry(ctx context.Context, task *workerpb.TaskRequest) (string, error) {
	var result string

	err := RetryOperation(ctx, c.logger, "weekly_report_generation", c.retryConfig, func() error {
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
	err := RetryOperation(ctx, c.logger, "update_task_progress", c.retryConfig, func() error {
		return c.doUpdateTaskProgress(ctx, taskID, progress, message)
	})

	if err != nil {
		c.logger.LogError(ctx, err, "Failed to update task progress after retries",
			logging.OperationField, "update_task_progress_with_retry",
			"task_id", taskID)
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
	err := RetryOperation(ctx, c.logger, "report_task_result", c.retryConfig, func() error {
		return c.doReportTaskResult(ctx, taskID, status, result, taskErr)
	})

	if err != nil {
		c.logger.LogError(ctx, err, "Failed to report task result after retries",
			logging.OperationField, "report_task_result_with_retry",
			"task_id", taskID,
			"status", status)
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

	c.logger.LogInfo(context.Background(), "Worker statistics",
		logging.OperationField, "log_statistics",
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
