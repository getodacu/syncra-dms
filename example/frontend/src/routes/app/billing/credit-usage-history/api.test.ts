import { describe, expect, it, vi } from "vitest";

import { fetchCreditUsageHistory } from "./api";

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init,
	});
}

describe("credit usage history client API", () => {
	it("fetches credit usage history lists", async () => {
		const result = {
			credit_usage_history: [
				{
					id: "entry-1",
					created_at: "2026-06-04T12:00:00Z",
					entry_type: "purchase",
					credits_delta: 1000,
				},
			],
			next_cursor: null,
		};
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(result));

		await expect(
			fetchCreditUsageHistory(fetchMock, { type: "purchase", size: 20, sort: "desc" })
		).resolves.toEqual(result);
		expect(fetchMock).toHaveBeenCalledWith(
			"/api/billing/credit-usage-history?type=purchase&size=20&sort=desc",
			{ method: "GET" }
		);
	});

	it("accepts optional related ids", async () => {
		const result = {
			credit_usage_history: [
				{
					id: "entry-1",
					created_at: "2026-06-04T12:00:00Z",
					entry_type: "purchase",
					credits_delta: 1000,
					related_order_id: "order-1",
					related_job_id: "job-1",
				},
			],
			next_cursor: "cursor-1",
		};
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(result));

		await expect(fetchCreditUsageHistory(fetchMock, {})).resolves.toEqual(result);
	});

	it("throws backend JSON messages for response errors", async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "invalid cursor" }, { status: 400 }));

		await expect(fetchCreditUsageHistory(fetchMock, { cursor: "bad" })).rejects.toThrow(
			"invalid cursor"
		);
	});

	it("throws fallback messages for non-JSON response errors", async () => {
		const fetchMock = vi.fn().mockResolvedValue(new Response("nope", { status: 500 }));

		await expect(fetchCreditUsageHistory(fetchMock, {})).rejects.toThrow(
			"Failed to load credit usage history"
		);
	});

	it("rejects invalid list responses", async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValue(jsonResponse({ credit_usage_history: [{}], next_cursor: null }));

		await expect(fetchCreditUsageHistory(fetchMock, {})).rejects.toThrow(
			"Invalid credit usage history response"
		);
	});

	it("rejects invalid entry types and deltas", async () => {
		const fetchMock = vi.fn().mockResolvedValue(
			jsonResponse({
				credit_usage_history: [
					{
						id: "entry-1",
						created_at: "2026-06-04T12:00:00Z",
						entry_type: "grant",
						credits_delta: Number.NaN,
					},
				],
				next_cursor: null,
			})
		);

		await expect(fetchCreditUsageHistory(fetchMock, {})).rejects.toThrow(
			"Invalid credit usage history response"
		);
	});
});
