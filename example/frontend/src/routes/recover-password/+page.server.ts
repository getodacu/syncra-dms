import { fail, redirect } from "@sveltejs/kit";
import { isAuthApiError, requestPasswordReset, resetPassword } from "$lib/server/auth";
import { privateEnv } from "$lib/server/internal-api";
import { safeError } from "$lib/server/logging";
import { sendPasswordResetEmail } from "$lib/server/mail";
import { publicErrorMessage, publicErrorStatus } from "$lib/server/public-errors";
import type { Actions, PageServerLoad } from "./$types";

const PASSWORD_RESET_REQUEST_SUCCESS = "If an account exists, we sent a reset link.";

function textValue(data: FormData, key: string) {
	const value = data.get(key);
	return typeof value === "string" ? value : "";
}

function appOrigin(url: URL) {
	return (privateEnv("SYNCRA_APP_ORIGIN") || url.origin).replace(/\/+$/, "");
}

function passwordResetURL(origin: string, email: string, token: string) {
	const url = new URL("/recover-password", origin);
	url.searchParams.set("email", email);
	url.searchParams.set("token", token);
	return url.toString();
}

export const load: PageServerLoad = ({ locals, url }) => {
	if (locals.user) redirect(303, "/app");

	return {
		email: url.searchParams.get("email")?.trim().toLowerCase() ?? "",
		token: url.searchParams.get("token")?.trim() ?? ""
	};
};

export const actions = {
	request: async ({ fetch, locals, request, url }) => {
		const data = await request.formData();
		const email = textValue(data, "email").trim().toLowerCase();
		const values = { email };

		if (!email) {
			return fail(400, {
				values,
				fieldErrors: { email: "Email is required." }
			});
		}

		let resetToken: string | undefined;
		let resetExpiresAt: string | undefined;
		try {
			const result = await requestPasswordReset(fetch, email);
			resetToken = result.resetToken;
			resetExpiresAt = result.resetExpiresAt;
		} catch (error) {
			if (isAuthApiError(error)) {
				return fail(publicErrorStatus(error.status), {
					values,
					error: publicErrorMessage(
						error.status,
						error.message,
						"Unable to request password reset. Please try again."
					)
				});
			}
			throw error;
		}

		if (resetToken) {
			try {
				await sendPasswordResetEmail({
					to: email,
					resetUrl: passwordResetURL(appOrigin(url), email, resetToken),
					expiresAt: resetExpiresAt
				});
			} catch (error) {
				locals.logger.error("auth.password_reset_email_failed", {
					email,
					error: safeError(error)
				});
			}
		}

		return {
			values,
			success: PASSWORD_RESET_REQUEST_SUCCESS
		};
	},
	reset: async ({ fetch, request }) => {
		const data = await request.formData();
		const email = textValue(data, "email").trim().toLowerCase();
		const token = textValue(data, "token").trim();
		const password = textValue(data, "password");
		const confirmPassword = textValue(data, "confirmPassword");
		const values = { email, token };
		const fieldErrors: Record<string, string> = {};

		if (!email) fieldErrors.email = "Email is required.";
		if (!token) fieldErrors.token = "Password reset token is required.";
		if (password.length < 8) fieldErrors.password = "Password must be at least 8 characters.";
		if (password !== confirmPassword) fieldErrors.confirmPassword = "Passwords do not match.";

		if (Object.keys(fieldErrors).length > 0) {
			return fail(400, { values, fieldErrors });
		}

		try {
			await resetPassword(fetch, { email, token, password });
		} catch (error) {
			if (isAuthApiError(error)) {
				return fail(publicErrorStatus(error.status), {
					values,
					error: publicErrorMessage(
						error.status,
						error.message,
						"Unable to reset password. Please try again."
					)
				});
			}
			throw error;
		}

		const params = new URLSearchParams({ email, reset: "1" });
		redirect(303, `/login?${params}`);
	}
} satisfies Actions;
