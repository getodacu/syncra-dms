import { beforeEach, describe, expect, it, vi } from "vitest";

import { AdminApiError } from "$lib/server/admin";
import { GET } from "./+server";
import type { RequestEvent } from "./$types";

const { listAdminBillingOrdersMock, AdminApiErrorMock } = vi.hoisted(() => {
	class MockAdminApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "AdminApiError";
			this.status = status;
		}
	}
	return {
		listAdminBillingOrdersMock: vi.fn(),
		AdminApiErrorMock: MockAdminApiError
	};
});

vi.mock("$lib/server/admin", () => ({
	listAdminBillingOrders: listAdminBillingOrdersMock,
	AdminApiError: AdminApiErrorMock,
	isAdminApiError: (error: unknown) => error instanceof AdminApiErrorMock
}));

function createEvent(
	url = "http://localhost/api/admin/billing/orders?user_id=user-1&status=paid&without_invoice=true&created_from=2026-06-04T00%3A00%3A00Z&created_to=2026-06-05T00%3A00%3A00Z&cursor=cursor-0&size=50&sort=desc",
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

describe("admin billing orders API endpoint", () => {
	beforeEach(() => {
		listAdminBillingOrdersMock.mockReset();
	});

	it("returns 401 for unauthenticated requests", async () => {
		const response = await GET(createEvent("http://localhost/api/admin/billing/orders", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listAdminBillingOrdersMock).not.toHaveBeenCalled();
	});

	it("returns 403 for authenticated non-admin users", async () => {
		const response = await GET(
			createEvent("http://localhost/api/admin/billing/orders", { id: "user-1", role: "user" })
		);

		expect(response.status).toBe(403);
		expect(await responseJson(response)).toEqual({ error: "admin access required" });
		expect(listAdminBillingOrdersMock).not.toHaveBeenCalled();
	});

	it("forwards admin order filters with the current cookie", async () => {
		const result = { orders: [], next_cursor: null };
		listAdminBillingOrdersMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await GET(event);

		expect(listAdminBillingOrdersMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=session-1", {
			userId: "user-1",
			status: "paid",
			withoutInvoice: true,
			createdFrom: "2026-06-04T00:00:00Z",
			createdTo: "2026-06-05T00:00:00Z",
			cursor: "cursor-0",
			size: "50",
			sort: "desc"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("rejects invalid without_invoice values", async () => {
		const response = await GET(createEvent("http://localhost/api/admin/billing/orders?without_invoice=nope"));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid without_invoice" });
		expect(listAdminBillingOrdersMock).not.toHaveBeenCalled();
	});

	it("normalizes backend admin errors", async () => {
		listAdminBillingOrdersMock.mockRejectedValue(new AdminApiError(503, "Admin service unavailable"));

		const response = await GET(createEvent("http://localhost/api/admin/billing/orders"));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
