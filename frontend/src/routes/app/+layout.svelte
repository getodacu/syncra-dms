<script lang="ts">
	import AppSidebar from '$lib/components/app-sidebar.svelte';
	import SiteHeader from '$lib/components/site-header.svelte';
	import { ConfirmDeleteDialog } from '$lib/components/ui/confirm-delete-dialog/index.js';
	import { Toaster } from '$lib/components/ui/sonner/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { Snippet } from 'svelte';
	import type { LayoutData } from './$types';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();
</script>

<Sidebar.Provider
	style="--sidebar-width: calc(var(--spacing) * 72); --header-height: calc(var(--spacing) * 12);"
>
	<AppSidebar variant="inset" user={data.user} permissions={data.permissions} />
	<Sidebar.Inset>
		<SiteHeader />
		<div class="flex flex-1 flex-col">
			{@render children()}
			<ConfirmDeleteDialog />
			<Toaster position="top-right" duration={3000} richColors />
		</div>
	</Sidebar.Inset>
</Sidebar.Provider>
