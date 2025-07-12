export interface Game {
  id: number;
  title: string;
  platform_id: number;
  platform?: Platform;
  year: number;
  genre: string;
  description: string;
  cover_art_url: string;
  rating: number;
  purchase_date: string;
  collection_formats: string[];
  file_locations: any[]; // Define a proper type for this later
  play_sessions: any[]; // Define a proper type for this later
  completion_status: string;
  completion_percentage: number;
}

export interface Platform {
  id: number;
  name: string;
  manufacturer: string;
  release_year: number;
}
