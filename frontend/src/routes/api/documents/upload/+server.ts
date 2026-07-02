import { json, type RequestHandler } from '@sveltejs/kit';
import { uploadDocument } from '$lib/server/documents';
import {
	cookieHeader,
	documentAPIErrorResponse,
	requireAuthenticatedUser
} from '../api.server';

const UPLOAD_ERROR_FALLBACK = 'Failed to upload document';

export const POST: RequestHandler = async ({ fetch, locals, request }) => {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;

	try {
		const formData = await request.formData();
		return json(await uploadDocument(fetch, cookieHeader(request), formData), { status: 201 });
	} catch (error) {
		return documentAPIErrorResponse(error, UPLOAD_ERROR_FALLBACK);
	}
};
