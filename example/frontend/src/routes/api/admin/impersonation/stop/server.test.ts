import { beforeEach, describe, expect, it, vi } from "vitest";

import { AdminApiError } from "$lib/server/admin";
import { POST } from "./+server";
import type { RequestEvent } from "./$types";

const { stopAdminImpersonationMock, AdminApiErrorMock } = vi.hoisted(() => {
	class MockAdminApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "AdminApiError";
			this.status = status;
		}
	}
	return {
		stopAdminImpersonationMock: vi.fn(),
		AdminApiErrorMock: MockAdminApiError
	};
});

vi.mock("$lib/server/admin", () => ({
	stopAdminImpersonation: stopAdminImpersonationMock,
	AdminApiError: AdminApiErrorMock,
	isAdminApiError: (error: unknown) => error instanceof AdminApiErrorMock
}));

function createEvent(user: unknown = { id: "admin-1", role: "admin" }) {
	const url = "http://localhost/api/admin/impersonation/stop";
	return {
		request: new Request(url, { method: "POST", headers: { cookie: "auth.session_token=session-1" } }),
		url: new URL(url),
		fetch: vi.fn(),
		locals: { user, adminUser: user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("admin impersonation stop API endpoint", () => {
	beforeEach(() => {
		stopAdminImpersonationMock.mockReset();
	});

	it("returns 401 for unauthenticated requests", async () => {
		const response = await POST(createEvent(null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(stopAdminImpersonationMock).not.toHaveBeenCalled();
	});

	it("forwards impersonation stop with the current cookie", async () => {
		const result = { session: { id: "session-1" }, user: { id: "admin-1" }, impersonation: null };
		stopAdminImpersonationMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await POST(event);

		expect(stopAdminImpersonationMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=session-1");
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("normalizes backend admin errors", async () => {
		stopAdminImpersonationMock.mockRejectedValue(new AdminApiError(503, "Admin service unavailable"));

		const response = await POST(createEvent());

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
