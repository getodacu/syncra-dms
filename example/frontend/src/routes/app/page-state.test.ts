import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

import type { DashboardSummaryResponse } from "./api";
import {
	DASHBOARD_RANGES,
	dashboardSummaryQueryKey,
	formatDashboardDate,
	formatInteger,
	formatPercent,
	normalizeDashboardRange,
	rangeLabel,
	shouldShowOnboarding,
	toChartData
} from "./page-state";

function summary(): DashboardSummaryResponse {
	return {
		range: {
			key: "30d",
			start_at: "2026-05-27T00:00:00Z",
			end_at: "2026-06-25T15:45:00Z",
			bucket: "day"
		},
		metrics: {
			documents_processed: 2,
			pages_processed: 5,
			jobs_completed: 1,
			jobs_failed: 1,
			jobs_processing: 0,
			completion_rate: 0.5,
			credits_spent: 7,
			dataset_count: 1,
			schema_count: 1
		},
		document_buckets: [
			{ date: "2026-06-24", documents_processed: 1 },
			{ date: "2026-06-25", documents_processed: 2 }
		],
		recent_documents: [],
		schema_throughput: [],
		dataset_summary: { total_count: 1, recent: [] },
		credit_summary: { available_credits: 43, credits_spent: 7, low_credit: false },
		onboarding: {
			has_schema: true,
			has_completed_document: true,
			has_dataset: true,
			has_api_key: false,
			has_webhook: false,
			show_onboarding: false
		},
		warnings: []
	};
}

describe("dashboard page state", () => {
	it("normalizes dashboard ranges", () => {
		expect(DASHBOARD_RANGES).toEqual(["7d", "30d", "90d"]);
		expect(normalizeDashboardRange("7d")).toBe("7d");
		expect(normalizeDashboardRange("90d")).toBe("90d");
		expect(normalizeDashboardRange("bad")).toBe("30d");
		expect(normalizeDashboardRange(null)).toBe("30d");
	});

	it("labels dashboard ranges", () => {
		expect(rangeLabel("7d")).toBe("Last 7 days");
		expect(rangeLabel("30d")).toBe("Last 30 days");
		expect(rangeLabel("90d")).toBe("Last 90 days");
	});

	it("builds stable query keys", () => {
		expect(dashboardSummaryQueryKey("30d")).toEqual(["dashboard", "summary", "30d"]);
	});

	it("formats dashboard values", () => {
		expect(formatInteger(12345)).toBe("12,345");
		expect(formatPercent(0.5)).toBe("50%");
		expect(formatPercent(Number.NaN)).toBe("0%");
		expect(formatDashboardDate("2026-06-25T12:00:00Z")).toBe("Jun 25, 2026");
	});

	it("maps document buckets to chart data", () => {
		expect(toChartData(summary())).toEqual([
			{ date: new Date("2026-06-24T00:00:00Z"), documents: 1 },
			{ date: new Date("2026-06-25T00:00:00Z"), documents: 2 }
		]);
	});

	it("uses backend onboarding state for dashboard mode", () => {
		const active = summary();
		expect(shouldShowOnboarding(active)).toBe(false);

		const onboarding = {
			...active,
			onboarding: { ...active.onboarding, show_onboarding: true }
		};

		expect(shouldShowOnboarding(onboarding)).toBe(true);
	});

	it("keeps the app dashboard browser client out of server-only modules", () => {
		const source = readFileSync(new URL("./api.ts", import.meta.url), "utf8");

		expect(source).not.toContain("$lib/server");
	});

	it("removes copied demo dashboard imports from the app dashboard", () => {
		const source = readFileSync(new URL("./+page.svelte", import.meta.url), "utf8");

		expect(source).not.toContain("SectionCards");
		expect(source).not.toContain("ChartAreaInteractive");
		expect(source).not.toContain("DataTable");
		expect(source).not.toContain("./data.js");
		expect(source).toContain("fetchDashboardSummary");
		expect(source).toContain("data.initialSummary");
		expect(source).toContain("initialData:");
	});
});
