<script lang="ts">
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import CheckCircle2Icon from "@lucide/svelte/icons/check-circle-2";
	import CircleIcon from "@lucide/svelte/icons/circle";
	import CoinsIcon from "@lucide/svelte/icons/coins";
	import DatabaseIcon from "@lucide/svelte/icons/database";
	import FileUpIcon from "@lucide/svelte/icons/file-up";
	import KeyRoundIcon from "@lucide/svelte/icons/key-round";
	import LibraryIcon from "@lucide/svelte/icons/library";
	import WebhookIcon from "@lucide/svelte/icons/webhook";

	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import type { DashboardSummaryResponse } from "./api";

	let { summary }: { summary: DashboardSummaryResponse } = $props();

	const steps = $derived([
		{
			label: m.dashboard_step_schema(),
			done: summary.onboarding.has_schema,
			href: "/app/schemas/library",
			icon: LibraryIcon
		},
		{
			label: m.dashboard_step_ocr_job(),
			done: summary.onboarding.has_completed_document,
			href: "/app/new-job",
			icon: FileUpIcon
		},
		{
			label: m.dashboard_step_dataset(),
			done: summary.onboarding.has_dataset,
			href: "/app/datasets",
			icon: DatabaseIcon
		},
		{
			label: m.dashboard_step_api_key(),
			done: summary.onboarding.has_api_key,
			href: "/app/developer-settings",
			icon: KeyRoundIcon
		},
		{
			label: m.dashboard_step_webhook(),
			done: summary.onboarding.has_webhook,
			href: "/app/developer-settings",
			icon: WebhookIcon
		}
	]);

	function creditsLabel(count: number) {
		return count === 1 ? m.dashboard_credits_one({ count }) : m.dashboard_credits_other({ count });
	}
</script>

<section class="px-4 lg:px-6">
	<Card.Root>
		<Card.Header class="gap-4">
			<div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
				<div>
					<Card.Title>{m.dashboard_onboarding_title()}</Card.Title>
					<Card.Description>{m.dashboard_onboarding_description()}</Card.Description>
				</div>
				<Button href="/app/new-job" class="w-full gap-2 sm:w-auto">
					{m.dashboard_new_ocr_job()}
					<ArrowRightIcon class="size-4" aria-hidden="true" />
				</Button>
			</div>
			<div class="flex flex-wrap gap-2">
				<Badge variant={summary.credit_summary.low_credit ? "destructive" : "secondary"} class="gap-1">
					<CoinsIcon class="size-3.5" aria-hidden="true" />
					{creditsLabel(summary.credit_summary.available_credits)}
				</Badge>
			</div>
		</Card.Header>
		<Card.Content>
			<div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-5">
				{#each steps as step (step.label)}
					<a href={step.href} class="group rounded-lg border p-3 transition-colors hover:bg-muted/50">
						<div class="flex items-center justify-between gap-3">
							<step.icon class="size-4 text-muted-foreground" aria-hidden="true" />
							{#if step.done}
								<CheckCircle2Icon class="size-4 text-emerald-500" aria-hidden="true" />
							{:else}
								<CircleIcon class="size-4 text-muted-foreground" aria-hidden="true" />
							{/if}
						</div>
						<div class="mt-3 text-sm font-medium">{step.label}</div>
						<div class="mt-1 text-xs text-muted-foreground">{step.done ? m.dashboard_step_ready() : m.dashboard_step_open()}</div>
					</a>
				{/each}
			</div>
		</Card.Content>
	</Card.Root>
</section>
