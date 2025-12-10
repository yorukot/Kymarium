<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { useSidebar } from '$lib/components/ui/sidebar/index.js';
	import { createAvatar } from '@dicebear/core';
	import { thumbs } from '@dicebear/collection';
	let { user }: { user: { name: string; email?: string; avatar?: string } } = $props();
	const sidebar = useSidebar();

	const displayName = $derived(user.name || 'User');
	const displayEmail = $derived(user.email || 'Signed in');
	const fallbackInitials = $derived(
		displayName
			.split(' ')
			.map((part) => part[0])
			.join('')
			.slice(0, 2)
			.toUpperCase()
	);

	const generatedAvatar = $derived(
		`data:image/svg+xml;utf8,${encodeURIComponent(
			createAvatar(thumbs, {
				seed: displayName
			}).toString()
		)}`
	);

	const avatarSrc = $derived(user.avatar || generatedAvatar);
</script>
<Sidebar.Menu>
  <Sidebar.MenuItem>
    <DropdownMenu.Root>
      <DropdownMenu.Trigger>
        {#snippet child({ props })}
          <Sidebar.MenuButton
            size="lg"
            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            {...props}
          >
            <Avatar.Root class="size-8 rounded-lg">
              <Avatar.Image src={avatarSrc} alt={displayName} />
              <Avatar.Fallback class="rounded-lg">{fallbackInitials}</Avatar.Fallback>
            </Avatar.Root>
            <div class="grid flex-1 text-start text-sm leading-tight">
              <span class="truncate font-medium">{displayName}</span>
              <span class="truncate text-xs">{displayEmail}</span>
            </div>
            <Icon icon="lucide:chevron-down" class="ms-auto size-4" />
          </Sidebar.MenuButton>
        {/snippet}
      </DropdownMenu.Trigger>
      <DropdownMenu.Content
        class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
        side={sidebar.isMobile ? "bottom" : "right"}
        align="end"
        sideOffset={4}
      >
        <DropdownMenu.Label class="p-0 font-normal">
          <div class="flex items-center gap-2 px-1 py-1.5 text-start text-sm">
            <Avatar.Root class="size-8 rounded-lg">
              <Avatar.Image src={avatarSrc} alt={displayName} />
              <Avatar.Fallback class="rounded-lg">{fallbackInitials}</Avatar.Fallback>
            </Avatar.Root>
            <div class="grid flex-1 text-start text-sm leading-tight">
              <span class="truncate font-medium">{displayName}</span>
              <span class="truncate text-xs">{displayEmail}</span>
            </div>
          </div>
        </DropdownMenu.Label>
        <DropdownMenu.Separator />
        <DropdownMenu.Group>
          <DropdownMenu.Item>
            <Icon icon="lucide:sparkles" />
            Upgrade to Pro
          </DropdownMenu.Item>
        </DropdownMenu.Group>
        <DropdownMenu.Separator />
        <DropdownMenu.Group>
          <DropdownMenu.Item>
            <Icon icon="lucide:badge-check" />
            Account
          </DropdownMenu.Item>
          <DropdownMenu.Item>
            <Icon icon="lucide:credit-card" />
            Billing
          </DropdownMenu.Item>
          <DropdownMenu.Item>
            <Icon icon="lucide:bell" />
            Notifications
          </DropdownMenu.Item>
        </DropdownMenu.Group>
        <DropdownMenu.Separator />
        <DropdownMenu.Item>
          <Icon icon="lucide:log-out" />
          Log out
        </DropdownMenu.Item>
      </DropdownMenu.Content>
    </DropdownMenu.Root>
  </Sidebar.MenuItem>
</Sidebar.Menu>
