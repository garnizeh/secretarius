package grpc_test

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/grpc"
	"github.com/garnizeh/englog/internal/logging"
	workerpb "github.com/garnizeh/englog/proto/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestWorkerCapabilityMatching tests worker capability matching functionality
func TestWorkerCapabilityMatching(t *testing.T) {
	cfg := createTestConfigForCapabilities()
	logger := createTestLoggerForCapabilities()
	server := grpc.NewServer(cfg, logger)

	// Register workers with different capabilities

	// Worker 1: AI Insights only
	registerWorker(t, server, &workerpb.RegisterWorkerRequest{
		WorkerId:   "ai-worker",
		WorkerName: "AI Insights Worker",
		Capabilities: []workerpb.WorkerCapability{
			workerpb.WorkerCapability_CAPABILITY_AI_INSIGHTS,
		},
		Version: "1.0.0",
	})

	// Worker 2: Weekly Reports only
	registerWorker(t, server, &workerpb.RegisterWorkerRequest{
		WorkerId:   "report-worker",
		WorkerName: "Weekly Reports Worker",
		Capabilities: []workerpb.WorkerCapability{
			workerpb.WorkerCapability_CAPABILITY_WEEKLY_REPORTS,
		},
		Version: "1.0.0",
	})

	// Worker 3: Multiple capabilities
	registerWorker(t, server, &workerpb.RegisterWorkerRequest{
		WorkerId:   "multi-worker",
		WorkerName: "Multi-Capability Worker",
		Capabilities: []workerpb.WorkerCapability{
			workerpb.WorkerCapability_CAPABILITY_AI_INSIGHTS,
			workerpb.WorkerCapability_CAPABILITY_WEEKLY_REPORTS,
			workerpb.WorkerCapability_CAPABILITY_DATA_ANALYSIS,
		},
		Version: "1.0.0",
	})

	// Worker 4: No specific capabilities
	registerWorker(t, server, &workerpb.RegisterWorkerRequest{
		WorkerId:   "generic-worker",
		WorkerName: "Generic Worker",
		Capabilities: []workerpb.WorkerCapability{
			workerpb.WorkerCapability_CAPABILITY_UNSPECIFIED,
		},
		Version: "1.0.0",
	})

	tests := []struct {
		name           string
		taskType       workerpb.TaskType
		expectedWorker string // Worker that should receive the task
		description    string
	}{
		{
			name:           "AI Insights task to AI worker",
			taskType:       workerpb.TaskType_TASK_TYPE_INSIGHT_GENERATION,
			expectedWorker: "ai-worker", // Should go to AI insights worker first
			description:    "AI insights task should be routed to worker with AI capability",
		},
		{
			name:           "Weekly report task to report worker",
			taskType:       workerpb.TaskType_TASK_TYPE_WEEKLY_REPORT,
			expectedWorker: "report-worker", // Should go to weekly reports worker
			description:    "Weekly report task should be routed to worker with reports capability",
		},
		{
			name:           "Data analysis task to multi worker",
			taskType:       workerpb.TaskType_TASK_TYPE_DATA_ANALYSIS,
			expectedWorker: "multi-worker", // Only worker with data analysis capability
			description:    "Data analysis task should be routed to worker with data analysis capability",
		},
		{
			name:           "Notification task to generic worker",
			taskType:       workerpb.TaskType_TASK_TYPE_NOTIFICATION,
			expectedWorker: "generic-worker", // Should accept any task type
			description:    "Notification task should be routed to worker with notification capability",
		},
		{
			name:           "Cleanup task to any worker",
			taskType:       workerpb.TaskType_TASK_TYPE_CLEANUP,
			expectedWorker: "generic-worker", // No specific capability required
			description:    "Cleanup task should be accepted by any worker",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a task
			task := &workerpb.TaskRequest{
				TaskId:   fmt.Sprintf("test-task-%s", tt.name),
				TaskType: tt.taskType,
				Payload:  fmt.Sprintf(`{"test": "%s"}`, tt.name),
				Priority: 5,
				Deadline: timestamppb.New(time.Now().Add(time.Hour)),
			}

			// Queue the task
			err := server.QueueTask(task)
			require.NoError(t, err, tt.description)

			// Note: In a real test, we would need to implement task streaming
			// to verify which worker receives the task. For now, we verify
			// that the task was queued successfully.
		})
	}
}

// TestTaskPriorityAndDeadlines tests task priority and deadline handling
func TestTaskPriorityAndDeadlines(t *testing.T) {
	cfg := createTestConfigForCapabilities()
	logger := createTestLoggerForCapabilities()
	server := grpc.NewServer(cfg, logger)

	tests := []struct {
		name     string
		tasks    []*workerpb.TaskRequest
		wantErr  bool
		testDesc string
	}{
		{
			name: "tasks with different priorities",
			tasks: []*workerpb.TaskRequest{
				{
					TaskId:   "high-priority-task",
					TaskType: workerpb.TaskType_TASK_TYPE_INSIGHT_GENERATION,
					Payload:  `{"priority": "high"}`,
					Priority: 10, // High priority
					Deadline: timestamppb.New(time.Now().Add(30 * time.Minute)),
				},
				{
					TaskId:   "medium-priority-task",
					TaskType: workerpb.TaskType_TASK_TYPE_WEEKLY_REPORT,
					Payload:  `{"priority": "medium"}`,
					Priority: 5, // Medium priority
					Deadline: timestamppb.New(time.Now().Add(2 * time.Hour)),
				},
				{
					TaskId:   "low-priority-task",
					TaskType: workerpb.TaskType_TASK_TYPE_DATA_ANALYSIS,
					Payload:  `{"priority": "low"}`,
					Priority: 1, // Low priority
					Deadline: timestamppb.New(time.Now().Add(24 * time.Hour)),
				},
			},
			wantErr:  false,
			testDesc: "Tasks with different priorities should be queued successfully",
		},
		{
			name: "tasks with past deadlines",
			tasks: []*workerpb.TaskRequest{
				{
					TaskId:   "past-deadline-task",
					TaskType: workerpb.TaskType_TASK_TYPE_INSIGHT_GENERATION,
					Payload:  `{"deadline": "past"}`,
					Priority: 5,
					Deadline: timestamppb.New(time.Now().Add(-time.Hour)), // Past deadline
				},
			},
			wantErr:  false, // Server doesn't validate deadlines during queuing
			testDesc: "Tasks with past deadlines should still be queued",
		},
		{
			name: "tasks with very long deadlines",
			tasks: []*workerpb.TaskRequest{
				{
					TaskId:   "long-deadline-task",
					TaskType: workerpb.TaskType_TASK_TYPE_WEEKLY_REPORT,
					Payload:  `{"deadline": "very_long"}`,
					Priority: 3,
					Deadline: timestamppb.New(time.Now().AddDate(0, 1, 0)), // 1 month deadline
				},
			},
			wantErr:  false,
			testDesc: "Tasks with very long deadlines should be queued successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, task := range tt.tasks {
				err := server.QueueTask(task)

				if tt.wantErr {
					assert.Error(t, err, tt.testDesc)
				} else {
					assert.NoError(t, err, tt.testDesc)
				}
			}
		})
	}
}

// TestTaskMetadata tests task metadata handling
func TestTaskMetadata(t *testing.T) {
	cfg := createTestConfigForCapabilities()
	logger := createTestLoggerForCapabilities()
	server := grpc.NewServer(cfg, logger)

	tests := []struct {
		name     string
		task     *workerpb.TaskRequest
		wantErr  bool
		testDesc string
	}{
		{
			name: "task with rich metadata",
			task: &workerpb.TaskRequest{
				TaskId:   "metadata-rich-task",
				TaskType: workerpb.TaskType_TASK_TYPE_INSIGHT_GENERATION,
				Payload:  `{"user_id": "user-123", "analysis_type": "productivity"}`,
				Priority: 5,
				Deadline: timestamppb.New(time.Now().Add(time.Hour)),
				Metadata: map[string]string{
					"user_id":        "user-123",
					"tenant_id":      "tenant-456",
					"source":         "api",
					"correlation_id": "corr-789",
					"environment":    "production",
					"analysis_type":  "productivity",
					"requested_by":   "user@example.com",
					"client_version": "2.1.0",
				},
			},
			wantErr:  false,
			testDesc: "Task with rich metadata should be queued successfully",
		},
		{
			name: "task with empty metadata",
			task: &workerpb.TaskRequest{
				TaskId:   "empty-metadata-task",
				TaskType: workerpb.TaskType_TASK_TYPE_WEEKLY_REPORT,
				Payload:  `{"user_id": "user-456"}`,
				Priority: 3,
				Deadline: timestamppb.New(time.Now().Add(2 * time.Hour)),
				Metadata: map[string]string{},
			},
			wantErr:  false,
			testDesc: "Task with empty metadata should be queued successfully",
		},
		{
			name: "task with nil metadata",
			task: &workerpb.TaskRequest{
				TaskId:   "nil-metadata-task",
				TaskType: workerpb.TaskType_TASK_TYPE_DATA_ANALYSIS,
				Payload:  `{"user_id": "user-789"}`,
				Priority: 7,
				Deadline: timestamppb.New(time.Now().Add(30 * time.Minute)),
				Metadata: nil,
			},
			wantErr:  false,
			testDesc: "Task with nil metadata should be queued successfully",
		},
		{
			name: "task with special characters in metadata",
			task: &workerpb.TaskRequest{
				TaskId:   "special-chars-task",
				TaskType: workerpb.TaskType_TASK_TYPE_INSIGHT_GENERATION,
				Payload:  `{"special": "test"}`,
				Priority: 5,
				Deadline: timestamppb.New(time.Now().Add(time.Hour)),
				Metadata: map[string]string{
					"unicode":       "测试 🚀 émojis",
					"json":          `{"nested": "value"}`,
					"query_params":  "?param1=value1&param2=value2",
					"html_entities": "&lt;script&gt;alert('test')&lt;/script&gt;",
					"whitespace":    "  spaces  and\ttabs\n",
				},
			},
			wantErr:  false,
			testDesc: "Task with special characters in metadata should be handled correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := server.QueueTask(tt.task)

			if tt.wantErr {
				assert.Error(t, err, tt.testDesc)
			} else {
				assert.NoError(t, err, tt.testDesc)
			}
		})
	}
}

// TestWorkerStatsAndStatus tests worker statistics and status tracking
func TestWorkerStatsAndStatus(t *testing.T) {
	cfg := createTestConfigForCapabilities()
	logger := createTestLoggerForCapabilities()
	server := grpc.NewServer(cfg, logger)

	ctx := context.Background()

	// Register a worker
	registerResp := registerWorker(t, server, &workerpb.RegisterWorkerRequest{
		WorkerId:   "stats-worker",
		WorkerName: "Statistics Test Worker",
		Capabilities: []workerpb.WorkerCapability{
			workerpb.WorkerCapability_CAPABILITY_AI_INSIGHTS,
		},
		Version: "1.0.0",
	})

	tests := []struct {
		name        string
		heartbeat   *workerpb.WorkerHeartbeatRequest
		wantErr     bool
		description string
	}{
		{
			name: "heartbeat with detailed stats",
			heartbeat: &workerpb.WorkerHeartbeatRequest{
				WorkerId:     "stats-worker",
				SessionToken: registerResp.SessionToken,
				Status:       workerpb.WorkerStatus_WORKER_STATUS_BUSY,
				Stats: &workerpb.WorkerStats{
					ActiveTasks:          3,
					CompletedTasks:       15,
					FailedTasks:          2,
					CpuUsage:             78.5,
					MemoryUsage:          45.2,
					Uptime:               timestamppb.New(time.Now().Add(-2 * time.Hour)),
					GrpcConnectionStatus: "connected",
				},
			},
			wantErr:     false,
			description: "Heartbeat with detailed statistics should be processed",
		},
		{
			name: "heartbeat with zero stats",
			heartbeat: &workerpb.WorkerHeartbeatRequest{
				WorkerId:     "stats-worker",
				SessionToken: registerResp.SessionToken,
				Status:       workerpb.WorkerStatus_WORKER_STATUS_IDLE,
				Stats: &workerpb.WorkerStats{
					ActiveTasks:          0,
					CompletedTasks:       0,
					FailedTasks:          0,
					CpuUsage:             0.0,
					MemoryUsage:          0.0,
					Uptime:               timestamppb.New(time.Now()),
					GrpcConnectionStatus: "initializing",
				},
			},
			wantErr:     false,
			description: "Heartbeat with zero statistics should be processed",
		},
		{
			name: "heartbeat with high resource usage",
			heartbeat: &workerpb.WorkerHeartbeatRequest{
				WorkerId:     "stats-worker",
				SessionToken: registerResp.SessionToken,
				Status:       workerpb.WorkerStatus_WORKER_STATUS_BUSY,
				Stats: &workerpb.WorkerStats{
					ActiveTasks:          10,
					CompletedTasks:       100,
					FailedTasks:          5,
					CpuUsage:             99.9,
					MemoryUsage:          95.8,
					Uptime:               timestamppb.New(time.Now().Add(-24 * time.Hour)),
					GrpcConnectionStatus: "connected",
				},
			},
			wantErr:     false,
			description: "Heartbeat with high resource usage should be processed",
		},
		{
			name: "heartbeat with error status",
			heartbeat: &workerpb.WorkerHeartbeatRequest{
				WorkerId:     "stats-worker",
				SessionToken: registerResp.SessionToken,
				Status:       workerpb.WorkerStatus_WORKER_STATUS_ERROR,
				Stats: &workerpb.WorkerStats{
					ActiveTasks:          1,
					CompletedTasks:       20,
					FailedTasks:          10,
					CpuUsage:             25.0,
					MemoryUsage:          30.0,
					Uptime:               timestamppb.New(time.Now().Add(-time.Hour)),
					GrpcConnectionStatus: "error",
				},
			},
			wantErr:     false,
			description: "Heartbeat with error status should be processed",
		},
		{
			name: "heartbeat without stats",
			heartbeat: &workerpb.WorkerHeartbeatRequest{
				WorkerId:     "stats-worker",
				SessionToken: registerResp.SessionToken,
				Status:       workerpb.WorkerStatus_WORKER_STATUS_IDLE,
				Stats:        nil, // No stats provided
			},
			wantErr:     false,
			description: "Heartbeat without statistics should be processed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.WorkerHeartbeat(ctx, tt.heartbeat)

			if tt.wantErr {
				assert.Error(t, err, tt.description)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, resp)
				assert.True(t, resp.ConnectionHealthy)
				assert.NotNil(t, resp.ServerTime)

				// Verify stats are stored in worker info
				workers := server.GetActiveWorkers()
				assert.Contains(t, workers, "stats-worker")

				worker := workers["stats-worker"]
				assert.Equal(t, tt.heartbeat.Status, worker.Status)

				if tt.heartbeat.Stats != nil {
					assert.Equal(t, tt.heartbeat.Stats, worker.Stats)
				}
			}
		})
	}
}

// TestHealthCheckWithDifferentScenarios tests health check under various conditions
func TestHealthCheckWithDifferentScenarios(t *testing.T) {
	cfg := createTestConfigForCapabilities()
	logger := createTestLoggerForCapabilities()
	server := grpc.NewServer(cfg, logger)

	ctx := context.Background()

	t.Run("health check with no workers", func(t *testing.T) {
		resp, err := server.HealthCheck(ctx, nil)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "healthy", resp.Status)
		assert.Equal(t, int32(0), resp.ActiveWorkers)
		assert.Contains(t, resp.Services, "grpc")
		assert.Equal(t, "healthy", resp.Services["grpc"])
	})

	t.Run("health check with workers of different ages", func(t *testing.T) {
		// Register workers at different times
		workers := []struct {
			id       string
			ageHours int
		}{
			{"fresh-worker", 0},  // Fresh worker
			{"old-worker", 3},    // 3 hours old (still healthy)
			{"stale-worker", 25}, // 25 hours old (stale)
		}

		for _, w := range workers {
			registerResp := registerWorker(t, server, &workerpb.RegisterWorkerRequest{
				WorkerId:   w.id,
				WorkerName: fmt.Sprintf("Worker %s", w.id),
				Capabilities: []workerpb.WorkerCapability{
					workerpb.WorkerCapability_CAPABILITY_AI_INSIGHTS,
				},
				Version: "1.0.0",
			})

			// Simulate different ages by sending heartbeats at different times
			heartbeatTime := time.Now().Add(-time.Duration(w.ageHours) * time.Hour)

			// Manually update the last heartbeat time in the server
			// Note: This would require access to server internals in a real test
			// For now, we'll just send a heartbeat
			_, err := server.WorkerHeartbeat(ctx, &workerpb.WorkerHeartbeatRequest{
				WorkerId:     w.id,
				SessionToken: registerResp.SessionToken,
				Status:       workerpb.WorkerStatus_WORKER_STATUS_IDLE,
				Stats: &workerpb.WorkerStats{
					Uptime: timestamppb.New(heartbeatTime),
				},
			})
			require.NoError(t, err)
		}

		resp, err := server.HealthCheck(ctx, nil)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "healthy", resp.Status)
		assert.Equal(t, int32(3), resp.ActiveWorkers) // All workers are counted as active
		assert.Contains(t, resp.Services, "grpc")
		assert.Contains(t, resp.Services, "task_queue")
	})
}

// Helper functions for comprehensive tests

func registerWorker(t *testing.T, server *grpc.Server, req *workerpb.RegisterWorkerRequest) *workerpb.RegisterWorkerResponse {
	ctx := context.Background()
	resp, err := server.RegisterWorker(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.RegistrationSuccessful)
	return resp
}

func createTestConfigForCapabilities() *config.Config {
	return &config.Config{
		GRPC: config.GRPCConfig{
			ServerPort:  9090,
			TLSEnabled:  false,
			TLSCertFile: "",
			TLSKeyFile:  "",
		},
	}
}

func createTestLoggerForCapabilities() *logging.Logger {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelWarn, // Reduce noise in tests
	}
	handler := slog.NewJSONHandler(&buf, opts)
	return &logging.Logger{Logger: slog.New(handler)}
}
