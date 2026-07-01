import type { RequestHandler } from '@sveltejs/kit';
import { suspendUser } from '$lib/server/rbac';
import { cookieHeader, permittedRbacJSON } from '../../../rbac/api.server';

const UPDATE_ERROR_FALLBACK = 'Failed to suspend user';

export const POST: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'user.suspend', UPDATE_ERROR_FALLBACK, () =>
		suspendUser(fetch, cookieHeader(request), params.id!)
	);
