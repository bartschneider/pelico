# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
- **Name**: Pelico - Video Game Collection Manager
- **Language**: Go 1.24
- **Type**: Web-based video game collection management system
- **Module**: `pelico`
- **Architecture**: REST API backend with web frontend
- **Database**: PostgreSQL with GORM ORM

## Development Commands

### Core Application Commands
```bash
# Run the server
go run cmd/server/main.go

# Build the server binary
go build -o pelico cmd/server/main.go

# Run with Make
make run
make build

# Docker development
make docker-up
make docker-down
make docker-logs
```

### Database Commands
```bash
# Reset database (development)
make db-reset

# Run migrations (automatic on startup)
# Migrations are handled by GORM AutoMigrate
```

### Testing Commands
```bash
# Run all tests
go test ./...

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/services/...
```

### Dependency Management
```bash
# Clean up dependencies
go mod tidy

# Download dependencies
make deps
```

## Project Structure

### Application Architecture
```
pelico/
├── cmd/server/          # Application entry point
├── internal/            # Private application packages
│   ├── api/            # HTTP server and routing
│   ├── config/         # Configuration management
│   ├── database/       # Database connection and migrations
│   ├── handlers/       # HTTP request handlers
│   ├── models/         # Database models (GORM)
│   └── services/       # Business logic services
├── web/                # Frontend assets
│   ├── templates/      # HTML templates
│   └── static/         # CSS, JS, images
├── scripts/            # Database init and utility scripts
└── docker-compose.yml  # Container orchestration
```

### Key Components
- **API Server** (`internal/api/`): Gin-based REST API
- **Database Models** (`internal/models/`): Game, Platform, FileLocation, PlaySession
- **ROM Scanner** (`internal/services/rom_scanner.go`): File system scanning
- **Metadata Service** (`internal/services/metadata_service.go`): External API integration
- **Web Frontend** (`web/`): Bootstrap-based responsive UI

## Database Schema

### Core Tables
- `platforms` - Gaming platforms (Nintendo, Sony, etc.)
- `games` - Game records with metadata
- `file_locations` - ROM file paths and server locations
- `play_sessions` - Gaming session tracking with notes

### Relationships
- Games belong to Platforms (many-to-one)
- Games have many FileLocations (one-to-many)
- Games have many PlaySessions (one-to-many)

## API Endpoints

### Games
- `GET /api/v1/games` - List all games
- `POST /api/v1/games` - Create game
- `GET /api/v1/games/:id` - Get specific game
- `PUT /api/v1/games/:id` - Update game
- `POST /api/v1/games/search` - Search games

### ROM Scanning
- `POST /api/v1/scan/directory` - Scan ROM directory
- `GET /api/v1/scan/duplicates` - Find duplicate files

### Play Sessions
- `GET /api/v1/games/:id/sessions` - Get game sessions
- `POST /api/v1/games/:id/sessions` - Create session

## Development Guidelines

### Code Organization
- Use package-based organization with clear separation of concerns
- Keep business logic in `services/` packages
- HTTP handling in `handlers/` packages
- Database operations through GORM models

### Error Handling
- Return proper HTTP status codes
- Use structured error responses
- Log errors appropriately

### Database Operations
- Use GORM for all database operations
- Leverage preloading for relationships
- Handle connection pooling automatically

### External Integrations
- TheGamesDB for metadata (no API key required)
- IGDB integration available (requires API key)
- Graceful fallback when external services fail

## Deployment

### Docker Development
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f pelico

# Stop services
docker-compose down
```

### Configuration
- Environment variables in `.env` file
- ROM paths configured in Docker volumes
- Database connection via DATABASE_URL

### Production Considerations
- Configure ROM directory mounts in docker-compose.yml
- Set up proper environment variables
- Ensure PostgreSQL persistence with volume mounts
- Configure reverse proxy for external access

## Supported File Types
- ROM files: `.rom`, `.bin`, `.iso`, `.cue`, `.img`
- Compressed: `.zip`, `.7z`, `.rar`
- Platform-specific: `.nes`, `.smc`, `.n64`, `.gba`, etc.

## External Services
- **TheGamesDB**: Free game metadata API
- **IGDB**: Comprehensive game database (requires API key)
- **File scanning**: Local filesystem integration

## Testing Strategy
- Unit tests for services and business logic
- Integration tests for database operations
- API endpoint testing with test database
- Mock external service calls for reliable testing