# Use a specific version of the Go image for reproducibility
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application using the full package path.
# This can resolve module path issues in certain edge cases.
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build \
    -o /app/pelico \
    -ldflags="-X pelico/internal/version.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X pelico/internal/version.GitCommit=$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" \
    pelico/cmd/server

# --- Final Stage ---
# Use a minimal base image for the final container
FROM alpine:latest

# Install necessary certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user and group for security
RUN addgroup -g 1001 -S pelico && \
    adduser -u 1001 -S pelico -G pelico

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/pelico .

# Copy the web assets required by the backend
COPY --from=builder /app/web ./web

# Create and set permissions for the data directory
RUN mkdir -p /data/roms && chown -R pelico:pelico /data

# Switch to the non-root user
USER pelico

# Expose the application port
EXPOSE 8080

# Set a health check to monitor the application's status
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# The command to run when the container starts
CMD ["./pelico"]