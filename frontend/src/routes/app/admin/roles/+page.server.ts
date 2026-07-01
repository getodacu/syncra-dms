import type { PageServerLoad } from './$types';

function hasPermission(permissions: string[], permission: string) {
	return permissions.includes('system.admin') || permissions.includes(permission);
}

export const load: PageServerLoad = ({ locals, url }) => {
	const permissions = locals.permissions ?? [];
	return {
		canViewRoles: hasPermission(permissions, 'role.view'),
		canManageRoles:
			permissions.includes('system.admin') ||
			['role.create', 'role.update', 'role.delete', 'role.assign_permissions'].some((permission) =>
				permissions.includes(permission)
			),
		canCreateRoles: hasPermission(permissions, 'role.create'),
		canUpdateRoles: hasPermission(permissions, 'role.update'),
		canDeleteRoles: hasPermission(permissions, 'role.delete'),
		canAssignRolePermissions: hasPermission(permissions, 'role.assign_permissions'),
		selectedRoleId: url.searchParams.get('selectedRoleId')
	};
};
