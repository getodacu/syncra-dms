<script lang="ts">
	import FolderIcon from '@lucide/svelte/icons/folder';
	import type { FlatDocumentFolderNode } from './tree';

	let {
		folders,
		selectedId,
		onSelect
	}: {
		folders: FlatDocumentFolderNode[];
		selectedId: string | null;
		onSelect: (id: string) => void;
	} = $props();
</script>

<div class="grid gap-1" role="tree" aria-label="Document folders">
	{#each folders as folder (folder.id)}
		<button
			type="button"
			role="treeitem"
			class="flex h-9 min-w-0 items-center gap-2 rounded-md px-2 text-left text-sm transition-colors hover:bg-muted data-[selected=true]:bg-secondary data-[selected=true]:text-secondary-foreground"
			style={`padding-left: ${0.5 + folder.depth * 1.25}rem`}
			data-selected={folder.id === selectedId}
			aria-level={folder.depth + 1}
			aria-selected={folder.id === selectedId}
			onclick={() => onSelect(folder.id)}
		>
			<FolderIcon class="size-3.5 shrink-0 text-muted-foreground" />
			<span class="min-w-0 truncate font-medium" title={folder.name}>{folder.name}</span>
		</button>
	{:else}
		<div
			class="flex h-16 items-center justify-center rounded-md border border-dashed bg-muted/20 px-3 text-sm text-muted-foreground"
		>
			No folders
		</div>
	{/each}
</div>
