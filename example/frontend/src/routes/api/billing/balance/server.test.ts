import { beforeEach, describe, expect, it, vi } from "vitest";

import { BillingApiError } from "$lib/server/billing";
import { GET } from "./+server";
import type { RequestEvent } from "./$types";

const { getCreditBalanceMock, BillingApiErrorMock } = vi.hoisted(() => {
	class MockBillingApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "BillingApiError";
			this.status = status;
		}
	}

	return {
		getCreditBalanceMock: vi.fn(),
		BillingApiErrorMock: MockBillingApiError,
	};
});

vi.mock("$lib/server/billing", () => ({
	getCreditBalance: getCreditBalanceMock,
	BillingApiError: BillingApiErrorMock,
	isBillingApiError: (error: unknown) => error instanceof BillingApiErrorMock,
}));

function createEvent(user: unknown = { id: "user-1" }) {
	return {
		request: new Request("http://localhost/api/billing/balance"),
		url: new URL("http://localhost/api/billing/balance"),
		fetch: vi.fn(),
		locals: { user },
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("billing balance API endpoint", () => {
	beforeEach(() => {
		getCreditBalanceMock.mockReset();
	});

	it("returns 401 when unauthenticated", async () => {
		const response = await GET(createEvent(null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(getCreditBalanceMock).not.toHaveBeenCalled();
	});

	it("loads the balance for the authenticated user id", async () => {
		const result = { user_id: "user-1", available_credits: 1234 };
		getCreditBalanceMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await GET(event);

		expect(getCreditBalanceMock).toHaveBeenCalledWith(event.fetch, "user-1");
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("preserves billing client errors and normalizes server errors", async () => {
		getCreditBalanceMock.mockRejectedValueOnce(new BillingApiError(400, "invalid user"));
		let response = await GET(createEvent());
		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid user" });

		getCreditBalanceMock.mockRejectedValueOnce(new BillingApiError(503, "Billing offline"));
		response = await GET(createEvent());
		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
