import { beforeEach, describe, expect, it, vi } from "vitest";

const { AuthApiErrorMock, listAuthAccountsMock, unlinkAuthAccountMock } = vi.hoisted(() => {
	class MockAuthApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "AuthApiError";
			this.status = status;
		}
	}

	return {
		AuthApiErrorMock: MockAuthApiError,
		listAuthAccountsMock: vi.fn(),
		unlinkAuthAccountMock: vi.fn()
	};
});

vi.mock("$lib/server/auth", () => ({
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock,
	listAuthAccounts: listAuthAccountsMock,
	unlinkAuthAccount: unlinkAuthAccountMock
}));

import { GET } from "./+server";
import { DELETE } from "./[provider]/+server";

function createEvent(
	method: "DELETE" | "GET",
	options: { provider?: string; user?: unknown; cookie?: string } = {}
) {
	const headers = new Headers();
	if (options.cookie) headers.set("cookie", options.cookie);
	const request = new Request("http://localhost/api/auth/accounts", { method, headers });

	return {
		request,
		fetch: vi.fn(),
		locals: { user: options.user === undefined ? { id: "user-1" } : options.user },
		params: { provider: options.provider ?? "google" }
	};
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("auth accounts proxy endpoints", () => {
	beforeEach(() => {
		listAuthAccountsMock.mockReset();
		unlinkAuthAccountMock.mockReset();
	});

	it("requires authentication before listing accounts", async () => {
		const response = await GET(createEvent("GET", { user: null }) as never);

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listAuthAccountsMock).not.toHaveBeenCalled();
	});

	it("lists accounts with the current cookie", async () => {
		const result = { accounts: [{ id: "account-1", providerId: "google" }] };
		listAuthAccountsMock.mockResolvedValue(result);
		const event = createEvent("GET", { cookie: "auth.session_token=token-1" });

		const response = await GET(event as never);

		expect(listAuthAccountsMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=token-1");
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("unlinks supported OAuth providers with the current cookie", async () => {
		unlinkAuthAccountMock.mockResolvedValue({ deleted_provider_id: "github", deleted_count: 1 });
		const event = createEvent("DELETE", {
			provider: "github",
			cookie: "auth.session_token=token-1"
		});

		const response = await DELETE(event as never);

		expect(unlinkAuthAccountMock).toHaveBeenCalledWith(
			event.fetch,
			"auth.session_token=token-1",
			"github"
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual({ deleted_provider_id: "github", deleted_count: 1 });
	});

	it("rejects unsupported providers before forwarding", async () => {
		const response = await DELETE(createEvent("DELETE", { provider: "credential" }) as never);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "only google and github accounts can be unlinked"
		});
		expect(unlinkAuthAccountMock).not.toHaveBeenCalled();
	});

	it("normalizes backend auth service errors", async () => {
		listAuthAccountsMock.mockRejectedValue(new AuthApiErrorMock(503, "auth service unavailable"));

		const response = await GET(createEvent("GET") as never);

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
