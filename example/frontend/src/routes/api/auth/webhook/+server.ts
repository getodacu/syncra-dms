import { json } from "@sveltejs/kit";

import { deleteWebhook, getWebhook, isAuthApiError, saveWebhook } from "$lib/server/auth";
import type { WebhookEvent } from "$lib/server/auth";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

const SUPPORTED_WEBHOOK_EVENTS = new Set<WebhookEvent>([
	"job.started",
	"job.failed",
	"job.succeeded"
]);

function authErrorResponse(error: unknown) {
	if (isAuthApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	throw error;
}

function cookieHeader(request: Request) {
	return request.headers.get("cookie");
}

function userId(locals: App.Locals) {
	return locals.user?.id ?? "";
}

function webhookUrl(value: unknown) {
	if (typeof value !== "string") {
		return { ok: false as const, error: "url must be an absolute http or https URL" };
	}

	const url = value.trim();
	try {
		const parsed = new URL(url);
		if (parsed.protocol !== "http:" && parsed.protocol !== "https:") {
			return { ok: false as const, error: "url must be an absolute http or https URL" };
		}
	} catch {
		return { ok: false as const, error: "url must be an absolute http or https URL" };
	}

	return { ok: true as const, url };
}

function webhookEvents(value: unknown) {
	if (value === undefined) {
		return { ok: true as const, eventsActive: [] };
	}
	if (!Array.isArray(value)) {
		return { ok: false as const, error: "events_active must be an array" };
	}

	if (!value.every((event): event is WebhookEvent => SUPPORTED_WEBHOOK_EVENTS.has(event))) {
		return { ok: false as const, error: "events_active contains unsupported events" };
	}

	return { ok: true as const, eventsActive: value };
}

export const GET: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) return json({ error: "authentication required" }, { status: 401 });

	try {
		return json(await getWebhook(fetch, cookieHeader(request), userId(locals)));
	} catch (error) {
		return authErrorResponse(error);
	}
};

export const POST: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) return json({ error: "authentication required" }, { status: 401 });

	let body: unknown;
	try {
		body = await request.json();
	} catch {
		return json({ error: "invalid webhook request" }, { status: 400 });
	}

	if (typeof body !== "object" || body === null || Array.isArray(body)) {
		return json({ error: "invalid webhook request" }, { status: 400 });
	}

	const urlResult = webhookUrl((body as Record<string, unknown>).url);
	if (!urlResult.ok) return json({ error: urlResult.error }, { status: 400 });
	const eventsResult = webhookEvents((body as Record<string, unknown>).events_active);
	if (!eventsResult.ok) return json({ error: eventsResult.error }, { status: 400 });

	try {
		const webhook = await saveWebhook(fetch, cookieHeader(request), {
			userId: userId(locals),
			url: urlResult.url,
			eventsActive: eventsResult.eventsActive
		});
		return json(webhook);
	} catch (error) {
		return authErrorResponse(error);
	}
};

export const DELETE: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) return json({ error: "authentication required" }, { status: 401 });

	try {
		return json(await deleteWebhook(fetch, cookieHeader(request), userId(locals)));
	} catch (error) {
		return authErrorResponse(error);
	}
};
