<script lang="ts">
	import Icon from '@iconify/svelte';
	import { goto } from '$app/navigation';
	import { buildOAuthUrl, login } from '$lib/api/auth.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import {
		FieldGroup,
		Field,
		FieldLabel,
		FieldDescription
	} from '$lib/components/ui/field/index.js';
	import { validator } from '@felte/validator-zod';
	import { reporter, ValidationMessage } from '@felte/reporter-svelte';
	import * as z from 'zod';
	import { createForm } from 'felte';

	const inBrowser = typeof window !== 'undefined';

	const loginSchema = z.object({
		email: z.email().max(255),
		password: z.string().min(8).max(255)
	});

	let redirectTo = '/';
	if (inBrowser) {
		const url = new URL(window.location.href);
		redirectTo = url.searchParams.get('next') ?? '/';
	}

	const { form, isSubmitting } = createForm({
		extend: [validator({ schema: loginSchema }), reporter()],
		onSubmit: async (values) => {
			try {
				await login(values.email, values.password);
				await goto(redirectTo);
			} catch (error) {
				if (
					error instanceof Error &&
					error.message.toLowerCase().includes('email not verified')
				) {
					await goto(`/auth/verify/sent?email=${encodeURIComponent(values.email)}`);
					return;
				}
				return {
					FORM_ERROR:
						error instanceof Error ? error.message : 'Unable to log in right now. Please try again.'
				};
			}
		}
	});

	const handleGoogle = () => {
		if (!inBrowser) return;
		window.location.href = buildOAuthUrl('google', redirectTo);
	};
</script>

<Card.Root class="mx-auto w-full max-w-sm">
	<Card.Header>
		<Card.Title class="text-2xl">Login</Card.Title>
		<Card.Description>Enter your email below to login to your account</Card.Description>
	</Card.Header>
	<Card.Content>
		<form use:form class="space-y-4">
			<FieldGroup>
				<Field>
					<FieldLabel for="email">Email</FieldLabel>
					<Input
						id="email"
						name="email"
						type="email"
						placeholder="m@example.com"
						autocomplete="email"
						required
					/>
					<ValidationMessage for="email" let:messages>
						{#if messages?.length}
							<FieldDescription class="text-destructive">
								{messages[0]}
							</FieldDescription>
						{/if}
					</ValidationMessage>
				</Field>

				<Field>
					<div class="flex items-center">
						<FieldLabel for="password">Password</FieldLabel>
						<a href="##" class="ms-auto inline-block text-sm underline"> Forgot your password? </a>
					</div>
					<Input
						id="password"
						name="password"
						type="password"
						autocomplete="current-password"
						required
					/>
					<ValidationMessage for="password" let:messages>
						{#if messages?.length}
							<FieldDescription class="text-destructive">
								{messages[0]}
							</FieldDescription>
						{/if}
					</ValidationMessage>
				</Field>

				<ValidationMessage for="FORM_ERROR" let:messages>
					{#if messages?.length}
						<FieldDescription class="text-destructive text-center">
							{messages[0]}
						</FieldDescription>
					{/if}
				</ValidationMessage>

				<Field>
					<Button type="submit" class="w-full" disabled={$isSubmitting}>
						{$isSubmitting ? 'Logging in...' : 'Login'}
					</Button>

					<Button type="button" variant="outline" class="w-full" onclick={handleGoogle}>
						<Icon icon="ri:google-fill" class="size-5" />
						Login with Google
					</Button>

					<FieldDescription class="text-center">
						Don't have an account?
						<a href="/auth/register">Sign up</a>
					</FieldDescription>
				</Field>
			</FieldGroup>
		</form>
	</Card.Content>
</Card.Root>
