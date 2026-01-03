<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import Icon from '@iconify/svelte';
	import CreateIncident from '$lib/components/incident/create-incident.svelte';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();
</script>

<div class="flex flex-col gap-5">
	<div class="flex justify-between items-start gap-3 flex-wrap">
		<h1 class="text-2xl font-bold">New incident</h1>
		<Button size="default" variant="ghost" href="../incidents">
			<Icon icon="lucide:arrow-left" />
			Back to incidents
		</Button>
	</div>

	<CreateIncident monitors={data.monitors} />
</div>
