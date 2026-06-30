import { describe, expect, it } from "vitest";

import {
	defaultBillingProfileForm,
	formFromBillingProfile,
	normalizeBillingProfileForm,
	validateBillingProfileForm
} from "./billing-profile-utils";

describe("billing profile form utilities", () => {
	it("defaults an empty profile from the account user", () => {
		expect(defaultBillingProfileForm({ name: "Radu Boncea", email: "radu@example.com" })).toMatchObject({
			entity_type: "individual",
			billing_name: "Radu Boncea",
			billing_email: "radu@example.com",
			country_code: "RO"
		});
	});

	it("normalizes whitespace and uppercases country", () => {
		expect(
			normalizeBillingProfileForm({
				entity_type: "company",
				billing_name: " ICI Bucuresti ",
				billing_email: " billing@example.com ",
				country_code: " ro ",
				address_line1: " Maresal Averescu 8-10 ",
				address_line2: " ",
				city: " Bucuresti ",
				region: " Sector 1 ",
				postal_code: " 011455 ",
				fiscal_code: " RO2785503 ",
				registration_number: " "
			})
		).toEqual({
			entity_type: "company",
			billing_name: "ICI Bucuresti",
			billing_email: "billing@example.com",
			country_code: "RO",
			address_line1: "Maresal Averescu 8-10",
			address_line2: "",
			city: "Bucuresti",
			region: "Sector 1",
			postal_code: "011455",
			fiscal_code: "RO2785503",
			registration_number: ""
		});
	});

	it("requires fiscal code for Romanian companies only", () => {
		const form = defaultBillingProfileForm({ name: "Radu", email: "radu@example.com" });
		form.entity_type = "company";
		form.billing_name = "ICI Bucuresti";
		form.country_code = "RO";
		form.address_line1 = "Maresal Averescu 8-10";
		form.city = "Bucuresti";
		form.postal_code = "011455";

		expect(validateBillingProfileForm(form)).toEqual("Fiscal code is required for Romanian companies.");
		form.country_code = "DE";
		expect(validateBillingProfileForm(form)).toBeNull();
	});

	it("maps saved API profiles back into form state", () => {
		expect(
			formFromBillingProfile({
				id: "profile-1",
				user_id: "user-1",
				entity_type: "company",
				billing_name: "ICI Bucuresti",
				billing_email: "billing@example.com",
				country_code: "RO",
				address_line1: "Maresal Averescu 8-10",
				city: "Bucuresti",
				postal_code: "011455",
				fiscal_code: "RO2785503",
				created_at: "2026-06-05T09:00:00Z",
				updated_at: "2026-06-05T09:00:00Z"
			})
		).toMatchObject({
			entity_type: "company",
			billing_name: "ICI Bucuresti",
			fiscal_code: "RO2785503",
			registration_number: ""
		});
	});

	it("validates required fields and invalid email in order", () => {
		const form = defaultBillingProfileForm({ name: "", email: "not-an-email" });

		expect(validateBillingProfileForm(form)).toBe("Full name is required.");
		form.billing_name = "Radu Boncea";
		expect(validateBillingProfileForm(form)).toBe("Billing email is invalid.");
		form.billing_email = "radu@example.com";
		form.country_code = "";
		expect(validateBillingProfileForm(form)).toBe("Country is required.");
		form.country_code = "RO";
		expect(validateBillingProfileForm(form)).toBe("Address line 1 is required.");
		form.address_line1 = "Maresal Averescu 8-10";
		expect(validateBillingProfileForm(form)).toBe("City is required.");
		form.city = "Bucuresti";
		expect(validateBillingProfileForm(form)).toBe("Postal code is required.");
		form.postal_code = "011455";
		expect(validateBillingProfileForm(form)).toBeNull();
	});

	it("uses company-specific name validation and normalizes unknown entity values", () => {
		const normalized = normalizeBillingProfileForm({
			...defaultBillingProfileForm({ name: "", email: "billing@example.com" }),
			entity_type: "firm"
		});

		expect(normalized.entity_type).toBe("individual");

		const companyForm = defaultBillingProfileForm({ name: "", email: "billing@example.com" });
		companyForm.entity_type = "company";
		expect(validateBillingProfileForm(companyForm)).toBe("Company name is required.");
	});

	it("rejects invalid ISO country codes", () => {
		const form = {
			...defaultBillingProfileForm({ name: "Radu Boncea", email: "radu@example.com" }),
			address_line1: "Maresal Averescu 8-10",
			city: "Bucuresti",
			postal_code: "011455"
		};

		for (const country_code of ["ROU", "R0", "1", "ZZ"]) {
			expect(validateBillingProfileForm({ ...form, country_code })).toBe("Country code is invalid.");
		}
	});

	it("enforces billing profile field length limits", () => {
		const form = {
			...defaultBillingProfileForm({ name: "Radu Boncea", email: "radu@example.com" }),
			address_line1: "Maresal Averescu 8-10",
			city: "Bucuresti",
			postal_code: "011455"
		};

		expect(validateBillingProfileForm({ ...form, billing_name: "a".repeat(256) })).toBe(
			"Full name must be 255 characters or fewer."
		);
		expect(
			validateBillingProfileForm({
				...form,
				billing_email: `${"a".repeat(309)}@example.com`
			})
		).toBe("Billing email must be 320 characters or fewer.");
		expect(validateBillingProfileForm({ ...form, address_line1: "a".repeat(256) })).toBe(
			"Address line 1 must be 255 characters or fewer."
		);
		expect(validateBillingProfileForm({ ...form, address_line2: "a".repeat(256) })).toBe(
			"Address line 2 must be 255 characters or fewer."
		);
		expect(validateBillingProfileForm({ ...form, city: "a".repeat(161) })).toBe(
			"City must be 160 characters or fewer."
		);
		expect(validateBillingProfileForm({ ...form, region: "a".repeat(161) })).toBe(
			"Region/state must be 160 characters or fewer."
		);
		expect(validateBillingProfileForm({ ...form, postal_code: "a".repeat(41) })).toBe(
			"Postal code must be 40 characters or fewer."
		);
		expect(
			validateBillingProfileForm({
				...form,
				entity_type: "company",
				country_code: "DE",
				fiscal_code: "a".repeat(81)
			})
		).toBe("Fiscal code must be 80 characters or fewer.");
		expect(
			validateBillingProfileForm({
				...form,
				entity_type: "company",
				country_code: "DE",
				registration_number: "a".repeat(121)
			})
		).toBe("Registration number must be 120 characters or fewer.");
	});
});
