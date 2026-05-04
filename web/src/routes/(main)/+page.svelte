<script>
	import { authState } from '$lib/client/state.svelte.js';
	import { Clock } from '$lib/components/custom/clock';

	let { data } = $props();
</script>

<Clock class="text-5xl font-bold tabular-nums" />
<pre>{JSON.stringify(data, null, 2)}</pre>
<pre>{JSON.stringify(authState, null, 2)}</pre>
