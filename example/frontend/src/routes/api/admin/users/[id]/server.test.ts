import { beforeEach, describe, expect, it, vi } from "vitest";

import { GET, PATCH } from "./+server";
import type { RequestEvent } from "./$types";

const { getAdminUserMock, updateAdminUserMock } = vi.hoisted(() => ({
	getAdminUserMock: vi.fn(),
	updateAdminUserMock: vi.fn()
}));

vi.mock("$lib/server/admin", () => ({
	getAdminUser: getAdminUserMock,
	updateAdminUser: updateAdminUserMock,
	isAdminApiError: () => false
}));

function adminUser() {
	return {
		id: "user-1",
		name: "Ada Lovelace",
		email: "ada@example.com",
		email_verified: true,
		role: "user",
		created_at: "2026-06-01T00:00:00Z",
		updated_at: "2026-06-02T00:00:00Z",
		last_login_at: null,
		billing_profile: null
	};
}

function createEvent(
	body?: unknown,
	user: unknown = { id: "admin-1", role: "admin" }
) {
	const init: RequestInit =
		body === undefined
			? { headers: { cookie: "auth.session_token=session-1" } }
			: {
					method: "PATCH",
					headers: {
						"content-type": "application/json",
						cookie: "auth.session_token=session-1"
					},
					body: JSON.stringify(body)
				};
	return {
		request: new Request("http://localhost/api/admin/users/user-1", init),
		url: new URL("http://localhost/api/admin/users/user-1"),
		params: { id: "user-1" },
		fetch: vi.fn(),
		locals: { user, adminUser: user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("admin user detail API endpoint", () => {
	beforeEach(() => {
		getAdminUserMock.mockReset();
		updateAdminUserMock.mockReset();
	});

	it("loads a user for admins", async () => {
		const result = adminUser();
		getAdminUserMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await GET(event);

		expect(getAdminUserMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=session-1", "user-1");
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("rejects role escalation payloads before forwarding", async () => {
		const response = await PATCH(createEvent({ name: "Ada", role: "admin" }));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid user update payload" });
		expect(updateAdminUserMock).not.toHaveBeenCalled();
	});

	it("updates allowed user profile fields for admins", async () => {
		const result = { ...adminUser(), name: "Ada Updated", email: "ada.updated@example.com" };
		updateAdminUserMock.mockResolvedValue(result);
		const event = createEvent({ name: "Ada Updated", email: "ada.updated@example.com" });

		const response = await PATCH(event);

		expect(updateAdminUserMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=session-1", "user-1", {
			name: "Ada Updated",
			email: "ada.updated@example.com"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("returns 403 for non-admin users", async () => {
		const response = await GET(createEvent(undefined, { id: "user-1", role: "user" }));

		expect(response.status).toBe(403);
		expect(await responseJson(response)).toEqual({ error: "admin access required" });
		expect(getAdminUserMock).not.toHaveBeenCalled();
	});
});
