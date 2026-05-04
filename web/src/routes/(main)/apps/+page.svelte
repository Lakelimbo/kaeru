<script lang="ts">
	import { AppCreate, AppEmpty, AppItem } from '$lib/components/custom/app';
	import { Button } from '$lib/components/ui/button';

	let { data } = $props();
</script>

<div class="flex flex-col gap-4 sm:flex-row md:items-center">
	<h1 class="flex-1 text-3xl font-bold">Installed apps</h1>
	<div class="flex items-center gap-2">
		<AppCreate class="flex-1" />
		<Button class="flex-1">App store</Button>
	</div>
</div>
<article class="mt-4">
	{#await data.apps}
		LOADING
	{:then apps}
		<div class="grid grid-cols-2 gap-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6">
			{#each apps.data as app (app.project_name)}
				<AppItem id={app.id} name={app.name} projectName={app.project_name} />
			{:else}
				<AppEmpty />
			{/each}
		</div>
	{/await}
</article>
