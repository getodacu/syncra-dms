import { describe, expect, it, vi } from 'vitest';

const { sendVerificationEmailMock, signUpEmailMock } = vi.hoisted(() => ({
	sendVerificationEmailMock: vi.fn(),
	signUpEmailMock: vi.fn()
}));

vi.mock('$lib/server/auth', () => ({
	isAuthApiError: (error: unknown) => error instanceof Error && 'status' in error,
	signUpEmail: signUpEmailMock
}));

vi.mock('$lib/server/mail', () => ({
	sendVerificationEmail: sendVerificationEmailMock,
	VERIFICATION_EMAIL_DELIVERY_ERROR: "We couldn't send the verification email. Please try again."
}));

import { actions } from './+page.server';

describe('signup page action', () => {
	it('emails the trusted verification code and redirects to confirmation', async () => {
		signUpEmailMock.mockResolvedValue({
			verificationCode: '123456',
			verificationExpiresAt: '2026-06-30T12:00:00Z'
		});
		sendVerificationEmailMock.mockResolvedValue(undefined);
		const data = new FormData();
		data.set('name', 'Ada Lovelace');
		data.set('email', 'ADA@EXAMPLE.COM');
		data.set('password', 'password123');
		data.set('confirmPassword', 'password123');

		await expect(
			actions.default({
				request: new Request('http://localhost/signup', { method: 'POST', body: data }),
				fetch: vi.fn(),
				locals: { logger: { info: vi.fn() } }
			} as never)
		).rejects.toMatchObject({
			status: 303,
			location: '/signup-confirmation?email=ada%40example.com'
		});

		expect(sendVerificationEmailMock).toHaveBeenCalledWith({
			to: 'ada@example.com',
			name: 'Ada Lovelace',
			code: '123456',
			expiresAt: '2026-06-30T12:00:00Z'
		});
	});
});
