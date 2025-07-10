package handlers

import (
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"time"
	"pelico/internal/errors"
	"pelico/internal/middleware"
	"pelico/internal/models"
	"pelico/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GameHandler struct {
	db              *gorm.DB
	metadataService *services.MetadataService
	cache           *services.CacheService
	logger          *services.LoggerService
}

func NewGameHandler(db *gorm.DB, clientID, clientSecret string, cache *services.CacheService, logger *services.LoggerService) *GameHandler {
	return &GameHandler{
		db:              db,
		metadataService: services.NewMetadataService(clientID, clientSecret),
		cache:           cache,
		logger:          logger,
	}
}

func (h *GameHandler) GetGames(c *gin.Context) {
	// Parse pagination parameters
	page := 1
	limit := 50
	sort := "title"
	
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	
	if s := c.Query("sort"); s != "" {
		allowedSorts := []string{"title", "year", "rating", "created_at"}
		for _, allowed := range allowedSorts {
			if s == allowed {
				sort = s
				break
			}
		}
	}
	
	// Parse filter parameters
	platformFilter := c.Query("platform")
	genreFilter := c.Query("genre")
	completionFilter := c.Query("completion_status")
	
	// Build base query with filters
	baseQuery := h.db.Model(&models.Game{})
	
	// Apply platform filter
	if platformFilter != "" && platformFilter != "all" {
		// Try to parse as ID first, then fallback to name lookup
		if platformID, err := strconv.ParseUint(platformFilter, 10, 32); err == nil {
			// Filter value is a valid ID
			baseQuery = baseQuery.Where("platform_id = ?", uint(platformID))
		} else {
			// Filter value is a name, look up the platform
			var platform models.Platform
			result := h.db.Where("name = ?", platformFilter).First(&platform)
			if result.Error == nil && platform.ID > 0 {
				baseQuery = baseQuery.Where("platform_id = ?", platform.ID)
			}
		}
	}
	
	// Apply genre filter
	if genreFilter != "" && genreFilter != "all" {
		baseQuery = baseQuery.Where("genre = ?", genreFilter)
	}
	
	// Apply completion status filter
	if completionFilter != "" && completionFilter != "all" {
		if completionFilter == "backlog" {
			baseQuery = baseQuery.Where("completion_status IN ?", []string{"not_started", "in_progress"})
		} else {
			baseQuery = baseQuery.Where("completion_status = ?", completionFilter)
		}
	}
	
	// Get total count with filters applied
	var total int64
	countResult := baseQuery.Count(&total)
	if countResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count games"})
		return
	}
	
	// Calculate offset
	offset := (page - 1) * limit
	
	// Get paginated games with filters
	var games []models.Game
	query := baseQuery.Preload("Platform").Preload("FileLocations").
		Offset(offset).Limit(limit)
	
	// Apply sorting
	switch sort {
	case "year":
		query = query.Order("year DESC")
	case "rating":
		query = query.Order("rating DESC")
	case "created_at":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("title ASC")
	}
	
	result := query.Find(&games)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch games"})
		return
	}
	
	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)
	
	// Log filter operation
	if h.logger != nil {
		// Determine platform filter type for logging
		platformFilterType := "none"
		if platformFilter != "" && platformFilter != "all" {
			if _, err := strconv.ParseUint(platformFilter, 10, 32); err == nil {
				platformFilterType = "id"
			} else {
				platformFilterType = "name"
			}
		}
		
		h.logger.LogWithContext(c, slog.LevelDebug, "games_filtered",
			slog.String("platform_filter", platformFilter),
			slog.String("platform_filter_type", platformFilterType),
			slog.String("genre_filter", genreFilter),
			slog.String("completion_filter", completionFilter),
			slog.Int64("filtered_total", total),
			slog.Int("page", page))
	}
	
	c.JSON(http.StatusOK, gin.H{
		"games": games,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    page < int(totalPages),
			"has_prev":    page > 1,
		},
		"filters": gin.H{
			"platform":   platformFilter,
			"genre":      genreFilter,
			"completion": completionFilter,
		},
	})
}

func (h *GameHandler) GetGame(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.ErrInvalidRequest, map[string]string{
			"parameter": "id",
			"expected": "positive integer",
			"received": c.Param("id"),
		})
		return
	}
	
	gameID := uint(id)
	
	// Try to get from cache first
	if cachedGame, found := h.cache.GetGame(gameID); found {
		h.logger.LogCacheOperation("get", "game", true, slog.Uint64("game_id", uint64(gameID)))
		c.JSON(http.StatusOK, cachedGame)
		return
	}
	
	h.logger.LogCacheOperation("get", "game", false, slog.Uint64("game_id", uint64(gameID)))
	
	var game models.Game
	result := h.db.Preload("Platform").Preload("FileLocations").Preload("PlaySessions").First(&game, id)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			errors.RespondWithError(c, errors.ErrGameNotFound, map[string]interface{}{
				"game_id": id,
			})
			return
		}
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "database_query",
			"error": result.Error.Error(),
		})
		return
	}
	
	// Cache the game for future requests
	h.cache.SetGame(&game)
	
	c.JSON(http.StatusOK, game)
}

func (h *GameHandler) CreateGame(c *gin.Context) {
	var req middleware.CreateGameRequest
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
	
	// Create game from validated request
	game := models.Game{
		Title:       req.Title,
		PlatformID:  req.PlatformID,
		Year:        req.Year,
		Genre:       req.Genre,
		Rating:      req.Rating,
		Description: req.Description,
		CoverArtURL: req.CoverArtURL,
	}
	
	result := h.db.Create(&game)
	if result.Error != nil {
		h.logger.LogGameOperation(c, "create", 0, 
			slog.String("error", result.Error.Error()),
			slog.Bool("success", false))
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "create_game",
			"error": result.Error.Error(),
		})
		return
	}
	
	// Log successful creation
	h.logger.LogGameOperation(c, "create", game.ID, 
		slog.String("title", game.Title),
		slog.Uint64("platform_id", uint64(game.PlatformID)),
		slog.Bool("success", true))
	
	// Invalidate related caches
	h.cache.InvalidateCompletionStats()
	h.cache.InvalidateRecentlyPlayed()
	
	c.JSON(http.StatusCreated, game)
}

func (h *GameHandler) UpdateGame(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.ErrInvalidRequest, map[string]string{
			"parameter": "id",
			"expected": "positive integer",
			"received": c.Param("id"),
		})
		return
	}
	
	var req middleware.UpdateGameRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}
	
	var game models.Game
	if result := h.db.First(&game, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			errors.RespondWithError(c, errors.ErrGameNotFound, map[string]interface{}{
				"game_id": id,
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
	if req.Title != "" {
		game.Title = req.Title
	}
	if req.PlatformID != 0 {
		// Verify platform exists if being updated
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
		game.PlatformID = req.PlatformID
	}
	if req.Year != 0 {
		game.Year = req.Year
	}
	if req.Genre != "" {
		game.Genre = req.Genre
	}
	if req.Rating != 0 {
		game.Rating = req.Rating
	}
	if req.Description != "" {
		game.Description = req.Description
	}
	if req.CoverArtURL != "" {
		game.CoverArtURL = req.CoverArtURL
	}
	
	result := h.db.Save(&game)
	if result.Error != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "update_game",
			"error": result.Error.Error(),
		})
		return
	}
	
	// Invalidate cache for this game and related data
	h.cache.InvalidateGame(uint(id))
	h.cache.InvalidateCompletionStats()
	h.cache.InvalidateRecentlyPlayed()
	
	c.JSON(http.StatusOK, game)
}

func (h *GameHandler) DeleteGame(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.ErrInvalidRequest, map[string]string{
			"parameter": "id",
			"expected": "positive integer",
			"received": c.Param("id"),
		})
		return
	}
	
	result := h.db.Delete(&models.Game{}, id)
	if result.Error != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "delete_game",
			"error": result.Error.Error(),
		})
		return
	}
	
	if result.RowsAffected == 0 {
		errors.RespondWithError(c, errors.ErrGameNotFound, map[string]interface{}{
			"game_id": id,
		})
		return
	}
	
	// Invalidate caches since game was deleted
	h.cache.InvalidateGame(uint(id))
	h.cache.InvalidateCompletionStats()
	h.cache.InvalidateRecentlyPlayed()
	
	c.JSON(http.StatusOK, gin.H{"message": "Game deleted successfully"})
}

func (h *GameHandler) SearchGames(c *gin.Context) {
	var searchParams struct {
		Title    string `json:"title"`
		Platform string `json:"platform"`
		Genre    string `json:"genre"`
		Year     int    `json:"year"`
	}
	
	if err := c.ShouldBindJSON(&searchParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	query := h.db.Preload("Platform").Preload("FileLocations")
	
	if searchParams.Title != "" {
		query = query.Where("title ILIKE ?", "%"+searchParams.Title+"%")
	}
	if searchParams.Genre != "" {
		query = query.Where("genre ILIKE ?", "%"+searchParams.Genre+"%")
	}
	if searchParams.Year != 0 {
		query = query.Where("year = ?", searchParams.Year)
	}
	if searchParams.Platform != "" {
		query = query.Joins("JOIN platforms ON games.platform_id = platforms.id").
			Where("platforms.name ILIKE ?", "%"+searchParams.Platform+"%")
	}
	
	var games []models.Game
	result := query.Find(&games)
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	
	c.JSON(http.StatusOK, games)
}

func (h *GameHandler) FetchMetadata(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}
	
	var game models.Game
	if result := h.db.Preload("Platform").First(&game, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	
	metadata, err := h.metadataService.FetchGameMetadata(game.Title, game.Platform.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch metadata: " + err.Error()})
		return
	}
	
	// Update game with fetched metadata
	if metadata.Description != "" {
		game.Description = metadata.Description
	}
	if metadata.Rating > 0 {
		game.Rating = metadata.Rating
	}
	if metadata.Genre != "" {
		game.Genre = metadata.Genre
	}
	if metadata.Year > 0 {
		game.Year = metadata.Year
	}
	if metadata.CoverArtURL != "" {
		game.CoverArtURL = metadata.CoverArtURL
	}
	if metadata.BoxArtURL != "" {
		game.BoxArtURL = metadata.BoxArtURL
	}
	if metadata.IGDBID > 0 {
		game.IGDBID = metadata.IGDBID
	}
	
	if result := h.db.Save(&game); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update game: " + result.Error.Error()})
		return
	}
	
	c.JSON(http.StatusOK, game)
}

func (h *GameHandler) SearchMetadata(c *gin.Context) {
	var req middleware.SearchGamesRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}
	
	results, err := h.metadataService.SearchGames(req.Title, req.Platform)
	if err != nil {
		errors.RespondWithError(c, errors.ErrMetadataAPIError, map[string]string{
			"operation": "search_games",
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"results": results})
}

func (h *GameHandler) CreateGameFromMetadata(c *gin.Context) {
	var request struct {
		PlatformID uint                          `json:"platform_id" binding:"required"`
		Metadata   services.GameMetadata `json:"metadata" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Create game with metadata
	game := models.Game{
		Title:       request.Metadata.Title,
		PlatformID:  request.PlatformID,
		Year:        request.Metadata.Year,
		Genre:       request.Metadata.Genre,
		Rating:      request.Metadata.Rating,
		Description: request.Metadata.Description,
		CoverArtURL: request.Metadata.CoverArtURL,
		BoxArtURL:   request.Metadata.BoxArtURL,
		IGDBID:      request.Metadata.IGDBID,
	}
	
	result := h.db.Create(&game)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	
	// Reload with platform
	h.db.Preload("Platform").Preload("FileLocations").First(&game, game.ID)
	
	c.JSON(http.StatusCreated, game)
}

func (h *GameHandler) GetRecentlyPlayedGames(c *gin.Context) {
	// Try to get from cache first
	if cachedGames, found := h.cache.GetRecentlyPlayed(); found {
		c.JSON(http.StatusOK, cachedGames)
		return
	}
	
	// Get games with recent play sessions (last 30 days)
	var games []models.Game
	
	// Query for games that have play sessions in the last 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	
	result := h.db.Preload("Platform").Preload("FileLocations").
		Preload("PlaySessions", func(db *gorm.DB) *gorm.DB {
			return db.Where("start_time >= ?", thirtyDaysAgo).Order("start_time DESC").Limit(1)
		}).
		Joins("JOIN play_sessions ON games.id = play_sessions.game_id").
		Where("play_sessions.start_time >= ?", thirtyDaysAgo).
		Group("games.id").
		Order("MAX(play_sessions.start_time) DESC").
		Limit(6).
		Find(&games)
	
	if result.Error != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "fetch_recently_played",
			"error": result.Error.Error(),
		})
		return
	}
	
	// Cache the results
	h.cache.SetRecentlyPlayed(games)
	
	c.JSON(http.StatusOK, games)
}

// UpdateCompletionStatus updates the completion status of a game
func (h *GameHandler) UpdateCompletionStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.ErrInvalidRequest, map[string]string{
			"parameter": "id",
			"expected": "positive integer",
			"received": c.Param("id"),
		})
		return
	}
	
	var req middleware.CompletionStatusRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}
	
	var game models.Game
	if err := h.db.First(&game, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.RespondWithError(c, errors.ErrGameNotFound, map[string]interface{}{
				"game_id": id,
			})
			return
		}
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "database_query",
			"error": err.Error(),
		})
		return
	}
	
	// Update completion fields
	game.CompletionStatus = req.Status
	game.CompletionPercentage = req.Percentage
	game.CompletionNotes = req.Notes
	
	// Set completion date when marked as completed
	if req.Status == "completed" || req.Status == "100_percent" {
		now := time.Now()
		game.CompletionDate = &now
	} else {
		game.CompletionDate = nil
	}
	
	if err := h.db.Save(&game).Error; err != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "update_completion",
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, game)
}

// GetCompletionStats returns completion statistics for all games
func (h *GameHandler) GetCompletionStats(c *gin.Context) {
	// Try to get from cache first
	if cachedStats, found := h.cache.GetCompletionStats(); found {
		c.JSON(http.StatusOK, cachedStats)
		return
	}
	
	var stats struct {
		TotalGames      int64   `json:"total_games"`
		Completed       int64   `json:"completed"`
		InProgress      int64   `json:"in_progress"`
		NotStarted      int64   `json:"not_started"`
		Abandoned       int64   `json:"abandoned"`
		HundredPercent  int64   `json:"hundred_percent"`
		CompletionRate  float64 `json:"completion_rate"`
		AverageProgress float64 `json:"average_progress"`
	}
	
	h.db.Model(&models.Game{}).Count(&stats.TotalGames)
	h.db.Model(&models.Game{}).Where("completion_status = ?", "completed").Count(&stats.Completed)
	h.db.Model(&models.Game{}).Where("completion_status = ?", "in_progress").Count(&stats.InProgress)
	h.db.Model(&models.Game{}).Where("completion_status = ?", "not_started").Count(&stats.NotStarted)
	h.db.Model(&models.Game{}).Where("completion_status = ?", "abandoned").Count(&stats.Abandoned)
	h.db.Model(&models.Game{}).Where("completion_status = ?", "100_percent").Count(&stats.HundredPercent)
	
	// Calculate completion rate (completed + 100% / total)
	if stats.TotalGames > 0 {
		completed := stats.Completed + stats.HundredPercent
		stats.CompletionRate = float64(completed) / float64(stats.TotalGames) * 100
		
		// Calculate average progress across all games
		var totalProgress struct {
			Sum float64
		}
		h.db.Model(&models.Game{}).Select("AVG(completion_percentage) as sum").Scan(&totalProgress)
		stats.AverageProgress = totalProgress.Sum
	}
	
	// Cache the stats
	h.cache.SetCompletionStats(stats)
	
	c.JSON(http.StatusOK, stats)
}

// GetGamesByCompletionStatus returns games filtered by completion status
func (h *GameHandler) GetGamesByCompletionStatus(c *gin.Context) {
	status := c.Param("status")
	
	// Validate status
	validStatuses := []string{"not_started", "in_progress", "completed", "abandoned", "100_percent", "backlog"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}
	
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid completion status"})
		return
	}
	
	var games []models.Game
	query := h.db.Preload("Platform").Preload("FileLocations")
	
	// Handle special "backlog" status (not_started + in_progress)
	if status == "backlog" {
		query = query.Where("completion_status IN ?", []string{"not_started", "in_progress"})
	} else {
		query = query.Where("completion_status = ?", status)
	}
	
	if err := query.Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, games)
}

// GetGenres returns all unique genres from the game collection
func (h *GameHandler) GetGenres(c *gin.Context) {
	var genres []string
	
	// Get distinct genres from all games
	result := h.db.Model(&models.Game{}).
		Distinct("genre").
		Where("genre != '' AND genre IS NOT NULL").
		Pluck("genre", &genres)
	
	if result.Error != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "fetch_genres",
			"error": result.Error.Error(),
		})
		return
	}
	
	// Sort genres alphabetically
	sort.Strings(genres)
	
	c.JSON(http.StatusOK, genres)
}