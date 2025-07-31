-- EngLog Log Entries Queries
-- Activity tracking and log entry management

-- name: CreateLogEntry :one
INSERT INTO log_entries (
    user_id, project_id, title, description, type,
    start_time, end_time, value_rating, impact_level
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetLogEntryByID :one
SELECT * FROM log_entries
WHERE id = $1;

-- name: GetLogEntriesByUser :many
SELECT * FROM log_entries
WHERE user_id = $1
ORDER BY start_time DESC
LIMIT $2 OFFSET $3;

-- name: GetLogEntriesByUserAndDateRange :many
SELECT * FROM log_entries
WHERE user_id = $1
  AND start_time >= $2
  AND end_time <= $3
ORDER BY start_time ASC;

-- name: GetLogEntriesByProject :many
SELECT * FROM log_entries
WHERE project_id = $1
ORDER BY start_time DESC;

-- name: GetLogEntriesByUserAndProject :many
SELECT * FROM log_entries
WHERE user_id = $1 AND project_id = $2
ORDER BY start_time DESC;

-- name: UpdateLogEntry :one
UPDATE log_entries
SET title = $2, description = $3, type = $4, project_id = $5,
    start_time = $6, end_time = $7, value_rating = $8,
    impact_level = $9, updated_at = NOW()
WHERE id = $1 AND user_id = $10
RETURNING *;

-- name: DeleteLogEntry :execrows
DELETE FROM log_entries
WHERE id = $1 AND user_id = $2;

-- name: GetUserActivitySummary :one
SELECT
    COUNT(*) as total_entries,
    SUM(duration_minutes) as total_minutes,
    AVG(duration_minutes) as avg_duration,
    COUNT(DISTINCT project_id) as projects_count,
    COUNT(DISTINCT DATE(start_time)) as active_days
FROM log_entries
WHERE user_id = $1
  AND start_time >= $2
  AND start_time <= $3;

-- name: GetRecentLogEntries :many
SELECT le.*, p.name as project_name, p.color as project_color
FROM log_entries le
LEFT JOIN projects p ON le.project_id = p.id
WHERE le.user_id = $1
ORDER BY le.start_time DESC
LIMIT $2;

-- name: GetLogEntriesByType :many
SELECT * FROM log_entries
WHERE user_id = $1 AND type = $2
ORDER BY start_time DESC;

-- name: GetHighValueEntries :many
SELECT * FROM log_entries
WHERE user_id = $1
  AND value_rating IN ('high', 'critical')
  AND start_time >= $2
ORDER BY start_time DESC;

-- name: GetLogEntriesWithTags :many
SELECT DISTINCT le.*,
       ARRAY_AGG(t.name ORDER BY t.name) as tag_names
FROM log_entries le
LEFT JOIN log_entry_tags let ON le.id = let.log_entry_id
LEFT JOIN tags t ON let.tag_id = t.id
WHERE le.user_id = $1
  AND le.start_time >= $2
  AND le.start_time <= $3
GROUP BY le.id
ORDER BY le.start_time DESC;

-- name: SearchLogEntries :many
SELECT * FROM log_entries
WHERE user_id = $1
  AND (
    title ILIKE '%' || $2 || '%' OR
    description ILIKE '%' || $2 || '%'
  )
ORDER BY start_time DESC
LIMIT $3;

-- name: GetDailyProductivityStats :many
SELECT
    DATE(start_time) as activity_date,
    COUNT(*) as entry_count,
    SUM(duration_minutes) as total_minutes,
    AVG(duration_minutes) as avg_duration,
    COUNT(CASE WHEN value_rating = 'critical' THEN 1 END) as critical_count,
    COUNT(CASE WHEN value_rating = 'high' THEN 1 END) as high_count
FROM log_entries
WHERE user_id = $1
  AND start_time >= $2
  AND start_time <= $3
GROUP BY DATE(start_time)
ORDER BY activity_date DESC;
