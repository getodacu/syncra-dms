import { fail } from '@sveltejs/kit';
import {
	archiveOrganizationUnit,
	createOrganizationUnit,
	getOrganizationUnitTree,
	isOrganizationUnitApiError,
	moveOrganizationUnit,
	updateOrganizationUnit
} from '$lib/server/organization-units';
import { publicErrorMessage, publicErrorStatus } from '$lib/server/public-errors';
import type { OrganizationUnitInput } from '$lib/server/organization-units';
import type { Actions, PageServerLoad } from './$types';

const LOAD_ERROR_FALLBACK = 'Failed to load organization units';
const CREATE_ERROR_FALLBACK = 'Failed to create organization unit';
const UPDATE_ERROR_FALLBACK = 'Failed to update organization unit';
const MOVE_ERROR_FALLBACK = 'Failed to move organization unit';
const ARCHIVE_ERROR_FALLBACK = 'Failed to archive organization unit';

export const load: PageServerLoad = async ({ fetch, locals, request, url }) => {
	const canManageOrganizationUnits = locals.user?.role === 'admin';
	const cookieHeader = request.headers.get('cookie');
	const selectedId = url.searchParams.get('selectedId');
	try {
		const tree = await getOrganizationUnitTree(fetch, cookieHeader);
		return {
			units: tree.units,
			loadError: null,
			canManageOrganizationUnits,
			selectedId
		};
	} catch (error) {
		if (isOrganizationUnitApiError(error)) {
			return {
				units: [],
				loadError: publicErrorMessage(error.status, error.message, LOAD_ERROR_FALLBACK),
				canManageOrganizationUnits,
				selectedId
			};
		}
		throw error;
	}
};

export const actions = {
	create: async ({ fetch, request }) => {
		const data = await request.formData();
		const selectedId = textValue(data, 'selectedId') || null;
		const input = organizationUnitInput(data);
		try {
			const unit = await createOrganizationUnit(fetch, cookieHeader(request), input);
			return { success: true, selectedId: unit.id };
		} catch (error) {
			return failKnownOrganizationUnitError(error, CREATE_ERROR_FALLBACK, 'create', selectedId, input);
		}
	},
	update: async ({ fetch, request }) => {
		const data = await request.formData();
		const id = textValue(data, 'id');
		const input = organizationUnitInput(data);
		try {
			await updateOrganizationUnit(fetch, cookieHeader(request), id, input);
			return { success: true, selectedId: id };
		} catch (error) {
			return failKnownOrganizationUnitError(error, UPDATE_ERROR_FALLBACK, 'update', id || null, {
				id,
				...input
			});
		}
	},
	move: async ({ fetch, request }) => {
		const data = await request.formData();
		const id = textValue(data, 'id');
		const parentId = parentIdValue(data);
		try {
			await moveOrganizationUnit(fetch, cookieHeader(request), id, parentId);
			return { success: true, selectedId: id };
		} catch (error) {
			return failKnownOrganizationUnitError(error, MOVE_ERROR_FALLBACK, 'move', id || null, {
				id,
				parentId
			});
		}
	},
	archive: async ({ fetch, request }) => {
		const data = await request.formData();
		const id = textValue(data, 'id');
		try {
			await archiveOrganizationUnit(fetch, cookieHeader(request), id);
			return { success: true, selectedId: null };
		} catch (error) {
			return failKnownOrganizationUnitError(error, ARCHIVE_ERROR_FALLBACK, 'archive', id || null, {
				id
			});
		}
	}
} satisfies Actions;

function cookieHeader(request: Request) {
	return request.headers.get('cookie');
}

function organizationUnitInput(data: FormData): OrganizationUnitInput {
	return {
		parentId: parentIdValue(data),
		name: textValue(data, 'name'),
		code: textValue(data, 'code'),
		description: textValue(data, 'description')
	};
}

function parentIdValue(data: FormData) {
	const parentId = textValue(data, 'parentId').trim();
	return parentId ? parentId : null;
}

function textValue(data: FormData, key: string) {
	const value = data.get(key);
	return typeof value === 'string' ? value : '';
}

function failKnownOrganizationUnitError(
	error: unknown,
	fallback: string,
	action: string,
	selectedId: string | null,
	values: Record<string, string | null>
) {
	if (isOrganizationUnitApiError(error)) {
		return fail(publicErrorStatus(error.status), {
			error: publicErrorMessage(error.status, error.message, fallback),
			action,
			selectedId,
			values
		});
	}
	throw error;
}
