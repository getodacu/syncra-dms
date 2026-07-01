import type { RequestHandler } from '@sveltejs/kit';
import { deactivateUser } from '$lib/server/rbac';
import { cookieHeader, permittedRbacJSON } from '../../../rbac/api.server';

const UPDATE_ERROR_FALLBACK = 'Failed to deactivate user';

export const POST: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'user.update', UPDATE_ERROR_FALLBACK, () =>
		deactivateUser(fetch, cookieHeader(request), params.id!)
	);
