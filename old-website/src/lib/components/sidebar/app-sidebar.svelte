<script lang="ts">
	import NavMain from './nav-main.svelte';
	import NavUser from './nav-user.svelte';
	import TeamSwitcher from './team-switcher.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { Team, User } from '$lib/types';
	import type { NavItem } from './type';

	const NavItems: NavItem[] = [
		{ title: 'Monitors', url: 'monitors', icon: 'lucide:monitor' },
		{ title: 'Status pages', url: 'status-pages', icon: 'lucide:layout-template' },
		{ title: 'Incidents', url: 'incidents', icon: 'lucide:alert-triangle' },
		{ title: 'Notifications', url: 'notifications', icon: 'lucide:bell' },
		{ title: 'Members', url: 'members', icon: 'lucide:users' },
		{ title: 'Settings', url: 'setting', icon: 'lucide:settings' }
	];

	let { user, teams }: { user: User; teams: Team[]; } =  $props();
</script>

<Sidebar.Root>
	<Sidebar.Header>
		<TeamSwitcher teams={teams}/>
	</Sidebar.Header>
	<Sidebar.Content>
		<NavMain items={NavItems} />
	</Sidebar.Content>
	<Sidebar.Footer>
		<NavUser user={user} />
	</Sidebar.Footer>
	<Sidebar.Rail />
</Sidebar.Root>
