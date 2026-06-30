<script lang="ts">
	import CameraIcon from "@tabler/icons-svelte/icons/camera";
	import LockIcon from "@tabler/icons-svelte/icons/lock";
	import MailIcon from "@tabler/icons-svelte/icons/mail";
	import ShieldCheckIcon from "@tabler/icons-svelte/icons/shield-check";
	import TrashIcon from "@tabler/icons-svelte/icons/trash";
	import UserCircleIcon from "@tabler/icons-svelte/icons/user-circle";
	import { invalidateAll } from "$app/navigation";
	import { toast } from "svelte-sonner";
	import {
		avatarFileError,
		clearedAccountSettingsPendingState,
		fileToDataURL,
		passwordConfirmationError,
		shouldApplyAccountSettingsResult,
		trimOptionalValue
	} from "./account-settings-utils";
	import * as Avatar from "$lib/components/ui/avatar/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { confirmDelete } from "$lib/components/ui/confirm-delete-dialog/index.js";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import * as Select from "$lib/components/ui/select/index.js";
	import * as ImageCropper from "$lib/components/ui/image-cropper/index.js";
	import { getFileFromUrl } from "$lib/components/ui/image-cropper/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import { getLocale } from "$lib/paraglide/runtime.js";
	import { cn } from "$lib/utils.js";

	type AuthUser = NonNullable<App.Locals["user"]>;
	type SettingsSection = "general" | "security" | "sessions" | "linked" | "danger";
	type Message = { kind: "success" | "error"; text: string } | null;
	type UpdatePayload = {
		name?: string;
		email?: string;
		image?: string | null;
		preferredLanguage?: "en" | "ro";
		password?: string;
	};
	type PreferredLanguage = NonNullable<UpdatePayload["preferredLanguage"]>;
	type AuthSessionListItem = {
		id: string;
		userId: string;
		expiresAt: string;
		ipAddress?: string;
		userAgent?: string;
		createdAt: string;
		updatedAt: string;
		current: boolean;
	};
	type AuthAccountListItem = {
		id: string;
		providerId: "credential" | "google" | "github" | string;
		createdAt: string;
		updatedAt: string;
	};
	type ManagedOAuthProvider = "google" | "github";
	type ProviderConfig = {
		id: ManagedOAuthProvider;
		label: string;
		description: string;
		connectHref: string;
	};

	type Props = {
		open: boolean;
		user: AuthUser | null;
		initialSection?: SettingsSection;
	};

	let { open = $bindable(false), user, initialSection = "general" }: Props = $props();

	let activeSection = $state<SettingsSection>("general");
	let avatarPreview = $state<string>("");
	let nameValue = $state("");
	let emailValue = $state("");
	let languageValue = $state<PreferredLanguage>("en");
	let passwordValue = $state("");
	let passwordConfirmationValue = $state("");
	let avatarPending = $state(false);
	let namePending = $state(false);
	let emailPending = $state(false);
	let languagePending = $state(false);
	let passwordPending = $state(false);
	let avatarMessage = $state<Message>(null);
	let nameMessage = $state<Message>(null);
	let emailMessage = $state<Message>(null);
	let passwordMessage = $state<Message>(null);
	let sessions = $state<AuthSessionListItem[]>([]);
	let sessionsPending = $state(false);
	let sessionsMessage = $state<Message>(null);
	let sessionsLoadedUserId = $state<string | null>(null);
	let sessionsAttemptedUserId = $state<string | null>(null);
	let sessionsRequestId = 0;
	let revokingSessionId = $state<string | null>(null);
	let accounts = $state<AuthAccountListItem[]>([]);
	let accountsPending = $state(false);
	let accountsMessage = $state<Message>(null);
	let accountsLoadedUserId = $state<string | null>(null);
	let accountsAttemptedUserId = $state<string | null>(null);
	let accountsRequestId = 0;
	let unlinkingProviderId = $state<ManagedOAuthProvider | null>(null);

	let wasOpen = false;
	let lastInitializedUserId: string | null = null;
	let lastAppliedInitialSection: SettingsSection | null = null;

	const avatarMessageId = "account-avatar-message";
	const nameMessageId = "account-name-message";
	const emailMessageId = "account-email-message";
	const passwordMessageId = "account-password-message";
	const sessionsMessageId = "account-sessions-message";
	const accountsMessageId = "account-linked-accounts-message";

	const languageOptions: Array<{ value: PreferredLanguage; label: string }> = [
		{ value: "en", label: m.common_english() },
		{ value: "ro", label: m.common_romanian() }
	];

	const dateFormatter = $derived(new Intl.DateTimeFormat(getLocale(), {
		dateStyle: "medium",
		timeStyle: "short"
	}));

	const displayName = $derived(user?.name?.trim() || m.account_settings_account_fallback());
	const displayEmail = $derived(user?.email?.trim() || m.account_settings_no_email_address());
	const fallbackInitials = $derived(getInitials(displayName, user?.email));
	const hasUser = $derived(Boolean(user));
	const credentialAccount = $derived(providerAccount("credential"));
	const signInMethodCount = $derived(new Set(accounts.map((account) => account.providerId)).size);
	const languageLabel = $derived(
		languageOptions.find((option) => option.value === languageValue)?.label ?? m.common_english()
	);

	const sections: Array<{ id: SettingsSection; label: string; disabled?: boolean }> = [
		{ id: "general", label: m.account_settings_general() },
		{ id: "security", label: m.account_settings_security() },
		{ id: "sessions", label: m.account_settings_sessions() },
		{ id: "linked", label: m.account_settings_linked_accounts() },
	];

	const oauthProviders: ProviderConfig[] = [
		{
			id: "google",
			label: "Google",
			description: m.account_settings_provider_google_description(),
			connectHref: "/api/auth/link/google"
		},
		{
			id: "github",
			label: "GitHub",
			description: m.account_settings_provider_github_description(),
			connectHref: "/api/auth/link/github"
		}
	];

	$effect(() => {
		const isOpen = open;
		const currentUser = user;
		const currentUserId = currentUser?.id ?? null;
		const requestedSection = initialSection;

		if (isOpen && !wasOpen) {
			resetFields(currentUser);
			resetManagementState();
			activeSection = requestedSection;
			passwordValue = "";
			passwordConfirmationValue = "";
			clearMessages();
			lastInitializedUserId = currentUserId;
			lastAppliedInitialSection = requestedSection;
			wasOpen = true;
			return;
		}

		if (isOpen && requestedSection !== lastAppliedInitialSection) {
			activeSection = requestedSection;
			lastAppliedInitialSection = requestedSection;
		}

		if (isOpen && currentUserId !== lastInitializedUserId) {
			resetFields(currentUser);
			resetManagementState();
			passwordValue = "";
			passwordConfirmationValue = "";
			clearMessages();
			lastInitializedUserId = currentUserId;
			return;
		}

		if (!isOpen) {
			resetFields(currentUser);
			lastInitializedUserId = currentUserId;
			lastAppliedInitialSection = null;
			wasOpen = false;
		}
	});

	$effect(() => {
		const isOpen = open;
		const currentUserId = user?.id ?? null;
		const section = activeSection;

		if (!isOpen || !currentUserId) return;

		if (
			section === "sessions" &&
			sessionsLoadedUserId !== currentUserId &&
			sessionsAttemptedUserId !== currentUserId
		) {
			void loadSessions();
		}

		if (
			section === "linked" &&
			accountsLoadedUserId !== currentUserId &&
			accountsAttemptedUserId !== currentUserId
		) {
			void loadAccounts();
		}
	});

	function resetFields(source: AuthUser | null) {
		nameValue = source?.name ?? "";
		emailValue = source?.email ?? "";
		languageValue = source?.preferredLanguage ?? "en";
		avatarPreview = source?.image ?? "";
		clearPending();
	}

	function resetManagementState() {
		sessions = [];
		sessionsPending = false;
		sessionsMessage = null;
		sessionsLoadedUserId = null;
		sessionsAttemptedUserId = null;
		sessionsRequestId += 1;
		revokingSessionId = null;
		accounts = [];
		accountsPending = false;
		accountsMessage = null;
		accountsLoadedUserId = null;
		accountsAttemptedUserId = null;
		accountsRequestId += 1;
		unlinkingProviderId = null;
	}

	function clearMessages() {
		avatarMessage = null;
		nameMessage = null;
		emailMessage = null;
		passwordMessage = null;
		sessionsMessage = null;
		accountsMessage = null;
	}

	function clearPending() {
		const pending = clearedAccountSettingsPendingState();
		avatarPending = pending.avatarPending;
		namePending = pending.namePending;
		emailPending = pending.emailPending;
		languagePending = false;
		passwordPending = pending.passwordPending;
	}

	function getInitials(name: string, email?: string | null) {
		const parts = name.trim().split(/\s+/).filter(Boolean);
		if (parts.length >= 2) return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
		if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
		return (email?.trim()[0] ?? "A").toUpperCase();
	}

	async function accountSettingsRequest<T>(url: string, init?: RequestInit) {
		const response = await fetch(url, init);
		const data = await response.json().catch(() => null);

		if (!response.ok) {
			const message =
				data && typeof data === "object" && "error" in data && typeof data.error === "string"
					? data.error
					: m.account_settings_update_error();
			throw new Error(message);
		}

		return data as T;
	}

	async function patchUser(payload: UpdatePayload) {
		const response = await fetch("/api/auth/user", {
			method: "PATCH",
			headers: {
				"content-type": "application/json"
			},
			body: JSON.stringify(payload)
		});
		const data = await response.json().catch(() => null);

		if (!response.ok) {
			const message =
				data && typeof data === "object" && "error" in data && typeof data.error === "string"
					? data.error
					: m.account_settings_save_error();
			throw new Error(message);
		}

		await invalidateAll();
		return data as AuthUser;
	}

	async function loadSessions(force = false) {
		if (!hasUser || sessionsPending) return;
		const requestUserId = user?.id ?? null;
		if (!requestUserId) return;
		if (!force && sessionsAttemptedUserId === requestUserId) return;

		const requestId = ++sessionsRequestId;
		sessionsPending = true;
		sessionsMessage = null;
		sessionsAttemptedUserId = requestUserId;
		try {
			const data = await accountSettingsRequest<{ sessions: AuthSessionListItem[] }>(
				"/api/auth/sessions"
			);
			if (requestId !== sessionsRequestId || !shouldApplyAccountSettingsResult(requestUserId, user?.id)) {
				return;
			}
			sessions = data.sessions;
			sessionsLoadedUserId = requestUserId;
		} catch (error) {
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			sessionsMessage = { kind: "error", text: errorMessage(error) };
		} finally {
			if (requestId === sessionsRequestId && shouldApplyAccountSettingsResult(requestUserId, user?.id)) {
				sessionsPending = false;
			}
		}
	}

	async function loadAccounts(force = false) {
		if (!hasUser || accountsPending) return;
		const requestUserId = user?.id ?? null;
		if (!requestUserId) return;
		if (!force && accountsAttemptedUserId === requestUserId) return;

		const requestId = ++accountsRequestId;
		accountsPending = true;
		accountsMessage = null;
		accountsAttemptedUserId = requestUserId;
		try {
			const data = await accountSettingsRequest<{ accounts: AuthAccountListItem[] }>(
				"/api/auth/accounts"
			);
			if (requestId !== accountsRequestId || !shouldApplyAccountSettingsResult(requestUserId, user?.id)) {
				return;
			}
			accounts = data.accounts;
			accountsLoadedUserId = requestUserId;
		} catch (error) {
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			accountsMessage = { kind: "error", text: errorMessage(error) };
		} finally {
			if (requestId === accountsRequestId && shouldApplyAccountSettingsResult(requestUserId, user?.id)) {
				accountsPending = false;
			}
		}
	}

	function confirmSessionRevoke(session: AuthSessionListItem) {
		if (session.current || revokingSessionId) return;
		confirmDelete({
			title: m.account_settings_revoke_session_title(),
			description: m.account_settings_revoke_session_description({ session: sessionTitle(session) }),
			confirm: { text: m.account_settings_revoke() },
			onConfirm: () => revokeSession(session)
		});
	}

	async function revokeSession(session: AuthSessionListItem) {
		const requestUserId = user?.id ?? null;
		revokingSessionId = session.id;
		sessionsMessage = null;
		try {
			await accountSettingsRequest<{ deleted_id: string; deleted_count: number }>(
				`/api/auth/sessions/${encodeURIComponent(session.id)}`,
				{ method: "DELETE" }
			);
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			sessions = sessions.filter((item) => item.id !== session.id);
			toast.success(m.account_settings_session_revoked());
		} catch (error) {
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			toast.error(errorMessage(error));
		} finally {
			if (shouldApplyAccountSettingsResult(requestUserId, user?.id)) {
				revokingSessionId = null;
			}
		}
	}

	function confirmProviderUnlink(provider: ProviderConfig) {
		const account = providerAccount(provider.id);
		if (!account || unlinkingProviderId || signInMethodCount <= 1) return;
		confirmDelete({
			title: m.account_settings_unlink_provider_title({ provider: provider.label }),
			description: m.account_settings_unlink_provider_description({ provider: provider.label }),
			confirm: { text: m.account_settings_unlink() },
			onConfirm: () => unlinkProvider(provider.id)
		});
	}

	async function unlinkProvider(providerId: ManagedOAuthProvider) {
		const requestUserId = user?.id ?? null;
		unlinkingProviderId = providerId;
		accountsMessage = null;
		try {
			await accountSettingsRequest<{ deleted_provider_id: string; deleted_count: number }>(
				`/api/auth/accounts/${providerId}`,
				{ method: "DELETE" }
			);
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			accounts = accounts.filter((account) => account.providerId !== providerId);
			toast.success(m.account_settings_linked_account_removed());
		} catch (error) {
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			toast.error(errorMessage(error));
		} finally {
			if (shouldApplyAccountSettingsResult(requestUserId, user?.id)) {
				unlinkingProviderId = null;
			}
		}
	}

	async function handleAvatarCropped(croppedUrl: string) {
		if (!hasUser || avatarPending) return;
		const requestUserId = user?.id ?? null;

		avatarPending = true;
		avatarMessage = null;
		try {
			const file = await getFileFromUrl(croppedUrl);
			const validationError = avatarFileError(file);
			if (validationError) {
				avatarMessage = { kind: "error", text: validationError };
				return;
			}

			const dataURL = await fileToDataURL(file);
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			const updated = await patchUser({ image: dataURL });
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			avatarPreview = updated.image ?? "";
			toast.success(m.account_settings_avatar_saved());
		} catch (error) {
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			toast.error(errorMessage(error));
		} finally {
			if (shouldApplyAccountSettingsResult(requestUserId, user?.id)) {
				avatarPending = false;
			}
		}
	}

	function handleUnsupportedFile(file: File) {
		const validationError = avatarFileError(file);
		if (validationError) {
			toast.error(validationError);
		}
	}

	async function saveName() {
		if (!hasUser || namePending) return;
		const requestUserId = user?.id ?? null;
		const name = trimOptionalValue(nameValue);
		namePending = true;
		try {
			const updated = await patchUser({ name });
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			nameValue = updated.name;
			toast.success(m.account_settings_name_saved());
		} catch (error) {
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			toast.error(errorMessage(error));
		} finally {
			if (shouldApplyAccountSettingsResult(requestUserId, user?.id)) {
				namePending = false;
			}
		}
	}

	async function saveEmail() {
		if (!hasUser || emailPending) return;
		const requestUserId = user?.id ?? null;
		const email = trimOptionalValue(emailValue);
		emailPending = true;
		emailMessage = null;
		try {
			const updated = await patchUser({ email });
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			emailValue = updated.email;
			toast.success(m.account_settings_email_saved());
		} catch (error) {
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			toast.error(errorMessage(error));
		} finally {
			if (shouldApplyAccountSettingsResult(requestUserId, user?.id)) {
				emailPending = false;
			}
		}
	}

	async function saveLanguage() {
		if (!hasUser || languagePending) return;
		const requestUserId = user?.id ?? null;
		const previousLanguage = user?.preferredLanguage ?? "en";
		languagePending = true;
		try {
			const updated = await patchUser({ preferredLanguage: languageValue });
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			languageValue = updated.preferredLanguage;
			toast.success(m.account_settings_language_saved());
			if (languageValue !== previousLanguage && typeof window !== "undefined") {
				window.location.reload();
			}
		} catch (error) {
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			toast.error(errorMessage(error));
		} finally {
			if (shouldApplyAccountSettingsResult(requestUserId, user?.id)) {
				languagePending = false;
			}
		}
	}

	async function savePassword() {
		if (!hasUser || passwordPending) return;
		const requestUserId = user?.id ?? null;
		const validationError = passwordConfirmationError(passwordValue, passwordConfirmationValue);
		if (validationError) {
			passwordMessage = { kind: "error", text: validationError };
			return;
		}

		passwordPending = true;
		passwordMessage = null;
		try {
			await patchUser({ password: passwordValue });
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			passwordValue = "";
			passwordConfirmationValue = "";
			passwordMessage = { kind: "success", text: m.account_settings_password_updated() };
			accountsAttemptedUserId = null;
			if (activeSection === "linked") void loadAccounts(true);
		} catch (error) {
			if (!shouldApplyAccountSettingsResult(requestUserId, user?.id)) return;
			passwordMessage = { kind: "error", text: errorMessage(error) };
		} finally {
			if (shouldApplyAccountSettingsResult(requestUserId, user?.id)) {
				passwordPending = false;
			}
		}
	}

	function providerAccount(providerId: string) {
		return accounts.find((account) => account.providerId === providerId) ?? null;
	}

	function sessionTitle(session: AuthSessionListItem) {
		if (session.current) return m.account_settings_current_session();
		const userAgent = session.userAgent?.trim();
		return userAgent ? userAgent : m.account_settings_browser_session();
	}

	function sessionMeta(session: AuthSessionListItem) {
		const ipAddress = session.ipAddress?.trim();
		const createdAt = formatDateTime(session.createdAt);
		return ipAddress
			? m.account_settings_session_ip_created_at({ ip: ipAddress, date: createdAt })
			: m.account_settings_session_created_at({ date: createdAt });
	}

	function formatDateTime(value: string) {
		const date = new Date(value);
		if (Number.isNaN(date.getTime())) return m.account_settings_unknown();
		return dateFormatter.format(date);
	}

	function errorMessage(error: unknown) {
		return error instanceof Error ? error.message : m.account_settings_save_error();
	}

	function messageClass(message: Message) {
		return message?.kind === "error" ? "text-destructive" : "text-muted-foreground";
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="flex h-[min(88vh,52rem)] w-full gap-0 overflow-hidden p-0 sm:max-w-[56rem]">
		<div class="flex w-full flex-col">
			<Dialog.Header class="border-b px-6 py-4">
				<Dialog.Title>{m.account_settings_title()}</Dialog.Title>
				<Dialog.Description>{m.account_settings_description()}</Dialog.Description>
			</Dialog.Header>

			<div class="grid min-h-0 flex-1 grid-cols-1 md:grid-cols-[14rem_1fr]">
				<nav class="border-b p-3 md:border-r md:border-b-0" aria-label={m.account_settings_nav_label()}>
					<div class="flex gap-1 overflow-x-auto md:grid md:overflow-visible">
						{#each sections as section (section.id)}
							<button
								type="button"
								class="text-muted-foreground hover:bg-muted hover:text-foreground disabled:text-muted-foreground/50 flex h-9 shrink-0 items-center rounded-md px-3 text-sm font-medium whitespace-nowrap transition-colors disabled:pointer-events-none {activeSection ===
								section.id
									? 'bg-muted text-foreground'
									: ''}"
								disabled={section.disabled}
								aria-current={activeSection === section.id ? "page" : undefined}
								onclick={() => (activeSection = section.id)}
							>
								{section.label}
							</button>
						{/each}
					</div>
				</nav>

				<div class="min-h-0 overflow-y-auto p-6">
					{#if activeSection === "general"}
						<div class="space-y-6">
							<section class="space-y-4">
										<div>
											<h3 class="flex items-center gap-2 text-sm font-medium">
												<UserCircleIcon class="size-4" aria-hidden="true" />
										{m.account_settings_general()}
									</h3>
								</div>

								<div class="grid gap-4 rounded-lg border p-4">
									<ImageCropper.Root
										bind:src={avatarPreview}
										onCropped={handleAvatarCropped}
										onUnsupportedFile={handleUnsupportedFile}
										accept="image/png,image/jpeg,image/jpg,image/gif,image/avif,image/apng,image/svg+xml,image/webp"
									>
										<ImageCropper.UploadTrigger
											class={cn(
												"border-border/60 bg-muted/20 hover:bg-muted/40 hover:border-border focus-within:ring-ring block w-full cursor-pointer space-y-4 rounded-xl border p-6 transition-all focus-within:ring-2 focus-within:ring-offset-2",
												(avatarPending || !hasUser) && "pointer-events-none opacity-50"
											)}
										>
											<div class="space-y-4">
												<div>
													<h4 class="text-foreground text-sm leading-none font-semibold tracking-tight">{m.account_settings_avatar()}</h4>
													<p class="text-muted-foreground mt-1.5 text-xs">
														{m.account_settings_avatar_description()}
													</p>
												</div>

												<div class="flex flex-col gap-3">
													<Avatar.Root class="border-border/80 size-20 rounded-full border shadow-sm">
														{#if avatarPreview}
															<Avatar.Image
																src={avatarPreview}
																alt={displayName}
																class="size-20 rounded-full object-cover"
															/>
														{/if}
														<Avatar.Fallback
															class="bg-muted text-muted-foreground flex items-center justify-center rounded-full text-lg font-medium"
														>
															{fallbackInitials}
														</Avatar.Fallback>
													</Avatar.Root>

													<div class="space-y-1">
														<p class="text-foreground text-sm font-semibold">
															{avatarPending ? m.account_settings_avatar_uploading() : m.account_settings_avatar_upload()}
														</p>
														<p class="text-muted-foreground text-xs font-normal">
															{m.account_settings_avatar_file_hint()}
														</p>
													</div>
												</div>
											</div>
										</ImageCropper.UploadTrigger>

										<ImageCropper.Dialog>
											<div class="flex flex-col space-y-1.5 text-center sm:text-left">
												<h2 class="text-lg leading-none font-semibold tracking-tight">{m.account_settings_crop_avatar()}</h2>
												<p class="text-muted-foreground text-sm">
													{m.account_settings_crop_avatar_description()}
												</p>
											</div>
											<div class="bg-muted relative h-80 w-full overflow-hidden rounded-lg border">
												<ImageCropper.Cropper />
											</div>
											<ImageCropper.Controls class="mt-2 flex justify-end gap-2">
												<ImageCropper.Cancel />
												<ImageCropper.Crop />
											</ImageCropper.Controls>
										</ImageCropper.Dialog>
									</ImageCropper.Root>
								</div>

								<form
									class="grid gap-4 rounded-lg border p-4"
									onsubmit={(event) => {
										event.preventDefault();
										void saveName();
									}}
								>
									<div class="grid gap-2">
										<Label for="account-name">{m.account_settings_display_name()}</Label>
										<div class="flex flex-wrap items-center gap-3">
											<Input
												id="account-name"
												class="min-w-0 flex-1"
												bind:value={nameValue}
												disabled={!hasUser}
												autocomplete="name"
											/>
											<Button
												type="submit"
												class="shrink-0"
												disabled={!hasUser || namePending || nameValue === (user?.name ?? "")}
											>
												{namePending ? m.common_saving() : m.account_settings_save_name()}
											</Button>
										</div>
									</div>
								</form>

								<form
									class="grid gap-4 rounded-lg border p-4"
									onsubmit={(event) => {
										event.preventDefault();
										void saveEmail();
									}}
								>
									<div class="grid gap-2">
										<Label for="account-email">{m.account_settings_email_address()}</Label>
										<div class="flex flex-wrap items-center gap-3">
											<Input
												id="account-email"
												type="email"
												class="min-w-0 flex-1"
												bind:value={emailValue}
												disabled={!hasUser}
												autocomplete="email"
											/>
											<Button
												type="submit"
												class="shrink-0"
												disabled={!hasUser || emailPending || emailValue === (user?.email ?? "")}
											>
												{emailPending ? m.common_saving() : m.account_settings_save_email()}
											</Button>
										</div>
									</div>
								</form>

								<form
									class="grid gap-4 rounded-lg border p-4"
									onsubmit={(event) => {
										event.preventDefault();
										void saveLanguage();
									}}
								>
									<div class="grid gap-2">
										<Label for="account-language">{m.account_settings_language()}</Label>
										<div class="flex flex-wrap items-center gap-3">
											<Select.Root type="single" bind:value={languageValue}>
												<Select.Trigger
													id="account-language"
													class="w-full justify-between sm:w-52"
													disabled={!hasUser}
												>
													<span data-slot="select-value">{languageLabel}</span>
												</Select.Trigger>
												<Select.Content>
													<Select.Item value="en">{m.common_english()}</Select.Item>
													<Select.Item value="ro">{m.common_romanian()}</Select.Item>
												</Select.Content>
											</Select.Root>
											<Button
												type="submit"
												class="shrink-0"
												disabled={!hasUser || languagePending || languageValue === (user?.preferredLanguage ?? "en")}
											>
												{languagePending ? m.common_saving() : m.account_settings_save_language()}
											</Button>
										</div>
									</div>
								</form>
							</section>
						</div>
					{:else if activeSection === "security"}
						<section class="space-y-4">
							<div>
								<h3 class="flex items-center gap-2 text-sm font-medium">
									<LockIcon class="size-4" aria-hidden="true" />
									{m.account_settings_security()}
								</h3>
								<p class="text-muted-foreground mt-1 text-sm">{m.account_settings_security_description()}</p>
							</div>

							<form
								class="grid gap-4 rounded-lg border p-4"
								onsubmit={(event) => {
									event.preventDefault();
									void savePassword();
								}}
							>
								<div class="grid gap-2">
									<Label for="account-password">{m.account_settings_new_password()}</Label>
									<Input
										id="account-password"
										type="password"
										bind:value={passwordValue}
										disabled={!hasUser}
										autocomplete="new-password"
										aria-invalid={passwordMessage?.kind === "error"}
										aria-describedby={passwordMessage ? passwordMessageId : undefined}
									/>
								</div>
								<div class="grid gap-2">
									<Label for="account-password-confirmation">{m.account_settings_confirm_password()}</Label>
									<Input
										id="account-password-confirmation"
										type="password"
										bind:value={passwordConfirmationValue}
										disabled={!hasUser}
										autocomplete="new-password"
										aria-invalid={passwordMessage?.kind === "error"}
										aria-describedby={passwordMessage ? passwordMessageId : undefined}
									/>
								</div>
								<div class="flex flex-wrap items-center justify-between gap-3">
									{#if passwordMessage}
										<p
											id={passwordMessageId}
											class="text-sm {messageClass(passwordMessage)}"
											role={passwordMessage.kind === "error" ? "alert" : "status"}
											aria-live={passwordMessage.kind === "error" ? "assertive" : "polite"}
										>
											{passwordMessage.text}
										</p>
									{:else}
										<span></span>
									{/if}
									<Button type="submit" disabled={!hasUser || passwordPending || !passwordValue}>
										{passwordPending ? m.common_saving() : m.account_settings_save_password()}
									</Button>
								</div>
							</form>
						</section>
					{:else if activeSection === "sessions"}
						<section class="space-y-4">
							<div class="flex flex-col justify-between gap-3 sm:flex-row sm:items-start">
								<div>
									<h3 class="flex items-center gap-2 text-sm font-medium">
										<ShieldCheckIcon class="size-4" aria-hidden="true" />
										{m.account_settings_sessions()}
									</h3>
									<p class="text-muted-foreground mt-1 text-sm">
										{m.account_settings_sessions_description()}
									</p>
								</div>
								<Button
									type="button"
									variant="outline"
									size="sm"
									disabled={!hasUser || sessionsPending}
									onclick={() => void loadSessions(true)}
								>
									{sessionsPending ? m.common_loading() : m.common_refresh()}
								</Button>
							</div>

							{#if sessionsMessage}
								<div
									id={sessionsMessageId}
									class="border-destructive/40 bg-destructive/5 text-destructive rounded-lg border p-3 text-sm"
									role="alert"
								>
									{sessionsMessage.text}
								</div>
							{/if}

							{#if sessionsPending && sessions.length === 0}
								<div class="text-muted-foreground rounded-lg border p-4 text-sm" role="status">
									{m.account_settings_loading_sessions()}
								</div>
							{:else if !sessionsMessage && sessions.length === 0}
								<div class="text-muted-foreground rounded-lg border p-4 text-sm">
									{m.account_settings_no_sessions()}
								</div>
							{:else}
								<div class="grid gap-3">
									{#each sessions as session (session.id)}
										<div class="grid gap-3 rounded-lg border p-4 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-start">
											<div class="min-w-0 space-y-1">
												<div class="flex flex-wrap items-center gap-2">
													<h4 class="min-w-0 truncate text-sm font-medium">{sessionTitle(session)}</h4>
													{#if session.current}
														<span class="bg-muted text-muted-foreground rounded-md px-2 py-0.5 text-xs font-medium">
															{m.account_settings_current()}
														</span>
													{/if}
												</div>
												<p class="text-muted-foreground text-xs">{sessionMeta(session)}</p>
												<p class="text-muted-foreground text-xs">
													{m.account_settings_expires({ date: formatDateTime(session.expiresAt) })}
												</p>
											</div>
											<Button
												type="button"
												variant="destructive"
												size="sm"
												disabled={session.current || revokingSessionId === session.id}
												aria-label={session.current ? m.account_settings_current_session_cannot_revoke() : m.account_settings_revoke_session()}
												onclick={() => confirmSessionRevoke(session)}
											>
												<TrashIcon class="size-4" aria-hidden="true" />
												{revokingSessionId === session.id ? m.account_settings_revoking() : m.account_settings_revoke()}
											</Button>
										</div>
									{/each}
								</div>
							{/if}
						</section>
					{:else if activeSection === "linked"}
						<section class="space-y-4">
							<div class="flex flex-col justify-between gap-3 sm:flex-row sm:items-start">
								<div>
									<h3 class="flex items-center gap-2 text-sm font-medium">
										<UserCircleIcon class="size-4" aria-hidden="true" />
										{m.account_settings_linked_accounts()}
									</h3>
									<p class="text-muted-foreground mt-1 text-sm">
										{m.account_settings_linked_accounts_description()}
									</p>
								</div>
								<Button
									type="button"
									variant="outline"
									size="sm"
									disabled={!hasUser || accountsPending}
									onclick={() => void loadAccounts(true)}
								>
									{accountsPending ? m.common_loading() : m.common_refresh()}
								</Button>
							</div>

							{#if accountsMessage}
								<div
									id={accountsMessageId}
									class="border-destructive/40 bg-destructive/5 text-destructive rounded-lg border p-3 text-sm"
									role="alert"
								>
									{accountsMessage.text}
								</div>
							{/if}

							{#if accountsPending && accounts.length === 0}
								<div class="text-muted-foreground rounded-lg border p-4 text-sm" role="status">
									{m.account_settings_loading_linked_accounts()}
								</div>
							{/if}

							{#if !accountsPending && !accountsMessage && accounts.length === 0}
								<div class="text-muted-foreground rounded-lg border p-4 text-sm">
									{m.account_settings_no_sign_in_methods()}
								</div>
							{/if}

							<div class="grid gap-3">
								<div class="grid gap-3 rounded-lg border p-4 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center">
									<div class="min-w-0 space-y-1">
										<h4 class="flex items-center gap-2 text-sm font-medium">
											<MailIcon class="size-4" aria-hidden="true" />
											{m.account_settings_email_password()}
										</h4>
										<p class="text-muted-foreground text-xs">
											{#if credentialAccount}
												{m.account_settings_password_enabled({ email: displayEmail })}
											{:else}
												{m.account_settings_add_password()}
											{/if}
										</p>
									</div>
									{#if credentialAccount}
										<span class="bg-muted text-muted-foreground rounded-md px-2 py-1 text-xs font-medium">
											{m.common_connected()}
										</span>
									{:else}
										<Button type="button" variant="outline" size="sm" onclick={() => (activeSection = "security")}>
											{m.account_settings_set_password()}
										</Button>
									{/if}
								</div>

								{#each oauthProviders as provider (provider.id)}
									{@const account = providerAccount(provider.id)}
									<div class="grid gap-3 rounded-lg border p-4 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center">
										<div class="min-w-0 space-y-1">
											<h4 class="text-sm font-medium">{provider.label}</h4>
											<p class="text-muted-foreground text-xs">{provider.description}</p>
											{#if account}
												<p class="text-muted-foreground text-xs">
													{m.account_settings_linked_at({ date: formatDateTime(account.createdAt) })}
												</p>
											{/if}
										</div>
										{#if account}
											<Button
												type="button"
												variant="destructive"
												size="sm"
												disabled={unlinkingProviderId === provider.id || signInMethodCount <= 1}
												onclick={() => confirmProviderUnlink(provider)}
											>
												{unlinkingProviderId === provider.id ? m.account_settings_unlinking() : m.account_settings_unlink()}
											</Button>
										{:else}
											<Button href={provider.connectHref} variant="outline" size="sm" disabled={!hasUser}>
												{m.common_connect()}
											</Button>
										{/if}
									</div>
								{/each}
							</div>
						</section>
					{:else}
						<section class="rounded-lg border p-4">
							<h3 class="text-sm font-medium">{m.account_settings_unavailable_title()}</h3>
							<p class="text-muted-foreground mt-1 text-sm">{m.account_settings_unavailable_body()}</p>
						</section>
					{/if}
				</div>
			</div>
		</div>
	</Dialog.Content>
</Dialog.Root>
