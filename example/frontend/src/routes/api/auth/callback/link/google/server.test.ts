import { beforeEach, describe, expect, it, vi } from "vitest";

const {
	AuthApiErrorMock,
	clearGoogleOAuthLinkStateCookieMock,
	linkGoogleAccountMock,
	privateEnvMock
} = vi.hoisted(() => {
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
		clearGoogleOAuthLinkStateCookieMock: vi.fn(),
		linkGoogleAccountMock: vi.fn(),
		privateEnvMock: vi.fn()
	};
});

vi.mock("$lib/server/auth", () => ({
	GOOGLE_OAUTH_LINK_STATE_COOKIE_NAME: "auth.google_oauth_link_state",
	clearGoogleOAuthLinkStateCookie: clearGoogleOAuthLinkStateCookieMock,
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock,
	linkGoogleAccount: linkGoogleAccountMock
}));

vi.mock("$lib/server/internal-api", () => ({
	privateEnv: privateEnvMock
}));

import { GET } from "./+server";

function event(
	path = "/api/auth/callback/link/google?code=code-1&state=state-1",
	stateCookie = "state-1",
	user: unknown = { id: "user-1" }
) {
	const url = new URL(`http://localhost:5173${path}`);
	const headers = new Headers({ cookie: "auth.session_token=token-1" });
	const request = new Request(url, { headers });
	const fetchMock = vi.fn();
	const cookies = { get: vi.fn(() => stateCookie), delete: vi.fn() };
	return {
		event: {
			url,
			request,
			fetch: fetchMock,
			cookies,
			locals: { user }
		},
		cookies,
		fetchMock
	};
}

describe("Google account link callback endpoint", () => {
	beforeEach(() => {
		clearGoogleOAuthLinkStateCookieMock.mockReset();
		linkGoogleAccountMock.mockReset();
		privateEnvMock.mockReset();
		linkGoogleAccountMock.mockResolvedValue({
			id: "account-1",
			providerId: "google",
			createdAt: "2026-06-13T12:00:00Z",
			updatedAt: "2026-06-13T12:00:00Z"
		});
	});

	it("validates state, links through the backend, clears state, and redirects", async () => {
		privateEnvMock.mockReturnValue("https://app.example.com/ignored");
		const { event: requestEvent, cookies, fetchMock } = event();

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location:
				"/app?account_settings=linked&account_link_provider=google&account_link_status=linked"
		});

		expect(linkGoogleAccountMock).toHaveBeenCalledWith(fetchMock, "auth.session_token=token-1", {
			code: "code-1",
			state: "state-1",
			redirectURI: "https://app.example.com/api/auth/callback/link/google"
		});
		expect(clearGoogleOAuthLinkStateCookieMock).toHaveBeenCalledWith(cookies);
	});

	it("rejects a mismatched state before calling the backend", async () => {
		const { event: requestEvent, cookies } = event(
			"/api/auth/callback/link/google?code=code-1&state=state-1",
			"different-state"
		);

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location:
				"/app?account_settings=linked&account_link_provider=google&account_link_status=invalid"
		});

		expect(linkGoogleAccountMock).not.toHaveBeenCalled();
		expect(clearGoogleOAuthLinkStateCookieMock).toHaveBeenCalledWith(cookies);
	});

	it("maps provider access denial to a denied linked-account status", async () => {
		const { event: requestEvent, cookies } = event(
			"/api/auth/callback/link/google?error=access_denied"
		);

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location:
				"/app?account_settings=linked&account_link_provider=google&account_link_status=denied"
		});
		expect(clearGoogleOAuthLinkStateCookieMock).toHaveBeenCalledWith(cookies);
	});

	it("maps provider conflicts to a conflict linked-account status", async () => {
		linkGoogleAccountMock.mockRejectedValueOnce(
			new AuthApiErrorMock(409, "google account is already linked to another user")
		);
		const { event: requestEvent, cookies } = event();

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location:
				"/app?account_settings=linked&account_link_provider=google&account_link_status=conflict"
		});
		expect(clearGoogleOAuthLinkStateCookieMock).toHaveBeenCalledWith(cookies);
	});

	it("clears link state and redirects unauthenticated callbacks to login", async () => {
		const { event: requestEvent, cookies } = event(
			"/api/auth/callback/link/google?code=code-1&state=state-1",
			"state-1",
			null
		);

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location: "/login?oauth_error=invalid"
		});
		expect(linkGoogleAccountMock).not.toHaveBeenCalled();
		expect(clearGoogleOAuthLinkStateCookieMock).toHaveBeenCalledWith(cookies);
	});
});
