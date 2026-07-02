<script lang="ts">
	import FileTextIcon from '@lucide/svelte/icons/file-text';
	import type { PageProps } from './$types';

	type DocumentsPageData = {
		canViewDocuments: boolean;
		canCreateDocuments: boolean;
		canUpdateDocuments: boolean;
		canDeleteDocuments: boolean;
		canDownloadDocuments: boolean;
		selectedOrganizationUnitId: string | null;
	};

	let { data }: PageProps = $props();

	const pageData = $derived(data as DocumentsPageData);
	const canManageDocuments = $derived(
		pageData.canCreateDocuments ||
			pageData.canUpdateDocuments ||
			pageData.canDeleteDocuments ||
			pageData.canDownloadDocuments
	);
	const hasDocumentAccess = $derived(pageData.canViewDocuments || canManageDocuments);
	const organizationUnitLabel = $derived(
		pageData.selectedOrganizationUnitId
			? `Organization unit ${pageData.selectedOrganizationUnitId}`
			: 'All organization units'
	);
</script>

<svelte:head>
	<title>Documents | Syncra DMS</title>
</svelte:head>

<div class="flex flex-1 flex-col gap-4 p-4 lg:p-6">
	<div class="min-w-0">
		<div class="flex items-center gap-2">
			<FileTextIcon class="size-5 text-primary" />
			<h2 class="truncate text-xl font-semibold tracking-normal">Documents</h2>
		</div>
		<p class="mt-1 text-sm text-muted-foreground">{organizationUnitLabel}</p>
	</div>

	{#if hasDocumentAccess}
		<section class="rounded-md border bg-muted/30 p-4">
			<h3 class="text-sm font-medium">Document repository</h3>
			<p class="mt-1 text-sm text-muted-foreground">
				Document repository access is available for this account.
			</p>
			{#if canManageDocuments}
				<p class="mt-3 text-xs font-medium uppercase tracking-normal text-muted-foreground">
					Document actions enabled
				</p>
			{/if}
		</section>
	{:else}
		<section class="rounded-md border border-dashed p-4" role="status">
			<h3 class="text-sm font-medium">No document access</h3>
			<p class="mt-1 text-sm text-muted-foreground">
				This account does not have permission to view documents.
			</p>
		</section>
	{/if}
</div>
