import { beforeEach, describe, expect, it, vi } from "vitest";

type FetchInput = Parameters<typeof fetch>[0];
type FetchInit = Parameters<typeof fetch>[1];

const { privateEnv } = vi.hoisted(() => ({
	privateEnv: {} as Record<string, string | undefined>
}));

vi.mock("$env/dynamic/private", () => ({ env: privateEnv }));

function balanceResponse() {
	return {
		user_id: "user-1",
		available_credits: 2500
	};
}

function orderResponse() {
	return {
		id: "order-1",
		user_id: "user-1",
		order_type: "credit_topup",
		status: "pending",
		provider: "stripe",
		pricing_tier: "tier_1",
		unit_amount_cents: 1000,
		credits: 1000,
		amount_cents: 1000,
		currency: "EUR",
		created_at: "2026-06-04T12:00:00Z",
		updated_at: "2026-06-04T12:00:00Z"
	};
}

function creditUsageHistoryResponse() {
	return {
		credit_usage_history: [
			{
				id: "entry-1",
				created_at: "2026-06-04T12:00:00Z",
				entry_type: "purchase",
				credits_delta: 1000,
				related_order_id: "order-1"
			},
			{
				id: "entry-2",
				created_at: "2026-06-04T13:00:00Z",
				entry_type: "debit",
				credits_delta: -3,
				related_job_id: "job-1"
			}
		],
		next_cursor: "cursor-1"
	};
}

function billingOrdersResponse() {
	return {
		orders: [
			{
				...orderResponse(),
				status: "paid",
				paid_at: "2026-06-04T12:05:00Z",
				invoice: {
					id: "invoice-1",
					invoice_serie: "SYN",
					invoice_number: 1,
					invoice_date: "2026-06-11",
					pdf_path: "/data/invoices/invoice-1.pdf"
				}
			}
		],
		next_cursor: "cursor-1"
	};
}

function profileResponse() {
	return {
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
		registration_number: "J40/1234/1999",
		created_at: "2026-06-05T09:00:00Z",
		updated_at: "2026-06-05T09:00:00Z"
	};
}

function invoiceResponse(overrides: Record<string, unknown> = {}) {
	return {
		id: "invoice-1",
		user_id: "user-1",
		order_id: "order-1",
		billing_profile_id: "profile-1",
		billing_name: "ICI Bucuresti",
		billing_email: "billing@example.com",
		billing_fiscal_code: "RO2785503",
		billing_profile_snapshot: { billing_name: "ICI Bucuresti" },
		lines: [
			{
				name: "SYNCRA SaaS 5000 credits",
				quantity: 1,
				unit_price: "47.50",
				vat_percentage: "0.00",
				total_vat_amount: "0.00",
				total_amount: "47.50"
			}
		],
		net_amount: "47.50",
		vat_amount: "0.00",
		total_amount: "47.50",
		invoice_date: "2026-06-13",
		invoice_serie: "SYN",
		invoice_number: 42,
		pdf_path: "/data/invoices/invoice-1.pdf",
		email_delivery_claimed_at: "2026-06-13T12:00:00Z",
		created_at: "2026-06-13T12:00:00Z",
		updated_at: "2026-06-13T12:00:00Z",
		...overrides
	};
}

describe("frontend billing server helper", () => {
	beforeEach(() => {
		vi.resetModules();
		for (const key of Object.keys(privateEnv)) delete privateEnv[key];
		process.env.SYNCRA_API_BASE_URL = "http://billing-api.test/";
		process.env.SYNCRA_INTERNAL_API_TOKEN = "internal-token";
		process.env.NODE_ENV = "test";
	});

	it("loads credit balance through the backend", async () => {
		const { getCreditBalance } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(balanceResponse()), { status: 200 });
		});

		await expect(getCreditBalance(fetchMock, "user-1")).resolves.toEqual(balanceResponse());
		expect(fetchMock).toHaveBeenCalledWith(
			"http://billing-api.test/api/billing/balance?user_id=user-1",
			expect.objectContaining({ method: "GET" })
		);
		expect(new Headers(fetchMock.mock.calls[0][1]?.headers).get("X-Syncra-Internal-Token")).toBe(
			"internal-token"
		);
	});

	it("creates credit orders with the internal API token", async () => {
		const { createCreditOrder } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(orderResponse()), { status: 201 });
		});

		await expect(createCreditOrder(fetchMock, { userId: "user-1", credits: 1000 })).resolves.toEqual(
			orderResponse()
		);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://billing-api.test/api/billing/orders",
			expect.objectContaining({
				method: "POST",
				headers: expect.objectContaining({
					"content-type": "application/json",
					"X-Syncra-Internal-Token": "internal-token"
				}),
				body: JSON.stringify({ user_id: "user-1", credits: 1000 })
			})
		);
	});

	it("attaches checkout sessions and marks orders paid through internal endpoints", async () => {
		const { attachBillingOrderCheckoutSession, markBillingOrderPaid } = await import("./billing");
		const fetchMock = vi.fn(async (input: FetchInput, _init?: FetchInit) => {
			const url = String(input);
			if (url.endsWith("/checkout-session")) return new Response(null, { status: 204 });
			return new Response(JSON.stringify({ ...orderResponse(), status: "paid" }), { status: 200 });
		});

		await attachBillingOrderCheckoutSession(fetchMock, {
			orderId: "order-1",
			checkoutSessionId: "cs_test_123"
		});
		await markBillingOrderPaid(fetchMock, {
			orderId: "order-1",
			checkoutSessionId: "cs_test_123",
			paymentIntentId: "pi_test_123",
			paidAt: "2026-06-04T12:05:00Z"
		});

		expect(fetchMock).toHaveBeenNthCalledWith(
			1,
			"http://billing-api.test/api/billing/orders/order-1/checkout-session",
			expect.objectContaining({
				method: "POST",
				headers: expect.objectContaining({ "X-Syncra-Internal-Token": "internal-token" }),
				body: JSON.stringify({ checkout_session_id: "cs_test_123" })
			})
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			2,
			"http://billing-api.test/api/billing/orders/order-1/paid",
			expect.objectContaining({
				method: "POST",
				headers: expect.objectContaining({ "X-Syncra-Internal-Token": "internal-token" }),
				body: JSON.stringify({
					checkout_session_id: "cs_test_123",
					payment_intent_id: "pi_test_123",
					paid_at: "2026-06-04T12:05:00Z"
				})
			})
		);
	});

	it("lists credit usage history through the backend with the internal API token", async () => {
		const { listCreditUsageHistory } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(creditUsageHistoryResponse()), { status: 200 });
		});

		await expect(
			listCreditUsageHistory(fetchMock, {
				userId: "user-1",
				type: "purchase",
				createdFrom: "2026-06-04T00:00:00Z",
				createdTo: "2026-06-05T00:00:00Z",
				cursor: "cursor-0",
				size: "50",
				sort: "desc"
			})
		).resolves.toEqual(creditUsageHistoryResponse());

		expect(fetchMock).toHaveBeenCalledWith(
			"http://billing-api.test/api/billing/credit-usage-history?user_id=user-1&type=purchase&created_from=2026-06-04T00%3A00%3A00Z&created_to=2026-06-05T00%3A00%3A00Z&cursor=cursor-0&size=50&sort=desc",
			expect.objectContaining({
				method: "GET",
				headers: {
					"X-Syncra-Internal-Token": "internal-token"
				}
			})
		);
	});

	it("lists billing orders through the backend with the internal API token", async () => {
		const { listBillingOrders } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(billingOrdersResponse()), { status: 200 });
		});

		await expect(
			listBillingOrders(fetchMock, {
				userId: "user-1",
				status: "paid",
				createdFrom: "2026-06-04T00:00:00Z",
				createdTo: "2026-06-05T00:00:00Z",
				cursor: "cursor-0",
				size: "50",
				sort: "desc"
			})
		).resolves.toEqual(billingOrdersResponse());

		expect(fetchMock).toHaveBeenCalledWith(
			"http://billing-api.test/api/billing/orders?user_id=user-1&status=paid&created_from=2026-06-04T00%3A00%3A00Z&created_to=2026-06-05T00%3A00%3A00Z&cursor=cursor-0&size=50&sort=desc",
			expect.objectContaining({
				method: "GET",
				headers: {
					"X-Syncra-Internal-Token": "internal-token"
				}
			})
		);
	});

	it("rejects invalid billing orders responses", async () => {
		const { listBillingOrders } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ orders: [{ status: "settled" }], next_cursor: null }), {
				status: 200
			});
		});

		await expect(listBillingOrders(fetchMock, { userId: "user-1" })).rejects.toMatchObject({
			status: 502,
			message: "Invalid billing orders response"
		});
	});

	it("rejects invalid billing order invoice metadata", async () => {
		const { listBillingOrders } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(
				JSON.stringify({
					orders: [{ ...orderResponse(), invoice: { id: "invoice-1", invoice_serie: "SYN" } }],
					next_cursor: null
				}),
				{ status: 200 }
			);
		});

		await expect(listBillingOrders(fetchMock, { userId: "user-1" })).rejects.toMatchObject({
			status: 502,
			message: "Invalid billing orders response"
		});
	});

	it("fetches billing invoice PDFs through the backend with the internal API token", async () => {
		const { fetchBillingInvoicePDF } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response("%PDF-test", {
				status: 200,
				headers: {
					"content-type": "application/pdf",
					"content-disposition": 'inline; filename="invoice-1.pdf"'
				}
			});
		});

		const result = await fetchBillingInvoicePDF(fetchMock, "user-1", "invoice-1", { download: true });

		expect(fetchMock).toHaveBeenCalledWith(
			"http://billing-api.test/api/billing/invoices/invoice-1/pdf?user_id=user-1&download=1",
			expect.objectContaining({ method: "GET" })
		);
		const headers = new Headers(fetchMock.mock.calls[0][1]?.headers);
		expect(headers.get("X-Syncra-Internal-Token")).toBe("internal-token");
		expect(result.status).toBe(200);
		expect(result.headers.get("content-type")).toBe("application/pdf");
		expect(result.headers.get("content-disposition")).toBe('inline; filename="invoice-1.pdf"');
		await expect(new Response(result.body).text()).resolves.toBe("%PDF-test");
	});

	it("claims invoice email delivery and marks invoice emails sent", async () => {
		const { claimBillingInvoiceEmailDelivery, markBillingInvoiceEmailSent } = await import("./billing");
		const fetchMock = vi.fn(async (input: FetchInput, _init?: FetchInit) => {
			const url = String(input);
			if (url.endsWith("/email-delivery/claim")) {
				return new Response(
					JSON.stringify({ status: "claimed", invoice: invoiceResponse() }),
					{ status: 200 }
				);
			}
			return new Response(
				JSON.stringify({
					status: "sent",
					invoice: invoiceResponse({ email_sent_at: "2026-06-13T12:01:00Z" })
				}),
				{ status: 200 }
			);
		});

		await expect(
			claimBillingInvoiceEmailDelivery(fetchMock, {
				invoiceId: "invoice-1",
				userId: "user-1"
			})
		).resolves.toEqual({ status: "claimed", invoice: invoiceResponse() });
		await expect(
			markBillingInvoiceEmailSent(fetchMock, {
				invoiceId: "invoice-1",
				userId: "user-1"
			})
		).resolves.toEqual({
			status: "sent",
			invoice: invoiceResponse({ email_sent_at: "2026-06-13T12:01:00Z" })
		});

		expect(fetchMock).toHaveBeenNthCalledWith(
			1,
			"http://billing-api.test/api/billing/invoices/invoice-1/email-delivery/claim",
			expect.objectContaining({
				method: "POST",
				headers: expect.objectContaining({
					"content-type": "application/json",
					"X-Syncra-Internal-Token": "internal-token"
				}),
				body: JSON.stringify({ user_id: "user-1" })
			})
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			2,
			"http://billing-api.test/api/billing/invoices/invoice-1/email-delivery/sent",
			expect.objectContaining({
				method: "POST",
				headers: expect.objectContaining({
					"content-type": "application/json",
					"X-Syncra-Internal-Token": "internal-token"
				}),
				body: JSON.stringify({ user_id: "user-1" })
			})
		);
	});

	it("rejects invalid invoice email delivery responses", async () => {
		const { claimBillingInvoiceEmailDelivery } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ status: "claimed", invoice: { id: "invoice-1" } }), {
				status: 200
			});
		});

		await expect(
			claimBillingInvoiceEmailDelivery(fetchMock, {
				invoiceId: "invoice-1",
				userId: "user-1"
			})
		).rejects.toMatchObject({
			status: 502,
			message: "Invalid invoice email delivery claim response"
		});
	});

	it("rejects invalid credit usage history responses", async () => {
		const { listCreditUsageHistory } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(
				JSON.stringify({ credit_usage_history: [{ entry_type: "credit" }], next_cursor: null }),
				{ status: 200 }
			);
		});

		await expect(listCreditUsageHistory(fetchMock, { userId: "user-1" })).rejects.toMatchObject({
			status: 502,
			message: "Invalid credit usage history response"
		});
	});

	it("loads billing profiles through the backend with the internal API token", async () => {
		const { getBillingProfile } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ profile: profileResponse() }), { status: 200 });
		});

		await expect(getBillingProfile(fetchMock, "user-1")).resolves.toEqual(profileResponse());
		expect(fetchMock).toHaveBeenCalledWith(
			"http://billing-api.test/api/billing/profile?user_id=user-1",
			expect.objectContaining({
				method: "GET",
				headers: {
					"X-Syncra-Internal-Token": "internal-token"
				}
			})
		);
	});

	it("returns null when the backend has no billing profile", async () => {
		const { getBillingProfile } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ profile: null }), { status: 200 });
		});

		await expect(getBillingProfile(fetchMock, "user-1")).resolves.toBeNull();
	});

	it("upserts billing profiles through the backend with snake case fields", async () => {
		const { upsertBillingProfile } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(profileResponse()), { status: 200 });
		});

		await expect(
			upsertBillingProfile(fetchMock, {
				userId: "user-1",
				entityType: "company",
				billingName: "ICI Bucuresti",
				billingEmail: "billing@example.com",
				countryCode: "RO",
				addressLine1: "Maresal Averescu 8-10",
				addressLine2: "Etaj 2",
				city: "Bucuresti",
				region: "Bucuresti",
				postalCode: "011455",
				fiscalCode: "RO2785503",
				registrationNumber: "J40/1234/1999"
			})
		).resolves.toEqual(profileResponse());

		expect(fetchMock).toHaveBeenCalledWith(
			"http://billing-api.test/api/billing/profile",
			expect.objectContaining({
				method: "PUT",
				headers: expect.objectContaining({
					"content-type": "application/json",
					"X-Syncra-Internal-Token": "internal-token"
				}),
				body: JSON.stringify({
					user_id: "user-1",
					entity_type: "company",
					billing_name: "ICI Bucuresti",
					billing_email: "billing@example.com",
					country_code: "RO",
					address_line1: "Maresal Averescu 8-10",
					address_line2: "Etaj 2",
					city: "Bucuresti",
					region: "Bucuresti",
					postal_code: "011455",
					fiscal_code: "RO2785503",
					registration_number: "J40/1234/1999"
				})
			})
		);
	});

	it("rejects invalid billing profile envelopes", async () => {
		const { getBillingProfile } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ billingProfile: profileResponse() }), { status: 200 });
		});

		await expect(getBillingProfile(fetchMock, "user-1")).rejects.toMatchObject({
			status: 502,
			message: "Invalid billing profile response"
		});
	});

	it("rejects invalid billing profile responses", async () => {
		const { upsertBillingProfile } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ ...profileResponse(), entity_type: "vendor" }), { status: 200 });
		});

		await expect(
			upsertBillingProfile(fetchMock, {
				userId: "user-1",
				entityType: "company",
				billingName: "ICI Bucuresti",
				billingEmail: "billing@example.com",
				countryCode: "RO",
				addressLine1: "Maresal Averescu 8-10",
				city: "Bucuresti",
				postalCode: "011455"
			})
		).rejects.toMatchObject({
			status: 502,
			message: "Invalid billing profile response"
		});
	});

	it("throws typed billing API errors from backend failures", async () => {
		const { BillingApiError, createCreditOrder, isBillingApiError } = await import("./billing");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ error: "credits must be a multiple of 1000" }), {
				status: 400
			});
		});

		await expect(createCreditOrder(fetchMock, { userId: "user-1", credits: 1500 })).rejects.toBeInstanceOf(
			BillingApiError
		);

		try {
			await createCreditOrder(fetchMock, { userId: "user-1", credits: 1500 });
		} catch (error) {
			expect(isBillingApiError(error)).toBe(true);
			expect(error).toMatchObject({
				status: 400,
				message: "credits must be a multiple of 1000"
			});
		}
	});

	it("rejects internal mutations when the internal API token is missing", async () => {
		const { createCreditOrder } = await import("./billing");
		delete process.env.SYNCRA_INTERNAL_API_TOKEN;
		const fetchMock = vi.fn();

		await expect(createCreditOrder(fetchMock, { userId: "user-1", credits: 1000 })).rejects.toMatchObject({
			status: 500,
			message: "Billing service is not configured"
		});
		expect(fetchMock).not.toHaveBeenCalled();
	});
});
