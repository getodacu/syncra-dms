import { INTERNAL_API_HEADER, apiBaseUrl, internalAPIHeaders } from "./internal-api";

type ServerFetch = typeof fetch;

export type BillingEntityType = "individual" | "company";

export type BillingProfileResponse = {
	id: string;
	user_id: string;
	entity_type: BillingEntityType;
	billing_name: string;
	billing_email: string;
	country_code: string;
	address_line1: string;
	address_line2?: string;
	city: string;
	region?: string;
	postal_code: string;
	fiscal_code?: string;
	registration_number?: string;
	created_at: string;
	updated_at: string;
};

export type CreditBalanceResponse = {
	user_id: string;
	available_credits: number;
};

export type BillingOrderResponse = {
	id: string;
	user_id: string;
	invoice?: BillingOrderInvoiceResponse | null;
	order_type: string;
	status: string;
	provider: string;
	pricing_tier: string;
	unit_amount_cents: number;
	credits: number;
	amount_cents: number;
	currency: string;
	provider_checkout_session_id?: string;
	provider_payment_intent_id?: string;
	created_at: string;
	updated_at: string;
	paid_at?: string;
	failed_at?: string;
	refunded_at?: string;
	canceled_at?: string;
};

export type BillingOrderInvoiceResponse = {
	id: string;
	invoice_serie: string;
	invoice_number: number;
	invoice_date: string;
	pdf_path?: string;
};

export type BillingInvoiceLineResponse = {
	name: string;
	quantity: number;
	unit_price: string;
	vat_percentage: string;
	total_vat_amount: string;
	total_amount: string;
};

export type BillingInvoiceResponse = {
	id: string;
	user_id?: string;
	order_id?: string;
	billing_profile_id?: string;
	billing_name: string;
	billing_email: string;
	billing_fiscal_code?: string;
	billing_profile_snapshot: unknown;
	lines: BillingInvoiceLineResponse[];
	net_amount: string;
	vat_amount: string;
	total_amount: string;
	invoice_date: string;
	invoice_serie: string;
	invoice_number: number;
	pdf_path?: string;
	email_delivery_claimed_at?: string;
	email_sent_at?: string;
	created_at: string;
	updated_at: string;
};

export type BillingInvoiceEmailDeliveryClaimStatus =
	| "claimed"
	| "claim_active"
	| "already_sent"
	| "not_ready";

export type BillingInvoiceEmailDeliveryClaimResponse = {
	status: BillingInvoiceEmailDeliveryClaimStatus;
	invoice?: BillingInvoiceResponse;
};

export type BillingInvoiceEmailDeliverySentResponse = {
	status: "sent";
	invoice: BillingInvoiceResponse;
};

export type BillingOrderStatus = "pending" | "paid" | "failed" | "refunded" | "canceled";

export type BillingOrderListResponse = {
	orders: BillingOrderResponse[];
	next_cursor: string | null;
};

export type BillingInvoicePDFResponse = {
	body: ReadableStream<Uint8Array> | null;
	headers: Headers;
	status: number;
};

export type FetchBillingInvoicePDFOptions = {
	download?: boolean;
};

export type BillingInvoiceEmailDeliveryInput = {
	invoiceId: string;
	userId: string;
};

export type CreditUsageHistoryEntryType = "purchase" | "debit";

export type CreditUsageHistoryEntryResponse = {
	id: string;
	created_at: string;
	entry_type: CreditUsageHistoryEntryType;
	credits_delta: number;
	related_order_id?: string;
	related_job_id?: string;
};

export type CreditUsageHistoryListResponse = {
	credit_usage_history: CreditUsageHistoryEntryResponse[];
	next_cursor: string | null;
};

export type CreateCreditOrderInput = {
	userId: string;
	credits: number;
};

export type UpsertBillingProfileInput = {
	userId: string;
	entityType: BillingEntityType;
	billingName: string;
	billingEmail: string;
	countryCode: string;
	addressLine1: string;
	addressLine2?: string;
	city: string;
	region?: string;
	postalCode: string;
	fiscalCode?: string;
	registrationNumber?: string;
};

export type ListCreditUsageHistoryOptions = {
	userId: string;
	type?: CreditUsageHistoryEntryType;
	createdFrom?: string;
	createdTo?: string;
	cursor?: string | null;
	size?: number | string;
	sort?: "asc" | "desc";
};

export type ListBillingOrdersOptions = {
	userId: string;
	status?: BillingOrderStatus;
	createdFrom?: string;
	createdTo?: string;
	cursor?: string | null;
	size?: number | string;
	sort?: "asc" | "desc";
};

export type AttachBillingOrderCheckoutSessionInput = {
	orderId: string;
	checkoutSessionId: string;
};

export type MarkBillingOrderPaidInput = {
	orderId: string;
	checkoutSessionId?: string;
	paymentIntentId?: string;
	paidAt: string;
};

export class BillingApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = "BillingApiError";
		this.status = status;
	}
}

function internalAPIToken() {
	const headers = internalAPIHeaders();
	const token = headers?.get(INTERNAL_API_HEADER);
	if (!token) {
		throw new BillingApiError(500, "Billing service is not configured");
	}
	return token;
}

function internalJSONHeaders() {
	return {
		"content-type": "application/json",
		[INTERNAL_API_HEADER]: internalAPIToken()
	};
}

function internalHeaders() {
	return {
		[INTERNAL_API_HEADER]: internalAPIToken()
	};
}

function isJsonObject(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function parseResponseJSON(text: string) {
	if (!text) return null;
	try {
		return JSON.parse(text) as unknown;
	} catch {
		return undefined;
	}
}

async function readResponseJSON(response: Response) {
	let text: string;
	try {
		text = await response.text();
	} catch {
		throw new BillingApiError(503, "Billing service unavailable");
	}

	return parseResponseJSON(text);
}

function errorMessage(data: unknown, fallback: string) {
	return data && typeof data === "object" && "error" in data && typeof data.error === "string"
		? data.error
		: fallback;
}

function isCreditBalanceResponse(value: unknown): value is CreditBalanceResponse {
	return (
		isJsonObject(value) &&
		typeof value.user_id === "string" &&
		typeof value.available_credits === "number" &&
		Number.isFinite(value.available_credits)
	);
}

function isBillingEntityType(value: unknown): value is BillingEntityType {
	return value === "individual" || value === "company";
}

function isBillingProfileResponse(value: unknown): value is BillingProfileResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.user_id === "string" &&
		isBillingEntityType(value.entity_type) &&
		typeof value.billing_name === "string" &&
		typeof value.billing_email === "string" &&
		typeof value.country_code === "string" &&
		typeof value.address_line1 === "string" &&
		!("address_line2" in value && typeof value.address_line2 !== "string") &&
		typeof value.city === "string" &&
		!("region" in value && typeof value.region !== "string") &&
		typeof value.postal_code === "string" &&
		!("fiscal_code" in value && typeof value.fiscal_code !== "string") &&
		!("registration_number" in value && typeof value.registration_number !== "string") &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
	);
}

function isBillingProfileEnvelope(
	value: unknown
): value is { profile: BillingProfileResponse | null } {
	return (
		isJsonObject(value) &&
		"profile" in value &&
		(value.profile === null || isBillingProfileResponse(value.profile))
	);
}

function isBillingOrderResponse(value: unknown): value is BillingOrderResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.user_id === "string" &&
		(!("invoice" in value) ||
			value.invoice === null ||
			isBillingOrderInvoiceResponse(value.invoice)) &&
		typeof value.order_type === "string" &&
		isBillingOrderStatus(value.status) &&
		typeof value.provider === "string" &&
		typeof value.pricing_tier === "string" &&
		typeof value.unit_amount_cents === "number" &&
		typeof value.credits === "number" &&
		typeof value.amount_cents === "number" &&
		typeof value.currency === "string" &&
		!("provider_checkout_session_id" in value && typeof value.provider_checkout_session_id !== "string") &&
		!("provider_payment_intent_id" in value && typeof value.provider_payment_intent_id !== "string") &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string" &&
		!("paid_at" in value && typeof value.paid_at !== "string") &&
		!("failed_at" in value && typeof value.failed_at !== "string") &&
		!("refunded_at" in value && typeof value.refunded_at !== "string") &&
		!("canceled_at" in value && typeof value.canceled_at !== "string")
	);
}

function isBillingOrderInvoiceResponse(value: unknown): value is BillingOrderInvoiceResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.invoice_serie === "string" &&
		typeof value.invoice_number === "number" &&
		Number.isFinite(value.invoice_number) &&
		typeof value.invoice_date === "string" &&
		(!("pdf_path" in value) || typeof value.pdf_path === "string")
	);
}

function isBillingInvoiceLineResponse(value: unknown): value is BillingInvoiceLineResponse {
	return (
		isJsonObject(value) &&
		typeof value.name === "string" &&
		typeof value.quantity === "number" &&
		Number.isFinite(value.quantity) &&
		typeof value.unit_price === "string" &&
		typeof value.vat_percentage === "string" &&
		typeof value.total_vat_amount === "string" &&
		typeof value.total_amount === "string"
	);
}

function isBillingInvoiceResponse(value: unknown): value is BillingInvoiceResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		(!("user_id" in value) || typeof value.user_id === "string") &&
		(!("order_id" in value) || typeof value.order_id === "string") &&
		(!("billing_profile_id" in value) || typeof value.billing_profile_id === "string") &&
		typeof value.billing_name === "string" &&
		typeof value.billing_email === "string" &&
		(!("billing_fiscal_code" in value) || typeof value.billing_fiscal_code === "string") &&
		"billing_profile_snapshot" in value &&
		Array.isArray(value.lines) &&
		value.lines.every(isBillingInvoiceLineResponse) &&
		typeof value.net_amount === "string" &&
		typeof value.vat_amount === "string" &&
		typeof value.total_amount === "string" &&
		typeof value.invoice_date === "string" &&
		typeof value.invoice_serie === "string" &&
		typeof value.invoice_number === "number" &&
		Number.isFinite(value.invoice_number) &&
		(!("pdf_path" in value) || typeof value.pdf_path === "string") &&
		(!("email_delivery_claimed_at" in value) ||
			typeof value.email_delivery_claimed_at === "string") &&
		(!("email_sent_at" in value) || typeof value.email_sent_at === "string") &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
	);
}

function isBillingInvoiceEmailDeliveryClaimStatus(
	value: unknown
): value is BillingInvoiceEmailDeliveryClaimStatus {
	return (
		value === "claimed" ||
		value === "claim_active" ||
		value === "already_sent" ||
		value === "not_ready"
	);
}

function isBillingInvoiceEmailDeliveryClaimResponse(
	value: unknown
): value is BillingInvoiceEmailDeliveryClaimResponse {
	return (
		isJsonObject(value) &&
		isBillingInvoiceEmailDeliveryClaimStatus(value.status) &&
		(!("invoice" in value) || isBillingInvoiceResponse(value.invoice))
	);
}

function isBillingInvoiceEmailDeliverySentResponse(
	value: unknown
): value is BillingInvoiceEmailDeliverySentResponse {
	return (
		isJsonObject(value) &&
		value.status === "sent" &&
		isBillingInvoiceResponse(value.invoice)
	);
}

function isBillingOrderStatus(value: unknown): value is BillingOrderStatus {
	return (
		value === "pending" ||
		value === "paid" ||
		value === "failed" ||
		value === "refunded" ||
		value === "canceled"
	);
}

function isBillingOrderListResponse(value: unknown): value is BillingOrderListResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.orders) &&
		value.orders.every(isBillingOrderResponse) &&
		(typeof value.next_cursor === "string" || value.next_cursor === null)
	);
}

function isCreditUsageHistoryEntryType(value: unknown): value is CreditUsageHistoryEntryType {
	return value === "purchase" || value === "debit";
}

function isCreditUsageHistoryEntryResponse(
	value: unknown
): value is CreditUsageHistoryEntryResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.created_at === "string" &&
		isCreditUsageHistoryEntryType(value.entry_type) &&
		typeof value.credits_delta === "number" &&
		Number.isFinite(value.credits_delta) &&
		!("related_order_id" in value && typeof value.related_order_id !== "string") &&
		!("related_job_id" in value && typeof value.related_job_id !== "string")
	);
}

function isCreditUsageHistoryListResponse(
	value: unknown
): value is CreditUsageHistoryListResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.credit_usage_history) &&
		value.credit_usage_history.every(isCreditUsageHistoryEntryResponse) &&
		(typeof value.next_cursor === "string" || value.next_cursor === null)
	);
}

function setOptionalSearchParam(url: URL, name: string, value: string | number | null | undefined) {
	if (value === null || value === undefined) return;
	const text = String(value).trim();
	if (text) url.searchParams.set(name, text);
}

async function requestBillingAPI(
	fetchFn: ServerFetch,
	url: string,
	init: RequestInit,
	fallbackError: string
) {
	let response: Response;
	try {
		response = await fetchFn(url, init);
	} catch {
		throw new BillingApiError(503, "Billing service unavailable");
	}

	const data = await readResponseJSON(response);
	if (!response.ok) {
		throw new BillingApiError(response.status, errorMessage(data, fallbackError));
	}

	return { response, data };
}

export function isBillingApiError(error: unknown): error is BillingApiError {
	return error instanceof BillingApiError;
}

export async function getCreditBalance(fetchFn: ServerFetch, userId: string) {
	const url = new URL(`${apiBaseUrl()}/api/billing/balance`);
	url.searchParams.set("user_id", userId);

	const { data } = await requestBillingAPI(
		fetchFn,
		url.toString(),
		{ method: "GET", headers: internalHeaders() },
		"Failed to load credit balance"
	);
	if (!isCreditBalanceResponse(data)) {
		throw new BillingApiError(502, "Invalid credit balance response");
	}
	return data;
}

export async function getBillingProfile(fetchFn: ServerFetch, userId: string) {
	const url = new URL(`${apiBaseUrl()}/api/billing/profile`);
	url.searchParams.set("user_id", userId);

	const { data } = await requestBillingAPI(
		fetchFn,
		url.toString(),
		{ method: "GET", headers: internalHeaders() },
		"Failed to load billing profile"
	);
	if (!isBillingProfileEnvelope(data)) {
		throw new BillingApiError(502, "Invalid billing profile response");
	}
	return data.profile;
}

export async function listCreditUsageHistory(
	fetchFn: ServerFetch,
	options: ListCreditUsageHistoryOptions
) {
	const url = new URL(`${apiBaseUrl()}/api/billing/credit-usage-history`);
	url.searchParams.set("user_id", options.userId);
	setOptionalSearchParam(url, "type", options.type);
	setOptionalSearchParam(url, "created_from", options.createdFrom);
	setOptionalSearchParam(url, "created_to", options.createdTo);
	setOptionalSearchParam(url, "cursor", options.cursor);
	setOptionalSearchParam(url, "size", options.size);
	setOptionalSearchParam(url, "sort", options.sort);

	const { data } = await requestBillingAPI(
		fetchFn,
		url.toString(),
		{ method: "GET", headers: internalHeaders() },
		"Failed to load credit usage history"
	);
	if (!isCreditUsageHistoryListResponse(data)) {
		throw new BillingApiError(502, "Invalid credit usage history response");
	}
	return data;
}

export async function listBillingOrders(fetchFn: ServerFetch, options: ListBillingOrdersOptions) {
	const url = new URL(`${apiBaseUrl()}/api/billing/orders`);
	url.searchParams.set("user_id", options.userId);
	setOptionalSearchParam(url, "status", options.status);
	setOptionalSearchParam(url, "created_from", options.createdFrom);
	setOptionalSearchParam(url, "created_to", options.createdTo);
	setOptionalSearchParam(url, "cursor", options.cursor);
	setOptionalSearchParam(url, "size", options.size);
	setOptionalSearchParam(url, "sort", options.sort);

	const { data } = await requestBillingAPI(
		fetchFn,
		url.toString(),
		{ method: "GET", headers: internalHeaders() },
		"Failed to load billing orders"
	);
	if (!isBillingOrderListResponse(data)) {
		throw new BillingApiError(502, "Invalid billing orders response");
	}
	return data;
}

export async function fetchBillingInvoicePDF(
	fetchFn: ServerFetch,
	userId: string,
	invoiceId: string,
	options: FetchBillingInvoicePDFOptions = {}
): Promise<BillingInvoicePDFResponse> {
	const url = new URL(`${apiBaseUrl()}/api/billing/invoices/${encodeURIComponent(invoiceId)}/pdf`);
	url.searchParams.set("user_id", userId);
	if (options.download) url.searchParams.set("download", "1");

	let response: Response;
	try {
		response = await fetchFn(url.toString(), {
			method: "GET",
			headers: internalHeaders()
		});
	} catch {
		throw new BillingApiError(503, "Billing service unavailable");
	}

	if (!response.ok) {
		const data = await readResponseJSON(response);
		throw new BillingApiError(response.status, errorMessage(data, "Failed to load invoice PDF"));
	}

	const headers = new Headers();
	for (const header of ["content-type", "content-disposition", "content-length"]) {
		const value = response.headers.get(header);
		if (value) headers.set(header, value);
	}

	return {
		body: response.body,
		headers,
		status: response.status
	};
}

export async function claimBillingInvoiceEmailDelivery(
	fetchFn: ServerFetch,
	input: BillingInvoiceEmailDeliveryInput
) {
	const { data } = await requestBillingAPI(
		fetchFn,
		`${apiBaseUrl()}/api/billing/invoices/${encodeURIComponent(input.invoiceId)}/email-delivery/claim`,
		{
			method: "POST",
			headers: internalJSONHeaders(),
			body: JSON.stringify({ user_id: input.userId })
		},
		"Failed to claim invoice email delivery"
	);
	if (!isBillingInvoiceEmailDeliveryClaimResponse(data)) {
		throw new BillingApiError(502, "Invalid invoice email delivery claim response");
	}
	return data;
}

export async function markBillingInvoiceEmailSent(
	fetchFn: ServerFetch,
	input: BillingInvoiceEmailDeliveryInput
) {
	const { data } = await requestBillingAPI(
		fetchFn,
		`${apiBaseUrl()}/api/billing/invoices/${encodeURIComponent(input.invoiceId)}/email-delivery/sent`,
		{
			method: "POST",
			headers: internalJSONHeaders(),
			body: JSON.stringify({ user_id: input.userId })
		},
		"Failed to mark invoice email sent"
	);
	if (!isBillingInvoiceEmailDeliverySentResponse(data)) {
		throw new BillingApiError(502, "Invalid invoice email sent response");
	}
	return data;
}

export async function upsertBillingProfile(fetchFn: ServerFetch, input: UpsertBillingProfileInput) {
	const { data } = await requestBillingAPI(
		fetchFn,
		`${apiBaseUrl()}/api/billing/profile`,
		{
			method: "PUT",
			headers: internalJSONHeaders(),
			body: JSON.stringify({
				user_id: input.userId,
				entity_type: input.entityType,
				billing_name: input.billingName,
				billing_email: input.billingEmail,
				country_code: input.countryCode,
				address_line1: input.addressLine1,
				address_line2: input.addressLine2,
				city: input.city,
				region: input.region,
				postal_code: input.postalCode,
				fiscal_code: input.fiscalCode,
				registration_number: input.registrationNumber
			})
		},
		"Failed to save billing profile"
	);
	if (!isBillingProfileResponse(data)) {
		throw new BillingApiError(502, "Invalid billing profile response");
	}
	return data;
}

export async function createCreditOrder(fetchFn: ServerFetch, input: CreateCreditOrderInput) {
	const { data } = await requestBillingAPI(
		fetchFn,
		`${apiBaseUrl()}/api/billing/orders`,
		{
			method: "POST",
			headers: internalJSONHeaders(),
			body: JSON.stringify({ user_id: input.userId, credits: input.credits })
		},
		"Failed to create credit order"
	);
	if (!isBillingOrderResponse(data)) {
		throw new BillingApiError(502, "Invalid billing order response");
	}
	return data;
}

export async function attachBillingOrderCheckoutSession(
	fetchFn: ServerFetch,
	input: AttachBillingOrderCheckoutSessionInput
) {
	await requestBillingAPI(
		fetchFn,
		`${apiBaseUrl()}/api/billing/orders/${encodeURIComponent(input.orderId)}/checkout-session`,
		{
			method: "POST",
			headers: internalJSONHeaders(),
			body: JSON.stringify({ checkout_session_id: input.checkoutSessionId })
		},
		"Failed to attach checkout session"
	);
}

export async function markBillingOrderPaid(fetchFn: ServerFetch, input: MarkBillingOrderPaidInput) {
	const { data } = await requestBillingAPI(
		fetchFn,
		`${apiBaseUrl()}/api/billing/orders/${encodeURIComponent(input.orderId)}/paid`,
		{
			method: "POST",
			headers: internalJSONHeaders(),
			body: JSON.stringify({
				checkout_session_id: input.checkoutSessionId,
				payment_intent_id: input.paymentIntentId,
				paid_at: input.paidAt
			})
		},
		"Failed to mark billing order paid"
	);
	if (!isBillingOrderResponse(data)) {
		throw new BillingApiError(502, "Invalid billing order response");
	}
	return data;
}

export async function markBillingOrderFailed(fetchFn: ServerFetch, orderId: string) {
	const { data } = await requestBillingAPI(
		fetchFn,
		`${apiBaseUrl()}/api/billing/orders/${encodeURIComponent(orderId)}/failed`,
		{
			method: "POST",
			headers: internalJSONHeaders(),
			body: "{}"
		},
		"Failed to mark billing order failed"
	);
	if (!isBillingOrderResponse(data)) {
		throw new BillingApiError(502, "Invalid billing order response");
	}
	return data;
}
