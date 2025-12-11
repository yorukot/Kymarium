<script lang="ts">
	import Icon from '@iconify/svelte';

	import Monitors from '$lib/components/monitors/monitors.svelte';
	import { Button } from '$lib/components/ui/button/index.js';

	const monitors = [
		{
			id: 'api-edge',
			name: 'Public API',
			target: 'https://api.knocker.dev/health',
			type: 'HTTP',
			status: 'down',
			regions: ['us-east-1', 'eu-west-2', 'ap-south-1'],
			frequency: 'Every 30s',
			uptime: '99.99%',
			responseTime: '182 ms',
			lastChecked: '2m ago',
			lastIncident: 'Resolved 2d ago',
			incident: {
				id: 'inc-104',
				status: 'investigating',
				severity: 'critical',
				updatedAt: '5m ago',
				summary: 'Public API returning 5xx for us-east-1 after rollout',
				link: '/incidents/inc-104'
			}
		},
		{
			id: 'marketing-site',
			name: 'Marketing site',
			target: 'https://knocker.dev',
			type: 'HTTP',
			status: 'degraded',
			regions: ['us-east-1', 'eu-west-2'],
			frequency: 'Every 1m',
			uptime: '99.80%',
			responseTime: '438 ms',
			lastChecked: 'Just now',
			lastIncident: 'Slow TTFB detected',
			incident: {
				id: 'inc-099',
				status: 'monitoring',
				severity: 'major',
				updatedAt: '12m ago',
				summary: 'Elevated TTFB while CDN purge finishes',
				link: '/incidents/inc-099'
			}
		},
		{
			id: 'bg-worker',
			name: 'Worker heartbeat',
			target: 'redis://internal:6379/ping',
			type: 'TCP',
			status: 'operational',
			regions: ['us-east-1'],
			frequency: 'Every 15s',
			uptime: '100%',
			responseTime: '41 ms',
			lastChecked: '45s ago',
			lastIncident: 'None in the last 7d'
		},
		{
			id: 'edge-latency',
			name: 'Edge latency',
			target: 'ping knocker-edge',
			type: 'Ping',
			status: 'paused',
			regions: ['eu-west-2', 'ap-south-1'],
			frequency: 'Paused',
			uptime: '—',
			responseTime: '—',
			lastChecked: '1h ago',
			lastIncident: 'Maintenance window'
		}
	];
</script>

<section class="space-y-6">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="space-y-1">
			<h1 class="text-2xl font-semibold leading-tight text-foreground">Monitors</h1>
		</div>
		<Button class="ms-auto" size="lg">
			<Icon icon="ri:add-line" class="size-4" />
			Create monitor
		</Button>
	</div>

	<Monitors {monitors} />
</section>
