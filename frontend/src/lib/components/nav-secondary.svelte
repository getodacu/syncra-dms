<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { Component } from 'svelte';

	type NavItem = {
		title: string;
		url: string;
		icon?: Component;
		disabled?: boolean;
	};

	let { items, class: className }: { items: NavItem[]; class?: string } = $props();
</script>

<Sidebar.Group class={className}>
	<Sidebar.GroupContent>
		<Sidebar.Menu>
			{#each items as item (item.title)}
				<Sidebar.MenuItem>
					<Sidebar.MenuButton
						tooltipContent={item.title}
						class={item.disabled ? 'pointer-events-none opacity-50' : undefined}
					>
						{#snippet child({ props })}
							{#if item.disabled}
								<span {...props} aria-disabled="true">
									{#if item.icon}
										<item.icon />
									{/if}
									<span>{item.title}</span>
								</span>
							{:else}
								<a {...props} href={item.url}>
									{#if item.icon}
										<item.icon />
									{/if}
									<span>{item.title}</span>
								</a>
							{/if}
						{/snippet}
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			{/each}
		</Sidebar.Menu>
	</Sidebar.GroupContent>
</Sidebar.Group>
