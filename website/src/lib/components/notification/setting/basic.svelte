<script lang="ts">
	import { notificationTypeMeta } from '$lib/utils/notification';
	import * as Sheet from '$lib/components/ui/sheet';
	import Icon from '@iconify/svelte';
	import DiscordForm from './discord.svelte';
	import SlackForm from './slack.svelte';
	import TelegramForm from './telegram.svelte';
	import EmailForm from './email.svelte';
	import type { NotificationType, Notification } from '../$lib/types';

	type SupportedNotificationType = Extract<NotificationType, 'discord' | 'slack' | 'telegram' | 'email'>;

	let {
		open = $bindable(false),
		selectedType = $bindable<SupportedNotificationType>('discord'),
		notification = null,
		onSaved,
		onDeleted
	}: {
		open: boolean;
		selectedType: SupportedNotificationType;
		notification?: Notification | null;
		onSaved?: (notification: Notification) => void;
		onDeleted?: (notification: Notification) => void;
	} = $props();

	$effect(() => {
		if (notification) {
			// lock selected type to the notification being edited
			selectedType = notification.type as SupportedNotificationType;
		}
	});
</script>

<Sheet.Root bind:open onOpenChange={(next) => (open = next)}>
	<Sheet.Content side="right" class="max-w-md w-full">
		<Sheet.Header class="p-4 border-b">
			<Sheet.Title>Create notification</Sheet.Title>
			<Sheet.Description>{notificationTypeMeta[selectedType].description}</Sheet.Description>
		</Sheet.Header>

		<div class="p-4 pt-2 flex flex-col gap-4 h-full">
			<div class="flex items-center gap-2 text-lg font-semibold">
				<Icon icon={notificationTypeMeta[selectedType].icon} class="size-8" />
				<span class="capitalize">{notificationTypeMeta[selectedType].label}</span>
			</div>

			{#if selectedType === 'discord'}
				<DiscordForm
					{notification}
					{onSaved}
					{onDeleted}
					onClose={() => {
						open = false;
						notification = null;
					}}
				/>
			{:else if selectedType === 'slack'}
				<SlackForm
					{notification}
					{onSaved}
					{onDeleted}
					onClose={() => {
						open = false;
						notification = null;
					}}
				/>
			{:else if selectedType === 'telegram'}
				<TelegramForm
					{notification}
					{onSaved}
					{onDeleted}
					onClose={() => {
						open = false;
						notification = null;
					}}
				/>
			{:else if selectedType === 'email'}
				<EmailForm
					{notification}
					{onSaved}
					{onDeleted}
					onClose={() => {
						open = false;
						notification = null;
					}}
				/>
			{/if}
		</div>
	</Sheet.Content>
</Sheet.Root>
