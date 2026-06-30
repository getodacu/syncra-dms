import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = ({ locals }) => ({
	user: locals.user
		? { name: locals.user.name, email: locals.user.email, role: locals.user.role }
		: null,
	session: locals.session ? { expiresAt: locals.session.expiresAt } : null
});
