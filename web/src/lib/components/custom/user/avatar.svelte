<script lang="ts">
	import { page } from '$app/state';
	import type { components } from '$lib/client/v1';
	import * as Avatar from '$lib/components/ui/avatar';
	import { LookTopFreeIcons } from '@hugeicons/core-free-icons';
	import { HugeiconsIcon } from '@hugeicons/svelte';

	const user: components['schemas']['auth.UserResponse'] = $derived(page.data.user);
</script>

<Avatar.Root>
	<!-- <Avatar.Image src={user.avatar} alt={user.username} /> -->
	<Avatar.Fallback>
		<HugeiconsIcon icon={LookTopFreeIcons} />
	</Avatar.Fallback>
</Avatar.Root>
