import { describe, expect, it } from "vitest";

import { apiErrorMessage, publicApiErrorMessage } from "./api-errors";

describe("api error message helpers", () => {
	it("uses API error text for non-server responses", () => {
		expect(publicApiErrorMessage(409, { error: "name already exists" }, "fallback")).toBe(
			"name already exists"
		);
	});

	it("uses fallback text for server responses", () => {
		expect(publicApiErrorMessage(503, { error: "database password leaked" }, "fallback")).toBe(
			"fallback"
		);
	});

	it("uses fallback text when no API error text exists", () => {
		expect(apiErrorMessage({ message: "ignored" }, "fallback")).toBe("fallback");
	});
});
