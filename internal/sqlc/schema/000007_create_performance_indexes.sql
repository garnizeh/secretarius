-- +goose Up
-- +goose StatementBegin
-- Additional performance indexes for complex queries
CREATE INDEX IF NOT EXISTS idx_log_entries_user_type_time ON log_entries (user_id, type, start_time DESC);

CREATE INDEX IF NOT EXISTS idx_log_entries_time_range ON log_entries (start_time, end_time);

CREATE INDEX IF NOT EXISTS idx_log_entries_full_text ON log_entries USING gin (
    to_tsvector (
        'english',
        title || ' ' || COALESCE(description, '')
    )
);

-- Composite index for analytics queries
CREATE INDEX IF NOT EXISTS idx_log_entries_analytics ON log_entries (
    user_id,
    start_time,
    value_rating,
    impact_level,
    duration_minutes
)
WHERE
    duration_minutes > 0;

-- Index for project statistics
CREATE INDEX IF NOT EXISTS idx_log_entries_project_stats ON log_entries (project_id, start_time, duration_minutes)
WHERE
    project_id IS NOT NULL;

-- Index for user activity patterns
CREATE INDEX IF NOT EXISTS idx_users_activity_pattern ON users (timezone, created_at)
WHERE
    timezone IS NOT NULL;

-- Index for active sessions
CREATE INDEX IF NOT EXISTS idx_sessions_active_recent ON user_sessions (is_active, last_activity DESC, user_id);

-- Index for tag popularity
CREATE INDEX IF NOT EXISTS idx_tags_popular ON tags (usage_count DESC, name)
WHERE
    usage_count > 0;

-- Composite index for insight generation queries
CREATE INDEX IF NOT EXISTS idx_insights_generation ON generated_insights (user_id, report_type, period_start DESC, status)
WHERE
    status = 'active';

-- Index for task queue processing
CREATE INDEX IF NOT EXISTS idx_tasks_queue ON tasks (status, priority ASC, scheduled_at ASC)
WHERE
    status IN ('pending', 'retrying');

-- Covering index for frequent log entry queries
CREATE INDEX IF NOT EXISTS idx_log_entries_cover ON log_entries (user_id, start_time DESC) INCLUDE (title, type, duration_minutes, value_rating);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_log_entries_cover;

DROP INDEX IF EXISTS idx_tasks_queue;

DROP INDEX IF EXISTS idx_insights_generation;

DROP INDEX IF EXISTS idx_tags_popular;

DROP INDEX IF EXISTS idx_sessions_active_recent;

DROP INDEX IF EXISTS idx_users_activity_pattern;

DROP INDEX IF EXISTS idx_log_entries_project_stats;

DROP INDEX IF EXISTS idx_log_entries_analytics;

DROP INDEX IF EXISTS idx_log_entries_full_text;

DROP INDEX IF EXISTS idx_log_entries_time_range;

DROP INDEX IF EXISTS idx_log_entries_user_type_time;

-- +goose StatementEnd