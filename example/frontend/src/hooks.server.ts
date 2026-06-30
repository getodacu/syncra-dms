import { randomUUID } from "node:crypto";
import {
	error,
	isHttpError,
	isRedirect,
	redirect,
	type Handle,
	type RequestEvent,
	type ResolveOptions
} from "@sveltejs/kit";
import { paraglideMiddleware } from "$lib/paraglide/server";
import { getTextDirection } from "$lib/paraglide/runtime";
import type { Locale } from "$lib/paraglide/runtime";
import { clearSessionCookie, getSession, hasSessionCookie, setPreferredLanguageCookie } from "$lib/server/auth";
import {
	cleanRequestId,
	requestIdHeader,
	rootLogger,
	safeError,
	userAgentClass,
	type Logger,
	type LogAttrs
} from "$lib/server/logging";

const guestOnlyRoutes = new Set(["/login", "/signup", "/signup-confirmation", "/recover-password"]);

function isProtectedRoute(pathname: string) {
	return pathname === "/app" || pathname.startsWith("/app/");
}

function isAdminPortalRoute(pathname: string) {
	return pathname === "/admin-portal" || pathname.startsWith("/admin-portal/");
}

function requiresSessionLoad(pathname: string) {
	return (
		isProtectedRoute(pathname) ||
		isAdminPortalRoute(pathname) ||
		pathname === "/api/admin" ||
		pathname.startsWith("/api/admin/") ||
		pathname === "/api/schemas" ||
		pathname.startsWith("/api/schemas/") ||
		pathname.startsWith("/api/json-recipes/") ||
		pathname === "/api/collections" ||
		pathname.startsWith("/api/collections/") ||
		pathname === "/api/billing/balance" ||
		pathname === "/api/billing/checkout" ||
		pathname === "/api/billing/profile" ||
		pathname === "/api/billing/credit-usage-history" ||
		pathname === "/api/dashboard/summary" ||
		pathname === "/api/auth/user" ||
		pathname === "/api/auth/apikeys" ||
		pathname === "/admin-impersonation/stop" ||
		pathname === "/api/ocr/documents" ||
		pathname.startsWith("/api/ocr/documents/") ||
		pathname.startsWith("/api/ocr/document/") ||
		pathname === "/api/ocr/jobs" ||
		pathname.startsWith("/api/ocr/jobs/")
	);
}

const appHandle: Handle = async ({ event, resolve }) => {
	const start = Date.now();
	const requestId = cleanRequestId(event.request.headers.get(requestIdHeader)) || randomUUID();
	const requestLogger = rootLogger.child({
		component: "http",
		request_id: requestId
	});
	event.locals.requestId = requestId;
	event.locals.logger = requestLogger;
	const cookieHeader = event.request.headers.get("cookie");

	try {
		try {
			const auth = await getSession(event.fetch, cookieHeader);
			event.locals.session = auth?.session ?? null;
			event.locals.user = auth?.user ?? null;
			event.locals.impersonation = auth?.impersonation ?? null;
			event.locals.adminUser = auth?.impersonation?.adminUser ?? (auth?.user?.role === "admin" ? auth.user : null);
			if (auth?.user?.preferredLanguage) {
				setPreferredLanguageCookie(event.cookies, auth.user.preferredLanguage);
			}
			if (!auth && hasSessionCookie(cookieHeader)) {
				clearSessionCookie(event.cookies);
			}
		} catch (sessionError) {
			event.locals.session = null;
			event.locals.user = null;
			event.locals.adminUser = null;
			event.locals.impersonation = null;
			if (requiresSessionLoad(event.url.pathname)) {
				requestLogger.error("auth.session_load_failed", {
					error: safeError(sessionError),
					path: event.url.pathname
				});
				error(503, "Authentication service unavailable");
			}
		}

		if (isProtectedRoute(event.url.pathname) && !event.locals.user) {
			redirect(303, "/login");
		}

		if (isAdminPortalRoute(event.url.pathname)) {
			if (!event.locals.user && !event.locals.adminUser) {
				redirect(303, "/login");
			}
			if (!event.locals.adminUser || event.locals.adminUser.role !== "admin") {
				error(403, "Admin access required");
			}
		}

		if (guestOnlyRoutes.has(event.url.pathname) && event.locals.user) {
			redirect(303, "/app");
		}

		const response = await resolve(event);
		response.headers.set(requestIdHeader, requestId);
		logRequestCompleted(requestLogger, event, start, response.status, response);
		return response;
	} catch (thrown) {
		const status = statusForThrown(thrown);
		if (!isRedirect(thrown) && !isHttpError(thrown)) {
			requestLogger.error("http.request_error", {
				...requestLogAttrs(event, start, status),
				error: safeError(thrown)
			});
		}
		logRequestCompleted(requestLogger, event, start, status);
		throw thrown;
	}
};

export const handle: Handle = ({ event, resolve }) =>
	paraglideMiddleware(event.request, ({ request: localizedRequest, locale }) => {
		event.request = localizedRequest;

		return appHandle({
			event,
			resolve: (event, options) =>
				resolve(event, {
					...options,
					transformPageChunk: localizePageAttributes(locale, options?.transformPageChunk)
				})
		});
	});

export function localizePageAttributes(
	locale: Locale,
	transformPageChunk?: ResolveOptions["transformPageChunk"]
): NonNullable<ResolveOptions["transformPageChunk"]> {
	return ({ html, done }) => {
		const transformed = html.replace("%lang%", locale).replace("%dir%", getTextDirection(locale));

		return transformPageChunk ? transformPageChunk({ html: transformed, done }) : transformed;
	};
}

function statusForThrown(thrown: unknown) {
	if (isRedirect(thrown) || isHttpError(thrown)) return thrown.status;
	return 500;
}

function logRequestCompleted(logger: Logger, event: RequestEvent, start: number, status: number, response?: Response) {
	const attrs = requestLogAttrs(event, start, status, response);
	if (status >= 500) {
		logger.error("http.request_completed", attrs);
	} else if (status >= 400) {
		logger.warn("http.request_completed", attrs);
	} else {
		logger.info("http.request_completed", attrs);
	}
}

function requestLogAttrs(event: RequestEvent, start: number, status: number, response?: Response): LogAttrs {
	const attrs: LogAttrs = {
		method: event.request.method,
		route: event.route.id ?? "unmatched",
		path: event.url.pathname,
		status,
		duration_ms: Date.now() - start,
		user_agent_class: userAgentClass(event.request.headers.get("user-agent"))
	};
	const requestBytes = contentLength(event.request.headers);
	if (requestBytes !== undefined) attrs.request_bytes = requestBytes;
	if (response) {
		const responseBytes = contentLength(response.headers);
		if (responseBytes !== undefined) attrs.response_bytes = responseBytes;
	}
	return attrs;
}

function contentLength(headers: Headers) {
	const value = headers.get("content-length");
	if (!value) return undefined;
	const bytes = Number(value);
	if (!Number.isSafeInteger(bytes) || bytes < 0) return undefined;
	return bytes;
}
