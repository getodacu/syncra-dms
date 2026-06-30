import { beforeEach, describe, expect, it, vi } from "vitest";

type FetchInput = Parameters<typeof fetch>[0];
type FetchInit = Parameters<typeof fetch>[1];

const INTERNAL_API_HEADER = "X-Syncra-Internal-Token";

const { privateEnv } = vi.hoisted(() => ({
	privateEnv: {} as Record<string, string | undefined>
}));

vi.mock("$env/dynamic/private", () => ({ env: privateEnv }));

describe("frontend auth server helper", () => {
	beforeEach(() => {
		vi.resetModules();
		for (const key of Object.keys(privateEnv)) delete privateEnv[key];
		process.env.SYNCRA_API_BASE_URL = "http://auth-api.test/";
		process.env.SYNCRA_INTERNAL_API_TOKEN = "internal-token";
		process.env.AUTH_DELIVERY_TOKEN = "trusted-delivery-token";
		process.env.NODE_ENV = "test";
		delete process.env.AUTH_COOKIE_SECURE;
	});

	it("pins the frontend session cookie name contract", async () => {
		const { AUTH_SESSION_COOKIE_NAME } = await import("./auth");

		expect(AUTH_SESSION_COOKIE_NAME).toBe("auth.session_token");
	});

	it("signs up through the backend with the trusted delivery header", async () => {
		const { signUpEmail } = await import("./auth");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(
				JSON.stringify({
					user: {
						id: "user-1",
						name: "Ada Lovelace",
						email: "ada@example.com",
						emailVerified: false,
						image: null,
						createdAt: "2026-05-26T00:00:00Z",
						updatedAt: "2026-05-26T00:00:00Z"
					},
					verificationRequired: true,
					verificationCode: "123456"
				})
			);
		});

		const result = await signUpEmail(fetchMock, {
			name: "Ada Lovelace",
			email: "ada@example.com",
			password: "password123"
		});

		expect(result.verificationCode).toBe("123456");
		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/sign-up/email",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({
					name: "Ada Lovelace",
					email: "ada@example.com",
					password: "password123"
				})
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
		expect(headers.get("X-Syncra-Auth-Delivery-Token")).toBe("trusted-delivery-token");
	});

	it("requests password resets through the backend with the trusted delivery header", async () => {
		const { requestPasswordReset } = await import("./auth");
		const response = {
			ok: true,
			resetToken: "reset-token",
			resetExpiresAt: "2026-06-11T13:30:00Z"
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(response));
		});

		await expect(requestPasswordReset(fetchMock, "ada@example.com")).resolves.toEqual(response);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/password-reset/request",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({ email: "ada@example.com" })
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
		expect(headers.get("X-Syncra-Auth-Delivery-Token")).toBe("trusted-delivery-token");
	});

	it("confirms password resets through the backend", async () => {
		const { resetPassword } = await import("./auth");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ ok: true }));
		});

		await expect(
			resetPassword(fetchMock, {
				email: "ada@example.com",
				token: "reset-token",
				password: "newpassword123"
			})
		).resolves.toEqual({ ok: true });

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/password-reset/confirm",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({
					email: "ada@example.com",
					token: "reset-token",
					password: "newpassword123"
				})
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
		expect(headers.get("X-Syncra-Auth-Delivery-Token")).toBeNull();
	});

	it("forwards cookies when loading the current session", async () => {
		const { AUTH_SESSION_COOKIE_NAME, getSession } = await import("./auth");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => new Response("null"));

		await getSession(fetchMock, `${AUTH_SESSION_COOKIE_NAME}=token-1`);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/get-session",
			expect.objectContaining({ method: "GET" })
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
		expect(headers.get("cookie")).toBe(`${AUTH_SESSION_COOKIE_NAME}=token-1`);
	});

	it("patches the current user through the backend with cookies", async () => {
		const { AUTH_SESSION_COOKIE_NAME, updateAuthUser } = await import("./auth");
		const user = {
			id: "user-1",
			name: "Ada Updated",
			email: "ada@example.com",
			emailVerified: true,
			image: "data:image/png;base64,iVBORw0KGgo=",
			preferredLanguage: "ro",
			createdAt: "2026-05-26T00:00:00Z",
			updatedAt: "2026-05-30T00:00:00Z"
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(user));
		});

		await expect(
			updateAuthUser(fetchMock, `${AUTH_SESSION_COOKIE_NAME}=token-1`, {
				name: "Ada Updated",
				preferredLanguage: "ro"
			})
		).resolves.toEqual(user);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/user",
			expect.objectContaining({
				method: "PATCH",
				body: JSON.stringify({ name: "Ada Updated", preferredLanguage: "ro" })
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
		expect(headers.get("cookie")).toBe(`${AUTH_SESSION_COOKIE_NAME}=token-1`);
	});

	it("sets the Paraglide locale cookie for supported preferred languages", async () => {
		const { setPreferredLanguageCookie } = await import("./auth");
		const cookies = { set: vi.fn() };

		expect(setPreferredLanguageCookie(cookies as never, "ro")).toBe(true);

		expect(cookies.set).toHaveBeenCalledWith("PARAGLIDE_LOCALE", "ro", {
			httpOnly: false,
			path: "/",
			maxAge: 34560000,
			sameSite: "lax"
		});
	});

	it("lists API keys through the backend with cookies", async () => {
		const { AUTH_SESSION_COOKIE_NAME, listAPIKeys } = await import("./auth");
		const apiKeys = {
			api_keys: [
				{
					id: "api-key-1",
					user_id: "user-1",
					name: "CLI",
					key_prefix: "abc12345",
					created_at: "2026-06-09T00:00:00Z",
					updated_at: "2026-06-09T00:00:00Z"
				}
			]
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(apiKeys));
		});

		await expect(
			listAPIKeys(fetchMock, `${AUTH_SESSION_COOKIE_NAME}=token-1`, "user-1")
		).resolves.toEqual(apiKeys);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/apikeys/user-1",
			expect.objectContaining({ method: "GET" })
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
		expect(headers.get("cookie")).toBe(`${AUTH_SESSION_COOKIE_NAME}=token-1`);
	});

	it("creates API keys through the backend with cookies and returns the one-time secret", async () => {
		const { AUTH_SESSION_COOKIE_NAME, createAPIKey } = await import("./auth");
		const apiKey = {
			id: "api-key-1",
			user_id: "user-1",
			name: "CLI",
			key_prefix: "abc12345",
			api_key: "abc12345secret",
			created_at: "2026-06-09T00:00:00Z",
			updated_at: "2026-06-09T00:00:00Z"
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(apiKey), { status: 201 });
		});

		await expect(
			createAPIKey(fetchMock, `${AUTH_SESSION_COOKIE_NAME}=token-1`, {
				userId: "user-1",
				name: "CLI"
			})
		).resolves.toEqual(apiKey);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/apikeys",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({ user_id: "user-1", name: "CLI" })
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
		expect(headers.get("cookie")).toBe(`${AUTH_SESSION_COOKIE_NAME}=token-1`);
	});

	it("passes API key expiration to the backend when provided", async () => {
		const { AUTH_SESSION_COOKIE_NAME, createAPIKey } = await import("./auth");
		const expiresAt = "2026-06-16T23:59:59.999Z";
		const apiKey = {
			id: "api-key-1",
			user_id: "user-1",
			name: "CLI",
			key_prefix: "abc12345",
			api_key: "abc12345secret",
			expires_at: expiresAt,
			created_at: "2026-06-09T00:00:00Z",
			updated_at: "2026-06-09T00:00:00Z"
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(apiKey), { status: 201 });
		});

		await expect(
			createAPIKey(fetchMock, `${AUTH_SESSION_COOKIE_NAME}=token-1`, {
				userId: "user-1",
				name: "CLI",
				expiresAt
			})
		).resolves.toEqual(apiKey);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/apikeys",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({ user_id: "user-1", name: "CLI", expires_at: expiresAt })
			})
		);
	});

	it("deletes API keys through the backend with cookies", async () => {
		const { AUTH_SESSION_COOKIE_NAME, deleteAPIKey } = await import("./auth");
		const result = { deleted_id: "api-key/1", deleted_count: 1 };
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(result));
		});

		await expect(
			deleteAPIKey(fetchMock, `${AUTH_SESSION_COOKIE_NAME}=token-1`, {
				userId: "user-1",
				apiKeyId: "api-key/1"
			})
		).resolves.toEqual(result);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/apikeys?user_id=user-1&api_key_id=api-key%2F1",
			expect.objectContaining({ method: "DELETE" })
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
		expect(headers.get("cookie")).toBe(`${AUTH_SESSION_COOKIE_NAME}=token-1`);
	});

	it("gets webhooks through the backend with cookies", async () => {
		const { AUTH_SESSION_COOKIE_NAME, getWebhook } = await import("./auth");
		const webhook = {
			webhook: {
				id: "webhook-1",
				user_id: "user-1",
				url: "https://example.com/webhook",
				events_active: ["job.started", "job.succeeded"],
				has_secret: true,
				secret_key: "whsec_once",
				created_at: "2026-06-09T00:00:00Z",
				updated_at: "2026-06-09T00:00:00Z"
			}
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(webhook));
		});

		await expect(
			getWebhook(fetchMock, `${AUTH_SESSION_COOKIE_NAME}=token-1`, "user-1")
		).resolves.toEqual(webhook);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/webhook/user-1",
			expect.objectContaining({ method: "GET" })
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
		expect(headers.get("cookie")).toBe(`${AUTH_SESSION_COOKIE_NAME}=token-1`);
	});

	it("saves webhooks through the backend with cookies", async () => {
		const { AUTH_SESSION_COOKIE_NAME, saveWebhook } = await import("./auth");
		const webhook = {
			id: "webhook-1",
			user_id: "user-1",
			url: "https://example.com/webhook",
			events_active: [],
			has_secret: true,
			created_at: "2026-06-09T00:00:00Z",
			updated_at: "2026-06-09T00:00:00Z"
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(webhook));
		});

		await expect(
			saveWebhook(fetchMock, `${AUTH_SESSION_COOKIE_NAME}=token-1`, {
				userId: "user-1",
				url: "https://example.com/webhook",
				eventsActive: []
			})
		).resolves.toEqual(webhook);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/webhook",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({
					user_id: "user-1",
					url: "https://example.com/webhook",
					events_active: []
				})
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
		expect(headers.get("cookie")).toBe(`${AUTH_SESSION_COOKIE_NAME}=token-1`);
	});

	it("lists sessions through the backend without accepting token leakage", async () => {
		const { AUTH_SESSION_COOKIE_NAME, listAuthSessions } = await import("./auth");
		const sessions = {
			sessions: [
				{
					id: "session-1",
					userId: "user-1",
					expiresAt: "2026-06-20T12:00:00Z",
					ipAddress: "203.0.113.1",
					userAgent: "Browser",
					createdAt: "2026-06-13T12:00:00Z",
					updatedAt: "2026-06-13T12:00:00Z",
					current: true
				}
			]
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(sessions));
		});

		await expect(
			listAuthSessions(fetchMock, `${AUTH_SESSION_COOKIE_NAME}=token-1`)
		).resolves.toEqual(sessions);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/sessions",
			expect.objectContaining({ method: "GET" })
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("cookie")).toBe(`${AUTH_SESSION_COOKIE_NAME}=token-1`);

		fetchMock.mockResolvedValueOnce(
			new Response(JSON.stringify({ sessions: [{ ...sessions.sessions[0], token: "leaked" }] }))
		);
		await expect(listAuthSessions(fetchMock, "auth.session_token=token-1")).rejects.toMatchObject({
			status: 502,
			message: "Invalid session list response"
		});
	});

	it("revokes sessions through the backend", async () => {
		const { revokeAuthSession } = await import("./auth");
		const result = { deleted_id: "session/1", deleted_count: 1 };
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(result));
		});

		await expect(
			revokeAuthSession(fetchMock, "auth.session_token=token-1", "session/1")
		).resolves.toEqual(result);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/sessions/session%2F1",
			expect.objectContaining({ method: "DELETE" })
		);
	});

	it("lists linked accounts through the backend without accepting sensitive fields", async () => {
		const { listAuthAccounts } = await import("./auth");
		const accounts = {
			accounts: [
				{
					id: "account-1",
					providerId: "google",
					createdAt: "2026-06-13T12:00:00Z",
					updatedAt: "2026-06-13T12:00:00Z"
				}
			]
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(accounts));
		});

		await expect(listAuthAccounts(fetchMock, "auth.session_token=token-1")).resolves.toEqual(
			accounts
		);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/accounts",
			expect.objectContaining({ method: "GET" })
		);

		fetchMock.mockResolvedValueOnce(
			new Response(JSON.stringify({ accounts: [{ ...accounts.accounts[0], accessToken: "secret" }] }))
		);
		await expect(listAuthAccounts(fetchMock, "auth.session_token=token-1")).rejects.toMatchObject({
			status: 502,
			message: "Invalid linked account list response"
		});
	});

	it("unlinks linked accounts through the backend", async () => {
		const { unlinkAuthAccount } = await import("./auth");
		const result = { deleted_provider_id: "github", deleted_count: 1 };
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(result));
		});

		await expect(
			unlinkAuthAccount(fetchMock, "auth.session_token=token-1", "github")
		).resolves.toEqual(result);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/accounts/github",
			expect.objectContaining({ method: "DELETE" })
		);
	});

	it("rejects invalid successful API key responses", async () => {
		const { listAPIKeys } = await import("./auth");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ api_keys: [{ id: "api-key-1" }] }));
		});

		await expect(listAPIKeys(fetchMock, "auth.session_token=token-1", "user-1")).rejects.toMatchObject({
			status: 502,
			message: "Invalid API key list response"
		});
	});

	it("rejects invalid successful webhook responses", async () => {
		const { getWebhook, saveWebhook } = await import("./auth");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(
				JSON.stringify({
					webhook: {
						id: "webhook-1",
						user_id: "user-1",
						url: "https://example.com/webhook",
						events_active: ["job.started"],
						secret: "legacy-field",
						created_at: "2026-06-09T00:00:00Z",
						updated_at: "2026-06-09T00:00:00Z"
					}
				})
			);
		});

		await expect(getWebhook(fetchMock, "auth.session_token=token-1", "user-1")).rejects.toMatchObject({
			status: 502,
			message: "Invalid webhook response"
		});

		fetchMock.mockResolvedValueOnce(
			new Response(
				JSON.stringify({
					id: "webhook-1",
					user_id: "user-1",
					url: "https://example.com/webhook",
					events_active: ["job.deleted"],
					has_secret: true,
					created_at: "2026-06-09T00:00:00Z",
					updated_at: "2026-06-09T00:00:00Z"
				})
			)
		);
		await expect(
			saveWebhook(fetchMock, "auth.session_token=token-1", {
				userId: "user-1",
				url: "https://example.com/webhook",
				eventsActive: []
			})
		).rejects.toMatchObject({
			status: 502,
			message: "Invalid webhook response"
		});
	});

	it("signs out through the backend with cookies", async () => {
		const { AUTH_SESSION_COOKIE_NAME, signOut } = await import("./auth");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ success: true }));
		});

		await expect(signOut(fetchMock, `${AUTH_SESSION_COOKIE_NAME}=token-1`)).resolves.toEqual({
			success: true
		});

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/sign-out",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({})
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
		expect(headers.get("cookie")).toBe(`${AUTH_SESSION_COOKIE_NAME}=token-1`);
	});

	it("starts Google OAuth through the backend", async () => {
		const { startGoogleOAuth } = await import("./auth");
		const response = {
			authorizationUrl: "https://accounts.google.com/o/oauth2/auth?state=state-1",
			state: "state-1",
			stateExpiresAt: "2026-06-13T12:00:00Z"
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(response));
		});

		await expect(
			startGoogleOAuth(fetchMock, "http://localhost:5173/api/auth/callback/google")
		).resolves.toEqual(response);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/oauth/google/start",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({
					redirectURI: "http://localhost:5173/api/auth/callback/google"
				})
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("completes Google OAuth through the backend", async () => {
		const { signInGoogleOAuth } = await import("./auth");
		const response = {
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
				createdAt: "2026-06-13T12:00:00Z",
				updatedAt: "2026-06-13T12:00:00Z"
			}
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(response));
		});

		await expect(
			signInGoogleOAuth(fetchMock, {
				code: "code-1",
				state: "state-1",
				redirectURI: "http://localhost:5173/api/auth/callback/google"
			})
		).resolves.toEqual(response);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/oauth/google/callback",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({
					code: "code-1",
					state: "state-1",
					redirectURI: "http://localhost:5173/api/auth/callback/google"
				})
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("starts and completes Google account linking through the backend with cookies", async () => {
		const { linkGoogleAccount, startGoogleAccountLink } = await import("./auth");
		const started = {
			authorizationUrl: "https://accounts.google.com/o/oauth2/auth?state=state-1",
			state: "state-1",
			stateExpiresAt: "2026-06-13T12:00:00Z"
		};
		const linked = {
			id: "account-1",
			providerId: "google",
			createdAt: "2026-06-13T12:00:00Z",
			updatedAt: "2026-06-13T12:00:00Z"
		};
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(new Response(JSON.stringify(started)))
			.mockResolvedValueOnce(new Response(JSON.stringify(linked)));

		await expect(
			startGoogleAccountLink(
				fetchMock,
				"auth.session_token=token-1",
				"http://localhost:5173/api/auth/callback/link/google"
			)
		).resolves.toEqual(started);
		await expect(
			linkGoogleAccount(fetchMock, "auth.session_token=token-1", {
				code: "code-1",
				state: "state-1",
				redirectURI: "http://localhost:5173/api/auth/callback/link/google"
			})
		).resolves.toEqual(linked);

		expect(fetchMock).toHaveBeenNthCalledWith(
			1,
			"http://auth-api.test/api/auth/accounts/google/start",
			expect.objectContaining({ method: "POST" })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			2,
			"http://auth-api.test/api/auth/accounts/google/callback",
			expect.objectContaining({ method: "POST" })
		);
	});

	it("starts GitHub OAuth through the backend", async () => {
		const { startGitHubOAuth } = await import("./auth");
		const response = {
			authorizationUrl: "https://github.com/login/oauth/authorize?state=state-1",
			state: "state-1",
			stateExpiresAt: "2026-06-13T12:00:00Z"
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(response));
		});

		await expect(
			startGitHubOAuth(fetchMock, "http://localhost:5173/api/auth/callback/github")
		).resolves.toEqual(response);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/oauth/github/start",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({
					redirectURI: "http://localhost:5173/api/auth/callback/github"
				})
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("completes GitHub OAuth through the backend", async () => {
		const { signInGitHubOAuth } = await import("./auth");
		const response = {
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
				createdAt: "2026-06-13T12:00:00Z",
				updatedAt: "2026-06-13T12:00:00Z"
			}
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(response));
		});

		await expect(
			signInGitHubOAuth(fetchMock, {
				code: "code-1",
				state: "state-1",
				redirectURI: "http://localhost:5173/api/auth/callback/github"
			})
		).resolves.toEqual(response);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://auth-api.test/api/auth/oauth/github/callback",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({
					code: "code-1",
					state: "state-1",
					redirectURI: "http://localhost:5173/api/auth/callback/github"
				})
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("starts and completes GitHub account linking through the backend with cookies", async () => {
		const { linkGitHubAccount, startGitHubAccountLink } = await import("./auth");
		const started = {
			authorizationUrl: "https://github.com/login/oauth/authorize?state=state-1",
			state: "state-1",
			stateExpiresAt: "2026-06-13T12:00:00Z"
		};
		const linked = {
			id: "account-1",
			providerId: "github",
			createdAt: "2026-06-13T12:00:00Z",
			updatedAt: "2026-06-13T12:00:00Z"
		};
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(new Response(JSON.stringify(started)))
			.mockResolvedValueOnce(new Response(JSON.stringify(linked)));

		await expect(
			startGitHubAccountLink(
				fetchMock,
				"auth.session_token=token-1",
				"http://localhost:5173/api/auth/callback/link/github"
			)
		).resolves.toEqual(started);
		await expect(
			linkGitHubAccount(fetchMock, "auth.session_token=token-1", {
				code: "code-1",
				state: "state-1",
				redirectURI: "http://localhost:5173/api/auth/callback/link/github"
			})
		).resolves.toEqual(linked);

		expect(fetchMock).toHaveBeenNthCalledWith(
			1,
			"http://auth-api.test/api/auth/accounts/github/start",
			expect.objectContaining({ method: "POST" })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			2,
			"http://auth-api.test/api/auth/accounts/github/callback",
			expect.objectContaining({ method: "POST" })
		);
	});

	it("sets and clears the Google OAuth state cookie", async () => {
		const {
			GOOGLE_OAUTH_STATE_COOKIE_NAME,
			clearGoogleOAuthStateCookie,
			setGoogleOAuthStateCookie
		} = await import("./auth");
		const cookies = { set: vi.fn(), delete: vi.fn() };

		setGoogleOAuthStateCookie(
			cookies as never,
			"state-1",
			new Date(Date.now() + 60_000).toISOString()
		);
		clearGoogleOAuthStateCookie(cookies as never);

		expect(cookies.set).toHaveBeenCalledWith(
			GOOGLE_OAUTH_STATE_COOKIE_NAME,
			"state-1",
			expect.objectContaining({
				httpOnly: true,
				sameSite: "lax",
				path: "/"
			})
		);
		expect(cookies.delete).toHaveBeenCalledWith(GOOGLE_OAUTH_STATE_COOKIE_NAME, {
			path: "/"
		});
	});

	it("sets and clears the GitHub OAuth state cookie", async () => {
		const {
			GITHUB_OAUTH_STATE_COOKIE_NAME,
			clearGitHubOAuthStateCookie,
			setGitHubOAuthStateCookie
		} = await import("./auth");
		const cookies = { set: vi.fn(), delete: vi.fn() };

		setGitHubOAuthStateCookie(
			cookies as never,
			"state-1",
			new Date(Date.now() + 60_000).toISOString()
		);
		clearGitHubOAuthStateCookie(cookies as never);

		expect(cookies.set).toHaveBeenCalledWith(
			GITHUB_OAUTH_STATE_COOKIE_NAME,
			"state-1",
			expect.objectContaining({
				httpOnly: true,
				sameSite: "lax",
				path: "/"
			})
		);
		expect(cookies.delete).toHaveBeenCalledWith(GITHUB_OAUTH_STATE_COOKIE_NAME, {
			path: "/"
		});
	});

	it("sets and clears OAuth link state cookies", async () => {
		const {
			GITHUB_OAUTH_LINK_STATE_COOKIE_NAME,
			GOOGLE_OAUTH_LINK_STATE_COOKIE_NAME,
			clearGitHubOAuthLinkStateCookie,
			clearGoogleOAuthLinkStateCookie,
			setGitHubOAuthLinkStateCookie,
			setGoogleOAuthLinkStateCookie
		} = await import("./auth");
		const cookies = { set: vi.fn(), delete: vi.fn() };
		const expiresAt = new Date(Date.now() + 60_000).toISOString();

		setGoogleOAuthLinkStateCookie(cookies as never, "google-state", expiresAt);
		setGitHubOAuthLinkStateCookie(cookies as never, "github-state", expiresAt);
		clearGoogleOAuthLinkStateCookie(cookies as never);
		clearGitHubOAuthLinkStateCookie(cookies as never);

		expect(cookies.set).toHaveBeenCalledWith(
			GOOGLE_OAUTH_LINK_STATE_COOKIE_NAME,
			"google-state",
			expect.objectContaining({ httpOnly: true, sameSite: "lax", path: "/" })
		);
		expect(cookies.set).toHaveBeenCalledWith(
			GITHUB_OAUTH_LINK_STATE_COOKIE_NAME,
			"github-state",
			expect.objectContaining({ httpOnly: true, sameSite: "lax", path: "/" })
		);
		expect(cookies.delete).toHaveBeenCalledWith(GOOGLE_OAUTH_LINK_STATE_COOKIE_NAME, {
			path: "/"
		});
		expect(cookies.delete).toHaveBeenCalledWith(GITHUB_OAUTH_LINK_STATE_COOKIE_NAME, {
			path: "/"
		});
	});

	it("rejects auth requests before calling Go when the internal API token is missing", async () => {
		const { signInEmail } = await import("./auth");
		delete process.env.SYNCRA_INTERNAL_API_TOKEN;
		const fetchMock = vi.fn();

		await expect(
			signInEmail(fetchMock, {
				email: "ada@example.com",
				password: "password123",
				rememberMe: false
			})
		).rejects.toMatchObject({
			status: 500,
			message: "Authentication service is not configured"
		});
		expect(fetchMock).not.toHaveBeenCalled();
	});

	it("throws typed auth API errors from backend error responses", async () => {
		const { isAuthApiError, signInEmail } = await import("./auth");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ error: "invalid email or password" }), {
				status: 401
			});
		});

		await expect(
			signInEmail(fetchMock, {
				email: "ada@example.com",
				password: "wrong-password",
				rememberMe: false
			})
		).rejects.toMatchObject({
			status: 401,
			message: "invalid email or password"
		});

		try {
			await signInEmail(fetchMock, {
				email: "ada@example.com",
				password: "wrong-password",
				rememberMe: false
			});
		} catch (error) {
			expect(isAuthApiError(error)).toBe(true);
		}
	});

	it("defaults session cookies to secure outside development", async () => {
		process.env.NODE_ENV = "production";
		const { AUTH_SESSION_COOKIE_NAME, setSessionCookie } = await import("./auth");
		const cookies = { set: vi.fn() };

		setSessionCookie(
			cookies as never,
			{
				id: "session-1",
				token: "token-1",
				userId: "user-1",
				expiresAt: new Date(Date.now() + 60_000).toISOString(),
				createdAt: "2026-05-26T00:00:00Z",
				updatedAt: "2026-05-26T00:00:00Z"
			},
			true
		);

		expect(cookies.set).toHaveBeenCalledWith(
			AUTH_SESSION_COOKIE_NAME,
			"token-1",
			expect.objectContaining({ secure: true })
		);
	});

	it("converts backend outages into typed auth API errors", async () => {
		const { isAuthApiError, signInEmail } = await import("./auth");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit): Promise<Response> => {
			throw new Error("ECONNREFUSED");
		});

		await expect(
			signInEmail(fetchMock, {
				email: "ada@example.com",
				password: "password123",
				rememberMe: false
			})
		).rejects.toMatchObject({
			status: 503,
			message: "Authentication service unavailable"
		});

		try {
			await signInEmail(fetchMock, {
				email: "ada@example.com",
				password: "password123",
				rememberMe: false
			});
		} catch (error) {
			expect(isAuthApiError(error)).toBe(true);
		}
	});

	it("uses fallback messages for non-JSON backend errors", async () => {
		const { signInEmail } = await import("./auth");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response("<h1>bad gateway</h1>", { status: 502 });
		});

		await expect(
			signInEmail(fetchMock, {
				email: "ada@example.com",
				password: "password123",
				rememberMe: false
			})
		).rejects.toMatchObject({
			status: 502,
			message: "Authentication request failed"
		});
	});
});
