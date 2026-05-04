<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Empty from '$lib/components/ui/empty';
	import { GridViewFreeIcons } from '@hugeicons/core-free-icons';
	import { HugeiconsIcon } from '@hugeicons/svelte';
	import { AppCreate } from '.';
</script>

<Empty.Root class="col-span-full">
	<Empty.Header>
		<Empty.Media variant="icon">
			<HugeiconsIcon icon={GridViewFreeIcons} />
		</Empty.Media>
		<Empty.Title>No apps installed</Empty.Title>
		<Empty.Description>Install an app from the store or manually create one!</Empty.Description>
		<Empty.Content class="flex flex-row items-center justify-center gap-2">
			<Button>App store</Button>
			<AppCreate />
		</Empty.Content>
	</Empty.Header>
</Empty.Root>
