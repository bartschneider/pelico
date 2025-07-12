<script lang="ts">
  import { onMount } from 'svelte';

  let stats: any = null;
  let loading = true;
  let error: string | null = null;

  onMount(async () => {
    try {
      const response = await fetch('/api/v1/stats');
      if (response.ok) {
        stats = await response.json();
      } else {
        error = 'Failed to load statistics from the backend.';
      }
    } catch (e) {
      error = 'Error connecting to the backend.';
      console.error(e);
    } finally {
      loading = false;
    }
  });
</script>

<div class="container-fluid mt-4">
  <h1>Statistics</h1>

  {#if loading}
    <p>Loading statistics...</p>
  {:else if error}
    <div class="alert alert-danger">{error}</div>
  {:else if stats}
    <div class="row">
      <div class="col-md-4">
        <div class="card text-center">
          <div class="card-body">
            <h5 class="card-title">Total Games</h5>
            <p class="card-text fs-1">{stats.total_games}</p>
          </div>
        </div>
      </div>
      <div class="col-md-4">
        <div class="card text-center">
          <div class="card-body">
            <h5 class="card-title">Total Platforms</h5>
            <p class="card-text fs-1">{stats.total_platforms}</p>
          </div>
        </div>
      </div>
      <div class="col-md-4">
        <div class="card text-center">
          <div class="card-body">
            <h5 class="card-title">Total Play Sessions</h5>
            <p class="card-text fs-1">{stats.total_sessions}</p>
          </div>
        </div>
      </div>
    </div>
    <div class="row mt-4">
      <div class="col-md-6">
        <div class="card text-center">
          <div class="card-body">
            <h5 class="card-title">Wishlist Items</h5>
            <p class="card-text fs-1">{stats.total_wishlist}</p>
          </div>
        </div>
      </div>
      <div class="col-md-6">
        <div class="card text-center">
          <div class="card-body">
            <h5 class="card-title">Shortlist Items</h5>
            <p class="card-text fs-1">{stats.total_shortlist}</p>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>
