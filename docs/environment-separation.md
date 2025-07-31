# Environment Configuration Changes

## Summary
"Separation is the key to clarity, like distinct instruments in an orchestra." 🎼

The project environment configuration has been separated from a single `.env.dev` file into service-specific environment files:

- `.env.api-dev` - API server configuration
- `.env.worker-dev` - Worker server configuration

## Changes Made

### 1. Environment Files Created

#### `.env.api-dev`
Contains API server specific configurations:
- Application settings (port, host, environment)
- Database configuration (PostgreSQL)
- Redis cache configuration
- JWT authentication settings
- gRPC server configuration
- CORS and security settings
- Rate limiting configuration
- Monitoring and metrics
- Feature flags
- Background job scheduling

#### `.env.worker-dev`
Contains Worker server specific configurations:
- Worker identification and health settings
- gRPC client configuration (to connect to API)
- Ollama AI service configuration
- Email service settings (SMTP)
- Background job processing settings
- Development settings

### 2. Docker Compose Updates

Updated the following files to use the new environment files:
- `docker-compose.dev.yml` - Uses `.env.api-dev` for API and `.env.worker-dev` for worker
- `docker-compose.api-dev.yml` - Uses `.env.api-dev` for API and `.env.worker-dev` for worker
- `docker-compose.worker-dev.yml` - Uses `.env.worker-dev` for worker

### 3. VS Code Launch Configuration

Updated `.vscode/launch.json` to use:
- `.env.api-dev` for API server debugging
- `.env.worker-dev` for worker server debugging

## Benefits

1. **Separation of Concerns**: Each service has its own configuration scope
2. **Easier Maintenance**: Clearer understanding of what each service needs
3. **Reduced Complexity**: Smaller, focused environment files
4. **Better Security**: Can apply different security policies per service
5. **Deployment Flexibility**: Can deploy services independently with their own configs

## Migration from `.env.dev`

If you have a custom `.env.dev` file, you need to:

1. Split your custom values between `.env.api-dev` and `.env.worker-dev`
2. API-related configs go to `.env.api-dev`
3. Worker-related configs go to `.env.worker-dev`
4. Update any deployment scripts if they reference `.env.dev`

## Development Usage

```bash
# Start API development server
make run-api  # Uses .env.api-dev

# Start worker development server
make run-worker  # Uses .env.worker-dev

# Start full development environment
make dev-up  # Uses both files via docker-compose
```

## Note

The unified configuration loading in `internal/config/config.go` remains unchanged and continues to work with both environment files, as Go loads all available environment variables.
