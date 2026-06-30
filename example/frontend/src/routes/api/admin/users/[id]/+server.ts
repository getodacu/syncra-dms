import { json } from "@sveltejs/kit";

import { getAdminUser, updateAdminUser } from "$lib/server/admin";
import type { RequestHandler } from "./$types";
import {
	adminApiErrorResponse,
	adminAuthError,
	readJsonObject,
	rejectUnknownKeys
} from "../../admin-route-utils";

const UPDATE_USER_KEYS = new Set(["name", "email"]);

export const GET: RequestHandler = async ({ request, params, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	try {
		const result = await getAdminUser(fetch, request.headers.get("cookie"), params.id);
		return json(result);
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};

export const PATCH: RequestHandler = async ({ request, params, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	const parsed = await readJsonObject(request, "invalid user update payload");
	if (parsed.error) return parsed.error;

	const unknown = rejectUnknownKeys(parsed.value, UPDATE_USER_KEYS, "invalid user update payload");
	if (unknown) return unknown;

	const input: { name?: string; email?: string } = {};
	if ("name" in parsed.value) {
		if (typeof parsed.value.name !== "string") {
			return json({ error: "invalid user update payload" }, { status: 400 });
		}
		input.name = parsed.value.name;
	}
	if ("email" in parsed.value) {
		if (typeof parsed.value.email !== "string") {
			return json({ error: "invalid user update payload" }, { status: 400 });
		}
		input.email = parsed.value.email;
	}

	try {
		const result = await updateAdminUser(fetch, request.headers.get("cookie"), params.id, input);
		return json(result);
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};
