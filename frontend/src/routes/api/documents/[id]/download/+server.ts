import { type RequestHandler } from '@sveltejs/kit';
import { downloadDocument } from '$lib/server/documents';
import {
	cookieHeader,
	documentAPIErrorResponse,
	jsonError,
	requireAuthenticatedUser
} from '../../api.server';

const DOWNLOAD_ERROR_FALLBACK = 'Failed to download document';
const UNSAFE_UPSTREAM_HEADERS = [
	'set-cookie',
	'connection',
	'keep-alive',
	'proxy-authenticate',
	'proxy-authorization',
	'te',
	'trailer',
	'transfer-encoding',
	'upgrade'
];

export const GET: RequestHandler = async ({ fetch, locals, params, request }) => {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;
	if (!params.id) return jsonError(400, 'invalid document id');

	try {
		const upstream = await downloadDocument(fetch, cookieHeader(request), params.id);
		const headers = new Headers(upstream.headers);
		for (const header of UNSAFE_UPSTREAM_HEADERS) {
			headers.delete(header);
		}
		if (!headers.has('content-type')) headers.set('content-type', 'application/octet-stream');
		if (!headers.has('content-disposition')) headers.set('content-disposition', 'attachment');

		return new Response(upstream.body, {
			status: upstream.status,
			headers
		});
	} catch (error) {
		return documentAPIErrorResponse(error, DOWNLOAD_ERROR_FALLBACK);
	}
};
