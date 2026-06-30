import { publicApiErrorMessage } from "$lib/client/api-errors";

type ClientFetch = typeof fetch;
type JsonObject = Record<string, unknown>;

export type DashboardRange = "7d" | "30d" | "90d";

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

export async function fetchDashboardSummary(
	fetchFn: ClientFetch,
	range: DashboardRange
): Promise<DashboardSummaryResponse> {
	const response = await fetchFn(`/api/dashboard/summary?range=${encodeURIComponent(range)}`, {
		method: "GET"
	});
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to load dashboard"));
	}
	if (!isDashboardSummaryResponse(json)) {
		throw new Error("Invalid dashboard response");
	}

	return json;
}

async function readResponseJSON(response: Response): Promise<unknown> {
	let text: string;
	try {
		text = await response.text();
	} catch {
		return null;
	}

	if (!text.trim()) return null;

	try {
		return JSON.parse(text) as unknown;
	} catch {
		return null;
	}
}

function isJsonObject(value: unknown): value is JsonObject {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isFiniteNumber(value: unknown): value is number {
	return typeof value === "number" && Number.isFinite(value);
}

function isNonNegativeInteger(value: unknown): value is number {
	return isFiniteNumber(value) && Number.isInteger(value) && value >= 0;
}

function isRatio(value: unknown): value is number {
	return isFiniteNumber(value) && value >= 0 && value <= 1;
}

function isNullableString(value: unknown): value is string | null {
	return value === null || typeof value === "string";
}

function isDashboardRange(value: unknown): value is DashboardRange {
	return value === "7d" || value === "30d" || value === "90d";
}

function isDashboardRangeResponse(value: unknown): value is DashboardSummaryResponse["range"] {
	return (
		isJsonObject(value) &&
		isDashboardRange(value.key) &&
		typeof value.start_at === "string" &&
		typeof value.end_at === "string" &&
		value.bucket === "day"
	);
}

function isDashboardMetricsResponse(
	value: unknown
): value is DashboardSummaryResponse["metrics"] {
	return (
		isJsonObject(value) &&
		isNonNegativeInteger(value.documents_processed) &&
		isNonNegativeInteger(value.pages_processed) &&
		isNonNegativeInteger(value.jobs_completed) &&
		isNonNegativeInteger(value.jobs_failed) &&
		isNonNegativeInteger(value.jobs_processing) &&
		isRatio(value.completion_rate) &&
		isNonNegativeInteger(value.credits_spent) &&
		isNonNegativeInteger(value.dataset_count) &&
		isNonNegativeInteger(value.schema_count)
	);
}

function isDashboardDocumentBucket(
	value: unknown
): value is DashboardSummaryResponse["document_buckets"][number] {
	return (
		isJsonObject(value) &&
		typeof value.date === "string" &&
		isNonNegativeInteger(value.documents_processed)
	);
}

function isDashboardRecentDocument(
	value: unknown
): value is DashboardSummaryResponse["recent_documents"][number] {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.original_filename === "string" &&
		isNullableString(value.schema_id) &&
		isNullableString(value.schema_name) &&
		isNonNegativeInteger(value.page_count) &&
		typeof value.created_at === "string"
	);
}

function isDashboardSchemaThroughput(
	value: unknown
): value is DashboardSummaryResponse["schema_throughput"][number] {
	return (
		isJsonObject(value) &&
		isNullableString(value.schema_id) &&
		typeof value.schema_name === "string" &&
		isNonNegativeInteger(value.documents_processed)
	);
}

function isDashboardRecentDataset(
	value: unknown
): value is DashboardSummaryResponse["dataset_summary"]["recent"][number] {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.name === "string" &&
		typeof value.schema_name === "string" &&
		isNonNegativeInteger(value.field_count) &&
		typeof value.created_at === "string"
	);
}

function isDashboardDatasetSummary(
	value: unknown
): value is DashboardSummaryResponse["dataset_summary"] {
	return (
		isJsonObject(value) &&
		isNonNegativeInteger(value.total_count) &&
		Array.isArray(value.recent) &&
		value.recent.every(isDashboardRecentDataset)
	);
}

function isDashboardCreditSummary(
	value: unknown
): value is DashboardSummaryResponse["credit_summary"] {
	return (
		isJsonObject(value) &&
		isNonNegativeInteger(value.available_credits) &&
		isNonNegativeInteger(value.credits_spent) &&
		typeof value.low_credit === "boolean"
	);
}

function isDashboardOnboarding(
	value: unknown
): value is DashboardSummaryResponse["onboarding"] {
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

function isDashboardWarningSection(
	value: unknown
): value is DashboardSummaryResponse["warnings"][number]["section"] {
	return (
		value === "recent_documents" ||
		value === "schema_throughput" ||
		value === "dataset_summary" ||
		value === "credit_summary"
	);
}

function isDashboardWarning(value: unknown): value is DashboardSummaryResponse["warnings"][number] {
	return (
		isJsonObject(value) &&
		isDashboardWarningSection(value.section) &&
		typeof value.message === "string"
	);
}

export function isDashboardSummaryResponse(
	value: unknown
): value is DashboardSummaryResponse {
	return (
		isJsonObject(value) &&
		isDashboardRangeResponse(value.range) &&
		isDashboardMetricsResponse(value.metrics) &&
		Array.isArray(value.document_buckets) &&
		value.document_buckets.every(isDashboardDocumentBucket) &&
		Array.isArray(value.recent_documents) &&
		value.recent_documents.every(isDashboardRecentDocument) &&
		Array.isArray(value.schema_throughput) &&
		value.schema_throughput.every(isDashboardSchemaThroughput) &&
		isDashboardDatasetSummary(value.dataset_summary) &&
		isDashboardCreditSummary(value.credit_summary) &&
		isDashboardOnboarding(value.onboarding) &&
		Array.isArray(value.warnings) &&
		value.warnings.every(isDashboardWarning)
	);
}
