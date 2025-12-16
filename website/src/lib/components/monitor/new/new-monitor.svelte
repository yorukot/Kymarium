<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import MultiSelect, { type MultiSelectOption } from '$lib/components/ui/multi-select';
	import type { MonitorType, Notification, NotificationType } from '../../../../types';
	import { Slider } from '$lib/components/ui/slider';
	import { decidedNotificationIcon } from '../utils';
	import * as Select from '$lib/components/ui/select';

	let selectedMonitorType = $state<MonitorType>('http');
	let selectedNotificationIds = $state<string[]>([]);

	type MonitorTypeSelect = {
		title: string;
		description: string;
		value: MonitorType;
	};

	const monitorTypeSelectData: readonly MonitorTypeSelect[] = [
		{
			title: 'HTTP',
			description: 'Monitor an HTTP or HTTPS endpoint and check response status and latency.',
			value: 'http'
		},
		{
			title: 'Ping',
			description: 'Monitor a host using ICMP ping to check network availability and latency.',
			value: 'ping'
		}
	];

	type IntervalOption = {
		label: string;
		seconds: number;
	};

	const intervalOptions: readonly IntervalOption[] = [
		{ label: '30s', seconds: 30 },
		{ label: '45s', seconds: 45 },
		{ label: '1m', seconds: 60 },
		{ label: '3m', seconds: 180 },
		{ label: '5m', seconds: 300 },
		{ label: '10m', seconds: 600 },
		{ label: '15m', seconds: 900 },
		{ label: '30m', seconds: 1800 },
		{ label: '1h', seconds: 3600 },
		{ label: '2h', seconds: 7200 }
	];
	let intervalIndex = $state<number>(3);

	type ThresholdOption = {
		label: string;
		value: number;
	};

	const thresholdOptions: readonly ThresholdOption[] = [
		{ label: 'Immediate (after 1 check)', value: 1 },
		{ label: 'After 2 checks', value: 2 },
		{ label: 'After 3 checks', value: 3 },
		{ label: 'After 4 checks', value: 4 },
		{ label: 'After 5 checks', value: 5 }
	];

	let failureThresholdValue = $state<string>('1');
	let recoveryThresholdValue = $state<string>('1');

	const failureThreshold = $derived(() => Number(failureThresholdValue));
	const recoveryThreshold = $derived(() => Number(recoveryThresholdValue));

	const failureThresholdLabel = $derived.by(() => {
		const match = thresholdOptions.find(
			(option) => option.value.toString() === failureThresholdValue
		);
		return match?.label ?? 'Select threshold';
	});
	const recoveryThresholdLabel = $derived.by(() => {
		const match = thresholdOptions.find(
			(option) => option.value.toString() === recoveryThresholdValue
		);
		return match?.label ?? 'Select threshold';
	});

	const failureThresholdHelper = $derived.by(() => {
		if (failureThresholdValue === '1') {
			return 'Immediately after the first failed check.';
		}

		return `After ${failureThresholdValue} consecutive failed checks.`;
	});
	const recoveryThresholdHelper = $derived.by(() => {
		if (recoveryThresholdValue === '1') {
			return 'After a single successful check.';
		}

		return `After ${recoveryThresholdValue} consecutive successful checks.`;
	});

	const notifications: Notification[] = [
		{
			id: 'notif_001',
			teamId: 'team_001',
			type: 'discord' as NotificationType,
			name: 'CI Failure Alert',
			config: {
				webhookUrl: 'https://discord.com/api/webhooks/mock/webhook-1'
			},
			createdAt: '2025-01-01T08:00:00.000Z',
			updatedAt: '2025-01-01T08:00:00.000Z'
		},
		{
			id: 'notif_002',
			teamId: 'team_001',
			type: 'telegram' as NotificationType,
			name: 'Deploy Success Notification',
			config: {
				botToken: '123456:mock-telegram-bot-token',
				chatId: '-1001234567890'
			},
			createdAt: '2025-01-02T09:30:00.000Z',
			updatedAt: '2025-01-03T10:15:00.000Z'
		},
		{
			id: 'notif_003',
			teamId: 'team_002',
			type: 'discord' as NotificationType,
			name: 'Security Alert',
			config: {
				webhookUrl: 'https://discord.com/api/webhooks/mock/webhook-2'
			},
			createdAt: '2025-01-04T12:00:00.000Z',
			updatedAt: '2025-01-05T18:45:00.000Z'
		}
	];

	const notificationOptions: MultiSelectOption[] = notifications.map((notification) => ({
		label: notification.name,
		value: notification.id,
		keywords: [notification.type, notification.name],
		icon: decidedNotificationIcon(notification)
	}));
</script>

<Card.Root class="mx-auto w-full">
	<Card.Content>
		<form class="space-y-4">
			<Field.Set>
				<Field.Label for="monitor-name">Monitor name</Field.Label>
				<Input id="monitor-name" name="monitor-name" type="text" placeholder="My Monitor" />
				<Field.Description>The name of your monitor.</Field.Description>

				<Field.Label>Monitor type</Field.Label>
				<Field.Description>Select the type of monitor you want to create.</Field.Description>
				<RadioGroup.Root
					bind:value={selectedMonitorType}
					class="grid gap-4 grid-cols-[repeat(auto-fit,minmax(240px,1fr))]"
				>
					{#each monitorTypeSelectData as option (option.value)}
						<Field.Label for={`monitor-type-${option.value}`}>
							<Field.Field orientation="horizontal">
								<Field.Content>
									<Field.Title>{option.title}</Field.Title>
									<Field.Description>
										{option.description}
									</Field.Description>
								</Field.Content>

								<RadioGroup.Item id={`monitor-type-${option.value}`} value={option.value} />
							</Field.Field>
						</Field.Label>
					{/each}
				</RadioGroup.Root>

				<Field.Label>Interval</Field.Label>
				<Field.Description>
					Monitor will run every
					<span class="font-medium">
						{intervalOptions[intervalIndex].label}
					</span>
				</Field.Description>

				<Slider
					type="single"
					min={0}
					max={intervalOptions.length - 1}
					step={1}
					bind:value={intervalIndex}
					class="mt-4"
				/>
				<div class="mt-2 flex justify-between text-xs text-muted-foreground">
					{#each intervalOptions as option (option.seconds)}
						<span class="w-6 text-center">
							{option.label}
						</span>
					{/each}
				</div>

				<Field.Label>Notifications</Field.Label>
				<Field.Description>
					Select which notifications to send when this monitor triggers.
				</Field.Description>
				<MultiSelect
					name="notificationIds"
					bind:value={selectedNotificationIds}
					options={notificationOptions}
					placeholder="Select notifications"
					emptyMessage="No matching notifications."
					class="mt-2"
				/>

				<Field.Label for="failure-threshold">Failure threshold</Field.Label>
				<Field.Description>{failureThresholdHelper}</Field.Description>
				<Select.Root type="single" bind:value={failureThresholdValue}>
					<Select.Trigger class="w-full justify-between">
						<span data-slot="select-value" class="text-sm font-medium">
							{failureThresholdLabel}
						</span>
					</Select.Trigger>
					<Select.Content>
						<Select.Group>
							<Select.Label>Failed checks</Select.Label>
							{#each thresholdOptions as option (option.value)}
								<Select.Item value={option.value.toString()}>
									{option.label}
								</Select.Item>
							{/each}
						</Select.Group>
					</Select.Content>
				</Select.Root>

				<Field.Label for="recovery-threshold">Recovery threshold</Field.Label>
				<Field.Description>{recoveryThresholdHelper}</Field.Description>
				<Select.Root type="single" bind:value={recoveryThresholdValue}>
					<Select.Trigger class="w-full justify-between">
						<span data-slot="select-value" class="text-sm font-medium">
							{recoveryThresholdLabel}
						</span>
					</Select.Trigger>
					<Select.Content>
						<Select.Group>
							<Select.Label>Successful checks</Select.Label>
							{#each thresholdOptions as option (option.value)}
								<Select.Item value={option.value.toString()}>
									{option.label}
								</Select.Item>
							{/each}
						</Select.Group>
					</Select.Content>
				</Select.Root>
			</Field.Set>
		</form>
	</Card.Content>
</Card.Root>
