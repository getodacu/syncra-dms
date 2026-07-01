import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals }) => {
	const permissions = locals.permissions ?? [];
	return {
		canViewPermissions: permissions.includes('system.admin') || permissions.includes('role.view'),
		canManagePermissions: false
	};
};
