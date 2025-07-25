-- EngLog Analytics Queries
-- Advanced analytics and reporting

-- name: GetUserActivitySummaryView :one
SELECT * FROM user_activity_summary
WHERE user_id = $1;

-- name: RefreshUserActivitySummary :exec
REFRESH MATERIALIZED VIEW user_activity_summary;

-- name: GetDailyActivityPattern :many
SELECT * FROM daily_activity_patterns
WHERE user_id = $1
  AND activity_date >= $2
  AND activity_date <= $3
ORDER BY activity_date DESC, hour_of_day ASC;

-- name: GetProjectPerformanceMetrics :many
SELECT * FROM project_performance_metrics
WHERE project_owner = $1
ORDER BY total_minutes DESC;

-- name: GetUserProductivityTrend :many
SELECT * FROM get_user_productivity_trend($1, $2);

-- name: GetActivityTypeDistribution :many
SELECT
    type,
    COUNT(*) as entry_count,
    SUM(duration_minutes) as total_minutes,
    AVG(duration_minutes) as avg_duration,
    ROUND(AVG(CASE
        WHEN value_rating = 'critical' THEN 4
        WHEN value_rating = 'high' THEN 3
        WHEN value_rating = 'medium' THEN 2
        WHEN value_rating = 'low' THEN 1
        ELSE 0
    END), 2) as avg_value_score
FROM log_entries
WHERE user_id = $1
  AND start_time >= $2
  AND start_time <= $3
GROUP BY type
ORDER BY total_minutes DESC;

-- name: GetValueRatingDistribution :many
SELECT
    value_rating,
    COUNT(*) as entry_count,
    SUM(duration_minutes) as total_minutes,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage
FROM log_entries
WHERE user_id = $1
  AND start_time >= $2
  AND start_time <= $3
GROUP BY value_rating
ORDER BY CASE value_rating
    WHEN 'critical' THEN 4
    WHEN 'high' THEN 3
    WHEN 'medium' THEN 2
    WHEN 'low' THEN 1
    ELSE 0
END DESC;

-- name: GetImpactLevelDistribution :many
SELECT
    impact_level,
    COUNT(*) as entry_count,
    SUM(duration_minutes) as total_minutes,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage
FROM log_entries
WHERE user_id = $1
  AND start_time >= $2
  AND start_time <= $3
GROUP BY impact_level
ORDER BY CASE impact_level
    WHEN 'company' THEN 4
    WHEN 'department' THEN 3
    WHEN 'team' THEN 2
    WHEN 'personal' THEN 1
    ELSE 0
END DESC;

-- name: GetWeeklyActivitySummary :many
SELECT
    DATE_TRUNC('week', start_time) as week_start,
    COUNT(*) as entry_count,
    SUM(duration_minutes) as total_minutes,
    AVG(duration_minutes) as avg_duration,
    COUNT(DISTINCT project_id) as projects_count,
    COUNT(DISTINCT DATE(start_time)) as active_days
FROM log_entries
WHERE user_id = $1
  AND start_time >= $2
  AND start_time <= $3
GROUP BY DATE_TRUNC('week', start_time)
ORDER BY week_start DESC;

-- name: GetMonthlyActivitySummary :many
SELECT
    DATE_TRUNC('month', start_time) as month_start,
    COUNT(*) as entry_count,
    SUM(duration_minutes) as total_minutes,
    AVG(duration_minutes) as avg_duration,
    COUNT(DISTINCT project_id) as projects_count,
    COUNT(DISTINCT DATE(start_time)) as active_days
FROM log_entries
WHERE user_id = $1
  AND start_time >= $2
  AND start_time <= $3
GROUP BY DATE_TRUNC('month', start_time)
ORDER BY month_start DESC;

-- name: GetTopProjectsByTime :many
SELECT
    p.id, p.name, p.color,
    COUNT(le.id) as entry_count,
    SUM(le.duration_minutes) as total_minutes,
    ROUND(SUM(le.duration_minutes) * 100.0 / SUM(SUM(le.duration_minutes)) OVER (), 2) as percentage
FROM log_entries le
JOIN projects p ON le.project_id = p.id
WHERE le.user_id = $1
  AND le.start_time >= $2
  AND le.start_time <= $3
GROUP BY p.id, p.name, p.color
ORDER BY total_minutes DESC
LIMIT $4;

-- name: GetProductivityByDayOfWeek :many
SELECT
    EXTRACT(DOW FROM start_time) as day_of_week,
    TO_CHAR(start_time, 'Day') as day_name,
    COUNT(*) as entry_count,
    SUM(duration_minutes) as total_minutes,
    AVG(duration_minutes) as avg_duration
FROM log_entries
WHERE user_id = $1
  AND start_time >= $2
  AND start_time <= $3
GROUP BY EXTRACT(DOW FROM start_time), TO_CHAR(start_time, 'Day')
ORDER BY day_of_week;

-- name: GetProductivityByHour :many
SELECT
    EXTRACT(HOUR FROM start_time) as hour_of_day,
    COUNT(*) as entry_count,
    SUM(duration_minutes) as total_minutes,
    AVG(duration_minutes) as avg_duration
FROM log_entries
WHERE user_id = $1
  AND start_time >= $2
  AND start_time <= $3
GROUP BY EXTRACT(HOUR FROM start_time)
ORDER BY hour_of_day;

-- name: GetComparisonStats :one
SELECT
    -- Current period
    COUNT(CASE WHEN start_time >= $2 AND start_time <= $3 THEN 1 END) as current_entries,
    SUM(CASE WHEN start_time >= $2 AND start_time <= $3 THEN duration_minutes ELSE 0 END) as current_minutes,

    -- Previous period
    COUNT(CASE WHEN start_time >= $4 AND start_time <= $5 THEN 1 END) as previous_entries,
    SUM(CASE WHEN start_time >= $4 AND start_time <= $5 THEN duration_minutes ELSE 0 END) as previous_minutes
FROM log_entries
WHERE user_id = $1;
