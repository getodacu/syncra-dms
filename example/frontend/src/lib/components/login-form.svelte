<script lang="ts">
	import { enhance } from "$app/forms";
	import { IconBrandGithub, IconBrandGoogle } from "@tabler/icons-svelte";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import {
		Field,
		FieldGroup,
		FieldLabel,
		FieldDescription,
		FieldSeparator,
		FieldError,
	} from "$lib/components/ui/field/index.js";
	import { cn } from "$lib/utils.js";
	import type { SubmitFunction } from "@sveltejs/kit";
	import type { HTMLAttributes } from "svelte/elements";

	type LoginData = {
		email: string;
		verified: boolean;
		reset: boolean;
		oauthError: string;
	};

	type LoginFormState = {
		values?: { email?: string };
		fieldErrors?: Record<string, string>;
		error?: string;
	};

	type LoginFormProps = HTMLAttributes<HTMLDivElement> & {
		data: LoginData;
		form?: LoginFormState | null;
	};

	let { data, form, class: className, ...restProps }: LoginFormProps = $props();

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

<div class={cn("mx-auto flex w-full max-w-sm flex-col gap-6", className)} {...restProps}>
	<Card.Root>
		<Card.Header class="text-center">
			<Card.Title class="text-xl">Welcome back</Card.Title>
			<Card.Description>Login with your GitHub or Google account</Card.Description>
		</Card.Header>
		<Card.Content>
			<form method="POST" use:enhance={submit}>
				<FieldGroup>
					<Field>
						<Button href="/api/auth/google" variant="outline" class="w-full">
							<IconBrandGoogle class="size-5" />
							Login with Google
						</Button>
						<Button href="/api/auth/github" variant="outline" class="w-full">
							<IconBrandGithub class="size-5" />
							Login with GitHub
						</Button>
					</Field>
					<FieldSeparator class="*:data-[slot=field-separator-content]:bg-card">
						Or continue with
					</FieldSeparator>
					{#if data.verified}
						<FieldDescription>Your email is verified. You can sign in now.</FieldDescription>
					{/if}
					{#if data.reset}
						<FieldDescription>Your password has been reset. You can sign in now.</FieldDescription>
					{/if}
					{#if data.oauthError}
						<FieldError>{data.oauthError}</FieldError>
					{/if}
					{#if form?.error}
						<FieldError>{form.error}</FieldError>
					{/if}
				<Field>
					<FieldLabel for="login-email">Email</FieldLabel>
					<Input
						id="login-email"
						name="email"
						type="email"
						placeholder="m@example.com"
						value={emailValue}
						required
					/>
					{#if form?.fieldErrors?.email}
						<FieldError>{form.fieldErrors.email}</FieldError>
					{/if}
				</Field>
				<Field>
					<div class="flex items-center justify-between gap-3">
						<FieldLabel for="login-password">Password</FieldLabel>
						<a href="/recover-password" class="text-sm font-medium text-primary hover:underline">
							Forgot password?
						</a>
					</div>
					<Input id="login-password" name="password" type="password" required />
					{#if form?.fieldErrors?.password}
						<FieldError>{form.fieldErrors.password}</FieldError>
					{/if}
				</Field>
				<Field>
					<Button type="submit" class="w-full" disabled={submitting}>
						{submitting ? "Logging in..." : "Login"}
					</Button>
					<FieldDescription class="text-center">
						Don't have an account? <a href="/signup">Sign up</a>
					</FieldDescription>
				</Field>
				</FieldGroup>
			</form>
		</Card.Content>
	</Card.Root>
	<FieldDescription class="px-6 text-center">
		By clicking continue, you agree to our <a href="/terms">Terms of Service</a> and
		<a href="/privacy">Privacy Policy</a>.
	</FieldDescription>
</div>
