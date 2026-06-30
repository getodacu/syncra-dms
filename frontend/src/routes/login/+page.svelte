<script lang="ts">
	import { enhance } from '$app/forms';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { SubmitFunction } from '@sveltejs/kit';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
	let submitting = $state(false);
	const emailValue = $derived(form?.values?.email ?? data.email);
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
		<Card.Header class="text-center">
			<Card.Title>Welcome back</Card.Title>
			<Card.Description>Sign in to Syncra DMS</Card.Description>
		</Card.Header>
		<Card.Content>
			<form method="POST" use:enhance={submit} class="space-y-4">
				<div class="grid gap-2">
					<Button href="/api/auth/google" variant="outline" class="w-full">Continue with Google</Button>
					<Button href="/api/auth/github" variant="outline" class="w-full">Continue with GitHub</Button>
				</div>
				<div class="h-px bg-border"></div>
				{#if data.verified}
					<p class="text-sm text-muted-foreground">Your email is verified. You can sign in now.</p>
				{/if}
				{#if data.reset}
					<p class="text-sm text-muted-foreground">Your password has been reset. You can sign in now.</p>
				{/if}
				{#if data.oauthError || form?.error}
					<p class="text-sm text-destructive">{form?.error ?? 'Social login failed. Please try again.'}</p>
				{/if}
				<label class="grid gap-2 text-sm font-medium">
					Email
					<input class="h-9 rounded-md border bg-background px-3 text-sm" name="email" type="email" value={emailValue} required />
					{#if form?.fieldErrors?.email}<span class="text-sm text-destructive">{form.fieldErrors.email}</span>{/if}
				</label>
				<label class="grid gap-2 text-sm font-medium">
					<span class="flex items-center justify-between gap-3">
						Password
						<a href="/recover-password" class="text-primary hover:underline">Forgot password?</a>
					</span>
					<input class="h-9 rounded-md border bg-background px-3 text-sm" name="password" type="password" required />
					{#if form?.fieldErrors?.password}<span class="text-sm text-destructive">{form.fieldErrors.password}</span>{/if}
				</label>
				<Button type="submit" class="w-full" disabled={submitting}>{submitting ? 'Signing in...' : 'Sign in'}</Button>
				<p class="text-center text-sm text-muted-foreground">No account? <a href="/signup" class="text-primary hover:underline">Sign up</a></p>
			</form>
		</Card.Content>
	</Card.Root>
</main>
