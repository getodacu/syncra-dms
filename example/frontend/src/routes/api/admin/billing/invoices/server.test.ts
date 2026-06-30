import { beforeEach, describe, expect, it, vi } from "vitest";
import type { RequestEvent } from "@sveltejs/kit";

import { AdminApiError } from "$lib/server/admin";
import { GET } from "./+server";

const { listAdminBillingInvoicesMock, AdminApiErrorMock } = vi.hoisted(() => {
	class MockAdminApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "AdminApiError";
			this.status = status;
		}
	}
	return {
		listAdminBillingInvoicesMock: vi.fn(),
		AdminApiErrorMock: MockAdminApiError
	};
});

vi.mock("$lib/server/admin", () => ({
	listAdminBillingInvoices: listAdminBillingInvoicesMock,
	AdminApiError: AdminApiErrorMock,
	isAdminApiError: (error: unknown) => error instanceof AdminApiErrorMock
}));

function createEvent(
	url = "http://localhost/api/admin/billing/invoices?search=SYN-00042&user_id=user-1&created_from=2026-06-11T00%3A00%3A00Z&created_to=2026-06-12T00%3A00%3A00Z&cursor=cursor-0&size=50&sort=desc",
	user: unknown = { id: "admin-1", role: "admin" }
) {
	return {
		url: new URL(url),
		request: new Request(url, { headers: { cookie: "auth.session_token=session-1" } }),
		fetch: vi.fn(),
		locals: { user, adminUser: user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("admin billing invoices API endpoint", () => {
	beforeEach(() => {
		listAdminBillingInvoicesMock.mockReset();
	});

	it("returns 401 for unauthenticated requests", async () => {
		const response = await GET(createEvent("http://localhost/api/admin/billing/invoices", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listAdminBillingInvoicesMock).not.toHaveBeenCalled();
	});

	it("returns 403 for authenticated non-admin users", async () => {
		const response = await GET(
			createEvent("http://localhost/api/admin/billing/invoices", { id: "user-1", role: "user" })
		);

		expect(response.status).toBe(403);
		expect(await responseJson(response)).toEqual({ error: "admin access required" });
		expect(listAdminBillingInvoicesMock).not.toHaveBeenCalled();
	});

	it("forwards admin invoice filters with the current cookie", async () => {
		const result = { invoices: [], next_cursor: null };
		listAdminBillingInvoicesMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await GET(event);

		expect(listAdminBillingInvoicesMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=session-1", {
			search: "SYN-00042",
			userId: "user-1",
			createdFrom: "2026-06-11T00:00:00Z",
			createdTo: "2026-06-12T00:00:00Z",
			cursor: "cursor-0",
			size: "50",
			sort: "desc"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("normalizes backend admin errors", async () => {
		listAdminBillingInvoicesMock.mockRejectedValue(new AdminApiError(503, "Admin service unavailable"));

		const response = await GET(createEvent("http://localhost/api/admin/billing/invoices"));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
