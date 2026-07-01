import type { RequestHandler } from '@sveltejs/kit';
import { removeGroupUser } from '$lib/server/rbac';
import { cookieHeader, permittedRbacJSON } from '../../../../rbac/api.server';

const REMOVE_ERROR_FALLBACK = 'Failed to remove group user';

export const DELETE: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'group.manage_users', REMOVE_ERROR_FALLBACK, () =>
		removeGroupUser(fetch, cookieHeader(request), params.id!, params.userId!)
	);
