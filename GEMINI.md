# Gemini Project Overview: Pelico

This document provides a high-level overview of the Pelico project for the Gemini AI assistant.

## Project Description

Pelico is a self-hosted, web-based application designed for managing and cataloging video game collections. It allows users to scan directories for game ROMs, fetch metadata from external APIs, track play sessions, and manage their game library through a user-friendly interface.

## Key Features

- **ROM Scanning**: Automatically scans user-specified directories to identify and import game files.
- **Metadata Fetching**: Integrates with IGDB and TheGamesDB to enrich the game library with cover art, descriptions, and other details.
- **Platform Management**: Supports a wide variety of gaming platforms.
- **Session Tracking**: Allows users to log and rate their gaming sessions.
- **Web UI**: A browser-based interface for easy access and management.
- **REST API**: Provides a comprehensive API for programmatic access to the collection.

## Technology Stack

- **Backend**: Go
  - **Web Framework**: Gin
  - **ORM**: Gorm
- **Database**: PostgreSQL
- **Frontend**: Standard HTML, CSS, and JavaScript.
- **Containerization**: Docker and Docker Compose for easy deployment.

## Project Structure

- `cmd/server/main.go`: The main application entry point.
- `internal/`: Contains the core application logic.
  - `api/`: Defines the Gin server and API routes.
  - `handlers/`: HTTP handlers for API endpoints.
  - `services/`: Business logic for features like ROM scanning and metadata fetching.
  - `models/`: Database models (Gorm structs).
  - `database/`: Database connection and migration logic.
- `web/`: Frontend assets, including HTML templates, CSS, and JavaScript.
- `go.mod`: Defines project dependencies.
- `docker-compose.yml`: Defines the services for running the application.
- `Dockerfile`: Used to build the application's Docker image.

## How to Run the Application

1.  **Using Docker (Recommended)**:
    ```bash
    docker-compose up -d
    ```
    The application will be available at `http://localhost:8080`.

2.  **Running from Source**:
    ```bash
    go run cmd/server/main.go
    ```

## How to Run Tests

```bash
go test ./...
```
