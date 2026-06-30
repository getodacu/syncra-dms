import { json } from "@sveltejs/kit";

import { unlinkAuthAccount } from "$lib/server/auth";
import { authErrorResponse, cookieHeader } from "../../auth-proxy-utils";
import type { RequestHandler } from "./$types";

const LINKABLE_PROVIDERS = new Set(["google", "github"]);

export const DELETE: RequestHandler = async ({ params, request, fetch, locals }) => {
	if (!locals.user) return json({ error: "authentication required" }, { status: 401 });

	const provider = params.provider?.trim() ?? "";
	if (!LINKABLE_PROVIDERS.has(provider)) {
		return json({ error: "only google and github accounts can be unlinked" }, { status: 400 });
	}

	try {
		return json(await unlinkAuthAccount(fetch, cookieHeader(request), provider));
	} catch (error) {
		return authErrorResponse(error);
	}
};
