import { describe, expect, it } from "vitest";

import type { OCRRecipeResponse } from "./api";
import { normalizeSearchText, summarizeRecipe } from "./recipe-summary";

function recipe(json: Record<string, unknown>): OCRRecipeResponse {
	return {
		id: "recipe-1",
		title: "Carte de identitate",
		description: "Romanian identity card extraction fields.",
		json,
		counter: 3,
		category_id: null,
		category: null,
		created_at: "2026-06-20T00:00:00Z",
		updated_at: "2026-06-20T00:00:00Z"
	};
}

describe("OCR recipe summary", () => {
	it("summarizes every field, required fields, nested arrays, objects, and enums", () => {
		const summary = summarizeRecipe(
			recipe({
				type: "object",
				properties: {
					name: { type: "string", description: "Full name" },
					lines: { type: "array", items: { type: "object" } },
					currency: { type: "string" },
					invoice_number: { type: "string" },
					status: { enum: ["paid", "open"] },
					meta: { properties: { source: { type: "string" } } },
					total: { type: "number" }
				},
				required: ["lines", "name"]
			})
		);

		expect(summary.fieldCount).toBe(7);
		expect(summary.requiredCount).toBe(2);
		expect(summary.fields.map((field) => [field.key, field.type, field.required])).toEqual([
			["lines", "object[]", true],
			["name", "string", true],
			["currency", "string", false],
			["invoice_number", "string", false],
			["meta", "object", false],
			["status", "enum", false],
			["total", "number", false]
		]);
		expect(summary.searchText).toContain("carte de identitate");
		expect(summary.searchText).toContain("full name");
		expect(summary.prettyJson).toContain('"properties"');
	});

	it("handles missing properties and unknown types", () => {
		const summary = summarizeRecipe(recipe({ type: "object", properties: { mystery: {} } }));

		expect(summary.fieldCount).toBe(1);
		expect(summary.requiredCount).toBe(0);
		expect(summary.fields[0]).toMatchObject({
			key: "mystery",
			type: "unknown",
			required: false
		});
	});

	it("normalizes search text for accents, separators, and whitespace", () => {
		expect(normalizeSearchText("  Număr_factură  CĂRȚI-test  ")).toBe("numar factura carti test");
	});

	it("includes localized category titles in search text", () => {
		const summary = summarizeRecipe({
			...recipe({ type: "object" }),
			category_id: "category-1",
			category: {
				id: "category-1",
				title: { en: "Invoices", ro: "Facturi" },
				created_at: "2026-06-19T00:00:00Z",
				updated_at: "2026-06-19T00:00:00Z"
			}
		});

		expect(summary.searchText).toContain("invoices");
		expect(summary.searchText).toContain("facturi");
	});
});
