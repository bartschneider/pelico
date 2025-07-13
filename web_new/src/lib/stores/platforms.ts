import { writable, derived } from 'svelte/store';
import type { Platform } from '../models';

export interface PlatformsState {
	platforms: Platform[];
	loading: boolean;
	error: string | null;
	selectedPlatform: Platform | null;
}

const initialState: PlatformsState = {
	platforms: [],
	loading: false,
	error: null,
	selectedPlatform: null
};

export const platformsState = writable<PlatformsState>(initialState);

// Derived stores for easy access
export const platforms = derived(platformsState, $state => $state.platforms);
export const platformsLoading = derived(platformsState, $state => $state.loading);
export const platformsError = derived(platformsState, $state => $state.error);

// Helper functions
export const setPlatforms = (platforms: Platform[]) => {
	platformsState.update(state => ({ ...state, platforms, error: null }));
};

export const addPlatform = (platform: Platform) => {
	platformsState.update(state => ({
		...state,
		platforms: [...state.platforms, platform]
	}));
};

export const updatePlatform = (updatedPlatform: Platform) => {
	platformsState.update(state => ({
		...state,
		platforms: state.platforms.map(platform =>
			platform.id === updatedPlatform.id ? updatedPlatform : platform
		)
	}));
};

export const removePlatform = (platformId: number) => {
	platformsState.update(state => ({
		...state,
		platforms: state.platforms.filter(platform => platform.id !== platformId)
	}));
};

export const setPlatformsLoading = (loading: boolean) => {
	platformsState.update(state => ({ ...state, loading }));
};

export const setPlatformsError = (error: string | null) => {
	platformsState.update(state => ({ ...state, error }));
};

export const setSelectedPlatform = (platform: Platform | null) => {
	platformsState.update(state => ({ ...state, selectedPlatform: platform }));
};