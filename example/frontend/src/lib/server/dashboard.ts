import { apiBaseUrl, internalAPIHeaders } from "./internal-api";

type ServerFetch = typeof fetch;
type JsonObject = Record<string, unknown>;

export type DashboardRange = "7d" | "30d" | "90d";

export type DashboardSummaryOptions = {
	userId: string;
	range?: DashboardRange | string;
};

export type DashboardSummaryResponse = {
	range: {
		key: DashboardRange;
		start_at: string;
		end_at: string;
		bucket: "day";
	};
	metrics: {
		documents_processed: number;
		pages_processed: number;
		jobs_completed: number;
		jobs_failed: number;
		jobs_processing: number;
		completion_rate: number;
		credits_spent: number;
		dataset_count: number;
		schema_count: number;
	};
	document_buckets: Array<{ date: string; documents_processed: number }>;
	recent_documents: Array<{
		id: string;
		original_filename: string;
		schema_id: string | null;
		schema_name: string | null;
		page_count: number;
		created_at: string;
	}>;
	schema_throughput: Array<{
		schema_id: string | null;
		schema_name: string;
		documents_processed: number;
	}>;
	dataset_summary: {
		total_count: number;
		recent: Array<{
			id: string;
			name: string;
			schema_name: string;
			field_count: number;
			created_at: string;
		}>;
	};
	credit_summary: {
		available_credits: number;
		credits_spent: number;
		low_credit: boolean;
	};
	onboarding: {
		has_schema: boolean;
		has_completed_document: boolean;
		has_dataset: boolean;
		has_api_key: boolean;
		has_webhook: boolean;
		show_onboarding: boolean;
	};
	warnings: Array<{
		section: "recent_documents" | "schema_throughput" | "dataset_summary" | "credit_summary";
		message: string;
	}>;
};

export class DashboardApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = "DashboardApiError";
		this.status = status;
	}
}

function isJsonObject(value: unknown): value is JsonObject {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isFiniteNumber(value: unknown): value is number {
	return typeof value === "number" && Number.isFinite(value);
}

function isNullableString(value: unknown): value is string | null {
	return value === null || typeof value === "string";
}

function isDashboardRange(value: unknown): value is DashboardRange {
	return value === "7d" || value === "30d" || value === "90d";
}

function parseResponseJSON(text: string) {
	if (!text) return null;
	try {
		return JSON.parse(text) as unknown;
	} catch {
		return undefined;
	}
}

async function readResponseJSON(response: Response) {
	let text: string;
	try {
		text = await response.text();
	} catch {
		throw new DashboardApiError(503, "Dashboard service unavailable");
	}

	return parseResponseJSON(text);
}

function errorMessage(data: unknown, fallback: string) {
	if (
		isJsonObject(data) &&
		typeof data.error === "string" &&
		data.error.trim() !== ""
	) {
		return data.error;
	}
	return fallback;
}

function dashboardSummaryUrl(options: DashboardSummaryOptions) {
	const url = new URL(`${apiBaseUrl()}/api/dashboard/summary`);
	url.searchParams.set("user_id", options.userId);
	url.searchParams.set("range", options.range ?? "30d");
	return url.toString();
}

function isRangeResponse(value: unknown): value is DashboardSummaryResponse["range"] {
	return (
		isJsonObject(value) &&
		isDashboardRange(value.key) &&
		typeof value.start_at === "string" &&
		typeof value.end_at === "string" &&
		value.bucket === "day"
	);
}

function isMetricsResponse(value: unknown): value is DashboardSummaryResponse["metrics"] {
	return (
		isJsonObject(value) &&
		isFiniteNumber(value.documents_processed) &&
		isFiniteNumber(value.pages_processed) &&
		isFiniteNumber(value.jobs_completed) &&
		isFiniteNumber(value.jobs_failed) &&
		isFiniteNumber(value.jobs_processing) &&
		isFiniteNumber(value.completion_rate) &&
		isFiniteNumber(value.credits_spent) &&
		isFiniteNumber(value.dataset_count) &&
		isFiniteNumber(value.schema_count)
	);
}

function isDocumentBucket(value: unknown): value is DashboardSummaryResponse["document_buckets"][number] {
	return (
		isJsonObject(value) &&
		typeof value.date === "string" &&
		isFiniteNumber(value.documents_processed)
	);
}

function isRecentDocument(value: unknown): value is DashboardSummaryResponse["recent_documents"][number] {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.original_filename === "string" &&
		isNullableString(value.schema_id) &&
		isNullableString(value.schema_name) &&
		isFiniteNumber(value.page_count) &&
		typeof value.created_at === "string"
	);
}

function isSchemaThroughput(value: unknown): value is DashboardSummaryResponse["schema_throughput"][number] {
	return (
		isJsonObject(value) &&
		isNullableString(value.schema_id) &&
		typeof value.schema_name === "string" &&
		isFiniteNumber(value.documents_processed)
	);
}

function isRecentDataset(value: unknown): value is DashboardSummaryResponse["dataset_summary"]["recent"][number] {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.name === "string" &&
		typeof value.schema_name === "string" &&
		isFiniteNumber(value.field_count) &&
		typeof value.created_at === "string"
	);
}

function isDatasetSummary(value: unknown): value is DashboardSummaryResponse["dataset_summary"] {
	return (
		isJsonObject(value) &&
		isFiniteNumber(value.total_count) &&
		Array.isArray(value.recent) &&
		value.recent.every(isRecentDataset)
	);
}

function isCreditSummary(value: unknown): value is DashboardSummaryResponse["credit_summary"] {
	return (
		isJsonObject(value) &&
		isFiniteNumber(value.available_credits) &&
		isFiniteNumber(value.credits_spent) &&
		typeof value.low_credit === "boolean"
	);
}

function isOnboarding(value: unknown): value is DashboardSummaryResponse["onboarding"] {
	return (
		isJsonObject(value) &&
		typeof value.has_schema === "boolean" &&
		typeof value.has_completed_document === "boolean" &&
		typeof value.has_dataset === "boolean" &&
		typeof value.has_api_key === "boolean" &&
		typeof value.has_webhook === "boolean" &&
		typeof value.show_onboarding === "boolean"
	);
}

function isWarningSection(value: unknown): value is DashboardSummaryResponse["warnings"][number]["section"] {
	return (
		value === "recent_documents" ||
		value === "schema_throughput" ||
		value === "dataset_summary" ||
		value === "credit_summary"
	);
}

function isWarning(value: unknown): value is DashboardSummaryResponse["warnings"][number] {
	return isJsonObject(value) && isWarningSection(value.section) && typeof value.message === "string";
}

function isDashboardSummaryResponse(value: unknown): value is DashboardSummaryResponse {
	return (
		isJsonObject(value) &&
		isRangeResponse(value.range) &&
		isMetricsResponse(value.metrics) &&
		Array.isArray(value.document_buckets) &&
		value.document_buckets.every(isDocumentBucket) &&
		Array.isArray(value.recent_documents) &&
		value.recent_documents.every(isRecentDocument) &&
		Array.isArray(value.schema_throughput) &&
		value.schema_throughput.every(isSchemaThroughput) &&
		isDatasetSummary(value.dataset_summary) &&
		isCreditSummary(value.credit_summary) &&
		isOnboarding(value.onboarding) &&
		Array.isArray(value.warnings) &&
		value.warnings.every(isWarning)
	);
}

export function isDashboardApiError(error: unknown): error is DashboardApiError {
	return error instanceof DashboardApiError;
}

export async function getDashboardSummary(
	fetchFn: ServerFetch,
	options: DashboardSummaryOptions
) {
	const headers = internalAPIHeaders();
	if (!headers) {
		throw new DashboardApiError(500, "Dashboard service is not configured");
	}

	let response: Response;
	try {
		response = await fetchFn(dashboardSummaryUrl(options), {
			method: "GET",
			headers
		});
	} catch {
		throw new DashboardApiError(503, "Dashboard service unavailable");
	}

	const data = await readResponseJSON(response);

	if (!response.ok) {
		throw new DashboardApiError(response.status, errorMessage(data, "Dashboard request failed"));
	}

	if (!isDashboardSummaryResponse(data)) {
		throw new DashboardApiError(502, "Invalid dashboard response");
	}

	return data;
}
