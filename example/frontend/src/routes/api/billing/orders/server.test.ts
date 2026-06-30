import { beforeEach, describe, expect, it, vi } from "vitest";

import { BillingApiError } from "$lib/server/billing";
import { GET } from "./+server";
import type { RequestEvent } from "./$types";

const { listBillingOrdersMock, BillingApiErrorMock } = vi.hoisted(() => {
	class MockBillingApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "BillingApiError";
			this.status = status;
		}
	}

	return {
		listBillingOrdersMock: vi.fn(),
		BillingApiErrorMock: MockBillingApiError
	};
});

vi.mock("$lib/server/billing", () => ({
	listBillingOrders: listBillingOrdersMock,
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

describe("billing orders API endpoint", () => {
	beforeEach(() => {
		listBillingOrdersMock.mockReset();
	});

	it("returns 401 when unauthenticated", async () => {
		const response = await GET(createEvent("http://localhost/api/billing/orders", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listBillingOrdersMock).not.toHaveBeenCalled();
	});

	it("injects the authenticated user id and forwards filters", async () => {
		const result = { orders: [], next_cursor: null };
		listBillingOrdersMock.mockResolvedValue(result);
		const event = createEvent(
			"http://localhost/api/billing/orders?user_id=attacker&status=paid&created_from=2026-06-04T00%3A00%3A00Z&created_to=2026-06-05T00%3A00%3A00Z&cursor=cursor-1&size=50&sort=asc"
		);

		const response = await GET(event);

		expect(listBillingOrdersMock).toHaveBeenCalledWith(event.fetch, {
			userId: "user-1",
			status: "paid",
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
		listBillingOrdersMock.mockRejectedValueOnce(new BillingApiError(400, "invalid cursor"));
		let response = await GET(createEvent("http://localhost/api/billing/orders?cursor=bad"));
		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid cursor" });

		listBillingOrdersMock.mockRejectedValueOnce(new BillingApiError(503, "Billing unavailable"));
		response = await GET(createEvent("http://localhost/api/billing/orders"));
		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
