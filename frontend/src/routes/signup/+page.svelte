<script lang="ts">
	import { enhance } from '$app/forms';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { SubmitFunction } from '@sveltejs/kit';
	import type { PageProps } from './$types';

	let { form }: PageProps = $props();
	let submitting = $state(false);
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
			<Card.Title>Create an account</Card.Title>
			<Card.Description>Use your work email to start Syncra DMS.</Card.Description>
		</Card.Header>
		<Card.Content>
			<form method="POST" use:enhance={submit} class="space-y-4">
				{#if form?.error}<p class="text-sm text-destructive">{form.error}</p>{/if}
				<label class="grid gap-2 text-sm font-medium">Full name<input class="h-9 rounded-md border bg-background px-3 text-sm" name="name" value={form?.values?.name ?? ''} required /></label>
				{#if form?.fieldErrors?.name}<p class="text-sm text-destructive">{form.fieldErrors.name}</p>{/if}
				<label class="grid gap-2 text-sm font-medium">Email<input class="h-9 rounded-md border bg-background px-3 text-sm" name="email" type="email" value={form?.values?.email ?? ''} required /></label>
				{#if form?.fieldErrors?.email}<p class="text-sm text-destructive">{form.fieldErrors.email}</p>{/if}
				<label class="grid gap-2 text-sm font-medium">Password<input class="h-9 rounded-md border bg-background px-3 text-sm" name="password" type="password" required /></label>
				{#if form?.fieldErrors?.password}<p class="text-sm text-destructive">{form.fieldErrors.password}</p>{/if}
				<label class="grid gap-2 text-sm font-medium">Confirm password<input class="h-9 rounded-md border bg-background px-3 text-sm" name="confirmPassword" type="password" required /></label>
				{#if form?.fieldErrors?.confirmPassword}<p class="text-sm text-destructive">{form.fieldErrors.confirmPassword}</p>{/if}
				<Button type="submit" class="w-full" disabled={submitting}>{submitting ? 'Creating account...' : 'Create account'}</Button>
				<p class="text-center text-sm text-muted-foreground">Already have an account? <a href="/login" class="text-primary hover:underline">Sign in</a></p>
			</form>
		</Card.Content>
	</Card.Root>
</main>
