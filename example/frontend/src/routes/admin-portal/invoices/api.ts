import { publicApiErrorMessage } from "$lib/client/api-errors";

export type AdminInvoicesSortDirection = "asc" | "desc";

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
	pdf_path?: string;
	created_at: string;
	updated_at: string;
};

export type AdminBillingInvoicesListResponse = {
	invoices: AdminBillingInvoiceResponse[];
	next_cursor: string | null;
};

export type AdminBillingInvoicesListQuery = {
	search?: string;
	userId?: string;
	createdFrom?: string;
	createdTo?: string;
	cursor?: string | null;
	size?: number | string;
	sort?: AdminInvoicesSortDirection;
};

type ClientFetch = typeof fetch;

export function buildAdminBillingInvoicesPath(query: AdminBillingInvoicesListQuery = {}) {
	const params = new URLSearchParams();

	setNonEmptyParam(params, "search", query.search);
	setNonEmptyParam(params, "user_id", query.userId);
	setNonEmptyParam(params, "created_from", query.createdFrom);
	setNonEmptyParam(params, "created_to", query.createdTo);
	setNonEmptyParam(params, "cursor", query.cursor);
	setNonEmptyParam(params, "size", query.size);
	setNonEmptyParam(params, "sort", query.sort);

	const queryString = params.toString();
	return queryString ? `/api/admin/billing/invoices?${queryString}` : "/api/admin/billing/invoices";
}

export function buildAdminBillingInvoicePDFPath(
	invoiceId: string,
	options: { download?: boolean } = {}
) {
	const path = `/api/admin/billing/invoices/${encodeURIComponent(invoiceId)}/pdf`;
	return options.download ? `${path}?download=1` : path;
}

export async function fetchAdminBillingInvoices(
	fetchFn: ClientFetch,
	query: AdminBillingInvoicesListQuery
): Promise<AdminBillingInvoicesListResponse> {
	const response = await fetchFn(buildAdminBillingInvoicesPath(query), { method: "GET" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(
			publicApiErrorMessage(response.status, json, "Failed to load admin billing invoices")
		);
	}
	if (!isAdminBillingInvoicesListResponse(json)) {
		throw new Error("Invalid admin billing invoices response");
	}

	return json;
}

export async function generateAdminBillingInvoicePDF(
	fetchFn: ClientFetch,
	invoiceId: string
): Promise<AdminBillingInvoiceResponse> {
	const response = await fetchFn(buildAdminBillingInvoicePDFPath(invoiceId), { method: "POST" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to generate invoice PDF"));
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

function isAdminBillingInvoicesListResponse(value: unknown): value is AdminBillingInvoicesListResponse {
	if (!isRecord(value)) return false;
	if (!Array.isArray(value.invoices)) return false;
	if (!(typeof value.next_cursor === "string" || value.next_cursor === null)) return false;

	return value.invoices.every(isAdminBillingInvoiceResponse);
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
		(value.pdf_path === undefined || typeof value.pdf_path === "string") &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
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
