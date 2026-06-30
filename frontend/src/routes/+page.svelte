<script lang="ts">
	import CircleCheckIcon from '@lucide/svelte/icons/circle-check';
	import CircleXIcon from '@lucide/svelte/icons/circle-x';
	import DatabaseIcon from '@lucide/svelte/icons/database';
	import GitBranchIcon from '@lucide/svelte/icons/git-branch';
	import RefreshCcwIcon from '@lucide/svelte/icons/refresh-ccw';
	import ServerIcon from '@lucide/svelte/icons/server';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const apiConnected = $derived(data.backend.version.ok);
	const databaseReady = $derived(data.backend.readiness.ok);
	const overallReady = $derived(apiConnected && databaseReady);
	const versionLabel = $derived(
		data.backend.version.ok ? data.backend.version.data.version : 'unavailable'
	);
	const versionError = $derived(data.backend.version.ok ? undefined : data.backend.version.error);
	const moduleLabel = $derived(
		data.backend.version.ok ? data.backend.version.data.module : 'ai.ro/syncra/dms'
	);
	const databaseMessage = $derived(
		data.backend.readiness.ok
			? data.backend.readiness.data.status
			: data.backend.readiness.error
	);
</script>

<main class="min-h-screen bg-background text-foreground">
	<section class="border-b bg-card">
		<div class="mx-auto flex max-w-6xl items-center justify-between gap-4 px-4 py-4">
			<div class="min-w-0">
				<p class="text-sm font-medium text-muted-foreground">Syncra DMS</p>
				<h1 class="text-2xl font-semibold tracking-normal">Environment status</h1>
			</div>
			<Badge variant={overallReady ? 'secondary' : 'destructive'} class="h-7">
				{#if overallReady}
					<CircleCheckIcon class="size-3.5" />
					Ready
				{:else}
					<CircleXIcon class="size-3.5" />
					Not ready
				{/if}
			</Badge>
		</div>
	</section>

	<section class="mx-auto grid max-w-6xl gap-4 px-4 py-6 lg:grid-cols-[1fr_20rem]">
		<div class="grid gap-4 md:grid-cols-3">
			<Card.Root size="sm">
				<Card.Header>
					<Card.Title class="flex items-center gap-2">
						<ServerIcon class="size-4 text-primary" />
						Backend API
					</Card.Title>
					<Card.Description>{data.backend.apiBaseUrl}</Card.Description>
				</Card.Header>
				<Card.Content>
					{@render StatusLine(apiConnected, 'Connected', versionError)}
				</Card.Content>
			</Card.Root>

			<Card.Root size="sm">
				<Card.Header>
					<Card.Title class="flex items-center gap-2">
						<DatabaseIcon class="size-4 text-primary" />
						Database
					</Card.Title>
					<Card.Description>Backend readiness probe</Card.Description>
				</Card.Header>
				<Card.Content>
					{@render StatusLine(databaseReady, databaseMessage, databaseMessage)}
				</Card.Content>
			</Card.Root>

			<Card.Root size="sm">
				<Card.Header>
					<Card.Title class="flex items-center gap-2">
						<GitBranchIcon class="size-4 text-primary" />
						Version
					</Card.Title>
					<Card.Description>{moduleLabel}</Card.Description>
				</Card.Header>
				<Card.Content>
					<p class="text-2xl font-semibold">{versionLabel}</p>
				</Card.Content>
			</Card.Root>
		</div>

		<Card.Root size="sm">
			<Card.Header>
				<Card.Title>Workspace</Card.Title>
				<Card.Description>Lean scaffold</Card.Description>
			</Card.Header>
			<Card.Content class="space-y-3 text-sm text-muted-foreground">
				<p>Go API, PostgreSQL, Atlas, SvelteKit, Tailwind, and Paraglide are wired.</p>
				<p>Organization Units are reserved for the first domain feature.</p>
			</Card.Content>
			<Card.Footer>
				<Button href="/" variant="outline" size="sm" class="gap-2">
					<RefreshCcwIcon class="size-4" />
					Refresh
				</Button>
			</Card.Footer>
		</Card.Root>
	</section>
</main>

{#snippet StatusLine(ok: boolean, okText: string | undefined, errorText: string | undefined)}
	<div class="flex items-center gap-2">
		{#if ok}
			<CircleCheckIcon class="size-4 text-emerald-600" />
			<span class="font-medium text-foreground">{okText}</span>
		{:else}
			<CircleXIcon class="size-4 text-destructive" />
			<span class="font-medium text-destructive">{errorText ?? 'Unavailable'}</span>
		{/if}
	</div>
{/snippet}
