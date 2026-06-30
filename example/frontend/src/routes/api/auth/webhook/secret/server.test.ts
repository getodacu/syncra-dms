import { beforeEach, describe, expect, it, vi } from "vitest";

import { AuthApiError } from "$lib/server/auth";
import { POST } from "./+server";
import type { RequestEvent } from "./$types";

const { regenerateWebhookSecretMock, AuthApiErrorMock } = vi.hoisted(() => {
	class MockAuthApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "AuthApiError";
			this.status = status;
		}
	}

	return {
		regenerateWebhookSecretMock: vi.fn(),
		AuthApiErrorMock: MockAuthApiError
	};
});

vi.mock("$lib/server/auth", () => ({
	regenerateWebhookSecret: regenerateWebhookSecretMock,
	AuthApiError: AuthApiErrorMock,
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock
}));

function webhookFixture() {
	return {
		id: "webhook-1",
		user_id: "user-1",
		url: "https://example.com/webhook",
		events_active: ["job.started"],
		secret: "whsec_new-secret",
		created_at: "2026-06-09T00:00:00Z",
		updated_at: "2026-06-09T00:00:00Z"
	};
}

function createEvent(
	options: {
		user?: unknown;
		cookie?: string;
	} = {}
) {
	const headers = new Headers();
	if (options.cookie) headers.set("cookie", options.cookie);

	const request = new Request("http://localhost/api/auth/webhook/secret", {
		method: "POST",
		headers
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

describe("webhook secret auth proxy endpoint", () => {
	beforeEach(() => {
		regenerateWebhookSecretMock.mockReset();
	});

	it("returns 401 for unauthenticated secret regeneration requests", async () => {
		const response = await POST(createEvent({ user: null }));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(regenerateWebhookSecretMock).not.toHaveBeenCalled();
	});

	it("regenerates the webhook secret for the current user and forwards the cookie", async () => {
		const result = webhookFixture();
		regenerateWebhookSecretMock.mockResolvedValue(result);
		const event = createEvent({ cookie: "auth.session_token=token-1" });

		const response = await POST(event);

		expect(regenerateWebhookSecretMock).toHaveBeenCalledWith(
			event.fetch,
			"auth.session_token=token-1",
			"user-1"
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("normalizes backend auth service errors", async () => {
		regenerateWebhookSecretMock.mockRejectedValue(
			new AuthApiError(503, "Authentication service unavailable")
		);

		const response = await POST(createEvent());

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
