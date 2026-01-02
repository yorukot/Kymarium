<script lang="ts">
	import type { MonitorWithIncidents } from '$lib/types';
	import { Card } from '$lib/components/ui/card';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import Icon from '@iconify/svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { deleteMonitor } from '$lib/api/monitor';
	import { toast } from 'svelte-sonner';
	import { page } from '$app/state';
	import DeleteMonitorDialog from './delete-monitor-dialog.svelte';
	import { goto } from '$app/navigation';
	import Button from '../ui/button/button.svelte';
	import { monitorTarget } from '$lib/utils/monitor';
	import { statusMeta } from '$lib/styles/status';

	let { monitors }: { monitors: MonitorWithIncidents[] } = $props();

	let confirmOpen = $state(false);
	let deleting = $state(false);
	let targetMonitor: MonitorWithIncidents | null = $state(null);

	function askDelete(monitor: MonitorWithIncidents) {
		targetMonitor = monitor;
		confirmOpen = true;
	}

	const formatRelativeTime = (isoDate: string | undefined) => {
		if (!isoDate) return '–';
		const parsed = new Date(isoDate);
		if (Number.isNaN(parsed.getTime())) return '–';

		const diffSeconds = Math.floor((Date.now() - parsed.getTime()) / 1000);
		if (diffSeconds < 0) return 'just now';
		if (diffSeconds < 60) return `${diffSeconds}s ago`;
		const diffMinutes = Math.floor(diffSeconds / 60);
		if (diffMinutes < 60) return `${diffMinutes}m ago`;
		const diffHours = Math.floor(diffMinutes / 60);
		if (diffHours < 24) return `${diffHours}h ago`;
		const diffDays = Math.floor(diffHours / 24);
		return `${diffDays}d ago`;
	};

	const formatPercent = (value: number | undefined) =>
		Number.isFinite(value) ? `${(value ?? 0).toFixed(2)}%` : '–';

	async function handleDelete() {
		if (!targetMonitor) return;
		const teamID = page.params.teamID;
		if (!teamID) {
			toast.error('Missing team id');
			return;
		}

		deleting = true;
		try {
			await deleteMonitor(teamID, targetMonitor.id);
			monitors = monitors.filter((m) => m.id !== targetMonitor?.id);
			toast.success('Monitor deleted');
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Failed to delete monitor';
			toast.error(message);
		} finally {
			deleting = false;
			confirmOpen = false;
			targetMonitor = null;
		}
	}
</script>

<div>
	<div class="flex flex-col gap-2">
		{#if monitors.length === 0}
			<Card class="py-6 px-4 text-sm text-muted-foreground">No monitors yet.</Card>
		{/if}
		{#each monitors as monitor, monitorIndex (monitor.id)}
			{@const openIncidents = (monitor.incidents ?? []).filter(
				(incident) => incident.status !== 'resolved'
			)}
			<Card class="p-0 gap-0 border-none">
				<Card id={'monitor-card-' + monitorIndex} class="py-2 px-4">
					<div class="flex justify-between items-center">
						<div class="flex-1 min-w-0">
							<div class="flex items-center gap-2 min-w-0">
								<span
									class={`w-3 h-3 rounded-full shrink-0 ${statusMeta[monitor.status]?.dot ?? 'bg-muted'}`}
								></span>

								<h2 class="text-md font-semibold truncate min-w-0 flex-1">
									{monitor.name}
								</h2>
							</div>

							<div class="min-w-0">
								<p class="text-sm text-muted-foreground truncate">
									<Badge variant="secondary" class="rounded-sm shrink-0">
										{monitor.type.toUpperCase()}
									</Badge>
									{monitorTarget(monitor)}
								</p>
							</div>
						</div>
						<div class="flex items-center gap-4">
							<div class="text-right hidden md:block">
								<div class="text-xs text-muted-foreground">Last check</div>
								<div class="text-sm font-medium">
									{formatRelativeTime(monitor.lastChecked)}
								</div>
							</div>
							<div class="text-right">
								<div class="text-xs text-muted-foreground">Uptime (30d)</div>
								<div class="text-sm font-medium">{formatPercent(monitor.uptimeSli30)}</div>
							</div>
							<DropdownMenu.Root>
								<DropdownMenu.Trigger class="shrink-0">
									<Button variant="ghost" size="icon">
										<Icon icon="lucide:more-vertical" />
									</Button>
								</DropdownMenu.Trigger>
								<DropdownMenu.Content>
									<DropdownMenu.Group>
										<DropdownMenu.Item onclick={() => goto(`monitors/${monitor.id}`)}>
											<Icon icon="lucide:eye" /> View details
										</DropdownMenu.Item>
										<DropdownMenu.Item onclick={() => goto(`monitors/${monitor.id}/edit`)}>
											<Icon icon="lucide:edit" /> Edit
										</DropdownMenu.Item>
										<DropdownMenu.Separator />
										<DropdownMenu.Item variant="destructive" onclick={() => askDelete(monitor)}>
											<Icon icon="lucide:trash" />
											Delete
										</DropdownMenu.Item>
									</DropdownMenu.Group>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</div>
					</div>
				</Card>
				{#if openIncidents.length > 0}
					<div class="px-4 py-2">
						<div class="mt-1 flex flex-col gap-1">
							{#each openIncidents as incident (incident.id)}
								<div class="flex items-center gap-2 text-xs justify-between">
									<div>
										<span class="text-lg text-foreground flex gap-2 items-center">
											<Icon icon="lucide:activity" class="text-destructive" />
											{incident.title ?? `Incident #${incident.id}`}
										</span>
										<span class="text-foreground/50 text-sm">
											Updated {formatRelativeTime(incident.updatedAt ?? incident.startedAt)}
										</span>
									</div>
									<Button
										variant="outline"
										class="text-destructive"
										href={`./incidents/${incident.id}`}
									>
										<Icon icon="lucide:arrow-right" /> Go to incident</Button
									>
								</div>
							{/each}
						</div>
					</div>
				{/if}
			</Card>
		{/each}
	</div>

	<DeleteMonitorDialog
		bind:open={confirmOpen}
		monitor={targetMonitor}
		onConfirm={handleDelete}
		loading={deleting}
	/>
</div>
