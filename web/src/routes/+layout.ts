import { browser } from '$app/environment';
import { client } from '$lib/client';
import { authState, logout } from '$lib/client/state.svelte';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async () => {
	if (!browser) return;

	try {
		if (!authState.accessToken) {
			return { user: null };
		}

		const { data, error, response } = await client.GET('/api/v1/auth/me');

		if (error || response.status === 401) {
			logout();
			return { user: null };
		}

		authState.isAuthenticated = true;

		return {
			user: data
		};
	} finally {
		authState.isLoading = false;
	}
};
