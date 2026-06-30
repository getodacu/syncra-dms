<script lang="ts">
	import { Badge } from "$lib/components/ui/badge/index.js";
	import type { OCRDocumentListItemResponse } from "./api";

	type Props = {
		document: OCRDocumentListItemResponse;
	};

	let { document }: Props = $props();
</script>

{#if document.collections.length === 0}
	<span class="whitespace-nowrap text-xs text-muted-foreground">-</span>
{:else}
	<div class="flex max-w-[18rem] flex-wrap gap-1.5">
		{#each document.collections as collection (collection.id)}
			<Badge variant="outline" class="max-w-[12rem] px-1.5 text-muted-foreground">
				<span class="truncate">{collection.name}</span>
			</Badge>
		{/each}
	</div>
{/if}
