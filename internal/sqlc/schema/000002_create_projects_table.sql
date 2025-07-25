-- +goose Up
-- +goose StatementBegin
-- Projects table for organizing activities
CREATE TABLE
    IF NOT EXISTS projects (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        name VARCHAR(200) NOT NULL,
        description TEXT,
        color VARCHAR(7) DEFAULT '#3498db', -- Hex color for UI
        status VARCHAR(20) DEFAULT 'active' CHECK (
            status IN ('active', 'completed', 'on_hold', 'cancelled')
        ),
        start_date DATE,
        end_date DATE,
        created_by UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        is_default BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT NOW (),
            updated_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT NOW ()
    );

-- Ensure each user has only one default project
CREATE UNIQUE INDEX IF NOT EXISTS idx_projects_user_default ON projects (created_by)
WHERE
    is_default = TRUE;

-- Index for project lookups by user
CREATE INDEX IF NOT EXISTS idx_projects_created_by ON projects (created_by);

-- Index for active projects
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects (status);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_projects_status;
DROP INDEX IF EXISTS idx_projects_created_by;
DROP INDEX IF EXISTS idx_projects_user_default;
DROP TABLE IF EXISTS projects CASCADE;

-- +goose StatementEnd