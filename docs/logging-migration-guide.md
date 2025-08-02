# EngLog Logging Migration Guide

> "Code is like humor. When you have to explain it, it's bad. But when you need to fix it, documentation is good." - Anonymous Developer 🔧

## Overview

This guide provides step-by-step instructions for migrating existing logging calls to the new standardized EngLog logging format.

**⚠️ IMPORTANT: All changes require explicit approval before implementation.**

## Approval Process

### 1. Change Proposal
For each logging call to be modified:
- Present the current code
- Propose the new standardized version
- Explain the benefits of the change
- Wait for explicit approval

### 2. Implementation
Only after receiving approval:
- Make the specific change
- Verify compilation
- Test functionality
- Report completion

### 3. Documentation
For each approved change:
- Document what was changed
- Note any issues encountered
- Confirm expected behavior

## Migration Patterns

### 1. Simple Log Replacement

**Before (Old Pattern):**
```go
s.logger.Info("Getting user profile", "user_id", userID)
```

**After (New Pattern):**
```go
s.logger.LogInfo(ctx, "Getting user profile",
    logging.UserIDField, userID,
    logging.OperationField, "get_profile")
```

### 2. Error Logging Migration

**Before:**
```go
s.logger.LogError(ctx, err, "Failed to update user profile", "user_id", userID)
```

**After:**
```go
s.logger.LogError(ctx, err, "Failed to update user profile",
    logging.UserIDField, userID,
    logging.OperationField, "update_profile")
```

### 3. Service Configuration Migration

**Before:**
```go
func NewUserService(db *database.DB, logger *logging.Logger) *UserService {
    return &UserService{
        db:     db,
        logger: logger.WithComponent("user_service"),
    }
}
```

**After:**
```go
func NewUserService(db *database.DB, logger *logging.Logger) *UserService {
    return &UserService{
        db:     db,
        logger: logger.WithService("user_service"),
    }
}
```

### 4. User Operation Logging

**Before:**
```go
s.logger.Info("User profile updated successfully",
    "user_id", userID,
    "email", profile.Email)
```

**After:**
```go
s.logger.LogUserOperation(ctx, "update_profile", userID, profile.Email, true,
    "fields_updated", []string{"first_name", "last_name", "timezone"})
```

## Service-Specific Migration Examples

### User Service Migration

**File: `internal/services/user.go`**

**Before:**
```go
s.logger.Info("Updating user profile",
    "service", "user_service",
    "user_id", userID,
    "first_name", req.FirstName,
    "last_name", req.LastName,
    "timezone", req.Timezone)
```

**After:**
```go
s.logger.LogInfo(ctx, "Updating user profile",
    logging.OperationField, "update_profile",
    logging.UserIDField, userID,
    "first_name", req.FirstName,
    "last_name", req.LastName,
    "timezone", req.Timezone)
```

### Auth Service Migration

**File: `internal/services/auth.go`**

**Before:**
```go
logger.Info("User login successful",
    "service", "auth_service",
    "user_id", user.ID,
    "email", user.Email,
    "client_ip", clientIP)
```

**After:**
```go
logger.LogAuthEvent(ctx, "user_login", user.ID, clientIP, true, map[string]any{
    "email": user.Email,
    "method": "email_password",
})
```

### Worker Service Migration

**File: `internal/worker/client.go`**

**Before:**
```go
c.logger.Info("Starting task processing",
    "service", "worker_client",
    "task_id", task.ID,
    "task_type", task.Type)
```

**After:**
```go
c.logger.LogInfo(ctx, "Starting task processing",
    logging.OperationField, "process_task",
    "task_id", task.ID,
    "task_type", task.Type)
```

## Manual Migration Approach

**All migrations must be done manually, case by case, with explicit approval for each change.**

The migration process follows these principles:
1. **Individual Review**: Each logging call is reviewed and updated individually
2. **Explicit Approval**: Every change must be approved before implementation
3. **Incremental Updates**: Changes are made in small, reviewable increments
4. **Testing After Each Change**: Compilation and functionality verified after each update## Manual Migration Checklist

### Step-by-Step Approval Process:

#### Phase 1: Service Configuration
For each service file, update the logger initialization first:

- [ ] **File**: `internal/services/user.go`
  - [ ] Review current `NewUserService()` method
  - [ ] Propose change from `WithComponent("user_service")` to `WithService("user_service")`
  - [ ] **APPROVAL REQUIRED** before implementation
  - [ ] Verify compilation after change

- [ ] **File**: `internal/services/auth.go`
  - [ ] Review current `NewAuthService()` method
  - [ ] Propose logger configuration change
  - [ ] **APPROVAL REQUIRED** before implementation
  - [ ] Verify compilation after change

#### Phase 2: Field Constants Migration
For each logging call, replace string literals individually:

- [ ] **User Service Logs**:
  - [ ] Identify each `"user_id"` → propose `logging.UserIDField`
  - [ ] **APPROVAL REQUIRED** for each change
  - [ ] Identify each `"operation"` → propose `logging.OperationField`
  - [ ] **APPROVAL REQUIRED** for each change
  - [ ] Continue field by field...

#### Phase 3: Method Signature Updates
Convert to specialized logging methods one by one:

- [ ] **Simple Info Logs** → `LogInfo(ctx, ...)`
  - [ ] Present each proposed change
  - [ ] **APPROVAL REQUIRED** before implementation

- [ ] **User Operations** → `LogUserOperation(...)`
  - [ ] Present each proposed change with parameters
  - [ ] **APPROVAL REQUIRED** before implementation

- [ ] **Error Logs** → Enhanced `LogError(...)`
  - [ ] Present each proposed change
  - [ ] **APPROVAL REQUIRED** before implementation

#### Phase 4: Context Addition
Ensure all logging methods receive context:

- [ ] **Review Each Log Call**:
  - [ ] Check if `ctx context.Context` is available
  - [ ] Propose context addition where missing
  - [ ] **APPROVAL REQUIRED** for each addition

### Priority Order for Manual Review:

1. **High Priority**: User service operations (most critical for application)
2. **Medium Priority**: Authentication and security logs
3. **Lower Priority**: Worker and background task logs

## Validation Steps

After each individual change, validate immediately:

### 1. Compilation Check (After Each Change)
```bash
cd /media/code/code/Go/garnizeh/englog
go build ./internal/services/user.go  # For user service changes
go build ./internal/services/auth.go  # For auth service changes
# etc.
```

### 2. Individual Function Testing
```bash
# Test specific functionality after each change
go test ./internal/services -run TestUserService_UpdateProfile
```

### 3. Manual Log Review
After implementing approved changes:
```bash
# Start application and manually trigger the updated functionality
make dev-api

# In another terminal, monitor specific log output
tail -f logs/api/api.log | grep "user_service"
```

### 4. Change-by-Change Verification
For each approved and implemented change:
- [ ] Verify compilation succeeds
- [ ] Check log output format matches expected pattern
- [ ] Ensure no functionality regression
- [ ] Confirm standardized fields are present

## Common Migration Issues

### Issue 1: Missing Context Parameter

**Problem:**
```go
s.logger.Info("Operation completed", "user_id", userID)
```

**Solution:**
```go
s.logger.LogInfo(ctx, "Operation completed", logging.UserIDField, userID)
```

### Issue 2: Inconsistent Service Names

**Problem:**
```go
// In different files
logger.WithService("userService")
logger.WithService("user-service")
logger.WithService("UserService")
```

**Solution:**
```go
// Standardize to snake_case
logger.WithService("user_service")
```

### Issue 3: Missing Field Constants Import

**Problem:**
```go
// Compile error: undefined: logging.UserIDField
s.logger.LogInfo(ctx, "msg", logging.UserIDField, userID)
```

**Solution:**
```go
import "github.com/garnizeh/englog/internal/logging"
```

### Issue 4: Over-logging in Loops

**Problem:**
```go
for _, user := range users {
    s.logger.LogInfo(ctx, "Processing user", logging.UserIDField, user.ID)
    // Process user...
}
```

**Solution:**
```go
s.logger.LogInfo(ctx, "Starting bulk user processing",
    "user_count", len(users),
    logging.OperationField, "bulk_process_users")

for _, user := range users {
    // Only log errors or every N iterations
    if err := s.processUser(ctx, user); err != nil {
        s.logger.LogError(ctx, err, "Failed to process user",
            logging.UserIDField, user.ID)
    }
}

s.logger.LogInfo(ctx, "Completed bulk user processing",
    "user_count", len(users),
    logging.OperationField, "bulk_process_users")
```

## Post-Migration Benefits

After completing the migration, you'll gain:

### 1. **Consistent Log Analysis**
```bash
# Find all user operations across services
jq 'select(.user_id and .operation)' logs/*.log

# Monitor service performance
jq 'select(.service == "user_service" and .duration_ms > 1000)' logs/*.log
```

### 2. **Automated Monitoring**
```bash
# Create alerts based on standardized fields
jq 'select(.service and .error and .duration_ms > 5000)' logs/*.log
```

### 3. **Better Debugging**
```bash
# Follow a user's journey across services
jq 'select(.user_id == "specific-user-id")' logs/*.log | sort_by(.timestamp)
```

### 4. **Performance Insights**
```bash
# Average response time by operation
jq -r 'select(.operation and .duration_ms) | "\(.operation) \(.duration_ms)"' logs/*.log | \
awk '{sum[$1]+=$2; count[$1]++} END {for(op in sum) print op, sum[op]/count[op] "ms"}'
```

This migration guide ensures a smooth transition to the new standardized logging format while maintaining full observability during the process.
