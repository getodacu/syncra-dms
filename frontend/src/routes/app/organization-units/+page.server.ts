import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals, url }) => ({
	canManageOrganizationUnits: locals.user?.role === 'admin',
	selectedId: url.searchParams.get('selectedId')
});
