# HTTP Handlers and API Endpoints Documentation

This document provides comprehensive documentation for all HTTP handlers and API endpoints implemented in Task 0070.

## Overview

The EngLog REST API provides complete functionality for managing activity logs, projects, analytics, tags, and user profiles. All endpoints follow RESTful conventions and implement proper authentication, validation, and error handling.

## Authentication

Most endpoints require authentication via JWT Bearer tokens. Include the token in the Authorization header:

```
Authorization: Bearer <access_token>
```

## Common Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "success": true,
  "data": <response_data>,
  "message": "Optional success message"
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message",
  "message": "Optional detailed message"
}
```

### Paginated Response
```json
{
  "success": true,
  "data": {
    "data": [<items>],
    "pagination": {
      "page": 1,
      "limit": 50,
      "total": 100,
      "total_pages": 2,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

## Endpoints

### Health Endpoints

#### GET /health
Health check endpoint (no authentication required)

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "uptime": "2h30m15s",
  "version": "1.0.0"
}
```

#### GET /ready
Readiness check endpoint (no authentication required)

**Response:**
```json
{
  "status": "ready",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Authentication Endpoints

#### POST /v1/auth/register
User registration (no authentication required)

#### POST /v1/auth/login
User login (no authentication required)

#### POST /v1/auth/refresh
Refresh access token (no authentication required)

#### POST /v1/auth/logout
User logout (no authentication required)

#### GET /v1/auth/me
Get current user profile (requires authentication)

### Log Entries

#### POST /v1/logs
Create a new log entry

**Authentication:** Required

**Request Body:**
```json
{
  "title": "Implemented OAuth 2.0 refresh tokens",
  "description": "Added token rotation and security improvements",
  "type": "development",
  "project_id": "550e8400-e29b-41d4-a716-446655440000",
  "start_time": "2024-01-15T09:00:00Z",
  "end_time": "2024-01-15T12:30:00Z",
  "value_rating": "high",
  "impact_level": "team",
  "tags": ["oauth", "security", "api"]
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "title": "Implemented OAuth 2.0 refresh tokens",
    "description": "Added token rotation and security improvements",
    "type": "development",
    "project_id": "550e8400-e29b-41d4-a716-446655440000",
    "start_time": "2024-01-15T09:00:00Z",
    "end_time": "2024-01-15T12:30:00Z",
    "duration_minutes": 210,
    "value_rating": "high",
    "impact_level": "team",
    "tags": ["oauth", "security", "api"],
    "created_at": "2024-01-15T12:35:00Z",
    "updated_at": "2024-01-15T12:35:00Z"
  }
}
```

#### GET /v1/logs
Get log entries with filtering and pagination

**Authentication:** Required

**Query Parameters:**
- `start_date` (string): Filter by start date (YYYY-MM-DD)
- `end_date` (string): Filter by end date (YYYY-MM-DD)
- `project_id` (string): Filter by project UUID
- `type` (string): Filter by activity type
- `value_rating` (string): Filter by value rating
- `impact_level` (string): Filter by impact level
- `tags` (string): Comma-separated list of tags
- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: 50, max: 100)

**Example:** `GET /v1/logs?start_date=2024-01-01&type=development&page=1&limit=10`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "data": [/* array of log entries */],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 45,
      "total_pages": 5,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

#### GET /v1/logs/:id
Get a specific log entry

**Authentication:** Required

**Response:** `200 OK` or `404 Not Found`

#### PUT /v1/logs/:id
Update a log entry

**Authentication:** Required

**Request Body:** Same as POST /v1/logs

**Response:** `200 OK` or `404 Not Found`

#### DELETE /v1/logs/:id
Delete a log entry

**Authentication:** Required

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Log entry deleted successfully"
}
```

#### POST /v1/logs/bulk
Bulk create log entries

**Authentication:** Required

**Request Body:**
```json
{
  "entries": [
    /* array of log entry objects (max 100) */
  ]
}
```

**Response:** `201 Created`, `207 Multi-Status`, or `400 Bad Request`
```json
{
  "success": true,
  "data": [/* array of results */],
  "summary": {
    "total": 5,
    "success": 4,
    "errors": 1
  }
}
```

### Projects

#### POST /v1/projects
Create a new project

**Authentication:** Required

**Request Body:**
```json
{
  "name": "Authentication System",
  "description": "OAuth 2.0 and JWT implementation",
  "color": "#FF6B6B",
  "status": "active",
  "start_date": "2024-01-01",
  "end_date": "2024-03-31",
  "is_default": false
}
```

#### GET /v1/projects
Get all user projects

**Authentication:** Required

#### GET /v1/projects/:id
Get a specific project

**Authentication:** Required

#### PUT /v1/projects/:id
Update a project

**Authentication:** Required

#### DELETE /v1/projects/:id
Delete a project

**Authentication:** Required

### Analytics

#### GET /v1/analytics/productivity
Get productivity metrics

**Authentication:** Required

**Query Parameters:**
- `start_date` (string): Start date (YYYY-MM-DD)
- `end_date` (string): End date (YYYY-MM-DD)

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "total_activities": 45,
    "total_minutes": 2100,
    "projects_worked": 3,
    "high_value_activities": 15,
    "activity_breakdown": {
      "development": 1200,
      "meeting": 450,
      "code_review": 300,
      "documentation": 150
    },
    "value_distribution": {
      "low": 5,
      "medium": 25,
      "high": 15
    },
    "impact_distribution": {
      "personal": 10,
      "team": 20,
      "department": 10,
      "company": 5
    }
  },
  "period": {
    "start_date": "2024-01-01",
    "end_date": "2024-01-31"
  }
}
```

#### GET /v1/analytics/summary
Get activity summary

**Authentication:** Required

**Query Parameters:** Same as productivity metrics

### Tags

#### POST /v1/tags
Create a new tag

**Authentication:** Required

**Request Body:**
```json
{
  "name": "security",
  "color": "#FF0000",
  "description": "Security-related tasks"
}
```

#### GET /v1/tags
Get all tags

**Authentication:** Required

#### GET /v1/tags/popular
Get popular tags

**Authentication:** Required

**Query Parameters:**
- `limit` (int): Maximum number of tags (default: 10, max: 50)

#### GET /v1/tags/recent
Get recently used tags for the current user

**Authentication:** Required

**Query Parameters:**
- `limit` (int): Maximum number of tags (default: 10, max: 50)

#### GET /v1/tags/search
Search tags by name

**Authentication:** Required

**Query Parameters:**
- `q` (string): Search query (required)
- `limit` (int): Maximum number of results (default: 20, max: 50)

#### GET /v1/tags/usage
Get tag usage statistics for the current user

**Authentication:** Required

#### GET /v1/tags/:id
Get a specific tag

**Authentication:** Required

#### PUT /v1/tags/:id
Update a tag

**Authentication:** Required

#### DELETE /v1/tags/:id
Delete a tag

**Authentication:** Required

### User Profile

#### GET /v1/users/profile
Get current user profile

**Authentication:** Required

#### PUT /v1/users/profile
Update user profile

**Authentication:** Required

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "timezone": "America/New_York",
  "preferences": {
    "theme": "dark",
    "notifications": true
  }
}
```

#### POST /v1/users/change-password
Change user password

**Authentication:** Required

**Request Body:**
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

#### DELETE /v1/users/account
Delete user account

**Authentication:** Required

## Error Handling

The API uses standard HTTP status codes:

- `200 OK` - Successful GET, PUT requests
- `201 Created` - Successful POST requests
- `204 No Content` - Successful DELETE requests
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required or invalid
- `403 Forbidden` - Access denied
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Validation errors
- `500 Internal Server Error` - Server errors

## Validation

All request bodies are validated according to the model definitions:

### Activity Types
- `development`, `meeting`, `code_review`, `debugging`, `documentation`, `testing`, `deployment`, `research`, `planning`, `learning`, `maintenance`, `support`, `other`

### Value Ratings
- `low`, `medium`, `high`, `critical`

### Impact Levels
- `personal`, `team`, `department`, `company`

### Project Status
- `active`, `completed`, `on_hold`, `cancelled`

## Rate Limiting

The API implements basic rate limiting to prevent abuse. In production, this should be enhanced with Redis-based rate limiting.

## CORS

Cross-Origin Resource Sharing (CORS) is enabled for all origins in development. In production, this should be restricted to specific domains.

## Implementation Notes

- All handlers use the standardized response utilities
- Consistent error handling across all endpoints
- Proper HTTP status codes for different scenarios
- Input validation and sanitization
- Authentication middleware integration
- Pagination support for list endpoints
- Filtering capabilities for log entries
- Comprehensive test coverage for critical paths
