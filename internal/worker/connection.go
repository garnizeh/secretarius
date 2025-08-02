package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/garnizeh/englog/internal/logging"
	workerpb "github.com/garnizeh/englog/proto/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

// ConnectionManager manages gRPC connection with automatic reconnection
type ConnectionManager struct {
	logger     *logging.Logger
	target     string
	tlsEnabled bool
	certFile   string
	serverName string

	conn   *grpc.ClientConn
	client workerpb.APIWorkerServiceClient
	mutex  sync.RWMutex

	retryConfig    *RetryConfig
	circuitBreaker *CircuitBreaker

	// Connection monitoring
	connected       bool
	connectionState connectivity.State
	lastConnectTime time.Time
	reconnectCount  int

	// Health monitoring
	healthCheckInterval time.Duration
	ctx                 context.Context
	cancel              context.CancelFunc
}

// ConnectionConfig contains configuration for connection management
type ConnectionConfig struct {
	Target              string
	TLSEnabled          bool
	CertFile            string
	ServerName          string
	HealthCheckInterval time.Duration
	RetryConfig         *RetryConfig
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(ctx context.Context, logger *logging.Logger, config *ConnectionConfig) *ConnectionManager {
	ctx, cancel := context.WithCancel(ctx)

	retryConfig := config.RetryConfig
	if retryConfig == nil {
		retryConfig = DefaultRetryConfig()
	}

	baseLogger := logger.WithServiceAndComponent("worker", "connection_manager")
	baseLogger.LogInfo(ctx, "Creating new connection manager",
		logging.OperationField, "new_connection_manager",
		"target", config.Target,
		"tls_enabled", config.TLSEnabled,
		"health_check_interval", config.HealthCheckInterval.String(),
		"retry_max_attempts", retryConfig.MaxAttempts)

	return &ConnectionManager{
		logger:              baseLogger,
		target:              config.Target,
		tlsEnabled:          config.TLSEnabled,
		certFile:            config.CertFile,
		serverName:          config.ServerName,
		retryConfig:         retryConfig,
		circuitBreaker:      NewCircuitBreaker(baseLogger, "grpc_connection", 5, 3, 30*time.Second),
		healthCheckInterval: config.HealthCheckInterval,
		ctx:                 ctx,
		cancel:              cancel,
	}
}

// Connect establishes the gRPC connection with retry logic
func (cm *ConnectionManager) Connect(ctx context.Context) error {
	return RetryOperation(ctx, cm.logger, "grpc_connect", cm.retryConfig, func() error {
		return cm.circuitBreaker.Execute(ctx, func() error {
			return cm.doConnect(ctx)
		})
	})
}

func (cm *ConnectionManager) doConnect(ctx context.Context) error {
	cm.logger.LogInfo(ctx, "Establishing gRPC connection",
		logging.OperationField, "do_connect")

	// Configure connection options
	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(4*1024*1024), // 4MB
			grpc.MaxCallSendMsgSize(4*1024*1024), // 4MB
		),
	}

	cm.logger.LogDebug(ctx, "Configured gRPC dial options",
		logging.OperationField, "do_connect",
		"keepalive_time", "10s",
		"keepalive_timeout", "3s",
		"max_recv_msg_size", "4MB",
		"max_send_msg_size", "4MB")

	// Configure TLS or insecure credentials
	if cm.tlsEnabled {
		cm.logger.LogDebug(ctx, "Configuring TLS credentials",
			logging.OperationField, "do_connect",
			"cert_file", cm.certFile,
			"server_name", cm.serverName)

		if cm.certFile != "" {
			creds, err := credentials.NewClientTLSFromFile(cm.certFile, cm.serverName)
			if err != nil {
				cm.logger.LogError(ctx, err, "Failed to load TLS credentials from file",
					logging.OperationField, "do_connect",
					"cert_file", cm.certFile)
				return fmt.Errorf("failed to load TLS credentials: %w", err)
			}
			opts = append(opts, grpc.WithTransportCredentials(creds))
			cm.logger.LogDebug(ctx, "TLS credentials loaded from file",
				logging.OperationField, "do_connect")
		} else {
			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
			cm.logger.LogDebug(ctx, "TLS credentials configured with system root CAs",
				logging.OperationField, "do_connect")
		}
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		cm.logger.LogDebug(ctx, "Configured insecure credentials",
			logging.OperationField, "do_connect")
	}

	cm.logger.LogInfo(ctx, "Initiating gRPC dial",
		logging.OperationField, "do_connect",
		"target", cm.target,
		"timeout", "30s")

	conn, err := grpc.NewClient(cm.target, opts...)
	if err != nil {
		cm.logger.LogError(ctx, err, "Failed to establish gRPC connection",
			logging.OperationField, "do_connect",
			"target", cm.target,
			"timeout", "30s")
		return fmt.Errorf("failed to connect to %s: %w", cm.target, err)
	}

	cm.logger.LogDebug(ctx, "Connection dial completed, waiting for ready state",
		logging.OperationField, "do_connect")

	// Wait for connection to be ready
	ctxReady, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if !conn.WaitForStateChange(ctxReady, connectivity.Connecting) {
		cm.logger.LogError(ctx, fmt.Errorf("connection timeout"), "Connection ready timeout",
			logging.OperationField, "do_connect",
			"target", cm.target,
			"timeout", "30s")
		conn.Close()
		return fmt.Errorf("connection timeout to %s", cm.target)
	}

	state := conn.GetState()
	cm.logger.LogDebug(ctx, "Connection state checked",
		logging.OperationField, "do_connect",
		"target", cm.target,
		"state", state.String())

	if state != connectivity.Ready && state != connectivity.Idle {
		cm.logger.LogError(ctx, fmt.Errorf("invalid state: %s", state), "Connection failed with invalid state",
			logging.OperationField, "do_connect",
			"target", cm.target,
			"state", state.String())
		conn.Close()
		return fmt.Errorf("connection failed, state: %s", state)
	}

	// Update connection state
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Close old connection if exists
	if cm.conn != nil {
		cm.conn.Close()
	}

	cm.conn = conn
	cm.client = workerpb.NewAPIWorkerServiceClient(conn)
	cm.connected = true
	cm.connectionState = state
	cm.lastConnectTime = time.Now()
	cm.reconnectCount++

	cm.logger.LogInfo(ctx, "gRPC connection established successfully",
		logging.OperationField, "do_connect",
		"target", cm.target,
		"state", state.String(),
		"last_connect_time", cm.lastConnectTime,
		"reconnect_count", cm.reconnectCount)

	// Start connection monitoring
	go cm.monitorConnection(ctx)

	return nil
}

// GetClient returns the gRPC client with connection health check
func (cm *ConnectionManager) GetClient(ctx context.Context) (workerpb.APIWorkerServiceClient, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	if cm.conn == nil || !cm.connected {
		cm.logger.LogWarn(ctx, "GetClient called on disconnected manager",
			logging.OperationField, "get_client",
			"target", cm.target,
			"connected", cm.connected)
		return nil, fmt.Errorf("not connected to %s", cm.target)
	}

	// Check connection state
	state := cm.conn.GetState()
	if state == connectivity.TransientFailure || state == connectivity.Shutdown {
		cm.logger.LogWarn(ctx, "GetClient called on unhealthy connection",
			logging.OperationField, "get_client",
			"target", cm.target,
			"state", state.String())
		return nil, fmt.Errorf("connection to %s is unhealthy: %s", cm.target, state)
	}

	cm.logger.LogDebug(ctx, "Returning healthy gRPC client",
		logging.OperationField, "get_client",
		"target", cm.target,
		"state", state.String())

	return cm.client, nil
}

// IsConnected returns whether the connection is healthy
func (cm *ConnectionManager) IsConnected(ctx context.Context) bool {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	connected := cm.connected
	cm.logger.LogDebug(ctx, "IsConnected check",
		logging.OperationField, "is_connected",
		"target", cm.target,
		"connected", connected)

	return connected
}

// GetConnectionState returns the current connection state
func (cm *ConnectionManager) GetConnectionState(ctx context.Context) connectivity.State {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	state := cm.connectionState
	cm.logger.LogDebug(ctx, "GetConnectionState check",
		logging.OperationField, "get_connection_state",
		"target", cm.target,
		"state", state.String())

	return state
}

// GetStats returns connection statistics
func (cm *ConnectionManager) GetStats(ctx context.Context) map[string]any {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	stats := map[string]any{
		"connected":         cm.connected,
		"connection_state":  cm.connectionState.String(),
		"last_connect_time": cm.lastConnectTime,
		"reconnect_count":   cm.reconnectCount,
		"target":            cm.target,
		"tls_enabled":       cm.tlsEnabled,
	}

	// Add circuit breaker stats
	cbStats := cm.circuitBreaker.GetStats()
	for k, v := range cbStats {
		stats["circuit_breaker_"+k] = v
	}

	cm.logger.LogDebug(ctx, "Generated connection statistics",
		logging.OperationField, "get_stats",
		"target", cm.target,
		"connected", cm.connected,
		"state", cm.connectionState.String(),
		"reconnect_count", cm.reconnectCount)

	return stats
}

// monitorConnection monitors the connection health and triggers reconnection if needed
func (cm *ConnectionManager) monitorConnection(ctx context.Context) {
	cm.logger.LogInfo(ctx, "Starting connection monitoring",
		logging.OperationField, "monitor_connection",
		"target", cm.target,
		"health_check_interval", cm.healthCheckInterval)

	ticker := time.NewTicker(cm.healthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-cm.ctx.Done():
			cm.logger.LogInfo(ctx, "Connection monitoring stopped",
				logging.OperationField, "monitor_connection",
				"target", cm.target,
				"reason", cm.ctx.Err())
			return
		case <-ticker.C:
			if err := cm.checkConnectionHealth(ctx); err != nil {
				cm.logger.LogWarn(ctx, "Connection health check failed",
					logging.OperationField, "monitor_connection",
					"target", cm.target,
					logging.ErrorField, err)
				cm.markDisconnected(ctx)
			}
		}
	}
}

func (cm *ConnectionManager) checkConnectionHealth(ctx context.Context) error {
	cm.mutex.RLock()
	conn := cm.conn
	target := cm.target
	cm.mutex.RUnlock()

	cm.logger.LogDebug(ctx, "Checking connection health",
		logging.OperationField, "check_connection_health",
		"target", target)

	if conn == nil {
		cm.logger.LogWarn(ctx, "Health check failed: no connection available",
			logging.OperationField, "check_connection_health",
			"target", target)
		return fmt.Errorf("no connection available")
	}

	state := conn.GetState()
	cm.mutex.Lock()
	cm.connectionState = state
	cm.mutex.Unlock()

	cm.logger.LogDebug(ctx, "Connection state checked",
		logging.OperationField, "check_connection_health",
		"target", target,
		"state", state.String())

	switch state {
	case connectivity.Ready, connectivity.Idle:
		cm.logger.LogDebug(ctx, "Connection health check passed",
			logging.OperationField, "check_connection_health",
			"target", target,
			"state", state.String())
		return nil
	case connectivity.Connecting:
		// Still connecting, not necessarily a failure
		cm.logger.LogDebug(ctx, "Connection still connecting",
			logging.OperationField, "check_connection_health",
			"target", target)
		return nil
	case connectivity.TransientFailure, connectivity.Shutdown:
		cm.logger.LogWarn(ctx, "Connection in unhealthy state",
			logging.OperationField, "check_connection_health",
			"target", target,
			"state", state.String())
		return fmt.Errorf("connection in bad state: %s", state)
	default:
		cm.logger.LogWarn(ctx, "Connection in unknown state",
			logging.OperationField, "check_connection_health",
			"target", target,
			"state", state.String())
		return fmt.Errorf("unknown connection state: %s", state)
	}
}

func (cm *ConnectionManager) markDisconnected(ctx context.Context) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	wasConnected := cm.connected
	cm.connected = false

	cm.logger.LogWarn(ctx, "Connection marked as disconnected",
		logging.OperationField, "mark_disconnected",
		"target", cm.target,
		"was_connected", wasConnected)
}

// Reconnect attempts to reconnect with exponential backoff
func (cm *ConnectionManager) Reconnect(ctx context.Context) error {
	cm.logger.LogInfo(ctx, "Attempting to reconnect gRPC connection",
		logging.OperationField, "reconnect",
		"target", cm.target)

	cm.markDisconnected(ctx)

	err := cm.Connect(ctx)
	if err != nil {
		cm.logger.LogError(ctx, err, "Reconnection failed",
			logging.OperationField, "reconnect",
			"target", cm.target)
	} else {
		cm.logger.LogInfo(ctx, "Reconnection successful",
			logging.OperationField, "reconnect",
			"target", cm.target)
	}

	return err
}

// Close closes the connection manager
func (cm *ConnectionManager) Close(ctx context.Context) error {
	cm.logger.LogInfo(ctx, "Closing connection manager",
		logging.OperationField, "close",
		"target", cm.target)

	cm.cancel()

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.conn != nil {
		cm.logger.LogDebug(ctx, "Closing gRPC connection",
			logging.OperationField, "close",
			"target", cm.target)
		err := cm.conn.Close()
		cm.conn = nil
		cm.client = nil
		cm.connected = false

		if err != nil {
			cm.logger.LogError(ctx, err, "Error closing gRPC connection",
				logging.OperationField, "close",
				"target", cm.target)
		} else {
			cm.logger.LogInfo(ctx, "gRPC connection closed successfully",
				logging.OperationField, "close",
				"target", cm.target)
		}

		return err
	}

	cm.logger.LogDebug(ctx, "Connection manager closed (no active connection)",
		logging.OperationField, "close",
		"target", cm.target)
	return nil
}

// ExecuteWithRetry executes a gRPC call with retry logic
func (cm *ConnectionManager) ExecuteWithRetry(ctx context.Context, operation string, fn func(client workerpb.APIWorkerServiceClient) error) error {
	return RetryOperation(ctx, cm.logger, operation, cm.retryConfig, func() error {
		client, err := cm.GetClient(ctx)
		if err != nil {
			// Try to reconnect if client is not available
			if reconnectErr := cm.Reconnect(ctx); reconnectErr != nil {
				return fmt.Errorf("failed to reconnect: %w (original error: %w)", reconnectErr, err)
			}

			client, err = cm.GetClient(ctx)
			if err != nil {
				return fmt.Errorf("client still unavailable after reconnect: %w", err)
			}
		}

		err = fn(client)
		if err != nil {
			// Check if error suggests connection issues
			if cm.isConnectionError(err) {
				cm.markDisconnected(ctx)
			}
		}

		return err
	})
}

func (cm *ConnectionManager) isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	switch st.Code() {
	case codes.Unavailable, codes.DeadlineExceeded, codes.Canceled:
		return true
	default:
		return false
	}
}
