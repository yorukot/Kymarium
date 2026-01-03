<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import Icon from '@iconify/svelte';
	import StatusPages from '$lib/components/status-page/status-pages.svelte';
	import { page } from '$app/state';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();
	const teamID = $state(page.params.teamID ?? '');
</script>

<header class="flex flex-col gap-3">
	<div class="flex justify-between">
		<div>
			<h1 class="text-2xl font-bold">Status pages</h1>
		</div>
		<Button size="default" href="status-pages/new">
			<Icon icon="lucide:plus" />
			New status page
		</Button>
	</div>
	<p class="text-sm text-muted-foreground">
		Create and manage status pages to provide users with real-time updates.
	</p>
</header>

<div class="mt-7">
	<StatusPages statusPages={data.statusPages} {teamID} />
</div>
