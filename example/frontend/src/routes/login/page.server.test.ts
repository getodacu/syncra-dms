import { beforeEach, describe, expect, it, vi } from "vitest";

const { isAuthApiErrorMock, setPreferredLanguageCookieMock, setSessionCookieMock, signInEmailMock } = vi.hoisted(
	() => ({
		isAuthApiErrorMock: vi.fn(),
		setPreferredLanguageCookieMock: vi.fn(),
		setSessionCookieMock: vi.fn(),
		signInEmailMock: vi.fn()
	})
);

vi.mock("$lib/server/auth", () => ({
	isAuthApiError: isAuthApiErrorMock,
	setPreferredLanguageCookie: setPreferredLanguageCookieMock,
	setSessionCookie: setSessionCookieMock,
	signInEmail: signInEmailMock
}));

import { actions, load } from "./+page.server";

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

describe("login page load", () => {
	beforeEach(() => {
		isAuthApiErrorMock.mockReset();
		setPreferredLanguageCookieMock.mockReset();
		setSessionCookieMock.mockReset();
		signInEmailMock.mockReset();
		signInEmailMock.mockResolvedValue(authPayload);
	});

	it("returns reset success state from the query string", () => {
		const result = load({
			locals: {},
			url: new URL("http://localhost/login?email=ADA%40EXAMPLE.COM&reset=1")
		} as never);

		expect(result).toEqual({
			email: "ada@example.com",
			verified: false,
			reset: true,
			oauthError: ""
		});
	});

	it("returns provider-neutral OAuth errors from the query string", () => {
		const result = load({
			locals: {},
			url: new URL("http://localhost/login?oauth_error=failed")
		} as never);

		expect(result).toMatchObject({
			oauthError: "Social login failed. Please try again."
		});
	});

	it("signs in email users with a browser-session cookie when remember me is absent", async () => {
		const formData = new FormData();
		formData.set("email", " ADA@EXAMPLE.COM ");
		formData.set("password", "password1234");
		const request = new Request("http://localhost/login", {
			method: "POST",
			body: formData
		});
		const cookies = {};
		const fetchMock = vi.fn();
		const locals = {};

		await expect(
			actions.default({
				cookies,
				fetch: fetchMock,
				locals,
				request
			} as never)
		).rejects.toMatchObject({
			status: 303,
			location: "/app"
		});

		expect(signInEmailMock).toHaveBeenCalledWith(fetchMock, {
			email: "ada@example.com",
			password: "password1234",
			rememberMe: false
		});
		expect(setSessionCookieMock).toHaveBeenCalledWith(cookies, authPayload.session, false);
		expect(setPreferredLanguageCookieMock).toHaveBeenCalledWith(cookies, "ro");
		expect(locals).toEqual({
			session: authPayload.session,
			user: authPayload.user
		});
	});

	it("returns generic form errors for auth service failures", async () => {
		const formData = new FormData();
		formData.set("email", "ada@example.com");
		formData.set("password", "password1234");
		const request = new Request("http://localhost/login", {
			method: "POST",
			body: formData
		});
		const authError = new Error("database connection failed") as Error & {
			status: number;
		};
		authError.status = 503;
		signInEmailMock.mockRejectedValue(authError);
		isAuthApiErrorMock.mockReturnValue(true);

		const result = await actions.default({
			cookies: {},
			fetch: vi.fn(),
			locals: {},
			request
		} as never);

		expect(result).toEqual({
			status: 502,
			data: {
				values: { email: "ada@example.com" },
				error: "Unable to sign in. Please try again."
			}
		});
	});
});
