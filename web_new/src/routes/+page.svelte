<script lang="ts">
  import { onMount } from 'svelte';
  import type { Game, Platform } from '$lib/models';
  import GameCard from '$lib/components/GameCard.svelte';
  import GameFormModal from '$lib/components/GameFormModal.svelte';

  let games: Game[] = [];
  let platforms: Platform[] = [];
  let loading = true;
  let error: string | null = null;
  let showModal = false;
  let selectedGame: Game | null = null;

  onMount(async () => {
    try {
      const [gamesResponse, platformsResponse] = await Promise.all([
        fetch('/api/v1/games'),
        fetch('/api/v1/platforms'),
      ]);

      if (gamesResponse.ok && platformsResponse.ok) {
        const gamesData = await gamesResponse.json();
        const platformsData = await platformsResponse.json();
        games = gamesData.games;
        platforms = platformsData;
      } else {
        error = 'Failed to load data from the backend.';
      }
    } catch (e) {
      error = 'Error connecting to the backend.';
      console.error(e);
    } finally {
      loading = false;
    }
  });

  function openAddGameModal() {
    selectedGame = null;
    showModal = true;
  }

  function openEditGameModal(event: CustomEvent<number>) {
    const gameId = event.detail;
    selectedGame = games.find((g) => g.id === gameId) || null;
    showModal = true;
  }

  async function handleFormSubmit(event: CustomEvent<Partial<Game>>) {
    const gameData = event.detail;
    const isEdit = !!selectedGame;
    const url = isEdit ? `/api/v1/games/${selectedGame.id}` : '/api/v1/games';
    const method = isEdit ? 'PUT' : 'POST';

    try {
      const response = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(gameData),
      });

      if (response.ok) {
        showModal = false;
        // Refresh the game list
        const gamesResponse = await fetch('/api/v1/games');
        const gamesData = await gamesResponse.json();
        games = gamesData.games;
      } else {
        const errorData = await response.json();
        alert(`Error: ${errorData.error}`);
      }
    } catch (e) {
      alert('An error occurred while saving the game.');
      console.error(e);
    }
  }
</script>

<div class="container-fluid mt-4">
  <div class="d-flex justify-content-between align-items-center mb-4">
    <h1>Game Collection</h1>
    <button class="btn btn-primary" on:click={openAddGameModal}>
      <i class="fas fa-plus me-1"></i> Add Game
    </button>
  </div>

  {#if loading}
    <p>Loading games...</p>
  {:else if error}
    <div class="alert alert-danger">{error}</div>
  {:else}
    <div class="row">
      {#each games as game (game.id)}
        <div class="col-lg-3 col-md-4 col-sm-6 mb-4">
          <GameCard {game} on:edit={openEditGameModal} />
        </div>
      {/each}
    </div>
  {/if}
</div>

{#if showModal}
  <GameFormModal game={selectedGame} {platforms} on:submit={handleFormSubmit} on:close={() => showModal = false} />
{/if}
