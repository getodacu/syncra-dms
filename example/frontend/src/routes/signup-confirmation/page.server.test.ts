import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

const {
	AuthApiErrorMock,
	deliveryErrorMessage,
	loggerMock,
	sendVerificationEmailMock,
	sendVerificationOTPMock,
	verifyEmailOTPMock
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
		deliveryErrorMessage: "We couldn't send the verification email. Please try again.",
		loggerMock: {
			info: vi.fn(),
			error: vi.fn(),
			warn: vi.fn(),
			debug: vi.fn(),
			child: vi.fn()
		},
		sendVerificationEmailMock: vi.fn(),
		sendVerificationOTPMock: vi.fn(),
		verifyEmailOTPMock: vi.fn()
	};
});

vi.mock("$lib/server/auth", () => ({
	AuthApiError: AuthApiErrorMock,
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock,
	sendVerificationOTP: sendVerificationOTPMock,
	verifyEmailOTP: verifyEmailOTPMock
}));

vi.mock("$lib/server/mail", () => ({
	sendVerificationEmail: sendVerificationEmailMock,
	VERIFICATION_EMAIL_DELIVERY_ERROR: deliveryErrorMessage
}));

import { actions } from "./+page.server";

function resendEvent(email = "ADA@EXAMPLE.COM") {
	const data = new FormData();
	data.set("email", email);
	const fetchMock = vi.fn();

	return {
		event: {
			request: new Request("http://localhost/signup-confirmation", {
				method: "POST",
				body: data
			}),
			fetch: fetchMock,
			locals: { logger: loggerMock }
		},
		fetchMock
	};
}

describe("signup confirmation page actions", () => {
	beforeEach(() => {
		sendVerificationEmailMock.mockReset();
		sendVerificationOTPMock.mockReset();
		verifyEmailOTPMock.mockReset();
		sendVerificationEmailMock.mockResolvedValue(undefined);
		loggerMock.info.mockReset();
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	it("uses named actions for both verify and resend submissions", () => {
		expect("default" in actions).toBe(false);
		expect(actions.verify).toEqual(expect.any(Function));
		expect(actions.resend).toEqual(expect.any(Function));
	});

	it("sends resent verification codes by email", async () => {
		sendVerificationOTPMock.mockResolvedValue({
			ok: true,
			verificationCode: "654321",
			verificationExpiresAt: "2026-06-02T13:30:00Z"
		});
		const { event, fetchMock } = resendEvent();

		const result = await actions.resend(event as never);

		expect(sendVerificationOTPMock).toHaveBeenCalledWith(fetchMock, "ada@example.com");
		expect(sendVerificationEmailMock).toHaveBeenCalledWith({
			to: "ada@example.com",
			code: "654321",
			expiresAt: "2026-06-02T13:30:00Z"
		});
		expect(loggerMock.info).toHaveBeenCalledWith("auth.verification_email_code_resent", {
			email: "ada@example.com"
		});
		expect(JSON.stringify(loggerMock.info.mock.calls)).not.toContain("654321");
		expect(result).toEqual({
			values: { email: "ada@example.com" },
			success: "A new verification code has been sent. Check your inbox and spam folder."
		});
		expect(JSON.stringify(result)).not.toContain("654321");
	});

	it("returns a generic delivery error when resend email delivery fails", async () => {
		sendVerificationOTPMock.mockResolvedValue({
			ok: true,
			verificationCode: "654321",
			verificationExpiresAt: "2026-06-02T13:30:00Z"
		});
		sendVerificationEmailMock.mockRejectedValue(new Error("smtp down"));
		const { event } = resendEvent();

		const result = await actions.resend(event as never);

		expect(result).toEqual({
			status: 502,
			data: {
				values: { email: "ada@example.com" },
				error: deliveryErrorMessage
			}
		});
		expect(JSON.stringify(result)).not.toContain("654321");
		expect(loggerMock.info).toHaveBeenCalledWith("auth.verification_email_code_resent", {
			email: "ada@example.com"
		});
		expect(JSON.stringify(loggerMock.info.mock.calls)).not.toContain("654321");
	});

	it("does not send email when resend returns no verification code", async () => {
		sendVerificationOTPMock.mockResolvedValue({ ok: true });
		const { event } = resendEvent();

		const result = await actions.resend(event as never);

		expect(sendVerificationEmailMock).not.toHaveBeenCalled();
		expect(loggerMock.info).not.toHaveBeenCalled();
		expect(result).toEqual({
			values: { email: "ada@example.com" },
			success: "A new verification code has been sent. Check your inbox and spam folder."
		});
	});

	it("maps backend resend auth errors without sending email", async () => {
		sendVerificationOTPMock.mockRejectedValue(new AuthApiErrorMock(400, "valid email is required"));
		const { event } = resendEvent("bad-email");

		const result = await actions.resend(event as never);

		expect(result).toEqual({
			status: 400,
			data: {
				values: { email: "bad-email" },
				error: "valid email is required"
			}
		});
		expect(sendVerificationEmailMock).not.toHaveBeenCalled();
		expect(loggerMock.info).not.toHaveBeenCalled();
	});

	it("returns generic resend form errors for auth service failures", async () => {
		sendVerificationOTPMock.mockRejectedValue(new AuthApiErrorMock(503, "database connection failed"));
		const { event } = resendEvent();

		const result = await actions.resend(event as never);

		expect(result).toEqual({
			status: 502,
			data: {
				values: { email: "ada@example.com" },
				error: "Unable to send verification code. Please try again."
			}
		});
		expect(sendVerificationEmailMock).not.toHaveBeenCalled();
		expect(loggerMock.info).not.toHaveBeenCalled();
	});
});
