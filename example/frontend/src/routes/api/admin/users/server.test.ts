import { beforeEach, describe, expect, it, vi } from "vitest";

import { AdminApiError } from "$lib/server/admin";
import { GET } from "./+server";
import type { RequestEvent } from "./$types";

const { listAdminUsersMock, AdminApiErrorMock } = vi.hoisted(() => {
	class MockAdminApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "AdminApiError";
			this.status = status;
		}
	}
	return {
		listAdminUsersMock: vi.fn(),
		AdminApiErrorMock: MockAdminApiError
	};
});

vi.mock("$lib/server/admin", () => ({
	listAdminUsers: listAdminUsersMock,
	AdminApiError: AdminApiErrorMock,
	isAdminApiError: (error: unknown) => error instanceof AdminApiErrorMock
}));

function createEvent(
	url = "http://localhost/api/admin/users?search=ada&sort=last_login_at&direction=desc&cursor=cursor-0&size=50",
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

describe("admin users API endpoint", () => {
	beforeEach(() => {
		listAdminUsersMock.mockReset();
	});

	it("returns 401 for unauthenticated requests", async () => {
		const response = await GET(createEvent("http://localhost/api/admin/users", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listAdminUsersMock).not.toHaveBeenCalled();
	});

	it("returns 403 for authenticated non-admin users", async () => {
		const response = await GET(createEvent("http://localhost/api/admin/users", { id: "user-1", role: "user" }));

		expect(response.status).toBe(403);
		expect(await responseJson(response)).toEqual({ error: "admin access required" });
		expect(listAdminUsersMock).not.toHaveBeenCalled();
	});

	it("forwards admin list filters with the current cookie", async () => {
		const result = { users: [], next_cursor: null };
		listAdminUsersMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await GET(event);

		expect(listAdminUsersMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=session-1", {
			search: "ada",
			sort: "last_login_at",
			direction: "desc",
			cursor: "cursor-0",
			size: "50"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("normalizes backend admin errors", async () => {
		listAdminUsersMock.mockRejectedValue(new AdminApiError(503, "Admin service unavailable"));

		const response = await GET(createEvent("http://localhost/api/admin/users"));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
