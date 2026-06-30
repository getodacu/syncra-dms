<script lang="ts">
	import { page } from '$app/state';
	import Building2Icon from '@lucide/svelte/icons/building-2';
	import HomeIcon from '@lucide/svelte/icons/home';
	import LogOutIcon from '@lucide/svelte/icons/log-out';
	import type { Snippet } from 'svelte';
	import { Button } from '$lib/components/ui/button';

	let { children }: { children: Snippet } = $props();

	const pathname = $derived(page.url.pathname);
</script>

<main class="min-h-screen bg-background text-foreground">
	<header class="border-b bg-card">
		<div
			class="mx-auto flex max-w-6xl flex-col gap-3 px-4 py-4 sm:flex-row sm:items-center sm:justify-between"
		>
			<div class="min-w-0">
				<p class="text-sm font-medium text-muted-foreground">Syncra DMS</p>
				<h1 class="truncate text-xl font-semibold tracking-normal">App</h1>
			</div>

			<nav class="flex flex-wrap items-center gap-2" aria-label="App navigation">
				<Button
					href="/app"
					variant="ghost"
					size="sm"
					class="gap-2 data-[active=true]:bg-muted"
					data-active={pathname === '/app'}
					aria-current={pathname === '/app' ? 'page' : undefined}
				>
					<HomeIcon class="size-4" />
					Dashboard
				</Button>
				<Button
					href="/app/organization-units"
					variant="ghost"
					size="sm"
					class="gap-2 data-[active=true]:bg-muted"
					data-active={pathname.startsWith('/app/organization-units')}
					aria-current={pathname.startsWith('/app/organization-units') ? 'page' : undefined}
				>
					<Building2Icon class="size-4" />
					Organization Units
				</Button>
				<form method="POST" action="/logout">
					<Button type="submit" variant="outline" size="sm" class="gap-2">
						<LogOutIcon class="size-4" />
						Logout
					</Button>
				</form>
			</nav>
		</div>
	</header>

	{@render children()}
</main>
