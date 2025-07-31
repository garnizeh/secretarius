# EngLog Architecture Overview

## System Architecture

EngLog implements a distributed two-machine architecture designed for simplicity, security, and maintainability.

```
┌─────────────────┐      ┌─────────────────┐
│   Machine 1     │      │   Machine 2     │
│   API Server    │ gRPC │  Worker Server  │
│   (Public)      │◄────►│   (Private)     │
│                 │      │                 │
│ • REST API      │      │ • LLM Processing│
│ • PostgreSQL    │      │ • Ollama        │
│ • Redis         │      │ • Background    │
│ • Authentication│      │   Tasks         │
└─────────────────┘      └─────────────────┘
```

## Components

### Machine 1: API Server (Public)
- **Purpose**: Client-facing API and data storage
- **Network**: Public IP with internet access
- **Technologies**: Go, Gin, PostgreSQL, Redis
- **Responsibilities**:
  - Handle HTTP requests
  - User authentication & authorization
  - Data validation and persistence
  - Session management
  - Serve as gRPC server

### Machine 2: Worker Server (Private)
- **Purpose**: Background processing and AI operations
- **Network**: Private IP, no direct internet access
- **Technologies**: Go, Ollama, gRPC
- **Responsibilities**:
  - LLM-powered insight generation
  - Background task processing
  - Report generation
  - Email notifications
  - Analytics processing

## Communication Flow

1. **Client Request**: Frontend/API client → API Server (HTTPS)
2. **Processing**: API Server → Worker Server (gRPC/TLS)
3. **AI Operations**: Worker Server → Ollama LLM
4. **Response**: Worker Server → API Server → Client

## Data Flow

### Log Entry Creation
```
Client → API Server → Database → Background Task Queue → Worker Server
```

### Insight Generation
```
Client Request → API Server → gRPC → Worker Server → Ollama → Database → Client
```

## Security Model

### Network Security
- API Server: Public HTTPS (443), Private gRPC (50051)
- Worker Server: Private network only, gRPC client
- Database: Local connections only

### Authentication
- JWT tokens with refresh mechanism
- Rate limiting per user
- Session denylist support

### Data Protection
- TLS encryption for all communications
- Environment-based configuration
- No sensitive data in Worker Server

## Scalability Considerations

### Current: Simple Deployment
- Two VMs/containers
- Direct gRPC communication
- Local database connections

### Future: Horizontal Scaling
- Multiple API server instances
- Load balancer for API servers
- Message queue for worker communication
- Database connection pooling

## Technology Stack

### Backend
- **Language**: Go 1.24+
- **Web Framework**: Gin
- **Database**: PostgreSQL 17+ with sqlc
- **Cache**: Redis 7+
- **Communication**: gRPC with Protocol Buffers

### AI/ML
- **LLM Platform**: Ollama
- **Models**: Local deployment (Llama 2/3, Mistral)
- **Processing**: Async task queue

### DevOps
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Monitoring**: Prometheus + Grafana
- **Deployment**: Simple two-machine setup

## Directory Structure

```
englog/
├── cmd/                    # Application entrypoints
│   ├── api/               # API server main
│   └── worker/            # Worker server main
├── internal/              # Private application code
│   ├── auth/             # Authentication logic
│   ├── handlers/         # HTTP handlers
│   ├── services/         # Business logic
│   ├── models/           # Data models
│   ├── sqlc/             # Generated DB code
│   └── config/           # Configuration
├── pkg/                   # Public packages
├── migrations/           # Database migrations
├── deployments/         # Docker configurations
├── scripts/             # Build/deploy scripts
└── docs/                # Documentation
```

## Deployment Model

### Development
```bash
# Start services
make dev-up

# Run API server
make run-api

# Run worker server
make run-worker
```

### Production
```bash
# Build and deploy
make docker-build
docker-compose up -d
```

This architecture provides a solid foundation that can grow from a simple two-machine setup to a full microservices architecture as needed.
