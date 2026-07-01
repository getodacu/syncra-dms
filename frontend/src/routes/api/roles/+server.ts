import { json, type RequestHandler } from '@sveltejs/kit';
import { createRole, listRoles } from '$lib/server/rbac';
import {
	cookieHeader,
	permittedRbacJSON,
	rbacAPIErrorResponse,
	readCreateRoleInput,
	requireLocalPermission
} from '../rbac/api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load roles';
const CREATE_ERROR_FALLBACK = 'Failed to create role';

export const GET: RequestHandler = async ({ fetch, locals, request }) =>
	permittedRbacJSON(locals, 'role.view', LOAD_ERROR_FALLBACK, () =>
		listRoles(fetch, cookieHeader(request))
	);

export const POST: RequestHandler = async ({ fetch, locals, request }) => {
	const authError = requireLocalPermission(locals, 'role.create');
	if (authError) return authError;

	const input = await readCreateRoleInput(request);
	if (input instanceof Response) return input;

	try {
		return json(await createRole(fetch, cookieHeader(request), input), { status: 201 });
	} catch (error) {
		return rbacAPIErrorResponse(error, CREATE_ERROR_FALLBACK);
	}
};
