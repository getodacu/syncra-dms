<script lang="ts">
	import CodeIcon from "@tabler/icons-svelte/icons/code";
	import CreditCardIcon from "@tabler/icons-svelte/icons/credit-card";
	import DashboardIcon from "@tabler/icons-svelte/icons/dashboard";
	import DatabaseIcon from "@tabler/icons-svelte/icons/database";
	import FileAiIcon from "@tabler/icons-svelte/icons/file-ai";
	import HelpIcon from "@tabler/icons-svelte/icons/help";
	import NavCollections from "./nav-collections.svelte";
	import NavDatasets from "./nav-datasets.svelte";
	import NavMain from "./nav-main.svelte";
	import NavSecondary from "./nav-secondary.svelte";
	import NavUser from "./nav-user.svelte";
	import SidebarSpaceSwitcher from "./sidebar-space-switcher.svelte";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import type { ComponentProps } from "svelte";

	type AuthUser = NonNullable<App.Locals["user"]>;

	const data = {
		navMain: [
			{
				title: m.nav_dashboard(),
				url: "/app",
				icon: DashboardIcon,
			},
			{
				title: m.nav_schemas(),
				url: "/app/schemas",
				icon: DatabaseIcon,
				plus: {
					url: "/app/schemas/new",
					title: m.nav_create_schema(),
				},
			},
			{
				title: m.nav_jobs(),
				url: "/app/jobs",
				icon: FileAiIcon,
				plus: {
					url: "/app/new-job",
					title: m.nav_create_job(),
				},
			},
			{
				title: m.nav_billing(),
				url: "/app/billing",
				icon: CreditCardIcon,
			},
		],
		navSecondary: [
			{
				title: m.nav_developer_settings(),
				url: "/app/developer-settings",
				icon: CodeIcon,
			},
			{
				title: m.nav_get_help(),
				url: "#",
				icon: HelpIcon,
			},
		],
	};

	let {
		user,
		...restProps
	}: ComponentProps<typeof Sidebar.Root> & { user: AuthUser | null } = $props();
</script>

<Sidebar.Root collapsible="offcanvas" {...restProps}>
	<Sidebar.Header>
		<SidebarSpaceSwitcher user={user} activeSpace="app" />
	</Sidebar.Header>
	<Sidebar.Content>
		<NavMain items={data.navMain} />
		<NavCollections />
		<NavDatasets />
		<NavSecondary items={data.navSecondary} class="mt-auto" />
	</Sidebar.Content>
	<Sidebar.Footer>
		<NavUser {user} />
	</Sidebar.Footer>
</Sidebar.Root>
