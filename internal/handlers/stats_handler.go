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

	// Completion statistics
	var completed, inProgress, notStarted, abandoned int64
	h.DB.Table("games").Where("completion_status = ?", "completed").Count(&completed)
	h.DB.Table("games").Where("completion_status = ?", "in_progress").Count(&inProgress)
	h.DB.Table("games").Where("completion_status = ? OR completion_status = ? OR completion_status IS NULL", "not_started", "").Count(&notStarted)
	h.DB.Table("games").Where("completion_status = ?", "abandoned").Count(&abandoned)

	// Playtime statistics
	var totalHours, avgHours float64
	h.DB.Table("play_sessions").Select("COALESCE(SUM(play_time_minutes), 0) / 60.0").Scan(&totalHours)
	if totalSessions > 0 {
		avgHours = totalHours / float64(totalSessions)
	}

	// Platform breakdown
	type PlatformStats struct {
		Platform string `json:"platform"`
		Count    int64  `json:"count"`
	}
	var platformBreakdown []PlatformStats
	h.DB.Table("games").
		Select("platforms.name as platform, COUNT(*) as count").
		Joins("LEFT JOIN platforms ON games.platform_id = platforms.id").
		Group("platforms.name").
		Scan(&platformBreakdown)

	c.JSON(http.StatusOK, gin.H{
		"total_games":     totalGames,
		"total_platforms": totalPlatforms,
		"total_sessions":  totalSessions,
		"total_wishlist":  totalWishlist,
		"total_shortlist": totalShortlist,
		"completion_stats": gin.H{
			"completed":   completed,
			"in_progress": inProgress,
			"not_started": notStarted,
			"abandoned":   abandoned,
		},
		"playtime_stats": gin.H{
			"total_hours":   totalHours,
			"average_hours": avgHours,
		},
		"platform_breakdown": platformBreakdown,
	})
}
