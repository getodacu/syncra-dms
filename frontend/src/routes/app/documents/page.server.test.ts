import { describe, expect, it } from 'vitest';

import { load } from './+page.server';

describe('documents page server load', () => {
	it('exposes document permission flags from locals', () => {
		const result = load(loadEvent(['document.view', 'document.create'], 'unit-id') as never);

		expect(result).toEqual({
			canViewDocuments: true,
			canCreateDocuments: true,
			canUpdateDocuments: false,
			canDeleteDocuments: false,
			canDownloadDocuments: false,
			selectedOrganizationUnitId: 'unit-id'
		});
	});

	it('returns full document access for system administrators', () => {
		const result = load(loadEvent(['system.admin']) as never);

		expect(result).toEqual({
			canViewDocuments: true,
			canCreateDocuments: true,
			canUpdateDocuments: true,
			canDeleteDocuments: true,
			canDownloadDocuments: true,
			selectedOrganizationUnitId: null
		});
	});

	it('returns no document access without document permissions, regardless of auth role', () => {
		const result = load(loadEvent([], undefined, { role: 'admin' }) as never);

		expect(result).toEqual({
			canViewDocuments: false,
			canCreateDocuments: false,
			canUpdateDocuments: false,
			canDeleteDocuments: false,
			canDownloadDocuments: false,
			selectedOrganizationUnitId: null
		});
	});
});

function loadEvent(
	permissions: string[],
	organizationUnitId?: string,
	user: { role: string } | null = { role: 'user' }
) {
	const url = new URL('http://localhost/app/documents');
	if (organizationUnitId) url.searchParams.set('organizationUnitId', organizationUnitId);
	return {
		locals: { user, permissions },
		url
	};
}
