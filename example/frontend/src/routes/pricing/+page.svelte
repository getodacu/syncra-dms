<script lang="ts">
	import { toggleMode } from "mode-watcher";
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import CheckIcon from "@lucide/svelte/icons/check";
	import CoinsIcon from "@lucide/svelte/icons/coins";
	import CreditCardIcon from "@lucide/svelte/icons/credit-card";
	import LayersIcon from "@lucide/svelte/icons/layers";
	import MoonIcon from "@lucide/svelte/icons/moon";
	import ShieldCheckIcon from "@lucide/svelte/icons/shield-check";
	import SunIcon from "@lucide/svelte/icons/sun";

	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import { Separator } from "$lib/components/ui/separator/index.js";
	import { CREDIT_PRICING_TIERS, CREDIT_RULES, checkoutHref } from "./pricing-data";

	let { data } = $props();
	const isLoggedIn = $derived(data.isLoggedIn);
	const ctaHref = $derived(checkoutHref(isLoggedIn));

	const rules = [
		{
			icon: CoinsIcon,
			title: "Page-based credits",
			description: CREDIT_RULES.creditConversion
		},
		{
			icon: LayersIcon,
			title: "Fixed purchase blocks",
			description: CREDIT_RULES.purchaseBlocks
		},
		{
			icon: CreditCardIcon,
			title: "Stripe checkout",
			description: CREDIT_RULES.noSubscriptions
		},
		{
			icon: ShieldCheckIcon,
			title: "Successful jobs only",
			description: CREDIT_RULES.successfulProcessing
		}
	];
</script>

<svelte:head>
	<title>Pricing | Syncra</title>
	<meta
		name="description"
		content="Syncra credit pricing for document processing. Buy credits in 1000-credit blocks with no monthly subscriptions."
	/>
</svelte:head>

<div class="min-h-screen bg-background text-foreground">
	<nav class="sticky top-0 z-50 border-b border-border bg-background/95 backdrop-blur-sm">
		<div class="mx-auto flex h-14 max-w-6xl items-center justify-between px-4">
			<a
				href="/"
				class="flex items-center gap-2 text-lg font-bold tracking-tight transition-opacity hover:opacity-90"
			>
				<span class="rounded-md bg-foreground px-1.5 py-0.5 text-sm font-extrabold text-background uppercase">
					S
				</span>
				<span>Syncra</span>
			</a>

			<div class="flex items-center gap-3">
				<a
					href="/ocr-recipes"
					class="hidden text-sm font-medium text-muted-foreground transition-colors hover:text-foreground sm:inline"
				>
					OCR Recipes
				</a>
				{#if isLoggedIn}
					<Button href="/app/billing" variant="ghost" size="sm">Billing</Button>
				{:else}
					<a
						href="/login"
						class="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
					>
						Log in
					</a>
					<Button href="/signup" size="sm">Sign up</Button>
				{/if}

				<Separator orientation="vertical" class="h-4" />

				<Button
					onclick={toggleMode}
					variant="ghost"
					size="icon"
					class="relative size-8"
					aria-label="Toggle theme"
				>
					<SunIcon
						class="h-[1.1rem] w-[1.1rem] scale-100 rotate-0 transition-all dark:scale-0 dark:-rotate-90"
					/>
					<MoonIcon
						class="absolute h-[1.1rem] w-[1.1rem] scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0"
					/>
				</Button>
			</div>
		</div>
	</nav>

	<main class="mx-auto flex w-full max-w-6xl flex-col gap-12 px-4 py-10 md:py-14">
		<section class="grid gap-8 lg:grid-cols-[1.1fr_0.9fr] lg:items-end">
			<div class="flex max-w-3xl flex-col gap-5">
				<Badge
					variant="secondary"
					class="w-fit rounded-full bg-muted/70 px-3 py-3 text-xs font-semibold tracking-wide uppercase"
				>
					Credit-only billing
				</Badge>
				<div class="flex flex-col gap-3">
					<h1 class="text-4xl font-extrabold leading-tight tracking-tight md:text-5xl">
						Buy credits when you need them
					</h1>
					<p class="max-w-2xl text-base leading-7 text-muted-foreground md:text-lg">
						One page consumes one credit. Credits are purchased through Stripe in 1000-credit
						blocks and are debited only after successful job processing.
					</p>
				</div>
				<div class="flex flex-wrap items-center gap-3">
					<Button href={ctaHref} class="gap-2">
						{isLoggedIn ? "Buy credits" : "Start free"}
						<ArrowRightIcon class="size-4" />
					</Button>
					{#if !isLoggedIn}
						<Button href="/login" variant="outline">Log in</Button>
					{/if}
				</div>
			</div>

			<div class="grid grid-cols-2 gap-3 rounded-[8px] border border-border bg-muted/20 p-4">
				<div>
					<p class="text-sm font-medium text-muted-foreground">Signup credits</p>
					<p class="mt-1 text-2xl font-bold">500</p>
				</div>
				<div>
					<p class="text-sm font-medium text-muted-foreground">Minimum purchase</p>
					<p class="mt-1 text-2xl font-bold">1,000</p>
				</div>
				<div>
					<p class="text-sm font-medium text-muted-foreground">Currency</p>
					<p class="mt-1 text-2xl font-bold">EUR</p>
				</div>
				<div>
					<p class="text-sm font-medium text-muted-foreground">Subscription</p>
					<p class="mt-1 text-2xl font-bold">None</p>
				</div>
			</div>
		</section>

		<section class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
			{#each CREDIT_PRICING_TIERS as tier (tier.id)}
				<Card.Root class="rounded-[8px] border-border">
					<Card.Header class="gap-3">
						<div class="flex items-start justify-between gap-3">
							<div>
								<Card.Title class="text-base">{tier.name}</Card.Title>
								<Card.Description>{tier.creditRange}</Card.Description>
							</div>
							<div class="flex size-9 items-center justify-center rounded-[8px] border border-border bg-background">
								<CoinsIcon class="size-4" />
							</div>
						</div>
						<div class="pt-2">
							<p class="text-2xl font-bold tracking-tight">{tier.unitPrice}</p>
							<p class="mt-1 text-sm text-muted-foreground">
								{tier.sampleCredits.toLocaleString()} credits: {tier.sampleTotal}
							</p>
						</div>
					</Card.Header>
					<Card.Content>
						<ul class="space-y-2 text-sm">
							<li class="flex items-center gap-2">
								<CheckIcon class="size-4 text-foreground" />
								<span>Purchased credits never expire</span>
							</li>
							<li class="flex items-center gap-2">
								<CheckIcon class="size-4 text-foreground" />
								<span>No recurring plan charge</span>
							</li>
						</ul>
					</Card.Content>
				</Card.Root>
			{/each}
		</section>

		<section class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
			{#each rules as rule (rule.title)}
				<div class="flex gap-3 rounded-[8px] border border-border bg-background p-4">
					<div class="flex size-9 shrink-0 items-center justify-center rounded-[8px] bg-muted">
						<rule.icon class="size-4" />
					</div>
					<div class="min-w-0">
						<h2 class="text-sm font-semibold">{rule.title}</h2>
						<p class="mt-1 text-sm leading-6 text-muted-foreground">{rule.description}</p>
					</div>
				</div>
			{/each}
		</section>
	</main>
</div>
