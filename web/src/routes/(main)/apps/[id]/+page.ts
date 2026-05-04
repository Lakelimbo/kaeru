import { client } from '$lib/client';
import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch }) => {
	const {
		data,
		error: err,
		response
	} = await client.GET('/api/v1/apps/{id}', {
		params: {
			path: {
				id: params.id
			}
		},
		fetch
	});

	if (!response.ok) {
		error(response.status, err);
	}

	return {
		app: data
	};
};
