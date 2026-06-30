import { publicApiErrorMessage } from "$lib/client/api-errors";

export type AdminBillingOrderStatus = "pending" | "paid" | "failed" | "refunded" | "canceled";
export type AdminOrdersSortDirection = "asc" | "desc";

export type AdminBillingOrderUserResponse = {
	id: string;
	name: string;
	email: string;
};

export type AdminBillingOrderInvoiceResponse = {
	id: string;
	invoice_serie: string;
	invoice_number: number;
	invoice_date: string;
};

export type AdminBillingOrderResponse = {
	id: string;
	user_id: string;
	user: AdminBillingOrderUserResponse;
	invoice: AdminBillingOrderInvoiceResponse | null;
	order_type: string;
	status: AdminBillingOrderStatus;
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

export type AdminBillingInvoiceLineResponse = {
	name: string;
	quantity: number;
	unit_price: string;
	vat_percentage: string;
	total_vat_amount: string;
	total_amount: string;
};

export type AdminBillingInvoiceResponse = {
	id: string;
	user_id?: string;
	order_id?: string;
	billing_profile_id?: string;
	billing_name: string;
	billing_email: string;
	billing_fiscal_code?: string;
	billing_profile_snapshot: unknown;
	lines: AdminBillingInvoiceLineResponse[];
	net_amount: string;
	vat_amount: string;
	total_amount: string;
	invoice_date: string;
	invoice_serie: string;
	invoice_number: number;
	created_at: string;
	updated_at: string;
};

export type AdminBillingOrdersListResponse = {
	orders: AdminBillingOrderResponse[];
	next_cursor: string | null;
};

export type AdminBillingOrdersListQuery = {
	userId?: string;
	status?: AdminBillingOrderStatus;
	withoutInvoice?: boolean;
	createdFrom?: string;
	createdTo?: string;
	cursor?: string | null;
	size?: number | string;
	sort?: AdminOrdersSortDirection;
};

type ClientFetch = typeof fetch;

export function buildAdminBillingOrdersPath(query: AdminBillingOrdersListQuery = {}) {
	const params = new URLSearchParams();

	setNonEmptyParam(params, "user_id", query.userId);
	setNonEmptyParam(params, "status", query.status);
	setNonEmptyParam(params, "without_invoice", query.withoutInvoice ? "true" : undefined);
	setNonEmptyParam(params, "created_from", query.createdFrom);
	setNonEmptyParam(params, "created_to", query.createdTo);
	setNonEmptyParam(params, "cursor", query.cursor);
	setNonEmptyParam(params, "size", query.size);
	setNonEmptyParam(params, "sort", query.sort);

	const queryString = params.toString();
	return queryString ? `/api/admin/billing/orders?${queryString}` : "/api/admin/billing/orders";
}

export function buildGenerateAdminBillingOrderInvoicePath(orderId: string) {
	return `/api/admin/billing/orders/${encodeURIComponent(orderId)}/invoice`;
}

export async function fetchAdminBillingOrders(
	fetchFn: ClientFetch,
	query: AdminBillingOrdersListQuery
): Promise<AdminBillingOrdersListResponse> {
	const response = await fetchFn(buildAdminBillingOrdersPath(query), { method: "GET" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(
			publicApiErrorMessage(response.status, json, "Failed to load admin billing orders")
		);
	}
	if (!isAdminBillingOrdersListResponse(json)) {
		throw new Error("Invalid admin billing orders response");
	}

	return json;
}

export async function generateAdminBillingOrderInvoice(
	fetchFn: ClientFetch,
	orderId: string
): Promise<AdminBillingInvoiceResponse> {
	const response = await fetchFn(buildGenerateAdminBillingOrderInvoicePath(orderId), { method: "POST" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to generate invoice"));
	}
	if (!isAdminBillingInvoiceResponse(json)) {
		throw new Error("Invalid admin billing invoice response");
	}

	return json;
}

async function readResponseJSON(response: Response): Promise<unknown> {
	const text = await response.text();
	if (!text.trim()) return null;

	try {
		return JSON.parse(text);
	} catch {
		return null;
	}
}

function isAdminBillingOrdersListResponse(value: unknown): value is AdminBillingOrdersListResponse {
	if (!isRecord(value)) return false;
	if (!Array.isArray(value.orders)) return false;
	if (!(typeof value.next_cursor === "string" || value.next_cursor === null)) return false;

	return value.orders.every(isAdminBillingOrderResponse);
}

function isAdminBillingOrderResponse(value: unknown): value is AdminBillingOrderResponse {
	if (!isRecord(value)) return false;

	return (
		typeof value.id === "string" &&
		typeof value.user_id === "string" &&
		isAdminBillingOrderUserResponse(value.user) &&
		(value.invoice === null || isAdminBillingOrderInvoiceResponse(value.invoice)) &&
		typeof value.order_type === "string" &&
		isAdminBillingOrderStatus(value.status) &&
		typeof value.provider === "string" &&
		typeof value.pricing_tier === "string" &&
		typeof value.unit_amount_cents === "number" &&
		Number.isFinite(value.unit_amount_cents) &&
		typeof value.credits === "number" &&
		Number.isFinite(value.credits) &&
		typeof value.amount_cents === "number" &&
		Number.isFinite(value.amount_cents) &&
		typeof value.currency === "string" &&
		(value.provider_checkout_session_id === undefined ||
			typeof value.provider_checkout_session_id === "string") &&
		(value.provider_payment_intent_id === undefined ||
			typeof value.provider_payment_intent_id === "string") &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string" &&
		(value.paid_at === undefined || typeof value.paid_at === "string") &&
		(value.failed_at === undefined || typeof value.failed_at === "string") &&
		(value.refunded_at === undefined || typeof value.refunded_at === "string") &&
		(value.canceled_at === undefined || typeof value.canceled_at === "string")
	);
}

function isAdminBillingOrderInvoiceResponse(value: unknown): value is AdminBillingOrderInvoiceResponse {
	return (
		isRecord(value) &&
		typeof value.id === "string" &&
		typeof value.invoice_serie === "string" &&
		typeof value.invoice_number === "number" &&
		Number.isFinite(value.invoice_number) &&
		typeof value.invoice_date === "string"
	);
}

function isAdminBillingInvoiceLineResponse(value: unknown): value is AdminBillingInvoiceLineResponse {
	return (
		isRecord(value) &&
		typeof value.name === "string" &&
		typeof value.quantity === "number" &&
		Number.isFinite(value.quantity) &&
		typeof value.unit_price === "string" &&
		typeof value.vat_percentage === "string" &&
		typeof value.total_vat_amount === "string" &&
		typeof value.total_amount === "string"
	);
}

function isAdminBillingInvoiceResponse(value: unknown): value is AdminBillingInvoiceResponse {
	return (
		isRecord(value) &&
		typeof value.id === "string" &&
		(value.user_id === undefined || typeof value.user_id === "string") &&
		(value.order_id === undefined || typeof value.order_id === "string") &&
		(value.billing_profile_id === undefined || typeof value.billing_profile_id === "string") &&
		typeof value.billing_name === "string" &&
		typeof value.billing_email === "string" &&
		(value.billing_fiscal_code === undefined || typeof value.billing_fiscal_code === "string") &&
		"billing_profile_snapshot" in value &&
		Array.isArray(value.lines) &&
		value.lines.every(isAdminBillingInvoiceLineResponse) &&
		typeof value.net_amount === "string" &&
		typeof value.vat_amount === "string" &&
		typeof value.total_amount === "string" &&
		typeof value.invoice_date === "string" &&
		typeof value.invoice_serie === "string" &&
		typeof value.invoice_number === "number" &&
		Number.isFinite(value.invoice_number) &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
	);
}

function isAdminBillingOrderUserResponse(value: unknown): value is AdminBillingOrderUserResponse {
	return (
		isRecord(value) &&
		typeof value.id === "string" &&
		typeof value.name === "string" &&
		typeof value.email === "string"
	);
}

function isAdminBillingOrderStatus(value: unknown): value is AdminBillingOrderStatus {
	return (
		value === "pending" ||
		value === "paid" ||
		value === "failed" ||
		value === "refunded" ||
		value === "canceled"
	);
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function setNonEmptyParam(
	params: URLSearchParams,
	key: string,
	value: string | number | null | undefined
) {
	if (value === null || value === undefined) return;

	const stringValue = String(value).trim();
	if (!stringValue) return;

	params.set(key, stringValue);
}
