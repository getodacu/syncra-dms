import { beforeEach, describe, expect, it, vi } from "vitest";

import { AuthApiError } from "$lib/server/auth";
import { DELETE, GET, POST } from "./+server";
import type { RequestEvent } from "./$types";

const { deleteWebhookMock, getWebhookMock, saveWebhookMock, AuthApiErrorMock } = vi.hoisted(
	() => {
		class MockAuthApiError extends Error {
			status: number;

			constructor(status: number, message: string) {
				super(message);
				this.name = "AuthApiError";
				this.status = status;
			}
		}

		return {
			deleteWebhookMock: vi.fn(),
			getWebhookMock: vi.fn(),
			saveWebhookMock: vi.fn(),
			AuthApiErrorMock: MockAuthApiError
		};
	}
);

vi.mock("$lib/server/auth", () => ({
	deleteWebhook: deleteWebhookMock,
	getWebhook: getWebhookMock,
	saveWebhook: saveWebhookMock,
	AuthApiError: AuthApiErrorMock,
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock
}));

function webhookFixture(overrides: Record<string, unknown> = {}) {
	return {
		id: "webhook-1",
		user_id: "user-1",
		url: "https://example.com/webhook",
		events_active: ["job.started", "job.succeeded"],
		created_at: "2026-06-09T00:00:00Z",
		updated_at: "2026-06-09T00:00:00Z",
		...overrides
	};
}

function createEvent(
	method: "DELETE" | "GET" | "POST",
	options: {
		body?: unknown;
		user?: unknown;
		cookie?: string;
	} = {}
) {
	const headers = new Headers();
	if (options.cookie) headers.set("cookie", options.cookie);
	if (options.body !== undefined) headers.set("content-type", "application/json");

	const request = new Request("http://localhost/api/auth/webhook", {
		method,
		headers,
		body: options.body === undefined ? undefined : JSON.stringify(options.body)
	});

	return {
		request,
		url: new URL(request.url),
		fetch: vi.fn(),
		locals: { user: options.user === undefined ? { id: "user-1" } : options.user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("webhook auth proxy endpoint", () => {
	beforeEach(() => {
		deleteWebhookMock.mockReset();
		getWebhookMock.mockReset();
		saveWebhookMock.mockReset();
	});

	it("returns 401 for unauthenticated get requests", async () => {
		const response = await GET(createEvent("GET", { user: null }));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(getWebhookMock).not.toHaveBeenCalled();
	});

	it("gets the webhook for the current user and forwards the cookie", async () => {
		const result = { webhook: webhookFixture() };
		getWebhookMock.mockResolvedValue(result);
		const event = createEvent("GET", { cookie: "auth.session_token=token-1" });

		const response = await GET(event);

		expect(getWebhookMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=token-1", "user-1");
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("saves webhooks with the current user id, trimmed URL, and supported events", async () => {
		const result = webhookFixture();
		saveWebhookMock.mockResolvedValue(result);
		const event = createEvent("POST", {
			body: {
				user_id: "different-user",
				url: "  https://example.com/webhook  ",
				events_active: ["job.started", "job.succeeded"]
			},
			cookie: "auth.session_token=token-1"
		});

		const response = await POST(event);

		expect(saveWebhookMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=token-1", {
			userId: "user-1",
			url: "https://example.com/webhook",
			eventsActive: ["job.started", "job.succeeded"]
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("allows empty webhook event lists", async () => {
		const result = webhookFixture({ events_active: [] });
		saveWebhookMock.mockResolvedValue(result);
		const event = createEvent("POST", {
			body: { url: "https://example.com/webhook", events_active: [] }
		});

		const response = await POST(event);

		expect(saveWebhookMock).toHaveBeenCalledWith(event.fetch, null, {
			userId: "user-1",
			url: "https://example.com/webhook",
			eventsActive: []
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("treats omitted webhook event lists as empty", async () => {
		const result = webhookFixture({ events_active: [] });
		saveWebhookMock.mockResolvedValue(result);
		const event = createEvent("POST", {
			body: { url: "https://example.com/webhook" }
		});

		const response = await POST(event);

		expect(saveWebhookMock).toHaveBeenCalledWith(event.fetch, null, {
			userId: "user-1",
			url: "https://example.com/webhook",
			eventsActive: []
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it.each([
		[
			"relative URL",
			{ url: "/webhook", events_active: [] },
			"url must be an absolute http or https URL"
		],
		[
			"unsupported URL protocol",
			{ url: "ftp://example.com/webhook", events_active: [] },
			"url must be an absolute http or https URL"
		],
		[
			"non-array events",
			{ url: "https://example.com/webhook", events_active: "job.started" },
			"events_active must be an array"
		],
		[
			"unsupported events",
			{ url: "https://example.com/webhook", events_active: ["job.started", "job.queued"] },
			"events_active contains unsupported events"
		]
	])("returns 400 for invalid %s", async (_name, body, error) => {
		const response = await POST(createEvent("POST", { body }));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error });
		expect(saveWebhookMock).not.toHaveBeenCalled();
	});

	it("deletes the webhook for the current user", async () => {
		deleteWebhookMock.mockResolvedValue({ deleted_id: "webhook-1", deleted_count: 1 });
		const event = createEvent("DELETE", { cookie: "auth.session_token=token-1" });

		const response = await DELETE(event);

		expect(deleteWebhookMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=token-1", "user-1");
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual({ deleted_id: "webhook-1", deleted_count: 1 });
	});

	it("normalizes backend auth service errors", async () => {
		getWebhookMock.mockRejectedValue(new AuthApiError(503, "Authentication service unavailable"));

		const response = await GET(createEvent("GET"));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
