#!/bin/bash
set -e

echo "Deploying EngLog API Server (Machine 1)..."

# Load environment variables
source .env.production

# Create necessary directories
mkdir -p logs/api logs/postgres logs/redis logs/caddy
mkdir -p deployments/caddy

# Set permissions for logs
chmod 755 logs/*

# Create Caddyfile if it doesn't exist
if [ ! -f deployments/caddy/Caddyfile ]; then
    echo "Warning: Caddyfile not found. Please create deployments/caddy/Caddyfile"
    exit 1
fi

# Pull latest images
docker compose -f docker-compose.api.yml pull

# Start services
docker compose -f docker-compose.api.yml up -d

# Wait for services to be healthy
echo "Waiting for services to be healthy..."
sleep 30

# Check health
docker compose -f docker-compose.api.yml exec api-server wget --quiet --tries=1 --spider http://localhost:8080/health

echo "Machine 1 deployment completed successfully!"
echo "Caddy will automatically obtain TLS certificates for domain: $DOMAIN_NAME"
