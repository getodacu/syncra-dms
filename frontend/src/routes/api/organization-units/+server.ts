import { json, type RequestHandler } from '@sveltejs/kit';
import { createOrganizationUnit } from '$lib/server/organization-units';
import {
	cookieHeader,
	organizationUnitAPIErrorResponse,
	readOrganizationUnitInput,
	requireAdminUser
} from './api.server';

const CREATE_ERROR_FALLBACK = 'Failed to create organization unit';

export const POST: RequestHandler = async ({ fetch, locals, request }) => {
	const authError = requireAdminUser(locals);
	if (authError) return authError;

	const input = await readOrganizationUnitInput(request);
	if (input instanceof Response) return input;

	try {
		return json(await createOrganizationUnit(fetch, cookieHeader(request), input), { status: 201 });
	} catch (error) {
		return organizationUnitAPIErrorResponse(error, CREATE_ERROR_FALLBACK);
	}
};
