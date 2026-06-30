import type { Cookies } from "@sveltejs/kit";
import { apiBaseUrl, internalAPIHeaders, privateEnv } from "./internal-api";

export const AUTH_SESSION_COOKIE_NAME = "auth.session_token";
export const GOOGLE_OAUTH_STATE_COOKIE_NAME = "auth.google_oauth_state";
export const GITHUB_OAUTH_STATE_COOKIE_NAME = "auth.github_oauth_state";
export const GOOGLE_OAUTH_LINK_STATE_COOKIE_NAME = "auth.google_oauth_link_state";
export const GITHUB_OAUTH_LINK_STATE_COOKIE_NAME = "auth.github_oauth_link_state";
export const PREFERRED_LANGUAGE_COOKIE_NAME = "PARAGLIDE_LOCALE";
export const PREFERRED_LANGUAGE_COOKIE_MAX_AGE = 34560000;

const AUTH_DELIVERY_HEADER = "X-Syncra-Auth-Delivery-Token";

type ServerFetch = typeof fetch;
export type PreferredLanguage = "en" | "ro";

export type AuthUser = {
	id: string;
	name: string;
	email: string;
	emailVerified: boolean;
	image: string | null;
	preferredLanguage: PreferredLanguage;
	role: "user" | "admin";
	lastLoginAt: string | null;
	createdAt: string;
	updatedAt: string;
};

export type AuthSession = {
	id: string;
	token: string;
	userId: string;
	expiresAt: string;
	ipAddress?: string;
	userAgent?: string;
	createdAt: string;
	updatedAt: string;
};

export type AuthSessionListItem = Omit<AuthSession, "token"> & {
	current: boolean;
};

export type AuthSessionListResponse = {
	sessions: AuthSessionListItem[];
};

export type DeleteAuthSessionResponse = {
	deleted_id: string;
	deleted_count: number;
};

export type AuthImpersonation = {
	adminUser: AuthUser;
	targetUser: AuthUser;
	startedAt: string;
};

export type AuthSessionPayload = {
	session: AuthSession;
	user: AuthUser;
	impersonation?: AuthImpersonation | null;
};

export type SignUpEmailResponse = {
	user: AuthUser;
	verificationRequired: boolean;
	verificationCode?: string;
	verificationExpiresAt?: string;
};

export type SendVerificationOTPResponse = {
	ok: boolean;
	verificationCode?: string;
	verificationExpiresAt?: string;
};

export type RequestPasswordResetResponse = {
	ok: boolean;
	resetToken?: string;
	resetExpiresAt?: string;
};

export type GoogleOAuthStartResponse = {
	authorizationUrl: string;
	state: string;
	stateExpiresAt: string;
};

export type GitHubOAuthStartResponse = GoogleOAuthStartResponse;

export type AuthAccountListItem = {
	id: string;
	providerId: "credential" | "google" | "github" | string;
	createdAt: string;
	updatedAt: string;
};

export type AuthAccountListResponse = {
	accounts: AuthAccountListItem[];
};

export type DeleteAuthAccountResponse = {
	deleted_provider_id: string;
	deleted_count: number;
};

export type APIKeyResponse = {
	id: string;
	user_id: string;
	name: string;
	key_prefix: string;
	api_key?: string;
	expires_at?: string;
	created_at: string;
	updated_at: string;
};

export type APIKeyListResponse = {
	api_keys: APIKeyResponse[];
};

export type CreateAPIKeyInput = {
	userId: string;
	name: string;
	expiresAt?: string;
};

export type DeleteAPIKeyInput = {
	userId: string;
	apiKeyId: string;
};

export type DeleteAPIKeyResponse = {
	deleted_id: string;
	deleted_count: number;
};

export type WebhookEvent = "job.started" | "job.failed" | "job.succeeded";

export type WebhookResponse = {
	id: string;
	user_id: string;
	url: string;
	events_active: WebhookEvent[];
	has_secret: boolean;
	secret_key?: string;
	created_at: string;
	updated_at: string;
};

export type WebhookEnvelopeResponse = {
	webhook: WebhookResponse | null;
};

export type SaveWebhookInput = {
	userId: string;
	url: string;
	eventsActive: WebhookEvent[];
};

export type DeleteWebhookResponse = {
	deleted_id: string;
	deleted_count: number;
};

export class AuthApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = "AuthApiError";
		this.status = status;
	}
}

function authDeliveryToken() {
	return (privateEnv("AUTH_DELIVERY_TOKEN") || "").trim();
}

function authCookieSecure() {
	const configured = privateEnv("AUTH_COOKIE_SECURE");
	if (configured !== undefined && configured !== "") {
		return configured === "true";
	}
	return privateEnv("NODE_ENV") !== "development";
}

function authInternalHeaders() {
	const headers = internalAPIHeaders();
	if (!headers) {
		throw new AuthApiError(500, "Authentication service is not configured");
	}
	return headers;
}

async function authRequest<T>(
	fetchFn: ServerFetch,
	path: string,
	options: {
		method?: "DELETE" | "GET" | "POST" | "PATCH";
		body?: unknown;
		cookieHeader?: string | null;
		trustedDelivery?: boolean;
	} = {}
) {
	const headers = authInternalHeaders();
	if (options.body !== undefined) headers.set("content-type", "application/json");
	if (options.cookieHeader) headers.set("cookie", options.cookieHeader);
	if (options.trustedDelivery) {
		const token = authDeliveryToken();
		if (token) headers.set(AUTH_DELIVERY_HEADER, token);
	}

	let response: Response;
	try {
		response = await fetchFn(`${apiBaseUrl()}${path}`, {
			method: options.method ?? "GET",
			headers,
			body: options.body === undefined ? undefined : JSON.stringify(options.body)
		});
	} catch {
		throw new AuthApiError(503, "Authentication service unavailable");
	}

	const text = await response.text();
	const data = parseResponseJSON(text);

	if (!response.ok) {
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "Authentication request failed";
		throw new AuthApiError(response.status, message);
	}

	if (data === undefined) {
		throw new AuthApiError(502, "Invalid authentication response");
	}

	return data as T;
}

function parseResponseJSON(text: string) {
	if (!text) return null;
	try {
		return JSON.parse(text) as unknown;
	} catch {
		return undefined;
	}
}

function isJsonObject(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isAuthSessionListItem(value: unknown): value is AuthSessionListItem {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.userId === "string" &&
		typeof value.expiresAt === "string" &&
		!("ipAddress" in value && typeof value.ipAddress !== "string") &&
		!("userAgent" in value && typeof value.userAgent !== "string") &&
		typeof value.createdAt === "string" &&
		typeof value.updatedAt === "string" &&
		typeof value.current === "boolean" &&
		!("token" in value)
	);
}

function isAuthSessionListResponse(value: unknown): value is AuthSessionListResponse {
	return isJsonObject(value) && Array.isArray(value.sessions) && value.sessions.every(isAuthSessionListItem);
}

function isDeleteAuthSessionResponse(value: unknown): value is DeleteAuthSessionResponse {
	return (
		isJsonObject(value) &&
		typeof value.deleted_id === "string" &&
		typeof value.deleted_count === "number" &&
		Number.isFinite(value.deleted_count)
	);
}

function isAuthAccountListItem(value: unknown): value is AuthAccountListItem {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.providerId === "string" &&
		typeof value.createdAt === "string" &&
		typeof value.updatedAt === "string" &&
		!("accountId" in value) &&
		!("accessToken" in value) &&
		!("refreshToken" in value) &&
		!("idToken" in value) &&
		!("password" in value)
	);
}

function isAuthAccountListResponse(value: unknown): value is AuthAccountListResponse {
	return isJsonObject(value) && Array.isArray(value.accounts) && value.accounts.every(isAuthAccountListItem);
}

function isDeleteAuthAccountResponse(value: unknown): value is DeleteAuthAccountResponse {
	return (
		isJsonObject(value) &&
		typeof value.deleted_provider_id === "string" &&
		typeof value.deleted_count === "number" &&
		Number.isFinite(value.deleted_count)
	);
}

function isAPIKeyResponse(value: unknown): value is APIKeyResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.user_id === "string" &&
		typeof value.name === "string" &&
		typeof value.key_prefix === "string" &&
		!("api_key" in value && typeof value.api_key !== "string") &&
		!("expires_at" in value && typeof value.expires_at !== "string") &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
	);
}

function isAPIKeyListResponse(value: unknown): value is APIKeyListResponse {
	return isJsonObject(value) && Array.isArray(value.api_keys) && value.api_keys.every(isAPIKeyResponse);
}

function isDeleteAPIKeyResponse(value: unknown): value is DeleteAPIKeyResponse {
	return (
		isJsonObject(value) &&
		typeof value.deleted_id === "string" &&
		typeof value.deleted_count === "number" &&
		Number.isFinite(value.deleted_count)
	);
}

function isWebhookEvent(value: unknown): value is WebhookEvent {
	return value === "job.started" || value === "job.failed" || value === "job.succeeded";
}

function isWebhookResponse(value: unknown): value is WebhookResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.user_id === "string" &&
		typeof value.url === "string" &&
		Array.isArray(value.events_active) &&
		value.events_active.every(isWebhookEvent) &&
		typeof value.has_secret === "boolean" &&
		!("secret_key" in value && typeof value.secret_key !== "string") &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
	);
}

function isWebhookEnvelopeResponse(value: unknown): value is WebhookEnvelopeResponse {
	return (
		isJsonObject(value) &&
		"webhook" in value &&
		(value.webhook === null || isWebhookResponse(value.webhook))
	);
}

function isDeleteWebhookResponse(value: unknown): value is DeleteWebhookResponse {
	return (
		isJsonObject(value) &&
		typeof value.deleted_id === "string" &&
		typeof value.deleted_count === "number" &&
		Number.isFinite(value.deleted_count)
	);
}

export function isAuthApiError(error: unknown): error is AuthApiError {
	return error instanceof AuthApiError;
}

export async function getSession(fetchFn: ServerFetch, cookieHeader: string | null) {
	return authRequest<AuthSessionPayload | null>(fetchFn, "/api/auth/get-session", {
		cookieHeader
	});
}

export async function listAuthSessions(fetchFn: ServerFetch, cookieHeader: string | null) {
	const data = await authRequest<unknown>(fetchFn, "/api/auth/sessions", {
		cookieHeader
	});
	if (!isAuthSessionListResponse(data)) {
		throw new AuthApiError(502, "Invalid session list response");
	}

	return data;
}

export async function revokeAuthSession(fetchFn: ServerFetch, cookieHeader: string | null, sessionId: string) {
	const data = await authRequest<unknown>(fetchFn, `/api/auth/sessions/${encodeURIComponent(sessionId)}`, {
		method: "DELETE",
		cookieHeader
	});
	if (!isDeleteAuthSessionResponse(data)) {
		throw new AuthApiError(502, "Invalid session delete response");
	}

	return data;
}

export type UpdateAuthUserInput = {
	name?: string;
	email?: string;
	image?: string | null;
	preferredLanguage?: PreferredLanguage;
	password?: string;
};

export async function updateAuthUser(fetchFn: ServerFetch, cookieHeader: string | null, input: UpdateAuthUserInput) {
	return authRequest<AuthUser>(fetchFn, "/api/auth/user", {
		method: "PATCH",
		cookieHeader,
		body: input
	});
}

export async function listAPIKeys(fetchFn: ServerFetch, cookieHeader: string | null, userId: string) {
	const data = await authRequest<unknown>(fetchFn, `/api/auth/apikeys/${encodeURIComponent(userId)}`, {
		cookieHeader
	});
	if (!isAPIKeyListResponse(data)) {
		throw new AuthApiError(502, "Invalid API key list response");
	}

	return data;
}

export async function createAPIKey(fetchFn: ServerFetch, cookieHeader: string | null, input: CreateAPIKeyInput) {
	const data = await authRequest<unknown>(fetchFn, "/api/auth/apikeys", {
		method: "POST",
		cookieHeader,
		body: {
			user_id: input.userId,
			name: input.name,
			...(input.expiresAt ? { expires_at: input.expiresAt } : {})
		}
	});
	if (!isAPIKeyResponse(data)) {
		throw new AuthApiError(502, "Invalid API key response");
	}

	return data;
}

export async function deleteAPIKey(fetchFn: ServerFetch, cookieHeader: string | null, input: DeleteAPIKeyInput) {
	const params = new URLSearchParams({
		user_id: input.userId,
		api_key_id: input.apiKeyId
	});
	const data = await authRequest<unknown>(fetchFn, `/api/auth/apikeys?${params.toString()}`, {
		method: "DELETE",
		cookieHeader
	});
	if (!isDeleteAPIKeyResponse(data)) {
		throw new AuthApiError(502, "Invalid API key delete response");
	}

	return data;
}

export async function getWebhook(fetchFn: ServerFetch, cookieHeader: string | null, userId: string) {
	const data = await authRequest<unknown>(fetchFn, `/api/auth/webhook/${encodeURIComponent(userId)}`, {
		cookieHeader
	});
	if (!isWebhookEnvelopeResponse(data)) {
		throw new AuthApiError(502, "Invalid webhook response");
	}

	return data;
}

export async function saveWebhook(fetchFn: ServerFetch, cookieHeader: string | null, input: SaveWebhookInput) {
	const data = await authRequest<unknown>(fetchFn, "/api/auth/webhook", {
		method: "POST",
		cookieHeader,
		body: {
			user_id: input.userId,
			url: input.url,
			events_active: input.eventsActive
		}
	});
	if (!isWebhookResponse(data)) {
		throw new AuthApiError(502, "Invalid webhook response");
	}

	return data;
}

export async function regenerateWebhookSecret(fetchFn: ServerFetch, cookieHeader: string | null, userId: string) {
	const data = await authRequest<unknown>(fetchFn, `/api/auth/webhook/${encodeURIComponent(userId)}/secret`, {
		method: "PATCH",
		cookieHeader
	});
	if (!isWebhookResponse(data)) {
		throw new AuthApiError(502, "Invalid webhook secret response");
	}

	return data;
}

export async function deleteWebhook(fetchFn: ServerFetch, cookieHeader: string | null, userId: string) {
	const params = new URLSearchParams({ user_id: userId });
	const data = await authRequest<unknown>(fetchFn, `/api/auth/webhook?${params.toString()}`, {
		method: "DELETE",
		cookieHeader
	});
	if (!isDeleteWebhookResponse(data)) {
		throw new AuthApiError(502, "Invalid webhook delete response");
	}

	return data;
}

export async function listAuthAccounts(fetchFn: ServerFetch, cookieHeader: string | null) {
	const data = await authRequest<unknown>(fetchFn, "/api/auth/accounts", {
		cookieHeader
	});
	if (!isAuthAccountListResponse(data)) {
		throw new AuthApiError(502, "Invalid linked account list response");
	}

	return data;
}

export async function unlinkAuthAccount(fetchFn: ServerFetch, cookieHeader: string | null, providerId: string) {
	const data = await authRequest<unknown>(fetchFn, `/api/auth/accounts/${encodeURIComponent(providerId)}`, {
		method: "DELETE",
		cookieHeader
	});
	if (!isDeleteAuthAccountResponse(data)) {
		throw new AuthApiError(502, "Invalid linked account delete response");
	}

	return data;
}

export async function signUpEmail(fetchFn: ServerFetch, input: { name: string; email: string; password: string }) {
	return authRequest<SignUpEmailResponse>(fetchFn, "/api/auth/sign-up/email", {
		method: "POST",
		trustedDelivery: true,
		body: input
	});
}

export async function sendVerificationOTP(fetchFn: ServerFetch, email: string) {
	return authRequest<SendVerificationOTPResponse>(fetchFn, "/api/auth/email-otp/send-verification-otp", {
		method: "POST",
		trustedDelivery: true,
		body: { email, type: "email-verification" }
	});
}

export async function verifyEmailOTP(fetchFn: ServerFetch, input: { email: string; otp: string }) {
	return authRequest<{ ok: boolean; user: AuthUser }>(fetchFn, "/api/auth/email-otp/verify-email", {
		method: "POST",
		body: input
	});
}

export async function requestPasswordReset(fetchFn: ServerFetch, email: string) {
	return authRequest<RequestPasswordResetResponse>(fetchFn, "/api/auth/password-reset/request", {
		method: "POST",
		trustedDelivery: true,
		body: { email }
	});
}

export async function resetPassword(fetchFn: ServerFetch, input: { email: string; token: string; password: string }) {
	return authRequest<{ ok: boolean }>(fetchFn, "/api/auth/password-reset/confirm", {
		method: "POST",
		body: input
	});
}

export async function signInEmail(
	fetchFn: ServerFetch,
	input: { email: string; password: string; rememberMe: boolean }
) {
	return authRequest<AuthSessionPayload>(fetchFn, "/api/auth/sign-in/email", {
		method: "POST",
		body: input
	});
}

export async function startGoogleOAuth(fetchFn: ServerFetch, redirectURI: string) {
	return authRequest<GoogleOAuthStartResponse>(fetchFn, "/api/auth/oauth/google/start", {
		method: "POST",
		body: { redirectURI }
	});
}

export async function startGoogleAccountLink(fetchFn: ServerFetch, cookieHeader: string | null, redirectURI: string) {
	return authRequest<GoogleOAuthStartResponse>(fetchFn, "/api/auth/accounts/google/start", {
		method: "POST",
		cookieHeader,
		body: { redirectURI }
	});
}

export async function signInGoogleOAuth(
	fetchFn: ServerFetch,
	input: { code: string; state: string; redirectURI: string }
) {
	return authRequest<AuthSessionPayload>(fetchFn, "/api/auth/oauth/google/callback", {
		method: "POST",
		body: input
	});
}

export async function linkGoogleAccount(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	input: { code: string; state: string; redirectURI: string }
) {
	const data = await authRequest<unknown>(fetchFn, "/api/auth/accounts/google/callback", {
		method: "POST",
		cookieHeader,
		body: input
	});
	if (!isAuthAccountListItem(data)) {
		throw new AuthApiError(502, "Invalid linked account response");
	}

	return data;
}

export async function startGitHubOAuth(fetchFn: ServerFetch, redirectURI: string) {
	return authRequest<GitHubOAuthStartResponse>(fetchFn, "/api/auth/oauth/github/start", {
		method: "POST",
		body: { redirectURI }
	});
}

export async function startGitHubAccountLink(fetchFn: ServerFetch, cookieHeader: string | null, redirectURI: string) {
	return authRequest<GitHubOAuthStartResponse>(fetchFn, "/api/auth/accounts/github/start", {
		method: "POST",
		cookieHeader,
		body: { redirectURI }
	});
}

export async function signInGitHubOAuth(
	fetchFn: ServerFetch,
	input: { code: string; state: string; redirectURI: string }
) {
	return authRequest<AuthSessionPayload>(fetchFn, "/api/auth/oauth/github/callback", {
		method: "POST",
		body: input
	});
}

export async function linkGitHubAccount(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	input: { code: string; state: string; redirectURI: string }
) {
	const data = await authRequest<unknown>(fetchFn, "/api/auth/accounts/github/callback", {
		method: "POST",
		cookieHeader,
		body: input
	});
	if (!isAuthAccountListItem(data)) {
		throw new AuthApiError(502, "Invalid linked account response");
	}

	return data;
}

export async function signOut(fetchFn: ServerFetch, cookieHeader: string | null) {
	return authRequest<{ success: boolean }>(fetchFn, "/api/auth/sign-out", {
		method: "POST",
		cookieHeader,
		body: {}
	});
}

export function setSessionCookie(cookies: Cookies, session: AuthSession, rememberMe: boolean) {
	const expires = new Date(session.expiresAt);
	const maxAge = Math.max(0, Math.floor((expires.getTime() - Date.now()) / 1000));

	cookies.set(AUTH_SESSION_COOKIE_NAME, session.token, {
		path: "/",
		httpOnly: true,
		sameSite: "lax",
		secure: authCookieSecure(),
		expires: rememberMe ? expires : undefined,
		maxAge: rememberMe ? maxAge : undefined
	});
}

export function isSupportedPreferredLanguage(value: unknown): value is PreferredLanguage {
	return value === "en" || value === "ro";
}

export function setPreferredLanguageCookie(cookies: Cookies, language: unknown) {
	if (!isSupportedPreferredLanguage(language)) return false;

	cookies.set(PREFERRED_LANGUAGE_COOKIE_NAME, language, {
		path: "/",
		httpOnly: false,
		sameSite: "lax",
		maxAge: PREFERRED_LANGUAGE_COOKIE_MAX_AGE
	});
	return true;
}

export function clearSessionCookie(cookies: Cookies) {
	cookies.delete(AUTH_SESSION_COOKIE_NAME, { path: "/" });
}

export function setGoogleOAuthStateCookie(cookies: Cookies, state: string, stateExpiresAt: string) {
	setOAuthStateCookie(cookies, GOOGLE_OAUTH_STATE_COOKIE_NAME, state, stateExpiresAt);
}

export function setGoogleOAuthLinkStateCookie(cookies: Cookies, state: string, stateExpiresAt: string) {
	setOAuthStateCookie(cookies, GOOGLE_OAUTH_LINK_STATE_COOKIE_NAME, state, stateExpiresAt);
}

export function clearGoogleOAuthStateCookie(cookies: Cookies) {
	cookies.delete(GOOGLE_OAUTH_STATE_COOKIE_NAME, { path: "/" });
}

export function clearGoogleOAuthLinkStateCookie(cookies: Cookies) {
	cookies.delete(GOOGLE_OAUTH_LINK_STATE_COOKIE_NAME, { path: "/" });
}

export function setGitHubOAuthStateCookie(cookies: Cookies, state: string, stateExpiresAt: string) {
	setOAuthStateCookie(cookies, GITHUB_OAUTH_STATE_COOKIE_NAME, state, stateExpiresAt);
}

export function setGitHubOAuthLinkStateCookie(cookies: Cookies, state: string, stateExpiresAt: string) {
	setOAuthStateCookie(cookies, GITHUB_OAUTH_LINK_STATE_COOKIE_NAME, state, stateExpiresAt);
}

export function clearGitHubOAuthStateCookie(cookies: Cookies) {
	cookies.delete(GITHUB_OAUTH_STATE_COOKIE_NAME, { path: "/" });
}

export function clearGitHubOAuthLinkStateCookie(cookies: Cookies) {
	cookies.delete(GITHUB_OAUTH_LINK_STATE_COOKIE_NAME, { path: "/" });
}

function setOAuthStateCookie(cookies: Cookies, cookieName: string, state: string, stateExpiresAt: string) {
	const expires = new Date(stateExpiresAt);
	const maxAge = Math.max(0, Math.floor((expires.getTime() - Date.now()) / 1000));

	cookies.set(cookieName, state, {
		path: "/",
		httpOnly: true,
		sameSite: "lax",
		secure: authCookieSecure(),
		expires,
		maxAge
	});
}

export function hasSessionCookie(cookieHeader: string | null) {
	return Boolean(cookieHeader?.includes(`${AUTH_SESSION_COOKIE_NAME}=`));
}
