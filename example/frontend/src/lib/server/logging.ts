import { env } from "$env/dynamic/private";

export const requestIdHeader = "X-Request-ID";

type LogLevel = "DEBUG" | "INFO" | "WARN" | "ERROR";

export type LogAttrs = Record<string, unknown>;

export type LogEntry = {
	time: string;
	level: LogLevel;
	msg: string;
	[key: string]: unknown;
};

export type LogSink = (entry: LogEntry) => void;

export interface Logger {
	debug(message: string, attrs?: LogAttrs): void;
	info(message: string, attrs?: LogAttrs): void;
	warn(message: string, attrs?: LogAttrs): void;
	error(message: string, attrs?: LogAttrs): void;
	child(attrs: LogAttrs): Logger;
}

export type JSONLoggerOptions = {
	debug?: boolean;
	attrs?: LogAttrs;
	now?: () => Date;
	sink?: LogSink;
};

type SanitizedValue =
	| string
	| number
	| boolean
	| null
	| SanitizedValue[]
	| { [key: string]: SanitizedValue };

const redacted = "[REDACTED]";

class JSONLogger implements Logger {
	readonly #debug: boolean;
	readonly #attrs: LogAttrs;
	readonly #now: () => Date;
	readonly #sink: LogSink;

	constructor(options: Required<JSONLoggerOptions>) {
		this.#debug = options.debug;
		this.#attrs = options.attrs;
		this.#now = options.now;
		this.#sink = options.sink;
	}

	debug(message: string, attrs: LogAttrs = {}) {
		this.#write("DEBUG", message, attrs);
	}

	info(message: string, attrs: LogAttrs = {}) {
		this.#write("INFO", message, attrs);
	}

	warn(message: string, attrs: LogAttrs = {}) {
		this.#write("WARN", message, attrs);
	}

	error(message: string, attrs: LogAttrs = {}) {
		this.#write("ERROR", message, attrs);
	}

	child(attrs: LogAttrs): Logger {
		return new JSONLogger({
			debug: this.#debug,
			attrs: { ...this.#attrs, ...attrs },
			now: this.#now,
			sink: this.#sink
		});
	}

	#write(level: LogLevel, message: string, attrs: LogAttrs) {
		if (level === "DEBUG" && !this.#debug) return;

		const sanitizedAttrs = sanitizeAttrs({ ...this.#attrs, ...attrs });
		const entry: LogEntry = {
			...sanitizedAttrs,
			time: this.#now().toISOString(),
			level,
			msg: message
		};

		try {
			this.#sink(entry);
		} catch {
			// Logging must not affect request handling.
		}
	}
}

export function createJSONLogger(options: JSONLoggerOptions = {}): Logger {
	return new JSONLogger({
		debug: options.debug ?? false,
		attrs: options.attrs ?? {},
		now: options.now ?? (() => new Date()),
		sink: options.sink ?? stdoutSink
	});
}

export const rootLogger = createJSONLogger({
	debug: debugFromEnv(),
	attrs: { service: "syncra-frontend" }
});

export function debugFromEnv() {
	const value = (privateEnv("DEBUG") ?? "").trim();
	if (value === "") return false;
	switch (value.toLowerCase()) {
		case "1":
		case "t":
		case "true":
			return true;
		case "0":
		case "f":
		case "false":
			return false;
		default:
			return false;
	}
}

export function cleanRequestId(raw: string | null | undefined) {
	const value = raw?.trim() ?? "";
	if (!value || value.length > 128) return "";
	for (const char of value) {
		if (/\s/.test(char) || char.charCodeAt(0) < 32 || char.charCodeAt(0) === 127) {
			return "";
		}
	}
	return value;
}

export function userAgentClass(raw: string | null | undefined) {
	const ua = raw?.trim().toLowerCase() ?? "";
	if (!ua) return "absent";
	if (ua.includes("mozilla/")) return "browser";
	if (ua.includes("curl/")) return "curl";
	if (ua.includes("go-http-client/")) return "go-http-client";
	return "api-client";
}

export function safeError(error: unknown) {
	if (error instanceof Error) {
		return redactSecretText(error.message || error.name);
	}
	if (typeof error === "string") return redactSecretText(error);
	return redactSecretText(String(error));
}

function sanitizeAttrs(attrs: LogAttrs) {
	const sanitized: Record<string, SanitizedValue> = {};
	for (const [key, value] of Object.entries(attrs)) {
		if (value === undefined) continue;
		sanitized[key] = sanitizeValue(value, key, new WeakSet<object>());
	}
	return sanitized;
}

function sanitizeValue(value: unknown, key: string, seen: WeakSet<object>): SanitizedValue {
	if (isSensitiveKey(key)) return redacted;
	if (value === undefined) return null;
	if (value === null || typeof value === "number" || typeof value === "boolean") return value;
	if (typeof value === "string") return redactSecretText(value);
	if (typeof value === "bigint") return value.toString();
	if (value instanceof Date) return value.toISOString();
	if (value instanceof Error) return safeError(value);
	if (Array.isArray(value)) {
		return value.map((item) => sanitizeValue(item, key, seen));
	}
	if (typeof value === "object") {
		if (seen.has(value)) return "[Circular]";
		seen.add(value);

		const sanitized: Record<string, SanitizedValue> = {};
		for (const [childKey, childValue] of Object.entries(value as Record<string, unknown>)) {
			if (childValue !== undefined) {
				sanitized[childKey] = sanitizeValue(childValue, childKey, seen);
			}
		}
		seen.delete(value);
		return sanitized;
	}
	return redactSecretText(String(value));
}

function isSensitiveKey(key: string) {
	const normalized = key.toLowerCase().replace(/[\s-]+/g, "_");
	return (
		normalized === "code" ||
		normalized === "otp" ||
		normalized.includes("verification_code") ||
		normalized.includes("reset_code") ||
		normalized.includes("password") ||
		normalized.includes("token") ||
		normalized.includes("secret") ||
		normalized.includes("authorization") ||
		normalized.includes("cookie") ||
		normalized.includes("api_key") ||
		normalized.includes("apikey")
	);
}

function redactSecretText(value: string) {
	return value
		.replace(
			/\b(password|token|secret|authorization|cookie|api[_-]?key)=([^&\s]+)/gi,
			"$1=[REDACTED]"
		)
		.replace(/\b(code|otp)=([^&\s]+)/gi, "$1=[REDACTED]");
}

function stdoutSink(entry: LogEntry) {
	const output = JSON.stringify(entry);
	if (!output) return;
	const stdout = nodeEnvProcess()?.stdout;
	if (stdout && typeof stdout.write === "function") {
		stdout.write(`${output}\n`);
	}
}

function privateEnv(key: string) {
	return env[key] || nodeEnvProcess()?.env?.[key];
}

function nodeEnvProcess() {
	if (typeof process !== "undefined") return process;
	return (
		globalThis as typeof globalThis & {
			process?: NodeJS.Process;
		}
	).process;
}
