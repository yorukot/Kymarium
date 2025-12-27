<script lang="ts">
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import { z } from 'zod';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import Icon from '@iconify/svelte';

	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import Textarea from '$lib/components/ui/textarea/textarea.svelte';
	import * as Field from '$lib/components/ui/field';
	import { Button } from '$lib/components/ui/button';
	import MultiSelect, { type MultiSelectOption } from '$lib/components/ui/multi-select';
	import * as Select from '$lib/components/ui/select';
	import { Switch } from '$lib/components/ui/switch';
	import { createIncident } from '$lib/api/incident';
	import type { Incident, Monitor } from '../../types';
	import { toast } from 'svelte-sonner';

	type FormValues = {
		title: string;
		status: Incident['status'];
		severity: Incident['severity'];
		message: string;
		startedAt: string;
		isPublic: boolean;
		autoResolve: boolean;
		monitorIds: string[];
	};

	const statusOptions: Array<{ label: string; value: Incident['status'] }> = [
		{ label: 'Detected', value: 'detected' },
		{ label: 'Investigating', value: 'investigating' },
		{ label: 'Identified', value: 'identified' },
		{ label: 'Monitoring', value: 'monitoring' },
		{ label: 'Resolved', value: 'resolved' }
	];

	const severityOptions: Array<{ label: string; value: Incident['severity'] }> = [
		{ label: 'Emergency', value: 'emergency' },
		{ label: 'Critical', value: 'critical' },
		{ label: 'Major', value: 'major' },
		{ label: 'Minor', value: 'minor' },
		{ label: 'Info', value: 'info' }
	];

	let { monitors = [] }: { monitors?: Monitor[] } = $props();

	const monitorOptions: MultiSelectOption[] = $derived.by(() =>
		(monitors ?? []).map((monitor) => ({
			label: monitor.name,
			value: monitor.id,
			keywords: [monitor.type, monitor.name],
			icon: monitor.type === 'ping' ? 'lucide:waveform' : 'lucide:globe'
		}))
	);

	const initialValues: FormValues = {
		title: '',
		status: 'detected',
		severity: 'major',
		message: '',
		startedAt: toLocalInputValue(new Date()),
		isPublic: true,
		autoResolve: false,
		monitorIds: []
	};

	const schema = z.object({
		title: z.string().max(255).optional().default(''),
		status: z.enum(['detected', 'investigating', 'identified', 'monitoring', 'resolved']),
		severity: z.enum(['emergency', 'critical', 'major', 'minor', 'info']),
		message: z.string().max(2000).optional().default(''),
		startedAt: z.string().min(1, 'Start time is required'),
		isPublic: z.boolean(),
		autoResolve: z.boolean(),
		monitorIds: z.array(z.string().min(1)).min(1, 'Select at least one monitor')
	});

	const { form, errors, isSubmitting } = createForm<FormValues>({
		initialValues,
		extend: validator({ schema }),
		onSubmit: handleSubmit
	});

	let status = $state<Incident['status']>(initialValues.status);
	let severity = $state<Incident['severity']>(initialValues.severity);
	let title = $state(initialValues.title);
	let message = $state(initialValues.message);
	let startedAt = $state(initialValues.startedAt);
	let isPublic = $state(initialValues.isPublic);
	let autoResolve = $state(initialValues.autoResolve);
	let monitorIds = $state<string[]>(initialValues.monitorIds);

	function toLocalInputValue(date: Date): string {
		// datetime-local expects no timezone; use local time and trim seconds
		const iso = new Date(date.getTime() - date.getTimezoneOffset() * 60000).toISOString();
		return iso.slice(0, 16);
	}

	function toIsoOrUndefined(value?: string): string | undefined {
		if (!value) return undefined;
		const parsed = new Date(value);
		if (Number.isNaN(parsed.getTime())) return undefined;
		return parsed.toISOString();
	}

	async function handleSubmit(values: FormValues) {
		const teamID = page.params.teamID;
		if (!teamID) {
			toast.error('Missing team id in route');
			return { FORM_ERROR: 'Missing team id in route.' };
		}

		try {
			const payload = {
				title: values.title?.trim() || undefined,
				status: values.status,
				severity: values.severity,
				message: values.message || undefined,
				started_at: toIsoOrUndefined(values.startedAt),
				public: values.isPublic,
				auto_resolve: values.autoResolve,
				monitor_ids: values.monitorIds
			};

			await createIncident(teamID, payload);
			toast.success('Incident created');
			await goto(`/${teamID}/incidents`);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : 'Failed to create incident. Please try again.';
			toast.error(message);
			return { FORM_ERROR: message };
		}
	}
</script>

<Card.Root class="w-full">
	<Card.Header>
		<Card.Title class="text-xl">Create incident</Card.Title>
		<Card.Description>
			Link at least one monitor and choose severity. Defaults mirror automatic incidents.
		</Card.Description>
	</Card.Header>

	<Card.Content>
		<form class="space-y-6" use:form>
			<Field.Set class="grid gap-4 md:grid-cols-2">
				<div class="space-y-2">
					<Field.Label>Status</Field.Label>
					<Select.Root type="single" bind:value={status}>
						<Select.Trigger class="justify-between">
							<span data-slot="select-value" class="text-sm font-medium capitalize">
								{status}
							</span>
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								{#each statusOptions as option (option.value)}
									<Select.Item value={option.value}>{option.label}</Select.Item>
								{/each}
							</Select.Group>
						</Select.Content>
					</Select.Root>
					<input type="hidden" name="status" value={status} />
					{#if $errors.status}
						<Field.Description class="text-destructive">{$errors.status[0]}</Field.Description>
					{/if}
				</div>

				<div class="space-y-2">
					<Field.Label>Severity</Field.Label>
					<Select.Root type="single" bind:value={severity}>
						<Select.Trigger class="justify-between">
							<span data-slot="select-value" class="text-sm font-medium capitalize">
								{severity}
							</span>
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								{#each severityOptions as option (option.value)}
									<Select.Item value={option.value}>{option.label}</Select.Item>
								{/each}
							</Select.Group>
						</Select.Content>
					</Select.Root>
					<input type="hidden" name="severity" value={severity} />
					{#if $errors.severity}
						<Field.Description class="text-destructive">{$errors.severity[0]}</Field.Description>
					{/if}
				</div>
			</Field.Set>

			<div class="space-y-2">
				<Field.Label for="title">Incident title (optional)</Field.Label>
				<Input
					id="title"
					name="title"
					placeholder="Short summary, e.g. API latency spike"
					bind:value={title}
				/>
				{#if $errors.title}
					<Field.Description class="text-destructive">{$errors.title[0]}</Field.Description>
				{/if}
			</div>

			<div class="space-y-2">
				<Field.Label for="message">Message (optional)</Field.Label>
				<Textarea
					id="message"
					name="message"
					placeholder="Describe what happened or why you are creating this incident."
					bind:value={message}
					rows={3}
				/>
				{#if $errors.message}
					<Field.Description class="text-destructive">{$errors.message[0]}</Field.Description>
				{/if}
			</div>

			<div class="space-y-2">
				<Field.Label>Monitors</Field.Label>
				<MultiSelect
					name="monitorIds"
					placeholder="Select monitors"
					options={monitorOptions}
					bind:value={monitorIds}
					emptyMessage="No monitors available"
				/>
				<Field.Description class={$errors.monitorIds ? 'text-destructive' : 'text-muted-foreground'}>
					{$errors.monitorIds ? $errors.monitorIds[0] : 'Select all monitors affected by this incident.'}
				</Field.Description>
			</div>

			<Field.Set class="grid gap-4 md:grid-cols-3">
				<div class="space-y-2">
					<Field.Label for="started-at">Started at</Field.Label>
					<Field.Description>Defaults to now; adjust to match detection time.</Field.Description>
					<Input
						id="started-at"
						name="startedAt"
						type="datetime-local"
						bind:value={startedAt}
						required
					/>
					{#if $errors.startedAt}
						<Field.Description class="text-destructive">{$errors.startedAt[0]}</Field.Description>
					{/if}
				</div>

				<div class="space-y-2">
					<Field.Label>Public status page</Field.Label>
					<div class="flex items-center gap-3 rounded-md border px-3 py-2">
						<Switch bind:checked={isPublic} />
						<input type="checkbox" class="hidden" name="isPublic" bind:checked={isPublic} />
						<div>
							<p class="text-sm font-medium">Visible on status page</p>
							<p class="text-xs text-muted-foreground">Toggle off to keep internal only.</p>
						</div>
					</div>
				</div>

				<div class="space-y-2">
					<Field.Label>Auto-resolve</Field.Label>
					<div class="flex items-center gap-3 rounded-md border px-3 py-2">
						<Switch bind:checked={autoResolve} />
						<input type="checkbox" class="hidden" name="autoResolve" bind:checked={autoResolve} />
						<div>
							<p class="text-sm font-medium">Resolve automatically</p>
							<p class="text-xs text-muted-foreground">Close once monitors recover.</p>
						</div>
					</div>
				</div>
			</Field.Set>

			{#if $errors.FORM_ERROR}
				<Field.Description class="text-destructive text-center">
					{$errors.FORM_ERROR}
				</Field.Description>
			{/if}

			<div class="flex justify-end gap-2">
				<Button type="button" variant="ghost" href="../incidents">
					<Icon icon="lucide:arrow-left" />
					Back
				</Button>
				<Button type="submit" disabled={$isSubmitting}>
					{#if $isSubmitting}
						Creating…
					{:else}
						Create incident
					{/if}
				</Button>
			</div>
		</form>
	</Card.Content>
</Card.Root>
