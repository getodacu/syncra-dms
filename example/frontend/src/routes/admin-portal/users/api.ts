import { publicApiErrorMessage } from "$lib/client/api-errors";

export type AdminUserRole = "user" | "admin";
export type AdminUserSort = "created_at" | "last_login_at";
export type AdminSortDirection = "asc" | "desc";

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

export type AdminBillingProfileResponse = {
	id: string;
	user_id: string;
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
	created_at: string;
	updated_at: string;
};

export type AdminUserDetailResponse = AdminUserResponse & {
	available_credits: number;
	billing_profile: AdminBillingProfileResponse | null;
};

export type AdminUserListResponse = {
	users: AdminUserResponse[];
	next_cursor: string | null;
};

export type AdminImpersonationResponse = {
	session: {
		id: string;
		token: string;
		userId: string;
		expiresAt: string;
		createdAt: string;
		updatedAt: string;
	};
	user: {
		id: string;
		name: string;
		email: string;
		role: AdminUserRole;
	};
	impersonation: {
		adminUser: { id: string; name: string; email: string; role: AdminUserRole };
		targetUser: { id: string; name: string; email: string; role: AdminUserRole };
		startedAt: string;
	} | null;
};

export type ListAdminUsersQuery = {
	search?: string;
	sort?: AdminUserSort;
	direction?: AdminSortDirection;
	cursor?: string | null;
	size?: number;
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

type ClientFetch = typeof fetch;

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function setOptionalSearchParam(url: URL, name: string, value: string | number | null | undefined) {
	if (value === null || value === undefined) return;
	const text = String(value).trim();
	if (text) url.searchParams.set(name, text);
}

export function buildAdminUsersPath(query: ListAdminUsersQuery) {
	const url = new URL("/api/admin/users", "http://localhost");
	setOptionalSearchParam(url, "search", query.search);
	setOptionalSearchParam(url, "sort", query.sort);
	setOptionalSearchParam(url, "direction", query.direction);
	setOptionalSearchParam(url, "cursor", query.cursor);
	setOptionalSearchParam(url, "size", query.size);
	return `${url.pathname}${url.search}`;
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

function responseErrorMessage(response: Response, value: unknown, fallback: string) {
	return publicApiErrorMessage(response.status, value, fallback);
}

function isAdminUserRole(value: unknown): value is AdminUserRole {
	return value === "user" || value === "admin";
}

function isAdminUserResponse(value: unknown): value is AdminUserResponse {
	return (
		isRecord(value) &&
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

function isBillingProfileResponse(value: unknown): value is AdminBillingProfileResponse {
	return (
		isRecord(value) &&
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
		isRecord(value) &&
		isAdminUserResponse(value) &&
		"available_credits" in value &&
		typeof value.available_credits === "number" &&
		"billing_profile" in value &&
		(value.billing_profile === null || isBillingProfileResponse(value.billing_profile))
	);
}

function isAdminUserListResponse(value: unknown): value is AdminUserListResponse {
	return (
		isRecord(value) &&
		Array.isArray(value.users) &&
		value.users.every(isAdminUserResponse) &&
		(typeof value.next_cursor === "string" || value.next_cursor === null)
	);
}

function isAuthUserShape(value: unknown): value is AdminImpersonationResponse["user"] {
	return (
		isRecord(value) &&
		typeof value.id === "string" &&
		typeof value.name === "string" &&
		typeof value.email === "string" &&
		isAdminUserRole(value.role)
	);
}

function isAdminImpersonationResponse(value: unknown): value is AdminImpersonationResponse {
	return (
		isRecord(value) &&
		isRecord(value.session) &&
		typeof value.session.id === "string" &&
		typeof value.session.token === "string" &&
		typeof value.session.userId === "string" &&
		typeof value.session.expiresAt === "string" &&
		typeof value.session.createdAt === "string" &&
		typeof value.session.updatedAt === "string" &&
		isAuthUserShape(value.user) &&
		(value.impersonation === null ||
			(isRecord(value.impersonation) &&
				isAuthUserShape(value.impersonation.adminUser) &&
				isAuthUserShape(value.impersonation.targetUser) &&
				typeof value.impersonation.startedAt === "string"))
	);
}

export async function fetchAdminUsers(fetchFn: ClientFetch, query: ListAdminUsersQuery) {
	const response = await fetchFn(buildAdminUsersPath(query), { method: "GET" });
	const json = await readResponseJSON(response);
	if (!response.ok) throw new Error(responseErrorMessage(response, json, "Failed to load users"));
	if (!isAdminUserListResponse(json)) throw new Error("Invalid admin users response");
	return json;
}

export async function fetchAdminUser(fetchFn: ClientFetch, userId: string) {
	const response = await fetchFn(`/api/admin/users/${encodeURIComponent(userId)}`, { method: "GET" });
	const json = await readResponseJSON(response);
	if (!response.ok) throw new Error(responseErrorMessage(response, json, "Failed to load user"));
	if (!isAdminUserDetailResponse(json)) throw new Error("Invalid admin user response");
	return json;
}

export async function startAdminUserImpersonation(fetchFn: ClientFetch, userId: string) {
	const response = await fetchFn(`/api/admin/users/${encodeURIComponent(userId)}/impersonation`, {
		method: "POST"
	});
	const json = await readResponseJSON(response);
	if (!response.ok) throw new Error(responseErrorMessage(response, json, "Failed to start impersonation"));
	if (!isAdminImpersonationResponse(json)) throw new Error("Invalid impersonation response");
	return json;
}

export async function updateAdminUser(fetchFn: ClientFetch, userId: string, input: { name?: string; email?: string }) {
	const response = await fetchFn(`/api/admin/users/${encodeURIComponent(userId)}`, {
		method: "PATCH",
		headers: { "content-type": "application/json" },
		body: JSON.stringify(input)
	});
	const json = await readResponseJSON(response);
	if (!response.ok) throw new Error(responseErrorMessage(response, json, "Failed to update user"));
	if (!isAdminUserResponse(json)) throw new Error("Invalid admin user response");
	return json;
}

export async function resetAdminUserPassword(fetchFn: ClientFetch, userId: string, password: string) {
	const response = await fetchFn(`/api/admin/users/${encodeURIComponent(userId)}/password`, {
		method: "POST",
		headers: { "content-type": "application/json" },
		body: JSON.stringify({ password })
	});
	const json = await readResponseJSON(response);
	if (!response.ok) throw new Error(responseErrorMessage(response, json, "Failed to reset password"));
	if (!isRecord(json) || json.ok !== true) throw new Error("Invalid password reset response");
	return json;
}

export async function adjustAdminUserBalance(fetchFn: ClientFetch, userId: string, creditsDelta: number) {
	const response = await fetchFn(`/api/admin/users/${encodeURIComponent(userId)}/balance-adjustment`, {
		method: "POST",
		headers: { "content-type": "application/json" },
		body: JSON.stringify({ credits_delta: creditsDelta })
	});
	const json = await readResponseJSON(response);
	if (!response.ok) throw new Error(responseErrorMessage(response, json, "Failed to adjust balance"));
	if (!isAdminUserDetailResponse(json)) throw new Error("Invalid admin user response");
	return json;
}

export async function upsertAdminUserBillingProfile(
	fetchFn: ClientFetch,
	userId: string,
	input: AdminBillingProfileInput
) {
	const response = await fetchFn(`/api/admin/users/${encodeURIComponent(userId)}/billing-profile`, {
		method: "PUT",
		headers: { "content-type": "application/json" },
		body: JSON.stringify(input)
	});
	const json = await readResponseJSON(response);
	if (!response.ok) throw new Error(responseErrorMessage(response, json, "Failed to save billing profile"));
	if (!isBillingProfileResponse(json)) throw new Error("Invalid billing profile response");
	return json;
}

export function formatAdminDate(value: string | null | undefined) {
	if (!value) return "Never";
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return "Invalid date";
	return new Intl.DateTimeFormat(undefined, {
		year: "numeric",
		month: "short",
		day: "2-digit",
		hour: "2-digit",
		minute: "2-digit"
	}).format(date);
}
