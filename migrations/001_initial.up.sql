-- Initial migration for EngLog
-- This creates a basic health check table
CREATE TABLE
    IF NOT EXISTS health_check (
        id SERIAL PRIMARY KEY,
        status VARCHAR(20) DEFAULT 'healthy',
        created_at TIMESTAMP DEFAULT NOW ()
    );