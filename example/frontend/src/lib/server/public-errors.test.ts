import { describe, expect, it } from "vitest";

import {
	jsonPublicErrorResponse,
	publicErrorMessage,
	publicErrorStatus
} from "./public-errors";

async function responseJson(response: Response) {
	return response.json();
}

describe("public server error helpers", () => {
	it("preserves non-server error messages and statuses", async () => {
		const response = jsonPublicErrorResponse(409, "schema conflict");

		expect(response.status).toBe(409);
		await expect(responseJson(response)).resolves.toEqual({ error: "schema conflict" });
	});

	it("uses fallback messages and 502 status for server errors", async () => {
		const response = jsonPublicErrorResponse(503, "database unavailable");

		expect(response.status).toBe(502);
		await expect(responseJson(response)).resolves.toEqual({
			error: "A server error occurred. Please try again."
		});
		expect(publicErrorMessage(500, "internal detail", "fallback")).toBe("fallback");
		expect(publicErrorStatus(500)).toBe(502);
	});
});
