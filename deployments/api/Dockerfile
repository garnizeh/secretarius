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
