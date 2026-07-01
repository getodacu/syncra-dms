import type { PageServerLoad } from './$types';

const managePermissions = [
	'user.create',
	'user.update',
	'user.delete',
	'user.activate',
	'user.suspend',
	'user.assign_role',
	'user.assign_group',
	'user.assign_unit'
];

function hasPermission(permissions: string[], permission: string) {
	return permissions.includes('system.admin') || permissions.includes(permission);
}

export const load: PageServerLoad = ({ locals, url }) => {
	const permissions = locals.permissions ?? [];
	return {
		canViewUsers: hasPermission(permissions, 'user.view'),
		canManageUsers:
			permissions.includes('system.admin') ||
			managePermissions.some((permission) => permissions.includes(permission)),
		canCreateUsers: hasPermission(permissions, 'user.create'),
		canUpdateUsers: hasPermission(permissions, 'user.update'),
		canDeleteUsers: hasPermission(permissions, 'user.delete'),
		canActivateUsers: hasPermission(permissions, 'user.activate'),
		canSuspendUsers: hasPermission(permissions, 'user.suspend'),
		canAssignUserRoles: hasPermission(permissions, 'user.assign_role'),
		canAssignUserGroups: hasPermission(permissions, 'user.assign_group'),
		canAssignUserUnits: hasPermission(permissions, 'user.assign_unit'),
		selectedUserId: url.searchParams.get('selectedUserId')
	};
};
