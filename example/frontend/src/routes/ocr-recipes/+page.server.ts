import { isSchemaApiError, listJsonRecipes } from "$lib/server/schemas";
import { publicErrorMessage } from "$lib/server/public-errors";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ fetch, locals }) => {
	try {
		const result = await listJsonRecipes(fetch, { size: 100, sort: "desc" });
		return {
			isLoggedIn: Boolean(locals.user),
			userId: locals.user?.id ?? null,
			recipes: result.recipes,
			nextCursor: result.next_cursor,
			loadError: null
		};
	} catch (error) {
		if (isSchemaApiError(error)) {
			return {
				isLoggedIn: Boolean(locals.user),
				userId: locals.user?.id ?? null,
				recipes: [],
				nextCursor: null,
				loadError: publicErrorMessage(error.status, error.message, "Unable to load OCR recipes.")
			};
		}
		throw error;
	}
};
