import { beforeEach, describe, expect, it, vi } from "vitest";

import { POST } from "./+server";
import type { RequestEvent } from "./$types";

const { adjustAdminUserBalanceMock } = vi.hoisted(() => ({
	adjustAdminUserBalanceMock: vi.fn()
}));

vi.mock("$lib/server/admin", () => ({
	adjustAdminUserBalance: adjustAdminUserBalanceMock,
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
		available_credits: 75,
		billing_profile: null
	};
}

function createEvent(body: unknown, user: unknown = { id: "admin-1", role: "admin" }) {
	return {
		request: new Request("http://localhost/api/admin/users/user-1/balance-adjustment", {
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

describe("admin user balance adjustment API endpoint", () => {
	beforeEach(() => {
		adjustAdminUserBalanceMock.mockReset();
	});

	it("adjusts a target user's balance for admins", async () => {
		const result = adminUser();
		adjustAdminUserBalanceMock.mockResolvedValue(result);
		const event = createEvent({ credits_delta: -50 });

		const response = await POST(event);

		expect(adjustAdminUserBalanceMock).toHaveBeenCalledWith(
			event.fetch,
			"auth.session_token=session-1",
			"user-1",
			-50
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("rejects zero, non-integer, and unknown fields", async () => {
		for (const body of [{ credits_delta: 0 }, { credits_delta: 1.5 }, { credits_delta: 10, user_id: "attacker" }]) {
			const response = await POST(createEvent(body));

			expect(response.status).toBe(400);
			expect(await responseJson(response)).toEqual({ error: "invalid balance adjustment payload" });
		}
		expect(adjustAdminUserBalanceMock).not.toHaveBeenCalled();
	});

	it("returns 403 for non-admin users", async () => {
		const response = await POST(createEvent({ credits_delta: 50 }, { id: "user-1", role: "user" }));

		expect(response.status).toBe(403);
		expect(await responseJson(response)).toEqual({ error: "admin access required" });
		expect(adjustAdminUserBalanceMock).not.toHaveBeenCalled();
	});
});
