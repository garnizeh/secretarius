#!/bin/bash
set -e

echo "Deploying EngLog Worker Server (Machine 2)..."

# Load environment variables
if [ -f ".env" ]; then
    source .env
elif [ -f "deployments/environments/production/.env.worker" ]; then
    source deployments/environments/production/.env.worker
else
    echo "Error: No environment file found. Please create .env or use 'make env-prod-worker'"
    exit 1
fi

# Create necessary directories
mkdir -p logs/worker logs/scheduler

# Verify Ollama server connectivity
echo "Checking Ollama server connectivity..."
if ! curl -f "${OLLAMA_URL:-http://localhost:11434}/api/tags" >/dev/null 2>&1; then
    echo "Warning: Cannot connect to Ollama server at ${OLLAMA_URL:-http://localhost:11434}"
    echo "Please ensure Ollama server is running and accessible"
fi

# Pull latest images
docker compose -f deployments/docker-compose/worker.yml pull

# Start services
docker compose -f deployments/docker-compose/worker.yml up -d

# Wait for services to be ready
echo "Waiting for worker services to be ready..."
sleep 30

echo "Machine 2 deployment completed successfully!"
