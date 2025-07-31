-- +goose Up
-- +goose StatementBegin
-- Tags table for flexible categorization
CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    color VARCHAR(7) DEFAULT '#6c757d', -- Hex color for UI
    description TEXT,
    usage_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Junction table for many-to-many relationship between log entries and tags
CREATE TABLE IF NOT EXISTS log_entry_tags (
    log_entry_id UUID NOT NULL REFERENCES log_entries(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    PRIMARY KEY (log_entry_id, tag_id)
);

-- Indexes for tag operations
CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_name ON tags(name);
CREATE INDEX IF NOT EXISTS idx_tags_usage_count ON tags(usage_count DESC);

-- Indexes for log_entry_tags junction table
CREATE INDEX IF NOT EXISTS idx_log_entry_tags_entry ON log_entry_tags(log_entry_id);
CREATE INDEX IF NOT EXISTS idx_log_entry_tags_tag ON log_entry_tags(tag_id);

-- Function to update tag usage count
CREATE OR REPLACE FUNCTION update_tag_usage_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE tags SET usage_count = usage_count + 1 WHERE id = NEW.tag_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE tags SET usage_count = GREATEST(usage_count - 1, 0) WHERE id = OLD.tag_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update tag usage counts
DO $$
BEGIN
    -- Check if the trigger already exists
    IF NOT EXISTS (
        SELECT 1
        FROM pg_trigger t
        JOIN pg_class c ON t.tgrelid = c.oid
        WHERE t.tgname = 'trigger_update_tag_usage_count'
          AND c.relname = 'log_entry_tags'
          -- Optionally, check the schema if tables with the same name exist in different schemas
          -- AND c.relnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'public') -- or your schema name
    ) THEN
        -- If the trigger does not exist, create it
        CREATE TRIGGER trigger_update_tag_usage_count
            AFTER INSERT OR DELETE ON log_entry_tags
            FOR EACH ROW
            EXECUTE FUNCTION update_tag_usage_count();
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_update_tag_usage_count ON log_entry_tags;
DROP FUNCTION IF EXISTS update_tag_usage_count();
DROP INDEX IF EXISTS idx_log_entry_tags_tag;
DROP INDEX IF EXISTS idx_log_entry_tags_entry;
DROP INDEX IF EXISTS idx_tags_usage_count;
DROP UNIQUE INDEX IF EXISTS idx_tags_name;
DROP TABLE IF EXISTS log_entry_tags CASCADE;
DROP TABLE IF EXISTS tags CASCADE;
-- +goose StatementEnd
