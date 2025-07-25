-- EngLog Insights Management Queries
-- AI-generated insights and analytics

-- name: CreateInsight :one
INSERT INTO generated_insights (
    user_id, report_type, period_start, period_end, title,
    content, summary, metadata, generation_model,
    generation_duration_ms, quality_score
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetInsightByID :one
SELECT * FROM generated_insights
WHERE id = $1;

-- name: GetInsightsByUser :many
SELECT * FROM generated_insights
WHERE user_id = $1 AND status = 'active'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetInsightsByUserAndType :many
SELECT * FROM generated_insights
WHERE user_id = $1 AND report_type = $2 AND status = 'active'
ORDER BY created_at DESC;

-- name: GetInsightsByPeriod :many
SELECT * FROM generated_insights
WHERE user_id = $1
  AND period_start >= $2
  AND period_end <= $3
  AND status = 'active'
ORDER BY period_start DESC;

-- name: UpdateInsight :one
UPDATE generated_insights
SET title = $2, content = $3, summary = $4,
    metadata = $5, updated_at = NOW()
WHERE id = $1 AND user_id = $6
RETURNING *;

-- name: ArchiveInsight :exec
UPDATE generated_insights
SET status = 'archived', updated_at = NOW()
WHERE id = $1 AND user_id = $2;

-- name: DeleteInsight :exec
DELETE FROM generated_insights
WHERE id = $1 AND user_id = $2;

-- name: GetLatestInsightByType :one
SELECT * FROM generated_insights
WHERE user_id = $1 AND report_type = $2 AND status = 'active'
ORDER BY created_at DESC
LIMIT 1;

-- name: GetInsightStats :one
SELECT
    COUNT(*) as total_insights,
    COUNT(CASE WHEN status = 'active' THEN 1 END) as active_insights,
    AVG(quality_score) as avg_quality_score,
    AVG(generation_duration_ms) as avg_generation_time
FROM generated_insights
WHERE user_id = $1;

-- name: GetInsightsByQuality :many
SELECT * FROM generated_insights
WHERE user_id = $1
  AND quality_score >= $2
  AND status = 'active'
ORDER BY quality_score DESC, created_at DESC;

-- name: SupersedeOldInsights :exec
UPDATE generated_insights
SET status = 'superseded', updated_at = NOW()
WHERE user_id = $1
  AND report_type = $2
  AND period_start = $3
  AND period_end = $4
  AND status = 'active'
  AND id != $5;

-- name: CleanupOldInsights :exec
DELETE FROM generated_insights
WHERE created_at < $1
  AND status IN ('archived', 'superseded');

-- name: GetInsightGenerationMetrics :many
SELECT
    report_type,
    COUNT(*) as total_generated,
    AVG(generation_duration_ms) as avg_generation_time,
    AVG(quality_score) as avg_quality,
    MAX(created_at) as last_generated
FROM generated_insights
WHERE user_id = $1
  AND created_at >= $2
GROUP BY report_type
ORDER BY total_generated DESC;
