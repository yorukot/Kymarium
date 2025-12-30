<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { useSidebar } from '$lib/components/ui/sidebar/index.js';
	import Icon from '@iconify/svelte';
	import { createAvatar } from '@dicebear/core';
	import { shapes } from '@dicebear/collection';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import type { Team } from '../../types';
	import AddTeamDialog from './add-team-dialog.svelte';

	const teamID = $derived.by(() => page.params.teamID);
	let { teams }: { teams: Team[] } = $props();

	const sidebar = useSidebar();

	const activeTeam = $derived.by(() => {
		return teams.find((team) => team.id === teamID) ?? teams[0];
	});

	const avatarMap = $derived.by(() => {
		return Object.fromEntries(
			teams.map((team) => {
				const key = team.id ?? team.name;

				const svg = createAvatar(shapes, { seed: key }).toString();
				const avatar = `data:image/svg+xml;utf8,${encodeURIComponent(svg)}`;

				return [key, avatar];
			})
		) as Record<string, string>;
	});

	const activeTeamAvatar = $derived.by(() => {
		const team = activeTeam;
		const key = team.id ?? team.name;
		return avatarMap[key];
	});

	let addTeamOpen = $state(false);
</script>

<Sidebar.Menu>
	<Sidebar.MenuItem>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Sidebar.MenuButton
						{...props}
						size="lg"
						class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground select-none"
					>
						<div
							class="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg"
						>
							<img
								src={activeTeamAvatar}
								alt={activeTeam.name}
								class="size-full object-cover rounded-lg"
								loading="lazy"
							/>
						</div>
						<div class="grid flex-1 text-start text-sm leading-tight">
							<span class="truncate font-medium">
								{activeTeam.name}
							</span>
						</div>
						<Icon icon="lucide:chevrons-up-down" class="ms-auto" />
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
				align="start"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				sideOffset={4}
			>
				<DropdownMenu.Label class="text-muted-foreground text-xs">Teams</DropdownMenu.Label>
				{#each teams as team (team.name)}
					<DropdownMenu.Item onSelect={() => goto(`/${team.id}`)} class="gap-2 p-2">
						<div class="flex size-6 items-center justify-center">
							<img src={avatarMap[team.id]} alt={team.name} class="rounded-sm" loading="lazy" />
						</div>
						<span class="truncate">{team.name}</span>
					</DropdownMenu.Item>
				{/each}
				<DropdownMenu.Separator />
				<DropdownMenu.Item class="gap-2 p-2" onSelect={() => (addTeamOpen = true)}>
					<div class="flex size-6 items-center justify-center rounded-md border bg-transparent">
						<Icon icon="lucide:plus" />
					</div>
					<div class="text-muted-foreground font-medium">Add team</div>
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>

<AddTeamDialog bind:open={addTeamOpen} />
