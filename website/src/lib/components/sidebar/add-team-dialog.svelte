<script lang="ts">
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import { reporter, ValidationMessage } from '@felte/reporter-svelte';
	import { z } from 'zod';
	import { goto } from '$app/navigation';
	import { createTeam } from '$lib/api/team';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Field from '$lib/components/ui/field';
	import { toast } from 'svelte-sonner';

	let { open = $bindable(false) }: { open?: boolean } = $props();

	const schema = z.object({
		name: z.string().min(1, 'Team name is required').max(255, 'Team name is too long')
	});

	const { form, isSubmitting } = createForm({
		extend: [validator({ schema }), reporter()],
		onSubmit: async (values) => {
			try {
				const response = await createTeam(values.name);
				const teamId = response.data.id;
				open = false;
				await goto(`/${teamId}`, { invalidateAll: true });
			} catch (error) {
				const message =
					error instanceof Error
						? error.message
						: 'Unable to create team right now. Please try again.';
				toast.error(message);
				return { FORM_ERROR: message };
			}
		}
	});
</script>

<Dialog.Root bind:open={open}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>Create a new team</Dialog.Title>
			<Dialog.Description>Give your team a name. You can add members later.</Dialog.Description>
		</Dialog.Header>
		<form use:form class="space-y-4">
			<Field.Field>
				<Field.Label for="team-name">Team name</Field.Label>
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
						<Field.Description class="text-destructive">
							{messages[0]}
						</Field.Description>
					{/if}
				</ValidationMessage>
			</Field.Field>

			<ValidationMessage for="FORM_ERROR" let:messages>
				{#if messages?.length}
					<Field.Description class="text-destructive text-center">
						{messages[0]}
					</Field.Description>
				{/if}
			</ValidationMessage>

			<Dialog.Footer class="gap-2">
				<Button type="button" variant="outline" disabled={$isSubmitting} onclick={() => (open = false)}>
					Cancel
				</Button>
				<Button type="submit" disabled={$isSubmitting}>
					{$isSubmitting ? 'Creating…' : 'Create team'}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
