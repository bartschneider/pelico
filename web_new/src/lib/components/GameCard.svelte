<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Game } from '$lib/models';

  export let game: Game;
  export let loading = false;

  const dispatch = createEventDispatcher();

  let imageError = false;

  function editGame(event: MouseEvent) {
    event.stopPropagation();
    dispatch('edit', game);
  }

  function deleteGame(event: MouseEvent) {
    event.stopPropagation();
    if (confirm(`Are you sure you want to delete "${game.title}"?`)) {
      dispatch('delete', game.id);
    }
  }

  function logSession(event: MouseEvent) {
    event.stopPropagation();
    dispatch('log', game.id);
  }

  function handleImageError() {
    imageError = true;
  }

  function viewGame() {
    dispatch('view', game);
  }

  function getCompletionStatusColor(status: string) {
    switch (status) {
      case 'completed': return 'text-success';
      case 'in_progress': return 'text-warning';
      case 'abandoned': return 'text-danger';
      default: return 'text-muted';
    }
  }

  function getCompletionStatusText(status: string) {
    switch (status) {
      case 'not_started': return 'Not Started';
      case 'in_progress': return 'In Progress';
      case 'completed': return 'Completed';
      case 'abandoned': return 'Abandoned';
      default: return 'Unknown';
    }
  }
</script>

{#if loading}
  <div class="card game-card loading">
    <div class="d-flex justify-content-center align-items-center" style="height: 200px;">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>
  </div>
{:else}
  <div class="card game-card" role="button" tabindex="0" 
       on:click={viewGame}
       on:keydown={(e) => e.key === 'Enter' && viewGame()}>
    <div class="position-relative">
      <div class="game-cover">
        {#if game.cover_art_url && !imageError}
          <img 
            src={game.cover_art_url} 
            alt={game.title}
            on:error={handleImageError}
            loading="lazy"
          />
        {:else}
          <div class="placeholder-cover d-flex align-items-center justify-content-center">
            <i class="fas fa-gamepad fa-2x text-muted"></i>
          </div>
        {/if}
      </div>
      
      <span class="platform-badge">
        {game.platform?.name || 'Unknown Platform'}
      </span>
      
      <div class="completion-badge">
        <span class="badge {getCompletionStatusColor(game.completion_status)}">
          {getCompletionStatusText(game.completion_status)}
        </span>
      </div>
      
      <div class="btn-group-actions position-absolute bottom-0 end-0 p-2">
        <button 
          class="btn btn-sm btn-outline-light shadow-sm" 
          on:click={editGame} 
          title="Edit Game"
          aria-label="Edit {game.title}"
        >
          <i class="fas fa-edit"></i>
        </button>
        <button 
          class="btn btn-sm btn-outline-danger shadow-sm" 
          on:click={deleteGame} 
          title="Delete Game"
          aria-label="Delete {game.title}"
        >
          <i class="fas fa-trash"></i>
        </button>
        <button 
          class="btn btn-sm btn-outline-success shadow-sm" 
          on:click={logSession} 
          title="Log Session"
          aria-label="Log session for {game.title}"
        >
          <i class="fas fa-play"></i>
        </button>
      </div>
    </div>
    
    <div class="card-body">
      <h6 class="card-title mb-1" title={game.title}>{game.title}</h6>
      <p class="card-text small text-muted mb-1">
        {#if game.release_date}
          {new Date(game.release_date).getFullYear()}
        {:else}
          Unknown Year
        {/if}
        {#if game.genre}
          &bull; {game.genre}
        {/if}
      </p>
      
      {#if game.rating}
        <div class="rating mb-1">
          {#each Array(5) as _, i}
            <i class="fas fa-star {i < game.rating ? 'text-warning' : 'text-light'}"></i>
          {/each}
        </div>
      {/if}
      
      {#if game.playtime_hours}
        <small class="text-muted">
          <i class="fas fa-clock"></i> {game.playtime_hours}h played
        </small>
      {/if}
    </div>
  </div>
{/if}

<style>
  .game-card {
    cursor: pointer;
    transition: transform 0.2s ease, box-shadow 0.2s ease;
    height: 100%;
  }

  .game-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  }

  .game-card.loading {
    cursor: default;
  }

  .game-cover {
    height: 200px;
    background: #f8f9fa;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    border-radius: 0.375rem 0.375rem 0 0;
  }

  .game-cover img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .placeholder-cover {
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  }

  .platform-badge {
    position: absolute;
    top: 8px;
    left: 8px;
    background: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 2px 8px;
    border-radius: 12px;
    font-size: 0.75rem;
    font-weight: 500;
  }

  .completion-badge {
    position: absolute;
    top: 8px;
    right: 8px;
  }

  .btn-group-actions {
    opacity: 0;
    transition: opacity 0.2s ease;
    background: linear-gradient(transparent, rgba(0, 0, 0, 0.5));
    border-radius: 0 0 0.375rem 0.375rem;
    left: 0;
    right: 0;
  }

  .game-card:hover .btn-group-actions {
    opacity: 1;
  }

  .btn-group-actions .btn {
    margin-left: 4px;
    backdrop-filter: blur(10px);
  }

  .card-title {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .rating .fa-star {
    font-size: 0.8rem;
  }
</style>
