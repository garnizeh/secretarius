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
func NewConnectionManager(logger *logging.Logger, config *ConnectionConfig) *ConnectionManager {
	ctx, cancel := context.WithCancel(context.Background())

	retryConfig := config.RetryConfig
	if retryConfig == nil {
		retryConfig = DefaultRetryConfig()
	}

	logger.Info("Creating new connection manager",
		"component", "connection_manager",
		"target", config.Target,
		"tls_enabled", config.TLSEnabled,
		"health_check_interval", config.HealthCheckInterval.String(),
		"retry_max_attempts", retryConfig.MaxAttempts)

	return &ConnectionManager{
		logger:              logger,
		target:              config.Target,
		tlsEnabled:          config.TLSEnabled,
		certFile:            config.CertFile,
		serverName:          config.ServerName,
		retryConfig:         retryConfig,
		circuitBreaker:      NewCircuitBreaker("grpc_connection", 5, 3, 30*time.Second),
		healthCheckInterval: config.HealthCheckInterval,
		ctx:                 ctx,
		cancel:              cancel,
	}
}

// Connect establishes the gRPC connection with retry logic
func (cm *ConnectionManager) Connect(ctx context.Context) error {
	return RetryOperation(ctx, "grpc_connect", cm.retryConfig, func() error {
		return cm.circuitBreaker.Execute(ctx, func() error {
			return cm.doConnect(ctx)
		})
	})
}

func (cm *ConnectionManager) doConnect(ctx context.Context) error {
	cm.logger.Info("Establishing gRPC connection")

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

	cm.logger.Debug("Configured gRPC dial options",
		"keepalive_time", "10s",
		"keepalive_timeout", "3s",
		"max_recv_msg_size", "4MB",
		"max_send_msg_size", "4MB")

	// Configure TLS or insecure credentials
	if cm.tlsEnabled {
		cm.logger.Debug("Configuring TLS credentials",
			"cert_file", cm.certFile,
			"server_name", cm.serverName)

		if cm.certFile != "" {
			creds, err := credentials.NewClientTLSFromFile(cm.certFile, cm.serverName)
			if err != nil {
				cm.logger.LogError(context.TODO(), err, "Failed to load TLS credentials from file",
					"cert_file", cm.certFile)
				return fmt.Errorf("failed to load TLS credentials: %w", err)
			}
			opts = append(opts, grpc.WithTransportCredentials(creds))
			cm.logger.Debug("TLS credentials loaded from file")
		} else {
			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
			cm.logger.Debug("TLS credentials configured with system root CAs")
		}
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		cm.logger.Debug("Configured insecure credentials")
	}

	cm.logger.Info("Initiating gRPC dial",
		"target", cm.target,
		"timeout", "30s")

	conn, err := grpc.NewClient(cm.target, opts...)
	if err != nil {
		cm.logger.LogError(context.TODO(), err, "Failed to establish gRPC connection",
			"target", cm.target,
			"timeout", "30s")
		return fmt.Errorf("failed to connect to %s: %w", cm.target, err)
	}

	cm.logger.Debug("Connection dial completed, waiting for ready state")

	// Wait for connection to be ready
	ctxReady, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if !conn.WaitForStateChange(ctxReady, connectivity.Connecting) {
		cm.logger.LogError(context.TODO(), fmt.Errorf("connection timeout"), "Connection ready timeout",
			"target", cm.target,
			"timeout", "30s")
		conn.Close()
		return fmt.Errorf("connection timeout to %s", cm.target)
	}

	state := conn.GetState()
	cm.logger.Debug("Connection state checked",
		"target", cm.target,
		"state", state.String())

	if state != connectivity.Ready && state != connectivity.Idle {
		cm.logger.LogError(context.TODO(), fmt.Errorf("invalid state: %s", state), "Connection failed with invalid state",
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

	cm.logger.Info("gRPC connection established successfully",
		"target", cm.target,
		"state", state.String(),
		"last_connect_time", cm.lastConnectTime,
		"reconnect_count", cm.reconnectCount)

	// Start connection monitoring
	go cm.monitorConnection()

	return nil
}

// GetClient returns the gRPC client with connection health check
func (cm *ConnectionManager) GetClient() (workerpb.APIWorkerServiceClient, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	if cm.conn == nil || !cm.connected {
		cm.logger.Warn("GetClient called on disconnected manager",
			"target", cm.target,
			"connected", cm.connected)
		return nil, fmt.Errorf("not connected to %s", cm.target)
	}

	// Check connection state
	state := cm.conn.GetState()
	if state == connectivity.TransientFailure || state == connectivity.Shutdown {
		cm.logger.Warn("GetClient called on unhealthy connection",
			"target", cm.target,
			"state", state.String())
		return nil, fmt.Errorf("connection to %s is unhealthy: %s", cm.target, state)
	}

	cm.logger.Debug("Returning healthy gRPC client",
		"target", cm.target,
		"state", state.String())

	return cm.client, nil
}

// IsConnected returns whether the connection is healthy
func (cm *ConnectionManager) IsConnected() bool {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	connected := cm.connected
	cm.logger.Debug("IsConnected check",
		"target", cm.target,
		"connected", connected)

	return connected
}

// GetConnectionState returns the current connection state
func (cm *ConnectionManager) GetConnectionState() connectivity.State {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	state := cm.connectionState
	cm.logger.Debug("GetConnectionState check",
		"target", cm.target,
		"state", state.String())

	return state
}

// GetStats returns connection statistics
func (cm *ConnectionManager) GetStats() map[string]any {
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

	cm.logger.Debug("Generated connection statistics",
		"target", cm.target,
		"connected", cm.connected,
		"state", cm.connectionState.String(),
		"reconnect_count", cm.reconnectCount)

	return stats
}

// monitorConnection monitors the connection health and triggers reconnection if needed
func (cm *ConnectionManager) monitorConnection() {
	cm.logger.Info("Starting connection monitoring",
		"target", cm.target,
		"health_check_interval", cm.healthCheckInterval)

	ticker := time.NewTicker(cm.healthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-cm.ctx.Done():
			cm.logger.Info("Connection monitoring stopped",
				"target", cm.target,
				"reason", cm.ctx.Err())
			return
		case <-ticker.C:
			if err := cm.checkConnectionHealth(); err != nil {
				cm.logger.Warn("Connection health check failed",
					"target", cm.target,
					"error", err)
				cm.markDisconnected()
			}
		}
	}
}

func (cm *ConnectionManager) checkConnectionHealth() error {
	cm.mutex.RLock()
	conn := cm.conn
	target := cm.target
	cm.mutex.RUnlock()

	cm.logger.Debug("Checking connection health",
		"target", target)

	if conn == nil {
		cm.logger.Warn("Health check failed: no connection available",
			"target", target)
		return fmt.Errorf("no connection available")
	}

	state := conn.GetState()
	cm.mutex.Lock()
	cm.connectionState = state
	cm.mutex.Unlock()

	cm.logger.Debug("Connection state checked",
		"target", target,
		"state", state.String())

	switch state {
	case connectivity.Ready, connectivity.Idle:
		cm.logger.Debug("Connection health check passed",
			"target", target,
			"state", state.String())
		return nil
	case connectivity.Connecting:
		// Still connecting, not necessarily a failure
		cm.logger.Debug("Connection still connecting",
			"target", target)
		return nil
	case connectivity.TransientFailure, connectivity.Shutdown:
		cm.logger.Warn("Connection in unhealthy state",
			"target", target,
			"state", state.String())
		return fmt.Errorf("connection in bad state: %s", state)
	default:
		cm.logger.Warn("Connection in unknown state",
			"target", target,
			"state", state.String())
		return fmt.Errorf("unknown connection state: %s", state)
	}
}

func (cm *ConnectionManager) markDisconnected() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	wasConnected := cm.connected
	cm.connected = false

	cm.logger.Warn("Connection marked as disconnected",
		"target", cm.target,
		"was_connected", wasConnected)
}

// Reconnect attempts to reconnect with exponential backoff
func (cm *ConnectionManager) Reconnect(ctx context.Context) error {
	cm.logger.Info("Attempting to reconnect gRPC connection",
		"target", cm.target)

	cm.markDisconnected()

	err := cm.Connect(ctx)
	if err != nil {
		cm.logger.LogError(ctx, err, "Reconnection failed",
			"target", cm.target)
	} else {
		cm.logger.Info("Reconnection successful",
			"target", cm.target)
	}

	return err
}

// Close closes the connection manager
func (cm *ConnectionManager) Close() error {
	cm.logger.Info("Closing connection manager",
		"target", cm.target)

	cm.cancel()

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.conn != nil {
		cm.logger.Debug("Closing gRPC connection",
			"target", cm.target)
		err := cm.conn.Close()
		cm.conn = nil
		cm.client = nil
		cm.connected = false

		if err != nil {
			cm.logger.LogError(context.TODO(), err, "Error closing gRPC connection",
				"target", cm.target)
		} else {
			cm.logger.Info("gRPC connection closed successfully",
				"target", cm.target)
		}

		return err
	}

	cm.logger.Debug("Connection manager closed (no active connection)",
		"target", cm.target)
	return nil
}

// ExecuteWithRetry executes a gRPC call with retry logic
func (cm *ConnectionManager) ExecuteWithRetry(ctx context.Context, operation string, fn func(client workerpb.APIWorkerServiceClient) error) error {
	return RetryOperation(ctx, operation, cm.retryConfig, func() error {
		client, err := cm.GetClient()
		if err != nil {
			// Try to reconnect if client is not available
			if reconnectErr := cm.Reconnect(ctx); reconnectErr != nil {
				return fmt.Errorf("failed to reconnect: %w (original error: %w)", reconnectErr, err)
			}

			client, err = cm.GetClient()
			if err != nil {
				return fmt.Errorf("client still unavailable after reconnect: %w", err)
			}
		}

		err = fn(client)
		if err != nil {
			// Check if error suggests connection issues
			if cm.isConnectionError(err) {
				cm.markDisconnected()
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
