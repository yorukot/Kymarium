<script lang="ts">
	import Icon from '@iconify/svelte';
	import { goto } from '$app/navigation';
	import { buildOAuthUrl, login } from '$lib/api/auth.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Spinner } from '$lib/components/ui/spinner';
	import { StatusCodes } from 'http-status-codes';
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

	const loginSchema = z.object({
		email: z.email().max(255),
		password: z.string().min(8).max(255)
	});

	let redirectTo = '/';
	const url = new URL(window.location.href);
	redirectTo = url.searchParams.get('next') ?? '/';

	const { form, isSubmitting, setErrors } = createForm({
		extend: [validator({ schema: loginSchema }), reporter()],
		onSubmit: async (values) => {
				const res = await login(values.email, values.password);
				if (res.statusCode != StatusCodes.OK)
				   return (res.message || 'Login failed');
	  },
		onError: (errors) => {
			console.error('Form error:', errors);
			
		}
	});

	const handleGoogle = () => {
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

				<div class="flex flex-col gap-3">
					<Button type="submit" class="w-full" disabled={$isSubmitting}>
						{#if $isSubmitting}
							<Spinner /> Loggin in...
						{:else}
							Login
						{/if}
					</Button>

					<Button type="button" variant="outline" class="w-full" onclick={handleGoogle}>
						<Icon icon="ri:google-fill" class="size-5" />
						Login with Google
					</Button>

					<FieldDescription class="text-center">
						Don't have an account?
						<a href="/auth/register">Sign up</a>
					</FieldDescription>
				</div>
			</FieldGroup>
		</form>
	</Card.Content>
</Card.Root>
