import { describe, expect, it } from "vitest";

import { CREDIT_PRICING_TIERS, CREDIT_RULES, checkoutHref } from "./pricing-data";

describe("Pricing Data", () => {
	it("contains exactly the four credit purchase tiers", () => {
		expect(CREDIT_PRICING_TIERS).toHaveLength(4);
		expect(CREDIT_PRICING_TIERS.map((tier) => tier.id)).toEqual([
			"tier_1",
			"tier_2",
			"tier_3",
			"tier_4"
		]);
	});

	it("uses the credit-only price list with no monthly plan metadata", () => {
		expect(CREDIT_PRICING_TIERS[0]).toMatchObject({
			id: "tier_1",
			creditRange: "1,000-4,000 credits",
			unitPrice: "€10.00 / 1000 credits"
		});
		expect(CREDIT_PRICING_TIERS[1]).toMatchObject({
			id: "tier_2",
			creditRange: "5,000-9,000 credits",
			unitPrice: "€9.50 / 1000 credits"
		});
		expect(CREDIT_PRICING_TIERS[2]).toMatchObject({
			id: "tier_3",
			creditRange: "10,000-19,000 credits",
			unitPrice: "€9.00 / 1000 credits"
		});
		expect(CREDIT_PRICING_TIERS[3]).toMatchObject({
			id: "tier_4",
			creditRange: "20,000+ credits",
			unitPrice: "€8.50 / 1000 credits"
		});
	});

	it("routes users to app billing when authenticated and signup otherwise", () => {
		expect(checkoutHref(true)).toBe("/app/billing");
		expect(checkoutHref(false)).toBe("/signup");
	});

	it("defines the credit ledger rules without subscription language", () => {
		expect(CREDIT_RULES.creditConversion).toContain("1 credit = 1 page");
		expect(CREDIT_RULES.signupBonus).toContain("500");
		expect(CREDIT_RULES.purchasedCredits).toContain("never expire");
		expect(CREDIT_RULES.purchaseBlocks).toContain("1000-credit blocks");
		expect(CREDIT_RULES.noSubscriptions).toContain("No monthly subscriptions");
		expect(CREDIT_RULES.successfulProcessing).toContain("successful job processing");
	});
});
