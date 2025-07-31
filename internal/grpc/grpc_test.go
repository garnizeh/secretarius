package grpc_test

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/logging"
	workerpb "github.com/garnizeh/englog/proto/worker"
	"google.golang.org/grpc/metadata"
)

// MockStream implements the APIWorkerService_StreamTasksServer interface for testing
type MockStream struct {
	ctx       context.Context
	sentTasks []*workerpb.TaskRequest
	mu        sync.Mutex
	sendError error
}

func NewMockStream(ctx context.Context) *MockStream {
	return &MockStream{
		ctx:       ctx,
		sentTasks: make([]*workerpb.TaskRequest, 0),
	}
}

func (m *MockStream) Send(task *workerpb.TaskRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.sendError != nil {
		return m.sendError
	}

	m.sentTasks = append(m.sentTasks, task)
	return nil
}

func (m *MockStream) Context() context.Context {
	return m.ctx
}

func (m *MockStream) GetSentTasks() []*workerpb.TaskRequest {
	m.mu.Lock()
	defer m.mu.Unlock()

	tasks := make([]*workerpb.TaskRequest, len(m.sentTasks))
	copy(tasks, m.sentTasks)
	return tasks
}

func (m *MockStream) SetSendError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sendError = err
}

func (m *MockStream) SendMsg(msg any) error {
	return nil
}

func (m *MockStream) RecvMsg(msg any) error {
	return nil
}

func (m *MockStream) SetHeader(md metadata.MD) error {
	return nil
}

func (m *MockStream) SendHeader(md metadata.MD) error {
	return nil
}

func (m *MockStream) SetTrailer(md metadata.MD) {
}

// Helper functions

// generateEntryIDs generates a slice of entry IDs for testing
func generateEntryIDs(count int) []string {
	entryIDs := make([]string, count)
	for i := 0; i < count; i++ {
		entryIDs[i] = fmt.Sprintf("entry-%d", i+1)
	}
	return entryIDs
}

// createTestLogger creates a test logger (same as in server_test.go)
func createTestLoggerForManager() *logging.Logger {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(&buf, opts)
	return &logging.Logger{Logger: slog.New(handler)}
}

// createTestConfig creates a test configuration (same as in server_test.go)
func createTestConfigForManager() *config.Config {
	return &config.Config{
		GRPC: config.GRPCConfig{
			ServerPort:  getFreePortForManager(), // Use dynamic port allocation
			TLSEnabled:  false,
			TLSCertFile: "",
			TLSKeyFile:  "",
		},
	}
}

func getFreePortForManager() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0 // Let the system choose
	}

	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	return port
}

// createTestConfig creates a test configuration
func createTestConfig() *config.Config {
	return &config.Config{
		GRPC: config.GRPCConfig{
			ServerPort:  9090,
			TLSEnabled:  false,
			TLSCertFile: "",
			TLSKeyFile:  "",
		},
	}
}
