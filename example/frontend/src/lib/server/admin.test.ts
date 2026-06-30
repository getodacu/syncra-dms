import { beforeEach, describe, expect, it, vi } from "vitest";

type FetchInput = Parameters<typeof fetch>[0];
type FetchInit = Parameters<typeof fetch>[1];

const { privateEnv } = vi.hoisted(() => ({
	privateEnv: {} as Record<string, string | undefined>
}));

vi.mock("$env/dynamic/private", () => ({ env: privateEnv }));

function adminUserResponse() {
	return {
		id: "user-1",
		name: "Ada Lovelace",
		email: "ada@example.com",
		email_verified: true,
		role: "user",
		image: null,
		created_at: "2026-06-01T00:00:00Z",
		updated_at: "2026-06-02T00:00:00Z",
		last_login_at: "2026-06-03T00:00:00Z"
	};
}

function adminUserDetailResponse() {
	return {
		...adminUserResponse(),
		available_credits: 125,
		billing_profile: null
	};
}

function authUserResponse(id: string, email: string, role: "user" | "admin" = "user") {
	return {
		id,
		name: role === "admin" ? "Admin User" : "Ada Lovelace",
		email,
		emailVerified: true,
		role,
		image: null,
		lastLoginAt: null,
		createdAt: "2026-06-01T00:00:00Z",
		updatedAt: "2026-06-02T00:00:00Z"
	};
}

function authSessionPayloadResponse(extra: Record<string, unknown> = {}) {
	return {
		session: {
			id: "session-1",
			token: "token-1",
			userId: "admin-1",
			expiresAt: "2026-07-01T00:00:00Z",
			createdAt: "2026-06-01T00:00:00Z",
			updatedAt: "2026-06-14T10:00:00Z"
		},
		user: authUserResponse("user-1", "ada@example.com", "user"),
		impersonation: {
			adminUser: authUserResponse("admin-1", "admin@example.com", "admin"),
			targetUser: authUserResponse("user-1", "ada@example.com", "user"),
			startedAt: "2026-06-14T10:00:00Z"
		},
		...extra
	};
}

function billingProfileResponse() {
	return {
		id: "profile-1",
		user_id: "user-1",
		entity_type: "company",
		billing_name: "Syncra SRL",
		billing_email: "billing@example.com",
		country_code: "RO",
		address_line1: "Main Street 1",
		city: "Bucharest",
		postal_code: "010101",
		fiscal_code: "RO123",
		created_at: "2026-06-05T00:00:00Z",
		updated_at: "2026-06-05T00:00:00Z"
	};
}

function adminBillingOrdersResponse() {
	return {
		orders: [
			{
				id: "order-1",
				user_id: "user-1",
				user: {
					id: "user-1",
					name: "Ada Lovelace",
					email: "ada@example.com"
				},
				invoice: null,
				order_type: "credit_topup",
				status: "paid",
				provider: "stripe",
				pricing_tier: "tier_2",
				unit_amount_cents: 950,
				credits: 5000,
				amount_cents: 4750,
				currency: "EUR",
				created_at: "2026-06-04T12:00:00Z",
				updated_at: "2026-06-04T12:05:00Z",
				paid_at: "2026-06-04T12:05:00Z"
			}
		],
		next_cursor: "cursor-1"
	};
}

function adminBillingInvoiceResponse() {
	return {
		id: "invoice-1",
		user_id: "user-1",
		order_id: "order-1",
		billing_profile_id: "profile-1",
		billing_name: "Ada Lovelace",
		billing_email: "ada@example.com",
		billing_profile_snapshot: {},
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
		invoice_date: "2026-06-11",
		invoice_serie: "SYN",
		invoice_number: 1,
		pdf_path: "/data/invoices/invoice-1.pdf",
		created_at: "2026-06-11T00:00:00Z",
		updated_at: "2026-06-11T00:00:00Z"
	};
}

function adminJsonRecipeResponse() {
	const category = adminJsonRecipeCategoryResponse();
	return {
		id: "recipe-1",
		title: "Invoice",
		description: "Invoice fields",
		json: { type: "object", properties: { number: { type: "string" } } },
		counter: 2,
		category_id: category.id,
		category,
		created_at: "2026-06-20T00:00:00Z",
		updated_at: "2026-06-21T00:00:00Z"
	};
}

function adminJsonRecipeCategoryResponse() {
	return {
		id: "category-1",
		title: { en: "Invoices", ro: "Facturi" },
		created_at: "2026-06-19T00:00:00Z",
		updated_at: "2026-06-19T00:00:00Z"
	};
}

describe("frontend admin server helper", () => {
	beforeEach(() => {
		vi.resetModules();
		for (const key of Object.keys(privateEnv)) delete privateEnv[key];
		process.env.SYNCRA_API_BASE_URL = "http://admin-api.test/";
		process.env.SYNCRA_INTERNAL_API_TOKEN = "internal-token";
		process.env.NODE_ENV = "test";
	});

	it("lists admin users through the backend with internal token and cookie", async () => {
		const { listAdminUsers } = await import("./admin");
		const response = { users: [adminUserResponse()], next_cursor: "cursor-1" };
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(response), { status: 200 });
		});

		await expect(
			listAdminUsers(fetchMock, "auth.session_token=session-1", {
				search: "ada",
				sort: "last_login_at",
				direction: "desc",
				cursor: "cursor-0",
				size: "50"
			})
		).resolves.toEqual(response);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://admin-api.test/api/admin/users?search=ada&sort=last_login_at&direction=desc&cursor=cursor-0&size=50",
			expect.objectContaining({ method: "GET" })
		);
		const headers = new Headers(fetchMock.mock.calls[0][1]?.headers);
		expect(headers.get("X-Syncra-Internal-Token")).toBe("internal-token");
		expect(headers.get("cookie")).toBe("auth.session_token=session-1");
	});

	it("lists admin billing orders through the backend with internal token and cookie", async () => {
		const { listAdminBillingOrders } = await import("./admin");
		const response = adminBillingOrdersResponse();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(response), { status: 200 });
		});

		await expect(
			listAdminBillingOrders(fetchMock, "auth.session_token=session-1", {
				userId: "user-1",
				status: "paid",
				withoutInvoice: true,
				createdFrom: "2026-06-04T00:00:00Z",
				createdTo: "2026-06-05T00:00:00Z",
				cursor: "cursor-0",
				size: "50",
				sort: "desc"
			})
		).resolves.toEqual(response);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://admin-api.test/api/admin/billing/orders?user_id=user-1&status=paid&without_invoice=true&created_from=2026-06-04T00%3A00%3A00Z&created_to=2026-06-05T00%3A00%3A00Z&cursor=cursor-0&size=50&sort=desc",
			expect.objectContaining({ method: "GET" })
		);
		const headers = new Headers(fetchMock.mock.calls[0][1]?.headers);
		expect(headers.get("X-Syncra-Internal-Token")).toBe("internal-token");
		expect(headers.get("cookie")).toBe("auth.session_token=session-1");
	});

	it("generates admin billing order invoices through the backend", async () => {
		const { generateAdminBillingOrderInvoice } = await import("./admin");
		const response = adminBillingInvoiceResponse();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(response), { status: 201 });
		});

		await expect(
			generateAdminBillingOrderInvoice(fetchMock, "auth.session_token=session-1", "order-1")
		).resolves.toEqual(response);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://admin-api.test/api/admin/billing/orders/order-1/invoice",
			expect.objectContaining({ method: "POST" })
		);
		const headers = new Headers(fetchMock.mock.calls[0][1]?.headers);
		expect(headers.get("X-Syncra-Internal-Token")).toBe("internal-token");
		expect(headers.get("cookie")).toBe("auth.session_token=session-1");
	});

	it("lists admin billing invoices through the backend with internal token and cookie", async () => {
		const { listAdminBillingInvoices } = await import("./admin");
		const response = { invoices: [adminBillingInvoiceResponse()], next_cursor: "cursor-1" };
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(response), { status: 200 });
		});

		await expect(
			listAdminBillingInvoices(fetchMock, "auth.session_token=session-1", {
				search: "SYN-00042",
				userId: "user-1",
				createdFrom: "2026-06-11T00:00:00Z",
				createdTo: "2026-06-12T00:00:00Z",
				cursor: "cursor-0",
				size: "50",
				sort: "desc"
			})
		).resolves.toEqual(response);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://admin-api.test/api/admin/billing/invoices?search=SYN-00042&user_id=user-1&created_from=2026-06-11T00%3A00%3A00Z&created_to=2026-06-12T00%3A00%3A00Z&cursor=cursor-0&size=50&sort=desc",
			expect.objectContaining({ method: "GET" })
		);
		const headers = new Headers(fetchMock.mock.calls[0][1]?.headers);
		expect(headers.get("X-Syncra-Internal-Token")).toBe("internal-token");
		expect(headers.get("cookie")).toBe("auth.session_token=session-1");
	});

	it("manages admin JSON recipes through the backend", async () => {
		const {
			createAdminJSONRecipeCategory,
			createAdminJSONRecipe,
			deleteAdminJSONRecipeCategory,
			deleteAdminJSONRecipe,
			getAdminJSONRecipeCategory,
			getAdminJSONRecipe,
			listAdminJSONRecipeCategories,
			listAdminJSONRecipes,
			updateAdminJSONRecipeCategory,
			updateAdminJSONRecipe
		} = await import("./admin");
		const recipe = adminJsonRecipeResponse();
		const category = adminJsonRecipeCategoryResponse();
		const list = { recipes: [recipe], next_cursor: "cursor-1" };
		const categories = { categories: [category] };
		const fetchMock = vi.fn(async (input: FetchInput, init?: FetchInit) => {
			const url = String(input);
			if (url.includes("/json-recipe-categories")) {
				if (init?.method === "GET" && url.endsWith("/json-recipe-categories")) {
					return new Response(JSON.stringify(categories));
				}
				if (init?.method === "DELETE") return new Response(null, { status: 204 });
				return new Response(JSON.stringify(category), { status: init?.method === "POST" ? 201 : 200 });
			}
			if (init?.method === "GET" && url.includes("?")) return new Response(JSON.stringify(list));
			if (init?.method === "DELETE") return new Response(null, { status: 204 });
			return new Response(JSON.stringify(recipe), { status: init?.method === "POST" ? 201 : 200 });
		});
		const input = {
			title: "Invoice",
			description: "Invoice fields",
			json: { type: "object" },
			category_id: "category-1"
		};
		const categoryInput = { title: { en: "Invoices", ro: "Facturi" } };

		await expect(
			listAdminJSONRecipes(fetchMock, "auth.session_token=session-1", {
				cursor: "cursor-0",
				size: 50,
				sort: "asc"
			})
		).resolves.toEqual(list);
		await expect(listAdminJSONRecipeCategories(fetchMock, "auth.session_token=session-1")).resolves.toEqual(
			categories
		);
		await expect(getAdminJSONRecipe(fetchMock, "auth.session_token=session-1", "recipe-1")).resolves.toEqual(
			recipe
		);
		await expect(
			getAdminJSONRecipeCategory(fetchMock, "auth.session_token=session-1", "category-1")
		).resolves.toEqual(category);
		await expect(createAdminJSONRecipe(fetchMock, "auth.session_token=session-1", input)).resolves.toEqual(
			recipe
		);
		await expect(
			createAdminJSONRecipeCategory(fetchMock, "auth.session_token=session-1", categoryInput)
		).resolves.toEqual(category);
		await expect(
			updateAdminJSONRecipe(fetchMock, "auth.session_token=session-1", "recipe-1", input)
		).resolves.toEqual(recipe);
		await expect(
			updateAdminJSONRecipeCategory(fetchMock, "auth.session_token=session-1", "category-1", categoryInput)
		).resolves.toEqual(category);
		await expect(deleteAdminJSONRecipe(fetchMock, "auth.session_token=session-1", "recipe-1")).resolves.toBeUndefined();
		await expect(
			deleteAdminJSONRecipeCategory(fetchMock, "auth.session_token=session-1", "category-1")
		).resolves.toBeUndefined();

		expect(fetchMock).toHaveBeenNthCalledWith(
			1,
			"http://admin-api.test/api/admin/json-recipes?cursor=cursor-0&size=50&sort=asc",
			expect.objectContaining({ method: "GET" })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			2,
			"http://admin-api.test/api/admin/json-recipe-categories",
			expect.objectContaining({ method: "GET" })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			3,
			"http://admin-api.test/api/admin/json-recipes/recipe-1",
			expect.objectContaining({ method: "GET" })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			4,
			"http://admin-api.test/api/admin/json-recipe-categories/category-1",
			expect.objectContaining({ method: "GET" })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			5,
			"http://admin-api.test/api/admin/json-recipes",
			expect.objectContaining({ method: "POST", body: JSON.stringify(input) })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			6,
			"http://admin-api.test/api/admin/json-recipe-categories",
			expect.objectContaining({ method: "POST", body: JSON.stringify(categoryInput) })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			7,
			"http://admin-api.test/api/admin/json-recipes/recipe-1",
			expect.objectContaining({ method: "PUT", body: JSON.stringify(input) })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			8,
			"http://admin-api.test/api/admin/json-recipe-categories/category-1",
			expect.objectContaining({ method: "PUT", body: JSON.stringify(categoryInput) })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			9,
			"http://admin-api.test/api/admin/json-recipes/recipe-1",
			expect.objectContaining({ method: "DELETE" })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			10,
			"http://admin-api.test/api/admin/json-recipe-categories/category-1",
			expect.objectContaining({ method: "DELETE" })
		);
		for (const call of fetchMock.mock.calls) {
			const headers = new Headers(call[1]?.headers);
			expect(headers.get("X-Syncra-Internal-Token")).toBe("internal-token");
			expect(headers.get("cookie")).toBe("auth.session_token=session-1");
		}
	});

	it("rejects invalid admin JSON recipe responses", async () => {
		const { listAdminJSONRecipes } = await import("./admin");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ recipes: [{ ...adminJsonRecipeResponse(), counter: "2" }], next_cursor: null }));
		});

		await expect(listAdminJSONRecipes(fetchMock, "auth.session_token=session-1", {})).rejects.toMatchObject({
			status: 502,
			message: "Invalid admin JSON recipe response"
		});
	});

	it("generates admin billing invoice PDFs through the backend", async () => {
		const { generateAdminBillingInvoicePDF } = await import("./admin");
		const response = adminBillingInvoiceResponse();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(response), { status: 200 });
		});

		await expect(
			generateAdminBillingInvoicePDF(fetchMock, "auth.session_token=session-1", "invoice-1")
		).resolves.toEqual(response);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://admin-api.test/api/billing/generate-invoice-pdf",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({ invoice_id: "invoice-1" })
			})
		);
		const headers = new Headers(fetchMock.mock.calls[0][1]?.headers);
		expect(headers.get("X-Syncra-Internal-Token")).toBe("internal-token");
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get("cookie")).toBe("auth.session_token=session-1");
	});

	it("fetches admin billing invoice PDFs through the backend", async () => {
		const { fetchAdminBillingInvoicePDF } = await import("./admin");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response("%PDF-test", {
				status: 200,
				headers: {
					"content-type": "application/pdf",
					"content-disposition": 'inline; filename="invoice-1.pdf"'
				}
			});
		});

		const result = await fetchAdminBillingInvoicePDF(
			fetchMock,
			"auth.session_token=session-1",
			"invoice-1",
			{ download: true }
		);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://admin-api.test/static/invoice/invoice-1.pdf?download=1",
			expect.objectContaining({ method: "GET" })
		);
		const headers = new Headers(fetchMock.mock.calls[0][1]?.headers);
		expect(headers.get("X-Syncra-Internal-Token")).toBe("internal-token");
		expect(headers.get("cookie")).toBe("auth.session_token=session-1");
		expect(result.status).toBe(200);
		expect(result.headers.get("content-type")).toBe("application/pdf");
		expect(result.headers.get("content-disposition")).toBe('inline; filename="invoice-1.pdf"');
		await expect(new Response(result.body).text()).resolves.toBe("%PDF-test");
	});

	it("rejects invalid admin billing invoice responses", async () => {
		const { generateAdminBillingOrderInvoice } = await import("./admin");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ ...adminBillingInvoiceResponse(), lines: [{}] }), {
				status: 201
			});
		});

		await expect(
			generateAdminBillingOrderInvoice(fetchMock, "auth.session_token=session-1", "order-1")
		).rejects.toMatchObject({
			status: 502,
			message: "Invalid admin billing invoice response"
		});
	});

	it("rejects invalid admin billing invoice pdf paths", async () => {
		const { listAdminBillingInvoices } = await import("./admin");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(
				JSON.stringify({
					invoices: [{ ...adminBillingInvoiceResponse(), pdf_path: 42 }],
					next_cursor: null
				}),
				{ status: 200 }
			);
		});

		await expect(listAdminBillingInvoices(fetchMock, "auth.session_token=session-1", {})).rejects.toMatchObject({
			status: 502,
			message: "Invalid admin billing invoices response"
		});
	});

	it("rejects invalid admin billing invoice list responses", async () => {
		const { listAdminBillingInvoices } = await import("./admin");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ invoices: [{ ...adminBillingInvoiceResponse(), lines: [{}] }], next_cursor: null }), {
				status: 200
			});
		});

		await expect(listAdminBillingInvoices(fetchMock, "auth.session_token=session-1", {})).rejects.toMatchObject({
			status: 502,
			message: "Invalid admin billing invoices response"
		});
	});

	it("rejects invalid admin billing order responses", async () => {
		const { listAdminBillingOrders } = await import("./admin");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ orders: [{ status: "settled" }], next_cursor: null }), {
				status: 200
			});
		});

		await expect(listAdminBillingOrders(fetchMock, "auth.session_token=session-1", {})).rejects.toMatchObject({
			status: 502,
			message: "Invalid admin billing orders response"
		});
	});

	it("loads and adjusts an admin user balance through the backend", async () => {
		const { adjustAdminUserBalance, getAdminUser } = await import("./admin");
		const detail = adminUserDetailResponse();
		const adjusted = { ...detail, available_credits: 75 };
		const fetchMock = vi.fn(async (input: FetchInput, _init?: FetchInit) => {
			const url = String(input);
			if (url.endsWith("/balance-adjustment")) {
				return new Response(JSON.stringify(adjusted), { status: 200 });
			}
			return new Response(JSON.stringify(detail), { status: 200 });
		});

		await expect(getAdminUser(fetchMock, "auth.session_token=session-1", "user-1")).resolves.toEqual(detail);
		await expect(adjustAdminUserBalance(fetchMock, "auth.session_token=session-1", "user-1", -50)).resolves.toEqual(
			adjusted
		);

		expect(fetchMock).toHaveBeenNthCalledWith(
			1,
			"http://admin-api.test/api/admin/users/user-1",
			expect.objectContaining({ method: "GET" })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			2,
			"http://admin-api.test/api/admin/users/user-1/balance-adjustment",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({ credits_delta: -50 })
			})
		);
	});

	it("updates users without accepting role fields", async () => {
		const { updateAdminUser } = await import("./admin");
		const updated = { ...adminUserResponse(), name: "Ada Updated" };
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(updated), { status: 200 });
		});

		await expect(
			updateAdminUser(fetchMock, "auth.session_token=session-1", "user-1", {
				name: "Ada Updated",
				email: "ada.updated@example.com"
			})
		).resolves.toEqual(updated);

		expect(fetchMock).toHaveBeenCalledWith(
			"http://admin-api.test/api/admin/users/user-1",
			expect.objectContaining({
				method: "PATCH",
				body: JSON.stringify({ name: "Ada Updated", email: "ada.updated@example.com" })
			})
		);
	});

	it("starts and stops admin impersonation through the backend", async () => {
		const { startAdminUserImpersonation, stopAdminImpersonation } = await import("./admin");
		const started = authSessionPayloadResponse();
		const stopped = authSessionPayloadResponse({
			user: authUserResponse("admin-1", "admin@example.com", "admin"),
			impersonation: null
		});
		const fetchMock = vi.fn(async (input: FetchInput, _init?: FetchInit) => {
			const url = String(input);
			if (url.endsWith("/api/admin/users/user-1/impersonation")) {
				return new Response(JSON.stringify(started), { status: 200 });
			}
			return new Response(JSON.stringify(stopped), { status: 200 });
		});

		await expect(
			startAdminUserImpersonation(fetchMock, "auth.session_token=session-1", "user-1")
		).resolves.toEqual(started);
		await expect(stopAdminImpersonation(fetchMock, "auth.session_token=session-1")).resolves.toEqual(stopped);

		expect(fetchMock).toHaveBeenNthCalledWith(
			1,
			"http://admin-api.test/api/admin/users/user-1/impersonation",
			expect.objectContaining({ method: "POST" })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			2,
			"http://admin-api.test/api/admin/impersonation/stop",
			expect.objectContaining({ method: "POST" })
		);
		const startHeaders = new Headers(fetchMock.mock.calls[0][1]?.headers);
		const stopHeaders = new Headers(fetchMock.mock.calls[1][1]?.headers);
		expect(startHeaders.get("X-Syncra-Internal-Token")).toBe("internal-token");
		expect(startHeaders.get("cookie")).toBe("auth.session_token=session-1");
		expect(stopHeaders.get("X-Syncra-Internal-Token")).toBe("internal-token");
		expect(stopHeaders.get("cookie")).toBe("auth.session_token=session-1");
	});

	it("resets user passwords and updates target billing profiles", async () => {
		const { resetAdminUserPassword, upsertAdminUserBillingProfile } = await import("./admin");
		const fetchMock = vi.fn(async (input: FetchInput, _init?: FetchInit) => {
			const url = String(input);
			if (url.endsWith("/password")) return new Response(JSON.stringify({ ok: true }), { status: 200 });
			return new Response(JSON.stringify(billingProfileResponse()), { status: 200 });
		});

		await expect(
			resetAdminUserPassword(fetchMock, "auth.session_token=session-1", "user-1", "newpassword123")
		).resolves.toEqual({ ok: true });
		await expect(
			upsertAdminUserBillingProfile(fetchMock, "auth.session_token=session-1", "user-1", {
				entity_type: "company",
				billing_name: "Syncra SRL",
				billing_email: "billing@example.com",
				country_code: "RO",
				address_line1: "Main Street 1",
				city: "Bucharest",
				postal_code: "010101",
				fiscal_code: "RO123"
			})
		).resolves.toEqual(billingProfileResponse());

		expect(fetchMock).toHaveBeenNthCalledWith(
			1,
			"http://admin-api.test/api/admin/users/user-1/password",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({ password: "newpassword123" })
			})
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			2,
			"http://admin-api.test/api/admin/users/user-1/billing-profile",
			expect.objectContaining({
				method: "PUT",
				body: JSON.stringify({
					entity_type: "company",
					billing_name: "Syncra SRL",
					billing_email: "billing@example.com",
					country_code: "RO",
					address_line1: "Main Street 1",
					city: "Bucharest",
					postal_code: "010101",
					fiscal_code: "RO123"
				})
			})
		);
	});
});
