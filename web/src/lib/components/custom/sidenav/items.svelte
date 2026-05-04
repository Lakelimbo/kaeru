<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { isPartOfURL } from '$lib/utils/strings';
	import { cn } from '$lib/utils/styles';
	import {
		Activity03FreeIcons,
		Folder02FreeIcons,
		GridViewFreeIcons,
		Notification01FreeIcons,
		ShoppingBag01FreeIcons
	} from '@hugeicons/core-free-icons';
	import { HugeiconsIcon, type IconSvgElement } from '@hugeicons/svelte';
	import { BrandIcon } from '../brand';
	import User from './user.svelte';

	type Item = {
		name: string;
		href: string;
		icon: IconSvgElement;
		class?: string;
	};

	const items: Item[] = [
		{
			name: 'My Apps',
			href: '/apps',
			icon: GridViewFreeIcons
		},
		{
			name: 'App Store',
			href: '#',
			icon: ShoppingBag01FreeIcons
		},
		{
			name: 'File explorer',
			href: '#',
			icon: Folder02FreeIcons
		},
		{
			name: 'Usage',
			href: '#',
			icon: Activity03FreeIcons
		},
		{
			name: 'Notifications',
			href: '#',
			icon: Notification01FreeIcons,
			class: 'mt-auto hidden md:block'
		}
	];
</script>

<Tooltip.Root>
	<Tooltip.Trigger class="flex-1 md:flex-0">
		<Button
			href="/"
			variant={isPartOfURL(page.url.pathname, '/', true) ? 'secondary' : 'ghost'}
			class="h-full w-full flex-col rounded-lg md:rounded-xl"
		>
			<BrandIcon class={cn('size-6', isPartOfURL(page.url.pathname, '/', true) && 'fill-brand')} />
			<p class="text-xs md:hidden">Home</p>
		</Button>
	</Tooltip.Trigger>
	<Tooltip.Content side="right">
		<p>Home</p>
	</Tooltip.Content>
</Tooltip.Root>
{#each items as item (item.name)}
	<Tooltip.Root>
		<Tooltip.Trigger class={cn('flex-1 md:flex-0', item.class)}>
			<Button
				href={item.href}
				variant={isPartOfURL(page.url.pathname, item.href) ? 'secondary' : 'ghost'}
				class={cn(
					'h-full w-full flex-col rounded-lg md:rounded-xl',
					isPartOfURL(page.url.pathname, item.href) && 'text-brand'
				)}
			>
				<HugeiconsIcon icon={item.icon} class="size-6" />
				<p class="text-xs md:hidden">{item.name}</p>
			</Button>
		</Tooltip.Trigger>
		<Tooltip.Content side="right" class="hidden md:block">
			<p>{item.name}</p>
		</Tooltip.Content>
	</Tooltip.Root>
{/each}
<User class="hidden md:inline-flex" />
