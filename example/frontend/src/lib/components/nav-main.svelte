<script lang="ts">
	import PlusIcon from "@lucide/svelte/icons/plus";
	import { page } from "$app/state";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import { cn } from "$lib/utils.js";
	import type { Icon } from "@tabler/icons-svelte";
	import StickyNotesIcon from "@lucide/svelte/icons/sticky-notes";
	import { goto } from "$app/navigation";
	import { m } from "$lib/paraglide/messages.js";

	let {
		items,
	}: {
		items: { title: string; url: string; icon?: Icon; plus?: { url: string; title: string } }[];
	} = $props();

	const pathname = $derived(page.url.pathname as string);
	const quickOCRLabel = m.nav_quick_ocr();

	function isActive(url: string) {
		if (url === "#") return false;
		if (url === "/app") return pathname === "/app";

		return pathname === url || pathname.startsWith(`${url}/`);
	}

	const quickOCRActive = $derived(isActive("/app/new-job"));
</script>

<Sidebar.Group>
	<Sidebar.GroupContent class="flex flex-col gap-2">
		<Sidebar.Menu>
			<Sidebar.MenuItem class="flex items-center gap-2">
				<Sidebar.MenuButton
					class="bg-primary text-primary-foreground hover:bg-primary/90 hover:text-primary-foreground active:bg-primary/90 active:text-primary-foreground min-w-8 duration-200 ease-linear"
					tooltipContent={quickOCRLabel}
					isActive={quickOCRActive}
				>
					{#snippet child({ props })}
						<a {...props} href="/app/new-job" aria-current={quickOCRActive ? "page" : undefined}>
							<StickyNotesIcon />
							<span>{quickOCRLabel}</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
				<Button
					size="icon"
					class="size-8 group-data-[collapsible=icon]:opacity-0"
					variant="outline"
					aria-label={m.nav_create_quick_ocr_job()}
					title={m.nav_create_quick_ocr_job()}
					onclick={() => {
						goto("/app/new-job", { keepFocus: true, noScroll: true });
					}}
				>
					<PlusIcon />
					<span class="sr-only">{m.nav_create_quick_ocr_job()}</span>
				</Button>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
		<Sidebar.Menu>
			{#each items as item (item.title)}
				{@const active = isActive(item.url)}
				<Sidebar.MenuItem>
					<Sidebar.MenuButton tooltipContent={item.title} isActive={active}>
						{#snippet child({ props })}
							<a {...props} href={item.url} aria-current={active ? "page" : undefined}>
								{#if item.icon}
									<item.icon />
								{/if}
								<span>{item.title}</span>
							</a>
						{/snippet}
					</Sidebar.MenuButton>
					{#if item.plus}
						{@const plus = item.plus}
						<Sidebar.MenuAction showOnHover>
							{#snippet child({ props })}
								<a
									{...props}
									href={plus.url}
									title={plus.title}
									aria-label={plus.title}
									class={cn(
										props.class as string,
										"hover:!bg-primary hover:!text-primary-foreground text-sidebar-foreground/80 hover:scale-105 active:scale-95 transition-all duration-200"
									)}
								>
									<PlusIcon />
								</a>
							{/snippet}
						</Sidebar.MenuAction>
					{/if}
				</Sidebar.MenuItem>
			{/each}
		</Sidebar.Menu>
	</Sidebar.GroupContent>
</Sidebar.Group>
