import type { PageServerLoad } from './$types';

function hasPermission(permissions: string[], permission: string) {
	return permissions.includes('system.admin') || permissions.includes(permission);
}

export const load: PageServerLoad = ({ locals, url }) => {
	const permissions = locals.permissions ?? [];
	return {
		canViewGroups: hasPermission(permissions, 'group.view'),
		canManageGroups:
			permissions.includes('system.admin') ||
			['group.create', 'group.update', 'group.delete', 'group.manage_users', 'group.assign_roles'].some(
				(permission) => permissions.includes(permission)
			),
		canCreateGroups: hasPermission(permissions, 'group.create'),
		canUpdateGroups: hasPermission(permissions, 'group.update'),
		canDeleteGroups: hasPermission(permissions, 'group.delete'),
		canManageGroupUsers: hasPermission(permissions, 'group.manage_users'),
		canAssignGroupRoles: hasPermission(permissions, 'group.assign_roles'),
		selectedGroupId: url.searchParams.get('selectedGroupId')
	};
};
