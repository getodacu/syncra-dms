<script lang="ts">
	import DashboardIcon from "@tabler/icons-svelte/icons/dashboard";
	import InnerShadowTopIcon from "@tabler/icons-svelte/icons/inner-shadow-top";
	import SelectorIcon from "@tabler/icons-svelte/icons/selector";
	import ShieldCheckIcon from "@tabler/icons-svelte/icons/shield-check";
	import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import { cn } from "$lib/utils.js";

	type AuthUser = NonNullable<App.Locals["user"]>;
	type ActiveSpace = "app" | "admin";

	let { user, activeSpace }: { user: AuthUser | null; activeSpace: ActiveSpace } = $props();

	const isAdmin = $derived(user?.role === "admin");
	const brandHref = $derived(activeSpace === "admin" ? "/admin-portal" : "/app");
	const brandTitle = $derived(activeSpace === "admin" ? m.sidebar_syncra_admin() : m.sidebar_syncra());
	const activeSpaceTitle = $derived(activeSpace === "admin" ? m.sidebar_admin_portal() : m.sidebar_user_space());
</script>

<Sidebar.Menu>
	<Sidebar.MenuItem>
		{#if isAdmin}
			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					{#snippet child({ props })}
						<Sidebar.MenuButton
							{...props}
							size="lg"
							class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
						>
							<InnerShadowTopIcon class="!size-5" />
							<div class="grid flex-1 text-start leading-tight">
								<span class="truncate text-base font-semibold">{activeSpaceTitle}</span>
								<span class="truncate text-xs text-muted-foreground">{m.sidebar_syncra()}</span>
							</div>
							<SelectorIcon class="ms-auto size-4 text-muted-foreground" aria-hidden="true" />
						</Sidebar.MenuButton>
					{/snippet}
				</DropdownMenu.Trigger>
				<DropdownMenu.Content class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg" side="bottom" align="start">
					<DropdownMenu.Label class="px-2 py-1.5 text-xs font-medium text-muted-foreground">
						{m.sidebar_switch_space()}
					</DropdownMenu.Label>
					<DropdownMenu.Group>
						<DropdownMenu.Item class={cn(activeSpace === "app" && "bg-accent text-accent-foreground")}>
							<a href="/app" aria-current={activeSpace === "app" ? "page" : undefined} class="flex w-full items-center gap-2">
								<DashboardIcon class="size-4" aria-hidden="true" />
								<span>{m.sidebar_user_space()}</span>
							</a>
						</DropdownMenu.Item>
						<DropdownMenu.Item class={cn(activeSpace === "admin" && "bg-accent text-accent-foreground")}>
							<a
								href="/admin-portal"
								aria-current={activeSpace === "admin" ? "page" : undefined}
								class="flex w-full items-center gap-2"
							>
								<ShieldCheckIcon class="size-4" aria-hidden="true" />
								<span>{m.sidebar_admin_portal()}</span>
							</a>
						</DropdownMenu.Item>
					</DropdownMenu.Group>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		{:else}
			<Sidebar.MenuButton class="data-[slot=sidebar-menu-button]:!p-1.5">
				{#snippet child({ props })}
					<a href={brandHref} {...props}>
						<InnerShadowTopIcon class="!size-5" />
						<span class="text-base font-semibold">{brandTitle}</span>
					</a>
				{/snippet}
			</Sidebar.MenuButton>
		{/if}
	</Sidebar.MenuItem>
</Sidebar.Menu>
