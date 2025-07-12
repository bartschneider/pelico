package handlers

import (
	"net/http"
	"strconv"

	"pelico/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WishlistHandler struct {
	DB *gorm.DB
}

func NewWishlistHandler(db *gorm.DB) *WishlistHandler {
	return &WishlistHandler{DB: db}
}

// GetWishlist retrieves the user's wishlist
func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	var wishlist []models.Wishlist
	if err := h.DB.Preload("Game.Platform").Find(&wishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve wishlist"})
		return
	}
	c.JSON(http.StatusOK, wishlist)
}

// AddToWishlist adds a game to the user's wishlist
func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	var body struct {
		GameID uint `json:"game_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	wishlistItem := models.Wishlist{
		GameID: body.GameID,
	}

	if err := h.DB.Create(&wishlistItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add game to wishlist"})
		return
	}

	c.JSON(http.StatusCreated, wishlistItem)
}

// RemoveFromWishlist removes a game from the user's wishlist
func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist item ID"})
		return
	}

	if err := h.DB.Delete(&models.Wishlist{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove game from wishlist"})
		return
	}

	c.Status(http.StatusNoContent)
}
