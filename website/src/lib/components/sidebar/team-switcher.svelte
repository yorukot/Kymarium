<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { useSidebar } from '$lib/components/ui/sidebar/index.js';
	import Icon from '@iconify/svelte';
	import { createAvatar } from '@dicebear/core';
	import { shapes } from '@dicebear/collection';

	type Team = { id?: string; name: string; logo?: string; plan?: string };

	let { teams, activeTeamId }: { teams: Team[]; activeTeamId?: string } = $props();
	const sidebar = useSidebar();

	let activeTeam = $derived(
		teams.find((team) => team.id === activeTeamId) ?? teams[0]
	);

	const avatarMap = $derived(
		teams.reduce((acc, team) => {
			const key = team.id ?? team.name;
			const avatarSvg = createAvatar(shapes, {
				seed: key
			}).toString();

			acc[key] = `data:image/svg+xml;utf8,${encodeURIComponent(avatarSvg)}`;
			return acc;
		}, {} as Record<string, string>)
	);

	const activeTeamAvatar = $derived(
		avatarMap[activeTeam.id ?? activeTeam.name]
	);
</script>

<Sidebar.Menu>
	<Sidebar.MenuItem>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Sidebar.MenuButton
						{...props}
						size="lg"
						class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
					>
						<div
							class="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg overflow-hidden"
						>
							{#if activeTeam.logo}
								<Icon icon={activeTeam.logo} class="size-4" />
							{:else}
								<img
									src={activeTeamAvatar}
									alt={activeTeam.name}
									class="size-full object-cover"
									loading="lazy"
								/>
							{/if}
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
				{#each teams as team, index (team.name)}
					<DropdownMenu.Item onSelect={() => (activeTeam = team)} class="gap-2 p-2">
						<div class="flex size-6 items-center justify-center">
							{#if team.logo}
								<Icon icon={team.logo} class="size-3.5 shrink-0" />
							{:else}
								<img
									src={avatarMap[team.id ?? team.name]}
									alt={team.name}
									class="rounded-sm"
									loading="lazy"
								/>
							{/if}
						</div>
						<span class="truncate">{team.name}</span>
						<DropdownMenu.Shortcut>⌘{index + 1}</DropdownMenu.Shortcut>
					</DropdownMenu.Item>
				{/each}
				<DropdownMenu.Separator />
				<DropdownMenu.Item class="gap-2 p-2">
					<div class="flex size-6 items-center justify-center rounded-md border bg-transparent">
						<Icon icon="lucide:plus" />
					</div>
					<div class="text-muted-foreground font-medium">Add team</div>
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>
