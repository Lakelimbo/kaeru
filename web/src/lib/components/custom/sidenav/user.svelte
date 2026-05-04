<script lang="ts">
	import { page } from '$app/state';
	import type { components } from '$lib/client/v1';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { cn } from '$lib/utils/styles';
	import {
		AccountSetting01FreeIcons,
		GibbousMoonFreeIcons,
		Logout05FreeIcons,
		Settings02FreeIcons
	} from '@hugeicons/core-free-icons';
	import { HugeiconsIcon, type IconSvgElement } from '@hugeicons/svelte';
	import { toggleMode } from 'mode-watcher';
	import type { MouseEventHandler } from 'svelte/elements';
	import { UserAvatar } from '../user';
	import { logout } from '$lib/client/state.svelte';

	let { class: className }: { class?: string } = $props();

	const user: components['schemas']['auth.User'] = $derived(page.data.user);

	type Content = {
		name: string;
		href?: string;
		onclick?: MouseEventHandler<HTMLButtonElement> & MouseEventHandler<HTMLAnchorElement>;
		icon: IconSvgElement;
		class?: string;
	};

	const contents: Content[][] = [
		[
			{
				name: 'Profile',
				icon: AccountSetting01FreeIcons
			},
			{
				name: 'Settings',
				icon: Settings02FreeIcons
			},
			{
				name: 'Toggle dark/light theme',
				icon: GibbousMoonFreeIcons,
				onclick: toggleMode
			}
		],
		[
			{
				name: 'Log out',
				icon: Logout05FreeIcons,
				class: 'text-destructive hover:bg-destructive/10!',
				onclick: logout
			}
		]
	];
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<Button variant="ghost" {...props} class={className}>
				<UserAvatar />
			</Button>
		{/snippet}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content class="w-56">
		<DropdownMenu.Label class="grid gap-2 text-sm">
			<span>{user.username}</span>
			<span>{user.email}</span>
		</DropdownMenu.Label>
		{#each contents as content, i (i)}
			<DropdownMenu.Separator />
			<DropdownMenu.Group>
				{#each content as item (item.name)}
					<DropdownMenu.Item class={cn('w-full justify-start', item.class)}>
						{#snippet child({ props })}
							<Button variant="ghost" {...props} onclick={item.onclick}
								><HugeiconsIcon icon={item.icon} /> {item.name}</Button
							>
						{/snippet}
					</DropdownMenu.Item>
				{/each}
			</DropdownMenu.Group>
		{/each}
	</DropdownMenu.Content>
</DropdownMenu.Root>
