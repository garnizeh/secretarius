-- +goose Up
-- +goose StatementBegin
-- Materialized view for user activity analytics
CREATE MATERIALIZED VIEW IF NOT EXISTS user_activity_summary AS
SELECT
    u.id AS user_id,
    u.email,
    u.timezone,
    COUNT(le.id) AS total_entries,
    SUM(le.duration_minutes) AS total_minutes,
    AVG(le.duration_minutes) AS avg_duration_minutes,
    COUNT(DISTINCT le.project_id) AS projects_count,
    COUNT(DISTINCT DATE(le.start_time)) AS active_days,
    MIN(le.start_time) AS first_entry_date,
    MAX(le.start_time) AS last_entry_date,

    -- Value rating distribution
    COUNT(CASE WHEN le.value_rating = 'critical' THEN 1 END) AS critical_entries,
    COUNT(CASE WHEN le.value_rating = 'high' THEN 1 END) AS high_entries,
    COUNT(CASE WHEN le.value_rating = 'medium' THEN 1 END) AS medium_entries,
    COUNT(CASE WHEN le.value_rating = 'low' THEN 1 END) AS low_entries,

    -- Impact level distribution
    COUNT(CASE WHEN le.impact_level = 'company' THEN 1 END) AS company_impact_entries,
    COUNT(CASE WHEN le.impact_level = 'department' THEN 1 END) AS department_impact_entries,
    COUNT(CASE WHEN le.impact_level = 'team' THEN 1 END) AS team_impact_entries,
    COUNT(CASE WHEN le.impact_level = 'personal' THEN 1 END) AS personal_impact_entries,

    -- Activity type distribution
    COUNT(CASE WHEN le.type = 'development' THEN 1 END) AS development_entries,
    COUNT(CASE WHEN le.type = 'meeting' THEN 1 END) AS meeting_entries,
    COUNT(CASE WHEN le.type = 'code_review' THEN 1 END) AS review_entries,
    COUNT(CASE WHEN le.type = 'debugging' THEN 1 END) AS debugging_entries,

    -- Time-based metrics
    EXTRACT(EPOCH FROM (MAX(le.start_time) - MIN(le.start_time))) / 86400 AS activity_span_days,

    -- Recent activity (last 30 days)
    COUNT(CASE WHEN le.start_time >= NOW() - INTERVAL '30 days' THEN 1 END) AS recent_entries_30d,
    SUM(CASE WHEN le.start_time >= NOW() - INTERVAL '30 days' THEN le.duration_minutes ELSE 0 END) AS recent_minutes_30d,

    NOW() AS refreshed_at
FROM users u
LEFT JOIN log_entries le ON u.id = le.user_id
GROUP BY u.id, u.email, u.timezone;

-- Create unique index on materialized view
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_activity_summary_user ON user_activity_summary(user_id);
CREATE INDEX IF NOT EXISTS idx_user_activity_summary_refreshed ON user_activity_summary(refreshed_at);

-- View for daily activity patterns
CREATE OR REPLACE VIEW daily_activity_patterns AS
SELECT
    user_id,
    DATE(start_time) AS activity_date,
    EXTRACT(DOW FROM start_time) AS day_of_week, -- 0=Sunday, 6=Saturday
    EXTRACT(HOUR FROM start_time) AS hour_of_day,
    COUNT(*) AS entry_count,
    SUM(duration_minutes) AS total_minutes,
    AVG(duration_minutes) AS avg_duration,
    STRING_AGG(DISTINCT type, ', ' ORDER BY type) AS activity_types,
    AVG(CASE
        WHEN value_rating = 'critical' THEN 4
        WHEN value_rating = 'high' THEN 3
        WHEN value_rating = 'medium' THEN 2
        WHEN value_rating = 'low' THEN 1
        ELSE 0
    END) AS avg_value_score
FROM log_entries
GROUP BY user_id, DATE(start_time), EXTRACT(DOW FROM start_time), EXTRACT(HOUR FROM start_time);

-- View for project performance metrics
CREATE OR REPLACE VIEW project_performance_metrics AS
SELECT
    p.id AS project_id,
    p.name AS project_name,
    p.status AS project_status,
    p.created_by AS project_owner,
    COUNT(le.id) AS total_entries,
    SUM(le.duration_minutes) AS total_minutes,
    AVG(le.duration_minutes) AS avg_entry_duration,
    COUNT(DISTINCT le.user_id) AS contributors_count,
    COUNT(DISTINCT DATE(le.start_time)) AS active_days,
    MIN(le.start_time) AS first_activity,
    MAX(le.start_time) AS last_activity,

    -- Value distribution
    ROUND(AVG(CASE
        WHEN le.value_rating = 'critical' THEN 4
        WHEN le.value_rating = 'high' THEN 3
        WHEN le.value_rating = 'medium' THEN 2
        WHEN le.value_rating = 'low' THEN 1
        ELSE 0
    END), 2) AS avg_value_score,

    -- Most common activity types
    MODE() WITHIN GROUP (ORDER BY le.type) AS most_common_activity,

    -- Recent activity
    COUNT(CASE WHEN le.start_time >= NOW() - INTERVAL '30 days' THEN 1 END) AS recent_entries_30d,

    p.created_at AS project_created_at
FROM projects p
LEFT JOIN log_entries le ON p.id = le.project_id
GROUP BY p.id, p.name, p.status, p.created_by, p.created_at;

-- Function to refresh materialized views
CREATE OR REPLACE FUNCTION refresh_analytics_views()
RETURNS VOID AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_activity_summary;
    -- Log the refresh
    INSERT INTO tasks (task_type, payload, status, completed_at)
    VALUES ('refresh_analytics', '{"view": "user_activity_summary"}', 'completed', NOW());
END;
$$ LANGUAGE plpgsql;

-- Function to get user productivity trends
CREATE OR REPLACE FUNCTION get_user_productivity_trend(user_uuid UUID, days_back INTEGER DEFAULT 30)
RETURNS TABLE(
    date DATE,
    total_minutes INTEGER,
    entry_count INTEGER,
    avg_value_score DECIMAL,
    productivity_score DECIMAL
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        DATE(le.start_time) AS date,
        SUM(le.duration_minutes)::INTEGER AS total_minutes,
        COUNT(le.id)::INTEGER AS entry_count,
        ROUND(AVG(CASE
            WHEN le.value_rating = 'critical' THEN 4
            WHEN le.value_rating = 'high' THEN 3
            WHEN le.value_rating = 'medium' THEN 2
            WHEN le.value_rating = 'low' THEN 1
            ELSE 0
        END), 2) AS avg_value_score,
        -- Simple productivity score: (total_minutes * avg_value_score) / 100
        ROUND((SUM(le.duration_minutes) * AVG(CASE
            WHEN le.value_rating = 'critical' THEN 4
            WHEN le.value_rating = 'high' THEN 3
            WHEN le.value_rating = 'medium' THEN 2
            WHEN le.value_rating = 'low' THEN 1
            ELSE 0
        END)) / 100.0, 2) AS productivity_score
    FROM log_entries le
    WHERE le.user_id = user_uuid
        AND le.start_time >= CURRENT_DATE - INTERVAL '1 day' * days_back
    GROUP BY DATE(le.start_time)
    ORDER BY DATE(le.start_time) DESC;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS get_user_productivity_trend(UUID, INTEGER);
DROP FUNCTION IF EXISTS refresh_analytics_views();
DROP VIEW IF EXISTS project_performance_metrics;
DROP VIEW IF EXISTS daily_activity_patterns;
DROP INDEX IF EXISTS idx_user_activity_summary_refreshed;
DROP INDEX IF EXISTS idx_user_activity_summary_user;
DROP MATERIALIZED VIEW IF EXISTS user_activity_summary;
-- +goose StatementEnd
