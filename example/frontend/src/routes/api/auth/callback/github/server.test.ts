import { beforeEach, describe, expect, it, vi } from "vitest";

const {
	AuthApiErrorMock,
	clearGitHubOAuthStateCookieMock,
	privateEnvMock,
	setPreferredLanguageCookieMock,
	setSessionCookieMock,
	signInGitHubOAuthMock
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
		clearGitHubOAuthStateCookieMock: vi.fn(),
		privateEnvMock: vi.fn(),
		setPreferredLanguageCookieMock: vi.fn(),
		setSessionCookieMock: vi.fn(),
		signInGitHubOAuthMock: vi.fn()
	};
});

vi.mock("$lib/server/auth", () => ({
	GITHUB_OAUTH_STATE_COOKIE_NAME: "auth.github_oauth_state",
	clearGitHubOAuthStateCookie: clearGitHubOAuthStateCookieMock,
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock,
	setPreferredLanguageCookie: setPreferredLanguageCookieMock,
	setSessionCookie: setSessionCookieMock,
	signInGitHubOAuth: signInGitHubOAuthMock
}));

vi.mock("$lib/server/internal-api", () => ({
	privateEnv: privateEnvMock
}));

import { GET } from "./+server";

const authPayload = {
	session: {
		id: "session-1",
		token: "token-1",
		userId: "user-1",
		expiresAt: "2026-06-20T12:00:00Z",
		createdAt: "2026-06-13T12:00:00Z",
		updatedAt: "2026-06-13T12:00:00Z"
	},
	user: {
		id: "user-1",
		name: "Ada",
		email: "ada@example.com",
		emailVerified: true,
		image: null,
		role: "user",
		lastLoginAt: "2026-06-13T12:00:00Z",
		preferredLanguage: "ro",
		createdAt: "2026-06-13T12:00:00Z",
		updatedAt: "2026-06-13T12:00:00Z"
	}
};

function event(path = "/api/auth/callback/github?code=code-1&state=state-1", stateCookie = "state-1") {
	const url = new URL(`http://localhost:5173${path}`);
	const fetchMock = vi.fn();
	const cookies = {
		get: vi.fn(() => stateCookie),
		set: vi.fn(),
		delete: vi.fn()
	};
	const locals = { session: null, user: null };
	return {
		event: {
			url,
			fetch: fetchMock,
			cookies,
			locals
		},
		cookies,
		fetchMock,
		locals
	};
}

describe("GitHub OAuth callback endpoint", () => {
	beforeEach(() => {
		clearGitHubOAuthStateCookieMock.mockReset();
		privateEnvMock.mockReset();
		setPreferredLanguageCookieMock.mockReset();
		setSessionCookieMock.mockReset();
		signInGitHubOAuthMock.mockReset();
		signInGitHubOAuthMock.mockResolvedValue(authPayload);
	});

	it("validates state, signs in with the backend, sets the session cookie, and redirects", async () => {
		privateEnvMock.mockReturnValue("https://app.example.com/ignored");
		const { event: requestEvent, cookies, fetchMock, locals } = event();

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location: "/app"
		});

		expect(signInGitHubOAuthMock).toHaveBeenCalledWith(fetchMock, {
			code: "code-1",
			state: "state-1",
			redirectURI: "https://app.example.com/api/auth/callback/github"
		});
		expect(setSessionCookieMock).toHaveBeenCalledWith(cookies, authPayload.session, true);
		expect(setPreferredLanguageCookieMock).toHaveBeenCalledWith(cookies, "ro");
		expect(clearGitHubOAuthStateCookieMock).toHaveBeenCalledWith(cookies);
		expect(locals).toEqual({
			session: authPayload.session,
			user: authPayload.user
		});
	});

	it("rejects a mismatched state before calling the backend", async () => {
		const { event: requestEvent, cookies } = event(
			"/api/auth/callback/github?code=code-1&state=state-1",
			"different-state"
		);

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location: "/login?oauth_error=invalid"
		});

		expect(signInGitHubOAuthMock).not.toHaveBeenCalled();
		expect(clearGitHubOAuthStateCookieMock).toHaveBeenCalledWith(cookies);
	});

	it("maps GitHub access denial to a cancelled login message", async () => {
		const { event: requestEvent, cookies } = event("/api/auth/callback/github?error=access_denied", "state-1");

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location: "/login?oauth_error=denied"
		});

		expect(clearGitHubOAuthStateCookieMock).toHaveBeenCalledWith(cookies);
	});

	it("maps backend configuration failures to the login configuration message", async () => {
		signInGitHubOAuthMock.mockRejectedValueOnce(new AuthApiErrorMock(503, "github oauth is not configured"));
		const { event: requestEvent, cookies } = event();

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location: "/login?oauth_error=configuration"
		});

		expect(clearGitHubOAuthStateCookieMock).toHaveBeenCalledWith(cookies);
	});
});
