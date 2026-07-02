import { json, type RequestHandler } from '@sveltejs/kit';
import { archiveDocument } from '$lib/server/documents';
import {
	cookieHeader,
	documentAPIErrorResponse,
	jsonError,
	requireAuthenticatedUser
} from '../../api.server';

const ARCHIVE_ERROR_FALLBACK = 'Failed to archive document';

export const POST: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;
	if (!params.id) return jsonError(400, 'invalid document id');

	try {
		return json(await archiveDocument(fetch, cookieHeader(request), params.id));
	} catch (error) {
		return documentAPIErrorResponse(error, ARCHIVE_ERROR_FALLBACK);
	}
};
