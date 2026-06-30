import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const navUserSource = () => readFileSync(new URL("./nav-user.svelte", import.meta.url), "utf8");
const normalizeSource = (value: string) => value.replace(/\s+/g, " ").trim();

describe("user navigation behavior", () => {
	it("wires logout through a confirmation dialog and POST form", () => {
		const source = normalizeSource(navUserSource());

		expect(source).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(source).toContain('import * as AlertDialog from "$lib/components/ui/alert-dialog/index.js";');
		expect(source).toContain("let logoutConfirmationOpen = $state(false);");
		expect(source).toContain("let logoutPending = $state(false);");
		expect(source).toContain("async function submitLogout(event: SubmitEvent) {");
		expect(source).toContain("event.preventDefault();");
		expect(source).toContain("if (logoutPending) return;");
		expect(source).toContain("const form = event.currentTarget as HTMLFormElement;");
		expect(source).toContain("const response = await fetch(form.action, { method: form.method });");
		expect(source).toContain("await goto(logoutRedirectPath(response), { invalidateAll: true });");
		expect(source).toContain("toast.error(m.nav_logout_failed());");
		expect(source).toContain("function logoutRedirectPath(response: Response) {");
		expect(source).toContain(
			'<DropdownMenu.Item variant="destructive" onclick={() => (logoutConfirmationOpen = true)}>'
		);
		expect(source).toContain("<AlertDialog.Root bind:open={logoutConfirmationOpen}>");
		expect(source).toContain('<form action="/logout" method="POST" class="flex flex-col gap-4" onsubmit={submitLogout}>');
		expect(source).toContain("<AlertDialog.Title>{m.nav_logout_title()}</AlertDialog.Title>");
		expect(source).toContain("{m.nav_logout_description()}");
		expect(source).toContain('<AlertDialog.Cancel type="button" disabled={logoutPending}>{m.common_cancel()}</AlertDialog.Cancel>');
		expect(source).toContain('<AlertDialog.Action type="submit" variant="destructive" loading={logoutPending} disabled={logoutPending}>');
	});

	it("uses Paraglide messages for account menu labels and feedback", () => {
		const source = normalizeSource(navUserSource());

		for (const messageCall of [
			"m.nav_account()",
			"m.nav_no_email_address()",
			"m.nav_billing()",
			"m.nav_notifications()",
			"m.nav_log_out()",
			"m.nav_account_linked({ provider: providerLabel })",
			"m.nav_account_link_conflict({ provider: providerLabel })",
			"m.nav_account_link_denied({ provider: providerLabel })",
			"m.nav_account_link_not_configured({ provider: providerLabel })",
			"m.nav_account_link_sign_in_again()",
			"m.nav_account_link_failed({ provider: providerLabel })"
		]) {
			expect(source).toContain(messageCall);
		}
	});

	it("opens the billing profile dialog from the Billing menu item", () => {
		const source = normalizeSource(navUserSource());

		expect(source).toContain('import BillingProfileDialog from "./billing-profile-dialog.svelte";');
		expect(source).toContain("let billingProfileOpen = $state(false);");
		expect(source).toContain('<DropdownMenu.Item onclick={() => (billingProfileOpen = true)}>');
		expect(source).toContain("<BillingProfileDialog bind:open={billingProfileOpen} {user} />");
	});

	it("opens account settings from OAuth link return parameters and cleans the URL", () => {
		const source = normalizeSource(navUserSource());

		expect(source).toContain('import { goto } from "$app/navigation";');
		expect(source).toContain('import { page } from "$app/state";');
		expect(source).toContain('import { toast } from "svelte-sonner";');
		expect(source).toContain('let accountSettingsInitialSection = $state<SettingsSection>("general");');
		expect(source).toContain('url.searchParams.get("account_settings")');
		expect(source).toContain('url.searchParams.get("account_link_provider")');
		expect(source).toContain('url.searchParams.get("account_link_status")');
		expect(source).toContain('accountSettingsInitialSection = settingsSection;');
		expect(source).toContain("replaceState: true");
		expect(source).toContain("initialSection={accountSettingsInitialSection}");
		expect(source).toContain("toast.success(m.nav_account_linked({ provider: providerLabel }));");
	});
});
