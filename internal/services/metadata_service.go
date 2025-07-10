package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type MetadataService struct {
	client      *http.Client
	igdbService *IGDBService
}

type GameMetadata struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Rating      float32 `json:"rating"`
	Genre       string  `json:"genre"`
	Year        int     `json:"year"`
	CoverArtURL string  `json:"cover_art_url"`
	BoxArtURL   string  `json:"box_art_url"`
	IGDBID      int     `json:"igdb_id"`
}

// TheGamesDB API structures
type TheGamesDBResponse struct {
	Data struct {
		Games []TheGamesDBGame `json:"games"`
	} `json:"data"`
	Include struct {
		Boxart map[string]TheGamesDBBoxart `json:"boxart"`
		Genres map[string]TheGamesDBGenre  `json:"genres"`
	} `json:"include"`
}

type TheGamesDBGame struct {
	ID          int    `json:"id"`
	GameTitle   string `json:"game_title"`
	ReleaseDate string `json:"release_date"`
	Overview    string `json:"overview"`
	Genres      []int  `json:"genres"`
	Rating      string `json:"rating"`
}

type TheGamesDBBoxart struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Side     string `json:"side"`
	Filename string `json:"filename"`
}

type TheGamesDBGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// RAWG API structures
type RAWGResponse struct {
	Count   int        `json:"count"`
	Results []RAWGGame `json:"results"`
}

type RAWGGame struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	Released       string        `json:"released"`
	BackgroundImage string       `json:"background_image"`
	Rating         float32       `json:"rating"`
	RatingTop      int           `json:"rating_top"`
	Metacritic     int           `json:"metacritic"`
	Genres         []RAWGGenre   `json:"genres"`
	Platforms      []RAWGPlatform `json:"platforms"`
	Description    string        `json:"description_raw"`
}

type RAWGGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RAWGPlatform struct {
	Platform struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"platform"`
}

func NewMetadataService(clientID, clientSecret string) *MetadataService {
	return &MetadataService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		igdbService: NewIGDBService(clientID, clientSecret),
	}
}

func (s *MetadataService) FetchGameMetadata(title, platform string) (*GameMetadata, error) {
	// Use IGDB if available
	if s.igdbService != nil {
		return s.igdbService.FetchGameMetadata(title, platform)
	}
	
	return nil, fmt.Errorf("no metadata service configured")
}

func (s *MetadataService) SearchGames(title, platform string) ([]GameMetadata, error) {
	// Use IGDB if available
	if s.igdbService != nil {
		return s.igdbService.SearchGames(title, platform)
	}
	
	return nil, fmt.Errorf("no metadata service configured")
}

func (s *MetadataService) fetchFromTheGamesDB(title, platform string) (*GameMetadata, error) {
	// TheGamesDB API endpoint
	baseURL := "https://api.thegamesdb.net/v1/Games/ByGameName"
	
	// Build query parameters
	params := url.Values{}
	params.Add("name", title)
	params.Add("fields", "players,publishers,genres,overview,last_updated,rating,platform")
	params.Add("include", "boxart,genres")
	
	// Make request
	resp, err := s.client.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}
	
	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	
	var apiResp TheGamesDBResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	
	// Find best match
	if len(apiResp.Data.Games) == 0 {
		return nil, fmt.Errorf("no games found for title: %s", title)
	}
	
	game := apiResp.Data.Games[0] // Take first result
	
	// Convert to our metadata format
	metadata := &GameMetadata{
		Title:       game.GameTitle,
		Description: game.Overview,
		IGDBID:      game.ID,
	}
	
	// Parse release year
	if game.ReleaseDate != "" {
		if releaseTime, err := time.Parse("2006-01-02", game.ReleaseDate); err == nil {
			metadata.Year = releaseTime.Year()
		}
	}
	
	// Get genre from TheGamesDB genres
	if len(game.Genres) > 0 {
		genreID := game.Genres[0] // Take first genre
		if genre, exists := apiResp.Include.Genres[fmt.Sprintf("%d", genreID)]; exists {
			metadata.Genre = genre.Name
		}
	}
	
	// Find cover art
	for _, boxart := range apiResp.Include.Boxart {
		if boxart.Type == "boxart" && boxart.Side == "front" {
			metadata.CoverArtURL = "https://cdn.thegamesdb.net/images/thumb/" + boxart.Filename
			break
		}
	}
	
	// Set default rating if not available
	if metadata.Rating == 0 {
		metadata.Rating = 7.5
	}
	
	return metadata, nil
}



func (s *MetadataService) searchTheGamesDB(title, platform string) ([]GameMetadata, error) {
	// Enhanced TheGamesDB search with multiple results
	baseURL := "https://api.thegamesdb.net/v1/Games/ByGameName"
	
	params := url.Values{}
	params.Add("name", title)
	params.Add("fields", "players,publishers,genres,overview,last_updated,rating,platform")
	params.Add("include", "boxart,genres")
	
	resp, err := s.client.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	
	var apiResp TheGamesDBResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	
	var results []GameMetadata
	for _, game := range apiResp.Data.Games {
		metadata := GameMetadata{
			Title:       game.GameTitle,
			Description: game.Overview,
			IGDBID:      game.ID,
		}
		
		// Parse release year
		if game.ReleaseDate != "" {
			if releaseTime, err := time.Parse("2006-01-02", game.ReleaseDate); err == nil {
				metadata.Year = releaseTime.Year()
			}
		}
		
		// Get genre from TheGamesDB genres
		if len(game.Genres) > 0 {
			genreID := game.Genres[0] // Take first genre
			if genre, exists := apiResp.Include.Genres[fmt.Sprintf("%d", genreID)]; exists {
				metadata.Genre = genre.Name
			}
		}
		
		// Find cover art
		for _, boxart := range apiResp.Include.Boxart {
			if boxart.Type == "boxart" && boxart.Side == "front" {
				metadata.CoverArtURL = "https://cdn.thegamesdb.net/images/thumb/" + boxart.Filename
				break
			}
		}
		
		// Set default rating if not available
		if metadata.Rating == 0 {
			metadata.Rating = 7.5 // Default rating
		}
		
		results = append(results, metadata)
	}
	
	return results, nil
}






// Future: Implement IGDB integration
func (s *MetadataService) fetchFromIGDB(title, platform, apiKey string) (*GameMetadata, error) {
	// TODO: Implement IGDB API integration
	// Requires API key and OAuth token
	return nil, fmt.Errorf("IGDB integration not implemented yet")
}