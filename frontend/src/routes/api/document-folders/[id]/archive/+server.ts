import { json, type RequestHandler } from '@sveltejs/kit';
import { archiveDocumentFolder } from '$lib/server/documents';
import {
	cookieHeader,
	documentAPIErrorResponse,
	jsonError,
	requireAuthenticatedUser
} from '../../../documents/api.server';

const ARCHIVE_ERROR_FALLBACK = 'Failed to archive document folder';

export const POST: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;
	if (!params.id) return jsonError(400, 'invalid document folder id');

	try {
		return json(await archiveDocumentFolder(fetch, cookieHeader(request), params.id));
	} catch (error) {
		return documentAPIErrorResponse(error, ARCHIVE_ERROR_FALLBACK);
	}
};
