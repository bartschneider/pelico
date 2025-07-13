package api

import (
	"log/slog"
	"net/http"
	"time"
	"pelico/internal/config"
	"pelico/internal/handlers"
	"pelico/internal/services"
	"pelico/internal/version"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	db     *gorm.DB
	config *config.Config
	cache  *services.CacheService
	logger *services.LoggerService
}

func NewServer(db *gorm.DB, cfg *config.Config) *Server {
	// Initialize logger service first
	logger := services.NewLoggerService(slog.LevelInfo)
	
	// Create router with custom middleware
	router := gin.New()
	
	// Add structured logging middleware
	router.Use(logger.RequestLoggingMiddleware())
	
	// Add recovery middleware with logging
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.GetLogger().Error("panic_recovered", 
			slog.Any("panic", recovered),
			slog.String("path", c.Request.URL.Path),
			slog.String("method", c.Request.Method))
		c.AbortWithStatus(500)
	}))
	
	// Initialize cache service with 30-minute TTL
	cache := services.NewCacheService(30 * time.Minute)
	
	server := &Server{
		router: router,
		db:     db,
		config: cfg,
		cache:  cache,
		logger: logger,
	}
	
	// Log server initialization
	logger.LogInfo("server_initialized", 
		slog.String("cache_ttl", "30m"),
		slog.String("log_level", "info"))
	
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Initialize handlers with cache service and logger
	gameHandler := handlers.NewGameHandler(s.db, s.config.TwitchClientID, s.config.TwitchClientSecret, s.cache, s.logger)
	platformHandler := handlers.NewPlatformHandler(s.db, s.cache)
	sessionHandler := handlers.NewSessionHandler(s.db, s.cache)
	scannerHandler := handlers.NewScannerHandler(s.db, s.config.TwitchClientID, s.config.TwitchClientSecret)
	directoryHandler := handlers.NewDirectoryHandler()
	backupHandler := handlers.NewBackupHandler(s.db, s.config)
	wishlistHandler := handlers.NewWishlistHandler(s.db)
	shortlistHandler := handlers.NewShortlistHandler(s.db)
	statsHandler := handlers.NewStatsHandler(s.db)
	
	// API routes
	api := s.router.Group("/api/v1")
	{
		// Games
		api.GET("/games", gameHandler.GetGames)
		api.GET("/games/recently-played", gameHandler.GetRecentlyPlayedGames)
		api.GET("/games/genres", gameHandler.GetGenres)
		api.GET("/games/:id", gameHandler.GetGame)
		api.POST("/games", gameHandler.CreateGame)
		api.POST("/games/from-metadata", gameHandler.CreateGameFromMetadata)
		api.PUT("/games/:id", gameHandler.UpdateGame)
		api.DELETE("/games/:id", gameHandler.DeleteGame)
		api.POST("/games/search", gameHandler.SearchGames)
		api.POST("/games/search-metadata", gameHandler.SearchMetadata)
		
		// Platforms
		api.GET("/platforms", platformHandler.GetPlatforms)
		api.GET("/platforms/:id", platformHandler.GetPlatform)
		api.POST("/platforms", platformHandler.CreatePlatform)
		api.PUT("/platforms/:id", platformHandler.UpdatePlatform)
		api.DELETE("/platforms/:id", platformHandler.DeletePlatform)
		
		// Play Sessions
		api.GET("/games/:id/sessions", sessionHandler.GetGameSessions)
		api.POST("/games/:id/sessions", sessionHandler.CreateSession)
		api.PUT("/sessions/:id", sessionHandler.UpdateSession)
		api.DELETE("/sessions/:id", sessionHandler.DeleteSession)
		api.GET("/sessions/active", sessionHandler.GetActiveSessions)
		api.POST("/sessions/:id/end", sessionHandler.EndSession)
		
		// ROM Scanning
		api.POST("/scan/directory", scannerHandler.ScanDirectory)
		api.POST("/scan/metadata-batch", scannerHandler.UpdateMetadataBatch)
		api.GET("/scan/duplicates", scannerHandler.FindDuplicates)
		
		// Directory Browser
		api.GET("/browse", directoryHandler.BrowseDirectory)
		api.GET("/browse/suggestions", directoryHandler.GetSuggestedPaths)
		
		// Metadata
		api.POST("/games/:id/fetch-metadata", gameHandler.FetchMetadata)
		
		// Completion tracking
		api.PUT("/games/:id/completion", gameHandler.UpdateCompletionStatus)
		api.GET("/games/stats/completion", gameHandler.GetCompletionStats)
		api.GET("/games/completion/:status", gameHandler.GetGamesByCompletionStatus)

		// Wishlist
		api.GET("/wishlist", wishlistHandler.GetWishlist)
		api.POST("/wishlist", wishlistHandler.AddToWishlist)
		api.DELETE("/wishlist/:id", wishlistHandler.RemoveFromWishlist)

		// Shortlist
		api.GET("/shortlist", shortlistHandler.GetShortlist)
		api.POST("/shortlist", shortlistHandler.AddToShortlist)
		api.DELETE("/shortlist/:id", shortlistHandler.RemoveFromShortlist)

		// Statistics
		api.GET("/stats", statsHandler.GetStats)
		
		// Backup/Restore
		api.GET("/backup/export", backupHandler.ExportDatabase)
		api.POST("/backup/import", backupHandler.ImportDatabase)
		api.GET("/backup/export/games", backupHandler.ExportGames)
		api.GET("/backup/info", backupHandler.GetBackupInfo)
		api.POST("/backup/nextcloud", backupHandler.BackupToNextcloud)
		api.GET("/backup/nextcloud/test", backupHandler.TestNextcloudConnection)
		
		
		// Cache status
		api.GET("/cache/stats", s.getCacheStats)
		api.POST("/cache/clear", s.clearCache)
		
		// Health check
		api.GET("/health", s.healthCheck)
		
		// Version info
		api.GET("/version", s.getVersion)
	}
	
	// Serve static files for web interface
	s.router.Static("/static", "./web/static")
	s.router.LoadHTMLGlob("web/templates/*")
	
	// Web interface routes
	s.router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Pelico - Game Collection Manager",
		})
	})
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// ServeHTTP implements http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// getCacheStats returns cache usage statistics
func (s *Server) getCacheStats(c *gin.Context) {
	stats := s.cache.GetCacheStats()
	c.JSON(200, gin.H{
		"cache_stats": stats,
		"cache_ttl": "30 minutes",
		"cleanup_interval": "5 minutes",
	})
}

// clearCache clears all cached data
func (s *Server) clearCache(c *gin.Context) {
	s.cache.Clear()
	s.logger.LogWithContext(c, slog.LevelInfo, "cache_cleared")
	c.JSON(200, gin.H{
		"message": "Cache cleared successfully",
	})
}

// healthCheck returns the health status of the application
func (s *Server) healthCheck(c *gin.Context) {
	// Test database connection
	sqlDB, err := s.db.DB()
	if err != nil {
		s.logger.LogWithContext(c, slog.LevelError, "health_check_failed", 
			slog.String("component", "database"),
			slog.String("error", err.Error()))
		c.JSON(503, gin.H{
			"status": "unhealthy",
			"checks": gin.H{
				"database": "failed",
				"cache": "ok",
			},
			"error": "Database connection failed",
		})
		return
	}
	
	if err := sqlDB.Ping(); err != nil {
		s.logger.LogWithContext(c, slog.LevelError, "health_check_failed", 
			slog.String("component", "database"),
			slog.String("error", err.Error()))
		c.JSON(503, gin.H{
			"status": "unhealthy",
			"checks": gin.H{
				"database": "failed",
				"cache": "ok",
			},
			"error": "Database ping failed",
		})
		return
	}
	
	// Get cache stats for health check
	cacheStats := s.cache.GetCacheStats()
	
	c.JSON(200, gin.H{
		"status": "healthy",
		"checks": gin.H{
			"database": "ok",
			"cache": "ok",
		},
		"cache_stats": cacheStats,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// getVersion returns the application version information
func (s *Server) getVersion(c *gin.Context) {
	versionInfo := version.GetInfo()
	c.JSON(200, versionInfo)
}
