import type { RequestHandler } from '@sveltejs/kit';
import { removeGroupRole } from '$lib/server/rbac';
import { cookieHeader, permittedRbacJSON } from '../../../../rbac/api.server';

const REMOVE_ERROR_FALLBACK = 'Failed to remove group role';

export const DELETE: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'group.assign_roles', REMOVE_ERROR_FALLBACK, () =>
		removeGroupRole(fetch, cookieHeader(request), params.id!, params.assignmentId!)
	);
