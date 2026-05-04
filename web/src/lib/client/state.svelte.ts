import { browser } from '$app/environment';
import { client } from '.';

const ACCESS_KEY = '.kaeru_access';
const REFRESH_KEY = '.kaeru_refresh';

export const authState = $state({
	isLoading: true,
	isAuthenticated: browser ? !!localStorage.getItem(ACCESS_KEY) : false,
	refreshPromise: null as Promise<string | null> | null,
	accessToken: browser ? localStorage.getItem(ACCESS_KEY) : null,
	refreshToken: browser ? localStorage.getItem(REFRESH_KEY) : null
});

export function setTokens(access: string | null, refresh: string | null) {
	authState.accessToken = access;
	authState.refreshToken = refresh;
	authState.isAuthenticated = !!access;

	if (browser) {
		if (access) {
			localStorage.setItem(ACCESS_KEY, access);
		} else {
			localStorage.removeItem(ACCESS_KEY);
		}

		if (refresh) {
			localStorage.setItem(REFRESH_KEY, refresh);
		} else {
			localStorage.removeItem(REFRESH_KEY);
		}
	}
}

export async function logout() {
	const { error, response } = await client.POST('/api/v1/auth/logout', {
		body: {
			refresh_token: authState.refreshToken!
		}
	});

	if (!response.ok) {
		throw new Error(error?.message);
	}

	setTokens(null, null);
	authState.isAuthenticated = false;
}
