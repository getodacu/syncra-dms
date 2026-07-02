<script lang="ts">
	import ArchiveIcon from '@lucide/svelte/icons/archive';
	import DownloadIcon from '@lucide/svelte/icons/download';
	import FileTextIcon from '@lucide/svelte/icons/file-text';
	import FolderOpenIcon from '@lucide/svelte/icons/folder-open';
	import PencilIcon from '@lucide/svelte/icons/pencil';
	import { Button } from '$lib/components/ui/button';
	import { documentDownloadHref } from './api';
	import type { RepositoryRow } from './tree';

	let {
		rows,
		canUpdate,
		canDelete,
		canDownload,
		isPending,
		onOpenFolder,
		onRenameFolder,
		onArchiveFolder,
		onRenameDocument,
		onArchiveDocument
	}: {
		rows: RepositoryRow[];
		canUpdate: boolean;
		canDelete: boolean;
		canDownload: boolean;
		isPending: boolean;
		onOpenFolder: (id: string) => void;
		onRenameFolder: (id: string, name: string) => Promise<void>;
		onArchiveFolder: (id: string) => Promise<void>;
		onRenameDocument: (id: string, displayName: string) => Promise<void>;
		onArchiveDocument: (id: string) => Promise<void>;
	} = $props();

	async function renameFolder(row: Extract<RepositoryRow, { type: 'folder' }>) {
		if (isPending) return;

		const name = prompt('Rename folder', row.name);
		if (!name?.trim() || name.trim() === row.name) return;

		try {
			await onRenameFolder(row.id, name.trim());
		} catch {
		}
	}

	async function archiveFolder(row: Extract<RepositoryRow, { type: 'folder' }>) {
		if (isPending) return;
		if (!confirm(`Archive folder "${row.name}"?`)) return;

		try {
			await onArchiveFolder(row.id);
		} catch {
		}
	}

	async function renameDocument(row: Extract<RepositoryRow, { type: 'document' }>) {
		if (isPending) return;

		const displayName = prompt('Rename document', row.name);
		if (!displayName?.trim() || displayName.trim() === row.name) return;

		try {
			await onRenameDocument(row.id, displayName.trim());
		} catch {
		}
	}

	async function archiveDocumentRow(row: Extract<RepositoryRow, { type: 'document' }>) {
		if (isPending) return;
		if (!confirm(`Archive document "${row.name}"?`)) return;

		try {
			await onArchiveDocument(row.id);
		} catch {
		}
	}

	function formatSize(bytes: number) {
		if (bytes < 1024) return `${bytes} B`;
		if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
		return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
	}
</script>

<div class="overflow-x-auto rounded-md border">
	<table class="w-full min-w-[42rem] text-sm">
		<thead class="bg-muted/40 text-xs font-medium text-muted-foreground">
			<tr>
				<th class="h-9 px-3 text-left">Name</th>
				<th class="h-9 px-3 text-left">Type</th>
				<th class="h-9 px-3 text-left">Size</th>
				<th class="h-9 px-3 text-left">Updated</th>
				<th class="h-9 px-3 text-right">Actions</th>
			</tr>
		</thead>
		<tbody>
			{#each rows as row (row.type + row.id)}
				<tr class="border-t">
					<td class="max-w-72 px-3 py-2">
						{#if row.type === 'folder'}
							<button
								type="button"
								class="flex min-w-0 items-center gap-2 text-left font-medium hover:text-primary"
								title={row.name}
								onclick={() => onOpenFolder(row.id)}
							>
								<FolderOpenIcon class="size-4 shrink-0 text-muted-foreground" />
								<span class="truncate">{row.name}</span>
							</button>
						{:else}
							<div class="flex min-w-0 items-center gap-2" title={row.name}>
								<FileTextIcon class="size-4 shrink-0 text-muted-foreground" />
								<span class="truncate font-medium">{row.name}</span>
							</div>
						{/if}
					</td>
					<td class="px-3 py-2 text-muted-foreground">
						{#if row.type === 'folder'}
							Folder
						{:else}
							{row.document.mimeType}
						{/if}
					</td>
					<td class="px-3 py-2 text-muted-foreground">
						{#if row.type === 'folder'}
							-
						{:else}
							{formatSize(row.document.sizeBytes)}
						{/if}
					</td>
					<td class="px-3 py-2 text-muted-foreground">
						{#if row.type === 'folder'}
							{row.folder.updatedAt}
						{:else}
							{row.document.updatedAt}
						{/if}
					</td>
					<td class="px-3 py-2">
						<div class="flex justify-end gap-1">
							{#if row.type === 'folder'}
								<Button
									type="button"
									variant="ghost"
									size="icon-xs"
									title={`Open ${row.name}`}
									aria-label={`Open ${row.name}`}
									disabled={isPending}
									onclick={() => onOpenFolder(row.id)}
								>
									<FolderOpenIcon class="size-3.5" />
								</Button>
								{#if canUpdate}
									<Button
										type="button"
										variant="ghost"
										size="icon-xs"
										title={`Rename ${row.name}`}
										aria-label={`Rename ${row.name}`}
										disabled={isPending}
										onclick={() => renameFolder(row)}
									>
										<PencilIcon class="size-3.5" />
									</Button>
								{/if}
								{#if canDelete}
									<Button
										type="button"
										variant="ghost"
										size="icon-xs"
										title={`Archive ${row.name}`}
										aria-label={`Archive ${row.name}`}
										disabled={isPending}
										onclick={() => archiveFolder(row)}
									>
										<ArchiveIcon class="size-3.5" />
									</Button>
								{/if}
							{:else}
								{#if canDownload}
									<a
										class="inline-flex size-6 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
										href={documentDownloadHref(row.id)}
										title={`Download ${row.name}`}
										aria-label={`Download ${row.name}`}
									>
										<DownloadIcon class="size-3.5" />
									</a>
								{/if}
								{#if canUpdate}
									<Button
										type="button"
										variant="ghost"
										size="icon-xs"
										title={`Rename ${row.name}`}
										aria-label={`Rename ${row.name}`}
										disabled={isPending}
										onclick={() => renameDocument(row)}
									>
										<PencilIcon class="size-3.5" />
									</Button>
								{/if}
								{#if canDelete}
									<Button
										type="button"
										variant="ghost"
										size="icon-xs"
										title={`Archive ${row.name}`}
										aria-label={`Archive ${row.name}`}
										disabled={isPending}
										onclick={() => archiveDocumentRow(row)}
									>
										<ArchiveIcon class="size-3.5" />
									</Button>
								{/if}
							{/if}
						</div>
					</td>
				</tr>
			{:else}
				<tr>
					<td class="h-24 px-3 text-center text-muted-foreground" colspan="5">
						{#if isPending}
							Loading repository
						{:else}
							No items
						{/if}
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
