-- EngLog Task Management Queries
-- Background task processing and job queue

-- name: CreateTask :one
INSERT INTO tasks (
    task_type, user_id, payload, priority, max_retries, scheduled_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = $1;

-- name: GetPendingTasks :many
SELECT * FROM tasks
WHERE status IN ('pending', 'retrying')
  AND scheduled_at <= NOW()
ORDER BY priority ASC, scheduled_at ASC
LIMIT $1;

-- name: GetTasksByUser :many
SELECT * FROM tasks
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTasksByType :many
SELECT * FROM tasks
WHERE task_type = $1
ORDER BY created_at DESC;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: StartTaskProcessing :one
UPDATE tasks
SET status = 'processing', started_at = NOW(), updated_at = NOW()
WHERE id = $1 AND status IN ('pending', 'retrying')
RETURNING *;

-- name: CompleteTask :one
UPDATE tasks
SET status = 'completed', result = $2, completed_at = NOW(),
    processing_duration_ms = EXTRACT(EPOCH FROM (NOW() - started_at)) * 1000,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: FailTask :one
UPDATE tasks
SET status = CASE
    WHEN retry_count < max_retries THEN 'retrying'
    ELSE 'failed'
    END,
    retry_count = retry_count + 1,
    error_message = $2,
    scheduled_at = CASE
        WHEN retry_count < max_retries THEN NOW() + INTERVAL '5 minutes' * (retry_count + 1)
        ELSE scheduled_at
    END,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CancelTask :exec
UPDATE tasks
SET status = 'cancelled', updated_at = NOW()
WHERE id = $1 AND status IN ('pending', 'retrying');

-- name: GetTaskQueue :many
SELECT
    task_type,
    status,
    COUNT(*) as task_count,
    AVG(EXTRACT(EPOCH FROM (NOW() - scheduled_at))) as avg_wait_time_seconds
FROM tasks
WHERE created_at >= NOW() - INTERVAL '24 hours'
GROUP BY task_type, status
ORDER BY task_type, status;

-- name: GetTaskStats :one
SELECT
    COUNT(*) as total_tasks,
    COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_tasks,
    COUNT(CASE WHEN status = 'processing' THEN 1 END) as processing_tasks,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_tasks,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_tasks,
    AVG(processing_duration_ms) as avg_processing_time
FROM tasks
WHERE created_at >= $1;

-- name: GetStuckTasks :many
SELECT * FROM tasks
WHERE status = 'processing'
  AND started_at < NOW() - INTERVAL '1 hour'
ORDER BY started_at ASC;

-- name: ResetStuckTasks :exec
UPDATE tasks
SET status = 'pending', started_at = NULL, updated_at = NOW()
WHERE status = 'processing'
  AND started_at < NOW() - INTERVAL '1 hour';

-- name: CleanupOldTasks :exec
DELETE FROM tasks
WHERE status IN ('completed', 'failed', 'cancelled')
  AND completed_at < $1;

-- name: GetTaskPerformanceMetrics :many
SELECT
    task_type,
    COUNT(*) as total_executed,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as successful,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed,
    AVG(processing_duration_ms) as avg_duration,
    AVG(retry_count) as avg_retries
FROM tasks
WHERE completed_at >= $1
GROUP BY task_type
ORDER BY total_executed DESC;

-- name: GetUserTaskHistory :many
SELECT * FROM tasks
WHERE user_id = $1
  AND created_at >= $2
ORDER BY created_at DESC;
