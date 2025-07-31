-- EngLog Database Health Queries
-- This file contains queries for health checking and system status
-- name: GetSystemHealth :one
SELECT
    'healthy' as status,
    NOW () as timestamp,
    current_database () as database_name,
    version () as database_version;