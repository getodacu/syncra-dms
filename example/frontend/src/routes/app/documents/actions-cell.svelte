<script lang="ts">
	import DotsVerticalIcon from "@tabler/icons-svelte/icons/dots-vertical";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import type { OCRDocumentListItemResponse } from "./api";

	let {
		document,
		onPreview,
		onRename,
		onDownload,
		onDelete,
		renamePending = false,
		downloadPending = false,
		deletePending = false,
	}: {
		document: OCRDocumentListItemResponse;
		onPreview: (document: OCRDocumentListItemResponse) => void;
		onRename: (document: OCRDocumentListItemResponse) => void;
		onDownload: (document: OCRDocumentListItemResponse) => void;
		onDelete: (document: OCRDocumentListItemResponse) => void;
		renamePending?: boolean;
		downloadPending?: boolean;
		deletePending?: boolean;
	} = $props();

	let menuOpen = $state(false);

	function selectRename(event: Event) {
		event.preventDefault();
		menuOpen = false;
		setTimeout(() => {
			onRename(document);
		}, 0);
	}
</script>

<div class="flex items-center justify-end">
	<DropdownMenu.Root bind:open={menuOpen}>
		<DropdownMenu.Trigger>
			{#snippet child({ props })}
				<Button
					type="button"
					variant="ghost"
					size="icon-sm"
					class="text-muted-foreground hover:text-foreground data-[state=open]:bg-muted transition-colors"
					aria-label={m.documents_open_actions_for({ name: document.original_filename })}
					{...props}
				>
					<DotsVerticalIcon class="size-4" aria-hidden="true" />
				</Button>
			{/snippet}
		</DropdownMenu.Trigger>
		<DropdownMenu.Content align="end" class="w-36">
			<DropdownMenu.Item onSelect={() => onPreview(document)}>
				{m.documents_preview()}
			</DropdownMenu.Item>
			<DropdownMenu.Item disabled={renamePending} onSelect={selectRename}>
				{m.documents_rename()}
			</DropdownMenu.Item>
			<DropdownMenu.Item disabled={downloadPending} onSelect={() => onDownload(document)}>
				{m.documents_download()}
			</DropdownMenu.Item>
			<DropdownMenu.Item
				variant="destructive"
				disabled={deletePending}
				onSelect={() => onDelete(document)}
			>
				{m.documents_delete()}
			</DropdownMenu.Item>
		</DropdownMenu.Content>
	</DropdownMenu.Root>
</div>
