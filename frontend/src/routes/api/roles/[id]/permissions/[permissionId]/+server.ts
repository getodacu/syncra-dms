import type { RequestHandler } from '@sveltejs/kit';
import { removeRolePermission } from '$lib/server/rbac';
import { cookieHeader, permittedRbacJSON } from '../../../../rbac/api.server';

const REMOVE_ERROR_FALLBACK = 'Failed to remove role permission';

export const DELETE: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'role.assign_permissions', REMOVE_ERROR_FALLBACK, () =>
		removeRolePermission(fetch, cookieHeader(request), params.id!, params.permissionId!)
	);
