package handlers

import (
	"net/http"
	"strconv"

	"pelico/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShortlistHandler struct {
	DB *gorm.DB
}

func NewShortlistHandler(db *gorm.DB) *ShortlistHandler {
	return &ShortlistHandler{DB: db}
}

// GetShortlist retrieves the user's shortlist
func (h *ShortlistHandler) GetShortlist(c *gin.Context) {
	var shortlist []models.Shortlist
	if err := h.DB.Preload("Game.Platform").Find(&shortlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve shortlist"})
		return
	}
	c.JSON(http.StatusOK, shortlist)
}

// AddToShortlist adds a game to the user's shortlist
func (h *ShortlistHandler) AddToShortlist(c *gin.Context) {
	var body struct {
		GameID uint `json:"game_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	shortlistItem := models.Shortlist{
		GameID: body.GameID,
	}

	if err := h.DB.Create(&shortlistItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add game to shortlist"})
		return
	}

	c.JSON(http.StatusCreated, shortlistItem)
}

// RemoveFromShortlist removes a game from the user's shortlist
func (h *ShortlistHandler) RemoveFromShortlist(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shortlist item ID"})
		return
	}

	if err := h.DB.Delete(&models.Shortlist{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove game from shortlist"})
		return
	}

	c.Status(http.StatusNoContent)
}
