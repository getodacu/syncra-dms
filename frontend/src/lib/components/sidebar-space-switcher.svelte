<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import LayersIcon from '@lucide/svelte/icons/layers';

	type AppShellUser = {
		role: 'user' | 'admin';
	};

	let { user }: { user: AppShellUser | null } = $props();

	const roleLabel = $derived(user?.role === 'admin' ? 'Admin workspace' : 'User workspace');
</script>

<Sidebar.Menu>
	<Sidebar.MenuItem>
		<Sidebar.MenuButton size="lg" class="data-[slot=sidebar-menu-button]:!p-1.5">
			{#snippet child({ props })}
				<a href="/app" {...props}>
					<LayersIcon class="!size-5" />
					<div class="grid flex-1 text-start leading-tight">
						<span class="truncate text-base font-semibold">Syncra DMS</span>
						<span class="truncate text-xs text-muted-foreground">{roleLabel}</span>
					</div>
				</a>
			{/snippet}
		</Sidebar.MenuButton>
	</Sidebar.MenuItem>
</Sidebar.Menu>
