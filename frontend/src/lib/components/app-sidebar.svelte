<script lang="ts">
	import Building2Icon from '@lucide/svelte/icons/building-2';
	import CircleHelpIcon from '@lucide/svelte/icons/circle-help';
	import LayoutDashboardIcon from '@lucide/svelte/icons/layout-dashboard';
	import NavMain from './nav-main.svelte';
	import NavSecondary from './nav-secondary.svelte';
	import NavUser from './nav-user.svelte';
	import SidebarSpaceSwitcher from './sidebar-space-switcher.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';

	type AppShellUser = {
		name: string;
		email: string;
		image?: string | null;
		role: 'user' | 'admin';
	};

	const data = {
		navMain: [
			{
				title: 'Dashboard',
				url: '/app',
				icon: LayoutDashboardIcon
			},
			{
				title: 'Organization Units',
				url: '/app/organization-units',
				icon: Building2Icon
			}
		],
		navSecondary: [
			{
				title: 'Help',
				url: '#',
				icon: CircleHelpIcon,
				disabled: true
			}
		]
	};

	let {
		user,
		...restProps
	}: ComponentProps<typeof Sidebar.Root> & { user: AppShellUser | null } = $props();
</script>

<Sidebar.Root collapsible="offcanvas" {...restProps}>
	<Sidebar.Header>
		<SidebarSpaceSwitcher {user} />
	</Sidebar.Header>
	<Sidebar.Content>
		<NavMain items={data.navMain} />
		<NavSecondary items={data.navSecondary} class="mt-auto" />
	</Sidebar.Content>
	<Sidebar.Footer>
		<NavUser {user} />
	</Sidebar.Footer>
</Sidebar.Root>
