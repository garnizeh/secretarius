-- EngLog Tags Management Queries
-- Tag CRUD operations and usage statistics

-- name: CreateTag :one
INSERT INTO tags (name, color, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTagByID :one
SELECT * FROM tags
WHERE id = $1;

-- name: GetTagByName :one
SELECT * FROM tags
WHERE name = $1;

-- name: GetAllTags :many
SELECT * FROM tags
ORDER BY usage_count DESC, name ASC;

-- name: GetPopularTags :many
SELECT * FROM tags
WHERE usage_count > 0
ORDER BY usage_count DESC, name ASC
LIMIT $1;

-- name: UpdateTag :one
UPDATE tags
SET name = $2, color = $3, description = $4
WHERE id = $1
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE id = $1;

-- name: SearchTags :many
SELECT * FROM tags
WHERE name ILIKE '%' || $1 || '%'
ORDER BY usage_count DESC, name ASC
LIMIT $2;

-- name: AddTagToLogEntry :exec
INSERT INTO log_entry_tags (log_entry_id, tag_id)
VALUES ($1, $2)
ON CONFLICT (log_entry_id, tag_id) DO NOTHING;

-- name: RemoveTagFromLogEntry :exec
DELETE FROM log_entry_tags
WHERE log_entry_id = $1 AND tag_id = $2;

-- name: GetTagsForLogEntry :many
SELECT t.* FROM tags t
JOIN log_entry_tags let ON t.id = let.tag_id
WHERE let.log_entry_id = $1
ORDER BY t.name;

-- name: GetLogEntriesForTag :many
SELECT le.* FROM log_entries le
JOIN log_entry_tags let ON le.id = let.log_entry_id
WHERE let.tag_id = $1
ORDER BY le.start_time DESC;

-- name: GetTagUsageStats :one
SELECT
    COUNT(*) as total_tags,
    COUNT(CASE WHEN usage_count > 0 THEN 1 END) as used_tags,
    AVG(usage_count) as avg_usage,
    MAX(usage_count) as max_usage
FROM tags;

-- name: GetUserTagUsage :many
SELECT
    t.id, t.name, t.color, t.description,
    COUNT(let.log_entry_id) as user_usage_count
FROM tags t
LEFT JOIN log_entry_tags let ON t.id = let.tag_id
LEFT JOIN log_entries le ON let.log_entry_id = le.id AND le.user_id = $1
GROUP BY t.id, t.name, t.color, t.description
HAVING COUNT(let.log_entry_id) > 0
ORDER BY user_usage_count DESC, t.name ASC;

-- name: GetRecentlyUsedTags :many
SELECT DISTINCT t.*, MAX(le.created_at) as last_used
FROM tags t
JOIN log_entry_tags let ON t.id = let.tag_id
JOIN log_entries le ON let.log_entry_id = le.id
WHERE le.user_id = $1
  AND le.created_at >= $2
GROUP BY t.id, t.name, t.color, t.description, t.usage_count, t.created_at
ORDER BY last_used DESC;

-- name: CleanupUnusedTags :exec
DELETE FROM tags
WHERE usage_count = 0
  AND created_at < NOW() - INTERVAL '30 days';
