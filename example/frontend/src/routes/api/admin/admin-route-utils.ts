import { json } from "@sveltejs/kit";

import { isAdminApiError } from "$lib/server/admin";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";

export function adminAuthError(locals: App.Locals) {
	if (!locals.user && !locals.adminUser) return json({ error: "authentication required" }, { status: 401 });
	if (!locals.adminUser || locals.adminUser.role !== "admin") {
		return json({ error: "admin access required" }, { status: 403 });
	}
	return null;
}

export function adminApiErrorResponse(error: unknown) {
	if (isAdminApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}
	throw error;
}

export function optionalQuery(url: URL, key: string) {
	const value = url.searchParams.get(key);
	return value && value.trim() ? value.trim() : undefined;
}

type OptionalBooleanQueryResult =
	| { ok: true; value: boolean | undefined }
	| { ok: false; error: Response };

export function optionalBooleanQuery(url: URL, key: string): OptionalBooleanQueryResult {
	const value = optionalQuery(url, key);
	if (value === undefined) return { ok: true, value: undefined };
	if (value === "true") return { ok: true, value: true };
	if (value === "false") return { ok: true, value: false };

	return { ok: false, error: json({ error: `invalid ${key}` }, { status: 400 }) };
}

export async function readJsonObject(request: Request, errorMessage: string) {
	let body: unknown;
	try {
		body = await request.json();
	} catch {
		return { error: json({ error: errorMessage }, { status: 400 }) };
	}
	if (typeof body !== "object" || body === null || Array.isArray(body)) {
		return { error: json({ error: errorMessage }, { status: 400 }) };
	}
	return { value: body as Record<string, unknown> };
}

export function rejectUnknownKeys(
	value: Record<string, unknown>,
	allowed: ReadonlySet<string>,
	errorMessage: string
) {
	for (const key of Object.keys(value)) {
		if (!allowed.has(key)) {
			return json({ error: errorMessage }, { status: 400 });
		}
	}
	return null;
}
