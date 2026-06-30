<script lang="ts">
	import { enhance } from '$app/forms';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { SubmitFunction } from '@sveltejs/kit';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
	let submitting = $state(false);
	let otp = $state('');
	const formValues = $derived((form?.values ?? {}) as Record<string, string>);
	const fieldErrors = $derived((form?.fieldErrors ?? {}) as Record<string, string>);
	const currentEmail = $derived(formValues.email ?? data.email);
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
			<Card.Title>Enter verification code</Card.Title>
			<Card.Description>We sent a 6-digit code to your email.</Card.Description>
		</Card.Header>
		<Card.Content>
			<form method="POST" action="?/verify" use:enhance={submit} class="space-y-4">
				{#if form?.error}<p class="text-sm text-destructive">{form.error}</p>{/if}
				{#if form?.success}<p class="text-sm text-muted-foreground">{form.success}</p>{/if}
				<label class="grid gap-2 text-sm font-medium">Email<input class="h-9 rounded-md border bg-background px-3 text-sm" name="email" type="email" value={currentEmail} required /></label>
				{#if fieldErrors.email}<p class="text-sm text-destructive">{fieldErrors.email}</p>{/if}
				<input type="hidden" name="otp" value={otp} />
				<label class="grid gap-2 text-sm font-medium">Verification code<input class="h-10 rounded-md border bg-background px-3 text-center text-lg tracking-[0.35em]" bind:value={otp} inputmode="numeric" maxlength="6" required /></label>
				{#if fieldErrors.otp}<p class="text-sm text-destructive">{fieldErrors.otp}</p>{/if}
				<Button type="submit" class="w-full" disabled={submitting}>{submitting ? 'Verifying...' : 'Verify'}</Button>
				<Button variant="outline" type="submit" formaction="?/resend" formnovalidate class="w-full" disabled={submitting}>Resend code</Button>
			</form>
		</Card.Content>
	</Card.Root>
</main>
