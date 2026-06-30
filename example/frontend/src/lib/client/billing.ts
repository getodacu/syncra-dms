import { publicApiErrorMessage } from "$lib/client/api-errors";

export type CreditBalanceResponse = {
	user_id: string;
	available_credits: number;
};

type ClientFetch = (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>;

export const CREDIT_BALANCE_QUERY_KEY = ["billing", "balance"] as const;

export async function fetchCreditBalance(
	fetchFn: ClientFetch
): Promise<CreditBalanceResponse> {
	const response = await fetchFn("/api/billing/balance", { method: "GET" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(
			publicApiErrorMessage(response.status, json, "Failed to load credit balance")
		);
	}
	if (!isCreditBalanceResponse(json)) {
		throw new Error("Invalid credit balance response");
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

function isCreditBalanceResponse(value: unknown): value is CreditBalanceResponse {
	return (
		isRecord(value) &&
		typeof value.user_id === "string" &&
		typeof value.available_credits === "number" &&
		Number.isFinite(value.available_credits)
	);
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}
