import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const accountSettingsDialogSource = () =>
	readFileSync(new URL("./account-settings-dialog.svelte", import.meta.url), "utf8");
const normalizeSource = (value: string) => value.replace(/\s+/g, " ").trim();

describe("account settings dialog source", () => {
	it("enables sessions and linked accounts sections", () => {
		const source = normalizeSource(accountSettingsDialogSource());

		expect(source).toContain('{ id: "sessions", label: m.account_settings_sessions() }');
		expect(source).toContain('{ id: "linked", label: m.account_settings_linked_accounts() }');
		expect(source).not.toContain('{ id: "sessions", label: "Sessions", disabled: true }');
		expect(source).not.toContain('{ id: "linked", label: "Linked Accounts", disabled: true }');
		expect(source).toContain("initialSection?: SettingsSection;");
		expect(source).toContain("activeSection = requestedSection;");
	});

	it("loads and revokes sessions through the account settings API", () => {
		const source = normalizeSource(accountSettingsDialogSource());

		expect(source).toContain('accountSettingsRequest<{ sessions: AuthSessionListItem[] }>( "/api/auth/sessions"');
		expect(source).toContain("`/api/auth/sessions/${encodeURIComponent(session.id)}`");
		expect(source).toContain("confirmSessionRevoke(session)");
		expect(source).toContain("disabled={session.current || revokingSessionId === session.id}");
		expect(source).toContain(
			"aria-label={session.current ? m.account_settings_current_session_cannot_revoke() : m.account_settings_revoke_session()}"
		);
		expect(source).toContain("toast.success(m.account_settings_session_revoked())");
	});

	it("lists, links, and unlinks OAuth accounts through the new routes", () => {
		const source = normalizeSource(accountSettingsDialogSource());

		expect(source).toContain('accountSettingsRequest<{ accounts: AuthAccountListItem[] }>( "/api/auth/accounts"');
		expect(source).toContain("`/api/auth/accounts/${providerId}`");
		expect(source).toContain('connectHref: "/api/auth/link/google"');
		expect(source).toContain('connectHref: "/api/auth/link/github"');
		expect(source).toContain("signInMethodCount <= 1");
		expect(source).toContain("confirmProviderUnlink(provider)");
		expect(source).toContain("toast.success(m.account_settings_linked_account_removed())");
	});

	it("adds a saved language preference to general settings", () => {
		const source = normalizeSource(accountSettingsDialogSource());

		expect(source).toContain('preferredLanguage?: "en" | "ro";');
		expect(source).toContain('languageValue = source?.preferredLanguage ?? "en";');
		expect(source).toContain("await patchUser({ preferredLanguage: languageValue });");
		expect(source).toContain("toast.success(m.account_settings_language_saved())");
		expect(source).toContain('<Select.Item value="en"');
		expect(source).toContain('<Select.Item value="ro"');
		expect(source).toContain("void saveLanguage();");
	});

	it("uses Paraglide messages for account settings labels and feedback", () => {
		const source = accountSettingsDialogSource();

		expect(source).toContain('import { m } from "$lib/paraglide/messages.js";');
		for (const messageCall of [
			"m.account_settings_title()",
			"m.account_settings_description()",
			"m.account_settings_avatar()",
			"m.account_settings_save_name()",
			"m.account_settings_save_email()",
			"m.account_settings_save_language()",
			"m.account_settings_save_password()",
			"m.account_settings_sessions_description()",
			"m.account_settings_linked_accounts_description()",
			"m.account_settings_unavailable_title()"
		]) {
			expect(source).toContain(messageCall);
		}
	});
});
