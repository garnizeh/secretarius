# EngLog Makefile
# "Make it so!" - Jean-Luc Picard 🚀

.PHONY: help build clean test test-unit test-integration test-e2e test-coverage test-race test-security test-performance test-clean test-docker-up test-docker-down test-fix generate-mocks lint dev-up dev-down dev-api-up dev-api-down dev-worker-up dev-worker-down dev-restart infra-up infra-down infra-logs infra-restart prod-api-up prod-api-down prod-worker-up prod-worker-down deploy-machine1 deploy-machine2 run-api run-worker health-api generate migrate-up migrate-down migrate-status migrate-reset migrate-create sqlc proto swagger docker-build docker-push git-clean-branches env-dev env-api-dev env-worker-dev env-test env-prod-api env-prod-worker env-check

# Default target
.DEFAULT_GOAL := help

# Variables
API_BINARY := bin/api
WORKER_BINARY := bin/worker
GO_FILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
DOCKER_REGISTRY := docker.io
IMAGE_NAME := englog
VERSION := $(shell git describe --tags --always --dirty)

# Environment variables for database connection
DB_USER := englog
DB_PASSWORD := englog_dev_password
DB_HOST := localhost
DB_PORT := 5432
DB_NAME := englog

## help: Show this help message
help:
	@echo "EngLog Development Commands"
	@echo "=========================="
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | sort

## build: Build all binaries
build: build-api build-worker

## build-api: Build API server binary
build-api:
	@echo "Building API server..."
	@mkdir -p bin
	@go build -ldflags="-X main.Version=$(VERSION)" -o $(API_BINARY) ./cmd/api

## build-worker: Build worker server binary
build-worker:
	@echo "Building worker server..."
	@mkdir -p bin
	@go build -ldflags="-X main.Version=$(VERSION)" -o $(WORKER_BINARY) ./cmd/worker

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf tmp/
	@go clean -cache
	@go clean -testcache

## git-clean-branches: Remove local branches that have been deleted from origin
git-clean-branches:
	@echo "Cleaning up local branches removed from origin..."
	@git remote prune origin
	@git branch -vv | grep ': gone]' | awk '{print $$1}' | xargs -r git branch -D
	@echo "Local branch cleanup completed!"

## env-dev: Setup development environment configuration
env-dev:
	@echo "Setting up development environment..."
	@cp deployments/environments/development/.env.dev .env
	@echo "✅ Development environment configured! Edit .env if needed."

## env-api-dev: Setup API development environment configuration
env-api-dev:
	@echo "Setting up API development environment..."
	@cp deployments/environments/development/.env.api-dev .env
	@echo "✅ API development environment configured! Edit .env if needed."

## env-worker-dev: Setup Worker development environment configuration
env-worker-dev:
	@echo "Setting up Worker development environment..."
	@cp deployments/environments/development/.env.worker-dev .env
	@echo "✅ Worker development environment configured! Edit .env if needed."

## env-test: Setup testing environment configuration
env-test:
	@echo "Setting up testing environment..."
	@cp deployments/environments/testing/.env.test .env.test
	@echo "✅ Testing environment configured!"

## env-prod-api: Setup production API environment template
env-prod-api:
	@echo "Setting up production API environment template..."
	@cp deployments/environments/production/.env.api.example .env
	@echo "⚠️  IMPORTANT: Edit .env with your production values before deploying!"

## env-prod-worker: Setup production Worker environment template
env-prod-worker:
	@echo "Setting up production Worker environment template..."
	@cp deployments/environments/production/.env.worker.example .env
	@echo "⚠️  IMPORTANT: Edit .env with your production values before deploying!"

## env-check: Check current environment configuration
env-check:
	@echo "Current Environment Configuration:"
	@echo "=================================="
	@if [ -f .env ]; then \
		echo "✅ Environment file found: .env"; \
		echo "Environment type: $$(grep APP_ENV .env 2>/dev/null || echo 'Not specified')"; \
	else \
		echo "❌ No .env file found. Run 'make env-dev' to set up development environment."; \
	fi

## test: Run all tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

## test-all: Run comprehensive test suite (unit, integration, e2e, race, security, performance)
test-all: test-unit test-integration test-e2e test-race test-security test-performance
	@echo "All tests completed successfully! 🎉"

## test-unit: Run unit tests only (fast, no race detection)
test-unit:
	@echo "Running unit tests..."
	go test -mod=mod -v -short -tags='!integration' ./internal/...

## test-integration: Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -mod=mod -v -tags=integration ./tests/integration/... ./internal/...

## test-e2e: Run end-to-end tests
test-e2e:
	@echo "Running end-to-end tests..."
	go test -mod=mod -v -tags=e2e ./tests/e2e/...

## test-race: Run race condition tests (unit tests with race detection)
test-race:
	@echo "Running race condition tests..."
	go test -v -race -short -tags='!integration' ./internal/...

## test-security: Run security tests
test-security:
	@echo "Running security tests..."
	@which gosec > /dev/null || go install github.com/securego/gosec/v2/cmd/gosec@latest
	@which govulncheck > /dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest
	gosec -fmt sarif -out gosec.sarif ./... || true
	govulncheck ./...

## test-performance: Run performance tests
test-performance:
	@echo "Running performance tests..."
	go test -mod=mod -v -tags=performance -bench=. -benchmem ./tests/performance/...

## test-clean: Clean test artifacts
test-clean:
	rm -f coverage.out coverage.html gosec.sarif
	go clean -testcache

## test-docker-up: Start Docker test environment
test-docker-up:
	docker compose -f deployments/docker-compose/test.yml up -d --build

## test-docker-down: Stop Docker test environment
test-docker-down:
	docker compose -f deployments/docker-compose/test.yml down -v

## test-fix: Fix failing handler tests
test-fix:
	@echo "Fixing failing handler tests..."
	go test -v ./internal/handlers/... -run="TestAnalytics"


## generate-mocks: Generate mocks
generate-mocks:
	@echo "Generating mocks..."
	go generate ./...

## test-coverage: Run tests with coverage report
test-coverage: test
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## lint: Run linting tools
lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@golangci-lint run

## format: Format Go code
format:
	@echo "Formatting code..."
	@go fmt ./...
	@which goimports > /dev/null || go install golang.org/x/tools/cmd/goimports@latest
	@goimports -w $(GO_FILES)

## generate: Generate all code (sqlc, protobuf, swagger)
generate: sqlc proto swagger

## certs: Generate TLS certificates for development
certs:
	@echo "Generating TLS certificates..."
	@./scripts/generate-certs.sh

## sqlc: Generate database code with sqlc
sqlc:
	@echo "Generating database code..."
	@which sqlc > /dev/null || (echo "Installing sqlc..." && go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest)
	@sqlc generate

## proto: Generate gRPC code from protobuf
proto:
	@echo "Generating gRPC code..."
	@which protoc > /dev/null || (echo "protoc not found. Please install Protocol Buffers compiler." && exit 1)
	@which protoc-gen-go > /dev/null || go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@which protoc-gen-go-grpc > /dev/null || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@mkdir -p proto/worker
	@protoc --go_out=. --go_opt=module=github.com/garnizeh/englog --go-grpc_out=. --go-grpc_opt=module=github.com/garnizeh/englog proto/worker.proto

## swagger: Generate Swagger documentation
swagger:
	@echo "Generating Swagger docs..."
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@latest
	@swag init -g cmd/api/main.go -o api/

## migrate-create: Create a new migration file (usage: make migrate-create NAME=migration_name)
migrate-create:
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=migration_name"; exit 1; fi
	@which goose > /dev/null || (echo "Installing goose..." && go install github.com/pressly/goose/v3/cmd/goose@latest)
	@goose -dir internal/sqlc/schema create $(NAME) sql

## dev-up: Start development environment with Docker Compose
dev-up:
	@echo "Starting development environment..."
	@docker compose -f deployments/docker-compose/dev.yml up -d --build

## dev-down: Stop development environment
dev-down:
	@echo "Stopping development environment..."
	@docker compose -f deployments/docker-compose/dev.yml down

## dev-logs: View development environment logs
dev-logs:
	@docker compose -f deployments/docker-compose/dev.yml logs -f

## dev-restart: Restart development environment
dev-restart: dev-down dev-up

## dev-api-up: Start API development environment only
dev-api-up:
	@echo "Starting API development environment..."
	@docker compose -f deployments/docker-compose/api-dev.yml up -d --build

## dev-api-down: Stop API development environment
dev-api-down:
	@echo "Stopping API development environment..."
	@docker compose -f deployments/docker-compose/api-dev.yml down

## dev-worker-up: Start Worker development environment only
dev-worker-up:
	@echo "Starting Worker development environment..."
	@docker compose -f deployments/docker-compose/worker-dev.yml up -d --build

## dev-worker-down: Stop Worker development environment
dev-worker-down:
	@echo "Stopping Worker development environment..."
	@docker compose -f deployments/docker-compose/worker-dev.yml down

## infra-up: Start development environment (infrastructure only) with Docker Compose
infra-up:
	@echo "Starting development environment..."
	@docker compose -f deployments/docker-compose/infra-dev.yml up -d --build

## infra-down: Stop development environment (infrastructure only)
infra-down:
	@echo "Stopping development environment..."
	@docker compose -f deployments/docker-compose/infra-dev.yml down

## infra-logs: View development environment (infrastructure only) logs
infra-logs:
	@docker compose -f deployments/docker-compose/infra-dev.yml logs -f

## infra-restart: Restart development environment (infrastructure only)
infra-restart: infra-down infra-up

## prod-api-up: Start production API server (Machine 1)
prod-api-up:
	@echo "Starting production API server (Machine 1)..."
	@docker compose -f deployments/docker-compose/api.yml up -d --build

## prod-api-down: Stop production API server (Machine 1)
prod-api-down:
	@echo "Stopping production API server (Machine 1)..."
	@docker compose -f deployments/docker-compose/api.yml down

## prod-api-logs: View production API server logs
prod-api-logs:
	@docker compose -f deployments/docker-compose/api.yml logs -f

## prod-worker-up: Start production worker server (Machine 2)
prod-worker-up:
	@echo "Starting production worker server (Machine 2)..."
	@docker compose -f deployments/docker-compose/worker.yml up -d --build

## prod-worker-down: Stop production worker server (Machine 2)
prod-worker-down:
	@echo "Stopping production worker server (Machine 2)..."
	@docker compose -f deployments/docker-compose/worker.yml down

## prod-worker-logs: View production worker server logs
prod-worker-logs:
	@docker compose -f deployments/docker-compose/worker.yml logs -f

## deploy-machine1: Deploy Machine 1 using script
deploy-machine1:
	@echo "Deploying Machine 1 (API Server)..."
	@./scripts/deploy-machine1.sh

## deploy-machine2: Deploy Machine 2 using script
deploy-machine2:
	@echo "Deploying Machine 2 (Worker Server)..."
	@./scripts/deploy-machine2.sh

## run-api: Run API server locally
run-api:
	@echo "Starting API server..."
	@go run ./cmd/api

## run-worker: Run worker server locally
run-worker:
	@echo "Starting worker server..."
	@go run ./cmd/worker

## watch-api: Run API server with live reload
watch-api:
	@echo "Starting API server with live reload..."
	@which air > /dev/null || go install github.com/air-verse/air@latest
	@air -c .air.api.toml

## watch-worker: Run worker server with live reload
watch-worker:
	@echo "Starting worker server with live reload..."
	@which air > /dev/null || go install github.com/air-verse/air@latest
	@air -c .air.worker.toml

## debug-api: Run API server with debug flags and live reload
debug-api:
	@echo "Starting API server with debug mode and live reload..."
	@which air > /dev/null || go install github.com/air-verse/air@latest
	@air -c .air.debug.toml

## dev-api: Alias for watch-api (backward compatibility)
dev-api: watch-api

## dev-worker: Alias for watch-worker (backward compatibility)
dev-worker: watch-worker

## air-dev: Interactive Air development helper
air-dev:
	@./scripts/air-dev.sh

## air-both: Start both API and Worker with Air
air-both:
	@./scripts/air-dev.sh both## health-api: Check if API server is running and healthy
health-api:
	@echo "Checking API server health..."
	@curl -f -s http://localhost:8080/health > /tmp/health_check.json && \
		echo "✅ API server is healthy and running!" && \
		echo "Response: $$(cat /tmp/health_check.json)" && \
		rm -f /tmp/health_check.json || \
		(echo "❌ API server is not responding. Make sure it's running with 'make dev-up' and 'make watch-api'" && exit 1)

## docker-build: Build Docker images
docker-build:
	@echo "Building Docker images..."
	@docker build -f deployments/api/Dockerfile -t $(DOCKER_REGISTRY)/$(IMAGE_NAME)-api:$(VERSION) .
	@docker build -f deployments/worker/Dockerfile -t $(DOCKER_REGISTRY)/$(IMAGE_NAME)-worker:$(VERSION) .

## docker-push: Push Docker images to registry
docker-push: docker-build
	@echo "Pushing Docker images..."
	@docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME)-api:$(VERSION)
	@docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME)-worker:$(VERSION)

## deps: Download and tidy dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/air-verse/air@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

## security: Run security checks
security:
	@echo "Running security checks..."
	@which gosec > /dev/null || go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@gosec ./...

## benchmark: Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

## vendor: Create vendor directory
vendor:
	@echo "Creating vendor directory..."
	@go mod vendor

## update: Update dependencies
update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

## check: Run all checks (lint, test, security)
check: lint test security

## release: Build release binaries for multiple platforms
release:
	@echo "Building release binaries..."
	@mkdir -p bin/release
	@GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o bin/release/englog-api-linux-amd64 ./cmd/api
	@GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o bin/release/englog-worker-linux-amd64 ./cmd/worker
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o bin/release/englog-api-darwin-amd64 ./cmd/api
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o bin/release/englog-worker-darwin-amd64 ./cmd/worker
	@GOOS=windows GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o bin/release/englog-api-windows-amd64.exe ./cmd/api
	@GOOS=windows GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o bin/release/englog-worker-windows-amd64.exe ./cmd/worker

# Load environment variables from .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif
