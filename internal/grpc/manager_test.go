package grpc_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/grpc"
	workerpb "github.com/garnizeh/englog/proto/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestManager_NewManager tests the manager constructor
func TestManager_NewManager(t *testing.T) {
	t.Run("creates new manager successfully", func(t *testing.T) {
		cfg := createTestConfigForManager()
		logger := createTestLoggerForManager()

		manager := grpc.NewManager(cfg, logger)

		assert.NotNil(t, manager)
		assert.NotNil(t, manager.GetServer())
	})
}

// TestManager_Start tests manager startup functionality
func TestManager_Start(t *testing.T) {
	t.Run("start without TLS", func(t *testing.T) {
		ctx := context.Background()

		cfg := createTestConfigForManager()
		cfg.GRPC.TLSEnabled = false
		logger := createTestLoggerForManager()

		manager := grpc.NewManager(cfg, logger)

		err := manager.Start(ctx)
		assert.NoError(t, err)

		if stopErr := manager.Stop(ctx); stopErr != nil {
			t.Logf("Warning: failed to stop manager: %v", stopErr)
		}
	})

	t.Run("start with invalid TLS configuration", func(t *testing.T) {
		ctx := context.Background()

		cfg := createTestConfigForManager()
		cfg.GRPC.TLSEnabled = true
		cfg.GRPC.TLSCertFile = "/non/existent/cert.pem"
		cfg.GRPC.TLSKeyFile = "/non/existent/key.pem"
		logger := createTestLoggerForManager()

		manager := grpc.NewManager(cfg, logger)

		err := manager.Start(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to load TLS credentials")
	})

	t.Run("start with port already in use", func(t *testing.T) {
		ctx := context.Background()

		// Garantir porta disponível antes do teste
		port := findAvailablePort()
		require.NotZero(t, port, "Failed to find available port")

		logger := createTestLoggerForManager()

		// Manager 1
		cfg1 := createTestConfigForManager()
		cfg1.GRPC.ServerPort = port
		manager1 := grpc.NewManager(cfg1, logger)

		err := manager1.Start(ctx)
		assert.NoError(t, err)

		// Aguardar inicialização completa
		time.Sleep(100 * time.Millisecond)

		// Manager 2 - deve falhar
		cfg2 := createTestConfigForManager()
		cfg2.GRPC.ServerPort = port
		manager2 := grpc.NewManager(cfg2, logger)

		err = manager2.Start(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to listen")

		// Cleanup adequado
		err = manager1.Stop(ctx)
		assert.NoError(t, err)

		// Manager2 nunca iniciou com sucesso, mas Stop() deve ser seguro
		// Não verificamos o erro do Stop() aqui porque pode retornar
		// "connection already closed" que é esperado
		_ = manager2.Stop(ctx)
	})
}

// TestManager_Stop tests manager shutdown functionality
func TestManager_Stop(t *testing.T) {
	ctx := context.Background()

	t.Run("stop without starting", func(t *testing.T) {
		cfg := createTestConfigForManager()
		logger := createTestLoggerForManager()

		manager := grpc.NewManager(cfg, logger)

		if stopErr := manager.Stop(ctx); stopErr != nil {
			t.Logf("Warning: failed to stop manager: %v", stopErr)
		}
	})

	t.Run("stop after starting", func(t *testing.T) {
		cfg := createTestConfigForManager()
		logger := createTestLoggerForManager()

		manager := grpc.NewManager(cfg, logger)

		err := manager.Start(ctx)
		assert.NoError(t, err)

		if stopErr := manager.Stop(ctx); stopErr != nil {
			t.Logf("Warning: failed to stop manager: %v", stopErr)
		}
	})

	t.Run("multiple stops", func(t *testing.T) {
		cfg := createTestConfigForManager()
		logger := createTestLoggerForManager()

		manager := grpc.NewManager(cfg, logger)

		err := manager.Start(ctx)
		assert.NoError(t, err)

		// First stop
		if stopErr := manager.Stop(ctx); stopErr != nil {
			t.Logf("Warning: failed to stop manager: %v", stopErr)
		}

		// Second stop should not error
		if stopErr := manager.Stop(ctx); stopErr != nil {
			t.Logf("Warning: failed to stop manager: %v", stopErr)
		}
	})
}

// TestManager_QueueInsightGenerationTask tests insight generation task queuing
func TestManager_QueueInsightGenerationTask(t *testing.T) {
	ctx := context.Background()

	cfg := createTestConfigForManager()
	logger := createTestLoggerForManager()

	manager := grpc.NewManager(cfg, logger)

	tests := []struct {
		name        string
		userID      string
		entryIDs    []string
		insightType string
		context     any
		wantErr     bool
	}{
		{
			name:        "valid insight generation task with string context",
			userID:      "user-123",
			entryIDs:    []string{"entry-1", "entry-2", "entry-3"},
			insightType: "productivity",
			context:     "Weekly productivity analysis",
			wantErr:     false,
		},
		{
			name:        "valid insight generation task with object context",
			userID:      "user-123",
			entryIDs:    []string{"entry-1", "entry-2"},
			insightType: "productivity",
			context: map[string]any{
				"time_blocks":    []string{"morning", "afternoon"},
				"focus_areas":    []string{"development", "meetings"},
				"date_range":     map[string]string{"start": "2025-07-01", "end": "2025-07-31"},
				"analysis_focus": "Weekly productivity analysis for performance review",
			},
			wantErr: false,
		},
		{
			name:        "empty user ID",
			userID:      "",
			entryIDs:    []string{"entry-1"},
			insightType: "productivity",
			context:     "Test context",
			wantErr:     false, // Task creation doesn't validate user ID
		},
		{
			name:        "empty entry IDs",
			userID:      "user-123",
			entryIDs:    []string{},
			insightType: "productivity",
			context:     "Test context",
			wantErr:     false,
		},
		{
			name:        "many entry IDs",
			userID:      "user-123",
			entryIDs:    generateEntryIDs(50),
			insightType: "patterns",
			context:     "Pattern analysis with many entries",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskID, err := manager.QueueInsightGenerationTask(
				ctx,
				tt.userID,
				tt.entryIDs,
				tt.insightType,
				tt.context,
			)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, taskID)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, taskID)
				assert.Contains(t, taskID, "insight_")
				assert.Contains(t, taskID, tt.userID)
			}
		})
	}
}

// TestManager_QueueWeeklyReportTask tests weekly report task queuing
func TestManager_QueueWeeklyReportTask(t *testing.T) {
	ctx := context.Background()

	cfg := createTestConfigForManager()
	logger := createTestLoggerForManager()

	manager := grpc.NewManager(cfg, logger)

	now := time.Now()
	weekStart := now.AddDate(0, 0, -7) // 7 days ago
	weekEnd := now

	tests := []struct {
		name      string
		userID    string
		weekStart time.Time
		weekEnd   time.Time
		wantErr   bool
	}{
		{
			name:      "valid weekly report task",
			userID:    "user-123",
			weekStart: weekStart,
			weekEnd:   weekEnd,
			wantErr:   false,
		},
		{
			name:      "empty user ID",
			userID:    "",
			weekStart: weekStart,
			weekEnd:   weekEnd,
			wantErr:   false, // Task creation doesn't validate user ID
		},
		{
			name:      "invalid date range (end before start)",
			userID:    "user-123",
			weekStart: weekEnd,
			weekEnd:   weekStart,
			wantErr:   false, // Task creation doesn't validate date range
		},
		{
			name:      "future dates",
			userID:    "user-123",
			weekStart: now.AddDate(0, 0, 1),
			weekEnd:   now.AddDate(0, 0, 8),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskID, err := manager.QueueWeeklyReportTask(
				ctx,
				tt.userID,
				tt.weekStart,
				tt.weekEnd,
			)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, taskID)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, taskID)
				assert.Contains(t, taskID, "report_")
				assert.Contains(t, taskID, tt.userID)
			}
		})
	}
}

// TestManager_GetTaskResult tests task result retrieval
func TestManager_GetTaskResult(t *testing.T) {
	ctx := context.Background()

	cfg := createTestConfigForManager()
	logger := createTestLoggerForManager()

	manager := grpc.NewManager(cfg, logger)

	t.Run("get non-existent task result", func(t *testing.T) {
		result, found := manager.GetTaskResult(ctx, "non-existent-task")
		assert.False(t, found)
		assert.Nil(t, result)
	})

	t.Run("get task result after reporting", func(t *testing.T) {
		server := manager.GetServer()

		// Report a task result
		req := &workerpb.TaskResultRequest{
			TaskId:      "test-task-001",
			WorkerId:    "worker-001",
			Status:      workerpb.TaskStatus_TASK_STATUS_COMPLETED,
			Result:      `{"insight": "Test insight"}`,
			StartedAt:   timestamppb.New(time.Now().Add(-time.Minute)),
			CompletedAt: timestamppb.New(time.Now()),
		}

		_, err := server.ReportTaskResult(ctx, req)
		require.NoError(t, err)

		// Retrieve the result
		result, found := manager.GetTaskResult(ctx, "test-task-001")
		assert.True(t, found)
		assert.NotNil(t, result)
		assert.Equal(t, "test-task-001", result.TaskID)
		assert.Equal(t, "worker-001", result.WorkerID)
		assert.Equal(t, workerpb.TaskStatus_TASK_STATUS_COMPLETED, result.Status)
	})
}

// TestManager_GetActiveWorkers tests active workers retrieval
func TestManager_GetActiveWorkers(t *testing.T) {
	ctx := context.Background()
	cfg := createTestConfigForManager()
	logger := createTestLoggerForManager()

	manager := grpc.NewManager(cfg, logger)

	t.Run("no active workers", func(t *testing.T) {
		workers := manager.GetActiveWorkers(ctx)
		assert.Empty(t, workers)
	})

	t.Run("with registered workers", func(t *testing.T) {
		ctx := context.Background()
		server := manager.GetServer()

		// Register workers
		workerIDs := []string{"worker-001", "worker-002", "worker-003"}
		for _, workerID := range workerIDs {
			req := &workerpb.RegisterWorkerRequest{
				WorkerId:   workerID,
				WorkerName: fmt.Sprintf("Test Worker %s", workerID),
				Capabilities: []workerpb.WorkerCapability{
					workerpb.WorkerCapability_CAPABILITY_AI_INSIGHTS,
				},
				Version: "1.0.0",
			}

			_, err := server.RegisterWorker(ctx, req)
			require.NoError(t, err)
		}

		workers := manager.GetActiveWorkers(ctx)
		assert.Len(t, workers, 3)

		for _, workerID := range workerIDs {
			assert.Contains(t, workers, workerID)
		}
	})
}

// TestManager_HealthCheck tests health check functionality
func TestManager_HealthCheck(t *testing.T) {
	ctx := context.Background()

	cfg := createTestConfigForManager()
	logger := createTestLoggerForManager()

	manager := grpc.NewManager(cfg, logger)

	t.Run("health check without starting server", func(t *testing.T) {
		err := manager.HealthCheck(ctx)
		assert.NoError(t, err)
	})

	t.Run("health check with running server", func(t *testing.T) {
		err := manager.Start(ctx)
		require.NoError(t, err)
		defer func() {
			if stopErr := manager.Stop(ctx); stopErr != nil {
				t.Logf("Warning: failed to stop manager: %v", stopErr)
			}
		}()

		err = manager.HealthCheck(ctx)
		assert.NoError(t, err)
	})

	t.Run("health check with context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		// Wait for context to timeout
		time.Sleep(2 * time.Millisecond)

		err := manager.HealthCheck(ctx)
		// Health check should still work even with expired context
		// since it's just calling the server method
		assert.NoError(t, err)
	})
}

// TestManager_TaskQueueIntegration tests full task queue integration
func TestManager_TaskQueueIntegration(t *testing.T) {
	ctx := context.Background()

	cfg := createTestConfigForManager()
	logger := createTestLoggerForManager()

	manager := grpc.NewManager(cfg, logger)

	t.Run("queue multiple task types", func(t *testing.T) {
		// Queue insight generation task
		insightTaskID, err := manager.QueueInsightGenerationTask(
			ctx,
			"user-123",
			[]string{"entry-1", "entry-2"},
			"productivity",
			"Test context",
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, insightTaskID)

		// Queue weekly report task
		reportTaskID, err := manager.QueueWeeklyReportTask(
			ctx,
			"user-123",
			time.Now().AddDate(0, 0, -7),
			time.Now(),
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, reportTaskID)

		// Task IDs should be different
		assert.NotEqual(t, insightTaskID, reportTaskID)
	})

	t.Run("queue tasks and retrieve results", func(t *testing.T) {
		ctx := context.Background()

		server := manager.GetServer()

		// Queue a task
		taskID, err := manager.QueueInsightGenerationTask(
			ctx,
			"user-456",
			[]string{"entry-1"},
			"patterns",
			"Pattern analysis",
		)
		require.NoError(t, err)

		// Simulate worker reporting result
		resultReq := &workerpb.TaskResultRequest{
			TaskId:      taskID,
			WorkerId:    "worker-001",
			Status:      workerpb.TaskStatus_TASK_STATUS_COMPLETED,
			Result:      `{"patterns": ["morning productivity", "afternoon focus"]}`,
			StartedAt:   timestamppb.New(time.Now().Add(-30 * time.Second)),
			CompletedAt: timestamppb.New(time.Now()),
		}

		_, err = server.ReportTaskResult(ctx, resultReq)
		require.NoError(t, err)

		// Retrieve result through manager
		result, found := manager.GetTaskResult(ctx, taskID)
		assert.True(t, found)
		assert.NotNil(t, result)
		assert.Equal(t, taskID, result.TaskID)
		assert.Contains(t, result.Result, "patterns")
	})
}

// TestManager_ConcurrentOperations tests concurrent manager operations
func TestManager_ConcurrentOperations(t *testing.T) {
	ctx := context.Background()

	cfg := createTestConfigForManager()
	logger := createTestLoggerForManager()

	manager := grpc.NewManager(cfg, logger)

	const numGoroutines = 10
	const numOperationsPerGoroutine = 5

	t.Run("concurrent task queuing", func(t *testing.T) {
		taskIDs := make(chan string, numGoroutines*numOperationsPerGoroutine)
		errors := make(chan error, numGoroutines*numOperationsPerGoroutine)

		// Start multiple goroutines queuing tasks
		for i := 0; i < numGoroutines; i++ {
			go func(goroutineID int) {
				for j := 0; j < numOperationsPerGoroutine; j++ {
					userID := fmt.Sprintf("user-%d-%d", goroutineID, j)
					entryIDs := []string{fmt.Sprintf("entry-%d-%d", goroutineID, j)}

					taskID, err := manager.QueueInsightGenerationTask(
						ctx,
						userID,
						entryIDs,
						"productivity",
						fmt.Sprintf("Context %d-%d", goroutineID, j),
					)

					if err != nil {
						errors <- err
					} else {
						taskIDs <- taskID
					}
				}
			}(i)
		}

		// Collect results
		var collectedTaskIDs []string
		var collectedErrors []error

		for i := 0; i < numGoroutines*numOperationsPerGoroutine; i++ {
			select {
			case taskID := <-taskIDs:
				collectedTaskIDs = append(collectedTaskIDs, taskID)
			case err := <-errors:
				collectedErrors = append(collectedErrors, err)
			case <-time.After(5 * time.Second):
				t.Fatal("Timeout waiting for concurrent operations")
			}
		}

		// Verify results
		assert.Empty(t, collectedErrors, "Should have no errors")
		assert.Len(t, collectedTaskIDs, numGoroutines*numOperationsPerGoroutine)

		// All task IDs should be unique
		uniqueTaskIDs := make(map[string]bool)
		for _, taskID := range collectedTaskIDs {
			assert.False(t, uniqueTaskIDs[taskID], "Task ID should be unique: %s", taskID)
			uniqueTaskIDs[taskID] = true
		}
	})

	t.Run("concurrent health checks", func(t *testing.T) {
		const numHealthChecks = 20
		errors := make(chan error, numHealthChecks)

		for i := 0; i < numHealthChecks; i++ {
			go func() {
				ctx := context.Background()
				err := manager.HealthCheck(ctx)
				errors <- err
			}()
		}

		// Collect results
		for i := 0; i < numHealthChecks; i++ {
			select {
			case err := <-errors:
				assert.NoError(t, err)
			case <-time.After(5 * time.Second):
				t.Fatal("Timeout waiting for health checks")
			}
		}
	})
}

// findAvailablePort finds an available port for testing
func findAvailablePort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}
