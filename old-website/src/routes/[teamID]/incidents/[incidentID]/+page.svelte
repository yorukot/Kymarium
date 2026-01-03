<script lang="ts">
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import * as Card from '$lib/components/ui/card';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Field from '$lib/components/ui/field';
	import * as Select from '$lib/components/ui/select';
	import { Switch } from '$lib/components/ui/switch';
	import { Input } from '$lib/components/ui/input';
	import Textarea from '$lib/components/ui/textarea/textarea.svelte';
	import { createIncidentEvent, updateIncident, updateIncidentStatus } from '$lib/api/incident';
	import { toast } from 'svelte-sonner';

	import type { Incident, IncidentEvent, IncidentEventType } from '../../../../lib/types/index.js';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();

	let incident = $derived<Incident>(data.incident);
	let events = $derived<IncidentEvent[]>(data.events ?? []);
	let monitorNames = $derived<string[]>(data.monitorNames ?? []);

	let status = $derived<Incident['status']>(data.incident.status);
	let statusMessage = $state('');

	let isPublic = $derived(data.incident.isPublic);
	let autoResolve = $derived(data.incident.autoResolve);
	let title = $derived(data.incident.title ?? '');
	let isSubmittingSettings = $state(false);

	let eventType = $state<IncidentEventType>('update');
	let eventMessage = $state('');

	let isSubmittingStatus = $state(false);
	let isSubmittingEvent = $state(false);

	const statusOptions: Array<{ label: string; value: Incident['status'] }> = [
		{ label: 'Detected', value: 'detected' },
		{ label: 'Investigating', value: 'investigating' },
		{ label: 'Identified', value: 'identified' },
		{ label: 'Monitoring', value: 'monitoring' },
		{ label: 'Resolved', value: 'resolved' }
	];

	const eventTypeOptions: Array<{ label: string; value: IncidentEventType }> = [
		{ label: 'Update', value: 'update' },
		{ label: 'Detected', value: 'detected' },
		{ label: 'Investigating', value: 'investigating' },
		{ label: 'Identified', value: 'identified' },
		{ label: 'Monitoring', value: 'monitoring' },
		{ label: 'Notification sent', value: 'notification_sent' },
		{ label: 'Published', value: 'published' },
		{ label: 'Unpublished', value: 'unpublished' },
		{ label: 'Auto resolved', value: 'auto_resolved' },
		{ label: 'Manually resolved', value: 'manually_resolved' }
	];

	const orderedEvents = $derived(
		[...events].sort(
			(a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
		)
	);

	function formatDate(value?: string) {
		if (!value) return '—';
		const d = new Date(value);
		return Number.isNaN(d.getTime()) ? '—' : d.toLocaleString();
	}

	function formatDuration(target: Incident) {
		const start = new Date(target.startedAt).getTime();
		const end = new Date(target.resolvedAt ?? Date.now()).getTime();
		if (!Number.isFinite(start) || !Number.isFinite(end) || end < start) return '—';
		const diffMs = end - start;
		const minutes = Math.floor(diffMs / 60000);
		if (minutes < 60) return `${minutes}m`;
		const hours = Math.floor(minutes / 60);
		const mins = minutes % 60;
		if (hours < 24) return `${hours}h ${mins}m`;
		const days = Math.floor(hours / 24);
		return `${days}d ${hours % 24}h`;
	}

	const severityMeta: Record<Incident['severity'], { label: string; badge: string; dot: string }> = {
		emergency: {
			label: 'Emergency',
			badge: '!bg-destructive !text-destructive-foreground border-transparent',
			dot: 'bg-destructive'
		},
		critical: {
			label: 'Critical',
			badge: '!bg-destructive/80 !text-destructive-foreground border-transparent',
			dot: 'bg-destructive/80'
		},
		major: {
			label: 'Major',
			badge: '!bg-amber-500 !text-amber-950 border-transparent',
			dot: 'bg-amber-500'
		},
		minor: {
			label: 'Minor',
			badge: '!bg-secondary !text-secondary-foreground border-transparent',
			dot: 'bg-secondary'
		},
		info: {
			label: 'Info',
			badge: '!bg-muted !text-foreground border-transparent',
			dot: 'bg-muted'
		}
	};

	const statusMeta: Record<Incident['status'], { label: string; color: string }> = {
		detected: { label: 'Detected', color: 'text-amber-600' },
		investigating: { label: 'Investigating', color: 'text-amber-700' },
		identified: { label: 'Identified', color: 'text-amber-700' },
		monitoring: { label: 'Monitoring', color: 'text-sky-700' },
		resolved: { label: 'Resolved', color: 'text-success' }
	};

	const eventTypeMeta: Record<IncidentEventType, { label: string; badge: string }> = {
		detected: { label: 'Detected', badge: '!bg-amber-100 !text-amber-900 border-transparent' },
		investigating: { label: 'Investigating', badge: '!bg-amber-100 !text-amber-900 border-transparent' },
		identified: { label: 'Identified', badge: '!bg-amber-100 !text-amber-900 border-transparent' },
		monitoring: { label: 'Monitoring', badge: '!bg-sky-100 !text-sky-900 border-transparent' },
		update: { label: 'Update', badge: '!bg-secondary !text-secondary-foreground border-transparent' },
		notification_sent: {
			label: 'Notification sent',
			badge: '!bg-muted !text-foreground border-transparent'
		},
		published: { label: 'Published', badge: '!bg-emerald-100 !text-emerald-900 border-transparent' },
		unpublished: { label: 'Unpublished', badge: '!bg-muted !text-foreground border-transparent' },
		auto_resolved: {
			label: 'Auto resolved',
			badge: '!bg-emerald-100 !text-emerald-900 border-transparent'
		},
		manually_resolved: {
			label: 'Manually resolved',
			badge: '!bg-emerald-100 !text-emerald-900 border-transparent'
		}
	};

	const eventTypeMetaSafe = (type: IncidentEventType) =>
		eventTypeMeta[type] ?? {
			label: type,
			badge: '!bg-muted !text-foreground border-transparent'
		};

	const dotClass = (severity: Incident['severity']) => severityMeta[severity]?.dot ?? 'bg-muted';
	const badgeClass = (severity: Incident['severity']) =>
		severityMeta[severity]?.badge ?? '!bg-muted !text-foreground border-transparent';

	async function handleStatusSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (isSubmittingStatus) return;

		const teamID = page.params.teamID;
		if (!teamID) {
			toast.error('Missing team id in route');
			return;
		}

		isSubmittingStatus = true;
		try {
			const res = await updateIncidentStatus(teamID, incident.id, {
				status,
				message: statusMessage.trim() || undefined
			});

			incident = res.data.incident;
			status = res.data.incident.status;
			isPublic = res.data.incident.isPublic;
			autoResolve = res.data.incident.autoResolve;
			title = res.data.incident.title ?? '';
			events = [res.data.event, ...events];
			statusMessage = '';
			toast.success('Incident status updated');
		} catch (error) {
			const message =
				error instanceof Error ? error.message : 'Failed to update incident status.';
			toast.error(message);
		} finally {
			isSubmittingStatus = false;
		}
	}

	async function handleEventSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (isSubmittingEvent) return;

		const teamID = page.params.teamID;
		if (!teamID) {
			toast.error('Missing team id in route');
			return;
		}

		if (!eventMessage.trim()) {
			toast.error('Add a timeline message to continue');
			return;
		}

		isSubmittingEvent = true;
		try {
			const res = await createIncidentEvent(teamID, incident.id, {
				message: eventMessage.trim(),
				event_type: eventType
			});

			events = [res.data, ...events];
			eventMessage = '';
			toast.success('Timeline update added');
		} catch (error) {
			const message =
				error instanceof Error ? error.message : 'Failed to add timeline update.';
			toast.error(message);
		} finally {
			isSubmittingEvent = false;
		}
	}

	async function handleSettingsSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (isSubmittingSettings) return;

		const teamID = page.params.teamID;
		if (!teamID) {
			toast.error('Missing team id in route');
			return;
		}

		isSubmittingSettings = true;
		try {
			const res = await updateIncident(teamID, incident.id, {
				title: title.trim() || null,
				public: isPublic,
				auto_resolve: autoResolve
			});

			incident = res.data;
			isPublic = res.data.isPublic;
			autoResolve = res.data.autoResolve;
			title = res.data.title ?? '';
			toast.success('Incident settings updated');
		} catch (error) {
			const message = error instanceof Error ? error.message : 'Failed to update incident settings.';
			toast.error(message);
		} finally {
			isSubmittingSettings = false;
		}
	}
</script>

<div class="flex flex-col gap-6">
	<header class="flex flex-col gap-3">
		<div class="flex items-start justify-between gap-3 flex-wrap">
			<div class="flex flex-col gap-1">
				<p class="text-sm text-muted-foreground">Incident detail</p>
				<div class="flex items-center gap-2 flex-wrap">
					<h1 class="text-2xl font-bold">
						{incident.title ?? `Incident #${incident.id}`}
					</h1>
					{#if incident.title}
						<span class="text-sm text-muted-foreground">#{incident.id}</span>
					{/if}
					<Badge class={badgeClass(incident.severity)}>{severityMeta[incident.severity].label}</Badge>
					<span class={`text-sm font-medium ${statusMeta[incident.status].color}`}>
						{statusMeta[incident.status].label}
					</span>
					{#if incident.isPublic}
						<Badge variant="outline" class="gap-1">
							<Icon icon="lucide:globe-2" class="size-3.5" /> Public
						</Badge>
					{/if}
					{#if incident.autoResolve}
						<Badge variant="outline" class="gap-1">
							<Icon icon="lucide:clock-3" class="size-3.5" /> Auto-resolve
						</Badge>
					{/if}
				</div>
			</div>
			<Button variant="ghost" href="../incidents">
				<Icon icon="lucide:arrow-left" />
				Back to incidents
			</Button>
		</div>
		<div class="text-sm text-muted-foreground flex gap-3 flex-wrap">
			<span class="flex items-center gap-1">
				<Icon icon="lucide:play-circle" class="size-4" />
				Started {formatDate(incident.startedAt)}
			</span>
			<span class="flex items-center gap-1">
				<Icon icon="lucide:flag" class="size-4" />
				Resolved {formatDate(incident.resolvedAt)}
			</span>
			<span class="flex items-center gap-1">
				<Icon icon="lucide:timer" class="size-4" />
				Duration {formatDuration(incident)}
			</span>
			<span class="flex items-center gap-1">
				<Icon icon="lucide:clock" class="size-4" />
				Updated {formatDate(incident.updatedAt)}
			</span>
		</div>
	</header>

	<div class="grid gap-6 lg:grid-cols-[1.1fr_0.9fr]">
		<div class="flex flex-col gap-6">
			<Card.Root class="p-6">
				<Card.Header class="p-0">
					<Card.Title class="text-lg">Affected monitors</Card.Title>
					<Card.Description>Linked monitors for this incident.</Card.Description>
				</Card.Header>
				<Card.Content class="p-0 flex flex-col gap-3">
					{#if monitorNames.length === 0}
						<p class="text-sm text-muted-foreground">No monitor links found for this incident.</p>
					{:else}
						<div class="flex flex-wrap gap-2">
							{#each monitorNames as name (name)}
								<Badge variant="secondary" class="truncate max-w-56">{name}</Badge>
							{/each}
						</div>
					{/if}
				</Card.Content>
			</Card.Root>

			<Card.Root class="p-6">
				<Card.Header class="p-0">
					<Card.Title class="text-lg">Timeline</Card.Title>
					<Card.Description>Latest updates for this incident.</Card.Description>
				</Card.Header>
				<Card.Content class="p-0">
					{#if orderedEvents.length === 0}
						<p class="text-sm text-muted-foreground">No timeline updates yet.</p>
					{:else}
						<div class="relative mt-4 space-y-4">
							{#each orderedEvents as event (event.id)}
									<div class="flex flex-col gap-1">
										<div class="flex flex-wrap items-center justify-between gap-2">
											<div class="flex flex-wrap items-center gap-2">
												<span class="font-semibold text-sm">
													{eventTypeMetaSafe(event.eventType).label}:
												</span>
												<span class="text-sm font-medium">{event.message || '—'}</span>
											</div>
											<span class="text-xs text-muted-foreground">
												{formatDate(event.createdAt)}
											</span>
										</div>
										<p class="text-xs text-muted-foreground">
											{event.createdBy ? `By ${event.createdBy}` : 'System update'}
										</p>
									</div>
							{/each}
						</div>
					{/if}
				</Card.Content>
			</Card.Root>
		</div>

		<div class="flex flex-col gap-6">
			<Card.Root class="p-6">
				<Card.Header class="p-0">
					<Card.Title class="text-lg">Update status</Card.Title>
					<Card.Description>Change the incident status and record a timeline event.</Card.Description>
				</Card.Header>
				<Card.Content class="p-0">
					<form class="space-y-4" onsubmit={handleStatusSubmit}>
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
						</div>

						<div class="space-y-2">
							<Field.Label for="status-message">Message (optional)</Field.Label>
							<Textarea
								id="status-message"
								name="statusMessage"
								rows={3}
								placeholder="Add context to the status update."
								bind:value={statusMessage}
							/>
						</div>

						<div class="flex justify-end">
							<Button type="submit" disabled={isSubmittingStatus}>
								{#if isSubmittingStatus}
									Updating…
								{:else}
									Update status
								{/if}
							</Button>
						</div>
					</form>
				</Card.Content>
			</Card.Root>

			<Card.Root class="p-6">
				<Card.Header class="p-0">
					<Card.Title class="text-lg">Visibility & automation</Card.Title>
					<Card.Description>Control incident title, visibility, and auto-resolve.</Card.Description>
				</Card.Header>
				<Card.Content class="p-0">
					<form class="space-y-4" onsubmit={handleSettingsSubmit}>
						<div class="space-y-2">
							<Field.Label for="incident-title">Title (optional)</Field.Label>
							<Input
								id="incident-title"
								name="title"
								placeholder="Short summary for this incident"
								bind:value={title}
							/>
						</div>

						<div class="flex flex-col gap-3">
							<div class="flex items-center justify-between gap-3 rounded-md border px-3 py-2">
								<div>
									<p class="text-sm font-medium">Public incident</p>
									<p class="text-xs text-muted-foreground">Show this incident on the status page.</p>
								</div>
								<Switch bind:checked={isPublic} />
							</div>

							<div class="flex items-center justify-between gap-3 rounded-md border px-3 py-2">
								<div>
									<p class="text-sm font-medium">Auto-resolve</p>
									<p class="text-xs text-muted-foreground">
										Resolve automatically when monitors recover.
									</p>
								</div>
								<Switch bind:checked={autoResolve} />
							</div>
						</div>

						<div class="flex justify-end">
							<Button type="submit" disabled={isSubmittingSettings}>
								{#if isSubmittingSettings}
									Saving…
								{:else}
									Save settings
								{/if}
							</Button>
						</div>
					</form>
				</Card.Content>
			</Card.Root>

			<Card.Root class="p-6">
				<Card.Header class="p-0">
					<Card.Title class="text-lg">Post update</Card.Title>
					<Card.Description>Add a timeline entry without changing status.</Card.Description>
				</Card.Header>
				<Card.Content class="p-0">
					<form class="space-y-4" onsubmit={handleEventSubmit}>
						<div class="space-y-2">
							<Field.Label>Event type</Field.Label>
							<Select.Root type="single" bind:value={eventType}>
								<Select.Trigger class="justify-between">
									<span data-slot="select-value" class="text-sm font-medium">
										{eventTypeMetaSafe(eventType).label}
									</span>
								</Select.Trigger>
								<Select.Content>
									<Select.Group>
										{#each eventTypeOptions as option (option.value)}
											<Select.Item value={option.value}>{option.label}</Select.Item>
										{/each}
									</Select.Group>
								</Select.Content>
							</Select.Root>
						</div>

						<div class="space-y-2">
							<Field.Label for="event-message">Message</Field.Label>
							<Textarea
								id="event-message"
								name="eventMessage"
								rows={4}
								placeholder="Share the latest update with your team."
								bind:value={eventMessage}
							/>
						</div>

						<div class="flex justify-end">
							<Button type="submit" disabled={isSubmittingEvent || !eventMessage.trim()}>
								{#if isSubmittingEvent}
									Posting…
								{:else}
									Add update
								{/if}
							</Button>
						</div>
					</form>
				</Card.Content>
			</Card.Root>
		</div>
	</div>
</div>
