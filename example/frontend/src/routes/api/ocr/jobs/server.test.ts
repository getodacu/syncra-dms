import { beforeEach, describe, expect, it, vi } from "vitest";

import { OCRApiError } from "$lib/server/ocr";
import { SchemaApiError } from "$lib/server/schemas";
import { DELETE, GET, POST } from "./+server";
import type { RequestEvent } from "./$types";

const MAX_OCR_FILE_BYTES = 20 << 20;
const MULTIPART_OVERHEAD_ALLOWANCE_BYTES = 64 * 1024;

const { createOCRJobMock, deleteOCRJobsMock, listOCRJobsMock, listSchemasMock, OCRApiErrorMock, SchemaApiErrorMock } =
	vi.hoisted(() => {
		class MockOCRApiError extends Error {
			status: number;

			constructor(status: number, message: string) {
				super(message);
				this.name = "OCRApiError";
				this.status = status;
			}
		}

		class MockSchemaApiError extends Error {
			status: number;

			constructor(status: number, message: string) {
				super(message);
				this.name = "SchemaApiError";
				this.status = status;
			}
		}

		return {
			createOCRJobMock: vi.fn(),
			deleteOCRJobsMock: vi.fn(),
			listOCRJobsMock: vi.fn(),
			listSchemasMock: vi.fn(),
			OCRApiErrorMock: MockOCRApiError,
			SchemaApiErrorMock: MockSchemaApiError
		};
	});

vi.mock("$lib/server/ocr", () => ({
	createOCRJob: createOCRJobMock,
	deleteOCRJobs: deleteOCRJobsMock,
	listOCRJobs: listOCRJobsMock,
	OCRApiError: OCRApiErrorMock,
	isOCRApiError: (error: unknown) => error instanceof OCRApiErrorMock
}));

vi.mock("$lib/server/schemas", () => ({
	listSchemas: listSchemasMock,
	SchemaApiError: SchemaApiErrorMock,
	isSchemaApiError: (error: unknown) => error instanceof SchemaApiErrorMock
}));

function createEvent(formData: FormData, user: unknown = { id: "user-1" }) {
	return {
		request: new Request("http://localhost/api/ocr/jobs", {
			method: "POST",
			body: formData
		}),
		fetch: vi.fn(),
		locals: { user }
	} as unknown as RequestEvent;
}

function createRequestEvent(request: Request, user: unknown = { id: "user-1" }) {
	return {
		request,
		fetch: vi.fn(),
		locals: { user }
	} as unknown as RequestEvent;
}

function createGetEvent(url = "http://localhost/api/ocr/jobs", user: unknown = { id: "user-1" }) {
	return {
		url: new URL(url),
		fetch: vi.fn(),
		locals: { user }
	} as unknown as RequestEvent;
}

function streamRequest(bodySize: number) {
	return new Request("http://localhost/api/ocr/jobs", {
		method: "POST",
		body: new ReadableStream<Uint8Array>({
			start(controller) {
				controller.enqueue(new Uint8Array(bodySize).fill(123));
				controller.close();
			}
		}),
		duplex: "half"
	} as RequestInit & { duplex: "half" });
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

function validSchema(id = "schema-1") {
	return {
		id,
		created_at: "2026-05-27T00:00:00Z",
		updated_at: "2026-05-27T00:00:00Z",
		user_id: "user-1",
		name: "Invoice",
		description: "Invoice extraction",
		strict: true,
		schema: { type: "object" }
	};
}

function validJob() {
	return {
		id: "job-1",
		created_at: "2026-05-27T00:00:00Z",
		original_filename: "invoice.pdf",
		mime_type: "application/pdf",
		file_size: 12,
		page_count: 2,
		schema_id: "schema-1",
		schema_name: "Invoice",
		has_inline_schema: false,
		document_id: null,
		status: "queued"
	};
}

function validJobList() {
	return {
		jobs: [validJob()],
		next_cursor: "cursor-1"
	};
}

function validFormData(
	file = new File([new Uint8Array([1])], "invoice.pdf", { type: "application/pdf" })
) {
	const formData = new FormData();
	formData.set("file", file);
	formData.set("schema_id", "schema-1");
	return formData;
}

function forwardedOCRJobInput() {
	const call = createOCRJobMock.mock.calls[0];
	expect(call).toBeDefined();
	return call as [
		RequestEvent["fetch"],
		{ file: File; schemaId: string | undefined; userId: string }
	];
}

async function expectForwardedFile(actual: File, expected: File) {
	expect(actual).toBeInstanceOf(File);
	expect(actual.name).toBe(expected.name);
	expect(actual.type).toBe(expected.type);
	expect(actual.size).toBe(expected.size);
	expect(new Uint8Array(await actual.arrayBuffer())).toEqual(
		new Uint8Array(await expected.arrayBuffer())
	);
}

describe("OCR jobs API endpoint", () => {
	beforeEach(() => {
		createOCRJobMock.mockReset();
		deleteOCRJobsMock.mockReset();
		listOCRJobsMock.mockReset();
		listSchemasMock.mockReset();
	});

	it("returns 401 for unauthenticated OCR job list requests", async () => {
		const response = await GET(createGetEvent("http://localhost/api/ocr/jobs", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listOCRJobsMock).not.toHaveBeenCalled();
	});

	it("lists owned OCR jobs with pagination parameters", async () => {
		const result = validJobList();
		listOCRJobsMock.mockResolvedValue(result);
		const event = createGetEvent(
			"http://localhost/api/ocr/jobs?cursor=cursor-0&size=50&sort=desc&status=processing",
			{ id: "user-1" }
		);

		const response = await GET(event);

		expect(listOCRJobsMock).toHaveBeenCalledWith(event.fetch, {
			userId: "user-1",
			cursor: "cursor-0",
			size: "50",
			sort: "desc",
			status: "processing"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("deletes owned OCR jobs", async () => {
		const result = { deleted_ids: ["job-1", "job-2"], deleted_count: 2 };
		deleteOCRJobsMock.mockResolvedValue(result);
		const request = new Request("http://localhost/api/ocr/jobs", {
			method: "DELETE",
			headers: { "content-type": "application/json" },
			body: JSON.stringify({ ids: ["job-1", "job-2"] })
		});
		const event = createRequestEvent(request, { id: "user-1" });

		const response = await DELETE(event);

		expect(deleteOCRJobsMock).toHaveBeenCalledWith(event.fetch, ["job-1", "job-2"], {
			userId: "user-1"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("returns 400 for invalid OCR job delete bodies", async () => {
		const request = new Request("http://localhost/api/ocr/jobs", {
			method: "DELETE",
			headers: { "content-type": "application/json" },
			body: JSON.stringify({ ids: "job-1" })
		});

		const response = await DELETE(createRequestEvent(request));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid OCR job delete request" });
		expect(deleteOCRJobsMock).not.toHaveBeenCalled();
	});

	it("returns 401 for unauthenticated OCR job requests", async () => {
		const response = await POST(createEvent(validFormData(), null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listSchemasMock).not.toHaveBeenCalled();
		expect(createOCRJobMock).not.toHaveBeenCalled();
	});

	it("returns 400 when file is missing", async () => {
		const formData = new FormData();
		formData.set("schema_id", "schema-1");

		const response = await POST(createEvent(formData));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "file is required" });
		expect(listSchemasMock).not.toHaveBeenCalled();
		expect(createOCRJobMock).not.toHaveBeenCalled();
	});

	it("returns 400 when content-length exceeds the upload limit plus multipart overhead", async () => {
		const response = await POST(
			createRequestEvent(
				new Request("http://localhost/api/ocr/jobs", {
					method: "POST",
					headers: {
						"content-length": String(
							MAX_OCR_FILE_BYTES + MULTIPART_OVERHEAD_ALLOWANCE_BYTES + 1
						)
					},
					body: validFormData()
				})
			)
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "request body too large" });
		expect(listSchemasMock).not.toHaveBeenCalled();
		expect(createOCRJobMock).not.toHaveBeenCalled();
	});

	it("returns 400 when a streamed body without content-length exceeds the request limit", async () => {
		const response = await POST(
			createRequestEvent(
				streamRequest(MAX_OCR_FILE_BYTES + MULTIPART_OVERHEAD_ALLOWANCE_BYTES + 1)
			)
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "request body too large" });
		expect(listSchemasMock).not.toHaveBeenCalled();
		expect(createOCRJobMock).not.toHaveBeenCalled();
	});

	it("returns 400 when file exceeds the upload limit", async () => {
		const formData = validFormData(
			new File([new Uint8Array(MAX_OCR_FILE_BYTES + 1)], "large.pdf", {
				type: "application/pdf"
			})
		);

		const response = await POST(createEvent(formData));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "file exceeds max upload size" });
		expect(listSchemasMock).not.toHaveBeenCalled();
		expect(createOCRJobMock).not.toHaveBeenCalled();
	});

	it("returns 400 for invalid form data", async () => {
		const response = await POST(
			createRequestEvent(
				new Request("http://localhost/api/ocr/jobs", {
					method: "POST",
					headers: { "content-type": "multipart/form-data; boundary=syncra" },
					body: "--not-the-declared-boundary\r\nbroken"
				})
			)
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid form data" });
		expect(listSchemasMock).not.toHaveBeenCalled();
		expect(createOCRJobMock).not.toHaveBeenCalled();
	});

	it("returns 400 when schema_id is not owned by the current user", async () => {
		listSchemasMock.mockResolvedValue([validSchema("schema-2")]);
		const event = createEvent(validFormData(), { id: "user-1" });

		const response = await POST(event);

		expect(listSchemasMock).toHaveBeenCalledWith(event.fetch, { userId: "user-1" });
		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid schema_id" });
		expect(createOCRJobMock).not.toHaveBeenCalled();
	});

	it("forwards file, schema id, and the authenticated user id for valid requests", async () => {
		const result = validJob();
		createOCRJobMock.mockResolvedValue(result);
		listSchemasMock.mockResolvedValue([validSchema()]);
		const file = new File([new Uint8Array([1, 2, 3])], "invoice.pdf", {
			type: "application/pdf"
		});
		const event = createEvent(validFormData(file), { id: "user-1" });

		const response = await POST(event);

		expect(listSchemasMock).toHaveBeenCalledWith(event.fetch, { userId: "user-1" });
		const [fetchFn, input] = forwardedOCRJobInput();
		expect(fetchFn).toBe(event.fetch);
		await expectForwardedFile(input.file, file);
		expect(input.schemaId).toBe("schema-1");
		expect(input.userId).toBe("user-1");
		expect(response.status).toBe(202);
		expect(await responseJson(response)).toEqual(result);
	});

	it("forwards OCR-only jobs when schema_id is omitted", async () => {
		const result = validJob();
		createOCRJobMock.mockResolvedValue(result);
		const file = new File([new Uint8Array([1, 2, 3])], "invoice.pdf", {
			type: "application/pdf"
		});
		const formData = new FormData();
		formData.set("file", file);
		const event = createEvent(formData, { id: "user-1" });

		const response = await POST(event);

		expect(listSchemasMock).not.toHaveBeenCalled();
		const [fetchFn, input] = forwardedOCRJobInput();
		expect(fetchFn).toBe(event.fetch);
		await expectForwardedFile(input.file, file);
		expect(input.schemaId).toBeUndefined();
		expect(input.userId).toBe("user-1");
		expect(response.status).toBe(202);
		expect(await responseJson(response)).toEqual(result);
	});

	it("returns 202 for OCR job creation responses", async () => {
		const result = { ...validJob(), status: "completed", document_id: "document-1" };
		createOCRJobMock.mockResolvedValue(result);
		listSchemasMock.mockResolvedValue([validSchema()]);

		const response = await POST(createEvent(validFormData(), { id: "user-1" }));

		expect(createOCRJobMock).toHaveBeenCalled();
		expect(response.status).toBe(202);
		expect(await responseJson(response)).toEqual(result);
	});

	it("ignores client supplied user_id", async () => {
		createOCRJobMock.mockResolvedValue(validJob());
		listSchemasMock.mockResolvedValue([validSchema()]);
		const file = new File([new Uint8Array([1])], "scan.png", { type: "image/png" });
		const formData = validFormData(file);
		formData.set("user_id", "attacker");
		const event = createEvent(formData, { id: "user-1" });

		await POST(event);

		const [fetchFn, input] = forwardedOCRJobInput();
		expect(fetchFn).toBe(event.fetch);
		await expectForwardedFile(input.file, file);
		expect(input.schemaId).toBe("schema-1");
		expect(input.userId).toBe("user-1");
	});

	it("preserves schema service client errors", async () => {
		listSchemasMock.mockRejectedValue(new SchemaApiError(409, "schema conflict"));

		const response = await POST(createEvent(validFormData()));

		expect(response.status).toBe(409);
		expect(await responseJson(response)).toEqual({ error: "schema conflict" });
		expect(createOCRJobMock).not.toHaveBeenCalled();
	});

	it("normalizes schema service server errors", async () => {
		listSchemasMock.mockRejectedValue(new SchemaApiError(503, "Schema service unavailable"));

		const response = await POST(createEvent(validFormData()));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
		expect(createOCRJobMock).not.toHaveBeenCalled();
	});

	it("preserves OCR job service client errors", async () => {
		listSchemasMock.mockResolvedValue([validSchema()]);
		createOCRJobMock.mockRejectedValue(new OCRApiError(422, "invalid document"));

		const response = await POST(createEvent(validFormData()));

		expect(response.status).toBe(422);
		expect(await responseJson(response)).toEqual({ error: "invalid document" });
	});

	it("preserves OCR job page-limit errors", async () => {
		createOCRJobMock.mockRejectedValue(
			new OCRApiError(400, "document must have at most 150 pages with a schema")
		);
		listSchemasMock.mockResolvedValue([validSchema()]);

		const response = await POST(createEvent(validFormData()));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "document must have at most 150 pages with a schema"
		});
	});

	it("normalizes OCR job service server errors", async () => {
		listSchemasMock.mockResolvedValue([validSchema()]);
		createOCRJobMock.mockRejectedValue(new OCRApiError(503, "OCR service unavailable"));

		const response = await POST(createEvent(validFormData()));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
