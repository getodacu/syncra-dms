import type { RequestHandler } from '@sveltejs/kit';
import { removeUserRole } from '$lib/server/rbac';
import { cookieHeader, permittedRbacJSON } from '../../../../rbac/api.server';

const REMOVE_ERROR_FALLBACK = 'Failed to remove user role';

export const DELETE: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'user.assign_role', REMOVE_ERROR_FALLBACK, () =>
		removeUserRole(fetch, cookieHeader(request), params.id!, params.assignmentId!)
	);
