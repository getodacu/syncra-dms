<script lang="ts">
	import ArrowDownIcon from "@lucide/svelte/icons/arrow-down";
	import ArrowUpIcon from "@lucide/svelte/icons/arrow-up";

	import { Button } from "$lib/components/ui/button/index.js";
	import type { SortDirection } from "../../app/billing/orders/table-utils";

	let {
		sortDirection,
		onToggle
	}: {
		sortDirection: SortDirection;
		onToggle: () => void;
	} = $props();

	const nextSortDirection = $derived(sortDirection === "desc" ? "ascending" : "descending");
</script>

<Button
	type="button"
	variant="ghost"
	size="sm"
	class="h-8 px-2"
	aria-label={`Sort by created date ${nextSortDirection}`}
	onclick={onToggle}
>
	Created
	{#if sortDirection === "desc"}
		<ArrowDownIcon class="size-4" aria-hidden="true" />
	{:else}
		<ArrowUpIcon class="size-4" aria-hidden="true" />
	{/if}
</Button>
