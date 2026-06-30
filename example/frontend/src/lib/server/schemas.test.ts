import { beforeEach, describe, expect, it, vi } from "vitest";

type FetchInput = Parameters<typeof fetch>[0];
type FetchInit = Parameters<typeof fetch>[1];
type FetchMock = (input: FetchInput, init?: FetchInit) => Promise<Response>;

const INTERNAL_API_HEADER = "X-Syncra-Internal-Token";

const validSchema = (id = "schema-1", userId?: string) => ({
	id,
	created_at: "2026-05-27T00:00:00Z",
	updated_at: "2026-05-27T00:00:00Z",
	...(userId ? { user_id: userId } : {}),
	name: "Invoice",
	description: "Invoice fields",
	strict: true,
	schema: { type: "object" }
});

const validDeleteResponse = {
	deleted_ids: ["schema-1"],
	deleted_count: 1
};

const validRecipe = {
	id: "recipe-1",
	title: "Invoice",
	description: "Invoice fields",
	json: { type: "object" },
	counter: 1,
	category_id: "category-1",
	category: {
		id: "category-1",
		title: { en: "Invoices", ro: "Facturi" },
		created_at: "2026-06-19T00:00:00Z",
		updated_at: "2026-06-19T00:00:00Z"
	},
	created_at: "2026-06-20T00:00:00Z",
	updated_at: "2026-06-21T00:00:00Z"
};

const { privateEnv } = vi.hoisted(() => ({
	privateEnv: {} as Record<string, string | undefined>
}));

vi.mock("$env/dynamic/private", () => ({ env: privateEnv }));

function requestHeaders(init: FetchInit | undefined) {
	return new Headers(init?.headers);
}

describe("frontend schema server helper", () => {
	beforeEach(() => {
		vi.resetModules();
		for (const key of Object.keys(privateEnv)) delete privateEnv[key];
		process.env.SYNCRA_API_BASE_URL = "http://schema-api.test/";
		process.env.SYNCRA_INTERNAL_API_TOKEN = "internal-token";
		process.env.NODE_ENV = "test";
	});

	it("creates schemas through the backend with JSON", async () => {
		const { createSchema } = await import("./schemas");
		const savedSchema = {
			id: "schema-1",
			created_at: "2026-05-26T00:00:00Z",
			updated_at: "2026-05-26T00:00:00Z",
			name: "Customer",
			description: "Customer profile payload",
			strict: true,
			schema: {
				type: "object",
				required: ["email"],
				properties: {
					email: { type: "string" }
				}
			}
		};
		const input = {
			name: "Customer",
			description: "Customer profile payload",
			strict: true,
			schema: savedSchema.schema
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(savedSchema));
		});

		const result = await createSchema(fetchMock, input);

		expect(result).toEqual(savedSchema);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://schema-api.test/api/ocr/schemas",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify(input)
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("creates user-owned schemas when a user id is provided", async () => {
		const { createSchema } = await import("./schemas");
		const savedSchema = {
			id: "schema-1",
			created_at: "2026-05-26T00:00:00Z",
			updated_at: "2026-05-26T00:00:00Z",
			user_id: "user-1",
			name: "Customer",
			description: "Customer profile payload",
			strict: true,
			schema: { type: "object" }
		};
		const input = {
			name: "Customer",
			description: "Customer profile payload",
			strict: true,
			schema: savedSchema.schema
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(savedSchema));
		});

		const result = await createSchema(fetchMock, input, { userId: "user-1" });

		expect(result).toEqual(savedSchema);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://schema-api.test/api/ocr/schemas",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({ ...input, user_id: "user-1" })
			})
		);
	});

	it("prefers Svelte private env over process env for the backend URL", async () => {
		privateEnv.SYNCRA_API_BASE_URL = "http://private-schema-api.test/";
		process.env.SYNCRA_API_BASE_URL = "http://process-schema-api.test/";
		const { createSchema } = await import("./schemas");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(
				JSON.stringify({
					id: "schema-1",
					created_at: "2026-05-26T00:00:00Z",
					updated_at: "2026-05-26T00:00:00Z",
					name: "Customer",
					description: "Customer profile payload",
					strict: true,
					schema: { type: "object" }
				})
			);
		});

		await createSchema(fetchMock, {
			name: "Customer",
			description: "Customer profile payload",
			strict: true,
			schema: { type: "object" }
		});

		expect(fetchMock).toHaveBeenCalledWith(
			"http://private-schema-api.test/api/ocr/schemas",
			expect.any(Object)
		);
	});

	it("throws typed schema API errors from backend error responses", async () => {
		const { createSchema, isSchemaApiError } = await import("./schemas");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ error: "schema name already exists" }), {
				status: 409
			});
		});
		const input = {
			name: "Customer",
			description: "Customer profile payload",
			strict: true,
			schema: { type: "object" }
		};

		await expect(createSchema(fetchMock, input)).rejects.toMatchObject({
			status: 409,
			message: "schema name already exists"
		});

		try {
			await createSchema(fetchMock, input);
		} catch (error) {
			expect(isSchemaApiError(error)).toBe(true);
		}
	});

	it("converts fetch outages into typed schema API errors", async () => {
		const { createSchema, isSchemaApiError } = await import("./schemas");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit): Promise<Response> => {
			throw new Error("ECONNREFUSED");
		});
		const input = {
			name: "Customer",
			description: "Customer profile payload",
			strict: true,
			schema: { type: "object" }
		};

		await expect(createSchema(fetchMock, input)).rejects.toMatchObject({
			status: 503,
			message: "Schema service unavailable"
		});

		try {
			await createSchema(fetchMock, input);
		} catch (error) {
			expect(isSchemaApiError(error)).toBe(true);
		}
	});

	it("rejects invalid successful schema responses", async () => {
		const { createSchema, isSchemaApiError } = await import("./schemas");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({}), { status: 200 });
		});
		const input = {
			name: "Customer",
			description: "Customer profile payload",
			strict: true,
			schema: { type: "object" }
		};

		await expect(createSchema(fetchMock, input)).rejects.toMatchObject({
			status: 502,
			message: "Invalid schema response"
		});

		try {
			await createSchema(fetchMock, input);
		} catch (error) {
			expect(isSchemaApiError(error)).toBe(true);
		}
	});

	it("converts body read failures into typed schema API errors", async () => {
		const { createSchema, isSchemaApiError } = await import("./schemas");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit): Promise<Response> => {
			return {
				ok: true,
				text: vi.fn(async () => {
					throw new Error("stream failure");
				})
			} as unknown as Response;
		});
		const input = {
			name: "Customer",
			description: "Customer profile payload",
			strict: true,
			schema: { type: "object" }
		};

		await expect(createSchema(fetchMock, input)).rejects.toMatchObject({
			status: 503,
			message: "Schema service unavailable"
		});

		try {
			await createSchema(fetchMock, input);
		} catch (error) {
			expect(isSchemaApiError(error)).toBe(true);
		}
	});

	it("lists system schemas through the backend", async () => {
		const { listSchemas } = await import("./schemas");
		const schemas = [validSchema()];
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ schemas, next_cursor: null }));
		});

		const result = await listSchemas(fetchMock);

		expect(result).toEqual(schemas);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://schema-api.test/api/ocr/schemas?size=100",
			expect.objectContaining({ method: "GET" })
		);
		expect(requestHeaders(fetchMock.mock.calls[0][1]).get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("lists user schemas through the backend with a user_id query", async () => {
		const { listSchemas } = await import("./schemas");
		const schemas = [validSchema("schema-1", "user-1")];
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ schemas, next_cursor: null }));
		});

		const result = await listSchemas(fetchMock, { userId: "user-1" });

		expect(result).toEqual(schemas);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://schema-api.test/api/ocr/schemas?user_id=user-1&size=100",
			expect.objectContaining({ method: "GET" })
		);
	});

	it("returns one schema list page with pagination query parameters", async () => {
		const { listSchemasPage } = await import("./schemas");
		const page = {
			schemas: [validSchema("schema-2", "user-1")],
			next_cursor: "cursor-2"
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(page));
		});

		const result = await listSchemasPage(fetchMock, {
			userId: "user-1",
			cursor: "cursor-1",
			size: 25,
			sort: "desc"
		});

		expect(result).toEqual(page);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://schema-api.test/api/ocr/schemas?user_id=user-1&cursor=cursor-1&size=25&sort=desc",
			expect.objectContaining({ method: "GET" })
		);
	});

	it("lists all schema pages for backward compatibility", async () => {
		const { listSchemas } = await import("./schemas");
		const pageOne = {
			schemas: [validSchema("schema-1", "user-1")],
			next_cursor: "cursor-2"
		};
		const pageTwo = {
			schemas: [validSchema("schema-2", "user-1")],
			next_cursor: null
		};
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(new Response(JSON.stringify(pageOne)))
			.mockResolvedValueOnce(new Response(JSON.stringify(pageTwo)));

		const result = await listSchemas(fetchMock, { userId: "user-1" });

		expect(result).toEqual([...pageOne.schemas, ...pageTwo.schemas]);
		expect(fetchMock).toHaveBeenNthCalledWith(
			1,
			"http://schema-api.test/api/ocr/schemas?user_id=user-1&size=100",
			expect.objectContaining({ method: "GET" })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			2,
			"http://schema-api.test/api/ocr/schemas?user_id=user-1&cursor=cursor-2&size=100",
			expect.objectContaining({ method: "GET" })
		);
	});

	it("rejects schema list pages with an empty next cursor", async () => {
		const { isSchemaApiError, listSchemas } = await import("./schemas");
		const pageOne = {
			schemas: [validSchema("schema-1", "user-1")],
			next_cursor: ""
		};
		const pageTwo = {
			schemas: [validSchema("schema-2", "user-1")],
			next_cursor: null
		};
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(new Response(JSON.stringify(pageOne)))
			.mockResolvedValueOnce(new Response(JSON.stringify(pageTwo)));

		let caught: unknown;
		try {
			await listSchemas(fetchMock, { userId: "user-1" });
		} catch (error) {
			caught = error;
		}

		expect(caught).toMatchObject({
			status: 502,
			message: "Invalid schema response"
		});
		expect(isSchemaApiError(caught)).toBe(true);
		expect(fetchMock).toHaveBeenCalledTimes(1);
	});

	it("rejects schema list pages with a repeated next cursor", async () => {
		const { isSchemaApiError, listSchemas } = await import("./schemas");
		const pageOne = {
			schemas: [validSchema("schema-1", "user-1")],
			next_cursor: "cursor-2"
		};
		const pageTwo = {
			schemas: [validSchema("schema-2", "user-1")],
			next_cursor: "cursor-2"
		};
		const pageThree = {
			schemas: [validSchema("schema-3", "user-1")],
			next_cursor: null
		};
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(new Response(JSON.stringify(pageOne)))
			.mockResolvedValueOnce(new Response(JSON.stringify(pageTwo)))
			.mockResolvedValueOnce(new Response(JSON.stringify(pageThree)));

		let caught: unknown;
		try {
			await listSchemas(fetchMock, { userId: "user-1" });
		} catch (error) {
			caught = error;
		}

		expect(caught).toMatchObject({
			status: 502,
			message: "Invalid schema pagination response"
		});
		expect(isSchemaApiError(caught)).toBe(true);
		expect(fetchMock).toHaveBeenCalledTimes(2);
	});

	it("gets a schema by id through the backend", async () => {
		const { getSchema } = await import("./schemas");
		const schema = validSchema("schema-1", "user-1");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(schema));
		});

		const result = await getSchema(fetchMock, "schema-1", { userId: "user-1" });

		expect(result).toEqual(schema);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://schema-api.test/api/ocr/schemas/schema-1?user_id=user-1",
			expect.objectContaining({ method: "GET" })
		);
	});

	it("updates a schema through the backend with JSON", async () => {
		const { updateSchema } = await import("./schemas");
		const input = {
			name: "Receipt",
			description: "Receipt fields",
			strict: false,
			schema: { type: "object" }
		};
		const savedSchema = {
			...validSchema("schema-1", "user-1"),
			...input
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(savedSchema));
		});

		const result = await updateSchema(fetchMock, "schema-1", input, { userId: "user-1" });

		expect(result).toEqual(savedSchema);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://schema-api.test/api/ocr/schemas/schema-1?user_id=user-1",
			expect.objectContaining({
				method: "PUT",
				body: JSON.stringify(input)
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("deletes a schema by id through the backend", async () => {
		const { deleteSchema } = await import("./schemas");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(validDeleteResponse));
		});

		const result = await deleteSchema(fetchMock, "schema-1", { userId: "user-1" });

		expect(result).toEqual(validDeleteResponse);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://schema-api.test/api/ocr/schemas/schema-1?user_id=user-1",
			expect.objectContaining({ method: "DELETE" })
		);
	});

	it("deploys JSON recipes into user-owned schemas", async () => {
		const { deployJsonRecipe } = await import("./schemas");
		const result = {
			recipe: validRecipe,
			schema: validSchema("schema-1", "user-1")
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(result), { status: 201 });
		});

		await expect(deployJsonRecipe(fetchMock, "recipe-1", { userId: "user-1" })).resolves.toEqual(result);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://schema-api.test/api/json-recipes/recipe-1/deploy",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({ user_id: "user-1" })
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("rejects invalid JSON recipe deploy responses", async () => {
		const { deployJsonRecipe } = await import("./schemas");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ recipe: { ...validRecipe, json: null }, schema: validSchema() }), {
				status: 201
			});
		});

		await expect(deployJsonRecipe(fetchMock, "recipe-1", { userId: "user-1" })).rejects.toMatchObject({
			status: 502,
			message: "Invalid JSON recipe deploy response"
		});
	});

	it("bulk deletes schemas through the backend with JSON", async () => {
		const { deleteSchemas } = await import("./schemas");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(validDeleteResponse));
		});

		const result = await deleteSchemas(fetchMock, ["schema-1"], { userId: "user-1" });

		expect(result).toEqual(validDeleteResponse);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://schema-api.test/api/ocr/schemas?user_id=user-1",
			expect.objectContaining({
				method: "DELETE",
				body: JSON.stringify({ ids: ["schema-1"] })
			})
		);
		const headers = fetchMock.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("rejects schema requests before calling Go when the internal API token is missing", async () => {
		const { listSchemasPage } = await import("./schemas");
		delete process.env.SYNCRA_INTERNAL_API_TOKEN;
		const fetchMock = vi.fn();

		await expect(listSchemasPage(fetchMock)).rejects.toMatchObject({
			status: 500,
			message: "Schema service is not configured"
		});
		expect(fetchMock).not.toHaveBeenCalled();
	});

	it("rejects invalid successful schema list responses", async () => {
		const { isSchemaApiError, listSchemas } = await import("./schemas");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ schemas: [{ id: "schema-1" }], next_cursor: null }), {
				status: 200
			});
		});

		await expect(listSchemas(fetchMock)).rejects.toMatchObject({
			status: 502,
			message: "Invalid schema response"
		});

		try {
			await listSchemas(fetchMock);
		} catch (error) {
			expect(isSchemaApiError(error)).toBe(true);
		}
	});

	it.each([
		{
			name: "listSchemasPage",
			call: async (fetchMock: FetchMock) => {
				const { listSchemasPage } = await import("./schemas");
				return listSchemasPage(fetchMock as unknown as typeof fetch);
			},
			invalidPayload: { schemas: [validSchema()], next_cursor: 12 }
		},
		{
			name: "getSchema",
			call: async (fetchMock: FetchMock) => {
				const { getSchema } = await import("./schemas");
				return getSchema(fetchMock as unknown as typeof fetch, "schema-1", { userId: "user-1" });
			},
			invalidPayload: { id: "schema-1" }
		},
		{
			name: "updateSchema",
			call: async (fetchMock: FetchMock) => {
				const { updateSchema } = await import("./schemas");
				return updateSchema(
					fetchMock as unknown as typeof fetch,
					"schema-1",
					{
						name: "Receipt",
						description: "Receipt fields",
						strict: false,
						schema: { type: "object" }
					},
					{ userId: "user-1" }
				);
			},
			invalidPayload: { id: "schema-1" }
		},
		{
			name: "deleteSchema",
			call: async (fetchMock: FetchMock) => {
				const { deleteSchema } = await import("./schemas");
				return deleteSchema(fetchMock as unknown as typeof fetch, "schema-1", { userId: "user-1" });
			},
			invalidPayload: { deleted_ids: ["schema-1"], deleted_count: "1" }
		},
		{
			name: "deleteSchemas",
			call: async (fetchMock: FetchMock) => {
				const { deleteSchemas } = await import("./schemas");
				return deleteSchemas(fetchMock as unknown as typeof fetch, ["schema-1"], {
					userId: "user-1"
				});
			},
			invalidPayload: { deleted_ids: ["schema-1"], deleted_count: "1" }
		}
	])(
		"converts fetch outages, backend errors, invalid payloads, and body read failures for $name into typed errors",
		async ({ call, invalidPayload }) => {
			const { isSchemaApiError } = await import("./schemas");
			const outageFetch = vi.fn(async (_input: FetchInput, _init?: FetchInit): Promise<Response> => {
				throw new Error("ECONNREFUSED");
			});
			const backendErrorFetch = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
				return new Response(JSON.stringify({ error: "schema conflict" }), { status: 409 });
			});
			const invalidPayloadFetch = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
				return new Response(JSON.stringify(invalidPayload), { status: 200 });
			});
			const bodyFailureFetch = vi.fn(async (_input: FetchInput, _init?: FetchInit): Promise<Response> => {
				return {
					ok: true,
					text: vi.fn(async () => {
						throw new Error("stream failure");
					})
				} as unknown as Response;
			});

			await expect(call(outageFetch)).rejects.toMatchObject({
				status: 503,
				message: "Schema service unavailable"
			});
			await expect(call(backendErrorFetch)).rejects.toMatchObject({
				status: 409,
				message: "schema conflict"
			});
			await expect(call(invalidPayloadFetch)).rejects.toMatchObject({
				status: 502,
				message: "Invalid schema response"
			});
			await expect(call(bodyFailureFetch)).rejects.toMatchObject({
				status: 503,
				message: "Schema service unavailable"
			});

			for (const fetchMock of [outageFetch, backendErrorFetch, invalidPayloadFetch, bodyFailureFetch]) {
				try {
					await call(fetchMock);
				} catch (error) {
					expect(isSchemaApiError(error)).toBe(true);
				}
			}
		}
	);
});
