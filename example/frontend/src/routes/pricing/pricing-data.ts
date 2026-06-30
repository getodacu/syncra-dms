import {
	CREDIT_PURCHASE_TIERS,
	formatCents,
	type CreditPurchaseTier
} from "$lib/billing/pricing";

export type CreditPricingTier = {
	id: CreditPurchaseTier["id"];
	name: string;
	creditRange: string;
	unitPrice: string;
	unitAmountCents: number;
	sampleCredits: number;
	sampleTotal: string;
};

function creditRange(tier: CreditPurchaseTier) {
	if (tier.maxCredits) {
		return `${tier.minCredits.toLocaleString()}-${tier.maxCredits.toLocaleString()} credits`;
	}
	return `${tier.minCredits.toLocaleString()}+ credits`;
}

export const CREDIT_PRICING_TIERS: CreditPricingTier[] = CREDIT_PURCHASE_TIERS.map((tier) => ({
	id: tier.id,
	name: tier.label,
	creditRange: creditRange(tier),
	unitPrice: `${formatCents(tier.unitAmountCents)} / 1000 credits`,
	unitAmountCents: tier.unitAmountCents,
	sampleCredits: tier.minCredits,
	sampleTotal: formatCents((tier.minCredits / 1000) * tier.unitAmountCents)
}));

export const CREDIT_RULES = {
	creditConversion: "1 credit = 1 page",
	signupBonus: "New users receive 500 one-time signup credits",
	purchasedCredits: "Purchased credits never expire",
	purchaseBlocks: "Credit purchases must be in 1000-credit blocks",
	noSubscriptions: "No monthly subscriptions or recurring plan allowances",
	successfulProcessing: "Credits are debited after successful job processing"
};

export function checkoutHref(isLoggedIn: boolean) {
	return isLoggedIn ? "/app/billing" : "/signup";
}
