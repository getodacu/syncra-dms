import { beforeEach, describe, expect, it, vi } from "vitest";

const { privateEnv } = vi.hoisted(() => ({
	privateEnv: {} as Record<string, string | undefined>
}));

vi.mock("$env/dynamic/private", () => ({ env: privateEnv }));

describe("internal API server helper", () => {
	beforeEach(() => {
		vi.resetModules();
		for (const key of Object.keys(privateEnv)) delete privateEnv[key];
		delete process.env.SYNCRA_API_BASE_URL;
		delete process.env.SYNCRA_INTERNAL_API_TOKEN;
	});

	it("builds internal API headers while preserving existing headers", async () => {
		const { INTERNAL_API_HEADER, internalAPIHeaders } = await import("./internal-api");
		privateEnv.SYNCRA_INTERNAL_API_TOKEN = " internal-token ";

		const headers = internalAPIHeaders({ "content-type": "application/json" });

		expect(headers).toBeInstanceOf(Headers);
		expect(headers?.get(INTERNAL_API_HEADER)).toBe("internal-token");
		expect(headers?.get("content-type")).toBe("application/json");
	});

	it("fails closed when the internal API token is not configured", async () => {
		const { internalAPIHeaders } = await import("./internal-api");

		expect(internalAPIHeaders()).toBeNull();
	});

	it("uses the private SvelteKit env before process env for the Go API base URL", async () => {
		const { apiBaseUrl } = await import("./internal-api");
		privateEnv.SYNCRA_API_BASE_URL = "http://private-api.test/";
		process.env.SYNCRA_API_BASE_URL = "http://process-api.test/";

		expect(apiBaseUrl()).toBe("http://private-api.test");
	});
});
