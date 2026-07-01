import type { RequestHandler } from '@sveltejs/kit';
import { activateUser } from '$lib/server/rbac';
import { cookieHeader, permittedRbacJSON } from '../../../rbac/api.server';

const UPDATE_ERROR_FALLBACK = 'Failed to activate user';

export const POST: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'user.activate', UPDATE_ERROR_FALLBACK, () =>
		activateUser(fetch, cookieHeader(request), params.id!)
	);
