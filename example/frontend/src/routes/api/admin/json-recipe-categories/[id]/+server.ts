import { json } from "@sveltejs/kit";

import {
	deleteAdminJSONRecipeCategory,
	getAdminJSONRecipeCategory,
	updateAdminJSONRecipeCategory
} from "$lib/server/admin";
import type { RequestHandler } from "./$types";
import {
	adminApiErrorResponse,
	adminAuthError,
	readJsonObject
} from "../../admin-route-utils";
import { _parseJSONRecipeCategoryPayload } from "../+server";

export const GET: RequestHandler = async ({ request, params, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	try {
		const result = await getAdminJSONRecipeCategory(fetch, request.headers.get("cookie"), params.id);
		return json(result);
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};

export const PUT: RequestHandler = async ({ request, params, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	const parsed = await readJsonObject(request, "invalid JSON recipe category payload");
	if (parsed.error) return parsed.error;

	const payload = _parseJSONRecipeCategoryPayload(parsed.value);
	if (payload.error) return payload.error;

	try {
		const result = await updateAdminJSONRecipeCategory(
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
		await deleteAdminJSONRecipeCategory(fetch, request.headers.get("cookie"), params.id);
		return new Response(null, { status: 204 });
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};
