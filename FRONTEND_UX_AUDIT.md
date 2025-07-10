# Pelico Frontend UX Audit

*Date: January 2025*  
*Auditor: Frontend UX Specialist*  
*Context: Comprehensive review of frontend user experience issues*

## Executive Summary

Despite backend improvements and some frontend fixes, the Pelico frontend suffers from **significant UX problems** that make the application frustrating to use. The core issue is **poor state management** - user actions often reset filters, views don't update consistently, and the interface doesn't respect user preferences.

**Critical UX Issues**:
1. **Filter resets on every major action** - Updating completion status wipes all filters
2. **Play Sessions view is fundamentally broken** - Inefficient data loading, poor performance
3. **State management is fragmented** - No central state, components don't communicate
4. **User preferences are ignored** - Settings reset frequently
5. **Modal actions don't preserve context** - Every modal close triggers full reload

**User Experience Rating**: 4/10 (Functional but frustrating)

## Update: Post-Analysis Findings

Based on comprehensive code analysis, the following critical issues have been confirmed:

### Active Sessions Display Bug - CONFIRMED
**Location**: `web/templates/index.html` line 120
```html
<div id="activeSessionsContainer" style="display: none;">
```
This hardcoded inline style prevents the active sessions from EVER being displayed, regardless of JavaScript attempts to show them. The JavaScript correctly tries to show/hide the container, but inline styles take precedence.

### Filter State Loss - CONFIRMED
**Location**: `web/static/js/app.js` line 2245
The `updateCompletionStatus()` function calls `loadGames()` without parameters, which resets all filters. This explains why users lose their carefully configured filters after updating game completion.

### Sessions View Performance - CONFIRMED CRITICAL
**Location**: `web/static/js/app.js` lines 414-427
The sessions view makes individual API calls for EACH game (up to 50 sequential requests). Additionally, it only loads sessions for games on the current page, meaning users never see sessions for games on other pages.

### Debouncing Not Connected - CONFIRMED
**Location**: `web/static/js/app.js` lines 880-888
The `debouncedFilterGames()` function exists but is never called. All filter controls still call `filterGames()` directly, causing unnecessary API calls on rapid filter changes.

## Detailed UX Analysis

### 1. Filter Reset Problem ⚠️ CRITICAL

**Issue**: When a user updates a game's completion status, ALL filters are reset.

**Code Evidence** (`app.js` line 2245):
```javascript
async function updateCompletionStatus() {
    // ... update completion ...
    
    // Close modal and reload games
    bootstrap.Modal.getInstance(document.getElementById('completionModal')).hide();
    await loadGames();  // ❌ This resets ALL filters!
    
    showToast('Completion status updated successfully!', 'success');
}
```

**User Impact**:
1. User filters by "Nintendo Switch" + "RPG" + "In Progress"
2. User updates one game to "Completed"
3. ALL filters are lost - back to showing all games
4. User must reapply all three filters
5. Pagination also resets to page 1

**Why This Happens**: `loadGames()` doesn't preserve current filter state when called without parameters.

### 2. Play Sessions View - Performance Disaster ⚠️ CRITICAL

**Issue**: The sessions view loads ALL sessions for ALL games sequentially.

**Code Evidence** (`app.js` lines 414-427):
```javascript
async function loadSessions() {
    try {
        // Load sessions for all games
        const allSessions = [];
        for (const game of games) {  // ❌ Only iterates through current 50 games!
            const gameSessions = await apiCall(`/games/${game.id}/sessions`);
            allSessions.push(...gameSessions.map(s => ({...s, game_title: game.title})));
        }
        sessions = allSessions;
        displaySessions(sessions);
    } catch (error) {
        console.error('Failed to load sessions:', error);
    }
}
```

**Problems**:
1. **Only shows sessions for current page** - If you have 200 games but only 50 are loaded, you only see sessions for those 50
2. **Sequential API calls** - 50 games = 50 sequential HTTP requests (could take 10+ seconds)
3. **No loading indicator** - User sees blank screen during load
4. **No caching** - Switching views reloads everything
5. **No pagination** - All sessions shown at once

**User Experience**:
- Click "Play Sessions" → Wait 5-10 seconds → See incomplete data
- Switch away and back → Wait another 5-10 seconds
- Never see sessions for games on other pages

### 3. State Management Chaos

**Issue**: No centralized state management leads to inconsistent behavior.

**Examples of State Fragmentation**:
```javascript
// Global variables scattered throughout (lines 6-16)
let currentView = 'collection';
let currentDisplayView = 'tiled';
let games = [];
let platforms = [];
let sessions = [];
let recentlyPlayedGames = [];
let allGenres = [];
let currentDetailGame = null;
let currentPage = 1;
let totalPages = 1;
let isLoading = false;
let activeSessions = [];  // Duplicate declaration!
```

**Problems**:
1. **No single source of truth** - Filter state lives in DOM, not JavaScript
2. **Components don't communicate** - Updating one part doesn't notify others
3. **Race conditions** - Multiple async operations can conflict
4. **Memory leaks** - Event listeners and intervals not cleaned up

### 4. Modal Workflow Issues

**Every modal action triggers unnecessary reloads**:

1. **Edit Game Modal**:
```javascript
// Line 1395
await loadGames();  // Reloads everything, loses filters
await loadRecentlyPlayedGames();  // Another unnecessary call
```

2. **Add Play Session Modal**:
```javascript
// Lines 1444-1446
await loadRecentlyPlayedGames();
await loadActiveSessions();
// Doesn't update the sessions view if that's active!
```

3. **Delete Game**:
```javascript
// Lines 1856-1858
await loadGames();  // Full reload
await loadRecentlyPlayedGames();  // Another reload
```

### 5. View Switching Problems

**Issue**: Switching between Collection/Sessions/Platforms doesn't preserve state.

**Code Evidence** (`app.js` lines 67-86):
```javascript
function showView(viewName) {
    // ... hide/show views ...
    
    // Load view-specific data
    switch(viewName) {
        case 'collection':
            loadGames();  // ❌ Loses current filters
            loadRecentlyPlayedGames();
            loadActiveSessions();
            loadCompletionStats();
            break;
        case 'sessions':
            loadSessions();  // ❌ Inefficient loading
            break;
        // ...
    }
}
```

**User Impact**:
- Apply filters in Collection view
- Switch to Sessions view
- Switch back to Collection → All filters lost

### 6. Inefficient Data Loading Patterns

**Multiple redundant API calls**:
```javascript
// On page load (lines 19-56)
await loadPlatforms();
await loadGenres();
await loadGames();
await loadActiveSessions();

// Each creates separate API calls instead of batching
```

**No intelligent caching**:
- Platforms rarely change but reload every time
- Genres reload from server despite having the data
- No use of browser storage for semi-static data

### 7. Poor Error Recovery

**Silent failures leave UI in broken state**:
```javascript
// Line 425
} catch (error) {
    console.error('Failed to load sessions:', error);
    // ❌ No user feedback, just empty view
}
```

### 8. Completion Filter Integration Issues

**Completion filters work differently than other filters**:
```javascript
// Line 147 in HTML - "All Games" button has inline logic
onclick="loadGames(); document.querySelectorAll('.completion-filter').forEach(b => b.classList.remove('active')); this.classList.add('active');"

// This bypasses the unified filter system!
```

## Root Cause Analysis

### Why These Issues Exist

1. **Incremental Development Without Refactoring**
   - Features added one by one
   - No overall state architecture
   - Quick fixes instead of systematic solutions

2. **DOM as State Storage**
   - Filter values read from DOM elements
   - No JavaScript state object
   - Makes state preservation difficult

3. **Tight Coupling**
   - UI actions directly call data loading
   - No separation of concerns
   - No event system for updates

4. **No Loading States**
   - Users don't know what's happening
   - No skeleton screens or spinners
   - Leads to confusion and repeated clicks

## User Journey Analysis

### Scenario: Managing Game Collection

1. **User opens app** → Waits for 4 API calls to complete
2. **Filters by "Nintendo Switch"** → Works correctly
3. **Adds genre filter "RPG"** → Works correctly  
4. **Clicks "Playing" status** → Previous filters maintained ✅
5. **Opens game to update completion** → Modal appears
6. **Saves completion status** → ❌ ALL FILTERS RESET
7. **User frustrated** → Must reapply all filters
8. **Switches to Sessions view** → Waits 10 seconds
9. **Goes back to Collection** → ❌ FILTERS LOST AGAIN
10. **User gives up**

### Pain Points
- Every significant action resets user's view preferences
- Long waits with no feedback
- Inconsistent behavior between similar actions
- No way to bookmark or share filtered views (URL updates but doesn't work reliably)

## Recommended Solutions

### Priority 1: Implement Proper State Management (1 week)

```javascript
// Create a central state manager
class AppState {
    constructor() {
        this.filters = {
            platform: '',
            genre: '',
            completionStatus: 'all',
            sort: 'title',
            page: 1,
            limit: 50
        };
        this.view = {
            current: 'collection',
            display: 'tiled'
        };
        this.data = {
            games: [],
            platforms: [],
            genres: [],
            sessions: [],
            activeSessions: []
        };
        this.ui = {
            isLoading: false,
            loadingMessage: ''
        };
    }
    
    // Methods to update state and notify listeners
    updateFilters(newFilters) {
        Object.assign(this.filters, newFilters);
        this.notifyListeners('filters');
    }
    
    preserveAndReload() {
        // Reload data while preserving all current state
        this.loadGames(this.filters);
    }
}
```

### Priority 2: Fix Modal Workflows (3 days)

```javascript
// Update completion without full reload
async function updateCompletionStatus() {
    try {
        const gameId = document.getElementById('completionGameId').value;
        const completionData = { /* ... */ };
        
        const updatedGame = await apiCall(`/games/${gameId}/completion`, 'PUT', completionData);
        
        // Update only the affected game in the current view
        const gameIndex = games.findIndex(g => g.id === updatedGame.id);
        if (gameIndex !== -1) {
            games[gameIndex] = updatedGame;
            updateGameDisplay(updatedGame);  // Update just this game's card
        }
        
        // Update stats without full reload
        await loadCompletionStats();
        
        bootstrap.Modal.getInstance(document.getElementById('completionModal')).hide();
        showToast('Completion status updated successfully!', 'success');
    } catch (error) {
        showToast('Failed to update completion status', 'error');
    }
}
```

### Priority 3: Implement Efficient Sessions View (2 days)

```javascript
// Add a dedicated sessions endpoint
// GET /api/v1/sessions?page=1&limit=50&sort=start_time

// Frontend implementation
async function loadSessions(page = 1) {
    try {
        showLoadingIndicator('Loading play sessions...');
        
        const response = await apiCall(`/sessions?page=${page}&limit=50&sort=start_time`);
        
        sessions = response.sessions;
        totalSessionPages = response.pagination.total_pages;
        
        displaySessions(sessions);
        updateSessionPagination(page, totalSessionPages);
    } catch (error) {
        showError('Failed to load play sessions');
    } finally {
        hideLoadingIndicator();
    }
}
```

### Priority 4: Add Loading States (1 day)

```javascript
// Create reusable loading component
function showLoadingState(container, message = 'Loading...') {
    container.innerHTML = `
        <div class="text-center py-5">
            <div class="spinner-border text-primary mb-3" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
            <p class="text-muted">${message}</p>
        </div>
    `;
}

// Add skeleton screens for better perceived performance
function showGameSkeleton(count = 12) {
    const grid = document.getElementById('gamesGrid');
    grid.innerHTML = Array(count).fill('').map(() => `
        <div class="col-lg-3 col-md-4 col-sm-6 mb-4">
            <div class="card">
                <div class="skeleton-box" style="height: 200px;"></div>
                <div class="card-body">
                    <div class="skeleton-line" style="width: 80%;"></div>
                    <div class="skeleton-line" style="width: 60%;"></div>
                </div>
            </div>
        </div>
    `).join('');
}
```

### Priority 5: Implement Smart Caching (1 week)

```javascript
class DataCache {
    constructor() {
        this.cache = new Map();
        this.maxAge = 5 * 60 * 1000; // 5 minutes
    }
    
    set(key, data) {
        this.cache.set(key, {
            data: data,
            timestamp: Date.now()
        });
    }
    
    get(key) {
        const item = this.cache.get(key);
        if (!item) return null;
        
        if (Date.now() - item.timestamp > this.maxAge) {
            this.cache.delete(key);
            return null;
        }
        
        return item.data;
    }
    
    invalidate(pattern) {
        // Invalidate cache entries matching pattern
        for (const key of this.cache.keys()) {
            if (key.includes(pattern)) {
                this.cache.delete(key);
            }
        }
    }
}

// Use cache for semi-static data
async function loadPlatforms(forceRefresh = false) {
    const cacheKey = 'platforms';
    
    if (!forceRefresh) {
        const cached = dataCache.get(cacheKey);
        if (cached) {
            platforms = cached;
            return;
        }
    }
    
    platforms = await apiCall('/platforms');
    dataCache.set(cacheKey, platforms);
}
```

## Implementation Roadmap

### Phase 1: Critical Fixes (1 week)
1. Implement state preservation for modals
2. Fix "All Games" button to use unified filters
3. Add loading indicators everywhere
4. Cache platforms and genres

### Phase 2: State Management (1 week)
1. Create AppState class
2. Migrate all global variables
3. Implement state persistence
4. Add state change notifications

### Phase 3: Performance (1 week)
1. Implement sessions pagination endpoint
2. Add request debouncing
3. Implement smart caching
4. Add skeleton screens

### Phase 4: Polish (3 days)
1. Add optimistic UI updates
2. Implement undo functionality
3. Add keyboard shortcuts
4. Improve error messages

## Success Metrics

After implementing these fixes, we should see:

1. **Filter Persistence**: 100% of user actions preserve filters
2. **Load Times**: < 1 second for view switches (with cache)
3. **API Calls**: 50% reduction in redundant calls
4. **User Satisfaction**: Target 8/10 experience rating

## Testing Checklist

- [ ] Update completion status → Filters remain
- [ ] Switch views → State preserved
- [ ] Add/Edit/Delete game → Minimal reload
- [ ] Load sessions → Fast with pagination
- [ ] Apply multiple filters → URL updates correctly
- [ ] Refresh page → All state restored
- [ ] Network error → Graceful degradation
- [ ] Rapid clicks → Debounced properly

## Conclusion

The Pelico frontend has solid features but poor execution. The constant filter resets and slow performance create a frustrating user experience. The root cause is the lack of proper state management and the practice of reloading everything after every action.

With the recommended fixes, the application can transform from "functional but frustrating" to "smooth and enjoyable." The key is to respect user intent - **never reset their carefully configured view unless they explicitly ask for it**.

**Current State**: Functional but requires patience and repeated actions  
**Target State**: Responsive, intelligent, and respectful of user preferences

The good news is that all these issues are fixable without major architectural changes. The backend is solid; the frontend just needs to be smarter about how it uses the available APIs.

## Critical Fixes Summary (In Order of Priority)

### 1. Active Sessions Display (1 minute fix)
**File**: `web/templates/index.html` line 120
**Change**: Remove `style="display: none;"` from the activeSessionsContainer div
**Impact**: Active sessions will finally be visible to users

### 2. Filter Preservation on Completion Update (5 minute fix)
**File**: `web/static/js/app.js` line 2245
**Change**: Replace `await loadGames();` with `await applyAllFilters();`
**Impact**: Filters will persist when updating game completion status

### 3. Connect Debouncing (10 minute fix)
**File**: `web/templates/index.html`
**Change**: Replace all `onchange="filterGames()"` with `onchange="debouncedFilterGames()"`
**Impact**: Reduced API calls and better performance during rapid filter changes

### 4. Sessions View Optimization (2 hours)
**Backend**: Create `/api/v1/sessions` endpoint with pagination
**Frontend**: Rewrite loadSessions() to use single paginated API call
**Impact**: From 50+ API calls to 1, massive performance improvement

These four fixes alone would improve the UX rating from 4/10 to 7/10.