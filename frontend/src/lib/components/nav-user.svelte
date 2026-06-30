<script lang="ts">
	import { goto } from '$app/navigation';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import CircleUserIcon from '@lucide/svelte/icons/circle-user';
	import CreditCardIcon from '@lucide/svelte/icons/credit-card';
	import LogOutIcon from '@lucide/svelte/icons/log-out';
	import MoreVerticalIcon from '@lucide/svelte/icons/more-vertical';
	import SettingsIcon from '@lucide/svelte/icons/settings';
	import { toast } from 'svelte-sonner';

	type AppShellUser = {
		name: string;
		email: string;
		image?: string | null;
		role: 'user' | 'admin';
	};

	let { user }: { user: AppShellUser | null } = $props();

	const sidebar = Sidebar.useSidebar();

	let logoutConfirmationOpen = $state(false);
	let logoutPending = $state(false);

	const displayName = $derived(user?.name?.trim() || 'Account');
	const displayEmail = $derived(user?.email?.trim() || 'No email address');
	const avatarSrc = $derived(user?.image ?? null);
	const fallbackInitials = $derived(getInitials(displayName, user?.email));

	function getInitials(name: string, email?: string | null) {
		const parts = name.trim().split(/\s+/).filter(Boolean);
		if (parts.length >= 2) return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
		if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
		return (email?.trim()[0] ?? 'A').toUpperCase();
	}

	async function submitLogout(event: SubmitEvent) {
		event.preventDefault();
		if (logoutPending) return;

		const form = event.currentTarget as HTMLFormElement;
		logoutPending = true;

		try {
			const response = await fetch(form.action, { method: form.method });
			if (!response.ok) throw new Error('Logout failed');
			await goto(logoutRedirectPath(response), { invalidateAll: true });
		} catch {
			logoutPending = false;
			toast.error('Logout failed');
		}
	}

	function logoutRedirectPath(response: Response) {
		if (!response.redirected) return '/login';

		const url = new URL(response.url);
		return `${url.pathname}${url.search}${url.hash}`;
	}
</script>

<Sidebar.Menu>
	<Sidebar.MenuItem>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Sidebar.MenuButton
						{...props}
						size="lg"
						class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
					>
						<Avatar.Root class="size-8 rounded-lg grayscale">
							{#if avatarSrc}
								<Avatar.Image src={avatarSrc} alt={displayName} />
							{/if}
							<Avatar.Fallback class="rounded-lg">{fallbackInitials}</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-start text-sm leading-tight">
							<span class="truncate font-medium">{displayName}</span>
							<span class="truncate text-xs text-muted-foreground">{displayEmail}</span>
						</div>
						<MoreVerticalIcon class="ms-auto size-4" />
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				align="end"
				sideOffset={4}
			>
				<DropdownMenu.Label class="p-0 font-normal">
					<div class="flex items-center gap-2 px-1 py-1.5 text-start text-sm">
						<Avatar.Root class="size-8 rounded-lg">
							{#if avatarSrc}
								<Avatar.Image src={avatarSrc} alt={displayName} />
							{/if}
							<Avatar.Fallback class="rounded-lg">{fallbackInitials}</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-start text-sm leading-tight">
							<span class="truncate font-medium">{displayName}</span>
							<span class="truncate text-xs text-muted-foreground">{displayEmail}</span>
						</div>
					</div>
				</DropdownMenu.Label>
				<DropdownMenu.Separator />
				<DropdownMenu.Group>
					<DropdownMenu.Item disabled>
						<CircleUserIcon />
						Account
					</DropdownMenu.Item>
					<DropdownMenu.Item disabled>
						<SettingsIcon />
						Settings
					</DropdownMenu.Item>
					<DropdownMenu.Item disabled>
						<CreditCardIcon />
						Billing
					</DropdownMenu.Item>
				</DropdownMenu.Group>
				<DropdownMenu.Separator />
				<DropdownMenu.Item variant="destructive" onclick={() => (logoutConfirmationOpen = true)}>
					<LogOutIcon />
					Log out
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>

<AlertDialog.Root bind:open={logoutConfirmationOpen}>
	<AlertDialog.Content>
		<form action="/logout" method="POST" class="flex flex-col gap-4" onsubmit={submitLogout}>
			<AlertDialog.Header>
				<AlertDialog.Title>Log out?</AlertDialog.Title>
				<AlertDialog.Description>
					You will need to sign in again to manage Syncra DMS.
				</AlertDialog.Description>
			</AlertDialog.Header>
			<AlertDialog.Footer>
				<AlertDialog.Cancel type="button" disabled={logoutPending}>Cancel</AlertDialog.Cancel>
				<AlertDialog.Action
					type="submit"
					variant="destructive"
					loading={logoutPending}
					disabled={logoutPending}
				>
					<LogOutIcon />
					Log out
				</AlertDialog.Action>
			</AlertDialog.Footer>
		</form>
	</AlertDialog.Content>
</AlertDialog.Root>
