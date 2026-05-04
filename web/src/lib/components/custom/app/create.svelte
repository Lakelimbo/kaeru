<script lang="ts">
	import { client } from '$lib/client';
	import { buttonVariants } from '$lib/components/ui/button';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Field, FieldGroup, FieldLabel } from '$lib/components/ui/field';
	import { Input } from '$lib/components/ui/input';
	import * as Sheet from '$lib/components/ui/sheet';
	import { cn } from '$lib/utils/styles';
	import { parse } from 'yaml';
	import type { KaeruCompose } from '.';
	import { Editor } from '../editor';
	import CreateDataOverview from './create-data-overview.svelte';

	type Props = {
		class?: string;
	};

	let { class: className }: Props = $props();
	const id = $props.id();

	let open = $state(false);
	let loading = $state(false);

	let compose: string = $state('');
	let parsed: KaeruCompose = $derived(parse(compose));

	let projectName: string = $state('');

	async function submit(event: SubmitEvent) {
		event.preventDefault();
		loading = true;

		const { response, error } = await client.POST('/api/v1/apps', {
			body: {
				compose_yaml: compose,
				name: `${projectName}_x`,
				origin: '',
				project_name: projectName
			}
		});

		if (!response.ok && error) {
			loading = false;
			throw new Error(error.message);
		}

		open = false;
		loading = false;

		compose = '';
		projectName = '';
	}
</script>

<Sheet.Root bind:open>
	<Sheet.Trigger class={cn(buttonVariants({ variant: 'outline' }), className)}
		>Create new</Sheet.Trigger
	>
	<Sheet.Content side="right" class="w-full max-w-none! overflow-y-auto md:w-180">
		<Sheet.Header>
			<Sheet.Title>Creating app manually</Sheet.Title>
			<Sheet.Description>
				Create a new stack of apps manually by pasting a Compose YAML declaration file.
			</Sheet.Description>
		</Sheet.Header>

		<form method="POST" class="grid h-full gap-2 md:grid-cols-2" onsubmit={submit}>
			<FieldGroup class="h-full px-6">
				<Field class="h-full">
					<FieldLabel for="compose-{id}">Compose YAML</FieldLabel>
					<Editor
						bind:value={compose}
						options={{
							language: 'yaml',
							minimap: { enabled: false },
							lineNumbers: 'off',
							automaticLayout: false
						}}
					/>
				</Field>
			</FieldGroup>

			<div>
				<h3 class="mt-2 mb-4 px-6 text-lg font-medium">Project data</h3>
				<FieldGroup class="gap-3 px-6">
					<FieldLabel for="project-name-{id}">Project name</FieldLabel>
					<Input
						id="project-name-{id}"
						bind:value={projectName}
						placeholder={parsed && parsed['x-kaeru'].project_name
							? parsed['x-kaeru'].project_name
							: undefined}
					/>
				</FieldGroup>

				{#if parsed}
					{#if (parsed.services && Object.isExtensible(parsed.services)) || (parsed.volumes && Object.isExtensible(parsed.volumes))}
						<CreateDataOverview services={parsed.services} volumes={parsed.volumes} />
					{/if}
				{/if}
			</div>

			<Sheet.Footer class="col-span-2">
				<Button disabled={loading} type="submit">Create</Button>
			</Sheet.Footer>
		</form>
	</Sheet.Content>
</Sheet.Root>
