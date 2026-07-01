import { json, type RequestHandler } from '@sveltejs/kit';
import { setUserPrimaryOrganizationUnit } from '$lib/server/rbac';
import {
	cookieHeader,
	rbacAPIErrorResponse,
	readPrimaryOrganizationUnitInput,
	requireLocalPermission
} from '../../../rbac/api.server';

const UPDATE_ERROR_FALLBACK = 'Failed to update primary organization unit';

export const POST: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireLocalPermission(locals, 'user.assign_unit');
	if (authError) return authError;

	const organizationUnitId = await readPrimaryOrganizationUnitInput(request);
	if (organizationUnitId instanceof Response) return organizationUnitId;

	try {
		return json(
			await setUserPrimaryOrganizationUnit(
				fetch,
				cookieHeader(request),
				params.id!,
				organizationUnitId
			)
		);
	} catch (error) {
		return rbacAPIErrorResponse(error, UPDATE_ERROR_FALLBACK);
	}
};
