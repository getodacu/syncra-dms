import { beforeEach, describe, expect, it, vi } from "vitest";
import type { SendMailOptions } from "nodemailer";

const { createTransportMock, privateEnv, sendMailMock } = vi.hoisted(() => ({
	createTransportMock: vi.fn(),
	privateEnv: {} as Record<string, string | undefined>,
	sendMailMock: vi.fn()
}));

vi.mock("$env/dynamic/private", () => ({ env: privateEnv }));
vi.mock("nodemailer", () => ({
	default: { createTransport: createTransportMock },
	createTransport: createTransportMock
}));

function setSMTPEnv(overrides: Record<string, string | undefined> = {}) {
	Object.assign(privateEnv, {
		MAIL_SMTP_HOST: "smtp.example.com",
		MAIL_SMTP_PORT: "587",
		MAIL_SMTP_USER: "mailer@example.com",
		MAIL_SMTP_PASSWORD: "smtp-password",
		MAIL_SMTP_FROM: "Syncra <no-reply@example.com>",
		MAIL_SMTP_TLS: "true",
		...overrides
	});
}

describe("verification email mailer", () => {
	beforeEach(() => {
		vi.resetModules();
		createTransportMock.mockReset();
		sendMailMock.mockReset();
		createTransportMock.mockReturnValue({ sendMail: sendMailMock });
		sendMailMock.mockResolvedValue({});
		for (const key of Object.keys(privateEnv)) delete privateEnv[key];
	});

	it("sends verification emails through configured SMTP with STARTTLS", async () => {
		setSMTPEnv();
		const { sendVerificationEmail } = await import("./mail");

		await sendVerificationEmail({
			to: "ada@example.com",
			name: "Ada <Admin>",
			code: "123456",
			expiresAt: "2026-06-02T13:30:00Z"
		});

		expect(createTransportMock).toHaveBeenCalledWith({
			host: "smtp.example.com",
			port: 587,
			secure: false,
			requireTLS: true,
			auth: {
				user: "mailer@example.com",
				pass: "smtp-password"
			}
		});
		const message = sendMailMock.mock.calls[0][0] as SendMailOptions;
		expect(message.to).toBe("ada@example.com");
		expect(message.from).toBe("Syncra <no-reply@example.com>");
		expect(message.subject).toBe("Your Syncra verification code");
		expect(message.text).toContain("123456");
		expect(message.text).toContain("Hi Ada <Admin>,");
		expect(message.html).toContain("123456");
		expect(message.html).toContain("Hi Ada &lt;Admin&gt;,");
	});

	it("uses implicit TLS for port 465", async () => {
		setSMTPEnv({ MAIL_SMTP_PORT: "465", MAIL_SMTP_TLS: "false" });
		const { sendVerificationEmail } = await import("./mail");

		await sendVerificationEmail({ to: "ada@example.com", code: "123456" });

		expect(createTransportMock).toHaveBeenCalledWith(
			expect.objectContaining({
				port: 465,
				secure: true,
				requireTLS: false
			})
		);
	});

	it("sends password reset emails with the reset link", async () => {
		setSMTPEnv();
		const { sendPasswordResetEmail } = await import("./mail");

		await sendPasswordResetEmail({
			to: "ada@example.com",
			resetUrl: "https://app.example.com/recover-password?email=ada%40example.com&token=abc123",
			expiresAt: "2026-06-11T13:30:00Z"
		});

		const message = sendMailMock.mock.calls[0][0] as SendMailOptions;
		expect(message.to).toBe("ada@example.com");
		expect(message.from).toBe("Syncra <no-reply@example.com>");
		expect(message.subject).toBe("Reset your Syncra password");
		expect(message.text).toContain("https://app.example.com/recover-password?email=ada%40example.com&token=abc123");
		expect(message.html).toContain("https://app.example.com/recover-password?email=ada%40example.com&amp;token=abc123");
		expect(message.text).toContain("Jun 11, 2026");
	});

	it("sends paid invoice emails with the invoice PDF attached", async () => {
		setSMTPEnv();
		const { sendInvoicePaidEmail } = await import("./mail");

		await sendInvoicePaidEmail({
			to: "billing@example.com",
			billingName: "Ada <Buyer>",
			invoiceLabel: "SYN-00042",
			invoiceDate: "2026-06-13",
			totalAmount: "47.50",
			currency: "eur",
			credits: 5000,
			orderId: "order-1",
			pdf: {
				filename: "SYN_00042_260613.pdf",
				content: Buffer.from("%PDF-test"),
				contentType: "application/pdf"
			}
		});

		const message = sendMailMock.mock.calls[0][0] as SendMailOptions;
		expect(message.to).toBe("billing@example.com");
		expect(message.from).toBe("Syncra <no-reply@example.com>");
		expect(message.subject).toBe("Your Syncra invoice SYN-00042");
		expect(message.text).toContain("Hi Ada <Buyer>,");
		expect(message.text).toContain("5000 Syncra credits");
		expect(message.text).toContain("47.50 EUR");
		expect(message.html).toContain("Ada &lt;Buyer&gt;");
		expect(message.html).toContain("SYN-00042");
		expect(message.html).toContain("47.50 EUR");
		expect(message.attachments).toEqual([
			{
				filename: "SYN_00042_260613.pdf",
				content: Buffer.from("%PDF-test"),
				contentType: "application/pdf"
			}
		]);
	});

	it("throws the generic delivery error for missing SMTP config", async () => {
		setSMTPEnv({ MAIL_SMTP_HOST: "" });
		const { sendVerificationEmail, VERIFICATION_EMAIL_DELIVERY_ERROR } = await import("./mail");

		await expect(
			sendVerificationEmail({ to: "ada@example.com", code: "123456" })
		).rejects.toThrow(VERIFICATION_EMAIL_DELIVERY_ERROR);
		expect(createTransportMock).not.toHaveBeenCalled();
	});

	it("throws the generic delivery error for invalid SMTP config", async () => {
		setSMTPEnv({ MAIL_SMTP_PORT: "not-a-port" });
		const { sendVerificationEmail, VERIFICATION_EMAIL_DELIVERY_ERROR } = await import("./mail");

		await expect(
			sendVerificationEmail({ to: "ada@example.com", code: "123456" })
		).rejects.toThrow(VERIFICATION_EMAIL_DELIVERY_ERROR);
		expect(createTransportMock).not.toHaveBeenCalled();
	});

	it("throws the generic delivery error when SMTP delivery fails", async () => {
		setSMTPEnv();
		sendMailMock.mockRejectedValue(new Error("smtp down"));
		const { sendVerificationEmail, VERIFICATION_EMAIL_DELIVERY_ERROR } = await import("./mail");

		await expect(
			sendVerificationEmail({ to: "ada@example.com", code: "123456" })
		).rejects.toThrow(VERIFICATION_EMAIL_DELIVERY_ERROR);
	});
});
