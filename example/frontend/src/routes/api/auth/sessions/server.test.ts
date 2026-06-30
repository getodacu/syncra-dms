import { beforeEach, describe, expect, it, vi } from "vitest";

const { AuthApiErrorMock, listAuthSessionsMock, revokeAuthSessionMock } = vi.hoisted(() => {
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
		listAuthSessionsMock: vi.fn(),
		revokeAuthSessionMock: vi.fn()
	};
});

vi.mock("$lib/server/auth", () => ({
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock,
	listAuthSessions: listAuthSessionsMock,
	revokeAuthSession: revokeAuthSessionMock
}));

import { GET } from "./+server";
import { DELETE } from "./[id]/+server";

function createEvent(
	method: "DELETE" | "GET",
	options: { id?: string; user?: unknown; cookie?: string } = {}
) {
	const headers = new Headers();
	if (options.cookie) headers.set("cookie", options.cookie);
	const request = new Request("http://localhost/api/auth/sessions", { method, headers });

	return {
		request,
		fetch: vi.fn(),
		locals: { user: options.user === undefined ? { id: "user-1" } : options.user },
		params: { id: options.id ?? "session-1" }
	};
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("auth sessions proxy endpoints", () => {
	beforeEach(() => {
		listAuthSessionsMock.mockReset();
		revokeAuthSessionMock.mockReset();
	});

	it("requires authentication before listing sessions", async () => {
		const response = await GET(createEvent("GET", { user: null }) as never);

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listAuthSessionsMock).not.toHaveBeenCalled();
	});

	it("lists sessions with the current cookie", async () => {
		const result = { sessions: [{ id: "session-1", current: true }] };
		listAuthSessionsMock.mockResolvedValue(result);
		const event = createEvent("GET", { cookie: "auth.session_token=token-1" });

		const response = await GET(event as never);

		expect(listAuthSessionsMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=token-1");
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("revokes sessions with the current cookie", async () => {
		revokeAuthSessionMock.mockResolvedValue({ deleted_id: "session-1", deleted_count: 1 });
		const event = createEvent("DELETE", {
			id: "session-1",
			cookie: "auth.session_token=token-1"
		});

		const response = await DELETE(event as never);

		expect(revokeAuthSessionMock).toHaveBeenCalledWith(
			event.fetch,
			"auth.session_token=token-1",
			"session-1"
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual({ deleted_id: "session-1", deleted_count: 1 });
	});

	it("normalizes backend auth service errors", async () => {
		listAuthSessionsMock.mockRejectedValue(new AuthApiErrorMock(503, "auth service unavailable"));

		const response = await GET(createEvent("GET") as never);

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
