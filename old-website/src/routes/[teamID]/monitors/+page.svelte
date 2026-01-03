<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import Icon from '@iconify/svelte';
	import Monitors from '$lib/components/monitor/monitors.svelte';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();
</script>

<header class="flex flex-col gap-3">
	<div class="flex justify-between">
		<h1 class="text-2xl font-bold">Monitors</h1>

		<Button size="default" href="monitors/new">
			<Icon icon="lucide:plus" />
			New monitor
		</Button>
	</div>
	<p class="text-sm text-muted-foreground">
		Create and manage status pages to provide users with real-time updates.
	</p>
</header>

<div class="mt-7">
	<Monitors monitors={data.monitors} />
</div>
