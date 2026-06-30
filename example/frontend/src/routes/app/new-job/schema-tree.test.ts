import { describe, expect, it } from "vitest";
import { buildSchemaTree } from "./schema-tree";

describe("new job schema tree", () => {
	it("builds a tree for object properties and required fields", () => {
		const tree = buildSchemaTree({
			type: "object",
			required: ["invoice_number"],
			properties: {
				invoice_number: {
					type: "string",
					description: "Invoice identifier",
					minLength: 3,
					maxLength: 32
				},
				total: {
					type: "number",
					minimum: 0,
					maximum: 100
				}
			}
		});

		expect(tree).toEqual([
			expect.objectContaining({
				name: "invoice_number",
				type: "string",
				required: true,
				description: "Invoice identifier",
				details: ["min length: 3", "max length: 32"]
			}),
			expect.objectContaining({
				name: "total",
				type: "number",
				required: false,
				details: ["min: 0", "max: 100"]
			})
		]);
	});

	it("includes nested object and array item details", () => {
		const tree = buildSchemaTree({
			type: "object",
			properties: {
				customer: {
					type: "object",
					properties: {
						email: { type: "string", format: "email" }
					}
				},
				lines: {
					type: "array",
					items: {
						type: "object",
						properties: {
							description: { type: "string" }
						}
					}
				}
			}
		});

		expect(tree[0].children[0]).toEqual(
			expect.objectContaining({ name: "email", type: "string", details: ["format: email"] })
		);
		expect(tree[1]).toEqual(expect.objectContaining({ name: "lines", type: "array" }));
		expect(tree[1].children[0]).toEqual(expect.objectContaining({ name: "items" }));
		expect(tree[1].children[0].children[0]).toEqual(
			expect.objectContaining({
				name: "description",
				type: "string",
				id: "/properties/lines/items/properties/description",
				path: "/properties/lines/items/properties/description"
			})
		);
	});

	it("returns a root node when a schema has no properties", () => {
		expect(buildSchemaTree({ type: "string", enum: [null, "paid", "open", 1, true, "void", "draft"] })).toEqual([
			expect.objectContaining({
				name: "schema",
				type: "string",
				details: ["enum: null, paid, open, 1, true, void"]
			})
		]);
	});

	it("uses JSON pointer paths so property names cannot collide", () => {
		const tree = buildSchemaTree({
			type: "object",
			properties: {
				"a.b": { type: "string" },
				a: {
					type: "object",
					properties: {
						b: { type: "number" },
						"c/d~e": { type: "boolean" }
					}
				}
			}
		});

		expect(tree[0]).toEqual(
			expect.objectContaining({ name: "a.b", id: "/properties/a.b", path: "/properties/a.b" })
		);
		expect(tree[1].children[0]).toEqual(
			expect.objectContaining({ name: "b", id: "/properties/a/properties/b", path: "/properties/a/properties/b" })
		);
		expect(tree[1].children[1]).toEqual(
			expect.objectContaining({
				name: "c/d~e",
				id: "/properties/a/properties/c~1d~0e",
				path: "/properties/a/properties/c~1d~0e"
			})
		);
		expect(tree[0].id).not.toBe(tree[1].children[0].id);
	});
});
