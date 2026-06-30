import { json, type RequestHandler } from '@sveltejs/kit';
import { archiveOrganizationUnit } from '$lib/server/organization-units';
import {
	cookieHeader,
	organizationUnitAPIErrorResponse,
	requireAdminUser
} from '../../api.server';

const ARCHIVE_ERROR_FALLBACK = 'Failed to archive organization unit';

export const POST: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireAdminUser(locals);
	if (authError) return authError;
	if (!params.id) return json({ error: 'invalid organization unit id' }, { status: 400 });

	try {
		return json(await archiveOrganizationUnit(fetch, cookieHeader(request), params.id));
	} catch (error) {
		return organizationUnitAPIErrorResponse(error, ARCHIVE_ERROR_FALLBACK);
	}
};
