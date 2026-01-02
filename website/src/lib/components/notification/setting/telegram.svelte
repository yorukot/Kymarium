<script lang="ts">
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import { reporter, ValidationMessage } from '@felte/reporter-svelte';
	import { z } from 'zod';
	import { page } from '$app/state';
	import {
		createNotification,
		updateNotification,
		deleteNotification
	} from '$lib/api/notification';
	import { Input } from '$lib/components/ui/input';
	import * as Field from '$lib/components/ui/field';
	import { Button } from '$lib/components/ui/button';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { cn } from '$lib/utils';
	import { toast } from 'svelte-sonner';
	import type { Notification, TelegramNotificationConfig } from '../$lib/types';

	let {
		notification = null,
		onSaved,
		onDeleted,
		onClose
	}: {
		notification?: Notification | null;
		onSaved?: (notification: Notification) => void;
		onDeleted?: (notification: Notification) => void;
		onClose: () => void;
	} = $props();

	const formSchema = z.object({
		type: z.literal('telegram'),
		name: z.string().min(1, 'Name is required'),
		config: z.object({
			botToken: z.string().min(1, 'Bot token is required'),
			chatId: z.string().min(1, 'Chat ID is required')
		})
	});

	type FormValues = z.infer<typeof formSchema>;

	function deriveInitialValues(): FormValues {
		if (notification) {
			const cfg = notification.config as Partial<TelegramNotificationConfig> & {
				bot_token?: string;
				chat_id?: string;
			};
			return {
				type: 'telegram',
				name: notification.name,
				config: {
					botToken: cfg.botToken ?? cfg.bot_token ?? '',
					chatId: cfg.chatId ?? cfg.chat_id ?? ''
				}
			};
		}
		return {
			type: 'telegram',
			name: '',
			config: { botToken: '', chatId: '' }
		};
	}

	const initialValues: FormValues = deriveInitialValues();

	const { form, isSubmitting, setFields, reset } = createForm<FormValues>({
		initialValues,
		extend: [validator({ schema: formSchema }), reporter()],
		onSubmit: handleSubmit
	});

	let deleteOpen = $state(false);
	let isDeleting = $state(false);

	function resetForm() {
		reset();
		setFields('name', '');
		setFields('config', { botToken: '', chatId: '' });
	}

	async function handleSubmit(values: FormValues) {
		const teamID = page.params.teamID;
		if (!teamID) {
			toast.error('Missing team id');
			return;
		}

		try {
			const payload = {
				name: values.name,
				config: {
					bot_token: values.config.botToken,
					chat_id: values.config.chatId
				}
			};

			let saved: Notification;
			if (notification) {
				const res = await updateNotification(teamID, notification.id, payload);
				saved = res.data;
				toast.success('Notification updated');
			} else {
				const res = await createNotification(teamID, { type: 'telegram', ...payload });
				saved = res.data;
				toast.success('Notification created');
			}

			resetForm();
			onClose();
			onSaved?.(saved);
		} catch (err) {
			const message =
				err instanceof Error
					? err.message
					: notification
						? 'Failed to update notification'
						: 'Failed to create notification';
			toast.error(message);
			return { FORM_ERROR: message };
		}
	}

	async function handleDelete() {
		if (!notification) return;
		const teamID = page.params.teamID;
		if (!teamID) {
			toast.error('Missing team id');
			return;
		}

		isDeleting = true;
		try {
			await deleteNotification(teamID, notification.id);
			toast.success('Notification deleted');
			onDeleted?.(notification);
			deleteOpen = false;
			onClose();
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Failed to delete notification';
			toast.error(message);
		} finally {
			isDeleting = false;
		}
	}

	$effect(() => {
		if (notification) {
			const cfg = notification.config as Partial<TelegramNotificationConfig> & {
				bot_token?: string;
				chat_id?: string;
			};
			setFields('name', notification.name);
			setFields('config', {
				botToken: cfg.botToken ?? cfg.bot_token ?? '',
				chatId: cfg.chatId ?? cfg.chat_id ?? ''
			} as FormValues['config']);
		} else {
			resetForm();
		}
	});
</script>

<form class="flex flex-col gap-4 h-full" use:form>
	<Field.Set>
		<div class="space-y-2">
			<Field.Label for="name">Name</Field.Label>
			<Input id="name" name="name" placeholder="On-call alerts" />
			<ValidationMessage for="name" let:messages>
				{#if messages?.length}
					<Field.Description class="text-destructive">{messages[0]}</Field.Description>
				{/if}
			</ValidationMessage>
		</div>

		<div class="space-y-2">
			<Field.Label for="config.botToken">Bot token</Field.Label>
			<Input id="config.botToken" name="config.botToken" placeholder="123456:ABC-DEF" />
			<ValidationMessage for="config.botToken" let:messages>
				{#if messages?.length}
					<Field.Description class="text-destructive">{messages[0]}</Field.Description>
				{/if}
			</ValidationMessage>
		</div>
		<div class="space-y-2">
			<Field.Label for="config.chatId">Chat ID</Field.Label>
			<Input id="config.chatId" name="config.chatId" placeholder="-1001234567890" />
			<ValidationMessage for="config.chatId" let:messages>
				{#if messages?.length}
					<Field.Description class="text-destructive">{messages[0]}</Field.Description>
				{/if}
			</ValidationMessage>
		</div>
	</Field.Set>

	<Sheet.Footer class="flex justify-end gap-2 mt-auto">
		{#if notification}
			<AlertDialog.Root bind:open={deleteOpen}>
				<AlertDialog.Trigger>
					<Button
						class="w-full"
						variant="destructive"
						type="button"
						disabled={isDeleting || $isSubmitting}
					>
						{isDeleting ? 'Deleting…' : 'Delete'}
					</Button>
				</AlertDialog.Trigger>
				<AlertDialog.Portal>
					<AlertDialog.Overlay />
					<AlertDialog.Content>
						<AlertDialog.Header>
							<AlertDialog.Title>Delete notification</AlertDialog.Title>
							<AlertDialog.Description>
								Are you sure you want to delete <strong>{notification.name}</strong>? This action
								cannot be undone.
							</AlertDialog.Description>
						</AlertDialog.Header>
						<AlertDialog.Footer>
							<AlertDialog.Cancel disabled={isDeleting}>Cancel</AlertDialog.Cancel>
							<AlertDialog.Action
								class={cn('bg-destructive text-destructive-foreground hover:bg-destructive/90')}
								disabled={isDeleting}
								onclick={handleDelete}
							>
								{isDeleting ? 'Deleting…' : 'Delete'}
							</AlertDialog.Action>
						</AlertDialog.Footer>
					</AlertDialog.Content>
				</AlertDialog.Portal>
			</AlertDialog.Root>
		{/if}
		<Button type="submit" disabled={$isSubmitting}>
			{$isSubmitting
				? notification
					? 'Saving…'
					: 'Creating…'
				: notification
					? 'Save changes'
					: 'Create'}
		</Button>
	</Sheet.Footer>
</form>
