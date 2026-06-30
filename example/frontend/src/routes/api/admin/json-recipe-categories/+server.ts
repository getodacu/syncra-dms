import { json } from "@sveltejs/kit";

import { createAdminJSONRecipeCategory, listAdminJSONRecipeCategories } from "$lib/server/admin";
import type { RequestHandler } from "./$types";
import {
	adminApiErrorResponse,
	adminAuthError,
	readJsonObject,
	rejectUnknownKeys
} from "../admin-route-utils";

const JSON_RECIPE_CATEGORY_KEYS = new Set(["title"]);
const JSON_RECIPE_CATEGORY_TITLE_KEYS = new Set(["en", "ro"]);

function isJsonObject(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

export function _parseJSONRecipeCategoryPayload(value: Record<string, unknown>) {
	const unknown = rejectUnknownKeys(value, JSON_RECIPE_CATEGORY_KEYS, "invalid JSON recipe category payload");
	if (unknown) return { error: unknown };
	if (!isJsonObject(value.title)) {
		return { error: json({ error: "invalid JSON recipe category payload" }, { status: 400 }) };
	}
	const unknownTitle = rejectUnknownKeys(
		value.title,
		JSON_RECIPE_CATEGORY_TITLE_KEYS,
		"invalid JSON recipe category payload"
	);
	if (unknownTitle) return { error: unknownTitle };
	if (typeof value.title.en !== "string" || typeof value.title.ro !== "string") {
		return { error: json({ error: "invalid JSON recipe category payload" }, { status: 400 }) };
	}
	return {
		input: {
			title: {
				en: value.title.en,
				ro: value.title.ro
			}
		}
	};
}

export const GET: RequestHandler = async ({ request, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	try {
		const result = await listAdminJSONRecipeCategories(fetch, request.headers.get("cookie"));
		return json(result);
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};

export const POST: RequestHandler = async ({ request, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	const parsed = await readJsonObject(request, "invalid JSON recipe category payload");
	if (parsed.error) return parsed.error;

	const payload = _parseJSONRecipeCategoryPayload(parsed.value);
	if (payload.error) return payload.error;

	try {
		const result = await createAdminJSONRecipeCategory(fetch, request.headers.get("cookie"), payload.input);
		return json(result, { status: 201 });
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};
