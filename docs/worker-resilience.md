# Error Handling, Retry Mechanisms, and Connection Resilience

## Overview

This document describes the comprehensive error handling, retry mechanisms, and connection resilience implemented in the Worker Client gRPC communication system. The implementation provides robust fault tolerance for distributed operations.

## Core Components

### 1. Retry Configuration (`internal/worker/retry.go`)

#### RetryConfig Structure
```go
type RetryConfig struct {
    MaxAttempts     int                    // Maximum retry attempts
    InitialDelay    time.Duration         // Initial delay before first retry
    MaxDelay        time.Duration         // Maximum delay cap
    BackoffFactor   float64               // Exponential backoff multiplier
    JitterEnabled   bool                  // Add randomness to prevent thundering herd
    RetriableErrors map[codes.Code]bool   // gRPC error codes that should trigger retries
}
```

#### Default Configuration
- **Max Attempts**: 5 retries
- **Initial Delay**: 100ms
- **Max Delay**: 30 seconds
- **Backoff Factor**: 2.0 (exponential)
- **Jitter**: Enabled (±10% randomization)
- **Retriable Errors**: Unavailable, DeadlineExceeded, ResourceExhausted, Aborted, Internal, Unknown

#### Retry Operation Function
```go
func RetryOperation(ctx context.Context, operation string, config *RetryConfig, fn func() error) error
```

**Features:**
- Context-aware cancellation
- Exponential backoff with jitter
- Configurable retriable error detection
- Comprehensive logging of retry attempts
- Operation-specific naming for debugging

### 2. Circuit Breaker Pattern (`internal/worker/retry.go`)

#### Circuit Breaker States
- **Closed**: Normal operation, requests pass through
- **Open**: Failures exceeded threshold, requests fail immediately
- **Half-Open**: Testing recovery, limited requests allowed

#### Configuration
```go
type CircuitBreaker struct {
    name                 string
    failureThreshold     int           // Failures to open circuit
    successThreshold     int           // Successes to close circuit
    timeout              time.Duration // Time before half-open attempt
    consecutiveFailures  int
    consecutiveSuccesses int
    state                CircuitBreakerState
    lastFailureTime      time.Time
}
```

**Default Settings:**
- **Failure Threshold**: 5 consecutive failures
- **Success Threshold**: 3 consecutive successes
- **Timeout**: 30 seconds before half-open

### 3. Connection Management (`internal/worker/connection.go`)

#### ConnectionManager Features
- **Automatic Reconnection**: Handles connection drops gracefully
- **Health Monitoring**: Continuous connection state checking
- **TLS Support**: Configurable secure/insecure connections
- **Keep-Alive**: Maintains persistent connections
- **Statistics**: Detailed connection metrics

#### Key Methods
```go
func (cm *ConnectionManager) Connect(ctx context.Context) error
func (cm *ConnectionManager) Reconnect(ctx context.Context) error
func (cm *ConnectionManager) ExecuteWithRetry(ctx context.Context, operation string, fn func(client workerpb.APIWorkerServiceClient) error) error
func (cm *ConnectionManager) GetClient() (workerpb.APIWorkerServiceClient, error)
```

#### Connection Configuration
```go
type ConnectionConfig struct {
    Target              string
    TLSEnabled          bool
    CertFile            string
    ServerName          string
    HealthCheckInterval time.Duration
    RetryConfig         *RetryConfig
}
```

### 4. Enhanced Worker Client (`internal/worker/client.go`)

#### Resilience Features

**Registration with Retry:**
```go
func (c *Client) registerWorkerWithRetry(ctx context.Context) error
```
- Circuit breaker protection
- Exponential backoff retry
- Session token management
- Connection state tracking

**Heartbeat with Resilience:**
```go
func (c *Client) sendHeartbeatWithRetry(ctx context.Context) error
```
- Regular health reporting
- Automatic reconnection on failure
- Statistics transmission
- Server health monitoring

**Task Streaming with Recovery:**
```go
func (c *Client) streamTasksWithRetry(ctx context.Context) error
```
- Stream recovery on disconnection
- Backoff delays between retries
- Concurrent task processing with semaphore
- Graceful stream closure

**Task Processing with Timeout:**
```go
func (c *Client) processTaskWithErrorHandling(ctx context.Context, task *workerpb.TaskRequest) error
```
- Individual task timeouts (5 minutes)
- Circuit breaker protection
- Progress reporting with retry
- Result reporting with retry

#### Background Routines

**Connection Health Monitoring:**
```go
func (c *Client) connectionHealthRoutine(ctx context.Context)
```
- Monitors connection manager state
- Triggers reconnection when needed
- Re-registers worker after reconnection
- 10-second monitoring interval

**Statistics Reporting:**
```go
func (c *Client) statisticsRoutine(ctx context.Context)
```
- Periodic statistics logging (5 minutes)
- Connection metrics
- Circuit breaker states
- Task processing statistics

### 5. AI Service Resilience (`internal/ai/ollama.go`)

#### Enhanced AI Operations

**Insight Generation with Retry:**
```go
func (s *OllamaService) GenerateInsight(ctx context.Context, prompt string) (*Insight, error)
```
- Input validation
- 3 retry attempts with 1-second base delay
- 60-second timeout per attempt
- Confidence scoring
- Context cancellation support

**Weekly Report Generation with Retry:**
```go
func (s *OllamaService) GenerateWeeklyReport(ctx context.Context, req *WeeklyReportRequest) (*WeeklyReport, error)
```
- Request validation
- 3 retry attempts with 2-second base delay
- 90-second timeout per attempt
- Intelligent content extraction
- Structured report generation

**Health Check with Retry:**
```go
func (s *OllamaService) HealthCheck(ctx context.Context) error
```
- 3 retry attempts with 500ms base delay
- Fast failure detection
- Service availability verification

## Error Categories and Handling

### 1. Network Errors
**Types:** Connection timeouts, DNS failures, network unreachable
**Handling:**
- Automatic retry with exponential backoff
- Connection manager reconnection
- Circuit breaker protection
- Health check validation

### 2. gRPC Service Errors
**Types:** Unavailable, DeadlineExceeded, ResourceExhausted
**Handling:**
- Configurable retry based on error codes
- Session re-establishment
- Stream recreation
- Service discovery retry

### 3. Application Errors
**Types:** Invalid requests, processing failures, business logic errors
**Handling:**
- Input validation
- Graceful error reporting
- Task failure handling
- Statistical tracking

### 4. Resource Exhaustion
**Types:** Memory pressure, CPU overload, connection limits
**Handling:**
- Semaphore-based concurrency control
- Task timeout enforcement
- Connection pooling
- Resource monitoring

## Configuration Examples

### Development Configuration
```go
retryConfig := &RetryConfig{
    MaxAttempts:   3,
    InitialDelay:  50 * time.Millisecond,
    MaxDelay:      5 * time.Second,
    BackoffFactor: 1.5,
    JitterEnabled: true,
}

connectionConfig := &ConnectionConfig{
    Target:              "localhost:50051",
    TLSEnabled:          false,
    HealthCheckInterval: 10 * time.Second,
    RetryConfig:         retryConfig,
}
```

### Production Configuration
```go
retryConfig := &RetryConfig{
    MaxAttempts:   5,
    InitialDelay:  100 * time.Millisecond,
    MaxDelay:      30 * time.Second,
    BackoffFactor: 2.0,
    JitterEnabled: true,
}

connectionConfig := &ConnectionConfig{
    Target:              "api-server.production.com:50051",
    TLSEnabled:          true,
    CertFile:            "/certs/server.crt",
    ServerName:          "api-server.production.com",
    HealthCheckInterval: 30 * time.Second,
    RetryConfig:         retryConfig,
}
```

## Monitoring and Observability

### Metrics Collected
- **Connection Statistics**: State, reconnection count, last connect time
- **Retry Statistics**: Attempts, success rate, failure patterns
- **Circuit Breaker Statistics**: State transitions, failure/success counts
- **Task Statistics**: Active, completed, failed task counts
- **Performance Metrics**: Processing times, queue sizes, resource usage

### Logging Levels
- **Info**: Successful operations, state changes, periodic statistics
- **Warn**: Retry attempts, circuit breaker state changes, recoverable failures
- **Error**: Non-recoverable failures, service unavailability, critical errors

### Health Check Endpoints
- **`/health`**: Overall service health including gRPC connection and AI service
- **`/readiness`**: Service readiness including connection state and capacity

## Testing Resilience

### Failure Scenarios Covered
1. **Network Partitions**: Connection drops, DNS failures
2. **Service Unavailability**: API server restarts, maintenance windows
3. **Resource Exhaustion**: High load, memory pressure
4. **Timeout Conditions**: Slow responses, processing delays
5. **Invalid Requests**: Malformed data, authentication failures

### Testing Commands
```bash
# Test worker with connection failures
./bin/worker # With API server stopped

# Test with network delays
tc qdisc add dev lo root netem delay 100ms

# Test with packet loss
tc qdisc add dev lo root netem loss 10%

# Monitor resilience metrics
curl http://localhost:8081/health
curl http://localhost:8081/readiness
```

## Best Practices

### 1. Context Usage
- Always pass context through the call chain
- Set appropriate timeouts for operations
- Handle context cancellation gracefully
- Use context for request tracing

### 2. Error Handling
- Distinguish between retriable and non-retriable errors
- Log errors with appropriate levels and context
- Provide meaningful error messages
- Track error patterns for debugging

### 3. Resource Management
- Use semaphores for concurrency control
- Implement proper cleanup in defer statements
- Monitor resource usage and limits
- Set reasonable timeouts for operations

### 4. Configuration
- Make retry parameters configurable
- Use environment-specific settings
- Validate configuration at startup
- Document configuration options

## Future Enhancements

### Planned Improvements
1. **Adaptive Retry**: Dynamic adjustment based on success rates
2. **Load Balancing**: Multiple API server endpoints
3. **Circuit Breaker Metrics**: Prometheus integration
4. **Distributed Tracing**: OpenTelemetry support
5. **Chaos Engineering**: Built-in failure injection

### Monitoring Integration
- Metrics export to Prometheus
- Distributed tracing with Jaeger
- Log aggregation with ELK stack
- Alerting on failure thresholds

This comprehensive error handling and resilience implementation ensures that the Worker Client can operate reliably in distributed environments with various failure modes and network conditions.
