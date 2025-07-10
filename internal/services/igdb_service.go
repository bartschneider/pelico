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

type IGDBService struct {
	client       *http.Client
	clientID     string
	clientSecret string
	accessToken  string
	tokenExpiry  time.Time
}

type TwitchOAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type IGDBGame struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	Summary      string      `json:"summary"`
	FirstReleaseDate int64   `json:"first_release_date"`
	Rating       float32     `json:"rating"`
	Cover        *IGDBCover  `json:"cover"`
	Genres       []IGDBGenre `json:"genres"`
	Platforms    []IGDBPlatform `json:"platforms"`
}

type IGDBCover struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type IGDBGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type IGDBPlatform struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewIGDBService(clientID, clientSecret string) *IGDBService {
	return &IGDBService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (s *IGDBService) authenticate() error {
	// Check if we have a valid token
	if s.accessToken != "" && time.Now().Before(s.tokenExpiry) {
		return nil
	}

	// Build OAuth URL
	authURL := "https://id.twitch.tv/oauth2/token"
	params := url.Values{}
	params.Add("client_id", s.clientID)
	params.Add("client_secret", s.clientSecret)
	params.Add("grant_type", "client_credentials")

	// Make POST request
	resp, err := s.client.PostForm(authURL, params)
	if err != nil {
		return fmt.Errorf("failed to authenticate with Twitch: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	var oauthResp TwitchOAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&oauthResp); err != nil {
		return fmt.Errorf("failed to decode OAuth response: %v", err)
	}

	// Store token and expiry
	s.accessToken = oauthResp.AccessToken
	s.tokenExpiry = time.Now().Add(time.Duration(oauthResp.ExpiresIn-300) * time.Second) // 5 min buffer

	return nil
}

func (s *IGDBService) SearchGames(title, platform string) ([]GameMetadata, error) {
	if err := s.authenticate(); err != nil {
		return nil, err
	}

	// Build IGDB query with platform filtering
	var query string
	if platform != "" {
		// Map common platform names to IGDB platform IDs for better matching
		platformFilter := s.getPlatformFilter(platform)
		if platformFilter != "" {
			query = fmt.Sprintf(`search "%s"; where platforms = (%s); fields name,summary,first_release_date,rating,cover.url,genres.name,platforms.name; limit 10;`, title, platformFilter)
		} else {
			// If no specific platform mapping, search all and filter results
			query = fmt.Sprintf(`search "%s"; fields name,summary,first_release_date,rating,cover.url,genres.name,platforms.name; limit 20;`, title)
		}
	} else {
		query = fmt.Sprintf(`search "%s"; fields name,summary,first_release_date,rating,cover.url,genres.name,platforms.name; limit 10;`, title)
	}

	// Make request to IGDB
	req, err := http.NewRequest("POST", "https://api.igdb.com/v4/games", strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Client-ID", s.clientID)
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make IGDB request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("IGDB request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var igdbGames []IGDBGame
	if err := json.NewDecoder(resp.Body).Decode(&igdbGames); err != nil {
		return nil, fmt.Errorf("failed to decode IGDB response: %v", err)
	}

	// Convert to our format and filter by platform if needed
	var results []GameMetadata
	for _, game := range igdbGames {
		// If platform filtering was requested but not applied in query, filter results
		if platform != "" && s.getPlatformFilter(platform) == "" {
			if !s.gameMatchesPlatform(game, platform) {
				continue
			}
		}
		metadata := GameMetadata{
			Title:       game.Name,
			Description: game.Summary,
			Rating:      game.Rating / 10.0, // IGDB uses 0-100, we use 0-10
			IGDBID:      game.ID,
		}

		// Parse release date
		if game.FirstReleaseDate > 0 {
			releaseTime := time.Unix(game.FirstReleaseDate, 0)
			metadata.Year = releaseTime.Year()
		}

		// Get genre
		if len(game.Genres) > 0 {
			metadata.Genre = game.Genres[0].Name
		}

		// Get cover art
		if game.Cover != nil && game.Cover.URL != "" {
			// IGDB returns URLs like "//images.igdb.com/igdb/image/upload/..."
			coverURL := game.Cover.URL
			if strings.HasPrefix(coverURL, "//") {
				coverURL = "https:" + coverURL
			}
			// Convert to high-res cover (replace t_thumb with t_cover_big)
			coverURL = strings.Replace(coverURL, "t_thumb", "t_cover_big", 1)
			metadata.CoverArtURL = coverURL
		}

		results = append(results, metadata)
	}

	return results, nil
}

// getPlatformFilter maps platform names to IGDB platform IDs for precise filtering
func (s *IGDBService) getPlatformFilter(platformName string) string {
	platformMap := map[string]string{
		"Nintendo Entertainment System": "18",
		"NES": "18",
		"Super Nintendo Entertainment System": "19",
		"SNES": "19",
		"Nintendo 64": "4",
		"N64": "4",
		"Nintendo GameCube": "21",
		"GameCube": "21",
		"Nintendo Wii": "5",
		"Wii": "5",
		"Game Boy": "33",
		"Game Boy Color": "22",
		"Game Boy Advance": "24",
		"Gameboy Advance": "24",
		"GBA": "24",
		"Nintendo DS": "20",
		"Nintendo 3DS": "37",
		"Sega Master System": "64",
		"Sega Genesis": "29",
		"Genesis": "29",
		"Mega Drive": "29",
		"Sega Saturn": "32",
		"Sega Dreamcast": "23",
		"Dreamcast": "23",
		"Sony PlayStation": "7",
		"PlayStation": "7",
		"PS1": "7",
		"Sony PlayStation 2": "8",
		"PlayStation 2": "8",
		"PS2": "8",
		"Sony PlayStation 3": "9",
		"PlayStation 3": "9",
		"PS3": "9",
		"Sony PlayStation 4": "48",
		"PlayStation 4": "48",
		"PS4": "48",
		"Sony PlayStation 5": "167",
		"PlayStation 5": "167",
		"PS5": "167",
		"Sony PlayStation Portable": "38",
		"PlayStation Portable": "38",
		"PSP": "38",
		"Sony PlayStation Vita": "46",
		"PlayStation Vita": "46",
		"PS Vita": "46",
		"Microsoft Xbox": "11",
		"Xbox": "11",
		"Microsoft Xbox 360": "12",
		"Xbox 360": "12",
		"Microsoft Xbox One": "49",
		"Xbox One": "49",
		"Microsoft Xbox Series X/S": "169",
		"Xbox Series X/S": "169",
		"Xbox Series X": "169",
		"Xbox Series S": "169",
		"Atari 2600": "59",
		"Atari 5200": "66",
		"Atari 7800": "60",
		"PC": "6",
		"Windows": "6",
		"Arcade": "52",
	}
	
	return platformMap[platformName]
}

// gameMatchesPlatform checks if a game is available on the specified platform
func (s *IGDBService) gameMatchesPlatform(game IGDBGame, platformName string) bool {
	platformLower := strings.ToLower(platformName)
	
	for _, platform := range game.Platforms {
		gamePlatformLower := strings.ToLower(platform.Name)
		
		// Check for exact matches or common abbreviations
		if gamePlatformLower == platformLower ||
			strings.Contains(gamePlatformLower, platformLower) ||
			strings.Contains(platformLower, gamePlatformLower) {
			return true
		}
		
		// Check for common platform name variations
		switch platformLower {
		case "nes", "nintendo entertainment system":
			if strings.Contains(gamePlatformLower, "nintendo") && strings.Contains(gamePlatformLower, "entertainment") {
				return true
			}
		case "snes", "super nintendo entertainment system", "super nintendo":
			if strings.Contains(gamePlatformLower, "super nintendo") {
				return true
			}
		case "genesis", "mega drive", "sega genesis":
			if strings.Contains(gamePlatformLower, "genesis") || strings.Contains(gamePlatformLower, "mega drive") {
				return true
			}
		case "playstation", "ps1", "sony playstation":
			if strings.Contains(gamePlatformLower, "playstation") && !strings.Contains(gamePlatformLower, "2") {
				return true
			}
		case "xbox", "microsoft xbox":
			if strings.Contains(gamePlatformLower, "xbox") && !strings.Contains(gamePlatformLower, "360") && !strings.Contains(gamePlatformLower, "one") {
				return true
			}
		}
	}
	
	return false
}

func (s *IGDBService) FetchGameMetadata(title, platform string) (*GameMetadata, error) {
	results, err := s.SearchGames(title, platform)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no game found with title: %s", title)
	}

	// Return the first (best) match
	return &results[0], nil
}