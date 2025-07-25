# EngLog API Documentation

## Overview

The EngLog API provides a comprehensive RESTful interface for managing work activity logs, projects, and generating insights for software engineers.

## Base URL

```
https://api.englog.dev/v1
```

For development:
```
http://localhost:8080/v1
```

## Authentication

All API endpoints (except authentication endpoints) require a valid JWT token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

### Authentication Endpoints

#### Login
```http
POST /v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your_password"
}
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe"
  }
}
```

## Core Endpoints

### Log Entries

#### Create Log Entry
```http
POST /v1/logs
Content-Type: application/json
Authorization: Bearer <token>

{
  "title": "Implemented user authentication",
  "description": "Added JWT-based authentication system with refresh tokens",
  "type": "development",
  "project_id": "uuid",
  "start_time": "2024-01-15T09:00:00Z",
  "end_time": "2024-01-15T11:30:00Z",
  "value_rating": "high",
  "impact_level": "team"
}
```

#### Get Log Entries
```http
GET /v1/logs?start_date=2024-01-01&end_date=2024-01-31&project_id=uuid
Authorization: Bearer <token>
```

### Projects

#### Create Project
```http
POST /v1/projects
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "EngLog API",
  "description": "Personal work activity tracker API",
  "color": "#3498db"
}
```

### Insights

#### Generate Insight
```http
POST /v1/insights/generate
Content-Type: application/json
Authorization: Bearer <token>

{
  "type": "weekly_summary",
  "period_start": "2024-01-01",
  "period_end": "2024-01-07"
}
```

## Error Responses

All errors follow the same format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request data",
    "details": [
      {
        "field": "email",
        "message": "Email is required"
      }
    ]
  }
}
```

## Rate Limiting

API requests are limited to 60 requests per minute per user. Rate limit headers are included in responses:

```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 59
X-RateLimit-Reset: 1642694400
```

## Interactive Documentation

For interactive API documentation, visit:
- Development: http://localhost:8080/swagger/
- Production: https://api.englog.dev/swagger/
