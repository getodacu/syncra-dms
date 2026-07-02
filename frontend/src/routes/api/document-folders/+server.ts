import { json, type RequestHandler } from '@sveltejs/kit';
import { createDocumentFolder } from '$lib/server/documents';
import {
	cookieHeader,
	documentAPIErrorResponse,
	readDocumentFolderInput,
	requireAuthenticatedUser
} from '../documents/api.server';

const CREATE_ERROR_FALLBACK = 'Failed to create document folder';

export const POST: RequestHandler = async ({ fetch, locals, request }) => {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;

	const input = await readDocumentFolderInput(request);
	if (input instanceof Response) return input;

	try {
		return json(await createDocumentFolder(fetch, cookieHeader(request), input), { status: 201 });
	} catch (error) {
		return documentAPIErrorResponse(error, CREATE_ERROR_FALLBACK);
	}
};
