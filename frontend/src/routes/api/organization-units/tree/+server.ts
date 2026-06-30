import { json, type RequestHandler } from '@sveltejs/kit';
import { getOrganizationUnitTree } from '$lib/server/organization-units';
import {
	cookieHeader,
	organizationUnitAPIErrorResponse,
	requireAuthenticatedUser
} from '../api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load organization units';

export const GET: RequestHandler = async ({ fetch, locals, request }) => {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;

	try {
		return json(await getOrganizationUnitTree(fetch, cookieHeader(request)));
	} catch (error) {
		return organizationUnitAPIErrorResponse(error, LOAD_ERROR_FALLBACK);
	}
};
