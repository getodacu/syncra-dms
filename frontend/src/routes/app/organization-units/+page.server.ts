import type { PageServerLoad } from './$types';

const organizationUnitManagePermissions = [
	'system.admin',
	'organization_unit.manage_hierarchy',
	'organization_unit.create',
	'organization_unit.update',
	'organization_unit.delete'
];

export const load: PageServerLoad = ({ locals, url }) => {
	const permissions = locals.permissions ?? [];
	return {
		canManageOrganizationUnits: organizationUnitManagePermissions.some((permission) =>
			permissions.includes(permission)
		),
		selectedId: url.searchParams.get('selectedId')
	};
};
