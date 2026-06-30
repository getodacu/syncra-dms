import { describe, expect, it } from "vitest";

import {
	CREDIT_PURCHASE_TIERS,
	formatCents,
	quoteCreditPurchase,
	validCreditPurchaseQuantity
} from "./pricing";

describe("credit purchase pricing", () => {
	it("quotes credit purchases using the configured tier boundaries", () => {
		expect(quoteCreditPurchase(1000)).toEqual({
			credits: 1000,
			blocks: 1,
			tier: "tier_1",
			tierLabel: "Tier 1",
			unitAmountCents: 1000,
			amountCents: 1000,
			currency: "EUR"
		});
		expect(quoteCreditPurchase(5000)).toMatchObject({
			tier: "tier_2",
			unitAmountCents: 950,
			amountCents: 4750
		});
		expect(quoteCreditPurchase(10000)).toMatchObject({
			tier: "tier_3",
			unitAmountCents: 900,
			amountCents: 9000
		});
		expect(quoteCreditPurchase(20000)).toMatchObject({
			tier: "tier_4",
			unitAmountCents: 850,
			amountCents: 17000
		});
		expect(quoteCreditPurchase(50000)).toMatchObject({
			tier: "tier_4",
			unitAmountCents: 850,
			amountCents: 42500
		});
	});

	it("rejects non-block credit quantities", () => {
		for (const credits of [0, -1000, 1, 999, 1001, 1500]) {
			expect(() => quoteCreditPurchase(credits)).toThrow("1000-credit blocks");
			expect(validCreditPurchaseQuantity(credits)).toBe(false);
		}
	});

	it("formats cents as EUR", () => {
		expect(formatCents(4750)).toBe("€47.50");
		expect(formatCents(1000)).toBe("€10.00");
	});

	it("exposes all four tiers", () => {
		expect(CREDIT_PURCHASE_TIERS).toHaveLength(4);
		expect(CREDIT_PURCHASE_TIERS.map((tier) => tier.id)).toEqual([
			"tier_1",
			"tier_2",
			"tier_3",
			"tier_4"
		]);
	});
});
