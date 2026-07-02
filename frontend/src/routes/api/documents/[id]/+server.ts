import { json, type RequestHandler } from '@sveltejs/kit';
import { getDocument, updateDocument } from '$lib/server/documents';
import {
	cookieHeader,
	documentAPIErrorResponse,
	jsonError,
	readDocumentUpdateInput,
	requireAuthenticatedUser
} from '../api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load document';
const UPDATE_ERROR_FALLBACK = 'Failed to update document';

export const GET: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;
	if (!params.id) return jsonError(400, 'invalid document id');

	try {
		return json(await getDocument(fetch, cookieHeader(request), params.id));
	} catch (error) {
		return documentAPIErrorResponse(error, LOAD_ERROR_FALLBACK);
	}
};

export const PATCH: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;
	if (!params.id) return jsonError(400, 'invalid document id');

	const input = await readDocumentUpdateInput(request);
	if (input instanceof Response) return input;

	try {
		return json(await updateDocument(fetch, cookieHeader(request), params.id, input));
	} catch (error) {
		return documentAPIErrorResponse(error, UPDATE_ERROR_FALLBACK);
	}
};
