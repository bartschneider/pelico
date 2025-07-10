package errors

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// APIError represents a structured API error response
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Status  int         `json:"status"`
}

// Error codes for different types of errors
const (
	// General errors
	ErrInvalidRequest    = "INVALID_REQUEST"
	ErrInternalServer    = "INTERNAL_SERVER_ERROR"
	ErrNotFound          = "NOT_FOUND"
	ErrUnauthorized      = "UNAUTHORIZED"
	ErrForbidden         = "FORBIDDEN"
	
	// Game-specific errors
	ErrGameNotFound      = "GAME_NOT_FOUND"
	ErrGameAlreadyExists = "GAME_ALREADY_EXISTS"
	ErrInvalidGameData   = "INVALID_GAME_DATA"
	
	// Platform-specific errors
	ErrPlatformNotFound      = "PLATFORM_NOT_FOUND"
	ErrPlatformHasGames      = "PLATFORM_HAS_GAMES"
	ErrInvalidPlatformData   = "INVALID_PLATFORM_DATA"
	
	// Session-specific errors
	ErrSessionNotFound       = "SESSION_NOT_FOUND"
	ErrSessionAlreadyEnded   = "SESSION_ALREADY_ENDED"
	ErrInvalidSessionData    = "INVALID_SESSION_DATA"
	
	// Scanner-specific errors
	ErrScanInProgress        = "SCAN_IN_PROGRESS"
	ErrInvalidDirectory      = "INVALID_DIRECTORY"
	ErrDirectoryNotFound     = "DIRECTORY_NOT_FOUND"
	ErrPermissionDenied      = "PERMISSION_DENIED"
	
	// External service errors
	ErrMetadataAPIError      = "METADATA_API_ERROR"
	ErrMetadataNotFound      = "METADATA_NOT_FOUND"
	ErrBackupServiceError    = "BACKUP_SERVICE_ERROR"
	ErrNextcloudError        = "NEXTCLOUD_ERROR"
	
	// Validation errors
	ErrValidationFailed      = "VALIDATION_FAILED"
	ErrMissingRequiredField  = "MISSING_REQUIRED_FIELD"
	ErrInvalidFormat         = "INVALID_FORMAT"
	ErrInvalidRange          = "INVALID_RANGE"
)

// Error message mappings
var ErrorMessages = map[string]string{
	// General errors
	ErrInvalidRequest:    "The request contains invalid parameters",
	ErrInternalServer:    "An internal server error occurred",
	ErrNotFound:          "The requested resource was not found",
	ErrUnauthorized:      "Authentication is required",
	ErrForbidden:         "Access to this resource is forbidden",
	
	// Game-specific errors
	ErrGameNotFound:      "Game not found",
	ErrGameAlreadyExists: "A game with this title already exists",
	ErrInvalidGameData:   "Invalid game data provided",
	
	// Platform-specific errors
	ErrPlatformNotFound:      "Platform not found",
	ErrPlatformHasGames:      "Cannot delete platform with associated games",
	ErrInvalidPlatformData:   "Invalid platform data provided",
	
	// Session-specific errors
	ErrSessionNotFound:       "Play session not found",
	ErrSessionAlreadyEnded:   "This play session has already ended",
	ErrInvalidSessionData:    "Invalid session data provided",
	
	// Scanner-specific errors
	ErrScanInProgress:        "A directory scan is already in progress",
	ErrInvalidDirectory:      "Invalid directory path provided",
	ErrDirectoryNotFound:     "Directory not found or not accessible",
	ErrPermissionDenied:      "Permission denied to access this directory",
	
	// External service errors
	ErrMetadataAPIError:      "Unable to fetch metadata from external service",
	ErrMetadataNotFound:      "No metadata found for this game",
	ErrBackupServiceError:    "Backup service is currently unavailable",
	ErrNextcloudError:        "Failed to connect to Nextcloud service",
	
	// Validation errors
	ErrValidationFailed:      "Request validation failed",
	ErrMissingRequiredField:  "Required field is missing",
	ErrInvalidFormat:         "Invalid format provided",
	ErrInvalidRange:          "Value is outside the allowed range",
}

// NewAPIError creates a new structured API error
func NewAPIError(code string, details ...interface{}) *APIError {
	message, exists := ErrorMessages[code]
	if !exists {
		message = "An unknown error occurred"
	}
	
	err := &APIError{
		Code:    code,
		Message: message,
		Status:  getHTTPStatusForCode(code),
	}
	
	if len(details) > 0 {
		err.Details = details[0]
	}
	
	return err
}

// getHTTPStatusForCode maps error codes to HTTP status codes
func getHTTPStatusForCode(code string) int {
	switch code {
	case ErrNotFound, ErrGameNotFound, ErrPlatformNotFound, ErrSessionNotFound, 
		 ErrDirectoryNotFound, ErrMetadataNotFound:
		return http.StatusNotFound
		
	case ErrInvalidRequest, ErrInvalidGameData, ErrInvalidPlatformData, 
		 ErrInvalidSessionData, ErrInvalidDirectory, ErrValidationFailed,
		 ErrMissingRequiredField, ErrInvalidFormat, ErrInvalidRange,
		 ErrSessionAlreadyEnded, ErrPlatformHasGames:
		return http.StatusBadRequest
		
	case ErrUnauthorized:
		return http.StatusUnauthorized
		
	case ErrForbidden, ErrPermissionDenied:
		return http.StatusForbidden
		
	case ErrScanInProgress:
		return http.StatusConflict
		
	case ErrMetadataAPIError, ErrBackupServiceError, ErrNextcloudError:
		return http.StatusServiceUnavailable
		
	default:
		return http.StatusInternalServerError
	}
}

// RespondWithError sends a structured error response
func RespondWithError(c *gin.Context, code string, details ...interface{}) {
	err := NewAPIError(code, details...)
	c.JSON(err.Status, err)
}

// RespondWithCustomError sends a custom error response
func RespondWithCustomError(c *gin.Context, status int, code, message string, details ...interface{}) {
	err := &APIError{
		Code:    code,
		Message: message,
		Status:  status,
	}
	
	if len(details) > 0 {
		err.Details = details[0]
	}
	
	c.JSON(status, err)
}

// ValidationError represents field-specific validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// RespondWithValidationError sends validation error response
func RespondWithValidationError(c *gin.Context, validationErrors []ValidationError) {
	RespondWithError(c, ErrValidationFailed, map[string]interface{}{
		"validation_errors": validationErrors,
	})
}