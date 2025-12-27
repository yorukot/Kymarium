<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { NavItem } from './type';
	import { page } from '$app/state';
	let { items }: { items: NavItem[] } = $props();
	const teamID = $derived.by(() => page.params.teamID);
</script>

<Sidebar.Group>
	<Sidebar.Menu>
		{#each items as item (item.title)}
			<a href={`/${teamID}/${item.url}`}>
				<Sidebar.MenuItem>
					<Sidebar.MenuButton tooltipContent={item.title} class="select-none">
						{#if item.icon}
							<Icon icon={item.icon} />
						{/if}
						<span>{item.title}</span>
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			</a>
		{/each}
	</Sidebar.Menu>
</Sidebar.Group>
