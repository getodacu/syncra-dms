import { beforeEach, describe, expect, it, vi } from "vitest";

import { POST } from "./+server";
import type { RequestEvent } from "./$types";

const { stopAdminImpersonationMock } = vi.hoisted(() => ({
	stopAdminImpersonationMock: vi.fn()
}));

vi.mock("$lib/server/admin", () => ({
	stopAdminImpersonation: stopAdminImpersonationMock
}));

function createEvent(locals: Record<string, unknown>) {
	const url = "http://localhost/admin-impersonation/stop";
	return {
		request: new Request(url, { method: "POST", headers: { cookie: "auth.session_token=session-1" } }),
		fetch: vi.fn(),
		locals
	} as unknown as RequestEvent;
}

describe("admin impersonation stop form route", () => {
	beforeEach(() => {
		stopAdminImpersonationMock.mockReset();
	});

	it("stops impersonation and redirects back to the target admin detail page", async () => {
		stopAdminImpersonationMock.mockResolvedValueOnce({});
		const event = createEvent({
			user: { id: "user-1", role: "user" },
			adminUser: { id: "admin-1", role: "admin" },
			impersonation: {
				targetUser: { id: "user-1", role: "user" }
			}
		});

		await expect(POST(event)).rejects.toMatchObject({
			status: 303,
			location: "/admin-portal/users/user-1"
		});
		expect(stopAdminImpersonationMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=session-1");
	});

	it("rejects authenticated non-admin users", async () => {
		await expect(
			POST(
				createEvent({
					user: { id: "user-1", role: "user" },
					adminUser: null,
					impersonation: null
				})
			)
		).rejects.toMatchObject({
			status: 403,
			body: { message: "Admin access required" }
		});
		expect(stopAdminImpersonationMock).not.toHaveBeenCalled();
	});
});
