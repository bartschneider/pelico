package config

import (
	"os"
)

type Config struct {
	Port           string
	DatabaseURL    string
	ROMPaths       []string
	TwitchClientID     string
	TwitchClientSecret string
	
	// Backup Configuration
	NextcloudURL      string
	NextcloudUsername string
	NextcloudPassword string
	NextcloudPath     string
	BackupEnabled     bool
	BackupSchedule    string
}

func Load() *Config {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://pelico:pelico@localhost/pelico?sslmode=disable"),
		TwitchClientID:     getEnv("TWITCH_CLIENT_ID", ""),
		TwitchClientSecret: getEnv("TWITCH_CLIENT_SECRET", ""),
		
		// Backup Configuration
		NextcloudURL:      getEnv("NEXTCLOUD_URL", ""),
		NextcloudUsername: getEnv("NEXTCLOUD_USERNAME", ""),
		NextcloudPassword: getEnv("NEXTCLOUD_PASSWORD", ""),
		NextcloudPath:     getEnv("NEXTCLOUD_PATH", "/Pelico/Backups"),
		BackupEnabled:     getEnv("BACKUP_ENABLED", "false") == "true",
		BackupSchedule:    getEnv("BACKUP_SCHEDULE", "daily"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}