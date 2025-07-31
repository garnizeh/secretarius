-- +goose Up
-- +goose StatementBegin
-- Log entries table for activity tracking
CREATE TABLE
    IF NOT EXISTS log_entries (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        project_id UUID REFERENCES projects (id) ON DELETE SET NULL,
        title VARCHAR(500) NOT NULL,
        description TEXT,
        type VARCHAR(50) NOT NULL CHECK (
            type IN (
                'development',
                'meeting',
                'code_review',
                'debugging',
                'documentation',
                'testing',
                'deployment',
                'research',
                'planning',
                'learning',
                'maintenance',
                'support',
                'other'
            )
        ),
        start_time TIMESTAMP
        WITH
            TIME ZONE NOT NULL,
            end_time TIMESTAMP
        WITH
            TIME ZONE NOT NULL,
            duration_minutes INTEGER GENERATED ALWAYS AS (
                EXTRACT(
                    EPOCH
                    FROM
                        (end_time - start_time)
                ) / 60
            ) STORED,
            value_rating VARCHAR(20) NOT NULL CHECK (
                value_rating IN ('low', 'medium', 'high', 'critical')
            ),
            impact_level VARCHAR(20) NOT NULL CHECK (
                impact_level IN ('personal', 'team', 'department', 'company')
            ),
            created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT NOW (),
            updated_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT NOW (),
            -- Constraints
            CONSTRAINT log_entries_time_check CHECK (end_time > start_time),
            CONSTRAINT log_entries_duration_check CHECK (
                duration_minutes > 0
                AND duration_minutes <= 1440
            ) -- Max 24 hours
    );

-- Performance indexes for common queries
CREATE INDEX IF NOT EXISTS idx_log_entries_user_time_desc ON log_entries (user_id, start_time DESC);

CREATE INDEX IF NOT EXISTS idx_log_entries_project_time ON log_entries (project_id, start_time);

CREATE INDEX IF NOT EXISTS idx_log_entries_type ON log_entries (type);

CREATE INDEX IF NOT EXISTS idx_log_entries_value_rating ON log_entries (value_rating);

CREATE INDEX IF NOT EXISTS idx_log_entries_impact_level ON log_entries (impact_level);

CREATE INDEX IF NOT EXISTS idx_log_entries_duration ON log_entries (duration_minutes);

-- Composite index for filtering and analytics
CREATE INDEX IF NOT EXISTS idx_log_entries_user_project_time ON log_entries (user_id, project_id, start_time);

CREATE INDEX IF NOT EXISTS idx_log_entries_value_impact ON log_entries (value_rating, impact_level);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_log_entries_value_impact;
DROP INDEX IF EXISTS idx_log_entries_user_project_time;
DROP INDEX IF EXISTS idx_log_entries_duration;
DROP INDEX IF EXISTS idx_log_entries_impact_level;
DROP INDEX IF EXISTS idx_log_entries_value_rating;
DROP INDEX IF EXISTS idx_log_entries_type;
DROP INDEX IF EXISTS idx_log_entries_project_time;
DROP INDEX IF EXISTS idx_log_entries_user_time_desc;
DROP TABLE IF EXISTS log_entries CASCADE;

-- +goose StatementEnd