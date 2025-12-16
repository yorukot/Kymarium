<script lang="ts">
	import NewMonitor from '$lib/components/monitor/setting/basic-monitor.svelte';
import { Button } from '$lib/components/ui/button';
	import Icon from '@iconify/svelte';
</script>

<div class="flex flex-col gap-5">
	<div class="flex justify-between">
		<h1 class="text-2xl font-bold">New monitor</h1>

		<Button size="default" variant="ghost">
			<Icon icon="lucide:arrow-left" />
			Back to list
		</Button>
	</div>
	<NewMonitor />
</div>
