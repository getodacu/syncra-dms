import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = ({ locals }) => ({
	user: locals.user
		? {
				id: locals.user.id,
				name: locals.user.name,
				email: locals.user.email,
				image: locals.user.image,
				role: locals.user.role
			}
		: null,
	session: locals.session ? { expiresAt: locals.session.expiresAt } : null
});
