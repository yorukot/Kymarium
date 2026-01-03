<script lang="ts">
	import NewMonitor from '$lib/components/monitor/setting/basic-monitor.svelte';
	import { Button } from '$lib/components/ui/button';
	import Icon from '@iconify/svelte';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();
</script>

<div class="flex flex-col gap-5">
	<div class="flex justify-between">
		<h1 class="text-2xl font-bold">
			Edit monitor{data.monitor ? `: ${data.monitor.name}` : ''}
		</h1>

		<Button size="default" variant="ghost" href="../../monitors">
			<Icon icon="lucide:arrow-left" />
			Back to list
		</Button>
	</div>
	<NewMonitor
		notifications={data.notifications}
		regions={data.regions}
		monitor={data.monitor}
	/>
</div>
