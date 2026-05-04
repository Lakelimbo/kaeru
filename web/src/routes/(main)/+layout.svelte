<script lang="ts">
	import { goto } from '$app/navigation';
	import { authState } from '$lib/client/state.svelte';
	import { MobileTopNav, Sidenav } from '$lib/components/custom/sidenav';
	import { resolve } from '$app/paths';

	let { children } = $props();

	$effect.pre(() => {
		if (!authState.isLoading && !authState.isAuthenticated) {
			goto(resolve('/login'));
		}
	});
</script>

<div class="flex min-h-[200vh] flex-col gap-4 md:flex-row">
	<MobileTopNav />
	<main class="flex-1 p-4 md:order-last">
		{@render children()}
	</main>
	<Sidenav />
</div>
