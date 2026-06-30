import { beforeEach, describe, expect, it, vi } from "vitest";

import { load } from "./+layout.server";
import type { LayoutServerLoadEvent } from "./$types";

const { getCreditBalanceMock, isBillingApiErrorMock, BillingApiErrorMock } = vi.hoisted(() => {
	class MockBillingApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "BillingApiError";
			this.status = status;
		}
	}

	return {
		getCreditBalanceMock: vi.fn(),
		isBillingApiErrorMock: (error: unknown) => error instanceof MockBillingApiError,
		BillingApiErrorMock: MockBillingApiError
	};
});

vi.mock("$lib/server/billing", () => ({
	getCreditBalance: getCreditBalanceMock,
	isBillingApiError: isBillingApiErrorMock,
	BillingApiError: BillingApiErrorMock
}));

type AppLayoutLoadData = {
	user: unknown;
	initialCreditBalance: { user_id: string; available_credits: number } | null;
	initialCreditBalanceError: string | null;
};

function loadEvent(user: unknown = { id: "user-1" }) {
	return {
		fetch: vi.fn(),
		locals: { user },
		url: new URL("http://localhost/app")
	} as unknown as LayoutServerLoadEvent;
}

describe("app layout load", () => {
	beforeEach(() => {
		getCreditBalanceMock.mockReset();
	});

	it("loads the current credit balance for the signed-in user", async () => {
		getCreditBalanceMock.mockResolvedValue({
			user_id: "user-1",
			available_credits: 1234
		});

		const event = loadEvent();
		const data = (await load(event)) as AppLayoutLoadData;

		expect(getCreditBalanceMock).toHaveBeenCalledWith(event.fetch, "user-1");
		expect(data.user).toEqual({ id: "user-1" });
		expect(data.initialCreditBalance).toEqual({
			user_id: "user-1",
			available_credits: 1234
		});
		expect(data.initialCreditBalanceError).toBeNull();
	});

	it("returns a layout error when balance lookup fails", async () => {
		getCreditBalanceMock.mockRejectedValue(new BillingApiErrorMock(503, "Billing offline"));

		const data = (await load(loadEvent())) as AppLayoutLoadData;

		expect(data.initialCreditBalance).toBeNull();
		expect(data.initialCreditBalanceError).toBe("Failed to load credit balance");
	});

	it("returns a safe null balance shape without a signed-in user", async () => {
		const data = (await load(loadEvent(null))) as AppLayoutLoadData;

		expect(getCreditBalanceMock).not.toHaveBeenCalled();
		expect(data.user).toBeNull();
		expect(data.initialCreditBalance).toBeNull();
		expect(data.initialCreditBalanceError).toBeNull();
	});
});
