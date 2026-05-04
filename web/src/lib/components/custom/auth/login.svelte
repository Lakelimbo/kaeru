<script lang="ts">
	import { goto } from '$app/navigation';
	import { client } from '$lib/client';
	import { setTokens } from '$lib/client/state.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Field, FieldError, FieldGroup, FieldLabel } from '$lib/components/ui/field';
	import { Input } from '$lib/components/ui/input';
	import { Spinner } from '$lib/components/ui/spinner';
	import { resolve } from '$app/paths';

	const id = $props.id();

	let email = $state('');
	let password = $state('');
	let errors: string | undefined = $state(undefined);
	let loading = $state(false);

	async function submit(event: SubmitEvent) {
		event.preventDefault();

		loading = true;
		const { data, response, error } = await client.POST('/api/v1/auth/login', {
			body: {
				email,
				password
			}
		});

		if (!response.ok && error) {
			loading = false;
			errors = error.message;
		} else {
			loading = false;
			setTokens(data!.access_token, data!.refresh_token);
			goto(resolve('/'));
		}
	}
</script>

<Card.Root class="mx-auto w-full max-w-sm">
	<Card.Header>
		<Card.Title class="text-2xl">Log in</Card.Title>
		<Card.Description>Log in to your Kaeru instance</Card.Description>
	</Card.Header>
	<Card.Content>
		<form method="POST" onsubmit={submit}>
			<FieldGroup>
				<Field>
					<FieldLabel for="email-{id}">Email</FieldLabel>
					<Input
						id="email-{id}"
						type="email"
						placeholder="majin@boo.com"
						required
						bind:value={email}
					/>
				</Field>
				<Field>
					<FieldLabel for="password-{id}">Password</FieldLabel>
					<Input id="password-{id}" type="password" required bind:value={password} />
				</Field>
			</FieldGroup>
			<Field class="mt-4">
				{#if errors}
					<FieldError>{errors}</FieldError>
				{/if}
				<Button type="submit" class="w-full" disabled={loading}>
					{#if loading}<Spinner />{/if}
					Login
				</Button>
			</Field>
		</form>
	</Card.Content>
</Card.Root>
