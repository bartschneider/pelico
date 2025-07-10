package models

import (
	"time"
	"gorm.io/gorm"
)

type Platform struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Name         string `json:"name" gorm:"unique;not null"`
	Manufacturer string `json:"manufacturer"`
	ReleaseYear  int    `json:"release_year"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Game struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	PlatformID  uint      `json:"platform_id"`
	Platform    Platform  `json:"platform" gorm:"foreignKey:PlatformID"`
	Year        int       `json:"year"`
	Genre       string    `json:"genre"`
	Rating      float32   `json:"rating"`
	Description string    `json:"description" gorm:"type:text"`
	CoverArtURL string    `json:"cover_art_url"`
	BoxArtURL   string    `json:"box_art_url"`
	PurchaseDate *time.Time `json:"purchase_date"`
	IGDBID      int       `json:"igdb_id"`
	
	// Completion tracking
	CompletionStatus     string     `json:"completion_status" gorm:"default:not_started"`
	CompletionDate       *time.Time `json:"completion_date"`
	CompletionPercentage int        `json:"completion_percentage" gorm:"default:0"`
	CompletionNotes      string     `json:"completion_notes" gorm:"type:text"`
	
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	FileLocations []FileLocation `json:"file_locations" gorm:"foreignKey:GameID"`
	PlaySessions  []PlaySession  `json:"play_sessions" gorm:"foreignKey:GameID"`
}

type FileLocation struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	GameID         uint   `json:"game_id"`
	Game           Game   `json:"game" gorm:"foreignKey:GameID"`
	ServerLocation string `json:"server_location"`
	FilePath       string `json:"file_path" gorm:"not null"`
	FileSize       int64  `json:"file_size"`
	FileHash       string `json:"file_hash"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PlaySession struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	GameID    uint       `json:"game_id"`
	Game      Game       `json:"game" gorm:"foreignKey:GameID"`
	StartTime time.Time  `json:"start_time" gorm:"not null"`
	EndTime   *time.Time `json:"end_time"`
	Duration  int        `json:"duration"` // in minutes
	Notes     string     `json:"notes" gorm:"type:text"`
	Rating    int        `json:"rating"` // 1-10 scale
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Platform{}, &Game{}, &FileLocation{}, &PlaySession{})
}