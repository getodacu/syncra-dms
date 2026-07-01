import type { RequestHandler } from '@sveltejs/kit';
import { removeUserGroup } from '$lib/server/rbac';
import { cookieHeader, permittedRbacJSON } from '../../../../rbac/api.server';

const REMOVE_ERROR_FALLBACK = 'Failed to remove user group';

export const DELETE: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'user.assign_group', REMOVE_ERROR_FALLBACK, () =>
		removeUserGroup(fetch, cookieHeader(request), params.id!, params.groupId!)
	);
