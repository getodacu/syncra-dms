import { describe, expect, it, vi } from "vitest";

const schemaResponse = (patch: Partial<Record<string, unknown>> = {}) => ({
	id: "schema-1",
	created_at: "2026-05-01T00:00:00.000Z",
	updated_at: "2026-05-01T00:00:00.000Z",
	user_id: "user-1",
	name: "Invoices",
	description: "Invoice extraction",
	strict: true,
	schema: { type: "object" },
	...patch,
});

const jsonResponse = (body: unknown, init: ResponseInit = {}) =>
	new Response(JSON.stringify(body), {
		status: 200,
		headers: { "content-type": "application/json" },
		...init,
	});

describe("schema option client helpers", () => {
	it("uses a dedicated personal schema options query key", async () => {
		const { PERSONAL_SCHEMA_OPTIONS_QUERY_KEY } = await import("./schemas");

		expect(PERSONAL_SCHEMA_OPTIONS_QUERY_KEY).toEqual([
			"schemas",
			"mine",
			"options",
			{ size: 100 },
		]);
	});

	it("fetches personal schema options from the paginated schema proxy", async () => {
		const { fetchPersonalSchemaOptions } = await import("./schemas");
		const fetchFn = vi.fn().mockResolvedValue(
			jsonResponse({
				schemas: [schemaResponse(), schemaResponse({ id: "schema-2", name: "Receipts" })],
				next_cursor: null,
			})
		);

		await expect(fetchPersonalSchemaOptions(fetchFn)).resolves.toEqual([
			{
				id: "schema-1",
				name: "Invoices",
				description: "Invoice extraction",
			},
			{
				id: "schema-2",
				name: "Receipts",
				description: "Invoice extraction",
			},
		]);
		expect(fetchFn).toHaveBeenCalledWith("/api/schemas?scope=mine&size=100");
	});

	it("rejects invalid personal schema option payloads", async () => {
		const { fetchPersonalSchemaOptions } = await import("./schemas");
		const fetchFn = vi.fn().mockResolvedValue(
			jsonResponse({
				schemas: [{ id: "schema-1", name: "Invoices" }],
				next_cursor: null,
			})
		);

		await expect(fetchPersonalSchemaOptions(fetchFn)).rejects.toThrow("Invalid schema response");
	});

	it("upserts saved schemas into the personal option list", async () => {
		const { upsertPersonalSchemaOption } = await import("./schemas");

		expect(upsertPersonalSchemaOption(undefined, schemaResponse())).toEqual([
			{
				id: "schema-1",
				name: "Invoices",
				description: "Invoice extraction",
			},
		]);

		expect(
			upsertPersonalSchemaOption(
				[
					{ id: "schema-1", name: "Old invoices", description: "Old" },
					{ id: "schema-2", name: "Receipts", description: "Receipt extraction" },
				],
				schemaResponse({ name: "Invoices v2" })
			)
		).toEqual([
			{ id: "schema-1", name: "Invoices v2", description: "Invoice extraction" },
			{ id: "schema-2", name: "Receipts", description: "Receipt extraction" },
		]);

		expect(
			upsertPersonalSchemaOption(
				[{ id: "schema-2", name: "Receipts", description: "Receipt extraction" }],
				schemaResponse()
			)
		).toEqual([
			{ id: "schema-1", name: "Invoices", description: "Invoice extraction" },
			{ id: "schema-2", name: "Receipts", description: "Receipt extraction" },
		]);
	});

	it("removes deleted schemas from the personal option list", async () => {
		const { removePersonalSchemaOptions } = await import("./schemas");

		expect(
			removePersonalSchemaOptions(
				[
					{ id: "schema-1", name: "Invoices", description: "Invoice extraction" },
					{ id: "schema-2", name: "Receipts", description: "Receipt extraction" },
					{ id: "schema-3", name: "Orders", description: "Order extraction" },
				],
				["schema-1", "schema-3"]
			)
		).toEqual([{ id: "schema-2", name: "Receipts", description: "Receipt extraction" }]);

		expect(removePersonalSchemaOptions(undefined, ["schema-1"])).toBeUndefined();
	});
});
