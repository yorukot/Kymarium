<script lang="ts">
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import { z } from 'zod';
	import { page } from '$app/state';
	import { createNotification } from '$lib/api/notification';
	import { Input } from '$lib/components/ui/input';
	import * as Field from '$lib/components/ui/field';
	import { Button } from '$lib/components/ui/button';
	import * as Sheet from '$lib/components/ui/sheet';
	import { toast } from 'svelte-sonner';
	import type { Notification } from '../../../../types';

	let {
		onCreated,
		onClose
	}: { onCreated?: (notification: Notification) => void; onClose: () => void } = $props();

	const formSchema = z.object({
		type: z.literal('discord'),
		name: z.string().min(1, 'Name is required'),
		config: z.object({
			webhookUrl: z.string().url('Enter a valid webhook URL')
		})
	});

	type FormValues = z.infer<typeof formSchema>;

	const initialValues: FormValues = {
		type: 'discord',
		name: '',
		config: { webhookUrl: '' }
	};

	const { form, errors, isSubmitting, setFields, reset } = createForm<FormValues>({
		initialValues,
		extend: validator({ schema: formSchema }),
		onSubmit: handleSubmit
	});

	function resetForm() {
		reset();
		setFields('name', '');
		setFields('config', { webhookUrl: '' });
	}

	async function handleSubmit(values: FormValues) {
		const teamID = page.params.teamID;
		if (!teamID) {
			toast.error('Missing team id');
			return;
		}

		try {
			const payload = {
				type: 'discord' as const,
				name: values.name,
				config: {
					webhook_url: values.config.webhookUrl
				}
			};

			const res = await createNotification(teamID, payload);
			const created = res.data as unknown as Notification;
			toast.success('Notification created');
			resetForm();
			onClose();
			onCreated?.(created);
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Failed to create notification';
			toast.error(message);
			return { FORM_ERROR: message };
		}
	}
</script>

<form class="flex flex-col gap-4 h-full" use:form>
	<Field.Set>
		<div class="space-y-2">
			<Field.Label for="name">Name</Field.Label>
			<Input id="name" name="name" placeholder="On-call alerts" />
			{#if $errors.name}
				<Field.Description class="text-destructive">{$errors.name[0]}</Field.Description>
			{/if}
		</div>

		<div class="space-y-2">
			<Field.Label for="config.webhookUrl">Webhook URL</Field.Label>
			<Input
				id="config.webhookUrl"
				name="config.webhookUrl"
				type="url"
				placeholder="https://discord.com/api/webhooks/..."
			/>
			{#if $errors.config?.webhookUrl}
				<Field.Description class="text-destructive"
					>{$errors.config.webhookUrl[0]}</Field.Description
				>
			{/if}
		</div>
	</Field.Set>

	<Sheet.Footer class="flex justify-end gap-2 mt-auto">
		<Button type="submit" disabled={$isSubmitting}>{$isSubmitting ? 'Creating…' : 'Create'}</Button>
	</Sheet.Footer>
</form>
