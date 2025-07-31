-- EngLog Project Management Queries
-- Project CRUD operations and statistics

-- name: CreateProject :one
INSERT INTO projects (
    name, description, color, status, start_date, end_date, created_by, is_default
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects
WHERE id = $1;

-- name: GetProjectsByUser :many
SELECT * FROM projects
WHERE created_by = $1
ORDER BY is_default DESC, name ASC;

-- name: GetActiveProjectsByUser :many
SELECT * FROM projects
WHERE created_by = $1 AND status = 'active'
ORDER BY is_default DESC, name ASC;

-- name: GetUserDefaultProject :one
SELECT * FROM projects
WHERE created_by = $1 AND is_default = true
LIMIT 1;

-- name: UpdateProject :one
UPDATE projects
SET name = $2, description = $3, color = $4, status = $5,
    start_date = $6, end_date = $7, is_default = $8, updated_at = NOW()
WHERE id = $1 AND created_by = $9
RETURNING *;

-- name: SetProjectAsDefault :exec
BEGIN;
-- First, unset all default projects for the user
UPDATE projects
SET is_default = false, updated_at = NOW()
WHERE created_by = $2 AND is_default = true;
-- Then set the specified project as default
UPDATE projects
SET is_default = true, updated_at = NOW()
WHERE id = $1 AND created_by = $2;
COMMIT;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1 AND created_by = $2;

-- name: GetProjectsWithActivity :many
SELECT
    p.*,
    COUNT(le.id) as entry_count,
    SUM(le.duration_minutes) as total_minutes
FROM projects p
LEFT JOIN log_entries le ON p.id = le.project_id
WHERE p.created_by = $1
GROUP BY p.id
ORDER BY p.is_default DESC, entry_count DESC;

-- name: GetProjectStats :one
SELECT
    COUNT(le.id) as total_entries,
    SUM(le.duration_minutes) as total_minutes,
    AVG(le.duration_minutes) as avg_duration,
    COUNT(DISTINCT le.user_id) as contributors_count,
    MIN(le.start_time) as first_activity,
    MAX(le.start_time) as last_activity
FROM log_entries le
WHERE le.project_id = $1;
