import { json } from "@sveltejs/kit";

import { createAdminJSONRecipe, listAdminJSONRecipes } from "$lib/server/admin";
import type { RequestHandler } from "./$types";
import {
	adminApiErrorResponse,
	adminAuthError,
	optionalQuery,
	readJsonObject,
	rejectUnknownKeys
} from "../admin-route-utils";

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

export const GET: RequestHandler = async ({ url, request, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	try {
		const result = await listAdminJSONRecipes(fetch, request.headers.get("cookie"), {
			cursor: optionalQuery(url, "cursor"),
			size: optionalQuery(url, "size"),
			sort: optionalQuery(url, "sort") as "asc" | "desc" | undefined
		});
		return json(result);
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};

export const POST: RequestHandler = async ({ request, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	const parsed = await readJsonObject(request, "invalid JSON recipe payload");
	if (parsed.error) return parsed.error;

	const payload = parseJSONRecipePayload(parsed.value);
	if (payload.error) return payload.error;

	try {
		const result = await createAdminJSONRecipe(fetch, request.headers.get("cookie"), payload.input);
		return json(result, { status: 201 });
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};
