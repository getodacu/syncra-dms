import { publicApiErrorMessage } from "$lib/client/api-errors";

import {
	buildCreditUsageHistoryQueryPath,
	type CreditUsageHistoryListQuery,
} from "./table-utils";

export type CreditUsageHistoryEntryType = "purchase" | "debit";

export type CreditUsageHistoryEntryResponse = {
	id: string;
	created_at: string;
	entry_type: CreditUsageHistoryEntryType;
	credits_delta: number;
	related_order_id?: string;
	related_job_id?: string;
};

export type CreditUsageHistoryListResponse = {
	credit_usage_history: CreditUsageHistoryEntryResponse[];
	next_cursor: string | null;
};

type ClientFetch = typeof fetch;

export async function fetchCreditUsageHistory(
	fetchFn: ClientFetch,
	query: CreditUsageHistoryListQuery
): Promise<CreditUsageHistoryListResponse> {
	const response = await fetchFn(buildCreditUsageHistoryQueryPath(query), { method: "GET" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(
			publicApiErrorMessage(response.status, json, "Failed to load credit usage history")
		);
	}
	if (!isCreditUsageHistoryListResponse(json)) {
		throw new Error("Invalid credit usage history response");
	}

	return json;
}

async function readResponseJSON(response: Response): Promise<unknown> {
	const text = await response.text();
	if (!text.trim()) return null;

	try {
		return JSON.parse(text);
	} catch {
		return null;
	}
}

function isCreditUsageHistoryListResponse(
	value: unknown
): value is CreditUsageHistoryListResponse {
	if (!isRecord(value)) return false;
	if (!Array.isArray(value.credit_usage_history)) return false;
	if (!(typeof value.next_cursor === "string" || value.next_cursor === null)) return false;

	return value.credit_usage_history.every(isCreditUsageHistoryEntryResponse);
}

function isCreditUsageHistoryEntryResponse(
	value: unknown
): value is CreditUsageHistoryEntryResponse {
	if (!isRecord(value)) return false;

	return (
		typeof value.id === "string" &&
		typeof value.created_at === "string" &&
		(value.entry_type === "purchase" || value.entry_type === "debit") &&
		typeof value.credits_delta === "number" &&
		Number.isFinite(value.credits_delta) &&
		(value.related_order_id === undefined || typeof value.related_order_id === "string") &&
		(value.related_job_id === undefined || typeof value.related_job_id === "string")
	);
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}
