import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const billingProfileDialogSource = () =>
	readFileSync(new URL("./billing-profile-dialog.svelte", import.meta.url), "utf8");
const normalizeSource = (value: string) => value.replace(/\s+/g, " ").trim();

describe("billing profile dialog source", () => {
	it("contains the required billing profile integration points", () => {
		const source = normalizeSource(billingProfileDialogSource());

		expect(source).toContain('fetch("/api/billing/profile"');
		expect(source).toContain("m.billing_profile_title()");
		expect(source).toContain("lastSavedForm");
		expect(source).toContain("storeSavedForm");
		expect(source).toContain("setDialogOpen");
		expect(source).toContain("if (!nextOpen && saving) return;");
		expect(source).toContain("bind:open={() => open, setDialogOpen}");
		expect(source).toContain("normalizedSavePayload");
		expect(source).toContain("clearCompanyFields");
		expect(source).toContain("requestGeneration === loadGeneration");
		expect(source).toContain("maxlength={255}");
		expect(source).toContain('import { CountrySelect } from "$lib/components/ui/country-select/index.js";');
		expect(source).toContain('id="billing-country-code" bind:value={form.country_code}');
		expect(source).toContain("m.billing_profile_country()");
		expect(source).not.toContain("Country code");
		expect(source).not.toContain("billing-country-code-help");
		expect(source).not.toContain("updateCountryCode");
		expect(source).toContain("m.billing_profile_fiscal_code()");
		expect(source).toContain("m.billing_profile_registration_number()");
		expect(source).toContain("toast.success(m.billing_profile_saved())");
	});

	it("uses Paraglide messages for billing profile labels and feedback", () => {
		const source = billingProfileDialogSource();

		expect(source).toContain('import { m } from "$lib/paraglide/messages.js";');
		for (const messageCall of [
			"m.billing_profile_load_error()",
			"m.billing_profile_save_error()",
			"m.billing_profile_company_name()",
			"m.billing_profile_full_name()",
			"m.billing_profile_description()",
			"m.billing_profile_billing_entity()",
			"m.billing_profile_general_details()",
			"m.billing_profile_billing_address()",
			"m.billing_profile_company_details()",
			"m.billing_profile_save_button()"
		]) {
			expect(source).toContain(messageCall);
		}
	});
});
