import { fail, redirect } from "@sveltejs/kit";
import { isAuthApiError, sendVerificationOTP, verifyEmailOTP } from "$lib/server/auth";
import { sendVerificationEmail, VERIFICATION_EMAIL_DELIVERY_ERROR } from "$lib/server/mail";
import { publicErrorMessage, publicErrorStatus } from "$lib/server/public-errors";
import type { Actions, PageServerLoad } from "./$types";

function textValue(data: FormData, key: string) {
	const value = data.get(key);
	return typeof value === "string" ? value : "";
}

export const load: PageServerLoad = ({ locals, url }) => {
	if (locals.user) redirect(303, "/app");

	return {
		email: url.searchParams.get("email")?.trim().toLowerCase() ?? ""
	};
};

export const actions = {
	verify: async ({ fetch, request }) => {
		const data = await request.formData();
		const email = textValue(data, "email").trim().toLowerCase();
		const otp = textValue(data, "otp").trim();
		const values = { email, otp };
		const fieldErrors: Record<string, string> = {};

		if (!email) fieldErrors.email = "Email is required.";
		if (!/^\d{6}$/.test(otp)) fieldErrors.otp = "Enter the 6-digit verification code.";

		if (Object.keys(fieldErrors).length > 0) {
			return fail(400, { values, fieldErrors });
		}

		try {
			await verifyEmailOTP(fetch, { email, otp });
		} catch (error) {
			if (isAuthApiError(error)) {
				return fail(publicErrorStatus(error.status), {
					values: { email, otp: "" },
					error: publicErrorMessage(
						error.status,
						error.message,
						"Unable to verify email. Please try again."
					)
				});
			}
			throw error;
		}

		const params = new URLSearchParams({ email, verified: "1" });
		redirect(303, `/login?${params}`);
	},
	resend: async ({ fetch, locals, request }) => {
		const data = await request.formData();
		const email = textValue(data, "email").trim().toLowerCase();

		if (!email) {
			return fail(400, {
				values: { email },
				fieldErrors: { email: "Email is required." }
			});
		}

		let verificationCode: string | undefined;
		let verificationExpiresAt: string | undefined;
		try {
			const result = await sendVerificationOTP(fetch, email);
			verificationCode = result.verificationCode;
			verificationExpiresAt = result.verificationExpiresAt;
		} catch (error) {
			if (isAuthApiError(error)) {
				return fail(publicErrorStatus(error.status), {
					values: { email },
					error: publicErrorMessage(
						error.status,
						error.message,
						"Unable to send verification code. Please try again."
					)
				});
			}
			throw error;
		}

		if (verificationCode) {
			locals.logger.info("auth.verification_email_code_resent", { email });
			try {
				await sendVerificationEmail({
					to: email,
					code: verificationCode,
					expiresAt: verificationExpiresAt
				});
			} catch {
				return fail(502, { values: { email }, error: VERIFICATION_EMAIL_DELIVERY_ERROR });
			}
		}

		return {
			values: { email },
			success: "A new verification code has been sent. Check your inbox and spam folder."
		};
	}
} satisfies Actions;
