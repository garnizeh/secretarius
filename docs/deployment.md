# EngLog Deployment Guide

*"Deploying software is like assembling IKEA furniture - the instructions look simple until you actually try to follow them!" 🚀*

## Overview

This guide provides step-by-step instructions for deploying EngLog's two-machine architecture:
- **Machine 1**: API Server, Database, Cache, and Monitoring
- **Machine 2**: Worker Server and Background Processing

## Prerequisites

### System Requirements

#### Machine 1 (API Server)
- **OS**: Ubuntu 22.04 LTS or similar
- **RAM**: 4GB minimum, 8GB recommended
- **CPU**: 2 cores minimum, 4 cores recommended
- **Storage**: 50GB SSD minimum
- **Network**: Static IP address recommended

#### Machine 2 (Worker Server)
- **OS**: Ubuntu 22.04 LTS or similar
- **RAM**: 2GB minimum, 4GB recommended
- **CPU**: 2 cores minimum
- **Storage**: 20GB SSD minimum
- **Network**: Can reach Machine 1

### Software Requirements

Both machines need:
- Docker Engine 24.0+
- Docker Compose 2.0+
- Git
- curl/wget

### Domain Setup
- Valid domain name pointing to Machine 1's IP
- DNS A record: `yourdomain.com` → Machine 1 IP
- Optional: CNAME record: `www.yourdomain.com` → `yourdomain.com`

## Quick Start with Makefile

The project includes a comprehensive Makefile with commands for easy deployment and management:

### Development Commands
```bash
make dev-up          # Start development environment
make dev-down        # Stop development environment
make dev-logs        # View development logs
make dev-restart     # Restart development environment
```

### Production Commands

**Machine 1 (API Server):**
```bash
make deploy-machine1   # Deploy using script
make prod-api-up      # Start API services
make prod-api-down    # Stop API services
make prod-api-logs    # View API logs
```

**Machine 2 (Worker Server):**
```bash
make deploy-machine2    # Deploy using script
make prod-worker-up     # Start worker services
make prod-worker-down   # Stop worker services
make prod-worker-logs   # View worker logs
```

### Development Tools
```bash
make build            # Build all binaries
make test             # Run tests
make lint             # Run linters
make watch-api        # Run API with hot reload
make watch-worker     # Run worker with hot reload
make install-tools    # Install development tools
```

### Other Useful Commands
```bash
make help             # Show all available commands
make clean            # Clean build artifacts
make generate         # Generate code (sqlc, proto, swagger)
make docker-build     # Build Docker images
```

## Installation Steps

### Step 1: Prepare Both Machines

#### 1.1 Install Docker on Both Machines

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install dependencies
sudo apt install -y apt-transport-https ca-certificates curl gnupg lsb-release

# Add Docker's official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Add Docker repository
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Add user to docker group
sudo usermod -aG docker $USER

# Start and enable Docker
sudo systemctl start docker
sudo systemctl enable docker

# Verify installation
docker --version
docker compose version
```

#### 1.2 Clone Repository on Both Machines

```bash
# Clone the repository
git clone https://github.com/garnizeh/englog.git
cd englog

# Make scripts executable
chmod +x scripts/*.sh
```

### Step 2: Configure Environment Variables

#### 2.1 Copy and Edit Production Environment File

```bash
# Copy environment template
cp .env.production .env

# Edit with your values
nano .env
```

#### 2.2 Required Environment Variables

Replace the following values in `.env`:

```bash
# Domain and TLS - REQUIRED
DOMAIN_NAME=yourdomain.com
ACME_EMAIL=admin@yourdomain.com

# Database - GENERATE SECURE PASSWORDS
DB_PASSWORD=your-very-secure-database-password-here

# Redis - GENERATE SECURE PASSWORD
REDIS_PASSWORD=your-very-secure-redis-password-here

# JWT - GENERATE 32+ CHARACTER SECRET
JWT_SECRET_KEY=your-very-secure-jwt-secret-key-minimum-32-characters

# Machine IPs - UPDATE WITH YOUR ACTUAL IPs
POSTGRES_HOST=10.0.1.10  # Machine 1 IP
REDIS_HOST=10.0.1.10     # Machine 1 IP
API_SERVER_GRPC_ADDRESS=10.0.1.10:9090  # Machine 1 IP
WORKER_GRPC_ADDRESS=10.0.1.20:9091      # Machine 2 IP

# External Services - UPDATE WITH YOUR OLLAMA SERVER
OLLAMA_URL=http://your-ollama-server:11434

# Monitoring - GENERATE SECURE PASSWORD
GRAFANA_PASSWORD=your-secure-grafana-password

# Email (Optional) - FOR NOTIFICATIONS
SMTP_HOST=smtp.your-provider.com
SMTP_PORT=587
SMTP_USERNAME=your-smtp-username
SMTP_PASSWORD=your-smtp-password
SMTP_FROM_EMAIL=noreply@yourdomain.com
```

#### 2.3 Generate Secure Passwords

```bash
# Generate secure passwords (run these commands)
echo "DB_PASSWORD=$(openssl rand -base64 32)"
echo "REDIS_PASSWORD=$(openssl rand -base64 32)"
echo "JWT_SECRET_KEY=$(openssl rand -base64 48)"
echo "GRAFANA_PASSWORD=$(openssl rand -base64 16)"
```

### Step 3: Setup External Dependencies

#### 3.1 Setup Ollama Server (Can be on either machine or separate server)

```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Start Ollama service
sudo systemctl start ollama
sudo systemctl enable ollama

# Pull required models
ollama pull llama2
ollama pull codellama

# Verify installation
curl http://localhost:11434/api/tags
```

### Step 4: Deploy Machine 1 (API Server)

#### 4.1 Verify Configuration

```bash
# Ensure you're in the project root
cd /path/to/englog

# Verify environment file exists
ls -la .env

# Verify Caddy configuration exists
ls -la deployments/caddy/Caddyfile

# Create log directories
mkdir -p logs/{api,postgres,redis,caddy}
```

### Step 4.2 Build and Start Services

```bash
# Method 1: Using deployment script (Recommended)
./scripts/deploy-machine1.sh

# Method 2: Using Makefile (Easy)
make deploy-machine1

# Method 3: Manual deployment with Makefile
make prod-api-up

# Method 4: Manual deployment with docker-compose
docker compose -f docker-compose.api.yml pull
docker compose -f docker-compose.api.yml up -d

# Monitor startup
make prod-api-logs
# OR
docker compose -f docker-compose.api.yml logs -f
```

#### 4.3 Verify Machine 1 Deployment

```bash
# Check all services are running
docker compose -f docker-compose.api.yml ps

# Check health endpoints
curl -f http://localhost:8080/health
curl -f http://localhost/health  # Through Caddy

# Check database connection
docker compose -f docker-compose.api.yml exec postgres psql -U englog -d englog -c "SELECT version();"

# Check Redis connection
docker compose -f docker-compose.api.yml exec redis redis-cli ping

# View logs
docker compose -f docker-compose.api.yml logs api-server
docker compose -f docker-compose.api.yml logs caddy
```

### Step 5: Deploy Machine 2 (Worker Server)

#### 5.1 Copy Environment File to Machine 2

```bash
# On Machine 1, copy the .env file to Machine 2
scp .env user@machine2-ip:/path/to/englog/.env
```

#### 5.2 Verify Network Connectivity

```bash
# On Machine 2, test connection to Machine 1
curl -f http://MACHINE1_IP:8080/health

# Test gRPC connectivity (if grpcurl is installed)
grpcurl -plaintext MACHINE1_IP:9090 list
```

#### 5.3 Build and Start Worker Services

```bash
# On Machine 2
cd /path/to/englog

# Method 1: Using deployment script (Recommended)
./scripts/deploy-machine2.sh

# Method 2: Using Makefile (Easy)
make deploy-machine2

# Method 3: Manual deployment with Makefile
make prod-worker-up

# Method 4: Manual deployment with docker-compose
docker compose -f docker-compose.worker.yml pull
docker compose -f docker-compose.worker.yml up -d

# Monitor startup
make prod-worker-logs
# OR
docker compose -f docker-compose.worker.yml logs -f
```

#### 5.4 Verify Machine 2 Deployment

```bash
# Check services are running
docker compose -f docker-compose.worker.yml ps

# Check worker health
curl -f http://localhost:9091/health

# View worker logs
docker compose -f docker-compose.worker.yml logs worker-server
docker compose -f docker-compose.worker.yml logs scheduler
```

### Step 6: Final Verification

#### 6.1 Test Full System Integration

```bash
# Test API endpoints through Caddy (from external)
curl https://yourdomain.com/health
curl https://yourdomain.com/v1/auth/health

# Test gRPC communication between machines
# This should be done from application logs, look for:
# - API server connecting to worker
# - Worker server receiving jobs
# - Database operations completing successfully
```

#### 6.2 Monitor System Health

```bash
# Check Caddy TLS certificate status
docker compose -f docker-compose.api.yml exec caddy caddy list-certificates

# Monitor resource usage
docker stats

# Check application logs
tail -f logs/api/*.log
tail -f logs/worker/*.log
```

#### 6.3 Access Monitoring Dashboards

- **Grafana**: `https://yourdomain.com:3000`
  - Username: `admin`
  - Password: Your `GRAFANA_PASSWORD`

- **Prometheus**: `https://yourdomain.com:9093`

## Troubleshooting

### Common Issues

#### Issue 1: Caddy Can't Obtain TLS Certificate

```bash
# Check DNS resolution
nslookup yourdomain.com

# Check Caddy logs
docker compose -f docker-compose.api.yml logs caddy

# Verify domain points to correct IP
dig yourdomain.com A

# Common solutions:
# 1. Ensure ports 80 and 443 are open
# 2. Verify DNS propagation (can take up to 24h)
# 3. Check domain ownership
```

#### Issue 2: Database Connection Failed

```bash
# Check PostgreSQL status
docker compose -f docker-compose.api.yml exec postgres pg_isready -U englog

# Check database logs
docker compose -f docker-compose.api.yml logs postgres

# Reset database (CAUTION: This deletes all data)
docker compose -f docker-compose.api.yml down -v
docker compose -f docker-compose.api.yml up -d postgres
```

#### Issue 3: Worker Can't Connect to API Server

```bash
# Check network connectivity
ping MACHINE1_IP
telnet MACHINE1_IP 9090

# Check firewall rules
sudo ufw status
sudo iptables -L

# Check API server gRPC endpoint
docker compose -f docker-compose.api.yml logs api-server | grep grpc
```

#### Issue 4: High Resource Usage

```bash
# Monitor resource usage
docker stats
htop

# Check disk usage
df -h
docker system df

# Clean up unused Docker resources
docker system prune -a
```

### Logs Location

```bash
# Application logs
logs/api/          # API server logs
logs/worker/       # Worker server logs
logs/caddy/        # Caddy access and error logs
logs/postgres/     # Database logs
logs/redis/        # Redis logs

# Docker logs
docker compose logs [service-name]
```

## Maintenance

### Regular Tasks

#### Daily
- Check service health: `make prod-api-logs` or `make prod-worker-logs`
- Monitor disk usage: `df -h`
- Review error logs: `tail -f logs/*/error.log`
- Check containers: `docker compose ps`

#### Weekly
- Update Docker images:
  ```bash
  # For API Server (Machine 1)
  docker compose -f docker-compose.api.yml pull && make prod-api-up

  # For Worker Server (Machine 2)
  docker compose -f docker-compose.worker.yml pull && make prod-worker-up
  ```
- Clean unused Docker resources: `docker system prune -f`
- Backup database: `pg_dump` commands in backup scripts

#### Monthly
- Review and rotate logs
- Update system packages: `sudo apt update && sudo apt upgrade`
- Review security settings and certificates
- Run security checks: `make security` (if in development machine)

### Backup Strategy

```bash
# Database backup
docker compose -f docker-compose.api.yml exec postgres pg_dump -U englog englog > backup_$(date +%Y%m%d).sql

# Redis backup
docker compose -f docker-compose.api.yml exec redis redis-cli --rdb backup_redis_$(date +%Y%m%d).rdb

# Volume backup
docker run --rm -v englog_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup_$(date +%Y%m%d).tar.gz /data
```

## Security Checklist

- [ ] All default passwords changed
- [ ] Firewall configured (UFW or iptables)
- [ ] SSH key-based authentication enabled
- [ ] Regular security updates applied
- [ ] TLS certificates auto-renewing
- [ ] Database access restricted to application only
- [ ] Redis password-protected
- [ ] Monitoring alerts configured
- [ ] Backup strategy implemented
- [ ] Rate limiting configured in Caddy

## Performance Optimization

### Machine 1 Optimizations
```bash
# PostgreSQL tuning in docker-compose.api.yml
# Add to postgres service environment:
# - POSTGRES_SHARED_BUFFERS=256MB
# - POSTGRES_EFFECTIVE_CACHE_SIZE=1GB
# - POSTGRES_WORK_MEM=64MB

# Redis optimization
# Add to redis command:
# --maxmemory 512mb --maxmemory-policy allkeys-lru
```

### Machine 2 Optimizations
```bash
# Worker concurrency tuning
# Adjust WORKER_CONCURRENCY based on CPU cores
# 2-4 workers per CPU core is typically optimal
```

## Scaling Considerations

### Horizontal Scaling
- Add more worker machines using the same `docker-compose.worker.yml`
- Configure load balancing for API servers
- Consider database read replicas for high read loads

### Vertical Scaling
- Increase machine resources (CPU, RAM, Storage)
- Adjust Docker resource limits
- Tune database and cache settings

## Makefile Reference

The project includes a comprehensive Makefile that simplifies deployment and development tasks:

### Complete Command List

```bash
# Development
make dev-up              # Start development environment
make dev-down            # Stop development environment
make dev-logs            # View development logs
make dev-restart         # Restart development environment

# Production Deployment
make deploy-machine1     # Deploy Machine 1 (API Server)
make deploy-machine2     # Deploy Machine 2 (Worker Server)

# Production Management - Machine 1
make prod-api-up         # Start API services
make prod-api-down       # Stop API services
make prod-api-logs       # View API logs

# Production Management - Machine 2
make prod-worker-up      # Start worker services
make prod-worker-down    # Stop worker services
make prod-worker-logs    # View worker logs

# Building and Testing
make build               # Build all binaries
make build-api           # Build API server binary
make build-worker        # Build worker server binary
make test                # Run all tests
make test-coverage       # Run tests with coverage
make lint                # Run linting tools
make format              # Format Go code

# Development Tools
make watch-api           # Run API with hot reload
make watch-worker        # Run worker with hot reload
make run-api             # Run API server locally
make run-worker          # Run worker server locally

# Code Generation
make generate            # Generate all code
make sqlc                # Generate database code
make proto               # Generate gRPC code
make swagger             # Generate API documentation

# Database
make migrate-create NAME=migration_name  # Create new migration
make migrate-up          # Run pending migrations
make migrate-down        # Rollback last migration

# Docker
make docker-build        # Build Docker images
make docker-push         # Push images to registry

# Utilities
make deps                # Download dependencies
make install-tools       # Install development tools
make clean               # Clean build artifacts
make security            # Run security checks
make help                # Show all commands
```

### Examples

**Quick development setup:**
```bash
make install-tools       # Install required tools
make dev-up             # Start development environment
make watch-api          # Run API with hot reload (in another terminal)
```

**Production deployment:**
```bash
# On Machine 1
make deploy-machine1

# On Machine 2
make deploy-machine2
```

**Troubleshooting:**
```bash
make prod-api-logs      # Check API logs
make prod-worker-logs   # Check worker logs
make dev-down && make dev-up  # Restart development
```

## Support

For issues and questions:
- Check logs first: `docker compose logs [service]`
- Review troubleshooting section above
- Check GitHub issues: https://github.com/garnizeh/englog/issues
- Documentation: See `docs/` directory

---

**Note**: This deployment guide assumes a production environment. For development, use `docker-compose.dev.yml` instead and follow the development setup instructions in the main README.md.
