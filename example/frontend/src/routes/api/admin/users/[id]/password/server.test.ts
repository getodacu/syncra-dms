import { beforeEach, describe, expect, it, vi } from "vitest";

import { POST } from "./+server";
import type { RequestEvent } from "./$types";

const { resetAdminUserPasswordMock } = vi.hoisted(() => ({
	resetAdminUserPasswordMock: vi.fn()
}));

vi.mock("$lib/server/admin", () => ({
	resetAdminUserPassword: resetAdminUserPasswordMock,
	isAdminApiError: () => false
}));

function createEvent(body: unknown, user: unknown = { id: "admin-1", role: "admin" }) {
	return {
		request: new Request("http://localhost/api/admin/users/user-1/password", {
			method: "POST",
			headers: {
				"content-type": "application/json",
				cookie: "auth.session_token=session-1"
			},
			body: typeof body === "string" ? body : JSON.stringify(body)
		}),
		params: { id: "user-1" },
		fetch: vi.fn(),
		locals: { user, adminUser: user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("admin user password API endpoint", () => {
	beforeEach(() => {
		resetAdminUserPasswordMock.mockReset();
	});

	it("sets a target user's password for admins", async () => {
		resetAdminUserPasswordMock.mockResolvedValue({ ok: true });
		const event = createEvent({ password: "newpassword123" });

		const response = await POST(event);

		expect(resetAdminUserPasswordMock).toHaveBeenCalledWith(
			event.fetch,
			"auth.session_token=session-1",
			"user-1",
			"newpassword123"
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual({ ok: true });
	});

	it("rejects invalid password bodies", async () => {
		const response = await POST(createEvent({ password: 123 }));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid password reset payload" });
		expect(resetAdminUserPasswordMock).not.toHaveBeenCalled();
	});

	it("returns 403 for non-admin users", async () => {
		const response = await POST(createEvent({ password: "newpassword123" }, { id: "user-1", role: "user" }));

		expect(response.status).toBe(403);
		expect(await responseJson(response)).toEqual({ error: "admin access required" });
		expect(resetAdminUserPasswordMock).not.toHaveBeenCalled();
	});
});
