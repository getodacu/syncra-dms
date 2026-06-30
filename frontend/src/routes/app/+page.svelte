<script lang="ts">
	import LogOutIcon from '@lucide/svelte/icons/log-out';
	import UserIcon from '@lucide/svelte/icons/user';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<main class="min-h-screen bg-background text-foreground">
	<section class="border-b bg-card">
		<div class="mx-auto flex max-w-5xl items-center justify-between gap-4 px-4 py-4">
			<div>
				<p class="text-sm text-muted-foreground">Syncra DMS</p>
				<h1 class="text-2xl font-semibold">App</h1>
			</div>
			<form method="POST" action="/logout">
				<Button type="submit" variant="outline" class="gap-2">
					<LogOutIcon class="size-4" />
					Logout
				</Button>
			</form>
		</div>
	</section>
	<section class="mx-auto grid max-w-5xl gap-4 px-4 py-6">
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2"><UserIcon class="size-4 text-primary" />Account</Card.Title>
				<Card.Description>Authenticated session</Card.Description>
			</Card.Header>
			<Card.Content class="grid gap-2 text-sm">
				<p><span class="font-medium">Name:</span> {data.user?.name}</p>
				<p><span class="font-medium">Email:</span> {data.user?.email}</p>
				<p><span class="font-medium">Session expires:</span> {data.session?.expiresAt}</p>
			</Card.Content>
		</Card.Root>
	</section>
</main>
