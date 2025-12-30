<script lang="ts">
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import { z } from 'zod';
	import { page } from '$app/state';
	import {
		createNotification,
		updateNotification,
		deleteNotification
	} from '$lib/api/notification';
	import * as Field from '$lib/components/ui/field';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { cn } from '$lib/utils';
	import { toast } from 'svelte-sonner';
	import Textarea from '$lib/components/ui/textarea/textarea.svelte';
	import type { EmailNotificationConfig, Notification } from '../../../types';

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

	const emailSchema = z.string().email();

	function parseEmailList(value: string): string[] {
		return value
			.split(/[,\n]/)
			.map((email) => email.trim())
			.filter((email) => email.length > 0);
	}

	const formSchema = z.object({
		type: z.literal('email'),
		name: z.string().min(1, 'Name is required'),
		config: z.object({
			emailAddress: z
				.string()
				.min(1, 'At least one email is required')
				.refine((value) => {
					const emails = parseEmailList(value);
					if (emails.length === 0) return false;
					return emails.every((email) => emailSchema.safeParse(email).success);
				}, 'Enter valid email addresses')
		})
	});

	type FormValues = z.infer<typeof formSchema>;

	function deriveInitialValues(): FormValues {
		if (notification) {
			const cfg = notification.config as Partial<EmailNotificationConfig> & {
				email_address?: string[];
			};
			const recipients = cfg.emailAddress ?? cfg.email_address ?? [];
			return {
				type: 'email',
				name: notification.name,
				config: {
					emailAddress: recipients.join('\n')
				}
			};
		}
		return {
			type: 'email',
			name: '',
			config: { emailAddress: '' }
		};
	}

	const initialValues: FormValues = deriveInitialValues();

	const { form, errors, isSubmitting, setFields, reset } = createForm<FormValues>({
		initialValues,
		extend: validator({ schema: formSchema }),
		onSubmit: handleSubmit
	});

	let deleteOpen = $state(false);
	let isDeleting = $state(false);

	function resetForm() {
		reset();
		setFields('name', '');
		setFields('config', { emailAddress: '' });
	}

	async function handleSubmit(values: FormValues) {
		const teamID = page.params.teamID;
		if (!teamID) {
			toast.error('Missing team id');
			return;
		}

		const recipients = parseEmailList(values.config.emailAddress);
		if (recipients.length === 0) {
			toast.error('Add at least one valid email');
			return;
		}

		try {
			const payload = {
				name: values.name,
				config: {
					email_address: recipients
				}
			};

			let saved: Notification;
			if (notification) {
				const res = await updateNotification(teamID, notification.id, payload);
				saved = res.data;
				toast.success('Notification updated');
			} else {
				const res = await createNotification(teamID, { type: 'email', ...payload });
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
			const cfg = notification.config as Partial<EmailNotificationConfig> & {
				email_address?: string[];
			};
			const recipients = cfg.emailAddress ?? cfg.email_address ?? [];
			setFields('name', notification.name);
			setFields('config', { emailAddress: recipients.join('\n') } as FormValues['config']);
		} else {
			resetForm();
		}
	});
</script>

<form class="flex flex-col gap-4 h-full" use:form>
	<Field.Set>
		<div class="space-y-2">
			<Field.Label for="name">Name</Field.Label>
			<Input id="name" name="name" placeholder="On-call email" />
			{#if $errors.name}
				<Field.Description class="text-destructive">{$errors.name[0]}</Field.Description>
			{/if}
		</div>

		<div class="space-y-2">
			<Field.Label for="config.emailAddress">Recipients</Field.Label>
			<Field.Description>One email per line or comma-separated.</Field.Description>
			<Textarea
				id="config.emailAddress"
				name="config.emailAddress"
				placeholder="alerts@example.com\nstatus@example.com"
				rows={4}
			/>
			{#if $errors.config?.emailAddress}
				<Field.Description class="text-destructive"
					>{$errors.config.emailAddress[0]}</Field.Description
				>
			{/if}
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
