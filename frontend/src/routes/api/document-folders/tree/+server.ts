import { json, type RequestHandler } from '@sveltejs/kit';
import { getDocumentFolderTree } from '$lib/server/documents';
import {
	cookieHeader,
	documentAPIErrorResponse,
	jsonError,
	requireAuthenticatedUser
} from '../../documents/api.server';

const LOAD_ERROR_FALLBACK = 'Failed to load document folders';

export const GET: RequestHandler = async ({ fetch, locals, request, url }) => {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;

	const organizationUnitId = url.searchParams.get('organizationUnitId')?.trim() ?? '';
	if (!organizationUnitId) return jsonError(400, 'organizationUnitId is required');

	try {
		return json(await getDocumentFolderTree(fetch, cookieHeader(request), organizationUnitId));
	} catch (error) {
		return documentAPIErrorResponse(error, LOAD_ERROR_FALLBACK);
	}
};
