<script lang="ts">
  import { onMount } from 'svelte';
  import type { Game } from '$lib/models';
  import { api, ApiError } from '$lib/api';
  import GameCard from '$lib/components/GameCard.svelte';
  import GameFormModal from '$lib/components/GameFormModal.svelte';

  interface Stats {
    total_games: number;
    total_platforms: number;
    total_sessions: number;
    total_wishlist: number;
    total_shortlist: number;
    completion_stats?: {
      completed: number;
      in_progress: number;
      not_started: number;
      abandoned: number;
    };
    playtime_stats?: {
      total_hours: number;
      average_hours: number;
    };
  }

  let stats: Stats | null = null;
  let recentlyPlayed: Game[] = [];
  let recentlyAdded: Game[] = [];
  let platforms: Platform[] = [];
  let loading = true;
  let error: string | null = null;
  let showModal = false;
  let selectedGame: Game | null = null;

  onMount(async () => {
    await loadData();
  });

  async function loadData() {
    try {
      loading = true;
      error = null;
      
      const [statsData, recentlyPlayedData, recentlyAddedData, platformsData] = await Promise.all([
        api.getStats().catch(() => null),
        api.getRecentlyPlayedGames().catch(() => []),
        api.getRecentlyAddedGames(5).catch(() => []),
        api.getPlatforms().catch(() => [])
      ]);

      stats = statsData;
      recentlyPlayed = recentlyPlayedData || [];
      recentlyAdded = recentlyAddedData || [];
      platforms = platformsData || [];
    } catch (e) {
      if (e instanceof ApiError) {
        error = `Failed to load data: ${e.message}`;
      } else {
        error = 'Error connecting to the backend. Please check if the server is running.';
      }
      console.error('Load data error:', e);
    } finally {
      loading = false;
    }
  }

  function openAddGameModal() {
    selectedGame = null;
    showModal = true;
  }

  function openEditGameModal(event: CustomEvent<Game>) {
    selectedGame = event.detail;
    showModal = true;
  }

  async function handleDeleteGame(event: CustomEvent<number>) {
    const gameId = event.detail;
    
    try {
      await api.deleteGame(gameId);
      // Reload data to refresh the lists
      await loadData();
    } catch (e) {
      if (e instanceof ApiError) {
        alert(`Failed to delete game: ${e.message}`);
      } else {
        alert('An error occurred while deleting the game.');
      }
      console.error('Delete game error:', e);
    }
  }

  async function handleFormSubmit(event: CustomEvent<Partial<Game>>) {
    const gameData = event.detail;
    const isEdit = !!selectedGame;

    try {
      if (isEdit && selectedGame) {
        await api.updateGame(selectedGame.id, gameData);
      } else {
        await api.createGame(gameData);
      }
      
      showModal = false;
      // Reload data to refresh the lists
      await loadData();
    } catch (e) {
      if (e instanceof ApiError) {
        alert(`Failed to save game: ${e.message}`);
      } else {
        alert('An error occurred while saving the game.');
      }
      console.error('Save game error:', e);
    }
  }

  async function handleLogSession(event: CustomEvent<number>) {
    const gameId = event.detail;
    
    try {
      // Create a simple session with current timestamp
      await api.createSession(gameId, {
        play_time_minutes: 60, // Default 1 hour
        notes: '',
        date_played: new Date().toISOString()
      });
      
      // Reload data to refresh recently played
      await loadData();
    } catch (e) {
      if (e instanceof ApiError) {
        alert(`Failed to log session: ${e.message}`);
      } else {
        alert('An error occurred while logging the session.');
      }
      console.error('Log session error:', e);
    }
  }

  function handleViewGame(event: CustomEvent<number>) {
    const gameId = event.detail;
    // TODO: Navigate to game detail view
    console.log('View game:', gameId);
  }

  function getCompletionPercentage(completed: number, total: number): number {
    return total > 0 ? Math.round((completed / total) * 100) : 0;
  }
</script>

<div class="container-fluid mt-4">
  <div class="d-flex justify-content-between align-items-center mb-4">
    <div>
      <h1>
        <i class="fas fa-home me-2 text-primary"></i>
        Dashboard
      </h1>
      <p class="text-muted mb-0">Welcome to your game collection</p>
    </div>
    <div class="d-flex gap-2">
      <button class="btn btn-outline-secondary" on:click={loadData} disabled={loading}>
        <i class="fas fa-sync-alt {loading ? 'fa-spin' : ''}"></i> Refresh
      </button>
      <button class="btn btn-primary" on:click={openAddGameModal}>
        <i class="fas fa-plus me-1"></i> Add Game
      </button>
    </div>
  </div>

  {#if loading}
    <div class="d-flex justify-content-center align-items-center" style="height: 300px;">
      <div class="text-center">
        <div class="spinner-border text-primary mb-3" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
        <p class="text-muted">Loading dashboard...</p>
      </div>
    </div>
  {:else if error}
    <div class="alert alert-danger d-flex align-items-center">
      <i class="fas fa-exclamation-triangle me-2"></i>
      <div>
        {error}
        <button class="btn btn-sm btn-outline-danger ms-2" on:click={loadData}>
          Try Again
        </button>
      </div>
    </div>
  {:else}
    <!-- Statistics Cards -->
    {#if stats}
      <div class="row mb-4">
        <div class="col-lg-2 col-md-4 col-sm-6 mb-3">
          <div class="card text-center h-100 bg-primary text-white">
            <div class="card-body">
              <i class="fas fa-gamepad fa-2x mb-2"></i>
              <h6 class="card-title">Total Games</h6>
              <p class="card-text fs-3 fw-bold mb-0">{stats.total_games}</p>
            </div>
          </div>
        </div>
        <div class="col-lg-2 col-md-4 col-sm-6 mb-3">
          <div class="card text-center h-100 bg-secondary text-white">
            <div class="card-body">
              <i class="fas fa-desktop fa-2x mb-2"></i>
              <h6 class="card-title">Platforms</h6>
              <p class="card-text fs-3 fw-bold mb-0">{stats.total_platforms}</p>
            </div>
          </div>
        </div>
        <div class="col-lg-2 col-md-4 col-sm-6 mb-3">
          <div class="card text-center h-100 bg-success text-white">
            <div class="card-body">
              <i class="fas fa-play fa-2x mb-2"></i>
              <h6 class="card-title">Play Sessions</h6>
              <p class="card-text fs-3 fw-bold mb-0">{stats.total_sessions}</p>
            </div>
          </div>
        </div>
        <div class="col-lg-2 col-md-4 col-sm-6 mb-3">
          <div class="card text-center h-100 bg-info text-white">
            <div class="card-body">
              <i class="fas fa-clock fa-2x mb-2"></i>
              <h6 class="card-title">Total Hours</h6>
              <p class="card-text fs-3 fw-bold mb-0">
                {stats.playtime_stats?.total_hours || 0}
              </p>
            </div>
          </div>
        </div>
        <div class="col-lg-2 col-md-4 col-sm-6 mb-3">
          <div class="card text-center h-100 bg-warning text-dark">
            <div class="card-body">
              <i class="fas fa-list fa-2x mb-2"></i>
              <h6 class="card-title">Shortlist</h6>
              <p class="card-text fs-3 fw-bold mb-0">{stats.total_shortlist}</p>
            </div>
          </div>
        </div>
        <div class="col-lg-2 col-md-4 col-sm-6 mb-3">
          <div class="card text-center h-100 bg-danger text-white">
            <div class="card-body">
              <i class="fas fa-heart fa-2x mb-2"></i>
              <h6 class="card-title">Wishlist</h6>
              <p class="card-text fs-3 fw-bold mb-0">{stats.total_wishlist}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Completion Progress -->
      {#if stats.completion_stats}
        <div class="row mb-4">
          <div class="col-12">
            <div class="card">
              <div class="card-header">
                <h5 class="mb-0">
                  <i class="fas fa-trophy me-2"></i>
                  Completion Progress
                </h5>
              </div>
              <div class="card-body">
                <div class="row text-center">
                  <div class="col-3">
                    <div class="completion-stat text-success">
                      <i class="fas fa-check-circle fa-2x mb-1"></i>
                      <h6>Completed</h6>
                      <span class="fs-4 fw-bold">{stats.completion_stats.completed}</span>
                    </div>
                  </div>
                  <div class="col-3">
                    <div class="completion-stat text-warning">
                      <i class="fas fa-play-circle fa-2x mb-1"></i>
                      <h6>In Progress</h6>
                      <span class="fs-4 fw-bold">{stats.completion_stats.in_progress}</span>
                    </div>
                  </div>
                  <div class="col-3">
                    <div class="completion-stat text-muted">
                      <i class="fas fa-pause-circle fa-2x mb-1"></i>
                      <h6>Not Started</h6>
                      <span class="fs-4 fw-bold">{stats.completion_stats.not_started}</span>
                    </div>
                  </div>
                  <div class="col-3">
                    <div class="completion-stat text-danger">
                      <i class="fas fa-times-circle fa-2x mb-1"></i>
                      <h6>Abandoned</h6>
                      <span class="fs-4 fw-bold">{stats.completion_stats.abandoned}</span>
                    </div>
                  </div>
                </div>
                
                <div class="text-center mt-3">
                  <h6>Overall Completion Rate</h6>
                  <div class="progress mb-2" style="height: 20px;">
                    <div 
                      class="progress-bar bg-success" 
                      role="progressbar" 
                      style="width: {getCompletionPercentage(stats.completion_stats.completed, stats.total_games)}%"
                    >
                      {getCompletionPercentage(stats.completion_stats.completed, stats.total_games)}%
                    </div>
                  </div>
                  <p class="text-muted mb-0">
                    {stats.completion_stats.completed} of {stats.total_games} games completed
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      {/if}
    {/if}

    <!-- Recently Played Games -->
    <div class="row mb-4">
      <div class="col-12">
        <div class="card">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">
              <i class="fas fa-history me-2"></i>
              Recently Played
            </h5>
            {#if recentlyPlayed.length > 0}
              <small class="text-muted">{recentlyPlayed.length} games</small>
            {/if}
          </div>
          <div class="card-body">
            {#if recentlyPlayed.length === 0}
              <div class="text-center py-4">
                <i class="fas fa-play fa-3x text-muted mb-3"></i>
                <h6>No recent sessions</h6>
                <p class="text-muted">Start playing games to see them here!</p>
              </div>
            {:else}
              <div class="row">
                {#each recentlyPlayed.slice(0, 5) as game (game.id)}
                  <div class="col-xl-2 col-lg-3 col-md-4 col-sm-6 mb-3">
                    <GameCard 
                      {game} 
                      on:edit={openEditGameModal}
                      on:delete={handleDeleteGame}
                      on:log={handleLogSession}
                      on:view={handleViewGame}
                    />
                  </div>
                {/each}
              </div>
              {#if recentlyPlayed.length > 5}
                <div class="text-center">
                  <a href="/collection" class="btn btn-outline-primary">
                    View All Recent Games
                  </a>
                </div>
              {/if}
            {/if}
          </div>
        </div>
      </div>
    </div>

    <!-- Recently Added Games -->
    <div class="row mb-4">
      <div class="col-12">
        <div class="card">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">
              <i class="fas fa-plus-circle me-2"></i>
              Recently Added
            </h5>
            {#if recentlyAdded.length > 0}
              <small class="text-muted">{recentlyAdded.length} games</small>
            {/if}
          </div>
          <div class="card-body">
            {#if recentlyAdded.length === 0}
              <div class="text-center py-4">
                <i class="fas fa-gamepad fa-3x text-muted mb-3"></i>
                <h6>No games yet</h6>
                <p class="text-muted">Add your first game to get started!</p>
                <button class="btn btn-primary" on:click={openAddGameModal}>
                  <i class="fas fa-plus me-1"></i> Add Game
                </button>
              </div>
            {:else}
              <div class="row">
                {#each recentlyAdded as game (game.id)}
                  <div class="col-xl-2 col-lg-3 col-md-4 col-sm-6 mb-3">
                    <GameCard 
                      {game} 
                      on:edit={openEditGameModal}
                      on:delete={handleDeleteGame}
                      on:log={handleLogSession}
                      on:view={handleViewGame}
                    />
                  </div>
                {/each}
              </div>
              <div class="text-center">
                <a href="/collection" class="btn btn-outline-primary">
                  <i class="fas fa-eye me-1"></i> View Full Collection
                </a>
              </div>
            {/if}
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>

{#if showModal}
  <GameFormModal 
    game={selectedGame} 
    {platforms}
    on:submit={handleFormSubmit} 
    on:close={() => showModal = false} 
  />
{/if}

<style>
  .completion-stat {
    padding: 1rem;
    border-radius: 0.5rem;
  }
</style>