# Docker Compose Setup for EngLog

This directory contains Docker Compose configurations for the EngLog two-machine architecture.

## Files Overview

- `docker-compose/api.yml` - Machine 1: API Server, Database, Cache, Monitoring
- `docker-compose/worker.yml` - Machine 2: Worker Server, Scheduler
- `docker-compose/dev.yml` - Development environment (single machine)
- `docker-compose/infra-dev.yml` - Infrastructure only for development
- `docker-compose/api-dev.yml` - API development environment
- `docker-compose/worker-dev.yml` - Worker development environment
- `docker-compose/test.yml` - Testing environment

## Quick Start

### Development Environment

```bash
# Start development environment
docker compose -f docker-compose/dev.yml up -d

# View logs
docker compose -f docker-compose/dev.yml logs -f

# Stop services
docker compose -f docker-compose/dev.yml down
```

### Production Deployment

#### Machine 1 (API Server)
```bash
# Copy environment file
cp .env.production .env

# Edit .env with your production values
vim .env

# Deploy Machine 1
./scripts/deploy-machine1.sh

# Or manually:
docker compose -f docker-compose/api.yml up -d
```

#### Machine 2 (Worker Server)
```bash
# Copy environment file
cp .env.production .env

# Deploy Machine 2
./scripts/deploy-machine2.sh

# Or manually:
docker compose -f docker-compose/worker.yml up -d
```

## Services

### Machine 1 Services
- **Caddy**: Reverse proxy with automatic TLS (ports 80, 443)
- **API Server**: Main application server (internal port 8080)
- **PostgreSQL**: Primary database (port 5432)
- **Redis**: Cache and session store (port 6379)
- **Prometheus**: Metrics collection (port 9093)
- **Grafana**: Monitoring dashboard (port 3000)

### Machine 2 Services
- **Worker Server**: Background job processing (port 9091)
- **Scheduler**: Cron-like job scheduling

## Environment Variables

Key environment variables for production:

```bash
# Domain and TLS
DOMAIN_NAME=yourdomain.com
ACME_EMAIL=admin@yourdomain.com

# Database
DB_PASSWORD=your-secure-password
POSTGRES_HOST=10.0.1.10

# Redis
REDIS_PASSWORD=your-redis-password
REDIS_HOST=10.0.1.10

# JWT
JWT_SECRET_KEY=your-32-character-secret

# Machine Communication
API_SERVER_GRPC_ADDRESS=10.0.1.10:9090
WORKER_GRPC_ADDRESS=10.0.1.20:9091

# External Services
OLLAMA_URL=http://your-ollama-server:11434
```

## Health Checks

All services include health checks. Monitor service health:

```bash
# Check service status
docker compose -f docker-compose/api.yml ps

# View service logs
docker compose -f docker-compose/api.yml logs api-server

# Check health endpoints
curl http://localhost/health
curl http://localhost:9091/health
```

## Volumes

### Persistent Data
- `postgres_data`: PostgreSQL database files
- `redis_data`: Redis persistence
- `caddy_data`: Caddy TLS certificates
- `prometheus_data`: Metrics data
- `grafana_data`: Grafana configuration

### Log Files
- `./logs/api`: API server logs
- `./logs/worker`: Worker server logs
- `./logs/caddy`: Caddy access/error logs
- `./logs/postgres`: PostgreSQL logs
- `./logs/redis`: Redis logs

## Networking

### Networks
- `public`: Internet-facing services (Caddy, Grafana, Prometheus)
- `private`: Internal services (Database, Redis, inter-service communication)
- `dev`: Development environment network

### Security
- Database and Redis are not exposed externally in production
- Caddy handles TLS termination
- Rate limiting configured for API endpoints
- Internal services communicate over private network

## Monitoring

### Prometheus Metrics
- API server metrics: `http://localhost:9093`
- Custom application metrics
- Infrastructure metrics

### Grafana Dashboards
- Grafana UI: `http://localhost:3000`
- Default credentials: admin/admin (change in production)
- Pre-configured Prometheus datasource

## Development Features

- Hot reloading with Air
- Database seeded with development data
- All ports exposed for debugging
- Volume mounts for live code editing

## Production Considerations

1. **TLS Certificates**: Caddy automatically obtains Let's Encrypt certificates
2. **Database Backups**: Configure regular PostgreSQL backups
3. **Log Rotation**: Logs are rotated automatically
4. **Resource Limits**: Set appropriate CPU/memory limits
5. **Secrets Management**: Use Docker secrets or external secret management
6. **Monitoring**: Set up alerting based on Prometheus metrics

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 80, 443, 5432, 6379 are available
2. **Permission issues**: Check log directory permissions
3. **Network connectivity**: Verify machine-to-machine communication
4. **TLS issues**: Check domain DNS configuration

### Logs
```bash
# View all service logs
docker compose logs -f

# View specific service logs
docker compose logs -f api-server

# Check health
docker compose exec api-server wget -q --spider http://localhost:8080/health
```

### Database Access
```bash
# Connect to PostgreSQL
docker compose exec postgres psql -U englog -d englog

# Connect to Redis
docker-compose exec redis redis-cli
```
