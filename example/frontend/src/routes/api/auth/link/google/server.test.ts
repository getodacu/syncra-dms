import { beforeEach, describe, expect, it, vi } from "vitest";

const {
	AuthApiErrorMock,
	privateEnvMock,
	setGoogleOAuthLinkStateCookieMock,
	startGoogleAccountLinkMock
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
		setGoogleOAuthLinkStateCookieMock: vi.fn(),
		startGoogleAccountLinkMock: vi.fn()
	};
});

vi.mock("$lib/server/auth", () => ({
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock,
	setGoogleOAuthLinkStateCookie: setGoogleOAuthLinkStateCookieMock,
	startGoogleAccountLink: startGoogleAccountLinkMock
}));

vi.mock("$lib/server/internal-api", () => ({
	privateEnv: privateEnvMock
}));

import { GET } from "./+server";

function event(path = "/api/auth/link/google", user: unknown = { id: "user-1" }) {
	const url = new URL(`http://localhost:5173${path}`);
	const headers = new Headers({ cookie: "auth.session_token=token-1" });
	const request = new Request(url, { headers });
	const fetchMock = vi.fn();
	const cookies = { set: vi.fn() };
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

describe("Google account link start endpoint", () => {
	beforeEach(() => {
		privateEnvMock.mockReset();
		setGoogleOAuthLinkStateCookieMock.mockReset();
		startGoogleAccountLinkMock.mockReset();
		startGoogleAccountLinkMock.mockResolvedValue({
			authorizationUrl: "https://accounts.google.com/o/oauth2/auth?state=state-1",
			state: "state-1",
			stateExpiresAt: "2026-06-13T12:00:00Z"
		});
	});

	it("requires an authenticated user", async () => {
		const { event: requestEvent } = event("/api/auth/link/google", null);

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location: "/login"
		});
		expect(startGoogleAccountLinkMock).not.toHaveBeenCalled();
	});

	it("starts Google linking, stores link state, and redirects to Google", async () => {
		privateEnvMock.mockReturnValue("https://app.example.com/some/path");
		const { event: requestEvent, cookies, fetchMock } = event();

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location: "https://accounts.google.com/o/oauth2/auth?state=state-1"
		});

		expect(startGoogleAccountLinkMock).toHaveBeenCalledWith(
			fetchMock,
			"auth.session_token=token-1",
			"https://app.example.com/api/auth/callback/link/google"
		);
		expect(setGoogleOAuthLinkStateCookieMock).toHaveBeenCalledWith(
			cookies,
			"state-1",
			"2026-06-13T12:00:00Z"
		);
	});

	it("redirects configuration failures back to linked accounts", async () => {
		startGoogleAccountLinkMock.mockRejectedValueOnce(
			new AuthApiErrorMock(503, "google oauth is not configured")
		);
		const { event: requestEvent } = event();

		await expect(GET(requestEvent as never)).rejects.toMatchObject({
			status: 303,
			location:
				"/app?account_settings=linked&account_link_provider=google&account_link_status=configuration"
		});
	});
});
