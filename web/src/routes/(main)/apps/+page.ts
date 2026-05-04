import { client } from '$lib/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	const apps = client.GET('/api/v1/apps');

	return {
		apps: apps
	};
};
