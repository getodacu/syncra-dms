<script lang="ts">
	import FileIcon from "@lucide/svelte/icons/file";
	import FileImageIcon from "@lucide/svelte/icons/file-image";
	import FileTextIcon from "@lucide/svelte/icons/file-text";
	import LoaderIcon from "@lucide/svelte/icons/loader-circle";
	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import type { OCRDocumentListItemResponse } from "./api";
	import { fileIconKind, truncateFilename } from "./table-utils";

	let {
		document,
		schemaName,
		editing = false,
		onEditingChange,
		onPreview,
		onRename,
		renamePending = false,
	}: {
		document: OCRDocumentListItemResponse;
		schemaName?: string;
		editing?: boolean | (() => boolean);
		onEditingChange?: (editing: boolean) => void;
		onPreview: (document: OCRDocumentListItemResponse) => void;
		onRename?: (originalFilename: string) => void | Promise<void>;
		renamePending?: boolean;
	} = $props();

	let draftFilename = $state("");
	let saving = $state(false);
	let saveError = $state<string | null>(null);
	let inputRef = $state<HTMLInputElement | null>(null);
	let editingDocumentId = $state<string | null>(null);

	const displayFilename = $derived(truncateFilename(document.original_filename, 60));
	const iconKind = $derived(fileIconKind(document.mime_type));
	const isEditing = $derived.by(() => (typeof editing === "function" ? editing() : editing));
	const editorPending = $derived(saving || renamePending);

	$effect(() => {
		if (isEditing && editingDocumentId !== document.id) {
			draftFilename = document.original_filename;
			saveError = null;
			editingDocumentId = document.id;
		}

		if (!isEditing) {
			editingDocumentId = null;
		}
	});

	$effect(() => {
		if (!isEditing || !inputRef) return;

		inputRef.focus();
		inputRef.select();
	});

	function setEditing(nextEditing: boolean) {
		onEditingChange?.(nextEditing);
	}

	function cancelEditing() {
		if (editorPending) return;

		draftFilename = document.original_filename;
		saveError = null;
		setEditing(false);
	}

	async function commitRename() {
		const nextFilename = draftFilename.trim();
		if (!nextFilename || nextFilename === document.original_filename) {
			cancelEditing();
			return;
		}

		saving = true;
		saveError = null;
		try {
			await onRename?.(nextFilename);
			setEditing(false);
		} catch (error) {
			saveError = error instanceof Error ? error.message : m.documents_failed_rename();
		} finally {
			saving = false;
		}
	}

	function handleEditorKeydown(event: KeyboardEvent) {
		if (event.key === "Enter") {
			event.preventDefault();
			void commitRename();
			return;
		}

		if (event.key === "Escape") {
			event.preventDefault();
			cancelEditing();
		}
	}
</script>

{#if isEditing}
	<div class="flex min-w-0 items-center gap-2" title={document.original_filename}>
		{#if iconKind === "pdf"}
			<FileTextIcon class="size-4 shrink-0 text-rose-500 dark:text-rose-400" aria-hidden="true" />
		{:else if iconKind === "image"}
			<FileImageIcon class="size-4 shrink-0 text-blue-500 dark:text-blue-400" aria-hidden="true" />
		{:else}
			<FileIcon class="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
		{/if}
		<Input
			bind:ref={inputRef}
			value={draftFilename}
			class="h-7 min-w-[12rem] max-w-full px-2 text-sm font-medium"
			disabled={editorPending}
			aria-label={m.documents_rename_file({ name: document.original_filename })}
			aria-invalid={saveError ? "true" : undefined}
			title={saveError ?? document.original_filename}
			oninput={(event) => {
				draftFilename = event.currentTarget.value;
			}}
			onkeydown={handleEditorKeydown}
			onblur={cancelEditing}
		/>
		{#if editorPending}
			<LoaderIcon class="size-3.5 shrink-0 animate-spin text-muted-foreground" aria-hidden="true" />
		{/if}
		{#if schemaName}
			<Badge
				variant="secondary"
				class="text-[10px] h-4.5 py-0 px-1.5 font-medium tracking-wide bg-blue-500/10 text-blue-600 border border-blue-500/20 dark:bg-blue-500/15 dark:text-blue-400 dark:border-blue-500/30 whitespace-nowrap rounded-sm shrink-0 uppercase"
			>
				{schemaName}
			</Badge>
		{/if}
	</div>
{:else}
	<div class="-mx-4 -my-3 flex w-[calc(100%+2rem)] min-w-0 items-center gap-2 px-4 py-3" title={document.original_filename}>
		{#if iconKind === "pdf"}
			<FileTextIcon class="size-4 shrink-0 text-rose-500 dark:text-rose-400" aria-hidden="true" />
		{:else if iconKind === "image"}
			<FileImageIcon class="size-4 shrink-0 text-blue-500 dark:text-blue-400" aria-hidden="true" />
		{:else}
			<FileIcon class="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
		{/if}
		<button
			type="button"
			class="min-w-0 max-w-full truncate rounded-sm text-left text-sm font-medium underline-offset-2 hover:underline focus-visible:ring-ring focus-visible:ring-2 focus-visible:outline-none disabled:cursor-default disabled:opacity-70"
			disabled={editorPending}
			aria-label={m.documents_preview_file({ name: document.original_filename })}
			onclick={() => onPreview(document)}
		>
			{displayFilename}
		</button>
		{#if editorPending}
			<LoaderIcon class="size-3.5 shrink-0 animate-spin text-muted-foreground" aria-hidden="true" />
		{/if}
		{#if schemaName}
			<Badge
				variant="secondary"
				class="text-[10px] h-4.5 py-0 px-1.5 font-medium tracking-wide bg-blue-500/10 text-blue-600 border border-blue-500/20 dark:bg-blue-500/15 dark:text-blue-400 dark:border-blue-500/30 whitespace-nowrap rounded-sm shrink-0 uppercase"
			>
				{schemaName}
			</Badge>
		{/if}
	</div>
{/if}
