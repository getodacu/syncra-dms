import type { RequestHandler } from '@sveltejs/kit';
import { listPermissionCategories } from '$lib/server/rbac';
import { cookieHeader, permittedRbacJSON } from '../../rbac/api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load permission categories';

export const GET: RequestHandler = async ({ fetch, locals, request }) =>
	permittedRbacJSON(locals, 'role.view', LOAD_ERROR_FALLBACK, () =>
		listPermissionCategories(fetch, cookieHeader(request))
	);
