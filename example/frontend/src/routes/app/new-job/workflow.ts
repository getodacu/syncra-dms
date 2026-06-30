export const MAX_UNKNOWN_STATUS_POLLS = 3;

export type CreateStatus = "pending" | "creating" | "queued" | "create_failed";
export type JobStatus = "queued" | "processing" | "completed" | "failed" | "unknown";

export type JobWorkflowRow = {
	createStatus: CreateStatus;
	jobStatus: JobStatus;
	jobId?: string;
	pollFailures: number;
};

export type JobDisplayStatus = "queueing" | "processing" | "completed" | "failed";

export type OCRJobStatusInput = {
	status: string;
	document_id: string | null;
	error_message?: string;
};

export type JobStatusPatch = {
	jobStatus: JobStatus;
	documentId?: string;
	errorMessage?: string;
	pollFailures: number;
};

export function isTerminalStatus(status: string) {
	return status === "completed" || status === "failed";
}

export function normalizeJobStatus(status: string): JobStatus {
	if (
		status === "queued" ||
		status === "processing" ||
		status === "completed" ||
		status === "failed"
	) {
		return status;
	}

	return "unknown";
}

export function isSubmittableRow(row: Pick<JobWorkflowRow, "createStatus">) {
	return row.createStatus === "pending" || row.createStatus === "create_failed";
}

export function getSubmittableRows<T extends Pick<JobWorkflowRow, "createStatus">>(rows: T[]) {
	return rows.filter(isSubmittableRow);
}

export function shouldShowSubmitButton(selectedSchemaId: string, fileCount: number) {
	void selectedSchemaId;
	return fileCount > 0;
}

export function canSubmitJobs(input: {
	selectedSchemaId: string;
	rows: Pick<JobWorkflowRow, "createStatus">[];
	isSubmitting: boolean;
	isPolling: boolean;
}) {
	return (
		shouldShowSubmitButton(input.selectedSchemaId, input.rows.length) &&
		getSubmittableRows(input.rows).length > 0 &&
		!input.isSubmitting &&
		!input.isPolling
	);
}

export function attemptedRows<T extends Pick<JobWorkflowRow, "createStatus">>(rows: T[]) {
	return rows.filter((row) => row.createStatus === "queued" || row.createStatus === "create_failed");
}

export function percentComplete(completed: number, total: number) {
	return total === 0 ? 0 : Math.round((completed / total) * 100);
}

export function isInsufficientCreditsError(message: string | undefined) {
	return message?.toLowerCase().includes("insufficient credits") ?? false;
}

export function displayJobStatus(
	row: Pick<JobWorkflowRow, "createStatus" | "jobStatus">
): JobDisplayStatus {
	if (row.createStatus === "create_failed" || row.jobStatus === "failed") return "failed";
	if (row.jobStatus === "completed") return "completed";
	if (row.createStatus === "pending" || row.createStatus === "creating") return "queueing";
	return "processing";
}

export function jobStatusPatch(
	job: OCRJobStatusInput,
	currentPollFailures = 0,
	maxUnknownStatusPolls = MAX_UNKNOWN_STATUS_POLLS
): JobStatusPatch {
	const jobStatus = normalizeJobStatus(job.status);
	if (jobStatus === "unknown") {
		const pollFailures = currentPollFailures + 1;
		return {
			jobStatus: pollFailures >= maxUnknownStatusPolls ? "failed" : "unknown",
			documentId: job.document_id ?? undefined,
			errorMessage: `Unexpected OCR job status: ${job.status}`,
			pollFailures,
		};
	}

	return {
		jobStatus,
		documentId: job.document_id ?? undefined,
		errorMessage: job.error_message || undefined,
		pollFailures: 0,
	};
}
