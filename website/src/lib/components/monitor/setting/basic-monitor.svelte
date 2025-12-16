<script lang="ts">
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import { z } from 'zod';

	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import MultiSelect, { type MultiSelectOption } from '$lib/components/ui/multi-select';
	import type { MonitorType, Notification, NotificationType } from '../../../../types';
	import { Slider } from '$lib/components/ui/slider';
	import { decidedNotificationIcon } from '../utils';
	import * as Select from '$lib/components/ui/select';
	import HttpMonitor from './http-monitor.svelte';
	import PingMonitor from './ping-monitor.svelte';
	import { Button } from '$lib/components/ui/button';

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

	const httpMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS'] as const;

	const httpConfigSchema = z.object({
		url: z.string().url('Must be a valid URL'),
		method: z.enum(httpMethods),
		maxRedirects: z.coerce.number().int().min(0).max(1000),
		requestTimeoutSeconds: z.coerce.number().int().min(0).max(120),
		headers: z.string().optional().default(''),
		bodyEncoding: z.enum(['json', 'xml', '']).default(''),
		body: z.string().optional().default(''),
		acceptedStatusCodes: z.array(z.string()).min(1, 'Select at least one status'),
		upsideDownMode: z.coerce.boolean(),
		certificateExpiryNotification: z.coerce.boolean(),
		ignoreTlsError: z.coerce.boolean()
	});

	const pingConfigSchema = z.object({
		host: z.string().min(1, 'Host is required'),
		timeoutSeconds: z.coerce.number().int().min(0).max(120),
		packetSize: z
			.union([z.coerce.number().int().min(1).max(65000), z.literal('')])
			.optional()
			.default('')
	});

	const baseSchema = z.object({
		name: z.string().min(1, 'Monitor name is required').max(255),
		type: z.enum(['http', 'ping']),
		interval: z.coerce.number().int().min(10, 'Interval must be at least 10 seconds'),
		failureThreshold: z.coerce.number().int().min(1).max(10),
		recoveryThreshold: z.coerce.number().int().min(1).max(10),
		notification: z.array(z.string())
	});

	const monitorFormSchema = z.discriminatedUnion('type', [
		baseSchema.extend({
			type: z.literal('http'),
			config: httpConfigSchema
		}),
		baseSchema.extend({
			type: z.literal('ping'),
			config: pingConfigSchema
		})
	]);

	const defaultHttpConfig = () => ({
		url: '',
		method: 'GET' as const,
		maxRedirects: 5,
		requestTimeoutSeconds: 10,
		headers: '',
		bodyEncoding: '' as const,
		body: '',
		acceptedStatusCodes: ['2xx'],
		upsideDownMode: false,
		certificateExpiryNotification: true,
		ignoreTlsError: false
	});

	const defaultPingConfig = () => ({
		host: '',
		timeoutSeconds: 5,
		packetSize: ''
	});

	const initialIntervalIndex = 3;
	const initialFailureThreshold = '1';
	const initialRecoveryThreshold = '1';

	let selectedMonitorType = $state<MonitorType>('http');
	let intervalIndex = $state<number>(initialIntervalIndex);
	let failureThresholdValue = $state<string>(initialFailureThreshold);
	let recoveryThresholdValue = $state<string>(initialRecoveryThreshold);
	let selectedNotificationIds = $state<string[]>([]);

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

	const initialValues = {
		name: '',
		type: 'http' as MonitorType,
		interval: intervalOptions[initialIntervalIndex].seconds,
		failureThreshold: Number(initialFailureThreshold),
		recoveryThreshold: Number(initialRecoveryThreshold),
		notification: [] as string[],
		config: defaultHttpConfig()
	};

	const { form, errors, data, isSubmitting, setFields } = createForm({
		initialValues,
		extend: validator({ schema: monitorFormSchema }),
		onSubmit: async (values) => {
			// TODO: wire to API
			console.log('submit monitor payload', values);
		}
	});

	$effect(() => {
		setFields('interval', intervalOptions[intervalIndex].seconds);
	});

	$effect(() => {
		setFields('failureThreshold', Number(failureThresholdValue));
	});

	$effect(() => {
		setFields('recoveryThreshold', Number(recoveryThresholdValue));
	});

	$effect(() => {
		setFields('notification', selectedNotificationIds);
	});

	function handleTypeChange(next: MonitorType) {
		selectedMonitorType = next;
		setFields('type', next);
		setFields('config', next === 'http' ? defaultHttpConfig() : defaultPingConfig());
	}

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
		if (failureThresholdValue === '1') return 'Immediately after the first failed check.';
		return `After ${failureThresholdValue} consecutive failed checks.`;
	});

	const recoveryThresholdHelper = $derived.by(() => {
		if (recoveryThresholdValue === '1') return 'After a single successful check.';
		return `After ${recoveryThresholdValue} consecutive successful checks.`;
	});
</script>

<form use:form class="space-y-4 w-full">

<Card.Root class="mx-auto w-full">
	<Card.Content>
			<Field.Set>
				<div class="space-y-2">
					<Field.Label for="monitor-name">Monitor name</Field.Label>
					<Input id="monitor-name" name="name" type="text" placeholder="My Monitor" required />
					<Field.Description>The name of your monitor.</Field.Description>
					{#if $errors.name}
						<Field.Description class="text-destructive">
							{$errors.name[0]}
						</Field.Description>
					{/if}
				</div>

				<div class="space-y-2">
					<Field.Label>Monitor type</Field.Label>
					<Field.Description>Select the type of monitor you want to create.</Field.Description>
					<RadioGroup.Root
						bind:value={selectedMonitorType}
						onValueChange={(value: MonitorType) => handleTypeChange(value)}
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
				</div>

				<div class="space-y-2">
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
					<input type="hidden" name="interval" value={intervalOptions[intervalIndex].seconds} />

					<div class="mt-2 flex justify-between text-xs text-muted-foreground">
						{#each intervalOptions as option (option.seconds)}
							<span class="w-6 text-center">
								{option.label}
							</span>
						{/each}
					</div>
				</div>

				<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
					<div class="space-y-2">
						<Field.Label>Failure threshold</Field.Label>
						<Field.Description>{failureThresholdHelper}</Field.Description>

						<Select.Root type="single" bind:value={failureThresholdValue}>
							<Select.Trigger class="w-full justify-between">
								<span data-slot="select-value" class="text-sm font-medium">
									{failureThresholdLabel}
								</span>
							</Select.Trigger>
							<Select.Content>
								<Select.Group>
									{#each thresholdOptions as option (option.value)}
										<Select.Item value={option.value.toString()}>{option.label}</Select.Item>
									{/each}
								</Select.Group>
							</Select.Content>
						</Select.Root>
						<input type="hidden" name="failureThreshold" value={failureThresholdValue} />
						{#if $errors.failureThreshold}
							<Field.Description class="text-destructive">
								{$errors.failureThreshold[0]}
							</Field.Description>
						{/if}
					</div>

					<div class="space-y-2">
						<Field.Label>Recovery threshold</Field.Label>
						<Field.Description>{recoveryThresholdHelper}</Field.Description>

						<Select.Root type="single" bind:value={recoveryThresholdValue}>
							<Select.Trigger class="w-full justify-between">
								<span data-slot="select-value" class="text-sm font-medium">
									{recoveryThresholdLabel}
								</span>
							</Select.Trigger>
							<Select.Content>
								<Select.Group>
									{#each thresholdOptions as option (option.value)}
										<Select.Item value={option.value.toString()}>{option.label}</Select.Item>
									{/each}
								</Select.Group>
							</Select.Content>
						</Select.Root>
						<input type="hidden" name="recoveryThreshold" value={recoveryThresholdValue} />
						{#if $errors.recoveryThreshold}
							<Field.Description class="text-destructive">
								{$errors.recoveryThreshold[0]}
							</Field.Description>
						{/if}
					</div>
				</div>

				<div class="space-y-2">
					<Field.Label>Notifications</Field.Label>
					<Field.Description>Choose where alerts should be sent.</Field.Description>
					<MultiSelect
						name="notification"
						bind:value={selectedNotificationIds}
						options={notificationOptions}
						placeholder="Select notifications"
						emptyMessage="No matching notification channels"
						maxBadges={4}
					/>
					{#if $errors.notification}
						<Field.Description class="text-destructive">
							{$errors.notification[0]}
						</Field.Description>
					{/if}
				</div>

				{#if $errors.FORM_ERROR}
					<Field.Description class="text-destructive text-center">
						{$errors.FORM_ERROR}
					</Field.Description>
				{/if}
			</Field.Set>
	</Card.Content>
</Card.Root>

{#if selectedMonitorType === 'http'}
	<HttpMonitor errors={$errors} />
{:else}
	<PingMonitor errors={$errors} />
{/if}

<div class="flex gap-2 justify-end">
	<Button size="default" variant="secondary">Cancel</Button>
	<Button size="default" variant="default">Create</Button>
</div>

</form>

