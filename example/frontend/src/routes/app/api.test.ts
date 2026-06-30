import { readFileSync } from "node:fs";
import { describe, expect, it, vi } from "vitest";

import {
	fetchDashboardSummary,
	isDashboardSummaryResponse,
	type DashboardSummaryResponse
} from "./api";

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init
	});
}

function summary(): DashboardSummaryResponse {
	return {
		range: {
			key: "7d",
			start_at: "2026-06-19T00:00:00Z",
			end_at: "2026-06-25T15:45:00Z",
			bucket: "day"
		},
		metrics: {
			documents_processed: 0,
			pages_processed: 0,
			jobs_completed: 0,
			jobs_failed: 0,
			jobs_processing: 0,
			completion_rate: 0,
			credits_spent: 0,
			dataset_count: 0,
			schema_count: 0
		},
		document_buckets: [],
		recent_documents: [],
		schema_throughput: [],
		dataset_summary: { total_count: 0, recent: [] },
		credit_summary: { available_credits: 0, credits_spent: 0, low_credit: true },
		onboarding: {
			has_schema: false,
			has_completed_document: false,
			has_dataset: false,
			has_api_key: false,
			has_webhook: false,
			show_onboarding: true
		},
		warnings: []
	};
}

describe("dashboard browser API", () => {
	it("fetches the selected range", async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(summary(), { status: 200 }));

		const result = await fetchDashboardSummary(fetchMock, "7d");

		expect(fetchMock).toHaveBeenCalledWith("/api/dashboard/summary?range=7d", {
			method: "GET"
		});
		expect(result.range.key).toBe("7d");
	});

	it("uses public-safe messages for failed requests", async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "internal details" }, { status: 503 }));

		await expect(fetchDashboardSummary(fetchMock, "30d")).rejects.toThrow(
			"Failed to load dashboard"
		);
	});

	it("rejects malformed payloads", async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ ok: true }, { status: 200 }));

		await expect(fetchDashboardSummary(fetchMock, "30d")).rejects.toThrow(
			"Invalid dashboard response"
		);
	});

	it("rejects malformed but well-shaped successful payloads from fetch", async () => {
		const invalid = summary();
		invalid.metrics.documents_processed = -1;
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(invalid, { status: 200 }));

		await expect(fetchDashboardSummary(fetchMock, "30d")).rejects.toThrow(
			"Invalid dashboard response"
		);
	});

	it("identifies valid dashboard summary responses", () => {
		expect(isDashboardSummaryResponse({ ok: true })).toBe(false);
		expect(isDashboardSummaryResponse(summary())).toBe(true);
	});

	it("rejects out-of-range dashboard completion rates", () => {
		for (const completionRate of [-0.1, 42]) {
			const invalid = summary();
			invalid.metrics.completion_rate = completionRate;

			expect(isDashboardSummaryResponse(invalid)).toBe(false);
		}
	});

	it("rejects invalid dashboard count-like fields", () => {
		const cases: Array<[string, (invalid: DashboardSummaryResponse) => void]> = [
			["metrics.documents_processed", (invalid) => (invalid.metrics.documents_processed = -1)],
			["metrics.pages_processed", (invalid) => (invalid.metrics.pages_processed = 1.5)],
			["metrics.jobs_completed", (invalid) => (invalid.metrics.jobs_completed = -1)],
			["metrics.jobs_failed", (invalid) => (invalid.metrics.jobs_failed = 1.5)],
			["metrics.jobs_processing", (invalid) => (invalid.metrics.jobs_processing = -1)],
			["metrics.credits_spent", (invalid) => (invalid.metrics.credits_spent = 1.5)],
			["metrics.dataset_count", (invalid) => (invalid.metrics.dataset_count = -1)],
			["metrics.schema_count", (invalid) => (invalid.metrics.schema_count = 1.5)],
			[
				"document_buckets.documents_processed",
				(invalid) =>
					(invalid.document_buckets = [{ date: "2026-06-25", documents_processed: -1 }])
			],
			[
				"recent_documents.page_count",
				(invalid) =>
					(invalid.recent_documents = [
						{
							id: "document-1",
							original_filename: "invoice.pdf",
							schema_id: null,
							schema_name: null,
							page_count: 1.5,
							created_at: "2026-06-25T12:00:00Z"
						}
					])
			],
			[
				"schema_throughput.documents_processed",
				(invalid) =>
					(invalid.schema_throughput = [
						{ schema_id: null, schema_name: "Uncategorized", documents_processed: -1 }
					])
			],
			["dataset_summary.total_count", (invalid) => (invalid.dataset_summary.total_count = 1.5)],
			[
				"dataset_summary.recent.field_count",
				(invalid) =>
					(invalid.dataset_summary.recent = [
						{
							id: "dataset-1",
							name: "Invoices",
							schema_name: "Invoice",
							field_count: -1,
							created_at: "2026-06-25T12:00:00Z"
						}
					])
			],
			[
				"credit_summary.available_credits",
				(invalid) => (invalid.credit_summary.available_credits = 1.5)
			],
			[
				"credit_summary.credits_spent",
				(invalid) => (invalid.credit_summary.credits_spent = -1)
			]
		];

		for (const [field, mutate] of cases) {
			const invalid = summary();
			mutate(invalid);

			expect(isDashboardSummaryResponse(invalid), field).toBe(false);
		}
	});

	it("exports only the intended dashboard browser API surface", () => {
		const source = readFileSync(new URL("./api.ts", import.meta.url), "utf8");
		const exportedNames = [
			...source.matchAll(/export\s+(?:async\s+)?(?:function|type|const|class)\s+(\w+)/g)
		]
			.map((match) => match[1])
			.sort();

		expect(exportedNames).toEqual([
			"DashboardRange",
			"DashboardSummaryResponse",
			"fetchDashboardSummary",
			"isDashboardSummaryResponse"
		]);
	});
});
