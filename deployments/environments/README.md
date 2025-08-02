# Environment Configuration

> "Configuration is the root of all evil." - Anonymous DevOps Engineer 🔧

This directory contains environment-specific configuration files for the EngLog application. Each environment has its own subdirectory with appropriate settings.

## Directory Structure

```
environments/
├── README.md                          # This file
├── .env.example                       # Base template for all environments
├── development/
│   ├── .env.dev                       # Complete development environment
│   ├── .env.api-dev                   # API server development only
│   └── .env.worker-dev                # Worker service development only
├── testing/
│   └── .env.test                      # Automated testing environment
└── production/
    ├── .env.api.example               # API server production template
    └── .env.worker.example            # Worker service production template
```

## Environment Types

### 📋 **Development**
- **Purpose**: Local development and debugging
- **Security**: Relaxed for convenience
- **Performance**: Debug-friendly settings
- **Database**: Local PostgreSQL/SQLite
- **Files**: `.env.dev`, `.env.api-dev`, `.env.worker-dev`

### 🧪 **Testing**
- **Purpose**: Automated testing and CI/CD
- **Security**: Minimal (test data only)
- **Performance**: Fast startup and execution
- **Database**: In-memory or test containers
- **Files**: `.env.test`

### 🚀 **Production**
- **Purpose**: Live production deployment
- **Security**: Maximum security settings
- **Performance**: Optimized for scale
- **Database**: Production PostgreSQL cluster
- **Files**: `.env.api.example`, `.env.worker.example`

## Usage Guide

### Development Setup

1. **Complete Development Environment**
   ```bash
   # Copy and configure main development environment
   cp deployments/environments/development/.env.dev .env

   # Edit with your local settings
   vim .env

   # Start full environment
   make dev-up
   ```

2. **API Development Only**
   ```bash
   # Copy API-specific environment
   cp deployments/environments/development/.env.api-dev .env

   # Start infrastructure + API
   make infra-up
   make dev-api-up
   ```

3. **Worker Development Only**
   ```bash
   # Copy Worker-specific environment
   cp deployments/environments/development/.env.worker-dev .env

   # Start infrastructure + Worker
   make infra-up
   make dev-worker-up
   ```

### Testing Setup

```bash
# Copy test environment
cp deployments/environments/testing/.env.test .env.test

# Run tests with test environment
make test-docker-up
make test
```

### Production Deployment

#### Machine 1 (API Server)
```bash
# Copy and customize API production template
cp deployments/environments/production/.env.api.example .env

# Edit with production values
vim .env

# Deploy API server
make prod-api-up
```

#### Machine 2 (Worker Server)
```bash
# Copy and customize Worker production template
cp deployments/environments/production/.env.worker.example .env

# Edit with production values
vim .env

# Deploy Worker server
make prod-worker-up
```

## Configuration Categories

### 🔐 **Security Settings**
- **JWT Secrets**: Strong, unique keys for production
- **Database Passwords**: Complex passwords with rotation
- **TLS Certificates**: Proper certificate management
- **CORS Origins**: Restricted to known domains
- **Rate Limiting**: Protection against abuse

### 📊 **Performance Settings**
- **Connection Pooling**: Optimized for environment load
- **Cache Configuration**: Memory and TTL settings
- **Worker Concurrency**: Based on available resources
- **Timeouts**: Appropriate for environment characteristics

### 🔍 **Monitoring & Logging**
- **Log Levels**: Debug for dev, info for production
- **Metrics Collection**: Enabled for production monitoring
- **Health Checks**: Environment-appropriate intervals
- **Error Tracking**: Structured logging configuration

## Environment Variables Reference

### Core Application
| Variable | Development | Testing | Production | Description |
|----------|-------------|---------|------------|-------------|
| `APP_ENV` | development | test | production | Application environment |
| `APP_PORT` | 8080 | 8080 | 8080 | HTTP server port |
| `LOG_LEVEL` | debug | debug | info | Logging verbosity |
| `TLS_ENABLED` | false | false | true | Enable HTTPS/TLS |

### Database
| Variable | Development | Testing | Production | Description |
|----------|-------------|---------|------------|-------------|
| `DB_HOST_READ_WRITE` | localhost:5432 | localhost:5433 | postgres:5432 | Primary database host |
| `DB_MAX_OPEN_CONNS` | 25 | 10 | 100 | Maximum open connections |
| `DB_CONN_MAX_LIFETIME` | 5m | 5m | 1h | Connection lifetime |

### Security
| Variable | Development | Testing | Production | Description |
|----------|-------------|---------|------------|-------------|
| `JWT_SECRET` | dev_secret | test_secret | SECURE_KEY | JWT signing key |
| `BCRYPT_COST` | 4 | 4 | 12 | Password hashing cost |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | 1000 | 1000 | 60 | API rate limiting |

### gRPC
| Variable | Development | Testing | Production | Description |
|----------|-------------|---------|------------|-------------|
| `GRPC_TLS_ENABLED` | false | false | true | Enable gRPC TLS |
| `GRPC_SERVER_PORT` | 50051 | 50052 | 50051 | gRPC server port |

## Security Best Practices

### 🔒 **Production Security Checklist**
- [ ] Strong, unique JWT secrets (32+ characters)
- [ ] Complex database passwords with special characters
- [ ] TLS enabled for all external communication
- [ ] Restricted CORS origins (no wildcards)
- [ ] Rate limiting enabled and properly configured
- [ ] Proper certificate management and rotation
- [ ] Environment variables injected securely (not in files)

### 🛡️ **Secret Management**
- **Development**: Simple passwords, shared secrets OK
- **Testing**: Temporary secrets, reset frequently
- **Production**: Use secret management systems (Vault, K8s secrets, etc.)

### 📝 **Configuration Validation**
Each environment file includes validation comments and examples. Always verify:
- Database connectivity
- External service availability
- Certificate validity
- Network accessibility
- Resource limits

## Docker Compose Integration

Environment files are automatically loaded by Docker Compose configurations:

```yaml
# Example docker-compose reference
services:
  api-server:
    env_file:
      - deployments/environments/production/.env.api.example
```

## Troubleshooting

### Common Issues

1. **Environment File Not Found**
   ```bash
   # Verify file exists and path is correct
   ls -la deployments/environments/development/.env.dev
   ```

2. **Invalid Environment Variables**
   ```bash
   # Validate environment loading
   docker-compose config
   ```

3. **Permission Issues**
   ```bash
   # Fix file permissions
   chmod 600 deployments/environments/production/.env.*
   ```

### Debug Commands

```bash
# Check current environment
env | grep -E "(APP_|DB_|REDIS_)"

# Validate Docker Compose with environment
docker-compose -f deployments/docker-compose/dev.yml config

# Test database connection
docker-compose exec api-server sh -c 'echo $DATABASE_URL'
```

## Migration Guide

If you have existing `.env` files in the root directory:

```bash
# Backup existing files
cp .env .env.backup

# Move to appropriate environment directory
cp .env deployments/environments/development/.env.dev

# Update references in scripts and documentation
grep -r "\.env" . --exclude-dir=.git
```

---

**Remember**: Never commit actual production environment files to version control. Always use `.example` templates and inject real values during deployment.
