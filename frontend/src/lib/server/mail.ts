import nodemailer from 'nodemailer';
import type { SendMailOptions } from 'nodemailer';
import type SMTPTransport from 'nodemailer/lib/smtp-transport';
import { privateEnv } from './internal-api';

export const VERIFICATION_EMAIL_DELIVERY_ERROR =
	"We couldn't send the verification email. Please try again.";

type MailConfig = {
	host: string;
	port: number;
	user: string;
	password: string;
	from: string;
	tls: boolean;
};

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

export async function sendVerificationEmail(input: VerificationEmailInput) {
	const config = mailConfig();
	const transport = nodemailer.createTransport(smtpTransportOptions(config));
	await transport.sendMail(verificationEmailMessage(config, input));
}

export async function sendPasswordResetEmail(input: PasswordResetEmailInput) {
	const config = mailConfig();
	const transport = nodemailer.createTransport(smtpTransportOptions(config));
	await transport.sendMail(passwordResetEmailMessage(config, input));
}

function mailConfig(): MailConfig {
	return {
		host: requiredMailEnv('MAIL_SMTP_HOST'),
		port: mailPort(requiredMailEnv('MAIL_SMTP_PORT')),
		user: requiredMailEnv('MAIL_SMTP_USER'),
		password: requiredMailEnv('MAIL_SMTP_PASSWORD'),
		from: requiredMailEnv('MAIL_SMTP_FROM'),
		tls: parseBooleanEnv(privateEnv('MAIL_SMTP_TLS'))
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

function verificationEmailMessage(config: MailConfig, input: VerificationEmailInput): SendMailOptions {
	const greetingName = input.name?.trim();
	const greeting = greetingName ? `Hi ${greetingName},` : 'Hi,';
	const expiryText = expiryTextFor(input.expiresAt);
	return {
		to: input.to,
		from: config.from,
		subject: 'Your Syncra DMS verification code',
		text: [
			greeting,
			'',
			`Your Syncra DMS verification code is ${input.code}.`,
			expiryText,
			'',
			'If you did not request this code, you can safely ignore this email.'
		].join('\n'),
		html: [
			'<!doctype html><html><body>',
			`<p>${escapeHTML(greeting)}</p>`,
			'<p>Your Syncra DMS verification code is:</p>',
			`<p><strong style="font-size:24px;letter-spacing:4px">${escapeHTML(input.code)}</strong></p>`,
			`<p>${escapeHTML(expiryText)}</p>`,
			'<p>If you did not request this code, you can safely ignore this email.</p>',
			'</body></html>'
		].join('')
	};
}

function passwordResetEmailMessage(config: MailConfig, input: PasswordResetEmailInput): SendMailOptions {
	const expiryText = expiryTextFor(input.expiresAt);
	return {
		to: input.to,
		from: config.from,
		subject: 'Reset your Syncra DMS password',
		text: [
			'Hi,',
			'',
			'Use this link to reset your Syncra DMS password:',
			input.resetUrl,
			expiryText,
			'',
			'If you did not request this reset, you can safely ignore this email.'
		].join('\n'),
		html: [
			'<!doctype html><html><body>',
			'<p>Hi,</p>',
			'<p>Use this link to reset your Syncra DMS password:</p>',
			`<p><a href="${escapeHTML(input.resetUrl)}">Reset your password</a></p>`,
			`<p>${escapeHTML(expiryText)}</p>`,
			'<p>If you did not request this reset, you can safely ignore this email.</p>',
			'</body></html>'
		].join('')
	};
}

function requiredMailEnv(key: string) {
	const value = privateEnv(key)?.trim();
	if (!value) throw new Error(`${key} is required`);
	return value;
}

function mailPort(value: string) {
	const port = Number(value);
	if (!Number.isInteger(port) || port <= 0) throw new Error('MAIL_SMTP_PORT must be a positive integer');
	return port;
}

function parseBooleanEnv(value: string | undefined) {
	return value === '1' || value?.toLowerCase() === 'true' || value?.toLowerCase() === 'yes';
}

function expiryTextFor(expiresAt?: string) {
	if (!expiresAt) return 'This code expires soon.';
	const date = new Date(expiresAt);
	if (Number.isNaN(date.getTime())) return 'This code expires soon.';
	return `This link or code expires at ${date.toLocaleString('en-US', {
		dateStyle: 'medium',
		timeStyle: 'short',
		timeZone: 'UTC'
	})} UTC.`;
}

function escapeHTML(value: string) {
	return value.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}
