<script lang="ts">
	import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
	import CircleDollarSignIcon from "@lucide/svelte/icons/circle-dollar-sign";
	import KeyRoundIcon from "@lucide/svelte/icons/key-round";
	import LogInIcon from "@lucide/svelte/icons/log-in";
	import MinusIcon from "@lucide/svelte/icons/minus";
	import PlusIcon from "@lucide/svelte/icons/plus";
	import SaveIcon from "@lucide/svelte/icons/save";
	import ShieldIcon from "@lucide/svelte/icons/shield";
	import UserIcon from "@lucide/svelte/icons/user";
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import * as Tabs from "$lib/components/ui/tabs/index.js";
	import { createMutation, createQuery, useQueryClient } from "@tanstack/svelte-query";
	import { toast } from "svelte-sonner";
	import {
		adjustAdminUserBalance,
		fetchAdminUser,
		formatAdminDate,
		resetAdminUserPassword,
		startAdminUserImpersonation,
		updateAdminUser,
		upsertAdminUserBillingProfile,
		type AdminBillingProfileInput,
		type AdminUserDetailResponse
	} from "../api";

	type ProfileForm = AdminBillingProfileInput;

	const queryClient = useQueryClient();
	const userId = $derived(page.params.id as string);

	let loadedUserId = $state("");
	let name = $state("");
	let email = $state("");
	let password = $state("");
	let passwordConfirmation = $state("");
	let passwordError = $state<string | null>(null);
	let balanceAmount = $state("");
	let balanceError = $state<string | null>(null);
	let profileForm = $state<ProfileForm>(emptyProfileForm());

	const userQuery = createQuery<AdminUserDetailResponse, Error>(() => ({
		queryKey: ["admin-user", userId],
		queryFn: () => fetchAdminUser(fetch, userId)
	}));

	const updateUserMutation = createMutation(() => ({
		mutationFn: () => updateAdminUser(fetch, userId, { name, email }),
		onSuccess: (updated) => {
			queryClient.setQueryData<AdminUserDetailResponse>(["admin-user", userId], (current) =>
				current ? { ...current, ...updated } : ({ ...updated, billing_profile: null } as AdminUserDetailResponse)
			);
			toast.success("User saved.");
		},
		onError: (error) => toast.error(error instanceof Error ? error.message : "Failed to save user.")
	}));

	const impersonationMutation = createMutation(() => ({
		mutationFn: () => startAdminUserImpersonation(fetch, userId),
		onSuccess: () => {
			toast.success("Impersonation started.");
			goto("/app");
		},
		onError: (error) => toast.error(error instanceof Error ? error.message : "Failed to start impersonation.")
	}));

	const passwordMutation = createMutation(() => ({
		mutationFn: () => resetAdminUserPassword(fetch, userId, password),
		onSuccess: () => {
			password = "";
			passwordConfirmation = "";
			passwordError = null;
			toast.success("Password reset.");
		},
		onError: (error) => toast.error(error instanceof Error ? error.message : "Failed to reset password.")
	}));

	const balanceMutation = createMutation(() => ({
		mutationFn: (creditsDelta: number) => adjustAdminUserBalance(fetch, userId, creditsDelta),
		onSuccess: (updated) => {
			queryClient.setQueryData<AdminUserDetailResponse>(["admin-user", userId], updated);
			balanceAmount = "";
			balanceError = null;
			toast.success("Balance adjusted.");
		},
		onError: (error) => toast.error(error instanceof Error ? error.message : "Failed to adjust balance.")
	}));

	const billingMutation = createMutation(() => ({
		mutationFn: () => upsertAdminUserBillingProfile(fetch, userId, profileForm),
		onSuccess: (profile) => {
			queryClient.setQueryData<AdminUserDetailResponse>(["admin-user", userId], (current) =>
				current ? { ...current, billing_profile: profile } : current
			);
			toast.success("Billing profile saved.");
		},
		onError: (error) => toast.error(error instanceof Error ? error.message : "Failed to save billing profile.")
	}));

	$effect(() => {
		const user = userQuery.data;
		if (!user || loadedUserId === user.id) return;
		loadedUserId = user.id;
		name = user.name;
		email = user.email;
		balanceAmount = "";
		balanceError = null;
		profileForm = profileFromUser(user);
	});

	function emptyProfileForm(): ProfileForm {
		return {
			entity_type: "individual",
			billing_name: "",
			billing_email: "",
			country_code: "",
			address_line1: "",
			address_line2: "",
			city: "",
			region: "",
			postal_code: "",
			fiscal_code: "",
			registration_number: ""
		};
	}

	function profileFromUser(user: AdminUserDetailResponse): ProfileForm {
		const profile = user.billing_profile;
		if (!profile) {
			return {
				...emptyProfileForm(),
				billing_name: user.name,
				billing_email: user.email
			};
		}
		return {
			entity_type: profile.entity_type,
			billing_name: profile.billing_name,
			billing_email: profile.billing_email,
			country_code: profile.country_code,
			address_line1: profile.address_line1,
			address_line2: profile.address_line2 ?? "",
			city: profile.city,
			region: profile.region ?? "",
			postal_code: profile.postal_code,
			fiscal_code: profile.fiscal_code ?? "",
			registration_number: profile.registration_number ?? ""
		};
	}

	function saveUser() {
		updateUserMutation.mutate();
	}

	function startImpersonation() {
		impersonationMutation.mutate();
	}

	function resetPassword() {
		if (password.length < 8 || password.length > 128) {
			passwordError = "Password must be between 8 and 128 characters.";
			return;
		}
		if (password !== passwordConfirmation) {
			passwordError = "Passwords do not match.";
			return;
		}
		passwordError = null;
		passwordMutation.mutate();
	}

	function adjustBalance(direction: 1 | -1) {
		const amount = Number(balanceAmount);
		if (!Number.isSafeInteger(amount) || amount <= 0) {
			balanceError = "Enter a positive whole number of credits.";
			return;
		}
		balanceError = null;
		balanceMutation.mutate(amount * direction);
	}

	function saveBillingProfile() {
		billingMutation.mutate();
	}

	function formatCredits(value: number) {
		return new Intl.NumberFormat().format(value);
	}

	function roleClass(role?: string) {
		if (role === "admin") return "text-blue-700 dark:text-blue-400";
		return "text-muted-foreground";
	}
</script>

<svelte:head>
	<title>User | Syncra Admin</title>
</svelte:head>

<div class="@container/main flex flex-1 flex-col gap-5 p-4 lg:p-6">
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
		<Button href="/admin-portal/users" variant="outline" size="sm" class="w-fit">
			<ArrowLeftIcon class="size-4" aria-hidden="true" />
			Users
		</Button>
	</div>

	{#if userQuery.isLoading}
		<div class="flex h-80 items-center justify-center">
			<Spinner class="size-16" />
		</div>
	{:else if userQuery.isError}
		<Card.Root>
			<Card.Content class="flex min-h-48 flex-col items-center justify-center gap-3 text-center">
				<p class="text-sm text-destructive">{userQuery.error.message}</p>
				<Button type="button" variant="outline" size="sm" onclick={() => userQuery.refetch()}>Retry</Button>
			</Card.Content>
		</Card.Root>
	{:else if userQuery.data}
		{@const user = userQuery.data}
		<div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_22rem]">
			<div class="flex min-w-0 flex-col gap-4">
				<Card.Root>
					<Card.Header class="gap-2">
						<div class="flex min-w-0 flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
							<div class="flex min-w-0 items-center gap-3">
								<div class="flex size-10 shrink-0 items-center justify-center rounded-md bg-muted text-muted-foreground">
									<UserIcon class="size-5" aria-hidden="true" />
								</div>
								<div class="min-w-0">
									<Card.Title class="truncate text-lg">{user.name || "Unnamed user"}</Card.Title>
									<Card.Description class="truncate">{user.email}</Card.Description>
								</div>
							</div>
							<div class="flex flex-wrap items-center gap-2">
								{#if user.role === "user"}
									<Button type="button" size="sm" onclick={startImpersonation} disabled={impersonationMutation.isPending}>
										<LogInIcon class="size-4" aria-hidden="true" />
										Log in as user
									</Button>
								{/if}
								<div class="flex w-fit max-w-full items-center gap-2 rounded-md border bg-muted/40 px-3 py-2 text-sm">
									<CircleDollarSignIcon class="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
									<span class="whitespace-nowrap font-medium">{formatCredits(user.available_credits)} credits</span>
								</div>
							</div>
						</div>
					</Card.Header>
				</Card.Root>

				<Tabs.Root value="profile" class="w-full">
					<Tabs.List>
						<Tabs.Trigger value="profile">Profile</Tabs.Trigger>
						<Tabs.Trigger value="password">Password</Tabs.Trigger>
						<Tabs.Trigger value="balance">Balance</Tabs.Trigger>
						<Tabs.Trigger value="billing">Billing</Tabs.Trigger>
					</Tabs.List>

					<Tabs.Content value="profile" class="mt-4">
						<Card.Root>
							<Card.Header>
								<Card.Title class="text-base">User Info</Card.Title>
								<Card.Description>Name and email changes do not modify role.</Card.Description>
							</Card.Header>
							<Card.Content>
								<form class="grid gap-4 sm:grid-cols-2" onsubmit={(event) => { event.preventDefault(); saveUser(); }}>
									<div class="grid gap-2">
										<Label for="admin-user-name">Name</Label>
										<Input id="admin-user-name" bind:value={name} autocomplete="off" />
									</div>
									<div class="grid gap-2">
										<Label for="admin-user-email">Email</Label>
										<Input id="admin-user-email" bind:value={email} autocomplete="off" />
									</div>
									<div class="sm:col-span-2">
										<Button type="submit" disabled={updateUserMutation.isPending}>
											<SaveIcon class="size-4" aria-hidden="true" />
											Save user
										</Button>
									</div>
								</form>
							</Card.Content>
						</Card.Root>
					</Tabs.Content>

					<Tabs.Content value="password" class="mt-4">
						<Card.Root>
							<Card.Header>
								<Card.Title class="text-base">Password Reset</Card.Title>
								<Card.Description>Existing sessions for this user are invalidated after reset.</Card.Description>
							</Card.Header>
							<Card.Content>
								<form class="grid max-w-xl gap-4" onsubmit={(event) => { event.preventDefault(); resetPassword(); }}>
									<div class="grid gap-2">
										<Label for="admin-user-password">New password</Label>
										<Input id="admin-user-password" type="password" bind:value={password} autocomplete="new-password" />
									</div>
									<div class="grid gap-2">
										<Label for="admin-user-password-confirmation">Confirm password</Label>
										<Input id="admin-user-password-confirmation" type="password" bind:value={passwordConfirmation} autocomplete="new-password" />
									</div>
									{#if passwordError}
										<p class="text-sm text-destructive" role="alert">{passwordError}</p>
									{/if}
									<div>
										<Button type="submit" variant="destructive" disabled={passwordMutation.isPending}>
											<KeyRoundIcon class="size-4" aria-hidden="true" />
											Reset password
										</Button>
									</div>
								</form>
							</Card.Content>
						</Card.Root>
					</Tabs.Content>

					<Tabs.Content value="balance" class="mt-4">
						<Card.Root>
							<Card.Header>
								<Card.Title class="text-base">Balance</Card.Title>
								<Card.Description>Current balance: {formatCredits(user.available_credits)} credits.</Card.Description>
							</Card.Header>
							<Card.Content>
								<form class="grid max-w-xl gap-4" onsubmit={(event) => { event.preventDefault(); adjustBalance(1); }}>
									<div class="grid gap-2">
										<Label for="admin-user-balance-amount">Credits</Label>
										<Input
											id="admin-user-balance-amount"
											type="number"
											min="1"
											step="1"
											inputmode="numeric"
											value={balanceAmount}
											oninput={(event) => (balanceAmount = event.currentTarget.value)}
										/>
									</div>
									{#if balanceError}
										<p class="text-sm text-destructive" role="alert">{balanceError}</p>
									{/if}
									<div class="flex flex-wrap gap-2">
										<Button type="submit" disabled={balanceMutation.isPending}>
											<PlusIcon class="size-4" aria-hidden="true" />
											Add credits
										</Button>
										<Button type="button" variant="outline" disabled={balanceMutation.isPending} onclick={() => adjustBalance(-1)}>
											<MinusIcon class="size-4" aria-hidden="true" />
											Subtract credits
										</Button>
									</div>
								</form>
							</Card.Content>
						</Card.Root>
					</Tabs.Content>

					<Tabs.Content value="billing" class="mt-4">
						<Card.Root>
							<Card.Header>
								<Card.Title class="text-base">Billing Profile</Card.Title>
								<Card.Description>Billing data is stored for the selected user.</Card.Description>
							</Card.Header>
							<Card.Content>
								<form class="grid gap-4 sm:grid-cols-2" onsubmit={(event) => { event.preventDefault(); saveBillingProfile(); }}>
									<div class="grid gap-2">
										<Label for="admin-billing-entity-type">Entity type</Label>
										<select
											id="admin-billing-entity-type"
											class="h-9 rounded-md border border-input bg-background px-3 text-sm shadow-xs outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50"
											bind:value={profileForm.entity_type}
										>
											<option value="individual">Individual</option>
											<option value="company">Company</option>
										</select>
									</div>
									<div class="grid gap-2">
										<Label for="admin-billing-name">Billing name</Label>
										<Input id="admin-billing-name" bind:value={profileForm.billing_name} />
									</div>
									<div class="grid gap-2">
										<Label for="admin-billing-email">Billing email</Label>
										<Input id="admin-billing-email" bind:value={profileForm.billing_email} />
									</div>
									<div class="grid gap-2">
										<Label for="admin-billing-country">Country code</Label>
										<Input id="admin-billing-country" bind:value={profileForm.country_code} maxlength={2} />
									</div>
									<div class="grid gap-2 sm:col-span-2">
										<Label for="admin-billing-address1">Address line 1</Label>
										<Input id="admin-billing-address1" bind:value={profileForm.address_line1} />
									</div>
									<div class="grid gap-2 sm:col-span-2">
										<Label for="admin-billing-address2">Address line 2</Label>
										<Input id="admin-billing-address2" bind:value={profileForm.address_line2} />
									</div>
									<div class="grid gap-2">
										<Label for="admin-billing-city">City</Label>
										<Input id="admin-billing-city" bind:value={profileForm.city} />
									</div>
									<div class="grid gap-2">
										<Label for="admin-billing-region">Region</Label>
										<Input id="admin-billing-region" bind:value={profileForm.region} />
									</div>
									<div class="grid gap-2">
										<Label for="admin-billing-postal">Postal code</Label>
										<Input id="admin-billing-postal" bind:value={profileForm.postal_code} />
									</div>
									<div class="grid gap-2">
										<Label for="admin-billing-fiscal">Fiscal code</Label>
										<Input id="admin-billing-fiscal" bind:value={profileForm.fiscal_code} />
									</div>
									<div class="grid gap-2">
										<Label for="admin-billing-registration">Registration number</Label>
										<Input id="admin-billing-registration" bind:value={profileForm.registration_number} />
									</div>
									<div class="sm:col-span-2">
										<Button type="submit" disabled={billingMutation.isPending}>
											<SaveIcon class="size-4" aria-hidden="true" />
											Save billing profile
										</Button>
									</div>
								</form>
							</Card.Content>
						</Card.Root>
					</Tabs.Content>
				</Tabs.Root>
			</div>

			<Card.Root class="h-fit">
				<Card.Header>
					<Card.Title class="text-base">Account Metadata</Card.Title>
				</Card.Header>
				<Card.Content class="grid gap-4 text-sm">
					<div class="grid gap-1">
						<span class="text-xs font-medium uppercase text-muted-foreground">Role</span>
						<span class={roleClass(user.role)}>
							{#if user.role === "admin"}
								<ShieldIcon class="mr-1 inline size-4 align-[-2px]" aria-hidden="true" />
							{/if}
							{user.role}
						</span>
					</div>
					<div class="grid gap-1">
						<span class="text-xs font-medium uppercase text-muted-foreground">Email verified</span>
						<span>{user.email_verified ? "Yes" : "No"}</span>
					</div>
					<div class="grid gap-1">
						<span class="text-xs font-medium uppercase text-muted-foreground">Created</span>
						<span>{formatAdminDate(user.created_at)}</span>
					</div>
					<div class="grid gap-1">
						<span class="text-xs font-medium uppercase text-muted-foreground">Last login</span>
						<span>{formatAdminDate(user.last_login_at)}</span>
					</div>
					<div class="grid gap-1">
						<span class="text-xs font-medium uppercase text-muted-foreground">User ID</span>
						<span class="break-all font-mono text-xs">{user.id}</span>
					</div>
				</Card.Content>
			</Card.Root>
		</div>
	{/if}
</div>
