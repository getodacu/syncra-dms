import type { RequestHandler } from '@sveltejs/kit';
import { getMyPermissions } from '$lib/server/rbac';
import { authenticatedRbacJSON, cookieHeader } from '../../rbac/api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load permissions';

export const GET: RequestHandler = async ({ fetch, locals, request }) =>
	authenticatedRbacJSON(locals, LOAD_ERROR_FALLBACK, () =>
		getMyPermissions(fetch, cookieHeader(request))
	);
