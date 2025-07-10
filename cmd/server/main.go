package main

import (
	"log"
	"pelico/internal/api"
	"pelico/internal/config"
	"pelico/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.Migrate(cfg.DatabaseURL); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Start API server
	server := api.NewServer(db, cfg)
	log.Printf("Starting server on port %s", cfg.Port)
	if err := server.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}