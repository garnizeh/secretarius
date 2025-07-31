package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/logging"
	workerpb "github.com/garnizeh/englog/proto/worker"
)

// Server implements the APIWorkerService gRPC server
type Server struct {
	workerpb.UnimplementedAPIWorkerServiceServer
	cfg          *config.Config
	logger       *logging.Logger
	workers      map[string]*WorkerInfo
	workersMutex sync.RWMutex
	taskQueue    chan *workerpb.TaskRequest
	taskResults  map[string]*TaskResult
	resultsMutex sync.RWMutex
}

// WorkerInfo holds information about a registered worker
type WorkerInfo struct {
	ID            string
	Name          string
	Capabilities  []workerpb.WorkerCapability
	Version       string
	Metadata      map[string]string
	SessionToken  string
	LastHeartbeat time.Time
	Status        workerpb.WorkerStatus
	Stats         *workerpb.WorkerStats
	TaskStream    workerpb.APIWorkerService_StreamTasksServer
}

// TaskResult holds the result of a completed task
type TaskResult struct {
	TaskID      string
	WorkerID    string
	Status      workerpb.TaskStatus
	Result      string
	ErrorMsg    string
	StartedAt   time.Time
	CompletedAt time.Time
}

// NewServer creates a new gRPC server instance
func NewServer(cfg *config.Config, logger *logging.Logger) *Server {
	serverLogger := logger.WithComponent("grpc-server")

	serverLogger.LogStartup("grpc-server", "v1.0.0", map[string]any{
		"task_queue_buffer":  100,
		"heartbeat_interval": "30s",
	})

	return &Server{
		cfg:         cfg,
		logger:      serverLogger,
		workers:     make(map[string]*WorkerInfo),
		taskQueue:   make(chan *workerpb.TaskRequest, 100), // Buffer for 100 tasks
		taskResults: make(map[string]*TaskResult),
	}
}

// RegisterWorker handles worker registration
func (s *Server) RegisterWorker(ctx context.Context, req *workerpb.RegisterWorkerRequest) (*workerpb.RegisterWorkerResponse, error) {
	start := time.Now()

	s.logger.WithContext(ctx).Info("Worker registration request",
		"worker_id", req.WorkerId,
		"worker_name", req.WorkerName,
		"capabilities", req.Capabilities,
		"version", req.Version)

	// Validate request
	if req.WorkerId == "" {
		err := status.Errorf(codes.InvalidArgument, "worker_id is required")
		s.logger.LogError(ctx, err, "Worker registration failed - missing worker ID")
		return nil, err
	}

	if req.WorkerName == "" {
		err := status.Errorf(codes.InvalidArgument, "worker_name is required")
		s.logger.LogError(ctx, err, "Worker registration failed - missing worker name")
		return nil, err
	}

	// Generate session token (simplified - in production use proper JWT or similar)
	sessionToken := fmt.Sprintf("session_%s_%d", req.WorkerId, time.Now().UnixNano())

	// Store worker info
	s.workersMutex.Lock()
	existingWorker, exists := s.workers[req.WorkerId]
	s.workers[req.WorkerId] = &WorkerInfo{
		ID:            req.WorkerId,
		Name:          req.WorkerName,
		Capabilities:  req.Capabilities,
		Version:       req.Version,
		Metadata:      req.Metadata,
		SessionToken:  sessionToken,
		LastHeartbeat: time.Now(),
		Status:        workerpb.WorkerStatus_WORKER_STATUS_IDLE,
	}
	s.workersMutex.Unlock()

	duration := time.Since(start)

	if exists {
		s.logger.WithContext(ctx).Info("Worker re-registered successfully",
			"worker_id", req.WorkerId,
			"previous_status", existingWorker.Status,
			"duration_ms", duration.Milliseconds())
	} else {
		s.logger.WithContext(ctx).Info("Worker registered successfully",
			"worker_id", req.WorkerId,
			"duration_ms", duration.Milliseconds())
	}

	return &workerpb.RegisterWorkerResponse{
		SessionToken:             sessionToken,
		HeartbeatIntervalSeconds: 30, // 30 seconds heartbeat interval
		RegistrationSuccessful:   true,
		Message:                  "Worker registered successfully",
	}, nil
}

// WorkerHeartbeat handles heartbeat from workers
func (s *Server) WorkerHeartbeat(ctx context.Context, req *workerpb.WorkerHeartbeatRequest) (*workerpb.WorkerHeartbeatResponse, error) {
	start := time.Now()

	s.workersMutex.Lock()
	worker, exists := s.workers[req.WorkerId]
	if !exists {
		s.workersMutex.Unlock()
		err := status.Errorf(codes.NotFound, "Worker not found: %s", req.WorkerId)
		s.logger.LogError(ctx, err, "Heartbeat failed - worker not found",
			"worker_id", req.WorkerId)
		return nil, err
	}

	// Validate session token
	if worker.SessionToken != req.SessionToken {
		s.workersMutex.Unlock()
		err := status.Errorf(codes.Unauthenticated, "Invalid session token")
		s.logger.LogError(ctx, err, "Heartbeat failed - invalid session token",
			"worker_id", req.WorkerId)
		return nil, err
	}

	// Update worker info
	previousStatus := worker.Status
	worker.LastHeartbeat = time.Now()
	worker.Status = req.Status
	worker.Stats = req.Stats
	s.workersMutex.Unlock()

	duration := time.Since(start)

	s.logger.WithContext(ctx).Debug("Worker heartbeat received",
		"worker_id", req.WorkerId,
		"status", req.Status,
		"previous_status", previousStatus,
		"duration_ms", duration.Milliseconds())

	// Log status changes
	if previousStatus != req.Status {
		s.logger.WithContext(ctx).Info("Worker status changed",
			"worker_id", req.WorkerId,
			"from_status", previousStatus,
			"to_status", req.Status)
	}

	return &workerpb.WorkerHeartbeatResponse{
		ConnectionHealthy: true,
		Message:           "Heartbeat received",
		ServerTime:        timestamppb.Now(),
	}, nil
}

// StreamTasks provides a stream of tasks to workers
func (s *Server) StreamTasks(req *workerpb.StreamTasksRequest, stream workerpb.APIWorkerService_StreamTasksServer) error {
	ctx := stream.Context()
	start := time.Now()

	s.logger.WithContext(ctx).Info("Worker requesting task stream",
		"worker_id", req.WorkerId)

	// Validate worker
	s.workersMutex.Lock()
	worker, exists := s.workers[req.WorkerId]
	if !exists {
		s.workersMutex.Unlock()
		err := status.Errorf(codes.NotFound, "Worker not found: %s", req.WorkerId)
		s.logger.LogError(ctx, err, "Task stream failed - worker not found",
			"worker_id", req.WorkerId)
		return err
	}

	// Validate session token
	if worker.SessionToken != req.SessionToken {
		s.workersMutex.Unlock()
		err := status.Errorf(codes.Unauthenticated, "Invalid session token")
		s.logger.LogError(ctx, err, "Task stream failed - invalid session token",
			"worker_id", req.WorkerId)
		return err
	}

	// Store stream reference
	worker.TaskStream = stream
	s.workersMutex.Unlock()

	s.logger.WithContext(ctx).Info("Task stream established",
		"worker_id", req.WorkerId,
		"setup_duration_ms", time.Since(start).Milliseconds())

	var tasksProcessed int

	// Listen for context cancellation and tasks
	for {
		select {
		case <-stream.Context().Done():
			duration := time.Since(start)
			s.logger.WithContext(ctx).Info("Worker disconnected from task stream",
				"worker_id", req.WorkerId,
				"tasks_processed", tasksProcessed,
				"connection_duration_ms", duration.Milliseconds())
			return stream.Context().Err()

		case task := <-s.taskQueue:
			// Check if worker has required capability for this task
			if s.workerHasCapability(req.WorkerId, task.TaskType) {
				tasksProcessed++
				s.logger.WithContext(ctx).Info("Sending task to worker",
					"task_id", task.TaskId,
					"worker_id", req.WorkerId,
					"task_type", task.TaskType,
					"tasks_processed", tasksProcessed)

				if err := stream.Send(task); err != nil {
					s.logger.LogError(ctx, err, "Failed to send task to worker",
						"worker_id", req.WorkerId,
						"task_id", task.TaskId,
						"tasks_processed", tasksProcessed)
					return err
				}
			} else {
				s.logger.WithContext(ctx).Debug("Task skipped - worker lacks capability",
					"task_id", task.TaskId,
					"worker_id", req.WorkerId,
					"task_type", task.TaskType,
					"worker_capabilities", worker.Capabilities)
			}
		}
	}
}

// ReportTaskResult handles task completion reports from workers
func (s *Server) ReportTaskResult(ctx context.Context, req *workerpb.TaskResultRequest) (*workerpb.TaskResultResponse, error) {
	start := time.Now()

	s.logger.WithContext(ctx).Info("Task result received",
		"task_id", req.TaskId,
		"worker_id", req.WorkerId,
		"status", req.Status)

	// Validate request
	if req.TaskId == "" {
		err := status.Errorf(codes.InvalidArgument, "task_id is required")
		s.logger.LogError(ctx, err, "Task result failed - missing task ID")
		return nil, err
	}

	if req.WorkerId == "" {
		err := status.Errorf(codes.InvalidArgument, "worker_id is required")
		s.logger.LogError(ctx, err, "Task result failed - missing worker ID",
			"task_id", req.TaskId)
		return nil, err
	}

	// Calculate task duration if timestamps are provided
	var taskDuration time.Duration
	if req.StartedAt != nil && req.CompletedAt != nil {
		taskDuration = req.CompletedAt.AsTime().Sub(req.StartedAt.AsTime())
	}

	// Store result
	s.resultsMutex.Lock()
	s.taskResults[req.TaskId] = &TaskResult{
		TaskID:      req.TaskId,
		WorkerID:    req.WorkerId,
		Status:      req.Status,
		Result:      req.Result,
		ErrorMsg:    req.ErrorMessage,
		StartedAt:   req.StartedAt.AsTime(),
		CompletedAt: req.CompletedAt.AsTime(),
	}
	s.resultsMutex.Unlock()

	duration := time.Since(start)

	// Log based on task status
	logAttrs := []any{
		"task_id", req.TaskId,
		"worker_id", req.WorkerId,
		"status", req.Status,
		"processing_duration_ms", duration.Milliseconds(),
	}

	if taskDuration > 0 {
		logAttrs = append(logAttrs, "task_duration_ms", taskDuration.Milliseconds())
	}

	switch req.Status {
	case workerpb.TaskStatus_TASK_STATUS_COMPLETED:
		s.logger.WithContext(ctx).Info("Task completed successfully", logAttrs...)
	case workerpb.TaskStatus_TASK_STATUS_FAILED:
		logAttrs = append(logAttrs, "error_message", req.ErrorMessage)
		s.logger.WithContext(ctx).Warn("Task failed", logAttrs...)
	default:
		s.logger.WithContext(ctx).Info("Task result processed", logAttrs...)
	}

	return &workerpb.TaskResultResponse{
		ResultReceived: true,
		Message:        "Task result received successfully",
	}, nil
}

// UpdateTaskProgress handles task progress updates from workers
func (s *Server) UpdateTaskProgress(ctx context.Context, req *workerpb.TaskProgressRequest) (*emptypb.Empty, error) {
	s.logger.WithContext(ctx).Debug("Task progress update",
		"task_id", req.TaskId,
		"worker_id", req.WorkerId,
		"progress", req.ProgressPercent)

	// Validate progress percentage
	if req.ProgressPercent < 0 || req.ProgressPercent > 100 {
		s.logger.WithContext(ctx).Warn("Invalid progress percentage received",
			"task_id", req.TaskId,
			"worker_id", req.WorkerId,
			"progress", req.ProgressPercent)
	}

	// In a real implementation, you might want to store progress updates
	// or notify interested parties

	return &emptypb.Empty{}, nil
}

// HealthCheck provides health status of the gRPC server
func (s *Server) HealthCheck(ctx context.Context, req *emptypb.Empty) (*workerpb.HealthCheckResponse, error) {
	start := time.Now()

	s.workersMutex.RLock()
	activeWorkers := len(s.workers)
	totalTasksQueued := len(s.taskQueue)

	// Count workers by status
	statusCounts := make(map[workerpb.WorkerStatus]int)
	healthyWorkers := 0

	for _, worker := range s.workers {
		statusCounts[worker.Status]++
		// Consider workers healthy if they've sent a heartbeat in the last 2 minutes
		if time.Since(worker.LastHeartbeat) < time.Minute*2 {
			healthyWorkers++
		}
	}

	services := make(map[string]string)
	services["grpc"] = "healthy"
	services["task_queue"] = "healthy"
	s.workersMutex.RUnlock()

	s.resultsMutex.RLock()
	totalTaskResults := len(s.taskResults)
	s.resultsMutex.RUnlock()

	duration := time.Since(start)

	s.logger.WithContext(ctx).Debug("Health check completed",
		"active_workers", activeWorkers,
		"healthy_workers", healthyWorkers,
		"tasks_queued", totalTasksQueued,
		"task_results", totalTaskResults,
		"duration_ms", duration.Milliseconds())

	return &workerpb.HealthCheckResponse{
		Status:        "healthy",
		Timestamp:     timestamppb.Now(),
		Services:      services,
		ActiveWorkers: int32(activeWorkers),
	}, nil
}

// Helper methods

// workerHasCapability checks if a worker has the required capability for a task type
func (s *Server) workerHasCapability(workerID string, taskType workerpb.TaskType) bool {
	s.workersMutex.RLock()
	defer s.workersMutex.RUnlock()

	worker, exists := s.workers[workerID]
	if !exists {
		return false
	}

	// Map task types to required capabilities
	requiredCapability := s.getRequiredCapability(taskType)
	if requiredCapability == workerpb.WorkerCapability_CAPABILITY_UNSPECIFIED {
		return true // No specific capability required
	}

	for _, capability := range worker.Capabilities {
		if capability == requiredCapability {
			return true
		}
	}

	return false
}

// getRequiredCapability maps task types to required capabilities
func (s *Server) getRequiredCapability(taskType workerpb.TaskType) workerpb.WorkerCapability {
	switch taskType {
	case workerpb.TaskType_TASK_TYPE_INSIGHT_GENERATION:
		return workerpb.WorkerCapability_CAPABILITY_AI_INSIGHTS
	case workerpb.TaskType_TASK_TYPE_WEEKLY_REPORT:
		return workerpb.WorkerCapability_CAPABILITY_WEEKLY_REPORTS
	case workerpb.TaskType_TASK_TYPE_DATA_ANALYSIS:
		return workerpb.WorkerCapability_CAPABILITY_DATA_ANALYSIS
	case workerpb.TaskType_TASK_TYPE_NOTIFICATION:
		return workerpb.WorkerCapability_CAPABILITY_NOTIFICATIONS
	default:
		return workerpb.WorkerCapability_CAPABILITY_UNSPECIFIED
	}
}

// QueueTask adds a task to the task queue
func (s *Server) QueueTask(task *workerpb.TaskRequest) error {
	start := time.Now()

	// Validate task
	if task.TaskId == "" {
		err := fmt.Errorf("task_id is required")
		s.logger.Error("Task queue failed - missing task ID", "error", err.Error())
		return err
	}

	select {
	case s.taskQueue <- task:
		duration := time.Since(start)
		queueSize := len(s.taskQueue)

		s.logger.Info("Task queued",
			"task_id", task.TaskId,
			"task_type", task.TaskType,
			"queue_size", queueSize,
			"duration_ms", duration.Milliseconds())
		return nil
	default:
		err := fmt.Errorf("task queue is full")
		s.logger.Error("Task queue failed - queue is full",
			"task_id", task.TaskId,
			"task_type", task.TaskType,
			"queue_capacity", cap(s.taskQueue),
			"error", err.Error())
		return err
	}
}

// GetTaskResult retrieves the result of a completed task
func (s *Server) GetTaskResult(taskID string) (*TaskResult, bool) {
	s.resultsMutex.RLock()
	defer s.resultsMutex.RUnlock()
	result, exists := s.taskResults[taskID]
	return result, exists
}

// GetActiveWorkers returns information about all active workers
func (s *Server) GetActiveWorkers() map[string]*WorkerInfo {
	start := time.Now()

	s.workersMutex.RLock()
	defer s.workersMutex.RUnlock()

	// Create a copy to avoid race conditions
	workers := make(map[string]*WorkerInfo)
	totalWorkers := len(s.workers)
	stalledWorkers := 0

	for id, worker := range s.workers {
		// Only include workers that have sent a heartbeat recently
		if time.Since(worker.LastHeartbeat) < time.Minute*2 {
			workers[id] = worker
		} else {
			stalledWorkers++
		}
	}

	duration := time.Since(start)
	activeWorkers := len(workers)

	s.logger.Debug("Active workers retrieved",
		"total_workers", totalWorkers,
		"active_workers", activeWorkers,
		"stalled_workers", stalledWorkers,
		"duration_ms", duration.Milliseconds())

	return workers
}

// Start starts the gRPC server
func (s *Server) Start(address string) error {
	start := time.Now()

	lis, err := net.Listen("tcp", address)
	if err != nil {
		s.logger.LogError(context.Background(), err, "Failed to listen on address",
			"address", address)
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	grpcServer := grpc.NewServer()
	workerpb.RegisterAPIWorkerServiceServer(grpcServer, s)

	setupDuration := time.Since(start)
	s.logger.Info("Starting gRPC server",
		"address", address,
		"setup_duration_ms", setupDuration.Milliseconds())

	// Log shutdown when server stops
	defer func() {
		totalDuration := time.Since(start)
		s.logger.LogShutdown("grpc-server", "serve_completed", true)
		s.logger.Info("gRPC server stopped",
			"address", address,
			"total_runtime_ms", totalDuration.Milliseconds())
	}()

	if err := grpcServer.Serve(lis); err != nil {
		s.logger.LogError(context.Background(), err, "gRPC server failed",
			"address", address)
		return err
	}

	return nil
}
