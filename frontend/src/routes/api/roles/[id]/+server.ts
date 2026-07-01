import { json, type RequestHandler } from '@sveltejs/kit';
import { deleteRole, getRole, updateRole } from '$lib/server/rbac';
import {
	cookieHeader,
	permittedRbacJSON,
	rbacAPIErrorResponse,
	readUpdateRoleInput,
	requireLocalPermission
} from '../../rbac/api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load role';
const UPDATE_ERROR_FALLBACK = 'Failed to update role';
const DELETE_ERROR_FALLBACK = 'Failed to delete role';

export const GET: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'role.view', LOAD_ERROR_FALLBACK, () =>
		getRole(fetch, cookieHeader(request), params.id!)
	);

export const PATCH: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireLocalPermission(locals, 'role.update');
	if (authError) return authError;

	const input = await readUpdateRoleInput(request);
	if (input instanceof Response) return input;

	try {
		return json(await updateRole(fetch, cookieHeader(request), params.id!, input));
	} catch (error) {
		return rbacAPIErrorResponse(error, UPDATE_ERROR_FALLBACK);
	}
};

export const DELETE: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'role.delete', DELETE_ERROR_FALLBACK, () =>
		deleteRole(fetch, cookieHeader(request), params.id!)
	);
