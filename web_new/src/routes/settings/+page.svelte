<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { showSuccess, showError, showInfo } from '$lib/stores/notifications';

	let backupInfo = {
		games_count: 0,
		platforms_count: 0,
		sessions_count: 0,
		files_count: 0
	};
	
	let versionInfo = {
		version: '',
		build_time: '',
		git_commit: '',
		go_version: ''
	};

	let nextcloudStatus = 'Checking...';
	let nextcloudMessage = '';
	let nextcloudAvailable = false;

	let backupFile: File | null = null;
	let importing = false;
	let exporting = false;
	let exportingGames = false;
	let backingUpToNextcloud = false;
	let exportingAndUploadingGames = false;
	let testingNextcloud = false;

	onMount(async () => {
		await Promise.all([
			loadBackupInfo(),
			loadVersionInfo(),
			checkNextcloudStatus()
		]);
	});

	async function loadBackupInfo() {
		try {
			backupInfo = await api.getBackupInfo();
		} catch (error: any) {
			showError('Error', 'Failed to load backup info');
		}
	}

	async function loadVersionInfo() {
		try {
			versionInfo = await api.getVersion();
		} catch (error: any) {
			showError('Error', 'Failed to load version info');
		}
	}

	async function checkNextcloudStatus() {
		try {
			const result = await api.testNextcloudConnection();
			nextcloudStatus = result.status;
			nextcloudMessage = result.message;
			nextcloudAvailable = result.status === 'success';
		} catch (error: any) {
			nextcloudStatus = 'Not configured';
			nextcloudMessage = 'Nextcloud connection not available';
			nextcloudAvailable = false;
		}
	}

	async function exportFullBackup() {
		exporting = true;
		try {
			const blob = await api.exportDatabase();
			downloadBlob(blob, `pelico-backup-${new Date().toISOString().slice(0, 10)}.json`);
			showSuccess('Export Complete', 'Database backup downloaded successfully');
		} catch (error: any) {
			showError('Export Failed', error.message);
		} finally {
			exporting = false;
		}
	}

	async function exportGamesOnly() {
		exportingGames = true;
		try {
			const blob = await api.exportGames();
			downloadBlob(blob, `pelico-games-${new Date().toISOString().slice(0, 10)}.json`);
			showSuccess('Export Complete', 'Games data downloaded successfully');
		} catch (error: any) {
			showError('Export Failed', error.message);
		} finally {
			exportingGames = false;
		}
	}

	async function importBackup() {
		if (!backupFile) return;
		
		if (!confirm('This will replace all existing data. Are you sure?')) {
			return;
		}

		importing = true;
		try {
			await api.importDatabase(backupFile);
			showSuccess('Import Complete', 'Database restored successfully. Page will reload.');
			setTimeout(() => window.location.reload(), 2000);
		} catch (error: any) {
			showError('Import Failed', error.message);
		} finally {
			importing = false;
		}
	}

	async function backupToNextcloud() {
		backingUpToNextcloud = true;
		try {
			const result = await api.backupToNextcloud();
			showSuccess('Backup Complete', result.message);
		} catch (error: any) {
			showError('Backup Failed', error.message);
		} finally {
			backingUpToNextcloud = false;
		}
	}

	async function exportAndUploadGames() {
		exportingAndUploadingGames = true;
		try {
			// First create a games-only backup
			const blob = await api.exportGames();
			
			// Create FormData to upload the backup
			const formData = new FormData();
			const filename = `pelico-games-${new Date().toISOString().slice(0, 10)}.json`;
			formData.append('backup', blob, filename);
			
			// Upload to Nextcloud via the backup endpoint
			const result = await api.backupToNextcloud();
			showSuccess('Games Upload Complete', `Games data uploaded to Nextcloud successfully`);
		} catch (error: any) {
			showError('Upload Failed', error.message);
		} finally {
			exportingAndUploadingGames = false;
		}
	}

	async function testNextcloudConnection() {
		testingNextcloud = true;
		try {
			const result = await api.testNextcloudConnection();
			if (result.status === 'success') {
				showSuccess('Connection Test', result.message);
			} else {
				showError('Connection Test', result.message);
			}
			nextcloudStatus = result.status;
			nextcloudMessage = result.message;
			nextcloudAvailable = result.status === 'success';
		} catch (error: any) {
			showError('Connection Test Failed', error.message);
			nextcloudStatus = 'error';
			nextcloudMessage = 'Connection failed';
			nextcloudAvailable = false;
		} finally {
			testingNextcloud = false;
		}
	}

	async function refreshAllData() {
		showInfo('Refreshing', 'Clearing cache and refreshing data...');
		try {
			await api.clearCache();
			window.location.reload();
		} catch (error: any) {
			showError('Refresh Failed', error.message);
		}
	}

	function downloadBlob(blob: Blob, filename: string) {
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = filename;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		URL.revokeObjectURL(url);
	}

	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		backupFile = target.files?.[0] || null;
	}
</script>

<div class="container-fluid py-4">
	<h1 class="h2 mb-4">
		<i class="fas fa-cog me-2 text-primary"></i>Settings
	</h1>
	<p class="text-muted mb-4">Manage your Pelico installation and data</p>

	<div class="row">
		<div class="col-lg-8">
			<!-- Backup & Restore Section -->
			<div class="card mb-4">
				<div class="card-header">
					<h5 class="mb-0">
						<i class="fas fa-cloud-download-alt me-2"></i>
						Backup & Restore
					</h5>
				</div>
				<div class="card-body">
					<p class="text-muted mb-4">
						Protect your game collection data with regular backups.
					</p>
					
					<div class="row">
						<div class="col-md-6">
							<h6>Export Data</h6>
							<div class="d-grid gap-2 mb-3">
								<button
									on:click={exportFullBackup}
									disabled={exporting}
									class="btn btn-primary"
								>
									{#if exporting}
										<i class="fas fa-spinner fa-spin me-1"></i>Exporting...
									{:else}
										<i class="fas fa-download me-1"></i>Download Full Backup
									{/if}
								</button>
								<button
									on:click={exportGamesOnly}
									disabled={exportingGames}
									class="btn btn-secondary"
								>
									{#if exportingGames}
										<i class="fas fa-spinner fa-spin me-1"></i>Exporting...
									{:else}
										<i class="fas fa-gamepad me-1"></i>Download Games Only
									{/if}
								</button>
							</div>
						</div>
						
						<div class="col-md-6">
							<h6>Import Data</h6>
							<div class="mb-3">
								<input
									type="file"
									accept=".json"
									on:change={handleFileSelect}
									class="form-control mb-2"
								/>
								<small class="text-muted">Select a Pelico backup file (.json)</small>
							</div>
							<button
								on:click={importBackup}
								disabled={importing || !backupFile}
								class="btn btn-warning w-100"
							>
								{#if importing}
									<i class="fas fa-spinner fa-spin me-1"></i>Importing...
								{:else}
									<i class="fas fa-upload me-1"></i>Restore Backup
								{/if}
							</button>
							<small class="text-danger d-block mt-1">
								<i class="fas fa-exclamation-triangle me-1"></i>
								Warning: This will replace all existing data!
							</small>
						</div>
					</div>
				</div>
			</div>
		</div>

		<div class="col-lg-4">
			<!-- Nextcloud Backup -->
			<div class="card mb-4">
				<div class="card-header">
					<h5 class="mb-0">
						<i class="fas fa-cloud me-2"></i>Nextcloud Backup
					</h5>
				</div>
				<div class="card-body">
					<div class="d-flex justify-content-between align-items-center mb-3">
						<span>Status:</span>
						<span class="badge {nextcloudAvailable ? 'bg-success' : 'bg-danger'}">
							{nextcloudStatus}
						</span>
					</div>
					
					{#if nextcloudMessage}
						<p class="small text-muted mb-3">{nextcloudMessage}</p>
					{/if}
					
					<div class="d-grid gap-2">
						{#if nextcloudAvailable}
							<button
								on:click={backupToNextcloud}
								disabled={backingUpToNextcloud}
								class="btn btn-primary btn-sm"
							>
								{#if backingUpToNextcloud}
									<i class="fas fa-spinner fa-spin me-1"></i>Creating & Uploading...
								{:else}
									<i class="fas fa-cloud-upload-alt me-1"></i>Create & Upload Backup
								{/if}
							</button>
							<button
								on:click={exportAndUploadGames}
								disabled={exportingAndUploadingGames}
								class="btn btn-outline-primary btn-sm"
							>
								{#if exportingAndUploadingGames}
									<i class="fas fa-spinner fa-spin me-1"></i>Uploading Games...
								{:else}
									<i class="fas fa-gamepad me-1"></i>Upload Games Only
								{/if}
							</button>
							<div class="small text-muted mb-2">
								<i class="fas fa-info-circle me-1"></i>
								Creates JSON backup and uploads to Nextcloud automatically. To restore, download from Nextcloud and use local restore above.
							</div>
						{/if}
						<button
							on:click={testNextcloudConnection}
							disabled={testingNextcloud}
							class="btn btn-outline-secondary btn-sm"
						>
							{#if testingNextcloud}
								<i class="fas fa-spinner fa-spin me-1"></i>Testing...
							{:else}
								<i class="fas fa-wifi me-1"></i>Test Connection
							{/if}
						</button>
					</div>
				</div>
			</div>

			<!-- Cache Management -->
			<div class="card mb-4">
				<div class="card-header">
					<h5 class="mb-0">
						<i class="fas fa-database me-2"></i>Cache Management
					</h5>
				</div>
				<div class="card-body">
					<p class="small text-muted mb-3">Clear cached data and refresh the application.</p>
					<button
						on:click={refreshAllData}
						class="btn btn-warning btn-sm w-100"
					>
						<i class="fas fa-sync-alt me-1"></i>Clear Cache & Refresh
					</button>
				</div>
			</div>

			<!-- Version Info -->
			<div class="card mb-4">
				<div class="card-header">
					<h5 class="mb-0">
						<i class="fas fa-info-circle me-2"></i>Version Info
					</h5>
				</div>
				<div class="card-body">
					<dl class="row small mb-3">
						<dt class="col-sm-5">Version:</dt>
						<dd class="col-sm-7">
							<span class="badge bg-primary">v{versionInfo.version || 'Unknown'}</span>
						</dd>
						
						<dt class="col-sm-5">Build:</dt>
						<dd class="col-sm-7">{versionInfo.build_time || 'Unknown'}</dd>
						
						<dt class="col-sm-5">Commit:</dt>
						<dd class="col-sm-7">
							<span class="font-monospace small">{versionInfo.git_commit || 'Unknown'}</span>
						</dd>
						
						<dt class="col-sm-5">Go Version:</dt>
						<dd class="col-sm-7">{versionInfo.go_version || 'Unknown'}</dd>
					</dl>

					<div class="d-grid gap-2">
						<dl class="row small mb-0">
							<dt class="col-sm-5">Games:</dt>
							<dd class="col-sm-7">
								<span class="badge bg-info">{backupInfo.games_count}</span>
							</dd>
							
							<dt class="col-sm-5">Platforms:</dt>
							<dd class="col-sm-7">
								<span class="badge bg-info">{backupInfo.platforms_count}</span>
							</dd>
							
							<dt class="col-sm-5">Sessions:</dt>
							<dd class="col-sm-7">
								<span class="badge bg-info">{backupInfo.sessions_count}</span>
							</dd>
						</dl>
					</div>
				</div>
			</div>

			<!-- Changelog -->
			<div class="card">
				<div class="card-header">
					<h5 class="mb-0">
						<i class="fas fa-history me-2"></i>Recent Changes
					</h5>
				</div>
				<div class="card-body">
					<h6 class="text-primary">Version 1.2.1 Changes</h6>
					<ul class="list-unstyled mb-3">
						<li class="mb-1"><small>✅ Fixed edit game modal functionality in game detail modals</small></li>
						<li class="mb-1"><small>✅ Complete Bootstrap styling for platforms page</small></li>
						<li class="mb-1"><small>✅ Enhanced Nextcloud backup instructions</small></li>
						<li class="mb-1"><small>✅ Improved modal transitions and state management</small></li>
					</ul>
					
					<h6 class="text-secondary">Version 1.2.0 Changes</h6>
					<ul class="list-unstyled mb-3">
						<li class="mb-1"><small>✅ SvelteKit frontend migration with improved performance</small></li>
						<li class="mb-1"><small>✅ Bootstrap UI framework for consistent styling</small></li>
						<li class="mb-1"><small>✅ Enhanced game detail modals with comprehensive metadata</small></li>
						<li class="mb-1"><small>✅ Improved responsive design for mobile and desktop</small></li>
						<li class="mb-1"><small>✅ Fixed cover art display and collection view</small></li>
						<li class="mb-1"><small>✅ Enhanced ROM scanner with better directory browsing</small></li>
					</ul>
					
					<h6 class="text-secondary">Version 1.1.0 Changes</h6>
					<ul class="list-unstyled mb-3">
						<li class="mb-1"><small>✅ Added Docker container logs display in settings page</small></li>
						<li class="mb-1"><small>✅ Proper versioning system with build information</small></li>
						<li class="mb-1"><small>✅ Fixed metadata updater HTTP 500 errors for Gameboy Advance games</small></li>
						<li class="mb-1"><small>✅ Enhanced platform filter system with robust ID-based mapping</small></li>
					</ul>

					<div class="text-end">
						<small class="text-muted">
							<i class="fas fa-clock me-1"></i>
							Last updated: {new Date().toLocaleDateString()}
						</small>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>