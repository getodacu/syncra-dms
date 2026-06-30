import { redirect, type Handle, type ResolveOptions } from '@sveltejs/kit';
import { paraglideMiddleware } from '$lib/paraglide/server';
import { getTextDirection } from '$lib/paraglide/runtime';
import type { Locale } from '$lib/paraglide/runtime';
import { clearSessionCookie, getSession, hasSessionCookie, setPreferredLanguageCookie } from '$lib/server/auth';

const guestOnlyRoutes = new Set(['/login', '/signup', '/signup-confirmation', '/recover-password']);

function isProtectedRoute(pathname: string) {
	return pathname === '/app' || pathname.startsWith('/app/');
}

const appHandle: Handle = async ({ event, resolve }) => {
	const cookieHeader = event.request.headers.get('cookie');
	try {
		const auth = await getSession(event.fetch, cookieHeader);
		event.locals.session = auth?.session ?? null;
		event.locals.user = auth?.user ?? null;
		if (auth?.user?.preferredLanguage) {
			setPreferredLanguageCookie(event.cookies, auth.user.preferredLanguage);
		}
		if (!auth && hasSessionCookie(cookieHeader)) {
			clearSessionCookie(event.cookies);
		}
	} catch {
		event.locals.session = null;
		event.locals.user = null;
		if (isProtectedRoute(event.url.pathname)) redirect(303, '/login');
	}

	if (isProtectedRoute(event.url.pathname) && !event.locals.user) {
		redirect(303, '/login');
	}
	if (guestOnlyRoutes.has(event.url.pathname) && event.locals.user) {
		redirect(303, '/app');
	}

	return resolve(event);
};

export const handle: Handle = ({ event, resolve }) =>
	paraglideMiddleware(event.request, ({ request: localizedRequest, locale }) => {
		event.request = localizedRequest;
		return appHandle({
			event,
			resolve: (event, options) =>
				resolve(event, {
					...options,
					transformPageChunk: localizePageAttributes(locale, options?.transformPageChunk)
				})
		});
	});

export function localizePageAttributes(
	locale: Locale,
	transformPageChunk?: ResolveOptions['transformPageChunk']
): NonNullable<ResolveOptions['transformPageChunk']> {
	return ({ html, done }) => {
		const transformed = html.replace('%lang%', locale).replace('%dir%', getTextDirection(locale));
		return transformPageChunk ? transformPageChunk({ html: transformed, done }) : transformed;
	};
}
