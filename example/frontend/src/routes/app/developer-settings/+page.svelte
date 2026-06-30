<script lang="ts">
	import { page } from "$app/state";
	import AlertCircleIcon from "@lucide/svelte/icons/alert-circle";
	import CalendarIcon from "@lucide/svelte/icons/calendar";
	import CopyIcon from "@lucide/svelte/icons/copy";
	import EyeOffIcon from "@lucide/svelte/icons/eye-off";
	import KeyRoundIcon from "@lucide/svelte/icons/key-round";
	import PlusIcon from "@lucide/svelte/icons/plus";
	import Trash2Icon from "@lucide/svelte/icons/trash-2";
	import WebhookIcon from "@lucide/svelte/icons/webhook";
	import { getLocalTimeZone, today, type DateValue } from "@internationalized/date";
	import { createMutation, createQuery, useQueryClient } from "@tanstack/svelte-query";
	import { toast } from "svelte-sonner";

	import * as Alert from "$lib/components/ui/alert/index.js";
	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Calendar from "$lib/components/ui/calendar/index.js";
	import { Checkbox } from "$lib/components/ui/checkbox/index.js";
	import { confirmDelete } from "$lib/components/ui/confirm-delete-dialog/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import * as Table from "$lib/components/ui/table/index.js";
	import {
		createAPIKey,
		deleteAPIKey,
		deleteWebhook,
		fetchAPIKeys,
		fetchWebhook,
		regenerateWebhookSecret,
		saveWebhook,
		type APIKeyResponse,
		type CreateAPIKeyInput,
		type WebhookEvent,
		type WebhookResponse
	} from "./api";
	import {
		apiKeyExpirationDateToISOString,
		apiKeyExpirationPresetDate,
		isSameCalendarDate,
		type APIKeyExpirationPreset
	} from "./expiration-utils";

	const API_KEYS_QUERY_KEY = ["auth", "api-keys"];
	const WEBHOOK_QUERY_KEY = ["auth", "webhook"];
	const API_KEY_COLUMN_COUNT = 5;
	const WEBHOOK_EVENT_OPTIONS: { value: WebhookEvent; label: string }[] = [
		{ value: "job.started", label: "Job started" },
		{ value: "job.failed", label: "Job failed" },
		{ value: "job.succeeded", label: "Job succeeded" }
	];
	const dateFormatter = new Intl.DateTimeFormat("en", {
		month: "short",
		day: "numeric",
		year: "numeric"
	});
	type WebhookQueryData = { webhook: WebhookResponse | null };

	let nameValue = $state("");
	let expirationDate = $state<DateValue | undefined>();
	let pendingExpirationDate = $state<DateValue | undefined>();
	let expirationPopoverOpen = $state(false);
	let createdAPIKey = $state<APIKeyResponse | null>(null);
	let webhookURL = $state("");
	let webhookEvents = $state<WebhookEvent[]>([]);
	let visibleWebhookSecret = $state("");
	let syncedWebhookCacheKey = $state("");

	const queryClient = useQueryClient();
	const activeTab = $derived(page.url.searchParams.get("tab") === "webhooks" ? "webhooks" : "api-keys");

	const apiKeysQuery = createQuery<{ api_keys: APIKeyResponse[] }, Error>(() => ({
		queryKey: API_KEYS_QUERY_KEY,
		queryFn: () => fetchAPIKeys(fetch)
	}));
	const webhookQuery = createQuery<WebhookQueryData, Error>(() => ({
		queryKey: WEBHOOK_QUERY_KEY,
		queryFn: () => fetchWebhook(fetch),
		enabled: activeTab === "webhooks"
	}));
	const createAPIKeyMutation = createMutation<APIKeyResponse, Error, CreateAPIKeyInput>(() => ({
		mutationKey: ["auth", "api-keys", "create"],
		mutationFn: (input) => createAPIKey(fetch, input),
		onSuccess: (result) => {
			createdAPIKey = result;
			nameValue = "";
			expirationDate = undefined;
			pendingExpirationDate = undefined;
			expirationPopoverOpen = false;
			queryClient.setQueryData<{ api_keys: APIKeyResponse[] }>(API_KEYS_QUERY_KEY, (current) => {
				const publicKey = hiddenAPIKey(result);
				const existing = current?.api_keys ?? [];
				return {
					api_keys: [publicKey, ...existing.filter((key) => key.id !== publicKey.id)]
				};
			});
			void queryClient.invalidateQueries({ queryKey: API_KEYS_QUERY_KEY });
			toast.success("API key created.");
		}
	}));
	const deleteMutation = createMutation<{ deleted_id: string; deleted_count: number }, Error, string>(
		() => ({
			mutationKey: ["auth", "api-keys", "delete"],
			mutationFn: (id) => deleteAPIKey(fetch, id),
			onSuccess: (result) => {
				queryClient.setQueryData<{ api_keys: APIKeyResponse[] }>(API_KEYS_QUERY_KEY, (current) => ({
					api_keys: (current?.api_keys ?? []).filter((key) => key.id !== result.deleted_id)
				}));
				if (createdAPIKey?.id === result.deleted_id) createdAPIKey = null;
				void queryClient.invalidateQueries({ queryKey: API_KEYS_QUERY_KEY });
				toast.success("API key deleted.");
			}
		})
	);
	const saveWebhookMutation = createMutation<WebhookResponse, Error, void>(() => ({
		mutationKey: ["auth", "webhook", "save"],
		mutationFn: () => saveWebhook(fetch, { url: webhookURL.trim(), events_active: webhookEvents }),
		onSuccess: (result) => {
			setWebhookCache(result);
			visibleWebhookSecret = result.secret_key ?? "";
			void queryClient.invalidateQueries({ queryKey: WEBHOOK_QUERY_KEY });
			toast.success("Webhook saved.");
		}
	}));
	const regenerateWebhookSecretMutation = createMutation<WebhookResponse, Error, void>(() => ({
		mutationKey: ["auth", "webhook", "secret", "regenerate"],
		mutationFn: () => regenerateWebhookSecret(fetch),
		onSuccess: (result) => {
			setWebhookCache(result);
			visibleWebhookSecret = result.secret_key ?? "";
			void queryClient.invalidateQueries({ queryKey: WEBHOOK_QUERY_KEY });
			toast.success("Webhook secret regenerated.");
		}
	}));
	const deleteWebhookMutation = createMutation<{ deleted_id: string; deleted_count: number }, Error, void>(
		() => ({
			mutationKey: ["auth", "webhook", "delete"],
			mutationFn: () => deleteWebhook(fetch),
			onSuccess: () => {
				webhookURL = "";
				webhookEvents = [];
				visibleWebhookSecret = "";
				queryClient.setQueryData<WebhookQueryData>(WEBHOOK_QUERY_KEY, { webhook: null });
				void queryClient.invalidateQueries({ queryKey: WEBHOOK_QUERY_KEY });
				toast.success("Webhook deleted.");
			}
		})
	);

	const apiKeys = $derived(apiKeysQuery.data?.api_keys ?? []);
	const webhook = $derived(webhookQuery.data?.webhook ?? null);
	const webhookExists = $derived(webhook !== null);
	const trimmedWebhookURL = $derived(webhookURL.trim());
	const trimmedName = $derived(nameValue.trim());
	const createError = $derived(createAPIKeyMutation.error?.message ?? "");
	const deleteError = $derived(deleteMutation.error?.message ?? "");
	const webhookError = $derived(
		webhookQuery.error?.message ??
			saveWebhookMutation.error?.message ??
			regenerateWebhookSecretMutation.error?.message ??
			deleteWebhookMutation.error?.message ??
			""
	);
	const webhookFormDisabled = $derived(
		webhookQuery.isLoading ||
			saveWebhookMutation.isPending ||
			regenerateWebhookSecretMutation.isPending ||
			deleteWebhookMutation.isPending
	);
	const webhookCacheKey = $derived.by(() => {
		if (webhookQuery.data === undefined) return "";
		const currentWebhook = webhookQuery.data.webhook;
		if (!currentWebhook) return "none";
		return `${currentWebhook.id}:${currentWebhook.url}:${currentWebhook.updated_at}:${currentWebhook.events_active.join(",")}`;
	});
	const expirationLabel = $derived.by(() =>
		expirationDate ? formatCalendarDate(expirationDate) : "No expiration"
	);
	const activeExpirationPreset = $derived.by(() => {
		const selectedDate = expirationPopoverOpen ? pendingExpirationDate : expirationDate;
		if (!selectedDate) return null;

		const baseDate = today(getLocalTimeZone());
		const presets: APIKeyExpirationPreset[] = ["week", "month", "quarter"];
		return (
			presets.find((preset) =>
				isSameCalendarDate(selectedDate, apiKeyExpirationPresetDate(preset, baseDate))
			) ?? null
		);
	});

	$effect(() => {
		if (!webhookCacheKey || syncedWebhookCacheKey === webhookCacheKey) return;

		syncedWebhookCacheKey = webhookCacheKey;
		webhookURL = webhook?.url ?? "";
		webhookEvents = webhook ? [...webhook.events_active] : [];
	});

	function hiddenAPIKey(key: APIKeyResponse): APIKeyResponse {
		const { api_key: _apiKey, ...publicKey } = key;
		return publicKey;
	}

	function hiddenWebhookSecret(currentWebhook: WebhookResponse): WebhookResponse {
		const { secret_key: _secretKey, ...publicWebhook } = currentWebhook;
		return publicWebhook;
	}

	function setWebhookCache(currentWebhook: WebhookResponse) {
		queryClient.setQueryData<WebhookQueryData>(WEBHOOK_QUERY_KEY, {
			webhook: hiddenWebhookSecret(currentWebhook)
		});
	}

	function maskedAPIKey(key: APIKeyResponse) {
		return `${key.key_prefix}${"*".repeat(24)}`;
	}

	function formatDate(value: string | undefined) {
		if (!value) return "Never";

		const date = new Date(value);
		if (Number.isNaN(date.getTime())) return "Invalid date";
		return dateFormatter.format(date);
	}

	function formatCalendarDate(value: DateValue) {
		return dateFormatter.format(new Date(value.year, value.month - 1, value.day));
	}

	function setExpirationPopoverOpen(open: boolean) {
		expirationPopoverOpen = open;
		if (open) pendingExpirationDate = expirationDate;
	}

	function setExpirationPreset(preset: APIKeyExpirationPreset) {
		const nextExpirationDate = apiKeyExpirationPresetDate(preset, today(getLocalTimeZone()));
		pendingExpirationDate = nextExpirationDate;
		expirationDate = nextExpirationDate;
		expirationPopoverOpen = false;
	}

	function applyExpirationDate() {
		expirationDate = pendingExpirationDate;
		expirationPopoverOpen = false;
	}

	function clearExpirationDate() {
		pendingExpirationDate = undefined;
		expirationDate = undefined;
		expirationPopoverOpen = false;
	}

	async function copyCreatedAPIKey() {
		const secret = createdAPIKey?.api_key;
		if (!secret) return;

		try {
			await navigator.clipboard.writeText(secret);
			toast.success("API key copied.");
		} catch {
			toast.error("Unable to copy API key.");
		}
	}

	async function copyWebhookSecret() {
		if (!visibleWebhookSecret) return;

		try {
			await navigator.clipboard.writeText(visibleWebhookSecret);
			toast.success("Webhook secret copied.");
		} catch {
			toast.error("Unable to copy webhook secret.");
		}
	}

	function submitCreate(event: SubmitEvent) {
		event.preventDefault();
		if (!trimmedName || createAPIKeyMutation.isPending) return;

		createAPIKeyMutation.reset();
		createAPIKeyMutation.mutate({
			name: trimmedName,
			...(expirationDate
				? { expires_at: apiKeyExpirationDateToISOString(expirationDate) }
				: {})
		});
	}

	function submitWebhook(event: SubmitEvent) {
		event.preventDefault();
		if (!trimmedWebhookURL || webhookFormDisabled) return;

		saveWebhookMutation.reset();
		regenerateWebhookSecretMutation.reset();
		deleteWebhookMutation.reset();
		saveWebhookMutation.mutate();
	}

	async function runDelete(id: string) {
		try {
			await deleteMutation.mutateAsync(id);
		} catch {
			// The mutation owns the error state; the page renders it above the table.
		}
	}

	async function runWebhookDelete() {
		try {
			await deleteWebhookMutation.mutateAsync();
		} catch {
			// The mutation owns the error state; the page renders it above the form.
		}
	}

	function isWebhookEventSelected(event: WebhookEvent) {
		return webhookEvents.includes(event);
	}

	function setWebhookEventSelected(event: WebhookEvent, checked: boolean | "indeterminate") {
		if (webhookFormDisabled) return;

		const selectedEvents = new Set(webhookEvents);
		if (checked === true) {
			selectedEvents.add(event);
		} else {
			selectedEvents.delete(event);
		}
		webhookEvents = WEBHOOK_EVENT_OPTIONS.map((option) => option.value).filter((value) =>
			selectedEvents.has(value)
		);
	}

	function regenerateWebhookSecretAction() {
		if (!webhookExists || webhookFormDisabled) return;

		saveWebhookMutation.reset();
		regenerateWebhookSecretMutation.reset();
		deleteWebhookMutation.reset();
		regenerateWebhookSecretMutation.mutate();
	}

	function confirmAPIKeyDelete(key: APIKeyResponse) {
		deleteMutation.reset();
		confirmDelete({
			title: "Delete API key?",
			description: `Delete "${key.name}"? Any clients using this key will stop working.`,
			confirm: { text: "Delete" },
			onConfirm: () => runDelete(key.id)
		});
	}

	function confirmWebhookDelete() {
		if (!webhookExists || webhookFormDisabled) return;

		saveWebhookMutation.reset();
		regenerateWebhookSecretMutation.reset();
		deleteWebhookMutation.reset();
		confirmDelete({
			title: "Delete webhook?",
			description: "Delete this webhook endpoint? Syncra will stop sending job events to it.",
			confirm: { text: "Delete" },
			onConfirm: () => runWebhookDelete()
		});
	}
</script>

<svelte:head>
	<title>Developer Settings | Syncra</title>
</svelte:head>

<div class="@container/main flex flex-1 flex-col gap-6 p-4 lg:p-6">
	<div class="grid min-h-0 gap-6 lg:grid-cols-[13rem_1fr]">
		<nav class="border-b pb-2 lg:border-r lg:border-b-0 lg:pr-3" aria-label="Developer settings">
			<div class="flex gap-1 overflow-x-auto lg:grid lg:overflow-visible">
				<a
					href="/app/developer-settings"
					aria-current={activeTab === "api-keys" ? "page" : undefined}
					class="{activeTab === 'api-keys'
						? 'bg-muted text-foreground'
						: 'text-muted-foreground hover:bg-muted/60 hover:text-foreground'} flex h-9 shrink-0 items-center gap-2 rounded-md px-3 text-sm font-medium whitespace-nowrap"
				>
					<KeyRoundIcon class="size-4" aria-hidden="true" />
					<span>API Keys</span>
				</a>
				<a
					href="/app/developer-settings?tab=webhooks"
					aria-current={activeTab === "webhooks" ? "page" : undefined}
					class="{activeTab === 'webhooks'
						? 'bg-muted text-foreground'
						: 'text-muted-foreground hover:bg-muted/60 hover:text-foreground'} flex h-9 shrink-0 items-center gap-2 rounded-md px-3 text-sm font-medium whitespace-nowrap"
				>
					<WebhookIcon class="size-4" aria-hidden="true" />
					<span>Webhook Settings</span>
				</a>
			</div>
		</nav>

		<main class="min-w-0 space-y-6">
			{#if activeTab === "api-keys"}
				<section class="space-y-4">
				<div class="flex flex-col justify-between gap-3 sm:flex-row sm:items-start">
					<div>
						<h2 class="flex items-center gap-2 text-base font-semibold">
							<KeyRoundIcon class="size-4" aria-hidden="true" />
							API Keys
						</h2>
						<p class="text-muted-foreground mt-1 text-sm">
							Create and revoke API keys for command-line tools and integrations.
						</p>
					</div>
				</div>

				<form class="grid gap-4 rounded-xl border bg-background p-4 shadow-xs" onsubmit={submitCreate}>
					<div class="grid gap-3 md:grid-cols-[minmax(0,1fr)_16rem_auto] md:items-end">
						<div class="grid gap-2">
							<Label for="api-key-name">Name</Label>
							<Input
								id="api-key-name"
								bind:value={nameValue}
								placeholder="My API key name"
								maxlength={255}
								autocomplete="off"
								disabled={createAPIKeyMutation.isPending}
							/>
						</div>

						<div class="grid gap-2">
							<Label>Expiration</Label>
							<Popover.Root bind:open={() => expirationPopoverOpen, setExpirationPopoverOpen}>
								<Popover.Trigger>
									{#snippet child({ props })}
										<Button
											type="button"
											variant="outline"
											class="w-full justify-start bg-background/50"
											aria-label={`Expiration: ${expirationLabel}`}
											disabled={createAPIKeyMutation.isPending}
											{...props}
										>
											<CalendarIcon class="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
											<span class="truncate">{expirationLabel}</span>
										</Button>
									{/snippet}
								</Popover.Trigger>
								<Popover.Content align="start" class="w-auto p-0">
									<div class="flex flex-col sm:flex-row">
										<div class="flex min-w-[132px] flex-row gap-1.5 border-b border-border bg-muted/5 p-3 sm:flex-col sm:border-r sm:border-b-0">
											<span class="hidden select-none px-2 py-1 text-[10px] font-bold tracking-wider text-muted-foreground/60 uppercase sm:inline-block">Presets</span>
											<Button
												type="button"
												variant={activeExpirationPreset === "week" ? "secondary" : "ghost"}
												class="h-8.5 flex-1 justify-center px-3 text-xs font-medium sm:flex-none sm:justify-start"
												onclick={() => setExpirationPreset("week")}
											>
												1 week
											</Button>
											<Button
												type="button"
												variant={activeExpirationPreset === "month" ? "secondary" : "ghost"}
												class="h-8.5 flex-1 justify-center px-3 text-xs font-medium sm:flex-none sm:justify-start"
												onclick={() => setExpirationPreset("month")}
											>
												1 month
											</Button>
											<Button
												type="button"
												variant={activeExpirationPreset === "quarter" ? "secondary" : "ghost"}
												class="h-8.5 flex-1 justify-center px-3 text-xs font-medium sm:flex-none sm:justify-start"
												onclick={() => setExpirationPreset("quarter")}
											>
												3 months
											</Button>
										</div>
										<div class="flex flex-col">
											<Calendar.Calendar
												type="single"
												bind:value={pendingExpirationDate}
												minValue={today(getLocalTimeZone())}
											/>
											<div class="flex justify-end gap-2 border-t bg-muted/20 p-3">
												<Button type="button" variant="ghost" size="sm" onclick={clearExpirationDate}>
													Clear
												</Button>
												<Button type="button" size="sm" onclick={applyExpirationDate}>Apply</Button>
											</div>
										</div>
									</div>
								</Popover.Content>
							</Popover.Root>
						</div>

						<div>
							<Button
								type="submit"
								class="w-full shrink-0 md:w-auto"
								disabled={!trimmedName || createAPIKeyMutation.isPending}
							>
								<PlusIcon class="size-4" aria-hidden="true" />
								{createAPIKeyMutation.isPending ? "Creating..." : "Create key"}
							</Button>
						</div>
					</div>
				</form>

				{#if createdAPIKey?.api_key}
					<Alert.Root>
						<KeyRoundIcon class="size-4" />
						<Alert.Title>Copy your new API key</Alert.Title>
						<Alert.Description>
							<div class="mt-3 space-y-3">
								<code
									class="bg-muted text-foreground block rounded-md border px-3 py-2 font-mono text-xs break-all"
								>
									{createdAPIKey.api_key}
								</code>
								<div class="flex flex-wrap gap-2">
									<Button type="button" size="sm" onclick={() => void copyCreatedAPIKey()}>
										<CopyIcon class="size-4" aria-hidden="true" />
										Copy
									</Button>
									<Button type="button" size="sm" variant="outline" onclick={() => (createdAPIKey = null)}>
										<EyeOffIcon class="size-4" aria-hidden="true" />
										Hide
									</Button>
								</div>
							</div>
						</Alert.Description>
					</Alert.Root>
				{/if}

				{#if createError || deleteError}
					<Alert.Root variant="destructive">
						<AlertCircleIcon class="size-4" />
						<Alert.Description>{createError || deleteError}</Alert.Description>
					</Alert.Root>
				{/if}

				<div class="overflow-x-auto rounded-xl border bg-background shadow-xs">
					<Table.Root>
						<Table.Header class="sticky top-0 z-10 border-b bg-muted/40">
							<Table.Row class="hover:bg-transparent">
								<Table.Head class="h-10 min-w-[220px] py-2.5 text-xs font-semibold uppercase text-muted-foreground/90">
									Name
								</Table.Head>
								<Table.Head class="h-10 min-w-[260px] py-2.5 text-xs font-semibold uppercase text-muted-foreground/90">
									Key
								</Table.Head>
								<Table.Head class="h-10 w-[140px] py-2.5 text-xs font-semibold uppercase text-muted-foreground/90">
									Created
								</Table.Head>
								<Table.Head class="h-10 w-[140px] py-2.5 text-xs font-semibold uppercase text-muted-foreground/90">
									Expires
								</Table.Head>
								<Table.Head class="h-10 w-[100px] py-2.5 text-right text-xs font-semibold uppercase text-muted-foreground/90">
									Actions
								</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#if apiKeysQuery.isLoading}
								<Table.Row>
									<Table.Cell colspan={API_KEY_COLUMN_COUNT} class="h-48 p-0">
										<div class="flex h-48 w-full items-center justify-center">
											<Spinner class="size-16 text-foreground dark:text-blue-500" />
										</div>
									</Table.Cell>
								</Table.Row>
							{:else if apiKeysQuery.isError}
								<Table.Row>
									<Table.Cell colspan={API_KEY_COLUMN_COUNT} class="h-24">
										<div class="mx-auto flex max-w-xl flex-wrap items-center justify-center gap-3 text-center text-sm text-destructive">
											<span>{apiKeysQuery.error.message}</span>
											<Button
												type="button"
												variant="outline"
												size="sm"
												onclick={() => apiKeysQuery.refetch()}
											>
												Retry
											</Button>
										</div>
									</Table.Cell>
								</Table.Row>
							{:else if apiKeys.length > 0}
								{#each apiKeys as key (key.id)}
									<Table.Row class="transition-colors duration-150 hover:bg-muted/40">
										<Table.Cell class="max-w-[320px] px-4 py-3">
											<div class="flex min-w-0 items-center gap-2">
												<KeyRoundIcon class="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
												<span class="truncate text-sm font-medium" title={key.name}>{key.name}</span>
											</div>
										</Table.Cell>
										<Table.Cell class="px-4 py-3">
											<code class="font-mono text-xs text-muted-foreground break-all">
												{maskedAPIKey(key)}
											</code>
										</Table.Cell>
										<Table.Cell class="whitespace-nowrap px-4 py-3 text-sm text-muted-foreground">
											{formatDate(key.created_at)}
										</Table.Cell>
										<Table.Cell class="whitespace-nowrap px-4 py-3 text-sm">
											{#if key.expires_at}
												{formatDate(key.expires_at)}
											{:else}
												<Badge variant="secondary">Never</Badge>
											{/if}
										</Table.Cell>
										<Table.Cell class="px-4 py-3">
											<div class="flex justify-end">
												<Button
													type="button"
													variant="ghost"
													size="icon-sm"
													disabled={deleteMutation.isPending}
													aria-label={`Delete ${key.name}`}
													onclick={() => confirmAPIKeyDelete(key)}
												>
													<Trash2Icon class="size-4" aria-hidden="true" />
												</Button>
											</div>
										</Table.Cell>
									</Table.Row>
								{/each}
							{:else}
								<Table.Row>
									<Table.Cell colspan={API_KEY_COLUMN_COUNT} class="h-48 text-center">
										<div class="mx-auto flex max-w-[340px] flex-col items-center justify-center gap-2 p-6">
											<div class="rounded-full bg-muted/50 p-3.5 text-muted-foreground/80">
												<KeyRoundIcon class="size-6" aria-hidden="true" />
											</div>
											<h3 class="mt-2 text-sm font-semibold text-foreground">No API keys</h3>
											<p class="text-muted-foreground text-sm">
												Create a key when you need programmatic access.
											</p>
										</div>
									</Table.Cell>
								</Table.Row>
							{/if}
						</Table.Body>
					</Table.Root>
				</div>
				</section>
			{:else}
				<section class="space-y-4">
					<div class="flex flex-col justify-between gap-3 sm:flex-row sm:items-start">
						<div>
							<h2 class="flex items-center gap-2 text-base font-semibold">
								<WebhookIcon class="size-4" aria-hidden="true" />
								Webhook Settings
							</h2>
							<p class="text-muted-foreground mt-1 text-sm">
								Configure a webhook endpoint for job status events.
							</p>
						</div>
					</div>

					{#if webhookQuery.isLoading}
						<div class="flex h-48 w-full items-center justify-center rounded-xl border bg-background shadow-xs">
							<Spinner class="size-16 text-foreground dark:text-blue-500" />
						</div>
					{:else}
						<form class="grid gap-5 rounded-xl border bg-background p-4 shadow-xs" onsubmit={submitWebhook}>
							<div class="grid gap-2">
								<Label for="webhook-url">URL</Label>
								<Input
									id="webhook-url"
									type="url"
									bind:value={webhookURL}
									placeholder="https://example.com/webhooks/syncra"
									autocomplete="url"
									required
									aria-invalid={Boolean(webhookError)}
									disabled={webhookFormDisabled}
								/>
							</div>

							<fieldset class="grid gap-3">
								<legend class="text-sm font-medium">Events</legend>
								<div class="grid gap-2 sm:grid-cols-3">
									{#each WEBHOOK_EVENT_OPTIONS as option (option.value)}
										<div class="flex min-w-0 items-start gap-2 rounded-md px-1 py-1.5">
											<Checkbox
												id={`webhook-event-${option.value}`}
												disabled={webhookFormDisabled}
												bind:checked={() =>
													isWebhookEventSelected(option.value),
													(checked) => setWebhookEventSelected(option.value, checked)}
											/>
											<Label
												for={`webhook-event-${option.value}`}
												class="min-w-0 cursor-pointer flex-col items-start gap-1 leading-normal"
											>
												<span>{option.label}</span>
												<code class="text-muted-foreground font-mono text-xs font-normal break-all">
													{option.value}
												</code>
											</Label>
										</div>
									{/each}
								</div>
							</fieldset>

							<div class="flex flex-col gap-2 sm:flex-row sm:flex-wrap">
								<Button
									type="submit"
									class="w-full justify-center sm:w-auto"
									disabled={!trimmedWebhookURL || webhookFormDisabled}
								>
									{#if saveWebhookMutation.isPending}
										<Spinner class="size-4" />
									{:else}
										<WebhookIcon class="size-4" aria-hidden="true" />
									{/if}
									{saveWebhookMutation.isPending ? "Saving..." : "Save webhook"}
								</Button>
								<Button
									type="button"
									variant="outline"
									class="w-full justify-center sm:w-auto"
									disabled={!webhookExists || webhookFormDisabled}
									onclick={regenerateWebhookSecretAction}
								>
									{#if regenerateWebhookSecretMutation.isPending}
										<Spinner class="size-4" />
									{:else}
										<KeyRoundIcon class="size-4" aria-hidden="true" />
									{/if}
									{regenerateWebhookSecretMutation.isPending
										? "Regenerating..."
										: "Regenerate secret"}
								</Button>
								<Button
									type="button"
									variant="destructive"
									class="w-full justify-center sm:w-auto"
									disabled={!webhookExists || webhookFormDisabled}
									onclick={confirmWebhookDelete}
								>
									{#if deleteWebhookMutation.isPending}
										<Spinner class="size-4" />
									{:else}
										<Trash2Icon class="size-4" aria-hidden="true" />
									{/if}
									{deleteWebhookMutation.isPending ? "Deleting..." : "Delete"}
								</Button>
							</div>
						</form>
					{/if}

					{#if visibleWebhookSecret}
						<Alert.Root>
							<KeyRoundIcon class="size-4" />
							<Alert.Title>Copy your webhook secret</Alert.Title>
							<Alert.Description>
								<div class="mt-3 space-y-3">
									<code
										class="bg-muted text-foreground block rounded-md border px-3 py-2 font-mono text-xs break-all"
									>
										{visibleWebhookSecret}
									</code>
									<div class="flex flex-wrap gap-2">
										<Button type="button" size="sm" onclick={() => void copyWebhookSecret()}>
											<CopyIcon class="size-4" aria-hidden="true" />
											Copy
										</Button>
										<Button
											type="button"
											size="sm"
											variant="outline"
											onclick={() => (visibleWebhookSecret = "")}
										>
											<EyeOffIcon class="size-4" aria-hidden="true" />
											Hide
										</Button>
									</div>
								</div>
							</Alert.Description>
						</Alert.Root>
					{/if}

					{#if webhookError}
						<Alert.Root variant="destructive">
							<AlertCircleIcon class="size-4" />
							<Alert.Description>{webhookError}</Alert.Description>
						</Alert.Root>
					{/if}
				</section>
			{/if}
		</main>
	</div>
</div>
