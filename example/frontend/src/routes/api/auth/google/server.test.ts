import { beforeEach, describe, expect, it, vi } from "vitest";

const {
	AuthApiErrorMock,
	privateEnvMock,
	setGoogleOAuthStateCookieMock,
	startGoogleOAuthMock
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
		setGoogleOAuthStateCookieMock: vi.fn(),
		startGoogleOAuthMock: vi.fn()
	};
});

vi.mock("$lib/server/auth", () => ({
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock,
	setGoogleOAuthStateCookie: setGoogleOAuthStateCookieMock,
	startGoogleOAuth: startGoogleOAuthMock
}));

vi.mock("$lib/server/internal-api", () => ({
	privateEnv: privateEnvMock
}));

import { GET } from "./+server";

function event(path = "/api/auth/google", user: unknown = null) {
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

describe("Google OAuth start endpoint", () => {
	beforeEach(() => {
		privateEnvMock.mockReset();
		setGoogleOAuthStateCookieMock.mockReset();
		startGoogleOAuthMock.mockReset();
		startGoogleOAuthMock.mockResolvedValue({
			authorizationUrl: "https://accounts.google.com/o/oauth2/auth?state=state-1",
			state: "state-1",
			stateExpiresAt: "2026-06-13T12:00:00Z"
		});
	});

	it("starts Google OAuth, stores the state cookie, and redirects to Google", async () => {
		privateEnvMock.mockReturnValue("https://app.example.com/some/path");
		const { event: requestEvent, cookies, fetchMock } = event();

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location: "https://accounts.google.com/o/oauth2/auth?state=state-1"
		});

		expect(startGoogleOAuthMock).toHaveBeenCalledWith(
			fetchMock,
			"https://app.example.com/api/auth/callback/google"
		);
		expect(setGoogleOAuthStateCookieMock).toHaveBeenCalledWith(
			cookies,
			"state-1",
			"2026-06-13T12:00:00Z"
		);
	});

	it("redirects authenticated users to the app without starting OAuth", async () => {
		const { event: requestEvent } = event("/api/auth/google", { id: "user-1" });

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location: "/app"
		});

		expect(startGoogleOAuthMock).not.toHaveBeenCalled();
	});

	it("redirects to login with a configuration error when Google OAuth is disabled", async () => {
		startGoogleOAuthMock.mockRejectedValueOnce(
			new AuthApiErrorMock(503, "google oauth is not configured")
		);
		const { event: requestEvent } = event();

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location: "/login?oauth_error=configuration"
		});
	});
});
