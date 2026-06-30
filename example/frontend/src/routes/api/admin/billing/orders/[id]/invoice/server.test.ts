import { beforeEach, describe, expect, it, vi } from "vitest";

import { AdminApiError } from "$lib/server/admin";
import { POST } from "./+server";
import type { RequestEvent } from "./$types";

const { generateAdminBillingOrderInvoiceMock, AdminApiErrorMock } = vi.hoisted(() => {
	class MockAdminApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "AdminApiError";
			this.status = status;
		}
	}
	return {
		generateAdminBillingOrderInvoiceMock: vi.fn(),
		AdminApiErrorMock: MockAdminApiError
	};
});

vi.mock("$lib/server/admin", () => ({
	generateAdminBillingOrderInvoice: generateAdminBillingOrderInvoiceMock,
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
		created_at: "2026-06-11T00:00:00Z",
		updated_at: "2026-06-11T00:00:00Z"
	};
}

function createEvent(user: unknown = { id: "admin-1", role: "admin" }) {
	const url = "http://localhost/api/admin/billing/orders/order-1/invoice";
	return {
		request: new Request(url, { method: "POST", headers: { cookie: "auth.session_token=session-1" } }),
		url: new URL(url),
		params: { id: "order-1" },
		fetch: vi.fn(),
		locals: { user, adminUser: user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("admin billing order invoice API endpoint", () => {
	beforeEach(() => {
		generateAdminBillingOrderInvoiceMock.mockReset();
	});

	it("returns 401 for unauthenticated requests", async () => {
		const response = await POST(createEvent(null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(generateAdminBillingOrderInvoiceMock).not.toHaveBeenCalled();
	});

	it("returns 403 for authenticated non-admin users", async () => {
		const response = await POST(createEvent({ id: "user-1", role: "user" }));

		expect(response.status).toBe(403);
		expect(await responseJson(response)).toEqual({ error: "admin access required" });
		expect(generateAdminBillingOrderInvoiceMock).not.toHaveBeenCalled();
	});

	it("forwards invoice generation with the current cookie", async () => {
		const result = invoiceResponse();
		generateAdminBillingOrderInvoiceMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await POST(event);

		expect(generateAdminBillingOrderInvoiceMock).toHaveBeenCalledWith(
			event.fetch,
			"auth.session_token=session-1",
			"order-1"
		);
		expect(response.status).toBe(201);
		expect(await responseJson(response)).toEqual(result);
	});

	it("normalizes backend admin errors", async () => {
		generateAdminBillingOrderInvoiceMock.mockRejectedValue(new AdminApiError(409, "already invoiced"));

		const response = await POST(createEvent());

		expect(response.status).toBe(409);
		expect(await responseJson(response)).toEqual({ error: "already invoiced" });
	});
});
