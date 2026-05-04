<script lang="ts">
	let { class: className = 'text-3xl font-bold tabular-nums' }: { class?: string } = $props();

	let date = $state(new Date());
	let time = $derived(
		date.toLocaleTimeString('en-UK', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
	);

	$effect(() => {
		setInterval(() => {
			date = new Date();
		}, 1000);
	});
</script>

<div class={className}>{time}</div>
