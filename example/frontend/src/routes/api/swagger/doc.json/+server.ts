import { env } from "$env/dynamic/private";
import { json } from "@sveltejs/kit";

import type { RequestHandler } from "./$types";

function apiBaseUrl() {
	return (privateEnv("SYNCRA_API_BASE_URL") || "http://localhost:8080").replace(/\/+$/, "");
}

function privateEnv(key: string) {
	return env[key] || nodeEnv()[key];
}

function nodeEnv() {
	if (typeof process !== "undefined") return process.env;
	return (
		globalThis as typeof globalThis & {
			process?: { env?: Record<string, string | undefined> };
		}
	).process?.env ?? {};
}

export const GET: RequestHandler = async ({ fetch }) => {
	let response: Response;
	try {
		response = await fetch(`${apiBaseUrl()}/swagger-public/doc.json`, {
			headers: { accept: "application/json" }
		});
	} catch {
		return json({ error: "Swagger service unavailable" }, { status: 502 });
	}

	let body: string;
	try {
		body = await response.text();
	} catch {
		return json({ error: "Swagger service unavailable" }, { status: 502 });
	}

	if (!response.ok) {
		return json(
			{ error: "Swagger document request failed" },
			{ status: response.status >= 500 ? 502 : response.status }
		);
	}

	return new Response(body, {
		headers: {
			"cache-control": "no-store",
			"content-type": response.headers.get("content-type") ?? "application/json; charset=utf-8"
		}
	});
};