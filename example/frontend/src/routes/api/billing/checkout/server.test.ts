import { beforeEach, describe, expect, it, vi } from "vitest";

import { POST } from "./+server";
import type { RequestEvent } from "./$types";

const {
	attachBillingOrderCheckoutSessionMock,
	createCreditOrderMock,
	createStripeCheckoutSessionMock,
	isBillingApiErrorMock,
	BillingApiErrorMock
} = vi.hoisted(() => {
	class MockBillingApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "BillingApiError";
			this.status = status;
		}
	}

	return {
		attachBillingOrderCheckoutSessionMock: vi.fn(),
		createCreditOrderMock: vi.fn(),
		createStripeCheckoutSessionMock: vi.fn(),
		isBillingApiErrorMock: (error: unknown) => error instanceof MockBillingApiError,
		BillingApiErrorMock: MockBillingApiError
	};
});

vi.mock("$lib/server/billing", () => ({
	attachBillingOrderCheckoutSession: attachBillingOrderCheckoutSessionMock,
	createCreditOrder: createCreditOrderMock,
	isBillingApiError: isBillingApiErrorMock,
	BillingApiError: BillingApiErrorMock
}));

vi.mock("$lib/server/stripe", () => ({
	createStripeCheckoutSession: createStripeCheckoutSessionMock
}));

function checkoutEvent(body: unknown, user: unknown = { id: "user-1" }) {
	return {
		request: new Request("http://localhost/api/billing/checkout", {
			method: "POST",
			headers: { "content-type": "application/json" },
			body: JSON.stringify(body)
		}),
		url: new URL("http://localhost/api/billing/checkout"),
		fetch: vi.fn(),
		locals: { user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

function orderResponse() {
	return {
		id: "order-1",
		user_id: "user-1",
		status: "pending",
		credits: 1000,
		amount_cents: 1000,
		currency: "EUR"
	};
}

describe("billing checkout API endpoint", () => {
	beforeEach(() => {
		attachBillingOrderCheckoutSessionMock.mockReset();
		createCreditOrderMock.mockReset();
		createStripeCheckoutSessionMock.mockReset();
	});

	it("returns 401 for unauthenticated checkout requests", async () => {
		const response = await POST(checkoutEvent({ credits: 1000 }, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(createCreditOrderMock).not.toHaveBeenCalled();
	});

	it("returns 400 for invalid credit quantities", async () => {
		const response = await POST(checkoutEvent({ credits: 1500 }));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "credits must be a multiple of 1000" });
		expect(createCreditOrderMock).not.toHaveBeenCalled();
	});

	it("creates an order, starts Stripe checkout, attaches the session, and returns the checkout URL", async () => {
		const order = orderResponse();
		createCreditOrderMock.mockResolvedValue(order);
		createStripeCheckoutSessionMock.mockResolvedValue({
			id: "cs_test_123",
			url: "https://checkout.stripe.test/session"
		});

		const event = checkoutEvent({ credits: 1000 }, { id: "user-1" });
		const response = await POST(event);

		expect(createCreditOrderMock).toHaveBeenCalledWith(event.fetch, {
			userId: "user-1",
			credits: 1000
		});
		expect(createStripeCheckoutSessionMock).toHaveBeenCalledWith({
			order
		});
		expect(attachBillingOrderCheckoutSessionMock).toHaveBeenCalledWith(event.fetch, {
			orderId: "order-1",
			checkoutSessionId: "cs_test_123"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual({ url: "https://checkout.stripe.test/session" });
	});

	it("preserves backend billing client errors", async () => {
		createCreditOrderMock.mockRejectedValue(
			new BillingApiErrorMock(402, "insufficient account state")
		);

		const response = await POST(checkoutEvent({ credits: 1000 }));

		expect(response.status).toBe(402);
		expect(await responseJson(response)).toEqual({ error: "insufficient account state" });
		expect(createStripeCheckoutSessionMock).not.toHaveBeenCalled();
	});
});
