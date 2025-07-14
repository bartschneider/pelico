import type { Game, Platform, PlaySession, WishlistItem, ShortlistItem } from './models';

export class ApiError extends Error {
	constructor(
		message: string,
		public status: number,
		public response?: Response
	) {
		super(message);
		this.name = 'ApiError';
	}
}

export class ApiClient {
	private baseUrl = '/api/v1';

	private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
		const url = `${this.baseUrl}${endpoint}`;
		
		const response = await fetch(url, {
			headers: {
				'Content-Type': 'application/json',
				...options.headers,
			},
			...options,
		});

		if (!response.ok) {
			const errorText = await response.text();
			throw new ApiError(
				errorText || `HTTP ${response.status}: ${response.statusText}`,
				response.status,
				response
			);
		}

		// Handle empty responses
		const text = await response.text();
		return text ? JSON.parse(text) : null;
	}

	async get<T>(endpoint: string): Promise<T> {
		return this.request<T>(endpoint, { method: 'GET' });
	}

	async post<T>(endpoint: string, data?: unknown): Promise<T> {
		return this.request<T>(endpoint, {
			method: 'POST',
			body: data ? JSON.stringify(data) : undefined,
		});
	}

	async put<T>(endpoint: string, data?: unknown): Promise<T> {
		return this.request<T>(endpoint, {
			method: 'PUT',
			body: data ? JSON.stringify(data) : undefined,
		});
	}

	async delete<T>(endpoint: string): Promise<T> {
		return this.request<T>(endpoint, { method: 'DELETE' });
	}

	// Game endpoints
	async getGames(): Promise<Game[]> {
		const response = await this.get<{games: Game[]}>('/games');
		return response.games || [];
	}

	async getGame(id: number): Promise<Game> {
		return this.get<Game>(`/games/${id}`);
	}

	async createGame(game: Partial<Game>): Promise<Game> {
		return this.post<Game>('/games', game);
	}

	async updateGame(id: number, game: Partial<Game>): Promise<Game> {
		return this.put<Game>(`/games/${id}`, game);
	}

	async deleteGame(id: number): Promise<void> {
		return this.delete<void>(`/games/${id}`);
	}

	async searchGames(query: string): Promise<Game[]> {
		return this.post<Game[]>('/games/search', { query });
	}

	async getRecentlyPlayedGames(): Promise<Game[]> {
		return this.get<Game[]>('/games/recently-played');
	}

	async getRecentlyAddedGames(limit = 5): Promise<Game[]> {
		const response = await this.get<{games: Game[]}>(`/games?limit=${limit}&sort=created_at&order=desc`);
		return response.games || [];
	}

	// Platform endpoints
	async getPlatforms(): Promise<Platform[]> {
		return this.get<Platform[]>('/platforms');
	}

	async getPlatform(id: number): Promise<Platform> {
		return this.get<Platform>(`/platforms/${id}`);
	}

	async createPlatform(platform: Partial<Platform>): Promise<Platform> {
		return this.post<Platform>('/platforms', platform);
	}

	async updatePlatform(id: number, platform: Partial<Platform>): Promise<Platform> {
		return this.put<Platform>(`/platforms/${id}`, platform);
	}

	async deletePlatform(id: number): Promise<void> {
		return this.delete<void>(`/platforms/${id}`);
	}

	// Session endpoints
	async getGameSessions(gameId: number): Promise<PlaySession[]> {
		return this.get<PlaySession[]>(`/games/${gameId}/sessions`);
	}

	async createSession(gameId: number, session: Partial<PlaySession>): Promise<PlaySession> {
		return this.post<PlaySession>(`/games/${gameId}/sessions`, session);
	}

	async updateSession(sessionId: number, session: Partial<PlaySession>): Promise<PlaySession> {
		return this.put<PlaySession>(`/sessions/${sessionId}`, session);
	}

	async deleteSession(sessionId: number): Promise<void> {
		return this.delete<void>(`/sessions/${sessionId}`);
	}

	async getActiveSessions(): Promise<PlaySession[]> {
		return this.get<PlaySession[]>('/sessions/active');
	}

	async getAllSessions(): Promise<PlaySession[]> {
		return this.get<PlaySession[]>('/sessions');
	}

	async endSession(sessionId: number): Promise<PlaySession> {
		return this.post<PlaySession>(`/sessions/${sessionId}/end`, {});
	}

	// Wishlist endpoints
	async getWishlist(): Promise<WishlistItem[]> {
		return this.get<WishlistItem[]>('/wishlist');
	}

	async addToWishlist(game: Partial<WishlistItem>): Promise<WishlistItem> {
		return this.post<WishlistItem>('/wishlist', game);
	}

	async removeFromWishlist(id: number): Promise<void> {
		return this.delete<void>(`/wishlist/${id}`);
	}

	// Shortlist endpoints
	async getShortlist(): Promise<ShortlistItem[]> {
		return this.get<ShortlistItem[]>('/shortlist');
	}

	async addToShortlist(game: Partial<ShortlistItem>): Promise<ShortlistItem> {
		return this.post<ShortlistItem>('/shortlist', game);
	}

	async removeFromShortlist(id: number): Promise<void> {
		return this.delete<void>(`/shortlist/${id}`);
	}

	// Stats endpoints
	async getStats(): Promise<any> {
		return this.get<any>('/stats');
	}

	// ROM Scanner endpoints
	async scanDirectory(data: {
		path: string;
		server_location: string;
		platform_id: number;
		recursive: boolean;
	}): Promise<{ message: string; games_found: number }> {
		return this.post<{ message: string; games_found: number }>('/scan/directory', data);
	}

	async updateMetadataBatch(): Promise<{ message: string; updated: number }> {
		return this.post<{ message: string; updated: number }>('/scan/metadata-batch', {});
	}

	async findDuplicates(): Promise<{ duplicates: any[] }> {
		return this.get<{ duplicates: any[] }>('/scan/duplicates');
	}

	// Directory browser endpoints
	async browseDirectory(path?: string): Promise<{
		current_path: string;
		directories: string[];
		parent?: string;
	}> {
		const params = path ? `?path=${encodeURIComponent(path)}` : '';
		return this.get<{
			current_path: string;
			directories: string[];
			parent?: string;
		}>(`/browse${params}`);
	}

	async getSuggestedPaths(): Promise<{ paths: string[] }> {
		return this.get<{ paths: string[] }>('/browse/suggestions');
	}

	// Metadata endpoints
	async fetchGameMetadata(gameId: number): Promise<Game> {
		return this.post<Game>(`/games/${gameId}/fetch-metadata`, {});
	}

	async searchMetadata(query: string): Promise<any[]> {
		return this.post<any[]>('/games/search-metadata', { query });
	}

	async createGameFromMetadata(metadata: any): Promise<Game> {
		return this.post<Game>('/games/from-metadata', metadata);
	}

	// Completion tracking endpoints
	async updateCompletionStatus(gameId: number, data: {
		completion_status: string;
		completion_percentage: number;
		completion_notes?: string;
	}): Promise<Game> {
		return this.put<Game>(`/games/${gameId}/completion`, data);
	}

	async getCompletionStats(): Promise<any> {
		return this.get<any>('/games/stats/completion');
	}

	async getGamesByCompletionStatus(status: string): Promise<Game[]> {
		return this.get<Game[]>(`/games/completion/${status}`);
	}

	// Backup and restore endpoints
	async exportDatabase(): Promise<Blob> {
		const response = await fetch(`${this.baseUrl}/backup/export`);
		return response.blob();
	}

	async exportGames(): Promise<Blob> {
		const response = await fetch(`${this.baseUrl}/backup/export/games`);
		return response.blob();
	}

	async importDatabase(file: File): Promise<{ message: string }> {
		const formData = new FormData();
		formData.append('backup', file);
		
		const response = await fetch(`${this.baseUrl}/backup/import`, {
			method: 'POST',
			body: formData,
		});
		
		if (!response.ok) {
			throw new ApiError(`HTTP ${response.status}: ${response.statusText}`, response.status, response);
		}
		
		return response.json();
	}

	async getBackupInfo(): Promise<{
		games_count: number;
		platforms_count: number;
		sessions_count: number;
		files_count: number;
	}> {
		return this.get<{
			games_count: number;
			platforms_count: number;
			sessions_count: number;
			files_count: number;
		}>('/backup/info');
	}

	async backupToNextcloud(): Promise<{ message: string }> {
		return this.post<{ message: string }>('/backup/nextcloud', {});
	}

	async testNextcloudConnection(): Promise<{ status: string; message: string }> {
		return this.get<{ status: string; message: string }>('/backup/nextcloud/test');
	}

	// Cache endpoints
	async getCacheStats(): Promise<any> {
		return this.get<any>('/cache/stats');
	}

	async clearCache(): Promise<{ message: string }> {
		return this.post<{ message: string }>('/cache/clear', {});
	}

	// Version endpoints
	async getVersion(): Promise<{
		version: string;
		build_time: string;
		git_commit: string;
		go_version: string;
	}> {
		return this.get<{
			version: string;
			build_time: string;
			git_commit: string;
			go_version: string;
		}>('/version');
	}

	// Genre endpoints
	async getGenres(): Promise<string[]> {
		return this.get<string[]>('/games/genres');
	}

	// Health check
	async healthCheck(): Promise<{ status: string }> {
		return this.get<{ status: string }>('/health');
	}
}

// Create singleton instance
export const api = new ApiClient();