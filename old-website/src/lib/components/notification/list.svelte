<script lang="ts">
	import { Card } from '$lib/components/ui/card';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import Icon from '@iconify/svelte';
  import * as AlertDialog from '$lib/components/ui/alert-dialog';
  import { buttonVariants } from '$lib/components/ui/button';
  import { cn } from '$lib/utils';
	import { decidedNotificationIcon, notificationTypeMeta } from '$lib/utils/notification';
	import Button from '$lib/components/ui/button/button.svelte';
	import type { EmailNotificationConfig, Notification } from '$lib/types';
	import { testNotification, deleteNotification } from '$lib/api/notification';
	import { toast } from 'svelte-sonner';
	import { page } from '$app/state';

  let {
    notifications,
    onEdit
  }: { notifications: Notification[]; onEdit?: (notification: Notification) => void } = $props();
	let testingId = $state<string | null>(null);
	let deleteOpen = $state(false);
	let deletingId = $state<string | null>(null);
	let targetNotification = $state<Notification | null>(null);

	function configSummary(notification: Notification): string {
		switch (notification.type) {
			case 'discord':
				return 'Webhook';
			case 'slack':
				return 'Webhook';
			case 'telegram':
				return 'Bot & Chat';
			case 'email':
				{
					const cfg = notification.config as Partial<EmailNotificationConfig> & {
						email_address?: string[];
					};
					const recipients = cfg.emailAddress ?? cfg.email_address ?? [];
					if (recipients.length === 0) {
						return 'Email';
					}
					return `${recipients.length} recipient${recipients.length > 1 ? 's' : ''}`;
				}
		}
	}

	async function handleTest(notification: Notification) {
		if (testingId) return;
		const teamID = page.params.teamID;
		if (!teamID) {
			toast.error('Missing team id');
			return;
		}

		testingId = notification.id;
		try {
			await testNotification(teamID, notification.id);
			toast.success('Test notification sent');
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Failed to send test notification';
			toast.error(message);
		} finally {
			testingId = null;
		}
	}

	function confirmDelete(notification: Notification) {
		targetNotification = notification;
		deleteOpen = true;
	}

	async function handleDelete() {
		if (!targetNotification) return;
		const teamID = page.params.teamID;
		if (!teamID) {
			toast.error('Missing team id');
			return;
		}

		deletingId = targetNotification.id;
		try {
			await deleteNotification(teamID, targetNotification.id);
			notifications = notifications.filter((n) => n.id !== targetNotification?.id);
			toast.success('Notification deleted');
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Failed to delete notification';
			toast.error(message);
		} finally {
			deletingId = null;
			deleteOpen = false;
			targetNotification = null;
		}
	}
</script>

<div class="flex flex-col gap-2">
	{#each notifications as notification (notification.id)}
		<Card class="py-2 px-4">
			<div class="flex items-center justify-between gap-3">
				<div class="flex items-start gap-3 min-w-0">
					<span class="shrink-0 inline-flex items-center justify-center rounded-full bg-muted p-2">
						<Icon icon={decidedNotificationIcon(notification)} class="size-5" />
					</span>
					<div class="min-w-0">
						<div class="flex items-center gap-2 min-w-0">
							<p class="font-semibold truncate">{notification.name}</p>
							<Badge variant="secondary" class="rounded-sm">
								{notificationTypeMeta[notification.type].label}
							</Badge>
						</div>
						<p class="text-sm text-muted-foreground truncate">
							{configSummary(notification)} • Updated {new Date(
								notification.updatedAt
							).toLocaleDateString()}
						</p>
					</div>
				</div>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger class="shrink-0">
						<Button variant="ghost" size="icon">
							<Icon icon="lucide:more-vertical" />
						</Button>
					</DropdownMenu.Trigger>
					<DropdownMenu.Content>
						<DropdownMenu.Group>
							<DropdownMenu.Item
								disabled={testingId === notification.id}
								onclick={() => handleTest(notification)}
							>
								{#if testingId === notification.id}
									<Icon icon="lucide:loader-2" class="animate-spin" />
								{:else}
									<Icon icon="lucide:bell-ring" />
								{/if}
								Send test
							</DropdownMenu.Item>
							<DropdownMenu.Item onclick={() => onEdit?.(notification)}>
								<Icon icon="lucide:edit" /> Edit
							</DropdownMenu.Item>
							<DropdownMenu.Separator />
							<DropdownMenu.Item variant="destructive" onclick={() => confirmDelete(notification)}>
								<Icon icon="lucide:trash" /> Delete
							</DropdownMenu.Item>
						</DropdownMenu.Group>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</div>
		</Card>
	{/each}
</div>

<AlertDialog.Root bind:open={deleteOpen}>
	<AlertDialog.Portal>
		<AlertDialog.Overlay />
		<AlertDialog.Content>
			<AlertDialog.Header>
				<AlertDialog.Title>Delete notification</AlertDialog.Title>
				<AlertDialog.Description>
					Are you sure you want to delete <strong>{targetNotification?.name}</strong>? This action
					cannot be undone.
				</AlertDialog.Description>
			</AlertDialog.Header>
			<AlertDialog.Footer>
				<button
					class={cn(buttonVariants({ variant: 'outline' }))}
					disabled={!!deletingId}
					onclick={() => {
						deleteOpen = false;
						targetNotification = null;
					}}
				>
					Cancel
				</button>
				<button
					class={cn(buttonVariants({ variant: 'destructive' }))}
					onclick={handleDelete}
					disabled={!!deletingId}
				>
					{deletingId ? 'Deleting…' : 'Delete'}
				</button>
			</AlertDialog.Footer>
		</AlertDialog.Content>
	</AlertDialog.Portal>
</AlertDialog.Root>
