import { error, redirect } from "@sveltejs/kit";

import { stopAdminImpersonation } from "$lib/server/admin";
import type { RequestHandler } from "./$types";

export const POST: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user && !locals.adminUser) {
		redirect(303, "/login");
	}
	if (!locals.adminUser || locals.adminUser.role !== "admin") {
		error(403, "Admin access required");
	}

	const targetUserId = locals.impersonation?.targetUser.id;
	await stopAdminImpersonation(fetch, request.headers.get("cookie"));
	redirect(303, targetUserId ? `/admin-portal/users/${targetUserId}` : "/admin-portal");
};
