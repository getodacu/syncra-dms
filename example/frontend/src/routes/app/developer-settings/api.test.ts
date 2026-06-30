import { describe, expect, it, vi } from "vitest";

import {
	createAPIKey,
	deleteAPIKey,
	deleteWebhook,
	fetchAPIKeys,
	fetchWebhook,
	regenerateWebhookSecret,
	saveWebhook
} from "./api";
import type { SaveWebhookInput, WebhookResponse } from "./api";

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init
	});
}

function apiKeyResponse() {
	return {
		id: "api-key-1",
		user_id: "user-1",
		name: "CLI",
		key_prefix: "abc12345",
		created_at: "2026-06-09T00:00:00Z",
		updated_at: "2026-06-09T00:00:00Z"
	};
}

function webhookInput(): SaveWebhookInput {
	return {
		url: "https://example.com/syncra-webhook",
		events_active: ["job.started", "job.failed"]
	};
}

function webhookResponse(overrides: Partial<WebhookResponse> = {}): WebhookResponse {
	return {
		id: "webhook-1",
		user_id: "user-1",
		url: "https://example.com/syncra-webhook",
		events_active: ["job.started", "job.failed", "job.succeeded"],
		has_secret: true,
		created_at: "2026-06-09T00:00:00Z",
		updated_at: "2026-06-09T00:00:00Z",
		...overrides
	};
}

describe("developer settings API key client", () => {
	it("fetches API keys through the SvelteKit proxy", async () => {
		const body = { api_keys: [apiKeyResponse()] };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchAPIKeys(fetchFn)).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/auth/apikeys", { method: "GET" });
	});

	it("creates API keys through the SvelteKit proxy", async () => {
		const body = { ...apiKeyResponse(), api_key: "abc12345secret" };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body, { status: 201 }));

		await expect(createAPIKey(fetchFn, { name: "CLI" })).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/auth/apikeys", {
			method: "POST",
			headers: { "content-type": "application/json" },
			body: JSON.stringify({ name: "CLI" })
		});
	});

	it("passes optional API key expiration through the SvelteKit proxy", async () => {
		const expiresAt = "2026-06-16T23:59:59.999Z";
		const body = { ...apiKeyResponse(), api_key: "abc12345secret", expires_at: expiresAt };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body, { status: 201 }));

		await expect(createAPIKey(fetchFn, { name: "CLI", expires_at: expiresAt })).resolves.toEqual(
			body
		);
		expect(fetchFn).toHaveBeenCalledWith("/api/auth/apikeys", {
			method: "POST",
			headers: { "content-type": "application/json" },
			body: JSON.stringify({ name: "CLI", expires_at: expiresAt })
		});
	});

	it("deletes API keys through the SvelteKit proxy", async () => {
		const body = { deleted_id: "api-key/1", deleted_count: 1 };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(deleteAPIKey(fetchFn, "api-key/1")).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/auth/apikeys?api_key_id=api-key%2F1", {
			method: "DELETE"
		});
	});

	it("throws backend JSON messages for create errors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "name is required" }, { status: 400 }));

		await expect(createAPIKey(fetchFn, { name: "" })).rejects.toThrow("name is required");
	});

	it("rejects invalid API key list responses", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ api_keys: [{ id: "api-key-1" }] }));

		await expect(fetchAPIKeys(fetchFn)).rejects.toThrow("Invalid API key list response");
	});
});

describe("developer settings webhook client", () => {
	it("fetches webhooks through the SvelteKit proxy", async () => {
		const body = { webhook: webhookResponse() };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchWebhook(fetchFn)).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/auth/webhook", { method: "GET" });
	});

	it("accepts an empty webhook envelope", async () => {
		const body = { webhook: null };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchWebhook(fetchFn)).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/auth/webhook", { method: "GET" });
	});

	it("saves webhooks through the SvelteKit proxy", async () => {
		const input = webhookInput();
		const body = webhookResponse({ secret_key: "whsec_created" });
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body, { status: 201 }));

		await expect(saveWebhook(fetchFn, input)).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/auth/webhook", {
			method: "POST",
			headers: { "content-type": "application/json" },
			body: JSON.stringify(input)
		});
	});

	it("regenerates webhook secrets through the SvelteKit proxy", async () => {
		const body = webhookResponse({ secret_key: "whsec_regenerated" });
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(regenerateWebhookSecret(fetchFn)).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/auth/webhook/secret", { method: "POST" });
	});

	it("deletes webhooks through the SvelteKit proxy", async () => {
		const body = { deleted_id: "webhook-1", deleted_count: 1 };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(deleteWebhook(fetchFn)).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/auth/webhook", { method: "DELETE" });
	});

	it("throws backend JSON messages for webhook save errors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "url is required" }, { status: 400 }));

		await expect(saveWebhook(fetchFn, webhookInput())).rejects.toThrow("url is required");
	});

	it("rejects invalid webhook envelope responses", async () => {
		const invalidBodies: Array<[string, unknown]> = [
			["missing webhook", {}],
			["null body", null],
			["invalid webhook object", { webhook: { id: "webhook-1" } }]
		];

		for (const [, body] of invalidBodies) {
			const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

			await expect(fetchWebhook(fetchFn)).rejects.toThrow("Invalid webhook envelope response");
		}
	});

	it("rejects invalid webhook responses", async () => {
		const invalidBodies: Array<[string, unknown]> = [
			["missing required field", { ...webhookResponse(), id: undefined }],
			["unsupported event", { ...webhookResponse(), events_active: ["job.deleted"] }],
			["wrong secret key type", { ...webhookResponse(), secret_key: 123 }]
		];

		for (const [, body] of invalidBodies) {
			const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

			await expect(saveWebhook(fetchFn, webhookInput())).rejects.toThrow(
				"Invalid webhook response"
			);
		}
	});

	it("rejects invalid webhook delete responses", async () => {
		const invalidBodies: Array<[string, unknown]> = [
			["missing deleted id", { deleted_count: 1 }],
			["bad deleted count", { deleted_id: "webhook-1", deleted_count: "1" }]
		];

		for (const [, body] of invalidBodies) {
			const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

			await expect(deleteWebhook(fetchFn)).rejects.toThrow("Invalid webhook delete response");
		}
	});
});
