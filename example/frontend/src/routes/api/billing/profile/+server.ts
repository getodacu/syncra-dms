import { json } from "@sveltejs/kit";

import { getBillingProfile, isBillingApiError, upsertBillingProfile } from "$lib/server/billing";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

function errorResponse(error: unknown) {
	if (isBillingApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	throw error;
}

function stringValue(value: unknown) {
	return typeof value === "string" ? value : "";
}

function optionalStringValue(value: unknown) {
	return typeof value === "string" ? value : undefined;
}

function entityTypeValue(value: unknown) {
	if (value === "individual" || value === "company") return value;
	return undefined;
}

export const GET: RequestHandler = async ({ fetch, locals }) => {
	if (!locals.user) return json({ error: "authentication required" }, { status: 401 });

	try {
		return json({ profile: await getBillingProfile(fetch, locals.user.id) });
	} catch (error) {
		return errorResponse(error);
	}
};

export const PUT: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) return json({ error: "authentication required" }, { status: 401 });

	let body: unknown;
	try {
		body = await request.json();
	} catch {
		return json({ error: "invalid billing profile request" }, { status: 400 });
	}

	if (typeof body !== "object" || body === null || Array.isArray(body)) {
		return json({ error: "invalid billing profile request" }, { status: 400 });
	}
	const value = body as Record<string, unknown>;
	const entityType = entityTypeValue(value.entity_type);
	if (!entityType) {
		return json({ error: "invalid billing profile request" }, { status: 400 });
	}

	try {
		const profile = await upsertBillingProfile(fetch, {
			userId: locals.user.id,
			entityType,
			billingName: stringValue(value.billing_name),
			billingEmail: stringValue(value.billing_email),
			countryCode: stringValue(value.country_code),
			addressLine1: stringValue(value.address_line1),
			addressLine2: optionalStringValue(value.address_line2),
			city: stringValue(value.city),
			region: optionalStringValue(value.region),
			postalCode: stringValue(value.postal_code),
			fiscalCode: optionalStringValue(value.fiscal_code),
			registrationNumber: optionalStringValue(value.registration_number)
		});
		return json(profile);
	} catch (error) {
		return errorResponse(error);
	}
};
