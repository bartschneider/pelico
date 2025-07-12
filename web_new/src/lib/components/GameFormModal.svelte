<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';
  import type { Game, Platform } from '$lib/models';
  import { bootstrapModal } from '$lib/actions/modal';

  export let game: Game | null = null;
  export let platforms: Platform[] = [];
  export let show: boolean;

  let title = '';
  let platform_id: number | null = null;
  let year: number | null = null;
  let genre = '';
  let description = '';
  let cover_art_url = '';
  let rating: number | null = null;
  let purchase_date = '';
  let collection_formats: string[] = [];

  const dispatch = createEventDispatcher();

  onMount(() => {
    if (game) {
      title = game.title;
      platform_id = game.platform_id;
      year = game.year;
      genre = game.genre;
      description = game.description;
      cover_art_url = game.cover_art_url;
      rating = game.rating;
      purchase_date = game.purchase_date ? new Date(game.purchase_date).toISOString().split('T')[0] : '';
      collection_formats = game.collection_formats || [];
    }
  });

  function handleSubmit() {
    const gameData = {
      title,
      platform_id,
      year,
      genre,
      description,
      cover_art_url,
      rating,
      purchase_date,
      collection_formats,
    };
    dispatch('submit', gameData);
  }
</script>

<div class="modal fade" tabindex="-1" use:bootstrapModal={{ show }} on:close={() => dispatch('close')}>
  <div class="modal-dialog modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title">{game ? 'Edit Game' : 'Add New Game'}</h5>
        <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
      </div>
      <div class="modal-body">
        <form on:submit|preventDefault={handleSubmit}>
          <div class="mb-3">
            <label for="gameTitle" class="form-label">Title</label>
            <input type="text" class="form-control" id="gameTitle" bind:value={title} required />
          </div>
          <div class="mb-3">
            <label for="gamePlatform" class="form-label">Platform</label>
            <select class="form-select" id="gamePlatform" bind:value={platform_id} required>
              <option value={null}>Select Platform</option>
              {#each platforms as platform}
                <option value={platform.id}>{platform.name}</option>
              {/each}
            </select>
          </div>
          <div class="row">
            <div class="col-md-6">
              <div class="mb-3">
                <label for="gameYear" class="form-label">Year</label>
                <input type="number" class="form-control" id="gameYear" bind:value={year} />
              </div>
            </div>
            <div class="col-md-6">
              <div class="mb-3">
                <label for="gameGenre" class="form-label">Genre</label>
                <input type="text" class="form-control" id="gameGenre" bind:value={genre} />
              </div>
            </div>
          </div>
          <div class="mb-3">
            <label for="gameDescription" class="form-label">Description</label>
            <textarea class="form-control" id="gameDescription" rows="3" bind:value={description}></textarea>
          </div>
          <div class="mb-3">
            <label for="gameCoverUrl" class="form-label">Cover Art URL</label>
            <input type="url" class="form-control" id="gameCoverUrl" bind:value={cover_art_url} />
          </div>
          <div class="row">
            <div class="col-md-6">
              <div class="mb-3">
                <label for="gameRating" class="form-label">Rating (0-10)</label>
                <input type="number" class="form-control" id="gameRating" min="0" max="10" step="0.1" bind:value={rating} />
              </div>
            </div>
            <div class="col-md-6">
              <div class="mb-3">
                <label for="gamePurchaseDate" class="form-label">Purchase Date</label>
                <input type="date" class="form-control" id="gamePurchaseDate" bind:value={purchase_date} />
              </div>
            </div>
          </div>
          <div class="mb-3" role="group" aria-labelledby="collectionFormatsLabel">
            <div id="collectionFormatsLabel" class="form-label">Collection Format(s)</div>
            <div class="form-check">
              <input class="form-check-input" type="checkbox" id="gameFormatPhysical" value="physical" bind:group={collection_formats} />
              <label class="form-check-label" for="gameFormatPhysical">
                <i class="fas fa-box me-1"></i>Physical Copy
              </label>
            </div>
            <div class="form-check">
              <input class="form-check-input" type="checkbox" id="gameFormatDigital" value="digital" bind:group={collection_formats} />
              <label class="form-check-label" for="gameFormatDigital">
                <i class="fas fa-download me-1"></i>Digital Purchase
              </label>
            </div>
            <div class="form-check">
              <input class="form-check-input" type="checkbox" id="gameFormatRom" value="rom" bind:group={collection_formats} />
              <label class="form-check-label" for="gameFormatRom">
                <i class="fas fa-hdd me-1"></i>ROM File
              </label>
            </div>
          </div>
        </form>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" on:click={() => dispatch('close')}>Cancel</button>
        <button type="button" class="btn btn-primary" on:click={handleSubmit}>{game ? 'Save Changes' : 'Add Game'}</button>
      </div>
    </div>
  </div>
</div>
