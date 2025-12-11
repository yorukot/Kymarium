<script lang="ts">
	import Icon from '@iconify/svelte';

	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardTitle } from '$lib/components/ui/card/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';

	type MonitorStatus = 'operational' | 'degraded' | 'down' | 'paused';
	type MonitorType = 'HTTP' | 'Ping' | 'TCP';
	type MonitorIncidentStatus = 'investigating' | 'identified' | 'monitoring' | 'resolved';
	type MonitorIncidentSeverity = 'critical' | 'major' | 'minor' | 'maintenance';

	type MonitorIncident = {
		id: string;
		status: MonitorIncidentStatus;
		severity: MonitorIncidentSeverity;
		updatedAt: string;
		summary: string;
		link: string;
	};

	type Monitor = {
		id: string;
		name: string;
		target: string;
		type: MonitorType;
		status: MonitorStatus;
		regions: string[];
		frequency: string;
		uptime: string;
		responseTime: string;
		lastChecked: string;
		lastIncident: string;
		incident?: MonitorIncident;
	};

	const statusStyles: Record<
		MonitorStatus,
		{
			label: string;
			class: string;
			dot: string;
		}
	> = {
		operational: {
			label: 'Operational',
			class:
				'border-emerald-200 bg-emerald-500/10 text-emerald-700 dark:border-emerald-500/30 dark:bg-emerald-500/20 dark:text-emerald-100',
			dot: 'bg-emerald-500 dark:bg-emerald-300'
		},
		degraded: {
			label: 'Degraded',
			class:
				'border-amber-200 bg-amber-500/10 text-amber-700 dark:border-amber-400/30 dark:bg-amber-400/15 dark:text-amber-50',
			dot: 'bg-amber-500 dark:bg-amber-300'
		},
		down: {
			label: 'Down',
			class:
				'border-destructive/40 bg-destructive/10 text-destructive dark:border-destructive/30 dark:bg-destructive/20 dark:text-destructive-foreground',
			dot: 'bg-destructive'
		},
		paused: {
			label: 'Paused',
			class:
				'border-muted-foreground/30 bg-muted text-muted-foreground dark:bg-input/40 dark:border-input',
			dot: 'bg-muted-foreground'
		}
	};

	const incidentStatusStyles: Record<
		MonitorIncidentStatus,
		{
			label: string;
			class: string;
		}
	> = {
		investigating: {
			label: 'Investigating',
			class:
				'border-amber-200 bg-amber-500/10 text-amber-700 dark:border-amber-400/30 dark:bg-amber-400/20 dark:text-amber-50'
		},
		identified: {
			label: 'Identified',
			class:
				'border-blue-200 bg-blue-500/10 text-blue-700 dark:border-blue-400/30 dark:bg-blue-400/20 dark:text-blue-50'
		},
		monitoring: {
			label: 'Monitoring',
			class:
				'border-purple-200 bg-purple-500/10 text-purple-700 dark:border-purple-400/30 dark:bg-purple-400/20 dark:text-purple-50'
		},
		resolved: {
			label: 'Resolved',
			class:
				'border-emerald-200 bg-emerald-500/10 text-emerald-700 dark:border-emerald-400/30 dark:bg-emerald-400/20 dark:text-emerald-50'
		}
	};

	const incidentSeverityStyles: Record<
		MonitorIncidentSeverity,
		{
			label: string;
			class: string;
		}
	> = {
		critical: {
			label: 'Critical',
			class:
				'border-destructive/40 bg-destructive/10 text-destructive dark:border-destructive/30 dark:bg-destructive/20 dark:text-destructive-foreground'
		},
		major: {
			label: 'Major',
			class:
				'border-orange-200 bg-orange-500/10 text-orange-700 dark:border-orange-400/30 dark:bg-orange-400/20 dark:text-orange-50'
		},
		minor: {
			label: 'Minor',
			class:
				'border-amber-200 bg-amber-100 text-amber-700 dark:border-amber-300/30 dark:bg-amber-300/20 dark:text-amber-50'
		},
		maintenance: {
			label: 'Maintenance',
			class:
				'border-muted-foreground/30 bg-muted text-muted-foreground dark:bg-input/40 dark:border-input'
		}
	};

	let { monitors = [] }: { monitors?: Monitor[] } = $props();
</script>

<svelte:head>
	<title>Monitors</title>
</svelte:head>

<div class="space-y-4">
	{#each monitors as monitor (monitor.id)}
		<Card class="p-0 border-0 bg-destructive/8 flex flex-col gap-0">
			<Card class="py-2 px-4">
				<CardContent class="p-0">
					<div class="flex gap-3 flex-row items-center justify-between">
						<div class="space-y-2">
							<div class="flex flex-wrap items-center gap-2">
								<span
									class={`inline-block size-2.5 rounded-full ${statusStyles[monitor.status].dot}`}
									aria-hidden="true"
								></span>
								<CardTitle class="text-base">{monitor.name}</CardTitle>
								<Badge variant="secondary" class="rounded-md">{monitor.type}</Badge>
							</div>
							<div class="flex items-center gap-2 text-sm text-muted-foreground">
								<Icon icon="lucide:link" class="size-4" />
								<span class="truncate">{monitor.target}</span>
							</div>
						</div>

						<div class="flex items-center gap-2 justify-end">
							<div class="hidden w-20 md:flex md:flex-col md:items-end md:justify-center text-left md:text-right">
								<p class="text-xs text-muted-foreground">Uptime</p>
								<p class="text-sm font-semibold leading-tight">{monitor.uptime}</p>
							</div>

							<div class="hidden w-20 md:flex md:flex-col md:items-end md:justify-center text-left md:text-right">
								<p class="text-xs text-muted-foreground">Last check</p>
								<p class="text-sm font-semibold leading-tight">{monitor.lastChecked}</p>
							</div>

							<DropdownMenu.Root>
								<DropdownMenu.Trigger>
									{#snippet child({ props })}
										<Button
											variant="ghost"
											size="icon-sm"
											class="-m-1"
											aria-label="More actions"
											{...props}
										>
											<Icon icon="lucide:ellipsis-vertical" class="size-4" />
										</Button>
									{/snippet}
								</DropdownMenu.Trigger>
								<DropdownMenu.Content class="w-44" align="end" sideOffset={6}>
									<DropdownMenu.Item>
										<Icon icon="lucide:route" class="size-4" />
										View details
									</DropdownMenu.Item>
									<DropdownMenu.Item>
										<Icon icon="lucide:radio-tower" class="size-4" />
										Edit monitor
									</DropdownMenu.Item>
									<DropdownMenu.Separator />
									<DropdownMenu.Item variant="destructive">
										<Icon icon="lucide:trash" class="size-4" />
										Delete
									</DropdownMenu.Item>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</div>
					</div>
				</CardContent>
			</Card>
			{#if monitor.incident}
				<div class="px-4 py-2">
					<div class="flex flex-col gap-2 rounded-lg px-3 py-2">
						<div class="flex flex-wrap items-start justify-between gap-3">
							<div class="flex items-start gap-2 text-sm text-foreground">
								<Icon icon="lucide:activity" class="mt-0.5 size-4 text-destructive" />
								<div class="space-y-2">
									<div class="space-y-1">
										<p class="text-sm font-semibold leading-tight">Active incident</p>
										<p class="text-sm text-muted-foreground">{monitor.incident.summary}</p>
									</div>
									<div class="flex flex-wrap items-center gap-2 text-xs">
										<Badge
											variant="outline"
											class={`rounded-md ${incidentSeverityStyles[monitor.incident.severity].class}`}
										>
											{incidentSeverityStyles[monitor.incident.severity].label}
										</Badge>
										<Badge
											variant="outline"
											class={`rounded-md ${incidentStatusStyles[monitor.incident.status].class}`}
										>
											{incidentStatusStyles[monitor.incident.status].label}
										</Badge>
										<span class="text-muted-foreground">Updated {monitor.incident.updatedAt}</span>
									</div>
								</div>
							</div>

							<div class="flex items-center gap-2">
								<Button
									variant="outline"
									size="sm"
									class="border-destructive/40 text-destructive hover:bg-destructive/10"
									href={monitor.incident.link}
								>
									<Icon icon="lucide:arrow-right" class="size-4" />
									Go to incident
								</Button>
							</div>
						</div>
					</div>
				</div>
			{/if}
		</Card>
	{/each}
</div>
