//go:build integration
// +build integration

package grpc_test

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/grpc"
	"github.com/garnizeh/englog/internal/logging"
	workerpb "github.com/garnizeh/englog/proto/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	grpcClient "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Integration tests for gRPC server and manager working together
func TestGRPCIntegration(t *testing.T) {
	ctx := context.Background()

	t.Run("full worker lifecycle", func(t *testing.T) {
		// Setup
		cfg := createTestConfigForIntegration()
		logger := createTestLoggerForIntegration()
		manager := grpc.NewManager(cfg, logger)

		// Start server
		err := manager.Start(ctx)
		require.NoError(t, err)
		defer func() {
			if stopErr := manager.Stop(ctx); stopErr != nil {
				t.Logf("Warning: failed to stop manager: %v", stopErr)
			}
		}()
		// Give server time to start
		time.Sleep(50 * time.Millisecond)

		// Create client connection
		conn, err := grpcClient.NewClient(
			fmt.Sprintf("localhost:%d", cfg.GRPC.ServerPort),
			grpcClient.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)
		defer conn.Close()

		client := workerpb.NewAPIWorkerServiceClient(conn)
		ctx := context.Background()

		// Test 1: Register worker
		registerReq := &workerpb.RegisterWorkerRequest{
			WorkerId:   "integration-worker-001",
			WorkerName: "Integration Test Worker",
			Capabilities: []workerpb.WorkerCapability{
				workerpb.WorkerCapability_CAPABILITY_AI_INSIGHTS,
				workerpb.WorkerCapability_CAPABILITY_WEEKLY_REPORTS,
			},
			Version: "1.0.0",
			Metadata: map[string]string{
				"env":  "integration-test",
				"type": "test-worker",
			},
		}

		registerResp, err := client.RegisterWorker(ctx, registerReq)
		require.NoError(t, err)
		require.NotNil(t, registerResp)
		assert.True(t, registerResp.RegistrationSuccessful)
		assert.NotEmpty(t, registerResp.SessionToken)

		sessionToken := registerResp.SessionToken

		// Test 2: Send heartbeat
		heartbeatReq := &workerpb.WorkerHeartbeatRequest{
			WorkerId:     "integration-worker-001",
			SessionToken: sessionToken,
			Status:       workerpb.WorkerStatus_WORKER_STATUS_IDLE,
			Stats: &workerpb.WorkerStats{
				ActiveTasks:          0,
				CompletedTasks:       0,
				FailedTasks:          0,
				CpuUsage:             15.5,
				MemoryUsage:          45.2,
				Uptime:               timestamppb.New(time.Now().Add(-time.Hour)),
				GrpcConnectionStatus: "connected",
			},
		}

		heartbeatResp, err := client.WorkerHeartbeat(ctx, heartbeatReq)
		require.NoError(t, err)
		require.NotNil(t, heartbeatResp)
		assert.True(t, heartbeatResp.ConnectionHealthy)

		// Test 3: Queue tasks through manager
		insightTaskID, err := manager.QueueInsightGenerationTask(
			ctx,
			"integration-user-001",
			[]string{"entry-1", "entry-2", "entry-3"},
			"productivity",
			"Integration test context",
		)
		require.NoError(t, err)
		assert.NotEmpty(t, insightTaskID)

		reportTaskID, err := manager.QueueWeeklyReportTask(
			ctx,
			"integration-user-001",
			time.Now().AddDate(0, 0, -7),
			time.Now(),
		)
		require.NoError(t, err)
		assert.NotEmpty(t, reportTaskID)

		// Test 4: Health check
		healthResp, err := client.HealthCheck(ctx, &emptypb.Empty{})
		require.NoError(t, err)
		require.NotNil(t, healthResp)
		assert.Equal(t, "healthy", healthResp.Status)
		assert.Equal(t, int32(1), healthResp.ActiveWorkers)

		// Test 5: Report task results
		taskResultReq := &workerpb.TaskResultRequest{
			TaskId:       insightTaskID,
			WorkerId:     "integration-worker-001",
			Status:       workerpb.TaskStatus_TASK_STATUS_COMPLETED,
			Result:       `{"insights": ["User is most productive in morning", "Focus decreases after lunch"], "confidence": 0.89}`,
			ErrorMessage: "",
			StartedAt:    timestamppb.New(time.Now().Add(-2 * time.Minute)),
			CompletedAt:  timestamppb.New(time.Now()),
		}

		taskResultResp, err := client.ReportTaskResult(ctx, taskResultReq)
		require.NoError(t, err)
		require.NotNil(t, taskResultResp)
		assert.True(t, taskResultResp.ResultReceived)

		// Test 6: Retrieve task result through manager
		result, found := manager.GetTaskResult(ctx, insightTaskID)
		assert.True(t, found)
		assert.NotNil(t, result)
		assert.Equal(t, insightTaskID, result.TaskID)
		assert.Equal(t, "integration-worker-001", result.WorkerID)
		assert.Equal(t, workerpb.TaskStatus_TASK_STATUS_COMPLETED, result.Status)
		assert.Contains(t, result.Result, "insights")

		// Test 7: Get active workers through manager
		workers := manager.GetActiveWorkers(ctx)
		assert.Len(t, workers, 1)
		assert.Contains(t, workers, "integration-worker-001")

		worker := workers["integration-worker-001"]
		assert.Equal(t, "Integration Test Worker", worker.Name)
		assert.Len(t, worker.Capabilities, 2)
	})

	t.Run("multiple workers task distribution", func(t *testing.T) {
		ctx := context.Background()

		// Setup
		cfg := createTestConfigForIntegration()
		logger := createTestLoggerForIntegration()
		manager := grpc.NewManager(cfg, logger)

		err := manager.Start(ctx)
		require.NoError(t, err)
		defer func() {
			if stopErr := manager.Stop(ctx); stopErr != nil {
				t.Logf("Warning: failed to stop manager: %v", stopErr)
			}
		}()

		time.Sleep(50 * time.Millisecond)

		// Create multiple worker connections
		const numWorkers = 3
		workers := make([]workerpb.APIWorkerServiceClient, numWorkers)
		sessionTokens := make([]string, numWorkers)

		for i := 0; i < numWorkers; i++ {
			conn, err := grpcClient.NewClient(
				fmt.Sprintf("localhost:%d", cfg.GRPC.ServerPort),
				grpcClient.WithTransportCredentials(insecure.NewCredentials()),
			)
			require.NoError(t, err)
			defer conn.Close()

			client := workerpb.NewAPIWorkerServiceClient(conn)
			workers[i] = client

			// Register worker
			registerReq := &workerpb.RegisterWorkerRequest{
				WorkerId:   fmt.Sprintf("multi-worker-%03d", i+1),
				WorkerName: fmt.Sprintf("Multi Test Worker %d", i+1),
				Capabilities: []workerpb.WorkerCapability{
					workerpb.WorkerCapability_CAPABILITY_AI_INSIGHTS,
				},
				Version: "1.0.0",
			}

			ctx := context.Background()
			registerResp, err := client.RegisterWorker(ctx, registerReq)
			require.NoError(t, err)
			sessionTokens[i] = registerResp.SessionToken
		}

		// Queue multiple tasks
		const numTasks = 5
		taskIDs := make([]string, numTasks)

		for i := 0; i < numTasks; i++ {
			taskID, err := manager.QueueInsightGenerationTask(
				ctx,
				fmt.Sprintf("user-%d", i+1),
				[]string{fmt.Sprintf("entry-%d", i+1)},
				"productivity",
				fmt.Sprintf("Multi-worker context %d", i+1),
			)
			require.NoError(t, err)
			taskIDs[i] = taskID
		}

		// Verify all workers are active
		activeWorkers := manager.GetActiveWorkers(ctx)
		assert.Len(t, activeWorkers, numWorkers)

		// Health check should show all workers
		conn, err := grpcClient.NewClient(
			fmt.Sprintf("localhost:%d", cfg.GRPC.ServerPort),
			grpcClient.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)
		defer conn.Close()

		client := workerpb.NewAPIWorkerServiceClient(conn)

		healthResp, err := client.HealthCheck(ctx, &emptypb.Empty{})
		require.NoError(t, err)
		assert.Equal(t, int32(numWorkers), healthResp.ActiveWorkers)
	})

	t.Run("worker disconnection and reconnection", func(t *testing.T) {
		ctx := context.Background()

		// Setup
		cfg := createTestConfigForIntegration()
		logger := createTestLoggerForIntegration()
		manager := grpc.NewManager(cfg, logger)

		err := manager.Start(ctx)
		require.NoError(t, err)
		defer func() {
			if stopErr := manager.Stop(ctx); stopErr != nil {
				t.Logf("Warning: failed to stop manager: %v", stopErr)
			}
		}()

		time.Sleep(50 * time.Millisecond)

		// Create and register worker
		conn, err := grpcClient.NewClient(
			fmt.Sprintf("localhost:%d", cfg.GRPC.ServerPort),
			grpcClient.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)

		client := workerpb.NewAPIWorkerServiceClient(conn)

		registerReq := &workerpb.RegisterWorkerRequest{
			WorkerId:   "disconnect-worker-001",
			WorkerName: "Disconnect Test Worker",
			Capabilities: []workerpb.WorkerCapability{
				workerpb.WorkerCapability_CAPABILITY_AI_INSIGHTS,
			},
			Version: "1.0.0",
		}

		registerResp, err := client.RegisterWorker(ctx, registerReq)
		require.NoError(t, err)
		sessionToken := registerResp.SessionToken

		// Verify worker is active
		workers := manager.GetActiveWorkers(ctx)
		assert.Len(t, workers, 1)

		// Simulate disconnection by closing connection
		conn.Close()

		// Wait a bit for disconnection to be processed
		time.Sleep(100 * time.Millisecond)

		// Reconnect with same worker ID
		newConn, err := grpcClient.NewClient(
			fmt.Sprintf("localhost:%d", cfg.GRPC.ServerPort),
			grpcClient.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)
		defer newConn.Close()

		newClient := workerpb.NewAPIWorkerServiceClient(newConn)

		// Re-register (should work)
		reRegisterResp, err := newClient.RegisterWorker(ctx, registerReq)
		require.NoError(t, err)
		assert.True(t, reRegisterResp.RegistrationSuccessful)

		// New session token should be different
		assert.NotEqual(t, sessionToken, reRegisterResp.SessionToken)

		// Verify worker is active again
		workers = manager.GetActiveWorkers(ctx)
		assert.Len(t, workers, 1)
	})

	t.Run("error handling and recovery", func(t *testing.T) {
		ctx := context.Background()

		// Setup
		cfg := createTestConfigForIntegration()
		logger := createTestLoggerForIntegration()
		manager := grpc.NewManager(cfg, logger)

		err := manager.Start(ctx)
		require.NoError(t, err)
		defer func() {
			if stopErr := manager.Stop(ctx); stopErr != nil {
				t.Logf("Warning: failed to stop manager: %v", stopErr)
			}
		}()

		time.Sleep(50 * time.Millisecond)

		conn, err := grpcClient.NewClient(
			fmt.Sprintf("localhost:%d", cfg.GRPC.ServerPort),
			grpcClient.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)
		defer conn.Close()

		client := workerpb.NewAPIWorkerServiceClient(conn)

		// Test invalid registration
		invalidReq := &workerpb.RegisterWorkerRequest{
			WorkerId:   "", // Invalid: empty worker ID
			WorkerName: "Test Worker",
		}

		_, err = client.RegisterWorker(ctx, invalidReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "worker_id is required")

		// Test valid registration after error
		validReq := &workerpb.RegisterWorkerRequest{
			WorkerId:   "error-recovery-worker",
			WorkerName: "Error Recovery Worker",
			Capabilities: []workerpb.WorkerCapability{
				workerpb.WorkerCapability_CAPABILITY_AI_INSIGHTS,
			},
			Version: "1.0.0",
		}

		registerResp, err := client.RegisterWorker(ctx, validReq)
		require.NoError(t, err)
		assert.True(t, registerResp.RegistrationSuccessful)

		// Test invalid heartbeat
		invalidHeartbeat := &workerpb.WorkerHeartbeatRequest{
			WorkerId:     "non-existent-worker",
			SessionToken: "invalid-token",
			Status:       workerpb.WorkerStatus_WORKER_STATUS_IDLE,
		}

		_, err = client.WorkerHeartbeat(ctx, invalidHeartbeat)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Worker not found")

		// Test valid heartbeat after error
		validHeartbeat := &workerpb.WorkerHeartbeatRequest{
			WorkerId:     "error-recovery-worker",
			SessionToken: registerResp.SessionToken,
			Status:       workerpb.WorkerStatus_WORKER_STATUS_IDLE,
		}

		heartbeatResp, err := client.WorkerHeartbeat(ctx, validHeartbeat)
		require.NoError(t, err)
		assert.True(t, heartbeatResp.ConnectionHealthy)
	})
}

// TestGRPCConcurrency tests concurrent operations in integration scenarios
func TestGRPCConcurrency(t *testing.T) {
	t.Run("concurrent worker registrations and task queuing", func(t *testing.T) {
		ctx := context.Background()

		// Setup
		cfg := createTestConfigForIntegration()
		logger := createTestLoggerForIntegration()
		manager := grpc.NewManager(cfg, logger)

		err := manager.Start(ctx)
		require.NoError(t, err)
		defer func() {
			if stopErr := manager.Stop(ctx); stopErr != nil {
				t.Logf("Warning: failed to stop manager: %v", stopErr)
			}
		}()

		time.Sleep(50 * time.Millisecond)

		const numWorkers = 5
		const numTasksPerWorker = 3

		var wg sync.WaitGroup
		errors := make(chan error, numWorkers*(numTasksPerWorker+1)) // +1 for registration
		taskIDs := make(chan string, numWorkers*numTasksPerWorker)

		// Launch concurrent workers
		for i := range numWorkers {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Create connection
				conn, err := grpcClient.NewClient(
					fmt.Sprintf("localhost:%d", cfg.GRPC.ServerPort),
					grpcClient.WithTransportCredentials(insecure.NewCredentials()),
				)
				if err != nil {
					errors <- err
					return
				}
				defer conn.Close()

				client := workerpb.NewAPIWorkerServiceClient(conn)

				// Register worker
				registerReq := &workerpb.RegisterWorkerRequest{
					WorkerId:   fmt.Sprintf("concurrent-worker-%03d", workerID),
					WorkerName: fmt.Sprintf("Concurrent Worker %d", workerID),
					Capabilities: []workerpb.WorkerCapability{
						workerpb.WorkerCapability_CAPABILITY_AI_INSIGHTS,
					},
					Version: "1.0.0",
				}

				registerResp, err := client.RegisterWorker(ctx, registerReq)
				if err != nil {
					errors <- err
					return
				}

				// Queue tasks concurrently
				for j := 0; j < numTasksPerWorker; j++ {
					taskID, err := manager.QueueInsightGenerationTask(
						ctx,
						fmt.Sprintf("user-%d-%d", workerID, j),
						[]string{fmt.Sprintf("entry-%d-%d", workerID, j)},
						"productivity",
						fmt.Sprintf("Concurrent context %d-%d", workerID, j),
					)

					if err != nil {
						errors <- err
					} else {
						taskIDs <- taskID
					}

					// Simulate some work
					time.Sleep(1 * time.Millisecond)

					// Send heartbeat
					heartbeatReq := &workerpb.WorkerHeartbeatRequest{
						WorkerId:     registerReq.WorkerId,
						SessionToken: registerResp.SessionToken,
						Status:       workerpb.WorkerStatus_WORKER_STATUS_BUSY,
						Stats: &workerpb.WorkerStats{
							ActiveTasks:    int32(j + 1),
							CompletedTasks: int32(j),
							FailedTasks:    0,
						},
					}

					_, heartbeatErr := client.WorkerHeartbeat(ctx, heartbeatReq)
					if heartbeatErr != nil {
						errors <- heartbeatErr
					}
				}
			}(i)
		}

		wg.Wait()
		close(errors)
		close(taskIDs)

		// Check for errors
		var collectedErrors []error
		for err := range errors {
			collectedErrors = append(collectedErrors, err)
		}
		assert.Empty(t, collectedErrors, "Should have no errors during concurrent operations")

		// Check task IDs
		var collectedTaskIDs []string
		for taskID := range taskIDs {
			collectedTaskIDs = append(collectedTaskIDs, taskID)
		}
		assert.Len(t, collectedTaskIDs, numWorkers*numTasksPerWorker)

		// Verify all workers are registered
		workers := manager.GetActiveWorkers(ctx)
		assert.Len(t, workers, numWorkers)
	})
}

// Helper functions for integration tests

func createTestConfigForIntegration() *config.Config {
	// Use port 0 to get a random available port
	return &config.Config{
		GRPC: config.GRPCConfig{
			ServerPort:  getFreePort(),
			TLSEnabled:  false,
			TLSCertFile: "",
			TLSKeyFile:  "",
		},
	}
}

func createTestLoggerForIntegration() *logging.Logger {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelWarn, // Reduce log noise in integration tests
	}
	handler := slog.NewJSONHandler(&buf, opts)
	return &logging.Logger{Logger: slog.New(handler)}
}

func getFreePort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 9090 // fallback port
	}

	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	return port
}
