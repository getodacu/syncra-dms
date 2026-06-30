import { afterEach, describe, expect, it, vi } from "vitest";

const { privateEnv } = vi.hoisted(() => ({
	privateEnv: {} as Record<string, string | undefined>
}));

vi.mock("$env/dynamic/private", () => ({ env: privateEnv }));

import {
	cleanRequestId,
	createJSONLogger,
	debugFromEnv,
	safeError,
	userAgentClass,
	type LogEntry
} from "./logging";

function testLogger(debug = false) {
	const entries: LogEntry[] = [];
	const logger = createJSONLogger({
		debug,
		attrs: { service: "syncra-frontend" },
		now: () => new Date("2026-06-14T10:00:00.000Z"),
		sink: (entry) => entries.push(entry)
	});
	return { entries, logger };
}

describe("server logging", () => {
	afterEach(() => {
		for (const key of Object.keys(privateEnv)) delete privateEnv[key];
		vi.unstubAllEnvs();
	});

	it("emits structured JSON info entries", () => {
		const { entries, logger } = testLogger();

		logger.info("frontend.started", { component: "test", port: 5173 });

		expect(entries).toEqual([
			{
				time: "2026-06-14T10:00:00.000Z",
				level: "INFO",
				msg: "frontend.started",
				service: "syncra-frontend",
				component: "test",
				port: 5173
			}
		]);
	});

	it("filters debug entries by default", () => {
		const { entries, logger } = testLogger();

		logger.debug("debug.hidden");
		logger.info("info.visible");

		expect(entries).toHaveLength(1);
		expect(entries[0]).toMatchObject({ level: "INFO", msg: "info.visible" });
	});

	it("emits debug entries when enabled", () => {
		const { entries, logger } = testLogger(true);

		logger.debug("debug.visible");

		expect(entries).toHaveLength(1);
		expect(entries[0]).toMatchObject({ level: "DEBUG", msg: "debug.visible" });
	});

	it("merges child attrs and lets call attrs win", () => {
		const { entries, logger } = testLogger();

		logger.child({ component: "http", request_id: "request-1" }).info("http.done", {
			component: "billing",
			status: 200
		});

		expect(entries[0]).toMatchObject({
			component: "billing",
			request_id: "request-1",
			status: 200
		});
	});

	it("serializes errors without stack traces and redacts secret-looking values", () => {
		const { entries, logger } = testLogger();

		logger.error("auth.failed", {
			error: new Error("request failed token=abc123"),
			password: "secret-password",
			nested: {
				code: "123456",
				country_code: "RO"
			}
		});

		expect(entries[0]).toMatchObject({
			error: "request failed token=[REDACTED]",
			password: "[REDACTED]",
			nested: {
				code: "[REDACTED]",
				country_code: "RO"
			}
		});
		expect(JSON.stringify(entries[0])).not.toContain("stack");
		expect(JSON.stringify(entries[0])).not.toContain("abc123");
		expect(JSON.stringify(entries[0])).not.toContain("123456");
	});

	it("parses DEBUG using server-compatible bool values", () => {
		privateEnv.DEBUG = "true";
		expect(debugFromEnv()).toBe(true);

		privateEnv.DEBUG = "1";
		expect(debugFromEnv()).toBe(true);

		privateEnv.DEBUG = "false";
		expect(debugFromEnv()).toBe(false);

		privateEnv.DEBUG = "verbose";
		expect(debugFromEnv()).toBe(false);
	});

	it("cleans request ids with the same constraints as the Go server", () => {
		expect(cleanRequestId(" request-123 ")).toBe("request-123");
		expect(cleanRequestId("bad request")).toBe("");
		expect(cleanRequestId("bad\nrequest")).toBe("");
		expect(cleanRequestId("x".repeat(129))).toBe("");
		expect(cleanRequestId("")).toBe("");
	});

	it("classifies user agents", () => {
		expect(userAgentClass("Mozilla/5.0")).toBe("browser");
		expect(userAgentClass("curl/8.0")).toBe("curl");
		expect(userAgentClass("Go-http-client/1.1")).toBe("go-http-client");
		expect(userAgentClass("custom-client")).toBe("api-client");
		expect(userAgentClass("")).toBe("absent");
	});

	it("redacts secret-looking text in safeError", () => {
		expect(safeError(new Error("bad password=hunter2 code=123456"))).toBe(
			"bad password=[REDACTED] code=[REDACTED]"
		);
	});
});
