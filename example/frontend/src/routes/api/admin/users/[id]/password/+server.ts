import { json } from "@sveltejs/kit";

import { resetAdminUserPassword } from "$lib/server/admin";
import type { RequestHandler } from "./$types";
import { adminApiErrorResponse, adminAuthError, readJsonObject } from "../../../admin-route-utils";

export const POST: RequestHandler = async ({ request, params, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	const parsed = await readJsonObject(request, "invalid password reset payload");
	if (parsed.error) return parsed.error;
	if (typeof parsed.value.password !== "string") {
		return json({ error: "invalid password reset payload" }, { status: 400 });
	}

	try {
		const result = await resetAdminUserPassword(
			fetch,
			request.headers.get("cookie"),
			params.id,
			parsed.value.password
		);
		return json(result);
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};
