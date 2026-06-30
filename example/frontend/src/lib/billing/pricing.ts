export type CreditPurchaseTierID = "tier_1" | "tier_2" | "tier_3" | "tier_4";

export type CreditPurchaseTier = {
	id: CreditPurchaseTierID;
	label: string;
	minCredits: number;
	maxCredits?: number;
	unitAmountCents: number;
};

export type CreditPurchaseQuote = {
	credits: number;
	blocks: number;
	tier: CreditPurchaseTierID;
	tierLabel: string;
	unitAmountCents: number;
	amountCents: number;
	currency: "EUR";
};

export const CREDIT_BLOCK_SIZE = 1000;

export const CREDIT_PURCHASE_TIERS: CreditPurchaseTier[] = [
	{
		id: "tier_1",
		label: "Tier 1",
		minCredits: 1000,
		maxCredits: 4000,
		unitAmountCents: 1000
	},
	{
		id: "tier_2",
		label: "Tier 2",
		minCredits: 5000,
		maxCredits: 9000,
		unitAmountCents: 950
	},
	{
		id: "tier_3",
		label: "Tier 3",
		minCredits: 10000,
		maxCredits: 19000,
		unitAmountCents: 900
	},
	{
		id: "tier_4",
		label: "Tier 4",
		minCredits: 20000,
		unitAmountCents: 850
	}
];

function tierForCredits(credits: number) {
	if (credits >= 20000) return CREDIT_PURCHASE_TIERS[3];
	if (credits >= 10000) return CREDIT_PURCHASE_TIERS[2];
	if (credits >= 5000) return CREDIT_PURCHASE_TIERS[1];
	return CREDIT_PURCHASE_TIERS[0];
}

export function validCreditPurchaseQuantity(credits: number) {
	return (
		Number.isInteger(credits) &&
		credits >= CREDIT_BLOCK_SIZE &&
		credits % CREDIT_BLOCK_SIZE === 0
	);
}

export function quoteCreditPurchase(credits: number): CreditPurchaseQuote {
	if (!validCreditPurchaseQuantity(credits)) {
		throw new Error(`Credit purchase must be in ${CREDIT_BLOCK_SIZE}-credit blocks`);
	}

	const tier = tierForCredits(credits);
	const blocks = credits / CREDIT_BLOCK_SIZE;

	return {
		credits,
		blocks,
		tier: tier.id,
		tierLabel: tier.label,
		unitAmountCents: tier.unitAmountCents,
		amountCents: blocks * tier.unitAmountCents,
		currency: "EUR"
	};
}

export function formatCents(amountCents: number, currency = "EUR") {
	return new Intl.NumberFormat("en-US", {
		style: "currency",
		currency,
		minimumFractionDigits: 2
	}).format(amountCents / 100);
}
