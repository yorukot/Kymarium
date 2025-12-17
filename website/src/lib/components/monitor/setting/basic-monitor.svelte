<script lang="ts">
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import { z } from 'zod';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import MultiSelect, { type MultiSelectOption } from '$lib/components/ui/multi-select';
	import type { Monitor, MonitorType } from '../../../../types';
	import type { Notification } from '../../../../types';
	import { Slider } from '$lib/components/ui/slider';
	import * as Select from '$lib/components/ui/select';
	import HttpMonitor from './http-monitor.svelte';
	import PingMonitor from './ping-monitor.svelte';
	import { Button } from '$lib/components/ui/button';
	import { createMonitor, updateMonitor } from '$lib/api/monitor';
	import { collapseStatusCodesToRanges, normalizeStatusCodes, parseHeaders } from './utils';
	import { intervalOptions, monitorTypeSelectData, thresholdOptions, httpMethods } from './setting';
	import { decidedNotificationIcon } from '$lib/utils/notification';

	const httpConfigSchema = z.object({
		url: z.url('Must be a valid URL'),
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

	type HttpConfig = z.infer<typeof httpConfigSchema>;
	type PingConfig = z.infer<typeof pingConfigSchema>;
	type MonitorFormValues = z.infer<typeof monitorFormSchema>;

	const defaultHttpConfig = (): HttpConfig => ({
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

	const defaultPingConfig = (): PingConfig => ({
		host: '',
		timeoutSeconds: 5,
		packetSize: ''
	});

	const initialIntervalIndex = 3;
	const initialFailureThreshold = '1';
	const initialRecoveryThreshold = '1';

	let {
		notifications,
		monitor = null
	}: { notifications: Notification[]; monitor?: Monitor | null } = $props();

	const isEdit = $derived.by(() => monitor !== null);

	let selectedMonitorType = $derived<MonitorType>(monitor?.type ?? 'http');
	const foundIntervalIndex = intervalOptions.findIndex(
		(opt) => opt.seconds === (monitor?.interval ?? -1)
	);
	let intervalIndex = $state<number>(
		foundIntervalIndex !== -1 ? foundIntervalIndex : initialIntervalIndex
	);
	let failureThresholdValue = $derived<string>(
		(monitor?.failureThreshold ?? Number(initialFailureThreshold)).toString()
	);
	let recoveryThresholdValue = $derived<string>(
		(monitor?.recoveryThreshold ?? Number(initialRecoveryThreshold)).toString()
	);
	let selectedNotificationIds = $derived<string[]>(monitor?.notification ?? []);

	const notificationOptions: MultiSelectOption[] = $derived.by(() =>
		notifications.map((notification) => ({
			label: notification.name,
			value: notification.id,
			keywords: [notification.type, notification.name],
			icon: decidedNotificationIcon(notification)
		}))
	);

	function toHeaderString(headers?: Record<string, string>): string {
		if (!headers) return '';
		return Object.entries(headers)
			.map(([key, value]) => `${key}: ${value}`)
			.join('\n');
	}

	function initialHttpFromMonitor(config: Record<string, unknown> | undefined): HttpConfig {
		return {
			url: (config?.url as string) ?? '',
			method: ((config?.method as HttpConfig['method']) ?? 'GET') as HttpConfig['method'],
			maxRedirects: (config?.maxRedirects as number) ?? defaultHttpConfig().maxRedirects,
			requestTimeoutSeconds:
				(config?.requestTimeout as number) ?? defaultHttpConfig().requestTimeoutSeconds,
			headers: toHeaderString(config?.headers as Record<string, string> | undefined),
			bodyEncoding: (config?.bodyEncoding as HttpConfig['bodyEncoding']) ?? '',
			body: (config?.body as string) ?? '',
			acceptedStatusCodes: (() => {
				const raw = config?.acceptedStatusCodes as number[] | undefined;
				if (raw && raw.length) {
					return collapseStatusCodesToRanges(raw);
				}
				return defaultHttpConfig().acceptedStatusCodes;
			})(),
			upsideDownMode: (config?.upsideDownMode as boolean) ?? defaultHttpConfig().upsideDownMode,
			certificateExpiryNotification:
				(config?.certificateExpiryNotification as boolean) ??
				defaultHttpConfig().certificateExpiryNotification,
			ignoreTlsError: (config?.ignoreTlsError as boolean) ?? defaultHttpConfig().ignoreTlsError
		};
	}

	function initialPingFromMonitor(config: Record<string, unknown> | undefined): PingConfig {
		return {
			host: (config?.host as string) ?? '',
			timeoutSeconds: (config?.timeoutSeconds as number) ?? defaultPingConfig().timeoutSeconds,
			packetSize: (config?.packetSize as number) ?? defaultPingConfig().packetSize ?? ''
		};
	}

	const initialValues: MonitorFormValues = (() => {
		if (!monitor) {
			return {
				name: '',
				type: 'http',
				interval: intervalOptions[initialIntervalIndex].seconds,
				failureThreshold: Number(initialFailureThreshold),
				recoveryThreshold: Number(initialRecoveryThreshold),
				notification: [] as string[],
				config: defaultHttpConfig()
			};
		}

		const base = {
			name: monitor.name,
			type: monitor.type,
			interval: monitor.interval,
			failureThreshold: monitor.failureThreshold,
			recoveryThreshold: monitor.recoveryThreshold,
			notification: monitor.notification
		};

		if (monitor.type === 'http') {
			return {
				...base,
				config: initialHttpFromMonitor(monitor.config as Record<string, unknown>)
			};
		}

		return {
			...base,
			config: initialPingFromMonitor(monitor.config as Record<string, unknown>)
		};
	})();

	const { form, errors, setFields, isSubmitting } = createForm<MonitorFormValues>({
		initialValues,
		extend: validator({ schema: monitorFormSchema }),
		onSubmit: async (values) => {
			try {
				const payload = buildPayload(values);
				console.log('monitor:create:payload', payload);
				const { params } = page;
				const teamID = params?.teamID;
				if (!teamID) {
					console.error('monitor:create missing teamID in route params', params);
					return { FORM_ERROR: 'Missing team id in route.' };
				}
				const response = isEdit
					? await updateMonitor(teamID, monitor!.id, payload)
					: await createMonitor(teamID, payload);
				console.log('monitor:save:success', response);
				await goto(`/${teamID}/monitors`);
			} catch (error) {
				console.error('monitor:save:error', error);
				const message =
					error instanceof Error ? error.message : 'Failed to save monitor. Please try again.';
				return { FORM_ERROR: message };
			}
		}
	});

	function buildPayload(values: MonitorFormValues) {
		const base = {
			name: values.name,
			type: values.type,
			interval: values.interval,
			failure_threshold: values.failureThreshold,
			recovery_threshold: values.recoveryThreshold,
			notification: values.notification
		};

		if (values.type === 'http') {
			const headers = parseHeaders(values.config.headers);
			return {
				...base,
				config: {
					url: values.config.url,
					method: values.config.method,
					max_redirects: values.config.maxRedirects,
					request_timeout: values.config.requestTimeoutSeconds,
					headers,
					body_encoding: values.config.bodyEncoding || undefined,
					body: values.config.body || undefined,
					upside_down_mode: values.config.upsideDownMode,
					certificate_expiry_notification: values.config.certificateExpiryNotification,
					ignore_tls_error: values.config.ignoreTlsError,
					accepted_status_codes: normalizeStatusCodes(values.config.acceptedStatusCodes)
				}
			};
		}

		// ping
		return {
			...base,
			config: {
				host: values.config.host,
				timeout_seconds: values.config.timeoutSeconds,
				...(values.config.packetSize === ''
					? {}
					: { packet_size: Number(values.config.packetSize) })
			}
		};
	}

	let formLevelError = $state<string | null>(null);

	$effect(() => {
		const keyed = $errors as unknown as { FORM_ERROR?: string[] | null };
		formLevelError = keyed.FORM_ERROR?.[0] ?? null;
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
		const configValue: MonitorFormValues['config'] =
			next === 'http' ? defaultHttpConfig() : defaultPingConfig();
		setFields('config', configValue);
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

	const initialHttpConfig = $derived(
		monitor?.type === 'http' ? (initialValues.config as HttpConfig) : undefined
	);
	const initialPingConfig = $derived(
		monitor?.type === 'ping' ? (initialValues.config as PingConfig) : undefined
	);
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
						onValueChange={(value: string) => handleTypeChange(value as MonitorType)}
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
			</Field.Set>
		</Card.Content>
	</Card.Root>

	{#if selectedMonitorType === 'http'}
		<HttpMonitor errors={$errors} initialConfig={initialHttpConfig} />
	{:else}
		<PingMonitor errors={$errors} initialConfig={initialPingConfig} />
	{/if}

	{#if formLevelError}
		<Field.Description class="text-destructive text-right">
			{formLevelError}
		</Field.Description>
	{/if}

	<div class="flex gap-2 justify-end">
		<Button type="submit" size="default" variant="default" disabled={$isSubmitting}>
			{$isSubmitting ? (isEdit ? 'Saving…' : 'Creating…') : isEdit ? 'Save changes' : 'Create'}
		</Button>
	</div>
</form>
