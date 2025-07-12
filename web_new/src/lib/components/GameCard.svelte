<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Game } from '$lib/models';

  export let game: Game;

  const dispatch = createEventDispatcher();

  function editGame(event: MouseEvent) {
    event.stopPropagation();
    dispatch('edit', game.id);
  }

  function deleteGame(event: MouseEvent) {
    event.stopPropagation();
    dispatch('delete', game.id);
  }

  function logSession(event: MouseEvent) {
    event.stopPropagation();
    dispatch('log', game.id);
  }
</script>

<button class="card game-card" on:click={() => dispatch('view', game.id)}>
  <div class="position-relative">
    <div class="game-cover">
      {#if game.cover_art_url}
        <img src={game.cover_art_url} alt={game.title} />
      {:else}
        <i class="fas fa-gamepad"></i>
      {/if}
    </div>
    <span class="platform-badge">{game.platform?.name || 'Unknown'}</span>
    <div class="btn-group-actions position-absolute bottom-0 end-0 p-2">
      <button class="btn btn-sm btn-primary" on:click={editGame} title="Edit Game">
        <i class="fas fa-edit"></i>
      </button>
      <button class="btn btn-sm btn-danger" on:click={deleteGame} title="Delete Game">
        <i class="fas fa-trash"></i>
      </button>
      <button class="btn btn-sm btn-success" on:click={logSession} title="Log Session">
        <i class="fas fa-play"></i>
      </button>
    </div>
  </div>
  <div class="card-body">
    <h6 class="card-title">{game.title}</h6>
    <p class="card-text small text-muted">
      {game.year || 'Unknown Year'} &bull; {game.genre || 'Unknown Genre'}
    </p>
  </div>
</button>

<style>
  .game-card {
    cursor: pointer;
    text-align: left;
    background: none;
    border: none;
    padding: 0;
  }
</style>
