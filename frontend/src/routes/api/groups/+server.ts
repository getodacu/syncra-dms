import { json, type RequestHandler } from '@sveltejs/kit';
import { createGroup, listGroups } from '$lib/server/rbac';
import {
	cookieHeader,
	permittedRbacJSON,
	rbacAPIErrorResponse,
	readCreateGroupInput,
	requireLocalPermission
} from '../rbac/api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load groups';
const CREATE_ERROR_FALLBACK = 'Failed to create group';

export const GET: RequestHandler = async ({ fetch, locals, request }) =>
	permittedRbacJSON(locals, 'group.view', LOAD_ERROR_FALLBACK, () =>
		listGroups(fetch, cookieHeader(request))
	);

export const POST: RequestHandler = async ({ fetch, locals, request }) => {
	const authError = requireLocalPermission(locals, 'group.create');
	if (authError) return authError;

	const input = await readCreateGroupInput(request);
	if (input instanceof Response) return input;

	try {
		return json(await createGroup(fetch, cookieHeader(request), input), { status: 201 });
	} catch (error) {
		return rbacAPIErrorResponse(error, CREATE_ERROR_FALLBACK);
	}
};
