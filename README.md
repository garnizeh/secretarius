# EngLog: Personal Work Activity Tracker for Software Engineers

> "The best way to predict the future is to create it" 📊

## Overview

EngLog is a specialized personal organizer designed for software engineers to capture, analyze, and derive insights from their daily work activities. The application serves as a comprehensive logging system that facilitates performance reviews, identifies professional growth patterns, and provides data-driven insights using LLM technology.

## Architecture

This project implements a distributed two-machine architecture:

- **Machine 1 (API Server)**: Public-facing Go REST API with PostgreSQL and Redis
- **Machine 2 (Worker Server)**: Private worker service with Ollama LLM integration

Communication between machines is handled via gRPC over TLS.

## Features

- **Activity Logging**: Comprehensive work activity tracking system
- **Project Management**: Organize activities by projects and teams
- **Intelligent Tagging**: Smart categorization and filtering
- **LLM-Powered Analytics**: AI-generated insights and reports
- **Data Export**: PDF/CSV exports for performance reviews
- **RESTful API**: Complete OpenAPI/Swagger documentation

## Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL 17+
- Redis 7+
- Docker & Docker Compose

### Installation

```bash
# Clone the repository
git clone https://github.com/garnizeh/englog.git
cd englog

# Install dependencies
go mod download

# Start development environment
make dev-up

# Run API server
make run-api

# Run worker server (separate terminal)
make run-worker
```

### Development

```bash
# Run tests
make test

# Run linting
make lint

# Generate code (sqlc, protobuf)
make generate

# Build binaries
make build

# Clean build artifacts
make clean
```

## Project Structure

```
.
├── cmd/
│   ├── api/          # API Server (Machine 1)
│   └── worker/       # Worker Server (Machine 2)
├── internal/
│   ├── auth/         # Authentication service
│   ├── models/       # Data models
│   ├── handlers/     # HTTP handlers
│   ├── services/     # Business logic
│   ├── sqlc/         # Generated database code
│   └── config/       # Configuration management
├── migrations/       # Database migrations
├── deployments/      # Docker configurations
├── pkg/              # Public packages
├── scripts/          # Build and deployment scripts
├── docs/             # Documentation
└── tests/            # Integration tests
```

## Technology Stack

### API Server (Machine 1)
- **Language**: Go 1.24+
- **Framework**: Gin HTTP Framework
- **Database**: PostgreSQL 17+ with sqlc
- **Cache**: Redis 7+ for sessions
- **Authentication**: JWT tokens
- **API Documentation**: Swagger/OpenAPI

### Worker Server (Machine 2)
- **Language**: Go 1.24+
- **LLM**: Ollama with local models
- **Queue**: In-memory with persistence
- **Communication**: gRPC client

## API Endpoints

- `POST /v1/auth/login` - User authentication
- `POST /v1/logs` - Create activity log
- `GET /v1/logs` - Retrieve activity logs
- `POST /v1/projects` - Create project
- `POST /v1/insights/generate` - Generate AI insights
- `GET /v1/analytics/dashboard` - Analytics dashboard

For complete API documentation, visit `/swagger/` when running the server.

## Configuration

Copy `.env.example` to `.env` and configure:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=englog
DB_USER=englog
DB_PASSWORD=your_password

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your_jwt_secret

# gRPC
GRPC_HOST=localhost
GRPC_PORT=50051
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) file for details.

## Documentation

- [API Documentation](docs/api.md)
- [Architecture Overview](docs/architecture.md)
- [Development Guide](docs/development.md)
- [Deployment Guide](docs/deployment.md)

## Support

For questions and support, please open an issue in the GitHub repository.
