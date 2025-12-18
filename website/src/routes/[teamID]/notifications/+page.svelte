<script lang="ts">
	import Icon from '@iconify/svelte';
	import NotificationsList from '$lib/components/notification/list.svelte';
	import * as Card from '$lib/components/ui/card';
	import type { Notification, NotificationType } from '../../../types/notification';
	import { notificationTypeMeta } from '$lib/utils/notification';
	import CreateNotificationSheet from '$lib/components/notification/setting/basic.svelte';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();

	let createOpen = $state(false);
	let editingNotification = $state<Notification | null>(null);
	type SupportedNotificationType = Extract<NotificationType, 'discord' | 'telegram'>;
	let selectedType = $state<SupportedNotificationType>('discord');
	let notifications = $derived<Notification[]>(data.notifications ?? []);

	const pickerOrder: SupportedNotificationType[] = ['telegram', 'discord'];
	const typeOptions = pickerOrder.map((type) => ({
		type,
		...notificationTypeMeta[type]
	}));

	function openCreator(type: SupportedNotificationType) {
		selectedType = type;
		editingNotification = null;
		createOpen = true;
	}

	function startEdit(notification: Notification) {
		if (notification.type === 'email') return;
		selectedType = notification.type as SupportedNotificationType;
		editingNotification = notification;
		createOpen = true;
	}

	$effect(() => {
		if (!createOpen) {
			editingNotification = null;
		}
	});
</script>

<div class="flex flex-col gap-2">
	<h1 class="text-xl font-semibold">Create a new notification</h1>
	<div
		class="grid gap-4
		       grid-cols-1
		       sm:grid-cols-2
		       lg:grid-cols-3
		       xl:grid-cols-4"
	>
		{#each typeOptions as option (option.type)}
			<Card.Root
				class="w-full cursor-pointer hover:border-primary transition p-4"
				onclick={() => openCreator(option.type)}
			>
				<Card.Content class="flex items-start gap-3 p-0">
					<div>
						<p class="font-semibold flex items-center gap-2">
							<Icon icon={option.icon} class="size-5" />
							{option.label}
						</p>
						<p class="text-sm text-muted-foreground">{option.description}</p>
					</div>
				</Card.Content>
			</Card.Root>
		{/each}
	</div>

	<h1 class="text-xl font-semibold">Notifications</h1>
	<div>
		<NotificationsList
			{notifications}
			onEdit={(notification: Notification) => {
				startEdit(notification);
			}}
		/>
	</div>
</div>

<CreateNotificationSheet
	bind:open={createOpen}
	bind:selectedType
	notification={editingNotification}
	onSaved={(notification: Notification) => {
		const idx = notifications.findIndex((n) => n.id === notification.id);
		if (idx >= 0) {
			notifications = [
				...notifications.slice(0, idx),
				notification,
				...notifications.slice(idx + 1)
			];
		} else {
			notifications = [notification, ...notifications];
		}
	}}
/>
