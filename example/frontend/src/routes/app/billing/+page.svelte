<script lang="ts">
	import { page } from "$app/state";
	import AlertCircleIcon from "@lucide/svelte/icons/alert-circle";
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import CheckCircle2Icon from "@lucide/svelte/icons/check-circle-2";
	import CheckCircleIcon from "@lucide/svelte/icons/check-circle";
	import CheckIcon from "@lucide/svelte/icons/check";
	import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
	import ClockIcon from "@lucide/svelte/icons/clock";
	import CoinsIcon from "@lucide/svelte/icons/coins";
	import CreditCardIcon from "@lucide/svelte/icons/credit-card";
	import HistoryIcon from "@lucide/svelte/icons/history";
	import ReceiptTextIcon from "@lucide/svelte/icons/receipt-text";
	import ShieldCheckIcon from "@lucide/svelte/icons/shield-check";
	import SparklesIcon from "@lucide/svelte/icons/sparkles";
	import ZapIcon from "@lucide/svelte/icons/zap";
	import { createQuery } from "@tanstack/svelte-query";
	import { toast } from "svelte-sonner";

	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as NumberField from "$lib/components/ui/number-field/index.js";
	import {
		CREDIT_BALANCE_QUERY_KEY,
		fetchCreditBalance,
		type CreditBalanceResponse,
	} from "$lib/client/billing";
	import { m } from "$lib/paraglide/messages.js";
	import { getLocale } from "$lib/paraglide/runtime.js";
	import {
		CREDIT_PURCHASE_TIERS,
		formatCents,
		quoteCreditPurchase,
		validCreditPurchaseQuantity
	} from "$lib/billing/pricing";
	import type { PageData } from "./$types";

	let { data }: { data: PageData } = $props();

	const quickAmounts = [1000, 5000, 10000, 20000, 50000];
	let credits = $state(5000);
	let pending = $state(false);
	let errorMessage = $state("");

	const creditBalanceQuery = createQuery<CreditBalanceResponse, Error>(() => ({
		queryKey: CREDIT_BALANCE_QUERY_KEY,
		queryFn: () => fetchCreditBalance(fetch),
		initialData: data.initialCreditBalance ?? undefined,
	}));

	const checkoutStatus = $derived(page.url.searchParams.get("checkout"));
	const balanceError = $derived(
		creditBalanceQuery.error?.message ?? data.initialCreditBalanceError
	);
	
	const quote = $derived.by(() => {
		try {
			return quoteCreditPurchase(credits);
		} catch {
			return null;
		}
	});

	const discountPercent = $derived.by(() => {
		if (!quote) return 0;
		const baseRate = 1000; // Tier 1 rate: 1000 cents per 1000 credits
		const unitRate = quote.unitAmountCents;
		if (unitRate >= baseRate) return 0;
		return Math.round(((baseRate - unitRate) / baseRate) * 100);
	});

	const savingsCents = $derived.by(() => {
		if (!quote) return 0;
		const baseRate = 1000;
		const baseAmountCents = (quote.credits / 1000) * baseRate;
		return baseAmountCents - quote.amountCents;
	});

	function balanceLabel() {
		if (!creditBalanceQuery.data) return m.billing_unavailable();
		return creditBalanceQuery.data.available_credits.toLocaleString(getLocale());
	}

	function updateCredits(value: string) {
		const parsed = Number(value);
		credits = Number.isFinite(parsed) ? parsed : 0;
		errorMessage = "";
	}

	function selectAmount(amount: number) {
		credits = amount;
		errorMessage = "";
	}

	async function startCheckout(event: SubmitEvent) {
		event.preventDefault();

		if (!validCreditPurchaseQuantity(credits)) {
			errorMessage = m.billing_credit_blocks_error();
			return;
		}

		pending = true;
		errorMessage = "";

		try {
			const response = await fetch("/api/billing/checkout", {
				method: "POST",
				headers: { "content-type": "application/json" },
				body: JSON.stringify({ credits })
			});
			const body = (await response.json().catch(() => null)) as unknown;
			const checkoutUrl =
				body && typeof body === "object" && "url" in body && typeof body.url === "string"
					? body.url
					: "";
			const responseError =
				body && typeof body === "object" && "error" in body && typeof body.error === "string"
					? body.error
					: m.billing_checkout_unavailable();

			if (!response.ok || !checkoutUrl) {
				errorMessage = responseError;
				toast.error(responseError);
				return;
			}

			window.location.assign(checkoutUrl);
		} catch {
			errorMessage = m.billing_checkout_unavailable();
			toast.error(errorMessage);
		} finally {
			pending = false;
		}
	}
</script>

<svelte:head>
	<title>{m.nav_billing()} | Syncra</title>
</svelte:head>

<div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
	{#if checkoutStatus === "success"}
		<div class="flex items-start gap-3.5 rounded-xl border border-emerald-200 bg-emerald-50/50 p-4 text-emerald-900 dark:border-emerald-900/60 dark:bg-emerald-950/20 dark:text-emerald-100 shadow-xs animate-in fade-in slide-in-from-top-2 duration-300">
			<div class="flex size-8 shrink-0 items-center justify-center rounded-lg bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
				<CheckCircle2Icon class="size-4" />
			</div>
			<div>
				<p class="text-sm font-semibold">{m.billing_payment_received_title()}</p>
				<p class="mt-1 text-xs opacity-90 leading-relaxed">
					{m.billing_payment_received_body()}
				</p>
			</div>
		</div>
	{:else if checkoutStatus === "canceled"}
		<div class="flex items-start gap-3.5 rounded-xl border border-border bg-muted/40 p-4 shadow-xs">
			<div class="flex size-8 shrink-0 items-center justify-center rounded-lg bg-muted text-muted-foreground">
				<AlertCircleIcon class="size-4" />
			</div>
			<div>
				<p class="text-sm font-semibold">{m.billing_checkout_canceled_title()}</p>
				<p class="mt-1 text-xs text-muted-foreground leading-relaxed">{m.billing_checkout_canceled_body()}</p>
			</div>
		</div>
	{/if}

	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end">
		<div class="flex flex-wrap items-center gap-2">
			<Button href="/app/billing/orders" variant="outline" size="sm" class="gap-2 font-medium cursor-pointer">
				<HistoryIcon class="size-4 text-muted-foreground" aria-hidden="true" />
				<span>{m.nav_billing_orders()}</span>
			</Button>
			<Button href="/app/billing/credit-usage-history" variant="outline" size="sm" class="gap-2 font-medium cursor-pointer">
				<CreditCardIcon class="size-4 text-muted-foreground" aria-hidden="true" />
				<span>{m.nav_credit_usage_history()}</span>
			</Button>
		</div>
	</div>

	<main class="grid grid-cols-1 lg:grid-cols-3 gap-8">
		<!-- Available Balance Column (h-full flex flex-col) -->
		<section class="bg-card border border-border p-6 rounded-2xl flex flex-col h-full shadow-xs relative overflow-hidden">
			<!-- Subtle glow effect for premium feel -->
			<div class="absolute -right-12 -top-12 size-36 rounded-full bg-indigo-500/10 blur-2xl"></div>
			<div class="absolute -left-12 -bottom-12 size-36 rounded-full bg-indigo-400/5 blur-2xl"></div>

			<div class="relative flex flex-col justify-between h-full min-h-[300px] gap-8">
				<div>
					<h2 class="text-muted-foreground text-sm uppercase tracking-wider mb-2 flex items-center gap-1.5 font-medium">
						{m.billing_available_balance()}
					</h2>
					<div class="text-5xl font-bold tracking-tight tabular-nums text-foreground mb-8">
						<CoinsIcon class="size-7 text-indigo-500 inline-block" />
						{balanceLabel()}
					</div>
				</div>

				{#if !creditBalanceQuery.data && balanceError}
					<p class="text-xs text-destructive-foreground bg-destructive/20 border border-destructive/30 rounded-xl px-3 py-2">
						{balanceError}
					</p>
				{/if}

				<div class="space-y-4 mt-auto">
					<div class="flex justify-between items-center text-sm p-3 bg-muted/40 dark:bg-zinc-950/45 rounded-lg border border-border/60">
						<span class="text-muted-foreground">{m.billing_conversion()}</span>
						<span class="font-semibold text-foreground">{m.billing_conversion_rate()}</span>
					</div>
					
					<div class="space-y-3 pt-4 border-t border-border">
						{#each [
							{ icon: ShieldCheckIcon, title: m.billing_balance_checked_upload() },
							{ icon: CheckCircleIcon, title: m.billing_debited_after_success() },
							{ icon: ZapIcon, title: m.billing_secure_stripe_checkout() }
						] as feat, i (i)}
							<div class="flex items-center gap-3">
								<feat.icon class="text-indigo-500 size-[18px]" />
								<span class="text-xs text-muted-foreground">{feat.title}</span>
							</div>
						{/each}
					</div>
				</div>
			</div>
		</section>

		<!-- Purchase Credits Column (Spans 2 columns) -->
		<section class="lg:col-span-2 bg-card border border-border p-6 rounded-2xl shadow-xs">
			<h2 class="text-xl font-semibold mb-6 text-foreground">{m.billing_purchase_credits()}</h2>
			
			<form onsubmit={startCheckout} class="grid grid-cols-1 md:grid-cols-2 gap-8">
				<!-- Left Form Panel -->
				<div class="space-y-6">
					<div class="space-y-2">
						<span class="block text-sm text-muted-foreground font-medium">{m.billing_credits_to_purchase()}</span>
						<NumberField.Root bind:value={credits} step={1000} min={1000} max={1000000}>
							<div class="flex flex-col items-center border border-border bg-muted/20 dark:bg-zinc-950/45 rounded-xl p-4">
								<div class="flex w-full items-center justify-between gap-4">
									<NumberField.Decrement variant="outline" class="rounded-full size-10" tabindex={null} />
									<span class="flex items-center gap-2 text-center text-3xl font-extrabold text-foreground tracking-tight select-none">
										<CoinsIcon class="size-7 text-indigo-500" />
										<span class="font-mono tabular-nums">
											{credits.toLocaleString()}
										</span>
									</span>
									<NumberField.Increment variant="outline" class="rounded-full size-10" tabindex={null} />
								</div>
							</div>
						</NumberField.Root>
						{#if errorMessage}
							<p class="text-xs text-destructive font-medium mt-1">{errorMessage}</p>
						{/if}
					</div>

					<div class="space-y-2">
						<span class="block text-sm text-muted-foreground font-medium">{m.billing_volume_discount_tiers()}</span>
						<div class="grid grid-cols-2 gap-2">
							{#each CREDIT_PURCHASE_TIERS as tier (tier.id)}
								{@const isSelected = quote?.tier === tier.id}
								<button 
									type="button"
									onclick={() => selectAmount(tier.minCredits)}
									class="p-3 rounded-lg border cursor-pointer transition-all duration-200 flex flex-col text-left group outline-none focus-visible:ring-2 focus-visible:ring-indigo-500/35
										{isSelected 
											? 'bg-indigo-500/10 border-indigo-500 dark:bg-indigo-900/20 dark:border-indigo-500 shadow-xs' 
											: 'bg-muted/30 dark:bg-zinc-950/40 border-border hover:border-indigo-500/40'}"
								>
									<div class="text-xs font-bold text-foreground group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">
										{tier.label}
									</div>
									<div class="text-[10px] text-muted-foreground">
										{formatCents(tier.unitAmountCents)}/k
									</div>
								</button>
							{/each}
						</div>
					</div>
				</div>

				<!-- Right Receipt Summary Panel -->
				<div class="flex flex-col justify-between bg-muted/30 dark:bg-zinc-950/45 p-6 rounded-xl border border-border min-h-[300px]">
					<div>
						<div class="text-muted-foreground text-sm mb-1 font-medium">{m.billing_total_to_pay()}</div>
						<div class="text-4xl font-bold mb-6 text-foreground tracking-tight">
							{quote ? formatCents(quote.amountCents) : "-"}
						</div>
						
						<div class="space-y-2 text-sm text-muted-foreground">
							<div class="flex justify-between">
								<span>{m.billing_base_price()}</span>
								<span>{quote ? formatCents((quote.credits / 1000) * 1000) : "-"}</span>
							</div>
							<div class="flex justify-between text-emerald-600 dark:text-emerald-400 font-medium">
								<span>{m.billing_volume_discount()}</span>
								<span>-{quote ? formatCents(savingsCents) : "-"}</span>
							</div>
						</div>
					</div>
					
					<Button 
						type="submit" 
						class="w-full mt-6 flex items-center justify-center gap-2 bg-indigo-600 hover:bg-indigo-500 dark:bg-indigo-600 dark:hover:bg-indigo-700 dark:text-white py-6 rounded-lg font-semibold transition cursor-pointer shadow-md hover:shadow-indigo-500/10 active:scale-98"
						disabled={pending || !quote}
					>
						{#if pending}
							{m.billing_starting_checkout()}
						{:else}
							{m.billing_secure_checkout()} <ChevronRightIcon class="size-[18px]" />
						{/if}
					</Button>
				</div>
			</form>
		</section>
	</main>
</div>
