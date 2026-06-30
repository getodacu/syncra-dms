<script lang="ts">
	import { toggleMode } from "mode-watcher";
	import MoonIcon from "@lucide/svelte/icons/moon";
	import SunIcon from "@lucide/svelte/icons/sun";

	import OcrRecipeLibrary from "$lib/ocr-recipes/ocr-recipe-library.svelte";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Separator } from "$lib/components/ui/separator/index.js";
	import { m } from "$lib/paraglide/messages.js";

	type PageData = import("./$types").PageData;

	let { data }: { data: PageData } = $props();

	const isLoggedIn = $derived(data.isLoggedIn && Boolean(data.userId));
</script>

<svelte:head>
	<title>{m.ocr_recipes_title()} | Syncra</title>
	<meta name="description" content={m.ocr_recipes_meta_description()} />
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

			<div class="hidden items-center gap-6 text-sm font-medium text-muted-foreground md:flex">
				<a href="/ocr-recipes" class="text-foreground">{m.ocr_recipes_nav()}</a>
				<a href="/pricing" class="transition-colors hover:text-foreground">Pricing</a>
				<a href="/apidoc" class="transition-colors hover:text-foreground">API Docs</a>
			</div>

			<div class="flex items-center gap-3">
				{#if isLoggedIn}
					<Button href="/app" variant="ghost" size="sm">{m.nav_dashboard()}</Button>
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
					class="relative size-8 cursor-pointer"
					aria-label={m.common_toggle_theme()}
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

	<main class="mx-auto flex w-full max-w-6xl flex-col gap-6 px-4 py-8 md:py-10">
		<section class="flex flex-col gap-4 border-b border-border pb-6 md:flex-row md:items-end md:justify-between">
			<div class="max-w-2xl space-y-2">
				<p class="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
					{m.ocr_recipes_eyebrow()}
				</p>
				<h1 class="text-3xl font-semibold leading-tight tracking-tight md:text-4xl">
					{m.ocr_recipes_hero_title()}
				</h1>
				<p class="text-sm leading-6 text-muted-foreground md:text-base">
					{m.ocr_recipes_hero_description()}
				</p>
			</div>

			<p class="flex items-center gap-1 text-sm text-muted-foreground">
				<span class="font-medium text-foreground">{data.recipes.length}</span>
				<span>{m.ocr_recipes_nav()}</span>
			</p>
		</section>

		<OcrRecipeLibrary
			recipes={data.recipes}
			loadError={data.loadError}
			{isLoggedIn}
			userId={data.userId}
		/>
	</main>
</div>
