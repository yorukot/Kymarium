<script lang="ts">
	import NavMain from './nav-main.svelte';
	import NavUser from './nav-user.svelte';
	import TeamSwitcher from './team-switcher.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';

	type SidebarTeam = {
		id: string;
		name: string;
		logo?: string;
		plan?: string;
	};

	type SidebarUser = {
		name: string;
		email?: string;
		avatar?: string;
	};

	type NavItem = {
		title: string;
		url: string;
		icon?: string;
	};

	const fallbackUser: SidebarUser = {
		name: 'User',
		email: 'Signed in'
	};

	const fallbackTeams: SidebarTeam[] = [
		{ id: 'placeholder', name: 'Team', logo: 'lucide:layout-panel-left', plan: 'Member' }
	];

	const fallbackNavItems: NavItem[] = [
		{ title: 'Monitors', url: '/monitors', icon: 'lucide:monitor' },
		{ title: 'Incidents', url: '/incidents', icon: 'lucide:alert-triangle' },
		{ title: 'Notifications', url: '/notifications', icon: 'lucide:bell' }
	];

let {
	user = fallbackUser,
	teams = fallbackTeams,
	navItems = fallbackNavItems,
		currentTeamId,
		ref = $bindable(null),
		collapsible = 'icon',
		...restProps
	}: {
		user?: SidebarUser;
		teams?: SidebarTeam[];
		navItems?: NavItem[];
		currentTeamId?: string;
	} & ComponentProps<typeof Sidebar.Root> = $props();

	const resolvedTeams = $derived(
		teams && teams.length > 0 ? teams : fallbackTeams
	);
	const baseNavItems = $derived(
		navItems && navItems.length > 0 ? navItems : fallbackNavItems
	);
	const resolvedNavItems = $derived(
		baseNavItems.map((item) => {
			if (
				!currentTeamId ||
				item.url.startsWith('http') ||
				item.url.startsWith(`/${currentTeamId}`)
			) {
				return item;
			}

			const normalizedPath = item.url.startsWith('/') ? item.url : `/${item.url}`;
			return { ...item, url: `/${currentTeamId}${normalizedPath}` };
		})
	);
</script>

<Sidebar.Root {collapsible} {...restProps}>
	<Sidebar.Header>
		<TeamSwitcher teams={resolvedTeams} activeTeamId={currentTeamId} />
	</Sidebar.Header>
	<Sidebar.Content>
		<NavMain items={resolvedNavItems} />
	</Sidebar.Content>
	<Sidebar.Footer>
		<NavUser user={user ?? fallbackUser} />
	</Sidebar.Footer>
	<Sidebar.Rail />
</Sidebar.Root>
