export { default as AppCreate } from './create.svelte';
export { default as AppEmpty } from './empty.svelte';
export { default as AppItem } from './item.svelte';

export type KaeruCompose = {
	'x-kaeru': {
		project_name: string;
		exposed_volumes?: string[];
	};
	services: {
		[key: string]: {
			image: string;
		};
	};
	volumes: {
		[key: string]: {
			exposed?: boolean;
		};
	};
};
