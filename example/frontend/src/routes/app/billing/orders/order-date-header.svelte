<script lang="ts">
	import ArrowDownIcon from "@lucide/svelte/icons/arrow-down";
	import ArrowUpIcon from "@lucide/svelte/icons/arrow-up";

	import { Button } from "$lib/components/ui/button/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import type { SortDirection } from "./table-utils";

	let {
		sortDirection,
		onToggle,
	}: {
		sortDirection: SortDirection;
		onToggle: () => void;
	} = $props();

	const nextSortAria = $derived(
		sortDirection === "desc"
			? m.billing_orders_sort_order_date_ascending()
			: m.billing_orders_sort_order_date_descending()
	);
</script>

<Button
	type="button"
	variant="ghost"
	size="sm"
	class="h-8 px-2"
	aria-label={nextSortAria}
	onclick={onToggle}
>
	{m.billing_orders_order_date_filter()}
	{#if sortDirection === "desc"}
		<ArrowDownIcon class="size-4" aria-hidden="true" />
	{:else}
		<ArrowUpIcon class="size-4" aria-hidden="true" />
	{/if}
</Button>
