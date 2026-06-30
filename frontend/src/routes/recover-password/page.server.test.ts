import { describe, expect, it, vi } from 'vitest';

const { requestPasswordResetMock, resetPasswordMock, sendPasswordResetEmailMock } = vi.hoisted(() => ({
	requestPasswordResetMock: vi.fn(),
	resetPasswordMock: vi.fn(),
	sendPasswordResetEmailMock: vi.fn()
}));

vi.mock('$lib/server/auth', () => ({
	isAuthApiError: (error: unknown) => error instanceof Error && 'status' in error,
	requestPasswordReset: requestPasswordResetMock,
	resetPassword: resetPasswordMock
}));

vi.mock('$lib/server/mail', () => ({
	sendPasswordResetEmail: sendPasswordResetEmailMock
}));

import { actions } from './+page.server';

describe('recover password actions', () => {
	it('sends password reset emails using the trusted backend token', async () => {
		vi.stubEnv('SYNCRA_APP_ORIGIN', 'http://localhost');
		requestPasswordResetMock.mockResolvedValue({
			resetToken: 'reset-token',
			resetExpiresAt: '2026-06-30T12:00:00Z'
		});
		sendPasswordResetEmailMock.mockResolvedValue(undefined);
		const data = new FormData();
		data.set('email', 'ADA@EXAMPLE.COM');

		const result = await actions.request({
			request: new Request('http://localhost/recover-password', { method: 'POST', body: data }),
			fetch: vi.fn(),
			url: new URL('http://localhost/recover-password'),
			locals: { logger: { error: vi.fn() } }
		} as never);

		expect(sendPasswordResetEmailMock).toHaveBeenCalledWith(
			expect.objectContaining({
				to: 'ada@example.com',
				resetUrl: 'http://localhost/recover-password?email=ada%40example.com&token=reset-token',
				expiresAt: '2026-06-30T12:00:00Z'
			})
		);
		expect(JSON.stringify(result)).not.toContain('reset-token');
	});

	it('resets a password and redirects to login', async () => {
		resetPasswordMock.mockResolvedValue({ ok: true });
		const data = new FormData();
		data.set('email', 'ADA@EXAMPLE.COM');
		data.set('token', 'reset-token');
		data.set('password', 'new-password');
		data.set('confirmPassword', 'new-password');

		await expect(
			actions.reset({
				request: new Request('http://localhost/recover-password', { method: 'POST', body: data }),
				fetch: vi.fn()
			} as never)
		).rejects.toMatchObject({
			status: 303,
			location: '/login?email=ada%40example.com&reset=1'
		});
	});
});
