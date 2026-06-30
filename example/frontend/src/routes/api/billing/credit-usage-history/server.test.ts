import { beforeEach, describe, expect, it, vi } from "vitest";

import { BillingApiError } from "$lib/server/billing";
import { GET } from "./+server";
import type { RequestEvent } from "./$types";

const { listCreditUsageHistoryMock, BillingApiErrorMock } = vi.hoisted(() => {
	class MockBillingApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "BillingApiError";
			this.status = status;
		}
	}

	return {
		listCreditUsageHistoryMock: vi.fn(),
		BillingApiErrorMock: MockBillingApiError
	};
});

vi.mock("$lib/server/billing", () => ({
	listCreditUsageHistory: listCreditUsageHistoryMock,
	BillingApiError: BillingApiErrorMock,
	isBillingApiError: (error: unknown) => error instanceof BillingApiErrorMock
}));

function createEvent(url: string, user: unknown = { id: "user-1" }) {
	return {
		url: new URL(url),
		request: new Request(url),
		fetch: vi.fn(),
		locals: { user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("credit usage history API endpoint", () => {
	beforeEach(() => {
		listCreditUsageHistoryMock.mockReset();
	});

	it("returns 401 when unauthenticated", async () => {
		const response = await GET(createEvent("http://localhost/api/billing/credit-usage-history", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listCreditUsageHistoryMock).not.toHaveBeenCalled();
	});

	it("injects the authenticated user id and forwards filters", async () => {
		const result = { credit_usage_history: [], next_cursor: null };
		listCreditUsageHistoryMock.mockResolvedValue(result);
		const event = createEvent(
			"http://localhost/api/billing/credit-usage-history?user_id=attacker&type=debit&created_from=2026-06-04T00%3A00%3A00Z&created_to=2026-06-05T00%3A00%3A00Z&cursor=cursor-1&size=50&sort=asc"
		);

		const response = await GET(event);

		expect(listCreditUsageHistoryMock).toHaveBeenCalledWith(event.fetch, {
			userId: "user-1",
			type: "debit",
			createdFrom: "2026-06-04T00:00:00Z",
			createdTo: "2026-06-05T00:00:00Z",
			cursor: "cursor-1",
			size: "50",
			sort: "asc"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("preserves billing client errors and normalizes server errors", async () => {
		listCreditUsageHistoryMock.mockRejectedValueOnce(new BillingApiError(400, "invalid cursor"));
		let response = await GET(createEvent("http://localhost/api/billing/credit-usage-history?cursor=bad"));
		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid cursor" });

		listCreditUsageHistoryMock.mockRejectedValueOnce(
			new BillingApiError(503, "Billing unavailable")
		);
		response = await GET(createEvent("http://localhost/api/billing/credit-usage-history"));
		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
