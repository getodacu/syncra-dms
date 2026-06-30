<script lang="ts">
	import ActivityIcon from "@lucide/svelte/icons/activity";
	import CheckCircle2Icon from "@lucide/svelte/icons/check-circle-2";
	import CoinsIcon from "@lucide/svelte/icons/coins";
	import FileTextIcon from "@lucide/svelte/icons/file-text";

	import * as Card from "$lib/components/ui/card/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import type { DashboardSummaryResponse } from "./api";
	import { formatInteger, formatPercent } from "./page-state";

	let {
		summary,
		pending = false
	}: { summary: DashboardSummaryResponse; pending?: boolean } = $props();

	function jobsInProgressLabel(count: number) {
		return count === 1
			? m.dashboard_jobs_in_progress_one({ count })
			: m.dashboard_jobs_in_progress_other({ count });
	}

	const cards = $derived([
		{
			label: m.dashboard_metric_documents_processed(),
			value: formatInteger(summary.metrics.documents_processed),
			subtitle: jobsInProgressLabel(summary.metrics.jobs_processing),
			icon: FileTextIcon
		},
		{
			label: m.dashboard_metric_pages_processed(),
			value: formatInteger(summary.metrics.pages_processed),
			subtitle: m.dashboard_pages_completed(),
			icon: ActivityIcon
		},
		{
			label: m.dashboard_metric_completion_rate(),
			value: formatPercent(summary.metrics.completion_rate),
			subtitle: m.dashboard_completion_summary({
				completed: formatInteger(summary.metrics.jobs_completed),
				failed: formatInteger(summary.metrics.jobs_failed)
			}),
			icon: CheckCircle2Icon
		},
		{
			label: m.dashboard_metric_credits_spent(),
			value: formatInteger(summary.metrics.credits_spent),
			subtitle: m.dashboard_credits_available_short({
				count: formatInteger(summary.credit_summary.available_credits)
			}),
			icon: CoinsIcon
		}
	]);
</script>

<section class="grid gap-3 px-4 sm:grid-cols-2 lg:grid-cols-4 lg:px-6" aria-label={m.dashboard_metrics_aria()}>
	{#each cards as card (card.label)}
		<Card.Root size="sm" class={pending ? "opacity-70 transition-opacity" : "transition-opacity"}>
			<Card.Header class="flex flex-row items-center justify-between gap-3 space-y-0 pb-0">
				<Card.Title class="text-sm font-medium text-muted-foreground">{card.label}</Card.Title>
				<card.icon class="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-semibold tabular-nums">{card.value}</div>
				<p class="mt-1 text-xs text-muted-foreground">{card.subtitle}</p>
			</Card.Content>
		</Card.Root>
	{/each}
</section>
