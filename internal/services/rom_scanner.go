package services

import (
	"crypto/md5"
	"fmt"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
	"pelico/internal/models"
	"strings"
)

type ROMScanner struct {
	db *gorm.DB
}

type ScanResult struct {
	FilesFound []string      `json:"files_found"`
	GamesAdded []models.Game `json:"games_added"`
	Errors     []string      `json:"errors"`
}

type DuplicateGroup struct {
	Hash  string                `json:"hash"`
	Files []models.FileLocation `json:"files"`
}

var supportedExtensions = []string{
	".rom", ".bin", ".iso", ".cue", ".img", ".zip", ".7z", ".rar",
	".nes", ".smc", ".sfc", ".gb", ".gbc", ".gba", ".n64", ".z64",
	".psx", ".ps2", ".gcm", ".wad", ".cia", ".3ds",
}

func NewROMScanner(db *gorm.DB) *ROMScanner {
	return &ROMScanner{db: db}
}

func (s *ROMScanner) ScanDirectory(dirPath, serverLocation string, platformID uint, recursive bool) (*ScanResult, error) {
	result := &ScanResult{
		FilesFound: make([]string, 0),
		GamesAdded: make([]models.Game, 0),
		Errors:     make([]string, 0),
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Error accessing %s: %v", path, err))
			return nil
		}

		// Skip directories unless recursive is enabled
		if info.IsDir() {
			if !recursive && path != dirPath {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file has supported extension
		if !s.isSupportedFile(path) {
			return nil
		}

		result.FilesFound = append(result.FilesFound, path)

		// Create file location record
		fileLocation, err := s.createFileLocation(path, serverLocation, info.Size())
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Error processing %s: %v", path, err))
			return nil
		}

		// Try to find existing game or create new one
		game, created, err := s.findOrCreateGame(path, platformID, fileLocation)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Error creating game for %s: %v", path, err))
			return nil
		}

		if created {
			result.GamesAdded = append(result.GamesAdded, *game)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %v", err)
	}

	return result, nil
}

func (s *ROMScanner) isSupportedFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, supportedExt := range supportedExtensions {
		if ext == supportedExt {
			return true
		}
	}
	return false
}

func (s *ROMScanner) createFileLocation(filePath, serverLocation string, fileSize int64) (*models.FileLocation, error) {
	// Calculate file hash
	hash, err := s.calculateFileHash(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %v", err)
	}

	fileLocation := &models.FileLocation{
		ServerLocation: serverLocation,
		FilePath:       filePath,
		FileSize:       fileSize,
		FileHash:       hash,
	}

	return fileLocation, nil
}

func (s *ROMScanner) calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (s *ROMScanner) findOrCreateGame(filePath string, platformID uint, fileLocation *models.FileLocation) (*models.Game, bool, error) {
	// Extract game title from filename
	filename := filepath.Base(filePath)
	title := s.extractGameTitle(filename)

	// Check if game already exists
	var existingGame models.Game
	result := s.db.Where("title = ? AND platform_id = ?", title, platformID).First(&existingGame)

	if result.Error == nil {
		// Game exists, add file location
		fileLocation.GameID = existingGame.ID
		if err := s.db.Create(fileLocation).Error; err != nil {
			return nil, false, err
		}
		return &existingGame, false, nil
	}

	if result.Error != gorm.ErrRecordNotFound {
		return nil, false, result.Error
	}

	// Create new game
	game := &models.Game{
		Title:      title,
		PlatformID: platformID,
	}

	if err := s.db.Create(game).Error; err != nil {
		return nil, false, err
	}

	// Add file location
	fileLocation.GameID = game.ID
	if err := s.db.Create(fileLocation).Error; err != nil {
		return nil, false, err
	}

	return game, true, nil
}

func (s *ROMScanner) extractGameTitle(filename string) string {
	// Remove extension
	title := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Remove common ROM tags and brackets
	title = strings.ReplaceAll(title, "_", " ")
	title = strings.ReplaceAll(title, "-", " ")

	// Remove content in parentheses and brackets (region codes, versions, etc.)
	title = removeTagsInBrackets(title, '(', ')')
	title = removeTagsInBrackets(title, '[', ']')

	// Clean up multiple spaces
	words := strings.Fields(title)
	return strings.Join(words, " ")
}

func removeTagsInBrackets(text string, open, close rune) string {
	var result strings.Builder
	depth := 0

	for _, char := range text {
		if char == open {
			depth++
		} else if char == close {
			depth--
		} else if depth == 0 {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func (s *ROMScanner) FindDuplicates() ([]DuplicateGroup, error) {
	var fileLocations []models.FileLocation
	result := s.db.Preload("Game").Find(&fileLocations)
	if result.Error != nil {
		return nil, result.Error
	}

	// Group files by hash
	hashGroups := make(map[string][]models.FileLocation)
	for _, file := range fileLocations {
		if file.FileHash != "" {
			hashGroups[file.FileHash] = append(hashGroups[file.FileHash], file)
		}
	}

	// Find groups with more than one file
	var duplicates []DuplicateGroup
	for hash, files := range hashGroups {
		if len(files) > 1 {
			duplicates = append(duplicates, DuplicateGroup{
				Hash:  hash,
				Files: files,
			})
		}
	}

	return duplicates, nil
}
