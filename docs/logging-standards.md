# EngLog Logging Standards Guide

> "The bitterness of poor quality remains long after the sweetness of low price is forgotten." - Benjamin Franklin 📊

## Table of Contents

1. [Overview](#overview)
2. [Field Standards](#field-standards)
3. [Logger Configuration](#logger-configuration)
4. [Usage Patterns](#usage-patterns)
5. [Best Practices](#best-practices)
6. [Examples](#examples)
7. [Troubleshooting](#troubleshooting)
8. [Performance Considerations](#performance-considerations)

## Overview

The EngLog logging system provides a standardized, structured approach to application logging that enables:

- **Consistent field naming** across all services
- **Easy log aggregation and analysis** with tools like ELK stack
- **Automated monitoring and alerting** based on structured data
- **Enhanced debugging capabilities** with contextual information
- **Performance tracking** with built-in duration measurements

### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
├─────────────────────────────────────────────────────────────┤
│  Services (user_service, auth_service, project_service)     │
│  ↓ Uses standardized logger methods                         │
├─────────────────────────────────────────────────────────────┤
│              EngLog Logger (Wrapper)                        │
│  • Standardized field names                                 │
│  • Service/Component identification                         │
│  • Specialized methods (LogUserOperation, etc.)             │
│  ↓ Wraps slog.Logger                                        │
├─────────────────────────────────────────────────────────────┤
│                 Go slog.Logger                              │
│  • JSON/Text formatting                                     │
│  • Level-based filtering                                    │
│  • Context propagation                                      │
│  ↓ Outputs to                                               │
├─────────────────────────────────────────────────────────────┤
│              Output Destinations                            │
│  • Console (development)                                    │
│  • Files (production)                                       │
│  • Log aggregation systems (ELK, etc.)                      │
└─────────────────────────────────────────────────────────────┘
```

## Field Standards

### Core Field Constants

All logging uses predefined field constants to ensure consistency:

```go
const (
    ServiceField    = "service"        // Identifies the service (user_service, auth_service)
    ComponentField  = "component"      // Identifies the component within service
    OperationField  = "operation"      // Specific operation being performed
    UserIDField     = "user_id"        // User identifier for user-related operations
    EmailField      = "email"          // User email for authentication/profile operations
    ErrorField      = "error"          // Error message for failed operations
    DurationField   = "duration_ms"    // Operation duration in milliseconds
    StatusField     = "status"         // Operation status (success, failed, pending)
    MethodField     = "method"         // HTTP method for request logging
    PathField       = "path"           // HTTP path for request logging
    StatusCodeField = "status_code"    // HTTP status code
    ClientIPField   = "client_ip"      // Client IP address
    TraceIDField    = "trace_id"       // Request tracing identifier
    TableField      = "table"          // Database table for DB operations
    EventField      = "event"          // Event type for security/auth events
    SuccessField    = "success"        // Boolean success indicator
    VersionField    = "version"        // Application version
    ReasonField     = "reason"         // Reason for operation (shutdown, failure)
    GracefulField   = "graceful"       // Boolean for graceful shutdown
)
```

### Field Usage Matrix

| Field | User Ops | Auth | HTTP | DB | Security | Performance |
|-------|----------|------|------|----|---------|-----------  |
| service | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| operation | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| user_id | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ |
| email | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ |
| duration_ms | ✅ | ✅ | ✅ | ✅ | ❌ | ✅ |
| error | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| client_ip | ❌ | ✅ | ✅ | ❌ | ✅ | ❌ |
| trace_id | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |

## Logger Configuration

### Basic Logger Creation

```go
// Create base logger from configuration
cfg := config.LoggingConfig{
    Level:     slog.LevelInfo,
    Format:    "json",
    AddSource: true,
}
baseLogger := logging.NewLogger(cfg)

// Configure for specific service
userLogger := baseLogger.WithService("user_service")

// Configure for specific component within service
profileLogger := userLogger.WithComponent("profile_manager")

// Or configure both at once
logger := baseLogger.WithServiceAndComponent("user_service", "profile_manager")
```

### Service Identification Pattern

Each service should create its logger with appropriate identification:

```go
// In internal/services/user.go
func NewUserService(db *database.DB, logger *logging.Logger) *UserService {
    return &UserService{
        db:     db,
        logger: logger.WithService("user_service"),
    }
}

// In internal/handlers/user.go
func NewUserHandler(userService *services.UserService, logger *logging.Logger) *UserHandler {
    return &UserHandler{
        userService: userService,
        logger:      logger.WithServiceAndComponent("user_service", "http_handler"),
    }
}
```

## Usage Patterns

### 1. User Operations

For operations involving users, use the specialized user operation methods:

```go
// Success case
logger.LogUserOperation(ctx, "update_profile", userID, email, true,
    "fields_changed", []string{"first_name", "timezone"},
    "previous_timezone", "UTC",
    "new_timezone", "America/New_York")

// Failure case with error details
logger.LogUserOperation(ctx, "change_password", userID, email, false,
    "reason", "current_password_incorrect",
    "attempt_count", 3)
```

### 2. Generic Operations with Timing

For operations that need duration tracking:

```go
// Using MeasureUserOperation for automatic timing
err := logger.MeasureUserOperation(ctx, "update_profile", userID, func() error {
    return s.updateUserProfile(ctx, userID, profileData)
})

// Using MeasureDuration for non-user operations
err := logger.MeasureDuration(ctx, "database_migration", func() error {
    return s.runMigration(ctx)
})
```

### 3. Error Logging

Standardized error logging with context:

```go
// Simple error logging
logger.LogError(ctx, err, "Failed to update user profile",
    "user_id", userID,
    "attempted_fields", []string{"first_name", "last_name"})

// Validation errors
logger.LogValidationError(ctx, "email", invalidEmail, "invalid_format", userID)

// Security events
logger.LogSecurityEvent(ctx, "suspicious_login_attempt", userID, clientIP, map[string]any{
    "failed_attempts": 5,
    "time_window": "5_minutes",
    "user_agent": "suspicious_bot",
})
```

### 4. HTTP Request Logging

For HTTP request/response logging:

```go
// In middleware or handler
logger.LogRequest(ctx, "POST", "/api/v1/users/profile", 200,
    time.Duration(150)*time.Millisecond, "192.168.1.100")
```

### 5. Database Operations

For database operation logging:

```go
start := time.Now()
result, err := s.db.UpdateUser(ctx, params)
duration := time.Since(start)

logger.LogDatabaseOperation(ctx, "UPDATE", "users", duration, err)
```

### 6. Authentication Events

For authentication and authorization events:

```go
// Successful login
logger.LogAuthEvent(ctx, "user_login", userID, clientIP, true, map[string]any{
    "method": "email_password",
    "session_duration": "24h",
})

// Failed login
logger.LogAuthEvent(ctx, "user_login", "", clientIP, false, map[string]any{
    "method": "email_password",
    "error": "invalid_credentials",
    "email": email,
})
```

## Best Practices

### 1. Service Naming Convention

Use consistent service names across the application:

- `user_service` - User management operations
- `auth_service` - Authentication and authorization
- `project_service` - Project management operations
- `analytics_service` - Analytics and reporting
- `log_entry_service` - Activity log management
- `tag_service` - Tag management operations
- `api_server` - HTTP API server
- `worker_service` - Background worker operations
- `ai_service` - AI/ML operations

### 2. Operation Naming Convention

Use descriptive, consistent operation names:

- `create_user` instead of `create` or `new_user`
- `update_profile` instead of `update` or `change_profile`
- `authenticate_user` instead of `auth` or `login`
- `validate_token` instead of `check_token` or `verify`

### 3. Error Context

Always include relevant context in error logs:

```go
// Good - includes context
logger.LogError(ctx, err, "Failed to update user profile",
    "user_id", userID,
    "operation", "update_profile",
    "fields_attempted", []string{"first_name", "timezone"},
    "validation_errors", validationErrs)

// Bad - lacks context
logger.LogError(ctx, err, "Update failed")
```

### 4. Sensitive Data Handling

Never log sensitive information directly:

```go
// Good - log email but not password
logger.LogUserOperation(ctx, "change_password", userID, email, true,
    "password_strength", "strong")

// Bad - logs sensitive data
logger.LogUserOperation(ctx, "change_password", userID, email, true,
    "new_password", newPassword) // ❌ NEVER DO THIS
```

### 5. Performance Logging

Use performance logging for operations that might be slow:

```go
// Critical operations that should be monitored
err := logger.MeasureUserOperation(ctx, "generate_weekly_report", userID, func() error {
    return s.generateWeeklyReport(ctx, userID)
})

// Database operations that might be slow
err := logger.MeasureDuration(ctx, "refresh_materialized_views", func() error {
    return s.refreshAnalyticsViews(ctx)
})
```

## Examples

### Complete Service Implementation

```go
package services

import (
    "context"
    "time"

    "github.com/garnizeh/englog/internal/logging"
    "github.com/garnizeh/englog/internal/models"
)

type UserService struct {
    logger *logging.Logger
    db     *database.DB
}

func NewUserService(db *database.DB, logger *logging.Logger) *UserService {
    return &UserService{
        db:     db,
        logger: logger.WithService("user_service"),
    }
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, req *models.ProfileRequest) (*models.UserProfile, error) {
    // Input validation with logging
    if err := s.validateProfileRequest(req); err != nil {
        s.logger.LogValidationError(ctx, "profile_request", req, err.Error(), userID)
        return nil, err
    }

    // Log operation start
    s.logger.LogInfo(ctx, "Starting profile update",
        logging.OperationField, "update_profile",
        logging.UserIDField, userID,
        "fields_to_update", []string{"first_name", "last_name", "timezone"})

    // Perform operation with timing
    var profile *models.UserProfile
    err := s.logger.MeasureUserOperation(ctx, "update_profile", userID, func() error {
        var updateErr error
        profile, updateErr = s.performProfileUpdate(ctx, userID, req)
        return updateErr
    })

    if err != nil {
        s.logger.LogUserOperation(ctx, "update_profile", userID, req.Email, false,
            "error_type", "database_error",
            "fields_attempted", []string{"first_name", "last_name", "timezone"})
        return nil, err
    }

    // Log successful completion
    s.logger.LogUserOperation(ctx, "update_profile", userID, profile.Email, true,
        "fields_updated", []string{"first_name", "last_name", "timezone"},
        "previous_timezone", req.PreviousTimezone,
        "new_timezone", profile.Timezone)

    return profile, nil
}
```

### Handler Implementation

```go
package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/garnizeh/englog/internal/logging"
)

type UserHandler struct {
    logger      *logging.Logger
    userService *services.UserService
}

func NewUserHandler(userService *services.UserService, logger *logging.Logger) *UserHandler {
    return &UserHandler{
        userService: userService,
        logger:      logger.WithServiceAndComponent("user_service", "http_handler"),
    }
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
    start := time.Now()
    userID := c.GetString("user_id")
    clientIP := c.ClientIP()

    // Log incoming request
    h.logger.LogInfo(c.Request.Context(), "Received profile update request",
        logging.OperationField, "update_profile",
        logging.UserIDField, userID,
        logging.ClientIPField, clientIP,
        logging.MethodField, c.Request.Method,
        logging.PathField, c.Request.URL.Path)

    var req models.ProfileRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        duration := time.Since(start)
        h.logger.LogRequest(c.Request.Context(), c.Request.Method, c.Request.URL.Path,
            http.StatusBadRequest, duration, clientIP)

        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    profile, err := h.userService.UpdateProfile(c.Request.Context(), userID, &req)
    duration := time.Since(start)

    if err != nil {
        h.logger.LogRequest(c.Request.Context(), c.Request.Method, c.Request.URL.Path,
            http.StatusInternalServerError, duration, clientIP)

        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
        return
    }

    // Log successful request
    h.logger.LogRequest(c.Request.Context(), c.Request.Method, c.Request.URL.Path,
        http.StatusOK, duration, clientIP)

    c.JSON(http.StatusOK, profile)
}
```

## Troubleshooting

### Common Log Queries

With the standardized fields, you can easily query logs:

**Find all user operations for a specific user:**
```bash
# Using jq for JSON logs
cat logs/api.log | jq 'select(.user_id == "123e4567-e89b-12d3-a456-426614174000")'

# Using grep for text logs
grep 'user_id=123e4567-e89b-12d3-a456-426614174000' logs/api.log
```

**Find slow operations (> 1000ms):**
```bash
cat logs/api.log | jq 'select(.duration_ms > 1000) | {operation, duration_ms, user_id}'
```

**Find all failed user operations:**
```bash
cat logs/api.log | jq 'select(.service == "user_service" and .success == false)'
```

**Find security events:**
```bash
cat logs/api.log | jq 'select(.event and (.event | contains("security") or contains("suspicious")))'
```

### Log Analysis Patterns

**Performance monitoring:**
```bash
# Average operation duration by operation type
cat logs/api.log | jq -r 'select(.operation and .duration_ms) | "\(.operation) \(.duration_ms)"' | \
awk '{sum[$1]+=$2; count[$1]++} END {for(op in sum) print op, sum[op]/count[op]}'
```

**Error analysis:**
```bash
# Most common errors by service
cat logs/api.log | jq -r 'select(.error) | "\(.service) \(.error)"' | sort | uniq -c | sort -nr
```

**User activity patterns:**
```bash
# Most active users by operation count
cat logs/api.log | jq -r 'select(.user_id and .operation) | .user_id' | sort | uniq -c | sort -nr | head -10
```

## Performance Considerations

### 1. Log Level Management

Configure appropriate log levels for different environments:

```go
// Development
cfg.Level = slog.LevelDebug  // Show all logs including debug

// Staging
cfg.Level = slog.LevelInfo   // Show info and above

// Production
cfg.Level = slog.LevelWarn   // Show warnings and errors only
```

### 2. Structured Data Size

Be mindful of log entry size, especially with large data structures:

```go
// Good - log summary information
logger.LogUserOperation(ctx, "bulk_import", userID, email, true,
    "records_processed", len(records),
    "import_type", "csv",
    "file_size_bytes", fileSize)

// Bad - log entire data structure
logger.LogUserOperation(ctx, "bulk_import", userID, email, true,
    "all_records", records) // ❌ Could be huge
```

### 3. Context Propagation

Always use context for trace ID propagation:

```go
// Good - maintains trace context
func (s *UserService) UpdateProfile(ctx context.Context, userID string, req *models.ProfileRequest) error {
    s.logger.LogInfo(ctx, "Starting profile update", ...)  // Uses ctx for trace ID
    return s.performUpdate(ctx, userID, req)
}

// Bad - loses trace context
func (s *UserService) UpdateProfile(ctx context.Context, userID string, req *models.ProfileRequest) error {
    s.logger.Info("Starting profile update", ...)  // No ctx = no trace ID
    return s.performUpdate(ctx, userID, req)
}
```

### 4. Async Operations

For background operations, ensure proper context and timing:

```go
// Background task with proper logging
go func() {
    bgCtx := context.Background()
    bgCtx = logging.SetTraceID(bgCtx, generateTraceID())

    err := logger.MeasureDuration(bgCtx, "cleanup_expired_sessions", func() error {
        return s.cleanupExpiredSessions(bgCtx)
    })

    if err != nil {
        logger.LogError(bgCtx, err, "Background cleanup failed",
            "operation", "cleanup_expired_sessions")
    }
}()
```

---

This standardized logging approach provides comprehensive observability into the EngLog application, enabling effective monitoring, debugging, and performance analysis across all services and components.
