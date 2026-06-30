import { describe, expect, it, vi } from "vitest";

import { buildBillingInvoicePDFPath, fetchBillingOrders } from "./api";

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init,
	});
}

function orderResponse() {
	return {
		id: "order-1",
		user_id: "user-1",
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
		created_at: "2026-06-04T12:00:00Z",
		updated_at: "2026-06-04T12:05:00Z",
		paid_at: "2026-06-04T12:05:00Z",
	};
}

describe("billing orders client API", () => {
	it("builds billing invoice PDF paths", () => {
		expect(buildBillingInvoicePDFPath("invoice-1")).toBe("/api/billing/invoices/invoice-1/pdf");
		expect(buildBillingInvoicePDFPath("invoice/1", { download: true })).toBe(
			"/api/billing/invoices/invoice%2F1/pdf?download=1"
		);
	});

	it("fetches billing order lists", async () => {
		const result = {
			orders: [orderResponse()],
			next_cursor: null,
		};
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(result));

		await expect(fetchBillingOrders(fetchMock, { status: "paid", size: 20, sort: "desc" })).resolves.toEqual(
			result
		);
		expect(fetchMock).toHaveBeenCalledWith(
			"/api/billing/orders?status=paid&size=20&sort=desc",
			{ method: "GET" }
		);
	});

	it("accepts optional provider ids and terminal timestamps", async () => {
		const result = {
			orders: [
				{
					...orderResponse(),
					provider_checkout_session_id: undefined,
					provider_payment_intent_id: undefined,
					paid_at: undefined,
					failed_at: "2026-06-04T12:05:00Z",
					refunded_at: "2026-06-04T12:06:00Z",
					canceled_at: "2026-06-04T12:07:00Z",
				},
			],
			next_cursor: "cursor-1",
		};
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(result));

		await expect(fetchBillingOrders(fetchMock, {})).resolves.toEqual(result);
	});

	it("accepts optional invoice metadata", async () => {
		const result = {
			orders: [
				{
					...orderResponse(),
					invoice: {
						id: "invoice-1",
						invoice_serie: "SYN",
						invoice_number: 1,
						invoice_date: "2026-06-11",
						pdf_path: "/data/invoices/invoice-1.pdf",
					},
				},
				{
					...orderResponse(),
					id: "order-2",
					invoice: null,
				},
			],
			next_cursor: null,
		};
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(result));

		await expect(fetchBillingOrders(fetchMock, {})).resolves.toEqual(result);
	});

	it("throws backend JSON messages for response errors", async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ error: "invalid cursor" }, { status: 400 }));

		await expect(fetchBillingOrders(fetchMock, { cursor: "bad" })).rejects.toThrow("invalid cursor");
	});

	it("throws fallback messages for non-JSON response errors", async () => {
		const fetchMock = vi.fn().mockResolvedValue(new Response("nope", { status: 500 }));

		await expect(fetchBillingOrders(fetchMock, {})).rejects.toThrow("Failed to load billing orders");
	});

	it("rejects invalid list responses", async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ orders: [{}], next_cursor: null }));

		await expect(fetchBillingOrders(fetchMock, {})).rejects.toThrow("Invalid billing orders response");
	});

	it("rejects invalid invoice metadata", async () => {
		const fetchMock = vi.fn().mockResolvedValue(
			jsonResponse({
				orders: [
					{
						...orderResponse(),
						invoice: {
							id: "invoice-1",
							invoice_serie: "SYN",
							invoice_number: 1,
							invoice_date: "2026-06-11",
							pdf_path: 42,
						},
					},
				],
				next_cursor: null,
			})
		);

		await expect(fetchBillingOrders(fetchMock, {})).rejects.toThrow("Invalid billing orders response");
	});

	it("rejects invalid statuses and monetary values", async () => {
		const fetchMock = vi.fn().mockResolvedValue(
			jsonResponse({
				orders: [
					{
						...orderResponse(),
						status: "settled",
						amount_cents: Number.NaN,
					},
				],
				next_cursor: null,
			})
		);

		await expect(fetchBillingOrders(fetchMock, {})).rejects.toThrow("Invalid billing orders response");
	});
});
