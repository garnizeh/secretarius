-- PostgreSQL initialization script for EngLog
-- This script runs when the PostgreSQL container starts for the first time
-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- Create additional users if needed (for read-only access, monitoring, etc.)
-- These can be uncommented and configured as needed
-- Read-only user for analytics
-- CREATE USER englog_readonly WITH PASSWORD 'readonly_password';
-- GRANT CONNECT ON DATABASE englog TO englog_readonly;
-- GRANT USAGE ON SCHEMA public TO englog_readonly;
-- GRANT SELECT ON ALL TABLES IN SCHEMA public TO englog_readonly;
-- ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO englog_readonly;
-- Monitoring user
-- CREATE USER englog_monitor WITH PASSWORD 'monitor_password';
-- GRANT CONNECT ON DATABASE englog TO englog_monitor;
-- GRANT USAGE ON SCHEMA public TO englog_monitor;
-- GRANT SELECT ON pg_stat_database, pg_stat_user_tables, pg_stat_user_indexes TO englog_monitor;
-- Set some PostgreSQL configuration for better performance
ALTER SYSTEM
SET
    shared_preload_libraries = 'pg_stat_statements';

ALTER SYSTEM
SET
    log_statement = 'mod';

ALTER SYSTEM
SET
    log_duration = on;

ALTER SYSTEM
SET
    log_min_duration_statement = 1000;

-- Note: Configuration changes require a restart to take effect
-- This is handled automatically in the Docker container