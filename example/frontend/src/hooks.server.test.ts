import { beforeEach, describe, expect, it, vi } from "vitest";

const {
	getSessionMock,
	clearSessionCookieMock,
	hasSessionCookieMock,
	setPreferredLanguageCookieMock,
	loggerMock,
	randomUUIDMock,
	rootLoggerMock
} = vi.hoisted(() => {
	const logger = {
		debug: vi.fn(),
		info: vi.fn(),
		warn: vi.fn(),
		error: vi.fn(),
		child: vi.fn()
	};
	logger.child.mockReturnValue(logger);

	return {
		getSessionMock: vi.fn(),
		clearSessionCookieMock: vi.fn(),
		hasSessionCookieMock: vi.fn(),
		setPreferredLanguageCookieMock: vi.fn(),
		loggerMock: logger,
		randomUUIDMock: vi.fn(),
		rootLoggerMock: {
			debug: vi.fn(),
			info: vi.fn(),
			warn: vi.fn(),
			error: vi.fn(),
			child: vi.fn(() => logger)
		}
	};
});

vi.mock("node:crypto", () => ({
	randomUUID: randomUUIDMock
}));

vi.mock("$lib/server/auth", () => ({
	getSession: getSessionMock,
	clearSessionCookie: clearSessionCookieMock,
	hasSessionCookie: hasSessionCookieMock,
	setPreferredLanguageCookie: setPreferredLanguageCookieMock
}));

vi.mock("$lib/server/logging", async () => {
	const actual = await vi.importActual<typeof import("$lib/server/logging")>("$lib/server/logging");
	return {
		...actual,
		rootLogger: rootLoggerMock
	};
});

import { requestIdHeader } from "$lib/server/logging";
import { handle, localizePageAttributes } from "./hooks.server";

function event(pathname: string, headers: HeadersInit = {}) {
	const url = new URL(`http://localhost${pathname}`);

	return {
		url,
		request: new Request(url, { headers }),
		fetch: vi.fn(),
		cookies: { delete: vi.fn(), set: vi.fn() },
		locals: {},
		route: { id: pathname }
	};
}

describe("server hooks", () => {
	beforeEach(() => {
		getSessionMock.mockReset();
		getSessionMock.mockResolvedValue(null);
		clearSessionCookieMock.mockReset();
		hasSessionCookieMock.mockReset();
		setPreferredLanguageCookieMock.mockReset();
		randomUUIDMock.mockReset();
		randomUUIDMock.mockReturnValue("generated-request-id");
		loggerMock.debug.mockReset();
		loggerMock.info.mockReset();
		loggerMock.warn.mockReset();
		loggerMock.error.mockReset();
		loggerMock.child.mockReset();
		loggerMock.child.mockReturnValue(loggerMock);
		rootLoggerMock.child.mockReset();
		rootLoggerMock.child.mockReturnValue(loggerMock);
	});

	it("adds generated request ids to locals and response headers", async () => {
		const requestEvent = event("/pricing");
		const response = await handle({
			event: requestEvent as never,
			resolve: vi.fn(async () => new Response("ok", { headers: { "content-length": "2" } })) as never
		});

		expect(requestEvent.locals).toMatchObject({
			logger: loggerMock,
			requestId: "generated-request-id"
		});
		expect(response.headers.get(requestIdHeader)).toBe("generated-request-id");
		expect(rootLoggerMock.child).toHaveBeenCalledWith({
			component: "http",
			request_id: "generated-request-id"
		});
		expect(loggerMock.info).toHaveBeenCalledWith(
			"http.request_completed",
			expect.objectContaining({
				method: "GET",
				route: "/pricing",
				path: "/pricing",
				status: 200,
				response_bytes: 2
			})
		);
	});

	it("reuses valid incoming request ids", async () => {
		const requestEvent = event("/pricing", {
			[requestIdHeader]: "request-123"
		});

		const response = await handle({
			event: requestEvent as never,
			resolve: vi.fn(async () => new Response("ok")) as never
		});

		expect(requestEvent.locals).toMatchObject({ requestId: "request-123" });
		expect(response.headers.get(requestIdHeader)).toBe("request-123");
		expect(randomUUIDMock).not.toHaveBeenCalled();
	});

	it("replaces invalid incoming request ids", async () => {
		const requestEvent = event("/pricing", {
			[requestIdHeader]: "bad request"
		});

		const response = await handle({
			event: requestEvent as never,
			resolve: vi.fn(async () => new Response("ok")) as never
		});

		expect(requestEvent.locals).toMatchObject({
			requestId: "generated-request-id"
		});
		expect(response.headers.get(requestIdHeader)).toBe("generated-request-id");
		expect(randomUUIDMock).toHaveBeenCalled();
	});

	it("uses the Paraglide locale cookie for document attributes", async () => {
		const requestEvent = event("/pricing", { cookie: "PARAGLIDE_LOCALE=ro" });
		const resolve = vi.fn(async (_event, options) => {
			if (!options?.transformPageChunk) {
				return new Response("missing transform", { status: 500 });
			}

			const html = await options.transformPageChunk({
				html: '<!doctype html><html lang="%lang%" dir="%dir%"><body>Pricing</body></html>',
				done: true
			});

			return new Response(html);
		});

		const response = await handle({
			event: requestEvent as never,
			resolve: resolve as never
		});

		await expect(response.text()).resolves.toContain('<html lang="ro" dir="ltr">');
		expect(response.status).toBe(200);
	});

	it("keeps same-URL routing when locale comes from the cookie strategy", async () => {
		const requestEvent = event("/pricing", { cookie: "PARAGLIDE_LOCALE=ro" });
		const resolve = vi.fn(async (resolvedEvent) => {
			expect(new URL(resolvedEvent.request.url).pathname).toBe("/pricing");
			expect(resolvedEvent.url.pathname).toBe("/pricing");

			return new Response("ok");
		});

		const response = await handle({
			event: requestEvent as never,
			resolve: resolve as never
		});

		expect(response.status).toBe(200);
		expect(response.headers.get("location")).toBeNull();
		expect(new URL(requestEvent.request.url).pathname).toBe("/pricing");
	});

	it("runs caller page transforms after locale document attributes are replaced", async () => {
		const callerTransform = vi.fn(({ html }: { html: string }) =>
			html.replace("</html>", "<!-- caller transform --></html>")
		);
		const transform = localizePageAttributes("ro", callerTransform);

		const html = await transform({
			html: '<html lang="%lang%" dir="%dir%"><body>Pricing</body></html>',
			done: true
		});

		expect(callerTransform).toHaveBeenCalledWith({
			html: '<html lang="ro" dir="ltr"><body>Pricing</body></html>',
			done: true
		});
		expect(html).toBe('<html lang="ro" dir="ltr"><body>Pricing</body><!-- caller transform --></html>');
	});

	it.each([
		[200, "info"],
		[404, "warn"],
		[503, "error"]
	] as const)("logs request completion at the expected level for HTTP %s", async (status, level) => {
		await handle({
			event: event("/pricing") as never,
			resolve: vi.fn(async () => new Response("ok", { status })) as never
		});

		expect(loggerMock[level]).toHaveBeenCalledWith("http.request_completed", expect.objectContaining({ status }));
	});

	it.each([
		"/api/auth/user",
		"/api/collections",
		"/api/collections/collection-1",
		"/api/billing/balance",
		"/api/billing/checkout",
		"/api/billing/profile",
		"/api/billing/credit-usage-history",
		"/api/dashboard/summary",
		"/api/ocr/documents/document-1",
		"/api/admin/users",
		"/admin-impersonation/stop"
	])("requires session load for %s", async (pathname) => {
		const sessionError = new Error("auth service down");
		getSessionMock.mockRejectedValueOnce(sessionError);

		await expect(
			handle({
				event: event(pathname) as never,
				resolve: vi.fn(async () => new Response("ok")) as never
			})
		).rejects.toMatchObject({
			status: 503,
			body: { message: "Authentication service unavailable" }
		});

		expect(loggerMock.error).toHaveBeenCalledWith("auth.session_load_failed", {
			error: "auth service down",
			path: pathname
		});
		expect(loggerMock.error).toHaveBeenCalledWith(
			"http.request_completed",
			expect.objectContaining({ status: 503, path: pathname })
		);
		expect(loggerMock.error).not.toHaveBeenCalledWith("http.request_error", expect.anything());
	});

	it("redirects unauthenticated admin portal requests to login", async () => {
		getSessionMock.mockResolvedValueOnce(null);

		await expect(
			handle({
				event: event("/admin-portal/users") as never,
				resolve: vi.fn(async () => new Response("ok")) as never
			})
		).rejects.toMatchObject({
			status: 303,
			location: "/login"
		});
		expect(loggerMock.error).not.toHaveBeenCalledWith("http.request_error", expect.anything());
		expect(loggerMock.info).toHaveBeenCalledWith("http.request_completed", expect.objectContaining({ status: 303 }));
	});

	it("rejects non-admin users from the admin portal", async () => {
		getSessionMock.mockResolvedValueOnce({
			session: { id: "session-1", token: "token", userId: "user-1" },
			user: { id: "user-1", email: "user@example.com", role: "user" }
		});

		await expect(
			handle({
				event: event("/admin-portal/users") as never,
				resolve: vi.fn(async () => new Response("ok")) as never
			})
		).rejects.toMatchObject({
			status: 403,
			body: { message: "Admin access required" }
		});
		expect(loggerMock.error).not.toHaveBeenCalledWith("http.request_error", expect.anything());
		expect(loggerMock.warn).toHaveBeenCalledWith("http.request_completed", expect.objectContaining({ status: 403 }));
	});

	it("allows admin users into the admin portal", async () => {
		const adminUser = {
			id: "admin-1",
			email: "admin@example.com",
			role: "admin",
			preferredLanguage: "ro"
		};
		getSessionMock.mockResolvedValueOnce({
			session: { id: "session-1", token: "token", userId: "admin-1" },
			user: adminUser
		});
		const resolve = vi.fn(async () => new Response("ok"));
		const requestEvent = event("/admin-portal/users");

		const response = await handle({
			event: requestEvent as never,
			resolve: resolve as never
		});

		expect(response.status).toBe(200);
		expect(requestEvent.locals).toMatchObject({ user: adminUser, adminUser });
		expect(setPreferredLanguageCookieMock).toHaveBeenCalledWith(requestEvent.cookies, "ro");
		expect(resolve).toHaveBeenCalled();
	});

	it("sets effective user and original admin while impersonating", async () => {
		const adminUser = {
			id: "admin-1",
			email: "admin@example.com",
			role: "admin"
		};
		const targetUser = {
			id: "user-1",
			email: "user@example.com",
			role: "user"
		};
		const impersonation = {
			adminUser,
			targetUser,
			startedAt: "2026-06-14T10:00:00Z"
		};
		getSessionMock.mockResolvedValueOnce({
			session: { id: "session-1", token: "token", userId: "admin-1" },
			user: targetUser,
			impersonation
		});
		const requestEvent = event("/api/collections");
		const resolve = vi.fn(async () => new Response("ok"));

		const response = await handle({
			event: requestEvent as never,
			resolve: resolve as never
		});

		expect(response.status).toBe(200);
		expect(requestEvent.locals).toMatchObject({
			user: targetUser,
			adminUser,
			impersonation
		});
		expect(resolve).toHaveBeenCalled();
	});

	it("allows admin portal access while the effective user is impersonated", async () => {
		const adminUser = {
			id: "admin-1",
			email: "admin@example.com",
			role: "admin"
		};
		const targetUser = {
			id: "user-1",
			email: "user@example.com",
			role: "user"
		};
		getSessionMock.mockResolvedValueOnce({
			session: { id: "session-1", token: "token", userId: "admin-1" },
			user: targetUser,
			impersonation: {
				adminUser,
				targetUser,
				startedAt: "2026-06-14T10:00:00Z"
			}
		});
		const resolve = vi.fn(async () => new Response("ok"));

		const response = await handle({
			event: event("/admin-portal/users") as never,
			resolve: resolve as never
		});

		expect(response.status).toBe(200);
		expect(resolve).toHaveBeenCalled();
	});

	it("redirects authenticated users away from recover password", async () => {
		getSessionMock.mockResolvedValueOnce({
			session: { id: "session-1", token: "token", userId: "user-1" },
			user: { id: "user-1", email: "user@example.com", role: "user" }
		});

		await expect(
			handle({
				event: event("/recover-password") as never,
				resolve: vi.fn(async () => new Response("ok")) as never
			})
		).rejects.toMatchObject({
			status: 303,
			location: "/app"
		});
	});
});
