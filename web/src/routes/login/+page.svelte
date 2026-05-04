<script>
	import { goto } from '$app/navigation';
	import { authState } from '$lib/client/state.svelte';
	import { AuthLogin } from '$lib/components/custom/auth';
	import { BrandIcon } from '$lib/components/custom/brand';
	import { resolve } from '$app/paths';

	$effect.pre(() => {
		if (authState.isAuthenticated) {
			goto(resolve('/'));
		}
	});
</script>

<main class="flex min-h-screen flex-col justify-center">
	<BrandIcon class="mx-auto mb-2 fill-brand" />
	<AuthLogin />
</main>
