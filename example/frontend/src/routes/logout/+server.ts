import { redirect, type RequestHandler } from "@sveltejs/kit";

import { clearSessionCookie, isAuthApiError, signOut } from "$lib/server/auth";
import { safeError } from "$lib/server/logging";

export const POST: RequestHandler = async ({ cookies, fetch, locals, request }) => {
	try {
		await signOut(fetch, request.headers.get("cookie"));
	} catch (error) {
		if (!isAuthApiError(error)) throw error;
		locals.logger.error("auth.signout_failed", {
			error: safeError(error),
			user_id: locals.user?.id ?? null
		});
	} finally {
		clearSessionCookie(cookies);
		locals.session = null;
		locals.user = null;
	}

	redirect(303, "/login");
};
