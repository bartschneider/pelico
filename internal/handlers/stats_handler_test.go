package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pelico/internal/handlers"
	"pelico/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupStatsTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&models.Platform{},
		&models.Game{},
		&models.PlaySession{},
		&models.WishlistItem{},
		&models.ShortlistItem{},
	)
	require.NoError(t, err)

	// Create test data
	platform := models.Platform{
		Name:         "Nintendo Switch",
		Manufacturer: "Nintendo",
		ReleaseYear:  2017,
	}
	db.Create(&platform)

	games := []models.Game{
		{Title: "Game 1", PlatformID: 1, CompletionStatus: "completed"},
		{Title: "Game 2", PlatformID: 1, CompletionStatus: "in_progress"},
		{Title: "Game 3", PlatformID: 1, CompletionStatus: "not_started"},
		{Title: "Game 4", PlatformID: 1, CompletionStatus: "abandoned"},
	}
	for _, game := range games {
		db.Create(&game)
	}

	sessions := []models.PlaySession{
		{GameID: 1, PlayTimeMinutes: 120, Notes: "Great session"},
		{GameID: 2, PlayTimeMinutes: 60, Notes: "Good session"},
		{GameID: 1, PlayTimeMinutes: 90, Notes: "Another session"},
	}
	for _, session := range sessions {
		db.Create(&session)
	}

	wishlist := []models.WishlistItem{
		{Title: "Wishlist Game 1", PlatformName: "Switch"},
		{Title: "Wishlist Game 2", PlatformName: "PC"},
	}
	for _, item := range wishlist {
		db.Create(&item)
	}

	shortlist := []models.ShortlistItem{
		{Title: "Shortlist Game 1", PlatformName: "Switch"},
	}
	for _, item := range shortlist {
		db.Create(&item)
	}

	return db
}

func setupStatsRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	statsHandler := handlers.NewStatsHandler(db)
	router.GET("/api/v1/stats", statsHandler.GetStats)
	
	return router
}

func TestStatsHandler_GetStats(t *testing.T) {
	db := setupStatsTestDB(t)
	router := setupStatsRouter(db)

	req := httptest.NewRequest("GET", "/api/v1/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		TotalGames     int64 `json:"total_games"`
		TotalPlatforms int64 `json:"total_platforms"`
		TotalSessions  int64 `json:"total_sessions"`
		TotalWishlist  int64 `json:"total_wishlist"`
		TotalShortlist int64 `json:"total_shortlist"`
		CompletionStats struct {
			Completed  int64 `json:"completed"`
			InProgress int64 `json:"in_progress"`
			NotStarted int64 `json:"not_started"`
			Abandoned  int64 `json:"abandoned"`
		} `json:"completion_stats"`
		PlaytimeStats struct {
			TotalHours   float64 `json:"total_hours"`
			AverageHours float64 `json:"average_hours"`
		} `json:"playtime_stats"`
		PlatformBreakdown []struct {
			Platform string `json:"platform"`
			Count    int64  `json:"count"`
		} `json:"platform_breakdown"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify basic counts
	assert.Equal(t, int64(4), response.TotalGames)
	assert.Equal(t, int64(1), response.TotalPlatforms)
	assert.Equal(t, int64(3), response.TotalSessions)
	assert.Equal(t, int64(2), response.TotalWishlist)
	assert.Equal(t, int64(1), response.TotalShortlist)

	// Verify completion stats
	assert.Equal(t, int64(1), response.CompletionStats.Completed)
	assert.Equal(t, int64(1), response.CompletionStats.InProgress)
	assert.Equal(t, int64(1), response.CompletionStats.NotStarted)
	assert.Equal(t, int64(1), response.CompletionStats.Abandoned)

	// Verify playtime stats (120 + 60 + 90 = 270 minutes = 4.5 hours)
	assert.Equal(t, 4.5, response.PlaytimeStats.TotalHours)
	assert.Equal(t, 1.5, response.PlaytimeStats.AverageHours) // 4.5 / 3 sessions

	// Verify platform breakdown
	assert.Len(t, response.PlatformBreakdown, 1)
	assert.Equal(t, "Nintendo Switch", response.PlatformBreakdown[0].Platform)
	assert.Equal(t, int64(4), response.PlatformBreakdown[0].Count)
}

func TestStatsHandler_GetStats_EmptyDatabase(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&models.Platform{},
		&models.Game{},
		&models.PlaySession{},
		&models.WishlistItem{},
		&models.ShortlistItem{},
	)
	require.NoError(t, err)

	router := setupStatsRouter(db)

	req := httptest.NewRequest("GET", "/api/v1/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		TotalGames     int64 `json:"total_games"`
		TotalPlatforms int64 `json:"total_platforms"`
		TotalSessions  int64 `json:"total_sessions"`
		TotalWishlist  int64 `json:"total_wishlist"`
		TotalShortlist int64 `json:"total_shortlist"`
		CompletionStats struct {
			Completed  int64 `json:"completed"`
			InProgress int64 `json:"in_progress"`
			NotStarted int64 `json:"not_started"`
			Abandoned  int64 `json:"abandoned"`
		} `json:"completion_stats"`
		PlaytimeStats struct {
			TotalHours   float64 `json:"total_hours"`
			AverageHours float64 `json:"average_hours"`
		} `json:"playtime_stats"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// All counts should be zero
	assert.Equal(t, int64(0), response.TotalGames)
	assert.Equal(t, int64(0), response.TotalPlatforms)
	assert.Equal(t, int64(0), response.TotalSessions)
	assert.Equal(t, int64(0), response.TotalWishlist)
	assert.Equal(t, int64(0), response.TotalShortlist)

	// Completion stats should be zero
	assert.Equal(t, int64(0), response.CompletionStats.Completed)
	assert.Equal(t, int64(0), response.CompletionStats.InProgress)
	assert.Equal(t, int64(0), response.CompletionStats.NotStarted)
	assert.Equal(t, int64(0), response.CompletionStats.Abandoned)

	// Playtime stats should be zero
	assert.Equal(t, float64(0), response.PlaytimeStats.TotalHours)
	assert.Equal(t, float64(0), response.PlaytimeStats.AverageHours)
}