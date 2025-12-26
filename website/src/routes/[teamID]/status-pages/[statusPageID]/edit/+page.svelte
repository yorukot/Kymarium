<script lang="ts">
	import DeleteStatusPage from '$lib/components/status-page/edit/delete-status-page.svelte';
	import { Button } from '$lib/components/ui/button';
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import {
		statusPageUpsertRequestSchema,
		type StatusPageUpsertValues
	} from '$lib/components/status-page/edit/schema';
	import { reporter } from '@felte/reporter-svelte';
	import { page } from '$app/state';
	import GeneralStatusPage from '$lib/components/status-page/edit/general-status-page.svelte';
	import { updateStatusPage } from '$lib/api/status-page';
	import { Spinner } from '$lib/components/ui/spinner';
	import type { StatusPage, StatusPageElement } from '$lib/types/status-page.js';
	import DndStatusPage from '$lib/components/status-page/edit/dnd-status-page.svelte';

	const schema = statusPageUpsertRequestSchema;

	/** @type {import('$types').PageProps} */
	let { data } = $props();

	// svelte-ignore state_referenced_locally
	let statusPage: StatusPage = data.statusPage.statusPage;
	// svelte-ignore state_referenced_locally
	const elements: StatusPageElement[] = data.statusPage.elements ?? [];

	const { form, setFields, isSubmitting } = createForm<StatusPageUpsertValues>({
		extend: [validator({ schema }), reporter()],
		initialValues: {
			name: statusPage.title,
			slug: statusPage.slug,
			elements: elements ?? []
		},
		onSubmit: async (values) => {
			try {
				const teamID = page.params.teamID;
				const statusPageID = page.params.statusPageID;
				if (!teamID || !statusPageID) {
					return { FORM_ERROR: 'Missing team or status page id' };
				}

				const payload = {
					title: values.name,
					slug: values.slug,
					icon: statusPage.icon ?? null,
					elements: (values.elements ?? []).map((element) => ({
						id: element.id,
						name: element.name,
						type: element.type,
						sort_order: element.sortOrder,
						monitor: element.monitor,
						monitor_id: element.monitor ? (element.monitorId ?? undefined) : undefined,
						monitors: element.monitor
							? []
							: (element.monitors ?? []).map((monitor) => ({
									id: monitor.id,
									monitor_id: monitor.monitorId,
									group_id: monitor.groupId ?? undefined,
									name: monitor.name,
									type: monitor.type,
									sort_order: monitor.sortOrder
								}))
					}))
				};

				console.log('status-page.save.values', values);
				console.log('status-page.save.payload', payload);

				await updateStatusPage(teamID, statusPageID, payload);
			} catch (error) {
				return {
					FORM_ERROR:
						error instanceof Error
							? error.message
							: 'Unable to save changes right now. Please try again.'
				};
			}
		}
	});
</script>

<form use:form>
	<div class="flex flex-col gap-6">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-2xl font-bold">Edit status page</h1>
				<p class="text-sm text-muted-foreground">Organize what your public status page shows.</p>
			</div>
		</div>
		<GeneralStatusPage />
		<DndStatusPage
			availableMonitors={data.monitors}
			statusPageId={statusPage.id}
			{elements}
			{setFields}
		/>
		<div class="flex items-center gap-2 justify-end">
			<DeleteStatusPage />
			<Button type="submit" disabled={$isSubmitting}>
				{#if $isSubmitting}
					<Spinner />
				{/if}
				Save changes
			</Button>
		</div>
	</div>
</form>
