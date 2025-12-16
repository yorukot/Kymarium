<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Accordion from '$lib/components/ui/accordion';

	const { errors = {} } = $props<{ errors?: any }>();

	let host = $state('');
	let timeoutSeconds = $state<number | ''>(5);
	let packetSize = $state<number | ''>('');

	const timeoutHelper = $derived.by(() =>
		timeoutSeconds === '' ? 'Default 5s timeout per ping.' : `${timeoutSeconds}s timeout per ping.`
	);

	const packetSizeHelper = $derived.by(
		() =>
			packetSize === ''
				? 'Leave blank to use the default 56-byte payload.'
				: `${packetSize} bytes per packet (1-65000).`
	);
</script>

<Card.Root class="mx-auto w-full">
	<Card.Header>
		<h2 class="text-lg font-bold">Ping monitor settings</h2>
		<p class="text-sm text-muted-foreground">
			Configure how Knocker will ping your host.
		</p>
	</Card.Header>
	<Card.Content>
		<Field.Set>
			<div class="space-y-2">
				<Field.Label for="ping-host">Host</Field.Label>
				<Field.Description>Hostname or IP address to ping.</Field.Description>
				<Input
					id="ping-host"
					name="config.host"
					type="text"
					bind:value={host}
					placeholder="example.com or 192.0.2.1"
					required
				/>
				{#if errors?.config?.host}
					<Field.Description class="text-destructive">
						{errors.config.host[0]}
					</Field.Description>
				{/if}
			</div>

			<Accordion.Root type="single">
				<Accordion.Item value="item-1">
					<Accordion.Trigger class="text-lg">Advance setting</Accordion.Trigger>
					<Accordion.Content class="space-y-6">
						<div class="space-y-2">
							<Field.Label for="ping-timeout">Timeout (seconds)</Field.Label>
							<Field.Description>{timeoutHelper}</Field.Description>
							<Input
								id="ping-timeout"
								name="config.timeoutSeconds"
								type="number"
								min="0"
								max="120"
								step="1"
								bind:value={timeoutSeconds}
								placeholder="5"
							/>
							{#if errors?.config?.timeoutSeconds}
								<Field.Description class="text-destructive">
									{errors.config.timeoutSeconds[0]}
								</Field.Description>
							{/if}
						</div>

						<div class="space-y-2">
							<Field.Label for="ping-packet-size">Packet size (bytes)</Field.Label>
							<Field.Description>{packetSizeHelper}</Field.Description>
							<Input
								id="ping-packet-size"
								name="config.packetSize"
								type="number"
								min="1"
								max="65000"
								step="1"
								bind:value={packetSize}
								placeholder="56"
							/>
							{#if errors?.config?.packetSize}
								<Field.Description class="text-destructive">
									{errors.config.packetSize[0]}
								</Field.Description>
							{/if}
						</div>
					</Accordion.Content>
				</Accordion.Item>
			</Accordion.Root>
		</Field.Set>
	</Card.Content>
</Card.Root>
