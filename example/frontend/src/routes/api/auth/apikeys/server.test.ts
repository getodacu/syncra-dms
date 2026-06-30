import { beforeEach, describe, expect, it, vi } from "vitest";

import { AuthApiError } from "$lib/server/auth";
import { DELETE, GET, POST } from "./+server";
import type { RequestEvent } from "./$types";

const { createAPIKeyMock, deleteAPIKeyMock, listAPIKeysMock, AuthApiErrorMock } = vi.hoisted(
	() => {
		class MockAuthApiError extends Error {
			status: number;

			constructor(status: number, message: string) {
				super(message);
				this.name = "AuthApiError";
				this.status = status;
			}
		}

		return {
			createAPIKeyMock: vi.fn(),
			deleteAPIKeyMock: vi.fn(),
			listAPIKeysMock: vi.fn(),
			AuthApiErrorMock: MockAuthApiError
		};
	}
);

vi.mock("$lib/server/auth", () => ({
	createAPIKey: createAPIKeyMock,
	deleteAPIKey: deleteAPIKeyMock,
	listAPIKeys: listAPIKeysMock,
	AuthApiError: AuthApiErrorMock,
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock
}));

function apiKeyFixture() {
	return {
		id: "api-key-1",
		user_id: "user-1",
		name: "CLI",
		key_prefix: "abc12345",
		created_at: "2026-06-09T00:00:00Z",
		updated_at: "2026-06-09T00:00:00Z"
	};
}

function createEvent(
	method: "DELETE" | "GET" | "POST",
	options: {
		body?: unknown;
		url?: string;
		user?: unknown;
		cookie?: string;
	} = {}
) {
	const headers = new Headers();
	if (options.cookie) headers.set("cookie", options.cookie);
	if (options.body !== undefined) headers.set("content-type", "application/json");

	const request = new Request(options.url ?? "http://localhost/api/auth/apikeys", {
		method,
		headers,
		body: options.body === undefined ? undefined : JSON.stringify(options.body)
	});

	return {
		request,
		url: new URL(request.url),
		fetch: vi.fn(),
		locals: { user: options.user === undefined ? { id: "user-1" } : options.user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("API keys auth proxy endpoint", () => {
	beforeEach(() => {
		createAPIKeyMock.mockReset();
		deleteAPIKeyMock.mockReset();
		listAPIKeysMock.mockReset();
	});

	it("returns 401 for unauthenticated list requests", async () => {
		const response = await GET(createEvent("GET", { user: null }));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listAPIKeysMock).not.toHaveBeenCalled();
	});

	it("lists API keys for the current user only", async () => {
		const result = { api_keys: [apiKeyFixture()] };
		listAPIKeysMock.mockResolvedValue(result);
		const event = createEvent("GET", { cookie: "auth.session_token=token-1" });

		const response = await GET(event);

		expect(listAPIKeysMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=token-1", "user-1");
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("returns 400 for invalid create bodies", async () => {
		const response = await POST(createEvent("POST", { body: { name: "   " } }));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "name is required" });
		expect(createAPIKeyMock).not.toHaveBeenCalled();
	});

	it("returns 400 for invalid create expiration values", async () => {
		const response = await POST(createEvent("POST", { body: { name: "CLI", expires_at: "tomorrow" } }));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "expires_at must be RFC3339" });
		expect(createAPIKeyMock).not.toHaveBeenCalled();
	});

	it("creates API keys with the current user id and trimmed name", async () => {
		const result = { ...apiKeyFixture(), api_key: "abc12345secret" };
		createAPIKeyMock.mockResolvedValue(result);
		const event = createEvent("POST", {
			body: { name: "  CLI  " },
			cookie: "auth.session_token=token-1"
		});

		const response = await POST(event);

		expect(createAPIKeyMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=token-1", {
			userId: "user-1",
			name: "CLI",
			expiresAt: undefined
		});
		expect(response.status).toBe(201);
		expect(await responseJson(response)).toEqual(result);
	});

	it("forwards API key expiration when provided", async () => {
		const expiresAt = "2026-06-16T23:59:59.999Z";
		const result = { ...apiKeyFixture(), api_key: "abc12345secret", expires_at: expiresAt };
		createAPIKeyMock.mockResolvedValue(result);
		const event = createEvent("POST", {
			body: { name: "CLI", expires_at: expiresAt },
			cookie: "auth.session_token=token-1"
		});

		const response = await POST(event);

		expect(createAPIKeyMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=token-1", {
			userId: "user-1",
			name: "CLI",
			expiresAt
		});
		expect(response.status).toBe(201);
		expect(await responseJson(response)).toEqual(result);
	});

	it("treats null API key expiration as no expiration", async () => {
		const result = { ...apiKeyFixture(), api_key: "abc12345secret" };
		createAPIKeyMock.mockResolvedValue(result);
		const event = createEvent("POST", {
			body: { name: "CLI", expires_at: null },
			cookie: "auth.session_token=token-1"
		});

		const response = await POST(event);

		expect(createAPIKeyMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=token-1", {
			userId: "user-1",
			name: "CLI",
			expiresAt: undefined
		});
		expect(response.status).toBe(201);
		expect(await responseJson(response)).toEqual(result);
	});

	it("returns 400 when delete is missing api_key_id", async () => {
		const response = await DELETE(createEvent("DELETE"));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "api_key_id is required" });
		expect(deleteAPIKeyMock).not.toHaveBeenCalled();
	});

	it("deletes API keys with the current user id", async () => {
		deleteAPIKeyMock.mockResolvedValue({ deleted_id: "api-key-1", deleted_count: 1 });
		const event = createEvent("DELETE", {
			url: "http://localhost/api/auth/apikeys?api_key_id=api-key-1",
			cookie: "auth.session_token=token-1"
		});

		const response = await DELETE(event);

		expect(deleteAPIKeyMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=token-1", {
			userId: "user-1",
			apiKeyId: "api-key-1"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual({ deleted_id: "api-key-1", deleted_count: 1 });
	});

	it("normalizes backend auth service errors", async () => {
		listAPIKeysMock.mockRejectedValue(new AuthApiError(503, "Authentication service unavailable"));

		const response = await GET(createEvent("GET"));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
