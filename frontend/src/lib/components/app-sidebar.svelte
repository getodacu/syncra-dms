<script lang="ts">
	import Building2Icon from '@lucide/svelte/icons/building-2';
	import CircleHelpIcon from '@lucide/svelte/icons/circle-help';
	import KeyRoundIcon from '@lucide/svelte/icons/key-round';
	import LayoutDashboardIcon from '@lucide/svelte/icons/layout-dashboard';
	import ShieldCheckIcon from '@lucide/svelte/icons/shield-check';
	import UserRoundCogIcon from '@lucide/svelte/icons/user-round-cog';
	import UsersIcon from '@lucide/svelte/icons/users';
	import NavMain from './nav-main.svelte';
	import NavSecondary from './nav-secondary.svelte';
	import NavUser from './nav-user.svelte';
	import SidebarSpaceSwitcher from './sidebar-space-switcher.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { Component, ComponentProps } from 'svelte';

	type AppShellUser = {
		name: string;
		email: string;
		image?: string | null;
		role: 'user' | 'admin';
	};

	type NavItem = {
		title: string;
		url: string;
		icon: Component;
	};

	const baseNav: NavItem[] = [
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
	];

	const data = {
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
		permissions = [],
		...restProps
	}: ComponentProps<typeof Sidebar.Root> & {
		user: AppShellUser | null;
		permissions?: string[];
	} = $props();

	const canViewAdminArea = (permission: string) =>
		permissions.includes('system.admin') || permissions.includes(permission);

	const adminNav = $derived([
		...(canViewAdminArea('user.view')
			? [
					{
						title: 'Users',
						url: '/app/admin/users',
						icon: UsersIcon
					}
				]
			: []),
		...(canViewAdminArea('role.view')
			? [
					{
						title: 'Roles',
						url: '/app/admin/roles',
						icon: ShieldCheckIcon
					},
					{
						title: 'Permissions',
						url: '/app/admin/permissions',
						icon: KeyRoundIcon
					}
				]
			: []),
		...(canViewAdminArea('group.view')
			? [
					{
						title: 'Groups',
						url: '/app/admin/groups',
						icon: UserRoundCogIcon
					}
				]
			: [])
	]);

	const navMain = $derived([...baseNav, ...adminNav]);
</script>

<Sidebar.Root collapsible="offcanvas" {...restProps}>
	<Sidebar.Header>
		<SidebarSpaceSwitcher {user} />
	</Sidebar.Header>
	<Sidebar.Content>
		<NavMain items={navMain} />
		<NavSecondary items={data.navSecondary} class="mt-auto" />
	</Sidebar.Content>
	<Sidebar.Footer>
		<NavUser {user} />
	</Sidebar.Footer>
</Sidebar.Root>
