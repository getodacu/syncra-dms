import { env } from "$env/dynamic/private";
import invoicePaidTemplate from "$lib/assets/mail-templates/invoice-paid.html?raw";
import nodemailer from "nodemailer";
import type SMTPTransport from "nodemailer/lib/smtp-transport";
import type { SendMailOptions } from "nodemailer";
import type { Buffer } from "node:buffer";

export const VERIFICATION_EMAIL_DELIVERY_ERROR =
	"We couldn't send the verification email. Please try again.";
export const INVOICE_EMAIL_DELIVERY_ERROR = "We couldn't send the invoice email.";

const VERIFICATION_EMAIL_SUBJECT = "Your Syncra verification code";
const PASSWORD_RESET_EMAIL_SUBJECT = "Reset your Syncra password";
const INVOICE_PAID_EMAIL_SUBJECT_PREFIX = "Your Syncra invoice";

type VerificationEmailInput = {
	to: string;
	name?: string;
	code: string;
	expiresAt?: string;
};

type PasswordResetEmailInput = {
	to: string;
	resetUrl: string;
	expiresAt?: string;
};

type InvoicePaidEmailInput = {
	to: string;
	billingName: string;
	invoiceLabel: string;
	invoiceDate: string;
	totalAmount: string;
	currency: string;
	credits: number;
	orderId: string;
	pdf: {
		filename: string;
		content: Buffer;
		contentType?: string;
	};
};

type MailConfig = {
	host: string;
	port: number;
	user: string;
	password: string;
	from: string;
	tls: boolean;
};

export async function sendVerificationEmail(input: VerificationEmailInput) {
	try {
		const config = mailConfig();
		const transport = nodemailer.createTransport(smtpTransportOptions(config));
		const message = verificationEmailMessage(config, input);

		await transport.sendMail(message);
	} catch (cause) {
		throw deliveryError(VERIFICATION_EMAIL_DELIVERY_ERROR, cause);
	}
}

export async function sendPasswordResetEmail(input: PasswordResetEmailInput) {
	try {
		const config = mailConfig();
		const transport = nodemailer.createTransport(smtpTransportOptions(config));
		const message = passwordResetEmailMessage(config, input);

		await transport.sendMail(message);
	} catch (cause) {
		throw deliveryError(VERIFICATION_EMAIL_DELIVERY_ERROR, cause);
	}
}

export async function sendInvoicePaidEmail(input: InvoicePaidEmailInput) {
	try {
		const config = mailConfig();
		const transport = nodemailer.createTransport(smtpTransportOptions(config));
		const message = invoicePaidEmailMessage(config, input);

		await transport.sendMail(message);
	} catch (cause) {
		throw deliveryError(INVOICE_EMAIL_DELIVERY_ERROR, cause);
	}
}

function mailConfig(): MailConfig {
	const host = requiredMailEnv("MAIL_SMTP_HOST");
	const port = mailPort(requiredMailEnv("MAIL_SMTP_PORT"));
	const user = requiredMailEnv("MAIL_SMTP_USER");
	const password = requiredMailEnv("MAIL_SMTP_PASSWORD");
	const from = requiredMailEnv("MAIL_SMTP_FROM");

	return {
		host,
		port,
		user,
		password,
		from,
		tls: parseBooleanEnv(privateEnv("MAIL_SMTP_TLS"))
	};
}

function smtpTransportOptions(config: MailConfig): SMTPTransport.Options {
	const secure = config.port === 465;

	return {
		host: config.host,
		port: config.port,
		secure,
		requireTLS: config.tls && !secure,
		auth: {
			user: config.user,
			pass: config.password
		}
	};
}

function verificationEmailMessage(
	config: MailConfig,
	input: VerificationEmailInput
): SendMailOptions {
	const greetingName = input.name?.trim();
	const greeting = greetingName ? `Hi ${greetingName},` : "Hi,";
	const expiryText = verificationExpiryText(input.expiresAt);
	const escapedGreeting = escapeHTML(greeting);
	const escapedCode = escapeHTML(input.code);
	const escapedExpiry = escapeHTML(expiryText);

	return {
		to: input.to,
		from: config.from,
		subject: VERIFICATION_EMAIL_SUBJECT,
		text: [
			greeting,
			"",
			`Your Syncra verification code is ${input.code}.`,
			expiryText,
			"",
			"If you did not request this code, you can safely ignore this email."
		].join("\n"),
		html: [
			"<!doctype html>",
			'<html lang="en">',
			"<body>",
			`<p>${escapedGreeting}</p>`,
			"<p>Your Syncra verification code is:</p>",
			`<p><strong style="font-size: 24px; letter-spacing: 4px;">${escapedCode}</strong></p>`,
			`<p>${escapedExpiry}</p>`,
			"<p>If you did not request this code, you can safely ignore this email.</p>",
			"</body>",
			"</html>"
		].join("")
	};
}

function passwordResetEmailMessage(
	config: MailConfig,
	input: PasswordResetEmailInput
): SendMailOptions {
	const expiryText = verificationExpiryText(input.expiresAt);
	const escapedResetUrl = escapeHTML(input.resetUrl);
	const escapedExpiry = escapeHTML(expiryText);

	return {
		to: input.to,
		from: config.from,
		subject: PASSWORD_RESET_EMAIL_SUBJECT,
		text: [
			"Hi,",
			"",
			"Use this link to reset your Syncra password:",
			input.resetUrl,
			expiryText,
			"",
			"If you did not request this reset, you can safely ignore this email."
		].join("\n"),
		html: [
			"<!doctype html>",
			'<html lang="en">',
			"<body>",
			"<p>Hi,</p>",
			"<p>Use this link to reset your Syncra password:</p>",
			`<p><a href="${escapedResetUrl}">Reset your password</a></p>`,
			`<p>${escapedExpiry}</p>`,
			"<p>If you did not request this reset, you can safely ignore this email.</p>",
			"</body>",
			"</html>"
		].join("")
	};
}

function invoicePaidEmailMessage(
	config: MailConfig,
	input: InvoicePaidEmailInput
): SendMailOptions {
	const billingName = input.billingName.trim() || "there";
	const currency = input.currency.trim().toUpperCase();
	const amountText = `${input.totalAmount} ${currency}`;
	const html = renderInvoicePaidTemplate({
		billingName,
		invoiceLabel: input.invoiceLabel,
		invoiceDate: input.invoiceDate,
		totalAmount: input.totalAmount,
		currency,
		credits: String(input.credits),
		orderId: input.orderId
	});

	return {
		to: input.to,
		from: config.from,
		subject: `${INVOICE_PAID_EMAIL_SUBJECT_PREFIX} ${input.invoiceLabel}`,
		text: [
			`Hi ${billingName},`,
			"",
			`We received your payment for ${input.credits} Syncra credits.`,
			`Invoice ${input.invoiceLabel}, dated ${input.invoiceDate}, for ${amountText} is attached as a PDF.`,
			`Order ID: ${input.orderId}`,
			"",
			"Thank you for using Syncra."
		].join("\n"),
		html,
		attachments: [
			{
				filename: input.pdf.filename,
				content: input.pdf.content,
				contentType: input.pdf.contentType || "application/pdf"
			}
		]
	};
}

function renderInvoicePaidTemplate(values: Record<string, string>) {
	return invoicePaidTemplate.replace(/\{\{(\w+)\}\}/g, (_match, key: string) =>
		escapeHTML(values[key] ?? "")
	);
}

function verificationExpiryText(expiresAt?: string) {
	if (!expiresAt) {
		return "This code expires soon.";
	}

	const date = new Date(expiresAt);
	if (Number.isNaN(date.getTime())) {
		return "This code expires soon.";
	}

	return `This code expires at ${date.toLocaleString("en-US", {
		dateStyle: "medium",
		timeStyle: "short",
		timeZone: "UTC"
	})} UTC.`;
}

function requiredMailEnv(key: string) {
	const value = privateEnv(key)?.trim();
	if (!value) throw deliveryError(VERIFICATION_EMAIL_DELIVERY_ERROR);
	return value;
}

function mailPort(value: string) {
	const port = Number(value);
	if (!Number.isInteger(port) || port < 1 || port > 65535) {
		throw deliveryError(VERIFICATION_EMAIL_DELIVERY_ERROR);
	}
	return port;
}

function parseBooleanEnv(value: string | undefined) {
	return value?.trim().toLowerCase() === "true";
}

function privateEnv(key: string) {
	return env[key] || nodeEnv()[key];
}

function nodeEnv() {
	return (
		globalThis as typeof globalThis & {
			process?: { env?: Record<string, string | undefined> };
		}
	).process?.env ?? {};
}

function escapeHTML(value: string) {
	return value.replace(/[&<>"']/g, (character) => {
		switch (character) {
			case "&":
				return "&amp;";
			case "<":
				return "&lt;";
			case ">":
				return "&gt;";
			case '"':
				return "&quot;";
			case "'":
				return "&#39;";
			default:
				return character;
		}
	});
}

function deliveryError(message: string, cause?: unknown) {
	return new Error(message, { cause });
}
