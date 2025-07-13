<script lang="ts">
  import { onMount } from 'svelte';
  import type { Game, Platform } from '$lib/models';
  import { api, ApiError } from '$lib/api';
  import GameCard from '$lib/components/GameCard.svelte';
  import GameFormModal from '$lib/components/GameFormModal.svelte';
  import GameDetailModal from '$lib/components/GameDetailModal.svelte';

  let games: Game[] = [];
  let platforms: Platform[] = [];
  let loading = true;
  let error: string | null = null;
  let showModal = false;
  let showDetailModal = false;
  let selectedGame: Game | null = null;
  let searchQuery = '';
  let filteredGames: Game[] = [];

  $: filteredGames = searchQuery 
    ? games.filter(game => 
        game.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        game.genre?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        game.platform?.name?.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : games;

  onMount(async () => {
    await loadData();
  });

  async function loadData() {
    try {
      loading = true;
      error = null;
      
      const [gamesData, platformsData] = await Promise.all([
        api.getGames(),
        api.getPlatforms(),
      ]);

      games = gamesData || [];
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
    showDetailModal = false;  // Close detail modal first
    showModal = true;         // Open edit modal
  }

  async function handleDeleteGame(event: CustomEvent<number>) {
    const gameId = event.detail;
    
    try {
      await api.deleteGame(gameId);
      // Remove the game from the local array
      games = games.filter(g => g.id !== gameId);
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
      let savedGame: Game;
      
      if (isEdit && selectedGame) {
        savedGame = await api.updateGame(selectedGame.id, gameData);
        // Update the game in the local array
        const index = games.findIndex(g => g.id === selectedGame.id);
        if (index !== -1) {
          games[index] = savedGame;
          games = [...games]; // Trigger reactivity
        }
      } else {
        savedGame = await api.createGame(gameData);
        // Add the new game to the local array
        games = [...games, savedGame];
      }
      
      showModal = false;
    } catch (e) {
      if (e instanceof ApiError) {
        alert(`Failed to save game: ${e.message}`);
      } else {
        alert('An error occurred while saving the game.');
      }
      console.error('Save game error:', e);
    }
  }

  function handleLogSession(event: CustomEvent<number>) {
    const gameId = event.detail;
    // TODO: Implement session logging modal
    console.log('Log session for game:', gameId);
  }

  function handleViewGame(event: CustomEvent<Game>) {
    selectedGame = event.detail;
    showDetailModal = true;
  }
</script>

<div class="container-fluid mt-4">
  <div class="d-flex justify-content-between align-items-center mb-4">
    <h1>Game Collection</h1>
    <div class="d-flex gap-2">
      <button class="btn btn-outline-secondary" on:click={loadData} disabled={loading}>
        <i class="fas fa-sync-alt {loading ? 'fa-spin' : ''}"></i> Refresh
      </button>
      <button class="btn btn-primary" on:click={openAddGameModal}>
        <i class="fas fa-plus me-1"></i> Add Game
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
          placeholder="Search games by title, genre, or platform..."
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
          Showing {filteredGames.length} of {games.length} games
        {:else}
          {games.length} games total
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
        <p class="text-muted">Loading games...</p>
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
  {:else if filteredGames.length === 0}
    <div class="text-center py-5">
      {#if searchQuery}
        <i class="fas fa-search fa-3x text-muted mb-3"></i>
        <h4>No games found</h4>
        <p class="text-muted">No games match your search "{searchQuery}"</p>
        <button class="btn btn-outline-primary" on:click={() => searchQuery = ''}>
          Clear Search
        </button>
      {:else}
        <i class="fas fa-gamepad fa-3x text-muted mb-3"></i>
        <h4>No games in your collection</h4>
        <p class="text-muted">Start building your collection by adding your first game!</p>
        <button class="btn btn-primary" on:click={openAddGameModal}>
          <i class="fas fa-plus me-1"></i> Add Your First Game
        </button>
      {/if}
    </div>
  {:else}
    <div class="row">
      {#each filteredGames as game (game.id)}
        <div class="col-xl-2 col-lg-3 col-md-4 col-sm-6 mb-4">
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
  {/if}
</div>

<GameFormModal 
  game={selectedGame} 
  {platforms} 
  show={showModal}
  on:submit={handleFormSubmit} 
  on:close={() => showModal = false} 
/>

{#if showDetailModal && selectedGame}
  <GameDetailModal 
    game={selectedGame} 
    show={showDetailModal}
    on:close={() => showDetailModal = false}
    on:edit={openEditGameModal}
    on:delete={handleDeleteGame}
    on:log={handleLogSession}
  />
{/if}
