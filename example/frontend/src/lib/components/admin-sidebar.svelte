<script lang="ts">
	import FileTextIcon from "@lucide/svelte/icons/file-text";
	import FileJsonIcon from "@lucide/svelte/icons/file-json";
	import ReceiptTextIcon from "@lucide/svelte/icons/receipt-text";
	import UsersIcon from "@lucide/svelte/icons/users";
	import LogoutIcon from "@tabler/icons-svelte/icons/logout";
	import { page } from "$app/state";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import SidebarSpaceSwitcher from "./sidebar-space-switcher.svelte";
	import type { ComponentProps } from "svelte";

	type AuthUser = NonNullable<App.Locals["user"]>;

	let {
		user,
		...restProps
	}: ComponentProps<typeof Sidebar.Root> & { user: AuthUser | null } = $props();

	const pathname = $derived(page.url.pathname as string);
	const usersActive = $derived(pathname === "/admin-portal/users" || pathname.startsWith("/admin-portal/users/"));
	const invoicesActive = $derived(pathname === "/admin-portal/invoices");
	const ordersActive = $derived(pathname === "/admin-portal/orders");
	const recipesActive = $derived(pathname === "/admin-portal/json-recipes" || pathname.startsWith("/admin-portal/json-recipes/"));
</script>

<Sidebar.Root collapsible="offcanvas" {...restProps}>
	<Sidebar.Header>
		<SidebarSpaceSwitcher user={user} activeSpace="admin" />
	</Sidebar.Header>
	<Sidebar.Content>
		<Sidebar.Group>
			<Sidebar.GroupContent>
				<Sidebar.Menu>
					<Sidebar.MenuItem>
						<Sidebar.MenuButton tooltipContent={m.admin_nav_users()} isActive={usersActive}>
							{#snippet child({ props })}
								<a {...props} href="/admin-portal/users" aria-current={usersActive ? "page" : undefined}>
									<UsersIcon />
									<span>{m.admin_nav_users()}</span>
								</a>
							{/snippet}
						</Sidebar.MenuButton>
					</Sidebar.MenuItem>
					<Sidebar.MenuItem>
						<Sidebar.MenuButton tooltipContent={m.admin_nav_invoices()} isActive={invoicesActive}>
							{#snippet child({ props })}
								<a {...props} href="/admin-portal/invoices" aria-current={invoicesActive ? "page" : undefined}>
									<FileTextIcon />
									<span>{m.admin_nav_invoices()}</span>
								</a>
							{/snippet}
						</Sidebar.MenuButton>
					</Sidebar.MenuItem>
					<Sidebar.MenuItem>
						<Sidebar.MenuButton tooltipContent={m.admin_nav_json_recipes()} isActive={recipesActive}>
							{#snippet child({ props })}
								<a {...props} href="/admin-portal/json-recipes" aria-current={recipesActive ? "page" : undefined}>
									<FileJsonIcon />
									<span>{m.admin_nav_json_recipes()}</span>
								</a>
							{/snippet}
						</Sidebar.MenuButton>
					</Sidebar.MenuItem>
					<Sidebar.MenuItem>
						<Sidebar.MenuButton tooltipContent={m.admin_nav_orders()} isActive={ordersActive}>
							{#snippet child({ props })}
								<a {...props} href="/admin-portal/orders" aria-current={ordersActive ? "page" : undefined}>
									<ReceiptTextIcon />
									<span>{m.admin_nav_orders()}</span>
								</a>
							{/snippet}
						</Sidebar.MenuButton>
					</Sidebar.MenuItem>
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>
	</Sidebar.Content>
	<Sidebar.Footer>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<div class="flex min-w-0 flex-col px-2 py-1.5 text-xs">
					<span class="truncate font-medium">{user?.name || m.admin_user_fallback()}</span>
					<span class="truncate text-muted-foreground">{user?.email || ""}</span>
				</div>
			</Sidebar.MenuItem>
			<Sidebar.MenuItem>
				<form action="/logout" method="POST">
					<Button type="submit" variant="ghost" class="h-8 w-full justify-start px-2 text-muted-foreground">
						<LogoutIcon class="size-4" aria-hidden="true" />
						<span>{m.nav_log_out()}</span>
					</Button>
				</form>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Footer>
</Sidebar.Root>
