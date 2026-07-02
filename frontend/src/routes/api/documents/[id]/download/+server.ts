import { type RequestHandler } from '@sveltejs/kit';
import { downloadDocument } from '$lib/server/documents';
import {
	cookieHeader,
	documentAPIErrorResponse,
	jsonError,
	requireAuthenticatedUser
} from '../../api.server';

const DOWNLOAD_ERROR_FALLBACK = 'Failed to download document';

export const GET: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;
	if (!params.id) return jsonError(400, 'invalid document id');

	try {
		const upstream = await downloadDocument(fetch, cookieHeader(request), params.id);
		return new Response(upstream.body, {
			status: upstream.status,
			headers: {
				'content-type': upstream.headers.get('content-type') ?? 'application/octet-stream',
				'content-disposition': upstream.headers.get('content-disposition') ?? 'attachment'
			}
		});
	} catch (error) {
		return documentAPIErrorResponse(error, DOWNLOAD_ERROR_FALLBACK);
	}
};
