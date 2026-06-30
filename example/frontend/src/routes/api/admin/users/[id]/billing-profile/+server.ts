import { json } from "@sveltejs/kit";

import { upsertAdminUserBillingProfile, type AdminBillingProfileInput } from "$lib/server/admin";
import type { RequestHandler } from "./$types";
import {
	adminApiErrorResponse,
	adminAuthError,
	readJsonObject,
	rejectUnknownKeys
} from "../../../admin-route-utils";

const BILLING_PROFILE_KEYS = new Set([
	"entity_type",
	"billing_name",
	"billing_email",
	"country_code",
	"address_line1",
	"address_line2",
	"city",
	"region",
	"postal_code",
	"fiscal_code",
	"registration_number"
]);

function stringValue(value: unknown) {
	return typeof value === "string" ? value : "";
}

function optionalStringValue(value: unknown) {
	return typeof value === "string" ? value : undefined;
}

function entityTypeValue(value: unknown): AdminBillingProfileInput["entity_type"] | undefined {
	if (value === "individual" || value === "company") return value;
	return undefined;
}

export const PUT: RequestHandler = async ({ request, params, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	const parsed = await readJsonObject(request, "invalid billing profile request");
	if (parsed.error) return parsed.error;
	const unknown = rejectUnknownKeys(parsed.value, BILLING_PROFILE_KEYS, "invalid billing profile request");
	if (unknown) return unknown;

	const entityType = entityTypeValue(parsed.value.entity_type);
	if (!entityType) {
		return json({ error: "invalid billing profile request" }, { status: 400 });
	}

	try {
		const result = await upsertAdminUserBillingProfile(fetch, request.headers.get("cookie"), params.id, {
			entity_type: entityType,
			billing_name: stringValue(parsed.value.billing_name),
			billing_email: stringValue(parsed.value.billing_email),
			country_code: stringValue(parsed.value.country_code),
			address_line1: stringValue(parsed.value.address_line1),
			address_line2: optionalStringValue(parsed.value.address_line2),
			city: stringValue(parsed.value.city),
			region: optionalStringValue(parsed.value.region),
			postal_code: stringValue(parsed.value.postal_code),
			fiscal_code: optionalStringValue(parsed.value.fiscal_code),
			registration_number: optionalStringValue(parsed.value.registration_number)
		});
		return json(result);
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};
