<script lang="ts">
	import { Badge } from "$lib/components/ui/badge/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import type { CreditUsageHistoryEntryType } from "./api";

	let { type }: { type: CreditUsageHistoryEntryType } = $props();

	const label = $derived(
		type === "purchase" ? m.credit_usage_type_purchase() : m.credit_usage_type_debit()
	);
	const className = $derived(
		type === "purchase"
			? "border-emerald-500/20 bg-emerald-500/10 text-emerald-700 dark:border-emerald-500/30 dark:bg-emerald-500/15 dark:text-emerald-300"
			: "border-sky-500/20 bg-sky-500/10 text-sky-700 dark:border-sky-500/30 dark:bg-sky-500/15 dark:text-sky-300"
	);
</script>

<Badge variant="outline" class={className}>
	{label}
</Badge>
