import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

const {
	AuthApiErrorMock,
	deliveryErrorMessage,
	loggerMock,
	sendVerificationEmailMock,
	signUpEmailMock
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
		signUpEmailMock: vi.fn()
	};
});

vi.mock("$lib/server/auth", () => ({
	AuthApiError: AuthApiErrorMock,
	isAuthApiError: (error: unknown) => error instanceof AuthApiErrorMock,
	signUpEmail: signUpEmailMock
}));

vi.mock("$lib/server/mail", () => ({
	sendVerificationEmail: sendVerificationEmailMock,
	VERIFICATION_EMAIL_DELIVERY_ERROR: deliveryErrorMessage
}));

import { actions } from "./+page.server";

function signupEvent(overrides: Record<string, string> = {}) {
	const data = new FormData();
	data.set("name", overrides.name ?? "Ada Lovelace");
	data.set("email", overrides.email ?? "ADA@EXAMPLE.COM");
	data.set("password", overrides.password ?? "password123");
	data.set("confirmPassword", overrides.confirmPassword ?? "password123");
	const fetchMock = vi.fn();

	return {
		event: {
			request: new Request("http://localhost/signup", {
				method: "POST",
				body: data
			}),
			fetch: fetchMock,
			locals: { logger: loggerMock }
		},
		fetchMock
	};
}

describe("signup page action", () => {
	beforeEach(() => {
		signUpEmailMock.mockReset();
		sendVerificationEmailMock.mockReset();
		sendVerificationEmailMock.mockResolvedValue(undefined);
		loggerMock.info.mockReset();
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	it("sends the backend verification code by email and redirects", async () => {
		signUpEmailMock.mockResolvedValue({
			verificationCode: "123456",
			verificationExpiresAt: "2026-06-02T13:30:00Z"
		});
		const { event, fetchMock } = signupEvent();

		await expect(actions.default(event as never)).rejects.toMatchObject({
			status: 303,
			location: "/signup-confirmation?email=ada%40example.com"
		});

		expect(signUpEmailMock).toHaveBeenCalledWith(fetchMock, {
			name: "Ada Lovelace",
			email: "ada@example.com",
			password: "password123"
		});
		expect(sendVerificationEmailMock).toHaveBeenCalledWith({
			to: "ada@example.com",
			name: "Ada Lovelace",
			code: "123456",
			expiresAt: "2026-06-02T13:30:00Z"
		});
		expect(loggerMock.info).toHaveBeenCalledWith("auth.verification_email_code_sent", {
			email: "ada@example.com"
		});
		expect(JSON.stringify(loggerMock.info.mock.calls)).not.toContain("123456");
	});

	it("returns a generic delivery error when SMTP send fails", async () => {
		signUpEmailMock.mockResolvedValue({
			verificationCode: "123456",
			verificationExpiresAt: "2026-06-02T13:30:00Z"
		});
		sendVerificationEmailMock.mockRejectedValue(new Error("smtp down"));
		const { event } = signupEvent();

		const result = await actions.default(event as never);

		expect(result).toEqual({
			status: 502,
			data: {
				values: { name: "Ada Lovelace", email: "ada@example.com" },
				error: deliveryErrorMessage
			}
		});
		expect(JSON.stringify(result)).not.toContain("123456");
		expect(loggerMock.info).toHaveBeenCalledWith("auth.verification_email_code_sent", {
			email: "ada@example.com"
		});
		expect(JSON.stringify(loggerMock.info.mock.calls)).not.toContain("123456");
	});

	it("does not send email when the backend returns no verification code", async () => {
		signUpEmailMock.mockResolvedValue({ verificationRequired: true });
		const { event } = signupEvent();

		await expect(actions.default(event as never)).rejects.toMatchObject({
			status: 303,
			location: "/signup-confirmation?email=ada%40example.com"
		});

		expect(sendVerificationEmailMock).not.toHaveBeenCalled();
		expect(loggerMock.info).not.toHaveBeenCalled();
	});

	it("maps backend auth errors without sending email", async () => {
		signUpEmailMock.mockRejectedValue(new AuthApiErrorMock(409, "email is already in use"));
		const { event } = signupEvent();

		const result = await actions.default(event as never);

		expect(result).toEqual({
			status: 409,
			data: {
				values: { name: "Ada Lovelace", email: "ada@example.com" },
				error: "email is already in use"
			}
		});
		expect(sendVerificationEmailMock).not.toHaveBeenCalled();
		expect(loggerMock.info).not.toHaveBeenCalled();
	});

	it("returns generic form errors for auth service failures", async () => {
		signUpEmailMock.mockRejectedValue(new AuthApiErrorMock(503, "database connection failed"));
		const { event } = signupEvent();

		const result = await actions.default(event as never);

		expect(result).toEqual({
			status: 502,
			data: {
				values: { name: "Ada Lovelace", email: "ada@example.com" },
				error: "Unable to create account. Please try again."
			}
		});
		expect(sendVerificationEmailMock).not.toHaveBeenCalled();
		expect(loggerMock.info).not.toHaveBeenCalled();
	});
});
