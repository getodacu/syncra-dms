import { json, type RequestHandler } from '@sveltejs/kit';
import { updateDocumentFolder } from '$lib/server/documents';
import {
	cookieHeader,
	documentAPIErrorResponse,
	jsonError,
	readDocumentFolderInput,
	requireAuthenticatedUser
} from '../../documents/api.server';

const UPDATE_ERROR_FALLBACK = 'Failed to update document folder';

export const PATCH: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;
	if (!params.id) return jsonError(400, 'invalid document folder id');

	const input = await readDocumentFolderInput(request);
	if (input instanceof Response) return input;

	try {
		return json(await updateDocumentFolder(fetch, cookieHeader(request), params.id, input));
	} catch (error) {
		return documentAPIErrorResponse(error, UPDATE_ERROR_FALLBACK);
	}
};
