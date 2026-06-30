import { beforeEach, describe, expect, it, vi } from "vitest";

const { getDashboardSummaryMock, publicErrorMessageMock, DashboardApiErrorMock } = vi.hoisted(() => {
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
		publicErrorMessageMock: vi.fn(),
		DashboardApiErrorMock: MockDashboardApiError
	};
});

vi.mock("$lib/server/dashboard", () => ({
	getDashboardSummary: getDashboardSummaryMock,
	DashboardApiError: DashboardApiErrorMock,
	isDashboardApiError: (error: unknown) => error instanceof DashboardApiErrorMock
}));

vi.mock("$lib/server/public-errors", () => ({
	publicErrorMessage: publicErrorMessageMock
}));

import { load } from "./+page.server";
import type { PageServerLoadEvent } from "./$types";

function loadEvent(user: unknown = { id: "user-1" }) {
	return {
		fetch: vi.fn(),
		locals: { user },
		url: new URL("http://localhost/app")
	} as unknown as PageServerLoadEvent;
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

describe("app dashboard page load", () => {
	beforeEach(() => {
		getDashboardSummaryMock.mockReset();
		publicErrorMessageMock.mockReset();
		publicErrorMessageMock.mockReturnValue("Failed to load dashboard summary");
	});

	it("returns an authentication error shape without a signed-in user", async () => {
		await expect(load(loadEvent(null))).resolves.toEqual({
			initialSummary: null,
			initialSummaryError: "Authentication required"
		});
		expect(getDashboardSummaryMock).not.toHaveBeenCalled();
	});

	it("loads the initial 30 day dashboard summary for the signed-in user", async () => {
		const result = dashboardSummary();
		getDashboardSummaryMock.mockResolvedValue(result);
		const event = loadEvent();

		await expect(load(event)).resolves.toEqual({
			initialSummary: result,
			initialSummaryError: null
		});
		expect(getDashboardSummaryMock).toHaveBeenCalledWith(event.fetch, {
			userId: "user-1",
			range: "30d"
		});
	});

	it("converts dashboard API failures to safe public load errors", async () => {
		getDashboardSummaryMock.mockRejectedValue(new DashboardApiErrorMock(503, "dashboard backend secret"));

		await expect(load(loadEvent())).resolves.toEqual({
			initialSummary: null,
			initialSummaryError: "Failed to load dashboard summary"
		});
		expect(publicErrorMessageMock).toHaveBeenCalledWith(
			503,
			"dashboard backend secret",
			"Failed to load dashboard summary"
		);
	});
});
