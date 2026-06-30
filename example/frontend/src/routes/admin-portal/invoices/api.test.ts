import { describe, expect, it, vi } from "vitest";

import {
	buildAdminBillingInvoicePDFPath,
	buildAdminBillingInvoicesPath,
	fetchAdminBillingInvoices,
	generateAdminBillingInvoicePDF
} from "./api";

type FetchInput = Parameters<typeof fetch>[0];
type FetchInit = Parameters<typeof fetch>[1];

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init
	});
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
		pdf_path: "/data/invoices/invoice-1.pdf",
		created_at: "2026-06-11T00:00:00Z",
		updated_at: "2026-06-11T00:00:00Z"
	};
}

describe("admin billing invoices client API", () => {
	it("builds admin billing invoice paths with only non-empty query parameters", () => {
		expect(buildAdminBillingInvoicesPath({ size: 20, sort: "desc" })).toBe(
			"/api/admin/billing/invoices?size=20&sort=desc"
		);
		expect(buildAdminBillingInvoicesPath({ search: " ada@example.com ", size: 20 })).toBe(
			"/api/admin/billing/invoices?search=ada%40example.com&size=20"
		);
		expect(
			buildAdminBillingInvoicesPath({
				userId: "user-1",
				cursor: "cursor-1",
				size: "50",
				sort: "asc"
			})
		).toBe("/api/admin/billing/invoices?user_id=user-1&cursor=cursor-1&size=50&sort=asc");
		expect(
			buildAdminBillingInvoicesPath({
				userId: " ",
				createdFrom: "2026-06-11T00:00:00.000Z",
				createdTo: "2026-06-11T23:59:59.999Z"
			})
		).toBe(
			"/api/admin/billing/invoices?created_from=2026-06-11T00%3A00%3A00.000Z&created_to=2026-06-11T23%3A59%3A59.999Z"
		);
		expect(buildAdminBillingInvoicesPath()).toBe("/api/admin/billing/invoices");
	});

	it("builds admin billing invoice PDF paths", () => {
		expect(buildAdminBillingInvoicePDFPath("invoice-1")).toBe(
			"/api/admin/billing/invoices/invoice-1/pdf"
		);
		expect(buildAdminBillingInvoicePDFPath("invoice/1", { download: true })).toBe(
			"/api/admin/billing/invoices/invoice%2F1/pdf?download=1"
		);
	});

	it("fetches admin billing invoice lists", async () => {
		const result = {
			invoices: [invoiceResponse()],
			next_cursor: null
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => jsonResponse(result));

		await expect(
			fetchAdminBillingInvoices(fetchMock, { search: "SYN-00001", userId: "user-1", size: 20, sort: "desc" })
		).resolves.toEqual(result);
		expect(fetchMock).toHaveBeenCalledWith(
			"/api/admin/billing/invoices?search=SYN-00001&user_id=user-1&size=20&sort=desc",
			{ method: "GET" }
		);
	});

	it("generates admin billing invoice PDFs", async () => {
		const result = invoiceResponse();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => jsonResponse(result));

		await expect(generateAdminBillingInvoicePDF(fetchMock, "invoice-1")).resolves.toEqual(result);
		expect(fetchMock).toHaveBeenCalledWith("/api/admin/billing/invoices/invoice-1/pdf", {
			method: "POST"
		});
	});

	it("throws backend JSON messages for response errors", async () => {
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) =>
			jsonResponse({ error: "invalid user_id" }, { status: 400 })
		);

		await expect(fetchAdminBillingInvoices(fetchMock, { userId: "bad" })).rejects.toThrow("invalid user_id");
	});

	it("rejects invalid list responses", async () => {
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) =>
			jsonResponse({ invoices: [{ ...invoiceResponse(), lines: [{}] }], next_cursor: null })
		);

		await expect(fetchAdminBillingInvoices(fetchMock, {})).rejects.toThrow(
			"Invalid admin billing invoices response"
		);
	});

	it("rejects invalid invoice PDF metadata", async () => {
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) =>
			jsonResponse({ invoices: [{ ...invoiceResponse(), pdf_path: 42 }], next_cursor: null })
		);

		await expect(fetchAdminBillingInvoices(fetchMock, {})).rejects.toThrow(
			"Invalid admin billing invoices response"
		);
	});
});
