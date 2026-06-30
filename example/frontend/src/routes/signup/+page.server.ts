import { fail, redirect } from "@sveltejs/kit";
import { isAuthApiError, signUpEmail } from "$lib/server/auth";
import { sendVerificationEmail, VERIFICATION_EMAIL_DELIVERY_ERROR } from "$lib/server/mail";
import { publicErrorMessage, publicErrorStatus } from "$lib/server/public-errors";
import type { Actions, PageServerLoad } from "./$types";

function textValue(data: FormData, key: string) {
	const value = data.get(key);
	return typeof value === "string" ? value : "";
}

export const load: PageServerLoad = ({ locals }) => {
	if (locals.user) redirect(303, "/app");
};

export const actions = {
	default: async ({ fetch, locals, request }) => {
		const data = await request.formData();
		const name = textValue(data, "name").trim();
		const email = textValue(data, "email").trim().toLowerCase();
		const password = textValue(data, "password");
		const confirmPassword = textValue(data, "confirmPassword");
		const values = { name, email };
		const fieldErrors: Record<string, string> = {};

		if (!name) fieldErrors.name = "Name is required.";
		if (!email) fieldErrors.email = "Email is required.";
		if (password.length < 8) fieldErrors.password = "Password must be at least 8 characters.";
		if (password !== confirmPassword) fieldErrors.confirmPassword = "Passwords do not match.";

		if (Object.keys(fieldErrors).length > 0) {
			return fail(400, { values, fieldErrors });
		}

		let verificationCode: string | undefined;
		let verificationExpiresAt: string | undefined;
		try {
			const result = await signUpEmail(fetch, { name, email, password });
			verificationCode = result.verificationCode;
			verificationExpiresAt = result.verificationExpiresAt;
		} catch (error) {
			if (isAuthApiError(error)) {
				return fail(publicErrorStatus(error.status), {
					values,
					error: publicErrorMessage(
						error.status,
						error.message,
						"Unable to create account. Please try again."
					)
				});
			}
			throw error;
		}

		if (verificationCode) {
			locals.logger.info("auth.verification_email_code_sent", { email });
			try {
				await sendVerificationEmail({
					to: email,
					name,
					code: verificationCode,
					expiresAt: verificationExpiresAt
				});
			} catch {
				return fail(502, { values, error: VERIFICATION_EMAIL_DELIVERY_ERROR });
			}
		}

		const params = new URLSearchParams({ email });
		redirect(303, `/signup-confirmation?${params}`);
	}
} satisfies Actions;
