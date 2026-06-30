import { Buffer } from "node:buffer";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { POST } from "./+server";
import type { RequestEvent } from "./$types";

const {
	claimBillingInvoiceEmailDeliveryMock,
	constructEventMock,
	fetchBillingInvoicePDFMock,
	isBillingApiErrorMock,
	isStripeCheckoutErrorMock,
	loggerMock,
	markBillingInvoiceEmailSentMock,
	markBillingOrderFailedMock,
	markBillingOrderPaidMock,
	sendInvoicePaidEmailMock,
	stripeMock,
	stripeWebhookSecretMock,
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

	const constructEvent = vi.fn();
	const logger = {
		debug: vi.fn(),
		info: vi.fn(),
		warn: vi.fn(),
		error: vi.fn(),
		child: vi.fn()
	};
	logger.child.mockReturnValue(logger);

	return {
		claimBillingInvoiceEmailDeliveryMock: vi.fn(),
		constructEventMock: constructEvent,
		fetchBillingInvoicePDFMock: vi.fn(),
		isBillingApiErrorMock: (error: unknown) => error instanceof MockBillingApiError,
		isStripeCheckoutErrorMock: () => false,
		loggerMock: logger,
		markBillingInvoiceEmailSentMock: vi.fn(),
		markBillingOrderFailedMock: vi.fn(),
		markBillingOrderPaidMock: vi.fn(),
		sendInvoicePaidEmailMock: vi.fn(),
		stripeMock: vi.fn(() => ({
			webhooks: {
				constructEvent
			}
		})),
		stripeWebhookSecretMock: vi.fn(() => "whsec_test"),
		BillingApiErrorMock: MockBillingApiError
	};
});

vi.mock("$lib/server/billing", () => ({
	claimBillingInvoiceEmailDelivery: claimBillingInvoiceEmailDeliveryMock,
	fetchBillingInvoicePDF: fetchBillingInvoicePDFMock,
	isBillingApiError: isBillingApiErrorMock,
	markBillingInvoiceEmailSent: markBillingInvoiceEmailSentMock,
	markBillingOrderFailed: markBillingOrderFailedMock,
	markBillingOrderPaid: markBillingOrderPaidMock,
	BillingApiError: BillingApiErrorMock
}));

vi.mock("$lib/server/mail", () => ({
	sendInvoicePaidEmail: sendInvoicePaidEmailMock
}));

vi.mock("$lib/server/stripe", () => ({
	isStripeCheckoutError: isStripeCheckoutErrorMock,
	stripe: stripeMock,
	stripeWebhookSecret: stripeWebhookSecretMock
}));

function webhookEvent(body = "{\"id\":\"evt_test\"}", signature = "sig_test") {
	const headers = new Headers({ "content-type": "application/json" });
	if (signature) headers.set("stripe-signature", signature);

	return {
		request: new Request("http://localhost/api/billing/stripe/webhook", {
			method: "POST",
			headers,
			body
		}),
		fetch: vi.fn(),
		locals: { logger: loggerMock }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

function completedEvent() {
	return {
		type: "checkout.session.completed",
		created: 1_777_777_777,
		data: {
			object: {
				id: "cs_test_123",
				object: "checkout.session",
				metadata: { billing_order_id: "order-1" } as Record<string, string>,
				payment_intent: "pi_test_123",
				payment_status: "paid"
			}
		}
	};
}

function paidOrderResponse(overrides: Record<string, unknown> = {}) {
	return {
		id: "order-1",
		user_id: "user-1",
		invoice: {
			id: "invoice-1",
			invoice_serie: "SYN",
			invoice_number: 42,
			invoice_date: "2026-06-13",
			pdf_path: "/data/invoices/invoice-1.pdf"
		},
		order_type: "credit_topup",
		status: "paid",
		provider: "stripe",
		pricing_tier: "tier_2",
		unit_amount_cents: 950,
		credits: 5000,
		amount_cents: 4750,
		currency: "EUR",
		provider_checkout_session_id: "cs_test_123",
		provider_payment_intent_id: "pi_test_123",
		created_at: "2026-06-13T11:00:00Z",
		updated_at: "2026-06-13T12:00:00Z",
		paid_at: "2026-06-13T12:00:00Z",
		...overrides
	};
}

function fullInvoiceResponse(overrides: Record<string, unknown> = {}) {
	return {
		id: "invoice-1",
		user_id: "user-1",
		order_id: "order-1",
		billing_profile_id: "profile-1",
		billing_name: "Ada <Buyer>",
		billing_email: "billing@example.com",
		billing_fiscal_code: "RO2785503",
		billing_profile_snapshot: { billing_name: "Ada <Buyer>" },
		lines: [
			{
				name: "SYNCRA SaaS 5000 credits",
				quantity: 1,
				unit_price: "47.50",
				vat_percentage: "0.00",
				total_vat_amount: "0.00",
				total_amount: "47.50"
			}
		],
		net_amount: "47.50",
		vat_amount: "0.00",
		total_amount: "47.50",
		invoice_date: "2026-06-13",
		invoice_serie: "SYN",
		invoice_number: 42,
		pdf_path: "/data/invoices/invoice-1.pdf",
		email_delivery_claimed_at: "2026-06-13T12:00:00Z",
		created_at: "2026-06-13T12:00:00Z",
		updated_at: "2026-06-13T12:00:00Z",
		...overrides
	};
}

describe("Stripe billing webhook endpoint", () => {
	beforeEach(() => {
		claimBillingInvoiceEmailDeliveryMock.mockReset();
		constructEventMock.mockReset();
		fetchBillingInvoicePDFMock.mockReset();
		loggerMock.debug.mockReset();
		loggerMock.info.mockReset();
		loggerMock.warn.mockReset();
		loggerMock.error.mockReset();
		loggerMock.child.mockReset();
		loggerMock.child.mockReturnValue(loggerMock);
		markBillingInvoiceEmailSentMock.mockReset();
		markBillingOrderFailedMock.mockReset();
		markBillingOrderPaidMock.mockReset();
		sendInvoicePaidEmailMock.mockReset();
		stripeWebhookSecretMock.mockClear();
		stripeMock.mockClear();
		claimBillingInvoiceEmailDeliveryMock.mockResolvedValue({ status: "not_ready" });
		markBillingOrderPaidMock.mockResolvedValue(paidOrderResponse());
		markBillingInvoiceEmailSentMock.mockResolvedValue({
			status: "sent",
			invoice: fullInvoiceResponse({ email_sent_at: "2026-06-13T12:01:00Z" })
		});
		sendInvoicePaidEmailMock.mockResolvedValue(undefined);
	});

	it("rejects requests without a Stripe signature", async () => {
		const response = await POST(webhookEvent("{}", ""));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "missing stripe signature" });
		expect(constructEventMock).not.toHaveBeenCalled();
	});

	it("rejects events with an invalid Stripe signature", async () => {
		constructEventMock.mockImplementation(() => {
			throw new Error("bad signature");
		});

		const response = await POST(webhookEvent("{\"id\":\"evt_bad\"}", "sig_bad"));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid stripe signature" });
		expect(constructEventMock).toHaveBeenCalledWith(
			"{\"id\":\"evt_bad\"}",
			"sig_bad",
			"whsec_test"
		);
		expect(loggerMock.error).toHaveBeenCalledWith("billing.stripe_webhook_signature_failed", {
			error: "bad signature"
		});
	});

	it("marks a billing order paid for completed checkout sessions", async () => {
		const event = completedEvent();
		constructEventMock.mockReturnValue(event);

		const requestEvent = webhookEvent("{\"id\":\"evt_paid\"}", "sig_paid");
		const response = await POST(requestEvent);

		expect(markBillingOrderPaidMock).toHaveBeenCalledWith(requestEvent.fetch, {
			orderId: "order-1",
			checkoutSessionId: "cs_test_123",
			paymentIntentId: "pi_test_123",
			paidAt: "2026-05-03T03:09:37.000Z"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual({ received: true });
	});

	it("sends an invoice email with the generated PDF for claimed paid orders", async () => {
		const event = completedEvent();
		constructEventMock.mockReturnValue(event);
		claimBillingInvoiceEmailDeliveryMock.mockResolvedValue({
			status: "claimed",
			invoice: fullInvoiceResponse()
		});
		fetchBillingInvoicePDFMock.mockResolvedValue({
			body: new Response("%PDF-test").body,
			headers: new Headers({
				"content-type": "application/pdf",
				"content-disposition": 'attachment; filename="SYN_00042_260613.pdf"'
			}),
			status: 200
		});

		const requestEvent = webhookEvent("{\"id\":\"evt_paid\"}", "sig_paid");
		const response = await POST(requestEvent);

		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual({ received: true });
		expect(claimBillingInvoiceEmailDeliveryMock).toHaveBeenCalledWith(requestEvent.fetch, {
			invoiceId: "invoice-1",
			userId: "user-1"
		});
		expect(fetchBillingInvoicePDFMock).toHaveBeenCalledWith(requestEvent.fetch, "user-1", "invoice-1", {
			download: true
		});
		expect(sendInvoicePaidEmailMock).toHaveBeenCalledWith({
			to: "billing@example.com",
			billingName: "Ada <Buyer>",
			invoiceLabel: "SYN-00042",
			invoiceDate: "2026-06-13",
			totalAmount: "47.50",
			currency: "EUR",
			credits: 5000,
			orderId: "order-1",
			pdf: {
				filename: "SYN_00042_260613.pdf",
				content: Buffer.from("%PDF-test"),
				contentType: "application/pdf"
			}
		});
		expect(markBillingInvoiceEmailSentMock).toHaveBeenCalledWith(requestEvent.fetch, {
			invoiceId: "invoice-1",
			userId: "user-1"
		});
	});

	it.each(["not_ready", "already_sent", "claim_active"])(
		"skips invoice email delivery when claim status is %s",
		async (status) => {
			const event = completedEvent();
			constructEventMock.mockReturnValue(event);
			claimBillingInvoiceEmailDeliveryMock.mockResolvedValue({ status });

			const response = await POST(webhookEvent());

			expect(response.status).toBe(200);
			expect(await responseJson(response)).toEqual({ received: true });
			expect(fetchBillingInvoicePDFMock).not.toHaveBeenCalled();
			expect(sendInvoicePaidEmailMock).not.toHaveBeenCalled();
			expect(markBillingInvoiceEmailSentMock).not.toHaveBeenCalled();
		}
	);

	it("logs invoice email delivery failures without failing the webhook", async () => {
		const event = completedEvent();
		constructEventMock.mockReturnValue(event);
		claimBillingInvoiceEmailDeliveryMock.mockResolvedValue({
			status: "claimed",
			invoice: fullInvoiceResponse()
		});
		fetchBillingInvoicePDFMock.mockResolvedValue({
			body: new Response("%PDF-test").body,
			headers: new Headers({ "content-type": "application/pdf" }),
			status: 200
		});
		sendInvoicePaidEmailMock.mockRejectedValue(new Error("smtp down"));

		const response = await POST(webhookEvent());

		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual({ received: true });
		expect(markBillingInvoiceEmailSentMock).not.toHaveBeenCalled();
		expect(loggerMock.error).toHaveBeenCalledWith("billing.invoice_email_delivery_failed", {
			order_id: "order-1",
			invoice_id: "invoice-1",
			error: "smtp down"
		});
	});

	it("does not mark unpaid completed checkout sessions paid", async () => {
		const event = completedEvent();
		event.data.object.payment_status = "unpaid";
		constructEventMock.mockReturnValue(event);

		const response = await POST(webhookEvent());

		expect(markBillingOrderPaidMock).not.toHaveBeenCalled();
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual({ received: true });
	});

	it("marks a billing order paid for async payment success sessions", async () => {
		const event = completedEvent();
		event.type = "checkout.session.async_payment_succeeded";
		constructEventMock.mockReturnValue(event);

		const requestEvent = webhookEvent("{\"id\":\"evt_async_paid\"}", "sig_async_paid");
		const response = await POST(requestEvent);

		expect(markBillingOrderPaidMock).toHaveBeenCalledWith(requestEvent.fetch, {
			orderId: "order-1",
			checkoutSessionId: "cs_test_123",
			paymentIntentId: "pi_test_123",
			paidAt: "2026-05-03T03:09:37.000Z"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual({ received: true });
	});

	it("marks a billing order failed for expired and async failed checkout sessions", async () => {
		constructEventMock.mockReturnValue({
			type: "checkout.session.expired",
			created: 1_777_777_777,
			data: {
				object: {
					id: "cs_expired",
					object: "checkout.session",
					metadata: { billing_order_id: "order-2" }
				}
			}
		});

		const requestEvent = webhookEvent("{\"id\":\"evt_expired\"}", "sig_expired");
		const response = await POST(requestEvent);

		expect(markBillingOrderFailedMock).toHaveBeenCalledWith(requestEvent.fetch, "order-2");
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual({ received: true });

		markBillingOrderFailedMock.mockClear();
		constructEventMock.mockReturnValue({
			type: "checkout.session.async_payment_failed",
			created: 1_777_777_777,
			data: {
				object: {
					id: "cs_failed",
					object: "checkout.session",
					metadata: { billing_order_id: "order-3" }
				}
			}
		});

		const failedEvent = webhookEvent("{\"id\":\"evt_failed\"}", "sig_failed");
		const failedResponse = await POST(failedEvent);

		expect(markBillingOrderFailedMock).toHaveBeenCalledWith(failedEvent.fetch, "order-3");
		expect(failedResponse.status).toBe(200);
		expect(await responseJson(failedResponse)).toEqual({ received: true });
	});

	it("acknowledges stale failed checkout sessions whose billing order no longer exists", async () => {
		constructEventMock.mockReturnValue({
			type: "checkout.session.expired",
			created: 1_777_777_777,
			data: {
				object: {
					id: "cs_expired",
					object: "checkout.session",
					metadata: { billing_order_id: "order-missing" }
				}
			}
		});
		markBillingOrderFailedMock.mockRejectedValue(
			new BillingApiErrorMock(404, "billing order not found")
		);

		const requestEvent = webhookEvent("{\"id\":\"evt_expired\"}", "sig_expired");
		const response = await POST(requestEvent);

		expect(markBillingOrderFailedMock).toHaveBeenCalledWith(
			requestEvent.fetch,
			"order-missing"
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual({ received: true });
		expect(loggerMock.warn).toHaveBeenCalledWith(
			"billing.stripe_webhook_stale_failed_session_ignored",
			{
				checkout_session_id: "cs_expired",
				order_id: "order-missing"
			}
		);
	});

	it("rejects paid sessions that cannot be reconciled to a billing order", async () => {
		const event = completedEvent();
		event.data.object.metadata = {};
		constructEventMock.mockReturnValue(event);

		const response = await POST(webhookEvent());

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "missing billing order metadata" });
		expect(markBillingOrderPaidMock).not.toHaveBeenCalled();
	});

	it("preserves backend billing client errors", async () => {
		constructEventMock.mockReturnValue(completedEvent());
		markBillingOrderPaidMock.mockRejectedValue(new BillingApiErrorMock(409, "order conflict"));

		const response = await POST(webhookEvent());

		expect(response.status).toBe(409);
		expect(await responseJson(response)).toEqual({ error: "order conflict" });
		expect(loggerMock.error).toHaveBeenCalledWith("billing.stripe_webhook_error", {
			error: "order conflict"
		});
	});
});
