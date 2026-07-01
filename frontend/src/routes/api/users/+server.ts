import { json, type RequestHandler } from '@sveltejs/kit';
import { createUser, listUsers } from '$lib/server/rbac';
import {
	cookieHeader,
	permittedRbacJSON,
	rbacAPIErrorResponse,
	readCreateUserInput,
	requireLocalPermission
} from '../rbac/api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load users';
const CREATE_ERROR_FALLBACK = 'Failed to create user';

export const GET: RequestHandler = async ({ fetch, locals, request }) =>
	permittedRbacJSON(locals, 'user.view', LOAD_ERROR_FALLBACK, () =>
		listUsers(fetch, cookieHeader(request))
	);

export const POST: RequestHandler = async ({ fetch, locals, request }) => {
	const authError = requireLocalPermission(locals, 'user.create');
	if (authError) return authError;

	const input = await readCreateUserInput(request);
	if (input instanceof Response) return input;

	try {
		return json(await createUser(fetch, cookieHeader(request), input), { status: 201 });
	} catch (error) {
		return rbacAPIErrorResponse(error, CREATE_ERROR_FALLBACK);
	}
};
