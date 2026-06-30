import { beforeEach, describe, expect, it, vi } from "vitest";

import { BillingApiError } from "$lib/server/billing";
import { GET } from "./+server";

const { fetchBillingInvoicePDFMock, BillingApiErrorMock } = vi.hoisted(() => {
	class MockBillingApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "BillingApiError";
			this.status = status;
		}
	}

	return {
		fetchBillingInvoicePDFMock: vi.fn(),
		BillingApiErrorMock: MockBillingApiError,
	};
});

vi.mock("$lib/server/billing", () => ({
	fetchBillingInvoicePDF: fetchBillingInvoicePDFMock,
	BillingApiError: BillingApiErrorMock,
	isBillingApiError: (error: unknown) => error instanceof BillingApiErrorMock,
}));

function createEvent(
	url = "http://localhost/api/billing/invoices/invoice-1/pdf",
	user: unknown = { id: "user-1" }
) {
	return {
		request: new Request(url),
		url: new URL(url),
		params: { id: "invoice-1" },
		fetch: vi.fn(),
		locals: { user },
	} as any;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("billing invoice PDF API endpoint", () => {
	beforeEach(() => {
		fetchBillingInvoicePDFMock.mockReset();
	});

	it("returns 401 when unauthenticated", async () => {
		const response = await GET(createEvent(undefined, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(fetchBillingInvoicePDFMock).not.toHaveBeenCalled();
	});

	it("proxies PDF preview and download responses for the authenticated user", async () => {
		fetchBillingInvoicePDFMock.mockResolvedValue({
			body: new Response("%PDF-test").body,
			headers: new Headers({
				"content-type": "application/pdf",
				"content-disposition": 'attachment; filename="invoice-1.pdf"',
			}),
			status: 200,
		});
		const event = createEvent("http://localhost/api/billing/invoices/invoice-1/pdf?download=1");

		const response = await GET(event);

		expect(fetchBillingInvoicePDFMock).toHaveBeenCalledWith(event.fetch, "user-1", "invoice-1", {
			download: true,
		});
		expect(response.status).toBe(200);
		expect(response.headers.get("content-type")).toBe("application/pdf");
		expect(response.headers.get("content-disposition")).toBe('attachment; filename="invoice-1.pdf"');
		await expect(response.text()).resolves.toBe("%PDF-test");
	});

	it("normalizes backend billing errors", async () => {
		fetchBillingInvoicePDFMock.mockRejectedValue(new BillingApiError(404, "billing invoice PDF not found"));

		const response = await GET(createEvent());

		expect(response.status).toBe(404);
		expect(await responseJson(response)).toEqual({ error: "billing invoice PDF not found" });
	});
});
