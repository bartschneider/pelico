# Pelico - Video Game Collection Manager

A comprehensive web-based application for managing and reviewing your video game collection in a homelab environment.

## Features

- **Multi-platform Support**: Manage games across different gaming platforms
- **ROM Collection Scanning**: Automatically scan directories for ROM files and add them to your collection
- **Metadata Integration**: Fetch game metadata, cover art, and descriptions from online sources
- **Play Session Tracking**: Log gaming sessions with notes and ratings
- **Duplicate Detection**: Find and manage duplicate games in your collection
- **Web Interface**: Modern, responsive web UI accessible on your local network
- **Self-hosted**: Runs entirely on your own infrastructure

## Technology Stack

- **Backend**: Go with Gin web framework
- **Database**: PostgreSQL
- **Frontend**: HTML/CSS/JavaScript with Bootstrap
- **Containerization**: Docker and Docker Compose
- **External APIs**: TheGamesDB, IGDB (optional)

## Quick Start

### Prerequisites

- Docker and Docker Compose
- At least 1GB RAM
- Storage space for database and cover art

### Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd pelico
   ```

2. **Set up environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your ROM paths and API keys
   ```

3. **Configure ROM paths**:
   Edit `docker-compose.yml` and `.env` to point to your ROM directories:
   ```bash
   # In .env
   ROM_PATH_1=/path/to/your/nes/roms
   ROM_PATH_2=/path/to/your/snes/roms
   ```

4. **Start the services**:
   ```bash
   docker-compose up -d
   ```

5. **Access the application**:
   Open http://localhost:8080 in your browser

## Configuration

### Environment Variables

- `PORT`: Server port (default: 8080)
- `DATABASE_URL`: PostgreSQL connection string
- `IGDB_API_KEY`: Optional IGDB API key for enhanced metadata
- `ROM_PATH_*`: Mount paths for your ROM collections

### ROM Directory Structure

Pelico supports various ROM file formats:
- `.rom`, `.bin`, `.iso`, `.cue`, `.img`
- `.zip`, `.7z`, `.rar` (compressed)
- Platform-specific: `.nes`, `.smc`, `.n64`, `.gba`, etc.

Example directory structure:
```
/data/roms/
├── nintendo/
│   ├── nes/
│   ├── snes/
│   └── n64/
├── sega/
│   ├── genesis/
│   └── saturn/
└── sony/
    ├── psx/
    └── ps2/
```

## Usage

### Adding Games

1. **Manual Entry**: Use the "Add Game" button to manually add games
2. **ROM Scanning**: Use the ROM Scanner to automatically detect games from file directories
3. **Metadata Fetching**: Click the download icon on any game to fetch metadata

### Managing Collection

- **Search**: Use the search bar to find games by title
- **Filter**: Browse by platform or genre
- **Edit**: Update game information and metadata
- **Sessions**: Track your gaming sessions with notes and ratings

### ROM Scanning

1. Navigate to the "ROM Scanner" tab
2. Enter the directory path (mounted in Docker)
3. Select the appropriate platform
4. Choose whether to scan recursively
5. Click "Start Scan"

### Finding Duplicates

1. Go to ROM Scanner tab
2. Click "Find Duplicates"
3. Review the results showing files with identical hashes
4. Manually remove unwanted duplicates

## API Documentation

The application provides a REST API accessible at `/api/v1/`:

### Games
- `GET /api/v1/games` - List all games
- `GET /api/v1/games/:id` - Get specific game
- `POST /api/v1/games` - Create new game
- `PUT /api/v1/games/:id` - Update game
- `DELETE /api/v1/games/:id` - Delete game
- `POST /api/v1/games/search` - Search games

### Platforms
- `GET /api/v1/platforms` - List platforms
- `POST /api/v1/platforms` - Create platform

### Play Sessions
- `GET /api/v1/games/:id/sessions` - Get game sessions
- `POST /api/v1/games/:id/sessions` - Create session
- `PUT /api/v1/sessions/:id` - Update session

### Scanning
- `POST /api/v1/scan/directory` - Scan ROM directory
- `GET /api/v1/scan/duplicates` - Find duplicates

## Development

### Building from Source

```bash
# Install dependencies
go mod download

# Run locally
go run cmd/server/main.go

# Build binary
go build -o pelico cmd/server/main.go
```

### Database Setup

```bash
# Start PostgreSQL
docker run -d --name pelico-postgres \
  -e POSTGRES_DB=pelico \
  -e POSTGRES_USER=pelico \
  -e POSTGRES_PASSWORD=pelico \
  -p 5432:5432 postgres:15-alpine

# Set DATABASE_URL
export DATABASE_URL="postgres://pelico:pelico@localhost:5432/pelico?sslmode=disable"
```

### Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Supported Platforms

The application comes pre-configured with common gaming platforms:

- Nintendo: NES, SNES, N64, GameCube, Wii, Game Boy family, DS family
- Sega: Master System, Genesis, Saturn, Dreamcast
- Sony: PlayStation 1-5, PSP, PS Vita
- Microsoft: Xbox, Xbox 360, Xbox One, Xbox Series X/S
- Atari: 2600, 5200, 7800
- PC and Arcade

## External Integrations

### TheGamesDB
Free metadata service (no API key required for basic usage)
- Game metadata and descriptions
- Cover art and screenshots
- Release information

### IGDB (Internet Game Database)
Comprehensive game database (requires free Twitch Developer account)
- Set `IGDB_API_KEY` in your environment
- Enhanced metadata and artwork
- More accurate game matching

## Troubleshooting

### Common Issues

1. **Database connection failed**:
   - Ensure PostgreSQL is running
   - Check DATABASE_URL format
   - Verify network connectivity

2. **ROM scanning not working**:
   - Check directory permissions
   - Verify mount paths in Docker
   - Ensure supported file extensions

3. **Missing cover art**:
   - Check internet connectivity
   - Verify API key configuration
   - Try manual metadata fetch

### Logs

```bash
# View application logs
docker-compose logs pelico

# Follow logs in real-time
docker-compose logs -f pelico
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions:
- Check the troubleshooting section
- Review Docker logs
- Create an issue on GitHub