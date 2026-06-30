import type { Cookies } from '@sveltejs/kit';
import { apiBaseUrl, internalAPIHeaders, privateEnv } from './internal-api';

export const AUTH_SESSION_COOKIE_NAME = 'auth.session_token';
export const GOOGLE_OAUTH_STATE_COOKIE_NAME = 'auth.google_oauth_state';
export const GITHUB_OAUTH_STATE_COOKIE_NAME = 'auth.github_oauth_state';
export const PREFERRED_LANGUAGE_COOKIE_NAME = 'PARAGLIDE_LOCALE';
export const PREFERRED_LANGUAGE_COOKIE_MAX_AGE = 34560000;

const AUTH_DELIVERY_HEADER = 'X-Syncra-Auth-Delivery-Token';

type ServerFetch = typeof fetch;

export type PreferredLanguage = 'en' | 'ro';

export type AuthUser = {
	id: string;
	name: string;
	email: string;
	emailVerified: boolean;
	image: string | null;
	preferredLanguage: PreferredLanguage;
	role: 'user' | 'admin';
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

export type AuthSessionPayload = {
	session: AuthSession;
	user: AuthUser;
};

export type SignUpEmailResponse = {
	user?: AuthUser;
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

export type OAuthStartResponse = {
	authorizationUrl: string;
	state: string;
	stateExpiresAt: string;
};

export class AuthApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = 'AuthApiError';
		this.status = status;
	}
}

function authDeliveryToken() {
	return (privateEnv('AUTH_DELIVERY_TOKEN') || '').trim();
}

function authCookieSecure() {
	const configured = privateEnv('AUTH_COOKIE_SECURE');
	if (configured !== undefined && configured !== '') return configured === 'true';
	return privateEnv('NODE_ENV') !== 'development';
}

function authInternalHeaders() {
	const headers = internalAPIHeaders();
	if (!headers) throw new AuthApiError(500, 'Authentication service is not configured');
	return headers;
}

async function authRequest<T>(
	fetchFn: ServerFetch,
	path: string,
	options: {
		method?: 'DELETE' | 'GET' | 'POST' | 'PATCH';
		body?: unknown;
		cookieHeader?: string | null;
		trustedDelivery?: boolean;
	} = {}
) {
	const headers = authInternalHeaders();
	if (options.body !== undefined) headers.set('content-type', 'application/json');
	if (options.cookieHeader) headers.set('cookie', options.cookieHeader);
	if (options.trustedDelivery) {
		const token = authDeliveryToken();
		if (token) headers.set(AUTH_DELIVERY_HEADER, token);
	}

	let response: Response;
	try {
		response = await fetchFn(`${apiBaseUrl()}${path}`, {
			method: options.method ?? 'GET',
			headers,
			body: options.body === undefined ? undefined : JSON.stringify(options.body)
		});
	} catch {
		throw new AuthApiError(503, 'Authentication service unavailable');
	}

	const text = await response.text();
	const data = parseResponseJSON(text);
	if (!response.ok) {
		const message =
			data && typeof data === 'object' && 'error' in data && typeof data.error === 'string'
				? data.error
				: 'Authentication request failed';
		throw new AuthApiError(response.status, message);
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

export function isAuthApiError(error: unknown): error is AuthApiError {
	return error instanceof AuthApiError;
}

export async function getSession(fetchFn: ServerFetch, cookieHeader: string | null) {
	return authRequest<AuthSessionPayload | null>(fetchFn, '/api/auth/get-session', {
		cookieHeader
	});
}

export async function signUpEmail(
	fetchFn: ServerFetch,
	input: { name: string; email: string; password: string }
) {
	return authRequest<SignUpEmailResponse>(fetchFn, '/api/auth/sign-up/email', {
		method: 'POST',
		trustedDelivery: true,
		body: input
	});
}

export async function sendVerificationOTP(fetchFn: ServerFetch, email: string) {
	return authRequest<SendVerificationOTPResponse>(
		fetchFn,
		'/api/auth/email-otp/send-verification-otp',
		{
			method: 'POST',
			trustedDelivery: true,
			body: { email, type: 'email-verification' }
		}
	);
}

export async function verifyEmailOTP(fetchFn: ServerFetch, input: { email: string; otp: string }) {
	return authRequest<{ ok: boolean; user: AuthUser }>(fetchFn, '/api/auth/email-otp/verify-email', {
		method: 'POST',
		body: input
	});
}

export async function requestPasswordReset(fetchFn: ServerFetch, email: string) {
	return authRequest<RequestPasswordResetResponse>(fetchFn, '/api/auth/password-reset/request', {
		method: 'POST',
		trustedDelivery: true,
		body: { email }
	});
}

export async function resetPassword(
	fetchFn: ServerFetch,
	input: { email: string; token: string; password: string }
) {
	return authRequest<{ ok: boolean }>(fetchFn, '/api/auth/password-reset/confirm', {
		method: 'POST',
		body: input
	});
}

export async function signInEmail(
	fetchFn: ServerFetch,
	input: { email: string; password: string; rememberMe: boolean }
) {
	return authRequest<AuthSessionPayload>(fetchFn, '/api/auth/sign-in/email', {
		method: 'POST',
		body: input
	});
}

export async function signOut(fetchFn: ServerFetch, cookieHeader: string | null) {
	return authRequest<{ success: boolean }>(fetchFn, '/api/auth/sign-out', {
		method: 'POST',
		cookieHeader,
		body: {}
	});
}

export async function startGoogleOAuth(fetchFn: ServerFetch, redirectURI: string) {
	return authRequest<OAuthStartResponse>(fetchFn, '/api/auth/oauth/google/start', {
		method: 'POST',
		body: { redirectURI }
	});
}

export async function signInGoogleOAuth(
	fetchFn: ServerFetch,
	input: { code: string; state: string; redirectURI: string }
) {
	return authRequest<AuthSessionPayload>(fetchFn, '/api/auth/oauth/google/callback', {
		method: 'POST',
		body: input
	});
}

export async function startGitHubOAuth(fetchFn: ServerFetch, redirectURI: string) {
	return authRequest<OAuthStartResponse>(fetchFn, '/api/auth/oauth/github/start', {
		method: 'POST',
		body: { redirectURI }
	});
}

export async function signInGitHubOAuth(
	fetchFn: ServerFetch,
	input: { code: string; state: string; redirectURI: string }
) {
	return authRequest<AuthSessionPayload>(fetchFn, '/api/auth/oauth/github/callback', {
		method: 'POST',
		body: input
	});
}

export function setSessionCookie(cookies: Cookies, session: AuthSession, rememberMe: boolean) {
	const expires = new Date(session.expiresAt);
	const maxAge = Math.max(0, Math.floor((expires.getTime() - Date.now()) / 1000));
	cookies.set(AUTH_SESSION_COOKIE_NAME, session.token, {
		path: '/',
		httpOnly: true,
		sameSite: 'lax',
		secure: authCookieSecure(),
		expires: rememberMe ? expires : undefined,
		maxAge: rememberMe ? maxAge : undefined
	});
}

export function clearSessionCookie(cookies: Cookies) {
	cookies.delete(AUTH_SESSION_COOKIE_NAME, { path: '/' });
}

export function setGoogleOAuthStateCookie(cookies: Cookies, state: string, stateExpiresAt: string) {
	setOAuthStateCookie(cookies, GOOGLE_OAUTH_STATE_COOKIE_NAME, state, stateExpiresAt);
}

export function clearGoogleOAuthStateCookie(cookies: Cookies) {
	cookies.delete(GOOGLE_OAUTH_STATE_COOKIE_NAME, { path: '/' });
}

export function setGitHubOAuthStateCookie(cookies: Cookies, state: string, stateExpiresAt: string) {
	setOAuthStateCookie(cookies, GITHUB_OAUTH_STATE_COOKIE_NAME, state, stateExpiresAt);
}

export function clearGitHubOAuthStateCookie(cookies: Cookies) {
	cookies.delete(GITHUB_OAUTH_STATE_COOKIE_NAME, { path: '/' });
}

function setOAuthStateCookie(cookies: Cookies, cookieName: string, state: string, stateExpiresAt: string) {
	const expires = new Date(stateExpiresAt);
	const maxAge = Math.max(0, Math.floor((expires.getTime() - Date.now()) / 1000));
	cookies.set(cookieName, state, {
		path: '/',
		httpOnly: true,
		sameSite: 'lax',
		secure: authCookieSecure(),
		expires,
		maxAge
	});
}

export function setPreferredLanguageCookie(cookies: Cookies, language: unknown) {
	if (language !== 'en' && language !== 'ro') return false;
	cookies.set(PREFERRED_LANGUAGE_COOKIE_NAME, language, {
		path: '/',
		httpOnly: false,
		sameSite: 'lax',
		maxAge: PREFERRED_LANGUAGE_COOKIE_MAX_AGE
	});
	return true;
}

export function hasSessionCookie(cookieHeader: string | null) {
	return Boolean(cookieHeader?.includes(`${AUTH_SESSION_COOKIE_NAME}=`));
}
