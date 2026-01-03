<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Field from '$lib/components/ui/field';
	import Input from '$lib/components/ui/input/input.svelte';
	import { ValidationMessage } from '@felte/reporter-svelte';
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>General Settings</Card.Title>
		<Card.Description>Configure the general settings for your status page.</Card.Description>
	</Card.Header>
	<Card.Content>
		<Field.Group class="flex gap-2 sm:flex-row">
			<Field.Field>
				<Field.Label for="status-page-name">Name</Field.Label>
				<Input
					id="status-page-name"
					type="text"
					name="name"
					placeholder="Acme Status"
					class="w-full"
				/>

				<Field.Description>
					<ValidationMessage for="name" let:messages>
						{#if messages?.length}
							<p class="text-sm text-destructive">{messages[0]}</p>
						{:else}
							Display name shown on the public status page header and metadata.
						{/if}
					</ValidationMessage>
				</Field.Description>
			</Field.Field>

			<Field.Field>
				<Field.Label for="status-page-slug">Slug</Field.Label>
				<Input
					id="status-page-slug"
					type="text"
					name="slug"
					placeholder="acme-status"
					class="w-full"
				/>

				<Field.Description>
					<ValidationMessage for="slug" let:messages>
						{#if messages?.length}
							<p class="text-sm text-destructive">{messages[0]}</p>
						{:else}
							Used in the public URL (e.g., status.kymarium.com/acme-status).
						{/if}
					</ValidationMessage>
				</Field.Description>
			</Field.Field>
		</Field.Group>
	</Card.Content>
</Card.Root>
