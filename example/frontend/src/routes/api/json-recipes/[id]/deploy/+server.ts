import { json } from "@sveltejs/kit";

import { deployJsonRecipe, isSchemaApiError } from "$lib/server/schemas";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

function isJsonObject(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

async function readDeployPayload(request: Request) {
	let body: unknown;
	try {
		body = await request.json();
	} catch {
		return { error: json({ error: "invalid JSON recipe deploy payload" }, { status: 400 }) };
	}
	if (!isJsonObject(body) || typeof body.user_id !== "string") {
		return { error: json({ error: "invalid JSON recipe deploy payload" }, { status: 400 }) };
	}
	for (const key of Object.keys(body)) {
		if (key !== "user_id") {
			return { error: json({ error: "invalid JSON recipe deploy payload" }, { status: 400 }) };
		}
	}
	return { userId: body.user_id };
}

export const POST: RequestHandler = async ({ request, params, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const payload = await readDeployPayload(request);
	if (payload.error) return payload.error;

	if (locals.adminUser?.role !== "admin" && payload.userId !== locals.user.id) {
		return json({ error: "cannot deploy recipe for another user" }, { status: 403 });
	}

	try {
		const result = await deployJsonRecipe(fetch, params.id, { userId: payload.userId });
		return json(result, { status: 201 });
	} catch (error) {
		if (isSchemaApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}
		throw error;
	}
};
