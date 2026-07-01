import { json, type RequestHandler } from '@sveltejs/kit';
import { deleteUser, getUser, updateUser } from '$lib/server/rbac';
import {
	cookieHeader,
	permittedRbacJSON,
	rbacAPIErrorResponse,
	readUpdateUserInput,
	requireLocalPermission
} from '../../rbac/api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load user';
const UPDATE_ERROR_FALLBACK = 'Failed to update user';
const DELETE_ERROR_FALLBACK = 'Failed to delete user';

export const GET: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'user.view', LOAD_ERROR_FALLBACK, () =>
		getUser(fetch, cookieHeader(request), params.id!)
	);

export const PATCH: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireLocalPermission(locals, 'user.update');
	if (authError) return authError;

	const input = await readUpdateUserInput(request);
	if (input instanceof Response) return input;

	try {
		return json(await updateUser(fetch, cookieHeader(request), params.id!, input));
	} catch (error) {
		return rbacAPIErrorResponse(error, UPDATE_ERROR_FALLBACK);
	}
};

export const DELETE: RequestHandler = async ({ fetch, locals, params, request }) =>
	permittedRbacJSON(locals, 'user.delete', DELETE_ERROR_FALLBACK, () =>
		deleteUser(fetch, cookieHeader(request), params.id!)
	);
