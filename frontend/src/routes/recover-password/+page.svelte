<script lang="ts">
	import { enhance } from '$app/forms';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { SubmitFunction } from '@sveltejs/kit';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
	let submitting = $state(false);
	const formValues = $derived((form?.values ?? {}) as Record<string, string>);
	const fieldErrors = $derived((form?.fieldErrors ?? {}) as Record<string, string>);
	const hasToken = $derived(Boolean(data.token || formValues.token));
	const submit: SubmitFunction = () => {
		submitting = true;
		return async ({ update }) => {
			submitting = false;
			await update({ reset: false });
		};
	};
</script>

<main class="flex min-h-screen items-center justify-center bg-background p-6 text-foreground">
	<Card.Root class="w-full max-w-sm">
		<Card.Header>
			<Card.Title>{hasToken ? 'Reset password' : 'Recover password'}</Card.Title>
			<Card.Description>{hasToken ? 'Choose a new password.' : 'Request a reset link by email.'}</Card.Description>
		</Card.Header>
		<Card.Content>
			<form method="POST" action={hasToken ? '?/reset' : '?/request'} use:enhance={submit} class="space-y-4">
				{#if form?.error}<p class="text-sm text-destructive">{form.error}</p>{/if}
				{#if form?.success}<p class="text-sm text-muted-foreground">{form.success}</p>{/if}
				<label class="grid gap-2 text-sm font-medium">Email<input class="h-9 rounded-md border bg-background px-3 text-sm" name="email" type="email" value={formValues.email ?? data.email} required /></label>
				{#if fieldErrors.email}<p class="text-sm text-destructive">{fieldErrors.email}</p>{/if}
				{#if hasToken}
					<input type="hidden" name="token" value={formValues.token ?? data.token} />
					{#if fieldErrors.token}<p class="text-sm text-destructive">{fieldErrors.token}</p>{/if}
					<label class="grid gap-2 text-sm font-medium">New password<input class="h-9 rounded-md border bg-background px-3 text-sm" name="password" type="password" required /></label>
					{#if fieldErrors.password}<p class="text-sm text-destructive">{fieldErrors.password}</p>{/if}
					<label class="grid gap-2 text-sm font-medium">Confirm password<input class="h-9 rounded-md border bg-background px-3 text-sm" name="confirmPassword" type="password" required /></label>
					{#if fieldErrors.confirmPassword}<p class="text-sm text-destructive">{fieldErrors.confirmPassword}</p>{/if}
				{/if}
				<Button type="submit" class="w-full" disabled={submitting}>{hasToken ? 'Reset password' : 'Send reset link'}</Button>
				<p class="text-center text-sm text-muted-foreground"><a href="/login" class="text-primary hover:underline">Back to login</a></p>
			</form>
		</Card.Content>
	</Card.Root>
</main>
