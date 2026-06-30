import { json } from "@sveltejs/kit";

import {
	deleteAdminJSONRecipe,
	getAdminJSONRecipe,
	updateAdminJSONRecipe
} from "$lib/server/admin";
import type { RequestHandler } from "./$types";
import {
	adminApiErrorResponse,
	adminAuthError,
	readJsonObject,
	rejectUnknownKeys
} from "../../admin-route-utils";

const JSON_RECIPE_KEYS = new Set(["title", "description", "json", "category_id"]);

function isJsonObject(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function parseJSONRecipePayload(value: Record<string, unknown>) {
	const unknown = rejectUnknownKeys(value, JSON_RECIPE_KEYS, "invalid JSON recipe payload");
	if (unknown) return { error: unknown };
	if (
		typeof value.title !== "string" ||
		typeof value.description !== "string" ||
		!isJsonObject(value.json) ||
		("category_id" in value && value.category_id !== null && typeof value.category_id !== "string")
	) {
		return { error: json({ error: "invalid JSON recipe payload" }, { status: 400 }) };
	}
	return {
		input: {
			title: value.title,
			description: value.description,
			json: value.json,
			...("category_id" in value ? { category_id: value.category_id as string | null } : {})
		}
	};
}

export const GET: RequestHandler = async ({ request, params, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	try {
		const result = await getAdminJSONRecipe(fetch, request.headers.get("cookie"), params.id);
		return json(result);
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};

export const PUT: RequestHandler = async ({ request, params, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	const parsed = await readJsonObject(request, "invalid JSON recipe payload");
	if (parsed.error) return parsed.error;

	const payload = parseJSONRecipePayload(parsed.value);
	if (payload.error) return payload.error;

	try {
		const result = await updateAdminJSONRecipe(
			fetch,
			request.headers.get("cookie"),
			params.id,
			payload.input
		);
		return json(result);
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};

export const DELETE: RequestHandler = async ({ request, params, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	try {
		await deleteAdminJSONRecipe(fetch, request.headers.get("cookie"), params.id);
		return new Response(null, { status: 204 });
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};
