import { json, type RequestHandler } from '@sveltejs/kit';
import { assignRolePermission, listRolePermissions } from '$lib/server/rbac';
import {
	cookieHeader,
	permittedRbacJSON,
	rbacAPIErrorResponse,
	readPermissionIdInput,
	requireLocalPermission
} from '../../../rbac/api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load role permissions';
const ASSIGN_ERROR_FALLBACK = 'Failed to assign role permission';

export const GET: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'role.view', LOAD_ERROR_FALLBACK, () =>
		listRolePermissions(fetch, cookieHeader(request), params.id!)
	);

export const POST: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireLocalPermission(locals, 'role.assign_permissions');
	if (authError) return authError;

	const permissionId = await readPermissionIdInput(request);
	if (permissionId instanceof Response) return permissionId;

	try {
		return json(await assignRolePermission(fetch, cookieHeader(request), params.id!, permissionId), {
			status: 201
		});
	} catch (error) {
		return rbacAPIErrorResponse(error, ASSIGN_ERROR_FALLBACK);
	}
};
