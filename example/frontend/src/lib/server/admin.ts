import { apiBaseUrl, internalAPIHeaders } from "./internal-api";
import type { BillingOrderResponse, BillingOrderStatus, BillingProfileResponse } from "./billing";
import type { AuthSessionPayload } from "./auth";

type ServerFetch = typeof fetch;

export type AdminUserSort = "created_at" | "last_login_at";
export type AdminSortDirection = "asc" | "desc";
export type AdminUserRole = "user" | "admin";

export type AdminUserResponse = {
	id: string;
	name: string;
	email: string;
	email_verified: boolean;
	role: AdminUserRole;
	image?: string | null;
	created_at: string;
	updated_at: string;
	last_login_at: string | null;
};

export type AdminUserDetailResponse = AdminUserResponse & {
	available_credits: number;
	billing_profile: BillingProfileResponse | null;
};

export type AdminUserListResponse = {
	users: AdminUserResponse[];
	next_cursor: string | null;
};

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

export type AdminBillingOrderResponse = BillingOrderResponse & {
	status: BillingOrderStatus;
	user: AdminBillingOrderUserResponse;
	invoice: AdminBillingOrderInvoiceResponse | null;
};

export type AdminBillingOrderListResponse = {
	orders: AdminBillingOrderResponse[];
	next_cursor: string | null;
};

export type AdminBillingInvoiceListResponse = {
	invoices: AdminBillingInvoiceResponse[];
	next_cursor: string | null;
};

export type AdminJSONRecipeResponse = {
	id: string;
	title: string;
	description: string;
	json: Record<string, unknown>;
	counter: number;
	category_id: string | null;
	category: AdminJSONRecipeCategoryResponse | null;
	created_at: string;
	updated_at: string;
};

export type AdminJSONRecipeCategoryTitle = {
	en: string;
	ro: string;
};

export type AdminJSONRecipeCategoryResponse = {
	id: string;
	title: AdminJSONRecipeCategoryTitle;
	created_at: string;
	updated_at: string;
};

export type AdminJSONRecipeCategoryListResponse = {
	categories: AdminJSONRecipeCategoryResponse[];
};

export type AdminJSONRecipeListResponse = {
	recipes: AdminJSONRecipeResponse[];
	next_cursor: string | null;
};

export type AdminJSONRecipeInput = {
	title: string;
	description: string;
	json: Record<string, unknown>;
	category_id?: string | null;
};

export type AdminJSONRecipeCategoryInput = {
	title: AdminJSONRecipeCategoryTitle;
};

export type ListAdminJSONRecipesOptions = {
	cursor?: string | null;
	size?: string | number;
	sort?: "asc" | "desc";
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
	pdf_path?: string;
	created_at: string;
	updated_at: string;
};

export type ListAdminUsersOptions = {
	search?: string;
	sort?: AdminUserSort;
	direction?: AdminSortDirection;
	cursor?: string | null;
	size?: string | number;
};

export type ListAdminBillingOrdersOptions = {
	userId?: string;
	status?: BillingOrderStatus;
	withoutInvoice?: boolean;
	createdFrom?: string;
	createdTo?: string;
	cursor?: string | null;
	size?: string | number;
	sort?: "asc" | "desc";
};

export type ListAdminBillingInvoicesOptions = {
	search?: string;
	userId?: string;
	createdFrom?: string;
	createdTo?: string;
	cursor?: string | null;
	size?: string | number;
	sort?: "asc" | "desc";
};

export type AdminBillingInvoicePDFResponse = {
	body: ReadableStream<Uint8Array> | null;
	headers: Headers;
	status: number;
};

export type FetchAdminBillingInvoicePDFOptions = {
	download?: boolean;
};

export type UpdateAdminUserInput = {
	name?: string;
	email?: string;
};

export type AdminBillingProfileInput = {
	entity_type: "individual" | "company";
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
};

export class AdminApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = "AdminApiError";
		this.status = status;
	}
}

export function isAdminApiError(error: unknown): error is AdminApiError {
	return error instanceof AdminApiError;
}

function adminHeaders(cookieHeader: string | null, json = false) {
	const headers = internalAPIHeaders(json ? { "content-type": "application/json" } : {});
	if (!headers) {
		throw new AdminApiError(500, "Admin service is not configured");
	}
	if (cookieHeader) headers.set("cookie", cookieHeader);
	return headers;
}

function setOptionalSearchParam(url: URL, name: string, value: string | number | null | undefined) {
	if (value === null || value === undefined) return;
	const text = String(value).trim();
	if (text) url.searchParams.set(name, text);
}

function parseResponseJSON(text: string) {
	if (!text) return null;
	try {
		return JSON.parse(text) as unknown;
	} catch {
		return undefined;
	}
}

function isJsonObject(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function errorMessage(data: unknown, fallback: string) {
	return data && typeof data === "object" && "error" in data && typeof data.error === "string"
		? data.error
		: fallback;
}

async function adminRequest(
	fetchFn: ServerFetch,
	path: string,
	cookieHeader: string | null,
	init: RequestInit,
	fallbackError: string
) {
	let response: Response;
	try {
		response = await fetchFn(`${apiBaseUrl()}${path}`, {
			...init,
			headers: init.headers
		});
	} catch {
		throw new AdminApiError(503, "Admin service unavailable");
	}

	const text = await response.text();
	const data = parseResponseJSON(text);
	if (!response.ok) {
		throw new AdminApiError(response.status, errorMessage(data, fallbackError));
	}
	return data;
}

function isAdminUserRole(value: unknown): value is AdminUserRole {
	return value === "user" || value === "admin";
}

function isAdminUserResponse(value: unknown): value is AdminUserResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.name === "string" &&
		typeof value.email === "string" &&
		typeof value.email_verified === "boolean" &&
		isAdminUserRole(value.role) &&
		(!("image" in value) || value.image === null || typeof value.image === "string") &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string" &&
		(value.last_login_at === null || typeof value.last_login_at === "string")
	);
}

function isBillingProfileResponse(value: unknown): value is BillingProfileResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.user_id === "string" &&
		(value.entity_type === "individual" || value.entity_type === "company") &&
		typeof value.billing_name === "string" &&
		typeof value.billing_email === "string" &&
		typeof value.country_code === "string" &&
		typeof value.address_line1 === "string" &&
		!("address_line2" in value && typeof value.address_line2 !== "string") &&
		typeof value.city === "string" &&
		!("region" in value && typeof value.region !== "string") &&
		typeof value.postal_code === "string" &&
		!("fiscal_code" in value && typeof value.fiscal_code !== "string") &&
		!("registration_number" in value && typeof value.registration_number !== "string")
	);
}

function isAdminUserDetailResponse(value: unknown): value is AdminUserDetailResponse {
	return (
		isJsonObject(value) &&
		isAdminUserResponse(value) &&
		"available_credits" in value &&
		typeof value.available_credits === "number" &&
		"billing_profile" in value &&
		(value.billing_profile === null || isBillingProfileResponse(value.billing_profile))
	);
}

function isAdminUserListResponse(value: unknown): value is AdminUserListResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.users) &&
		value.users.every(isAdminUserResponse) &&
		(typeof value.next_cursor === "string" || value.next_cursor === null)
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

function isAdminBillingOrderUserResponse(value: unknown): value is AdminBillingOrderUserResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.name === "string" &&
		typeof value.email === "string"
	);
}

function isAdminBillingOrderInvoiceResponse(value: unknown): value is AdminBillingOrderInvoiceResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.invoice_serie === "string" &&
		typeof value.invoice_number === "number" &&
		Number.isFinite(value.invoice_number) &&
		typeof value.invoice_date === "string"
	);
}

function isAdminBillingOrderResponse(value: unknown): value is AdminBillingOrderResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.user_id === "string" &&
		isAdminBillingOrderUserResponse(value.user) &&
		(value.invoice === null || isAdminBillingOrderInvoiceResponse(value.invoice)) &&
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
		(!("provider_checkout_session_id" in value) ||
			typeof value.provider_checkout_session_id === "string") &&
		(!("provider_payment_intent_id" in value) ||
			typeof value.provider_payment_intent_id === "string") &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string" &&
		(!("paid_at" in value) || typeof value.paid_at === "string") &&
		(!("failed_at" in value) || typeof value.failed_at === "string") &&
		(!("refunded_at" in value) || typeof value.refunded_at === "string") &&
		(!("canceled_at" in value) || typeof value.canceled_at === "string")
	);
}

function isAdminBillingInvoiceLineResponse(value: unknown): value is AdminBillingInvoiceLineResponse {
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

function isAdminBillingInvoiceResponse(value: unknown): value is AdminBillingInvoiceResponse {
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
		value.lines.every(isAdminBillingInvoiceLineResponse) &&
		typeof value.net_amount === "string" &&
		typeof value.vat_amount === "string" &&
		typeof value.total_amount === "string" &&
		typeof value.invoice_date === "string" &&
		typeof value.invoice_serie === "string" &&
		typeof value.invoice_number === "number" &&
		Number.isFinite(value.invoice_number) &&
		(!("pdf_path" in value) || typeof value.pdf_path === "string") &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
	);
}

function isAdminBillingOrderListResponse(value: unknown): value is AdminBillingOrderListResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.orders) &&
		value.orders.every(isAdminBillingOrderResponse) &&
		(typeof value.next_cursor === "string" || value.next_cursor === null)
	);
}

function isAdminBillingInvoiceListResponse(value: unknown): value is AdminBillingInvoiceListResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.invoices) &&
		value.invoices.every(isAdminBillingInvoiceResponse) &&
		(typeof value.next_cursor === "string" || value.next_cursor === null)
	);
}

function isAdminJSONRecipeResponse(value: unknown): value is AdminJSONRecipeResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.title === "string" &&
		typeof value.description === "string" &&
		isJsonObject(value.json) &&
		typeof value.counter === "number" &&
		Number.isFinite(value.counter) &&
		(typeof value.category_id === "string" || value.category_id === null) &&
		(value.category === null || isAdminJSONRecipeCategoryResponse(value.category)) &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
	);
}

function isAdminJSONRecipeCategoryTitle(value: unknown): value is AdminJSONRecipeCategoryTitle {
	return isJsonObject(value) && typeof value.en === "string" && typeof value.ro === "string";
}

function isAdminJSONRecipeCategoryResponse(value: unknown): value is AdminJSONRecipeCategoryResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		isAdminJSONRecipeCategoryTitle(value.title) &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
	);
}

function isAdminJSONRecipeListResponse(value: unknown): value is AdminJSONRecipeListResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.recipes) &&
		value.recipes.every(isAdminJSONRecipeResponse) &&
		(typeof value.next_cursor === "string" || value.next_cursor === null)
	);
}

function isAdminJSONRecipeCategoryListResponse(value: unknown): value is AdminJSONRecipeCategoryListResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.categories) &&
		value.categories.every(isAdminJSONRecipeCategoryResponse)
	);
}

function isPasswordResetResponse(value: unknown): value is { ok: true } {
	return isJsonObject(value) && value.ok === true;
}

function isAuthUserResponse(value: unknown) {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.name === "string" &&
		typeof value.email === "string" &&
		typeof value.emailVerified === "boolean" &&
		(value.role === "user" || value.role === "admin") &&
		(!("image" in value) || value.image === null || typeof value.image === "string") &&
		(!("lastLoginAt" in value) || value.lastLoginAt === null || typeof value.lastLoginAt === "string") &&
		typeof value.createdAt === "string" &&
		typeof value.updatedAt === "string"
	);
}

function isAuthSessionPayload(value: unknown): value is AuthSessionPayload {
	return (
		isJsonObject(value) &&
		isJsonObject(value.session) &&
		typeof value.session.id === "string" &&
		typeof value.session.token === "string" &&
		typeof value.session.userId === "string" &&
		typeof value.session.expiresAt === "string" &&
		typeof value.session.createdAt === "string" &&
		typeof value.session.updatedAt === "string" &&
		isAuthUserResponse(value.user) &&
		(!("impersonation" in value) ||
			value.impersonation === null ||
			(isJsonObject(value.impersonation) &&
				isAuthUserResponse(value.impersonation.adminUser) &&
				isAuthUserResponse(value.impersonation.targetUser) &&
				typeof value.impersonation.startedAt === "string"))
	);
}

export async function listAdminUsers(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	options: ListAdminUsersOptions
) {
	const url = new URL(`${apiBaseUrl()}/api/admin/users`);
	setOptionalSearchParam(url, "search", options.search);
	setOptionalSearchParam(url, "sort", options.sort);
	setOptionalSearchParam(url, "direction", options.direction);
	setOptionalSearchParam(url, "cursor", options.cursor);
	setOptionalSearchParam(url, "size", options.size);

	const data = await adminRequest(
		fetchFn,
		`${url.pathname}${url.search}`,
		cookieHeader,
		{ method: "GET", headers: adminHeaders(cookieHeader) },
		"Failed to load admin users"
	);
	if (!isAdminUserListResponse(data)) {
		throw new AdminApiError(502, "Invalid admin user list response");
	}
	return data;
}

export async function listAdminBillingOrders(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	options: ListAdminBillingOrdersOptions
) {
	const url = new URL(`${apiBaseUrl()}/api/admin/billing/orders`);
	setOptionalSearchParam(url, "user_id", options.userId);
	setOptionalSearchParam(url, "status", options.status);
	setOptionalSearchParam(url, "without_invoice", options.withoutInvoice ? "true" : undefined);
	setOptionalSearchParam(url, "created_from", options.createdFrom);
	setOptionalSearchParam(url, "created_to", options.createdTo);
	setOptionalSearchParam(url, "cursor", options.cursor);
	setOptionalSearchParam(url, "size", options.size);
	setOptionalSearchParam(url, "sort", options.sort);

	const data = await adminRequest(
		fetchFn,
		`${url.pathname}${url.search}`,
		cookieHeader,
		{ method: "GET", headers: adminHeaders(cookieHeader) },
		"Failed to load admin billing orders"
	);
	if (!isAdminBillingOrderListResponse(data)) {
		throw new AdminApiError(502, "Invalid admin billing orders response");
	}
	return data;
}

export async function listAdminBillingInvoices(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	options: ListAdminBillingInvoicesOptions
) {
	const url = new URL(`${apiBaseUrl()}/api/admin/billing/invoices`);
	setOptionalSearchParam(url, "search", options.search);
	setOptionalSearchParam(url, "user_id", options.userId);
	setOptionalSearchParam(url, "created_from", options.createdFrom);
	setOptionalSearchParam(url, "created_to", options.createdTo);
	setOptionalSearchParam(url, "cursor", options.cursor);
	setOptionalSearchParam(url, "size", options.size);
	setOptionalSearchParam(url, "sort", options.sort);

	const data = await adminRequest(
		fetchFn,
		`${url.pathname}${url.search}`,
		cookieHeader,
		{ method: "GET", headers: adminHeaders(cookieHeader) },
		"Failed to load admin billing invoices"
	);
	if (!isAdminBillingInvoiceListResponse(data)) {
		throw new AdminApiError(502, "Invalid admin billing invoices response");
	}
	return data;
}

export async function listAdminJSONRecipes(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	options: ListAdminJSONRecipesOptions = {}
) {
	const url = new URL(`${apiBaseUrl()}/api/admin/json-recipes`);
	setOptionalSearchParam(url, "cursor", options.cursor);
	setOptionalSearchParam(url, "size", options.size);
	setOptionalSearchParam(url, "sort", options.sort);

	const data = await adminRequest(
		fetchFn,
		`${url.pathname}${url.search}`,
		cookieHeader,
		{ method: "GET", headers: adminHeaders(cookieHeader) },
		"Failed to load admin JSON recipes"
	);
	if (!isAdminJSONRecipeListResponse(data)) {
		throw new AdminApiError(502, "Invalid admin JSON recipe response");
	}
	return data;
}

export async function listAdminJSONRecipeCategories(
	fetchFn: ServerFetch,
	cookieHeader: string | null
) {
	const data = await adminRequest(
		fetchFn,
		"/api/admin/json-recipe-categories",
		cookieHeader,
		{ method: "GET", headers: adminHeaders(cookieHeader) },
		"Failed to load admin JSON recipe categories"
	);
	if (!isAdminJSONRecipeCategoryListResponse(data)) {
		throw new AdminApiError(502, "Invalid admin JSON recipe categories response");
	}
	return data;
}

export async function getAdminJSONRecipe(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	recipeId: string
) {
	const data = await adminRequest(
		fetchFn,
		`/api/admin/json-recipes/${encodeURIComponent(recipeId)}`,
		cookieHeader,
		{ method: "GET", headers: adminHeaders(cookieHeader) },
		"Failed to load admin JSON recipe"
	);
	if (!isAdminJSONRecipeResponse(data)) {
		throw new AdminApiError(502, "Invalid admin JSON recipe response");
	}
	return data;
}

export async function getAdminJSONRecipeCategory(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	categoryId: string
) {
	const data = await adminRequest(
		fetchFn,
		`/api/admin/json-recipe-categories/${encodeURIComponent(categoryId)}`,
		cookieHeader,
		{ method: "GET", headers: adminHeaders(cookieHeader) },
		"Failed to load admin JSON recipe category"
	);
	if (!isAdminJSONRecipeCategoryResponse(data)) {
		throw new AdminApiError(502, "Invalid admin JSON recipe category response");
	}
	return data;
}

export async function createAdminJSONRecipe(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	input: AdminJSONRecipeInput
) {
	const data = await adminRequest(
		fetchFn,
		"/api/admin/json-recipes",
		cookieHeader,
		{
			method: "POST",
			headers: adminHeaders(cookieHeader, true),
			body: JSON.stringify(input)
		},
		"Failed to create admin JSON recipe"
	);
	if (!isAdminJSONRecipeResponse(data)) {
		throw new AdminApiError(502, "Invalid admin JSON recipe response");
	}
	return data;
}

export async function createAdminJSONRecipeCategory(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	input: AdminJSONRecipeCategoryInput
) {
	const data = await adminRequest(
		fetchFn,
		"/api/admin/json-recipe-categories",
		cookieHeader,
		{
			method: "POST",
			headers: adminHeaders(cookieHeader, true),
			body: JSON.stringify(input)
		},
		"Failed to create admin JSON recipe category"
	);
	if (!isAdminJSONRecipeCategoryResponse(data)) {
		throw new AdminApiError(502, "Invalid admin JSON recipe category response");
	}
	return data;
}

export async function updateAdminJSONRecipe(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	recipeId: string,
	input: AdminJSONRecipeInput
) {
	const data = await adminRequest(
		fetchFn,
		`/api/admin/json-recipes/${encodeURIComponent(recipeId)}`,
		cookieHeader,
		{
			method: "PUT",
			headers: adminHeaders(cookieHeader, true),
			body: JSON.stringify(input)
		},
		"Failed to update admin JSON recipe"
	);
	if (!isAdminJSONRecipeResponse(data)) {
		throw new AdminApiError(502, "Invalid admin JSON recipe response");
	}
	return data;
}

export async function updateAdminJSONRecipeCategory(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	categoryId: string,
	input: AdminJSONRecipeCategoryInput
) {
	const data = await adminRequest(
		fetchFn,
		`/api/admin/json-recipe-categories/${encodeURIComponent(categoryId)}`,
		cookieHeader,
		{
			method: "PUT",
			headers: adminHeaders(cookieHeader, true),
			body: JSON.stringify(input)
		},
		"Failed to update admin JSON recipe category"
	);
	if (!isAdminJSONRecipeCategoryResponse(data)) {
		throw new AdminApiError(502, "Invalid admin JSON recipe category response");
	}
	return data;
}

export async function deleteAdminJSONRecipe(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	recipeId: string
) {
	await adminRequest(
		fetchFn,
		`/api/admin/json-recipes/${encodeURIComponent(recipeId)}`,
		cookieHeader,
		{ method: "DELETE", headers: adminHeaders(cookieHeader) },
		"Failed to delete admin JSON recipe"
	);
}

export async function deleteAdminJSONRecipeCategory(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	categoryId: string
) {
	await adminRequest(
		fetchFn,
		`/api/admin/json-recipe-categories/${encodeURIComponent(categoryId)}`,
		cookieHeader,
		{ method: "DELETE", headers: adminHeaders(cookieHeader) },
		"Failed to delete admin JSON recipe category"
	);
}

export async function generateAdminBillingOrderInvoice(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	orderId: string
) {
	const data = await adminRequest(
		fetchFn,
		`/api/admin/billing/orders/${encodeURIComponent(orderId)}/invoice`,
		cookieHeader,
		{ method: "POST", headers: adminHeaders(cookieHeader) },
		"Failed to generate invoice"
	);
	if (!isAdminBillingInvoiceResponse(data)) {
		throw new AdminApiError(502, "Invalid admin billing invoice response");
	}
	return data;
}

export async function generateAdminBillingInvoicePDF(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	invoiceId: string
) {
	const data = await adminRequest(
		fetchFn,
		"/api/billing/generate-invoice-pdf",
		cookieHeader,
		{
			method: "POST",
			headers: adminHeaders(cookieHeader, true),
			body: JSON.stringify({ invoice_id: invoiceId })
		},
		"Failed to generate invoice PDF"
	);
	if (!isAdminBillingInvoiceResponse(data)) {
		throw new AdminApiError(502, "Invalid admin billing invoice response");
	}
	return data;
}

export async function fetchAdminBillingInvoicePDF(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	invoiceId: string,
	options: FetchAdminBillingInvoicePDFOptions = {}
): Promise<AdminBillingInvoicePDFResponse> {
	const url = new URL(`${apiBaseUrl()}/static/invoice/${encodeURIComponent(invoiceId)}.pdf`);
	if (options.download) url.searchParams.set("download", "1");

	let response: Response;
	try {
		response = await fetchFn(url.toString(), {
			method: "GET",
			headers: adminHeaders(cookieHeader)
		});
	} catch {
		throw new AdminApiError(503, "Admin service unavailable");
	}

	if (!response.ok) {
		const text = await response.text();
		const data = parseResponseJSON(text);
		throw new AdminApiError(response.status, errorMessage(data, "Failed to load invoice PDF"));
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

export async function getAdminUser(fetchFn: ServerFetch, cookieHeader: string | null, userId: string) {
	const data = await adminRequest(
		fetchFn,
		`/api/admin/users/${encodeURIComponent(userId)}`,
		cookieHeader,
		{ method: "GET", headers: adminHeaders(cookieHeader) },
		"Failed to load admin user"
	);
	if (!isAdminUserDetailResponse(data)) {
		throw new AdminApiError(502, "Invalid admin user response");
	}
	return data;
}

export async function updateAdminUser(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	userId: string,
	input: UpdateAdminUserInput
) {
	const data = await adminRequest(
		fetchFn,
		`/api/admin/users/${encodeURIComponent(userId)}`,
		cookieHeader,
		{
			method: "PATCH",
			headers: adminHeaders(cookieHeader, true),
			body: JSON.stringify(input)
		},
		"Failed to update admin user"
	);
	if (!isAdminUserResponse(data)) {
		throw new AdminApiError(502, "Invalid admin user response");
	}
	return data;
}

export async function startAdminUserImpersonation(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	userId: string
) {
	const data = await adminRequest(
		fetchFn,
		`/api/admin/users/${encodeURIComponent(userId)}/impersonation`,
		cookieHeader,
		{ method: "POST", headers: adminHeaders(cookieHeader) },
		"Failed to start impersonation"
	);
	if (!isAuthSessionPayload(data)) {
		throw new AdminApiError(502, "Invalid impersonation response");
	}
	return data;
}

export async function stopAdminImpersonation(fetchFn: ServerFetch, cookieHeader: string | null) {
	const data = await adminRequest(
		fetchFn,
		"/api/admin/impersonation/stop",
		cookieHeader,
		{ method: "POST", headers: adminHeaders(cookieHeader) },
		"Failed to stop impersonation"
	);
	if (!isAuthSessionPayload(data)) {
		throw new AdminApiError(502, "Invalid impersonation response");
	}
	return data;
}

export async function resetAdminUserPassword(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	userId: string,
	password: string
) {
	const data = await adminRequest(
		fetchFn,
		`/api/admin/users/${encodeURIComponent(userId)}/password`,
		cookieHeader,
		{
			method: "POST",
			headers: adminHeaders(cookieHeader, true),
			body: JSON.stringify({ password })
		},
		"Failed to reset password"
	);
	if (!isPasswordResetResponse(data)) {
		throw new AdminApiError(502, "Invalid password reset response");
	}
	return data;
}

export async function adjustAdminUserBalance(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	userId: string,
	creditsDelta: number
) {
	const data = await adminRequest(
		fetchFn,
		`/api/admin/users/${encodeURIComponent(userId)}/balance-adjustment`,
		cookieHeader,
		{
			method: "POST",
			headers: adminHeaders(cookieHeader, true),
			body: JSON.stringify({ credits_delta: creditsDelta })
		},
		"Failed to adjust user balance"
	);
	if (!isAdminUserDetailResponse(data)) {
		throw new AdminApiError(502, "Invalid admin user response");
	}
	return data;
}

export async function upsertAdminUserBillingProfile(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	userId: string,
	input: AdminBillingProfileInput
) {
	const data = await adminRequest(
		fetchFn,
		`/api/admin/users/${encodeURIComponent(userId)}/billing-profile`,
		cookieHeader,
		{
			method: "PUT",
			headers: adminHeaders(cookieHeader, true),
			body: JSON.stringify(input)
		},
		"Failed to save billing profile"
	);
	if (!isBillingProfileResponse(data)) {
		throw new AdminApiError(502, "Invalid billing profile response");
	}
	return data;
}
