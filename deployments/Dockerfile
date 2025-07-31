# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates for HTTPS requests
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod and sum files for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o api-server ./cmd/api

# Build worker as well
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o worker-server ./cmd/worker

# Development stage for testing and development
FROM golang:1.24-alpine AS development

# Install development tools and runtime dependencies
RUN apk add --no-cache git ca-certificates tzdata curl

WORKDIR /app

# Copy go mod and sum files for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build applications in development mode (with debug info)
RUN go build -o ./bin/api ./cmd/api
RUN go build -o ./bin/worker ./cmd/worker

# Expose ports for API and Worker
EXPOSE 8080 9090 9091

# Default command (can be overridden in docker-compose)
CMD ["./bin/api"]

# Final stage for API server
FROM scratch AS api

# Copy ca-certificates from builder for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the API binary
COPY --from=builder /app/api-server /api-server

# Expose API port
EXPOSE 8080

# Health check for API server
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/api-server", "--health-check"]

# Run the API binary
ENTRYPOINT ["/api-server"]

# Final stage for Worker server
FROM scratch AS worker

# Copy ca-certificates from builder for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the worker binary
COPY --from=builder /app/worker-server /worker-server

# Expose gRPC port for worker
EXPOSE 9090

# Run the worker binary
ENTRYPOINT ["/worker-server"]
