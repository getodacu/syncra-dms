import { json, type RequestHandler } from '@sveltejs/kit';
import { moveDocumentFolder } from '$lib/server/documents';
import {
	cookieHeader,
	documentAPIErrorResponse,
	jsonError,
	readMoveParentId,
	requireAuthenticatedUser
} from '../../../documents/api.server';

const MOVE_ERROR_FALLBACK = 'Failed to move document folder';

export const PATCH: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;
	if (!params.id) return jsonError(400, 'invalid document folder id');

	const parentId = await readMoveParentId(request);
	if (parentId instanceof Response) return parentId;

	try {
		return json(await moveDocumentFolder(fetch, cookieHeader(request), params.id, parentId));
	} catch (error) {
		return documentAPIErrorResponse(error, MOVE_ERROR_FALLBACK);
	}
};
