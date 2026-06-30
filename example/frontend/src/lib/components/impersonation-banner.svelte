<script lang="ts">
	import ArrowUpRightIcon from "@lucide/svelte/icons/arrow-up-right";
	import ShieldAlertIcon from "@lucide/svelte/icons/shield-alert";
	import XCircleIcon from "@lucide/svelte/icons/x-circle";
	import { Button } from "$lib/components/ui/button/index.js";

	type AuthUser = NonNullable<App.Locals["user"]>;
	type AuthImpersonation = {
		adminUser: AuthUser;
		targetUser: AuthUser;
		startedAt: string;
	};

	let { impersonation }: { impersonation: AuthImpersonation | null } = $props();

	const targetLabel = $derived(userLabel(impersonation?.targetUser));
	const adminLabel = $derived(userLabel(impersonation?.adminUser));

	function userLabel(user: AuthUser | null | undefined) {
		if (!user) return "";
		return user.name?.trim() || user.email;
	}
</script>

{#if impersonation}
	<div class="border-b border-amber-200 bg-amber-50 text-amber-950 dark:border-amber-900/70 dark:bg-amber-950/25 dark:text-amber-100">
		<div class="flex min-h-12 flex-col gap-3 px-4 py-3 sm:flex-row sm:items-center sm:px-6">
			<div class="flex min-w-0 flex-1 items-start gap-2 sm:items-center">
				<ShieldAlertIcon class="mt-0.5 size-4 shrink-0 text-amber-600 dark:text-amber-300 sm:mt-0" />
				<div class="min-w-0 text-sm leading-5">
					<span class="font-medium">Viewing as {targetLabel}</span>
					<span class="text-amber-800 dark:text-amber-200"> from {adminLabel}</span>
				</div>
			</div>
			<div class="flex shrink-0 flex-wrap items-center gap-2">
				<Button
					href={`/admin-portal/users/${impersonation.targetUser.id}`}
					variant="outline"
					size="sm"
					class="border-amber-300 bg-white/70 text-amber-950 hover:bg-white dark:border-amber-800 dark:bg-amber-950/40 dark:text-amber-100 dark:hover:bg-amber-950/70"
				>
					<ArrowUpRightIcon />
					Admin portal
				</Button>
				<form action="/admin-impersonation/stop" method="POST">
					<Button
						type="submit"
						variant="outline"
						size="sm"
						class="border-amber-300 bg-white/70 text-amber-950 hover:bg-white dark:border-amber-800 dark:bg-amber-950/40 dark:text-amber-100 dark:hover:bg-amber-950/70"
					>
						<XCircleIcon />
						Stop impersonating
					</Button>
				</form>
			</div>
		</div>
	</div>
{/if}
