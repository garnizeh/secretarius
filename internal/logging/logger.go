package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/garnizeh/englog/internal/config"
)

// Logger wraps slog.Logger with additional convenience methods
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new structured logger based on configuration
func NewLogger(cfg config.LoggingConfig) *Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}

	var output io.Writer = os.Stdout

	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	logger := slog.New(handler)
	return &Logger{Logger: logger}
}

// NewTestLogger creates a logger suitable for testing with minimal output
func NewTestLogger() *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelWarn, // Only show warnings and errors in tests
	}

	handler := slog.NewTextHandler(io.Discard, opts) // Discard output in tests
	logger := slog.New(handler)
	return &Logger{Logger: logger}
}

// WithContext adds context information to logger
func (l *Logger) WithContext(ctx context.Context) *slog.Logger {
	return l.Logger.With("trace_id", getTraceID(ctx))
}

// WithComponent adds component information to logger and returns a new Logger instance
func (l *Logger) WithComponent(component string) *Logger {
	newLogger := l.Logger.With("component", component)
	return &Logger{Logger: newLogger}
}

// WithFields adds custom fields to logger
func (l *Logger) WithFields(fields map[string]any) *slog.Logger {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return l.Logger.With(args...)
}

// HTTP request logging helpers

// LogRequest logs HTTP request information with structured data
func (l *Logger) LogRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration, clientIP string) {
	l.WithContext(ctx).Info("HTTP Request",
		"method", method,
		"path", path,
		"status_code", statusCode,
		"duration_ms", duration.Milliseconds(),
		"client_ip", clientIP,
	)
}

// LogError logs errors with structured context
func (l *Logger) LogError(ctx context.Context, err error, msg string, attrs ...any) {
	allAttrs := append(attrs, "error", err.Error())
	l.WithContext(ctx).Error(msg, allAttrs...)
}

// LogDatabaseOperation logs database operations
func (l *Logger) LogDatabaseOperation(ctx context.Context, operation, table string, duration time.Duration, err error) {
	attrs := []any{
		"operation", operation,
		"table", table,
		"duration_ms", duration.Milliseconds(),
	}

	if err != nil {
		attrs = append(attrs, "error", err.Error())
		l.WithContext(ctx).Error("Database operation failed", attrs...)
	} else {
		l.WithContext(ctx).Debug("Database operation completed", attrs...)
	}
}

// LogServiceOperation logs service layer operations
func (l *Logger) LogServiceOperation(ctx context.Context, service, operation string, userID string, duration time.Duration, err error) {
	attrs := []any{
		"service", service,
		"operation", operation,
		"duration_ms", duration.Milliseconds(),
	}

	if userID != "" {
		attrs = append(attrs, "user_id", userID)
	}

	if err != nil {
		attrs = append(attrs, "error", err.Error())
		l.WithContext(ctx).Error("Service operation failed", attrs...)
	} else {
		l.WithContext(ctx).Info("Service operation completed", attrs...)
	}
}

// LogAuthEvent logs authentication-related events
func (l *Logger) LogAuthEvent(ctx context.Context, event, userID, clientIP string, success bool, details map[string]any) {
	attrs := []any{
		"event", event,
		"user_id", userID,
		"client_ip", clientIP,
		"success", success,
	}

	for k, v := range details {
		attrs = append(attrs, k, v)
	}

	if success {
		l.WithContext(ctx).Info("Authentication event", attrs...)
	} else {
		l.WithContext(ctx).Warn("Authentication event failed", attrs...)
	}
}

// Application lifecycle logging

// LogStartup logs application startup information
func (l *Logger) LogStartup(component string, version string, config map[string]any) {
	attrs := []any{
		"component", component,
		"version", version,
	}

	for k, v := range config {
		attrs = append(attrs, k, v)
	}

	l.Info("Application starting", attrs...)
}

// LogShutdown logs application shutdown information
func (l *Logger) LogShutdown(component string, reason string, graceful bool) {
	l.Info("Application shutting down",
		"component", component,
		"reason", reason,
		"graceful", graceful,
	)
}

// Helper functions

// getTraceID extracts trace ID from context
func getTraceID(ctx context.Context) string {
	if traceID := ctx.Value("trace_id"); traceID != nil {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}

// Performance measurement helpers

// MeasureDuration measures operation duration and logs it
func (l *Logger) MeasureDuration(ctx context.Context, operation string, fn func() error) error {
	start := time.Now()
	err := fn()
	duration := time.Since(start)

	if err != nil {
		l.WithContext(ctx).Error("Operation failed",
			"operation", operation,
			"duration_ms", duration.Milliseconds(),
			"error", err.Error(),
		)
	} else {
		l.WithContext(ctx).Debug("Operation completed",
			"operation", operation,
			"duration_ms", duration.Milliseconds(),
		)
	}

	return err
}
