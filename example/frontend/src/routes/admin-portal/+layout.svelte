<script lang="ts">
	import type { Snippet } from "svelte";
	import type { LayoutData } from "./$types";
	import AdminHeader from "$lib/components/admin-header.svelte";
	import AdminSidebar from "$lib/components/admin-sidebar.svelte";
	import { Toaster } from "$lib/components/ui/sonner/index.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";

	let { data, children }: { data: LayoutData; children: Snippet } = $props();
</script>

<Sidebar.Provider
	style="--sidebar-width: calc(var(--spacing) * 64); --header-height: calc(var(--spacing) * 12);"
>
	<AdminSidebar variant="inset" user={data.user} />
	<Sidebar.Inset>
		<AdminHeader />
		<div class="flex flex-1 flex-col">
			{@render children()}
			<Toaster position="top-right" duration={3000} richColors />
		</div>
	</Sidebar.Inset>
</Sidebar.Provider>
