import { json } from "@sveltejs/kit";

import { adjustAdminUserBalance } from "$lib/server/admin";
import type { RequestHandler } from "./$types";
import {
	adminApiErrorResponse,
	adminAuthError,
	readJsonObject,
	rejectUnknownKeys
} from "../../../admin-route-utils";

const BALANCE_ADJUSTMENT_KEYS = new Set(["credits_delta"]);

export const POST: RequestHandler = async ({ request, params, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	const parsed = await readJsonObject(request, "invalid balance adjustment payload");
	if (parsed.error) return parsed.error;

	const unknown = rejectUnknownKeys(parsed.value, BALANCE_ADJUSTMENT_KEYS, "invalid balance adjustment payload");
	if (unknown) return unknown;

	const creditsDelta = parsed.value.credits_delta;
	if (typeof creditsDelta !== "number" || !Number.isSafeInteger(creditsDelta) || creditsDelta === 0) {
		return json({ error: "invalid balance adjustment payload" }, { status: 400 });
	}

	try {
		const result = await adjustAdminUserBalance(
			fetch,
			request.headers.get("cookie"),
			params.id,
			creditsDelta
		);
		return json(result);
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};
