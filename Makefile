# EngLog Makefile
# "Make it so!" - Jean-Luc Picard 🚀

.PHONY: help build clean test lint dev-up dev-down run-api run-worker generate migrate-up migrate-down migrate-status migrate-reset migrate-create sqlc proto swagger docker-build docker-push

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

## test: Run all tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

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
	@mkdir -p pkg/proto
	@protoc --go_out=pkg/proto --go-grpc_out=pkg/proto proto/*.proto

## swagger: Generate Swagger documentation
swagger:
	@echo "Generating Swagger docs..."
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@latest
	@swag init -g cmd/api/main.go -o docs/

## migrate-create: Create a new migration file (usage: make migrate-create NAME=migration_name)
migrate-create:
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=migration_name"; exit 1; fi
	@which goose > /dev/null || (echo "Installing goose..." && go install github.com/pressly/goose/v3/cmd/goose@latest)
	@goose -dir internal/sqlc/schema create $(NAME) sql

## dev-up: Start development environment with Docker Compose
dev-up:
	@echo "Starting development environment..."
	@docker compose -f deployments/docker-compose.dev.yml up -d

## dev-down: Stop development environment
dev-down:
	@echo "Stopping development environment..."
	@docker compose -f deployments/docker-compose.dev.yml down

## dev-logs: View development environment logs
dev-logs:
	@docker compose -f deployments/docker-compose.dev.yml logs -f

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
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	@air -c .air.api.toml

## watch-worker: Run worker server with live reload
watch-worker:
	@echo "Starting worker server with live reload..."
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	@air -c .air.worker.toml

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
	@go install github.com/cosmtrek/air@latest
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
