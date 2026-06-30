import { env } from "$env/dynamic/private";
import Stripe from "stripe";

import { CREDIT_BLOCK_SIZE, type CreditPurchaseTierID } from "$lib/billing/pricing";
import type { BillingOrderResponse } from "./billing";

let stripeClient: Stripe | null = null;

const STRIPE_PRICE_ENV_BY_TIER = {
	tier_1: "STRIPE_PRICE_ID_TIER_1",
	tier_2: "STRIPE_PRICE_ID_TIER_2",
	tier_3: "STRIPE_PRICE_ID_TIER_3",
	tier_4: "STRIPE_PRICE_ID_TIER_4"
} satisfies Record<CreditPurchaseTierID, string>;

export class StripeCheckoutError extends Error {
	constructor(message: string) {
		super(message);
		this.name = "StripeCheckoutError";
	}
}

function privateEnv(key: string) {
	return env[key] || nodeEnv()[key];
}

function nodeEnv() {
	if (typeof process !== "undefined") return process.env;
	return (
		globalThis as typeof globalThis & {
			process?: { env?: Record<string, string | undefined> };
		}
	).process?.env ?? {};
}

function stripeSecretKey() {
	const key = privateEnv("STRIPE_SECRET_KEY")?.trim();
	if (!key) throw new StripeCheckoutError("Billing checkout is not configured");
	return key;
}

function isCreditPurchaseTierID(value: string): value is CreditPurchaseTierID {
	return value in STRIPE_PRICE_ENV_BY_TIER;
}

function stripePriceIdForTier(tier: string) {
	if (!isCreditPurchaseTierID(tier)) {
		throw new StripeCheckoutError("Billing checkout price tier is not configured");
	}

	const priceId = privateEnv(STRIPE_PRICE_ENV_BY_TIER[tier])?.trim();
	if (!priceId) throw new StripeCheckoutError(`Stripe price ID is not configured for ${tier}`);
	return priceId;
}

function checkoutQuantity(credits: number) {
	if (
		!Number.isInteger(credits) ||
		credits < CREDIT_BLOCK_SIZE ||
		credits % CREDIT_BLOCK_SIZE !== 0
	) {
		throw new StripeCheckoutError("Billing checkout credit quantity is invalid");
	}
	return credits;
}

function appOrigin() {
	const origin = privateEnv("SYNCRA_APP_ORIGIN")?.trim();
	if (!origin) throw new StripeCheckoutError("Billing checkout app origin is not configured");

	let url: URL;
	try {
		url = new URL(origin);
	} catch {
		throw new StripeCheckoutError("Billing checkout app origin must be an http(s) URL");
	}

	if (url.protocol !== "http:" && url.protocol !== "https:") {
		throw new StripeCheckoutError("Billing checkout app origin must be an http(s) URL");
	}

	return `${url.protocol}//${url.host}`;
}

export function stripeWebhookSecret() {
	const secret = privateEnv("STRIPE_WEBHOOK_SECRET")?.trim();
	if (!secret) throw new StripeCheckoutError("Stripe webhook is not configured");
	return secret;
}

export function stripe() {
	if (!stripeClient) {
		stripeClient = new Stripe(stripeSecretKey());
	}
	return stripeClient;
}

export type CreateStripeCheckoutSessionInput = {
	order: Pick<BillingOrderResponse, "id" | "user_id" | "pricing_tier" | "credits">;
};

export async function createStripeCheckoutSession(input: CreateStripeCheckoutSessionInput) {
	const metadata = {
		billing_order_id: input.order.id,
		user_id: input.order.user_id,
		credits: String(input.order.credits)
	};
	const origin = appOrigin();
	const price = stripePriceIdForTier(input.order.pricing_tier);
	const quantity = checkoutQuantity(input.order.credits);
	const session = await stripe().checkout.sessions.create({
		mode: "payment",
		client_reference_id: input.order.id,
		success_url: `${origin}/app/billing?checkout=success&session_id={CHECKOUT_SESSION_ID}`,
		cancel_url: `${origin}/app/billing?checkout=canceled`,
		line_items: [
			{
				price,
				quantity
			}
		],
		metadata,
		payment_intent_data: { metadata }
	});

	if (!session.url) {
		throw new StripeCheckoutError("Stripe checkout session did not include a URL");
	}

	return { id: session.id, url: session.url };
}

export function isStripeCheckoutError(error: unknown): error is StripeCheckoutError {
	return error instanceof StripeCheckoutError;
}
