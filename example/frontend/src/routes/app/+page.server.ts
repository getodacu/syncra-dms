import { getDashboardSummary, isDashboardApiError } from "$lib/server/dashboard";
import { publicErrorMessage } from "$lib/server/public-errors";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ fetch, locals }) => {
	if (!locals.user) {
		return {
			initialSummary: null,
			initialSummaryError: "Authentication required"
		};
	}

	try {
		const summary = await getDashboardSummary(fetch, {
			userId: locals.user.id,
			range: "30d"
		});
		return {
			initialSummary: summary,
			initialSummaryError: null
		};
	} catch (error) {
		if (isDashboardApiError(error)) {
			return {
				initialSummary: null,
				initialSummaryError: publicErrorMessage(
					error.status,
					error.message,
					"Failed to load dashboard summary"
				)
			};
		}

		throw error;
	}
};
