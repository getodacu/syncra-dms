import { fail, redirect } from "@sveltejs/kit";
import { isAuthApiError, setPreferredLanguageCookie, setSessionCookie, signInEmail } from "$lib/server/auth";
import { publicErrorMessage, publicErrorStatus } from "$lib/server/public-errors";
import type { Actions, PageServerLoad } from "./$types";

function textValue(data: FormData, key: string) {
	const value = data.get(key);
	return typeof value === "string" ? value : "";
}

function oauthErrorMessage(value: string | null) {
	switch (value) {
		case "configuration":
			return "Social login is not configured yet.";
		case "denied":
			return "Social login was cancelled.";
		case "invalid":
			return "Social login could not be verified. Please try again.";
		case "failed":
			return "Social login failed. Please try again.";
		default:
			return "";
	}
}

export const load: PageServerLoad = ({ locals, url }) => {
	if (locals.user) redirect(303, "/app");

	return {
		email: url.searchParams.get("email")?.trim().toLowerCase() ?? "",
		verified: url.searchParams.get("verified") === "1",
		reset: url.searchParams.get("reset") === "1",
		oauthError: oauthErrorMessage(url.searchParams.get("oauth_error"))
	};
};

export const actions = {
	default: async ({ cookies, fetch, locals, request }) => {
		const data = await request.formData();
		const email = textValue(data, "email").trim().toLowerCase();
		const password = textValue(data, "password");
		const rememberMe = false;
		const values = { email };
		const fieldErrors: Record<string, string> = {};

		if (!email) fieldErrors.email = "Email is required.";
		if (!password) fieldErrors.password = "Password is required.";

		if (Object.keys(fieldErrors).length > 0) {
			return fail(400, { values, fieldErrors });
		}

		try {
			const auth = await signInEmail(fetch, { email, password, rememberMe });
			setSessionCookie(cookies, auth.session, rememberMe);
			setPreferredLanguageCookie(cookies, auth.user.preferredLanguage);
			locals.session = auth.session;
			locals.user = auth.user;
		} catch (error) {
			if (isAuthApiError(error)) {
				const message =
					error.status === 403
						? "Verify your email before signing in."
						: error.status === 401
							? "Invalid email or password."
							: publicErrorMessage(error.status, error.message, "Unable to sign in. Please try again.");
				return fail(publicErrorStatus(error.status), {
					values,
					error: message
				});
			}
			throw error;
		}

		redirect(303, "/app");
	}
} satisfies Actions;
