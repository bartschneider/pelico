<script lang="ts">
  import { onMount } from 'svelte';
  import type { Shortlist, Game } from '$lib/models';
  import GameCard from '$lib/components/GameCard.svelte';

  let shortlist: Shortlist[] = [];
  let loading = true;
  let error: string | null = null;

  onMount(async () => {
    try {
      const response = await fetch('/api/v1/shortlist');
      if (response.ok) {
        shortlist = await response.json();
      } else {
        error = 'Failed to load shortlist from the backend.';
      }
    } catch (e) {
      error = 'Error connecting to the backend.';
      console.error(e);
    } finally {
      loading = false;
    }
  });

  async function removeFromShortlist(gameId: number) {
    try {
      const response = await fetch(`/api/v1/shortlist/${gameId}`, {
        method: 'DELETE',
      });

      if (response.ok) {
        shortlist = shortlist.filter((item) => item.game.id !== gameId);
      } else {
        alert('Failed to remove game from shortlist.');
      }
    } catch (e) {
      alert('An error occurred while removing the game from the shortlist.');
      console.error(e);
    }
  }
</script>

<div class="container-fluid mt-4">
  <h1>Shortlist</h1>

  {#if loading}
    <p>Loading shortlist...</p>
  {:else if error}
    <div class="alert alert-danger">{error}</div>
  {:else}
    <div class="row">
      {#each shortlist as item (item.id)}
        <div class="col-lg-3 col-md-4 col-sm-6 mb-4">
          <GameCard game={item.game} on:delete={() => removeFromShortlist(item.game.id)} />
        </div>
      {/each}
    </div>
  {/if}
</div>
