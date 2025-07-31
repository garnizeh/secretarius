-- +goose Up
-- +goose StatementBegin
-- Generated insights table for storing AI-generated analysis
CREATE TABLE IF NOT EXISTS generated_insights (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    report_type VARCHAR(50) NOT NULL CHECK (report_type IN (
        'daily_summary', 'weekly_summary', 'monthly_summary',
        'quarterly_summary', 'project_analysis', 'productivity_trends',
        'time_distribution', 'performance_review', 'goal_progress', 'custom'
    )),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL, -- AI-generated insight content
    summary TEXT, -- Short summary for quick viewing
    metadata JSONB DEFAULT '{}', -- Additional structured data
    generation_model VARCHAR(50), -- Which LLM model was used
    generation_duration_ms INTEGER, -- Time taken to generate
    quality_score DECIMAL(3,2), -- AI confidence score (0.00-1.00)
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'archived', 'superseded')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Constraints
    CONSTRAINT insights_period_check CHECK (period_end >= period_start),
    CONSTRAINT insights_quality_score_check CHECK (quality_score IS NULL OR (quality_score >= 0.00 AND quality_score <= 1.00))
);

-- Background tasks table for async processing
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_type VARCHAR(50) NOT NULL CHECK (task_type IN (
        'generate_insight', 'send_email', 'export_data', 'cleanup_data',
        'process_analytics', 'generate_report', 'backup_data', 'custom'
    )),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    payload JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN (
        'pending', 'processing', 'completed', 'failed', 'cancelled', 'retrying'
    )),
    priority INTEGER DEFAULT 5 CHECK (priority >= 1 AND priority <= 10), -- 1 = highest, 10 = lowest
    max_retries INTEGER DEFAULT 3,
    retry_count INTEGER DEFAULT 0,
    scheduled_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    result JSONB,
    error_message TEXT,
    processing_duration_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Constraints
    CONSTRAINT tasks_retry_check CHECK (retry_count <= max_retries)
);

-- Indexes for insights
CREATE INDEX IF NOT EXISTS idx_insights_user_period ON generated_insights(user_id, period_start, period_end);
CREATE INDEX IF NOT EXISTS idx_insights_report_type ON generated_insights(report_type);
CREATE INDEX IF NOT EXISTS idx_insights_status ON generated_insights(status);
CREATE INDEX IF NOT EXISTS idx_insights_created ON generated_insights(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_insights_quality ON generated_insights(quality_score DESC) WHERE quality_score IS NOT NULL;

-- Indexes for tasks
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_scheduled ON tasks(scheduled_at) WHERE status IN ('pending', 'retrying');
CREATE INDEX IF NOT EXISTS idx_tasks_user ON tasks(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_tasks_type_status ON tasks(task_type, status);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority, scheduled_at) WHERE status IN ('pending', 'retrying');

-- Function to automatically update task status
CREATE OR REPLACE FUNCTION update_task_timestamps()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();

    -- Set started_at when status changes to processing
    IF OLD.status != 'processing' AND NEW.status = 'processing' THEN
        NEW.started_at = NOW();
    END IF;

    -- Set completed_at when task finishes
    IF OLD.status IN ('pending', 'processing', 'retrying') AND NEW.status IN ('completed', 'failed', 'cancelled') THEN
        NEW.completed_at = NOW();

        -- Calculate processing duration if started_at is set
        IF NEW.started_at IS NOT NULL THEN
            NEW.processing_duration_ms = EXTRACT(EPOCH FROM (NOW() - NEW.started_at)) * 1000;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for task timestamps
DO $$
BEGIN
    -- Check if the trigger already exists
    IF NOT EXISTS (
        SELECT 1
        FROM pg_trigger t
        JOIN pg_class c ON t.tgrelid = c.oid
        WHERE t.tgname = 'trigger_update_task_timestamps'
          AND c.relname = 'tasks'
          -- Optionally, check the schema if tables with the same name exist in different schemas
          -- AND c.relnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'public') -- or your schema name
    ) THEN
        -- If the trigger does not exist, create it
        CREATE TRIGGER trigger_update_task_timestamps
            AFTER INSERT OR DELETE ON tasks
            FOR EACH ROW
            EXECUTE FUNCTION update_task_timestamps();
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_update_task_timestamps ON tasks;
DROP FUNCTION IF EXISTS update_task_timestamps();
DROP INDEX IF EXISTS idx_tasks_priority;
DROP INDEX IF EXISTS idx_tasks_type_status;
DROP INDEX IF EXISTS idx_tasks_user;
DROP INDEX IF EXISTS idx_tasks_scheduled;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_insights_quality;
DROP INDEX IF EXISTS idx_insights_created;
DROP INDEX IF EXISTS idx_insights_status;
DROP INDEX IF EXISTS idx_insights_report_type;
DROP INDEX IF EXISTS idx_insights_user_period;
DROP TABLE IF EXISTS tasks CASCADE;
DROP TABLE IF EXISTS generated_insights CASCADE;
-- +goose StatementEnd
