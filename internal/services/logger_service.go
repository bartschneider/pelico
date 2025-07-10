package services

import (
	"log/slog"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggerService provides structured logging capabilities
type LoggerService struct {
	logger *slog.Logger
}

// NewLoggerService creates a new structured logger
func NewLoggerService(logLevel slog.Level) *LoggerService {
	// Create JSON handler for structured logging
	opts := &slog.HandlerOptions{
		Level: logLevel,
		AddSource: true,
	}
	
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	
	return &LoggerService{
		logger: logger,
	}
}

// GetLogger returns the underlying slog.Logger
func (ls *LoggerService) GetLogger() *slog.Logger {
	return ls.logger
}

// RequestLoggingMiddleware creates Gin middleware for request logging
func (ls *LoggerService) RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		
		// Store request ID in context for use in handlers
		c.Set("request_id", requestID)
		
		// Start time
		start := time.Now()
		
		// Log request start
		ls.logger.Info("request_start",
			slog.String("request_id", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
			slog.String("remote_addr", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
		)
		
		// Process request
		c.Next()
		
		// Calculate latency
		latency := time.Since(start)
		
		// Log request completion
		ls.logger.Info("request_complete",
			slog.String("request_id", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", latency),
			slog.Int("response_size", c.Writer.Size()),
		)
		
		// Log errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				ls.logger.Error("request_error",
					slog.String("request_id", requestID),
					slog.String("error", err.Error()),
					slog.Int("type", int(err.Type)),
				)
			}
		}
	}
}

// LogInfo logs an info message with optional attributes
func (ls *LoggerService) LogInfo(msg string, attrs ...slog.Attr) {
	ls.logger.LogAttrs(nil, slog.LevelInfo, msg, attrs...)
}

// LogError logs an error message with optional attributes
func (ls *LoggerService) LogError(msg string, err error, attrs ...slog.Attr) {
	allAttrs := append([]slog.Attr{slog.String("error", err.Error())}, attrs...)
	ls.logger.LogAttrs(nil, slog.LevelError, msg, allAttrs...)
}

// LogWarn logs a warning message with optional attributes
func (ls *LoggerService) LogWarn(msg string, attrs ...slog.Attr) {
	ls.logger.LogAttrs(nil, slog.LevelWarn, msg, attrs...)
}

// LogDebug logs a debug message with optional attributes
func (ls *LoggerService) LogDebug(msg string, attrs ...slog.Attr) {
	ls.logger.LogAttrs(nil, slog.LevelDebug, msg, attrs...)
}

// LogWithContext logs a message with request context
func (ls *LoggerService) LogWithContext(c *gin.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	requestID, exists := c.Get("request_id")
	if exists {
		attrs = append([]slog.Attr{slog.String("request_id", requestID.(string))}, attrs...)
	}
	
	ls.logger.LogAttrs(nil, level, msg, attrs...)
}

// LogGameOperation logs game-related operations
func (ls *LoggerService) LogGameOperation(c *gin.Context, operation string, gameID uint, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("operation", operation),
		slog.Uint64("game_id", uint64(gameID)),
		slog.String("entity", "game"),
	}
	
	allAttrs := append(baseAttrs, attrs...)
	ls.LogWithContext(c, slog.LevelInfo, "game_operation", allAttrs...)
}

// LogPlatformOperation logs platform-related operations
func (ls *LoggerService) LogPlatformOperation(c *gin.Context, operation string, platformID uint, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("operation", operation),
		slog.Uint64("platform_id", uint64(platformID)),
		slog.String("entity", "platform"),
	}
	
	allAttrs := append(baseAttrs, attrs...)
	ls.LogWithContext(c, slog.LevelInfo, "platform_operation", allAttrs...)
}

// LogSessionOperation logs session-related operations
func (ls *LoggerService) LogSessionOperation(c *gin.Context, operation string, sessionID uint, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("operation", operation),
		slog.Uint64("session_id", uint64(sessionID)),
		slog.String("entity", "session"),
	}
	
	allAttrs := append(baseAttrs, attrs...)
	ls.LogWithContext(c, slog.LevelInfo, "session_operation", allAttrs...)
}

// LogCacheOperation logs cache-related operations
func (ls *LoggerService) LogCacheOperation(operation string, key string, hit bool, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("operation", operation),
		slog.String("cache_key", key),
		slog.Bool("cache_hit", hit),
		slog.String("component", "cache"),
	}
	
	allAttrs := append(baseAttrs, attrs...)
	ls.logger.LogAttrs(nil, slog.LevelDebug, "cache_operation", allAttrs...)
}

// LogDatabaseOperation logs database-related operations
func (ls *LoggerService) LogDatabaseOperation(c *gin.Context, operation string, table string, duration time.Duration, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("operation", operation),
		slog.String("table", table),
		slog.Duration("duration", duration),
		slog.String("component", "database"),
	}
	
	allAttrs := append(baseAttrs, attrs...)
	ls.LogWithContext(c, slog.LevelDebug, "database_operation", allAttrs...)
}

// LogMetadataOperation logs external metadata service operations
func (ls *LoggerService) LogMetadataOperation(c *gin.Context, operation string, service string, success bool, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("operation", operation),
		slog.String("service", service),
		slog.Bool("success", success),
		slog.String("component", "metadata"),
	}
	
	allAttrs := append(baseAttrs, attrs...)
	level := slog.LevelInfo
	if !success {
		level = slog.LevelWarn
	}
	
	ls.LogWithContext(c, level, "metadata_operation", allAttrs...)
}

// LogBackupOperation logs backup/restore operations
func (ls *LoggerService) LogBackupOperation(c *gin.Context, operation string, destination string, success bool, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("operation", operation),
		slog.String("destination", destination),
		slog.Bool("success", success),
		slog.String("component", "backup"),
	}
	
	allAttrs := append(baseAttrs, attrs...)
	level := slog.LevelInfo
	if !success {
		level = slog.LevelError
	}
	
	ls.LogWithContext(c, level, "backup_operation", allAttrs...)
}