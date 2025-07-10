# Pelico Changelog

All notable changes to this project will be documented in this file.

## [1.1.0] - 2025-07-10

### Added
- Docker container logs display in settings page
- Proper versioning system with build information
- Version and changelog display in settings
- Docker logs modal with container selection and line count options
- Real-time log viewing with color-coded log levels

### Fixed
- Metadata updater HTTP 500 errors for Gameboy Advance games
- Platform mapping for "Gameboy Advance" (one word) in IGDB service
- Missing platform mapping causing game metadata lookup failures

### Technical
- Added GetDockerLogs function to backup_handler.go
- Added /api/v1/docker/logs API endpoint
- Added showDockerLogs(), refreshDockerLogs() JavaScript functions
- Added Docker logs modal to index.html
- Enhanced platform filter system with robust ID-based mapping

## [1.0.0] - 2024-12-01

### Initial Release
- Game collection management
- Platform management
- ROM scanning functionality
- Play session tracking
- IGDB metadata integration
- Backup and restore functionality
- Nextcloud backup support
- PostgreSQL database backend
- Docker containerization