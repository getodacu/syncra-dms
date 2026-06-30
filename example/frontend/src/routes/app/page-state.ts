import { m } from "$lib/paraglide/messages.js";
import { getLocale, type Locale } from "$lib/paraglide/runtime.js";
import type { DashboardRange, DashboardSummaryResponse } from "./api";

export const DASHBOARD_RANGES = ["7d", "30d", "90d"] as const;

export function normalizeDashboardRange(value: string | null | undefined): DashboardRange {
	if (value === "7d" || value === "30d" || value === "90d") return value;
	return "30d";
}

export function dashboardSummaryQueryKey(range: DashboardRange) {
	return ["dashboard", "summary", range] as const;
}

export function rangeLabel(range: DashboardRange) {
	if (range === "7d") return m.dashboard_range_7d();
	if (range === "90d") return m.dashboard_range_90d();
	return m.dashboard_range_30d();
}

export function formatInteger(value: number, locale: Locale = getLocale()) {
	if (!Number.isFinite(value)) return "0";
	return new Intl.NumberFormat(locale).format(Math.round(value));
}

export function formatPercent(value: number, locale: Locale = getLocale()) {
	if (!Number.isFinite(value) || value <= 0) return "0%";
	return new Intl.NumberFormat(locale, {
		style: "percent",
		maximumFractionDigits: 0
	}).format(value);
}

export function formatDashboardDate(value: string, locale: Locale = getLocale()) {
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return m.common_unknown();

	return new Intl.DateTimeFormat(locale, {
		month: "short",
		day: "numeric",
		year: "numeric",
		timeZone: "UTC"
	}).format(date);
}

export function toChartData(summary: DashboardSummaryResponse) {
	return summary.document_buckets.map((bucket) => ({
		date: new Date(`${bucket.date}T00:00:00Z`),
		documents: bucket.documents_processed
	}));
}

export function shouldShowOnboarding(summary: DashboardSummaryResponse) {
	return summary.onboarding.show_onboarding;
}
