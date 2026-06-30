import { describe, expect, it } from "vitest";

import { buildDatasetFieldTree } from "./field-tree";

describe("dataset field tree", () => {
	it("builds selectable scalar nodes from object schema properties", () => {
		const tree = buildDatasetFieldTree({
			type: "object",
			properties: {
				invoice_number: {
					type: "string",
					description: "Invoice identifier",
				},
				total: { type: "number" },
			},
		});

		expect(tree).toEqual([
			{
				id: "/invoice_number",
				name: "invoice_number",
				path: "/invoice_number",
				key: "invoice_number",
				label: "invoice_number",
				type: "string",
				jsonCell: false,
				description: "Invoice identifier",
				children: [],
			},
			{
				id: "/total",
				name: "total",
				path: "/total",
				key: "total",
				label: "total",
				type: "number",
				jsonCell: false,
				children: [],
			},
		]);
	});

	it("keeps object nodes selectable while exposing nested object children", () => {
		const tree = buildDatasetFieldTree({
			type: "object",
			properties: {
				vendor: {
					type: "object",
					description: "Vendor details",
					properties: {
						name: { type: "string" },
						address: {
							type: "object",
							properties: {
								city: { type: "string" },
							},
						},
					},
				},
			},
		});

		expect(tree).toEqual([
			{
				id: "/vendor",
				name: "vendor",
				path: "/vendor",
				key: "vendor",
				label: "vendor",
				type: "object",
				jsonCell: true,
				description: "Vendor details",
				children: [
					{
						id: "/vendor/name",
						name: "name",
						path: "/vendor/name",
						key: "vendor_name",
						label: "name",
						type: "string",
						jsonCell: false,
						children: [],
					},
					{
						id: "/vendor/address",
						name: "address",
						path: "/vendor/address",
						key: "vendor_address",
						label: "address",
						type: "object",
						jsonCell: true,
						children: [
							{
								id: "/vendor/address/city",
								name: "city",
								path: "/vendor/address/city",
								key: "vendor_address_city",
								label: "city",
								type: "string",
								jsonCell: false,
								children: [],
							},
						],
					},
				],
			},
		]);
	});

	it("treats arrays as JSON cells without exposing item children", () => {
		const tree = buildDatasetFieldTree({
			type: "object",
			properties: {
				line_items: {
					type: "array",
					description: "Invoice line items",
					items: {
						type: "object",
						properties: {
							description: { type: "string" },
							amount: { type: "number" },
						},
					},
				},
			},
		});

		expect(tree).toEqual([
			{
				id: "/line_items",
				name: "line_items",
				path: "/line_items",
				key: "line_items",
				label: "line_items",
				type: "array",
				jsonCell: true,
				description: "Invoice line items",
				children: [],
			},
		]);
	});

	it("uses JSON Pointer escaping for slash and tilde property names", () => {
		const tree = buildDatasetFieldTree({
			type: "object",
			properties: {
				"a/b": { type: "string" },
				"tilde~name": {
					type: "object",
					properties: {
						"deep/value~": { type: "number" },
					},
				},
			},
		});

		expect(tree[0]).toMatchObject({
			id: "/a~1b",
			name: "a/b",
			path: "/a~1b",
			key: "a_b",
			label: "a/b",
		});
		expect(tree[1]).toMatchObject({
			id: "/tilde~0name",
			name: "tilde~name",
			path: "/tilde~0name",
			key: "tilde_name",
			label: "tilde~name",
		});
		expect(tree[1]?.children[0]).toMatchObject({
			id: "/tilde~0name/deep~1value~0",
			name: "deep/value~",
			path: "/tilde~0name/deep~1value~0",
			key: "tilde_name_deep_value",
			label: "deep/value~",
		});
	});

	it("makes top-level escaped paths distinct from nested paths with the same readable key", () => {
		const first = buildDatasetFieldTree({
			type: "object",
			properties: {
				"a/b": { type: "string" },
				a: {
					type: "object",
					properties: {
						b: { type: "string" },
					},
				},
			},
		});
		const second = buildDatasetFieldTree({
			type: "object",
			properties: {
				"a/b": { type: "string" },
				a: {
					type: "object",
					properties: {
						b: { type: "string" },
					},
				},
			},
		});

		const escapedTopLevel = first[0];
		const nested = first[1]?.children[0];

		expect(escapedTopLevel?.path).toBe("/a~1b");
		expect(nested?.path).toBe("/a/b");
		expect(escapedTopLevel?.key).toMatch(/^a_b__[a-z0-9]+$/);
		expect(nested?.key).toMatch(/^a_b__[a-z0-9]+$/);
		expect(escapedTopLevel?.key).not.toBe(nested?.key);
		expect(first).toEqual(second);
	});

	it("makes sibling properties with the same sanitized key distinct", () => {
		const tree = buildDatasetFieldTree({
			type: "object",
			properties: {
				"a b": { type: "string" },
				"a-b": { type: "string" },
				a_b: { type: "string" },
			},
		});
		const keys = tree.map((node) => node.key);

		expect(new Set(keys).size).toBe(3);
		for (const key of keys) {
			expect(key).toMatch(/^a_b__[a-z0-9]+$/);
		}
	});

	it("gives empty property names usable labels and keys without colliding with field", () => {
		const tree = buildDatasetFieldTree({
			type: "object",
			properties: {
				"": { type: "string" },
				field: { type: "string" },
			},
		});
		const emptyNameNode = tree[0];
		const fieldNode = tree[1];

		expect(emptyNameNode).toMatchObject({
			name: "",
			path: "/",
			label: "(empty)",
		});
		expect(emptyNameNode?.key).toBeTruthy();
		expect(emptyNameNode?.key).not.toMatch(/\s/);
		expect(emptyNameNode?.key).not.toBe(fieldNode?.key);
		expect(fieldNode?.key).toBeTruthy();
	});

	it("keeps generated keys within the backend limit and free of whitespace", () => {
		const longName = ` ${"long ".repeat(40)}field `;
		const tree = buildDatasetFieldTree({
			type: "object",
			properties: {
				[longName]: { type: "string" },
			},
		});

		expect(Array.from(tree[0]?.key ?? "")).toHaveLength(120);
		expect(tree[0]?.key).not.toMatch(/\s/);
	});

	it("caps generated labels for long Unicode property names without changing the path", () => {
		const longName = `invoice_${"😀".repeat(180)}`;
		const tree = buildDatasetFieldTree({
			type: "object",
			properties: {
				[longName]: { type: "string" },
			},
		});
		const node = tree[0];

		expect(node?.path).toBe(`/${longName}`);
		expect(Array.from(node?.label ?? "")).toHaveLength(160);
		expect(node?.label).toBe(Array.from(longName).slice(0, 160).join(""));
		expect(Array.from(node?.key ?? "").length).toBeLessThanOrEqual(120);
	});

	it("returns an empty tree for invalid or non-object schema input", () => {
		expect(buildDatasetFieldTree(null)).toEqual([]);
		expect(buildDatasetFieldTree([])).toEqual([]);
		expect(buildDatasetFieldTree({ type: "string" })).toEqual([]);
		expect(buildDatasetFieldTree({ type: "object", properties: [] })).toEqual([]);
	});
});
