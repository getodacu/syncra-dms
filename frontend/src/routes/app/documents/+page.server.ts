import type { PageServerLoad } from './$types';

const documentPermissions = {
	view: ['system.admin', 'document.view'],
	create: ['system.admin', 'document.create'],
	update: ['system.admin', 'document.update'],
	delete: ['system.admin', 'document.delete'],
	download: ['system.admin', 'document.download']
} as const;

function hasAny(permissions: string[], required: readonly string[]) {
	return required.some((permission) => permissions.includes(permission));
}

export const load: PageServerLoad = ({ locals, url }) => {
	const permissions = locals.permissions ?? [];
	return {
		canViewDocuments: hasAny(permissions, documentPermissions.view),
		canCreateDocuments: hasAny(permissions, documentPermissions.create),
		canUpdateDocuments: hasAny(permissions, documentPermissions.update),
		canDeleteDocuments: hasAny(permissions, documentPermissions.delete),
		canDownloadDocuments: hasAny(permissions, documentPermissions.download),
		selectedOrganizationUnitId: url.searchParams.get('organizationUnitId')
	};
};
