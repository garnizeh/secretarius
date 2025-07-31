# Build stage
FROM golang:1.24.5-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the worker binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker ./cmd/worker

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget

# Create a non-root user
RUN addgroup -g 1000 worker && \
    adduser -D -s /bin/sh -u 1000 -G worker worker

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/worker .

# Change ownership
RUN chown worker:worker /app/worker

# Switch to non-root user
USER worker

# Expose health check port
EXPOSE 8091

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8091/health || exit 1

# Run the worker
CMD ["./worker"]
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the worker binary
COPY --from=builder /app/worker-server /worker-server

# Expose gRPC port for worker
EXPOSE 9090

# Run the worker binary
ENTRYPOINT ["/worker-server"]
