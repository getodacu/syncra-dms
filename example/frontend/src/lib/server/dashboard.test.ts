import { beforeEach, describe, expect, it, vi } from "vitest";

type FetchInput = Parameters<typeof fetch>[0];
type FetchInit = Parameters<typeof fetch>[1];

const INTERNAL_API_HEADER = "X-Syncra-Internal-Token";

const { privateEnv } = vi.hoisted(() => ({
	privateEnv: {} as Record<string, string | undefined>
}));

vi.mock("$env/dynamic/private", () => ({ env: privateEnv }));

function dashboardSummaryFixture() {
	return {
		range: {
			key: "30d",
			start_at: "2026-05-27T00:00:00Z",
			end_at: "2026-06-25T23:59:59Z",
			bucket: "day"
		},
		metrics: {
			documents_processed: 12,
			pages_processed: 48,
			jobs_completed: 10,
			jobs_failed: 1,
			jobs_processing: 2,
			completion_rate: 0.91,
			credits_spent: 24,
			dataset_count: 3,
			schema_count: 2
		},
		document_buckets: [{ date: "2026-06-25", documents_processed: 4 }],
		recent_documents: [
			{
				id: "document-1",
				original_filename: "invoice.pdf",
				schema_id: "schema-1",
				schema_name: "Invoice",
				page_count: 4,
				created_at: "2026-06-25T12:00:00Z"
			},
			{
				id: "document-2",
				original_filename: "inline.pdf",
				schema_id: null,
				schema_name: null,
				page_count: 2,
				created_at: "2026-06-25T13:00:00Z"
			}
		],
		schema_throughput: [
			{ schema_id: "schema-1", schema_name: "Invoice", documents_processed: 8 },
			{ schema_id: null, schema_name: "No schema", documents_processed: 4 }
		],
		dataset_summary: {
			total_count: 3,
			recent: [
				{
					id: "dataset-1",
					name: "Invoices",
					schema_name: "Invoice",
					field_count: 5,
					created_at: "2026-06-25T14:00:00Z"
				}
			]
		},
		credit_summary: {
			available_credits: 1200,
			credits_spent: 24,
			low_credit: false
		},
		onboarding: {
			has_schema: true,
			has_completed_document: true,
			has_dataset: true,
			has_api_key: false,
			has_webhook: true,
			show_onboarding: false
		},
		warnings: [{ section: "credit_summary", message: "Credits are running low" }]
	};
}

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init
	});
}

describe("frontend dashboard server helper", () => {
	beforeEach(() => {
		vi.resetModules();
		for (const key of Object.keys(privateEnv)) delete privateEnv[key];
		process.env.SYNCRA_API_BASE_URL = "http://dashboard-api.test/";
		process.env.SYNCRA_INTERNAL_API_TOKEN = "internal-token";
		process.env.NODE_ENV = "test";
	});

	it("loads a dashboard summary through the backend with the internal API token", async () => {
		const { getDashboardSummary } = await import("./dashboard");
		const summary = dashboardSummaryFixture();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return jsonResponse(summary);
		});

		await expect(
			getDashboardSummary(fetchMock, { userId: "user-1", range: "30d" })
		).resolves.toEqual(summary);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://dashboard-api.test/api/dashboard/summary?user_id=user-1&range=30d",
			expect.objectContaining({ method: "GET" })
		);
		expect(new Headers(fetchMock.mock.calls[0]?.[1]?.headers).get(INTERNAL_API_HEADER)).toBe(
			"internal-token"
		);
	});

	it("rejects before calling the backend when the internal API token is missing", async () => {
		const { getDashboardSummary, DashboardApiError } = await import("./dashboard");
		delete process.env.SYNCRA_INTERNAL_API_TOKEN;
		const fetchMock = vi.fn();

		await expect(getDashboardSummary(fetchMock, { userId: "user-1", range: "7d" })).rejects.toBeInstanceOf(
			DashboardApiError
		);
		await expect(getDashboardSummary(fetchMock, { userId: "user-1", range: "7d" })).rejects.toMatchObject({
			status: 500,
			message: "Dashboard service is not configured"
		});
		expect(fetchMock).not.toHaveBeenCalled();
	});

	it("maps backend error JSON to a typed dashboard API error", async () => {
		const { getDashboardSummary } = await import("./dashboard");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return jsonResponse({ error: "invalid range" }, { status: 400 });
		});

		await expect(getDashboardSummary(fetchMock, { userId: "user-1", range: "bad-range" })).rejects.toMatchObject({
			status: 400,
			message: "invalid range"
		});
	});

	it("maps fetch failures to service unavailable", async () => {
		const { getDashboardSummary } = await import("./dashboard");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			throw new Error("network down");
		});

		await expect(getDashboardSummary(fetchMock, { userId: "user-1", range: "30d" })).rejects.toMatchObject({
			status: 503,
			message: "Dashboard service unavailable"
		});
	});

	it("maps response body read failures to service unavailable", async () => {
		const { getDashboardSummary } = await import("./dashboard");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return {
				ok: true,
				status: 200,
				text: async () => {
					throw new Error("body unavailable");
				}
			} as unknown as Response;
		});

		await expect(getDashboardSummary(fetchMock, { userId: "user-1", range: "30d" })).rejects.toMatchObject({
			status: 503,
			message: "Dashboard service unavailable"
		});
	});

	it("rejects malformed successful responses", async () => {
		const { getDashboardSummary } = await import("./dashboard");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return jsonResponse({ ...dashboardSummaryFixture(), metrics: { documents_processed: "12" } });
		});

		await expect(getDashboardSummary(fetchMock, { userId: "user-1", range: "30d" })).rejects.toMatchObject({
			status: 502,
			message: "Invalid dashboard response"
		});
	});

	it("identifies only dashboard API error instances", async () => {
		const { DashboardApiError, isDashboardApiError } = await import("./dashboard");

		expect(isDashboardApiError(new DashboardApiError(400, "invalid range"))).toBe(true);
		expect(isDashboardApiError(new Error("invalid range"))).toBe(false);
		expect(isDashboardApiError({ status: 400, message: "invalid range" })).toBe(false);
	});
});
