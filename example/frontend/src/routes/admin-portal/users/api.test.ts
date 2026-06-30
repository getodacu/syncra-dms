import { describe, expect, it, vi } from "vitest";

import { adjustAdminUserBalance, fetchAdminUser, startAdminUserImpersonation } from "./api";

type FetchInput = Parameters<typeof fetch>[0];
type FetchInit = Parameters<typeof fetch>[1];

function adminUserDetailResponse(extra: Record<string, unknown> = {}) {
	return {
		id: "user-1",
		name: "Ada Lovelace",
		email: "ada@example.com",
		email_verified: true,
		role: "user",
		created_at: "2026-06-01T00:00:00Z",
		updated_at: "2026-06-02T00:00:00Z",
		last_login_at: null,
		available_credits: 125,
		billing_profile: null,
		...extra
	};
}

function impersonationResponse() {
	return {
		session: {
			id: "session-1",
			token: "token-1",
			userId: "admin-1",
			expiresAt: "2026-07-01T00:00:00Z",
			createdAt: "2026-06-01T00:00:00Z",
			updatedAt: "2026-06-14T10:00:00Z"
		},
		user: {
			id: "user-1",
			name: "Ada Lovelace",
			email: "ada@example.com",
			role: "user"
		},
		impersonation: {
			adminUser: {
				id: "admin-1",
				name: "Admin User",
				email: "admin@example.com",
				role: "admin"
			},
			targetUser: {
				id: "user-1",
				name: "Ada Lovelace",
				email: "ada@example.com",
				role: "user"
			},
			startedAt: "2026-06-14T10:00:00Z"
		}
	};
}

describe("admin users client API", () => {
	it("loads admin user details with available credits", async () => {
		const result = adminUserDetailResponse();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(result), { status: 200 });
		});

		await expect(fetchAdminUser(fetchMock, "user-1")).resolves.toEqual(result);

		expect(fetchMock).toHaveBeenCalledWith("/api/admin/users/user-1", { method: "GET" });
	});

	it("rejects admin user details without available credits", async () => {
		const { available_credits: _availableCredits, ...invalid } = adminUserDetailResponse();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(invalid), { status: 200 });
		});

		await expect(fetchAdminUser(fetchMock, "user-1")).rejects.toThrow("Invalid admin user response");
	});

	it("adjusts admin user balance", async () => {
		const result = adminUserDetailResponse({ available_credits: 75 });
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(result), { status: 200 });
		});

		await expect(adjustAdminUserBalance(fetchMock, "user-1", -50)).resolves.toEqual(result);

		expect(fetchMock).toHaveBeenCalledWith(
			"/api/admin/users/user-1/balance-adjustment",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({ credits_delta: -50 })
			})
		);
	});

	it("starts admin user impersonation", async () => {
		const result = impersonationResponse();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(result), { status: 200 });
		});

		await expect(startAdminUserImpersonation(fetchMock, "user-1")).resolves.toEqual(result);

		expect(fetchMock).toHaveBeenCalledWith("/api/admin/users/user-1/impersonation", {
			method: "POST"
		});
	});
});
