import type { RequestHandler } from '@sveltejs/kit';
import { listPermissions } from '$lib/server/rbac';
import { cookieHeader, permittedRbacJSON } from '../rbac/api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load permissions';

export const GET: RequestHandler = async ({ fetch, locals, request }) =>
	permittedRbacJSON(locals, 'role.view', LOAD_ERROR_FALLBACK, () =>
		listPermissions(fetch, cookieHeader(request))
	);
