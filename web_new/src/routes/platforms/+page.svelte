<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { showSuccess, showError } from '$lib/stores/notifications';
	import type { Platform } from '$lib/models';

	let platforms: Platform[] = [];
	let loading = false;
	let showAddModal = false;
	let showEditModal = false;
	let editingPlatform: Platform | null = null;

	// Form data
	let platformName = '';
	let platformManufacturer = '';
	let platformYear: number | undefined;

	onMount(async () => {
		await loadPlatforms();
	});

	async function loadPlatforms() {
		loading = true;
		try {
			platforms = await api.getPlatforms();
		} catch (error: any) {
			showError('Error', 'Failed to load platforms');
		} finally {
			loading = false;
		}
	}

	function openAddModal() {
		platformName = '';
		platformManufacturer = '';
		platformYear = undefined;
		showAddModal = true;
	}

	function openEditModal(platform: Platform) {
		editingPlatform = platform;
		platformName = platform.name;
		platformManufacturer = platform.manufacturer || '';
		platformYear = platform.release_year || undefined;
		showEditModal = true;
	}

	function closeModals() {
		showAddModal = false;
		showEditModal = false;
		editingPlatform = null;
		platformName = '';
		platformManufacturer = '';
		platformYear = undefined;
	}

	async function handleAddPlatform() {
		if (!platformName.trim()) {
			showError('Validation Error', 'Platform name is required');
			return;
		}

		try {
			const newPlatform = await api.createPlatform({
				name: platformName.trim(),
				manufacturer: platformManufacturer.trim() || undefined,
				release_year: platformYear || 0
			});
			
			platforms = [...platforms, newPlatform];
			showSuccess('Success', 'Platform added successfully');
			closeModals();
		} catch (error: any) {
			showError('Error', error.message);
		}
	}

	async function handleUpdatePlatform() {
		if (!editingPlatform || !platformName.trim()) {
			showError('Validation Error', 'Platform name is required');
			return;
		}

		try {
			const updatedPlatform = await api.updatePlatform(editingPlatform.id, {
				name: platformName.trim(),
				manufacturer: platformManufacturer.trim() || undefined,
				release_year: platformYear || 0
			});
			
			platforms = platforms.map(p => 
				p.id === editingPlatform.id ? updatedPlatform : p
			);
			showSuccess('Success', 'Platform updated successfully');
			closeModals();
		} catch (error: any) {
			showError('Error', error.message);
		}
	}

	async function handleDeletePlatform(platform: Platform) {
		if (!confirm(`Are you sure you want to delete "${platform.name}"? This action cannot be undone.`)) {
			return;
		}

		try {
			await api.deletePlatform(platform.id);
			platforms = platforms.filter(p => p.id !== platform.id);
			showSuccess('Success', 'Platform deleted successfully');
		} catch (error: any) {
			showError('Error', error.message);
		}
	}
</script>

<div class="container-fluid py-4">
	<div class="d-flex justify-content-between align-items-center mb-4">
		<div>
			<h1 class="h2 mb-1">
				<i class="fas fa-desktop me-2 text-primary"></i>Platforms
			</h1>
			<p class="text-muted mb-0">Manage gaming platforms in your collection.</p>
		</div>
		<button
			on:click={openAddModal}
			class="btn btn-primary"
		>
			<i class="fas fa-plus me-1"></i>Add Platform
		</button>
	</div>

	{#if loading}
		<div class="text-center py-5">
			<div class="spinner-border text-primary mb-3" role="status">
				<span class="visually-hidden">Loading...</span>
			</div>
			<p class="text-muted">Loading platforms...</p>
		</div>
	{:else if platforms.length === 0}
		<div class="text-center py-5">
			<i class="fas fa-desktop fa-3x text-muted mb-3"></i>
			<h3 class="h5 mb-2">No Platforms Found</h3>
			<p class="text-muted mb-4">Get started by adding your first gaming platform.</p>
			<button
				on:click={openAddModal}
				class="btn btn-primary"
			>
				<i class="fas fa-plus me-1"></i>Add Platform
			</button>
		</div>
	{:else}
		<div class="row">
			{#each platforms as platform}
				<div class="col-lg-4 col-md-6 mb-4">
					<div class="card h-100">
						<div class="card-body">
							<div class="d-flex justify-content-between align-items-start mb-3">
								<div class="flex-grow-1">
									<h5 class="card-title mb-1">{platform.name}</h5>
									{#if platform.manufacturer}
										<p class="card-text text-muted mb-1">{platform.manufacturer}</p>
									{/if}
									{#if platform.release_year && platform.release_year > 0}
										<small class="text-muted">{platform.release_year}</small>
									{/if}
								</div>
								<div class="dropdown">
									<button class="btn btn-sm btn-outline-secondary dropdown-toggle" type="button" data-bs-toggle="dropdown">
										<i class="fas fa-ellipsis-v"></i>
									</button>
									<ul class="dropdown-menu">
										<li>
											<button 
												class="dropdown-item" 
												on:click={() => openEditModal(platform)}
											>
												<i class="fas fa-edit me-2"></i>Edit Platform
											</button>
										</li>
										<li><hr class="dropdown-divider"></li>
										<li>
											<button 
												class="dropdown-item text-danger" 
												on:click={() => handleDeletePlatform(platform)}
											>
												<i class="fas fa-trash me-2"></i>Delete Platform
											</button>
										</li>
									</ul>
								</div>
							</div>
						</div>
						
						<div class="card-footer bg-light">
							<div class="row small text-muted">
								<div class="col-6">
									<span>Platform ID:</span>
								</div>
								<div class="col-6 text-end">
									<span>#{platform.id}</span>
								</div>
							</div>
							<div class="row small text-muted">
								<div class="col-6">
									<span>Added:</span>
								</div>
								<div class="col-6 text-end">
									<span>{new Date(platform.created_at).toLocaleDateString()}</span>
								</div>
							</div>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Add Platform Modal -->
{#if showAddModal}
	<div class="modal fade show d-block" tabindex="-1" style="background-color: rgba(0,0,0,0.5);">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title">
						<i class="fas fa-plus me-2"></i>Add New Platform
					</h5>
					<button type="button" class="btn-close" on:click={closeModals}></button>
				</div>

				<form on:submit|preventDefault={handleAddPlatform}>
					<div class="modal-body">
						<div class="mb-3">
							<label for="platformName" class="form-label">
								Name <span class="text-danger">*</span>
							</label>
							<input
								type="text"
								id="platformName"
								bind:value={platformName}
								required
								class="form-control"
								placeholder="e.g., Nintendo Switch"
							/>
						</div>

						<div class="mb-3">
							<label for="platformManufacturer" class="form-label">
								Manufacturer
							</label>
							<input
								type="text"
								id="platformManufacturer"
								bind:value={platformManufacturer}
								class="form-control"
								placeholder="e.g., Nintendo"
							/>
						</div>

						<div class="mb-3">
							<label for="platformYear" class="form-label">
								Release Year
							</label>
							<input
								type="number"
								id="platformYear"
								bind:value={platformYear}
								min="1970"
								max="2030"
								class="form-control"
								placeholder="e.g., 2017"
							/>
						</div>
					</div>

					<div class="modal-footer">
						<button
							type="button"
							on:click={closeModals}
							class="btn btn-secondary"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="btn btn-primary"
						>
							<i class="fas fa-save me-1"></i>Add Platform
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}

<!-- Edit Platform Modal -->
{#if showEditModal && editingPlatform}
	<div class="modal fade show d-block" tabindex="-1" style="background-color: rgba(0,0,0,0.5);">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title">
						<i class="fas fa-edit me-2"></i>Edit Platform
					</h5>
					<button type="button" class="btn-close" on:click={closeModals}></button>
				</div>

				<form on:submit|preventDefault={handleUpdatePlatform}>
					<div class="modal-body">
						<div class="mb-3">
							<label for="editPlatformName" class="form-label">
								Name <span class="text-danger">*</span>
							</label>
							<input
								type="text"
								id="editPlatformName"
								bind:value={platformName}
								required
								class="form-control"
							/>
						</div>

						<div class="mb-3">
							<label for="editPlatformManufacturer" class="form-label">
								Manufacturer
							</label>
							<input
								type="text"
								id="editPlatformManufacturer"
								bind:value={platformManufacturer}
								class="form-control"
							/>
						</div>

						<div class="mb-3">
							<label for="editPlatformYear" class="form-label">
								Release Year
							</label>
							<input
								type="number"
								id="editPlatformYear"
								bind:value={platformYear}
								min="1970"
								max="2030"
								class="form-control"
							/>
						</div>
					</div>

					<div class="modal-footer">
						<button
							type="button"
							on:click={closeModals}
							class="btn btn-secondary"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="btn btn-primary"
						>
							<i class="fas fa-save me-1"></i>Update Platform
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}