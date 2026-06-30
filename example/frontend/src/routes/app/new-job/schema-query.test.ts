import { describe, expect, it } from "vitest";

import {
	buildNewJobPath,
	isSelectedSchemaIdUnavailable,
	schemaIdFromSearchParams,
} from "./schema-query";

describe("new job schema query helpers", () => {
	it("builds new-job paths with encoded schema ids", () => {
		expect(buildNewJobPath("schema-1")).toBe("/app/new-job?schema_id=schema-1");
		expect(buildNewJobPath(" schema:/next page ")).toBe(
			"/app/new-job?schema_id=schema%3A%2Fnext+page"
		);
	});

	it("omits blank schema ids from new-job paths", () => {
		expect(buildNewJobPath("")).toBe("/app/new-job");
		expect(buildNewJobPath("   ")).toBe("/app/new-job");
	});

	it("parses and normalizes schema ids from search params", () => {
		expect(schemaIdFromSearchParams(new URLSearchParams("schema_id=schema-1"))).toBe("schema-1");
		expect(schemaIdFromSearchParams(new URLSearchParams("schema_id=%20schema-2%20"))).toBe(
			"schema-2"
		);
		expect(schemaIdFromSearchParams(new URLSearchParams("schema_id=%20"))).toBe("");
		expect(schemaIdFromSearchParams(new URLSearchParams("other=value"))).toBe("");
	});

	it("detects unavailable selected schema ids only after schemas load", () => {
		const schemas = [{ id: "schema-1" }, { id: "schema-2" }];

		expect(
			isSelectedSchemaIdUnavailable({
				selectedSchemaId: "missing",
				schemas,
				schemasLoaded: false,
			})
		).toBe(false);
		expect(
			isSelectedSchemaIdUnavailable({
				selectedSchemaId: "",
				schemas,
				schemasLoaded: true,
			})
		).toBe(false);
		expect(
			isSelectedSchemaIdUnavailable({
				selectedSchemaId: "schema-1",
				schemas,
				schemasLoaded: true,
			})
		).toBe(false);
		expect(
			isSelectedSchemaIdUnavailable({
				selectedSchemaId: "missing",
				schemas,
				schemasLoaded: true,
			})
		).toBe(true);
	});
});
