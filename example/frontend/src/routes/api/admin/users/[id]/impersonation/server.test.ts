import { beforeEach, describe, expect, it, vi } from "vitest";

import { AdminApiError } from "$lib/server/admin";
import { POST } from "./+server";
import type { RequestEvent } from "./$types";

const { startAdminUserImpersonationMock, AdminApiErrorMock } = vi.hoisted(() => {
	class MockAdminApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "AdminApiError";
			this.status = status;
		}
	}
	return {
		startAdminUserImpersonationMock: vi.fn(),
		AdminApiErrorMock: MockAdminApiError
	};
});

vi.mock("$lib/server/admin", () => ({
	startAdminUserImpersonation: startAdminUserImpersonationMock,
	AdminApiError: AdminApiErrorMock,
	isAdminApiError: (error: unknown) => error instanceof AdminApiErrorMock
}));

function createEvent(user: unknown = { id: "admin-1", role: "admin" }) {
	const url = "http://localhost/api/admin/users/user-1/impersonation";
	return {
		request: new Request(url, { method: "POST", headers: { cookie: "auth.session_token=session-1" } }),
		url: new URL(url),
		params: { id: "user-1" },
		fetch: vi.fn(),
		locals: { user, adminUser: user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("admin user impersonation API endpoint", () => {
	beforeEach(() => {
		startAdminUserImpersonationMock.mockReset();
	});

	it("returns 401 for unauthenticated requests", async () => {
		const response = await POST(createEvent(null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(startAdminUserImpersonationMock).not.toHaveBeenCalled();
	});

	it("forwards impersonation start with the current cookie", async () => {
		const result = { session: { id: "session-1" }, user: { id: "user-1" }, impersonation: {} };
		startAdminUserImpersonationMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await POST(event);

		expect(startAdminUserImpersonationMock).toHaveBeenCalledWith(
			event.fetch,
			"auth.session_token=session-1",
			"user-1"
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("normalizes backend admin errors", async () => {
		startAdminUserImpersonationMock.mockRejectedValue(new AdminApiError(409, "impersonation already active"));

		const response = await POST(createEvent());

		expect(response.status).toBe(409);
		expect(await responseJson(response)).toEqual({ error: "impersonation already active" });
	});
});
