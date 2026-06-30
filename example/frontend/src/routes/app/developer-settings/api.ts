import { publicApiErrorMessage } from "$lib/client/api-errors";

type ClientFetch = typeof fetch;

const WEBHOOK_EVENTS: readonly WebhookEvent[] = [
	"job.started",
	"job.failed",
	"job.succeeded"
];

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
	name: string;
	expires_at?: string;
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
	url: string;
	events_active: WebhookEvent[];
};

export type SaveWebhookResponse = WebhookResponse;

export type RegenerateWebhookSecretResponse = WebhookResponse;

export type DeleteWebhookResponse = {
	deleted_id: string;
	deleted_count: number;
};

export async function fetchAPIKeys(fetchFn: ClientFetch): Promise<APIKeyListResponse> {
	const response = await fetchFn("/api/auth/apikeys", { method: "GET" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to load API keys"));
	}
	if (!isAPIKeyListResponse(json)) {
		throw new Error("Invalid API key list response");
	}

	return json;
}

export async function createAPIKey(
	fetchFn: ClientFetch,
	input: CreateAPIKeyInput
): Promise<APIKeyResponse> {
	const response = await fetchFn("/api/auth/apikeys", {
		method: "POST",
		headers: { "content-type": "application/json" },
		body: JSON.stringify(input)
	});
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to create API key"));
	}
	if (!isAPIKeyResponse(json)) {
		throw new Error("Invalid API key response");
	}

	return json;
}

export async function deleteAPIKey(
	fetchFn: ClientFetch,
	apiKeyId: string
): Promise<DeleteAPIKeyResponse> {
	const params = new URLSearchParams({ api_key_id: apiKeyId });
	const response = await fetchFn(`/api/auth/apikeys?${params.toString()}`, { method: "DELETE" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to delete API key"));
	}
	if (!isDeleteAPIKeyResponse(json)) {
		throw new Error("Invalid API key delete response");
	}

	return json;
}

export async function fetchWebhook(fetchFn: ClientFetch): Promise<WebhookEnvelopeResponse> {
	const response = await fetchFn("/api/auth/webhook", { method: "GET" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to load webhook"));
	}
	if (!isWebhookEnvelopeResponse(json)) {
		throw new Error("Invalid webhook envelope response");
	}

	return json;
}

export async function saveWebhook(
	fetchFn: ClientFetch,
	input: SaveWebhookInput
): Promise<SaveWebhookResponse> {
	const response = await fetchFn("/api/auth/webhook", {
		method: "POST",
		headers: { "content-type": "application/json" },
		body: JSON.stringify(input)
	});
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to save webhook"));
	}
	if (!isWebhookResponse(json)) {
		throw new Error("Invalid webhook response");
	}

	return json;
}

export async function regenerateWebhookSecret(
	fetchFn: ClientFetch
): Promise<RegenerateWebhookSecretResponse> {
	const response = await fetchFn("/api/auth/webhook/secret", { method: "POST" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(
			publicApiErrorMessage(response.status, json, "Failed to regenerate webhook secret")
		);
	}
	if (!isWebhookResponse(json)) {
		throw new Error("Invalid webhook response");
	}

	return json;
}

export async function deleteWebhook(fetchFn: ClientFetch): Promise<DeleteWebhookResponse> {
	const response = await fetchFn("/api/auth/webhook", { method: "DELETE" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to delete webhook"));
	}
	if (!isDeleteWebhookResponse(json)) {
		throw new Error("Invalid webhook delete response");
	}

	return json;
}

async function readResponseJSON(response: Response): Promise<unknown> {
	const text = await response.text();
	if (!text.trim()) return null;

	try {
		return JSON.parse(text);
	} catch {
		return null;
	}
}

function isAPIKeyListResponse(value: unknown): value is APIKeyListResponse {
	return isRecord(value) && Array.isArray(value.api_keys) && value.api_keys.every(isAPIKeyResponse);
}

function isAPIKeyResponse(value: unknown): value is APIKeyResponse {
	return (
		isRecord(value) &&
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

function isDeleteAPIKeyResponse(value: unknown): value is DeleteAPIKeyResponse {
	return (
		isRecord(value) &&
		typeof value.deleted_id === "string" &&
		typeof value.deleted_count === "number" &&
		Number.isFinite(value.deleted_count)
	);
}

function isWebhookEnvelopeResponse(value: unknown): value is WebhookEnvelopeResponse {
	return (
		isRecord(value) &&
		"webhook" in value &&
		(value.webhook === null || isWebhookResponse(value.webhook))
	);
}

function isWebhookResponse(value: unknown): value is WebhookResponse {
	return (
		isRecord(value) &&
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

function isWebhookEvent(value: unknown): value is WebhookEvent {
	return WEBHOOK_EVENTS.some((event) => event === value);
}

function isDeleteWebhookResponse(value: unknown): value is DeleteWebhookResponse {
	return (
		isRecord(value) &&
		typeof value.deleted_id === "string" &&
		typeof value.deleted_count === "number" &&
		Number.isInteger(value.deleted_count) &&
		value.deleted_count >= 0
	);
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}
