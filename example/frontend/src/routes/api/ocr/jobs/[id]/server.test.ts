import { beforeEach, describe, expect, it, vi } from "vitest";

import { OCRApiError } from "$lib/server/ocr";
import { GET } from "./+server";
import type { RequestEvent } from "./$types";

const { getOCRJobMock, OCRApiErrorMock } = vi.hoisted(() => {
	class MockOCRApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "OCRApiError";
			this.status = status;
		}
	}

	return { getOCRJobMock: vi.fn(), OCRApiErrorMock: MockOCRApiError };
});

vi.mock("$lib/server/ocr", () => ({
	getOCRJob: getOCRJobMock,
	OCRApiError: OCRApiErrorMock,
	isOCRApiError: (error: unknown) => error instanceof OCRApiErrorMock
}));

function createEvent(id = "job-1", user: unknown = { id: "user-1" }) {
	return {
		params: { id },
		fetch: vi.fn(),
		locals: { user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

function validJob(overrides: Record<string, unknown> = {}) {
	return {
		id: "job-1",
		file_size: 12,
		page_count: 2,
		document_id: null,
		status: "queued",
		...overrides
	};
}

describe("OCR job API endpoint", () => {
	beforeEach(() => {
		getOCRJobMock.mockReset();
	});

	it("returns 401 for unauthenticated OCR job requests", async () => {
		const response = await GET(createEvent("job-1", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(getOCRJobMock).not.toHaveBeenCalled();
	});

	it("returns an owned job", async () => {
		const result = validJob();
		getOCRJobMock.mockResolvedValue(result);
		const event = createEvent("job-1", { id: "user-1" });

		const response = await GET(event);

		expect(getOCRJobMock).toHaveBeenCalledWith(event.fetch, "job-1", { userId: "user-1" });
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("preserves backend 404 errors", async () => {
		getOCRJobMock.mockRejectedValue(new OCRApiError(404, "OCR job not found"));

		const response = await GET(createEvent());

		expect(response.status).toBe(404);
		expect(await responseJson(response)).toEqual({ error: "OCR job not found" });
	});

	it("normalizes backend server errors", async () => {
		getOCRJobMock.mockRejectedValue(new OCRApiError(503, "OCR service unavailable"));

		const response = await GET(createEvent());

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});

});
