<script lang="ts">
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import { reporter, ValidationMessage } from '@felte/reporter-svelte';
	import { z } from 'zod';
	import { goto } from '$app/navigation';

	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { FieldGroup, Field, FieldLabel, FieldDescription } from '$lib/components/ui/field';
	import { createTeam } from '$lib/api/team';

	const schema = z.object({
		name: z.string().min(1, 'Team name is required').max(255, 'Team name is too long')
	});

	const { form, isSubmitting } = createForm({
		extend: [validator({ schema }), reporter()],
		onSubmit: async (values) => {
			try {
				const response = await createTeam(values.name);
				const teamId = response.data.id;
				await goto(`/${teamId}`);
			} catch (error) {
				return {
					FORM_ERROR:
						error instanceof Error
							? error.message
							: 'Unable to create team right now. Please try again.'
				};
			}
		}
	});
</script>

<div>
<div class="min-h-screen px-4 flex items-center justify-center">
	<Card.Root class="w-full max-w-md">
		<Card.Header>
			<Card.Title class="text-2xl">Create a new team</Card.Title>
			<Card.Description>Give your team a name. You can add members later.</Card.Description>
		</Card.Header>
		<Card.Content>
			<form use:form class="space-y-4">
				<FieldGroup>
					<Field>
						<FieldLabel for="team-name">Team name</FieldLabel>
						<Input
							id="team-name"
							name="name"
							type="text"
							placeholder="Acme SRE"
							autocomplete="organization"
							required
						/>
						<ValidationMessage for="name" let:messages>
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
							{$isSubmitting ? 'Creating…' : 'Create team'}
						</Button>
					</Field>
				</FieldGroup>
			</form>
		</Card.Content>
	</Card.Root>
</div>
</div>
