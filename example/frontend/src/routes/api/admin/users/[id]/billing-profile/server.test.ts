import { beforeEach, describe, expect, it, vi } from "vitest";

import { PUT } from "./+server";
import type { RequestEvent } from "./$types";

const { upsertAdminUserBillingProfileMock } = vi.hoisted(() => ({
	upsertAdminUserBillingProfileMock: vi.fn()
}));

vi.mock("$lib/server/admin", () => ({
	upsertAdminUserBillingProfile: upsertAdminUserBillingProfileMock,
	isAdminApiError: () => false
}));

function profileBody(extra: Record<string, unknown> = {}) {
	return {
		entity_type: "company",
		billing_name: "Syncra SRL",
		billing_email: "billing@example.com",
		country_code: "RO",
		address_line1: "Main Street 1",
		city: "Bucharest",
		postal_code: "010101",
		fiscal_code: "RO123",
		...extra
	};
}

function createEvent(body: unknown, user: unknown = { id: "admin-1", role: "admin" }) {
	return {
		request: new Request("http://localhost/api/admin/users/user-1/billing-profile", {
			method: "PUT",
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

describe("admin user billing profile API endpoint", () => {
	beforeEach(() => {
		upsertAdminUserBillingProfileMock.mockReset();
	});

	it("saves a billing profile for the target user", async () => {
		const result = { id: "profile-1", user_id: "user-1", ...profileBody() };
		upsertAdminUserBillingProfileMock.mockResolvedValue(result);
		const event = createEvent(profileBody());

		const response = await PUT(event);

		expect(upsertAdminUserBillingProfileMock).toHaveBeenCalledWith(
			event.fetch,
			"auth.session_token=session-1",
			"user-1",
			profileBody()
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("rejects browser-controlled user ids", async () => {
		const response = await PUT(createEvent(profileBody({ user_id: "attacker" })));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid billing profile request" });
		expect(upsertAdminUserBillingProfileMock).not.toHaveBeenCalled();
	});

	it("returns 403 for non-admin users", async () => {
		const response = await PUT(createEvent(profileBody(), { id: "user-1", role: "user" }));

		expect(response.status).toBe(403);
		expect(await responseJson(response)).toEqual({ error: "admin access required" });
		expect(upsertAdminUserBillingProfileMock).not.toHaveBeenCalled();
	});
});
