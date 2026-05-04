<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Notification01FreeIcons } from '@hugeicons/core-free-icons';
	import { HugeiconsIcon } from '@hugeicons/svelte';
	import { BrandIcon } from '../brand';
	import User from './user.svelte';
</script>

<div class="sticky flex flex-row items-center p-2 md:hidden">
	<BrandIcon class="mr-auto fill-brand" />
	<Button variant="ghost"><HugeiconsIcon icon={Notification01FreeIcons} class="size-6" /></Button>
	<User />
</div>
