import { json, type RequestHandler } from '@sveltejs/kit';
import { moveOrganizationUnit } from '$lib/server/organization-units';
import {
	cookieHeader,
	organizationUnitAPIErrorResponse,
	readMoveParentId,
	requireAdminUser
} from '../../api.server';

const MOVE_ERROR_FALLBACK = 'Failed to move organization unit';

export const PATCH: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireAdminUser(locals);
	if (authError) return authError;
	if (!params.id) return json({ error: 'invalid organization unit id' }, { status: 400 });

	const parentId = await readMoveParentId(request);
	if (parentId instanceof Response) return parentId;

	try {
		return json(await moveOrganizationUnit(fetch, cookieHeader(request), params.id, parentId));
	} catch (error) {
		return organizationUnitAPIErrorResponse(error, MOVE_ERROR_FALLBACK);
	}
};
