-- +goose Up
-- +goose StatementBegin
-- Development seed data for testing
-- This migration should only be run in development environments
-- Insert sample users
INSERT INTO
    users (
        id,
        email,
        password_hash,
        first_name,
        last_name,
        timezone,
        preferences
    )
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440001',
        'john.doe@englog.dev',
        '$2a$12$LQv3c1yqBwAzWx4bSqiJCuFe.0TrVZOSmE7Gj.K8Xr5ZLGq8QJLQ6',
        'John',
        'Doe',
        'America/New_York',
        '{"theme": "dark", "notifications": true}'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440002',
        'jane.smith@englog.dev',
        '$2a$12$LQv3c1yqBwAzWx4bSqiJCuFe.0TrVZOSmE7Gj.K8Xr5ZLGq8QJLQ6',
        'Jane',
        'Smith',
        'Europe/London',
        '{"theme": "light", "notifications": false}'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440003',
        'alex.johnson@englog.dev',
        '$2a$12$LQv3c1yqBwAzWx4bSqiJCuFe.0TrVZOSmE7Gj.K8Xr5ZLGq8QJLQ6',
        'Alex',
        'Johnson',
        'UTC',
        '{"theme": "auto"}'
    )
ON CONFLICT DO NOTHING;

-- Insert sample projects
INSERT INTO
    projects (
        id,
        name,
        description,
        color,
        created_by,
        is_default
    )
VALUES
    (
        '660e8400-e29b-41d4-a716-446655440001',
        'EngLog API',
        'Personal work activity tracker API development',
        '#3498db',
        '550e8400-e29b-41d4-a716-446655440001',
        true
    ),
    (
        '660e8400-e29b-41d4-a716-446655440002',
        'Mobile App',
        'React Native mobile application',
        '#e74c3c',
        '550e8400-e29b-41d4-a716-446655440001',
        false
    ),
    (
        '660e8400-e29b-41d4-a716-446655440003',
        'Data Analytics',
        'Machine learning and analytics pipeline',
        '#2ecc71',
        '550e8400-e29b-41d4-a716-446655440002',
        true
    ),
    (
        '660e8400-e29b-41d4-a716-446655440004',
        'DevOps Infrastructure',
        'CI/CD and infrastructure automation',
        '#f39c12',
        '550e8400-e29b-41d4-a716-446655440003',
        true
    );

-- Insert sample tags
INSERT INTO
    tags (id, name, color, description)
VALUES
    (
        '770e8400-e29b-41d4-a716-446655440001',
        'backend',
        '#3498db',
        'Backend development tasks'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440002',
        'frontend',
        '#e74c3c',
        'Frontend development tasks'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440003',
        'database',
        '#2ecc71',
        'Database related work'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440004',
        'testing',
        '#f39c12',
        'Testing and QA activities'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440005',
        'documentation',
        '#9b59b6',
        'Documentation writing'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440006',
        'bug-fix',
        '#e67e22',
        'Bug fixing activities'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440007',
        'feature',
        '#1abc9c',
        'New feature development'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440008',
        'refactoring',
        '#34495e',
        'Code refactoring tasks'
    );

-- Insert sample log entries (last 7 days of activity)
INSERT INTO
    log_entries (
        id,
        user_id,
        project_id,
        title,
        description,
        type,
        start_time,
        end_time,
        value_rating,
        impact_level
    )
VALUES
    -- John Doe's activities
    (
        '880e8400-e29b-41d4-a716-446655440001',
        '550e8400-e29b-41d4-a716-446655440001',
        '660e8400-e29b-41d4-a716-446655440001',
        'Implement user authentication',
        'Added JWT-based authentication system with refresh tokens',
        'development',
        NOW () - INTERVAL '1 day' - INTERVAL '2 hours',
        NOW () - INTERVAL '1 day' + INTERVAL '30 minutes',
        'high',
        'team'
    ),
    (
        '880e8400-e29b-41d4-a716-446655440002',
        '550e8400-e29b-41d4-a716-446655440001',
        '660e8400-e29b-41d4-a716-446655440001',
        'Database schema design',
        'Designed and implemented PostgreSQL schema for activity tracking',
        'development',
        NOW () - INTERVAL '2 days' - INTERVAL '3 hours',
        NOW () - INTERVAL '2 days' - INTERVAL '30 minutes',
        'critical',
        'company'
    ),
    (
        '880e8400-e29b-41d4-a716-446655440003',
        '550e8400-e29b-41d4-a716-446655440001',
        '660e8400-e29b-41d4-a716-446655440001',
        'API documentation',
        'Created comprehensive Swagger documentation for REST endpoints',
        'documentation',
        NOW () - INTERVAL '3 days' - INTERVAL '1 hour',
        NOW () - INTERVAL '3 days' + INTERVAL '1 hour',
        'medium',
        'team'
    ),
    (
        '880e8400-e29b-41d4-a716-446655440004',
        '550e8400-e29b-41d4-a716-446655440001',
        '660e8400-e29b-41d4-a716-446655440002',
        'React Native setup',
        'Initial setup and configuration of React Native project',
        'development',
        NOW () - INTERVAL '4 days' - INTERVAL '2 hours',
        NOW () - INTERVAL '4 days' + INTERVAL '1 hour',
        'medium',
        'personal'
    ),
    -- Jane Smith's activities
    (
        '880e8400-e29b-41d4-a716-446655440005',
        '550e8400-e29b-41d4-a716-446655440002',
        '660e8400-e29b-41d4-a716-446655440003',
        'Data pipeline implementation',
        'Built ETL pipeline for processing user activity data',
        'development',
        NOW () - INTERVAL '1 day' - INTERVAL '4 hours',
        NOW () - INTERVAL '1 day' - INTERVAL '1 hour',
        'high',
        'department'
    ),
    (
        '880e8400-e29b-41d4-a716-446655440006',
        '550e8400-e29b-41d4-a716-446655440002',
        '660e8400-e29b-41d4-a716-446655440003',
        'ML model training',
        'Trained productivity prediction model using historical data',
        'research',
        NOW () - INTERVAL '2 days' - INTERVAL '5 hours',
        NOW () - INTERVAL '2 days' - INTERVAL '2 hours',
        'critical',
        'company'
    ),
    (
        '880e8400-e29b-41d4-a716-446655440007',
        '550e8400-e29b-41d4-a716-446655440002',
        '660e8400-e29b-41d4-a716-446655440003',
        'Team sync meeting',
        'Weekly team synchronization and planning meeting',
        'meeting',
        NOW () - INTERVAL '3 days' - INTERVAL '1 hour',
        NOW () - INTERVAL '3 days' - INTERVAL '30 minutes',
        'medium',
        'team'
    ),
    -- Alex Johnson's activities
    (
        '880e8400-e29b-41d4-a716-446655440008',
        '550e8400-e29b-41d4-a716-446655440003',
        '660e8400-e29b-41d4-a716-446655440004',
        'Docker containerization',
        'Created Docker containers for API and worker services',
        'deployment',
        NOW () - INTERVAL '1 day' - INTERVAL '3 hours',
        NOW () - INTERVAL '1 day' - INTERVAL '1 hour',
        'high',
        'team'
    ),
    (
        '880e8400-e29b-41d4-a716-446655440009',
        '550e8400-e29b-41d4-a716-446655440003',
        '660e8400-e29b-41d4-a716-446655440004',
        'CI/CD pipeline setup',
        'Configured GitHub Actions for automated testing and deployment',
        'deployment',
        NOW () - INTERVAL '2 days' - INTERVAL '4 hours',
        NOW () - INTERVAL '2 days' - INTERVAL '2 hours',
        'critical',
        'company'
    ),
    (
        '880e8400-e29b-41d4-a716-446655440010',
        '550e8400-e29b-41d4-a716-446655440003',
        '660e8400-e29b-41d4-a716-446655440004',
        'Security audit',
        'Performed security audit and vulnerability assessment',
        'testing',
        NOW () - INTERVAL '3 days' - INTERVAL '2 hours',
        NOW () - INTERVAL '3 days' + INTERVAL '30 minutes',
        'high',
        'company'
    );

-- Insert sample log entry tags
INSERT INTO
    log_entry_tags (log_entry_id, tag_id)
VALUES
    (
        '880e8400-e29b-41d4-a716-446655440001',
        '770e8400-e29b-41d4-a716-446655440001'
    ), -- backend
    (
        '880e8400-e29b-41d4-a716-446655440001',
        '770e8400-e29b-41d4-a716-446655440007'
    ), -- feature
    (
        '880e8400-e29b-41d4-a716-446655440002',
        '770e8400-e29b-41d4-a716-446655440003'
    ), -- database
    (
        '880e8400-e29b-41d4-a716-446655440003',
        '770e8400-e29b-41d4-a716-446655440005'
    ), -- documentation
    (
        '880e8400-e29b-41d4-a716-446655440004',
        '770e8400-e29b-41d4-a716-446655440002'
    ), -- frontend
    (
        '880e8400-e29b-41d4-a716-446655440005',
        '770e8400-e29b-41d4-a716-446655440001'
    ), -- backend
    (
        '880e8400-e29b-41d4-a716-446655440006',
        '770e8400-e29b-41d4-a716-446655440001'
    ), -- backend
    (
        '880e8400-e29b-41d4-a716-446655440008',
        '770e8400-e29b-41d4-a716-446655440001'
    ), -- backend
    (
        '880e8400-e29b-41d4-a716-446655440009',
        '770e8400-e29b-41d4-a716-446655440001'
    ), -- backend
    (
        '880e8400-e29b-41d4-a716-446655440010',
        '770e8400-e29b-41d4-a716-446655440004'
    );

-- testing
-- Insert sample generated insights
INSERT INTO
    generated_insights (
        id,
        user_id,
        report_type,
        period_start,
        period_end,
        title,
        content,
        summary,
        metadata,
        generation_model,
        quality_score
    )
VALUES
    (
        '990e8400-e29b-41d4-a716-446655440001',
        '550e8400-e29b-41d4-a716-446655440001',
        'weekly_summary',
        CURRENT_DATE - INTERVAL '7 days',
        CURRENT_DATE,
        'Weekly Productivity Summary',
        'This week you focused heavily on backend development, spending 67% of your time on the EngLog API project. Your productivity was highest on Tuesday with 4.5 hours of high-value work. Consider balancing feature development with more testing activities.',
        'Strong week with focus on API development. Productivity peaked mid-week.',
        '{"total_hours": 18.5, "top_activity": "development", "productivity_trend": "increasing"}',
        'llama3.1',
        0.87
    ),
    (
        '990e8400-e29b-41d4-a716-446655440002',
        '550e8400-e29b-41d4-a716-446655440002',
        'monthly_summary',
        CURRENT_DATE - INTERVAL '30 days',
        CURRENT_DATE,
        'Monthly Analytics Review',
        'Your data science work this month shows excellent progress in machine learning model development. You spent 45 hours on research and development, with consistent daily contributions. The ML model training task was your highest-impact activity.',
        'Consistent month with strong focus on ML development.',
        '{"total_hours": 72, "research_percentage": 35, "impact_score": 3.4}',
        'llama3.1',
        0.92
    );

-- Insert sample background tasks
INSERT INTO
    tasks (id, task_type, user_id, payload, status, priority)
VALUES
    (
        'aa0e8400-e29b-41d4-a716-446655440001',
        'generate_insight',
        '550e8400-e29b-41d4-a716-446655440001',
        '{"report_type": "daily_summary", "date": "2024-12-01"}',
        'completed',
        3
    ),
    (
        'aa0e8400-e29b-41d4-a716-446655440002',
        'send_email',
        '550e8400-e29b-41d4-a716-446655440002',
        '{"template": "weekly_report", "recipient": "jane.smith@englog.dev"}',
        'pending',
        5
    ),
    (
        'aa0e8400-e29b-41d4-a716-446655440003',
        'export_data',
        '550e8400-e29b-41d4-a716-446655440003',
        '{"format": "csv", "date_range": "last_month"}',
        'processing',
        7
    );

-- Refresh the materialized view with the new data
REFRESH MATERIALIZED VIEW user_activity_summary;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Remove all seed data
DELETE FROM tasks
WHERE
    task_type IN ('generate_insight', 'send_email', 'export_data');

DELETE FROM generated_insights
WHERE
    user_id IN (
        '550e8400-e29b-41d4-a716-446655440001',
        '550e8400-e29b-41d4-a716-446655440002',
        '550e8400-e29b-41d4-a716-446655440003'
    );

DELETE FROM log_entry_tags
WHERE
    log_entry_id LIKE '880e8400-e29b-41d4-a716-44665544%';

DELETE FROM log_entries
WHERE
    user_id IN (
        '550e8400-e29b-41d4-a716-446655440001',
        '550e8400-e29b-41d4-a716-446655440002',
        '550e8400-e29b-41d4-a716-446655440003'
    );

DELETE FROM tags
WHERE
    id LIKE '770e8400-e29b-41d4-a716-44665544%';

DELETE FROM projects
WHERE
    created_by IN (
        '550e8400-e29b-41d4-a716-446655440001',
        '550e8400-e29b-41d4-a716-446655440002',
        '550e8400-e29b-41d4-a716-446655440003'
    );

DELETE FROM users
WHERE
    id IN (
        '550e8400-e29b-41d4-a716-446655440001',
        '550e8400-e29b-41d4-a716-446655440002',
        '550e8400-e29b-41d4-a716-446655440003'
    );

REFRESH MATERIALIZED VIEW user_activity_summary;

-- +goose StatementEnd