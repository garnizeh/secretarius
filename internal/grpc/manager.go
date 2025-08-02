package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/logging"
	workerpb "github.com/garnizeh/englog/proto/worker"
)

// Manager manages the gRPC server lifecycle
type Manager struct {
	server     *Server
	grpcServer *grpc.Server
	config     *config.Config
	logger     *logging.Logger
	listener   net.Listener
	mu         sync.Mutex
	stopped    bool
}

// NewManager creates a new gRPC manager
func NewManager(cfg *config.Config, logger *logging.Logger) *Manager {
	managerLogger := logger.WithComponent("grpc-manager")

	return &Manager{
		server: NewServer(cfg, logger),
		config: cfg,
		logger: managerLogger,
	}
}

// Start starts the gRPC server
func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Verificar se já foi iniciado
	if m.grpcServer != nil {
		return fmt.Errorf("gRPC server already started")
	}

	start := time.Now()
	address := fmt.Sprintf(":%d", m.config.GRPC.ServerPort)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		m.logger.LogError(context.Background(), err, "Failed to listen on address",
			"address", address)
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}
	m.listener = lis

	// Configure gRPC server options
	var opts []grpc.ServerOption

	// Add TLS if enabled
	if m.config.GRPC.TLSEnabled {
		creds, err := credentials.NewServerTLSFromFile(
			m.config.GRPC.TLSCertFile,
			m.config.GRPC.TLSKeyFile,
		)
		if err != nil {
			m.logger.LogError(context.Background(), err, "Failed to load TLS credentials",
				"cert_file", m.config.GRPC.TLSCertFile,
				"key_file", m.config.GRPC.TLSKeyFile)
			return fmt.Errorf("failed to load TLS credentials: %w", err)
		}
		opts = append(opts, grpc.Creds(creds))
		m.logger.Info("gRPC server configured with TLS",
			"cert", m.config.GRPC.TLSCertFile)
	} else {
		m.logger.Warn("gRPC server running without TLS - not recommended for production")
	}

	// Create gRPC server
	m.grpcServer = grpc.NewServer(opts...)

	// Register our service
	workerpb.RegisterAPIWorkerServiceServer(m.grpcServer, m.server)

	// Reset stopped flag
	m.stopped = false

	setupDuration := time.Since(start)
	m.logger.Info("Starting gRPC server",
		"address", address,
		"tls_enabled", m.config.GRPC.TLSEnabled,
		"setup_duration_ms", setupDuration.Milliseconds())

	// Start server in goroutine - capture grpcServer to avoid race condition
	grpcServer := m.grpcServer
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			m.logger.LogError(context.Background(), err, "gRPC server failed",
				"address", address)
		}
	}()

	// Start periodic worker cleanup
	m.server.StartPeriodicCleanup(context.Background())

	return nil
}

// Stop gracefully stops the gRPC server
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Verificar se já foi parado
	if m.stopped {
		return nil // Já parado, não é erro
	}

	start := time.Now()
	m.stopped = true

	if m.grpcServer != nil {
		m.logger.Info("Stopping gRPC server")
		m.grpcServer.GracefulStop()

		duration := time.Since(start)
		m.logger.LogShutdown("grpc-manager", "graceful_stop", true)
		m.logger.Info("gRPC server stopped gracefully",
			"shutdown_duration_ms", duration.Milliseconds())
		m.grpcServer = nil
	}

	if m.listener != nil {
		if err := m.listener.Close(); err != nil {
			// Log o erro mas não retorne, pois já marcamos como stopped
			m.logger.LogError(context.Background(), err, "Error closing listener")
		}
		m.listener = nil
	}

	return nil
}

// GetServer returns the underlying gRPC server implementation
func (m *Manager) GetServer() *Server {
	return m.server
}

// QueueInsightGenerationTask queues an insight generation task
func (m *Manager) QueueInsightGenerationTask(userID string, entryIDs []string, insightType string, contextData any) (string, error) {
	start := time.Now()
	taskID := fmt.Sprintf("insight_%s_%d", userID, time.Now().Unix())

	m.logger.Info("Queuing insight generation task",
		"task_id", taskID,
		"user_id", userID,
		"entry_count", len(entryIDs),
		"insight_type", insightType)

	// Create task payload
	payload := map[string]any{
		"user_id":      userID,
		"entry_ids":    entryIDs,
		"insight_type": insightType,
		"context":      contextData,
	}

	payloadJSON, err := jsonMarshal(payload)
	if err != nil {
		m.logger.LogError(context.Background(), err, "Failed to marshal insight task payload",
			"task_id", taskID,
			"user_id", userID)
		return "", fmt.Errorf("failed to marshal task payload: %w", err)
	}

	// CHECK: How to use entryIDs in the task request? Adjust the deadlline and priority as needed
	task := &workerpb.TaskRequest{
		TaskId:   taskID,
		TaskType: workerpb.TaskType_TASK_TYPE_INSIGHT_GENERATION,
		Payload:  string(payloadJSON),
		Priority: 5, // Medium priority
		Deadline: timestamppb.New(time.Now().Add(5 * time.Minute)),
		Metadata: map[string]string{
			"user_id":      userID,
			"insight_type": insightType,
		},
	}

	err = m.server.QueueTask(task)
	duration := time.Since(start)

	if err != nil {
		m.logger.LogError(context.Background(), err, "Failed to queue insight generation task",
			"task_id", taskID,
			"user_id", userID,
			"duration_ms", duration.Milliseconds())
		return "", err
	}

	m.logger.Info("Insight generation task queued successfully",
		"task_id", taskID,
		"user_id", userID,
		"duration_ms", duration.Milliseconds())

	return taskID, nil
}

// QueueWeeklyReportTask queues a weekly report generation task
func (m *Manager) QueueWeeklyReportTask(userID string, weekStart, weekEnd time.Time) (string, error) {
	start := time.Now()
	taskID := fmt.Sprintf("report_%s_%d", userID, time.Now().Unix())

	m.logger.Info("Queuing weekly report task",
		"task_id", taskID,
		"user_id", userID,
		"week_start", weekStart.Format("2006-01-02"),
		"week_end", weekEnd.Format("2006-01-02"))

	// Create task payload
	payload := map[string]any{
		"user_id":    userID,
		"week_start": weekStart.Format(time.RFC3339),
		"week_end":   weekEnd.Format(time.RFC3339),
	}

	payloadJSON, err := jsonMarshal(payload)
	if err != nil {
		m.logger.LogError(context.Background(), err, "Failed to marshal weekly report task payload",
			"task_id", taskID,
			"user_id", userID)
		return "", fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := &workerpb.TaskRequest{
		TaskId:   taskID,
		TaskType: workerpb.TaskType_TASK_TYPE_WEEKLY_REPORT,
		Payload:  string(payloadJSON),
		Priority: 3, // Lower priority than insights
		Deadline: timestamppb.New(time.Now().Add(15 * time.Minute)),
		Metadata: map[string]string{
			"user_id": userID,
			"period":  fmt.Sprintf("%s_to_%s", weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02")),
		},
	}

	err = m.server.QueueTask(task)
	duration := time.Since(start)

	if err != nil {
		m.logger.LogError(context.Background(), err, "Failed to queue weekly report task",
			"task_id", taskID,
			"user_id", userID,
			"duration_ms", duration.Milliseconds())
		return "", err
	}

	m.logger.Info("Weekly report task queued successfully",
		"task_id", taskID,
		"user_id", userID,
		"duration_ms", duration.Milliseconds())

	return taskID, nil
}

// GetTaskResult retrieves the result of a completed task
func (m *Manager) GetTaskResult(taskID string) (*TaskResult, bool) {
	start := time.Now()

	result, found := m.server.GetTaskResult(taskID)
	duration := time.Since(start)

	if found {
		m.logger.Debug("Task result retrieved",
			"task_id", taskID,
			"worker_id", result.WorkerID,
			"status", result.Status,
			"duration_ms", duration.Milliseconds())
	} else {
		m.logger.Debug("Task result not found",
			"task_id", taskID,
			"duration_ms", duration.Milliseconds())
	}

	return result, found
}

// GetActiveWorkers returns information about active workers
func (m *Manager) GetActiveWorkers() map[string]*WorkerInfo {
	return m.server.GetActiveWorkers()
}

// HealthCheck performs a health check of the gRPC server
func (m *Manager) HealthCheck(ctx context.Context) error {
	if m.server == nil {
		return fmt.Errorf("gRPC server not initialized")
	}

	// Try to call the health check endpoint
	_, err := m.server.HealthCheck(ctx, nil)
	return err
}

// Helper function for JSON marshaling
func jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
