<script lang="ts">
	import { notificationTypeMeta } from '$lib/utils/notification';
	import * as Sheet from '$lib/components/ui/sheet';
	import Icon from '@iconify/svelte';
	import DiscordForm from './discord.svelte';
	import TelegramForm from './telegram.svelte';
	import type { NotificationType, Notification } from '../../../../types';

	type SupportedNotificationType = Extract<NotificationType, 'discord' | 'telegram'>;

	let {
		open = $bindable(false),
		selectedType = $bindable<SupportedNotificationType>('discord'),
		onCreated
	}: {
		open: boolean;
		selectedType: SupportedNotificationType;
		onCreated?: (notification: Notification) => void;
	} = $props();
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
					{onCreated}
					onClose={() => {
						open = false;
					}}
				/>
			{:else if selectedType === 'telegram'}
				<TelegramForm
					{onCreated}
					onClose={() => {
						open = false;
					}}
				/>
			{/if}
		</div>
	</Sheet.Content>
</Sheet.Root>
