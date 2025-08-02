# Docker Compose Configurations

> "Docker is a way to isolate your code from the world's complexities." - Solomon Hykes 🐳

This directory contains all Docker Compose configurations for different deployment environments and scenarios.

## Available Configurations

### Development Environments

- **`dev.yml`** - Complete development environment (API + Worker + Infrastructure)
- **`api-dev.yml`** - API server development environment only
- **`worker-dev.yml`** - Worker service development environment only
- **`infra-dev.yml`** - Infrastructure only (PostgreSQL + Redis)

### Production Environments

- **`api.yml`** - Production API server deployment (Machine 1)
- **`worker.yml`** - Production Worker service deployment (Machine 2)

### Testing Environment

- **`test.yml`** - Testing environment with isolated containers

## Usage Examples

### Development Workflow

```bash
# Start complete development environment
docker-compose -f deployments/docker-compose/dev.yml up -d

# Start infrastructure only
docker-compose -f deployments/docker-compose/infra-dev.yml up -d

# Start API development
docker-compose -f deployments/docker-compose/api-dev.yml up -d

# Start Worker development
docker-compose -f deployments/docker-compose/worker-dev.yml up -d
```

### Production Deployment

```bash
# Deploy API server (Machine 1)
docker-compose -f deployments/docker-compose/api.yml up -d

# Deploy Worker server (Machine 2)
docker-compose -f deployments/docker-compose/worker.yml up -d
```

### Testing

```bash
# Run tests with isolated environment
docker-compose -f deployments/docker-compose/test.yml up --abort-on-container-exit
```

## Environment Variables

Each configuration requires appropriate environment files:

- Development: `.env.dev`, `.env.api-dev`, `.env.worker-dev`
- Production: Environment variables injected via deployment system
- Testing: `.env.test` (auto-generated during test runs)

## Network Configuration

All configurations use the `englog-network` for service communication:

- **API Server**: Port 8080 (HTTP), 50051 (gRPC)
- **Worker Service**: Port 8091 (Health), gRPC client
- **PostgreSQL**: Port 5432 (internal)
- **Redis**: Port 6379 (internal)

## Volume Mounts

### Development
- Source code hot reload via volume mounts
- Persistent database data
- Log file persistence

### Production
- Optimized container images (no source mounts)
- Persistent data volumes
- Log aggregation ready

## Health Checks

All services include comprehensive health checks:

- **Startup probes**: Initial service availability
- **Liveness probes**: Ongoing service health
- **Readiness probes**: Service ready to accept traffic

## Security Considerations

### Development
- Local certificates for TLS testing
- Development-friendly security settings
- Debug logging enabled

### Production
- Production TLS certificates
- Hardened security configurations
- Structured logging only
- No debug endpoints exposed

## Troubleshooting

### Common Issues

1. **Port Conflicts**: Ensure ports 8080, 8091, 5432, 6379 are available
2. **Volume Permissions**: Check Docker daemon permissions for volume mounts
3. **Environment Variables**: Verify all required environment files exist
4. **Network Issues**: Ensure Docker daemon networking is properly configured

### Debug Commands

```bash
# Check service logs
docker-compose -f deployments/docker-compose/dev.yml logs [service-name]

# Check service health
docker-compose -f deployments/docker-compose/dev.yml ps

# Execute commands in running container
docker-compose -f deployments/docker-compose/dev.yml exec [service-name] sh
```

## Integration with Makefile

All configurations are integrated with the project Makefile:

```bash
make dev-up           # Uses dev.yml
make infra-up         # Uses infra-dev.yml
make prod-api-up      # Uses api.yml
make prod-worker-up   # Uses worker.yml
```

See the main Makefile for complete list of available commands.
