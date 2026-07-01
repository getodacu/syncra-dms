import { json, type RequestHandler } from '@sveltejs/kit';
import { assignUserRole } from '$lib/server/rbac';
import {
	cookieHeader,
	rbacAPIErrorResponse,
	readScopedRoleAssignmentInput,
	requireLocalPermission
} from '../../../rbac/api.server';

const ASSIGN_ERROR_FALLBACK = 'Failed to assign user role';

export const POST: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireLocalPermission(locals, 'user.assign_role');
	if (authError) return authError;

	const input = await readScopedRoleAssignmentInput(request);
	if (input instanceof Response) return input;

	try {
		return json(await assignUserRole(fetch, cookieHeader(request), params.id!, input), {
			status: 201
		});
	} catch (error) {
		return rbacAPIErrorResponse(error, ASSIGN_ERROR_FALLBACK);
	}
};
