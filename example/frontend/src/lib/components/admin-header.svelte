<script lang="ts">
	import { page } from "$app/state";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Separator } from "$lib/components/ui/separator/index.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import MoonIcon from "@lucide/svelte/icons/moon";
	import SunIcon from "@lucide/svelte/icons/sun";
	import { toggleMode } from "mode-watcher";

	const pathname = $derived(page.url.pathname as string);
	const title = $derived.by(() => {
		if (pathname === "/admin-portal/invoices") return m.admin_nav_invoices();
		if (pathname === "/admin-portal/orders") return m.admin_nav_orders();
		if (pathname === "/admin-portal/json-recipes" || pathname.startsWith("/admin-portal/json-recipes/")) return m.admin_nav_json_recipes();
		if (pathname === "/admin-portal/users") return m.admin_nav_users();
		if (pathname.startsWith("/admin-portal/users/")) return m.admin_nav_user();
		return m.admin_nav_admin();
	});
</script>

<header
	class="flex h-(--header-height) shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-(--header-height)"
>
	<div class="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
		<Sidebar.Trigger class="-ms-1" />
		<Separator orientation="vertical" class="mx-2 data-[orientation=vertical]:h-4" />
		<h1 class="text-base font-medium">{title}</h1>
		<div class="ms-auto flex items-center gap-2">
			<Button onclick={toggleMode} variant="ghost" size="icon" class="relative">
				<SunIcon
					class="h-[1.2rem] w-[1.2rem] scale-100 rotate-0 transition-all! dark:scale-0 dark:-rotate-90"
				/>
				<MoonIcon
					class="absolute h-[1.2rem] w-[1.2rem] scale-0 rotate-90 transition-all! dark:scale-100 dark:rotate-0"
				/>
				<span class="sr-only">{m.common_toggle_theme()}</span>
			</Button>
		</div>
	</div>
</header>
