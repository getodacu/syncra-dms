<script lang="ts">
	import ArrowDownAZIcon from "@lucide/svelte/icons/arrow-down-a-z";
	import ChevronLeftIcon from "@lucide/svelte/icons/chevron-left";
	import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
	import SearchIcon from "@lucide/svelte/icons/search";
	import ShieldIcon from "@lucide/svelte/icons/shield";
	import UserIcon from "@lucide/svelte/icons/user";
	import { createQuery } from "@tanstack/svelte-query";

	import * as Alert from "$lib/components/ui/alert/index.js";
	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import * as Table from "$lib/components/ui/table/index.js";
	import {
		fetchAdminUsers,
		formatAdminDate,
		type AdminSortDirection,
		type AdminUserResponse,
		type AdminUserSort
	} from "./api";

	const PAGE_SIZE_OPTIONS = [10, 20, 50, 100];

	let pendingSearch = $state("");
	let search = $state("");
	let sort = $state<AdminUserSort>("created_at");
	let direction = $state<AdminSortDirection>("desc");
	let pageSize = $state(20);
	let currentCursor = $state<string | null>(null);
	let cursorHistory = $state<string[]>([]);

	const usersQuery = createQuery(() => ({
		queryKey: ["admin-users", { search, sort, direction, pageSize, currentCursor }],
		queryFn: () =>
			fetchAdminUsers(fetch, {
				search,
				sort,
				direction,
				size: pageSize,
				cursor: currentCursor
			})
	}));

	const users = $derived(usersQuery.data?.users ?? []);
	const nextCursor = $derived(usersQuery.data?.next_cursor ?? null);

	function resetPagination() {
		currentCursor = null;
		cursorHistory = [];
	}

	function applySearch() {
		search = pendingSearch.trim();
		resetPagination();
	}

	function setSort(value: Event) {
		const target = value.currentTarget as HTMLSelectElement;
		if (target.value === "created_at" || target.value === "last_login_at") {
			sort = target.value;
			resetPagination();
		}
	}

	function toggleDirection() {
		direction = direction === "desc" ? "asc" : "desc";
		resetPagination();
	}

	function setPageSize(value: Event) {
		const target = value.currentTarget as HTMLSelectElement;
		const next = Number(target.value);
		if (PAGE_SIZE_OPTIONS.includes(next)) {
			pageSize = next;
			resetPagination();
		}
	}

	function goNext() {
		if (!nextCursor) return;
		if (currentCursor) cursorHistory = [...cursorHistory, currentCursor];
		else cursorHistory = [...cursorHistory, ""];
		currentCursor = nextCursor;
	}

	function goPrevious() {
		const previous = cursorHistory.at(-1);
		if (previous === undefined) return;
		cursorHistory = cursorHistory.slice(0, -1);
		currentCursor = previous || null;
	}

	function roleClass(role: AdminUserResponse["role"]) {
		if (role === "admin") {
			return "border-blue-500/20 bg-blue-500/10 text-blue-700 dark:text-blue-400";
		}
		return "border-muted-foreground/20 bg-muted/50 text-muted-foreground";
	}
</script>

<svelte:head>
	<title>Users | Syncra Admin</title>
</svelte:head>

<div class="@container/main flex flex-1 flex-col gap-4 p-4 lg:p-6">
	<div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
		<form class="flex min-w-0 flex-1 gap-2" onsubmit={(event) => { event.preventDefault(); applySearch(); }}>
			<div class="relative min-w-0 flex-1 lg:max-w-md">
				<SearchIcon class="pointer-events-none absolute left-2.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" aria-hidden="true" />
				<Input
					class="h-9 pl-8"
					placeholder="Search name or email"
					value={pendingSearch}
					oninput={(event) => (pendingSearch = event.currentTarget.value)}
				/>
			</div>
			<Button type="submit" size="sm" class="h-9 px-4">Search</Button>
		</form>
		<div class="flex flex-wrap items-center gap-2">
			<select
				class="h-9 rounded-md border border-input bg-background px-3 text-sm shadow-xs outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50"
				value={sort}
				onchange={setSort}
				aria-label="Sort users"
			>
				<option value="created_at">Created date</option>
				<option value="last_login_at">Last login</option>
			</select>
			<Button type="button" variant="outline" size="sm" class="h-9" onclick={toggleDirection}>
				<ArrowDownAZIcon class="size-4" aria-hidden="true" />
				{direction === "desc" ? "Newest" : "Oldest"}
			</Button>
		</div>
	</div>

	{#if usersQuery.isError}
		<Alert.Root variant="destructive">
			<Alert.Description>{usersQuery.error.message}</Alert.Description>
		</Alert.Root>
	{/if}

	<div class="overflow-x-auto rounded-xl border bg-background shadow-xs">
		<Table.Root>
			<Table.Header class="border-b bg-muted/40">
				<Table.Row class="hover:bg-transparent">
					<Table.Head class="min-w-[260px]">User</Table.Head>
					<Table.Head class="w-[120px]">Role</Table.Head>
					<Table.Head class="w-[180px]">Created</Table.Head>
					<Table.Head class="w-[180px]">Last login</Table.Head>
					<Table.Head class="w-[90px]"></Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#if usersQuery.isLoading}
					<Table.Row>
						<Table.Cell colspan={5} class="h-56 p-0">
							<div class="flex h-56 items-center justify-center">
								<Spinner class="size-16" />
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if usersQuery.isError}
					<Table.Row>
						<Table.Cell colspan={5} class="h-24 text-center">
							<Button type="button" variant="outline" size="sm" onclick={() => usersQuery.refetch()}>Retry</Button>
						</Table.Cell>
					</Table.Row>
				{:else if users.length === 0}
					<Table.Row>
						<Table.Cell colspan={5} class="h-40 text-center text-sm text-muted-foreground">
							No users found.
						</Table.Cell>
					</Table.Row>
				{:else}
					{#each users as user (user.id)}
						<Table.Row class="transition-colors hover:bg-muted/40">
							<Table.Cell class="px-4 py-3">
								<div class="flex min-w-0 items-center gap-3">
									<div class="flex size-8 shrink-0 items-center justify-center rounded-md bg-muted text-muted-foreground">
										<UserIcon class="size-4" aria-hidden="true" />
									</div>
									<div class="min-w-0">
										<a href={`/admin-portal/users/${user.id}`} class="block truncate text-sm font-medium hover:underline">
											{user.name || "Unnamed user"}
										</a>
										<div class="truncate text-xs text-muted-foreground">{user.email}</div>
									</div>
								</div>
							</Table.Cell>
							<Table.Cell>
								<Badge variant="outline" class={roleClass(user.role)}>
									{#if user.role === "admin"}
										<ShieldIcon class="size-3" aria-hidden="true" />
									{/if}
									{user.role}
								</Badge>
							</Table.Cell>
							<Table.Cell class="text-sm text-muted-foreground">{formatAdminDate(user.created_at)}</Table.Cell>
							<Table.Cell class="text-sm text-muted-foreground">{formatAdminDate(user.last_login_at)}</Table.Cell>
							<Table.Cell>
								<Button href={`/admin-portal/users/${user.id}`} variant="ghost" size="sm">Open</Button>
							</Table.Cell>
						</Table.Row>
					{/each}
				{/if}
			</Table.Body>
		</Table.Root>
	</div>

	<div class="flex flex-col gap-3 px-1 text-xs text-muted-foreground sm:flex-row sm:items-center sm:justify-between">
		<div>
			{#if users.length}
				Showing <span class="font-semibold text-foreground">{users.length}</span> user{users.length === 1 ? "" : "s"}.
			{:else}
				No users to show.
			{/if}
		</div>
		<div class="flex items-center gap-2">
			<select
				class="h-8 rounded-md border border-input bg-background px-2 text-xs"
				value={pageSize}
				onchange={setPageSize}
				aria-label="Rows per page"
			>
				{#each PAGE_SIZE_OPTIONS as option (option)}
					<option value={option}>{option}</option>
				{/each}
			</select>
			<Button type="button" variant="outline" size="icon-sm" disabled={cursorHistory.length === 0 || usersQuery.isFetching} onclick={goPrevious} aria-label="Previous page">
				<ChevronLeftIcon class="size-4" aria-hidden="true" />
			</Button>
			<Button type="button" variant="outline" size="icon-sm" disabled={!nextCursor || usersQuery.isFetching} onclick={goNext} aria-label="Next page">
				<ChevronRightIcon class="size-4" aria-hidden="true" />
			</Button>
		</div>
	</div>
</div>
