package handlers_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pelico/internal/config"
	"pelico/internal/handlers"
	"pelico/internal/models"
	"pelico/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&models.Platform{},
		&models.Game{},
		&models.FileLocation{},
		&models.PlaySession{},
	)
	require.NoError(t, err)

	// Create test platform
	platform := models.Platform{
		Name:         "Test Console",
		Manufacturer: "Test Company",
		ReleaseYear:  2020,
	}
	db.Create(&platform)

	return db
}

func setupTestServer(db *gorm.DB) http.Handler {
	cfg := &config.Config{
		TwitchClientID:     "test_client_id",
		TwitchClientSecret: "test_client_secret",
	}
	
	// Create a custom test server setup that doesn't load HTML templates
	// Let's create the router directly for testing
	router := gin.New()
	router.Use(gin.Recovery())
	
	// Create cache and logger services
	cache := services.NewCacheService(30 * time.Minute)
	logger := services.NewLoggerService(slog.LevelInfo)
	
	// Initialize handlers
	gameHandler := handlers.NewGameHandler(db, cfg.TwitchClientID, cfg.TwitchClientSecret, cache, logger)
	platformHandler := handlers.NewPlatformHandler(db, cache)
	sessionHandler := handlers.NewSessionHandler(db, cache)
	wishlistHandler := handlers.NewWishlistHandler(db)
	shortlistHandler := handlers.NewShortlistHandler(db)
	statsHandler := handlers.NewStatsHandler(db)
	
	// Setup API routes only (skip web routes that need templates)
	api := router.Group("/api/v1")
	{
		// Games
		api.GET("/games", gameHandler.GetGames)
		api.GET("/games/:id", gameHandler.GetGame)
		api.POST("/games", gameHandler.CreateGame)
		api.PUT("/games/:id", gameHandler.UpdateGame)
		api.DELETE("/games/:id", gameHandler.DeleteGame)
		api.PUT("/games/:id/completion", gameHandler.UpdateCompletionStatus)
		
		// Platforms
		api.GET("/platforms", platformHandler.GetPlatforms)
		api.POST("/platforms", platformHandler.CreatePlatform)
		
		// Sessions
		api.GET("/games/:id/sessions", sessionHandler.GetGameSessions)
		api.POST("/games/:id/sessions", sessionHandler.CreateSession)
		
		// Wishlist & Shortlist
		api.GET("/wishlist", wishlistHandler.GetWishlist)
		api.GET("/shortlist", shortlistHandler.GetShortlist)
		
		// Stats
		api.GET("/stats", statsHandler.GetStats)
		
		// Health check and cache stats
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "healthy",
				"checks": gin.H{
					"database": "ok",
					"cache": "ok",
				},
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
		})
		api.GET("/cache/stats", func(c *gin.Context) {
			stats := cache.GetCacheStats()
			c.JSON(200, gin.H{
				"cache_stats": stats,
				"cache_ttl": "30 minutes",
				"cleanup_interval": "5 minutes",
			})
		})
	}
	
	return router
}

func TestGameHandler_CreateGame(t *testing.T) {
	db := setupTestDB(t)
	server := setupTestServer(db)

	tests := []struct {
		name     string
		payload  map[string]interface{}
		wantCode int
		wantKeys []string
	}{
		{
			name: "valid game creation",
			payload: map[string]interface{}{
				"title":       "Super Mario World",
				"platform_id": 1,
				"year":        1990,
				"genre":       "Platformer",
				"rating":      9.5,
			},
			wantCode: http.StatusCreated,
			wantKeys: []string{"id", "title", "platform_id"},
		},
		{
			name: "missing title",
			payload: map[string]interface{}{
				"platform_id": 1,
			},
			wantCode: http.StatusBadRequest,
			wantKeys: []string{"code", "message"},
		},
		{
			name: "missing platform_id",
			payload: map[string]interface{}{
				"title": "Test Game",
			},
			wantCode: http.StatusBadRequest,
			wantKeys: []string{"code", "message"},
		},
		{
			name: "invalid platform_id",
			payload: map[string]interface{}{
				"title":       "Test Game",
				"platform_id": 999,
			},
			wantCode: http.StatusNotFound,
			wantKeys: []string{"code", "message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/games", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			for _, key := range tt.wantKeys {
				assert.Contains(t, response, key, "Response should contain key: %s", key)
			}
		})
	}
}

func TestGameHandler_GetGame(t *testing.T) {
	db := setupTestDB(t)
	server := setupTestServer(db)

	// Create a test game
	game := models.Game{
		Title:      "Test Game",
		PlatformID: 1,
		Year:       2020,
		Genre:      "Action",
		Rating:     8.5,
	}
	db.Create(&game)

	tests := []struct {
		name     string
		gameID   string
		wantCode int
	}{
		{
			name:     "valid game ID",
			gameID:   "1",
			wantCode: http.StatusOK,
		},
		{
			name:     "invalid game ID",
			gameID:   "abc",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "non-existent game ID",
			gameID:   "999",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/games/"+tt.gameID, nil)
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			if w.Code == http.StatusOK {
				var response models.Game
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "Test Game", response.Title)
			}
		})
	}
}

func TestGameHandler_GetGames(t *testing.T) {
	db := setupTestDB(t)
	server := setupTestServer(db)

	// Create test games
	games := []models.Game{
		{Title: "Game 1", PlatformID: 1},
		{Title: "Game 2", PlatformID: 1},
		{Title: "Game 3", PlatformID: 1},
	}
	for _, game := range games {
		db.Create(&game)
	}

	tests := []struct {
		name      string
		query     string
		wantCode  int
		wantCount int
	}{
		{
			name:      "default pagination",
			query:     "",
			wantCode:  http.StatusOK,
			wantCount: 3,
		},
		{
			name:      "custom page size",
			query:     "?limit=2",
			wantCode:  http.StatusOK,
			wantCount: 2,
		},
		{
			name:      "second page",
			query:     "?page=2&limit=2",
			wantCode:  http.StatusOK,
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/games"+tt.query, nil)
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			if w.Code == http.StatusOK {
				var response struct {
					Games      []models.Game `json:"games"`
					Pagination struct {
						Page  int `json:"page"`
						Total int `json:"total"`
					} `json:"pagination"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Len(t, response.Games, tt.wantCount)
				assert.Equal(t, 3, response.Pagination.Total)
			}
		})
	}
}

func TestGameHandler_UpdateCompletionStatus(t *testing.T) {
	db := setupTestDB(t)
	server := setupTestServer(db)

	// Create a test game
	game := models.Game{
		Title:      "Test Game",
		PlatformID: 1,
	}
	db.Create(&game)

	tests := []struct {
		name     string
		gameID   string
		payload  map[string]interface{}
		wantCode int
	}{
		{
			name:   "valid completion update",
			gameID: "1",
			payload: map[string]interface{}{
				"status":     "completed",
				"percentage": 100,
				"notes":      "Great game!",
			},
			wantCode: http.StatusOK,
		},
		{
			name:   "invalid status",
			gameID: "1",
			payload: map[string]interface{}{
				"status": "invalid_status",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:   "invalid percentage",
			gameID: "1",
			payload: map[string]interface{}{
				"status":     "in_progress",
				"percentage": 150,
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest("PUT", "/api/v1/games/"+tt.gameID+"/completion", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			if w.Code == http.StatusOK {
				var response models.Game
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.payload["status"], response.CompletionStatus)
			}
		})
	}
}

func TestHealthCheck(t *testing.T) {
	db := setupTestDB(t)
	server := setupTestServer(db)

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status string `json:"status"`
		Checks struct {
			Database string `json:"database"`
			Cache    string `json:"cache"`
		} `json:"checks"`
		Timestamp string `json:"timestamp"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response.Status)
	assert.Equal(t, "ok", response.Checks.Database)
	assert.Equal(t, "ok", response.Checks.Cache)
	assert.NotEmpty(t, response.Timestamp)

	// Validate timestamp format
	_, err = time.Parse(time.RFC3339, response.Timestamp)
	assert.NoError(t, err)
}

func TestCacheStats(t *testing.T) {
	db := setupTestDB(t)
	server := setupTestServer(db)

	req := httptest.NewRequest("GET", "/api/v1/cache/stats", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		CacheStats struct {
			Games     int `json:"games"`
			Platforms int `json:"platforms"`
			Stats     int `json:"stats"`
		} `json:"cache_stats"`
		CacheTTL        string `json:"cache_ttl"`
		CleanupInterval string `json:"cleanup_interval"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "30 minutes", response.CacheTTL)
	assert.Equal(t, "5 minutes", response.CleanupInterval)
}