<script lang="ts">
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { useSidebar } from '$lib/components/ui/sidebar/index.js';
	import { logout } from '$lib/api/auth.js';
	import InviteDialog from './invite-dialog.svelte';
	import { createAvatar } from '@dicebear/core';
	import { thumbs } from '@dicebear/collection';
	import { setMode, userPrefersMode } from 'mode-watcher';
	import type { User } from '../../types';

	const sidebar = useSidebar();

	let { user }: { user: User } = $props();
	let inviteDialogOpen = $state(false);

	const teamID = $derived.by(() => page.params.teamID);
	const themePreference = $derived(userPrefersMode.current ?? 'system');

	const generatedAvatar = $derived(
		`data:image/svg+xml;utf8,${encodeURIComponent(
			createAvatar(thumbs, {
				seed: user.id
			}).toString()
		)}`
	);

	const avatarSrc = $derived(user.avatar || generatedAvatar);

	const handleLogout = async () => {
		try {
			await logout();
		} catch {
			// Ignore logout errors; we'll still clear client state by redirecting.
		}
		await goto('/auth/login');
	};

	const handleAccount = async () => {
		if (!teamID) return;
		await goto(`/${teamID}/account`);
	};

	const handleThemeChange = (value: 'light' | 'dark' | 'system') => {
		setMode(value);
	};
</script>

<Sidebar.Menu>
	<Sidebar.MenuItem>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Sidebar.MenuButton
						size="lg"
						class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground select-none"
						{...props}
					>
						<Avatar.Root class="size-8 rounded-lg">
							<Avatar.Image src={avatarSrc} alt={user.displayName} />
							<Avatar.Fallback class="rounded-lg">{user.displayName}</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-start text-sm leading-tight">
							<span class="truncate font-medium">{user.displayName}</span>
						</div>
						<Icon icon="lucide:chevron-down" class="ms-auto size-4" />
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				align="end"
				sideOffset={4}
			>
				<DropdownMenu.Label class="p-0 font-normal">
					<div class="flex items-center gap-2 px-1 py-1.5 text-start text-sm">
						<Avatar.Root class="size-8 rounded-lg">
							<Avatar.Image src={avatarSrc} alt={user.displayName} />
							<Avatar.Fallback class="rounded-lg">{user.displayName}</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-start text-sm leading-tight">
							<span class="truncate font-medium">{user.displayName}</span>
						</div>
					</div>
				</DropdownMenu.Label>
				<DropdownMenu.Separator />
				<DropdownMenu.Group>
					<DropdownMenu.Item onclick={handleAccount}>
						<Icon icon="lucide:badge-check" />
						Account
					</DropdownMenu.Item>
					<DropdownMenu.Item onclick={() => (inviteDialogOpen = true)}>
						<Icon icon="lucide:mail" />
						Invitations
					</DropdownMenu.Item>
				</DropdownMenu.Group>
				<DropdownMenu.Group>
					<DropdownMenu.Sub>
						<DropdownMenu.SubTrigger>
							<Icon icon="lucide:palette" />
							Theme
						</DropdownMenu.SubTrigger>
						<DropdownMenu.SubContent>
							<DropdownMenu.Item onclick={() => handleThemeChange('light')}>
								<Icon
									icon="lucide:check"
									class={themePreference === 'light' ? '' : 'opacity-0'}
								/>
								Light
							</DropdownMenu.Item>
							<DropdownMenu.Item onclick={() => handleThemeChange('dark')}>
								<Icon
									icon="lucide:check"
									class={themePreference === 'dark' ? '' : 'opacity-0'}
								/>
								Dark
							</DropdownMenu.Item>
							<DropdownMenu.Item onclick={() => handleThemeChange('system')}>
								<Icon
									icon="lucide:check"
									class={themePreference === 'system' ? '' : 'opacity-0'}
								/>
								System (default)
							</DropdownMenu.Item>
						</DropdownMenu.SubContent>
					</DropdownMenu.Sub>
				</DropdownMenu.Group>
				<DropdownMenu.Separator />
				<DropdownMenu.Item onclick={handleLogout} variant="destructive">
					<Icon icon="lucide:log-out" />
					Log out
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>

<InviteDialog bind:open={inviteDialogOpen} />
