<script lang="ts">
	import { scaleUtc } from "d3-scale";
	import { curveNatural } from "d3-shape";
	import { Area, AreaChart } from "layerchart";

	import * as Card from "$lib/components/ui/card/index.js";
	import * as Chart from "$lib/components/ui/chart/index.js";
	import * as Select from "$lib/components/ui/select/index.js";
	import * as ToggleGroup from "$lib/components/ui/toggle-group/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import { getLocale } from "$lib/paraglide/runtime.js";
	import type { DashboardRange, DashboardSummaryResponse } from "./api";
	import { DASHBOARD_RANGES, toChartData } from "./page-state";

	let {
		summary,
		range,
		onRangeChange,
		pending = false
	}: {
		summary: DashboardSummaryResponse;
		range: DashboardRange;
		onRangeChange: (range: DashboardRange) => void;
		pending?: boolean;
	} = $props();

	const chartData = $derived(toChartData(summary));
	const chartConfig = {
		documents: { label: m.dashboard_chart_documents_label(), color: "var(--primary)" }
	} satisfies Chart.ChartConfig;
	const chartDateFormatter = $derived(
		new Intl.DateTimeFormat(getLocale(), {
			month: "short",
			day: "numeric",
			timeZone: "UTC"
		})
	);
	const tooltipDateFormatter = $derived(
		new Intl.DateTimeFormat(getLocale(), {
			month: "short",
			day: "numeric",
			year: "numeric",
			timeZone: "UTC"
		})
	);

	function setRange(value: string) {
		if (value === "") {
			onRangeChange(range);
			return;
		}
		if (value === "7d" || value === "30d" || value === "90d") onRangeChange(value);
	}

	function formatDocumentTick(value: number) {
		return Number.isInteger(value) ? String(value) : "";
	}

	function dashboardRangeLabel(option: DashboardRange) {
		if (option === "7d") return m.dashboard_range_7d();
		if (option === "90d") return m.dashboard_range_90d();
		return m.dashboard_range_30d();
	}
</script>

<Card.Root class="@container/card">
	<Card.Header class="gap-3">
		<div>
			<Card.Title>{m.dashboard_documents_processed_title()}</Card.Title>
			<Card.Description>{dashboardRangeLabel(range)}</Card.Description>
		</div>
		<Card.Action>
			<ToggleGroup.Root
				type="single"
				value={range}
				onValueChange={setRange}
				variant="outline"
				aria-label={m.dashboard_select_range()}
				class="hidden *:data-[slot=toggle-group-item]:!px-4 @[720px]/card:flex"
			>
				{#each DASHBOARD_RANGES as option (option)}
					<ToggleGroup.Item value={option}>{dashboardRangeLabel(option)}</ToggleGroup.Item>
				{/each}
			</ToggleGroup.Root>
			<Select.Root type="single" value={range} onValueChange={setRange}>
				<Select.Trigger
					size="sm"
					class="flex w-36 **:data-[slot=select-value]:block **:data-[slot=select-value]:truncate @[720px]/card:hidden"
					aria-label={m.dashboard_select_range()}
				>
					<span data-slot="select-value">{dashboardRangeLabel(range)}</span>
				</Select.Trigger>
				<Select.Content>
					{#each DASHBOARD_RANGES as option (option)}
						<Select.Item value={option}>{dashboardRangeLabel(option)}</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</Card.Action>
	</Card.Header>
	<Card.Content
		class={pending
			? "px-2 pt-4 opacity-70 transition-opacity sm:px-6"
			: "px-2 pt-4 transition-opacity sm:px-6"}
	>
		<Chart.Container config={chartConfig} class="aspect-auto h-64 w-full">
			<AreaChart
				data={chartData}
				x="date"
				xScale={scaleUtc()}
				series={[
					{
						key: "documents",
						label: m.dashboard_chart_documents_label(),
						color: chartConfig.documents.color
					}
				]}
				props={{
					xAxis: {
						ticks: range === "7d" ? 7 : undefined,
						format: (value: Date) => chartDateFormatter.format(value)
					},
					yAxis: { format: formatDocumentTick }
				}}
			>
				{#snippet marks({ context })}
					<defs>
						<linearGradient id="dashboardDocumentsFill" x1="0" y1="0" x2="0" y2="1">
							<stop offset="5%" stop-color="var(--color-documents)" stop-opacity={0.7} />
							<stop offset="95%" stop-color="var(--color-documents)" stop-opacity={0.08} />
						</linearGradient>
					</defs>
					{#each context.series.visibleSeries as s (s.key)}
						<Area
							seriesKey={s.key}
							curve={curveNatural}
							fillOpacity={0.4}
							line={{ class: "stroke-2" }}
							motion="tween"
							{...s.props}
							fill="url(#dashboardDocumentsFill)"
						/>
					{/each}
				{/snippet}
				{#snippet tooltip()}
					<Chart.Tooltip
						labelFormatter={(value: Date) => tooltipDateFormatter.format(value)}
						indicator="line"
					/>
				{/snippet}
			</AreaChart>
		</Chart.Container>
	</Card.Content>
</Card.Root>
