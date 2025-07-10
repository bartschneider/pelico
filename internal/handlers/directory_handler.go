package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"pelico/internal/errors"
	"github.com/gin-gonic/gin"
)

type DirectoryHandler struct{
	allowedPaths []string
}

// Common safe directories for ROM browsing
var defaultAllowedPaths = []string{
	"/media",
	"/mnt", 
	"/home",
	"/opt",
	"/data",
	"/storage",
	"/var/lib",
	"/usr/local/share",
}

type DirectoryEntry struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsDir    bool   `json:"is_dir"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
}

func NewDirectoryHandler() *DirectoryHandler {
	return &DirectoryHandler{
		allowedPaths: defaultAllowedPaths,
	}
}

// isPathAllowed checks if the given path is within allowed directories
func (h *DirectoryHandler) isPathAllowed(path string) bool {
	cleanPath := filepath.Clean(path)
	
	// Always allow root for initial navigation
	if cleanPath == "/" {
		return true
	}
	
	for _, allowedPath := range h.allowedPaths {
		// Check if the path starts with an allowed path
		if strings.HasPrefix(cleanPath, allowedPath) {
			return true
		}
	}
	
	return false
}

func (h *DirectoryHandler) BrowseDirectory(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = "/"
	}
	
	// Security: Prevent path traversal attacks
	cleanPath := filepath.Clean(path)
	if !filepath.IsAbs(cleanPath) {
		cleanPath = filepath.Join("/", cleanPath)
	}
	
	// Security: Check if path is allowed
	if !h.isPathAllowed(cleanPath) {
		errors.RespondWithError(c, errors.ErrPermissionDenied, map[string]interface{}{
			"requested_path": cleanPath,
			"allowed_paths": h.allowedPaths,
			"message": "Access restricted to safe ROM directories only",
		})
		return
	}
	
	// Check if directory exists and is accessible
	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			errors.RespondWithError(c, errors.ErrDirectoryNotFound, map[string]string{
				"path": cleanPath,
			})
		} else if os.IsPermission(err) {
			errors.RespondWithError(c, errors.ErrPermissionDenied, map[string]string{
				"path": cleanPath,
				"error": err.Error(),
			})
		} else {
			errors.RespondWithError(c, errors.ErrInvalidDirectory, map[string]string{
				"path": cleanPath,
				"error": err.Error(),
			})
		}
		return
	}
	
	if !info.IsDir() {
		errors.RespondWithError(c, errors.ErrInvalidDirectory, map[string]string{
			"path": cleanPath,
			"reason": "Path is not a directory",
		})
		return
	}
	
	// Read directory contents
	files, err := os.ReadDir(cleanPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read directory: " + err.Error()})
		return
	}
	
	var entries []DirectoryEntry
	
	// Add parent directory entry if not at root
	if cleanPath != "/" {
		parentPath := filepath.Dir(cleanPath)
		entries = append(entries, DirectoryEntry{
			Name:  "..",
			Path:  parentPath,
			IsDir: true,
		})
	}
	
	// Add directory contents
	for _, file := range files {
		fullPath := filepath.Join(cleanPath, file.Name())
		
		// Skip hidden files and directories (starting with .)
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		
		fileInfo, err := file.Info()
		if err != nil {
			continue // Skip files we can't read
		}
		
		entry := DirectoryEntry{
			Name:     file.Name(),
			Path:     fullPath,
			IsDir:    file.IsDir(),
			Size:     fileInfo.Size(),
			Modified: fileInfo.ModTime().Format("2006-01-02 15:04:05"),
		}
		
		entries = append(entries, entry)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"current_path": cleanPath,
		"entries":      entries,
	})
}

func (h *DirectoryHandler) GetSuggestedPaths(c *gin.Context) {
	// Provide common ROM directory suggestions
	suggestions := []string{
		"/home",
		"/media",
		"/mnt",
		"/opt",
		"/usr/local",
		"/var",
	}
	
	// Filter suggestions to only include existing directories
	var validSuggestions []string
	for _, path := range suggestions {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			validSuggestions = append(validSuggestions, path)
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"suggestions": validSuggestions,
	})
}