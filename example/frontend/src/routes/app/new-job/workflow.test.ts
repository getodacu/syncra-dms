import { describe, expect, it } from "vitest";
import {
	attemptedRows,
	canSubmitJobs,
	displayJobStatus,
	getSubmittableRows,
	isInsufficientCreditsError,
	jobStatusPatch,
	percentComplete,
	shouldShowSubmitButton,
	type JobWorkflowRow,
} from "./workflow";

const row = (patch: Partial<JobWorkflowRow>): JobWorkflowRow => ({
	createStatus: "pending",
	jobStatus: "unknown",
	pollFailures: 0,
	...patch,
});

describe("new job workflow", () => {
	it("shows submit after files are selected, with or without a schema", () => {
		expect(shouldShowSubmitButton("", 0)).toBe(false);
		expect(shouldShowSubmitButton("schema-1", 0)).toBe(false);
		expect(shouldShowSubmitButton("", 1)).toBe(true);
		expect(shouldShowSubmitButton("schema-1", 1)).toBe(true);
	});

	it("submits only pending and create-failed rows", () => {
		const rows = [
			row({ createStatus: "pending" }),
			row({ createStatus: "creating" }),
			row({ createStatus: "queued", jobStatus: "completed" }),
			row({ createStatus: "create_failed", jobStatus: "failed" }),
		];

		expect(getSubmittableRows(rows).map((item) => item.createStatus)).toEqual([
			"pending",
			"create_failed",
		]);
		expect(
			canSubmitJobs({
				selectedSchemaId: "",
				rows,
				isSubmitting: false,
				isPolling: false,
			})
		).toBe(true);
		expect(
			canSubmitJobs({
				selectedSchemaId: "schema-1",
				rows: [row({ createStatus: "queued", jobStatus: "completed" })],
				isSubmitting: false,
				isPolling: false,
			})
		).toBe(false);
		expect(
			canSubmitJobs({
				selectedSchemaId: "schema-1",
				rows,
				isSubmitting: true,
				isPolling: false,
			})
		).toBe(false);
		expect(
			canSubmitJobs({
				selectedSchemaId: "schema-1",
				rows,
				isSubmitting: false,
				isPolling: true,
			})
		).toBe(false);
	});

	it("keeps unknown backend statuses polling until the retry cap", () => {
		expect(jobStatusPatch({ status: "waiting", document_id: null }, 0, 3)).toEqual({
			jobStatus: "unknown",
			documentId: undefined,
			errorMessage: "Unexpected OCR job status: waiting",
			pollFailures: 1,
		});
		expect(jobStatusPatch({ status: "waiting", document_id: null }, 2, 3)).toEqual({
			jobStatus: "failed",
			documentId: undefined,
			errorMessage: "Unexpected OCR job status: waiting",
			pollFailures: 3,
		});
	});

	it("resets failures when a known status arrives", () => {
		expect(
			jobStatusPatch(
				{ status: "processing", document_id: "doc-1", error_message: "old transient detail" },
				2,
				3
			)
		).toEqual({
			jobStatus: "processing",
			documentId: "doc-1",
			errorMessage: "old transient detail",
			pollFailures: 0,
		});
	});

	it("reports attempt progress from queued and create-failed rows", () => {
		const rows = [
			row({ createStatus: "pending" }),
			row({ createStatus: "queued" }),
			row({ createStatus: "create_failed" }),
		];

		expect(attemptedRows(rows)).toHaveLength(2);
		expect(percentComplete(attemptedRows(rows).length, rows.length)).toBe(67);
	});

	it("maps workflow state to visual job status", () => {
		expect(displayJobStatus(row({ createStatus: "pending" }))).toBe("queueing");
		expect(displayJobStatus(row({ createStatus: "creating" }))).toBe("queueing");
		expect(displayJobStatus(row({ createStatus: "queued", jobStatus: "queued" }))).toBe(
			"processing"
		);
		expect(displayJobStatus(row({ createStatus: "queued", jobStatus: "processing" }))).toBe(
			"processing"
		);
		expect(displayJobStatus(row({ createStatus: "queued", jobStatus: "unknown" }))).toBe(
			"processing"
		);
		expect(displayJobStatus(row({ createStatus: "queued", jobStatus: "completed" }))).toBe(
			"completed"
		);
		expect(displayJobStatus(row({ createStatus: "create_failed", jobStatus: "unknown" }))).toBe(
			"failed"
		);
		expect(displayJobStatus(row({ createStatus: "queued", jobStatus: "failed" }))).toBe("failed");
	});

	it("identifies insufficient credit failures", () => {
		expect(isInsufficientCreditsError("insufficient credits")).toBe(true);
		expect(isInsufficientCreditsError("Insufficient credits. Buy credits to continue.")).toBe(true);
		expect(isInsufficientCreditsError("Failed to create OCR job")).toBe(false);
		expect(isInsufficientCreditsError(undefined)).toBe(false);
	});
});
