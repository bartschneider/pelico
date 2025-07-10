package handlers

import (
	"net/http"
	"strconv"
	"pelico/internal/errors"
	"pelico/internal/middleware"
	"pelico/internal/models"
	"pelico/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlatformHandler struct {
	db    *gorm.DB
	cache *services.CacheService
}

func NewPlatformHandler(db *gorm.DB, cache *services.CacheService) *PlatformHandler {
	return &PlatformHandler{
		db:    db,
		cache: cache,
	}
}

func (h *PlatformHandler) GetPlatforms(c *gin.Context) {
	// Try to get from cache first
	if cachedPlatforms, found := h.cache.GetPlatforms(); found {
		c.JSON(http.StatusOK, cachedPlatforms)
		return
	}
	
	var platforms []models.Platform
	result := h.db.Find(&platforms)
	
	if result.Error != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "fetch_platforms",
			"error": result.Error.Error(),
		})
		return
	}
	
	// Cache the platforms
	h.cache.SetPlatforms(platforms)
	
	c.JSON(http.StatusOK, platforms)
}

func (h *PlatformHandler) GetPlatform(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.ErrInvalidRequest, map[string]string{
			"parameter": "id",
			"expected": "positive integer",
			"received": c.Param("id"),
		})
		return
	}
	
	var platform models.Platform
	result := h.db.First(&platform, id)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			errors.RespondWithError(c, errors.ErrPlatformNotFound, map[string]interface{}{
				"platform_id": id,
			})
			return
		}
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "database_query",
			"error": result.Error.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, platform)
}

func (h *PlatformHandler) CreatePlatform(c *gin.Context) {
	var req middleware.CreatePlatformRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}
	
	platform := models.Platform{
		Name:         req.Name,
		Manufacturer: req.Manufacturer,
		ReleaseYear:  req.ReleaseYear,
	}
	
	result := h.db.Create(&platform)
	if result.Error != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "create_platform",
			"error": result.Error.Error(),
		})
		return
	}
	
	// Invalidate platforms cache
	h.cache.InvalidatePlatforms()
	
	c.JSON(http.StatusCreated, platform)
}

func (h *PlatformHandler) UpdatePlatform(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.ErrInvalidRequest, map[string]string{
			"parameter": "id",
			"expected": "positive integer",
			"received": c.Param("id"),
		})
		return
	}
	
	var req middleware.UpdatePlatformRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}
	
	var platform models.Platform
	if result := h.db.First(&platform, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			errors.RespondWithError(c, errors.ErrPlatformNotFound, map[string]interface{}{
				"platform_id": id,
			})
			return
		}
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "database_query",
			"error": result.Error.Error(),
		})
		return
	}
	
	// Update fields if provided in request
	if req.Name != "" {
		platform.Name = req.Name
	}
	if req.Manufacturer != "" {
		platform.Manufacturer = req.Manufacturer
	}
	if req.ReleaseYear != 0 {
		platform.ReleaseYear = req.ReleaseYear
	}
	
	result := h.db.Save(&platform)
	if result.Error != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "update_platform",
			"error": result.Error.Error(),
		})
		return
	}
	
	// Invalidate platforms cache
	h.cache.InvalidatePlatforms()
	
	c.JSON(http.StatusOK, platform)
}

func (h *PlatformHandler) DeletePlatform(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.ErrInvalidRequest, map[string]string{
			"parameter": "id",
			"expected": "positive integer",
			"received": c.Param("id"),
		})
		return
	}
	
	// Check if platform has associated games
	var gameCount int64
	h.db.Model(&models.Game{}).Where("platform_id = ?", id).Count(&gameCount)
	
	if gameCount > 0 {
		errors.RespondWithError(c, errors.ErrPlatformHasGames, map[string]interface{}{
			"platform_id": id,
			"game_count": gameCount,
		})
		return
	}
	
	result := h.db.Delete(&models.Platform{}, id)
	if result.Error != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "delete_platform",
			"error": result.Error.Error(),
		})
		return
	}
	
	if result.RowsAffected == 0 {
		errors.RespondWithError(c, errors.ErrPlatformNotFound, map[string]interface{}{
			"platform_id": id,
		})
		return
	}
	
	// Invalidate platforms cache
	h.cache.InvalidatePlatforms()
	
	c.JSON(http.StatusOK, gin.H{"message": "Platform deleted successfully"})
}