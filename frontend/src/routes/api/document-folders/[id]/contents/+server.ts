import { json, type RequestHandler } from '@sveltejs/kit';
import { getDocumentFolderContents } from '$lib/server/documents';
import {
	cookieHeader,
	documentAPIErrorResponse,
	jsonError,
	requireAuthenticatedUser
} from '../../../documents/api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load document folder contents';

export const GET: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;
	if (!params.id) return jsonError(400, 'invalid document folder id');

	try {
		return json(await getDocumentFolderContents(fetch, cookieHeader(request), params.id));
	} catch (error) {
		return documentAPIErrorResponse(error, LOAD_ERROR_FALLBACK);
	}
};
