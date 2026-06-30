import { redirect, type RequestHandler } from '@sveltejs/kit';
import { clearSessionCookie, isAuthApiError, signOut } from '$lib/server/auth';

export const POST: RequestHandler = async ({ cookies, fetch, locals, request }) => {
	try {
		await signOut(fetch, request.headers.get('cookie'));
	} catch (error) {
		if (!isAuthApiError(error)) throw error;
	} finally {
		clearSessionCookie(cookies);
		locals.session = null;
		locals.user = null;
	}
	redirect(303, '/login');
};
