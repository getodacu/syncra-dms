import { describe, expect, it, vi } from "vitest";

import {
	buildAdminBillingOrdersPath,
	buildGenerateAdminBillingOrderInvoicePath,
	fetchAdminBillingOrders,
	generateAdminBillingOrderInvoice
} from "./api";

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init
	});
}

function orderResponse() {
	return {
		id: "order-1",
		user_id: "user-1",
		user: {
			id: "user-1",
			name: "Ada Lovelace",
			email: "ada@example.com"
		},
		invoice: null,
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
		paid_at: "2026-06-04T12:05:00Z"
	};
}

function invoiceResponse() {
	return {
		id: "invoice-1",
		user_id: "user-1",
		order_id: "order-1",
		billing_profile_id: "profile-1",
		billing_name: "Ada Lovelace",
		billing_email: "ada@example.com",
		billing_profile_snapshot: {},
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
		invoice_date: "2026-06-11",
		invoice_serie: "SYN",
		invoice_number: 1,
		created_at: "2026-06-11T00:00:00Z",
		updated_at: "2026-06-11T00:00:00Z"
	};
}

describe("admin billing orders client API", () => {
	it("builds admin billing orders paths with only non-empty query parameters", () => {
		expect(buildAdminBillingOrdersPath({ size: 20, sort: "desc" })).toBe(
			"/api/admin/billing/orders?size=20&sort=desc"
		);
		expect(
			buildAdminBillingOrdersPath({
				userId: "user-1",
				status: "paid",
				withoutInvoice: true,
				cursor: "cursor-1",
				size: "50",
				sort: "asc"
			})
		).toBe(
			"/api/admin/billing/orders?user_id=user-1&status=paid&without_invoice=true&cursor=cursor-1&size=50&sort=asc"
		);
		expect(
			buildAdminBillingOrdersPath({
				userId: " ",
				createdFrom: "2026-06-04T00:00:00.000Z",
				createdTo: "2026-06-04T23:59:59.999Z"
			})
		).toBe(
			"/api/admin/billing/orders?created_from=2026-06-04T00%3A00%3A00.000Z&created_to=2026-06-04T23%3A59%3A59.999Z"
		);
		expect(buildAdminBillingOrdersPath()).toBe("/api/admin/billing/orders");
		expect(buildGenerateAdminBillingOrderInvoicePath("order 1")).toBe(
			"/api/admin/billing/orders/order%201/invoice"
		);
	});

	it("fetches admin billing order lists", async () => {
		const result = {
			orders: [
				{
					...orderResponse(),
					invoice: {
						id: "invoice-1",
						invoice_serie: "SYN",
						invoice_number: 1,
						invoice_date: "2026-06-11"
					}
				}
			],
			next_cursor: null
		};
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(result));

		await expect(
			fetchAdminBillingOrders(fetchMock, {
				userId: "user-1",
				status: "paid",
				withoutInvoice: true,
				size: 20,
				sort: "desc"
			})
		).resolves.toEqual(result);
		expect(fetchMock).toHaveBeenCalledWith(
			"/api/admin/billing/orders?user_id=user-1&status=paid&without_invoice=true&size=20&sort=desc",
			{ method: "GET" }
		);
	});

	it("throws backend JSON messages for response errors", async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ error: "invalid user_id" }, { status: 400 }));

		await expect(fetchAdminBillingOrders(fetchMock, { userId: "bad" })).rejects.toThrow("invalid user_id");
	});

	it("rejects invalid list responses", async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ orders: [{ ...orderResponse(), user: null }], next_cursor: null }));

		await expect(fetchAdminBillingOrders(fetchMock, {})).rejects.toThrow(
			"Invalid admin billing orders response"
		);
	});

	it("rejects invalid invoice metadata in list responses", async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
			orders: [{ ...orderResponse(), invoice: { id: "invoice-1", invoice_serie: "SYN" } }],
			next_cursor: null
		}));

		await expect(fetchAdminBillingOrders(fetchMock, {})).rejects.toThrow(
			"Invalid admin billing orders response"
		);
	});

	it("generates admin billing order invoices", async () => {
		const result = invoiceResponse();
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(result, { status: 201 }));

		await expect(generateAdminBillingOrderInvoice(fetchMock, "order-1")).resolves.toEqual(result);
		expect(fetchMock).toHaveBeenCalledWith("/api/admin/billing/orders/order-1/invoice", {
			method: "POST"
		});
	});

	it("throws backend JSON messages for invoice generation errors", async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ error: "billing order already has an invoice" }, { status: 409 }));

		await expect(generateAdminBillingOrderInvoice(fetchMock, "order-1")).rejects.toThrow(
			"billing order already has an invoice"
		);
	});

	it("rejects invalid invoice generation responses", async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ ...invoiceResponse(), lines: [{}] }, { status: 201 }));

		await expect(generateAdminBillingOrderInvoice(fetchMock, "order-1")).rejects.toThrow(
			"Invalid admin billing invoice response"
		);
	});
});
