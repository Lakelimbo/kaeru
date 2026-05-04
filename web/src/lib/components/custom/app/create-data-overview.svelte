<script lang="ts">
	import * as Accordion from '$lib/components/ui/accordion';
	import * as Table from '$lib/components/ui/table';
	import {
		Database01FreeIcons,
		MultiplicationSignFreeIcons,
		Package01FreeIcons,
		Tick02FreeIcons
	} from '@hugeicons/core-free-icons';
	import { HugeiconsIcon } from '@hugeicons/svelte';
	import type { KaeruCompose } from '.';

	type Props = {
		services?: KaeruCompose['services'];
		volumes?: KaeruCompose['volumes'];
	};

	let { services, volumes }: Props = $props();
</script>

<div class="mx-6 mt-4">
	<Accordion.Root type="single">
		{#if services}
			<Accordion.Item value="services">
				<Accordion.Trigger><HugeiconsIcon icon={Package01FreeIcons} /> Services</Accordion.Trigger>
				<Accordion.Content>
					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Name</Table.Head>
								<Table.Head>Image</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each Object.entries(services) as service, i (i)}
								<Table.Row>
									<Table.Cell>{service[0]}</Table.Cell>
									<Table.Cell class="font-mono">{service[1].image}</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</Accordion.Content>
			</Accordion.Item>
		{/if}
		{#if volumes}
			<Accordion.Item value="volumes">
				<Accordion.Trigger><HugeiconsIcon icon={Database01FreeIcons} /> Volumes</Accordion.Trigger>
				<Accordion.Content>
					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Name</Table.Head>
								<Table.Head>Explorable</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each Object.entries(volumes) as volume, i (i)}
								<Table.Row>
									<Table.Cell>{volume[0]}</Table.Cell>
									<Table.Cell class="font-mono">
										{#if volume[1] && volume[1].exposed}
											<HugeiconsIcon icon={Tick02FreeIcons} class="text-brand" />
										{:else}
											<HugeiconsIcon icon={MultiplicationSignFreeIcons} class="text-destructive" />
										{/if}
									</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</Accordion.Content>
			</Accordion.Item>
		{/if}
	</Accordion.Root>
</div>
<pre>{JSON.stringify(volumes, null, 2)}</pre>
