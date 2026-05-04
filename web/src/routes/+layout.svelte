<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import { authState } from '$lib/client/state.svelte';
	import { Spinner } from '$lib/components/ui/spinner';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { ModeWatcher } from 'mode-watcher';
	import './layout.css';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<ModeWatcher />
{#if authState.isLoading}
	<div class="flex h-screen items-center">
		<Spinner class="m-auto size-8" />
	</div>
{:else}
	<Tooltip.Provider>
		{@render children()}
	</Tooltip.Provider>
{/if}
