<script lang="ts">
  import { onMount } from 'svelte';
  import type { Wishlist, Game } from '$lib/models';
  import GameCard from '$lib/components/GameCard.svelte';

  let wishlist: Wishlist[] = [];
  let loading = true;
  let error: string | null = null;

  onMount(async () => {
    try {
      const response = await fetch('/api/v1/wishlist');
      if (response.ok) {
        wishlist = await response.json();
      } else {
        error = 'Failed to load wishlist from the backend.';
      }
    } catch (e) {
      error = 'Error connecting to the backend.';
      console.error(e);
    } finally {
      loading = false;
    }
  });

  async function removeFromWishlist(gameId: number) {
    try {
      const response = await fetch(`/api/v1/wishlist/${gameId}`, {
        method: 'DELETE',
      });

      if (response.ok) {
        wishlist = wishlist.filter((item) => item.game.id !== gameId);
      } else {
        alert('Failed to remove game from wishlist.');
      }
    } catch (e) {
      alert('An error occurred while removing the game from the wishlist.');
      console.error(e);
    }
  }
</script>

<div class="container-fluid mt-4">
  <div class="d-flex justify-content-between align-items-center mb-4">
    <h1>Wishlist</h1>
    <button class="btn btn-primary">
      <i class="fas fa-plus me-1"></i> Add Game to Wishlist
    </button>
  </div>

  {#if loading}
    <p>Loading wishlist...</p>
  {:else if error}
    <div class="alert alert-danger">{error}</div>
  {:else}
    <div class="row">
      {#each wishlist as item (item.id)}
        <div class="col-lg-3 col-md-4 col-sm-6 mb-4">
          <GameCard game={item.game} on:delete={() => removeFromWishlist(item.game.id)} />
        </div>
      {/each}
    </div>
  {/if}
</div>
