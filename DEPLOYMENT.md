# Pelico Enhanced Deployment Guide

## What's New in This Version

The enhanced Pelico application now includes:

✅ **Edit Game Functionality**: Edit existing games with comprehensive form fields
✅ **Play Session Logging**: Log and track play sessions with ratings and notes  
✅ **Platform & Genre Filters**: Filter your game collection by platform and genre
✅ **Recently Played Games Section**: View recently played games at the top of the collection
✅ **Enhanced UI**: Improved buttons, modals, and user experience

## Features Implemented

### 1. Edit Game Modal
- Comprehensive edit form with all game fields
- Platform selection dropdown
- Rating input (0-10 scale)
- Purchase date tracking
- Cover art URL editing

### 2. Play Session Logging
- Modal interface for logging play sessions
- Start and end time tracking with datetime inputs
- Session rating (1-10 scale)
- Notes field for session details
- Integration with recently played games

### 3. Collection Filtering
- Platform filter dropdown (populated from your platforms)
- Genre filter dropdown (populated from your games)
- Sort options: Title, Year, Rating, Recently Played

### 4. Recently Played Games
- Displays games played in the last 30 days
- Shows last played date
- Compact card layout
- Quick play session logging buttons

## Deployment to Homelab

### Manual Deployment Steps

Since SSH authentication needs to be configured, here are the manual steps:

1. **Transfer the Docker image** to your homelab server:
   ```bash
   # On your homelab server (192.168.1.52)
   docker load < pelico-enhanced.tar
   ```

2. **Stop existing containers**:
   ```bash
   cd ~/pelico
   docker compose down
   ```

3. **Update the source code** (copy these enhanced files):
   - `web/static/js/app.js` - Enhanced JavaScript with new features
   - `web/templates/index.html` - Updated HTML with new modals
   - `internal/handlers/game_handler.go` - Enhanced with recently played endpoint
   - `internal/api/server.go` - Updated API routes

4. **Rebuild and start**:
   ```bash
   docker compose up -d --build
   ```

### Alternative: Direct File Updates

If you prefer to update just the frontend without rebuilding:

1. **Update JavaScript**: Replace the content of `web/static/js/app.js` on your server
2. **Update HTML**: Replace the content of `web/templates/index.html` on your server  
3. **Restart the application**: `docker compose restart pelico`

## API Endpoints Added

- `GET /api/v1/games/recently-played` - Get recently played games
- All existing endpoints remain unchanged

## User Interface Enhancements

### Game Cards
- Added "Edit" button (pencil icon)
- Play session button now opens the logging modal
- Enhanced metadata display

### Modals
- **Edit Game Modal**: Comprehensive editing form
- **Play Session Modal**: Session logging with datetime controls
- **Enhanced Add Game Modal**: Improved search and manual entry

### Filters
- Platform dropdown filter
- Genre dropdown filter  
- Sort by multiple criteria

### Recently Played Section
- Appears at top of collection when you have recent sessions
- Compact game cards with last played dates
- Quick access to log new sessions

## Testing the Features

1. **Add a platform** if you don't have any
2. **Add a game** using the search or manual entry
3. **Edit the game** using the edit button on the game card
4. **Log a play session** using the play button
5. **Test the filters** - filter by platform or genre
6. **Check recently played** - should appear after logging sessions

The enhanced Pelico application maintains all existing functionality while adding these powerful new features for better game collection management!