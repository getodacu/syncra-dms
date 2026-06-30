import { describe, expect, it, vi } from "vitest";

const jsonResponse = (body: unknown, init: ResponseInit = {}) =>
	new Response(JSON.stringify(body), {
		status: 200,
		headers: { "content-type": "application/json" },
		...init,
	});

describe("billing client helpers", () => {
	it("uses a shared credit balance query key", async () => {
		const { CREDIT_BALANCE_QUERY_KEY } = await import("./billing");

		expect(CREDIT_BALANCE_QUERY_KEY).toEqual(["billing", "balance"]);
	});

	it("fetches the authenticated credit balance through the frontend proxy", async () => {
		const { fetchCreditBalance } = await import("./billing");
		const fetchFn = vi.fn().mockResolvedValue(
			jsonResponse({
				user_id: "user-1",
				available_credits: 1234,
			})
		);

		await expect(fetchCreditBalance(fetchFn)).resolves.toEqual({
			user_id: "user-1",
			available_credits: 1234,
		});
		expect(fetchFn).toHaveBeenCalledWith("/api/billing/balance", { method: "GET" });
	});

	it("uses a fallback message for server errors", async () => {
		const { fetchCreditBalance } = await import("./billing");
		const fetchFn = vi.fn().mockResolvedValue(
			jsonResponse({ error: "Billing offline" }, { status: 502 })
		);

		await expect(fetchCreditBalance(fetchFn)).rejects.toThrow("Failed to load credit balance");
	});

	it("uses a fallback error when the response is not JSON", async () => {
		const { fetchCreditBalance } = await import("./billing");
		const fetchFn = vi
			.fn()
			.mockResolvedValue(new Response("not json", { status: 502 }));

		await expect(fetchCreditBalance(fetchFn)).rejects.toThrow(
			"Failed to load credit balance"
		);
	});

	it("rejects invalid credit balance payloads", async () => {
		const { fetchCreditBalance } = await import("./billing");
		const fetchFn = vi.fn().mockResolvedValue(
			jsonResponse({
				user_id: "user-1",
				available_credits: "1234",
			})
		);

		await expect(fetchCreditBalance(fetchFn)).rejects.toThrow(
			"Invalid credit balance response"
		);
	});
});
