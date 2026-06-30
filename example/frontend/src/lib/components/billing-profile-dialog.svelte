<script lang="ts">
	import { toast } from "svelte-sonner";
	import { slide } from "svelte/transition";

	import { Button } from "$lib/components/ui/button/index.js";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { CountrySelect } from "$lib/components/ui/country-select/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import { Switch } from "$lib/components/ui/switch/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import {
		defaultBillingProfileForm,
		formFromBillingProfile,
		normalizeBillingProfileForm,
		validateBillingProfileForm,
		type BillingProfileForm
	} from "./billing-profile-utils";

	import ReceiptTextIcon from "@lucide/svelte/icons/receipt-text";
	import UserIcon from "@lucide/svelte/icons/user";
	import Building2Icon from "@lucide/svelte/icons/building-2";
	import MapPinIcon from "@lucide/svelte/icons/map-pin";
	import AlertTriangleIcon from "@lucide/svelte/icons/alert-triangle";
	import LoaderCircleIcon from "@lucide/svelte/icons/loader-circle";

	type AuthUser = NonNullable<App.Locals["user"]>;

	type Props = {
		open: boolean;
		user: AuthUser | null;
	};

	type ApiRecord = Record<string, unknown>;

	let { open = $bindable(false), user }: Props = $props();

	let form = $state<BillingProfileForm>(defaultBillingProfileForm());
	let loaded = $state(false);
	let loading = $state(false);
	let saving = $state(false);
	let error = $state<string | null>(null);

	let wasOpen = false;
	let loadedUserId: string | null | undefined = undefined;
	let loadAttemptedUserId: string | null | undefined = undefined;
	let loadGeneration = 0;
	let lastSavedForm = $state<BillingProfileForm | null>(null);
	let lastSavedUserId: string | null | undefined = undefined;

	const errorId = "billing-profile-error";
	const loadingId = "billing-profile-loading";

	const hasUser = $derived(Boolean(user));
	const formDisabled = $derived(loading || saving || !hasUser);
	const saveDisabled = $derived(saving || loading || !loaded || !hasUser);
	const billingNameLabel = $derived(
		form.entity_type === "company" ? m.billing_profile_company_name() : m.billing_profile_full_name()
	);

	function setDialogOpen(nextOpen: boolean) {
		if (!nextOpen && saving) return;
		open = nextOpen;
	}

	$effect(() => {
		const isOpen = open;
		const currentUserId = user?.id ?? null;

		if (!isOpen) {
			if (wasOpen) {
				error = null;
				loaded = false;
				loadedUserId = undefined;
				loadAttemptedUserId = undefined;
				loadGeneration += 1;
				loading = false;
				saving = false;
				form = restoredForm(currentUserId);
			}
			wasOpen = false;
			return;
		}

		if (isOpen && (!wasOpen || loadedUserId !== currentUserId)) {
			if (loadedUserId !== undefined && loadedUserId !== currentUserId) {
				loaded = false;
				form = restoredForm(currentUserId);
			}

			wasOpen = true;
			if (!loading && !loaded && loadAttemptedUserId !== currentUserId) {
				void loadBillingProfile();
			}
			return;
		}

		if (isOpen && !loading && !loaded && loadAttemptedUserId !== currentUserId) {
			void loadBillingProfile();
			return;
		}
	});

	function isRecord(value: unknown): value is ApiRecord {
		return typeof value === "object" && value !== null && !Array.isArray(value);
	}

	function responseErrorMessage(data: unknown, fallback: string) {
		if (isRecord(data) && typeof data.error === "string" && data.error.trim()) {
			return data.error;
		}

		return fallback;
	}

	function caughtErrorMessage(caught: unknown, fallback: string) {
		return caught instanceof Error ? caught.message : fallback;
	}

	function isCurrentUser(requestUserId: string | null) {
		return requestUserId === (user?.id ?? null);
	}

	function canApplyRequest(requestUserId: string | null, requestGeneration: number) {
		return open && requestGeneration === loadGeneration && isCurrentUser(requestUserId);
	}

	async function readResponseJson(response: Response) {
		return response.json().catch(() => null) as Promise<unknown>;
	}

	function copyForm(source: BillingProfileForm): BillingProfileForm {
		return { ...source };
	}

	function restoredForm(currentUserId: string | null) {
		if (lastSavedForm && lastSavedUserId === currentUserId) {
			return copyForm(lastSavedForm);
		}

		return defaultBillingProfileForm(user ?? {});
	}

	function storeSavedForm(savedForm: BillingProfileForm, requestUserId: string | null) {
		lastSavedForm = copyForm(savedForm);
		lastSavedUserId = requestUserId;
	}

	function clearCompanyFields(targetForm: BillingProfileForm) {
		targetForm.fiscal_code = "";
		targetForm.registration_number = "";
	}

	function normalizedSavePayload() {
		const normalized = normalizeBillingProfileForm(form);
		if (normalized.entity_type === "individual") {
			clearCompanyFields(normalized);
		}

		return normalized;
	}

	async function loadBillingProfile() {
		if (loading) return;

		const requestUserId = user?.id ?? null;
		const requestGeneration = loadGeneration;
		loadAttemptedUserId = requestUserId;
		loading = true;
		error = null;

		try {
			const response = await fetch("/api/billing/profile");
			const data = await readResponseJson(response);

			if (!response.ok) {
				throw new Error(responseErrorMessage(data, m.billing_profile_load_error()));
			}

			if (!canApplyRequest(requestUserId, requestGeneration)) return;
			if (!isRecord(data) || !("profile" in data)) {
				throw new Error(m.billing_profile_load_error());
			}

			const profile = data.profile;
			let loadedForm: BillingProfileForm;
			if (profile === null) {
				loadedForm = defaultBillingProfileForm(user ?? {});
			} else if (isRecord(profile)) {
				loadedForm = formFromBillingProfile(profile);
			} else {
				throw new Error(m.billing_profile_load_error());
			}

			form = loadedForm;
			storeSavedForm(loadedForm, requestUserId);
			loaded = true;
			loadedUserId = requestUserId;
		} catch (caught) {
			if (canApplyRequest(requestUserId, requestGeneration)) {
				error = caughtErrorMessage(caught, m.billing_profile_load_error());
			}
		} finally {
			if (requestGeneration === loadGeneration) {
				loading = false;
			}
		}
	}

	async function saveBillingProfile() {
		if (saveDisabled) return;

		const normalized = normalizedSavePayload();
		const validationError = validateBillingProfileForm(normalized);
		if (validationError) {
			error = validationError;
			return;
		}

		const requestUserId = user?.id ?? null;
		const requestGeneration = loadGeneration;
		saving = true;
		error = null;

		try {
			const response = await fetch("/api/billing/profile", {
				method: "PUT",
				headers: {
					"content-type": "application/json"
				},
				body: JSON.stringify(normalized)
			});
			const data = await readResponseJson(response);

			if (!response.ok) {
				throw new Error(responseErrorMessage(data, m.billing_profile_save_error()));
			}

			if (!canApplyRequest(requestUserId, requestGeneration)) return;
			if (!isRecord(data)) {
				throw new Error(m.billing_profile_save_error());
			}

			const profile = isRecord(data.profile) ? data.profile : data;
			const savedForm = formFromBillingProfile(profile);
			form = savedForm;
			storeSavedForm(savedForm, requestUserId);
			loaded = true;
			loadedUserId = requestUserId;
			loadAttemptedUserId = requestUserId;
			toast.success(m.billing_profile_saved());
		} catch (caught) {
			if (canApplyRequest(requestUserId, requestGeneration)) {
				error = caughtErrorMessage(caught, m.billing_profile_save_error());
			}
		} finally {
			if (requestGeneration === loadGeneration) {
				saving = false;
			}
		}
	}

	function handleEntityTypeChange(checked: boolean) {
		form.entity_type = checked ? "company" : "individual";
		if (form.entity_type === "individual") {
			clearCompanyFields(form);
		}
	}

	function retryLoad() {
		if (loading) return;
		void loadBillingProfile();
	}

	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		void saveBillingProfile();
	}
</script>

<Dialog.Root bind:open={() => open, setDialogOpen}>
	<Dialog.Content class="w-full overflow-hidden p-0 sm:max-w-2xl border border-border/80 shadow-lg rounded-2xl">
		<form
			class="flex max-h-[90vh] min-h-0 flex-col bg-background"
			aria-describedby={error ? errorId : loading ? loadingId : undefined}
			onsubmit={handleSubmit}
		>
			<Dialog.Header class="border-b px-6 py-5 bg-muted/5">
				<div class="flex items-center gap-3">
					<div class="flex size-10 items-center justify-center rounded-xl bg-primary/10 text-primary ring-1 ring-primary/20">
						<ReceiptTextIcon class="size-5" />
					</div>
					<div>
						<Dialog.Title class="text-lg font-bold tracking-tight text-foreground">{m.billing_profile_title()}</Dialog.Title>
						<Dialog.Description class="text-xs text-muted-foreground">{m.billing_profile_description()}</Dialog.Description>
					</div>
				</div>
			</Dialog.Header>

			<div class="min-h-0 overflow-y-auto px-6 py-6">
				<div class="flex flex-col gap-6">
					{#if error && loaded}
						<div
							id={errorId}
							role="alert"
							class="flex gap-3 rounded-xl border border-destructive/20 bg-destructive/5 p-4 text-sm text-destructive"
						>
							<AlertTriangleIcon class="size-5 shrink-0 text-destructive mt-0.5" />
							<div class="flex-1 space-y-1.5">
								<p class="font-medium leading-none">{m.billing_profile_error_title()}</p>
								<p class="text-xs text-destructive/80 leading-normal">{error}</p>
							</div>
						</div>
					{/if}

					{#if loading && !loaded}
						<div id={loadingId} class="flex flex-col items-center justify-center py-16 text-center" role="status" aria-live="polite">
							<LoaderCircleIcon class="size-8 animate-spin text-primary mb-3" />
							<p class="text-sm font-medium text-foreground">{m.billing_profile_loading()}</p>
							<p class="text-xs text-muted-foreground mt-1">{m.billing_profile_loading_body()}</p>
						</div>
					{:else if !loaded && error}
						<div
							id={errorId}
							role="alert"
							class="flex flex-col items-center justify-center py-16 text-center"
						>
							<AlertTriangleIcon class="size-10 text-destructive mb-3" />
							<p class="text-sm font-semibold text-foreground">{m.billing_profile_failed_load()}</p>
							<p class="text-xs text-muted-foreground mt-1 max-w-sm">{error}</p>
							<Button
								type="button"
								variant="outline"
								class="mt-4 h-9 shadow-xs"
								disabled={!hasUser}
								onclick={retryLoad}
							>
								{m.billing_profile_retry_loading()}
							</Button>
						</div>
					{:else}
						<!-- Profile Type Section -->
						<div class="flex flex-col gap-2 rounded-xl border border-border/80 bg-muted/15 p-4">
							<div class="flex items-center justify-between">
								<div class="space-y-1">
									<span class="text-xs font-bold uppercase tracking-wider text-muted-foreground/90">{m.billing_profile_billing_entity()}</span>
									<p class="text-[11px] text-muted-foreground leading-normal">{m.billing_profile_entity_description()}</p>
								</div>
								<div class="flex items-center gap-3">
									<span class="text-xs font-semibold {form.entity_type === 'individual' ? 'text-foreground font-bold' : 'text-muted-foreground'} transition-colors duration-200">{m.billing_profile_individual()}</span>
									<Switch
										id="entity-type-switch"
										checked={form.entity_type === "company"}
										onCheckedChange={handleEntityTypeChange}
										disabled={formDisabled}
										class="cursor-pointer"
									/>
									<span class="text-xs font-semibold {form.entity_type === 'company' ? 'text-foreground font-bold' : 'text-muted-foreground'} transition-colors duration-200">{m.billing_profile_company()}</span>
								</div>
							</div>
						</div>

						<!-- Contact Details Section -->
						<div class="space-y-4">
							<div class="flex items-center gap-2 border-b border-border/60 pb-1.5">
								<UserIcon class="size-3.5 text-muted-foreground" />
								<span class="text-xs font-bold uppercase tracking-wider text-muted-foreground/90">{m.billing_profile_general_details()}</span>
							</div>
							<div class="grid gap-4 sm:grid-cols-2">
								<div class="grid gap-2">
									<Label for="billing-name" class="text-xs font-semibold text-foreground/90">{billingNameLabel}</Label>
									<Input
										id="billing-name"
										bind:value={form.billing_name}
										disabled={formDisabled}
										maxlength={255}
										autocomplete={form.entity_type === "company" ? "organization" : "name"}
										class="h-9 rounded-lg text-sm focus-visible:ring-primary/20"
									/>
								</div>
								<div class="grid gap-2">
									<Label for="billing-email" class="text-xs font-semibold text-foreground/90">{m.billing_profile_billing_email()}</Label>
									<Input
										id="billing-email"
										type="email"
										bind:value={form.billing_email}
										disabled={formDisabled}
										maxlength={320}
										autocomplete="email"
										class="h-9 rounded-lg text-sm focus-visible:ring-primary/20"
									/>
								</div>
							</div>
						</div>

						<!-- Address Section -->
						<div class="space-y-4">
							<div class="flex items-center gap-2 border-b border-border/60 pb-1.5">
								<MapPinIcon class="size-3.5 text-muted-foreground" />
								<span class="text-xs font-bold uppercase tracking-wider text-muted-foreground/90">{m.billing_profile_billing_address()}</span>
							</div>
							<div class="grid gap-4 sm:grid-cols-2">
								<div class="grid gap-2 sm:col-span-2">
									<Label for="billing-address-line1" class="text-xs font-semibold text-foreground/90">{m.billing_profile_address_line1()}</Label>
									<Input
										id="billing-address-line1"
										bind:value={form.address_line1}
										disabled={formDisabled}
										maxlength={255}
										autocomplete="address-line1"
										class="h-9 rounded-lg text-sm focus-visible:ring-primary/20"
									/>
								</div>
								<div class="grid gap-2 sm:col-span-2">
									<Label for="billing-address-line2" class="text-xs font-semibold text-foreground/90">{m.billing_profile_address_line2()}</Label>
									<Input
										id="billing-address-line2"
										bind:value={form.address_line2}
										disabled={formDisabled}
										maxlength={255}
										autocomplete="address-line2"
										class="h-9 rounded-lg text-sm focus-visible:ring-primary/20"
									/>
								</div>
								<div class="grid gap-2">
									<Label for="billing-city" class="text-xs font-semibold text-foreground/90">{m.billing_profile_city()}</Label>
									<Input
										id="billing-city"
										bind:value={form.city}
										disabled={formDisabled}
										maxlength={160}
										autocomplete="address-level2"
										class="h-9 rounded-lg text-sm focus-visible:ring-primary/20"
									/>
								</div>
								<div class="grid gap-2">
									<Label for="billing-region" class="text-xs font-semibold text-foreground/90">{m.billing_profile_region_state()}</Label>
									<Input
										id="billing-region"
										bind:value={form.region}
										disabled={formDisabled}
										maxlength={160}
										autocomplete="address-level1"
										class="h-9 rounded-lg text-sm focus-visible:ring-primary/20"
									/>
								</div>
								<div class="grid gap-2">
									<Label for="billing-country-code" class="text-xs font-semibold text-foreground/90">{m.billing_profile_country()}</Label>
									<CountrySelect
										id="billing-country-code"
										bind:value={form.country_code}
										disabled={formDisabled}
										class="h-9 rounded-lg text-sm focus-visible:ring-primary/20"
									/>
								</div>
								<div class="grid gap-2">
									<Label for="billing-postal-code" class="text-xs font-semibold text-foreground/90">{m.billing_profile_postal_code()}</Label>
									<Input
										id="billing-postal-code"
										bind:value={form.postal_code}
										disabled={formDisabled}
										maxlength={40}
										autocomplete="postal-code"
										class="h-9 rounded-lg text-sm focus-visible:ring-primary/20"
									/>
								</div>
							</div>
						</div>

						<!-- Company Details Section (Smooth expansion) -->
						{#if form.entity_type === "company"}
							<div transition:slide={{ duration: 200 }} class="space-y-4">
								<div class="flex items-center gap-2 border-b border-border/60 pb-1.5">
									<Building2Icon class="size-3.5 text-muted-foreground" />
									<span class="text-xs font-bold uppercase tracking-wider text-muted-foreground/90">{m.billing_profile_company_details()}</span>
								</div>
								<div class="grid gap-4 sm:grid-cols-2">
									<div class="grid gap-2">
										<Label for="billing-fiscal-code" class="text-xs font-semibold text-foreground/90">{m.billing_profile_fiscal_code()}</Label>
										<Input
											id="billing-fiscal-code"
											bind:value={form.fiscal_code}
											disabled={formDisabled}
											maxlength={80}
											class="h-9 rounded-lg text-sm focus-visible:ring-primary/20"
										/>
									</div>
									<div class="grid gap-2">
										<Label for="billing-registration-number" class="text-xs font-semibold text-foreground/90">{m.billing_profile_registration_number()}</Label>
										<Input
											id="billing-registration-number"
											bind:value={form.registration_number}
											disabled={formDisabled}
											maxlength={120}
											class="h-9 rounded-lg text-sm focus-visible:ring-primary/20"
										/>
									</div>
								</div>
							</div>
						{/if}
					{/if}
				</div>
			</div>

			<Dialog.Footer class="border-t px-6 py-4 bg-muted/5 flex items-center justify-end">
				<Button type="submit" disabled={saveDisabled} class="shadow-sm cursor-pointer min-w-[160px] h-10 transition-all">
					{#if saving}
						<span class="mr-2 animate-spin size-4 border-2 border-primary-foreground border-t-transparent rounded-full"></span>
						{m.common_saving()}
					{:else}
						{m.billing_profile_save_button()}
					{/if}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
