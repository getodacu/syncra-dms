import { beforeEach, describe, expect, it, vi } from "vitest";

import { AuthApiError } from "$lib/server/auth";
import { PATCH } from "./+server";
import type { RequestEvent } from "./$types";

const { updateAuthUserMock, setPreferredLanguageCookieMock, AuthApiErrorMock } = vi.hoisted(() => {
	class MockAuthApiError extends Error {
		status: number;
		constructor(status: number, message: string) {
			super(message);
			this.name = "AuthApiError";
			this.status = status;
		}
	}
	return {
		updateAuthUserMock: vi.fn(),
		setPreferredLanguageCookieMock: vi.fn(),
		AuthApiErrorMock: MockAuthApiError
	};
});

vi.mock("$lib/server/auth", () => ({
	updateAuthUser: updateAuthUserMock,
	setPreferredLanguageCookie: setPreferredLanguageCookieMock,
	isSupportedPreferredLanguage: (value: unknown) => value === "en" || value === "ro",
	AuthApiError: AuthApiErrorMock,
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock
}));

function event(body: unknown, user: unknown = { id: "user-1" }, cookie = "auth.session_token=token-1") {
	const cookies = { set: vi.fn() };
	return {
		request: new Request("http://localhost/api/auth/user", {
			method: "PATCH",
			headers: { cookie },
			body: JSON.stringify(body)
		}),
		fetch: vi.fn(),
		locals: { user },
		cookies
	} as unknown as RequestEvent;
}

function rawEvent(body: string, user: unknown = { id: "user-1" }, cookie = "auth.session_token=token-1") {
	const cookies = { set: vi.fn() };
	return {
		request: new Request("http://localhost/api/auth/user", {
			method: "PATCH",
			headers: { cookie },
			body
		}),
		fetch: vi.fn(),
		locals: { user },
		cookies
	} as unknown as RequestEvent;
}

function requestEvent(request: Request, user: unknown = { id: "user-1" }) {
	const cookies = { set: vi.fn() };
	return {
		request,
		fetch: vi.fn(),
		locals: { user },
		cookies
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("auth user API endpoint", () => {
	beforeEach(() => {
		updateAuthUserMock.mockReset();
		setPreferredLanguageCookieMock.mockReset();
	});

	it("returns 401 without an authenticated user", async () => {
		const response = await PATCH(event({ name: "Ada" }, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({
			error: "authentication required"
		});
		expect(updateAuthUserMock).not.toHaveBeenCalled();
	});

	it("returns 400 when the request body exceeds the account update limit", async () => {
		const response = await PATCH(
			requestEvent(
				new Request("http://localhost/api/auth/user", {
					method: "PATCH",
					headers: {
						cookie: "auth.session_token=token-1",
						"content-length": "99999999"
					},
					body: "{}"
				})
			)
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "request body too large"
		});
		expect(updateAuthUserMock).not.toHaveBeenCalled();
	});

	it("returns 400 for malformed JSON", async () => {
		const response = await PATCH(rawEvent("{"));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "invalid JSON body"
		});
		expect(updateAuthUserMock).not.toHaveBeenCalled();
	});

	it("returns 400 for non-object JSON", async () => {
		const response = await PATCH(event(null));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "invalid user update payload"
		});
		expect(updateAuthUserMock).not.toHaveBeenCalled();
	});

	it("forwards patch payload and cookies", async () => {
		const updated = {
			id: "user-1",
			name: "Ada",
			email: "ada@example.com",
			emailVerified: true,
			image: null,
			preferredLanguage: "ro",
			createdAt: "2026-05-26T00:00:00Z",
			updatedAt: "2026-05-30T00:00:00Z"
		};
		updateAuthUserMock.mockResolvedValue(updated);

		const requestEvent = event({ name: "Ada", preferredLanguage: "ro" });
		const response = await PATCH(requestEvent);

		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(updated);
		expect(updateAuthUserMock).toHaveBeenCalledWith(expect.any(Function), "auth.session_token=token-1", {
			name: "Ada",
			preferredLanguage: "ro"
		});
		expect(setPreferredLanguageCookieMock).toHaveBeenCalledWith(
			(requestEvent as unknown as { cookies: unknown }).cookies,
			"ro"
		);
	});

	it("maps backend auth errors", async () => {
		const error = new AuthApiError(409, "email is already in use");
		updateAuthUserMock.mockImplementationOnce(async () => {
			throw error;
		});

		const response = await PATCH(event({ email: "taken@example.com" }));

		expect(response.status).toBe(409);
		expect(await responseJson(response)).toEqual({
			error: "email is already in use"
		});
	});

	it("rethrows non-auth errors", async () => {
		const error = new Error("unexpected failure");
		updateAuthUserMock.mockImplementationOnce(async () => {
			throw error;
		});

		await expect(PATCH(event({ name: "Ada" }))).rejects.toBe(error);
	});
});
