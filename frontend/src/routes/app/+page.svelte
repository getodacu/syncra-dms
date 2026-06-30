<script lang="ts">
	import Building2Icon from '@lucide/svelte/icons/building-2';
	import UserIcon from '@lucide/svelte/icons/user';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<section class="mx-auto grid max-w-6xl gap-4 px-4 py-6 md:grid-cols-2">
	<Card.Root>
		<Card.Header>
			<Card.Title class="flex items-center gap-2">
				<UserIcon class="size-4 text-primary" />
				Account
			</Card.Title>
			<Card.Description>Authenticated session</Card.Description>
		</Card.Header>
		<Card.Content class="grid gap-2 text-sm">
			<p><span class="font-medium">Name:</span> {data.user?.name}</p>
			<p><span class="font-medium">Email:</span> {data.user?.email}</p>
			<p><span class="font-medium">Session expires:</span> {data.session?.expiresAt}</p>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title class="flex items-center gap-2">
				<Building2Icon class="size-4 text-primary" />
				Organization Units
			</Card.Title>
			<Card.Description>Company structure</Card.Description>
		</Card.Header>
		<Card.Footer>
			<Button href="/app/organization-units" variant="outline" size="sm" class="gap-2">
				<Building2Icon class="size-4" />
				Open
			</Button>
		</Card.Footer>
	</Card.Root>
</section>
