package middleware

import (
	"strings"
	"pelico/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidationMiddleware creates a middleware for request validation
func ValidationMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()
	})
}

// ValidateAndBind validates and binds JSON request to struct
func ValidateAndBind(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		var validationErrors []errors.ValidationError
		
		// Check if it's a validation error
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrs {
				validationErrors = append(validationErrors, errors.ValidationError{
					Field:   fieldErr.Field(),
					Message: getValidationMessage(fieldErr),
					Value:   fieldErr.Value(),
				})
			}
			errors.RespondWithValidationError(c, validationErrors)
		} else {
			// It's a binding error (JSON parsing, etc.)
			errors.RespondWithError(c, errors.ErrValidationFailed, map[string]string{
				"binding_error": err.Error(),
			})
		}
		return false
	}
	return true
}

// getValidationMessage returns user-friendly validation messages
func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Must be at least " + fe.Param() + " characters long"
	case "max":
		return "Must be at most " + fe.Param() + " characters long"
	case "len":
		return "Must be exactly " + fe.Param() + " characters long"
	case "numeric":
		return "Must be a number"
	case "alpha":
		return "Must contain only letters"
	case "alphanum":
		return "Must contain only letters and numbers"
	case "url":
		return "Must be a valid URL"
	case "oneof":
		return "Must be one of: " + strings.ReplaceAll(fe.Param(), " ", ", ")
	case "gte":
		return "Must be greater than or equal to " + fe.Param()
	case "lte":
		return "Must be less than or equal to " + fe.Param()
	case "gt":
		return "Must be greater than " + fe.Param()
	case "lt":
		return "Must be less than " + fe.Param()
	default:
		return "Invalid value"
	}
}

// Common validation request structs

// CreateGameRequest represents the request to create a game
type CreateGameRequest struct {
	Title       string  `json:"title" binding:"required,min=1,max=255"`
	PlatformID  uint    `json:"platform_id" binding:"required,gt=0"`
	Year        int     `json:"year" binding:"omitempty,gte=1970,lte=2030"`
	Genre       string  `json:"genre" binding:"omitempty,max=100"`
	Rating      float32 `json:"rating" binding:"omitempty,gte=0,lte=10"`
	Description string  `json:"description" binding:"omitempty,max=2000"`
	CoverArtURL string  `json:"cover_art_url" binding:"omitempty,url"`
}

// UpdateGameRequest represents the request to update a game
type UpdateGameRequest struct {
	Title       string  `json:"title" binding:"omitempty,min=1,max=255"`
	PlatformID  uint    `json:"platform_id" binding:"omitempty,gt=0"`
	Year        int     `json:"year" binding:"omitempty,gte=1970,lte=2030"`
	Genre       string  `json:"genre" binding:"omitempty,max=100"`
	Rating      float32 `json:"rating" binding:"omitempty,gte=0,lte=10"`
	Description string  `json:"description" binding:"omitempty,max=2000"`
	CoverArtURL string  `json:"cover_art_url" binding:"omitempty,url"`
}

// CreatePlatformRequest represents the request to create a platform
type CreatePlatformRequest struct {
	Name         string `json:"name" binding:"required,min=1,max=100"`
	Manufacturer string `json:"manufacturer" binding:"omitempty,max=100"`
	ReleaseYear  int    `json:"release_year" binding:"omitempty,gte=1970,lte=2030"`
}

// UpdatePlatformRequest represents the request to update a platform
type UpdatePlatformRequest struct {
	Name         string `json:"name" binding:"omitempty,min=1,max=100"`
	Manufacturer string `json:"manufacturer" binding:"omitempty,max=100"`
	ReleaseYear  int    `json:"release_year" binding:"omitempty,gte=1970,lte=2030"`
}

// CreateSessionRequest represents the request to create a play session
type CreateSessionRequest struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"omitempty"`
	Notes     string `json:"notes" binding:"omitempty,max=1000"`
	Rating    int    `json:"rating" binding:"omitempty,gte=1,lte=10"`
}

// UpdateSessionRequest represents the request to update a play session
type UpdateSessionRequest struct {
	StartTime string `json:"start_time" binding:"omitempty"`
	EndTime   string `json:"end_time" binding:"omitempty"`
	Notes     string `json:"notes" binding:"omitempty,max=1000"`
	Rating    int    `json:"rating" binding:"omitempty,gte=1,lte=10"`
}

// ScanDirectoryRequest represents the request to scan a directory
type ScanDirectoryRequest struct {
	DirectoryPath  string `json:"directory_path" binding:"required,min=1"`
	ServerLocation string `json:"server_location" binding:"required,min=1,max=100"`
	PlatformID     uint   `json:"platform_id" binding:"required,gt=0"`
	Recursive      bool   `json:"recursive"`
}

// CompletionStatusRequest represents the request to update completion status
type CompletionStatusRequest struct {
	Status     string `json:"status" binding:"required,oneof=not_started in_progress completed abandoned 100_percent"`
	Percentage int    `json:"percentage" binding:"gte=0,lte=100"`
	Notes      string `json:"notes" binding:"omitempty,max=1000"`
}

// SearchGamesRequest represents the request to search for games
type SearchGamesRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=255"`
	Platform string `json:"platform" binding:"omitempty,max=100"`
}

// ListGamesQuery represents query parameters for listing games
type ListGamesQuery struct {
	Page     int    `form:"page" binding:"omitempty,gte=1"`
	Limit    int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
	Sort     string `form:"sort" binding:"omitempty,oneof=title year rating created_at"`
	Platform uint   `form:"platform" binding:"omitempty,gt=0"`
	Genre    string `form:"genre" binding:"omitempty,max=100"`
}