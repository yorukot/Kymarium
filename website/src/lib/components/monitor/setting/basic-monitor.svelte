<script lang="ts">
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import { reporter, ValidationMessage } from '@felte/reporter-svelte';
	import { z } from 'zod';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import MultiSelect, { type MultiSelectOption } from '$lib/components/ui/multi-select';
	import type { Monitor, MonitorType, Notification, Region } from '../$lib/types';
	import { Slider } from '$lib/components/ui/slider';
	import * as Select from '$lib/components/ui/select';
	import HttpMonitor from './http-monitor.svelte';
	import PingMonitor from './ping-monitor.svelte';
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { createMonitor, updateMonitor, deleteMonitor } from '$lib/api/monitor';
	import { collapseStatusCodesToRanges, normalizeStatusCodes, parseHeaders } from './utils';
	import { intervalOptions, monitorTypeSelectData, thresholdOptions, httpMethods } from './setting';
	import { decidedNotificationIcon } from '$lib/utils/notification';
	import { regionFlagIcon } from '$lib/utils/region';
	import DeleteMonitorDialog from '../delete-monitor-dialog.svelte';
	import { toast } from 'svelte-sonner';

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
		regions: z.array(z.string()).min(1, 'Select at least one region'),
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
		regions,
		monitor = null
	}: {
		notifications?: Notification[] | null;
		regions?: Region[] | null;
		monitor?: Monitor | null;
	} = $props();

	const isEdit = $derived.by(() => monitor !== null);
	let confirmOpen = $state(false);
	let isDeleting = $state(false);

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
	let selectedRegionIds = $derived<string[]>(monitor?.regions ?? []);

	const safeNotifications = $derived(notifications ?? []);
	const notificationOptions: MultiSelectOption[] = $derived.by(() =>
		safeNotifications.map((notification) => ({
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

	function initialHttpFromMonitor(config: unknown): HttpConfig {
		const cfg = (config ?? {}) as Record<string, unknown>;
		return {
			url: (cfg.url as string) ?? '',
			method: ((cfg.method as HttpConfig['method']) ?? 'GET') as HttpConfig['method'],
			maxRedirects: (cfg.max_redirects as number) ?? defaultHttpConfig().maxRedirects,
			requestTimeoutSeconds:
				(cfg.request_timeout as number) ?? defaultHttpConfig().requestTimeoutSeconds,
			headers: toHeaderString(cfg.headers as Record<string, string> | undefined),
			bodyEncoding: (cfg.body_encoding as HttpConfig['bodyEncoding']) ?? '',
			body: (cfg.body as string) ?? '',
			acceptedStatusCodes: (() => {
				const raw = cfg.accepted_status_codes as number[] | undefined;
				if (raw && raw.length) {
					return collapseStatusCodesToRanges(raw);
				}
				return defaultHttpConfig().acceptedStatusCodes;
			})(),
			upsideDownMode: (cfg.upside_down_mode as boolean) ?? defaultHttpConfig().upsideDownMode,
			certificateExpiryNotification:
				(cfg.certificate_expiry_notification as boolean) ??
				defaultHttpConfig().certificateExpiryNotification,
			ignoreTlsError: (cfg.ignore_tls_error as boolean) ?? defaultHttpConfig().ignoreTlsError
		};
	}

	function initialPingFromMonitor(config: unknown): PingConfig {
		const cfg = (config ?? {}) as Record<string, unknown>;
		return {
			host: (cfg.host as string) ?? '',
			timeoutSeconds: (cfg.timeout_seconds as number) ?? defaultPingConfig().timeoutSeconds,
			packetSize: (cfg.packet_size as number) ?? defaultPingConfig().packetSize ?? ''
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
				regions: [] as string[],
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
			regions: monitor.regions,
			notification: monitor.notification
		};

		if (monitor.type === 'http') {
			return {
				...base,
				config: initialHttpFromMonitor(monitor.config)
			};
		}

		return {
			...base,
			config: initialPingFromMonitor(monitor.config)
		};
	})();

	const { form, setFields, isSubmitting } = createForm<MonitorFormValues>({
		initialValues,
		extend: [validator({ schema: monitorFormSchema }), reporter()],
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
			regions: values.regions,
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

	$effect(() => {
		setFields('regions', selectedRegionIds);
	});

	function handleTypeChange(next: MonitorType) {
		selectedMonitorType = next;
		setFields('type', next);
		const configValue: MonitorFormValues['config'] =
			next === 'http' ? defaultHttpConfig() : defaultPingConfig();
		setFields('config', configValue);
	}

	function toggleRegion(regionId: string) {
		selectedRegionIds = selectedRegionIds.includes(regionId)
			? selectedRegionIds.filter((id) => id !== regionId)
			: [...selectedRegionIds, regionId];
		setFields('regions', selectedRegionIds);
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

	function askDelete() {
		if (!monitor) return;
		confirmOpen = true;
	}

	async function handleDelete() {
		if (!monitor) return;
		const teamID = page.params.teamID;
		if (!teamID) {
			toast.error('Missing team id');
			return;
		}

		isDeleting = true;
		try {
			await deleteMonitor(teamID, monitor.id);
			toast.success('Monitor deleted');
			await goto(`/${teamID}/monitors`);
		} catch (error) {
			console.error('monitor:delete:error', error);
			const message =
				error instanceof Error ? error.message : 'Failed to delete monitor. Please try again.';
			toast.error(message);
		} finally {
			isDeleting = false;
			confirmOpen = false;
		}
	}
</script>

<form use:form class="space-y-4 w-full">
	<Card.Root class="mx-auto w-full">
		<Card.Content>
			<Field.Set>
				<div class="space-y-2">
					<Field.Label for="monitor-name">Monitor name</Field.Label>
					<Input id="monitor-name" name="name" type="text" placeholder="My Monitor" required />
					<Field.Description>The name of your monitor.</Field.Description>
					<ValidationMessage for="name" let:messages>
						{#if messages?.length}
							<Field.Description class="text-destructive">
								{messages[0]}
							</Field.Description>
						{/if}
					</ValidationMessage>
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
						<ValidationMessage for="failureThreshold" let:messages>
							{#if messages?.length}
								<Field.Description class="text-destructive">
									{messages[0]}
								</Field.Description>
							{/if}
						</ValidationMessage>
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
						<ValidationMessage for="recoveryThreshold" let:messages>
							{#if messages?.length}
								<Field.Description class="text-destructive">
									{messages[0]}
								</Field.Description>
							{/if}
						</ValidationMessage>
					</div>
				</div>

				<div class="space-y-2">
					<Field.Label>Regions</Field.Label>
					<Field.Description>Select at least one region to run checks from.</Field.Description>
					<div class="grid gap-2 sm:grid-cols-2">
						{#each regions as region (region.id)}
							{@const flagIcon = regionFlagIcon(region)}
							<button
								type="button"
								class="flex items-start gap-3 rounded-md border px-3 py-2 text-left transition hover:border-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
								onclick={() => toggleRegion(region.id)}
							>
								<Checkbox checked={selectedRegionIds.includes(region.id)} class="mt-1" />
								{#if flagIcon}
									<Icon icon={flagIcon} class="mt-0.5 size-5 shrink-0" aria-hidden="true" />
								{/if}
								<div class="space-y-0.5">
									<Field.Title>{region.displayName ?? region.name}</Field.Title>
									<Field.Description class="text-xs text-muted-foreground">
										{region.name}
									</Field.Description>
								</div>
							</button>
						{/each}
					</div>
					<ValidationMessage for="regions" let:messages>
						{#if messages?.length}
							<Field.Description class="text-destructive">
								{messages[0]}
							</Field.Description>
						{/if}
					</ValidationMessage>
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
					<ValidationMessage for="notification" let:messages>
						{#if messages?.length}
							<Field.Description class="text-destructive">
								{messages[0]}
							</Field.Description>
						{/if}
					</ValidationMessage>
				</div>
			</Field.Set>
		</Card.Content>
	</Card.Root>

	{#if selectedMonitorType === 'http'}
		<HttpMonitor initialConfig={initialHttpConfig} />
	{:else}
		<PingMonitor initialConfig={initialPingConfig} />
	{/if}

	<ValidationMessage for="FORM_ERROR" let:messages>
		{#if messages?.length}
			<Field.Description class="text-destructive text-right">
				{messages[0]}
			</Field.Description>
		{/if}
	</ValidationMessage>

	<div class="flex gap-2 justify-end">
		{#if isEdit}
			<Button
				type="button"
				variant="destructive"
				disabled={$isSubmitting || isDeleting}
				onclick={askDelete}
			>
				{isDeleting ? 'Deleting…' : 'Delete'}
			</Button>
		{/if}
		<Button type="submit" size="default" variant="default" disabled={$isSubmitting}>
			{$isSubmitting ? (isEdit ? 'Saving…' : 'Creating…') : isEdit ? 'Save changes' : 'Create'}
		</Button>
	</div>

	{#if isEdit}
		<DeleteMonitorDialog
			bind:open={confirmOpen}
			{monitor}
			onConfirm={handleDelete}
			loading={isDeleting}
		/>
	{/if}
</form>
