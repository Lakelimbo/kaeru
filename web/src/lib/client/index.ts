import createClient from 'openapi-fetch';
import { authState, logout, setTokens } from './state.svelte';
import type { paths } from './v1.d';

export const BASE_URL = 'http://localhost:4040';

export const client = createClient<paths>({
	baseUrl: BASE_URL,
	credentials: 'include'
});

client.use({
	async onRequest({ request }) {
		if (authState.accessToken) {
			request.headers.set('Authorization', `Bearer ${authState.accessToken}`);
		}

		return request;
	},
	async onResponse({ response, request }) {
		if (response.status === 401 && !authState.refreshToken) {
			if (!authState.refreshPromise) {
				authState.refreshPromise = (async () => {
					try {
						const res = await fetch(`${BASE_URL}/api/v1/auth/refresh`, {
							method: 'POST',
							body: JSON.stringify({
								refresh_token: authState.refreshToken
							})
						});

						if (!res.ok) {
							throw new Error('Refresh failed');
						}

						const data = await res.json();
						setTokens(data.access_token, data.refresh_token);
						return data.access_token;
					} catch {
						logout();
						return null;
					} finally {
						authState.refreshPromise = null;
					}
				})();
			}

			const newToken = await authState.refreshPromise;
			if (newToken) {
				const retry = new Request(request.url, request);
				retry.headers.set('Authorization', `Bearer ${newToken}`);

				return fetch(retry);
			}
		}

		return response;
	}
});
