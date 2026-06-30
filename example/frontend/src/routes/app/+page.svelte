<script lang="ts">
	import AlertCircleIcon from "@lucide/svelte/icons/alert-circle";
	import PlusIcon from "@lucide/svelte/icons/plus";
	import RefreshCwIcon from "@lucide/svelte/icons/refresh-cw";
	import { createQuery } from "@tanstack/svelte-query";
	import type { PageProps } from "./$types";

	import * as Alert from "$lib/components/ui/alert/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import { fetchDashboardSummary, type DashboardRange, type DashboardSummaryResponse } from "./api";
	import CreditSummaryPanel from "./credit-summary-panel.svelte";
	import DatasetSummaryPanel from "./dataset-summary-panel.svelte";
	import MetricCards from "./metric-cards.svelte";
	import OnboardingCockpit from "./onboarding-cockpit.svelte";
	import RecentDocumentsPanel from "./recent-documents-panel.svelte";
	import SchemaThroughputPanel from "./schema-throughput-panel.svelte";
	import ThroughputChart from "./throughput-chart.svelte";
	import { dashboardSummaryQueryKey, normalizeDashboardRange, shouldShowOnboarding } from "./page-state";

	let { data }: PageProps = $props();

	let selectedRange = $state<DashboardRange>("30d");

	const initialSummary = $derived(
		selectedRange === data.initialSummary?.range.key ? data.initialSummary : undefined
	);
	const dashboardQuery = createQuery<DashboardSummaryResponse, Error>(() => ({
		queryKey: dashboardSummaryQueryKey(selectedRange),
		queryFn: () => fetchDashboardSummary(fetch, selectedRange),
		initialData: initialSummary
	}));
	const summary = $derived(dashboardQuery.data ?? null);
	const initialLoading = $derived(dashboardQuery.isLoading && summary === null);
	const pending = $derived(dashboardQuery.isFetching && summary !== null);
	const dashboardError = $derived(
		dashboardQuery.isError
			? dashboardQuery.error.message
			: summary
				? null
				: data.initialSummaryError
	);

	function setRange(range: DashboardRange) {
		selectedRange = normalizeDashboardRange(range);
	}

	function retry() {
		void dashboardQuery.refetch();
	}
</script>

<svelte:head>
	<title>{m.nav_dashboard()} | Syncra</title>
</svelte:head>

<div class="@container/main flex flex-1 flex-col gap-4 py-4 md:gap-6 md:py-6">
	<header class="flex flex-col gap-3 px-4 sm:flex-row sm:items-center sm:justify-between lg:px-6">
		<div class="min-w-0">
			<h1 class="truncate text-xl font-semibold tracking-tight">{m.nav_dashboard()}</h1>
			<p class="mt-1 text-sm text-muted-foreground">{m.dashboard_page_description()}</p>
		</div>
		<div class="flex flex-wrap items-center gap-2">
			{#if pending}
				<div class="flex items-center gap-2 text-xs text-muted-foreground" aria-live="polite">
					<Spinner class="size-3.5" aria-label={m.dashboard_refreshing()} />
					<span>{m.dashboard_refreshing()}</span>
				</div>
			{/if}
			<Button href="/app/new-job" size="sm" class="gap-2">
				<PlusIcon class="size-4" aria-hidden="true" />
				{m.dashboard_new_ocr_job()}
			</Button>
		</div>
	</header>

	{#if initialLoading}
		<section class="px-4 lg:px-6">
			<div class="flex min-h-96 items-center justify-center rounded-lg border border-dashed bg-muted/20">
				<div class="flex flex-col items-center gap-3 text-center">
					<Spinner class="size-8 text-muted-foreground" aria-label={m.dashboard_loading_title()} />
					<div>
						<h2 class="font-medium">{m.dashboard_loading_title()}</h2>
						<p class="mt-1 text-sm text-muted-foreground">{m.dashboard_loading_description()}</p>
					</div>
				</div>
			</div>
		</section>
	{:else}
		{#if dashboardError}
			<section class="px-4 lg:px-6">
				<Alert.Root variant="destructive">
					<AlertCircleIcon aria-hidden="true" />
					<Alert.Title>{m.dashboard_unavailable_title()}</Alert.Title>
					<Alert.Description>
						{dashboardError || m.dashboard_unavailable_default()}
					</Alert.Description>
					<Alert.Action>
						<Button type="button" variant="outline" size="sm" onclick={retry}>
							<RefreshCwIcon class="size-4" aria-hidden="true" />
							{m.common_retry()}
						</Button>
					</Alert.Action>
				</Alert.Root>
			</section>
		{/if}

		{#if summary}
			{#if summary.warnings.length > 0}
				<section class="px-4 lg:px-6">
					<Alert.Root>
						<AlertCircleIcon aria-hidden="true" />
						<Alert.Title>{m.dashboard_warning_title()}</Alert.Title>
						<Alert.Description>
							<ul class="list-disc space-y-1 pl-4">
								{#each summary.warnings as warning (warning.section)}
									<li>{warning.message}</li>
								{/each}
							</ul>
						</Alert.Description>
					</Alert.Root>
				</section>
			{/if}

			{#if shouldShowOnboarding(summary)}
				<OnboardingCockpit {summary} />
			{/if}

			<MetricCards {summary} {pending} />

			<section class="px-4 lg:px-6">
				<ThroughputChart {summary} range={selectedRange} onRangeChange={setRange} {pending} />
			</section>

			<section class="grid gap-4 px-4 lg:grid-cols-2 lg:px-6 2xl:grid-cols-4">
				<RecentDocumentsPanel {summary} />
				<SchemaThroughputPanel {summary} />
				<DatasetSummaryPanel {summary} />
				<CreditSummaryPanel {summary} />
			</section>
		{:else if !dashboardError}
			<section class="px-4 lg:px-6">
				<Alert.Root variant="destructive">
					<AlertCircleIcon aria-hidden="true" />
					<Alert.Title>{m.dashboard_unavailable_title()}</Alert.Title>
					<Alert.Description>{m.dashboard_unavailable_default()}</Alert.Description>
					<Alert.Action>
						<Button type="button" variant="outline" size="sm" onclick={retry}>
							<RefreshCwIcon class="size-4" aria-hidden="true" />
							{m.common_retry()}
						</Button>
					</Alert.Action>
				</Alert.Root>
			</section>
		{/if}
	{/if}
</div>
