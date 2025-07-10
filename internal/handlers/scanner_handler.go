package handlers

import (
	"fmt"
	"net/http"
	"time"
	"pelico/internal/errors"
	"pelico/internal/middleware"
	"pelico/internal/models"
	"pelico/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ScannerHandler struct {
	db              *gorm.DB
	scanner         *services.ROMScanner
	metadataService *services.MetadataService
}

func NewScannerHandler(db *gorm.DB, clientID, clientSecret string) *ScannerHandler {
	return &ScannerHandler{
		db:              db,
		scanner:         services.NewROMScanner(db),
		metadataService: services.NewMetadataService(clientID, clientSecret),
	}
}

func (h *ScannerHandler) ScanDirectory(c *gin.Context) {
	var req middleware.ScanDirectoryRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}
	
	// Verify platform exists
	var platform models.Platform
	if err := h.db.First(&platform, req.PlatformID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.RespondWithError(c, errors.ErrPlatformNotFound, map[string]interface{}{
				"platform_id": req.PlatformID,
			})
			return
		}
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "platform_lookup",
			"error": err.Error(),
		})
		return
	}
	
	results, err := h.scanner.ScanDirectory(req.DirectoryPath, req.ServerLocation, req.PlatformID, req.Recursive)
	if err != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "scan_directory",
			"directory": req.DirectoryPath,
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Directory scan completed",
		"results":    results,
		"files_found": len(results.FilesFound),
		"games_added": len(results.GamesAdded),
		"errors":     results.Errors,
	})
}

func (h *ScannerHandler) FindDuplicates(c *gin.Context) {
	duplicates, err := h.scanner.FindDuplicates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find duplicates: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"duplicates": duplicates,
		"count":      len(duplicates),
	})
}

func (h *ScannerHandler) UpdateMetadataBatch(c *gin.Context) {
	var request struct {
		GameIDs      []uint `json:"game_ids"`
		BatchSize    int    `json:"batch_size"`    // Default: 5
		DelaySeconds int    `json:"delay_seconds"` // Default: 2
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Set defaults
	if request.BatchSize == 0 {
		request.BatchSize = 5
	}
	if request.DelaySeconds == 0 {
		request.DelaySeconds = 2
	}
	
	// If no specific game IDs provided, get all games without metadata
	var gameIDs []uint
	if len(request.GameIDs) == 0 {
		var games []models.Game
		result := h.db.Where("description IS NULL OR description = '' OR cover_art_url IS NULL OR cover_art_url = ''").
			Preload("Platform").Find(&games)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch games: " + result.Error.Error()})
			return
		}
		
		for _, game := range games {
			gameIDs = append(gameIDs, game.ID)
		}
	} else {
		gameIDs = request.GameIDs
	}
	
	if len(gameIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No games need metadata updates",
			"updated": 0,
			"errors":  []string{},
		})
		return
	}
	
	// Start the batch update process
	go h.processBatchMetadataUpdate(gameIDs, request.BatchSize, request.DelaySeconds)
	
	c.JSON(http.StatusAccepted, gin.H{
		"message":     "Batch metadata update started",
		"total_games": len(gameIDs),
		"batch_size":  request.BatchSize,
		"delay":       request.DelaySeconds,
	})
}

func (h *ScannerHandler) processBatchMetadataUpdate(gameIDs []uint, batchSize, delaySeconds int) {
	updated := 0
	errors := []string{}
	
	for i := 0; i < len(gameIDs); i += batchSize {
		end := i + batchSize
		if end > len(gameIDs) {
			end = len(gameIDs)
		}
		
		batch := gameIDs[i:end]
		
		for _, gameID := range batch {
			// Fetch game with platform
			var game models.Game
			if err := h.db.Preload("Platform").First(&game, gameID).Error; err != nil {
				errors = append(errors, fmt.Sprintf("Game %d: Failed to fetch - %v", gameID, err))
				continue
			}
			
			// Fetch metadata
			metadata, err := h.metadataService.FetchGameMetadata(game.Title, game.Platform.Name)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Game %d (%s): Failed to fetch metadata - %v", gameID, game.Title, err))
				continue
			}
			
			// Update game with metadata
			updateData := map[string]interface{}{}
			if metadata.Description != "" {
				updateData["description"] = metadata.Description
			}
			if metadata.Rating > 0 {
				updateData["rating"] = metadata.Rating
			}
			if metadata.Genre != "" {
				updateData["genre"] = metadata.Genre
			}
			if metadata.Year > 0 {
				updateData["year"] = metadata.Year
			}
			if metadata.CoverArtURL != "" {
				updateData["cover_art_url"] = metadata.CoverArtURL
			}
			if metadata.BoxArtURL != "" {
				updateData["box_art_url"] = metadata.BoxArtURL
			}
			if metadata.IGDBID > 0 {
				updateData["igdb_id"] = metadata.IGDBID
			}
			
			if len(updateData) > 0 {
				if err := h.db.Model(&game).Updates(updateData).Error; err != nil {
					errors = append(errors, fmt.Sprintf("Game %d (%s): Failed to update - %v", gameID, game.Title, err))
					continue
				}
				updated++
			}
		}
		
		// Rate limiting delay between batches
		if i+batchSize < len(gameIDs) {
			time.Sleep(time.Duration(delaySeconds) * time.Second)
		}
	}
	
	// Log completion (in a real app, you might want to store this result somewhere)
	fmt.Printf("Batch metadata update completed: %d updated, %d errors\n", updated, len(errors))
	for _, err := range errors {
		fmt.Printf("Error: %s\n", err)
	}
}