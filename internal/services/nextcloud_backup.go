package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"pelico/internal/config"
	"pelico/internal/models"
	"gorm.io/gorm"
)

type NextcloudBackup struct {
	baseURL  string
	username string
	password string
	path     string
	client   *http.Client
	db       *gorm.DB
}

func NewNextcloudBackup(cfg *config.Config, db *gorm.DB) *NextcloudBackup {
	return &NextcloudBackup{
		baseURL:  cfg.NextcloudURL,
		username: cfg.NextcloudUsername,
		password: cfg.NextcloudPassword,
		path:     cfg.NextcloudPath,
		client:   &http.Client{Timeout: 30 * time.Second},
		db:       db,
	}
}

func (n *NextcloudBackup) IsConfigured() bool {
	return n.baseURL != "" && n.username != "" && n.password != ""
}

func (n *NextcloudBackup) TestConnection() error {
	if !n.IsConfigured() {
		return fmt.Errorf("Nextcloud backup not configured")
	}

	// Test by creating the backup directory
	url := fmt.Sprintf("%s/remote.php/dav/files/%s%s",
		n.baseURL, n.username, n.path)

	req, err := http.NewRequest("MKCOL", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(n.username, n.password)

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	// 201 = created, 405 = already exists (method not allowed)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusMethodNotAllowed {
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}

	return nil
}

func (n *NextcloudBackup) CreateBackup() (*BackupData, error) {
	backup := &BackupData{
		Version:   "1.0",
		Timestamp: time.Now(),
	}

	// Export all data with associations
	if err := n.db.Preload("Platform").Preload("FileLocations").Preload("PlaySessions").Find(&backup.Games).Error; err != nil {
		return nil, fmt.Errorf("failed to export games: %w", err)
	}

	if err := n.db.Find(&backup.Platforms).Error; err != nil {
		return nil, fmt.Errorf("failed to export platforms: %w", err)
	}

	if err := n.db.Find(&backup.Sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to export sessions: %w", err)
	}

	if err := n.db.Find(&backup.FileLocations).Error; err != nil {
		return nil, fmt.Errorf("failed to export file locations: %w", err)
	}

	return backup, nil
}

func (n *NextcloudBackup) UploadBackup(backup *BackupData) error {
	if !n.IsConfigured() {
		return fmt.Errorf("Nextcloud backup not configured")
	}

	// Serialize backup to JSON
	data, err := json.Marshal(backup)
	if err != nil {
		return fmt.Errorf("failed to serialize backup: %w", err)
	}

	// Generate filename
	filename := fmt.Sprintf("pelico_backup_%s.json", backup.Timestamp.Format("20060102_150405"))
	
	url := fmt.Sprintf("%s/remote.php/dav/files/%s%s/%s",
		n.baseURL, n.username, n.path, filename)

	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create upload request: %w", err)
	}

	req.SetBasicAuth(n.username, n.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("upload failed with status: %s", resp.Status)
	}

	return nil
}

func (n *NextcloudBackup) PerformAutomaticBackup() error {
	if !n.IsConfigured() {
		return fmt.Errorf("automatic backup skipped: Nextcloud not configured")
	}

	// Create backup
	backup, err := n.CreateBackup()
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Upload to Nextcloud
	if err := n.UploadBackup(backup); err != nil {
		return fmt.Errorf("failed to upload backup: %w", err)
	}

	return nil
}

// BackupData represents the complete backup structure
type BackupData struct {
	Version       string                 `json:"version"`
	Timestamp     time.Time              `json:"timestamp"`
	Games         []models.Game          `json:"games"`
	Platforms     []models.Platform      `json:"platforms"`
	Sessions      []models.PlaySession   `json:"sessions"`
	FileLocations []models.FileLocation  `json:"file_locations"`
}

// BackupStats provides information about backup contents
type BackupStats struct {
	TotalGames     int `json:"total_games"`
	TotalPlatforms int `json:"total_platforms"`
	TotalSessions  int `json:"total_sessions"`
	TotalFiles     int `json:"total_files"`
}

func (b *BackupData) GetStats() BackupStats {
	return BackupStats{
		TotalGames:     len(b.Games),
		TotalPlatforms: len(b.Platforms),
		TotalSessions:  len(b.Sessions),
		TotalFiles:     len(b.FileLocations),
	}
}