package services

import (
	"sync"
	"time"
	"pelico/internal/models"
)

// CacheEntry holds cached data with expiration
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

// CacheService provides in-memory caching for frequently accessed data
type CacheService struct {
	games      sync.Map // map[uint]*CacheEntry
	platforms  sync.Map // map[string]*CacheEntry (key: "platforms")
	stats      sync.Map // map[string]*CacheEntry (key: "completion_stats", etc.)
	ttl        time.Duration
	mutex      sync.RWMutex
}

// NewCacheService creates a new cache service with specified TTL
func NewCacheService(ttl time.Duration) *CacheService {
	cache := &CacheService{
		ttl: ttl,
	}
	
	// Start cleanup goroutine
	go cache.cleanup()
	
	return cache
}

// GetGame retrieves a cached game by ID
func (c *CacheService) GetGame(id uint) (*models.Game, bool) {
	if val, ok := c.games.Load(id); ok {
		entry := val.(*CacheEntry)
		if time.Now().Before(entry.ExpiresAt) {
			return entry.Data.(*models.Game), true
		}
		c.games.Delete(id)
	}
	return nil, false
}

// SetGame caches a game with TTL expiration
func (c *CacheService) SetGame(game *models.Game) {
	entry := &CacheEntry{
		Data:      game,
		ExpiresAt: time.Now().Add(c.ttl),
	}
	c.games.Store(game.ID, entry)
}

// InvalidateGame removes a game from cache
func (c *CacheService) InvalidateGame(id uint) {
	c.games.Delete(id)
}

// GetPlatforms retrieves cached platforms list
func (c *CacheService) GetPlatforms() ([]models.Platform, bool) {
	if val, ok := c.platforms.Load("platforms"); ok {
		entry := val.(*CacheEntry)
		if time.Now().Before(entry.ExpiresAt) {
			return entry.Data.([]models.Platform), true
		}
		c.platforms.Delete("platforms")
	}
	return nil, false
}

// SetPlatforms caches the platforms list
func (c *CacheService) SetPlatforms(platforms []models.Platform) {
	entry := &CacheEntry{
		Data:      platforms,
		ExpiresAt: time.Now().Add(c.ttl),
	}
	c.platforms.Store("platforms", entry)
}

// InvalidatePlatforms removes platforms from cache
func (c *CacheService) InvalidatePlatforms() {
	c.platforms.Delete("platforms")
}

// GetCompletionStats retrieves cached completion statistics
func (c *CacheService) GetCompletionStats() (interface{}, bool) {
	if val, ok := c.stats.Load("completion_stats"); ok {
		entry := val.(*CacheEntry)
		if time.Now().Before(entry.ExpiresAt) {
			return entry.Data, true
		}
		c.stats.Delete("completion_stats")
	}
	return nil, false
}

// SetCompletionStats caches completion statistics
func (c *CacheService) SetCompletionStats(stats interface{}) {
	entry := &CacheEntry{
		Data:      stats,
		ExpiresAt: time.Now().Add(c.ttl),
	}
	c.stats.Store("completion_stats", entry)
}

// InvalidateCompletionStats removes completion stats from cache
func (c *CacheService) InvalidateCompletionStats() {
	c.stats.Delete("completion_stats")
}

// GetRecentlyPlayed retrieves cached recently played games
func (c *CacheService) GetRecentlyPlayed() ([]models.Game, bool) {
	if val, ok := c.stats.Load("recently_played"); ok {
		entry := val.(*CacheEntry)
		if time.Now().Before(entry.ExpiresAt) {
			return entry.Data.([]models.Game), true
		}
		c.stats.Delete("recently_played")
	}
	return nil, false
}

// SetRecentlyPlayed caches recently played games
func (c *CacheService) SetRecentlyPlayed(games []models.Game) {
	entry := &CacheEntry{
		Data:      games,
		ExpiresAt: time.Now().Add(time.Minute * 5), // Shorter TTL for recent data
	}
	c.stats.Store("recently_played", entry)
}

// InvalidateRecentlyPlayed removes recently played games from cache
func (c *CacheService) InvalidateRecentlyPlayed() {
	c.stats.Delete("recently_played")
}

// Clear removes all cached data
func (c *CacheService) Clear() {
	c.games = sync.Map{}
	c.platforms = sync.Map{}
	c.stats = sync.Map{}
}

// cleanup periodically removes expired entries
func (c *CacheService) cleanup() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()
	
	for range ticker.C {
		now := time.Now()
		
		// Clean games cache
		c.games.Range(func(key, value interface{}) bool {
			entry := value.(*CacheEntry)
			if now.After(entry.ExpiresAt) {
				c.games.Delete(key)
			}
			return true
		})
		
		// Clean platforms cache
		c.platforms.Range(func(key, value interface{}) bool {
			entry := value.(*CacheEntry)
			if now.After(entry.ExpiresAt) {
				c.platforms.Delete(key)
			}
			return true
		})
		
		// Clean stats cache
		c.stats.Range(func(key, value interface{}) bool {
			entry := value.(*CacheEntry)
			if now.After(entry.ExpiresAt) {
				c.stats.Delete(key)
			}
			return true
		})
	}
}

// GetCacheStats returns cache usage statistics
func (c *CacheService) GetCacheStats() map[string]int {
	stats := map[string]int{
		"games":     0,
		"platforms": 0,
		"stats":     0,
	}
	
	c.games.Range(func(key, value interface{}) bool {
		stats["games"]++
		return true
	})
	
	c.platforms.Range(func(key, value interface{}) bool {
		stats["platforms"]++
		return true
	})
	
	c.stats.Range(func(key, value interface{}) bool {
		stats["stats"]++
		return true
	})
	
	return stats
}