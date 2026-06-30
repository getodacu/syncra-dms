import { getCreditBalance, isBillingApiError } from "$lib/server/billing";
import { publicErrorMessage } from "$lib/server/public-errors";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ fetch, locals }) => {
	if (!locals.user) {
		return {
			user: null,
			impersonation: null,
			initialCreditBalance: null,
			initialCreditBalanceError: null
		};
	}

	try {
		return {
			user: locals.user,
			impersonation: locals.impersonation,
			initialCreditBalance: await getCreditBalance(fetch, locals.user.id),
			initialCreditBalanceError: null
		};
	} catch (error) {
		if (isBillingApiError(error)) {
			return {
				user: locals.user,
				impersonation: locals.impersonation,
				initialCreditBalance: null,
				initialCreditBalanceError: publicErrorMessage(
					error.status,
					error.message,
					"Failed to load credit balance"
				)
			};
		}

		throw error;
	}
};
