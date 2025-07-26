-- +goose Up
-- +goose StatementBegin
-- Users table with authentication and profile information
CREATE TABLE
    IF NOT EXISTS users (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        email VARCHAR(255) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        first_name VARCHAR(100) NOT NULL,
        last_name VARCHAR(100) NOT NULL,
        timezone VARCHAR(50) DEFAULT 'UTC',
        preferences JSONB DEFAULT '{}',
        last_login_at TIMESTAMP
        WITH
            TIME ZONE,
            created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT NOW (),
            updated_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT NOW ()
    );

-- Index for email lookups
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email);

-- Index for user timezone queries
CREATE INDEX IF NOT EXISTS idx_users_timezone ON users (timezone);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_timezone;

DROP INDEX IF EXISTS idx_users_email;

DROP TABLE IF EXISTS users CASCADE;

DROP EXTENSION IF EXISTS "pg_stat_statements";

DROP EXTENSION IF EXISTS "uuid-ossp";

-- +goose StatementEnd