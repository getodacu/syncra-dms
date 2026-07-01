import { json, type RequestHandler } from '@sveltejs/kit';
import { addGroupUser } from '$lib/server/rbac';
import {
	cookieHeader,
	rbacAPIErrorResponse,
	readUserIdInput,
	requireLocalPermission
} from '../../../rbac/api.server';

const ASSIGN_ERROR_FALLBACK = 'Failed to add group user';

export const POST: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireLocalPermission(locals, 'group.manage_users');
	if (authError) return authError;

	const userId = await readUserIdInput(request);
	if (userId instanceof Response) return userId;

	try {
		return json(await addGroupUser(fetch, cookieHeader(request), params.id!, userId), {
			status: 201
		});
	} catch (error) {
		return rbacAPIErrorResponse(error, ASSIGN_ERROR_FALLBACK);
	}
};
