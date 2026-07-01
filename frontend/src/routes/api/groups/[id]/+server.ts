import { json, type RequestHandler } from '@sveltejs/kit';
import { deleteGroup, getGroup, updateGroup } from '$lib/server/rbac';
import {
	cookieHeader,
	permittedRbacJSON,
	rbacAPIErrorResponse,
	readUpdateGroupInput,
	requireLocalPermission
} from '../../rbac/api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load group';
const UPDATE_ERROR_FALLBACK = 'Failed to update group';
const DELETE_ERROR_FALLBACK = 'Failed to delete group';

export const GET: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'group.view', LOAD_ERROR_FALLBACK, () =>
		getGroup(fetch, cookieHeader(request), params.id!)
	);

export const PATCH: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireLocalPermission(locals, 'group.update');
	if (authError) return authError;

	const input = await readUpdateGroupInput(request);
	if (input instanceof Response) return input;

	try {
		return json(await updateGroup(fetch, cookieHeader(request), params.id!, input));
	} catch (error) {
		return rbacAPIErrorResponse(error, UPDATE_ERROR_FALLBACK);
	}
};

export const DELETE: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'group.delete', DELETE_ERROR_FALLBACK, () =>
		deleteGroup(fetch, cookieHeader(request), params.id!)
	);
