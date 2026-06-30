import { beforeEach, describe, expect, it, vi } from "vitest";

import { BillingApiError } from "$lib/server/billing";
import { GET, PUT } from "./+server";
import type { RequestEvent } from "./$types";

const { getBillingProfileMock, upsertBillingProfileMock, BillingApiErrorMock } =
	vi.hoisted(() => {
		class MockBillingApiError extends Error {
			status: number;

			constructor(status: number, message: string) {
				super(message);
				this.name = "BillingApiError";
				this.status = status;
			}
		}

		return {
			getBillingProfileMock: vi.fn(),
			upsertBillingProfileMock: vi.fn(),
			BillingApiErrorMock: MockBillingApiError
		};
	});

vi.mock("$lib/server/billing", () => ({
	getBillingProfile: getBillingProfileMock,
	upsertBillingProfile: upsertBillingProfileMock,
	BillingApiError: BillingApiErrorMock,
	isBillingApiError: (error: unknown) => error instanceof BillingApiErrorMock
}));

function createEvent(body?: unknown, user: unknown = { id: "user-1" }) {
	const init =
		body === undefined
			? undefined
			: {
					method: "PUT",
					headers: { "content-type": "application/json" },
					body: typeof body === "string" ? body : JSON.stringify(body)
				};

	return {
		request: new Request("http://localhost/api/billing/profile", init),
		url: new URL("http://localhost/api/billing/profile"),
		fetch: vi.fn(),
		locals: { user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

function savedProfile() {
	return {
		id: "profile-1",
		user_id: "user-1",
		entity_type: "company",
		billing_name: "Syncra SRL",
		billing_email: "billing@example.com",
		country_code: "RO",
		address_line1: "Main Street 1",
		address_line2: "Floor 2",
		city: "Bucharest",
		region: "B",
		postal_code: "010101",
		fiscal_code: "RO123",
		registration_number: "J40/1/2026",
		created_at: "2026-06-05T00:00:00Z",
		updated_at: "2026-06-05T00:00:00Z"
	};
}

describe("billing profile API endpoint", () => {
	beforeEach(() => {
		getBillingProfileMock.mockReset();
		upsertBillingProfileMock.mockReset();
	});

	it("returns 401 for unauthenticated GET requests", async () => {
		const response = await GET(createEvent(undefined, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(getBillingProfileMock).not.toHaveBeenCalled();
		expect(upsertBillingProfileMock).not.toHaveBeenCalled();
	});

	it("returns 401 for unauthenticated PUT requests", async () => {
		const response = await PUT(createEvent({ billing_name: "Syncra SRL" }, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(getBillingProfileMock).not.toHaveBeenCalled();
		expect(upsertBillingProfileMock).not.toHaveBeenCalled();
	});

	it("loads the authenticated user's billing profile", async () => {
		getBillingProfileMock.mockResolvedValue(null);
		const event = createEvent();

		const response = await GET(event);

		expect(getBillingProfileMock).toHaveBeenCalledWith(event.fetch, "user-1");
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual({ profile: null });
	});

	it("saves the authenticated user's billing profile and ignores browser user ids", async () => {
		const profile = savedProfile();
		upsertBillingProfileMock.mockResolvedValue(profile);
		const event = createEvent({
			user_id: "attacker",
			entity_type: "company",
			billing_name: "Syncra SRL",
			billing_email: "billing@example.com",
			country_code: "RO",
			address_line1: "Main Street 1",
			address_line2: "Floor 2",
			city: "Bucharest",
			region: "B",
			postal_code: "010101",
			fiscal_code: "RO123",
			registration_number: "J40/1/2026"
		});

		const response = await PUT(event);

		expect(upsertBillingProfileMock).toHaveBeenCalledWith(event.fetch, {
			userId: "user-1",
			entityType: "company",
			billingName: "Syncra SRL",
			billingEmail: "billing@example.com",
			countryCode: "RO",
			addressLine1: "Main Street 1",
			addressLine2: "Floor 2",
			city: "Bucharest",
			region: "B",
			postalCode: "010101",
			fiscalCode: "RO123",
			registrationNumber: "J40/1/2026"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(profile);
	});

	it("preserves billing profile validation errors", async () => {
		upsertBillingProfileMock.mockRejectedValue(
			new BillingApiError(400, "billing name is required")
		);

		const response = await PUT(createEvent({ entity_type: "individual", billing_name: "" }));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "billing name is required" });
	});

	it("normalizes billing service errors", async () => {
		upsertBillingProfileMock.mockRejectedValue(
			new BillingApiError(503, "Billing unavailable")
		);

		const response = await PUT(createEvent({ entity_type: "individual" }));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});

	it("returns 400 for invalid billing profile requests", async () => {
		let response = await PUT(createEvent("{"));
		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "invalid billing profile request"
		});

		response = await PUT(createEvent(null));
		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "invalid billing profile request"
		});

		response = await PUT(createEvent([]));
		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "invalid billing profile request"
		});
		expect(upsertBillingProfileMock).not.toHaveBeenCalled();
	});

	it("returns 400 for malformed entity types", async () => {
		upsertBillingProfileMock.mockResolvedValue(savedProfile());

		let response = await PUT(createEvent({ entity_type: "vendor" }));
		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "invalid billing profile request"
		});

		response = await PUT(createEvent({ entity_type: " company " }));
		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "invalid billing profile request"
		});
		expect(upsertBillingProfileMock).not.toHaveBeenCalled();
	});
});
