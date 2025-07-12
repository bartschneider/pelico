package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StatsHandler struct {
	DB *gorm.DB
}

func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{DB: db}
}

// GetStats retrieves various statistics about the game collection
func (h *StatsHandler) GetStats(c *gin.Context) {
	var totalGames int64
	h.DB.Table("games").Count(&totalGames)

	var totalPlatforms int64
	h.DB.Table("platforms").Count(&totalPlatforms)

	var totalSessions int64
	h.DB.Table("play_sessions").Count(&totalSessions)

	var totalWishlist int64
	h.DB.Table("wishlist").Count(&totalWishlist)

	var totalShortlist int64
	h.DB.Table("shortlist").Count(&totalShortlist)

	c.JSON(http.StatusOK, gin.H{
		"total_games":     totalGames,
		"total_platforms": totalPlatforms,
		"total_sessions":  totalSessions,
		"total_wishlist":  totalWishlist,
		"total_shortlist": totalShortlist,
	})
}
