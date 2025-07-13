export interface Game {
  id: number;
  title: string;
  platform_id?: number;
  platform?: Platform;
  release_date?: string;
  genre?: string;
  completion_status: 'not_started' | 'in_progress' | 'completed' | 'abandoned';
  rating?: number;
  notes?: string;
  playtime_hours?: number;
  cover_image_url?: string;
  description?: string;
  developer?: string;
  publisher?: string;
  created_at?: string;
  updated_at?: string;
  file_locations?: FileLocation[];
  play_sessions?: PlaySession[];
}

export interface Platform {
  id: number;
  name: string;
  manufacturer?: string;
  release_year?: number;
  short_name?: string;
  icon_url?: string;
  created_at?: string;
  updated_at?: string;
}

export interface FileLocation {
  id: number;
  game_id: number;
  file_path: string;
  file_size?: number;
  file_hash?: string;
  server_name?: string;
  created_at?: string;
  updated_at?: string;
}

export interface PlaySession {
  id: number;
  game_id: number;
  start_time: string;
  end_time?: string;
  duration_minutes?: number;
  notes?: string;
  created_at?: string;
  updated_at?: string;
}

export interface WishlistItem {
  id: number;
  title: string;
  platform_name?: string;
  priority: 'low' | 'medium' | 'high';
  notes?: string;
  estimated_price?: number;
  release_date?: string;
  created_at?: string;
  updated_at?: string;
}

export interface ShortlistItem {
  id: number;
  game_id?: number;
  game?: Game;
  title: string;
  platform_name?: string;
  priority: 'low' | 'medium' | 'high';
  reason?: string;
  created_at?: string;
  updated_at?: string;
}

export interface GameFormData {
  title: string;
  platform_id?: number;
  release_date?: string;
  genre?: string;
  completion_status: Game['completion_status'];
  rating?: number;
  notes?: string;
  description?: string;
  developer?: string;
  publisher?: string;
}

export interface SearchResult {
  id: number;
  name: string;
  platform?: string;
  release_date?: string;
  cover_url?: string;
  description?: string;
  external_id?: string;
  source: 'igdb' | 'thegamesdb';
}

export interface ApiResponse<T> {
  data: T;
  message?: string;
  success: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}
