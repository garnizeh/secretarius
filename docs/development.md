# EngLog Development Guide

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:
- Go 1.24+
- PostgreSQL 17+
- Redis 7+
- Docker & Docker Compose
- Make (for running Makefile commands)

### Setup Development Environment

1. **Clone the repository**
   ```bash
   git clone https://github.com/garnizeh/englog.git
   cd englog
   ```

2. **Install development tools**
   ```bash
   make install-tools
   ```

3. **Setup environment variables**
   ```bash
   # Copy environment templates for development
   cp .env.example .env.api-dev
   cp .env.example .env.worker-dev
   # Edit .env.api-dev and .env.worker-dev with your configuration
   # Or use the provided .env.api-dev and .env.worker-dev files which are pre-configured for development
   ```

4. **Start local services**
   ```bash
   make dev-up
   ```

5. **Run database migrations**
   ```bash
   make migrate-up
   ```

6. **Start the API server**
   ```bash
   make run-api
   ```

7. **Start the worker server** (in another terminal)
   ```bash
   make run-worker
   ```

The API will be available at http://localhost:8080

## Development Workflow

### Making Changes

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Follow Go conventions and best practices
   - Add tests for new functionality
   - Update documentation as needed

3. **Run tests and checks**
   ```bash
   make check  # Runs lint, test, and security checks
   ```

4. **Generate code if needed**
   ```bash
   make generate  # Generates sqlc, protobuf, swagger
   ```

5. **Commit and push**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   git push origin feature/your-feature-name
   ```

### Code Generation

EngLog uses several code generation tools:

- **sqlc**: Generates type-safe Go code from SQL queries
- **protobuf**: Generates gRPC service definitions
- **swagger**: Generates API documentation

Run all generators:
```bash
make generate
```

Or run individually:
```bash
make sqlc    # Generate database code
make proto   # Generate gRPC code
make swagger # Generate API docs
```

### Database Development

#### Creating Migrations
```bash
make migrate-create NAME=add_user_table
```

This creates two files:
- `migrations/NNNN_add_user_table.up.sql`
- `migrations/NNNN_add_user_table.down.sql`

#### Running Migrations
```bash
make migrate-up    # Apply all pending migrations
make migrate-down  # Rollback last migration
```

#### Adding SQL Queries
1. Add SQL queries to files in `internal/sqlc/queries/`
2. Run `make sqlc` to generate Go code
3. Use the generated code in your services

### Testing

#### Running Tests
```bash
make test              # Run all tests
make test-coverage     # Run tests with coverage report
make benchmark         # Run benchmarks
```

#### Writing Tests
- Unit tests: `*_test.go` files alongside your code
- Integration tests: Place in `tests/` directory
- Use table-driven tests where appropriate
- Mock external dependencies

Example test:
```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateUserRequest
        want    *User
        wantErr bool
    }{
        {
            name: "valid user",
            input: CreateUserRequest{
                Email:     "test@example.com",
                FirstName: "John",
                LastName:  "Doe",
            },
            want: &User{
                Email:     "test@example.com",
                FirstName: "John",
                LastName:  "Doe",
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Debugging

#### Using Delve Debugger
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug API server
dlv debug ./cmd/api

# Debug worker server
dlv debug ./cmd/worker
```

#### Debugging with VS Code
Add to `.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug API",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "./cmd/api",
            "env": {
                "APP_ENV": "development"
            }
        }
    ]
}
```

### Live Reload

For faster development, use Air for live reloading:

```bash
make watch-api     # Start API with live reload
make watch-worker  # Start worker with live reload
```

### Docker Development

#### Build Images
```bash
make docker-build
```

#### Development with Docker Compose
```bash
# Start all services
docker-compose -f deployments/docker-compose.dev.yml up

# View logs
make dev-logs

# Stop services
make dev-down
```

## Code Style and Standards

### Go Conventions
- Follow standard Go conventions
- Use `gofmt` and `goimports` (included in `make format`)
- Write meaningful variable and function names
- Add comments for exported functions

### Project Structure
- `cmd/`: Application entry points
- `internal/`: Private application code (not importable)
- `pkg/`: Public packages (importable by other projects)
- `migrations/`: Database schema changes
- `docs/`: Documentation

### Error Handling
```go
// Good: Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// Good: Use custom error types for domain errors
var ErrUserNotFound = errors.New("user not found")
```

### Logging
```go
// Use structured logging
log.Info("creating user",
    "email", user.Email,
    "id", user.ID)

// Log errors with context
log.Error("failed to save user",
    "error", err,
    "user_id", userID)
```

## Common Tasks

### Adding a New API Endpoint
1. Define the handler in `internal/handlers/`
2. Add route in the router setup
3. Add service logic in `internal/services/`
4. Add database queries if needed
5. Write tests
6. Update API documentation

### Adding a New Background Task
1. Define task in worker server
2. Add gRPC service definition if needed
3. Implement task processor
4. Add tests
5. Update documentation

### Adding Database Changes
1. Create migration with `make migrate-create`
2. Write SQL for schema changes
3. Add corresponding sqlc queries
4. Generate code with `make sqlc`
5. Update Go code to use new schema

## Troubleshooting

### Common Issues

1. **Build fails**: Run `go mod tidy` to ensure dependencies are correct
2. **Tests fail**: Check database connection and migrations
3. **Generated code outdated**: Run `make generate`
4. **Docker issues**: Check if ports are available and services are running

### Getting Help

- Check existing issues in the GitHub repository
- Review documentation in `docs/` directory
- Ask questions in team chat or create an issue

## Performance Considerations

- Use connection pooling for database connections
- Implement caching for frequently accessed data
- Profile your code with `go tool pprof`
- Monitor memory usage and goroutine leaks
- Use context for timeouts and cancellation
