import { json } from "@sveltejs/kit";

import { createAPIKey, deleteAPIKey, isAuthApiError, listAPIKeys } from "$lib/server/auth";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

const MAX_API_KEY_NAME_CHARACTERS = 255;
const RFC3339_DATE_TIME_PATTERN =
	/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$/;

function authErrorResponse(error: unknown) {
	if (isAuthApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	throw error;
}

function cookieHeader(request: Request) {
	return request.headers.get("cookie");
}

function userId(locals: App.Locals) {
	return locals.user?.id ?? "";
}

function apiKeyName(value: unknown) {
	if (typeof value !== "string") return { ok: false as const, error: "name is required" };

	const name = value.trim();
	if (!name) return { ok: false as const, error: "name is required" };
	if (Array.from(name).length > MAX_API_KEY_NAME_CHARACTERS) {
		return { ok: false as const, error: "name must be at most 255 characters" };
	}

	return { ok: true as const, name };
}

function apiKeyExpiration(value: unknown) {
	if (value === undefined || value === null) return { ok: true as const, expiresAt: undefined };
	if (typeof value !== "string") {
		return { ok: false as const, error: "expires_at must be RFC3339" };
	}

	const expiresAt = value.trim();
	if (
		!RFC3339_DATE_TIME_PATTERN.test(expiresAt) ||
		Number.isNaN(new Date(expiresAt).getTime())
	) {
		return { ok: false as const, error: "expires_at must be RFC3339" };
	}

	return { ok: true as const, expiresAt };
}

export const GET: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) return json({ error: "authentication required" }, { status: 401 });

	try {
		return json(await listAPIKeys(fetch, cookieHeader(request), userId(locals)));
	} catch (error) {
		return authErrorResponse(error);
	}
};

export const POST: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) return json({ error: "authentication required" }, { status: 401 });

	let body: unknown;
	try {
		body = await request.json();
	} catch {
		return json({ error: "invalid API key request" }, { status: 400 });
	}

	if (typeof body !== "object" || body === null || Array.isArray(body)) {
		return json({ error: "invalid API key request" }, { status: 400 });
	}

	const nameResult = apiKeyName((body as Record<string, unknown>).name);
	if (!nameResult.ok) return json({ error: nameResult.error }, { status: 400 });
	const expirationResult = apiKeyExpiration((body as Record<string, unknown>).expires_at);
	if (!expirationResult.ok) return json({ error: expirationResult.error }, { status: 400 });

	try {
		const apiKey = await createAPIKey(fetch, cookieHeader(request), {
			userId: userId(locals),
			name: nameResult.name,
			expiresAt: expirationResult.expiresAt
		});
		return json(apiKey, { status: 201 });
	} catch (error) {
		return authErrorResponse(error);
	}
};

export const DELETE: RequestHandler = async ({ request, url, fetch, locals }) => {
	if (!locals.user) return json({ error: "authentication required" }, { status: 401 });

	const apiKeyId = url.searchParams.get("api_key_id")?.trim() ?? "";
	if (!apiKeyId) return json({ error: "api_key_id is required" }, { status: 400 });

	try {
		return json(
			await deleteAPIKey(fetch, cookieHeader(request), {
				userId: userId(locals),
				apiKeyId
			})
		);
	} catch (error) {
		return authErrorResponse(error);
	}
};
