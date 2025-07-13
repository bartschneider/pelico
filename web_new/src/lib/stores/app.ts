import { writable } from 'svelte/store';

export interface AppState {
	currentView: string;
	loading: boolean;
	sidebarOpen: boolean;
	viewMode: 'grid' | 'list';
}

const initialState: AppState = {
	currentView: 'collection',
	loading: false,
	sidebarOpen: false,
	viewMode: 'grid'
};

export const appState = writable<AppState>(initialState);

// Helper functions for common operations
export const setCurrentView = (view: string) => {
	appState.update(state => ({ ...state, currentView: view }));
};

export const setLoading = (loading: boolean) => {
	appState.update(state => ({ ...state, loading }));
};

export const toggleSidebar = () => {
	appState.update(state => ({ ...state, sidebarOpen: !state.sidebarOpen }));
};

export const setViewMode = (mode: 'grid' | 'list') => {
	appState.update(state => ({ ...state, viewMode: mode }));
	// Persist to localStorage
	if (typeof localStorage !== 'undefined') {
		localStorage.setItem('pelico_viewMode', mode);
	}
};

// Initialize view mode from localStorage on client
if (typeof localStorage !== 'undefined') {
	const savedViewMode = localStorage.getItem('pelico_viewMode') as 'grid' | 'list' | null;
	if (savedViewMode) {
		setViewMode(savedViewMode);
	}
}