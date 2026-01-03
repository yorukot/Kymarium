<script lang="ts">
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import { reporter, ValidationMessage } from '@felte/reporter-svelte';
	import { z } from 'zod';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { FieldGroup, Field, FieldLabel, FieldDescription } from '$lib/components/ui/field';
	import { createStatusPage } from '$lib/api/status-page';
	import Icon from '@iconify/svelte';

	const schema = z.object({
		title: z.string().min(11, 'Title is required').max(255, 'Title is too long'),
		slug: z.string().min(3, 'Slug is required').max(255, 'Slug is too long')
	});

	const { form, isSubmitting } = createForm({
		extend: [validator({ schema }), reporter()],
		onSubmit: async (values) => {
			try {
				const teamID = page.params.teamID;
				if (!teamID) {
					return { FORM_ERROR: 'Missing team id' };
				}

				const response = await createStatusPage(teamID, {
					title: values.title,
					slug: values.slug,
					icon: null,
					groups: [],
					monitors: []
				});

				await goto(`/${teamID}/status-pages/${response.data.statusPage.id}/edit`);
			} catch (error) {
				return {
					FORM_ERROR:
						error instanceof Error
							? error.message
							: 'Unable to create status page right now. Please try again.'
				};
			}
		}
	});

	function slugify(value: string) {
		return value
			.toLowerCase()
			.trim()
			.replace(/[^a-z0-9\s-]/g, '')
			.replace(/\s+/g, '-')
			.replace(/-+/g, '-');
	}
</script>

<div class="flex flex-col gap-5">
	<header class="flex flex-col gap-3">
		<div class="flex justify-between">
			<div>
				<h1 class="text-2xl font-bold">Status pages</h1>
			</div>
			<Button size="default" variant="ghost" href="../status-pages">
				<Icon icon="lucide:arrow-left" />
				Back to list
			</Button>
		</div>
		<p class="text-sm text-muted-foreground">
			Start with a title and slug. You can add monitors later.
		</p>
	</header>

	<form use:form class="space-y-4">
		<Card.Root class="w-full">
			<Card.Content>
				<FieldGroup>
					<Field>
						<FieldLabel for="status-page-title">Title</FieldLabel>
						<Input
							id="status-page-title"
							name="title"
							type="text"
							placeholder="Public Status"
							autocomplete="off"
							required
							inputmode="text"
							oninput={(event) => {
								const target = event.currentTarget as HTMLInputElement;
								const slugInput = document.getElementById(
									'status-page-slug'
								) as HTMLInputElement | null;
								if (slugInput && !slugInput.value) {
									slugInput.value = slugify(target.value);
									slugInput.dispatchEvent(new Event('input'));
								}
							}}
						/>
						<ValidationMessage for="title" let:messages>
							{#if messages?.length}
								<FieldDescription class="text-destructive">
									{messages[0]}
								</FieldDescription>
							{/if}
						</ValidationMessage>
					</Field>

					<Field>
						<FieldLabel for="status-page-slug">Slug</FieldLabel>
						<Input
							id="status-page-slug"
							name="slug"
							type="text"
							placeholder="public-status"
							autocomplete="off"
							required
							inputmode="url"
							onblur={(event) => {
								const target = event.currentTarget as HTMLInputElement;
								target.value = slugify(target.value);
								target.dispatchEvent(new Event('input'));
							}}
						/>
						<ValidationMessage for="slug" let:messages>
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
				</FieldGroup>
			</Card.Content>
		</Card.Root>
		<div class="w-full flex justify-end">
			<Field class="w-fit">
				<Button type="submit" disabled={$isSubmitting}>
					{$isSubmitting ? 'Creating…' : 'Create status page'}
				</Button>
			</Field>
		</div>
	</form>
</div>
