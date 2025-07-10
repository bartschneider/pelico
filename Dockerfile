FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with version information
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-X pelico/internal/version.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
              -X pelico/internal/version.GitCommit=$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" \
    -o pelico ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for SSL/TLS and Docker CLI for logs functionality
RUN apk --no-cache add ca-certificates tzdata docker-cli

# Create app user (staying as root for Docker socket access)
RUN addgroup -g 1001 -S pelico && \
    adduser -u 1001 -S pelico -G pelico

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/pelico .

# Copy web assets
COPY --from=builder /app/web ./web

# Create data directory for ROM scanning
RUN mkdir -p /data/roms && chown -R pelico:pelico /data

# Stay as root for Docker socket access
# USER pelico

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Run the application
CMD ["./pelico"]