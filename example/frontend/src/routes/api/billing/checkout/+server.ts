import { json } from "@sveltejs/kit";

import {
	attachBillingOrderCheckoutSession,
	createCreditOrder,
	isBillingApiError
} from "$lib/server/billing";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import { createStripeCheckoutSession, isStripeCheckoutError } from "$lib/server/stripe";
import type { RequestHandler } from "./$types";

/**
 * Formats a catchable checkout error into a standard JSON response.
 * Maps backend API and Stripe errors to appropriate HTTP status codes.
 * 
 * @param error - The error thrown during checkout processing
 * @returns A JSON Response with error message and status code
 */
function checkoutErrorResponse(error: unknown) {
	if (isBillingApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	if (isStripeCheckoutError(error)) {
		return jsonPublicErrorResponse(502, error.message, "Failed to start checkout");
	}

	throw error;
}

/**
 * Validates the requested credit amount to buy.
 * Credits must be a positive integer and a multiple of 1000.
 * 
 * @param value - The raw credits amount input from client request
 */
function creditQuantity(value: unknown) {
	if (typeof value !== "number" || !Number.isInteger(value) || value <= 0) {
		return { ok: false as const, error: "credits must be a positive integer" };
	}
	if (value % 1000 !== 0) {
		return { ok: false as const, error: "credits must be a multiple of 1000" };
	}
	return { ok: true as const, credits: value };
}

/**
 * POST /api/billing/checkout
 * Initializes a new credit purchase checkout flow via Stripe Checkout.
 * 
 * @returns {Promise<Response>} 200 with checkout URL, or error status code
 */
export const POST: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	let body: unknown;
	try {
		body = await request.json();
	} catch {
		return json({ error: "invalid checkout request" }, { status: 400 });
	}

	const quantity = creditQuantity(
		typeof body === "object" && body !== null && "credits" in body ? body.credits : undefined
	);
	if (!quantity.ok) {
		return json({ error: quantity.error }, { status: 400 });
	}

	try {
		const order = await createCreditOrder(fetch, {
			userId: locals.user.id,
			credits: quantity.credits
		});
		const session = await createStripeCheckoutSession({
			order
		});
		await attachBillingOrderCheckoutSession(fetch, {
			orderId: order.id,
			checkoutSessionId: session.id
		});

		return json({ url: session.url });
	} catch (error) {
		return checkoutErrorResponse(error);
	}
};
