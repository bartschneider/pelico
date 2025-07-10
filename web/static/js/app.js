// Pelico Frontend JavaScript

const API_BASE = '/api/v1';

// Modal Management System
class ModalManager {
    constructor() {
        this.activeModals = [];
    }
    
    async openModal(modalId, setupFunc = null) {
        // Close any currently open modals first
        await this.closeAll();
        
        const modalElement = document.getElementById(modalId);
        if (!modalElement) {
            console.error(`Modal ${modalId} not found`);
            return null;
        }
        
        const modal = new bootstrap.Modal(modalElement);
        this.activeModals.push({ id: modalId, instance: modal, element: modalElement });
        
        // Run setup function if provided
        if (setupFunc) {
            try {
                await setupFunc();
            } catch (error) {
                console.error('Modal setup failed:', error);
            }
        }
        
        modal.show();
        return modal;
    }
    
    async closeAll() {
        for (const modalData of this.activeModals) {
            modalData.instance.hide();
            await this.waitForClose(modalData.instance);
        }
        this.activeModals = [];
    }
    
    waitForClose(modal) {
        return new Promise(resolve => {
            const element = modal._element;
            const handler = () => {
                element.removeEventListener('hidden.bs.modal', handler);
                resolve();
            };
            element.addEventListener('hidden.bs.modal', handler);
        });
    }
    
    getCurrentModal() {
        return this.activeModals.length > 0 ? this.activeModals[this.activeModals.length - 1] : null;
    }
}

// Global modal manager instance
const modalManager = new ModalManager();

// Request Management System
class RequestManager {
    constructor() {
        this.activeRequests = new Map();
    }
    
    async fetch(key, url, options = {}) {
        // Cancel any existing request with same key
        this.cancel(key);
        
        const controller = new AbortController();
        this.activeRequests.set(key, controller);
        
        try {
            const response = await fetch(url, {
                ...options,
                signal: controller.signal
            });
            
            this.activeRequests.delete(key);
            return response;
        } catch (error) {
            this.activeRequests.delete(key);
            if (error.name === 'AbortError') {
                console.log(`Request ${key} was cancelled`);
                return null;
            }
            throw error;
        }
    }
    
    cancel(key) {
        const controller = this.activeRequests.get(key);
        if (controller) {
            controller.abort();
            this.activeRequests.delete(key);
        }
    }
    
    cancelAll() {
        for (const [key, controller] of this.activeRequests) {
            controller.abort();
        }
        this.activeRequests.clear();
    }
}

// Global request manager instance
const requestManager = new RequestManager();

// Global state
let currentView = 'collection';
let currentDisplayView = 'tiled'; // 'tiled' or 'list'
let games = [];
let platforms = [];
let sessions = [];
let recentlyPlayedGames = [];
let allGenres = [];
let currentDetailGame = null;
let currentPage = 1;
let totalPages = 1;
let isLoading = false;

// Initialize app
document.addEventListener('DOMContentLoaded', async function() {
    // Load state from URL first
    const params = new URLSearchParams(window.location.search);
    
    await loadPlatforms();
    await loadGenres();
    
    // Restore filters from URL parameters
    if (params.get('platform')) {
        const platformFilter = document.getElementById('platformFilter');
        if (platformFilter) {
            const platformValue = params.get('platform');
            // Try to match by ID first, then fallback to name for backward compatibility
            const matchingOption = Array.from(platformFilter.options).find(opt => 
                opt.value === platformValue || opt.dataset.name === platformValue
            );
            if (matchingOption) {
                platformFilter.value = matchingOption.value;
            } else {
                // Direct assignment for backward compatibility
                platformFilter.value = platformValue;
            }
        }
    }
    if (params.get('genre')) {
        const genreFilter = document.getElementById('genreFilter');
        if (genreFilter) {
            genreFilter.value = params.get('genre');
        }
    }
    if (params.get('completion_status')) {
        const completionButton = document.querySelector(`[data-status="${params.get('completion_status')}"]`);
        if (completionButton) {
            // Remove active from all buttons first
            document.querySelectorAll('.completion-filter').forEach(btn => btn.classList.remove('active'));
            completionButton.classList.add('active');
        }
    }
    
    // Load games with restored filters
    await loadGames();
    await loadActiveSessions();
    showView('collection');
    
    // Restore preferred view
    const preferredView = localStorage.getItem('preferredView') || 'tiled';
    switchView(preferredView);
});

// View management
function showView(viewName) {
    // Hide all views
    document.querySelectorAll('.view-content').forEach(view => {
        view.style.display = 'none';
    });
    
    // Show selected view
    document.getElementById(viewName + 'View').style.display = 'block';
    
    // Update nav
    document.querySelectorAll('.nav-link').forEach(link => {
        link.classList.remove('active');
    });
    event.target.classList.add('active');
    
    currentView = viewName;
    
    // Load view-specific data
    switch(viewName) {
        case 'collection':
            loadGames();
            loadRecentlyPlayedGames();
            loadActiveSessions();
            loadCompletionStats();
            break;
        case 'scanner':
            loadPlatformsForScanner();
            break;
        case 'sessions':
            loadSessions();
            break;
        case 'platforms':
            loadPlatformsView();
            break;
        case 'settings':
            loadSettingsView();
            break;
    }
}

// API calls
async function apiCall(endpoint, method = 'GET', data = null) {
    try {
        const options = {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            }
        };
        
        if (data) {
            options.body = JSON.stringify(data);
        }
        
        const response = await fetch(API_BASE + endpoint, options);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error('API call failed:', error);
        showToast('API Error: ' + error.message, 'error');
        throw error;
    }
}

// Load functions
async function loadGames(page = 1, append = false) {
    if (isLoading) return;
    
    try {
        isLoading = true;
        showLoadingIndicator();
        
        // Check if we have active filters and should use them
        const platformFilter = document.getElementById('platformFilter')?.value;
        const genreFilter = document.getElementById('genreFilter')?.value;
        const completionFilter = document.querySelector('.completion-filter.active')?.dataset.status;
        const sortBy = document.getElementById('sortBy')?.value || 'title';
        
        // Build API query parameters
        const params = new URLSearchParams({
            page: page,
            limit: 50,
            sort: sortBy
        });
        
        if (platformFilter && platformFilter !== '') {
            params.append('platform', platformFilter);
        }
        if (genreFilter && genreFilter !== '') {
            params.append('genre', genreFilter);
        }
        if (completionFilter && completionFilter !== 'all') {
            params.append('completion_status', completionFilter);
        }
        
        const response = await apiCall(`/games?${params.toString()}`);
        
        if (append) {
            games = [...games, ...response.games];
        } else {
            games = response.games;
        }
        
        currentPage = response.pagination.page;
        totalPages = response.pagination.total_pages;
        
        extractGenresFromCurrentGames();
        updateFilterOptions(true); // Preserve current selections
        displayGames(games);
        updatePaginationControls();
    } catch (error) {
        console.error('Failed to load games:', error);
        showToast('Failed to load games', 'error');
    } finally {
        isLoading = false;
        hideLoadingIndicator();
    }
}

function showLoadingIndicator() {
    // Add loading indicator if it doesn't exist
    if (!document.getElementById('loadingIndicator')) {
        const indicator = document.createElement('div');
        indicator.id = 'loadingIndicator';
        indicator.className = 'text-center py-3';
        indicator.innerHTML = '<i class="fas fa-spinner fa-spin fa-2x text-primary"></i>';
        
        const grid = document.getElementById('gamesGrid');
        grid.appendChild(indicator);
    }
}

function hideLoadingIndicator() {
    const indicator = document.getElementById('loadingIndicator');
    if (indicator) {
        indicator.remove();
    }
}

async function loadRecentlyPlayedGames() {
    try {
        recentlyPlayedGames = await apiCall('/games/recently-played');
        displayRecentlyPlayedGames(recentlyPlayedGames);
    } catch (error) {
        console.error('Failed to load recently played games:', error);
        // Hide section if no recently played games
        document.getElementById('recentlyPlayedSection').style.display = 'none';
    }
}

let activeSessions = [];
let activeSessionsRefreshInterval;

// Debouncing for performance optimization
let filterDebounceTimeout;
const FILTER_DEBOUNCE_DELAY = 300; // 300ms delay

async function loadActiveSessions() {
    try {
        console.log('Loading active sessions...');
        activeSessions = await apiCall('/sessions/active');
        
        // Validate response
        if (!Array.isArray(activeSessions)) {
            console.warn('Active sessions response is not an array, defaulting to empty array');
            activeSessions = [];
        }
        
        displayActiveSessions();
        
        // Set up auto-refresh if there are active sessions
        if (activeSessions.length > 0) {
            if (activeSessionsRefreshInterval) {
                clearInterval(activeSessionsRefreshInterval);
            }
            activeSessionsRefreshInterval = setInterval(loadActiveSessions, 60000); // Refresh every minute
            console.log(`Auto-refresh set up for ${activeSessions.length} active sessions`);
        } else {
            if (activeSessionsRefreshInterval) {
                clearInterval(activeSessionsRefreshInterval);
                activeSessionsRefreshInterval = null;
                console.log('No active sessions, auto-refresh disabled');
            }
        }
    } catch (error) {
        console.error('Failed to load active sessions:', error);
        // Don't hide the container immediately - might be a temporary network issue
        // showToast('Unable to load active sessions', 'warning');
        
        // Set empty array and display (which will hide container)
        activeSessions = [];
        displayActiveSessions();
    }
}

function displayActiveSessions() {
    const container = document.getElementById('activeSessionsContainer');
    const content = document.getElementById('activeSessionsContent');
    
    // Safety checks for DOM elements
    if (!container || !content) {
        console.warn('Active sessions container or content not found in DOM');
        return;
    }
    
    if (!activeSessions || activeSessions.length === 0) {
        container.style.display = 'none';
        container.classList.add('d-none');
        return;
    }
    
    // ALWAYS show container if we have sessions
    container.style.display = 'block';
    container.classList.remove('d-none');
    
    content.innerHTML = activeSessions.map(session => `
        <div class="col-md-6 mb-2">
            <div class="card bg-primary text-white">
                <div class="card-body p-3">
                    <div class="d-flex justify-content-between align-items-center">
                        <div>
                            <h6 class="mb-1">${session.game?.title || 'Unknown Game'}</h6>
                            <small>Started: ${new Date(session.start_time).toLocaleTimeString()}</small>
                            <div class="mt-1">
                                <i class="fas fa-clock"></i> 
                                ${formatDuration(session.current_duration || 0)}
                            </div>
                            ${session.game?.platform ? `<small><i class="fas fa-desktop"></i> ${session.game.platform.name}</small>` : ''}
                        </div>
                        <div>
                            <button class="btn btn-light btn-sm" 
                                    onclick="endActiveSession(${session.id})"
                                    title="End Session">
                                <i class="fas fa-stop"></i> End
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `).join('');
}

async function endActiveSession(sessionId) {
    try {
        await apiCall(`/sessions/${sessionId}/end`, 'POST');
        await loadActiveSessions();
        await loadRecentlyPlayedGames();
        showToast('Session ended successfully', 'success');
    } catch (error) {
        console.error('Failed to end session:', error);
        showToast('Failed to end session', 'error');
    }
}

function formatDuration(minutes) {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    
    if (hours > 0) {
        return `${hours}h ${mins}m`;
    }
    return `${mins}m`;
}

async function loadGenres() {
    try {
        allGenres = await apiCall('/games/genres');
    } catch (error) {
        console.error('Failed to load genres:', error);
        // Fallback to extracting from current games
        extractGenresFromCurrentGames();
    }
}

function extractGenresFromCurrentGames() {
    const genres = new Set();
    games.forEach(game => {
        if (game.genre) {
            genres.add(game.genre);
        }
    });
    allGenres = Array.from(genres).sort();
}

function updateFilterOptions(preserveSelection = false) {
    // Save current selections if we need to preserve them
    const currentPlatform = preserveSelection ? document.getElementById('platformFilter')?.value : '';
    const currentGenre = preserveSelection ? document.getElementById('genreFilter')?.value : '';
    
    // Update platform filter
    const platformFilter = document.getElementById('platformFilter');
    if (platformFilter) {
        platformFilter.innerHTML = '<option value="">All Platforms</option>';
        platforms.forEach(platform => {
            const option = document.createElement('option');
            option.value = platform.id;  // Use ID for robust identification
            option.textContent = platform.name;
            option.dataset.name = platform.name;  // Store name for backward compatibility
            platformFilter.appendChild(option);
        });
        
        // Restore selection if preserving
        if (preserveSelection && currentPlatform) {
            // Try to match by ID first, then fallback to name for backward compatibility
            const matchingOption = Array.from(platformFilter.options).find(opt => 
                opt.value === currentPlatform || opt.dataset.name === currentPlatform
            );
            if (matchingOption) {
                platformFilter.value = matchingOption.value;
            }
        }
    }
    
    // Update genre filter
    const genreFilter = document.getElementById('genreFilter');
    if (genreFilter) {
        genreFilter.innerHTML = '<option value="">All Genres</option>';
        allGenres.forEach(genre => {
            const option = document.createElement('option');
            option.value = genre;
            option.textContent = genre;
            genreFilter.appendChild(option);
        });
        
        // Restore selection if preserving
        if (preserveSelection && currentGenre) {
            genreFilter.value = currentGenre;
        }
    }
}

async function loadPlatforms() {
    try {
        platforms = await apiCall('/platforms');
        updatePlatformSelects();
    } catch (error) {
        console.error('Failed to load platforms:', error);
    }
}

async function loadSessions() {
    try {
        // Load sessions for all games
        const allSessions = [];
        for (const game of games) {
            const gameSessions = await apiCall(`/games/${game.id}/sessions`);
            allSessions.push(...gameSessions.map(s => ({...s, game_title: game.title})));
        }
        sessions = allSessions;
        displaySessions(sessions);
    } catch (error) {
        console.error('Failed to load sessions:', error);
    }
}

// View switching
function switchView(viewType) {
    currentDisplayView = viewType;
    
    // Update button states
    document.getElementById('tiledViewBtn').classList.toggle('active', viewType === 'tiled');
    document.getElementById('listViewBtn').classList.toggle('active', viewType === 'list');
    
    // Show/hide appropriate container using Bootstrap classes
    const grid = document.getElementById('gamesGrid');
    const list = document.getElementById('gamesList');
    
    if (viewType === 'tiled') {
        // Use Bootstrap utility classes to maintain flexbox grid
        grid.classList.remove('d-none');
        grid.classList.add('d-flex', 'flex-wrap');
        list.classList.add('d-none');
        list.classList.remove('d-flex', 'flex-wrap');
    } else {
        // Switch to list view
        grid.classList.add('d-none');
        grid.classList.remove('d-flex', 'flex-wrap');
        list.classList.remove('d-none');
        list.classList.add('d-flex', 'flex-column');
    }
    
    // Store view preference
    localStorage.setItem('preferredView', viewType);
    
    // Redisplay games in new view
    displayGames(games);
}

// Collection Format Functions
function renderCollectionFormats(game) {
    if (!game.collection_formats || !Array.isArray(game.collection_formats) || game.collection_formats.length === 0) {
        return '';
    }
    
    const formatIcons = {
        'physical': { icon: 'fa-box', color: 'success', label: 'Physical' },
        'digital': { icon: 'fa-download', color: 'info', label: 'Digital' },
        'rom': { icon: 'fa-hdd', color: 'warning', label: 'ROM' }
    };
    
    return game.collection_formats.map(format => {
        const config = formatIcons[format];
        if (!config) return '';
        
        return `<span class="badge bg-${config.color} me-1" title="${config.label}">
            <i class="fas ${config.icon}"></i>
        </span>`;
    }).join('');
}

function renderCollectionFormatsDetailed(game) {
    if (!game.collection_formats || !Array.isArray(game.collection_formats) || game.collection_formats.length === 0) {
        return '<span class="text-muted">No format specified</span>';
    }
    
    const formatIcons = {
        'physical': { icon: 'fa-box', color: 'success', label: 'Physical Copy' },
        'digital': { icon: 'fa-download', color: 'info', label: 'Digital Purchase' },
        'rom': { icon: 'fa-hdd', color: 'warning', label: 'ROM File' }
    };
    
    return game.collection_formats.map(format => {
        const config = formatIcons[format];
        if (!config) return '';
        
        return `<span class="badge bg-${config.color} me-1">
            <i class="fas ${config.icon} me-1"></i>${config.label}
        </span>`;
    }).join('');
}

// Completion Status Functions
function renderCompletionStatus(game) {
    const statuses = {
        'not_started': { icon: 'fa-circle', color: 'secondary', label: 'Not Started' },
        'in_progress': { icon: 'fa-spinner', color: 'primary', label: 'Playing' },
        'completed': { icon: 'fa-check-circle', color: 'success', label: 'Completed' },
        'abandoned': { icon: 'fa-times-circle', color: 'warning', label: 'Abandoned' },
        '100_percent': { icon: 'fa-trophy', color: 'warning', label: '100%', style: 'color: gold;' }
    };
    
    const status = statuses[game.completion_status || 'not_started'];
    const percentage = game.completion_percentage || 0;
    
    return `
        <div class="completion-status position-absolute top-0 start-0 m-2" 
             onclick="event.stopPropagation(); showCompletionModal(${game.id})" 
             title="Click to update completion status"
             style="cursor: pointer;">
            <span class="badge bg-${status.color}" ${status.style ? `style="${status.style}"` : ''}>
                <i class="fas ${status.icon}"></i> ${status.label}
                ${percentage > 0 ? ` (${percentage}%)` : ''}
            </span>
        </div>
    `;
}

function renderCompletionStatusBadge(game) {
    const statuses = {
        'not_started': { icon: 'fa-circle', color: 'secondary', label: 'Not Started' },
        'in_progress': { icon: 'fa-spinner', color: 'primary', label: 'Playing' },
        'completed': { icon: 'fa-check-circle', color: 'success', label: 'Completed' },
        'abandoned': { icon: 'fa-times-circle', color: 'warning', label: 'Abandoned' },
        '100_percent': { icon: 'fa-trophy', color: 'warning', label: '100%', style: 'color: gold; background-color: gold; color: black;' }
    };
    
    const status = statuses[game.completion_status || 'not_started'];
    const percentage = game.completion_percentage || 0;
    
    return `<span class="badge bg-${status.color}" ${status.style ? `style="${status.style}"` : ''}>
        <i class="fas ${status.icon}"></i> ${status.label}
        ${percentage > 0 ? ` ${percentage}%` : ''}
    </span>`;
}

// Display functions
function displayGames(gamesList) {
    if (currentDisplayView === 'tiled') {
        displayGamesTiled(gamesList);
    } else {
        displayGamesList(gamesList);
    }
}

function displayGamesTiled(gamesList) {
    const grid = document.getElementById('gamesGrid');
    grid.innerHTML = '';
    
    if (gamesList.length === 0) {
        grid.innerHTML = `
            <div class="col-12 text-center py-5">
                <i class="fas fa-gamepad fa-3x text-muted mb-3"></i>
                <h4 class="text-muted">No games in collection</h4>
                <p>Add games manually or scan ROM directories to get started</p>
            </div>
        `;
        return;
    }
    
    gamesList.forEach(game => {
        const gameCard = createGameCard(game);
        grid.appendChild(gameCard);
    });
}

function displayGamesList(gamesList) {
    const list = document.getElementById('gamesList');
    list.innerHTML = '';
    
    if (gamesList.length === 0) {
        list.innerHTML = `
            <div class="text-center py-5">
                <i class="fas fa-gamepad fa-3x text-muted mb-3"></i>
                <h4 class="text-muted">No games in collection</h4>
                <p>Add games manually or scan ROM directories to get started</p>
            </div>
        `;
        return;
    }
    
    const table = document.createElement('table');
    table.className = 'table table-hover';
    
    table.innerHTML = `
        <thead>
            <tr>
                <th style="width: 60px;"></th>
                <th>Title</th>
                <th>Platform</th>
                <th>Year</th>
                <th>Genre</th>
                <th>Format</th>
                <th>Rating</th>
                <th>Completion</th>
                <th>Files</th>
                <th style="width: 120px;">Actions</th>
            </tr>
        </thead>
        <tbody id="gamesTableBody">
        </tbody>
    `;
    
    const tbody = table.querySelector('#gamesTableBody');
    gamesList.forEach(game => {
        const row = createGameListRow(game);
        tbody.appendChild(row);
    });
    
    list.appendChild(table);
}

function createGameCard(game) {
    const col = document.createElement('div');
    col.className = 'col-lg-3 col-md-4 col-sm-6 mb-4';
    
    col.innerHTML = `
        <div class="card game-card" style="cursor: pointer;" onclick="showGameDetails(${game.id})">
            <div class="position-relative">
                <div class="game-cover">
                    ${game.cover_art_url ? 
                        `<img src="${game.cover_art_url}" alt="${game.title}">` : 
                        `<i class="fas fa-gamepad"></i>`
                    }
                </div>
                <span class="platform-badge">${game.platform?.name || 'Unknown'}</span>
                ${renderCompletionStatus(game)}
                <div class="btn-group-actions position-absolute bottom-0 end-0 p-2" onclick="event.stopPropagation();">
                    <button class="btn btn-sm btn-primary" onclick="editGame(${game.id})" title="Edit Game">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-sm btn-info" onclick="fetchMetadata(${game.id})" title="Update Metadata">
                        <i class="fas fa-download"></i>
                    </button>
                    <button class="btn btn-sm btn-success" onclick="showPlaySessionModal(${game.id}, '${game.title}')" title="Log Session">
                        <i class="fas fa-play"></i>
                    </button>
                    <button class="btn btn-sm btn-warning" onclick="showCompletionModal(${game.id})" title="Update Completion">
                        <i class="fas fa-trophy"></i>
                    </button>
                    <button class="btn btn-sm btn-danger" onclick="deleteGame(${game.id}, '${game.title}')" title="Delete Game">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            </div>
            <div class="card-body">
                <h6 class="card-title">${game.title}</h6>
                <p class="card-text small text-muted">
                    ${game.year ? game.year : 'Unknown Year'} • ${game.genre || 'Unknown Genre'}
                </p>
                <div class="mb-2">
                    ${renderCollectionFormats(game)}
                </div>
                ${game.rating ? `
                    <div class="rating-stars">
                        ${generateStars(game.rating)}
                        <small class="text-muted">(${game.rating.toFixed(1)}/10)</small>
                    </div>
                ` : ''}
                ${game.file_locations && game.file_locations.length > 0 ? `
                    <small class="text-success">
                        <i class="fas fa-hdd"></i> ROM: ${game.file_locations.length} file(s)
                    </small>
                ` : ''}
            </div>
        </div>
    `;
    
    return col;
}

function createGameListRow(game) {
    const row = document.createElement('tr');
    row.style.cursor = 'pointer';
    row.onclick = () => showGameDetails(game.id);
    
    const fileCount = game.file_locations ? game.file_locations.length : 0;
    const ratingStars = game.rating ? generateStars(game.rating) : '-';
    
    row.innerHTML = `
        <td>
            ${game.cover_art_url ? 
                `<img src="${game.cover_art_url}" alt="${game.title}" style="width: 40px; height: 40px; object-fit: cover;" class="rounded">` : 
                `<div class="bg-light rounded d-flex align-items-center justify-content-center" style="width: 40px; height: 40px;"><i class="fas fa-gamepad text-muted"></i></div>`
            }
        </td>
        <td>
            <strong>${game.title}</strong>
            ${fileCount > 0 ? '<br><small class="text-success"><i class="fas fa-hdd"></i> ROM</small>' : ''}
        </td>
        <td>${game.platform?.name || 'Unknown'}</td>
        <td>${game.year || '-'}</td>
        <td>${game.genre || '-'}</td>
        <td>${renderCollectionFormats(game)}</td>
        <td>
            ${game.rating ? `
                <div class="rating-stars-small">
                    ${ratingStars}
                    <small class="text-muted">(${game.rating.toFixed(1)}/10)</small>
                </div>
            ` : '-'}
        </td>
        <td onclick="event.stopPropagation(); showCompletionModal(${game.id})" style="cursor: pointer;" title="Click to update completion">
            ${renderCompletionStatusBadge(game)}
        </td>
        <td>
            ${fileCount > 0 ? `<span class="badge bg-success">${fileCount}</span>` : '<span class="badge bg-secondary">0</span>'}
        </td>
        <td onclick="event.stopPropagation();">
            <div class="btn-group btn-group-sm">
                <button class="btn btn-outline-primary" onclick="editGame(${game.id})" title="Edit">
                    <i class="fas fa-edit"></i>
                </button>
                <button class="btn btn-outline-info" onclick="fetchMetadata(${game.id})" title="Update Metadata">
                    <i class="fas fa-download"></i>
                </button>
                <button class="btn btn-outline-success" onclick="showPlaySessionModal(${game.id}, '${game.title}')" title="Log Session">
                    <i class="fas fa-play"></i>
                </button>
                <button class="btn btn-outline-warning" onclick="showCompletionModal(${game.id})" title="Update Completion">
                    <i class="fas fa-trophy"></i>
                </button>
                <button class="btn btn-outline-danger" onclick="deleteGame(${game.id}, '${game.title}')" title="Delete">
                    <i class="fas fa-trash"></i>
                </button>
            </div>
        </td>
    `;
    
    return row;
}

function generateStars(rating) {
    const stars = Math.round(rating / 2); // Convert 0-10 to 0-5 stars
    let html = '';
    for (let i = 1; i <= 5; i++) {
        html += `<i class="fas fa-star${i <= stars ? '' : '-o'}"></i>`;
    }
    return html;
}

function displaySessions(sessionsList) {
    const content = document.getElementById('sessionsContent');
    content.innerHTML = '';
    
    if (sessionsList.length === 0) {
        content.innerHTML = `
            <div class="text-center py-5">
                <i class="fas fa-clock fa-3x text-muted mb-3"></i>
                <h4 class="text-muted">No play sessions recorded</h4>
                <p>Start playing games to track your sessions</p>
            </div>
        `;
        return;
    }
    
    sessionsList.forEach(session => {
        const sessionCard = createSessionCard(session);
        content.appendChild(sessionCard);
    });
}

function createSessionCard(session) {
    const div = document.createElement('div');
    div.className = 'card session-card mb-3';
    
    const duration = session.duration || 0;
    const hours = Math.floor(duration / 60);
    const minutes = duration % 60;
    
    div.innerHTML = `
        <div class="card-body">
            <div class="row">
                <div class="col-md-8">
                    <h6 class="card-title">${session.game_title}</h6>
                    <p class="session-time">
                        <i class="fas fa-calendar"></i> ${new Date(session.start_time).toLocaleDateString()}
                        <i class="fas fa-clock ms-3"></i> ${new Date(session.start_time).toLocaleTimeString()}
                        ${session.end_time ? ` - ${new Date(session.end_time).toLocaleTimeString()}` : ' (In Progress)'}
                    </p>
                    ${session.notes ? `<p class="card-text">${session.notes}</p>` : ''}
                </div>
                <div class="col-md-4 text-end">
                    <div class="session-duration">${hours}h ${minutes}m</div>
                    ${session.rating ? `
                        <div class="rating-stars mt-1">
                            ${generateStars(session.rating)}
                        </div>
                    ` : ''}
                    <div class="btn-group btn-group-sm mt-2">
                        <button class="btn btn-outline-primary" onclick="editSession(${session.id})" title="Edit Session">
                            <i class="fas fa-edit"></i>
                        </button>
                        <button class="btn btn-outline-danger" onclick="deleteSession(${session.id})" title="Delete Session">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    `;
    
    return div;
}

function displayRecentlyPlayedGames(gamesList) {
    const grid = document.getElementById('recentlyPlayedGames');
    const section = document.getElementById('recentlyPlayedSection');
    
    if (!gamesList || gamesList.length === 0) {
        section.style.display = 'none';
        return;
    }
    
    section.style.display = 'block';
    grid.innerHTML = '';
    
    gamesList.slice(0, 6).forEach(game => {
        const gameCard = createCompactGameCard(game);
        grid.appendChild(gameCard);
    });
}

function createCompactGameCard(game) {
    const col = document.createElement('div');
    col.className = 'col-lg-2 col-md-3 col-sm-4 mb-3';
    
    // Get last play session for display
    let lastPlayed = '';
    if (game.play_sessions && game.play_sessions.length > 0) {
        // Sort sessions by start_time descending to get the most recent
        const sortedSessions = [...game.play_sessions].sort((a, b) => new Date(b.start_time) - new Date(a.start_time));
        const lastSession = sortedSessions[0];
        lastPlayed = new Date(lastSession.start_time).toLocaleDateString();
    }
    
    col.innerHTML = `
        <div class="card game-card recently-played-card">
            <div class="position-relative">
                <div class="game-cover" style="height: 120px;">
                    ${game.cover_art_url ? 
                        `<img src="${game.cover_art_url}" alt="${game.title}">` : 
                        `<i class="fas fa-gamepad"></i>`
                    }
                </div>
                <div class="btn-group-actions position-absolute bottom-0 end-0 p-1">
                    <button class="btn btn-sm btn-success" onclick="showPlaySessionModal(${game.id}, '${game.title}')" title="Log Session">
                        <i class="fas fa-play"></i>
                    </button>
                </div>
            </div>
            <div class="card-body p-2">
                <h6 class="card-title mb-1" style="font-size: 0.8rem;">${game.title}</h6>
                <small class="text-muted">
                    Last played: ${lastPlayed || 'Unknown'}
                </small>
            </div>
        </div>
    `;
    
    return col;
}

// Unified filtering and sorting functions
async function applyAllFilters() {
    // Prevent multiple simultaneous filter operations
    if (isLoading) {
        console.log('Filter operation already in progress, skipping...');
        return;
    }
    
    try {
        const filters = {
            platform: document.getElementById('platformFilter')?.value,
            genre: document.getElementById('genreFilter')?.value,
            format: document.getElementById('formatFilter')?.value,
            completion: document.querySelector('.completion-filter.active')?.dataset.status,
            sort: document.getElementById('sortBy')?.value || 'title',
            page: currentPage,
            limit: 50
        };
        
        // Reset to first page when filtering (unless pagination navigation)
        if (filters.page === 1) {
            currentPage = 1;
            filters.page = 1;
        }
        
        // Build query params
        const params = new URLSearchParams({ 
            page: filters.page, 
            limit: filters.limit,
            sort: filters.sort 
        });
        
        if (filters.platform && filters.platform !== '') {
            params.append('platform', filters.platform);
        }
        if (filters.genre && filters.genre !== '') {
            params.append('genre', filters.genre);
        }
        if (filters.format && filters.format !== '') {
            params.append('collection_format', filters.format);
        }
        if (filters.completion && filters.completion !== 'all') {
            params.append('completion_status', filters.completion);
        }
        
        // Save state to URL for bookmarking/sharing
        const currentUrl = new URL(window.location);
        currentUrl.search = params.toString();
        window.history.replaceState(filters, '', currentUrl.toString());
        
        await loadGamesWithParams(params);
    } catch (error) {
        console.error('Failed to apply filters:', error);
        showToast('Failed to apply filters: ' + error.message, 'error');
    }
}

// Debounced filter function for performance
function debouncedFilterGames() {
    if (filterDebounceTimeout) {
        clearTimeout(filterDebounceTimeout);
    }
    
    filterDebounceTimeout = setTimeout(() => {
        applyAllFilters();
    }, FILTER_DEBOUNCE_DELAY);
}

// Helper function to get platform ID by name for URL compatibility
function getPlatformIdByName(platformName) {
    const platform = platforms.find(p => p.name === platformName);
    return platform ? platform.id : null;
}

// Helper function to get platform name by ID 
function getPlatformNameById(platformId) {
    const platform = platforms.find(p => p.id.toString() === platformId.toString());
    return platform ? platform.name : null;
}

// Legacy function for backward compatibility
function filterGames() {
    applyAllFilters();
}

async function loadGamesWithParams(params) {
    if (isLoading) return;
    
    try {
        isLoading = true;
        showLoadingIndicator();
        
        // Use request manager to handle cancellation for filter operations
        const response = await requestManager.fetch(
            'filter-games',
            `${API_BASE}/games?${params.toString()}`,
            {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' }
            }
        );
        
        // Handle cancelled request
        if (response === null) {
            return; // Request was cancelled
        }
        
        const data = await response.json();
        
        // Validate response structure
        if (!data || !data.games || !data.pagination) {
            throw new Error('Invalid response structure from server');
        }
        
        games = data.games;
        currentPage = data.pagination.page;
        totalPages = data.pagination.total_pages;
        
        extractGenresFromCurrentGames();
        updateFilterOptions(true); // Preserve current filter selections
        displayGames(games);
        updatePaginationControls();
        
        // Show success feedback for filtered results
        if (params.get('platform') || params.get('genre') || params.get('completion_status')) {
            const activeFilters = [];
            if (params.get('platform')) activeFilters.push('platform');
            if (params.get('genre')) activeFilters.push('genre');
            if (params.get('completion_status')) activeFilters.push('status');
            
            console.log(`Loaded ${games.length} games with ${activeFilters.join(', ')} filters applied`);
        }
    } catch (error) {
        console.error('Failed to load games:', error);
        showToast('Failed to load games: ' + error.message, 'error');
        
        // Show fallback state
        games = [];
        displayGames(games);
    } finally {
        isLoading = false;
        hideLoadingIndicator();
    }
}

// Modal functions
function showAddGameModal() {
    updatePlatformSelects();
    resetAddGameModal();
    new bootstrap.Modal(document.getElementById('addGameModal')).show();
}

function resetAddGameModal() {
    // Reset to search step
    document.getElementById('gameSearchStep').style.display = 'block';
    document.getElementById('manualGameStep').style.display = 'none';
    document.getElementById('addGameBtn').style.display = 'none';
    
    // Clear inputs
    document.getElementById('gameSearchInput').value = '';
    document.getElementById('searchPlatform').value = '';
    document.getElementById('gameSearchResults').style.display = 'none';
    document.getElementById('searchResultsList').innerHTML = '';
    
    // Clear manual form
    if (document.getElementById('addGameForm')) {
        document.getElementById('addGameForm').reset();
    }
}

function showGameSearchStep() {
    document.getElementById('gameSearchStep').style.display = 'block';
    document.getElementById('manualGameStep').style.display = 'none';
    document.getElementById('addGameBtn').style.display = 'none';
}

function showManualGameForm() {
    document.getElementById('gameSearchStep').style.display = 'none';
    document.getElementById('manualGameStep').style.display = 'block';
    document.getElementById('addGameBtn').style.display = 'block';
    updatePlatformSelects();
}

function showAddPlatformModal() {
    new bootstrap.Modal(document.getElementById('addPlatformModal')).show();
}

// Search functionality
let searchTimeout;
async function searchGamesAsYouType() {
    const query = document.getElementById('gameSearchInput').value.trim();
    const platform = document.getElementById('searchPlatform').value;
    
    // Clear previous timeout
    if (searchTimeout) {
        clearTimeout(searchTimeout);
    }
    
    // Hide results if query is too short
    if (query.length < 3) {
        document.getElementById('gameSearchResults').style.display = 'none';
        return;
    }
    
    // Debounce search
    searchTimeout = setTimeout(async () => {
        try {
            const results = await apiCall('/games/search-metadata', 'POST', { 
                title: query, 
                platform: platform 
            });
            
            displaySearchResults(results.results);
        } catch (error) {
            console.error('Search failed:', error);
        }
    }, 500);
}

function displaySearchResults(results) {
    const resultsList = document.getElementById('searchResultsList');
    const resultsContainer = document.getElementById('gameSearchResults');
    
    resultsList.innerHTML = '';
    
    if (!results || results.length === 0) {
        resultsList.innerHTML = '<div class="text-muted">No results found. Try different search terms or add manually.</div>';
        resultsContainer.style.display = 'block';
        return;
    }
    
    results.forEach((game, index) => {
        const resultCard = document.createElement('div');
        resultCard.className = 'card mb-2 search-result-card';
        resultCard.style.cursor = 'pointer';
        
        resultCard.innerHTML = `
            <div class="card-body">
                <div class="row">
                    <div class="col-md-2">
                        ${game.cover_art_url ? 
                            `<img src="${game.cover_art_url}" class="img-fluid rounded" style="max-height: 80px;">` :
                            '<div class="bg-secondary rounded d-flex align-items-center justify-content-center" style="height: 80px; width: 60px;"><i class="fas fa-gamepad text-white"></i></div>'
                        }
                    </div>
                    <div class="col-md-10">
                        <h6 class="card-title mb-1">${game.title}</h6>
                        <div class="row">
                            <div class="col-md-6">
                                <small class="text-muted">
                                    <strong>Year:</strong> ${game.year || 'Unknown'}<br>
                                    <strong>Genre:</strong> ${game.genre || 'Unknown'}<br>
                                    <strong>Rating:</strong> ${game.rating ? game.rating.toFixed(1) + '/10' : 'N/A'}
                                </small>
                            </div>
                            <div class="col-md-6">
                                <small class="text-muted">
                                    ${game.description ? game.description.substring(0, 150) + '...' : 'No description available'}
                                </small>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
        
        resultCard.addEventListener('click', () => selectGameFromSearch(game));
        resultsList.appendChild(resultCard);
    });
    
    resultsContainer.style.display = 'block';
}

function selectGameFromSearch(gameMetadata) {
    const platformSelect = document.getElementById('searchPlatform');
    const selectedPlatformId = platformSelect.value;
    
    if (!selectedPlatformId) {
        showToast('Please select a platform first', 'warning');
        return;
    }
    
    // Create game with metadata
    createGameFromMetadata(parseInt(selectedPlatformId), gameMetadata);
}

async function createGameFromMetadata(platformId, metadata) {
    try {
        const gameData = {
            platform_id: platformId,
            metadata: metadata
        };
        
        await apiCall('/games/from-metadata', 'POST', gameData);
        
        // Close modal and reload games
        bootstrap.Modal.getInstance(document.getElementById('addGameModal')).hide();
        await applyAllFilters();
        
        showToast('Game added successfully with metadata!', 'success');
    } catch (error) {
        console.error('Failed to add game from metadata:', error);
        showToast('Failed to add game: ' + error.message, 'error');
    }
}

// CRUD operations
async function addGame() {
    try {
        // Collect collection formats
        const collectionFormats = [];
        if (document.getElementById('gameFormatPhysical').checked) {
            collectionFormats.push('physical');
        }
        if (document.getElementById('gameFormatDigital').checked) {
            collectionFormats.push('digital');
        }
        if (document.getElementById('gameFormatRom').checked) {
            collectionFormats.push('rom');
        }
        
        const gameData = {
            title: document.getElementById('gameTitle').value,
            platform_id: parseInt(document.getElementById('gamePlatform').value),
            year: parseInt(document.getElementById('gameYear').value) || 0,
            genre: document.getElementById('gameGenre').value,
            description: document.getElementById('gameDescription').value,
            collection_formats: collectionFormats
        };
        
        await apiCall('/games', 'POST', gameData);
        
        // Close modal and reload games
        bootstrap.Modal.getInstance(document.getElementById('addGameModal')).hide();
        document.getElementById('addGameForm').reset();
        await applyAllFilters();
        
        showToast('Game added successfully!', 'success');
    } catch (error) {
        console.error('Failed to add game:', error);
    }
}

async function addPlatform() {
    try {
        const platformData = {
            name: document.getElementById('platformName').value,
            manufacturer: document.getElementById('platformManufacturer').value,
            release_year: parseInt(document.getElementById('platformYear').value) || 0
        };
        
        await apiCall('/platforms', 'POST', platformData);
        
        // Close modal and reload platforms
        bootstrap.Modal.getInstance(document.getElementById('addPlatformModal')).hide();
        document.getElementById('addPlatformForm').reset();
        await loadPlatforms();
        loadPlatformsView();
        
        showToast('Platform added successfully!', 'success');
    } catch (error) {
        console.error('Failed to add platform:', error);
        showToast('Failed to add platform: ' + error.message, 'error');
    }
}

async function editPlatform(platformId) {
    try {
        const platform = await apiCall(`/platforms/${platformId}`);
        
        // Populate edit form
        document.getElementById('editPlatformId').value = platform.id;
        document.getElementById('editPlatformName').value = platform.name || '';
        document.getElementById('editPlatformManufacturer').value = platform.manufacturer || '';
        document.getElementById('editPlatformYear').value = platform.release_year || '';
        
        // Show modal
        new bootstrap.Modal(document.getElementById('editPlatformModal')).show();
    } catch (error) {
        console.error('Failed to load platform for editing:', error);
        showToast('Failed to load platform details', 'error');
    }
}

async function updatePlatform() {
    try {
        const platformId = document.getElementById('editPlatformId').value;
        const platformData = {
            name: document.getElementById('editPlatformName').value,
            manufacturer: document.getElementById('editPlatformManufacturer').value,
            release_year: parseInt(document.getElementById('editPlatformYear').value) || 0
        };
        
        await apiCall(`/platforms/${platformId}`, 'PUT', platformData);
        
        // Close modal and reload
        bootstrap.Modal.getInstance(document.getElementById('editPlatformModal')).hide();
        await loadPlatforms();
        loadPlatformsView();
        
        showToast('Platform updated successfully!', 'success');
    } catch (error) {
        console.error('Failed to update platform:', error);
        showToast('Failed to update platform: ' + error.message, 'error');
    }
}

async function fetchMetadata(gameId) {
    try {
        showToast('Fetching metadata...', 'info');
        await apiCall(`/games/${gameId}/fetch-metadata`, 'POST');
        await applyAllFilters();
        showToast('Metadata updated successfully!', 'success');
    } catch (error) {
        console.error('Failed to fetch metadata:', error);
    }
}

// Scanner functions
async function scanDirectory(event) {
    event.preventDefault();
    
    try {
        const scanData = {
            directory_path: document.getElementById('directoryPath').value,
            server_location: document.getElementById('serverLocation').value,
            platform_id: parseInt(document.getElementById('platformSelect').value),
            recursive: document.getElementById('recursive').checked
        };
        
        showToast('Starting directory scan...', 'info');
        const result = await apiCall('/scan/directory', 'POST', scanData);
        
        showToast(`Scan completed! Found ${result.files_found} files, added ${result.games_added} games`, 'success');
        
        if (result.errors && result.errors.length > 0) {
            console.warn('Scan errors:', result.errors);
        }
        
        await loadGames();
    } catch (error) {
        console.error('Failed to scan directory:', error);
    }
}

async function findDuplicates() {
    try {
        showToast('Finding duplicates...', 'info');
        const result = await apiCall('/scan/duplicates');
        
        const duplicatesDiv = document.getElementById('duplicatesResult');
        if (result.duplicates.length === 0) {
            duplicatesDiv.innerHTML = '<div class="alert alert-success">No duplicates found!</div>';
        } else {
            let html = `<div class="alert alert-warning">Found ${result.count} duplicate groups:</div>`;
            result.duplicates.forEach(group => {
                html += `
                    <div class="duplicate-group card mb-2">
                        <div class="card-body">
                            <h6>Hash: ${group.hash}</h6>
                            ${group.files.map(file => `
                                <div class="file-path">${file.server_location}: ${file.file_path}</div>
                            `).join('')}
                        </div>
                    </div>
                `;
            });
            duplicatesDiv.innerHTML = html;
        }
    } catch (error) {
        console.error('Failed to find duplicates:', error);
    }
}

async function updateMetadataBatch() {
    try {
        showToast('Starting batch metadata update...', 'info');
        
        const result = await apiCall('/scan/metadata-batch', 'POST', {
            batch_size: 5,
            delay_seconds: 2
        });
        
        const metadataDiv = document.getElementById('metadataUpdateResult');
        metadataDiv.innerHTML = `
            <div class="alert alert-info">
                <h6><i class="fas fa-cloud-download-alt"></i> Batch Metadata Update Started</h6>
                <p><strong>Total games:</strong> ${result.total_games}</p>
                <p><strong>Batch size:</strong> ${result.batch_size} games at a time</p>
                <p><strong>Delay:</strong> ${result.delay} seconds between batches</p>
                <p class="mb-0"><small>This process runs in the background and respects IGDB rate limits. Check the games collection to see updated metadata as it processes.</small></p>
            </div>
        `;
        
        showToast(`Batch update started for ${result.total_games} games`, 'success');
    } catch (error) {
        console.error('Failed to start batch metadata update:', error);
        const metadataDiv = document.getElementById('metadataUpdateResult');
        metadataDiv.innerHTML = `
            <div class="alert alert-danger">
                <h6><i class="fas fa-exclamation-triangle"></i> Update Failed</h6>
                <p>Failed to start batch metadata update: ${error.message}</p>
            </div>
        `;
    }
}

// Search function
async function searchGames(event) {
    event.preventDefault();
    
    const query = document.getElementById('searchInput').value.trim();
    
    if (!query) {
        // Clear search and reapply current filters
        await applyAllFilters();
        return;
    }
    
    try {
        // Build search params including current filters
        const searchParams = {
            title: query,
            platform: document.getElementById('platformFilter')?.value || '',
            genre: document.getElementById('genreFilter')?.value || '',
            format: document.getElementById('formatFilter')?.value || '',
            completion_status: document.querySelector('.completion-filter.active')?.dataset.status || 'all'
        };
        
        // Remove empty values
        Object.keys(searchParams).forEach(key => {
            if (!searchParams[key] || searchParams[key] === 'all' || searchParams[key] === '') {
                delete searchParams[key];
            }
        });
        
        const results = await apiCall('/games/search', 'POST', searchParams);
        displayGames(results);
        
        // Update URL to include search
        const currentUrl = new URL(window.location);
        currentUrl.searchParams.set('search', query);
        window.history.replaceState(null, '', currentUrl.toString());
        
    } catch (error) {
        console.error('Search failed:', error);
        showToast('Search failed', 'error');
    }
}

// Utility functions
function updatePlatformSelects() {
    const selects = ['gamePlatform', 'platformSelect', 'searchPlatform', 'editGamePlatform'];
    
    selects.forEach(selectId => {
        const select = document.getElementById(selectId);
        if (select) {
            // Clear existing options except the first one
            select.innerHTML = select.children[0].outerHTML;
            
            platforms.forEach(platform => {
                const option = document.createElement('option');
                option.value = platform.id;
                option.textContent = platform.name;
                select.appendChild(option);
            });
        }
    });
}

// Edit game functionality
async function editGame(gameId) {
    await loadGameForEdit(gameId);
    new bootstrap.Modal(document.getElementById('editGameModal')).show();
}

async function loadGameForEdit(gameId) {
    try {
        const game = await apiCall(`/games/${gameId}`);
        
        // Populate edit form
        document.getElementById('editGameId').value = game.id;
        document.getElementById('editGameTitle').value = game.title || '';
        document.getElementById('editGamePlatform').value = game.platform_id || '';
        document.getElementById('editGameYear').value = game.year || '';
        document.getElementById('editGameGenre').value = game.genre || '';
        document.getElementById('editGameRating').value = game.rating || '';
        document.getElementById('editGameDescription').value = game.description || '';
        document.getElementById('editGameCoverUrl').value = game.cover_art_url || '';
        
        // Handle purchase date
        if (game.purchase_date) {
            const date = new Date(game.purchase_date);
            document.getElementById('editGamePurchaseDate').value = date.toISOString().split('T')[0];
        } else {
            document.getElementById('editGamePurchaseDate').value = '';
        }
        
        // Handle collection formats
        document.getElementById('editGameFormatPhysical').checked = false;
        document.getElementById('editGameFormatDigital').checked = false;
        document.getElementById('editGameFormatRom').checked = false;
        
        if (game.collection_formats && Array.isArray(game.collection_formats)) {
            game.collection_formats.forEach(format => {
                if (format === 'physical') {
                    document.getElementById('editGameFormatPhysical').checked = true;
                } else if (format === 'digital') {
                    document.getElementById('editGameFormatDigital').checked = true;
                } else if (format === 'rom') {
                    document.getElementById('editGameFormatRom').checked = true;
                }
            });
        }
        
        // Update platform options
        updatePlatformSelects();
    } catch (error) {
        console.error('Failed to load game for editing:', error);
        showToast('Failed to load game details', 'error');
        throw error;
    }
}

async function updateGame() {
    try {
        const gameId = document.getElementById('editGameId').value;
        const gameData = {
            title: document.getElementById('editGameTitle').value,
            platform_id: parseInt(document.getElementById('editGamePlatform').value),
            year: parseInt(document.getElementById('editGameYear').value) || 0,
            genre: document.getElementById('editGameGenre').value,
            rating: parseFloat(document.getElementById('editGameRating').value) || 0,
            description: document.getElementById('editGameDescription').value,
            cover_art_url: document.getElementById('editGameCoverUrl').value,
        };
        
        // Handle purchase date
        const purchaseDate = document.getElementById('editGamePurchaseDate').value;
        if (purchaseDate) {
            gameData.purchase_date = new Date(purchaseDate).toISOString();
        }
        
        // Collect collection formats
        const collectionFormats = [];
        if (document.getElementById('editGameFormatPhysical').checked) {
            collectionFormats.push('physical');
        }
        if (document.getElementById('editGameFormatDigital').checked) {
            collectionFormats.push('digital');
        }
        if (document.getElementById('editGameFormatRom').checked) {
            collectionFormats.push('rom');
        }
        gameData.collection_formats = collectionFormats;
        
        await apiCall(`/games/${gameId}`, 'PUT', gameData);
        
        // Check if we came from game details modal by checking if currentDetailGame is set
        const returnToDetails = currentDetailGame && currentDetailGame.id == gameId;
        
        // Close edit modal
        const editModal = bootstrap.Modal.getInstance(document.getElementById('editGameModal'));
        editModal.hide();
        
        // Reload data preserving filters and pagination
        await applyAllFilters();
        await loadRecentlyPlayedGames();
        
        showToast('Game updated successfully!', 'success');
        
        // If we came from game details modal, wait for edit modal to close then reopen details
        if (returnToDetails) {
            const editModalElement = document.getElementById('editGameModal');
            editModalElement.addEventListener('hidden.bs.modal', function onEditHidden() {
                editModalElement.removeEventListener('hidden.bs.modal', onEditHidden);
                // Refresh the game details modal with updated data
                showGameDetails(gameId);
            });
        }
    } catch (error) {
        console.error('Failed to update game:', error);
        showToast('Failed to update game', 'error');
    }
}

// Play session functionality
function showPlaySessionModal(gameId, gameTitle) {
    setupPlaySessionModal(gameId, gameTitle);
    new bootstrap.Modal(document.getElementById('addSessionModal')).show();
}

function setupPlaySessionModal(gameId, gameTitle) {
    document.getElementById('sessionGameId').value = gameId;
    document.getElementById('sessionGameTitle').textContent = gameTitle;
    
    // Set default start time to now
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    document.getElementById('sessionStartTime').value = now.toISOString().slice(0, 16);
    
    // Clear other fields
    document.getElementById('sessionEndTime').value = '';
    document.getElementById('sessionRating').value = '';
    document.getElementById('sessionNotes').value = '';
}

async function addPlaySession() {
    try {
        const gameId = document.getElementById('sessionGameId').value;
        const sessionData = {
            start_time: new Date(document.getElementById('sessionStartTime').value).toISOString(),
            notes: document.getElementById('sessionNotes').value,
        };
        
        // Add end time if provided
        const endTime = document.getElementById('sessionEndTime').value;
        if (endTime) {
            sessionData.end_time = new Date(endTime).toISOString();
        }
        
        // Add rating if provided
        const rating = document.getElementById('sessionRating').value;
        if (rating) {
            sessionData.rating = parseInt(rating);
        }
        
        await apiCall(`/games/${gameId}/sessions`, 'POST', sessionData);
        
        // Close modal and reload
        bootstrap.Modal.getInstance(document.getElementById('addSessionModal')).hide();
        await loadRecentlyPlayedGames();
        await loadActiveSessions();
        
        showToast('Play session logged successfully!', 'success');
    } catch (error) {
        console.error('Failed to add play session:', error);
        showToast('Failed to log play session', 'error');
    }
}

function loadPlatformsForScanner() {
    updatePlatformSelects();
}

function loadPlatformsView() {
    const list = document.getElementById('platformsList');
    list.innerHTML = '';
    
    platforms.forEach(platform => {
        const card = document.createElement('div');
        card.className = 'card mb-3';
        card.innerHTML = `
            <div class="card-body">
                <div class="row">
                    <div class="col">
                        <h5>${platform.name}</h5>
                        <p class="text-muted">${platform.manufacturer || 'Unknown Manufacturer'} • ${platform.release_year || 'Unknown Year'}</p>
                    </div>
                    <div class="col-auto">
                        <div class="btn-group btn-group-sm">
                            <button class="btn btn-outline-primary" onclick="editPlatform(${platform.id})" title="Edit Platform">
                                <i class="fas fa-edit"></i>
                            </button>
                            <button class="btn btn-outline-danger" onclick="deletePlatform(${platform.id}, '${platform.name}')" title="Delete Platform">
                                <i class="fas fa-trash"></i>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `;
        list.appendChild(card);
    });
}

function showToast(message, type = 'info') {
    // Create toast container if it doesn't exist
    let container = document.querySelector('.toast-container');
    if (!container) {
        container = document.createElement('div');
        container.className = 'toast-container';
        document.body.appendChild(container);
    }
    
    const toast = document.createElement('div');
    toast.className = `toast align-items-center text-white bg-${type === 'error' ? 'danger' : type} border-0`;
    toast.setAttribute('role', 'alert');
    
    toast.innerHTML = `
        <div class="d-flex">
            <div class="toast-body">${message}</div>
            <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
        </div>
    `;
    
    container.appendChild(toast);
    
    const bsToast = new bootstrap.Toast(toast, { delay: 3000 });
    bsToast.show();
    
    // Remove toast element after it's hidden
    toast.addEventListener('hidden.bs.toast', () => {
        toast.remove();
    });
}

// Directory browser functionality
let currentBrowserPath = '/';

function showDirectoryBrowser() {
    currentBrowserPath = '/data/roms';
    new bootstrap.Modal(document.getElementById('directoryBrowserModal')).show();
    loadDirectoryContents('/data/roms');
}

async function loadDirectoryContents(path) {
    try {
        document.getElementById('directoryEntries').innerHTML = `
            <div class="text-center py-4">
                <i class="fas fa-spinner fa-spin fa-2x"></i>
                <p class="mt-2">Loading directories...</p>
            </div>
        `;
        
        const response = await apiCall(`/browse?path=${encodeURIComponent(path)}`);
        currentBrowserPath = response.current_path;
        document.getElementById('currentPath').value = currentBrowserPath;
        
        displayDirectoryEntries(response.entries);
    } catch (error) {
        document.getElementById('directoryEntries').innerHTML = `
            <div class="alert alert-danger">
                <i class="fas fa-exclamation-triangle"></i> 
                Failed to load directory: ${error.message}
            </div>
        `;
    }
}

function displayDirectoryEntries(entries) {
    const container = document.getElementById('directoryEntries');
    container.innerHTML = '';
    
    if (entries.length === 0) {
        container.innerHTML = `
            <div class="text-center py-4 text-muted">
                <i class="fas fa-folder-open fa-2x mb-2"></i>
                <p>This directory is empty</p>
            </div>
        `;
        return;
    }
    
    entries.forEach(entry => {
        if (entry.is_dir) {
            const entryDiv = document.createElement('div');
            entryDiv.className = 'directory-entry d-flex align-items-center p-2 border-bottom';
            entryDiv.style.cursor = 'pointer';
            
            entryDiv.innerHTML = `
                <div class="me-3">
                    <i class="fas ${entry.name === '..' ? 'fa-arrow-up' : 'fa-folder'} text-primary"></i>
                </div>
                <div class="flex-grow-1">
                    <div class="fw-medium">${entry.name}</div>
                    <small class="text-muted">${entry.path}</small>
                </div>
                <div class="text-end">
                    ${entry.modified ? `<small class="text-muted">${entry.modified}</small>` : ''}
                </div>
            `;
            
            entryDiv.addEventListener('click', () => navigateToPath(entry.path));
            entryDiv.addEventListener('mouseenter', () => entryDiv.classList.add('bg-light'));
            entryDiv.addEventListener('mouseleave', () => entryDiv.classList.remove('bg-light'));
            
            container.appendChild(entryDiv);
        }
    });
}

function navigateToPath(path) {
    loadDirectoryContents(path);
}

function selectCurrentDirectory() {
    document.getElementById('directoryPath').value = currentBrowserPath;
    bootstrap.Modal.getInstance(document.getElementById('directoryBrowserModal')).hide();
    showToast(`Selected directory: ${currentBrowserPath}`, 'success');
}

// Game details modal functionality
async function showGameDetails(gameId) {
    try {
        const game = await apiCall(`/games/${gameId}`);
        currentDetailGame = game;
        
        // Populate modal
        document.getElementById('detailGameTitle').textContent = game.title;
        document.getElementById('detailPlatform').textContent = game.platform?.name || 'Unknown';
        document.getElementById('detailYear').textContent = game.year || '-';
        document.getElementById('detailGenre').textContent = game.genre || '-';
        document.getElementById('detailDescription').textContent = game.description || 'No description available';
        
        // Handle cover art
        const coverImage = document.getElementById('detailCoverImage');
        const noCover = document.getElementById('detailNoCover');
        if (game.cover_art_url) {
            coverImage.src = game.cover_art_url;
            coverImage.style.display = 'block';
            noCover.style.display = 'none';
        } else {
            coverImage.style.display = 'none';
            noCover.style.display = 'flex';
        }
        
        // Handle rating
        const ratingDiv = document.getElementById('detailRating');
        if (game.rating) {
            ratingDiv.innerHTML = `
                <div class="rating-stars">
                    ${generateStars(game.rating)}
                    <span class="text-muted">(${game.rating.toFixed(1)}/10)</span>
                </div>
            `;
        } else {
            ratingDiv.textContent = '-';
        }
        
        // Handle purchase date
        const purchaseDate = document.getElementById('detailPurchaseDate');
        if (game.purchase_date) {
            const date = new Date(game.purchase_date);
            purchaseDate.textContent = date.toLocaleDateString();
        } else {
            purchaseDate.textContent = '-';
        }
        
        // Handle files
        const filesDiv = document.getElementById('detailFiles');
        if (game.file_locations && game.file_locations.length > 0) {
            filesDiv.innerHTML = `
                <span class="badge bg-success">${game.file_locations.length} file(s)</span>
                <div class="mt-1">
                    ${game.file_locations.map(file => `
                        <small class="text-muted d-block">${file.server_location}: ${file.file_path}</small>
                    `).join('')}
                </div>
            `;
        } else {
            filesDiv.innerHTML = '<span class="badge bg-secondary">No files</span>';
        }
        
        // Handle collection formats
        const formatsDiv = document.getElementById('detailFormats');
        if (game.collection_formats && Array.isArray(game.collection_formats) && game.collection_formats.length > 0) {
            formatsDiv.innerHTML = renderCollectionFormatsDetailed(game);
        } else {
            formatsDiv.innerHTML = '<span class="text-muted">No format specified</span>';
        }
        
        // Load and display play sessions
        await loadGameSessions(gameId);
        
        // Show modal
        new bootstrap.Modal(document.getElementById('gameDetailsModal')).show();
    } catch (error) {
        console.error('Failed to load game details:', error);
        showToast('Failed to load game details', 'error');
    }
}

async function loadGameSessions(gameId) {
    try {
        const sessions = await apiCall(`/games/${gameId}/sessions`);
        const sessionsList = document.getElementById('detailSessionsList');
        
        if (sessions.length === 0) {
            sessionsList.innerHTML = '<p class="text-muted">No play sessions recorded</p>';
            return;
        }
        
        // Show last 5 sessions
        const recentSessions = sessions.slice(0, 5);
        sessionsList.innerHTML = recentSessions.map(session => {
            const duration = session.duration || 0;
            const hours = Math.floor(duration / 60);
            const minutes = duration % 60;
            const startDate = new Date(session.start_time);
            
            return `
                <div class="border-bottom py-2">
                    <div class="d-flex justify-content-between align-items-start">
                        <div>
                            <small class="text-muted">
                                <i class="fas fa-calendar me-1"></i>${startDate.toLocaleDateString()}
                                <i class="fas fa-clock ms-2 me-1"></i>${startDate.toLocaleTimeString()}
                            </small>
                            ${session.notes ? `<div class="small">${session.notes}</div>` : ''}
                        </div>
                        <div class="text-end">
                            <small class="text-muted">${hours}h ${minutes}m</small>
                            ${session.rating ? `
                                <div class="rating-stars-small">
                                    ${generateStars(session.rating)}
                                </div>
                            ` : ''}
                        </div>
                    </div>
                </div>
            `;
        }).join('');
        
        if (sessions.length > 5) {
            sessionsList.innerHTML += `<small class="text-muted">...and ${sessions.length - 5} more sessions</small>`;
        }
    } catch (error) {
        console.error('Failed to load sessions:', error);
        document.getElementById('detailSessionsList').innerHTML = '<p class="text-muted">Failed to load sessions</p>';
    }
}

async function editGameFromDetails() {
    if (currentDetailGame) {
        // Use modal manager to safely transition between modals
        await modalManager.openModal('editGameModal', async () => {
            await loadGameForEdit(currentDetailGame.id);
        });
    }
}

async function fetchMetadataFromDetails() {
    if (currentDetailGame) {
        try {
            showToast('Updating metadata...', 'info');
            await fetchMetadata(currentDetailGame.id);
            // Refresh the details modal
            showGameDetails(currentDetailGame.id);
        } catch (error) {
            console.error('Failed to update metadata:', error);
        }
    }
}

async function logSessionFromDetails() {
    if (currentDetailGame) {
        // Use modal manager to safely transition between modals
        await modalManager.openModal('addSessionModal', async () => {
            setupPlaySessionModal(currentDetailGame.id, currentDetailGame.title);
        });
    }
}

function deleteGameFromDetails() {
    if (currentDetailGame) {
        bootstrap.Modal.getInstance(document.getElementById('gameDetailsModal')).hide();
        deleteGame(currentDetailGame.id, currentDetailGame.title);
    }
}

// Pagination controls
function updatePaginationControls() {
    const paginationInfo = document.getElementById('paginationInfo');
    const paginationControls = document.getElementById('paginationControls');
    
    if (!paginationInfo || !paginationControls) return;
    
    // Update info - use actual pagination data if available
    const start = ((currentPage - 1) * 50) + 1;
    const end = Math.min(currentPage * 50, games.length);
    const total = games.length;
    
    paginationInfo.innerHTML = `Showing ${start}-${end} of ${total} games`;
    
    // Update controls
    paginationControls.innerHTML = '';
    
    if (totalPages <= 1) {
        paginationControls.style.display = 'none';
        return;
    }
    
    paginationControls.style.display = 'flex';
    
    // Previous button
    const prevLi = document.createElement('li');
    prevLi.className = `page-item ${currentPage === 1 ? 'disabled' : ''}`;
    prevLi.innerHTML = `
        <a class="page-link" href="#" onclick="navigateToPage(${currentPage - 1}); return false;">
            <i class="fas fa-chevron-left"></i>
        </a>
    `;
    paginationControls.appendChild(prevLi);
    
    // Page numbers
    const startPage = Math.max(1, currentPage - 2);
    const endPage = Math.min(totalPages, currentPage + 2);
    
    if (startPage > 1) {
        const firstLi = document.createElement('li');
        firstLi.className = 'page-item';
        firstLi.innerHTML = `<a class="page-link" href="#" onclick="navigateToPage(1); return false;">1</a>`;
        paginationControls.appendChild(firstLi);
        
        if (startPage > 2) {
            const ellipsisLi = document.createElement('li');
            ellipsisLi.className = 'page-item disabled';
            ellipsisLi.innerHTML = '<span class="page-link">...</span>';
            paginationControls.appendChild(ellipsisLi);
        }
    }
    
    for (let i = startPage; i <= endPage; i++) {
        const pageLi = document.createElement('li');
        pageLi.className = `page-item ${i === currentPage ? 'active' : ''}`;
        pageLi.innerHTML = `<a class="page-link" href="#" onclick="navigateToPage(${i}); return false;">${i}</a>`;
        paginationControls.appendChild(pageLi);
    }
    
    if (endPage < totalPages) {
        if (endPage < totalPages - 1) {
            const ellipsisLi = document.createElement('li');
            ellipsisLi.className = 'page-item disabled';
            ellipsisLi.innerHTML = '<span class="page-link">...</span>';
            paginationControls.appendChild(ellipsisLi);
        }
        
        const lastLi = document.createElement('li');
        lastLi.className = 'page-item';
        lastLi.innerHTML = `<a class="page-link" href="#" onclick="navigateToPage(${totalPages}); return false;">${totalPages}</a>`;
        paginationControls.appendChild(lastLi);
    }
    
    // Next button
    const nextLi = document.createElement('li');
    nextLi.className = `page-item ${currentPage === totalPages ? 'disabled' : ''}`;
    nextLi.innerHTML = `
        <a class="page-link" href="#" onclick="navigateToPage(${currentPage + 1}); return false;">
            <i class="fas fa-chevron-right"></i>
        </a>
    `;
    paginationControls.appendChild(nextLi);
}

// Navigation function for pagination with filters
async function navigateToPage(page) {
    // Validate page number
    if (!page || page < 1 || page > totalPages) {
        console.warn(`Invalid page number: ${page}. Valid range: 1-${totalPages}`);
        return;
    }
    
    // Prevent navigation during loading
    if (isLoading) {
        console.log('Navigation blocked: loading in progress');
        return;
    }
    
    currentPage = page;
    await applyAllFilters();
}

// Delete operations
async function deleteGame(gameId, gameTitle) {
    if (!confirm(`Are you sure you want to delete "${gameTitle}"?\n\nThis action cannot be undone and will also delete all associated play sessions and file locations.`)) {
        return;
    }
    
    try {
        await apiCall(`/games/${gameId}`, 'DELETE');
        await applyAllFilters();
        await loadRecentlyPlayedGames();
        showToast('Game deleted successfully', 'success');
    } catch (error) {
        console.error('Failed to delete game:', error);
        showToast('Failed to delete game: ' + error.message, 'error');
    }
}

async function deletePlatform(platformId, platformName) {
    if (!confirm(`Are you sure you want to delete platform "${platformName}"?\n\nThis action cannot be undone. Note: Platforms with associated games cannot be deleted.`)) {
        return;
    }
    
    try {
        await apiCall(`/platforms/${platformId}`, 'DELETE');
        await loadPlatforms();
        loadPlatformsView();
        showToast('Platform deleted successfully', 'success');
    } catch (error) {
        console.error('Failed to delete platform:', error);
        showToast('Failed to delete platform: ' + error.message, 'error');
    }
}

// Session management functions
async function editSession(sessionId) {
    try {
        // Get session data - we need to find it in the sessions array or make an API call
        let session = sessions.find(s => s.id === sessionId);
        
        if (!session) {
            // If not found in current sessions, we might need to fetch it
            showToast('Session not found', 'error');
            return;
        }
        
        // Find the game title for this session
        const game = games.find(g => g.id === session.game_id);
        
        // Populate edit form
        document.getElementById('editSessionId').value = session.id;
        document.getElementById('editSessionGameTitle').textContent = game ? game.title : 'Unknown Game';
        
        // Format dates for datetime-local input
        if (session.start_time) {
            const startDate = new Date(session.start_time);
            startDate.setMinutes(startDate.getMinutes() - startDate.getTimezoneOffset());
            document.getElementById('editSessionStartTime').value = startDate.toISOString().slice(0, 16);
        }
        
        if (session.end_time) {
            const endDate = new Date(session.end_time);
            endDate.setMinutes(endDate.getMinutes() - endDate.getTimezoneOffset());
            document.getElementById('editSessionEndTime').value = endDate.toISOString().slice(0, 16);
        } else {
            document.getElementById('editSessionEndTime').value = '';
        }
        
        document.getElementById('editSessionRating').value = session.rating || '';
        document.getElementById('editSessionNotes').value = session.notes || '';
        
        // Show modal
        new bootstrap.Modal(document.getElementById('editSessionModal')).show();
    } catch (error) {
        console.error('Failed to load session for editing:', error);
        showToast('Failed to load session details', 'error');
    }
}

async function updatePlaySession() {
    try {
        const sessionId = document.getElementById('editSessionId').value;
        const sessionData = {
            start_time: new Date(document.getElementById('editSessionStartTime').value).toISOString(),
            notes: document.getElementById('editSessionNotes').value,
        };
        
        // Add end time if provided
        const endTime = document.getElementById('editSessionEndTime').value;
        if (endTime) {
            sessionData.end_time = new Date(endTime).toISOString();
        }
        
        // Add rating if provided
        const rating = document.getElementById('editSessionRating').value;
        if (rating) {
            sessionData.rating = parseInt(rating);
        }
        
        await apiCall(`/sessions/${sessionId}`, 'PUT', sessionData);
        
        // Close modal and reload
        bootstrap.Modal.getInstance(document.getElementById('editSessionModal')).hide();
        await loadSessions();
        await loadRecentlyPlayedGames();
        
        showToast('Session updated successfully!', 'success');
    } catch (error) {
        console.error('Failed to update session:', error);
        showToast('Failed to update session: ' + error.message, 'error');
    }
}

async function deleteSession(sessionId) {
    if (!confirm('Are you sure you want to delete this play session?\n\nThis action cannot be undone.')) {
        return;
    }
    
    try {
        await apiCall(`/sessions/${sessionId}`, 'DELETE');
        await loadSessions();
        await loadRecentlyPlayedGames();
        showToast('Session deleted successfully', 'success');
    } catch (error) {
        console.error('Failed to delete session:', error);
        showToast('Failed to delete session: ' + error.message, 'error');
    }
}

// Settings and backup functions
async function loadSettingsView() {
    await loadBackupInfo();
    await loadVersionInfo();
}

async function loadBackupInfo() {
    try {
        const info = await apiCall('/backup/info');
        const stats = info.database_stats;
        
        document.getElementById('backupStats').innerHTML = `
            <div class="col-3">
                <div class="fs-4 text-primary">${stats.total_games}</div>
                <small class="text-muted">Games</small>
            </div>
            <div class="col-3">
                <div class="fs-4 text-success">${stats.total_platforms}</div>
                <small class="text-muted">Platforms</small>
            </div>
            <div class="col-3">
                <div class="fs-4 text-info">${stats.total_sessions}</div>
                <small class="text-muted">Sessions</small>
            </div>
            <div class="col-3">
                <div class="fs-4 text-warning">${stats.total_files}</div>
                <small class="text-muted">Files</small>
            </div>
        `;
        
        // Update Nextcloud status
        updateNextcloudStatus(info.nextcloud);
        
    } catch (error) {
        console.error('Failed to load backup info:', error);
    }
}

function updateNextcloudStatus(nextcloudInfo) {
    const statusDiv = document.getElementById('nextcloudStatus');
    const controlsDiv = document.getElementById('nextcloudControls');
    
    if (nextcloudInfo.configured) {
        statusDiv.innerHTML = '<span class="badge bg-success"><i class="fas fa-check me-1"></i>Configured</span>';
        controlsDiv.style.display = 'block';
    } else {
        statusDiv.innerHTML = '<span class="badge bg-warning"><i class="fas fa-exclamation-triangle me-1"></i>Not Configured</span>';
        controlsDiv.style.display = 'none';
        
        document.getElementById('nextcloudMessage').innerHTML = `
            <div class="alert alert-info">
                <small><strong>To enable Nextcloud backup:</strong><br>
                Add these environment variables to your Docker deployment:<br>
                <code>NEXTCLOUD_URL</code>, <code>NEXTCLOUD_USERNAME</code>, <code>NEXTCLOUD_PASSWORD</code>
                </small>
            </div>
        `;
    }
}

async function exportFullBackup() {
    try {
        showToast('Creating backup...', 'info');
        
        const response = await fetch(`${API_BASE}/backup/export`);
        if (!response.ok) {
            throw new Error('Backup export failed');
        }
        
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `pelico_backup_${new Date().toISOString().split('T')[0]}.json`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
        
        showToast('Backup downloaded successfully!', 'success');
    } catch (error) {
        console.error('Failed to export backup:', error);
        showToast('Failed to export backup: ' + error.message, 'error');
    }
}

async function exportGamesOnly() {
    try {
        showToast('Creating games export...', 'info');
        
        const response = await fetch(`${API_BASE}/backup/export/games`);
        if (!response.ok) {
            throw new Error('Games export failed');
        }
        
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `pelico_games_export_${new Date().toISOString().split('T')[0]}.json`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
        
        showToast('Games export downloaded successfully!', 'success');
    } catch (error) {
        console.error('Failed to export games:', error);
        showToast('Failed to export games: ' + error.message, 'error');
    }
}

async function importBackup() {
    const fileInput = document.getElementById('backupFile');
    const file = fileInput.files[0];
    
    if (!file) {
        showToast('Please select a backup file first', 'warning');
        return;
    }
    
    if (!file.name.endsWith('.json')) {
        showToast('Please select a valid JSON backup file', 'error');
        return;
    }
    
    if (!confirm('⚠️ WARNING: This will REPLACE ALL EXISTING DATA!\n\nAre you sure you want to restore from this backup?\n\nThis action cannot be undone.')) {
        return;
    }
    
    try {
        showToast('Reading backup file...', 'info');
        
        const fileContent = await file.text();
        const backupData = JSON.parse(fileContent);
        
        // Validate backup structure
        if (!backupData.version || !backupData.games) {
            throw new Error('Invalid backup file format');
        }
        
        showToast('Restoring backup...', 'info');
        
        const result = await apiCall('/backup/import', 'POST', backupData);
        
        showToast(`Backup restored successfully!\n\nRestored: ${result.stats.games} games, ${result.stats.platforms} platforms, ${result.stats.sessions} sessions`, 'success');
        
        // Clear file input
        fileInput.value = '';
        
        // Reload all data
        await refreshAllData();
        
    } catch (error) {
        console.error('Failed to import backup:', error);
        showToast('Failed to restore backup: ' + error.message, 'error');
    }
}

async function refreshAllData() {
    try {
        showToast('Refreshing all data...', 'info');
        
        await loadPlatforms();
        await loadGames();
        await loadRecentlyPlayedGames();
        await loadSessions();
        
        if (currentView === 'platforms') {
            loadPlatformsView();
        }
        if (currentView === 'settings') {
            await loadBackupInfo();
        }
        
        showToast('All data refreshed successfully!', 'success');
    } catch (error) {
        console.error('Failed to refresh data:', error);
        showToast('Failed to refresh some data', 'error');
    }
}

async function updateAllMetadata() {
    if (!confirm('This will update metadata for all games using IGDB. This may take a while and will respect rate limits.\n\nContinue?')) {
        return;
    }
    
    try {
        showToast('Starting metadata update for all games...', 'info');
        
        await updateMetadataBatch();
        
        showToast('Metadata update started in background', 'success');
    } catch (error) {
        console.error('Failed to start metadata update:', error);
        showToast('Failed to start metadata update: ' + error.message, 'error');
    }
}

// Nextcloud backup functions
async function backupToNextcloud() {
    try {
        showToast('Uploading backup to Nextcloud...', 'info');
        await apiCall('/backup/nextcloud', 'POST');
        showToast('Backup uploaded to Nextcloud successfully!', 'success');
    } catch (error) {
        showToast('Failed to upload to Nextcloud: ' + error.message, 'error');
    }
}

async function testNextcloudConnection() {
    try {
        const result = await apiCall('/backup/nextcloud/test');
        showToast('Nextcloud connection successful!', 'success');
    } catch (error) {
        showToast('Nextcloud connection failed: ' + error.message, 'error');
    }
}

// Completion tracking functions (renderCompletionStatus defined earlier in file)

async function showCompletionModal(gameId) {
    try {
        const game = await apiCall(`/games/${gameId}`);
        
        // Populate completion modal
        document.getElementById('completionGameId').value = game.id;
        document.getElementById('completionGameTitle').textContent = game.title;
        document.getElementById('completionStatus').value = game.completion_status || 'not_started';
        document.getElementById('completionPercentage').value = game.completion_percentage || 0;
        document.getElementById('percentageLabel').textContent = (game.completion_percentage || 0) + '%';
        document.getElementById('completionNotes').value = game.completion_notes || '';
        
        // Show modal
        new bootstrap.Modal(document.getElementById('completionModal')).show();
    } catch (error) {
        console.error('Failed to load game for completion update:', error);
        showToast('Failed to load game details', 'error');
    }
}

async function updateCompletionStatus() {
    try {
        const gameId = document.getElementById('completionGameId').value;
        const completionData = {
            status: document.getElementById('completionStatus').value,
            percentage: parseInt(document.getElementById('completionPercentage').value),
            notes: document.getElementById('completionNotes').value
        };
        
        await apiCall(`/games/${gameId}/completion`, 'PUT', completionData);
        
        // Close modal and refresh with current filters preserved
        bootstrap.Modal.getInstance(document.getElementById('completionModal')).hide();
        await applyAllFilters();
        
        showToast('Completion status updated successfully!', 'success');
    } catch (error) {
        console.error('Failed to update completion status:', error);
        showToast('Failed to update completion status: ' + error.message, 'error');
    }
}

function updatePercentageLabel() {
    const percentage = document.getElementById('completionPercentage').value;
    document.getElementById('percentageLabel').textContent = percentage + '%';
}

async function loadCompletionStats() {
    try {
        const stats = await apiCall('/games/stats/completion');
        
        document.getElementById('completionStats').innerHTML = `
            <div class="row text-center">
                <div class="col-md-2">
                    <div class="fs-4 text-success">${stats.completed + stats.hundred_percent}</div>
                    <small class="text-muted">Completed</small>
                </div>
                <div class="col-md-2">
                    <div class="fs-4 text-primary">${stats.in_progress}</div>
                    <small class="text-muted">Playing</small>
                </div>
                <div class="col-md-2">
                    <div class="fs-4 text-secondary">${stats.not_started}</div>
                    <small class="text-muted">Backlog</small>
                </div>
                <div class="col-md-2">
                    <div class="fs-4 text-warning">${stats.abandoned}</div>
                    <small class="text-muted">Abandoned</small>
                </div>
                <div class="col-md-2">
                    <div class="fs-4 text-info">${stats.completion_rate.toFixed(1)}%</div>
                    <small class="text-muted">Completion Rate</small>
                </div>
                <div class="col-md-2">
                    <div class="fs-4 text-dark">${stats.average_progress.toFixed(0)}%</div>
                    <small class="text-muted">Avg Progress</small>
                </div>
            </div>
        `;
        
    } catch (error) {
        console.error('Failed to load completion stats:', error);
    }
}

async function filterByCompletionStatus(status) {
    try {
        // Reset pagination when changing completion filter
        currentPage = 1;
        
        // Update active filter button first
        document.querySelectorAll('.completion-filter').forEach(btn => btn.classList.remove('active'));
        const targetButton = document.querySelector(`[data-status="${status}"]`);
        if (targetButton) {
            targetButton.classList.add('active');
        }
        
        // Use unified filter system which integrates with other filters
        await applyAllFilters();
        
        // Show user feedback
        const statusLabel = status === 'all' ? 'all' :
                           status === 'backlog' ? 'backlog' : 
                           status === 'in_progress' ? 'currently playing' :
                           status === '100_percent' ? '100% completed' :
                           status.replace('_', ' ');
        
        if (status === 'all') {
            showToast('Showing all games', 'info');
        } else {
            showToast(`Filtered by ${statusLabel} status`, 'info');
        }
    } catch (error) {
        console.error('Failed to filter games by completion status:', error);
        showToast('Failed to filter games', 'error');
    }
}


// Helper function to escape HTML
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Version information functions
async function loadVersionInfo() {
    try {
        const version = await apiCall('/version');
        const versionDiv = document.getElementById('versionInfo');
        const changelogDiv = document.getElementById('changelogInfo');
        
        // Display version information
        versionDiv.innerHTML = `
            <div class="row">
                <div class="col-md-6">
                    <dl class="row">
                        <dt class="col-sm-6">Version</dt>
                        <dd class="col-sm-6">
                            <span class="badge bg-primary">v${version.version}</span>
                        </dd>
                        <dt class="col-sm-6">Build Time</dt>
                        <dd class="col-sm-6">
                            <small class="text-muted">${version.build_time || 'Development'}</small>
                        </dd>
                    </dl>
                </div>
                <div class="col-md-6">
                    <dl class="row">
                        <dt class="col-sm-6">Git Commit</dt>
                        <dd class="col-sm-6">
                            <small class="text-muted font-monospace">${version.git_commit || 'Unknown'}</small>
                        </dd>
                        <dt class="col-sm-6">Updated</dt>
                        <dd class="col-sm-6">
                            <button class="btn btn-sm btn-outline-secondary" onclick="loadVersionInfo()">
                                <i class="fas fa-sync"></i> Refresh
                            </button>
                        </dd>
                    </dl>
                </div>
            </div>
        `;
        
        // Display recent changes
        const recentChanges = [
            '✅ Fixed completion status update preserving active filters',
            '✅ Enhanced search to respect platform/genre/format filters',
            '✅ Implemented proper modal management (no more stacking bugs)',
            '✅ Added request cancellation for rapid filter changes',
            '✅ Fixed rating display precision (now shows 8.6/10 instead of 8.567/10)',
            '✅ Fixed last played display showing proper recent sessions'
        ];
        
        changelogDiv.innerHTML = `
            <div class="card">
                <div class="card-body">
                    <h6 class="card-title">Version ${version.version} Changes</h6>
                    <ul class="list-unstyled mb-0">
                        ${recentChanges.map(change => `<li class="mb-1"><small>${change}</small></li>`).join('')}
                    </ul>
                    <div class="mt-2">
                        <small class="text-muted">
                            <i class="fas fa-clock"></i> Last updated: ${new Date(version.timestamp).toLocaleString()}
                        </small>
                    </div>
                </div>
            </div>
        `;
        
    } catch (error) {
        console.error('Failed to load version info:', error);
        document.getElementById('versionInfo').innerHTML = `
            <div class="alert alert-warning">
                <small>Failed to load version information</small>
            </div>
        `;
        document.getElementById('changelogInfo').innerHTML = `
            <div class="alert alert-warning">
                <small>Failed to load changelog</small>
            </div>
        `;
    }
}