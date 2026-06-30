import { beforeEach, describe, expect, it, vi } from "vitest";

const { checkoutSessionCreateMock, privateEnv } = vi.hoisted(() => ({
	checkoutSessionCreateMock: vi.fn(),
	privateEnv: {} as Record<string, string | undefined>
}));

vi.mock("$env/dynamic/private", () => ({ env: privateEnv }));

vi.mock("stripe", () => ({
	default: class MockStripe {
		checkout = {
			sessions: {
				create: checkoutSessionCreateMock
			}
		};

		constructor(_secretKey: string) {}
	}
}));

const stripePriceIds = {
	tier_1: "price_1Tf3VtFamDrIBYmv7bV08cYy",
	tier_2: "price_1Tf3VtFamDrIBYmvQ1GCa7O1",
	tier_3: "price_1Tf3fIFamDrIBYmvZcM0CXie",
	tier_4: "price_1Tf3feFamDrIBYmvLjpEvUqO"
};

type TestBillingOrder = {
	id: string;
	user_id: string;
	pricing_tier: string;
	credits: number;
	amount_cents: number;
	currency: string;
};

function configureCheckoutEnv() {
	privateEnv.STRIPE_SECRET_KEY = "sk_test_secret";
	privateEnv.SYNCRA_APP_ORIGIN = "https://app.syncra.test/some/path?ignored=true";
	privateEnv.STRIPE_PRICE_ID_TIER_1 = stripePriceIds.tier_1;
	privateEnv.STRIPE_PRICE_ID_TIER_2 = stripePriceIds.tier_2;
	privateEnv.STRIPE_PRICE_ID_TIER_3 = stripePriceIds.tier_3;
	privateEnv.STRIPE_PRICE_ID_TIER_4 = stripePriceIds.tier_4;
}

function orderResponse(overrides: Partial<TestBillingOrder> = {}): TestBillingOrder {
	return {
		id: "order-1",
		user_id: "user-1",
		pricing_tier: "tier_2",
		credits: 5000,
		amount_cents: 4750,
		currency: "EUR",
		...overrides
	};
}

describe("Stripe server helper", () => {
	beforeEach(() => {
		vi.resetModules();
		checkoutSessionCreateMock.mockReset();
		for (const key of Object.keys(privateEnv)) delete privateEnv[key];
		delete process.env.STRIPE_SECRET_KEY;
		delete process.env.SYNCRA_APP_ORIGIN;
		delete process.env.STRIPE_PRICE_ID_TIER_1;
		delete process.env.STRIPE_PRICE_ID_TIER_2;
		delete process.env.STRIPE_PRICE_ID_TIER_3;
		delete process.env.STRIPE_PRICE_ID_TIER_4;
		process.env.NODE_ENV = "test";
	});

	it("creates checkout redirect URLs from the configured app origin", async () => {
		configureCheckoutEnv();
		checkoutSessionCreateMock.mockResolvedValue({
			id: "cs_test_123",
			url: "https://checkout.stripe.test/session"
		});
		const { createStripeCheckoutSession } = await import("./stripe");

		await expect(createStripeCheckoutSession({ order: orderResponse() })).resolves.toEqual({
			id: "cs_test_123",
			url: "https://checkout.stripe.test/session"
		});

		expect(checkoutSessionCreateMock).toHaveBeenCalledWith(
			expect.objectContaining({
				success_url:
					"https://app.syncra.test/app/billing?checkout=success&session_id={CHECKOUT_SESSION_ID}",
				cancel_url: "https://app.syncra.test/app/billing?checkout=canceled"
			})
		);
	});

	it.each([
		{
			credits: 1000,
			pricing_tier: "tier_1",
			price: stripePriceIds.tier_1,
			amount_cents: 1000
		},
		{
			credits: 5000,
			pricing_tier: "tier_2",
			price: stripePriceIds.tier_2,
			amount_cents: 4750
		},
		{
			credits: 10000,
			pricing_tier: "tier_3",
			price: stripePriceIds.tier_3,
			amount_cents: 9000
		},
		{
			credits: 20000,
			pricing_tier: "tier_4",
			price: stripePriceIds.tier_4,
			amount_cents: 17000
		}
	])(
		"uses the configured $pricing_tier Stripe Price ID for $credits credits",
		async ({ credits, pricing_tier, price, amount_cents }) => {
			configureCheckoutEnv();
			checkoutSessionCreateMock.mockResolvedValue({
				id: "cs_test_123",
				url: "https://checkout.stripe.test/session"
			});
			const { createStripeCheckoutSession } = await import("./stripe");

			await createStripeCheckoutSession({
				order: orderResponse({
					credits,
					pricing_tier,
					amount_cents
				})
			});

			const checkoutInput = checkoutSessionCreateMock.mock.calls[0]?.[0];
			expect(checkoutInput).toMatchObject({
				line_items: [{ price, quantity: credits }]
			});
			expect(checkoutInput.line_items[0]).not.toHaveProperty("price_data");
		}
	);

	it("rejects checkout creation when the tier Stripe Price ID is missing", async () => {
		configureCheckoutEnv();
		privateEnv.STRIPE_PRICE_ID_TIER_2 = " ";
		const { createStripeCheckoutSession } = await import("./stripe");

		await expect(createStripeCheckoutSession({ order: orderResponse() })).rejects.toMatchObject({
			message: "Stripe price ID is not configured for tier_2"
		});
		expect(checkoutSessionCreateMock).not.toHaveBeenCalled();
	});

	it("rejects checkout creation when the canonical app origin is missing", async () => {
		privateEnv.STRIPE_SECRET_KEY = "sk_test_secret";
		const { createStripeCheckoutSession } = await import("./stripe");

		await expect(createStripeCheckoutSession({ order: orderResponse() })).rejects.toMatchObject({
			message: "Billing checkout app origin is not configured"
		});
		expect(checkoutSessionCreateMock).not.toHaveBeenCalled();
	});

	it("rejects checkout creation when the canonical app origin is invalid", async () => {
		privateEnv.STRIPE_SECRET_KEY = "sk_test_secret";
		privateEnv.SYNCRA_APP_ORIGIN = "ftp://app.syncra.test";
		const { createStripeCheckoutSession } = await import("./stripe");

		await expect(createStripeCheckoutSession({ order: orderResponse() })).rejects.toMatchObject({
			message: "Billing checkout app origin must be an http(s) URL"
		});
		expect(checkoutSessionCreateMock).not.toHaveBeenCalled();
	});
});
