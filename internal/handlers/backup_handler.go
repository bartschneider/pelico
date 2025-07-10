package handlers

import (
	"fmt"
	"net/http"
	"time"

	"pelico/internal/config"
	"pelico/internal/models"
	"pelico/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BackupHandler struct {
	db              *gorm.DB
	nextcloudBackup *services.NextcloudBackup
}

func NewBackupHandler(db *gorm.DB, cfg *config.Config) *BackupHandler {
	return &BackupHandler{
		db:              db,
		nextcloudBackup: services.NewNextcloudBackup(cfg, db),
	}
}

func (h *BackupHandler) ExportDatabase(c *gin.Context) {
	backup, err := h.nextcloudBackup.CreateBackup()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create backup: " + err.Error()})
		return
	}

	// Set appropriate headers for download
	filename := fmt.Sprintf("pelico_backup_%s.json", backup.Timestamp.Format("20060102_150405"))
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/json")
	
	c.JSON(http.StatusOK, backup)
}

func (h *BackupHandler) ImportDatabase(c *gin.Context) {
	var backup services.BackupData
	if err := c.ShouldBindJSON(&backup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup data: " + err.Error()})
		return
	}

	// Validate backup version
	if backup.Version != "1.0" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported backup version: " + backup.Version})
		return
	}

	// Start transaction for atomic restore
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Clear existing data (be very careful here)
	if err := tx.Exec("DELETE FROM play_sessions").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear sessions: " + err.Error()})
		return
	}

	if err := tx.Exec("DELETE FROM file_locations").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear file locations: " + err.Error()})
		return
	}

	if err := tx.Exec("DELETE FROM games").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear games: " + err.Error()})
		return
	}

	if err := tx.Exec("DELETE FROM platforms").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear platforms: " + err.Error()})
		return
	}

	// Import platforms first (required for games)
	for _, platform := range backup.Platforms {
		if err := tx.Create(&platform).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import platform: " + err.Error()})
			return
		}
	}

	// Import games
	for _, game := range backup.Games {
		// Clear associations to avoid conflicts
		game.FileLocations = nil
		game.PlaySessions = nil
		
		if err := tx.Create(&game).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import game: " + err.Error()})
			return
		}
	}

	// Import file locations
	for _, fileLocation := range backup.FileLocations {
		if err := tx.Create(&fileLocation).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import file location: " + err.Error()})
			return
		}
	}

	// Import sessions
	for _, session := range backup.Sessions {
		if err := tx.Create(&session).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import session: " + err.Error()})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit backup restore: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Backup restored successfully",
		"stats": gin.H{
			"platforms":      len(backup.Platforms),
			"games":          len(backup.Games),
			"sessions":       len(backup.Sessions),
			"file_locations": len(backup.FileLocations),
			"backup_date":    backup.Timestamp,
		},
	})
}

func (h *BackupHandler) ExportGames(c *gin.Context) {
	var games []models.Game
	if err := h.db.Preload("Platform").Preload("FileLocations").Preload("PlaySessions").Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export games: " + err.Error()})
		return
	}

	exportData := struct {
		Version   string        `json:"version"`
		Timestamp time.Time     `json:"timestamp"`
		Games     []models.Game `json:"games"`
	}{
		Version:   "1.0",
		Timestamp: time.Now(),
		Games:     games,
	}

	filename := fmt.Sprintf("pelico_games_export_%s.json", time.Now().Format("20060102_150405"))
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, exportData)
}

func (h *BackupHandler) GetBackupInfo(c *gin.Context) {
	var stats struct {
		TotalGames     int64 `json:"total_games"`
		TotalPlatforms int64 `json:"total_platforms"`
		TotalSessions  int64 `json:"total_sessions"`
		TotalFiles     int64 `json:"total_files"`
	}

	h.db.Model(&models.Game{}).Count(&stats.TotalGames)
	h.db.Model(&models.Platform{}).Count(&stats.TotalPlatforms)
	h.db.Model(&models.PlaySession{}).Count(&stats.TotalSessions)
	h.db.Model(&models.FileLocation{}).Count(&stats.TotalFiles)

	c.JSON(http.StatusOK, gin.H{
		"database_stats": stats,
		"backup_info": gin.H{
			"supported_version": "1.0",
			"last_backup":       nil, // Could be enhanced to track last backup time
			"backup_size_mb":    0,   // Could be enhanced to calculate size
		},
		"nextcloud": gin.H{
			"configured": h.nextcloudBackup.IsConfigured(),
			"status":     "ready",
		},
	})
}

func (h *BackupHandler) TestNextcloudConnection(c *gin.Context) {
	if !h.nextcloudBackup.IsConfigured() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nextcloud backup is not configured. Please set NEXTCLOUD_URL, NEXTCLOUD_USERNAME, and NEXTCLOUD_PASSWORD environment variables.",
		})
		return
	}

	if err := h.nextcloudBackup.TestConnection(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":     "Failed to connect to Nextcloud",
			"details":   err.Error(),
			"configured": true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Nextcloud connection successful",
		"configured": true,
		"status":     "connected",
	})
}

func (h *BackupHandler) BackupToNextcloud(c *gin.Context) {
	if !h.nextcloudBackup.IsConfigured() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nextcloud backup is not configured",
		})
		return
	}

	if err := h.nextcloudBackup.PerformAutomaticBackup(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to backup to Nextcloud",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Backup successfully uploaded to Nextcloud",
		"timestamp": time.Now(),
	})
}