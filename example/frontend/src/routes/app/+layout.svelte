<script lang="ts">
	import type { Snippet } from "svelte";
	import type { LayoutData } from "./$types";
	import { ConfirmDeleteDialog } from "$lib/components/ui/confirm-delete-dialog/index.js";
	import { Toaster } from "$lib/components/ui/sonner/index.js";
	import AppSidebar from "$lib/components/app-sidebar.svelte";
	import ImpersonationBanner from "$lib/components/impersonation-banner.svelte";
	import SiteHeader from "$lib/components/site-header.svelte";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";

	let { data, children }: { data: LayoutData; children: Snippet } = $props();
</script>

<Sidebar.Provider
	style="--sidebar-width: calc(var(--spacing) * 72); --header-height: calc(var(--spacing) * 12);"
>
	<AppSidebar variant="inset" user={data.user} />
	<Sidebar.Inset>
		<SiteHeader
			initialCreditBalance={data.initialCreditBalance}
			initialCreditBalanceError={data.initialCreditBalanceError}
		/>
		<ImpersonationBanner impersonation={data.impersonation} />
		<div class="flex flex-1 flex-col">
			{@render children()}
			<ConfirmDeleteDialog />
			<Toaster position="top-right" duration={3000} richColors />
		</div>
	</Sidebar.Inset>
</Sidebar.Provider>
