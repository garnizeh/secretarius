package worker

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/garnizeh/englog/internal/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RetryConfig contains configuration for retry mechanisms
type RetryConfig struct {
	MaxAttempts     int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	BackoffFactor   float64
	JitterEnabled   bool
	RetriableErrors map[codes.Code]bool
}

// DefaultRetryConfig returns a sensible default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:   5,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
		JitterEnabled: true,
		RetriableErrors: map[codes.Code]bool{
			codes.Unavailable:       true,
			codes.DeadlineExceeded:  true,
			codes.ResourceExhausted: true,
			codes.Aborted:           true,
			codes.Internal:          true,
			codes.Unknown:           true,
		},
	}
}

// IsRetriableError checks if a gRPC error is retriable
func (rc *RetryConfig) IsRetriableError(err error) bool {
	if err == nil {
		return false
	}

	st, ok := status.FromError(err)
	if !ok {
		// Non-gRPC errors are generally retriable
		return true
	}

	return rc.RetriableErrors[st.Code()]
}

// CalculateDelay calculates the delay for the next retry attempt
func (rc *RetryConfig) CalculateDelay(attempt int) time.Duration {
	if attempt <= 0 {
		return rc.InitialDelay
	}

	// Exponential backoff
	delay := float64(rc.InitialDelay) * math.Pow(rc.BackoffFactor, float64(attempt-1))

	// Cap at max delay
	if delay > float64(rc.MaxDelay) {
		delay = float64(rc.MaxDelay)
	}

	// Add jitter to prevent thundering herd
	if rc.JitterEnabled {
		jitter := delay * 0.1 * (2.0*rand.Float64() - 1.0)
		delay += jitter
	}

	return time.Duration(delay)
}

// RetryOperation executes an operation with retry logic
func RetryOperation(ctx context.Context, logger *logging.Logger, operation string, config *RetryConfig, fn func() error) error {
	var lastError error

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation %s cancelled: %w", operation, ctx.Err())
		default:
		}

		err := fn()
		if err == nil {
			if attempt > 1 {
				logger.LogInfo(ctx, "Operation succeeded after retry",
					logging.OperationField, operation,
					"attempt", attempt,
					"total_attempts", config.MaxAttempts)
			}
			return nil
		}

		lastError = err

		if !config.IsRetriableError(err) {
			logger.LogError(ctx, err, "Operation failed with non-retriable error",
				logging.OperationField, operation,
				"attempt", attempt)
			return fmt.Errorf("operation %s failed (non-retriable): %w", operation, err)
		}

		if attempt == config.MaxAttempts {
			logger.LogError(ctx, err, "Operation failed after all retry attempts",
				logging.OperationField, operation,
				"attempts", config.MaxAttempts)
			break
		}

		delay := config.CalculateDelay(attempt)
		logger.LogWarn(ctx, "Operation failed, retrying",
			logging.OperationField, operation,
			"attempt", attempt,
			"max_attempts", config.MaxAttempts,
			"delay", delay,
			logging.ErrorField, err)

		select {
		case <-ctx.Done():
			return fmt.Errorf("operation %s cancelled during retry delay: %w", operation, ctx.Err())
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return fmt.Errorf("operation %s failed after %d attempts: %w", operation, config.MaxAttempts, lastError)
}

// CircuitBreakerState represents the state of a circuit breaker
type CircuitBreakerState int

const (
	CircuitBreakerClosed CircuitBreakerState = iota
	CircuitBreakerOpen
	CircuitBreakerHalfOpen
)

// CircuitBreaker implements circuit breaker pattern for fault tolerance
type CircuitBreaker struct {
	logger               *logging.Logger
	name                 string
	failureThreshold     int
	successThreshold     int
	timeout              time.Duration
	consecutiveFailures  int
	consecutiveSuccesses int
	state                CircuitBreakerState
	lastFailureTime      time.Time
	mutex                *sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(logger *logging.Logger, name string, failureThreshold, successThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		logger:           logger,
		name:             name,
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		timeout:          timeout,
		state:            CircuitBreakerClosed,
		mutex:            &sync.RWMutex{},
	}
}

// Execute executes an operation through the circuit breaker
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Check if circuit breaker should transition from Open to Half-Open
	if cb.state == CircuitBreakerOpen {
		if time.Since(cb.lastFailureTime) >= cb.timeout {
			cb.state = CircuitBreakerHalfOpen
			cb.consecutiveSuccesses = 0
			cb.logger.LogInfo(ctx, "Circuit breaker transitioning to half-open",
				logging.OperationField, "circuit_breaker_transition",
				"name", cb.name,
				"timeout_elapsed", time.Since(cb.lastFailureTime))
		} else {
			return fmt.Errorf("circuit breaker %s is open", cb.name)
		}
	}

	// Execute the operation
	err := fn()

	if err != nil {
		cb.onFailure(ctx)
		return err
	}

	cb.onSuccess(ctx)
	return nil
}

func (cb *CircuitBreaker) onSuccess(ctx context.Context) {
	switch cb.state {
	case CircuitBreakerClosed:
		cb.consecutiveFailures = 0
	case CircuitBreakerHalfOpen:
		cb.consecutiveSuccesses++
		if cb.consecutiveSuccesses >= cb.successThreshold {
			cb.state = CircuitBreakerClosed
			cb.consecutiveFailures = 0
			cb.logger.LogInfo(ctx, "Circuit breaker closed after successful recovery",
				logging.OperationField, "circuit_breaker_recovery",
				"name", cb.name,
				"consecutive_successes", cb.consecutiveSuccesses)
		}
	}
}

func (cb *CircuitBreaker) onFailure(ctx context.Context) {
	cb.lastFailureTime = time.Now()

	switch cb.state {
	case CircuitBreakerClosed:
		cb.consecutiveFailures++
		if cb.consecutiveFailures >= cb.failureThreshold {
			cb.state = CircuitBreakerOpen
			cb.logger.LogWarn(ctx, "Circuit breaker opened due to consecutive failures",
				logging.OperationField, "circuit_breaker_failure",
				"name", cb.name,
				"consecutive_failures", cb.consecutiveFailures,
				"threshold", cb.failureThreshold)
		}
	case CircuitBreakerHalfOpen:
		cb.state = CircuitBreakerOpen
		cb.consecutiveSuccesses = 0
		cb.logger.LogWarn(ctx, "Circuit breaker opened from half-open state",
			logging.OperationField, "circuit_breaker_failure",
			"name", cb.name)
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// GetStats returns current statistics of the circuit breaker
func (cb *CircuitBreaker) GetStats() map[string]any {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	return map[string]any{
		"name":                  cb.name,
		"state":                 cb.state,
		"consecutive_failures":  cb.consecutiveFailures,
		"consecutive_successes": cb.consecutiveSuccesses,
		"last_failure_time":     cb.lastFailureTime,
		"failure_threshold":     cb.failureThreshold,
		"success_threshold":     cb.successThreshold,
	}
}
