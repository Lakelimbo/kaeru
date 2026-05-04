<script lang="ts">
	import { resolve } from '$app/paths';
	import * as Item from '$lib/components/ui/item';

	type Props = {
		id: string;
		name: string;
		projectName: string;
	};

	let { id, name, projectName }: Props = $props();
</script>

<Item.Root>
	{#snippet child({ props })}
		<a href={resolve('/(main)/apps/[id]', { id: id })} {...props}>
			<Item.Media>
				<div class="size-16 rounded-lg bg-secondary"></div>
			</Item.Media>
			<Item.Content>
				<Item.Title>{name}</Item.Title>
				<Item.Description>{projectName}</Item.Description>
			</Item.Content>
		</a>
	{/snippet}
</Item.Root>
