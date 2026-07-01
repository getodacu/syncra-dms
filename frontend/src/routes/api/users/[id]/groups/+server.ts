import { json, type RequestHandler } from '@sveltejs/kit';
import { assignUserGroup } from '$lib/server/rbac';
import {
	cookieHeader,
	rbacAPIErrorResponse,
	readGroupIdInput,
	requireLocalPermission
} from '../../../rbac/api.server';

const ASSIGN_ERROR_FALLBACK = 'Failed to assign user group';

export const POST: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireLocalPermission(locals, 'user.assign_group');
	if (authError) return authError;

	const groupId = await readGroupIdInput(request);
	if (groupId instanceof Response) return groupId;

	try {
		return json(await assignUserGroup(fetch, cookieHeader(request), params.id!, groupId), {
			status: 201
		});
	} catch (error) {
		return rbacAPIErrorResponse(error, ASSIGN_ERROR_FALLBACK);
	}
};
