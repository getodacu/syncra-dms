import { describe, expect, it } from "vitest";

import { schemaHasFields } from "./schema-validation";

describe("new schema validation", () => {
	it("requires at least one top-level JSON schema field", () => {
		expect(schemaHasFields({ type: "object", properties: { total: { type: "number" } } })).toBe(
			true
		);
		expect(schemaHasFields({ type: "object", properties: {} })).toBe(false);
		expect(schemaHasFields({ type: "object" })).toBe(false);
		expect(schemaHasFields({})).toBe(false);
		expect(schemaHasFields({ type: "object", properties: [] })).toBe(false);
	});
});
