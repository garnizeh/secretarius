# EngLog: Personal Work Activity Tracker for Software Engineers

## Executive Summary

EngLog is a specialized personal organizer designed for software engineers to capture, analyze, and derive insights from their daily work activities. The application serves as a comprehensive logging system that facilitates performance reviews, identifies professional growth patterns, and provides data-driven insights using LLM technology.

## Table of Contents

1. [System Architecture Overview](#system-architecture-overview)
2. [Product Overview](#product-overview)
3. [Core Features & Requirements](#core-features--requirements)
4. [Technical Architecture](#technical-architecture)
5. [Implementation Details](#implementation-details)
6. [Security & Privacy](#security--privacy)
7. [Performance & Scalability](#performance--scalability)
8. [Deployment & Operations](#deployment--operations)
9. [Example Usage Scenarios](#example-usage-scenarios)
10. [Future Roadmap](#future-roadmap)

---

## System Architecture Overview

### Distributed Architecture Design

EngLog is designed as a distributed system with clear separation between client and server responsibilities, optimized for scalability and maintainability.

#### Client Side
- **Current Phase**: REST API with Swagger documentation
- **Future Phase**: React-based web application frontend
- **API Documentation**: Comprehensive Swagger/OpenAPI specification
- **Authentication**: JWT token-based authentication

#### Server Side Architecture

The server-side consists of two distinct services communicating via gRPC:

##### 1. API Server (Public Machine)
- **Purpose**: Handles all client requests and API endpoints
- **Network**: Public IP address with direct internet access
- **Responsibilities**:
  - User authentication and authorization
  - REST API endpoint handling
  - Request validation and response formatting
  - Client session management
  - Database operations (PostgreSQL and Redis)
- **Database Access**: Direct PostgreSQL and Redis connections
- **Role**: gRPC Server (receives requests from Worker)

##### 2. Worker Server (Private Machine)
- **Purpose**: Background processing and LLM operations
- **Network**: Private IP address, accessible only from API Server
- **Responsibilities**:
  - Ollama LLM integration and AI processing
  - Background task processing (reports, emails, analytics)
  - Scheduled job execution
  - Data processing and computation-heavy operations
  - Email notifications and external integrations
- **Database Access**: Read-only access to PostgreSQL for data processing
- **Role**: gRPC Client (connects to API Server)

##### Inter-Machine Communication
- **Protocol**: gRPC over private network
- **Direction**: Worker Server → API Server (Worker as gRPC client)
- **Security**: Simple TLS encryption between machines
- **Purpose**: Worker requests data and sends processing results back to API Server

### Architecture Benefits

1. **Simplicity**: Two-machine setup easy to deploy and maintain
2. **Security**: Worker machine isolated from public internet
3. **Cost-Effective**: Minimal infrastructure requirements
4. **Performance**: Dedicated machines for specific functions
5. **Scalability**: Clear path for future horizontal scaling

---

## Product Overview

### Mission Statement
Empower software engineers to track their professional activities systematically, enabling data-driven career development and streamlined performance evaluations.

### Target Users
- Software Engineers (Junior to Senior)
- Engineering Managers
- Technical Leads
- DevOps Engineers
- QA Engineers

### Key Value Propositions
1. **Automated Performance Review Generation**: Transform daily logs into comprehensive performance summaries
2. **Pattern Recognition**: Identify productivity patterns and areas for improvement
3. **Time Allocation Analysis**: Understand how time is distributed across different activities
4. **Career Growth Tracking**: Monitor skill development and project contributions over time

---

## Core Features & Requirements

### 1. Activity Logging System

#### Log Entry Model
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "project_id": "uuid",
  "title": "string (required, max 200 chars)",
  "description": "string (optional, max 2000 chars)",
  "type": "enum (Meeting, Code, Debug, Docs, Mentoring, Review, Planning, Learning, Ops, Admin)",
  "tags": ["array of strings"],
  "start_time": "timestamp (required)",
  "end_time": "timestamp (required)",
  "duration_minutes": "integer (computed)",
  "value_rating": "enum (Low, Medium, High) - perceived value of the activity",
  "impact_level": "enum (Personal, Team, Product, Company) - scope of impact",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

#### Activity Types
- **Meeting**: Stand-ups, Sprint Planning, Reviews, 1-on-1s, Architecture Reviews
- **Code**: Feature development, bug fixes, refactoring, prototyping
- **Debug**: Production issues, troubleshooting, incident response
- **Docs**: API documentation, architecture diagrams, technical writing
- **Mentoring**: Pair programming, code reviews for juniors, knowledge sharing
- **Review**: Code reviews, PR reviews, design reviews
- **Planning**: Sprint planning, backlog refinement, technical discovery
- **Learning**: Training, research, self-directed learning
- **Ops**: Deployments, monitoring, infrastructure work
- **Admin**: Jira updates, administrative tasks, process work

#### Value & Impact Assessment

##### Value Rating
Subjective assessment of how valuable the activity was:
- **Low**: Routine tasks, administrative work, or activities with minimal learning/output
- **Medium**: Standard development work, regular meetings, normal problem-solving
- **High**: Critical bug fixes, breakthrough solutions, high-impact features, mentoring success

##### Impact Level
Scope of who benefits from this activity:
- **Personal**: Self-improvement, individual learning, personal skill development
- **Team**: Helps team members, improves team processes, knowledge sharing within team
- **Product**: Directly affects product quality, user experience, or product capabilities
- **Company**: Influences company-wide processes, standards, or strategic initiatives

### 2. Project Management

#### Project Model
```json
{
  "id": "uuid",
  "name": "string (required)",
  "description": "string",
  "status": "enum (Active, Completed, On Hold, Cancelled)",
  "start_date": "date",
  "end_date": "date (optional)",
  "team_id": "uuid (optional)",
  "created_by": "uuid",
  "created_at": "timestamp"
}
```

### 3. Intelligent Tagging System

#### Predefined Tag Categories
- **Priority**: `p0-critical`, `p1-high`, `p2-medium`, `p3-low`
- **Technology**: `api-design`, `database`, `backend`, `devops`
- **Activity Focus**: `bug-fix`, `feature`, `refactor`, `performance`, `security`, `testing`
- **Collaboration**: `junior-dev`, `cross-team`, `external-stakeholder`, `solo-work`
- **Skills**: `leadership`, `technical-design`, `problem-solving`, `communication`

### 4. LLM-Powered Analytics

#### Report Types
1. **Daily Summary**: End-of-day activity recap
2. **Weekly Report**: Comprehensive week overview with key achievements
3. **Sprint Analysis**: Sprint performance and contribution analysis
4. **Monthly Insights**: Pattern analysis and productivity trends
5. **Quarterly Review**: Performance review ready document
6. **Annual Assessment**: Comprehensive career growth analysis

#### Insight Categories
- **Time Allocation**: Distribution across activity types and projects
- **Productivity Patterns**: Peak performance times and efficiency trends
- **Skill Development**: Growth in different technical and soft skills
- **Collaboration Impact**: Mentoring and cross-team contributions
- **Goal Progress**: Alignment with personal and team objectives
- **Value Distribution**: Analysis of high-value vs low-value activities ratio
- **Impact Assessment**: Breakdown of personal, team, product, and company-level contributions

### 5. Data Export & Integration

#### Export Formats
- PDF reports for performance reviews
- CSV data for external analysis
- JSON API for integration with other tools
- Calendar integration (Google Calendar, Outlook)

#### Import Capabilities
- Calendar events import
- Jira ticket integration
- GitHub activity sync
- Slack status integration

---

## Technical Architecture

### Distributed System Design

EngLog implements a microservices architecture with two primary components:

#### 1. API Server (Public-Facing)
```
┌─────────────────────────────┐
│        API Server           │
│     (Public Network)        │
├─────────────────────────────┤
│ • REST API Endpoints        │
│ • Authentication Service    │
│ • Request Validation        │
│ • Rate Limiting             │
│ • Database Access           │
│ • Session Management        │
│ • gRPC Server               │
└─────────────────────────────┘
```

#### 2. Worker Server (Private Network)
```
┌─────────────────────────────┐
│       Worker Server         │
│     (Private Network)       │
├─────────────────────────────┤
│ • Ollama LLM Integration    │
│ • Background Task Queue     │
│ • Email Notifications       │
│ • Report Generation         │
│ • Analytics Processing      │
│ • Scheduled Jobs            │
│ • gRPC Client               │
└─────────────────────────────┘
```

#### Communication Flow
```
Client Request → API Server → Worker Server (via gRPC) → Processing → Response
```

### Technology Stack

#### Client Side
- **Current**: Swagger/OpenAPI Documentation
- **Future**: React.js Frontend Application
- **Authentication**: JWT Bearer Tokens
- **API Communication**: REST over HTTPS

#### API Server (Public)
- **Language**: Go 1.24.5+
- **Framework**: Gin HTTP Framework
- **Database**: PostgreSQL 17+ (Read/Write)
- **Cache**: Redis 7+ for sessions and caching
- **Database Driver**: github.com/jackc/pgx/v5
- **Query Builder**: sqlc for type-safe SQL queries
- **Authentication**: JWT ftokens
- **gRPC**: Protocol Buffers + gRPC Server
- **API Documentation**: Swagger with swaggo
- **Validation**: go-playground/validator

#### Worker Server (Private)
- **Language**: Go 1.24.5+
- **Framework**: Standalone gRPC client application
- **Database**: PostgreSQL 17+ (Read-Only)
- **LLM**: Ollama with local model execution
- **Task Queue**: In-memory queue with persistence
- **gRPC**: Protocol Buffers + gRPC Client
- **Background Jobs**: Cron-based scheduling
- **Email**: SMTP integration for notifications

#### Inter-Service Communication
- **Protocol**: gRPC over HTTP/2
- **Security**: Mutual TLS (mTLS) authentication
- **Message Format**: Protocol Buffers (protobuf)
- **Error Handling**: gRPC status codes and error details
- **Load Balancing**: Client-side load balancing for multiple workers

### Simplified Network Architecture

```
Internet
    │
    │ HTTPS
    ▼
┌─────────────────┐
│ Machine 1       │ Public IP: 203.0.113.10
│ API Server      │
│ - Go API        │
│ - PostgreSQL    │
│ - Redis         │
│ - gRPC Server   │
└─────────┬───────┘
          │
          │ gRPC/TLS
          │ Private Network
          ▼
┌─────────────────┐
│ Machine 2       │ Private IP: 10.0.1.20
│ Worker Server   │
│ - Go Worker     │
│ - Ollama LLM    │
│ - gRPC Client   │
│ - Task Engine   │
└─────────────────┘
```

#### Network Configuration

- **Machine 1 (API Server)**:
  - Public IP: Accessible from internet
  - Ports: 80 (HTTP redirect), 443 (HTTPS), 5432 (PostgreSQL), 6379 (Redis)
  - Security Group: Allow HTTP/HTTPS from anywhere, gRPC from Machine 2 only

- **Machine 2 (Worker Server)**:
  - Private IP: No direct internet access
  - Ports: 50051 (gRPC client), 11434 (Ollama)
  - Security Group: Allow gRPC connections to Machine 1 only

- **Communication**:
  - Client ↔ API Server: HTTPS (port 443)
  - API Server ↔ Worker: gRPC with TLS (port 50051)
  - API Server ↔ Database: Local connections

#### Database Schema

```sql
-- Users and Authentication
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    timezone VARCHAR(50) DEFAULT 'UTC',
    preferences JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    access_token_hash VARCHAR(255) NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Refresh token denylist for JWT security
CREATE TABLE refresh_token_denylist (
    jti VARCHAR(255) PRIMARY KEY,
    denylisted_at TIMESTAMP DEFAULT NOW()
);

-- Add index for efficient denylist lookups
CREATE INDEX idx_refresh_token_denylist_time ON refresh_token_denylist(denylisted_at);

-- Add TTL-like cleanup (optional, can be done via cron job)
-- This is a PostgreSQL specific approach using a trigger or scheduled job

-- Teams and Projects
CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_teams (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, team_id)
);

CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'active',
    start_date DATE,
    end_date DATE,
    team_id UUID REFERENCES teams(id),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Activity Logging
CREATE TYPE activity_type AS ENUM (
    'meeting', 'code', 'debug', 'docs', 'mentoring',
    'review', 'planning', 'learning', 'ops', 'admin'
);

CREATE TYPE value_rating AS ENUM ('low', 'medium', 'high');

CREATE TYPE impact_level AS ENUM ('personal', 'team', 'product', 'company');

CREATE TABLE log_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    type activity_type NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    duration_minutes INTEGER GENERATED ALWAYS AS (
        EXTRACT(EPOCH FROM (end_time - start_time)) / 60
    ) STORED,
    value_rating value_rating DEFAULT 'medium',
    impact_level impact_level DEFAULT 'personal',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT valid_time_range CHECK (end_time > start_time)
);

CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    category VARCHAR(50),
    color VARCHAR(7), -- Hex color code
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE log_entry_tags (
    log_entry_id UUID REFERENCES log_entries(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (log_entry_id, tag_id)
);

-- Analytics and Insights
CREATE TABLE generated_insights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    report_type VARCHAR(50) NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    content TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for Performance
CREATE INDEX idx_log_entries_user_time ON log_entries(user_id, start_time DESC);
CREATE INDEX idx_log_entries_project ON log_entries(project_id);
CREATE INDEX idx_log_entries_type ON log_entries(type);
CREATE INDEX idx_log_entries_duration ON log_entries(duration_minutes);
CREATE INDEX idx_log_entries_value ON log_entries(value_rating);
CREATE INDEX idx_log_entries_impact ON log_entries(impact_level);
CREATE INDEX idx_log_entries_value_impact ON log_entries(value_rating, impact_level);
CREATE INDEX idx_log_entry_tags_entry ON log_entry_tags(log_entry_id);
CREATE INDEX idx_log_entry_tags_tag ON log_entry_tags(tag_id);
CREATE INDEX idx_insights_user_period ON generated_insights(user_id, period_start, period_end);
```

### SQLC Configuration

The application uses sqlc for type-safe SQL query generation. Here are the main queries needed:

#### sqlc.yaml Configuration
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/queries"
    schema: "migrations"
    gen:
      go:
        package: "sqlc"
        out: "internal/sqlc"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: false
        emit_exact_table_names: false
```

#### Core Queries (queries/log_entries.sql)
```sql
-- name: CreateLogEntry :one
INSERT INTO log_entries (
    user_id, project_id, title, description, type,
    start_time, end_time, value_rating, impact_level
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetLogEntriesByUserAndDateRange :many
SELECT * FROM log_entries
WHERE user_id = $1
  AND start_time >= $2
  AND end_time <= $3
ORDER BY start_time ASC;

-- name: GetLogEntriesByUser :many
SELECT * FROM log_entries
WHERE user_id = $1
ORDER BY start_time DESC;

-- name: UpdateLogEntry :one
UPDATE log_entries
SET title = $2, description = $3, type = $4,
    start_time = $5, end_time = $6, value_rating = $7,
    impact_level = $8, updated_at = NOW()
WHERE id = $1 AND user_id = $9
RETURNING *;

-- name: DeleteLogEntry :exec
DELETE FROM log_entries
WHERE id = $1 AND user_id = $2;

-- name: DeleteLogEntriesByUser :exec
DELETE FROM log_entries
WHERE user_id = $1;
```

#### User & Authentication Queries (queries/users.sql)
```sql
-- name: CreateUser :one
INSERT INTO users (
    id, email, password_hash, first_name, last_name, timezone
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserProfile :one
UPDATE users
SET first_name = $2, last_name = $3, timezone = $4,
    preferences = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetAllUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: GetActiveUserCount :one
SELECT COUNT(DISTINCT user_id) as count
FROM log_entries
WHERE created_at > $1;
```

#### Refresh Token Blacklist Queries (queries/auth.sql)
```sql
-- name: CreateRefreshTokenDenylist :exec
INSERT INTO refresh_token_denylist (jti, denylisted_at)
VALUES ($1, NOW());

-- name: IsRefreshTokenDenylisted :one
SELECT EXISTS(
    SELECT 1 FROM refresh_token_denylist
    WHERE jti = $1
) as exists;

-- name: CleanupExpiredDenylistedTokens :exec
DELETE FROM refresh_token_denylist
WHERE denylisted_at < NOW() - INTERVAL '30 days';

-- name: GetDenylistedTokenCount :one
SELECT COUNT(*) as count
FROM refresh_token_denylist;
```

#### Analytics & Metrics Queries (queries/analytics.sql)
```sql
-- name: CreateUser :one
INSERT INTO users (
    email, password_hash, first_name, last_name, timezone
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetAllUsers :many
SELECT * FROM users ORDER BY created_at;

-- name: UpdateUserProfile :one
UPDATE users
SET first_name = $2, last_name = $3, timezone = $4,
    preferences = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: CreateRefreshTokenDenylist :exec
INSERT INTO refresh_token_denylist (jti, created_at)
VALUES ($1, NOW());

-- name: IsRefreshTokenDenylisted :one
SELECT EXISTS(SELECT 1 FROM refresh_token_denylist WHERE jti = $1);
```

#### Insights & Analytics Queries (queries/insights.sql)
```sql
-- name: CreateInsight :exec
INSERT INTO generated_insights (
    user_id, report_type, period_start, period_end, content, metadata
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: GetInsightsByUser :many
SELECT * FROM generated_insights
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteInsightsByUser :exec
DELETE FROM generated_insights
WHERE user_id = $1;

-- name: DeleteOldInsights :exec
DELETE FROM generated_insights
WHERE created_at < $1;
```

#### Task Management Queries (queries/tasks.sql)
```sql
-- name: CreateTask :exec
INSERT INTO tasks (
    id, type, user_id, payload, status, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: UpdateTaskStatus :exec
UPDATE tasks
SET status = $2, result = $3, error = $4, updated_at = NOW()
WHERE id = $1;

-- name: DeleteOldTasks :exec
DELETE FROM tasks
WHERE created_at < $1 AND status IN ('completed', 'failed');
```

#### GDPR Queries (queries/gdpr.sql)
```sql
-- name: CreateScheduledDeletion :exec
INSERT INTO scheduled_deletions (
    user_id, scheduled_at, status
) VALUES (
    $1, $2, $3
);

-- name: GetPendingScheduledDeletions :many
SELECT * FROM scheduled_deletions
WHERE scheduled_at <= NOW() AND status = 'pending';

-- name: UpdateScheduledDeletionStatus :exec
UPDATE scheduled_deletions
SET status = $2, completed_at = CASE WHEN $2 = 'completed' THEN NOW() ELSE NULL END
WHERE id = $1;

-- name: DeleteLogEntryTagsByUser :exec
DELETE FROM log_entry_tags
WHERE log_entry_id IN (SELECT id FROM log_entries WHERE user_id = $1);

-- name: DeleteUserTeamsByUser :exec
DELETE FROM user_teams WHERE user_id = $1;

-- name: DeleteUserSessionsByUser :exec
DELETE FROM user_sessions WHERE user_id = $1;

-- name: AnonymizeLogDescriptions :exec
UPDATE log_entries
SET description = $2
WHERE user_id = $1;
```

### API Design

#### Core Endpoints

```python
# Authentication
POST   /v1/auth/login
POST   /v1/auth/refresh
POST   /v1/auth/logout
GET    /v1/auth/me

# Log Entries
POST   /v1/logs
GET    /v1/logs?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&project_id=uuid&type=activity_type&tags=tag1,tag2&value_rating=high&impact_level=company
GET    /v1/logs/{log_id}
PUT    /v1/logs/{log_id}
DELETE /v1/logs/{log_id}
POST   /v1/logs/bulk  # Bulk operations

# Projects
GET    /v1/projects
POST   /v1/projects
GET    /v1/projects/{project_id}
PUT    /v1/projects/{project_id}
DELETE /v1/projects/{project_id}
GET    /v1/projects/{project_id}/stats

# Tags
GET    /v1/tags
POST   /v1/tags
GET    /v1/tags/popular
DELETE /v1/tags/{tag_id}

# Analytics & Insights
POST   /v1/insights/generate
GET    /v1/insights?type=report_type&period_start=date&period_end=date
GET    /v1/insights/{insight_id}
GET    /v1/analytics/dashboard
GET    /v1/analytics/time-distribution
GET    /v1/analytics/productivity-trends

# Data Export
GET    /v1/export/pdf?period_start=date&period_end=date
GET    /v1/export/csv?period_start=date&period_end=date
POST   /v1/import/calendar

# Admin & Settings
GET    /v1/users/profile
PUT    /v1/users/profile
GET    /v1/users/preferences
PUT    /v1/users/preferences
```

#### Request/Response Models

```go
package models

import (
    "time"
    "github.com/go-playground/validator/v10"
)

type ActivityType string

const (
    ActivityMeeting   ActivityType = "meeting"
    ActivityCode      ActivityType = "code"
    ActivityDebug     ActivityType = "debug"
    ActivityDocs      ActivityType = "docs"
    ActivityMentoring ActivityType = "mentoring"
    ActivityReview    ActivityType = "review"
    ActivityPlanning  ActivityType = "planning"
    ActivityLearning  ActivityType = "learning"
    ActivityOps       ActivityType = "ops"
    ActivityAdmin     ActivityType = "admin"
)

type ValueRating string

const (
    ValueLow    ValueRating = "low"
    ValueMedium ValueRating = "medium"
    ValueHigh   ValueRating = "high"
)

type ImpactLevel string

const (
    ImpactPersonal ImpactLevel = "personal"
    ImpactTeam     ImpactLevel = "team"
    ImpactProduct  ImpactLevel = "product"
    ImpactCompany  ImpactLevel = "company"
)

type LogEntryCreate struct {
    Title       string       `json:"title" validate:"required,max=200"`
    Description *string      `json:"description,omitempty" validate:"omitempty,max=2000"`
    Type        ActivityType `json:"type" validate:"required,oneof=meeting code debug docs mentoring review planning learning ops admin"`
    ProjectID   *string      `json:"project_id,omitempty"`
    Tags        []string     `json:"tags"`
    StartTime   time.Time    `json:"start_time" validate:"required"`
    EndTime     time.Time    `json:"end_time" validate:"required,gtfield=StartTime"`
    ValueRating ValueRating  `json:"value_rating" validate:"omitempty,oneof=low medium high"`
    ImpactLevel ImpactLevel  `json:"impact_level" validate:"omitempty,oneof=personal team product company"`
}

type LogEntryResponse struct {
    ID              string       `json:"id"`
    Title           string       `json:"title"`
    Description     *string      `json:"description,omitempty"`
    Type            ActivityType `json:"type"`
    ProjectID       *string      `json:"project_id,omitempty"`
    ProjectName     *string      `json:"project_name,omitempty"`
    Tags            []string     `json:"tags"`
    StartTime       time.Time    `json:"start_time"`
    EndTime         time.Time    `json:"end_time"`
    DurationMinutes int          `json:"duration_minutes"`
    ValueRating     ValueRating  `json:"value_rating"`
    ImpactLevel     ImpactLevel  `json:"impact_level"`
    CreatedAt       time.Time    `json:"created_at"`
    UpdatedAt       time.Time    `json:"updated_at"`
}

type InsightRequest struct {
    ReportType    string     `json:"report_type" validate:"required,oneof=daily weekly sprint monthly quarterly annual"`
    PeriodStart   time.Time  `json:"period_start" validate:"required"`
    PeriodEnd     time.Time  `json:"period_end" validate:"required,gtfield=PeriodStart"`
    IncludeProjects []string `json:"include_projects,omitempty"`
    FocusAreas    []string   `json:"focus_areas,omitempty"`
}

type APIResponse struct {
    Success bool        `json:"success"`
    Data    any `json:"data,omitempty"`
    Error   *string     `json:"error,omitempty"`
    Message string      `json:"message,omitempty"`
}
```

---

## Implementation Details

### System Architecture Diagrams

#### Database Entity Relationships

```mermaid
erDiagram
    USERS {
        uuid id PK
        string email UK
        string password_hash
        string first_name
        string last_name
        string timezone
        jsonb preferences
        timestamp created_at
        timestamp updated_at
    }

    PROJECTS {
        uuid id PK
        uuid creator_id FK
        string name
        string description
        string color
        boolean is_default
        timestamp created_at
        timestamp updated_at
    }

    LOG_ENTRIES {
        uuid id PK
        uuid user_id FK
        uuid project_id FK
        string title
        text description
        activity_type type
        timestamp start_time
        timestamp end_time
        integer duration_minutes
        value_rating value_rating
        impact_level impact_level
        timestamp created_at
        timestamp updated_at
    }

    TAGS {
        uuid id PK
        string name UK
        string color
        timestamp created_at
    }

    LOG_ENTRY_TAGS {
        uuid log_entry_id FK
        uuid tag_id FK
    }

    GENERATED_INSIGHTS {
        uuid id PK
        uuid user_id FK
        string report_type
        timestamp period_start
        timestamp period_end
        text content
        jsonb metadata
        timestamp created_at
    }

    TASKS {
        uuid id PK
        string type
        uuid user_id FK
        jsonb payload
        string status
        text result
        text error
        timestamp created_at
        timestamp updated_at
    }

    SCHEDULED_DELETIONS {
        uuid id PK
        uuid user_id FK
        timestamp scheduled_at
        string status
        timestamp completed_at
        timestamp created_at
    }

    REFRESH_TOKEN_DENYLIST {
        string jti PK
        timestamp created_at
    }

    USERS ||--o{ PROJECTS : creates
    USERS ||--o{ LOG_ENTRIES : owns
    USERS ||--o{ GENERATED_INSIGHTS : receives
    USERS ||--o{ TASKS : triggers
    USERS ||--o{ SCHEDULED_DELETIONS : schedules

    PROJECTS ||--o{ LOG_ENTRIES : contains

    LOG_ENTRIES ||--o{ LOG_ENTRY_TAGS : tagged_with
    TAGS ||--o{ LOG_ENTRY_TAGS : applied_to
```

#### Simplified Two-Machine Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        SWAGGER[Swagger Documentation]
        FUTURE_WEB[Future: React Web App]
        API_CLIENT[API Clients]
    end

    subgraph "Machine 1: API Server (Public IP)"
        direction TB
        API_GATEWAY[API Gateway]
        AUTH_SVC[Authentication Service]
        USER_SVC[User Management]
        LOG_SVC[Log Entry Service]
        PROJECT_SVC[Project Service]
        GRPC_SERVER[gRPC Server]
    end

    subgraph "Machine 2: Worker Server (Private IP)"
        direction TB
        GRPC_CLIENT[gRPC Client]
        TASK_ENGINE[Task Processing Engine]
        INSIGHT_ENGINE[Insight Generation]
        OLLAMA_CLIENT[Ollama LLM Client]
        EMAIL_SVC[Email Service]
        SCHEDULER[Background Scheduler]
        ANALYTICS_PROC[Analytics Processor]
    end

    subgraph "LLM Infrastructure (Machine 2)"
        OLLAMA[Ollama Server]
        LLM_MODELS[Local LLM Models]
    end

    subgraph "Data Layer (Machine 1)"
        POSTGRES[(PostgreSQL)]
        REDIS[(Redis Cache)]
    end

    %% Client to API Server
    SWAGGER --> API_GATEWAY
    FUTURE_WEB --> API_GATEWAY
    API_CLIENT --> API_GATEWAY

    %% API Server Internal Flow
    API_GATEWAY --> AUTH_SVC
    API_GATEWAY --> USER_SVC
    API_GATEWAY --> LOG_SVC
    API_GATEWAY --> PROJECT_SVC

    %% gRPC Communication
    GRPC_CLIENT -.->|gRPC/TLS| GRPC_SERVER

    %% Worker Server Internal Flow
    GRPC_CLIENT --> TASK_ENGINE
    GRPC_CLIENT --> INSIGHT_ENGINE
    TASK_ENGINE --> OLLAMA_CLIENT
    INSIGHT_ENGINE --> OLLAMA_CLIENT
    SCHEDULER --> ANALYTICS_PROC
    SCHEDULER --> EMAIL_SVC

    %% LLM Integration
    OLLAMA_CLIENT --> OLLAMA
    OLLAMA --> LLM_MODELS

    %% Database Connections
    AUTH_SVC --> POSTGRES
    USER_SVC --> POSTGRES
    LOG_SVC --> POSTGRES
    PROJECT_SVC --> POSTGRES
    USER_SVC --> REDIS
    LOG_SVC --> REDIS

    %% Styling
    classDef publicMachine fill:#e1f5fe
    classDef privateMachine fill:#f3e5f5
    classDef dataStore fill:#fff3e0
    classDef client fill:#e8f5e8

    class API_GATEWAY,AUTH_SVC,USER_SVC,LOG_SVC,PROJECT_SVC,GRPC_SERVER publicMachine
    class GRPC_CLIENT,TASK_ENGINE,INSIGHT_ENGINE,OLLAMA_CLIENT,EMAIL_SVC,SCHEDULER,ANALYTICS_PROC,OLLAMA,LLM_MODELS privateMachine
    class POSTGRES,REDIS dataStore
```

#### gRPC Service Communication

```mermaid
sequenceDiagram
    participant Client
    participant API_Server as API Server (Machine 1)
    participant Worker_Server as Worker Server (Machine 2)
    participant Ollama

    Client->>API_Server: POST /api/insights/generate
    API_Server->>API_Server: Validate request & authenticate
    API_Server->>API_Server: Store in database

    Note over API_Server, Worker_Server: gRPC over TLS
    API_Server->>Worker_Server: GenerateInsightRequest(userID, period)

    Worker_Server->>Worker_Server: Fetch data via gRPC
    Worker_Server->>Ollama: Generate insight with LLM
    Ollama-->>Worker_Server: AI-generated insight

    Worker_Server->>Worker_Server: Process and format insight
    Worker_Server-->>API_Server: InsightResponse(insight_data)

    API_Server->>API_Server: Save insight to database
    API_Server-->>Client: 201 Created + insight_id

    Note over API_Server, Worker_Server: Background processing continues
    Worker_Server->>API_Server: NotifyInsightReady(insight_id)
    API_Server->>API_Server: Update insight status
```

#### Component Responsibilities

```mermaid
graph LR
    subgraph "API Server Responsibilities"
        A1[Request Handling]
        A2[Authentication]
        A3[Data Validation]
        A4[Database Operations]
        A5[Response Formatting]
        A6[Rate Limiting]
        A7[Session Management]
    end

    subgraph "Worker Server Responsibilities"
        W1[LLM Processing]
        W2[Background Tasks]
        W3[Email Sending]
        W4[Report Generation]
        W5[Analytics Processing]
        W6[Scheduled Jobs]
        W7[Heavy Computations]
    end

    subgraph "Shared Data"
        D1[(User Data)]
        D2[(Log Entries)]
        D3[(Projects)]
        D4[(Generated Insights)]
    end

    A1 --> A2
    A2 --> A3
    A3 --> A4
    A4 --> A5

    W1 --> W2
    W2 --> W3
    W3 --> W4
    W4 --> W5
    W5 --> W6

    A4 -.-> D1
    A4 -.-> D2
    A4 -.-> D3
    A4 -.-> D4

    W1 -.-> D2
    W4 -.-> D4
    W5 -.-> D1
    W5 -.-> D2
```

    SQLC --> POSTGRES

    INSIGHT_SVC --> OLLAMA
    OLLAMA --> LLM_MODELS

    TASK_QUEUE --> POSTGRES
    SCHEDULER --> TASK_QUEUE

    AUTH --> METRICS
    LOG_SVC --> METRICS
    INSIGHT_SVC --> METRICS
    ANALYTICS --> METRICS

    HEALTH --> POSTGRES
    HEALTH --> REDIS
    HEALTH --> OLLAMA
```

#### Authentication Flow

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant AuthService
    participant Database
    participant Redis

    Client->>API: POST /auth/login
    API->>AuthService: validateCredentials(email, password)
    AuthService->>Database: getUserByEmail(email)
    Database-->>AuthService: user data
    AuthService->>AuthService: checkPassword(password, hash)

    alt Valid credentials
        AuthService->>AuthService: createAccessToken(userID)
        AuthService->>AuthService: createRefreshToken(userID)
        AuthService->>Redis: cache user session
        AuthService-->>API: tokens
        API-->>Client: 200 OK + tokens
    else Invalid credentials
        AuthService-->>API: error
        API-->>Client: 401 Unauthorized
    end

    Note over Client: Subsequent requests
    Client->>API: GET /api/logs (with Bearer token)
    API->>AuthService: validateToken(accessToken)

    alt Valid token
        AuthService-->>API: user claims
        API->>API: process request
        API-->>Client: 200 OK + data
    else Expired/Invalid token
        AuthService-->>API: error
        API-->>Client: 401 Unauthorized

        Note over Client: Token refresh flow
        Client->>API: POST /auth/refresh (with refresh token)
        API->>AuthService: rotateRefreshToken(oldToken)
        AuthService->>Database: blacklist old token
        AuthService->>AuthService: create new tokens
        AuthService-->>API: new tokens
        API-->>Client: 200 OK + new tokens
    end
```

#### JWT Token Lifecycle & Security

```mermaid
flowchart TD
    START([User Registration/Login]) --> CREATE_TOKENS[Create JWT Tokens]
    CREATE_TOKENS --> ACCESS_TOKEN[Access Token<br/>15 min TTL<br/>No JTI]
    CREATE_TOKENS --> REFRESH_TOKEN[Refresh Token<br/>30 days TTL<br/>With JTI]

    ACCESS_TOKEN --> USE_ACCESS[Use Access Token<br/>for API requests]
    USE_ACCESS --> CHECK_EXPIRED{Token Expired?}

    CHECK_EXPIRED -->|No| API_REQUEST[Process API Request]
    CHECK_EXPIRED -->|Yes| NEED_REFRESH[Need Token Refresh]

    NEED_REFRESH --> USE_REFRESH[Use Refresh Token]
    USE_REFRESH --> VALIDATE_REFRESH{Refresh Token Valid?}

    VALIDATE_REFRESH -->|No| LOGIN_REQUIRED[Redirect to Login]
    VALIDATE_REFRESH -->|Yes| BLACKLIST_OLD[Blacklist Old Refresh Token]

    BLACKLIST_OLD --> CREATE_NEW_TOKENS[Create New Token Pair]
    CREATE_NEW_TOKENS --> NEW_ACCESS[New Access Token]
    CREATE_NEW_TOKENS --> NEW_REFRESH[New Refresh Token]

    NEW_ACCESS --> USE_ACCESS
    NEW_REFRESH --> REFRESH_TOKEN

    REFRESH_TOKEN --> LOGOUT{User Logout?}
    LOGOUT -->|Yes| BLACKLIST_REFRESH[Blacklist Refresh Token]
    LOGOUT -->|No| WAIT_EXPIRY[Wait for Expiry]

    BLACKLIST_REFRESH --> TOKEN_INVALID[Token Invalid]
    WAIT_EXPIRY --> AUTO_CLEANUP[Auto Cleanup<br/>After 30 days]
    AUTO_CLEANUP --> TOKEN_INVALID

    LOGIN_REQUIRED --> START
    API_REQUEST --> SUCCESS([Successful Response])
    TOKEN_INVALID --> END([Session Ended])
```

#### Token Security Features

```mermaid
graph LR
    subgraph "Access Token (15 min)"
        AT[JWT Payload]
        AT_SUB[sub: userID]
        AT_TYPE[type: access]
        AT_EXP[exp: 15 min]
        AT_IAT[iat: issued time]

        AT --> AT_SUB
        AT --> AT_TYPE
        AT --> AT_EXP
        AT --> AT_IAT
    end

    subgraph "Refresh Token (30 days)"
        RT[JWT Payload]
        RT_SUB[sub: userID]
        RT_TYPE[type: refresh]
        RT_JTI[jti: unique ID]
        RT_EXP[exp: 30 days]
        RT_IAT[iat: issued time]

        RT --> RT_SUB
        RT --> RT_TYPE
        RT --> RT_JTI
        RT --> RT_EXP
        RT --> RT_IAT
    end

    subgraph "Security Measures"
        BLACKLIST[(Token Blacklist<br/>Database)]
        ROTATION[Automatic Rotation]
        VALIDATION[Signature Validation]
        EXPIRY[Expiration Checking]

        RT_JTI --> BLACKLIST
        ROTATION --> BLACKLIST
        AT --> VALIDATION
        RT --> VALIDATION
        AT_EXP --> EXPIRY
        RT_EXP --> EXPIRY
    end

    subgraph "Attack Prevention"
        REPLAY[Replay Attack<br/>Prevention]
        STOLEN[Stolen Token<br/>Mitigation]
        BRUTE[Brute Force<br/>Protection]

        BLACKLIST --> REPLAY
        ROTATION --> STOLEN
        EXPIRY --> BRUTE
    end
```

#### Log Entry Creation Flow

```mermaid
sequenceDiagram
    participant Client
    participant LogService
    participant ValidationLayer
    participant Database
    participant TaskQueue
    participant MetricsCollector

    Client->>LogService: POST /api/logs/entries
    LogService->>ValidationLayer: validateLogEntry(data)

    alt Valid data
        ValidationLayer-->>LogService: validation passed
        LogService->>Database: createLogEntry(entry)
        Database-->>LogService: created entry

        LogService->>MetricsCollector: recordLogEntryCreated()
        LogService->>TaskQueue: enqueue(updateUserMetrics)

        LogService-->>Client: 201 Created + entry data
    else Invalid data
        ValidationLayer-->>LogService: validation errors
        LogService-->>Client: 400 Bad Request + errors
    end

    Note over TaskQueue: Background processing
    TaskQueue->>TaskQueue: process updateUserMetrics task
    TaskQueue->>Database: update user analytics
```

#### Insight Generation Workflow

```mermaid
flowchart TD
    START([User requests insight]) --> CHECK_CACHE{Check Redis Cache}
    CHECK_CACHE -->|Cache Hit| RETURN_CACHED[Return cached insight]
    CHECK_CACHE -->|Cache Miss| FETCH_LOGS[Fetch log entries from DB]

    FETCH_LOGS --> VALIDATE_DATA{Sufficient data?}
    VALIDATE_DATA -->|No| ERROR[Return insufficient data error]
    VALIDATE_DATA -->|Yes| ENQUEUE_TASK[Enqueue insight generation task]

    ENQUEUE_TASK --> TASK_WORKER[Task worker picks up job]
    TASK_WORKER --> ANALYZE_LOGS[Analyze log entries]
    ANALYZE_LOGS --> BUILD_CONTEXT[Build context for LLM]

    BUILD_CONTEXT --> CALL_LLM[Call Ollama API]
    CALL_LLM --> LLM_RESPONSE{LLM Response OK?}

    LLM_RESPONSE -->|Error| RETRY{Retry count < 3?}
    RETRY -->|Yes| CALL_LLM
    RETRY -->|No| TASK_FAILED[Mark task as failed]

    LLM_RESPONSE -->|Success| PROCESS_RESPONSE[Process LLM response]
    PROCESS_RESPONSE --> SAVE_INSIGHT[Save insight to database]
    SAVE_INSIGHT --> CACHE_RESULT[Cache result in Redis]
    CACHE_RESULT --> NOTIFY_USER[Notify user of completion]

    RETURN_CACHED --> END([Return insight to user])
    ERROR --> END
    TASK_FAILED --> END
    NOTIFY_USER --> END
```

#### Data Processing Pipeline

```mermaid
graph LR
    subgraph "Data Ingestion"
        INPUT[Log Entry Input]
        VALIDATION[Data Validation]
        ENRICHMENT[Data Enrichment]
    end

    subgraph "Real-time Processing"
        METRICS[Metrics Calculation]
        AGGREGATION[Data Aggregation]
        TREND[Trend Analysis]
    end

    subgraph "Background Processing"
        DAILY[Daily Reports]
        WEEKLY[Weekly Summaries]
        CLEANUP[Data Cleanup]
    end

    subgraph "AI Processing"
        CONTEXT[Context Building]
        LLM[LLM Generation]
        INSIGHT[Insight Creation]
    end

    subgraph "Storage"
        PRIMARY[(Primary DB)]
        CACHE[(Redis Cache)]
        ARCHIVE[(Long-term Storage)]
    end

    INPUT --> VALIDATION
    VALIDATION --> ENRICHMENT
    ENRICHMENT --> PRIMARY

    PRIMARY --> METRICS
    METRICS --> AGGREGATION
    AGGREGATION --> TREND
    TREND --> CACHE

    PRIMARY --> DAILY
    DAILY --> WEEKLY
    WEEKLY --> CLEANUP
    CLEANUP --> ARCHIVE

    PRIMARY --> CONTEXT
    CONTEXT --> LLM
    LLM --> INSIGHT
    INSIGHT --> PRIMARY
    INSIGHT --> CACHE
```

#### Simple Two-Machine Deployment

```mermaid
graph TB
    subgraph "Internet"
        CLIENT[API Clients]
        SWAGGER[Swagger UI]
    end

    subgraph "Machine 1: Public Server"
        direction TB
        API[API Server]
        DB[(PostgreSQL)]
        CACHE[(Redis)]
        GRPC_SRV[gRPC Server]
    end

    subgraph "Machine 2: Private Worker"
        direction TB
        WORKER[Worker Server]
        GRPC_CLI[gRPC Client]
        OLLAMA[Ollama LLM]
        MODELS[(LLM Models)]
    end

    subgraph "External Services"
        SMTP[Email SMTP]
    end

    %% Client connections
    CLIENT --> API
    SWAGGER --> API

    %% Internal connections Machine 1
    API --> DB
    API --> CACHE
    API <--> GRPC_SRV

    %% Inter-machine communication
    GRPC_SRV <-.->|gRPC/TLS| GRPC_CLI

    %% Internal connections Machine 2
    GRPC_CLI --> WORKER
    WORKER --> OLLAMA
    OLLAMA --> MODELS

    %% External services
    WORKER --> SMTP

    %% Styling
    classDef publicMachine fill:#e1f5fe
    classDef privateMachine fill:#f3e5f5
    classDef dataStore fill:#fff3e0
    classDef external fill:#fff8e1

    class API,DB,CACHE,GRPC_SRV publicMachine
    class WORKER,GRPC_CLI,OLLAMA,MODELS privateMachine
```

### Docker Compose Configuration

#### Machine 1: API Server (docker-compose.api.yml)

```yaml
version: '3.8'

services:
  api-server:
    build:
      context: .
      dockerfile: cmd/api/Dockerfile
    ports:
      - "80:8080"   # HTTP redirect
      - "443:8443"  # HTTPS
      - "50051:50051"  # gRPC server
    environment:
      - DB_HOST=postgres
      - DB_USER=englog
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=englog
      - REDIS_HOST=redis
      - GRPC_PORT=50051
      - WORKER_HOST=10.0.1.20
      - TLS_CERT_FILE=/certs/server.crt
      - TLS_KEY_FILE=/certs/server.key
    volumes:
      - ./certs:/certs:ro
      - ./logs:/app/logs
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

  postgres:
    image: postgres:17
    environment:
      - POSTGRES_DB=englog
      - POSTGRES_USER=englog
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
```

#### Machine 2: Worker Server (docker-compose.worker.yml)

```yaml
version: '3.8'

services:
  worker-server:
    build:
      context: .
      dockerfile: cmd/worker/Dockerfile
    environment:
      - API_SERVER_HOST=203.0.113.10
      - API_SERVER_PORT=50051
      - OLLAMA_HOST=ollama
      - OLLAMA_PORT=11434
      - TLS_CERT_FILE=/certs/client.crt
      - TLS_KEY_FILE=/certs/client.key
      - TLS_CA_FILE=/certs/ca.crt
    volumes:
      - ./certs:/certs:ro
      - ./logs:/app/logs
    depends_on:
      - ollama
    restart: unless-stopped

  ollama:
    image: ollama/ollama:latest
    ports:
      - "11434:11434"
    volumes:
      - ollama_data:/root/.ollama
    environment:
      - OLLAMA_HOST=0.0.0.0
    restart: unless-stopped
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: all
              capabilities: [gpu]

volumes:
  ollama_data:
```

#### Deployment Commands

```bash
# Machine 1 (API Server)
export DB_PASSWORD="your_secure_password"
docker-compose -f docker-compose.api.yml up -d

# Machine 2 (Worker Server)
docker-compose -f docker-compose.worker.yml up -d

# Initial Ollama model setup on Machine 2
docker exec -it worker_ollama_1 ollama pull llama3.2
```

    %% API Server connections
    SSL --> API1
    SSL --> API2
    SSL --> API3

    %% gRPC connections (Private network)
    API1 -.->|gRPC| WORKER1
    API2 -.->|gRPC| WORKER2
    API3 -.->|gRPC| WORKER1

    %% Database connections
    API1 --> DB_PRIMARY
    API2 --> DB_PRIMARY
    API3 --> DB_PRIMARY

    API1 --> DB_REPLICA1
    API2 --> DB_REPLICA2
    API3 --> DB_REPLICA1

    %% Cache connections
    API1 --> REDIS_PRIMARY
    API2 --> REDIS_PRIMARY
    API3 --> REDIS_PRIMARY

    %% Worker to LLM
    WORKER1 --> OLLAMA1
    WORKER2 --> OLLAMA2

    OLLAMA1 --> MODEL_STORAGE
    OLLAMA2 --> MODEL_STORAGE

    %% Worker data access (read-only)
    WORKER1 --> DB_REPLICA1
    WORKER2 --> DB_REPLICA2

    %% External service connections
    WORKER1 --> SMTP
    WORKER2 --> SMTP
    SMTP --> EMAIL_PROVIDER

    %% Database replication
    DB_PRIMARY --> DB_REPLICA1
    DB_PRIMARY --> DB_REPLICA2
    REDIS_PRIMARY --> REDIS_REPLICA

    %% Monitoring connections
    API1 --> PROMETHEUS
    API2 --> PROMETHEUS
    API3 --> PROMETHEUS
    WORKER1 --> PROMETHEUS
    WORKER2 --> PROMETHEUS

    PROMETHEUS --> GRAFANA
    PROMETHEUS --> ALERTMANAGER

    %% Tracing
    API1 --> JAEGER
    API2 --> JAEGER
    API3 --> JAEGER
    WORKER1 --> JAEGER
    WORKER2 --> JAEGER
```

#### High Availability Setup

```mermaid
graph TB
    subgraph "Multi-Zone Deployment"
        subgraph "Zone A"
            LB_A[Load Balancer A]
            API_A[API Server Cluster A]
            WORKER_A[Worker Server A]
            DB_A[(Database Primary)]
            REDIS_A[(Redis Primary)]
        end

        subgraph "Zone B"
            LB_B[Load Balancer B]
            API_B[API Server Cluster B]
            WORKER_B[Worker Server B]
            DB_B[(Database Replica)]
            REDIS_B[(Redis Replica)]
        end

        subgraph "Zone C"
            API_C[API Server Cluster C]
            WORKER_C[Worker Server C]
            DB_C[(Database Replica)]
        end
    end

    subgraph "Shared Services"
        OLLAMA_CLUSTER[Ollama LLM Cluster]
        MONITORING[Monitoring Stack]
        EMAIL[Email Service]
    end

    %% Load balancer failover
    LB_A -.->|Failover| LB_B

    %% Cross-zone communication
    API_A -.->|gRPC| WORKER_B
    API_B -.->|gRPC| WORKER_A
    API_C -.->|gRPC| WORKER_A

    %% Database replication
    DB_A --> DB_B
    DB_A --> DB_C
    REDIS_A --> REDIS_B

    %% Shared service access
    WORKER_A --> OLLAMA_CLUSTER
    WORKER_B --> OLLAMA_CLUSTER
    WORKER_C --> OLLAMA_CLUSTER

    API_A --> MONITORING
    API_B --> MONITORING
    API_C --> MONITORING
    WORKER_A --> MONITORING
    WORKER_B --> MONITORING
    WORKER_C --> MONITORING

    WORKER_A --> EMAIL
    WORKER_B --> EMAIL
    WORKER_C --> EMAIL
```

#### Security Network Topology

```mermaid
graph TB
    subgraph "Internet (Untrusted)"
        INTERNET[Internet Traffic]
    end

    subgraph "DMZ (Semi-Trusted)"
        WAF[Web Application Firewall]
        LB[Load Balancer]
        BASTION[Bastion Host]
    end

    subgraph "Public App Network (Trusted)"
        API_SERVERS[API Server Cluster]
        API_DB[(API Database)]
        API_CACHE[(API Cache)]
    end

    subgraph "Private Worker Network (Highly Trusted)"
        WORKER_SERVERS[Worker Server Cluster]
        OLLAMA_CLUSTER[Ollama LLM Cluster]
        WORKER_DB[(Worker Database - Read Only)]
    end

    subgraph "Management Network (Admin)"
        MONITORING[Monitoring Services]
        BACKUP[Backup Services]
        LOG_AGGREGATION[Log Aggregation]
    end

    %% Traffic flow with security controls
    INTERNET --> WAF
    WAF -->|HTTPS Only| LB
    LB -->|Internal TLS| API_SERVERS

    %% Secure admin access
    INTERNET -.->|SSH/VPN| BASTION
    BASTION -.->|Admin Access| MONITORING

    %% Inter-service communication
    API_SERVERS -->|Read/Write| API_DB
    API_SERVERS -->|Session Cache| API_CACHE
    API_SERVERS -.->|mTLS gRPC| WORKER_SERVERS

    WORKER_SERVERS -->|Read Only| WORKER_DB
    WORKER_SERVERS -->|Local API| OLLAMA_CLUSTER

    %% Monitoring and logging
    API_SERVERS --> MONITORING
    WORKER_SERVERS --> MONITORING
    API_SERVERS --> LOG_AGGREGATION
    WORKER_SERVERS --> LOG_AGGREGATION

    %% Backup processes
    API_DB --> BACKUP
    WORKER_DB --> BACKUP
```

### gRPC Service Contracts

#### Protocol Buffer Definitions

```protobuf
syntax = "proto3";

package englog.v1;

option go_package = "github.com/garnizeh/englog/proto/v1";

// Core data types
message User {
    string id = 1;
    string email = 2;
    string first_name = 3;
    string last_name = 4;
    string timezone = 5;
}

message LogEntry {
    string id = 1;
    string user_id = 2;
    string project_id = 3;
    string title = 4;
    string description = 5;
    string type = 6; // ActivityType enum
    int64 start_time = 7; // Unix timestamp
    int64 end_time = 8;   // Unix timestamp
    int32 duration_minutes = 9;
    string value_rating = 10; // ValueRating enum
    string impact_level = 11; // ImpactLevel enum
}

// Task processing requests
message GenerateInsightRequest {
    string user_id = 1;
    string report_type = 2; // daily, weekly, monthly, etc.
    int64 period_start = 3; // Unix timestamp
    int64 period_end = 4;   // Unix timestamp
    repeated string focus_areas = 5;
}

message GenerateInsightResponse {
    string task_id = 1;
    string status = 2; // pending, processing, completed, failed
    string insight_content = 3;
    string error_message = 4;
}

message SendEmailRequest {
    string recipient_email = 1;
    string subject = 2;
    string html_content = 3;
    string text_content = 4;
    string email_type = 5; // welcome, report, notification
}

message SendEmailResponse {
    string email_id = 1;
    string status = 2;
    string error_message = 3;
}

message ProcessAnalyticsRequest {
    string user_id = 1;
    int64 period_start = 2;
    int64 period_end = 3;
    repeated string metrics = 4; // specific metrics to calculate
}

message ProcessAnalyticsResponse {
    string user_id = 1;
    map<string, double> metrics = 2; // metric_name -> value
    string status = 3;
}

message ScheduledTaskRequest {
    string task_type = 1;
    string user_id = 2;
    map<string, string> parameters = 3;
    int64 scheduled_time = 4;
}

message ScheduledTaskResponse {
    string task_id = 1;
    string status = 2;
}

// Health check
message HealthCheckRequest {}
message HealthCheckResponse {
    string status = 1;
    map<string, string> service_status = 2;
}

// Main service definition
service WorkerService {
    // Insight generation
    rpc GenerateInsight(GenerateInsightRequest) returns (GenerateInsightResponse);

    // Email notifications
    rpc SendEmail(SendEmailRequest) returns (SendEmailResponse);

    // Analytics processing
    rpc ProcessAnalytics(ProcessAnalyticsRequest) returns (ProcessAnalyticsResponse);

    // Background task scheduling
    rpc ScheduleTask(ScheduledTaskRequest) returns (ScheduledTaskResponse);

    // Health monitoring
    rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);

    // Streaming responses for long-running tasks
    rpc StreamTaskStatus(stream string) returns (stream GenerateInsightResponse);
}
```

#### API Server gRPC Implementation

```go
package grpc

import (
    "context"
    "log"
    "net"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/keepalive"

    pb "github.com/garnizeh/englog/proto/v1"
    "github.com/garnizeh/englog/internal/services"
)

type APIServerGRPC struct {
    pb.UnimplementedWorkerServiceServer
    workerClient pb.WorkerServiceClient

    // Services
    userService    *services.UserService
    logService     *services.LogService
    insightService *services.InsightService
}

func NewAPIServerGRPC(workerClient pb.WorkerServiceClient, userSvc *services.UserService, logSvc *services.LogService) *APIServerGRPC {
    return &APIServerGRPC{
        workerClient:   workerClient,
        userService:    userSvc,
        logService:     logSvc,
        insightService: services.NewInsightService(),
    }
}

func (s *APIServerGRPC) StartGRPCServer(port string, certFile, keyFile string) error {
    lis, err := net.Listen("tcp", ":"+port)
    if err != nil {
        return err
    }

    // TLS credentials
    creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
    if err != nil {
        return err
    }

    // Server options
    opts := []grpc.ServerOption{
        grpc.Creds(creds),
        grpc.KeepaliveParams(keepalive.ServerParameters{
            MaxConnectionIdle:     15 * time.Second,
            MaxConnectionAge:      30 * time.Second,
            MaxConnectionAgeGrace: 5 * time.Second,
            Time:                  5 * time.Second,
            Timeout:               1 * time.Second,
        }),
        grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
            MinTime:             5 * time.Second,
            PermitWithoutStream: true,
        }),
    }

    server := grpc.NewServer(opts...)
    pb.RegisterWorkerServiceServer(server, s)

    log.Printf("gRPC server listening on port %s", port)
    return server.Serve(lis)
}

// Handle insight generation requests from REST API
func (s *APIServerGRPC) RequestInsightGeneration(ctx context.Context, userID string, reportType string, periodStart, periodEnd time.Time) (string, error) {
    req := &pb.GenerateInsightRequest{
        UserId:      userID,
        ReportType:  reportType,
        PeriodStart: periodStart.Unix(),
        PeriodEnd:   periodEnd.Unix(),
    }

    resp, err := s.workerClient.GenerateInsight(ctx, req)
    if err != nil {
        return "", err
    }

    return resp.TaskId, nil
}

// Handle email sending requests
func (s *APIServerGRPC) RequestEmailSend(ctx context.Context, recipientEmail, subject, htmlContent, textContent, emailType string) (string, error) {
    req := &pb.SendEmailRequest{
        RecipientEmail: recipientEmail,
        Subject:        subject,
        HtmlContent:    htmlContent,
        TextContent:    textContent,
        EmailType:      emailType,
    }

    resp, err := s.workerClient.SendEmail(ctx, req)
    if err != nil {
        return "", err
    }

    return resp.EmailId, nil
}
```

#### Worker Server gRPC Implementation

```go
package worker

import (
    "context"
    "crypto/tls"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"

    pb "github.com/garnizeh/englog/proto/v1"
    "github.com/garnizeh/englog/worker/internal/services"
)

type WorkerServerGRPC struct {
    pb.UnimplementedWorkerServiceServer

    // Services
    ollamaService    *services.OllamaService
    emailService     *services.EmailService
    analyticsService *services.AnalyticsService
    taskService      *services.TaskService
}

func NewWorkerServerGRPC() *WorkerServerGRPC {
    return &WorkerServerGRPC{
        ollamaService:    services.NewOllamaService(),
        emailService:     services.NewEmailService(),
        analyticsService: services.NewAnalyticsService(),
        taskService:      services.NewTaskService(),
    }
}

func (w *WorkerServerGRPC) ConnectToAPIServer(apiServerAddr string, certFile string) (pb.WorkerServiceClient, error) {
    // TLS configuration
    config := &tls.Config{
        ServerName: "api-server", // Must match certificate
    }

    creds := credentials.NewTLS(config)

    // Connection options
    opts := []grpc.DialOption{
        grpc.WithTransportCredentials(creds),
        grpc.WithKeepaliveParams(keepalive.ClientParameters{
            Time:                10 * time.Second,
            Timeout:             time.Second,
            PermitWithoutStream: true,
        }),
    }

    conn, err := grpc.NewClient(apiServerAddr, opts...)
    if err != nil {
        return nil, err
    }

    return pb.NewWorkerServiceClient(conn), nil
}

func (w *WorkerServerGRPC) GenerateInsight(ctx context.Context, req *pb.GenerateInsightRequest) (*pb.GenerateInsightResponse, error) {
    log.Printf("Generating insight for user %s, type %s", req.UserId, req.ReportType)

    // Create task
    taskID := w.taskService.CreateTask("insight_generation", req.UserId)

    // Process in background
    go func() {
        periodStart := time.Unix(req.PeriodStart, 0)
        periodEnd := time.Unix(req.PeriodEnd, 0)

        insight, err := w.ollamaService.GenerateInsight(context.Background(), req.UserId, req.ReportType, periodStart, periodEnd)
        if err != nil {
            w.taskService.MarkTaskFailed(taskID, err.Error())
            return
        }

        w.taskService.MarkTaskCompleted(taskID, insight)
    }()

    return &pb.GenerateInsightResponse{
        TaskId: taskID,
        Status: "pending",
    }, nil
}

func (w *WorkerServerGRPC) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
    log.Printf("Sending email to %s, type %s", req.RecipientEmail, req.EmailType)

    emailID, err := w.emailService.SendEmail(ctx, req.RecipientEmail, req.Subject, req.HtmlContent, req.TextContent)
    if err != nil {
        return &pb.SendEmailResponse{
            Status:       "failed",
            ErrorMessage: err.Error(),
        }, nil
    }

    return &pb.SendEmailResponse{
        EmailId: emailID,
        Status:  "sent",
    }, nil
}

func (w *WorkerServerGRPC) ProcessAnalytics(ctx context.Context, req *pb.ProcessAnalyticsRequest) (*pb.ProcessAnalyticsResponse, error) {
    log.Printf("Processing analytics for user %s", req.UserId)

    periodStart := time.Unix(req.PeriodStart, 0)
    periodEnd := time.Unix(req.PeriodEnd, 0)

    metrics, err := w.analyticsService.CalculateMetrics(ctx, req.UserId, periodStart, periodEnd, req.Metrics)
    if err != nil {
        return &pb.ProcessAnalyticsResponse{
            UserId: req.UserId,
            Status: "failed",
        }, nil
    }

    return &pb.ProcessAnalyticsResponse{
        UserId:  req.UserId,
        Status:  "completed",
        Metrics: metrics,
    }, nil
}

func (w *WorkerServerGRPC) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
    serviceStatus := make(map[string]string)

    // Check Ollama
    if err := w.ollamaService.HealthCheck(); err != nil {
        serviceStatus["ollama"] = "unhealthy: " + err.Error()
    } else {
        serviceStatus["ollama"] = "healthy"
    }

    // Check email service
    if err := w.emailService.HealthCheck(); err != nil {
        serviceStatus["email"] = "unhealthy: " + err.Error()
    } else {
        serviceStatus["email"] = "healthy"
    }

    status := "healthy"
    for _, svcStatus := range serviceStatus {
        if svcStatus != "healthy" {
            status = "degraded"
            break
        }
    }

    return &pb.HealthCheckResponse{
        Status:        status,
        ServiceStatus: serviceStatus,
    }, nil
}
```

### LLM Integration Architecture

```go
package services

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/garnizeh/englog/models"
)

type OllamaClient struct {
    BaseURL string
    Client  *http.Client
}

type OllamaRequest struct {
    Model   string                 `json:"model"`
    Prompt  string                 `json:"prompt"`
    Stream  bool                   `json:"stream"`
    Options map[string]any `json:"options,omitempty"`
}

type OllamaResponse struct {
    Response string `json:"response"`
    Done     bool   `json:"done"`
}

func NewOllamaClient(baseURL string) *OllamaClient {
    return &OllamaClient{
        BaseURL: baseURL,
        Client: &http.Client{
            Timeout: 60 * time.Second,
        },
    }
}

func (c *OllamaClient) GenerateInsight(ctx context.Context, logs []models.LogEntryResponse, reportType string, context map[string]any) (string, error) {
    prompt := c.buildPrompt(logs, reportType, context)

    reqBody := OllamaRequest{
        Model:  "llama3",
        Prompt: prompt,
        Stream: false,
        Options: map[string]any{
            "temperature": 0.3,
            "top_p":       0.9,
            "max_tokens":  2000,
        },
    }

    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return "", err
    }

    req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/api/generate", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.Client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var ollamaResp OllamaResponse
    if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
        return "", err
    }

    return ollamaResp.Response, nil
}

func (c *OllamaClient) buildPrompt(logs []models.LogEntryResponse, reportType string, context map[string]any) string {
    // Calculate summary statistics
    totalHours := c.calculateTotalHours(logs)
    activityDistribution := c.calculateActivityDistribution(logs)
    projectBreakdown := c.calculateProjectBreakdown(logs)

    var builder strings.Builder

    builder.WriteString(fmt.Sprintf(`
As a senior software engineering manager, analyze the following work log and generate a professional %s summary.

## Work Log Summary
- Total Hours Logged: %.1f hours
- Number of Activities: %d
- Date Range: %s to %s

## Activity Distribution
%s

## Project Breakdown
%s

## Value & Impact Analysis
%s

## Detailed Activities
%s

## Analysis Requirements
Generate a summary that includes:
1. Key accomplishments and deliverables
2. Time allocation analysis and productivity insights
3. Technical contributions and skill development
4. Collaboration and mentorship activities
5. Value assessment: highlight high-value activities and their outcomes
6. Impact analysis: contributions at personal, team, product, and company levels
7. Areas of focus and potential improvements
8. Notable patterns or observations

Format the response in clear sections with bullet points for easy reading.
Focus on professional achievements and quantifiable impacts.
Pay special attention to high-value, high-impact activities.
`, reportType, totalHours, len(logs),
        logs[0].StartTime.Format("2006-01-02"),
        logs[len(logs)-1].StartTime.Format("2006-01-02"),
        c.formatDistribution(activityDistribution),
        c.formatProjectBreakdown(projectBreakdown),
        c.formatValueImpactAnalysis(logs),
        c.formatDetailedActivities(logs)))

    return builder.String()
}

func (c *OllamaClient) calculateTotalHours(logs []models.LogEntryResponse) float64 {
    total := 0
    for _, log := range logs {
        total += log.DurationMinutes
    }
    return float64(total) / 60.0
}

func (c *OllamaClient) calculateActivityDistribution(logs []models.LogEntryResponse) map[models.ActivityType]float64 {
    totalMinutes := 0
    for _, log := range logs {
        totalMinutes += log.DurationMinutes
    }

    distribution := make(map[models.ActivityType]float64)
    activityMinutes := make(map[models.ActivityType]int)

    for _, log := range logs {
        activityMinutes[log.Type] += log.DurationMinutes
    }

    for activityType, minutes := range activityMinutes {
        if minutes > 0 {
            distribution[activityType] = (float64(minutes) / float64(totalMinutes)) * 100
        }
    }

    return distribution
}

func (c *OllamaClient) calculateProjectBreakdown(logs []models.LogEntryResponse) map[string]float64 {
    totalMinutes := 0
    for _, log := range logs {
        totalMinutes += log.DurationMinutes
    }

    projects := make(map[string]int)
    for _, log := range logs {
        project := "No Project"
        if log.ProjectName != nil {
            project = *log.ProjectName
        }
        projects[project] += log.DurationMinutes
    }

    breakdown := make(map[string]float64)
    for project, minutes := range projects {
        breakdown[project] = (float64(minutes) / float64(totalMinutes)) * 100
    }

    return breakdown
}

func (c *OllamaClient) formatValueImpactAnalysis(logs []models.LogEntryResponse) string {
    valueDistribution := make(map[models.ValueRating]int)
    impactDistribution := make(map[models.ImpactLevel]int)
    totalMinutes := 0

    for _, log := range logs {
        valueDistribution[log.ValueRating] += log.DurationMinutes
        impactDistribution[log.ImpactLevel] += log.DurationMinutes
        totalMinutes += log.DurationMinutes
    }

    var result strings.Builder

    result.WriteString("Value Distribution:\n")
    for value, minutes := range valueDistribution {
        percentage := (float64(minutes) / float64(totalMinutes)) * 100
        result.WriteString(fmt.Sprintf("- %s: %.1f%% (%.1f hours)\n",
            strings.Title(string(value)), percentage, float64(minutes)/60.0))
    }

    result.WriteString("\nImpact Distribution:\n")
    for impact, minutes := range impactDistribution {
        percentage := (float64(minutes) / float64(totalMinutes)) * 100
        result.WriteString(fmt.Sprintf("- %s: %.1f%% (%.1f hours)\n",
            strings.Title(string(impact)), percentage, float64(minutes)/60.0))
    }

    return result.String()
}

type InsightService struct {
    ollama *OllamaClient
    db     *pgxpool.Pool
    queries *sqlc.Queries
}

func NewInsightService(ollama *OllamaClient, db *pgxpool.Pool) *InsightService {
    return &InsightService{
        ollama:  ollama,
        db:      db,
        queries: sqlc.New(db),
    }
}

func (s *InsightService) GenerateWeeklySummary(ctx context.Context, userID string, weekStart time.Time) (string, error) {
    weekEnd := weekStart.AddDate(0, 0, 7)

    logs, err := s.fetchLogs(ctx, userID, weekStart, weekEnd)
    if err != nil {
        return "", err
    }

    if len(logs) == 0 {
        return "No activities logged for this week.", nil
    }

    // Generate context
    context := map[string]any{
        "week_number": weekStart.ISOWeek(),
        "year":        weekStart.Year(),
    }

    // Add previous week comparison if needed
    prevWeekComparison, _ := s.getPreviousWeekComparison(ctx, userID, weekStart)
    if prevWeekComparison != nil {
        context["previous_week_comparison"] = prevWeekComparison
    }

    return s.ollama.GenerateInsight(ctx, logs, "weekly", context)
}

func (s *InsightService) GeneratePerformanceReview(ctx context.Context, userID string, periodStart, periodEnd time.Time) (string, error) {
    logs, err := s.fetchLogs(ctx, userID, periodStart, periodEnd)
    if err != nil {
        return "", err
    }

    // Enhance context for performance review
    context := map[string]any{
        "period_type": s.determinePeriodType(periodStart, periodEnd),
    }

    // Add additional analysis
    goalAlignment := s.analyzeGoalAlignment(logs)
    skillDevelopment := s.analyzeSkillDevelopment(logs)
    leadershipActivities := s.identifyLeadershipActivities(logs)

    context["goal_alignment"] = goalAlignment
    context["skill_development"] = skillDevelopment
    context["leadership_activities"] = leadershipActivities

    return s.ollama.GenerateInsight(ctx, logs, "performance_review", context)
}

func (s *InsightService) fetchLogs(ctx context.Context, userID string, start, end time.Time) ([]models.LogEntryResponse, error) {
    logs, err := s.queries.GetLogEntriesByUserAndDateRange(ctx, sqlc.GetLogEntriesByUserAndDateRangeParams{
        UserID:    userID,
        StartTime: start,
        EndTime:   end,
    })
    if err != nil {
        return nil, err
    }

    // Convert to response models
    var responses []models.LogEntryResponse
    for _, log := range logs {
        responses = append(responses, s.convertToResponse(log))
    }

    return responses, nil
}
```

### Advanced Analytics Features

```go
package analytics

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/garnizeh/englog/internal/sqlc"
    "github.com/garnizeh/englog/models"
)

type ProductivityMetrics struct {
    AverageDailyHours     float64 `json:"average_daily_hours"`
    PeakProductivityHour  int     `json:"peak_productivity_hour"`
    MostProductiveDay     string  `json:"most_productive_day"`
    ActivityEfficiencyScore float64 `json:"activity_efficiency_score"`
    FocusTimePercentage   float64 `json:"focus_time_percentage"`
    CollaborationRatio    float64 `json:"collaboration_ratio"`
    HighValuePercentage   float64 `json:"high_value_percentage"`
    CompanyImpactRatio    float64 `json:"company_impact_ratio"`
    ValueImpactScore      float64 `json:"value_impact_score"`
}

type AnalyticsEngine struct {
    db      *pgxpool.Pool
    queries *sqlc.Queries
}

func NewAnalyticsEngine(db *pgxpool.Pool) *AnalyticsEngine {
    return &AnalyticsEngine{
        db:      db,
        queries: sqlc.New(db),
    }
}

func (a *AnalyticsEngine) CalculateProductivityMetrics(ctx context.Context, userID string, periodStart, periodEnd time.Time) (*ProductivityMetrics, error) {
    logs, err := a.fetchLogsForAnalysis(ctx, userID, periodStart, periodEnd)
    if err != nil {
        return nil, err
    }

    if len(logs) == 0 {
        return &ProductivityMetrics{}, nil
    }

    return &ProductivityMetrics{
        AverageDailyHours:     a.calculateAverageDailyHours(logs),
        PeakProductivityHour:  a.findPeakProductivityHour(logs),
        MostProductiveDay:     a.findMostProductiveDay(logs),
        ActivityEfficiencyScore: a.calculateEfficiencyScore(logs),
        FocusTimePercentage:   a.calculateFocusTimePercentage(logs),
        CollaborationRatio:    a.calculateCollaborationRatio(logs),
        HighValuePercentage:   a.calculateHighValuePercentage(logs),
        CompanyImpactRatio:    a.calculateCompanyImpactRatio(logs),
        ValueImpactScore:      a.calculateValueImpactScore(logs),
    }, nil
}

func (a *AnalyticsEngine) calculateAverageDailyHours(logs []models.LogEntry) float64 {
    if len(logs) == 0 {
        return 0
    }

    // Group by date and calculate daily totals
    dailyMinutes := make(map[string]int)
    for _, log := range logs {
        date := log.StartTime.Format("2006-01-02")
        dailyMinutes[date] += log.DurationMinutes
    }

    totalHours := 0.0
    for _, minutes := range dailyMinutes {
        totalHours += float64(minutes) / 60.0
    }

    return totalHours / float64(len(dailyMinutes))
}

func (a *AnalyticsEngine) findPeakProductivityHour(logs []models.LogEntry) int {
    hourlyActivity := make(map[int]int)

    for _, log := range logs {
        hour := log.StartTime.Hour()
        hourlyActivity[hour] += log.DurationMinutes
    }

    maxMinutes := 0
    peakHour := 0
    for hour, minutes := range hourlyActivity {
        if minutes > maxMinutes {
            maxMinutes = minutes
            peakHour = hour
        }
    }

    return peakHour
}

func (a *AnalyticsEngine) calculateEfficiencyScore(logs []models.LogEntry) float64 {
    focusedWork := 0
    totalTime := 0

    focusedTypes := map[models.ActivityType]bool{
        models.ActivityCode:  true,
        models.ActivityDebug: true,
        models.ActivityDocs:  true,
    }

    for _, log := range logs {
        totalTime += log.DurationMinutes
        if focusedTypes[log.Type] {
            focusedWork += log.DurationMinutes
        }
    }

    if totalTime == 0 {
        return 0
    }

    return (float64(focusedWork) / float64(totalTime)) * 100
}

func (a *AnalyticsEngine) calculateHighValuePercentage(logs []models.LogEntry) float64 {
    totalTime := 0
    highValueTime := 0

    for _, log := range logs {
        totalTime += log.DurationMinutes
        if log.ValueRating == models.ValueHigh {
            highValueTime += log.DurationMinutes
        }
    }

    if totalTime == 0 {
        return 0
    }

    return (float64(highValueTime) / float64(totalTime)) * 100
}

func (a *AnalyticsEngine) calculateCompanyImpactRatio(logs []models.LogEntry) float64 {
    totalTime := 0
    companyImpactTime := 0

    for _, log := range logs {
        totalTime += log.DurationMinutes
        if log.ImpactLevel == models.ImpactCompany || log.ImpactLevel == models.ImpactProduct {
            companyImpactTime += log.DurationMinutes
        }
    }

    if totalTime == 0 {
        return 0
    }

    return (float64(companyImpactTime) / float64(totalTime)) * 100
}

func (a *AnalyticsEngine) calculateValueImpactScore(logs []models.LogEntry) float64 {
    if len(logs) == 0 {
        return 0
    }

    totalScore := 0.0
    totalActivities := len(logs)

    // Weight matrix for value and impact combination
    valueWeights := map[models.ValueRating]float64{
        models.ValueLow:    1.0,
        models.ValueMedium: 2.0,
        models.ValueHigh:   3.0,
    }

    impactWeights := map[models.ImpactLevel]float64{
        models.ImpactPersonal: 1.0,
        models.ImpactTeam:     2.0,
        models.ImpactProduct:  3.0,
        models.ImpactCompany:  4.0,
    }

    for _, log := range logs {
        valueWeight := valueWeights[log.ValueRating]
        impactWeight := impactWeights[log.ImpactLevel]

        // Calculate weighted score (max possible is 3 * 4 = 12)
        activityScore := (valueWeight * impactWeight) / 12.0 * 100
        totalScore += activityScore
    }

    return totalScore / float64(totalActivities)
}

func (a *AnalyticsEngine) GenerateTrendAnalysis(ctx context.Context, userID string, weeks int) (map[string][]float64, error) {
    trends := map[string][]float64{
        "weekly_hours":     make([]float64, 0, weeks),
        "code_percentage":  make([]float64, 0, weeks),
        "meeting_percentage": make([]float64, 0, weeks),
        "efficiency_score": make([]float64, 0, weeks),
    }

    for weekOffset := 0; weekOffset < weeks; weekOffset++ {
        weekStart := time.Now().AddDate(0, 0, -weekOffset*7)
        weekEnd := weekStart.AddDate(0, 0, 7)

        metrics, err := a.CalculateProductivityMetrics(ctx, userID, weekStart, weekEnd)
        if err != nil {
            continue
        }

        trends["weekly_hours"] = append(trends["weekly_hours"], metrics.AverageDailyHours*7)
        trends["efficiency_score"] = append(trends["efficiency_score"], metrics.ActivityEfficiencyScore)

        // Calculate activity type percentages
        logs, _ := a.fetchLogsForAnalysis(ctx, userID, weekStart, weekEnd)
        codePercentage, meetingPercentage := a.calculateActivityPercentages(logs)
        trends["code_percentage"] = append(trends["code_percentage"], codePercentage)
        trends["meeting_percentage"] = append(trends["meeting_percentage"], meetingPercentage)
    }

    return trends, nil
}

func (a *AnalyticsEngine) calculateActivityPercentages(logs []models.LogEntry) (float64, float64) {
    totalMinutes := 0
    codeMinutes := 0
    meetingMinutes := 0

    for _, log := range logs {
        totalMinutes += log.DurationMinutes
        switch log.Type {
        case models.ActivityCode:
            codeMinutes += log.DurationMinutes
        case models.ActivityMeeting:
            meetingMinutes += log.DurationMinutes
        }
    }

    if totalMinutes == 0 {
        return 0, 0
    }

    codePercentage := (float64(codeMinutes) / float64(totalMinutes)) * 100
    meetingPercentage := (float64(meetingMinutes) / float64(totalMinutes)) * 100

    return codePercentage, meetingPercentage
}

func (a *AnalyticsEngine) fetchLogsForAnalysis(ctx context.Context, userID string, start, end time.Time) ([]models.LogEntry, error) {
    logs, err := a.queries.GetLogEntriesByUserAndDateRange(ctx, sqlc.GetLogEntriesByUserAndDateRangeParams{
        UserID:    userID,
        StartTime: start,
        EndTime:   end,
    })
    if err != nil {
        return nil, err
    }

    // Convert sqlc generated types to models.LogEntry
    var modelLogs []models.LogEntry
    for _, log := range logs {
        modelLogs = append(modelLogs, models.LogEntry{
            ID:              log.ID,
            UserID:          log.UserID,
            Title:           log.Title,
            Description:     log.Description,
            Type:            models.ActivityType(log.Type),
            StartTime:       log.StartTime,
            EndTime:         log.EndTime,
            DurationMinutes: log.DurationMinutes,
            ValueRating:     models.ValueRating(log.ValueRating),
            ImpactLevel:     models.ImpactLevel(log.ImpactLevel),
            CreatedAt:       log.CreatedAt,
            UpdatedAt:       log.UpdatedAt,
        })
    }

    return modelLogs, nil
}
```

---

## Security & Privacy

### JWT Authentication Architecture

EngLog implementa um sistema robusto de autenticação baseado em **JWT (JSON Web Tokens)** com **refresh tokens** para garantir máxima segurança e experiência do usuário otimizada.

#### Características Principais do Sistema JWT

**1. Dual Token Strategy**
- **Access Token**: Vida curta (15 minutos), usado para requisições da API
- **Refresh Token**: Vida longa (30 dias), usado apenas para renovar access tokens

**2. Security Features**
- **Token Rotation**: Refresh tokens são automaticamente rotacionados a cada uso
- **Blacklist System**: Refresh tokens invalidados são armazenados em blacklist
- **JTI (JWT ID)**: Identificador único para rastreamento de refresh tokens
- **Automatic Cleanup**: Tokens expirados são automaticamente removidos

**3. Attack Prevention**
- **Replay Attack Protection**: Blacklist previne reutilização de tokens
- **Token Hijacking Mitigation**: Rotação automática limita janela de exposição
- **Brute Force Protection**: Tokens de vida curta reduzem superficie de ataque

#### Fluxo de Autenticação Detalhado

**Registro/Login:**
1. Cliente fornece credenciais válidas
2. Sistema gera par de tokens (access + refresh)
3. Access token usado para requisições subsequentes
4. Refresh token armazenado seguramente pelo cliente

**Refresh Token Flow:**
1. Access token expira após 15 minutos
2. Cliente usa refresh token para obter novos tokens
3. Sistema valida refresh token e verifica blacklist
4. Refresh token antigo é invalidado (blacklist)
5. Novos tokens são gerados e retornados
6. Processo repete automaticamente

**Logout/Security:**
1. Cliente envia refresh token para logout
2. Sistema adiciona token à blacklist
3. Todas as sessões do usuário são invalidadas
4. Cliente remove tokens do armazenamento local

#### Configurações de Segurança

```yaml
# Configurações recomendadas para produção
JWT_ACCESS_TOKEN_TTL: "15m"     # 15 minutos
JWT_REFRESH_TOKEN_TTL: "720h"   # 30 dias
JWT_SECRET_KEY: "complex-secret-key-minimum-32-chars"
JWT_BLACKLIST_CLEANUP_INTERVAL: "24h"  # Limpeza diária
```

#### Endpoints de Autenticação

| Endpoint | Método | Propósito | Tokens Necessários |
|----------|---------|-----------|-------------------|
| `/v1/auth/register` | POST | Registro de usuário | Nenhum |
| `/v1/auth/login` | POST | Login do usuário | Nenhum |
| `/v1/auth/refresh` | POST | Renovar tokens | Refresh Token |
| `/v1/auth/logout` | POST | Logout do usuário | Refresh Token |
| `/v1/auth/me` | GET | Dados do usuário | Access Token |

### Authentication & Authorization

#### JWT Token Strategy
```go
package auth

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/garnizeh/englog/internal/sqlc"
)

// Authentication errors
var (
    ErrInvalidTokenType     = errors.New("invalid token type")
    ErrInvalidSigningMethod = errors.New("invalid signing method")
    ErrInvalidToken         = errors.New("invalid token")
    ErrTokenBlacklisted     = errors.New("token is blacklisted")
)

type AuthService struct {
    db               *pgxpool.Pool
    queries          *sqlc.Queries
    secretKey        []byte
    accessTokenTTL   time.Duration
    refreshTokenTTL  time.Duration
}

type Claims struct {
    UserID    string `json:"sub"`
    TokenType string `json:"type"`
    JTI       string `json:"jti,omitempty"` // JWT ID for refresh tokens
    jwt.RegisteredClaims
}

func NewAuthService(db *pgxpool.Pool, secretKey string) *AuthService {
    return &AuthService{
        db:               db,
        queries:          sqlc.New(db),
        secretKey:        []byte(secretKey),
        accessTokenTTL:   15 * time.Minute,
        refreshTokenTTL:  30 * 24 * time.Hour,
    }
}

func (a *AuthService) CreateAccessToken(userID string) (string, error) {
    claims := &Claims{
        UserID:    userID,
        TokenType: "access",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.accessTokenTTL)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(a.secretKey)
}

func (a *AuthService) CreateRefreshToken(userID string) (string, error) {
    jti, err := generateJTI()
    if err != nil {
        return "", err
    }

    claims := &Claims{
        UserID:    userID,
        TokenType: "refresh",
        JTI:       jti,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.refreshTokenTTL)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(a.secretKey)
}

func (a *AuthService) RotateRefreshToken(oldRefreshToken string) (string, string, error) {
    claims, err := a.ValidateToken(oldRefreshToken)
    if err != nil {
        return "", "", err
    }

    if claims.TokenType != "refresh" {
        return "", "", ErrInvalidTokenType
    }

    // Invalidate old refresh token
    if err := a.invalidateRefreshToken(claims.JTI); err != nil {
        return "", "", err
    }

    // Create new tokens
    newAccessToken, err := a.CreateAccessToken(claims.UserID)
    if err != nil {
        return "", "", err
    }

    newRefreshToken, err := a.CreateRefreshToken(claims.UserID)
    if err != nil {
        return "", "", err
    }

    return newAccessToken, newRefreshToken, nil
}

func (a *AuthService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidSigningMethod
        }
        return a.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, ErrInvalidToken
    }

    // Check if refresh token is denylisted
    if claims.TokenType == "refresh" {
        if denylisted, err := a.isRefreshTokenDenylisted(claims.JTI); err != nil || denylisted {
            return nil, ErrTokenDenylisted
        }
    }

    return claims, nil
}

func (a *AuthService) HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func (a *AuthService) CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func generateJTI() (string, error) {
    bytes := make([]byte, 16)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

func (a *AuthService) invalidateRefreshToken(jti string) error {
    // Add to denylist table using sqlc
    return a.queries.CreateRefreshTokenDenylist(context.Background(), jti)
}

func (a *AuthService) isRefreshTokenDenylisted(jti string) (bool, error) {
    exists, err := a.queries.IsRefreshTokenDenylisted(context.Background(), jti)
    return exists, err
}

// Middleware for protecting routes
func (a *AuthService) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        if !strings.HasPrefix(authHeader, "Bearer ") {
            c.JSON(401, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }

        tokenString := authHeader[7:] // Remove "Bearer " prefix
        claims, err := a.ValidateToken(tokenString)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        if claims.TokenType != "access" {
            c.JSON(401, gin.H{"error": "Access token required"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Next()
    }
}

// HTTP Handlers for Authentication
func (a *AuthService) LoginHandler(c *gin.Context) {
    var loginReq struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&loginReq); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request format", "details": err.Error()})
        return
    }

    // Get user from database
    user, err := a.queries.GetUserByEmail(c.Request.Context(), loginReq.Email)
    if err != nil {
        c.JSON(401, gin.H{"error": "Invalid credentials"})
        return
    }

    // Check password
    if !a.CheckPassword(loginReq.Password, user.PasswordHash) {
        c.JSON(401, gin.H{"error": "Invalid credentials"})
        return
    }

    // Create tokens
    accessToken, err := a.CreateAccessToken(user.ID)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to create access token"})
        return
    }

    refreshToken, err := a.CreateRefreshToken(user.ID)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to create refresh token"})
        return
    }

    // Update last login
    a.queries.UpdateUserLastLogin(c.Request.Context(), user.ID)

    c.JSON(200, gin.H{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
        "expires_in":    int(a.accessTokenTTL.Seconds()),
        "token_type":    "Bearer",
        "user": gin.H{
            "id":         user.ID,
            "email":      user.Email,
            "first_name": user.FirstName,
            "last_name":  user.LastName,
        },
    })
}

func (a *AuthService) RefreshHandler(c *gin.Context) {
    var refreshReq struct {
        RefreshToken string `json:"refresh_token" binding:"required"`
    }

    if err := c.ShouldBindJSON(&refreshReq); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request format", "details": err.Error()})
        return
    }

    // Rotate refresh token
    newAccessToken, newRefreshToken, err := a.RotateRefreshToken(refreshReq.RefreshToken)
    if err != nil {
        c.JSON(401, gin.H{"error": "Invalid or expired refresh token"})
        return
    }

    c.JSON(200, gin.H{
        "access_token":  newAccessToken,
        "refresh_token": newRefreshToken,
        "expires_in":    int(a.accessTokenTTL.Seconds()),
        "token_type":    "Bearer",
    })
}

func (a *AuthService) LogoutHandler(c *gin.Context) {
    var logoutReq struct {
        RefreshToken string `json:"refresh_token" binding:"required"`
    }

    if err := c.ShouldBindJSON(&logoutReq); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request format"})
        return
    }

    // Validate and blacklist refresh token
    claims, err := a.ValidateToken(logoutReq.RefreshToken)
    if err != nil {
        c.JSON(401, gin.H{"error": "Invalid refresh token"})
        return
    }

    if claims.TokenType != "refresh" {
        c.JSON(400, gin.H{"error": "Expected refresh token"})
        return
    }

    // Blacklist the refresh token
    if err := a.invalidateRefreshToken(claims.JTI); err != nil {
        c.JSON(500, gin.H{"error": "Failed to logout"})
        return
    }

    c.JSON(200, gin.H{"message": "Successfully logged out"})
}

func (a *AuthService) MeHandler(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(401, gin.H{"error": "Unauthorized"})
        return
    }

    user, err := a.queries.GetUserByID(c.Request.Context(), userID.(string))
    if err != nil {
        c.JSON(404, gin.H{"error": "User not found"})
        return
    }

    c.JSON(200, gin.H{
        "id":         user.ID,
        "email":      user.Email,
        "first_name": user.FirstName,
        "last_name":  user.LastName,
        "timezone":   user.Timezone,
        "created_at": user.CreatedAt,
        "updated_at": user.UpdatedAt,
    })
}

func (a *AuthService) RegisterHandler(c *gin.Context) {
    var registerReq struct {
        Email     string `json:"email" binding:"required,email"`
        Password  string `json:"password" binding:"required,min=8"`
        FirstName string `json:"first_name" binding:"required"`
        LastName  string `json:"last_name" binding:"required"`
        Timezone  string `json:"timezone"`
    }

    if err := c.ShouldBindJSON(&registerReq); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request format", "details": err.Error()})
        return
    }

    // Check if user already exists
    _, err := a.queries.GetUserByEmail(c.Request.Context(), registerReq.Email)
    if err == nil {
        c.JSON(409, gin.H{"error": "User already exists"})
        return
    }

    // Hash password
    hashedPassword, err := a.HashPassword(registerReq.Password)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to process password"})
        return
    }

    // Create user
    userID := generateUserID()
    timezone := registerReq.Timezone
    if timezone == "" {
        timezone = "UTC"
    }

    _, err = a.queries.CreateUser(c.Request.Context(), sqlc.CreateUserParams{
        ID:           userID,
        Email:        registerReq.Email,
        PasswordHash: hashedPassword,
        FirstName:    registerReq.FirstName,
        LastName:     registerReq.LastName,
        Timezone:     timezone,
    })
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to create user"})
        return
    }

    // Create tokens for immediate login
    accessToken, err := a.CreateAccessToken(userID)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to create access token"})
        return
    }

    refreshToken, err := a.CreateRefreshToken(userID)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to create refresh token"})
        return
    }

    c.JSON(201, gin.H{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
        "expires_in":    int(a.accessTokenTTL.Seconds()),
        "token_type":    "Bearer",
        "user": gin.H{
            "id":         userID,
            "email":      registerReq.Email,
            "first_name": registerReq.FirstName,
            "last_name":  registerReq.LastName,
        },
    })
}

func generateUserID() string {
    bytes := make([]byte, 16)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}
```

### Data Privacy Measures

#### Data Encryption
```go
package security

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "io"
)

type DataEncryption struct {
    gcm cipher.AEAD
}

func NewDataEncryption(key []byte) (*DataEncryption, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    return &DataEncryption{gcm: gcm}, nil
}

func (d *DataEncryption) EncryptSensitiveData(data string) (string, error) {
    if data == "" {
        return "", nil
    }

    nonce := make([]byte, d.gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := d.gcm.Seal(nonce, nonce, []byte(data), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (d *DataEncryption) DecryptSensitiveData(encryptedData string) (string, error) {
    if encryptedData == "" {
        return "", nil
    }

    data, err := base64.StdEncoding.DecodeString(encryptedData)
    if err != nil {
        return "", err
    }

    nonceSize := d.gcm.NonceSize()
    if len(data) < nonceSize {
        return "", fmt.Errorf("ciphertext too short")
    }

    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    plaintext, err := d.gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}

// Data encryption service for sensitive fields
type EncryptionService struct {
    encryption *DataEncryption
}

func NewEncryptionService(key []byte) (*EncryptionService, error) {
    enc, err := NewDataEncryption(key)
    if err != nil {
        return nil, err
    }
    return &EncryptionService{encryption: enc}, nil
}

func (es *EncryptionService) EncryptDescription(description string) (string, error) {
    return es.encryption.EncryptSensitiveData(description)
}

func (es *EncryptionService) DecryptDescription(encryptedDesc string) (string, error) {
    return es.encryption.DecryptSensitiveData(encryptedDesc)
}
```

#### GDPR Compliance Features
```go
package gdpr

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
    "github.com/garnizeh/englog/internal/sqlc"
    "github.com/garnizeh/englog/models"
)

type GDPRService struct {
    db      *pgxpool.Pool
    queries *sqlc.Queries
}

type UserDataExport struct {
    UserProfile []models.User         `json:"user_profile"`
    LogEntries  []models.LogEntry     `json:"log_entries"`
    Projects    []models.Project      `json:"projects"`
    Insights    []models.Insight      `json:"insights"`
    Preferences map[string]any `json:"preferences"`
    ExportedAt  time.Time             `json:"exported_at"`
}

func NewGDPRService(db *pgxpool.Pool) *GDPRService {
    return &GDPRService{
        db:      db,
        queries: sqlc.New(db),
    }
}

func (g *GDPRService) ExportUserData(ctx context.Context, userID string) (*UserDataExport, error) {
    export := &UserDataExport{
        ExportedAt: time.Now(),
    }

    // Get user profile
    user, err := g.queries.GetUserByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    export.UserProfile = []models.User{convertSqlcUserToModel(user)}

    // Get log entries
    logs, err := g.queries.GetLogEntriesByUser(ctx, userID)
    if err != nil {
        return nil, err
    }
    export.LogEntries = convertSqlcLogsToModels(logs)

    // Get projects (where user is creator)
    projects, err := g.queries.GetProjectsByCreator(ctx, userID)
    if err != nil {
        return nil, err
    }
    export.Projects = convertSqlcProjectsToModels(projects)

    // Get insights
    insights, err := g.queries.GetInsightsByUser(ctx, userID)
    if err != nil {
        return nil, err
    }
    export.Insights = convertSqlcInsightsToModels(insights)

    // Parse user preferences
    if user.Preferences != nil {
        export.Preferences = user.Preferences.(map[string]any)
    }

    return export, nil
}

func (g *GDPRService) AnonymizeUserData(ctx context.Context, userID string) error {
    // Start transaction for atomicity
    tx, err := g.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer func() {
        if err != nil {
            tx.Rollback(ctx)
        } else {
            tx.Commit(ctx)
        }
    }()

    // Generate anonymous ID
    anonymousID, err := generateAnonymousID()
    if err != nil {
        return err
    }

    // Update user profile using sqlc
    _, err = g.queries.WithTx(tx).UpdateUserProfile(ctx, sqlc.UpdateUserProfileParams{
        ID:        userID,
        FirstName: "Anonymous",
        LastName:  "User",
        Timezone:  "UTC",
        Preferences: map[string]any{
            "privacy": "anonymized",
        },
    })
    if err != nil {
        return fmt.Errorf("failed to anonymize user profile: %w", err)
    }
                "preferences": nil,
            }).Error
        if err != nil {
            return err
        }

        // Anonymize log entry descriptions
        err = g.anonymizeLogDescriptions(ctx, tx, userID)
        if err != nil {
            return err
        }

        // Remove sensitive insights
        err = tx.Where("user_id = ?", userID).Delete(&models.Insight{}).Error
        if err != nil {
            return err
        }

        return nil
    })
}

func (g *GDPRService) DeleteUserData(ctx context.Context, userID string) error {
    tx, err := g.db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    qtx := g.queries.WithTx(tx)

    // Delete in correct order to respect foreign key constraints
    if err := qtx.DeleteLogEntryTagsByUser(ctx, userID); err != nil {
        return err
    }

    if err := qtx.DeleteLogEntriesByUser(ctx, userID); err != nil {
        return err
    }

    if err := qtx.DeleteInsightsByUser(ctx, userID); err != nil {
        return err
    }

    if err := qtx.DeleteUserTeamsByUser(ctx, userID); err != nil {
        return err
    }

    if err := qtx.DeleteUserSessionsByUser(ctx, userID); err != nil {
        return err
    }

    if err := qtx.DeleteUser(ctx, userID); err != nil {
        return err
    }

    return tx.Commit(ctx)
}

func (g *GDPRService) ScheduleDataDeletion(ctx context.Context, userID string, deleteAfter time.Time) error {
    return g.queries.CreateScheduledDeletion(ctx, sqlc.CreateScheduledDeletionParams{
        UserID:      userID,
        ScheduledAt: deleteAfter,
        Status:      "pending",
    })
}

func (g *GDPRService) ProcessScheduledDeletions(ctx context.Context) error {
    deletions, err := g.queries.GetPendingScheduledDeletions(ctx)
    if err != nil {
        return err
    }

    for _, deletion := range deletions {
        err := g.DeleteUserData(ctx, deletion.UserID)
        if err != nil {
            // Mark as failed and continue
            g.queries.UpdateScheduledDeletionStatus(ctx, sqlc.UpdateScheduledDeletionStatusParams{
                ID:     deletion.ID,
                Status: "failed",
            })
            continue
        }

        // Mark as completed
        g.queries.UpdateScheduledDeletionStatus(ctx, sqlc.UpdateScheduledDeletionStatusParams{
            ID:     deletion.ID,
            Status: "completed",
        })
    }

    return nil
}

func (g *GDPRService) anonymizeLogDescriptions(ctx context.Context, queries *sqlc.Queries, userID string) error {
    // Update all log entries to have empty descriptions
    _, err := queries.AnonymizeLogDescriptions(ctx, sqlc.AnonymizeLogDescriptionsParams{
        UserID:      userID,
        Description: "",
    })
    return err
}

func generateAnonymousID() (string, error) {
    bytes := make([]byte, 4)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return "anon_" + hex.EncodeToString(bytes), nil
}

// Background job for processing scheduled deletions
func (g *GDPRService) StartDeletionWorker(ctx context.Context) {
    ticker := time.NewTicker(24 * time.Hour) // Check daily
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            if err := g.ProcessScheduledDeletions(ctx); err != nil {
                // Log error but continue
                fmt.Printf("Error processing scheduled deletions: %v\n", err)
            }
        }
    }
}
```

---

## Performance & Scalability

### Database Optimization

#### Query Optimization Strategies
```sql
-- Partitioning for large datasets
CREATE TABLE log_entries_2024 PARTITION OF log_entries
FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

-- Optimized indexes for common queries
CREATE INDEX CONCURRENTLY idx_log_entries_user_time_desc
ON log_entries (user_id, start_time DESC)
WHERE start_time >= '2024-01-01';

-- Materialized view for analytics
CREATE MATERIALIZED VIEW user_activity_summary AS
SELECT
    user_id,
    DATE_TRUNC('week', start_time) as week,
    type,
    SUM(duration_minutes) as total_minutes,
    COUNT(*) as activity_count
FROM log_entries
WHERE start_time >= NOW() - INTERVAL '1 year'
GROUP BY user_id, week, type;

-- Refresh strategy
CREATE INDEX ON user_activity_summary (user_id, week);
```

#### Caching Strategy
```python
from redis import Redis
import json
from functools import wraps

class CacheService:
    def __init__(self):
        self.redis = Redis.from_url(os.getenv("REDIS_URL"))
        self.default_ttl = 3600  # 1 hour

    def cache_key(self, prefix: str, **kwargs) -> str:
        """Generate cache key from parameters."""
        key_parts = [prefix] + [f"{k}:{v}" for k, v in sorted(kwargs.items())]
        return ":".join(key_parts)

    async def get_or_set(self, key: str, func, ttl: int = None):
        """Get from cache or execute function and cache result."""
        cached = await self.redis.get(key)
        if cached:
            return json.loads(cached)

        result = await func()
        await self.redis.setex(
            key,
            ttl or self.default_ttl,
            json.dumps(result, default=str)
        )
        return result

def cached_query(ttl: int = 3600):
    """Decorator for caching database queries."""
    def decorator(func):
        @wraps(func)
        async def wrapper(*args, **kwargs):
            cache_key = f"query:{func.__name__}:{hash(str(args) + str(kwargs))}"
            return await cache_service.get_or_set(cache_key, lambda: func(*args, **kwargs), ttl)
        return wrapper
    return decorator

# Usage example
@cached_query(ttl=1800)  # Cache for 30 minutes
async def get_user_activity_summary(user_id: str, period_start: datetime):
    # Expensive query implementation
    pass
```

### Background Processing

#### Task Processing with Go Routines
```go
package tasks

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "sync"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/garnizeh/englog/internal/sqlc"
    "github.com/garnizeh/englog/models"
    "github.com/garnizeh/englog/services"
)

type TaskType string

const (
    TaskGenerateWeeklyReport    TaskType = "generate_weekly_report"
    TaskBatchImportCalendar     TaskType = "batch_import_calendar"
    TaskCleanupOldInsights      TaskType = "cleanup_old_insights"
    TaskUpdateUserMetrics       TaskType = "update_user_metrics"
)

type Task struct {
    ID        string                 `json:"id"`
    Type      TaskType               `json:"type"`
    UserID    string                 `json:"user_id,omitempty"`
    Payload   map[string]any `json:"payload"`
    CreatedAt time.Time              `json:"created_at"`
    Status    string                 `json:"status"` // pending, running, completed, failed
    Result    *string                `json:"result,omitempty"`
    Error     *string                `json:"error,omitempty"`
}

type TaskQueue struct {
    tasks    chan *Task
    workers  int
    db       *pgxpool.Pool
    queries  *sqlc.Queries
    wg       sync.WaitGroup
    ctx      context.Context
    cancel   context.CancelFunc

    // Services
    insightService   *services.InsightService
    analyticsService *services.AnalyticsService
}

func NewTaskQueue(db *pgxpool.Pool, workers int, insightService *services.InsightService) *TaskQueue {
    ctx, cancel := context.WithCancel(context.Background())

    return &TaskQueue{
        tasks:            make(chan *Task, 100), // Buffer for 100 tasks
        workers:          workers,
        db:               db,
        queries:          sqlc.New(db),
        ctx:              ctx,
        cancel:           cancel,
        insightService:   insightService,
        analyticsService: services.NewAnalyticsService(db),
    }
}

func (tq *TaskQueue) Start() {
    for i := 0; i < tq.workers; i++ {
        tq.wg.Add(1)
        go tq.worker(i)
    }

    // Start cleanup routine
    tq.wg.Add(1)
    go tq.cleanupWorker()

    log.Printf("Started task queue with %d workers", tq.workers)
}

func (tq *TaskQueue) Stop() {
    tq.cancel()
    close(tq.tasks)
    tq.wg.Wait()
    log.Println("Task queue stopped")
}

func (tq *TaskQueue) Enqueue(task *Task) error {
    task.Status = "pending"
    task.CreatedAt = time.Now()

    // Save to database using sqlc
    err := tq.queries.CreateTask(context.Background(), sqlc.CreateTaskParams{
        ID:        task.ID,
        Type:      string(task.Type),
        UserID:    task.UserID,
        Payload:   task.Payload,
        Status:    task.Status,
        CreatedAt: task.CreatedAt,
    })
    if err != nil {
        return err
    }

    select {
    case tq.tasks <- task:
        return nil
    case <-tq.ctx.Done():
        return fmt.Errorf("task queue is shutting down")
    default:
        return fmt.Errorf("task queue is full")
    }
}

func (tq *TaskQueue) worker(id int) {
    defer tq.wg.Done()

    log.Printf("Worker %d started", id)

    for {
        select {
        case task := <-tq.tasks:
            if task == nil {
                return // Channel closed
            }
            tq.processTask(task)

        case <-tq.ctx.Done():
            return
        }
    }
}

func (tq *TaskQueue) processTask(task *Task) {
    log.Printf("Processing task %s of type %s", task.ID, task.Type)

    // Update status to running
    tq.updateTaskStatus(task.ID, "running", nil, nil)

    start := time.Now()
    var result any
    var err error

    switch task.Type {
    case TaskGenerateWeeklyReport:
        result, err = tq.generateWeeklyReport(task)
    case TaskBatchImportCalendar:
        result, err = tq.batchImportCalendar(task)
    case TaskCleanupOldInsights:
        result, err = tq.cleanupOldInsights(task)
    case TaskUpdateUserMetrics:
        result, err = tq.updateUserMetrics(task)
    default:
        err = fmt.Errorf("unknown task type: %s", task.Type)
    }

    duration := time.Since(start)
    log.Printf("Task %s completed in %v", task.ID, duration)

    // Update task with result
    if err != nil {
        errorMsg := err.Error()
        tq.updateTaskStatus(task.ID, "failed", nil, &errorMsg)
    } else {
        resultStr := tq.serializeResult(result)
        tq.updateTaskStatus(task.ID, "completed", &resultStr, nil)
    }
}

func (tq *TaskQueue) generateWeeklyReport(task *Task) (any, error) {
    userID, ok := task.Payload["user_id"].(string)
    if !ok {
        return nil, fmt.Errorf("user_id is required")
    }

    weekStartStr, ok := task.Payload["week_start"].(string)
    if !ok {
        return nil, fmt.Errorf("week_start is required")
    }

    weekStart, err := time.Parse(time.RFC3339, weekStartStr)
    if err != nil {
        return nil, fmt.Errorf("invalid week_start format: %v", err)
    }

    report, err := tq.insightService.GenerateWeeklySummary(tq.ctx, userID, weekStart)
    if err != nil {
        return nil, err
    }

    // Store the generated report using sqlc
    err = tq.queries.CreateInsight(tq.ctx, sqlc.CreateInsightParams{
        UserID:      userID,
        ReportType:  "weekly",
        PeriodStart: weekStart,
        PeriodEnd:   weekStart.AddDate(0, 0, 7),
        Content:     report,
    })
    if err != nil {
        return nil, err
    }

    return map[string]any{
        "status": "completed",
        "length": len(report),
    }, nil
}

func (tq *TaskQueue) batchImportCalendar(task *Task) (any, error) {
    userID, ok := task.Payload["user_id"].(string)
    if !ok {
        return nil, fmt.Errorf("user_id is required")
    }

    calendarData, ok := task.Payload["calendar_data"].([]any)
    if !ok {
        return nil, fmt.Errorf("calendar_data is required")
    }

    processed := 0
    errors := []map[string]any{}

    for i, eventData := range calendarData {
        eventMap, ok := eventData.(map[string]any)
        if !ok {
            errors = append(errors, map[string]any{
                "index": i,
                "error": "invalid event format",
            })
            continue
        }

        logEntry, err := tq.convertCalendarEvent(eventMap)
        if err != nil {
            errors = append(errors, map[string]any{
                "index": i,
                "error": err.Error(),
            })
            continue
        }

        logEntry.UserID = userID
        _, err = tq.queries.CreateLogEntry(context.Background(), sqlc.CreateLogEntryParams{
            UserID:          logEntry.UserID,
            Title:           logEntry.Title,
            Description:     logEntry.Description,
            Type:            string(logEntry.Type),
            StartTime:       logEntry.StartTime,
            EndTime:         logEntry.EndTime,
            DurationMinutes: logEntry.DurationMinutes,
            ValueRating:     string(logEntry.ValueRating),
            ImpactLevel:     string(logEntry.ImpactLevel),
        })
        if err != nil {
            errors = append(errors, map[string]any{
                "index": i,
                "error": err.Error(),
            })
            continue
        }

        processed++
    }

    return map[string]any{
        "processed":     processed,
        "errors":        len(errors),
        "error_details": errors,
    }, nil
}

func (tq *TaskQueue) convertCalendarEvent(eventData map[string]any) (*models.LogEntry, error) {
    title, ok := eventData["title"].(string)
    if !ok {
        return nil, fmt.Errorf("missing or invalid title")
    }

    startTimeStr, ok := eventData["start_time"].(string)
    if !ok {
        return nil, fmt.Errorf("missing or invalid start_time")
    }

    endTimeStr, ok := eventData["end_time"].(string)
    if !ok {
        return nil, fmt.Errorf("missing or invalid end_time")
    }

    startTime, err := time.Parse(time.RFC3339, startTimeStr)
    if err != nil {
        return nil, fmt.Errorf("invalid start_time format: %v", err)
    }

    endTime, err := time.Parse(time.RFC3339, endTimeStr)
    if err != nil {
        return nil, fmt.Errorf("invalid end_time format: %v", err)
    }

    duration := int(endTime.Sub(startTime).Minutes())
    if duration <= 0 {
        return nil, fmt.Errorf("invalid duration")
    }

    logEntry := &models.LogEntry{
        Title:           title,
        Type:            models.ActivityMeeting, // Default for calendar events
        StartTime:       startTime,
        EndTime:         endTime,
        DurationMinutes: duration,
        ValueRating:     models.ValueMedium,
        ImpactLevel:     models.ImpactTeam,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }

    // Optional fields
    if description, ok := eventData["description"].(string); ok {
        logEntry.Description = &description
    }

    return logEntry, nil
}

func (tq *TaskQueue) cleanupOldInsights(task *Task) (any, error) {
    cutoffDate := time.Now().AddDate(0, -6, 0) // 6 months ago

    deletedCount, err := tq.queries.DeleteOldInsights(tq.ctx, cutoffDate)
    if err != nil {
        return nil, err
    }

    return map[string]any{
        "deleted_count": deletedCount,
    }, nil
}

func (tq *TaskQueue) updateTaskStatus(taskID, status string, result, errorMsg *string) {
    params := sqlc.UpdateTaskStatusParams{
        ID:     taskID,
        Status: status,
    }

    if result != nil {
        params.Result = result
    }

    if errorMsg != nil {
        params.Error = errorMsg
    }

    tq.queries.UpdateTaskStatus(context.Background(), params)
}

func (tq *TaskQueue) serializeResult(result any) string {
    if result == nil {
        return ""
    }

    jsonData, err := json.Marshal(result)
    if err != nil {
        return fmt.Sprintf("Error serializing result: %v", err)
    }

    return string(jsonData)
}

func (tq *TaskQueue) cleanupWorker() {
    defer tq.wg.Done()

    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-tq.ctx.Done():
            return
        case <-ticker.C:
            // Clean up old completed/failed tasks (older than 30 days)
            cutoff := time.Now().AddDate(0, 0, -30)
            tq.queries.DeleteOldTasks(context.Background(), cutoff)
        }
    }
}

// Scheduled tasks runner
type Scheduler struct {
    taskQueue *TaskQueue
    db        *pgxpool.Pool
    queries   *sqlc.Queries
}

func NewScheduler(taskQueue *TaskQueue, db *pgxpool.Pool) *Scheduler {
    return &Scheduler{
        taskQueue: taskQueue,
        db:        db,
        queries:   sqlc.New(db),
    }
}

func (s *Scheduler) Start(ctx context.Context) {
    // Weekly report generation
    go s.scheduleWeeklyReports(ctx)

    // Daily cleanup
    go s.scheduleDailyCleanup(ctx)
}

func (s *Scheduler) scheduleWeeklyReports(ctx context.Context) {
    // Run every Monday at 9 AM
    for {
        now := time.Now()
        next := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())

        // If we've passed 9 AM today, schedule for next Monday
        if next.Before(now) || next.Weekday() != time.Monday {
            days := (7 - int(now.Weekday()) + int(time.Monday)) % 7
            if days == 0 {
                days = 7
            }
            next = next.AddDate(0, 0, days)
        }

        select {
        case <-ctx.Done():
            return
        case <-time.After(time.Until(next)):
            s.generateWeeklyReportsForAllUsers()
        }
    }
}

func (s *Scheduler) generateWeeklyReportsForAllUsers() {
    users, err := s.queries.GetAllUsers(context.Background())
    if err != nil {
        log.Printf("Error fetching users for weekly reports: %v", err)
        return
    }

    weekStart := time.Now().AddDate(0, 0, -7) // Last week

    for _, user := range users {
        task := &Task{
            ID:   fmt.Sprintf("weekly-report-%s-%s", user.ID, weekStart.Format("2006-01-02")),
            Type: TaskGenerateWeeklyReport,
            UserID: user.ID,
            Payload: map[string]any{
                "user_id":    user.ID,
                "week_start": weekStart.Format(time.RFC3339),
            },
        }

        if err := s.taskQueue.Enqueue(task); err != nil {
            log.Printf("Error enqueuing weekly report for user %s: %v", user.ID, err)
        }
    }
}
```

---

## Deployment & Operations

### Simple Two-Machine Infrastructure Setup

#### Prerequisites

- **Machine 1**: Public IP, 4+ GB RAM, 50+ GB storage
- **Machine 2**: Private IP, 8+ GB RAM (for LLM), 100+ GB storage
- **Docker & Docker Compose** installed on both machines
- **Network**: Private connection between machines
- **SSL Certificates**: For gRPC communication

#### API Server Dockerfile (Machine 1)
```dockerfile
# cmd/api/Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the API server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-server ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/api-server .

# Create non-root user
RUN addgroup -g 1000 englog && \
    adduser -D -s /bin/sh -u 1000 -G englog englog

USER englog

EXPOSE 8080 8443 50051

CMD ["./api-server"]
```

#### Worker Server Dockerfile (Machine 2)
```dockerfile
# cmd/worker/Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the worker server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker ./cmd/worker

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl
WORKDIR /root/

COPY --from=builder /app/worker .

RUN addgroup -g 1000 englog && \
    adduser -D -s /bin/sh -u 1000 -G englog englog

USER englog

CMD ["./worker"]
```

### SSL/TLS Setup

#### Generate Certificates for gRPC Communication

```bash
# Create CA certificate
openssl genrsa -out ca.key 4096
openssl req -new -x509 -key ca.key -sha256 -subj "/C=US/ST=CA/O=EngLog/CN=ca" -days 3650 -out ca.crt

# Create API server certificate (Machine 1)
openssl genrsa -out server.key 4096
openssl req -new -key server.key -out server.csr -subj "/C=US/ST=CA/O=EngLog/CN=203.0.113.10"
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256

# Create worker client certificate (Machine 2)
openssl genrsa -out client.key 4096
openssl req -new -key client.key -out client.csr -subj "/C=US/ST=CA/O=EngLog/CN=worker"
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365 -sha256
```

### Step-by-Step Deployment

#### 1. Prepare Machine 1 (API Server)

```bash
# Clone repository
git clone https://github.com/your-org/englog.git
cd englog

# Set environment variables
export DB_PASSWORD="your_secure_password_here"

# Create certificates directory
mkdir -p certs
# Copy server.crt, server.key, ca.crt to ./certs/

# Deploy API server
docker-compose -f docker-compose.api.yml up -d

# Verify services
docker-compose -f docker-compose.api.yml ps
```

#### 2. Prepare Machine 2 (Worker Server)

```bash
# Clone repository
git clone https://github.com/your-org/englog.git
cd englog

# Create certificates directory
mkdir -p certs
# Copy client.crt, client.key, ca.crt to ./certs/

# Deploy worker server
docker-compose -f docker-compose.worker.yml up -d

# Setup Ollama models
docker exec -it englog_ollama_1 ollama pull llama3.2
docker exec -it englog_ollama_1 ollama pull codellama

# Verify services
docker-compose -f docker-compose.worker.yml ps
```
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/ssl/certs
    networks:
      - public
    depends_on:
      - api-server

  # API Server (Public Network)
  api-server:
    build:
      context: .
      dockerfile: Dockerfile.api-server
    ports:
      - "8080:8080"  # HTTP API
      - "9090:9090"  # gRPC Server
    environment:
      - DATABASE_URL=postgresql://englog:password@postgres:5432/englog
      - REDIS_URL=redis://redis:6379
      - GRPC_SERVER_PORT=9090
      - WORKER_GRPC_ADDRESS=worker-server:9091
      - JWT_SECRET_KEY=${JWT_SECRET_KEY}
      - TLS_CERT_FILE=/app/certs/server.crt
      - TLS_KEY_FILE=/app/certs/server.key
    volumes:
      - ./logs:/app/logs
      - ./certs:/app/certs
    networks:
      - public
      - private
    depends_on:
      - postgres
      - redis
      - worker-server

  # Worker Server (Private Network Only)
  worker-server:
    build:
      context: .
      dockerfile: Dockerfile.worker-server
    ports:
      - "9091:9091"  # gRPC Client (private)
    environment:
      - API_SERVER_GRPC_ADDRESS=api-server:9090
      - OLLAMA_URL=http://ollama:11434
      - POSTGRES_READ_URL=postgresql://englog:password@postgres:5432/englog
      - SMTP_HOST=${SMTP_HOST}
      - SMTP_PORT=${SMTP_PORT}
      - SMTP_USERNAME=${SMTP_USERNAME}
      - SMTP_PASSWORD=${SMTP_PASSWORD}
      - TLS_CERT_FILE=/app/certs/client.crt
      - TLS_KEY_FILE=/app/certs/client.key
    volumes:
      - ./logs:/app/logs
      - ./certs:/app/certs
    networks:
      - private
    depends_on:
      - ollama
      - postgres

  # Database (Shared)
  postgres:
    image: postgres:17
    environment:
      - POSTGRES_DB=englog
      - POSTGRES_USER=englog
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d/
    networks:
      - private
    ports:
      - "5432:5432"  # For development only

  # Redis Cache (API Server Only)
  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    networks:
      - private

  # Ollama LLM (Worker Server Only)
  ollama:
    image: ollama/ollama:latest
    volumes:
      - ollama_data:/root/.ollama
    networks:
      - private
    environment:
      - OLLAMA_HOST=0.0.0.0
    command: ollama serve

  # Prometheus Monitoring
  prometheus:
    image: prom/prometheus
    ports:
      - "9093:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    networks:
      - public
      - private

  # Grafana Dashboard
  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana
    networks:
      - public

volumes:
  postgres_data:
  redis_data:
  ollama_data:
  prometheus_data:
  grafana_data:
```

```

### Monitoring and Maintenance

#### Health Checks

```bash
# Check API Server health (Machine 1)
curl https://203.0.113.10/health

# Check database connection
docker exec englog_postgres_1 pg_isready -U englog

# Check Redis
docker exec englog_redis_1 redis-cli ping

# Check Worker Server (from Machine 1)
curl -k https://10.0.1.20:50051/health

# Check Ollama service
curl http://10.0.1.20:11434/api/tags
```

#### Basic Monitoring Scripts

```bash
#!/bin/bash
# monitor.sh - Simple monitoring script

# Machine 1 monitoring
echo "=== API Server Status ==="
docker-compose -f docker-compose.api.yml ps

echo "=== Database Size ==="
docker exec englog_postgres_1 psql -U englog -c "SELECT pg_size_pretty(pg_database_size('englog'));"

echo "=== Redis Memory ==="
docker exec englog_redis_1 redis-cli info memory | grep used_memory_human

# Machine 2 monitoring
echo "=== Worker Server Status ==="
docker-compose -f docker-compose.worker.yml ps

echo "=== Ollama Models ==="
docker exec englog_ollama_1 ollama list

echo "=== Disk Usage ==="
df -h /var/lib/docker
```

#### Backup and Recovery

```bash
#!/bin/bash
# backup.sh - Simple backup script for Machine 1

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/englog_$DATE"

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup PostgreSQL
docker exec englog_postgres_1 pg_dump -U englog englog > $BACKUP_DIR/database.sql

# Backup Redis (if needed)
docker exec englog_redis_1 redis-cli BGSAVE
docker cp englog_redis_1:/data/dump.rdb $BACKUP_DIR/

# Backup application logs
cp -r ./logs $BACKUP_DIR/

# Create archive
tar -czf /backup/englog_backup_$DATE.tar.gz $BACKUP_DIR
rm -rf $BACKUP_DIR

echo "Backup completed: englog_backup_$DATE.tar.gz"
```

#### Log Management

```bash
# View API server logs
docker-compose -f docker-compose.api.yml logs -f api-server

# View worker logs
docker-compose -f docker-compose.worker.yml logs -f worker-server

# View database logs
docker-compose -f docker-compose.api.yml logs postgres

# Cleanup old logs (weekly cron job)
find ./logs -name "*.log" -mtime +30 -delete
```

### Troubleshooting

#### Common Issues

1. **gRPC Connection Failed**
   ```bash
   # Check network connectivity
   ping 10.0.1.20

   # Verify certificates
   openssl x509 -in certs/server.crt -text -noout

   # Check firewall rules
   sudo ufw status
   ```

2. **Database Connection Issues**
   ```bash
   # Check PostgreSQL status
   docker exec englog_postgres_1 pg_isready

   # View connection logs
   docker logs englog_postgres_1
   ```

3. **Ollama Model Loading**
   ```bash
   # Check available models
   docker exec englog_ollama_1 ollama list

   # Check GPU availability
   docker exec englog_ollama_1 nvidia-smi

   # View Ollama logs
   docker logs englog_ollama_1
   ```

#### Performance Tuning

```bash
# PostgreSQL tuning (add to docker-compose.api.yml)
# environment:
#   - POSTGRES_SHARED_BUFFERS=256MB
#   - POSTGRES_EFFECTIVE_CACHE_SIZE=1GB
#   - POSTGRES_WORK_MEM=4MB

# Redis tuning
# command: redis-server --maxmemory 512mb --maxmemory-policy allkeys-lru

# Ollama GPU optimization
# Deploy with GPU support for better LLM performance
```

---

#### Monitoring & Metrics
```go
package monitoring

import (
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next()

        status := c.Writer.Status()
        method := c.Request.Method
        endpoint := c.FullPath()

        httpRequestsTotal.WithLabelValues(method, endpoint, fmt.Sprintf("%d", status)).Inc()
        httpRequestDuration.WithLabelValues(method, endpoint).Observe(time.Since(start).Seconds())
    }
}

// RecordLogEntryCreated increments the log entries counter
func RecordLogEntryCreated() {
    logEntriesCreated.Inc()
}

// RecordInsightGeneration records the time taken to generate insights
func RecordInsightGeneration(duration time.Duration) {
    insightGenerationTime.Observe(duration.Seconds())
}

// UpdateActiveUsers updates the active users gauge
func UpdateActiveUsers(count float64) {
    activeUsers.Set(count)
}

// UpdateDatabaseConnections updates the database connections gauge
func UpdateDatabaseConnections(count float64) {
    databaseConnections.Set(count)
}

// Handler returns the Prometheus metrics handler
func Handler() gin.HandlerFunc {
    h := promhttp.Handler()
    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}

// Background metrics collector
type MetricsCollector struct {
    pool    *pgxpool.Pool
    queries *sqlc.Queries
}

func NewMetricsCollector(pool *pgxpool.Pool, queries *sqlc.Queries) *MetricsCollector {
    return &MetricsCollector{
        pool:    pool,
        queries: queries,
    }
}

func (m *MetricsCollector) Start(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            m.collectMetrics(ctx)
        }
    }
}

func (m *MetricsCollector) collectMetrics(ctx context.Context) {
    // Collect active users (users with activity in last 24h)
    yesterday := time.Now().Add(-24 * time.Hour)

    activeUserCount, err := m.queries.GetActiveUserCount(ctx, yesterday)
    if err == nil {
        UpdateActiveUsers(float64(activeUserCount))
    }

    // Collect database connection metrics
    if m.pool != nil {
        stats := m.pool.Stat()
        UpdateDatabaseConnections(float64(stats.TotalConns()))
    }
}
```

#### Health Checks
```go
package health

import (
    "context"
    "database/sql"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "github.com/jackc/pgx/v5/pgxpool"
)

type HealthStatus struct {
    Status    string            `json:"status"`
    Timestamp string            `json:"timestamp"`
    Services  map[string]string `json:"services"`
}

type HealthChecker struct {
    pool        *pgxpool.Pool
    redis       *redis.Client
    ollamaURL   string
    httpClient  *http.Client
}

func NewHealthChecker(pool *pgxpool.Pool, redisClient *redis.Client, ollamaURL string) *HealthChecker {
    return &HealthChecker{
        pool:      pool,
        redis:     redisClient,
        ollamaURL: ollamaURL,
        httpClient: &http.Client{
            Timeout: 5 * time.Second,
        },
    }
}
        db:        db,
        redis:     redisClient,
        ollamaURL: ollamaURL,
        httpClient: &http.Client{
            Timeout: 5 * time.Second,
        },
    }
}

func (h *HealthChecker) HealthCheck(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
    defer cancel()

    healthStatus := &HealthStatus{
        Status:    "healthy",
        Timestamp: time.Now().UTC().Format(time.RFC3339),
        Services:  make(map[string]string),
    }

    // Database check
    if err := h.checkDatabase(ctx); err != nil {
        healthStatus.Services["database"] = "unhealthy: " + err.Error()
        healthStatus.Status = "unhealthy"
    } else {
        healthStatus.Services["database"] = "healthy"
    }

    // Redis check
    if err := h.checkRedis(ctx); err != nil {
        healthStatus.Services["redis"] = "unhealthy: " + err.Error()
        healthStatus.Status = "unhealthy"
    } else {
        healthStatus.Services["redis"] = "healthy"
    }

    // Ollama check
    if err := h.checkOllama(ctx); err != nil {
        healthStatus.Services["ollama"] = "unhealthy: " + err.Error()
        // Ollama being down is not critical, so only mark as degraded
        if healthStatus.Status == "healthy" {
            healthStatus.Status = "degraded"
        }
    } else {
        healthStatus.Services["ollama"] = "healthy"
    }

    statusCode := http.StatusOK
    if healthStatus.Status == "unhealthy" {
        statusCode = http.StatusServiceUnavailable
    }

    c.JSON(statusCode, healthStatus)
}

func (h *HealthChecker) ReadinessCheck(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()

    // For readiness, we only check critical services
    if err := h.checkDatabase(ctx); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "not ready",
            "reason": "database unavailable",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "ready",
    })
}

func (h *HealthChecker) LivenessCheck(c *gin.Context) {
    // Simple liveness check - if we can respond, we're alive
    c.JSON(http.StatusOK, gin.H{
        "status": "alive",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
    })
}

func (h *HealthChecker) checkDatabase(ctx context.Context) error {
    sqlDB, err := h.db.DB()
    if err != nil {
        return err
    }

    return sqlDB.PingContext(ctx)
}

func (h *HealthChecker) checkRedis(ctx context.Context) error {
    return h.redis.Ping(ctx).Err()
}

func (h *HealthChecker) checkOllama(ctx context.Context) error {
    req, err := http.NewRequestWithContext(ctx, "GET", h.ollamaURL+"/api/tags", nil)
    if err != nil {
        return err
    }

    resp, err := h.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("ollama returned status %d", resp.StatusCode)
    }

    return nil
}

// Middleware to check if service is healthy before processing requests
func (h *HealthChecker) HealthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip health check for health endpoints
        if c.Request.URL.Path == "/health" ||
           c.Request.URL.Path == "/ready" ||
           c.Request.URL.Path == "/live" {
            c.Next()
            return
        }

        ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
        defer cancel()

        // Quick database check
        if err := h.checkDatabase(ctx); err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{
                "error": "Service temporarily unavailable",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// Background health monitor
func (h *HealthChecker) StartHealthMonitor(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            h.performHealthCheck(ctx)
        }
    }
}

func (h *HealthChecker) performHealthCheck(ctx context.Context) {
    checkCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // Log health status
    dbStatus := "healthy"
    if err := h.checkDatabase(checkCtx); err != nil {
        dbStatus = "unhealthy: " + err.Error()
        // Could send alert here
    }

    redisStatus := "healthy"
    if err := h.checkRedis(checkCtx); err != nil {
        redisStatus = "unhealthy: " + err.Error()
        // Could send alert here
    }

    ollamaStatus := "healthy"
    if err := h.checkOllama(checkCtx); err != nil {
        ollamaStatus = "degraded: " + err.Error()
    }

    // Log or send metrics
    fmt.Printf("Health Check - DB: %s, Redis: %s, Ollama: %s\n",
        dbStatus, redisStatus, ollamaStatus)
}
```

---

## Example Usage Scenarios

### Scenario 1: Daily Activity Logging

```json
{
  "title": "Implement OAuth 2.0 refresh token rotation",
  "description": "Added security enhancement to prevent token replay attacks by implementing automatic token rotation on each refresh request",
  "type": "code",
  "project_id": "auth-modernization-project",
  "tags": ["security", "authentication", "api-enhancement", "p1-feature"],
  "start_time": "2024-07-15T09:00:00Z",
  "end_time": "2024-07-15T13:30:00Z"
}
```

### Scenario 2: Sprint Performance Analysis

**Generated Insight Example:**
```markdown
# Sprint 16 Performance Summary

## Key Accomplishments
- **Major Feature Delivery**: Successfully implemented and deployed OAuth 2.0 refresh token rotation, enhancing application security
- **Performance Optimization**: Achieved 90% reduction in geofencing query latency through PostGIS spatial indexing
- **Team Collaboration**: Conducted 3 code reviews for junior developers and mentored 2 team members

## Time Allocation Analysis
- **Development Work**: 68% (54.4 hours)
  - Feature coding: 32 hours
  - Bug fixes and optimizations: 12.4 hours
  - Testing and validation: 10 hours
- **Meetings & Collaboration**: 22% (17.6 hours)
  - Daily standups: 5 hours
  - Sprint ceremonies: 6 hours
  - Cross-team syncs: 6.6 hours
- **Documentation & Planning**: 10% (8 hours)

## Technical Contributions
- Introduced spatial database optimizations that will benefit all location-based queries
- Enhanced security posture through advanced authentication mechanisms
- Established performance benchmarking processes for database optimizations

## Areas of Excellence
- **Problem-Solving**: Quickly identified and resolved production geofencing performance issues
- **Architecture Design**: Led architectural review and gained approval for complex database schema changes
- **Knowledge Sharing**: Effectively mentored junior developers on authentication flows and testing practices

## Recommendations for Next Sprint
- Continue focus on performance optimization initiatives
- Increase involvement in early architectural discussions for the Social Check-ins epic
- Consider dedicating time blocks for deep focus work to maintain high coding productivity
```

### Scenario 3: Annual Performance Review

**Generated Annual Review Document:**
```markdown
# Annual Performance Review - Senior Software Engineer

## Executive Summary
This review covers contributions from January 1, 2024 to December 31, 2024, representing 2,080 logged work hours across 1,247 activities. Demonstrated consistent high performance in technical delivery, team collaboration, and professional development.

## Major Achievements

### Technical Delivery
- **OAuth 2.0 Modernization**: Led complete authentication system overhaul, improving security posture and user experience
- **Database Performance Initiative**: Delivered 3 major optimizations resulting in 75% average query performance improvement
- **API Architecture Evolution**: Designed and implemented 12 new API endpoints with comprehensive documentation

### Project Leadership
- **Geofencing Optimization Project**: Technical lead for critical performance initiative, delivered 2 weeks ahead of schedule
- **Social Check-ins Epic**: Provided architectural guidance and mentorship for 6-month feature development

### Team Impact
- **Mentorship**: Guided 4 junior developers through 47 code reviews and 23 pairing sessions
- **Knowledge Sharing**: Conducted 8 technical presentations and created 15 architectural decision documents
- **Cross-team Collaboration**: Facilitated 34 cross-team syncs, improving mobile-backend alignment

### Performance Metrics

#### Productivity Analysis
- **Average Weekly Hours**: 42.3 hours (consistent and sustainable)
- **Code Quality**: 96% of PRs approved on first review
- **Delivery Reliability**: 100% of commitments delivered on time or early

#### Skill Development
- **Technical Growth**: Advanced PostgreSQL optimization, spatial databases, advanced authentication patterns
- **Leadership Skills**: Increased mentorship capacity by 40% year-over-year
- **Communication**: Led 12 architectural reviews with stakeholder alignment

### Goal Achievement
- ✅ **Technical Excellence**: Exceeded expectations on code quality and system performance
- ✅ **Team Leadership**: Successfully mentored 4 junior developers with positive feedback
- ✅ **Innovation**: Delivered 2 architectural improvements adopted company-wide
- ✅ **Professional Development**: Completed advanced database certification and security training

### Areas for Continued Growth
- **Product Strategy**: Increase involvement in product planning and system architecture design
- **External Visibility**: Consider speaking at conferences or writing technical blog posts
- **Cross-functional Skills**: Develop deeper understanding of distributed systems and microservices

### Recommendations for Next Year
- **Senior+ Promotion Track**: Ready for increased technical leadership responsibilities
- **Architecture Council**: Nominate for participation in company-wide architecture decisions
- **Mentorship Program**: Formalize mentorship role with structured junior developer assignments
```

---

## Future Roadmap

### Phase 1: Simple Two-Machine Foundation (Months 1-3) ✅
- ✅ Simple two-machine distributed architecture
- ✅ API Server (Machine 1) - Public access with PostgreSQL and Redis
- ✅ Worker Server (Machine 2) - Private LLM processing with Ollama
- ✅ Basic gRPC communication with TLS
- ✅ Docker Compose deployment for both machines
- ✅ Core activity logging REST API
- ✅ Project and tag management API
- ✅ Swagger/OpenAPI documentation as initial client
- ✅ Basic authentication and authorization
- ✅ PostgreSQL database with sqlc queries

### Phase 2: Core API Enhancement (Months 4-6)
- 📋 Enhanced productivity metrics API endpoints
- 📋 Background task processing via Worker Server
- 📋 Email notification system
- 📋 Basic trend analysis and reporting APIs
- 📋 Improved error handling and logging
- 📋 API security hardening (rate limiting, input validation)
- 📋 Simple backup and monitoring scripts
- 📋 Data export features for GDPR compliance

### Phase 3: React Frontend Migration (Months 7-9)
- 📋 Replace Swagger UI with React.js web application
- 📋 Modern dashboard for activity visualization
- 📋 Interactive charts using chart.js or recharts
- 📋 Responsive design for mobile and desktop
- 📋 Client-side authentication and routing
- 📋 Real-time activity logging interface
- 📋 Progressive Web App (PWA) for offline access
- 📋 Activity management and editing interface

### Phase 4: AI-Enhanced User Experience (Months 10-12)
- 📋 Smart activity categorization using LLM
- 📋 AI-generated productivity insights dashboard
- 📋 Automatic time tracking suggestions
- 📋 Personalized recommendation system
- 📋 Enhanced LLM prompting and context management
- 📋 Voice activity logging (Web Speech API)
- 📋 Predictive analytics visualization
- 📋 Custom productivity goal tracking

### Phase 5: Scaling and Performance (Months 13-18)
- 📋 Machine performance optimization
- 📋 Database performance tuning
- 📋 Enhanced Redis caching strategies
- 📋 Optional horizontal scaling (additional Worker machines)
- 📋 Load balancer introduction (when traffic demands)
- 📋 Advanced monitoring and alerting
- 📋 Multi-user team features
- 📋 Admin panel for user management

### Phase 6: Enterprise Features (Months 19-24)
- 📋 Calendar synchronization (Google, Outlook)
- 📋 Development tool integrations (GitHub, Jira, Slack)
- 📋 Mobile app development (React Native)
- 📋 Third-party API integrations
- 📋 Advanced analytics and reporting
- 📋 Multi-language support (i18n)
- 📋 Optional Kubernetes migration (for large deployments)
- 📋 Plugin architecture and ecosystem

### Deployment Evolution Timeline

#### Current: Simple Two-Machine Setup
```
Initial Deployment
├── Machine 1: API Server (Public IP)
│   ├── Go API service
│   ├── PostgreSQL database
│   ├── Redis cache
│   └── gRPC server
└── Machine 2: Worker Server (Private IP)
    ├── Go worker service
    ├── Ollama LLM
    ├── gRPC client
    └── Background processing
```

#### Phase 3-4: React Web Interface
```
Frontend Development
├── React.js replacing Swagger UI
├── Modern dashboard interface
├── Interactive activity management
├── Real-time data visualization
├── Progressive Web App features
└── Mobile-responsive design
```

#### Phase 5-6: Scalable Platform
```
Enterprise Growth
├── Optional load balancer
├── Additional worker machines
├── Enhanced monitoring
├── Team collaboration features
├── Third-party integrations
└── Optional Kubernetes migration
```

### Architecture Philosophy

#### Current: Simplicity First
- Two-machine distributed architecture
- Docker Compose deployment
- Direct gRPC communication
- Minimal operational complexity

#### Future: Scalable Growth
- Gradual scaling based on actual needs
- Preserve distributed benefits
- Clear upgrade paths without major rewrites

---

## Conclusion

EngLog represents a **simple yet powerful** solution for software engineers to systematically track, analyze, and optimize their professional activities. The **two-machine distributed architecture** ensures security and maintainability while providing a clear evolution path from basic deployment to enterprise-scale platform.

### Key Architectural Decisions

1. **Simple Distributed Design**: Two-machine setup balances distribution benefits with operational simplicity
2. **API-First Approach**: Starting with comprehensive REST API and Swagger documentation ensures solid foundation
3. **Direct gRPC Communication**: Simple TLS communication between machines without complex service mesh
4. **Docker Compose Deployment**: Easy setup and maintenance without Kubernetes complexity
5. **Security by Design**: Network isolation, TLS encryption, and proper authentication

### Value Proposition

By combining structured data collection with AI-powered insights, delivered through a **simple yet scalable architecture**, EngLog enables:

- **Cost-effective deployment** with just two machines for initial version
- **Data-driven career development** through comprehensive activity tracking
- **Streamlined performance management** with AI-generated insights
- **Clear scaling path** from simple setup to enterprise platform
- **Future-proof foundation** ready for growth without major rewrites

The roadmap ensures **immediate value delivery** with the simple two-machine setup while building toward advanced features that will set new standards for professional activity tracking in the software engineering field.
