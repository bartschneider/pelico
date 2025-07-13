<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { showSuccess, showError } from '$lib/stores/notifications';
	import type { PlaySession, Game } from '$lib/models';

	let activeSessions: PlaySession[] = [];
	let allGames: Game[] = [];
	let loading = false;

	// Form for adding new session
	let showAddModal = false;
	let selectedGameId = '';
	let startTime = '';
	let endTime = '';
	let sessionRating = '';
	let sessionNotes = '';

	// Edit session
	let showEditModal = false;
	let editingSession: PlaySession | null = null;

	onMount(async () => {
		await loadData();
		
		// Set default start time to now
		const now = new Date();
		startTime = now.toISOString().slice(0, 16);
	});

	async function loadData() {
		loading = true;
		try {
			const [sessions, games] = await Promise.all([
				api.getActiveSessions(),
				api.getGames()
			]);
			activeSessions = sessions;
			allGames = games;
		} catch (error: any) {
			showError('Error', 'Failed to load data');
		} finally {
			loading = false;
		}
	}

	function openAddModal() {
		selectedGameId = '';
		const now = new Date();
		startTime = now.toISOString().slice(0, 16);
		endTime = '';
		sessionRating = '';
		sessionNotes = '';
		showAddModal = true;
	}

	function openEditModal(session: PlaySession) {
		editingSession = session;
		selectedGameId = session.game_id?.toString() || '';
		
		if (session.start_time) {
			startTime = new Date(session.start_time).toISOString().slice(0, 16);
		}
		if (session.end_time) {
			endTime = new Date(session.end_time).toISOString().slice(0, 16);
		} else {
			endTime = '';
		}
		
		sessionRating = session.rating?.toString() || '';
		sessionNotes = session.notes || '';
		showEditModal = true;
	}

	function closeModals() {
		showAddModal = false;
		showEditModal = false;
		editingSession = null;
		selectedGameId = '';
		startTime = '';
		endTime = '';
		sessionRating = '';
		sessionNotes = '';
	}

	async function handleAddSession() {
		if (!selectedGameId || !startTime) {
			showError('Validation Error', 'Please select a game and start time');
			return;
		}

		const sessionData: Partial<PlaySession> = {
			start_time: new Date(startTime).toISOString(),
			end_time: endTime ? new Date(endTime).toISOString() : undefined,
			rating: sessionRating ? parseInt(sessionRating) : undefined,
			notes: sessionNotes || undefined
		};

		try {
			const newSession = await api.createSession(parseInt(selectedGameId), sessionData);
			await loadData(); // Refresh the list
			showSuccess('Success', 'Play session logged successfully');
			closeModals();
		} catch (error: any) {
			showError('Error', error.message);
		}
	}

	async function handleUpdateSession() {
		if (!editingSession || !selectedGameId || !startTime) {
			showError('Validation Error', 'Please fill in required fields');
			return;
		}

		const sessionData: Partial<PlaySession> = {
			start_time: new Date(startTime).toISOString(),
			end_time: endTime ? new Date(endTime).toISOString() : undefined,
			rating: sessionRating ? parseInt(sessionRating) : undefined,
			notes: sessionNotes || undefined
		};

		try {
			await api.updateSession(editingSession.id, sessionData);
			await loadData(); // Refresh the list
			showSuccess('Success', 'Session updated successfully');
			closeModals();
		} catch (error: any) {
			showError('Error', error.message);
		}
	}

	async function handleEndSession(session: PlaySession) {
		try {
			await api.endSession(session.id);
			await loadData(); // Refresh the list
			showSuccess('Success', 'Session ended');
		} catch (error: any) {
			showError('Error', error.message);
		}
	}

	async function handleDeleteSession(session: PlaySession) {
		if (!confirm('Are you sure you want to delete this play session?')) {
			return;
		}

		try {
			await api.deleteSession(session.id);
			await loadData(); // Refresh the list
			showSuccess('Success', 'Session deleted');
		} catch (error: any) {
			showError('Error', error.message);
		}
	}

	function formatDuration(start: string, end?: string): string {
		const startDate = new Date(start);
		const endDate = end ? new Date(end) : new Date();
		const diffMs = endDate.getTime() - startDate.getTime();
		
		const hours = Math.floor(diffMs / (1000 * 60 * 60));
		const minutes = Math.floor((diffMs % (1000 * 60 * 60)) / (1000 * 60));
		
		if (hours > 0) {
			return `${hours}h ${minutes}m`;
		}
		return `${minutes}m`;
	}

	function getGameTitle(gameId: number): string {
		const game = allGames.find(g => g.id === gameId);
		return game ? game.title : 'Unknown Game';
	}

	function isActiveSession(session: PlaySession): boolean {
		return !session.end_time;
	}
</script>

<div class="container-fluid py-4">
	<div class="d-flex justify-content-between align-items-center mb-4">
		<div>
			<h1 class="h2 mb-1">
				<i class="fas fa-clock me-2 text-primary"></i>Play Sessions
			</h1>
			<p class="text-muted mb-0">Track and manage your gaming sessions.</p>
		</div>
		<button
			on:click={openAddModal}
			class="btn btn-success"
		>
			<i class="fas fa-play me-1"></i>Log Session
		</button>
	</div>

	{#if loading}
		<div class="text-center py-5">
			<div class="spinner-border text-primary mb-3" role="status">
				<span class="visually-hidden">Loading...</span>
			</div>
			<p class="text-muted">Loading sessions...</p>
		</div>
	{:else}
		<!-- Active Sessions -->
		{#if activeSessions.filter(session => isActiveSession(session)).length > 0}
			<div class="mb-4">
				<h2 class="h4 mb-3">
					<i class="fas fa-play-circle text-success me-2"></i>Currently Playing
				</h2>
				<div class="row">
					{#each activeSessions.filter(session => isActiveSession(session)) as session}
						<div class="col-md-6 col-lg-4 mb-3">
							<div class="card border-success">
								<div class="card-body">
									<div class="d-flex justify-content-between align-items-start mb-2">
										<h5 class="card-title mb-0">{getGameTitle(session.game_id)}</h5>
										<span class="badge bg-success">ACTIVE</span>
									</div>
									<p class="card-text small text-muted mb-2">
										<i class="fas fa-clock me-1"></i>
										Started: {new Date(session.start_time).toLocaleString()}
									</p>
									<p class="card-text small text-muted mb-3">
										<i class="fas fa-stopwatch me-1"></i>
										Duration: {formatDuration(session.start_time)}
									</p>
									<div class="d-flex gap-2">
										<button
											on:click={() => handleEndSession(session)}
											class="btn btn-sm btn-danger"
										>
											<i class="fas fa-stop me-1"></i>End Session
										</button>
										<button
											on:click={() => openEditModal(session)}
											class="btn btn-sm btn-primary"
										>
											<i class="fas fa-edit me-1"></i>Edit
										</button>
									</div>
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		<!-- All Sessions -->
		<div>
			<h2 class="h4 mb-3">
				<i class="fas fa-history me-2"></i>Session History
			</h2>
			
			{#if activeSessions.length === 0}
				<div class="text-center py-5">
					<i class="fas fa-clock fa-3x text-muted mb-3"></i>
					<h3 class="h5 mb-2">No Sessions Found</h3>
					<p class="text-muted mb-4">Start tracking your gaming sessions to see your play history.</p>
					<button
						on:click={openAddModal}
						class="btn btn-success"
					>
						<i class="fas fa-play me-1"></i>Log First Session
					</button>
				</div>
			{:else}
				<div class="card">
					<div class="table-responsive">
						<table class="table table-hover mb-0">
							<thead class="table-light">
								<tr>
									<th>Game</th>
									<th>Start Time</th>
									<th>Duration</th>
									<th>Rating</th>
									<th>Status</th>
									<th>Actions</th>
								</tr>
							</thead>
							<tbody>
								{#each activeSessions as session}
									<tr>
										<td>
											<div class="fw-medium">{getGameTitle(session.game_id)}</div>
											{#if session.notes}
												<div class="small text-muted text-truncate" style="max-width: 200px;">{session.notes}</div>
											{/if}
										</td>
										<td class="small">
											{new Date(session.start_time).toLocaleString()}
										</td>
										<td class="small">
											{formatDuration(session.start_time, session.end_time)}
										</td>
										<td>
											{#if session.rating}
												<div class="d-flex align-items-center">
													<span class="small me-2">{session.rating}/10</span>
													<div>
														{#each Array(5) as _, i}
															<i class="fas fa-star {i < (session.rating || 0) / 2 ? 'text-warning' : 'text-light'}" style="font-size: 0.7rem;"></i>
														{/each}
													</div>
												</div>
											{:else}
												<span class="text-muted">-</span>
											{/if}
										</td>
										<td>
											{#if isActiveSession(session)}
												<span class="badge bg-success">
													<i class="fas fa-play me-1"></i>Active
												</span>
											{:else}
												<span class="badge bg-secondary">
													<i class="fas fa-check me-1"></i>Completed
												</span>
											{/if}
										</td>
										<td>
											<div class="d-flex gap-1">
												{#if isActiveSession(session)}
													<button
														on:click={() => handleEndSession(session)}
														class="btn btn-sm btn-outline-danger"
														title="End Session"
													>
														<i class="fas fa-stop"></i>
													</button>
												{/if}
												<button
													on:click={() => openEditModal(session)}
													class="btn btn-sm btn-outline-primary"
													title="Edit Session"
												>
													<i class="fas fa-edit"></i>
												</button>
												<button
													on:click={() => handleDeleteSession(session)}
													class="btn btn-sm btn-outline-danger"
													title="Delete Session"
												>
													<i class="fas fa-trash"></i>
												</button>
											</div>
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div>

<!-- Add Session Modal -->
{#if showAddModal}
	<div class="modal fade show d-block" tabindex="-1" style="background-color: rgba(0,0,0,0.5);">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title">
						<i class="fas fa-play me-2"></i>Log Play Session
					</h5>
					<button type="button" class="btn-close" on:click={closeModals}></button>
				</div>

				<form on:submit|preventDefault={handleAddSession}>
					<div class="modal-body">
						<div class="mb-3">
							<label for="gameSelect" class="form-label">
								Game <span class="text-danger">*</span>
							</label>
							<select
								id="gameSelect"
								bind:value={selectedGameId}
								required
								class="form-select"
							>
								<option value="">Select a game</option>
								{#each allGames as game}
									<option value={game.id}>{game.title}</option>
								{/each}
							</select>
						</div>

						<div class="mb-3">
							<label for="startTime" class="form-label">
								Start Time <span class="text-danger">*</span>
							</label>
							<input
								type="datetime-local"
								id="startTime"
								bind:value={startTime}
								required
								class="form-control"
							/>
						</div>

						<div class="mb-3">
							<label for="endTime" class="form-label">
								End Time (optional)
							</label>
							<input
								type="datetime-local"
								id="endTime"
								bind:value={endTime}
								class="form-control"
							/>
							<div class="form-text">Leave empty for ongoing session</div>
						</div>

						<div class="mb-3">
							<label for="sessionRating" class="form-label">
								Session Rating (1-10)
							</label>
							<select
								id="sessionRating"
								bind:value={sessionRating}
								class="form-select"
							>
								<option value="">No rating</option>
								<option value="10">10 - Excellent</option>
								<option value="9">9 - Great</option>
								<option value="8">8 - Very Good</option>
								<option value="7">7 - Good</option>
								<option value="6">6 - Fair</option>
								<option value="5">5 - Average</option>
								<option value="4">4 - Below Average</option>
								<option value="3">3 - Poor</option>
								<option value="2">2 - Bad</option>
								<option value="1">1 - Terrible</option>
							</select>
						</div>

						<div class="mb-3">
							<label for="sessionNotes" class="form-label">
								Notes
							</label>
							<textarea
								id="sessionNotes"
								bind:value={sessionNotes}
								rows="3"
								class="form-control"
								placeholder="How was this gaming session? Any achievements, progress, or thoughts?"
							></textarea>
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
							class="btn btn-success"
						>
							<i class="fas fa-save me-1"></i>Log Session
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}

<!-- Edit Session Modal -->
{#if showEditModal && editingSession}
	<div class="modal fade show d-block" tabindex="-1" style="background-color: rgba(0,0,0,0.5);">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title">
						<i class="fas fa-edit me-2"></i>Edit Play Session
					</h5>
					<button type="button" class="btn-close" on:click={closeModals}></button>
				</div>

				<form on:submit|preventDefault={handleUpdateSession}>
					<div class="modal-body">
						<div class="mb-3">
							<label for="editGameSelect" class="form-label">
								Game <span class="text-danger">*</span>
							</label>
							<select
								id="editGameSelect"
								bind:value={selectedGameId}
								required
								class="form-select"
							>
								<option value="">Select a game</option>
								{#each allGames as game}
									<option value={game.id}>{game.title}</option>
								{/each}
							</select>
						</div>

						<div class="mb-3">
							<label for="editStartTime" class="form-label">
								Start Time <span class="text-danger">*</span>
							</label>
							<input
								type="datetime-local"
								id="editStartTime"
								bind:value={startTime}
								required
								class="form-control"
							/>
						</div>

						<div class="mb-3">
							<label for="editEndTime" class="form-label">
								End Time (optional)
							</label>
							<input
								type="datetime-local"
								id="editEndTime"
								bind:value={endTime}
								class="form-control"
							/>
							<div class="form-text">Leave empty for ongoing session</div>
						</div>

						<div class="mb-3">
							<label for="editSessionRating" class="form-label">
								Session Rating (1-10)
							</label>
							<select
								id="editSessionRating"
								bind:value={sessionRating}
								class="form-select"
							>
								<option value="">No rating</option>
								<option value="10">10 - Excellent</option>
								<option value="9">9 - Great</option>
								<option value="8">8 - Very Good</option>
								<option value="7">7 - Good</option>
								<option value="6">6 - Fair</option>
								<option value="5">5 - Average</option>
								<option value="4">4 - Below Average</option>
								<option value="3">3 - Poor</option>
								<option value="2">2 - Bad</option>
								<option value="1">1 - Terrible</option>
							</select>
						</div>

						<div class="mb-3">
							<label for="editSessionNotes" class="form-label">
								Notes
							</label>
							<textarea
								id="editSessionNotes"
								bind:value={sessionNotes}
								rows="3"
								class="form-control"
								placeholder="How was this gaming session? Any achievements, progress, or thoughts?"
							></textarea>
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
							<i class="fas fa-save me-1"></i>Update Session
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}