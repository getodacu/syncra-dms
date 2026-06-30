import { describe, expect, it, vi } from 'vitest';

const { sendVerificationEmailMock, sendVerificationOTPMock, verifyEmailOTPMock } = vi.hoisted(() => ({
	sendVerificationEmailMock: vi.fn(),
	sendVerificationOTPMock: vi.fn(),
	verifyEmailOTPMock: vi.fn()
}));

vi.mock('$lib/server/auth', () => ({
	isAuthApiError: (error: unknown) => error instanceof Error && 'status' in error,
	sendVerificationOTP: sendVerificationOTPMock,
	verifyEmailOTP: verifyEmailOTPMock
}));

vi.mock('$lib/server/mail', () => ({
	sendVerificationEmail: sendVerificationEmailMock,
	VERIFICATION_EMAIL_DELIVERY_ERROR: "We couldn't send the verification email. Please try again."
}));

import { actions } from './+page.server';

describe('signup confirmation actions', () => {
	it('verifies a 6-digit OTP and redirects to login', async () => {
		verifyEmailOTPMock.mockResolvedValue({ ok: true });
		const data = new FormData();
		data.set('email', 'ADA@EXAMPLE.COM');
		data.set('otp', '123456');

		await expect(
			actions.verify({
				request: new Request('http://localhost/signup-confirmation', { method: 'POST', body: data }),
				fetch: vi.fn()
			} as never)
		).rejects.toMatchObject({
			status: 303,
			location: '/login?email=ada%40example.com&verified=1'
		});
	});

	it('resends OTP email without exposing the code in the form result', async () => {
		sendVerificationOTPMock.mockResolvedValue({
			ok: true,
			verificationCode: '654321',
			verificationExpiresAt: '2026-06-30T12:00:00Z'
		});
		sendVerificationEmailMock.mockResolvedValue(undefined);
		const data = new FormData();
		data.set('email', 'ADA@EXAMPLE.COM');

		const result = await actions.resend({
			request: new Request('http://localhost/signup-confirmation', { method: 'POST', body: data }),
			fetch: vi.fn(),
			locals: { logger: { info: vi.fn() } }
		} as never);

		expect(sendVerificationEmailMock).toHaveBeenCalledWith({
			to: 'ada@example.com',
			code: '654321',
			expiresAt: '2026-06-30T12:00:00Z'
		});
		expect(JSON.stringify(result)).not.toContain('654321');
	});
});
