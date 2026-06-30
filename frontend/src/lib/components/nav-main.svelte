<script lang="ts">
	import { page } from '$app/state';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { Component } from 'svelte';

	type NavItem = {
		title: string;
		url: string;
		icon?: Component;
	};

	let { items }: { items: NavItem[] } = $props();

	const pathname = $derived(page.url.pathname as string);

	function isActive(url: string) {
		if (url === '/app') return pathname === '/app';
		return pathname === url || pathname.startsWith(`${url}/`);
	}
</script>

<Sidebar.Group>
	<Sidebar.GroupContent class="flex flex-col gap-2">
		<Sidebar.Menu>
			{#each items as item (item.title)}
				{@const active = isActive(item.url)}
				<Sidebar.MenuItem>
					<Sidebar.MenuButton tooltipContent={item.title} isActive={active}>
						{#snippet child({ props })}
							<a {...props} href={item.url} aria-current={active ? 'page' : undefined}>
								{#if item.icon}
									<item.icon />
								{/if}
								<span>{item.title}</span>
							</a>
						{/snippet}
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			{/each}
		</Sidebar.Menu>
	</Sidebar.GroupContent>
</Sidebar.Group>
