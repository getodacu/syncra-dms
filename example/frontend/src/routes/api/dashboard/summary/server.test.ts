import { beforeEach, describe, expect, it, vi } from "vitest";

import { DashboardApiError } from "$lib/server/dashboard";
import { GET } from "./+server";
import type { RequestEvent } from "./$types";

const { getDashboardSummaryMock, DashboardApiErrorMock } = vi.hoisted(() => {
	class MockDashboardApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "DashboardApiError";
			this.status = status;
		}
	}

	return {
		getDashboardSummaryMock: vi.fn(),
		DashboardApiErrorMock: MockDashboardApiError
	};
});

vi.mock("$lib/server/dashboard", () => ({
	getDashboardSummary: getDashboardSummaryMock,
	DashboardApiError: DashboardApiErrorMock,
	isDashboardApiError: (error: unknown) => error instanceof DashboardApiErrorMock
}));

function createEvent(path = "http://localhost/api/dashboard/summary", user: unknown = { id: "user-1" }) {
	return {
		request: new Request(path),
		url: new URL(path),
		fetch: vi.fn(),
		locals: { user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

function dashboardSummary() {
	return {
		range: {
			key: "30d",
			start_at: "2026-05-26T00:00:00Z",
			end_at: "2026-06-25T00:00:00Z",
			bucket: "day"
		},
		metrics: {
			documents_processed: 10,
			pages_processed: 42,
			jobs_completed: 8,
			jobs_failed: 1,
			jobs_processing: 1,
			completion_rate: 0.8,
			credits_spent: 120,
			dataset_count: 3,
			schema_count: 2
		},
		document_buckets: [],
		recent_documents: [],
		schema_throughput: [],
		dataset_summary: {
			total_count: 3,
			recent: []
		},
		credit_summary: {
			available_credits: 880,
			credits_spent: 120,
			low_credit: false
		},
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

describe("dashboard summary API endpoint", () => {
	beforeEach(() => {
		getDashboardSummaryMock.mockReset();
	});

	it("returns 401 when unauthenticated", async () => {
		const response = await GET(createEvent("http://localhost/api/dashboard/summary", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(getDashboardSummaryMock).not.toHaveBeenCalled();
	});

	it("loads the summary for the authenticated user id and requested range", async () => {
		const result = dashboardSummary();
		getDashboardSummaryMock.mockResolvedValue(result);
		const event = createEvent("http://localhost/api/dashboard/summary?range=90d");

		const response = await GET(event);

		expect(getDashboardSummaryMock).toHaveBeenCalledWith(event.fetch, {
			userId: "user-1",
			range: "90d"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("returns 400 before backend lookup when the requested range is invalid", async () => {
		const response = await GET(createEvent("http://localhost/api/dashboard/summary?range=14d"));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid range" });
		expect(getDashboardSummaryMock).not.toHaveBeenCalled();
	});

	it.each([
		["leading and trailing spaces", "http://localhost/api/dashboard/summary?range=%207d%20"],
		["only whitespace", "http://localhost/api/dashboard/summary?range=%20%20"]
	])("returns 400 when the requested range contains %s", async (_label, path) => {
		const response = await GET(createEvent(path));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid range" });
		expect(getDashboardSummaryMock).not.toHaveBeenCalled();
	});

	it("ignores user_id query params and forwards the authenticated user id", async () => {
		const result = dashboardSummary();
		getDashboardSummaryMock.mockResolvedValue(result);
		const event = createEvent("http://localhost/api/dashboard/summary?range=7d&user_id=attacker");

		const response = await GET(event);

		expect(getDashboardSummaryMock).toHaveBeenCalledWith(event.fetch, {
			userId: "user-1",
			range: "7d"
		});
		expect(response.status).toBe(200);
	});

	it("passes through public dashboard client errors", async () => {
		getDashboardSummaryMock.mockRejectedValueOnce(new DashboardApiError(400, "invalid dashboard filter"));

		const response = await GET(createEvent());

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid dashboard filter" });
	});

	it("normalizes dashboard client server errors", async () => {
		getDashboardSummaryMock.mockRejectedValueOnce(new DashboardApiError(503, "Dashboard database offline"));

		const response = await GET(createEvent());

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({
			error: "A server error occurred. Please try again."
		});
	});

	it("omits range when the range query is not provided", async () => {
		const result = dashboardSummary();
		getDashboardSummaryMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await GET(event);

		expect(getDashboardSummaryMock).toHaveBeenCalledWith(event.fetch, {
			userId: "user-1"
		});
		expect(response.status).toBe(200);
	});
});
