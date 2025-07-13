<script lang="ts">
  import { onMount } from 'svelte';
  import { api, ApiError } from '$lib/api';

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
    platform_breakdown?: Array<{
      platform: string;
      count: number;
    }>;
  }

  let stats: Stats | null = null;
  let loading = true;
  let error: string | null = null;

  onMount(async () => {
    await loadStats();
  });

  async function loadStats() {
    try {
      loading = true;
      error = null;
      stats = await api.getStats();
    } catch (e) {
      if (e instanceof ApiError) {
        error = `Failed to load statistics: ${e.message}`;
      } else {
        error = 'Error connecting to the backend. Please check if the server is running.';
      }
      console.error('Load stats error:', e);
    } finally {
      loading = false;
    }
  }

  function getCompletionPercentage(completed: number, total: number): number {
    return total > 0 ? Math.round((completed / total) * 100) : 0;
  }

  function getProgressBarColor(percentage: number): string {
    if (percentage >= 75) return 'bg-success';
    if (percentage >= 50) return 'bg-warning';
    if (percentage >= 25) return 'bg-info';
    return 'bg-danger';
  }
</script>

<div class="container-fluid mt-4">
  <div class="d-flex justify-content-between align-items-center mb-4">
    <h1>
      <i class="fas fa-chart-bar me-2 text-info"></i>
      Statistics
    </h1>
    <button class="btn btn-outline-secondary" on:click={loadStats} disabled={loading}>
      <i class="fas fa-sync-alt {loading ? 'fa-spin' : ''}"></i> Refresh
    </button>
  </div>

  {#if loading}
    <div class="d-flex justify-content-center align-items-center" style="height: 400px;">
      <div class="text-center">
        <div class="spinner-border text-primary mb-3" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
        <p class="text-muted">Loading statistics...</p>
      </div>
    </div>
  {:else if error}
    <div class="alert alert-danger d-flex align-items-center">
      <i class="fas fa-exclamation-triangle me-2"></i>
      <div>
        {error}
        <button class="btn btn-sm btn-outline-danger ms-2" on:click={loadStats}>
          Try Again
        </button>
      </div>
    </div>
  {:else if stats}
    <!-- Main Statistics Cards -->
    <div class="row mb-4">
      <div class="col-lg-2 col-md-4 col-sm-6 mb-3">
        <div class="card text-center h-100 bg-primary text-white">
          <div class="card-body">
            <i class="fas fa-gamepad fa-2x mb-2"></i>
            <h5 class="card-title">Total Games</h5>
            <p class="card-text fs-2 fw-bold">{stats.total_games}</p>
          </div>
        </div>
      </div>
      <div class="col-lg-2 col-md-4 col-sm-6 mb-3">
        <div class="card text-center h-100 bg-secondary text-white">
          <div class="card-body">
            <i class="fas fa-desktop fa-2x mb-2"></i>
            <h5 class="card-title">Platforms</h5>
            <p class="card-text fs-2 fw-bold">{stats.total_platforms}</p>
          </div>
        </div>
      </div>
      <div class="col-lg-2 col-md-4 col-sm-6 mb-3">
        <div class="card text-center h-100 bg-success text-white">
          <div class="card-body">
            <i class="fas fa-play fa-2x mb-2"></i>
            <h5 class="card-title">Play Sessions</h5>
            <p class="card-text fs-2 fw-bold">{stats.total_sessions}</p>
          </div>
        </div>
      </div>
      <div class="col-lg-2 col-md-4 col-sm-6 mb-3">
        <div class="card text-center h-100 bg-danger text-white">
          <div class="card-body">
            <i class="fas fa-heart fa-2x mb-2"></i>
            <h5 class="card-title">Wishlist</h5>
            <p class="card-text fs-2 fw-bold">{stats.total_wishlist}</p>
          </div>
        </div>
      </div>
      <div class="col-lg-2 col-md-4 col-sm-6 mb-3">
        <div class="card text-center h-100 bg-warning text-dark">
          <div class="card-body">
            <i class="fas fa-list fa-2x mb-2"></i>
            <h5 class="card-title">Shortlist</h5>
            <p class="card-text fs-2 fw-bold">{stats.total_shortlist}</p>
          </div>
        </div>
      </div>
      <div class="col-lg-2 col-md-4 col-sm-6 mb-3">
        <div class="card text-center h-100 bg-info text-white">
          <div class="card-body">
            <i class="fas fa-clock fa-2x mb-2"></i>
            <h5 class="card-title">Total Hours</h5>
            <p class="card-text fs-2 fw-bold">
              {stats.playtime_stats?.total_hours || 0}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Completion Statistics -->
    {#if stats.completion_stats}
      <div class="row mb-4">
        <div class="col-12">
          <div class="card">
            <div class="card-header">
              <h5 class="mb-0">
                <i class="fas fa-trophy me-2"></i>
                Game Completion Progress
              </h5>
            </div>
            <div class="card-body">
              <div class="row">
                <div class="col-md-6">
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
                </div>
                <div class="col-md-6">
                  {#if stats.completion_stats}
                    {@const completionRate = getCompletionPercentage(stats.completion_stats.completed, stats.total_games)}
                    <div class="text-center">
                      <h6>Completion Rate</h6>
                      <div class="progress mb-2" style="height: 20px;">
                        <div 
                          class="progress-bar {getProgressBarColor(completionRate)}" 
                          role="progressbar" 
                          style="width: {completionRate}%"
                          aria-valuenow={completionRate}
                          aria-valuemin="0" 
                          aria-valuemax="100"
                        >
                          {completionRate}%
                        </div>
                      </div>
                      <p class="text-muted mb-0">
                        {stats.completion_stats.completed} of {stats.total_games} games completed
                      </p>
                    </div>
                  {/if}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    {/if}

    <!-- Platform Breakdown -->
    {#if stats.platform_breakdown && stats.platform_breakdown.length > 0}
      <div class="row mb-4">
        <div class="col-12">
          <div class="card">
            <div class="card-header">
              <h5 class="mb-0">
                <i class="fas fa-pie-chart me-2"></i>
                Games by Platform
              </h5>
            </div>
            <div class="card-body">
              <div class="row">
                {#each stats.platform_breakdown as platform (platform.platform)}
                  <div class="col-lg-3 col-md-4 col-sm-6 mb-3">
                    <div class="platform-stat text-center">
                      <h6 class="fw-bold">{platform.platform}</h6>
                      <span class="badge bg-primary fs-6">{platform.count} games</span>
                      <div class="progress mt-2" style="height: 8px;">
                        <div 
                          class="progress-bar" 
                          style="width: {(platform.count / stats.total_games) * 100}%"
                        ></div>
                      </div>
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          </div>
        </div>
      </div>
    {/if}

    <!-- Playtime Statistics -->
    {#if stats.playtime_stats}
      <div class="row">
        <div class="col-md-6 mb-4">
          <div class="card text-center">
            <div class="card-body">
              <i class="fas fa-clock fa-3x text-info mb-3"></i>
              <h5 class="card-title">Total Playtime</h5>
              <p class="card-text">
                <span class="fs-2 fw-bold text-info">{stats.playtime_stats.total_hours}</span>
                <br>
                <small class="text-muted">hours logged</small>
              </p>
            </div>
          </div>
        </div>
        <div class="col-md-6 mb-4">
          <div class="card text-center">
            <div class="card-body">
              <i class="fas fa-chart-line fa-3x text-success mb-3"></i>
              <h5 class="card-title">Average Playtime</h5>
              <p class="card-text">
                <span class="fs-2 fw-bold text-success">{stats.playtime_stats.average_hours.toFixed(1)}</span>
                <br>
                <small class="text-muted">hours per game</small>
              </p>
            </div>
          </div>
        </div>
      </div>
    {/if}
  {:else}
    <div class="text-center py-5">
      <i class="fas fa-chart-bar fa-3x text-muted mb-3"></i>
      <h4>No statistics available</h4>
      <p class="text-muted">Add some games to your collection to see statistics!</p>
    </div>
  {/if}
</div>

<style>
  .completion-stat {
    padding: 1rem;
    border-radius: 0.5rem;
    margin-bottom: 1rem;
  }

  .platform-stat {
    padding: 1rem;
    border: 1px solid #dee2e6;
    border-radius: 0.5rem;
    background-color: #f8f9fa;
  }

  .platform-stat:hover {
    background-color: #e9ecef;
    transition: background-color 0.2s ease;
  }
</style>
