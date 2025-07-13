<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Game } from '$lib/models';

	export let game: Game;
	export let show = false;

	const dispatch = createEventDispatcher();

	function closeModal() {
		show = false;
		dispatch('close');
	}

	function editGame() {
		dispatch('edit', game);
		// Don't close modal here - let parent handle modal state
	}

	function deleteGame() {
		if (confirm(`Are you sure you want to delete "${game.title}"?`)) {
			dispatch('delete', game.id);
			closeModal();
		}
	}

	function logSession() {
		dispatch('log', game.id);
		closeModal();
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

{#if show}
	<div class="modal fade show d-block" tabindex="-1" style="background-color: rgba(0,0,0,0.5);">
		<div class="modal-dialog modal-lg">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title">
						<i class="fas fa-gamepad me-2"></i>{game.title}
					</h5>
					<button type="button" class="btn-close" on:click={closeModal}></button>
				</div>
				<div class="modal-body">
					<div class="row">
						<div class="col-md-4">
							<!-- Game Cover -->
							<div class="text-center mb-3">
								{#if game.cover_art_url}
									<img 
										src={game.cover_art_url} 
										alt={game.title}
										class="img-fluid rounded shadow"
										style="max-height: 300px;"
									/>
								{:else}
									<div class="bg-light rounded d-flex align-items-center justify-content-center" style="height: 300px;">
										<i class="fas fa-gamepad fa-3x text-muted"></i>
									</div>
								{/if}
							</div>

							<!-- Quick Actions -->
							<div class="d-grid gap-2">
								<button class="btn btn-primary" on:click={logSession}>
									<i class="fas fa-play me-1"></i>Log Play Session
								</button>
								<button class="btn btn-outline-secondary" on:click={editGame}>
									<i class="fas fa-edit me-1"></i>Edit Game
								</button>
								<button class="btn btn-outline-danger" on:click={deleteGame}>
									<i class="fas fa-trash me-1"></i>Delete Game
								</button>
							</div>
						</div>

						<div class="col-md-8">
							<!-- Game Details -->
							<div class="row mb-3">
								<div class="col-sm-6">
									<h6 class="text-muted">Platform</h6>
									<p class="mb-0">
										<i class="fas fa-desktop me-1"></i>
										{game.platform?.name || 'Unknown'}
									</p>
								</div>
								<div class="col-sm-6">
									<h6 class="text-muted">Year</h6>
									<p class="mb-0">
										<i class="fas fa-calendar me-1"></i>
										{game.year || 'Unknown'}
									</p>
								</div>
							</div>

							<div class="row mb-3">
								<div class="col-sm-6">
									<h6 class="text-muted">Genre</h6>
									<p class="mb-0">
										<i class="fas fa-tag me-1"></i>
										{game.genre || 'Unknown'}
									</p>
								</div>
								<div class="col-sm-6">
									<h6 class="text-muted">Rating</h6>
									<p class="mb-0">
										<i class="fas fa-star me-1"></i>
										{#if game.rating && game.rating > 0}
											{game.rating.toFixed(1)}/10
										{:else}
											Not rated
										{/if}
									</p>
								</div>
							</div>

							<div class="row mb-3">
								<div class="col-sm-6">
									<h6 class="text-muted">Completion Status</h6>
									<p class="mb-0 {getCompletionStatusColor(game.completion_status)}">
										<i class="fas fa-flag me-1"></i>
										{getCompletionStatusText(game.completion_status)}
									</p>
								</div>
								{#if game.completion_percentage && game.completion_percentage > 0}
									<div class="col-sm-6">
										<h6 class="text-muted">Completion</h6>
										<div class="progress" style="height: 20px;">
											<div 
												class="progress-bar" 
												role="progressbar" 
												style="width: {game.completion_percentage}%"
											>
												{game.completion_percentage}%
											</div>
										</div>
									</div>
								{/if}
							</div>

							{#if game.description}
								<div class="mb-3">
									<h6 class="text-muted">Description</h6>
									<p class="text-justify">{game.description}</p>
								</div>
							{/if}

							{#if game.file_locations && game.file_locations.length > 0}
								<div class="mb-3">
									<h6 class="text-muted">File Locations</h6>
									{#each game.file_locations as location}
										<div class="card card-body bg-light mb-2">
											<div class="d-flex justify-content-between align-items-start">
												<div>
													<small class="text-muted">Server:</small>
													<div>{location.server_location || 'Unknown'}</div>
													<small class="text-muted">Path:</small>
													<div class="font-monospace small">{location.file_path}</div>
												</div>
												<div class="text-end">
													{#if location.file_size}
														<small class="text-muted">
															{(location.file_size / 1024 / 1024).toFixed(1)} MB
														</small>
													{/if}
												</div>
											</div>
										</div>
									{/each}
								</div>
							{/if}

							{#if game.developer || game.publisher}
								<div class="row mb-3">
									{#if game.developer}
										<div class="col-sm-6">
											<h6 class="text-muted">Developer</h6>
											<p class="mb-0">{game.developer}</p>
										</div>
									{/if}
									{#if game.publisher}
										<div class="col-sm-6">
											<h6 class="text-muted">Publisher</h6>
											<p class="mb-0">{game.publisher}</p>
										</div>
									{/if}
								</div>
							{/if}

							{#if game.igdb_id}
								<div class="mb-3">
									<h6 class="text-muted">External IDs</h6>
									<p class="mb-0">
										<small class="text-muted">IGDB ID:</small> {game.igdb_id}
									</p>
								</div>
							{/if}
						</div>
					</div>
				</div>
				<div class="modal-footer">
					<button type="button" class="btn btn-secondary" on:click={closeModal}>
						Close
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}