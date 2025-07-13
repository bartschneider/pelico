<script lang="ts">
  import { onMount } from 'svelte';
  import type { WishlistItem } from '$lib/models';
  import { api, ApiError } from '$lib/api';

  let wishlist: WishlistItem[] = [];
  let loading = true;
  let error: string | null = null;
  let showAddModal = false;
  let searchQuery = '';
  let filteredWishlist: WishlistItem[] = [];

  $: filteredWishlist = searchQuery 
    ? wishlist.filter(item => 
        item.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        item.platform_name?.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : wishlist;

  onMount(async () => {
    await loadWishlist();
  });

  async function loadWishlist() {
    try {
      loading = true;
      error = null;
      wishlist = await api.getWishlist();
    } catch (e) {
      if (e instanceof ApiError) {
        error = `Failed to load wishlist: ${e.message}`;
      } else {
        error = 'Error connecting to the backend. Please check if the server is running.';
      }
      console.error('Load wishlist error:', e);
    } finally {
      loading = false;
    }
  }

  async function removeFromWishlist(itemId: number) {
    try {
      await api.removeFromWishlist(itemId);
      wishlist = wishlist.filter(item => item.id !== itemId);
    } catch (e) {
      if (e instanceof ApiError) {
        alert(`Failed to remove from wishlist: ${e.message}`);
      } else {
        alert('An error occurred while removing the game from the wishlist.');
      }
      console.error('Remove from wishlist error:', e);
    }
  }

  function getPriorityColor(priority: string) {
    switch (priority) {
      case 'high': return 'text-danger';
      case 'medium': return 'text-warning';
      case 'low': return 'text-success';
      default: return 'text-muted';
    }
  }

  function getPriorityIcon(priority: string) {
    switch (priority) {
      case 'high': return 'fas fa-exclamation-circle';
      case 'medium': return 'fas fa-exclamation-triangle';
      case 'low': return 'fas fa-info-circle';
      default: return 'fas fa-circle';
    }
  }
</script>

<div class="container-fluid mt-4">
  <div class="d-flex justify-content-between align-items-center mb-4">
    <h1>
      <i class="fas fa-heart me-2 text-danger"></i>
      Wishlist
    </h1>
    <div class="d-flex gap-2">
      <button class="btn btn-outline-secondary" on:click={loadWishlist} disabled={loading}>
        <i class="fas fa-sync-alt {loading ? 'fa-spin' : ''}"></i> Refresh
      </button>
      <button class="btn btn-primary" on:click={() => showAddModal = true}>
        <i class="fas fa-plus me-1"></i> Add to Wishlist
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
          placeholder="Search wishlist by title or platform..."
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
          Showing {filteredWishlist.length} of {wishlist.length} items
        {:else}
          {wishlist.length} items on your wishlist
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
        <p class="text-muted">Loading wishlist...</p>
      </div>
    </div>
  {:else if error}
    <div class="alert alert-danger d-flex align-items-center">
      <i class="fas fa-exclamation-triangle me-2"></i>
      <div>
        {error}
        <button class="btn btn-sm btn-outline-danger ms-2" on:click={loadWishlist}>
          Try Again
        </button>
      </div>
    </div>
  {:else if filteredWishlist.length === 0}
    <div class="text-center py-5">
      {#if searchQuery}
        <i class="fas fa-search fa-3x text-muted mb-3"></i>
        <h4>No items found</h4>
        <p class="text-muted">No wishlist items match your search "{searchQuery}"</p>
        <button class="btn btn-outline-primary" on:click={() => searchQuery = ''}>
          Clear Search
        </button>
      {:else}
        <i class="fas fa-heart fa-3x text-muted mb-3"></i>
        <h4>Your wishlist is empty</h4>
        <p class="text-muted">Start adding games you want to play to your wishlist!</p>
        <button class="btn btn-primary" on:click={() => showAddModal = true}>
          <i class="fas fa-plus me-1"></i> Add Your First Game
        </button>
      {/if}
    </div>
  {:else}
    <div class="row">
      {#each filteredWishlist as item (item.id)}
        <div class="col-xl-3 col-lg-4 col-md-6 mb-4">
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
              
              {#if item.release_date}
                <p class="text-muted small mb-2">
                  <i class="fas fa-calendar me-1"></i>
                  {new Date(item.release_date).toLocaleDateString()}
                </p>
              {/if}
              
              {#if item.estimated_price}
                <p class="text-success small mb-2">
                  <i class="fas fa-tag me-1"></i>
                  ~${item.estimated_price}
                </p>
              {/if}
              
              {#if item.notes}
                <p class="card-text small text-muted mb-3">{item.notes}</p>
              {/if}
              
              <div class="d-flex justify-content-between align-items-center">
                <small class="text-muted">
                  Added {new Date(item.created_at || '').toLocaleDateString()}
                </small>
                <button 
                  class="btn btn-sm btn-outline-danger"
                  on:click={() => removeFromWishlist(item.id)}
                  title="Remove from wishlist"
                >
                  <i class="fas fa-trash"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<!-- TODO: Add modal for adding new wishlist items -->
{#if showAddModal}
  <div class="modal fade show d-block" style="background-color: rgba(0,0,0,0.5);">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Add to Wishlist</h5>
          <button type="button" class="btn-close" on:click={() => showAddModal = false}></button>
        </div>
        <div class="modal-body">
          <p class="text-muted">Feature coming soon! Add new games to your wishlist.</p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" on:click={() => showAddModal = false}>Close</button>
        </div>
      </div>
    </div>
  </div>
{/if}
