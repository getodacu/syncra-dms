import type { RequestHandler } from '@sveltejs/kit';
import { getMe } from '$lib/server/rbac';
import { authenticatedRbacJSON, cookieHeader } from '../rbac/api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load profile';

export const GET: RequestHandler = async ({ fetch, locals, request }) =>
	authenticatedRbacJSON(locals, LOAD_ERROR_FALLBACK, () => getMe(fetch, cookieHeader(request)));
