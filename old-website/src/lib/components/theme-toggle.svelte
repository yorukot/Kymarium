<script lang="ts">
	import Icon from '@iconify/svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { mode, setMode, toggleMode, userPrefersMode } from 'mode-watcher';

	type ModeValue = 'light' | 'dark' | 'system';

	const activeMode = $derived(mode.current ?? 'light');
	const currentPreference = $derived(userPrefersMode.current ?? 'system');
	const nextMode = $derived(activeMode === 'dark' ? 'light' : 'dark');

	function handleSet(modeValue: ModeValue) {
		setMode(modeValue);
	}

	const label = $derived.by(() => {
		if (currentPreference === 'system') {
			return `System (${activeMode} now)`;
		}

		return currentPreference === 'light' ? 'Light' : 'Dark';
	});
</script>

<div class="flex items-center gap-2">
	<Button
		variant="ghost"
		size="icon"
		class="relative"
		aria-label={`Switch to ${nextMode} mode`}
		aria-pressed={activeMode === 'dark'}
		onclick={toggleMode}
	>
		<Icon
			icon="lucide:sun"
			class="size-5 rotate-0 scale-100 transition-transform duration-300 dark:-rotate-90 dark:scale-0"
		/>
		<Icon
			icon="lucide:moon"
			class="size-5 absolute rotate-90 scale-0 transition-transform duration-300 dark:rotate-0 dark:scale-100"
		/>
	<span class="sr-only">Toggle theme</span>
</Button>

	<Button
		variant="outline"
		size="sm"
		class="flex items-center gap-2"
		aria-pressed={currentPreference === 'system'}
		onclick={() => handleSet('system')}
	>
		<Icon icon="lucide:monitor" class="size-4" />
	<span>System</span>
</Button>

	<span class="text-sm text-muted-foreground">{label}</span>
</div>
