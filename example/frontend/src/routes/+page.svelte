<script lang="ts">
	import { toggleMode } from "mode-watcher";
	import { fade, slide } from "svelte/transition";
	
	// Icon imports
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import CheckIcon from "@lucide/svelte/icons/check";
	import CoinsIcon from "@lucide/svelte/icons/coins";
	import MoonIcon from "@lucide/svelte/icons/moon";
	import SunIcon from "@lucide/svelte/icons/sun";
	import FileTextIcon from "@lucide/svelte/icons/file-text";
	import DatabaseIcon from "@lucide/svelte/icons/database";
	import CodeIcon from "@lucide/svelte/icons/code";
	import CpuIcon from "@lucide/svelte/icons/cpu";
	import SparklesIcon from "@lucide/svelte/icons/sparkles";
	import ShieldCheckIcon from "@lucide/svelte/icons/shield-check";
	import TerminalIcon from "@lucide/svelte/icons/terminal";
	import LayersIcon from "@lucide/svelte/icons/layers";
	import SearchIcon from "@lucide/svelte/icons/search";
	import SettingsIcon from "@lucide/svelte/icons/settings";
	import CheckCircle2Icon from "@lucide/svelte/icons/check-circle-2";
	import ServerIcon from "@lucide/svelte/icons/server";
	import LockIcon from "@lucide/svelte/icons/lock";
	import NetworkIcon from "@lucide/svelte/icons/network";

	// UI Component imports
	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import { Separator } from "$lib/components/ui/separator/index.js";

	// Data loaded from server
	let { data } = $props();
	const isLoggedIn = $derived(data.isLoggedIn);

	// Interactive Simulator Setup
	type TabType = "invoice" | "receipt" | "w2";
	let activeTab = $state<TabType>("invoice");
	let hoveredField = $state<string | null>(null);

	const documentTemplates = {
		invoice: {
			title: "Invoice PDF",
			header: "ACME INDUSTRIAL CO.",
			sub: "INVOICE #INV-2026-0042",
			rawText: [
				{ id: "vendor_name", label: "Acme Corp", x: "12%", y: "15%", w: "45%", h: "9%" },
				{ id: "invoice_number", label: "INV-2026-0042", x: "12%", y: "28%", w: "55%", h: "9%" },
				{ id: "invoice_date", label: "June 12, 2026", x: "12%", y: "41%", w: "48%", h: "9%" },
				{ id: "total_amount", label: "$1,450.00", x: "12%", y: "68%", w: "42%", h: "12%" }
			],
			json: {
				"vendor_name": "Acme Corp",
				"invoice_number": "INV-2026-0042",
				"invoice_date": "2026-06-12",
				"total_amount": 1450.00,
				"currency": "USD"
			},
			jsonMapping: {
				"vendor_name": "vendor_name",
				"invoice_number": "invoice_number",
				"invoice_date": "invoice_date",
				"total_amount": "total_amount",
				"currency": null
			}
		},
		receipt: {
			title: "Store Receipt",
			header: "ROAST & BREW COFFEE",
			sub: "TERMINAL #03 - TXN 841",
			rawText: [
				{ id: "merchant_name", label: "Roast & Brew Coffee", x: "12%", y: "14%", w: "68%", h: "9%" },
				{ id: "timestamp", label: "09:42 AM", x: "12%", y: "26%", w: "35%", h: "9%" },
				{ id: "items", label: "2x Espresso ($8.00)", x: "12%", y: "45%", w: "60%", h: "9%" },
				{ id: "tax_amount", label: "$0.72", x: "12%", y: "58%", w: "30%", h: "9%" },
				{ id: "total_amount", label: "$8.72", x: "12%", y: "72%", w: "40%", h: "11%" }
			],
			json: {
				"merchant_name": "Roast & Brew Coffee",
				"timestamp": "09:42 AM",
				"items": [
					{ "item_name": "Espresso", "quantity": 2, "price": 4.00 }
				],
				"tax_amount": 0.72,
				"total_amount": 8.72
			},
			jsonMapping: {
				"merchant_name": "merchant_name",
				"timestamp": "timestamp",
				"items": "items",
				"tax_amount": "tax_amount",
				"total_amount": "total_amount"
			}
		},
		w2: {
			title: "Tax W-2 Form",
			header: "Form W-2 Wage Statement",
			sub: "OMB No. 1545-0008",
			rawText: [
				{ id: "employer_name", label: "TechStart Inc", x: "15%", y: "16%", w: "58%", h: "10%" },
				{ id: "employee_name", label: "Jane Doe", x: "15%", y: "32%", w: "45%", h: "10%" },
				{ id: "wages_tips_compensation", label: "$84,200.00", x: "15%", y: "54%", w: "65%", h: "10%" },
				{ id: "federal_income_tax_withheld", label: "$12,630.00", x: "15%", y: "70%", w: "65%", h: "10%" }
			],
			json: {
				"employer_name": "TechStart Inc",
				"employee_name": "Jane Doe",
				"wages_tips_compensation": 84200.00,
				"federal_income_tax_withheld": 12630.00
			},
			jsonMapping: {
				"employer_name": "employer_name",
				"employee_name": "employee_name",
				"wages_tips_compensation": "wages_tips_compensation",
				"federal_income_tax_withheld": "federal_income_tax_withheld"
			}
		}
	};

	const currentTemplate = $derived(documentTemplates[activeTab]);

	// Steps for "How it works"
	const workflowSteps = [
		{
			number: "01",
			title: "Define Target Schema",
			description: "Create custom JSON schemas for documents using our visual schema designer. Define exactly what keys, types, and validations you expect in the final output."
		},
		{
			number: "02",
			title: "Upload & OCR Ingestion",
			description: "Submit documents via our high-speed API or drop them directly into the dashboard. Syncra handles OCR text extraction and image deskewing instantly."
		},
		{
			number: "03",
			title: "AI Semantic Extraction",
			description: "Our context-aware LLMs parse the raw OCR streams, locating and structuring fields into the exact JSON format specified by your schema."
		},
		{
			number: "04",
			title: "Sync & Automate",
			description: "Export structured records to Datasets, fire webhooks to downstream databases, or fetch data dynamically through developer-friendly REST endpoints."
		}
	];

	// On-Premise Showcase State
	type OnPremTab = "compose" | "k8s" | "security";
	let activeOnPremTab = $state<OnPremTab>("compose");
</script>

<svelte:head>
	<title>Syncra | Intelligent Document Processing & OCR</title>
	<meta
		name="description"
		content="Syncra turns unstructured invoices, receipts, and forms into clean, validated JSON datasets using AI-powered OCR semantic extraction pipelines."
	/>
</svelte:head>

<div class="min-h-screen bg-background text-foreground overflow-x-hidden">
	<!-- NAVBAR -->
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

			<!-- Center navigation links -->
			<div class="hidden md:flex items-center gap-6 text-sm font-medium text-muted-foreground">
				<a href="#features" class="transition-colors hover:text-foreground">Features</a>
				<a href="#how-it-works" class="transition-colors hover:text-foreground">How it Works</a>
				<a href="#on-premise" class="transition-colors hover:text-foreground">On-Premise</a>
				<a href="/ocr-recipes" class="transition-colors hover:text-foreground">OCR Recipes</a>
				<a href="/pricing" class="transition-colors hover:text-foreground">Pricing</a>
				<a href="/apidoc" class="transition-colors hover:text-foreground">API Docs</a>
			</div>

			<div class="flex items-center gap-3">
				{#if isLoggedIn}
					<Button href="/app" variant="ghost" size="sm">Dashboard</Button>
					<Button href="/app/billing" size="sm">Billing</Button>
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

	<!-- HERO SECTION -->
	<header class="relative py-20 md:py-28 bg-[linear-gradient(to_right,#8080800a_1px,transparent_1px),linear-gradient(to_bottom,#8080800a_1px,transparent_1px)] bg-[size:14px_24px]">
		<!-- Decorative light blur behind hero content -->
		<div class="absolute inset-0 top-12 -z-10 flex items-center justify-center pointer-events-none overflow-hidden">
			<div class="w-[30rem] h-[30rem] rounded-full bg-primary/15 blur-[100px] dark:bg-primary/20"></div>
		</div>

		<div class="mx-auto max-w-6xl px-4 flex flex-col items-center text-center gap-6">
			<!-- Small update announcement badge -->
			<Badge
				variant="outline"
				class="w-fit rounded-full bg-muted/50 px-3 py-3 text-xs font-semibold tracking-wide text-primary border-primary/20 hover:bg-muted transition-all cursor-pointer"
			>
				<SparklesIcon class="size-3 text-amber-500 animate-pulse mr-1" />
				AI-Powered OCR Pipelines
			</Badge>

			<h1 class="max-w-4xl text-4xl font-extrabold tracking-tight sm:text-5xl md:text-6xl lg:text-7xl leading-[1.15]">
				Turn unstructured documents <br />
				into clean, <span class="bg-gradient-to-r from-primary to-blue-400 bg-clip-text text-transparent">structured JSON</span>
			</h1>

			<p class="max-w-2xl text-base leading-relaxed text-muted-foreground md:text-xl">
				Syncra is a developer-first OCR and intelligent document processing engine. Define custom schemas, run high-volume processing tasks, and pipe clean data into your databases instantly.
			</p>

			<div class="flex flex-wrap justify-center items-center gap-4 mt-2">
				{#if isLoggedIn}
					<Button href="/app" size="lg" class="gap-2 font-semibold shadow-md cursor-pointer">
						Go to Dashboard
						<ArrowRightIcon class="size-4" />
					</Button>
				{:else}
					<Button href="/signup" size="lg" class="gap-2 font-semibold shadow-md cursor-pointer">
						Get Started for Free
						<ArrowRightIcon class="size-4" />
					</Button>
				{/if}
				<Button href="/apidoc" variant="outline" size="lg" class="gap-2 cursor-pointer">
					<TerminalIcon class="size-4" />
					Explore API Docs
				</Button>
			</div>

			<!-- Quick stats bar -->
			<div class="grid grid-cols-2 sm:grid-cols-3 gap-6 sm:gap-12 mt-12 text-center border-y border-border/60 py-6 w-full max-w-3xl">
				<div>
					<p class="text-3xl font-extrabold text-foreground">99.8%</p>
					<p class="text-sm text-muted-foreground mt-0.5">OCR Extraction Accuracy</p>
				</div>
				<div>
					<p class="text-3xl font-extrabold text-foreground">&lt; 1.5s</p>
					<p class="text-sm text-muted-foreground mt-0.5">Average Job Run Time</p>
				</div>
				<div class="col-span-2 sm:col-span-1">
					<p class="text-3xl font-extrabold text-foreground">500 Free</p>
					<p class="text-sm text-muted-foreground mt-0.5">Credits upon signup</p>
				</div>
			</div>
		</div>
	</header>

	<!-- INTERACTIVE SIMULATOR SHOWCASE -->
	<section class="py-12 bg-muted/10 border-y border-border/50">
		<div class="mx-auto max-w-6xl px-4">
			<div class="text-center flex flex-col items-center gap-3 mb-10">
				<Badge variant="secondary" class="uppercase tracking-widest text-[10px] px-3 py-3">Interactive Demo</Badge>
				<h2 class="text-3xl font-bold tracking-tight">Watch the engine extract structured fields</h2>
				<p class="text-muted-foreground max-w-xl">
					Hover over document bounding boxes on the left or JSON keys on the right to see the semantic mapping in real-time.
				</p>
			</div>

			<!-- Simulator Container -->
			<div class="grid gap-6 lg:grid-cols-12 items-stretch max-w-5xl mx-auto">
				<!-- Tab Selectors (Left side of simulator on lg screens) -->
				<div class="lg:col-span-2 flex lg:flex-col gap-2 justify-center lg:justify-start">
					{#each Object.entries(documentTemplates) as [key, value]}
						<button
							onclick={() => {
								activeTab = key as TabType;
								hoveredField = null;
							}}
							class="flex-1 lg:flex-none text-left px-3 py-2.5 rounded-md text-xs font-semibold tracking-wide border transition-all cursor-pointer uppercase
							{activeTab === key
								? 'bg-primary text-primary-foreground border-primary shadow-xs'
								: 'bg-background border-border hover:bg-muted text-muted-foreground hover:text-foreground'}"
						>
							{value.title}
						</button>
					{/each}
				</div>

				<!-- Visual simulator split grid -->
				<div class="lg:col-span-10 grid md:grid-cols-2 gap-6">
					<!-- Simulated Scanned Document -->
					<Card.Root class="rounded-lg overflow-hidden border border-border bg-card shadow-sm flex flex-col min-h-[360px]">
						<Card.Header class="py-3 px-4 border-b border-border bg-muted/30">
							<div class="flex items-center gap-2">
								<span class="size-2 rounded-full bg-red-500 animate-pulse"></span>
								<span class="text-xs font-mono tracking-tight text-muted-foreground uppercase">{currentTemplate.title} Source</span>
							</div>
						</Card.Header>
						<Card.Content class="relative flex-1 p-6 flex flex-col justify-between bg-white dark:bg-zinc-950 font-sans text-black dark:text-zinc-100 overflow-hidden">
							<!-- Mock scanned document representation -->
							<div class="border border-zinc-200 dark:border-zinc-800 rounded-md p-4 flex-1 flex flex-col bg-zinc-50/50 dark:bg-zinc-900/50 select-none">
								<div class="flex justify-between items-start border-b border-zinc-200 dark:border-zinc-800 pb-2 mb-4 shrink-0">
									<div>
										<p class="text-[11px] font-extrabold tracking-wider text-primary">{currentTemplate.header}</p>
										<p class="text-[9px] text-zinc-400 dark:text-zinc-500 font-mono mt-0.5">{currentTemplate.sub}</p>
									</div>
									<div class="text-[9px] text-zinc-400 dark:text-zinc-500 font-mono text-right">
										<p>SYSTEM OK</p>
										<p>SCAN ID: 940E-F32</p>
									</div>
								</div>

								<!-- Document Body Container (isolated relative container to avoid header overlap) -->
								<div class="relative flex-1 w-full min-h-[220px]">
									<!-- Dotted background rows simulating other unextracted layout lines -->
									<div class="w-full h-2 bg-zinc-200 dark:bg-zinc-800/80 rounded-[2px] opacity-40"></div>
									<div class="w-2/3 h-2 bg-zinc-200 dark:bg-zinc-800/80 rounded-[2px] opacity-40 mb-3"></div>

									<div class="w-full h-1.5 bg-zinc-200 dark:bg-zinc-800/80 rounded-[2px] opacity-25"></div>
									<div class="w-4/5 h-1.5 bg-zinc-200 dark:bg-zinc-800/80 rounded-[2px] opacity-25 mb-4"></div>

									<!-- Absolute positioned OCR highlighted boxes -->
									{#each currentTemplate.rawText as box}
										<div
											role="button"
											tabindex="0"
											onmouseenter={() => hoveredField = box.id}
											onmouseleave={() => hoveredField = null}
											style="left: {box.x}; top: {box.y}; width: {box.w}; height: {box.h};"
											class="absolute rounded-[4px] border border-dashed text-[10px] font-mono flex items-center px-1.5 font-bold transition-all duration-150 cursor-pointer
											{hoveredField === box.id
												? 'border-primary bg-primary/10 dark:bg-primary/20 text-primary scale-[1.02] shadow-sm z-10'
												: 'border-zinc-300 dark:border-zinc-700 bg-background text-zinc-400 dark:text-zinc-500 hover:border-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200'}"
										>
											{box.label}
										</div>
									{/each}
								</div>
							</div>
						</Card.Content>
					</Card.Root>

					<!-- Extracted JSON Output -->
					<Card.Root class="rounded-lg overflow-hidden border border-border bg-card shadow-sm flex flex-col min-h-[360px]">
						<Card.Header class="py-3 px-4 border-b border-border bg-muted/30 flex flex-row items-center justify-between">
							<span class="text-xs font-mono tracking-tight text-muted-foreground uppercase">Extracted Dataset JSON</span>
							<Badge variant="secondary" class="bg-blue-500/10 text-blue-500 hover:bg-blue-500/20 text-[10px] border-none font-semibold">
								<CpuIcon class="size-3 mr-1" /> Schema Mapped
							</Badge>
						</Card.Header>
						<Card.Content class="flex-1 p-6 bg-zinc-950 dark:bg-black font-mono text-xs text-zinc-300 overflow-y-auto leading-relaxed flex flex-col justify-start">
							<div>
								<span class="text-zinc-500">&#x7b;</span>
								<div class="pl-4 space-y-1.5 my-1">
									{#each Object.entries(currentTemplate.json) as [key, value]}
										<!-- svelte-ignore a11y_no_static_element_interactions -->
										<div
											onmouseenter={() => hoveredField = currentTemplate.jsonMapping[key as keyof typeof currentTemplate.jsonMapping] || key}
											onmouseleave={() => hoveredField = null}
											class="py-1 px-2 rounded-md transition-all duration-150 cursor-pointer border-l-2
											{hoveredField === key || hoveredField === currentTemplate.jsonMapping[key as keyof typeof currentTemplate.jsonMapping]
												? 'bg-primary/20 dark:bg-primary/25 border-primary text-white scale-[1.01]'
												: 'hover:bg-zinc-900 border-transparent text-zinc-300'}"
										>
											<span class="text-blue-400 font-semibold">"{key}"</span>: 
											{#if typeof value === 'object'}
												<span class="text-zinc-400">[</span>
												<div class="pl-4">
													<span class="text-zinc-500">&#x7b;</span>
													<div class="pl-4">
														{#each Object.entries(value[0]) as [subKey, subVal]}
															<div>
																<span class="text-zinc-400">"{subKey}"</span>: 
																{#if typeof subVal === 'number'}
																	<span class="text-amber-400">{subVal}</span>
																{:else}
																	<span class="text-emerald-400">"{subVal}"</span>
																{/if}
															</div>
														{/each}
													</div>
													<span class="text-zinc-500">&#x7d;</span>
												</div>
												<span class="text-zinc-400">]</span>
											{:else if typeof value === 'number'}
												<span class="text-amber-400 font-semibold">{value}</span>
											{:else}
												<span class="text-emerald-400 font-semibold">"{value}"</span>
											{/if}
										</div>
									{/each}
								</div>
								<span class="text-zinc-500">&#x7d;</span>
							</div>
						</Card.Content>
					</Card.Root>
				</div>
			</div>
		</div>
	</section>

	<!-- FEATURES GRID SECTION -->
	<section id="features" class="py-20 max-w-6xl mx-auto px-4">
		<div class="text-center flex flex-col items-center gap-3 mb-16">
			<Badge variant="outline" class="w-fit rounded-full px-3 py-3 text-xs font-semibold tracking-wide border-primary/20">
				Platform Capabilities
			</Badge>
			<h2 class="text-3xl font-extrabold tracking-tight sm:text-4xl md:text-5xl">
				Engineered for precise extraction
			</h2>
			<p class="max-w-2xl text-muted-foreground text-base sm:text-lg">
				Syncra combines state-of-the-art OCR engines with semantic AI transformers to ensure your documents map exactly to your relational schemas.
			</p>
		</div>

		<div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
			<!-- Custom Schemas Feature Card -->
			<Card.Root class="rounded-lg border border-border bg-card shadow-sm hover:shadow-md transition-all duration-300 hover:-translate-y-0.5 group">
				<Card.Header class="gap-2">
					<div class="flex size-10 items-center justify-center rounded-lg bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-all duration-300">
						<CodeIcon class="size-5" />
					</div>
					<Card.Title class="text-xl font-bold mt-2">Custom JSON Schemas</Card.Title>
					<Card.Description class="text-sm leading-relaxed text-muted-foreground">
						Build strict target structures using our visual schema designer. Map strings, integers, arrays, and nested objects with built-in format validations.
					</Card.Description>
				</Card.Header>
			</Card.Root>

			<!-- OCR Extraction Pipelines Feature Card -->
			<Card.Root class="rounded-lg border border-border bg-card shadow-sm hover:shadow-md transition-all duration-300 hover:-translate-y-0.5 group">
				<Card.Header class="gap-2">
					<div class="flex size-10 items-center justify-center rounded-lg bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-all duration-300">
						<FileTextIcon class="size-5" />
					</div>
					<Card.Title class="text-xl font-bold mt-2">High-Fidelity OCR</Card.Title>
					<Card.Description class="text-sm leading-relaxed text-muted-foreground">
						Our OCR layer reads low-resolution scans, multi-page PDFs, and complex layouts (invoices, receipts, statement forms) with precision.
					</Card.Description>
				</Card.Header>
			</Card.Root>

			<!-- Datasets Feature Card -->
			<Card.Root class="rounded-lg border border-border bg-card shadow-sm hover:shadow-md transition-all duration-300 hover:-translate-y-0.5 group">
				<Card.Header class="gap-2">
					<div class="flex size-10 items-center justify-center rounded-lg bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-all duration-300">
						<DatabaseIcon class="size-5" />
					</div>
					<Card.Title class="text-xl font-bold mt-2">Structured Datasets</Card.Title>
					<Card.Description class="text-sm leading-relaxed text-muted-foreground">
						Organize processing outputs into queryable datasets. Filter, search, and export extracted data as JSON or CSV files instantly.
					</Card.Description>
				</Card.Header>
			</Card.Root>

			<!-- Developer-First API Card -->
			<Card.Root class="rounded-lg border border-border bg-card shadow-sm hover:shadow-md transition-all duration-300 hover:-translate-y-0.5 group">
				<Card.Header class="gap-2">
					<div class="flex size-10 items-center justify-center rounded-lg bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-all duration-300">
						<TerminalIcon class="size-5" />
					</div>
					<Card.Title class="text-xl font-bold mt-2">API Keys & Webhooks</Card.Title>
					<Card.Description class="text-sm leading-relaxed text-muted-foreground">
						Seamlessly integrate extraction into your codebase. Generate client API keys and configure webhook targets to notify your backend when extraction completes.
					</Card.Description>
				</Card.Header>
			</Card.Root>

			<!-- Job History & Analytics Card -->
			<Card.Root class="rounded-lg border border-border bg-card shadow-sm hover:shadow-md transition-all duration-300 hover:-translate-y-0.5 group">
				<Card.Header class="gap-2">
					<div class="flex size-10 items-center justify-center rounded-lg bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-all duration-300">
						<CpuIcon class="size-5" />
					</div>
					<Card.Title class="text-xl font-bold mt-2">Batch Processing & Logs</Card.Title>
					<Card.Description class="text-sm leading-relaxed text-muted-foreground">
						Submit large document pools concurrently. Review granular processing logs, execution durations, and step-by-step extraction audits for debugging.
					</Card.Description>
				</Card.Header>
			</Card.Root>

			<!-- Pay-Per-Success Billing Card -->
			<Card.Root class="rounded-lg border border-border bg-card shadow-sm hover:shadow-md transition-all duration-300 hover:-translate-y-0.5 group">
				<Card.Header class="gap-2">
					<div class="flex size-10 items-center justify-center rounded-lg bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-all duration-300">
						<CoinsIcon class="size-5" />
					</div>
					<Card.Title class="text-xl font-bold mt-2">Pay-as-you-Go Billing</Card.Title>
					<Card.Description class="text-sm leading-relaxed text-muted-foreground">
						Buy credits via Stripe blocks when needed. No monthly subscriptions. Best of all: credits are only deducted for successfully parsed document pages.
					</Card.Description>
				</Card.Header>
			</Card.Root>
		</div>
	</section>

	<!-- HOW IT WORKS (STEP-BY-STEP) -->
	<section id="how-it-works" class="py-20 bg-muted/20 border-y border-border/60">
		<div class="mx-auto max-w-6xl px-4">
			<div class="text-center flex flex-col items-center gap-3 mb-16">
				<Badge variant="outline" class="w-fit rounded-full px-3 py-3 text-xs font-semibold tracking-wide border-primary/20">
					Developer Workflow
				</Badge>
				<h2 class="text-3xl font-extrabold tracking-tight sm:text-4xl">
					Integrate in four simple steps
				</h2>
				<p class="max-w-xl text-muted-foreground text-sm sm:text-base">
					Syncra fits cleanly into your backend infrastructure, transforming documents asynchronously without custom parsing logic.
				</p>
			</div>

			<div class="grid gap-8 md:grid-cols-2 lg:grid-cols-4">
				{#each workflowSteps as step}
					<div class="flex flex-col gap-4 relative">
						<div class="text-5xl font-extrabold text-primary/20 dark:text-primary/10 select-none">
							{step.number}
						</div>
						<div class="flex flex-col gap-2">
							<h3 class="text-lg font-bold">{step.title}</h3>
							<p class="text-sm leading-relaxed text-muted-foreground">{step.description}</p>
						</div>
					</div>
				{/each}
			</div>

			<!-- Mock code execution preview -->
			<div class="mt-16 max-w-3xl mx-auto rounded-lg border border-border bg-zinc-950 dark:bg-black shadow-lg overflow-hidden">
				<div class="px-4 py-2 bg-zinc-900 border-b border-border/80 flex items-center justify-between">
					<div class="flex gap-1.5">
						<span class="size-3 rounded-full bg-red-500/80"></span>
						<span class="size-3 rounded-full bg-amber-500/80"></span>
						<span class="size-3 rounded-full bg-green-500/80"></span>
					</div>
					<span class="text-xs font-mono text-zinc-500">POST /api/v1/jobs</span>
				</div>
				<div class="p-5 font-mono text-[11px] sm:text-xs text-zinc-300 leading-relaxed overflow-x-auto">
					<p><span class="text-zinc-500"># Trigger a document parsing job via CURL</span></p>
					<p>
						<span class="text-blue-400">curl</span> -X POST <span class="text-emerald-400">"https://api.syncra.io/v1/jobs"</span> \
					</p>
					<p class="pl-4">
						-H <span class="text-emerald-400">"Authorization: Bearer sc_live_8f3c4d2e"</span> \
					</p>
					<p class="pl-4">
						-F <span class="text-emerald-400">"file=@invoice.pdf"</span> \
					</p>
					<p class="pl-4">
						-F <span class="text-emerald-400">"schema_id=sch_invoice_basic"</span> \
					</p>
					<p class="pl-4">
						-F <span class="text-emerald-400">"webhook_url=https://myapi.com/webhooks/syncra"</span>
					</p>
					<br />
					<p><span class="text-zinc-500"># Response (Status 202 Accepted)</span></p>
					<p class="text-zinc-400">
						&#x7b; <span class="text-blue-400">"job_id"</span>: <span class="text-emerald-400">"job_9281a4cf"</span>, <span class="text-blue-400">"status"</span>: <span class="text-emerald-400">"queued"</span>, <span class="text-blue-400">"pages"</span>: <span class="text-amber-400">1</span> &#x7d;
					</p>
				</div>
			</div>
		</div>
	</section>

	<!-- ON-PREMISE & SELF-HOSTED SECTION -->
	<section id="on-premise" class="py-20 bg-muted/10 border-y border-border/50">
		<div class="mx-auto max-w-6xl px-4">
			<div class="grid gap-12 lg:grid-cols-12 items-center">
				<!-- Left Column: Content -->
				<div class="lg:col-span-5 flex flex-col gap-6">
					<Badge variant="outline" class="w-fit rounded-full px-3 py-3 text-xs font-semibold tracking-wide border-primary/20 bg-primary/5 text-primary">
						<ServerIcon class="size-3 mr-1" />
						Enterprise Sovereignty
					</Badge>
					<h2 class="text-3xl font-extrabold tracking-tight sm:text-4xl">
						Full Sovereignty with On-Premise Deployment
					</h2>
					<p class="text-muted-foreground text-sm sm:text-base leading-relaxed">
						Keep your document data inside your security perimeter. Deploy Syncra's high-speed OCR and LLM extraction pipelines on your private infrastructure, with no external data transmission.
					</p>
					
					<div class="flex flex-col gap-4">
						<div class="flex gap-3 items-start">
							<div class="flex size-6 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary mt-0.5">
								<ShieldCheckIcon class="size-4" />
							</div>
							<div>
								<h4 class="text-sm font-bold text-foreground">100% Air-Gapped Ready</h4>
								<p class="text-xs text-muted-foreground mt-0.5">Operate entirely offline behind private VPCs. Verify license leasing via local cryptographically signed tokens.</p>
							</div>
						</div>
						
						<div class="flex gap-3 items-start">
							<div class="flex size-6 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary mt-0.5">
								<LockIcon class="size-4" />
							</div>
							<div>
								<h4 class="text-sm font-bold text-foreground">Strict Compliance & Privacy</h4>
								<p class="text-xs text-muted-foreground mt-0.5">Achieve full GDPR, HIPAA, and SOC 2 data processing compliance since document text never leaves your environment.</p>
							</div>
						</div>
						
						<div class="flex gap-3 items-start">
							<div class="flex size-6 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary mt-0.5">
								<NetworkIcon class="size-4" />
							</div>
							<div>
								<h4 class="text-sm font-bold text-foreground">High-Throughput GPU Scale</h4>
								<p class="text-xs text-muted-foreground mt-0.5">Leverage native CUDA hardware acceleration to process millions of documents at maximum throughput on your own clusters.</p>
							</div>
						</div>
					</div>

					<div class="mt-2">
						<Button href="mailto:sales@syncra.io?subject=Syncra%20On-Premise%20Inquiry" class="gap-2 cursor-pointer font-semibold shadow-xs">
							Contact Enterprise Sales
							<ArrowRightIcon class="size-4" />
						</Button>
					</div>
				</div>

				<!-- Right Column: Interactive Code & Spec Deck -->
				<div class="lg:col-span-7 flex flex-col">
					<Card.Root class="rounded-lg overflow-hidden border border-border bg-card shadow-sm flex flex-col min-h-[420px]">
						<!-- Interactive Tab Header -->
						<Card.Header class="py-3 px-4 border-b border-border bg-muted/30">
							<div class="flex flex-wrap items-center justify-between gap-4">
								<div class="flex gap-2 flex-wrap">
									<button
										onclick={() => activeOnPremTab = "compose"}
										class="px-3 py-1.5 rounded-md text-xs font-semibold border transition-all cursor-pointer uppercase tracking-wider
										{activeOnPremTab === 'compose'
											? 'bg-primary text-primary-foreground border-primary shadow-xs'
											: 'bg-background border-border text-muted-foreground hover:text-foreground'}"
									>
										docker-compose.yml
									</button>
									<button
										onclick={() => activeOnPremTab = "k8s"}
										class="px-3 py-1.5 rounded-md text-xs font-semibold border transition-all cursor-pointer uppercase tracking-wider
										{activeOnPremTab === 'k8s'
											? 'bg-primary text-primary-foreground border-primary shadow-xs'
											: 'bg-background border-border text-muted-foreground hover:text-foreground'}"
									>
										Helm Install
									</button>
									<button
										onclick={() => activeOnPremTab = "security"}
										class="px-3 py-1.5 rounded-md text-xs font-semibold border transition-all cursor-pointer uppercase tracking-wider
										{activeOnPremTab === 'security'
											? 'bg-primary text-primary-foreground border-primary shadow-xs'
											: 'bg-background border-border text-muted-foreground hover:text-foreground'}"
									>
										Security Specs
									</button>
								</div>
								<div class="flex items-center gap-1.5 text-xs text-zinc-500 font-mono">
									<span class="size-2 rounded-full bg-emerald-500"></span>
									<span>Self-Hosted Mode</span>
								</div>
							</div>
						</Card.Header>
						
						<!-- Tab Contents -->
						<Card.Content class="flex-1 p-6 bg-zinc-950 dark:bg-black font-mono text-[11px] sm:text-xs text-zinc-300 overflow-y-auto leading-relaxed flex flex-col justify-start">
							{#if activeOnPremTab === 'compose'}
								<div transition:fade={{ duration: 150 }}>
									<p class="text-zinc-500 mb-2"># Spin up Syncra OCR and Semantic LLM locally with GPU acceleration</p>
									<p><span class="text-blue-400">version</span>: <span class="text-emerald-400">"3.8"</span></p>
									<p><span class="text-blue-400">services</span>:</p>
									<p class="pl-4"><span class="text-blue-400">syncra-ocr-engine</span>:</p>
									<p class="pl-8"><span class="text-blue-400">image</span>: <span class="text-emerald-400">syncra.io/enterprise/ocr-engine:v2.4</span></p>
									<p class="pl-8"><span class="text-blue-400">ports</span>:</p>
									<p class="pl-12">- <span class="text-emerald-400">"8080:8080"</span></p>
									<p class="pl-8"><span class="text-blue-400">environment</span>:</p>
									<p class="pl-12">- <span class="text-blue-400">AIR_GAPPED</span>=<span class="text-emerald-400">"true"</span></p>
									<p class="pl-12">- <span class="text-blue-400">DB_CONNECTION</span>=<span class="text-emerald-400">"postgresql://db:5432/syncra"</span></p>
									<br />
									<p class="pl-4"><span class="text-blue-400">syncra-semantic-parser</span>:</p>
									<p class="pl-8"><span class="text-blue-400">image</span>: <span class="text-emerald-400">syncra.io/enterprise/parser-gpu-cuda:latest</span></p>
									<p class="pl-8"><span class="text-blue-400">deploy</span>:</p>
									<p class="pl-12"><span class="text-blue-400">resources</span>:</p>
									<p class="pl-16"><span class="text-blue-400">reservations</span>:</p>
									<p class="pl-20"><span class="text-blue-400">devices</span>:</p>
									<p class="pl-24">- <span class="text-blue-400">driver</span>: <span class="text-emerald-400">nvidia</span></p>
									<p class="pl-28"><span class="text-blue-400">count</span>: <span class="text-amber-400">all</span></p>
									<p class="pl-28"><span class="text-blue-400">capabilities</span>: [<span class="text-emerald-400">gpu</span>]</p>
								</div>
							{:else if activeOnPremTab === 'k8s'}
								<div transition:fade={{ duration: 150 }} class="space-y-4">
									<div>
										<p class="text-zinc-500"># Add the Syncra Enterprise chart repository</p>
										<p><span class="text-blue-400">helm</span> repo add syncra-enterprise https://charts.syncra.io/enterprise</p>
										<p><span class="text-blue-400">helm</span> repo update</p>
									</div>
									
									<div>
										<p class="text-zinc-500"># Deploy the architecture inside your secure namespace</p>
										<p><span class="text-blue-400">helm</span> install syncra-core syncra-enterprise/syncra-core \</p>
										<p class="pl-4">--namespace syncra-vpc \</p>
										<p class="pl-4">--create-namespace \</p>
										<p class="pl-4">--set global.airGapped=<span class="text-emerald-400">true</span> \</p>
										<p class="pl-4">--set global.replicaCount=<span class="text-amber-400">3</span> \</p>
										<p class="pl-4">--set auth.licenseKey=<span class="text-emerald-400">"sc_license_92c3482d"</span></p>
									</div>

									<div class="text-zinc-500 text-[10px] mt-4 border-t border-zinc-800/80 pt-3 font-sans">
										Note: Compatible with EKS, GKE, AKS, and bare-metal Kubernetes clusters (v1.26+).
									</div>
								</div>
							{:else if activeOnPremTab === 'security'}
								<div transition:fade={{ duration: 150 }} class="grid gap-4 sm:grid-cols-2 font-sans py-2">
									<div class="bg-zinc-900/40 p-4 rounded-md border border-zinc-800/80">
										<p class="text-xs font-bold text-primary flex items-center gap-1.5">
											<span class="size-1.5 rounded-full bg-primary"></span>
											Zero Outbound Traffic
										</p>
										<p class="text-[11px] text-zinc-400 mt-1 leading-relaxed">
											All operations, logic, and intelligence run local to the instance. Outbound queries to public domains are fully blocked.
										</p>
									</div>
									<div class="bg-zinc-900/40 p-4 rounded-md border border-zinc-800/80">
										<p class="text-xs font-bold text-primary flex items-center gap-1.5">
											<span class="size-1.5 rounded-full bg-primary"></span>
											Audit Log Coverage
										</p>
										<p class="text-[11px] text-zinc-400 mt-1 leading-relaxed">
											Generate unified execution audits directly to syslog, AWS CloudWatch, or Datadog, tracing every schema invocation.
										</p>
									</div>
									<div class="bg-zinc-900/40 p-4 rounded-md border border-zinc-800/80">
										<p class="text-xs font-bold text-primary flex items-center gap-1.5">
											<span class="size-1.5 rounded-full bg-primary"></span>
											Local Licensing
										</p>
										<p class="text-[11px] text-zinc-400 mt-1 leading-relaxed">
											Uses a secure RSA cryptographic lease check. No connection to central Syncra auth servers is required.
										</p>
									</div>
									<div class="bg-zinc-900/40 p-4 rounded-md border border-zinc-800/80">
										<p class="text-xs font-bold text-primary flex items-center gap-1.5">
											<span class="size-1.5 rounded-full bg-primary"></span>
											Custom LLM Models
										</p>
										<p class="text-[11px] text-zinc-400 mt-1 leading-relaxed">
											Bring your own weights or use our pre-packaged local Llama-3 parsing models tuned for invoice and form extraction.
										</p>
									</div>
								</div>
							{/if}
						</Card.Content>
					</Card.Root>
				</div>
			</div>
		</div>
	</section>

	<!-- CREDITS & PRICING TEASER -->
	<section class="py-20 max-w-6xl mx-auto px-4">
		<div class="rounded-xl border border-border bg-card shadow-sm p-8 md:p-12 grid md:grid-cols-12 gap-8 items-center bg-[linear-gradient(135deg,rgba(59,130,246,0.02)_0%,rgba(59,130,246,0.06)_100%)]">
			<div class="md:col-span-8 flex flex-col gap-4">
				<Badge variant="secondary" class="w-fit rounded-full bg-blue-500/10 text-blue-500 border-none font-semibold px-3 py-3 text-xs">
					Zero subscription risk
				</Badge>
				<h2 class="text-3xl font-extrabold tracking-tight sm:text-4xl">
					Simple, credit-based pricing
				</h2>
				<p class="text-muted-foreground text-sm sm:text-base leading-relaxed max-w-2xl">
					We don't lock you into recurring monthly fees. You purchase credits in blocks through Stripe, and we only debit them when a page has been successfully processed. Every signup includes 500 free credits to test the API.
				</p>
				<div class="flex flex-wrap gap-4 mt-2">
					<Button href="/pricing" class="gap-2 cursor-pointer font-semibold">
						View Pricing Details
						<ArrowRightIcon class="size-4" />
					</Button>
				</div>
			</div>
			<div class="md:col-span-4 grid grid-cols-2 gap-3 bg-background/50 backdrop-blur-xs rounded-lg border border-border p-5">
				<div>
					<p class="text-xs text-muted-foreground">Signup Bonus</p>
					<p class="text-xl font-bold mt-0.5">500 Credits</p>
				</div>
				<div>
					<p class="text-xs text-muted-foreground">Minimum Block</p>
					<p class="text-xl font-bold mt-0.5">1,000 Pgs</p>
				</div>
				<div>
					<p class="text-xs text-muted-foreground">Failed Runs</p>
					<p class="text-xl font-bold mt-0.5 text-emerald-500 dark:text-emerald-400">Free</p>
				</div>
				<div>
					<p class="text-xs text-muted-foreground">Currency</p>
					<p class="text-xl font-bold mt-0.5">EUR</p>
				</div>
			</div>
		</div>
	</section>

	<!-- FINAL CALL-TO-ACTION -->
	<section class="py-20 bg-muted/10 border-t border-border/50 text-center relative overflow-hidden">
		<!-- Background glowing accent -->
		<div class="absolute inset-x-0 bottom-0 -z-10 flex justify-center pointer-events-none">
			<div class="w-[50rem] h-[15rem] rounded-full bg-primary/10 blur-[100px]"></div>
		</div>

		<div class="mx-auto max-w-4xl px-4 flex flex-col items-center gap-6">
			<h2 class="text-3xl font-extrabold sm:text-4xl md:text-5xl tracking-tight">
				Ready to automate your document flows?
			</h2>
			<p class="max-w-xl text-muted-foreground text-sm sm:text-base leading-relaxed">
				Join hundreds of developers using Syncra to pipeline invoices, receipts, and structured tax reports into their systems.
			</p>
			<div class="flex flex-wrap justify-center items-center gap-4 mt-2">
				{#if isLoggedIn}
					<Button href="/app" size="lg" class="gap-2 font-semibold cursor-pointer">
						Go to Dashboard
						<ArrowRightIcon class="size-4" />
					</Button>
				{:else}
					<Button href="/signup" size="lg" class="gap-2 font-semibold cursor-pointer">
						Create Free Account
						<ArrowRightIcon class="size-4" />
					</Button>
					<Button href="/login" variant="outline" size="lg" class="cursor-pointer">
						Sign In
					</Button>
				{/if}
			</div>
		</div>
	</section>

	<!-- FOOTER -->
	<footer class="border-t border-border bg-background py-12 text-sm text-muted-foreground">
		<div class="mx-auto max-w-6xl px-4 flex flex-col md:flex-row items-center justify-between gap-6">
			<div class="flex items-center gap-2 font-semibold text-foreground">
				<span class="rounded bg-foreground px-1 py-0.2 text-xs font-bold text-background uppercase">
					S
				</span>
				<span>Syncra</span>
			</div>
			<div class="flex flex-wrap gap-x-6 gap-y-2 justify-center">
				<a href="#features" class="hover:text-foreground transition-colors">Features</a>
				<a href="#how-it-works" class="hover:text-foreground transition-colors">How it Works</a>
				<a href="#on-premise" class="hover:text-foreground transition-colors">On-Premise</a>
				<a href="/pricing" class="hover:text-foreground transition-colors">Pricing</a>
				<a href="/apidoc" class="hover:text-foreground transition-colors">API Docs</a>
			</div>
			<div>
				&copy; {new Date().getFullYear()} Syncra Technologies. All rights reserved.
			</div>
		</div>
	</footer>
</div>
