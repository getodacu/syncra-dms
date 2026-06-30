<script lang="ts">
	import { page } from "$app/state";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Separator } from "$lib/components/ui/separator/index.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import {
		CREDIT_BALANCE_QUERY_KEY,
		fetchCreditBalance,
		type CreditBalanceResponse,
	} from "$lib/client/billing";
	import { m } from "$lib/paraglide/messages.js";
	import { getLocale } from "$lib/paraglide/runtime.js";
	import CoinsIcon from "@lucide/svelte/icons/coins";
	import MoonIcon from "@lucide/svelte/icons/moon";
	import SunIcon from "@lucide/svelte/icons/sun";
	import { createQuery } from "@tanstack/svelte-query";
	import { toggleMode } from "mode-watcher";

	let {
		initialCreditBalance,
		initialCreditBalanceError
	}: {
		initialCreditBalance: CreditBalanceResponse | null;
		initialCreditBalanceError: string | null;
	} = $props();

	const creditBalanceQuery = createQuery<CreditBalanceResponse, Error>(() => ({
		queryKey: CREDIT_BALANCE_QUERY_KEY,
		queryFn: () => fetchCreditBalance(fetch),
		initialData: initialCreditBalance ?? undefined,
	}));

	const pathname = $derived(page.url.pathname as string);
	const title = $derived.by(() => {
		if (pathname === "/app") return m.nav_dashboard();
		if (pathname === "/app/schemas") return m.nav_schemas();
		if (pathname === "/app/schemas/library") return m.schemas_library();
		if (pathname === "/app/schemas/new") return m.nav_new_schema();
		if (pathname.startsWith("/app/schemas/edit/")) return m.nav_edit_schema();
		if (pathname === "/app/new-job") return m.nav_new_job();
		if (pathname === "/app/jobs") return m.nav_jobs();
		if (pathname === "/app/billing") return m.nav_billing();
		if (pathname === "/app/billing/orders") return m.nav_billing_orders();
		if (pathname === "/app/billing/credit-usage-history") return m.nav_credit_usage_history();
		if (pathname === "/app/developer-settings") return m.nav_developer_settings();
		if (pathname === "/app/datasets" || pathname.startsWith("/app/datasets/")) return m.datasets_page_title();
		if (pathname === "/app/documents") return m.documents_page_title();
		return m.documents_page_title();
	});
	const creditBalanceLabel = $derived.by(() => {
		if (!creditBalanceQuery.data) return m.header_credits_unavailable();
		return m.header_credits({ count: creditBalanceQuery.data.available_credits.toLocaleString(getLocale()) });
	});
	const creditBalanceAriaLabel = $derived.by(() => {
		const errorMessage = creditBalanceQuery.error?.message ?? initialCreditBalanceError;
		if (!creditBalanceQuery.data && errorMessage) {
			return m.header_credit_balance_unavailable({ message: errorMessage });
		}
		return creditBalanceLabel;
	});
</script>

<header
	class="flex h-(--header-height) shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-(--header-height)"
>
	<div class="flex w-full min-w-0 items-center gap-1 px-4 lg:gap-2 lg:px-6">
		<Sidebar.Trigger class="-ms-1" />
		<Separator orientation="vertical" class="mx-2 data-[orientation=vertical]:h-4" />
		<h1 class="min-w-0 truncate text-base font-medium">{title}</h1>
		<div class="ms-auto flex shrink-0 items-center gap-2">
			<Button
				href="/app/billing"
				variant="outline"
				size="sm"
				class="max-w-[8.5rem] gap-1.5 px-2.5 text-xs sm:max-w-none"
				aria-label={creditBalanceAriaLabel}
				title={creditBalanceAriaLabel}
			>
				<CoinsIcon class="size-4 text-indigo-500" aria-hidden="true" />
				<span class="truncate tabular-nums">{creditBalanceLabel}</span>
			</Button>
			<Button
				onclick={toggleMode}
				variant="ghost"
				size="icon"
				class="relative"
			>
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
