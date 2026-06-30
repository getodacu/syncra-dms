import { publicApiErrorMessage } from "$lib/client/api-errors";

import { buildBillingOrdersQueryPath, type BillingOrdersListQuery } from "./table-utils";

export type BillingOrderStatus = "pending" | "paid" | "failed" | "refunded" | "canceled";

export type BillingOrderInvoiceResponse = {
	id: string;
	invoice_serie: string;
	invoice_number: number;
	invoice_date: string;
	pdf_path?: string;
};

export type BillingOrderResponse = {
	id: string;
	user_id: string;
	invoice?: BillingOrderInvoiceResponse | null;
	order_type: string;
	status: BillingOrderStatus;
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

export type BillingOrdersListResponse = {
	orders: BillingOrderResponse[];
	next_cursor: string | null;
};

type ClientFetch = typeof fetch;

export function buildBillingInvoicePDFPath(invoiceId: string, options: { download?: boolean } = {}) {
	const path = `/api/billing/invoices/${encodeURIComponent(invoiceId)}/pdf`;
	return options.download ? `${path}?download=1` : path;
}

export async function fetchBillingOrders(
	fetchFn: ClientFetch,
	query: BillingOrdersListQuery
): Promise<BillingOrdersListResponse> {
	const response = await fetchFn(buildBillingOrdersQueryPath(query), { method: "GET" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to load billing orders"));
	}
	if (!isBillingOrdersListResponse(json)) {
		throw new Error("Invalid billing orders response");
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

function isBillingOrdersListResponse(value: unknown): value is BillingOrdersListResponse {
	if (!isRecord(value)) return false;
	if (!Array.isArray(value.orders)) return false;
	if (!(typeof value.next_cursor === "string" || value.next_cursor === null)) return false;

	return value.orders.every(isBillingOrderResponse);
}

function isBillingOrderInvoiceResponse(value: unknown): value is BillingOrderInvoiceResponse {
	return (
		isRecord(value) &&
		typeof value.id === "string" &&
		typeof value.invoice_serie === "string" &&
		typeof value.invoice_number === "number" &&
		Number.isFinite(value.invoice_number) &&
		typeof value.invoice_date === "string" &&
		(value.pdf_path === undefined || typeof value.pdf_path === "string")
	);
}

function isBillingOrderResponse(value: unknown): value is BillingOrderResponse {
	if (!isRecord(value)) return false;

	return (
		typeof value.id === "string" &&
		typeof value.user_id === "string" &&
		(value.invoice === undefined ||
			value.invoice === null ||
			isBillingOrderInvoiceResponse(value.invoice)) &&
		typeof value.order_type === "string" &&
		isBillingOrderStatus(value.status) &&
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

function isBillingOrderStatus(value: unknown): value is BillingOrderStatus {
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
