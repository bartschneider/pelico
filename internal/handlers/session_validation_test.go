package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"pelico/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Missing start_time field",
			requestBody: map[string]interface{}{
				"notes":  "Test session",
				"rating": 5,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "StartTime",
		},
		{
			name: "Empty start_time string",
			requestBody: map[string]interface{}{
				"start_time": "",
				"notes":      "Test session",
				"rating":     5,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "StartTime",
		},
		{
			name: "Valid start_time",
			requestBody: map[string]interface{}{
				"start_time": time.Now().Format(time.RFC3339),
				"notes":      "Test session",
				"rating":     5,
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Invalid rating (too high)",
			requestBody: map[string]interface{}{
				"start_time": time.Now().Format(time.RFC3339),
				"rating":     11,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Rating",
		},
		{
			name: "Invalid rating (too low)",
			requestBody: map[string]interface{}{
				"start_time": time.Now().Format(time.RFC3339),
				"rating":     0,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Rating",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			
			// Simple test endpoint that validates the request
			router.POST("/test", func(c *gin.Context) {
				var req middleware.CreateSessionRequest
				if !middleware.ValidateAndBind(c, &req) {
					return
				}
				c.JSON(http.StatusOK, gin.H{"success": true})
			})
			
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				
				responseStr := string(w.Body.Bytes())
				assert.Contains(t, responseStr, tt.expectedError)
			}
		})
	}
}

// Test demonstrates the exact error seen in the bug report
func TestReproduceBugReport(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.POST("/api/v1/games/:id/sessions", func(c *gin.Context) {
		var req middleware.CreateSessionRequest
		if !middleware.ValidateAndBind(c, &req) {
			return
		}
		c.JSON(http.StatusCreated, gin.H{"success": true})
	})
	
	// This simulates what happens when the frontend sends an empty request
	// or a request without the start_time field
	reqBody := map[string]interface{}{
		// start_time is missing or empty
	}
	
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/games/1/sessions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Should get 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Check the exact error structure
	assert.Equal(t, "VALIDATION_FAILED", response["code"])
	assert.Equal(t, "Request validation failed", response["message"])
	assert.Equal(t, float64(400), response["status"])
	
	// Check validation error details
	details := response["details"].(map[string]interface{})
	validationErrors := details["validation_errors"].([]interface{})
	require.Len(t, validationErrors, 1)
	
	firstError := validationErrors[0].(map[string]interface{})
	assert.Equal(t, "StartTime", firstError["field"])
	assert.Equal(t, "This field is required", firstError["message"])
	assert.Equal(t, "", firstError["value"])
	
	// This matches the exact error from the bug report:
	// {"code":"VALIDATION_FAILED","message":"Request validation failed","details":{"validation_errors":[{"field":"StartTime","message":"This field is required","value":""}]},"status":400}
}