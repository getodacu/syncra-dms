import { beforeEach, describe, expect, it, vi } from "vitest";

const {
	AuthApiErrorMock,
	privateEnvMock,
	setGitHubOAuthStateCookieMock,
	startGitHubOAuthMock
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
		privateEnvMock: vi.fn(),
		setGitHubOAuthStateCookieMock: vi.fn(),
		startGitHubOAuthMock: vi.fn()
	};
});

vi.mock("$lib/server/auth", () => ({
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock,
	setGitHubOAuthStateCookie: setGitHubOAuthStateCookieMock,
	startGitHubOAuth: startGitHubOAuthMock
}));

vi.mock("$lib/server/internal-api", () => ({
	privateEnv: privateEnvMock
}));

import { GET } from "./+server";

function event(path = "/api/auth/github", user: unknown = null) {
	const url = new URL(`http://localhost:5173${path}`);
	const fetchMock = vi.fn();
	const cookies = { set: vi.fn() };
	return {
		event: {
			url,
			fetch: fetchMock,
			cookies,
			locals: { user }
		},
		cookies,
		fetchMock
	};
}

describe("GitHub OAuth start endpoint", () => {
	beforeEach(() => {
		privateEnvMock.mockReset();
		setGitHubOAuthStateCookieMock.mockReset();
		startGitHubOAuthMock.mockReset();
		startGitHubOAuthMock.mockResolvedValue({
			authorizationUrl: "https://github.com/login/oauth/authorize?state=state-1",
			state: "state-1",
			stateExpiresAt: "2026-06-13T12:00:00Z"
		});
	});

	it("starts GitHub OAuth, stores the state cookie, and redirects to GitHub", async () => {
		privateEnvMock.mockReturnValue("https://app.example.com/some/path");
		const { event: requestEvent, cookies, fetchMock } = event();

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location: "https://github.com/login/oauth/authorize?state=state-1"
		});

		expect(startGitHubOAuthMock).toHaveBeenCalledWith(
			fetchMock,
			"https://app.example.com/api/auth/callback/github"
		);
		expect(setGitHubOAuthStateCookieMock).toHaveBeenCalledWith(
			cookies,
			"state-1",
			"2026-06-13T12:00:00Z"
		);
	});

	it("redirects authenticated users to the app without starting OAuth", async () => {
		const { event: requestEvent } = event("/api/auth/github", { id: "user-1" });

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location: "/app"
		});

		expect(startGitHubOAuthMock).not.toHaveBeenCalled();
	});

	it("redirects to login with a configuration error when GitHub OAuth is disabled", async () => {
		startGitHubOAuthMock.mockRejectedValueOnce(
			new AuthApiErrorMock(503, "github oauth is not configured")
		);
		const { event: requestEvent } = event();

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location: "/login?oauth_error=configuration"
		});
	});
});
