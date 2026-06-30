import { beforeEach, describe, expect, it, vi } from "vitest";

import { AdminApiError } from "$lib/server/admin";
import { GET, POST } from "./+server";
import type { RequestEvent } from "./$types";

const { generateAdminBillingInvoicePDFMock, fetchAdminBillingInvoicePDFMock, AdminApiErrorMock } =
	vi.hoisted(() => {
		class MockAdminApiError extends Error {
			status: number;

			constructor(status: number, message: string) {
				super(message);
				this.name = "AdminApiError";
				this.status = status;
			}
		}
		return {
			generateAdminBillingInvoicePDFMock: vi.fn(),
			fetchAdminBillingInvoicePDFMock: vi.fn(),
			AdminApiErrorMock: MockAdminApiError
		};
	});

vi.mock("$lib/server/admin", () => ({
	generateAdminBillingInvoicePDF: generateAdminBillingInvoicePDFMock,
	fetchAdminBillingInvoicePDF: fetchAdminBillingInvoicePDFMock,
	AdminApiError: AdminApiErrorMock,
	isAdminApiError: (error: unknown) => error instanceof AdminApiErrorMock
}));

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

function createEvent(
	url = "http://localhost/api/admin/billing/invoices/invoice-1/pdf",
	user: unknown = { id: "admin-1", role: "admin" }
) {
	return {
		request: new Request(url, { method: "POST", headers: { cookie: "auth.session_token=session-1" } }),
		url: new URL(url),
		params: { id: "invoice-1" },
		fetch: vi.fn(),
		locals: { user, adminUser: user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("admin billing invoice PDF API endpoint", () => {
	beforeEach(() => {
		generateAdminBillingInvoicePDFMock.mockReset();
		fetchAdminBillingInvoicePDFMock.mockReset();
	});

	it("returns 401 for unauthenticated POST requests", async () => {
		const response = await POST(createEvent(undefined, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(generateAdminBillingInvoicePDFMock).not.toHaveBeenCalled();
	});

	it("returns 403 for authenticated non-admin GET requests", async () => {
		const response = await GET(createEvent(undefined, { id: "user-1", role: "user" }));

		expect(response.status).toBe(403);
		expect(await responseJson(response)).toEqual({ error: "admin access required" });
		expect(fetchAdminBillingInvoicePDFMock).not.toHaveBeenCalled();
	});

	it("forwards PDF generation with the current cookie", async () => {
		const result = invoiceResponse();
		generateAdminBillingInvoicePDFMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await POST(event);

		expect(generateAdminBillingInvoicePDFMock).toHaveBeenCalledWith(
			event.fetch,
			"auth.session_token=session-1",
			"invoice-1"
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("proxies PDF preview and download responses", async () => {
		fetchAdminBillingInvoicePDFMock.mockResolvedValue({
			body: new Response("%PDF-test").body,
			headers: new Headers({
				"content-type": "application/pdf",
				"content-disposition": 'attachment; filename="invoice-1.pdf"'
			}),
			status: 200
		});
		const event = createEvent("http://localhost/api/admin/billing/invoices/invoice-1/pdf?download=1");

		const response = await GET(event);

		expect(fetchAdminBillingInvoicePDFMock).toHaveBeenCalledWith(
			event.fetch,
			"auth.session_token=session-1",
			"invoice-1",
			{ download: true }
		);
		expect(response.status).toBe(200);
		expect(response.headers.get("content-type")).toBe("application/pdf");
		expect(response.headers.get("content-disposition")).toBe('attachment; filename="invoice-1.pdf"');
		await expect(response.text()).resolves.toBe("%PDF-test");
	});

	it("normalizes backend admin errors", async () => {
		generateAdminBillingInvoicePDFMock.mockRejectedValue(new AdminApiError(502, "failed to convert invoice PDF"));

		const response = await POST(createEvent());

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
