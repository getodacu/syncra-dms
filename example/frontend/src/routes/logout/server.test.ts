import { beforeEach, describe, expect, it, vi } from "vitest";

const { AuthApiErrorMock, clearSessionCookieMock, loggerMock, signOutMock } = vi.hoisted(() => {
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
		clearSessionCookieMock: vi.fn(),
		loggerMock: {
			info: vi.fn(),
			error: vi.fn(),
			warn: vi.fn(),
			debug: vi.fn(),
			child: vi.fn()
		},
		signOutMock: vi.fn()
	};
});

vi.mock("$lib/server/auth", () => ({
	clearSessionCookie: clearSessionCookieMock,
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock,
	signOut: signOutMock
}));

import { POST } from "./+server";

function logoutEvent(cookie = "auth.session_token=token-1") {
	const request = new Request("http://localhost/logout", {
		method: "POST",
		headers: { cookie }
	});
	const fetchMock = vi.fn();
	const cookies = { delete: vi.fn() };
	const locals = {
		session: { id: "session-1" },
		user: { id: "user-1" },
		logger: loggerMock
	};

	return {
		event: {
			request,
			fetch: fetchMock,
			cookies,
			locals
		},
		cookies,
		fetchMock,
		locals
	};
}

describe("logout endpoint", () => {
	beforeEach(() => {
		clearSessionCookieMock.mockReset();
		loggerMock.error.mockReset();
		signOutMock.mockReset();
		signOutMock.mockResolvedValue({ success: true });
	});

	it("signs out with the backend, clears the local session, and redirects to login", async () => {
		const { event, cookies, fetchMock, locals } = logoutEvent();

		await expect(POST(event as never)).rejects.toMatchObject({
			status: 303,
			location: "/login"
		});

		expect(signOutMock).toHaveBeenCalledWith(fetchMock, "auth.session_token=token-1");
		expect(clearSessionCookieMock).toHaveBeenCalledWith(cookies);
		expect(locals).toMatchObject({
			session: null,
			user: null
		});
	});

	it("still clears the local session when the backend sign-out request fails", async () => {
		const backendError = new AuthApiErrorMock(503, "Authentication service unavailable");
		signOutMock.mockRejectedValueOnce(backendError);
		const { event, cookies, locals } = logoutEvent();

		await expect(POST(event as never)).rejects.toMatchObject({
			status: 303,
			location: "/login"
		});

		expect(clearSessionCookieMock).toHaveBeenCalledWith(cookies);
		expect(locals).toMatchObject({
			session: null,
			user: null
		});
		expect(loggerMock.error).toHaveBeenCalledWith("auth.signout_failed", {
			error: "Authentication service unavailable",
			user_id: "user-1"
		});
	});
});
