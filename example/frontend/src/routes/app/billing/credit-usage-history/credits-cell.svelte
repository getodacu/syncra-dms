<script lang="ts">
	import { cn } from "$lib/utils.js";
	import { formatCreditsDelta } from "./table-utils";

	let { credits }: { credits: number } = $props();

	const value = $derived(formatCreditsDelta(credits));
	const className = $derived(
		cn(
			"whitespace-nowrap text-sm font-semibold tabular-nums",
			credits > 0 && "text-emerald-700 dark:text-emerald-300",
			credits < 0 && "text-destructive",
			credits === 0 && "text-muted-foreground"
		)
	);
</script>

<span class={className}>{value}</span>
