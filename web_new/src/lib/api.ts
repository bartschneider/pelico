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
		return this.get<Game[]>('/games');
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

	// Platform endpoints
	async getPlatforms(): Promise<Platform[]> {
		return this.get<Platform[]>('/platforms');
	}

	async createPlatform(platform: Partial<Platform>): Promise<Platform> {
		return this.post<Platform>('/platforms', platform);
	}

	// Session endpoints
	async getGameSessions(gameId: number): Promise<PlaySession[]> {
		return this.get<PlaySession[]>(`/games/${gameId}/sessions`);
	}

	async createSession(gameId: number, session: Partial<PlaySession>): Promise<PlaySession> {
		return this.post<PlaySession>(`/games/${gameId}/sessions`, session);
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

	// Health check
	async healthCheck(): Promise<{ status: string }> {
		return this.get<{ status: string }>('/health');
	}
}

// Create singleton instance
export const api = new ApiClient();