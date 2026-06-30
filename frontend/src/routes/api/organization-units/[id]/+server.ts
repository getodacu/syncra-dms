import { json, type RequestHandler } from '@sveltejs/kit';
import { updateOrganizationUnit } from '$lib/server/organization-units';
import {
	cookieHeader,
	organizationUnitAPIErrorResponse,
	readOrganizationUnitInput,
	requireAdminUser
} from '../api.server';

const UPDATE_ERROR_FALLBACK = 'Failed to update organization unit';

export const PATCH: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireAdminUser(locals);
	if (authError) return authError;
	if (!params.id) return json({ error: 'invalid organization unit id' }, { status: 400 });

	const input = await readOrganizationUnitInput(request);
	if (input instanceof Response) return input;

	try {
		return json(await updateOrganizationUnit(fetch, cookieHeader(request), params.id, input));
	} catch (error) {
		return organizationUnitAPIErrorResponse(error, UPDATE_ERROR_FALLBACK);
	}
};
