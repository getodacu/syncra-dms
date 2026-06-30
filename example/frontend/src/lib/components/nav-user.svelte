<script lang="ts">
	import CreditCardIcon from "@tabler/icons-svelte/icons/credit-card";
	import DotsVerticalIcon from "@tabler/icons-svelte/icons/dots-vertical";
	import LogoutIcon from "@tabler/icons-svelte/icons/logout";
	import NotificationIcon from "@tabler/icons-svelte/icons/notification";
	import UserCircleIcon from "@tabler/icons-svelte/icons/user-circle";
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { toast } from "svelte-sonner";
	import AccountSettingsDialog from "./account-settings-dialog.svelte";
	import BillingProfileDialog from "./billing-profile-dialog.svelte";
	import * as AlertDialog from "$lib/components/ui/alert-dialog/index.js";
	import * as Avatar from "$lib/components/ui/avatar/index.js";
	import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import { m } from "$lib/paraglide/messages.js";

	type AuthUser = NonNullable<App.Locals["user"]>;
	type SettingsSection = "general" | "security" | "sessions" | "linked" | "danger";

	let { user }: { user: AuthUser | null } = $props();

	const sidebar = Sidebar.useSidebar();

	let accountSettingsOpen = $state(false);
	let accountSettingsInitialSection = $state<SettingsSection>("general");
	let billingProfileOpen = $state(false);
	let logoutConfirmationOpen = $state(false);
	let logoutPending = $state(false);
	let handledAccountSettingsRequest = $state("");

	const displayName = $derived(user?.name?.trim() || m.nav_account());
	const displayEmail = $derived(user?.email?.trim() || m.nav_no_email_address());
	const avatarSrc = $derived(user?.image ?? null);
	const fallbackInitials = $derived(getInitials(displayName, user?.email));

	$effect(() => {
		const url = page.url;
		const settingsSection = url.searchParams.get("account_settings")?.trim() ?? "";
		const linkProvider = url.searchParams.get("account_link_provider")?.trim() ?? "";
		const linkStatus = url.searchParams.get("account_link_status")?.trim() ?? "";
		const requestKey = `${url.pathname}?${url.searchParams.toString()}`;

		if (!settingsSection && !linkStatus) return;
		if (requestKey === handledAccountSettingsRequest) return;
		handledAccountSettingsRequest = requestKey;

		if (isSettingsSection(settingsSection)) {
			accountSettingsInitialSection = settingsSection;
			accountSettingsOpen = true;
		}

		if (linkProvider && linkStatus) {
			showAccountLinkToast(linkProvider, linkStatus);
		}

		void goto(cleanAccountSettingsURL(url), {
			replaceState: true,
			keepFocus: true,
			noScroll: true
		});
	});

	function getInitials(name: string, email?: string | null) {
		const parts = name.trim().split(/\s+/).filter(Boolean);
		if (parts.length >= 2) return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
		if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
		return (email?.trim()[0] ?? "A").toUpperCase();
	}

	function isSettingsSection(value: string): value is SettingsSection {
		return ["general", "security", "sessions", "linked", "danger"].includes(value);
	}

	function cleanAccountSettingsURL(url: URL) {
		const params = new URLSearchParams(url.searchParams);
		params.delete("account_settings");
		params.delete("account_link_provider");
		params.delete("account_link_status");
		const query = params.toString();
		return `${url.pathname}${query ? `?${query}` : ""}${url.hash}`;
	}

	function showAccountLinkToast(provider: string, status: string) {
		const providerLabel = provider === "github" ? "GitHub" : provider === "google" ? "Google" : m.nav_account();
		switch (status) {
			case "linked":
				toast.success(m.nav_account_linked({ provider: providerLabel }));
				break;
			case "conflict":
				toast.error(m.nav_account_link_conflict({ provider: providerLabel }));
				break;
			case "denied":
				toast.error(m.nav_account_link_denied({ provider: providerLabel }));
				break;
			case "configuration":
				toast.error(m.nav_account_link_not_configured({ provider: providerLabel }));
				break;
			case "auth":
			case "invalid":
				toast.error(m.nav_account_link_sign_in_again());
				break;
			default:
				toast.error(m.nav_account_link_failed({ provider: providerLabel }));
				break;
		}
	}

	async function submitLogout(event: SubmitEvent) {
		event.preventDefault();
		if (logoutPending) return;

		const form = event.currentTarget as HTMLFormElement;
		logoutPending = true;

		try {
			const response = await fetch(form.action, { method: form.method });
			if (!response.ok) throw new Error(m.nav_logout_failed());
			await goto(logoutRedirectPath(response), { invalidateAll: true });
		} catch {
			logoutPending = false;
			toast.error(m.nav_logout_failed());
		}
	}

	function logoutRedirectPath(response: Response) {
		if (!response.redirected) return "/login";

		const url = new URL(response.url);
		return `${url.pathname}${url.search}${url.hash}`;
	}
</script>

<Sidebar.Menu>
	<Sidebar.MenuItem>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Sidebar.MenuButton
						{...props}
						size="lg"
						class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
					>
						<Avatar.Root class="size-8 rounded-lg grayscale">
							{#if avatarSrc}
								<Avatar.Image src={avatarSrc} alt={displayName} />
							{/if}
							<Avatar.Fallback class="rounded-lg">{fallbackInitials}</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-start text-sm leading-tight">
							<span class="truncate font-medium">{displayName}</span>
							<span class="text-muted-foreground truncate text-xs">
								{displayEmail}
							</span>
						</div>
						<DotsVerticalIcon class="ms-auto size-4" />
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
				side={sidebar.isMobile ? "bottom" : "right"}
				align="end"
				sideOffset={4}
			>
				<DropdownMenu.Label class="p-0 font-normal">
					<div class="flex items-center gap-2 px-1 py-1.5 text-start text-sm">
						<Avatar.Root class="size-8 rounded-lg">
							{#if avatarSrc}
								<Avatar.Image src={avatarSrc} alt={displayName} />
							{/if}
							<Avatar.Fallback class="rounded-lg">{fallbackInitials}</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-start text-sm leading-tight">
							<span class="truncate font-medium">{displayName}</span>
							<span class="text-muted-foreground truncate text-xs">
								{displayEmail}
							</span>
						</div>
					</div>
				</DropdownMenu.Label>
				<DropdownMenu.Separator />
				<DropdownMenu.Group>
					<DropdownMenu.Item
						onclick={() => {
							accountSettingsInitialSection = "general";
							accountSettingsOpen = true;
						}}
					>
						<UserCircleIcon />
						{m.nav_account()}
					</DropdownMenu.Item>
					<DropdownMenu.Item onclick={() => (billingProfileOpen = true)}>
						<CreditCardIcon />
						{m.nav_billing()}
					</DropdownMenu.Item>
					<DropdownMenu.Item>
						<NotificationIcon />
						{m.nav_notifications()}
					</DropdownMenu.Item>
				</DropdownMenu.Group>
				<DropdownMenu.Separator />
				<DropdownMenu.Item variant="destructive" onclick={() => (logoutConfirmationOpen = true)}>
					<LogoutIcon />
					{m.nav_log_out()}
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>

<AccountSettingsDialog
	bind:open={accountSettingsOpen}
	initialSection={accountSettingsInitialSection}
	{user}
/>
<BillingProfileDialog bind:open={billingProfileOpen} {user} />

<AlertDialog.Root bind:open={logoutConfirmationOpen}>
	<AlertDialog.Content>
		<form action="/logout" method="POST" class="flex flex-col gap-4" onsubmit={submitLogout}>
			<AlertDialog.Header>
				<AlertDialog.Title>{m.nav_logout_title()}</AlertDialog.Title>
				<AlertDialog.Description>
					{m.nav_logout_description()}
				</AlertDialog.Description>
			</AlertDialog.Header>
			<AlertDialog.Footer>
				<AlertDialog.Cancel type="button" disabled={logoutPending}>{m.common_cancel()}</AlertDialog.Cancel>
				<AlertDialog.Action type="submit" variant="destructive" loading={logoutPending} disabled={logoutPending}>
					<LogoutIcon />
					{m.nav_log_out()}
				</AlertDialog.Action>
			</AlertDialog.Footer>
		</form>
	</AlertDialog.Content>
</AlertDialog.Root>
