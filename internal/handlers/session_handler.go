package handlers

import (
	"net/http"
	"strconv"
	"time"
	"pelico/internal/errors"
	"pelico/internal/middleware"
	"pelico/internal/models"
	"pelico/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SessionHandler struct {
	db    *gorm.DB
	cache *services.CacheService
}

func NewSessionHandler(db *gorm.DB, cache *services.CacheService) *SessionHandler {
	return &SessionHandler{
		db:    db,
		cache: cache,
	}
}

func (h *SessionHandler) GetGameSessions(c *gin.Context) {
	gameID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.ErrInvalidRequest, map[string]string{
			"parameter": "id",
			"expected": "positive integer",
			"received": c.Param("id"),
		})
		return
	}
	
	var sessions []models.PlaySession
	result := h.db.Where("game_id = ?", gameID).Order("start_time DESC").Find(&sessions)
	
	if result.Error != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "fetch_sessions",
			"error": result.Error.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, sessions)
}

func (h *SessionHandler) CreateSession(c *gin.Context) {
	gameID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.ErrInvalidRequest, map[string]string{
			"parameter": "id",
			"expected": "positive integer",
			"received": c.Param("id"),
		})
		return
	}
	
	var req middleware.CreateSessionRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}
	
	// Parse start time
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		startTime = time.Now()
	}
	
	// Parse end time if provided
	var endTime *time.Time
	if req.EndTime != "" {
		if parsed, err := time.Parse(time.RFC3339, req.EndTime); err == nil {
			endTime = &parsed
		}
	}
	
	session := models.PlaySession{
		GameID:    uint(gameID),
		StartTime: startTime,
		EndTime:   endTime,
		Notes:     req.Notes,
		Rating:    req.Rating,
	}
	
	// Calculate duration if end time is provided
	if endTime != nil {
		duration := endTime.Sub(startTime)
		session.Duration = int(duration.Minutes())
	}
	
	result := h.db.Create(&session)
	if result.Error != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "create_session",
			"error": result.Error.Error(),
		})
		return
	}
	
	// Invalidate recently played cache since session affects it
	h.cache.InvalidateRecentlyPlayed()
	
	c.JSON(http.StatusCreated, session)
}

func (h *SessionHandler) UpdateSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.ErrInvalidRequest, map[string]string{
			"parameter": "id",
			"expected": "positive integer",
			"received": c.Param("id"),
		})
		return
	}
	
	var req middleware.UpdateSessionRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}
	
	var session models.PlaySession
	if result := h.db.First(&session, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			errors.RespondWithError(c, errors.ErrSessionNotFound, map[string]interface{}{
				"session_id": id,
			})
			return
		}
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "database_query",
			"error": result.Error.Error(),
		})
		return
	}
	
	// Update fields if provided
	if req.StartTime != "" {
		if parsed, err := time.Parse(time.RFC3339, req.StartTime); err == nil {
			session.StartTime = parsed
		}
	}
	if req.EndTime != "" {
		if parsed, err := time.Parse(time.RFC3339, req.EndTime); err == nil {
			session.EndTime = &parsed
			// Recalculate duration
			duration := parsed.Sub(session.StartTime)
			session.Duration = int(duration.Minutes())
		}
	}
	if req.Notes != "" {
		session.Notes = req.Notes
	}
	if req.Rating != 0 {
		session.Rating = req.Rating
	}
	
	result := h.db.Save(&session)
	if result.Error != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "update_session",
			"error": result.Error.Error(),
		})
		return
	}
	
	// Invalidate recently played cache since session affects it
	h.cache.InvalidateRecentlyPlayed()
	
	c.JSON(http.StatusOK, session)
}

func (h *SessionHandler) DeleteSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.ErrInvalidRequest, map[string]string{
			"parameter": "id",
			"expected": "positive integer",
			"received": c.Param("id"),
		})
		return
	}
	
	result := h.db.Delete(&models.PlaySession{}, id)
	if result.Error != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "delete_session",
			"error": result.Error.Error(),
		})
		return
	}
	
	if result.RowsAffected == 0 {
		errors.RespondWithError(c, errors.ErrSessionNotFound, map[string]interface{}{
			"session_id": id,
		})
		return
	}
	
	// Invalidate recently played cache since session affects it
	h.cache.InvalidateRecentlyPlayed()
	
	c.JSON(http.StatusOK, gin.H{"message": "Session deleted successfully"})
}

// GetAllSessions returns all sessions with game info
func (h *SessionHandler) GetAllSessions(c *gin.Context) {
	var sessions []models.PlaySession
	
	// Get all sessions with preloaded game and platform info
	result := h.db.Preload("Game").
		Preload("Game.Platform").
		Order("start_time DESC").
		Find(&sessions)
	
	if result.Error != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "fetch_all_sessions",
			"error": result.Error.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, sessions)
}

// GetActiveSessions returns all sessions without end_time (currently active)
func (h *SessionHandler) GetActiveSessions(c *gin.Context) {
	var sessions []models.PlaySession
	
	// Find all sessions without end_time
	result := h.db.Where("end_time IS NULL").
		Preload("Game").
		Preload("Game.Platform").
		Order("start_time DESC").
		Find(&sessions)
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	
	// Calculate current duration for each active session
	type ActiveSession struct {
		models.PlaySession
		CurrentDuration int `json:"current_duration"` // minutes
	}
	
	activeSessions := make([]ActiveSession, len(sessions))
	now := time.Now()
	
	for i, session := range sessions {
		duration := int(now.Sub(session.StartTime).Minutes())
		activeSessions[i] = ActiveSession{
			PlaySession:     session,
			CurrentDuration: duration,
		}
	}
	
	c.JSON(http.StatusOK, activeSessions)
}

// EndSession sets the end_time for an active session
func (h *SessionHandler) EndSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.ErrInvalidRequest, map[string]string{
			"parameter": "id",
			"expected": "positive integer",
			"received": c.Param("id"),
		})
		return
	}
	
	var session models.PlaySession
	if err := h.db.First(&session, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.RespondWithError(c, errors.ErrSessionNotFound, map[string]interface{}{
				"session_id": id,
			})
			return
		}
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "database_query",
			"error": err.Error(),
		})
		return
	}
	
	// Check if session is already ended
	if session.EndTime != nil {
		errors.RespondWithError(c, errors.ErrSessionAlreadyEnded, map[string]interface{}{
			"session_id": id,
			"end_time": session.EndTime,
		})
		return
	}
	
	// Update end_time to now and calculate duration
	now := time.Now()
	session.EndTime = &now
	session.Duration = int(now.Sub(session.StartTime).Minutes())
	
	if err := h.db.Save(&session).Error; err != nil {
		errors.RespondWithError(c, errors.ErrInternalServer, map[string]string{
			"operation": "end_session",
			"error": err.Error(),
		})
		return
	}
	
	// Invalidate recently played cache since session affects it
	h.cache.InvalidateRecentlyPlayed()
	
	c.JSON(http.StatusOK, session)
}