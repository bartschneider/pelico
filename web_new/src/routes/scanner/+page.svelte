<script lang="ts">
	import { onMount } from 'svelte';
	import type { Platform } from '$lib/models';
	import { api } from '$lib/api';
	import { showSuccess, showError, showInfo } from '$lib/stores/notifications';

	let allPlatforms: Platform[] = [];
	let directoryPath = '';
	let serverLocation = '';
	let selectedPlatformId: number | null = null;
	let recursive = true;
	let scanning = false;
	let scanResult: any = null;
	let showDirectoryBrowser = false;
	let availableDirectories: string[] = [];
	let loadingDirectories = false;

	onMount(async () => {
		await loadPlatforms();
		await loadSuggestedPaths();
	});

	async function loadPlatforms() {
		try {
			allPlatforms = await api.getPlatforms();
		} catch (error) {
			showError('Error', 'Failed to load platforms');
		}
	}

	async function loadSuggestedPaths() {
		try {
			const suggestions = await api.getSuggestedPaths();
			availableDirectories = suggestions.paths || [];
		} catch (error) {
			console.warn('Could not load suggested paths');
		}
	}

	async function scanDirectory() {
		if (!directoryPath || !serverLocation || !selectedPlatformId) {
			showError('Validation Error', 'Please fill in all required fields');
			return;
		}

		scanning = true;
		scanResult = null;

		try {
			const result = await api.scanDirectory({
				path: directoryPath,
				server_location: serverLocation,
				platform_id: selectedPlatformId,
				recursive
			});

			scanResult = result;
			showSuccess('Scan Complete', `Found ${result.games_found} games`);
		} catch (error: any) {
			showError('Scan Failed', error.message);
		} finally {
			scanning = false;
		}
	}

	function openDirectoryBrowser() {
		showDirectoryBrowser = true;
	}

	function selectDirectory(dir: string) {
		directoryPath = dir;
		showDirectoryBrowser = false;
	}

	function closeDirectoryBrowser() {
		showDirectoryBrowser = false;
	}
</script>

<div class="container-fluid py-4">
	<div class="row">
		<div class="col-12">
			<div class="d-flex justify-content-between align-items-center mb-4">
				<div>
					<h1 class="h2 mb-1">
						<i class="fas fa-search me-2 text-primary"></i>ROM Scanner
					</h1>
					<p class="text-muted mb-0">Scan directories for ROM files and add them to your collection</p>
				</div>
			</div>
		</div>
	</div>

	<div class="row">
		<div class="col-lg-8">
			<!-- Scan Configuration -->
			<div class="card mb-4">
				<div class="card-header">
					<h5 class="mb-0">
						<i class="fas fa-cog me-2"></i>
						Scan Configuration
					</h5>
				</div>
				<div class="card-body">
					<form on:submit|preventDefault={scanDirectory}>
						<div class="row">
							<div class="col-md-8">
								<div class="mb-3">
									<label for="directoryPath" class="form-label">
										Directory Path <span class="text-danger">*</span>
									</label>
									<div class="input-group">
										<input
											type="text"
											id="directoryPath"
											bind:value={directoryPath}
											required
											readonly
											placeholder="/path/to/roms"
											class="form-control"
										/>
										<button
											type="button"
											on:click={openDirectoryBrowser}
											class="btn btn-outline-secondary"
										>
											<i class="fas fa-folder-open"></i> Browse
										</button>
									</div>
								</div>

								<div class="mb-3">
									<label for="serverLocation" class="form-label">
										Server Location <span class="text-danger">*</span>
									</label>
									<input
										type="text"
										id="serverLocation"
										bind:value={serverLocation}
										required
										placeholder="Home Server"
										class="form-control"
									/>
									<div class="form-text">Name of the server or storage device</div>
								</div>
							</div>

							<div class="col-md-4">
								<div class="mb-3">
									<label for="platformSelect" class="form-label">
										Platform <span class="text-danger">*</span>
									</label>
									<select
										id="platformSelect"
										bind:value={selectedPlatformId}
										required
										class="form-select"
									>
										<option value="">Select Platform</option>
										{#each allPlatforms as platform}
											<option value={platform.id}>{platform.name}</option>
										{/each}
									</select>
								</div>

								<div class="form-check">
									<input
										class="form-check-input"
										type="checkbox"
										id="recursiveCheck"
										bind:checked={recursive}
									/>
									<label class="form-check-label" for="recursiveCheck">
										Scan subdirectories
									</label>
								</div>
							</div>
						</div>

						<div class="d-grid gap-2 d-md-flex justify-content-md-end">
							<button
								type="submit"
								disabled={scanning}
								class="btn btn-primary"
							>
								{#if scanning}
									<i class="fas fa-spinner fa-spin me-1"></i>Scanning...
								{:else}
									<i class="fas fa-search me-1"></i>Start Scan
								{/if}
							</button>
						</div>
					</form>
				</div>
			</div>

			<!-- Scan Results -->
			{#if scanResult}
				<div class="card">
					<div class="card-header">
						<h5 class="mb-0">
							<i class="fas fa-check-circle me-2 text-success"></i>
							Scan Results
						</h5>
					</div>
					<div class="card-body">
						<div class="row text-center">
							<div class="col-md-4">
								<div class="border rounded p-3">
									<i class="fas fa-gamepad fa-2x text-primary mb-2"></i>
									<h4 class="mb-0">{scanResult.games_found}</h4>
									<small class="text-muted">Games Found</small>
								</div>
							</div>
							<div class="col-md-4">
								<div class="border rounded p-3">
									<i class="fas fa-plus-circle fa-2x text-success mb-2"></i>
									<h4 class="mb-0">{scanResult.new_games || 0}</h4>
									<small class="text-muted">New Games</small>
								</div>
							</div>
							<div class="col-md-4">
								<div class="border rounded p-3">
									<i class="fas fa-sync-alt fa-2x text-info mb-2"></i>
									<h4 class="mb-0">{scanResult.updated_games || 0}</h4>
									<small class="text-muted">Updated Games</small>
								</div>
							</div>
						</div>

						{#if scanResult.message}
							<div class="alert alert-info mt-3">
								<i class="fas fa-info-circle me-2"></i>
								{scanResult.message}
							</div>
						{/if}
					</div>
				</div>
			{/if}
		</div>

		<div class="col-lg-4">
			<!-- Quick Actions -->
			<div class="card mb-4">
				<div class="card-header">
					<h5 class="mb-0">
						<i class="fas fa-bolt me-2"></i>Quick Actions
					</h5>
				</div>
				<div class="card-body">
					<div class="d-grid gap-2">
						<a href="/platforms" class="btn btn-outline-primary">
							<i class="fas fa-desktop me-1"></i>Manage Platforms
						</a>
						<a href="/collection" class="btn btn-outline-secondary">
							<i class="fas fa-gamepad me-1"></i>View Collection
						</a>
						<button 
							type="button" 
							class="btn btn-outline-info"
							on:click={loadSuggestedPaths}
						>
							<i class="fas fa-refresh me-1"></i>Refresh Paths
						</button>
					</div>
				</div>
			</div>

			<!-- Scan Tips -->
			<div class="card">
				<div class="card-header">
					<h5 class="mb-0">
						<i class="fas fa-lightbulb me-2"></i>Scanning Tips
					</h5>
				</div>
				<div class="card-body">
					<ul class="list-unstyled mb-0">
						<li class="mb-2">
							<i class="fas fa-check text-success me-2"></i>
							<small>Use recursive scan for nested folders</small>
						</li>
						<li class="mb-2">
							<i class="fas fa-check text-success me-2"></i>
							<small>Organize ROMs by platform for better results</small>
						</li>
						<li class="mb-2">
							<i class="fas fa-check text-success me-2"></i>
							<small>Scanner supports ZIP and 7Z archives</small>
						</li>
						<li class="mb-0">
							<i class="fas fa-check text-success me-2"></i>
							<small>Metadata is fetched automatically</small>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- Directory Browser Modal -->
{#if showDirectoryBrowser}
	<div class="modal fade show d-block" tabindex="-1" style="background-color: rgba(0,0,0,0.5);">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title">
						<i class="fas fa-folder me-2"></i>Browse Directories
					</h5>
					<button type="button" class="btn-close" on:click={closeDirectoryBrowser}></button>
				</div>
				<div class="modal-body">
					{#if availableDirectories.length > 0}
						<p class="text-muted">Select a directory from the available paths:</p>
						<div class="list-group">
							{#each availableDirectories as dir}
								<button
									type="button"
									class="list-group-item list-group-item-action"
									on:click={() => selectDirectory(dir)}
								>
									<i class="fas fa-folder me-2"></i>{dir}
								</button>
							{/each}
						</div>
					{:else}
						<div class="text-center py-4">
							<i class="fas fa-folder-open fa-3x text-muted mb-3"></i>
							<p class="text-muted">No suggested directories available.</p>
							<p class="small">You can manually enter a directory path in the form.</p>
						</div>
					{/if}
				</div>
				<div class="modal-footer">
					<button type="button" class="btn btn-secondary" on:click={closeDirectoryBrowser}>
						Cancel
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}