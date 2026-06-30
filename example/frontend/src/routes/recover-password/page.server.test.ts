import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

const {
	AuthApiErrorMock,
	loggerMock,
	privateEnv,
	requestPasswordResetMock,
	resetPasswordMock,
	sendPasswordResetEmailMock
} = vi.hoisted(() => {
	class MockAuthApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "AuthApiError";
			this.status = status;
		}
	}

	return {
		AuthApiErrorMock: MockAuthApiError,
		loggerMock: {
			info: vi.fn(),
			error: vi.fn(),
			warn: vi.fn(),
			debug: vi.fn(),
			child: vi.fn()
		},
		privateEnv: {} as Record<string, string | undefined>,
		requestPasswordResetMock: vi.fn(),
		resetPasswordMock: vi.fn(),
		sendPasswordResetEmailMock: vi.fn()
	};
});

vi.mock("$env/dynamic/private", () => ({ env: privateEnv }));

vi.mock("$lib/server/auth", () => ({
	AuthApiError: AuthApiErrorMock,
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock,
	requestPasswordReset: requestPasswordResetMock,
	resetPassword: resetPasswordMock
}));

vi.mock("$lib/server/mail", () => ({
	sendPasswordResetEmail: sendPasswordResetEmailMock
}));

import { actions } from "./+page.server";

const GENERIC_RESET_MESSAGE = "If an account exists, we sent a reset link.";

function requestEvent(email = " ADA@EXAMPLE.COM ") {
	const data = new FormData();
	data.set("email", email);
	const fetchMock = vi.fn();

	return {
		event: {
			request: new Request("http://localhost/recover-password", {
				method: "POST",
				body: data
			}),
			fetch: fetchMock,
			locals: { logger: loggerMock },
			url: new URL("http://localhost/recover-password")
		},
		fetchMock
	};
}

function resetEvent(overrides: Record<string, string> = {}) {
	const data = new FormData();
	data.set("email", overrides.email ?? "ADA@EXAMPLE.COM");
	data.set("token", overrides.token ?? " reset-token ");
	data.set("password", overrides.password ?? "newpassword123");
	data.set("confirmPassword", overrides.confirmPassword ?? "newpassword123");
	const fetchMock = vi.fn();

	return {
		event: {
			request: new Request("http://localhost/recover-password", {
				method: "POST",
				body: data
			}),
			fetch: fetchMock,
			locals: { logger: loggerMock },
			url: new URL("http://localhost/recover-password")
		},
		fetchMock
	};
}

describe("recover password page actions", () => {
	beforeEach(() => {
		requestPasswordResetMock.mockReset();
		resetPasswordMock.mockReset();
		sendPasswordResetEmailMock.mockReset();
		sendPasswordResetEmailMock.mockResolvedValue(undefined);
		loggerMock.error.mockReset();
		privateEnv.SYNCRA_APP_ORIGIN = "https://app.example.com/";
	});

	afterEach(() => {
		vi.restoreAllMocks();
		for (const key of Object.keys(privateEnv)) delete privateEnv[key];
	});

	it("uses named actions for requesting and confirming password resets", () => {
		expect("default" in actions).toBe(false);
		expect(actions.request).toEqual(expect.any(Function));
		expect(actions.reset).toEqual(expect.any(Function));
	});

	it("sends trusted reset links by email and returns a generic success message", async () => {
		requestPasswordResetMock.mockResolvedValue({
			ok: true,
			resetToken: "reset-token",
			resetExpiresAt: "2026-06-11T13:30:00Z"
		});
		const { event, fetchMock } = requestEvent();

		const result = await actions.request(event as never);

		expect(requestPasswordResetMock).toHaveBeenCalledWith(fetchMock, "ada@example.com");
		expect(sendPasswordResetEmailMock).toHaveBeenCalledWith({
			to: "ada@example.com",
			resetUrl:
				"https://app.example.com/recover-password?email=ada%40example.com&token=reset-token",
			expiresAt: "2026-06-11T13:30:00Z"
		});
		expect(result).toEqual({
			values: { email: "ada@example.com" },
			success: GENERIC_RESET_MESSAGE
		});
		expect(JSON.stringify(result)).not.toContain("reset-token");
	});

	it("keeps request success generic when reset email delivery fails", async () => {
		requestPasswordResetMock.mockResolvedValue({
			ok: true,
			resetToken: "reset-token",
			resetExpiresAt: "2026-06-11T13:30:00Z"
		});
		sendPasswordResetEmailMock.mockRejectedValue(new Error("smtp down"));
		const { event } = requestEvent();

		const result = await actions.request(event as never);

		expect(result).toEqual({
			values: { email: "ada@example.com" },
			success: GENERIC_RESET_MESSAGE
		});
		expect(JSON.stringify(result)).not.toContain("reset-token");
		expect(loggerMock.error).toHaveBeenCalledWith("auth.password_reset_email_failed", {
			email: "ada@example.com",
			error: "smtp down"
		});
		expect(JSON.stringify(loggerMock.error.mock.calls)).not.toContain("reset-token");
	});

	it("does not send email when the backend returns no reset token", async () => {
		requestPasswordResetMock.mockResolvedValue({ ok: true });
		const { event } = requestEvent();

		const result = await actions.request(event as never);

		expect(sendPasswordResetEmailMock).not.toHaveBeenCalled();
		expect(result).toEqual({
			values: { email: "ada@example.com" },
			success: GENERIC_RESET_MESSAGE
		});
	});

	it("resets the password and redirects to login", async () => {
		resetPasswordMock.mockResolvedValue({ ok: true });
		const { event, fetchMock } = resetEvent();

		await expect(actions.reset(event as never)).rejects.toMatchObject({
			status: 303,
			location: "/login?email=ada%40example.com&reset=1"
		});
		expect(resetPasswordMock).toHaveBeenCalledWith(fetchMock, {
			email: "ada@example.com",
			token: "reset-token",
			password: "newpassword123"
		});
	});

	it("validates password confirmation before resetting", async () => {
		const { event } = resetEvent({ confirmPassword: "different" });

		const result = await actions.reset(event as never);

		expect(result).toEqual({
			status: 400,
			data: {
				values: { email: "ada@example.com", token: "reset-token" },
				fieldErrors: { confirmPassword: "Passwords do not match." }
			}
		});
		expect(resetPasswordMock).not.toHaveBeenCalled();
	});

	it("returns generic reset errors for auth service failures", async () => {
		resetPasswordMock.mockRejectedValue(new AuthApiErrorMock(503, "database connection failed"));
		const { event } = resetEvent();

		const result = await actions.reset(event as never);

		expect(result).toEqual({
			status: 502,
			data: {
				values: { email: "ada@example.com", token: "reset-token" },
				error: "Unable to reset password. Please try again."
			}
		});
	});
});
