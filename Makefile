# Pelico Makefile

# Deployment variables - CONFIGURE THESE VIA ENVIRONMENT VARIABLES
SERVER_HOST ?= 192.168.1.52
SERVER_USER ?= bartosz
SERVER_PASSWORD ?= $(shell echo "$$SERVER_PASSWORD")
PROJECT_DIR = pelico

.PHONY: help build run test clean docker-build docker-up docker-down deps lint deploy homelab-status homelab-logs

# Default target
help:
	@echo "Available commands:"
	@echo "  build          - Build the Go binary"
	@echo "  run            - Run the application locally"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-up      - Start Docker Compose services"
	@echo "  docker-down    - Stop Docker Compose services"
	@echo "  deps           - Download Go dependencies"
	@echo "  lint           - Run linter (requires golangci-lint)"
	@echo ""
	@echo "Homelab deployment:"
	@echo "  deploy         - Build and deploy to homelab server"
	@echo "  homelab-status - Show container status on homelab"
	@echo "  homelab-logs   - Show application logs from homelab"

# Build the Go binary
build:
	@echo "Building Pelico..."
	go build -o bin/pelico cmd/server/main.go

# Run the application locally
run:
	@echo "Starting Pelico..."
	go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	docker system prune -f

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t pelico:latest .

docker-up:
	@echo "Starting Docker Compose services..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker Compose services..."
	docker-compose down

docker-logs:
	@echo "Following Docker logs..."
	docker-compose logs -f

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	cp .env.example .env
	@echo "Please edit .env with your configuration"

# Database operations
db-reset:
	@echo "Resetting database..."
	docker-compose down postgres
	docker volume rm pelico_postgres_data
	docker-compose up -d postgres

# Full reset (nuclear option)
reset-all: clean docker-down
	@echo "Performing full reset..."
	docker system prune -af
	docker volume prune -f

# Homelab deployment commands
deploy: docker-build
	@echo "Deploying to homelab server..."
	@echo "Stopping containers on server..."
	sshpass -p '$(SERVER_PASSWORD)' ssh $(SERVER_USER)@$(SERVER_HOST) 'cd $(PROJECT_DIR) && docker compose down'
	@echo "Transferring Docker image..."
	docker save pelico:latest | sshpass -p '$(SERVER_PASSWORD)' ssh $(SERVER_USER)@$(SERVER_HOST) 'docker load'
	@echo "Starting containers on server..."
	sshpass -p '$(SERVER_PASSWORD)' ssh $(SERVER_USER)@$(SERVER_HOST) 'cd $(PROJECT_DIR) && docker compose up -d'
	@echo "✅ Deployment complete! Application available at: http://$(SERVER_HOST):8081"

homelab-status:
	@echo "Checking homelab container status..."
	sshpass -p '$(SERVER_PASSWORD)' ssh $(SERVER_USER)@$(SERVER_HOST) 'docker ps --filter name=pelico'

homelab-logs:
	@echo "Following homelab application logs..."
	sshpass -p '$(SERVER_PASSWORD)' ssh $(SERVER_USER)@$(SERVER_HOST) 'cd $(PROJECT_DIR) && docker compose logs -f pelico'