package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/garnizeh/englog/internal/config"
)

// Define a custom type for context keys to avoid collisions
type contextKey string

const traceIDKey contextKey = "trace_id"

// Standard field names for consistent logging
const (
	ServiceField    = "service"
	ComponentField  = "component"
	OperationField  = "operation"
	UserIDField     = "user_id"
	EmailField      = "email"
	ErrorField      = "error"
	DurationField   = "duration_ms"
	StatusField     = "status"
	MethodField     = "method"
	PathField       = "path"
	StatusCodeField = "status_code"
	ClientIPField   = "client_ip"
	TraceIDField    = "trace_id"
	TableField      = "table"
	EventField      = "event"
	SuccessField    = "success"
	VersionField    = "version"
	ReasonField     = "reason"
	GracefulField   = "graceful"
)

// Logger wraps slog.Logger with additional convenience methods and standardized field names
type Logger struct {
	*slog.Logger
	service   string
	component string
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

// WithService creates a new logger instance with service identification
func (l *Logger) WithService(service string) *Logger {
	newLogger := l.Logger.With(ServiceField, service)
	return &Logger{
		Logger:  newLogger,
		service: service,
	}
}

// WithComponent creates a new logger instance with component identification
func (l *Logger) WithComponent(component string) *Logger {
	newLogger := l.Logger.With(ComponentField, component)
	return &Logger{
		Logger:    newLogger,
		service:   l.service,
		component: component,
	}
}

// WithServiceAndComponent creates a new logger instance with both service and component
func (l *Logger) WithServiceAndComponent(service, component string) *Logger {
	newLogger := l.Logger.With(ServiceField, service, ComponentField, component)
	return &Logger{
		Logger:    newLogger,
		service:   service,
		component: component,
	}
}

// WithContext adds context information to logger
func (l *Logger) WithContext(ctx context.Context) *slog.Logger {
	return l.Logger.With(TraceIDField, getTraceID(ctx))
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
		MethodField, method,
		PathField, path,
		StatusCodeField, statusCode,
		DurationField, duration.Milliseconds(),
		ClientIPField, clientIP,
	)
}

// LogError logs errors with structured context
func (l *Logger) LogError(ctx context.Context, err error, msg string, attrs ...any) {
	allAttrs := append(attrs, ErrorField, err.Error())
	l.WithContext(ctx).Error(msg, allAttrs...)
}

// LogOperation logs a generic operation with standard fields
func (l *Logger) LogOperation(ctx context.Context, operation string, userID string, duration time.Duration, err error, attrs ...any) {
	logAttrs := []any{
		OperationField, operation,
		DurationField, duration.Milliseconds(),
	}

	if userID != "" {
		logAttrs = append(logAttrs, UserIDField, userID)
	}

	// Add custom attributes
	logAttrs = append(logAttrs, attrs...)

	if err != nil {
		logAttrs = append(logAttrs, ErrorField, err.Error())
		l.WithContext(ctx).Error("Operation failed", logAttrs...)
	} else {
		l.WithContext(ctx).Info("Operation completed", logAttrs...)
	}
}

// LogUserOperation logs user-related operations with standard fields
func (l *Logger) LogUserOperation(ctx context.Context, operation string, userID string, email string, success bool, attrs ...any) {
	logAttrs := []any{
		OperationField, operation,
		UserIDField, userID,
		SuccessField, success,
	}

	if email != "" {
		logAttrs = append(logAttrs, EmailField, email)
	}

	// Add custom attributes
	logAttrs = append(logAttrs, attrs...)

	if success {
		l.WithContext(ctx).Info("User operation completed", logAttrs...)
	} else {
		l.WithContext(ctx).Warn("User operation failed", logAttrs...)
	}
}

// LogDatabaseOperation logs database operations
func (l *Logger) LogDatabaseOperation(ctx context.Context, operation, table string, duration time.Duration, err error) {
	attrs := []any{
		OperationField, operation,
		TableField, table,
		DurationField, duration.Milliseconds(),
	}

	if err != nil {
		attrs = append(attrs, ErrorField, err.Error())
		l.WithContext(ctx).Error("Database operation failed", attrs...)
	} else {
		l.WithContext(ctx).Debug("Database operation completed", attrs...)
	}
}

// LogServiceOperation logs service layer operations
func (l *Logger) LogServiceOperation(ctx context.Context, service, operation string, userID string, duration time.Duration, err error) {
	attrs := []any{
		ServiceField, service,
		OperationField, operation,
		DurationField, duration.Milliseconds(),
	}

	if userID != "" {
		attrs = append(attrs, UserIDField, userID)
	}

	if err != nil {
		attrs = append(attrs, ErrorField, err.Error())
		l.WithContext(ctx).Error("Service operation failed", attrs...)
	} else {
		l.WithContext(ctx).Info("Service operation completed", attrs...)
	}
}

// LogAuthEvent logs authentication-related events
func (l *Logger) LogAuthEvent(ctx context.Context, event, userID, clientIP string, success bool, details map[string]any) {
	attrs := []any{
		EventField, event,
		UserIDField, userID,
		ClientIPField, clientIP,
		SuccessField, success,
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
		ComponentField, component,
		VersionField, version,
	}

	for k, v := range config {
		attrs = append(attrs, k, v)
	}

	l.Info("Application starting", attrs...)
}

// LogShutdown logs application shutdown information
func (l *Logger) LogShutdown(component string, reason string, graceful bool) {
	l.Info("Application shutting down",
		ComponentField, component,
		ReasonField, reason,
		GracefulField, graceful,
	)
}

// Business logic logging helpers

// LogInfo logs informational messages with service context
func (l *Logger) LogInfo(ctx context.Context, msg string, attrs ...any) {
	l.WithContext(ctx).Info(msg, attrs...)
}

// LogWarn logs warning messages with service context
func (l *Logger) LogWarn(ctx context.Context, msg string, attrs ...any) {
	l.WithContext(ctx).Warn(msg, attrs...)
}

// LogDebug logs debug messages with service context
func (l *Logger) LogDebug(ctx context.Context, msg string, attrs ...any) {
	l.WithContext(ctx).Debug(msg, attrs...)
}

// LogValidationError logs validation errors with standard format
func (l *Logger) LogValidationError(ctx context.Context, field string, value any, reason string, userID string) {
	attrs := []any{
		"field", field,
		"value", value,
		ReasonField, reason,
	}

	if userID != "" {
		attrs = append(attrs, UserIDField, userID)
	}

	l.WithContext(ctx).Warn("Validation error", attrs...)
}

// LogSecurityEvent logs security-related events
func (l *Logger) LogSecurityEvent(ctx context.Context, event string, userID string, clientIP string, details map[string]any) {
	attrs := []any{
		EventField, event,
		UserIDField, userID,
		ClientIPField, clientIP,
	}

	for k, v := range details {
		attrs = append(attrs, k, v)
	}

	l.WithContext(ctx).Warn("Security event", attrs...)
}

// Helper functions

// GetTraceIDKey returns the context key for trace ID
func GetTraceIDKey() contextKey {
	return traceIDKey
}

// SetTraceID sets the trace ID in context
func SetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// getTraceID extracts trace ID from context
func getTraceID(ctx context.Context) string {
	if traceID := ctx.Value(traceIDKey); traceID != nil {
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
			OperationField, operation,
			DurationField, duration.Milliseconds(),
			ErrorField, err.Error(),
		)
	} else {
		l.WithContext(ctx).Debug("Operation completed",
			OperationField, operation,
			DurationField, duration.Milliseconds(),
		)
	}

	return err
}

// MeasureUserOperation measures user-related operations with detailed logging
func (l *Logger) MeasureUserOperation(ctx context.Context, operation string, userID string, fn func() error) error {
	start := time.Now()
	err := fn()
	duration := time.Since(start)

	attrs := []any{
		OperationField, operation,
		UserIDField, userID,
		DurationField, duration.Milliseconds(),
	}

	if err != nil {
		attrs = append(attrs, ErrorField, err.Error())
		l.WithContext(ctx).Error("User operation failed", attrs...)
	} else {
		l.WithContext(ctx).Info("User operation completed", attrs...)
	}

	return err
}
