<script lang="ts">
  import { onMount } from 'svelte';
  import type { ShortlistItem } from '$lib/models';
  import { api, ApiError } from '$lib/api';
  import GameCard from '$lib/components/GameCard.svelte';

  let shortlist: ShortlistItem[] = [];
  let loading = true;
  let error: string | null = null;
  let showAddModal = false;
  let searchQuery = '';
  let filteredShortlist: ShortlistItem[] = [];

  $: filteredShortlist = searchQuery 
    ? shortlist.filter(item => 
        item.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        item.platform_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        item.reason?.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : shortlist;

  onMount(async () => {
    await loadShortlist();
  });

  async function loadShortlist() {
    try {
      loading = true;
      error = null;
      shortlist = await api.getShortlist();
    } catch (e) {
      if (e instanceof ApiError) {
        error = `Failed to load shortlist: ${e.message}`;
      } else {
        error = 'Error connecting to the backend. Please check if the server is running.';
      }
      console.error('Load shortlist error:', e);
    } finally {
      loading = false;
    }
  }

  async function removeFromShortlist(itemId: number) {
    try {
      await api.removeFromShortlist(itemId);
      shortlist = shortlist.filter(item => item.id !== itemId);
    } catch (e) {
      if (e instanceof ApiError) {
        alert(`Failed to remove from shortlist: ${e.message}`);
      } else {
        alert('An error occurred while removing the game from the shortlist.');
      }
      console.error('Remove from shortlist error:', e);
    }
  }

  function getPriorityColor(priority: string) {
    switch (priority) {
      case 'high': return 'bg-danger';
      case 'medium': return 'bg-warning';
      case 'low': return 'bg-success';
      default: return 'bg-secondary';
    }
  }

  function getPriorityIcon(priority: string) {
    switch (priority) {
      case 'high': return 'fas fa-arrow-up';
      case 'medium': return 'fas fa-minus';
      case 'low': return 'fas fa-arrow-down';
      default: return 'fas fa-circle';
    }
  }

  function handleGameEdit(event: CustomEvent) {
    // If the shortlist item has an associated game, show the game
    const item = event.detail;
    if (item.game) {
      // TODO: Navigate to game detail or edit view
      console.log('Edit game from shortlist:', item.game);
    }
  }

  function handleGameView(event: CustomEvent) {
    const gameId = event.detail;
    // TODO: Navigate to game detail view
    console.log('View game from shortlist:', gameId);
  }
</script>

<div class="container-fluid mt-4">
  <div class="d-flex justify-content-between align-items-center mb-4">
    <h1>
      <i class="fas fa-list me-2 text-warning"></i>
      Shortlist
    </h1>
    <div class="d-flex gap-2">
      <button class="btn btn-outline-secondary" on:click={loadShortlist} disabled={loading}>
        <i class="fas fa-sync-alt {loading ? 'fa-spin' : ''}"></i> Refresh
      </button>
      <button class="btn btn-primary" on:click={() => showAddModal = true}>
        <i class="fas fa-plus me-1"></i> Add to Shortlist
      </button>
    </div>
  </div>

  <!-- Search Bar -->
  <div class="row mb-4">
    <div class="col-md-6">
      <div class="input-group">
        <span class="input-group-text">
          <i class="fas fa-search"></i>
        </span>
        <input 
          type="text" 
          class="form-control" 
          placeholder="Search shortlist by title, platform, or reason..."
          bind:value={searchQuery}
        />
        {#if searchQuery}
          <button 
            class="btn btn-outline-secondary" 
            type="button"
            on:click={() => searchQuery = ''}
          >
            <i class="fas fa-times"></i>
          </button>
        {/if}
      </div>
    </div>
    <div class="col-md-6">
      <p class="text-muted mt-2 mb-0">
        {#if searchQuery}
          Showing {filteredShortlist.length} of {shortlist.length} items
        {:else}
          {shortlist.length} games to play next
        {/if}
      </p>
    </div>
  </div>

  {#if loading}
    <div class="d-flex justify-content-center align-items-center" style="height: 300px;">
      <div class="text-center">
        <div class="spinner-border text-primary mb-3" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
        <p class="text-muted">Loading shortlist...</p>
      </div>
    </div>
  {:else if error}
    <div class="alert alert-danger d-flex align-items-center">
      <i class="fas fa-exclamation-triangle me-2"></i>
      <div>
        {error}
        <button class="btn btn-sm btn-outline-danger ms-2" on:click={loadShortlist}>
          Try Again
        </button>
      </div>
    </div>
  {:else if filteredShortlist.length === 0}
    <div class="text-center py-5">
      {#if searchQuery}
        <i class="fas fa-search fa-3x text-muted mb-3"></i>
        <h4>No items found</h4>
        <p class="text-muted">No shortlist items match your search "{searchQuery}"</p>
        <button class="btn btn-outline-primary" on:click={() => searchQuery = ''}>
          Clear Search
        </button>
      {:else}
        <i class="fas fa-list fa-3x text-muted mb-3"></i>
        <h4>Your shortlist is empty</h4>
        <p class="text-muted">Add games from your collection that you want to play next!</p>
        <button class="btn btn-primary" on:click={() => showAddModal = true}>
          <i class="fas fa-plus me-1"></i> Add Your First Game
        </button>
      {/if}
    </div>
  {:else}
    <div class="row">
      {#each filteredShortlist as item (item.id)}
        <div class="col-xl-3 col-lg-4 col-md-6 mb-4">
          {#if item.game}
            <!-- If item has associated game, show GameCard -->
            <GameCard 
              game={item.game} 
              on:edit={handleGameEdit}
              on:view={handleGameView}
              on:delete={() => removeFromShortlist(item.id)}
            />
          {:else}
            <!-- If item is standalone, show custom card -->
            <div class="card h-100">
              <div class="card-body">
                <div class="d-flex justify-content-between align-items-start mb-2">
                  <h6 class="card-title mb-0" title={item.title}>{item.title}</h6>
                  <span class="badge {getPriorityColor(item.priority)} ms-2">
                    <i class="{getPriorityIcon(item.priority)} me-1"></i>
                    {item.priority.charAt(0).toUpperCase() + item.priority.slice(1)}
                  </span>
                </div>
                
                {#if item.platform_name}
                  <p class="text-muted small mb-2">
                    <i class="fas fa-gamepad me-1"></i>
                    {item.platform_name}
                  </p>
                {/if}
                
                {#if item.reason}
                  <p class="card-text small text-muted mb-3">
                    <strong>Why play:</strong> {item.reason}
                  </p>
                {/if}
                
                <div class="d-flex justify-content-between align-items-center">
                  <small class="text-muted">
                    Added {new Date(item.created_at || '').toLocaleDateString()}
                  </small>
                  <button 
                    class="btn btn-sm btn-outline-danger"
                    on:click={() => removeFromShortlist(item.id)}
                    title="Remove from shortlist"
                  >
                    <i class="fas fa-trash"></i>
                  </button>
                </div>
              </div>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>

<!-- TODO: Add modal for adding new shortlist items -->
{#if showAddModal}
  <div class="modal fade show d-block" style="background-color: rgba(0,0,0,0.5);">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Add to Shortlist</h5>
          <button type="button" class="btn-close" on:click={() => showAddModal = false}></button>
        </div>
        <div class="modal-body">
          <p class="text-muted">Feature coming soon! Add games from your collection to your shortlist.</p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" on:click={() => showAddModal = false}>Close</button>
        </div>
      </div>
    </div>
  </div>
{/if}
