<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import Textarea from '$lib/components/ui/textarea/textarea.svelte';
	import * as Select from '$lib/components/ui/select';
	import { Switch } from '$lib/components/ui/switch';
	import MultiSelect from '$lib/components/ui/multi-select';
	import * as Accordion from '$lib/components/ui/accordion';
	import { ValidationMessage } from '@felte/reporter-svelte';

	import type { BodyEncoding, HTTPMethod } from '../$lib/type/monitor-config';
	import { bodyEncodingOptions, acceptedStatusOptions, httpMethods } from './setting';
	import { condenseStatusTokens } from './utils';

	
	const { initialConfig } = $props<{
		initialConfig?: {
			url: string;
			method: HTTPMethod;
			maxRedirects: number;
			requestTimeoutSeconds: number;
			headers?: string;
			bodyEncoding?: BodyEncoding | '';
			body?: string;
			acceptedStatusCodes?: string[];
			upsideDownMode: boolean;
			ignoreTlsError: boolean;
			certificateExpiryNotification: boolean;
		};
	}>();

	let url = $derived(initialConfig?.url ?? '');
	let method = $derived<HTTPMethod>(initialConfig?.method ?? 'GET');
	let requestTimeoutSeconds = $derived<number | ''>(initialConfig?.requestTimeoutSeconds ?? 10);
	let maxRedirects = $derived<number | ''>(initialConfig?.maxRedirects ?? 5);
	let headers = $derived(initialConfig?.headers ?? '');
	let bodyEncoding = $derived<BodyEncoding | ''>(initialConfig?.bodyEncoding ?? '');
	let body = $derived(initialConfig?.body ?? '');
	let acceptedStatusCodes = $derived<string[]>(initialConfig?.acceptedStatusCodes ?? ['2xx']);
	let upsideDownMode = $derived(initialConfig?.upsideDownMode ?? false);
	let ignoreTlsError = $derived(initialConfig?.ignoreTlsError ?? false);
	let certificateExpiryNotification = $derived(
		initialConfig?.certificateExpiryNotification ?? true
	);

	const timeoutHelper = $derived.by(() =>
		requestTimeoutSeconds === ''
			? 'No timeout (not recommended).'
			: `${requestTimeoutSeconds}s timeout per request.`
	);

	const redirectsHelper = $derived.by(() =>
		maxRedirects === ''
			? 'Will stop after the first response.'
			: `Follow up to ${maxRedirects} redirects.`
	);

	const acceptedStatusHelper =
		'Choose specific codes or broad ranges (2xx/4xx/5xx). Leave empty to accept any 2xx.';

	$effect(() => {
		const condensed = condenseStatusTokens(acceptedStatusCodes ?? []);
		if (condensed.join('|') !== (acceptedStatusCodes ?? []).join('|')) {
			acceptedStatusCodes = condensed;
		}
	});
</script>

<Card.Root class="mx-auto w-full">
	<Card.Header>
		<h2 class="text-lg font-bold">HTTP monitor settings</h2>
		<p class="text-sm text-muted-foreground">
			Configure how Knocker will probe your HTTP endpoint.
		</p>
	</Card.Header>
	<Card.Content>
		<Field.Set>
			<div class="space-y-2">
				<Field.Label for="http-url">URL</Field.Label>
				<Field.Description>Full URL (http or https) to check.</Field.Description>
				<Input
					id="http-url"
					name="config.url"
					type="url"
					bind:value={url}
					placeholder="https://example.com/health"
					required
				/>
				<ValidationMessage for="config.url" let:messages>
					{#if messages?.length}
						<Field.Description class="text-destructive">
							{messages[0]}
						</Field.Description>
					{/if}
				</ValidationMessage>
			</div>
			<Accordion.Root type="single">
				<Accordion.Item value="item-1">
					<Accordion.Trigger class="text-lg">Advance setting</Accordion.Trigger>
					<Accordion.Content class="space-y-6">
						<div class="space-y-2">
							<Field.Label>Method</Field.Label>
							<Select.Root type="single" bind:value={method}>
								<Select.Trigger class="w-full justify-between">
									<span data-slot="select-value" class="text-sm font-medium">
										{method}
									</span>
								</Select.Trigger>
								<Select.Content>
									<Select.Group>
										{#each httpMethods as opt (opt)}
											<Select.Item value={opt}>{opt}</Select.Item>
										{/each}
									</Select.Group>
								</Select.Content>
							</Select.Root>
							<input type="hidden" name="config.method" value={method} />
						</div>

						<div class="space-y-2">
							<Field.Label for="http-timeout">Request timeout (seconds)</Field.Label>
							<Field.Description>{timeoutHelper}</Field.Description>
							<Input
								id="http-timeout"
								name="config.requestTimeoutSeconds"
								type="number"
								min="0"
								step="1"
								bind:value={requestTimeoutSeconds}
								placeholder="10"
							/>
							<ValidationMessage for="config.requestTimeoutSeconds" let:messages>
								{#if messages?.length}
									<Field.Description class="text-destructive">
										{messages[0]}
									</Field.Description>
								{/if}
							</ValidationMessage>
						</div>

						<div class="space-y-2">
							<Field.Label for="http-redirects">Max redirects</Field.Label>
							<Field.Description>{redirectsHelper}</Field.Description>
							<Input
								id="http-redirects"
								name="config.maxRedirects"
								type="number"
								min="0"
								step="1"
								bind:value={maxRedirects}
								placeholder="5"
							/>
							<ValidationMessage for="config.maxRedirects" let:messages>
								{#if messages?.length}
									<Field.Description class="text-destructive">
										{messages[0]}
									</Field.Description>
								{/if}
							</ValidationMessage>
						</div>

						<div class="space-y-2">
							<Field.Label for="http-status-codes">Accepted status codes</Field.Label>
							<Field.Description>{acceptedStatusHelper}</Field.Description>
							<MultiSelect
								name="config.acceptedStatusCodes"
								bind:value={acceptedStatusCodes}
								options={acceptedStatusOptions}
								placeholder="Default: any 2xx"
								emptyMessage="No matching codes"
								maxBadges={4}
							/>
							<ValidationMessage for="config.acceptedStatusCodes" let:messages>
								{#if messages?.length}
									<Field.Description class="text-destructive">
										{messages[0]}
									</Field.Description>
								{/if}
							</ValidationMessage>
						</div>

						<div class="space-y-2">
							<Field.Label for="http-headers">Headers (optional)</Field.Label>
							<Field.Description>
								One header per line using <code>Key: Value</code>.
							</Field.Description>
							<Textarea
								id="http-headers"
								name="config.headers"
								bind:value={headers}
								placeholder="Authorization: Bearer token"
							/>
						</div>

						<div class="space-y-2">
							<h2 class="font-bold text-lg">Body setting</h2>
							<Field.Label>Body encoding</Field.Label>
							<Select.Root type="single" bind:value={bodyEncoding}>
								<Select.Trigger class="w-full justify-between">
									<span data-slot="select-value" class="text-sm font-medium">
										{bodyEncoding ? bodyEncoding.toUpperCase() : 'None'}
									</span>
								</Select.Trigger>
								<Select.Content>
									<Select.Group>
										{#each bodyEncodingOptions as opt (opt.value)}
											<Select.Item value={opt.value}>{opt.label}</Select.Item>
										{/each}
									</Select.Group>
								</Select.Content>
							</Select.Root>
							<input type="hidden" name="config.bodyEncoding" value={bodyEncoding} />
						</div>

						<div class="space-y-2">
							<Field.Label for="http-body">Request body (optional)</Field.Label>
							<Field.Description>Shown when an encoding is selected.</Field.Description>
							<Textarea
								id="http-body"
								name="config.body"
								disabled={!bodyEncoding}
								bind:value={body}
								placeholder="Example: status ok"
							/>
						</div>

						<Field.Field class="gap-1">
							<Field.Label class="font-medium">Upside-down mode</Field.Label>
							<Field.Content class="flex flex-row items-center gap-2">
								<Switch bind:checked={upsideDownMode} />
								<input
									type="checkbox"
									name="config.upsideDownMode"
									class="hidden"
									bind:checked={upsideDownMode}
								/>
								<Field.Description>
									Mark monitor as failed when status <em>is</em> accepted; useful for maintenance pages.
								</Field.Description>
							</Field.Content>
						</Field.Field>

						<Field.Field class="gap-1">
							<Field.Label class="font-medium">Ignore TLS errors</Field.Label>
							<Field.Content class="flex flex-row items-center gap-2">
								<Switch bind:checked={ignoreTlsError} />
								<input
									type="checkbox"
									name="config.ignoreTlsError"
									class="hidden"
									bind:checked={ignoreTlsError}
								/>
								<Field.Description>
									Skip certificate validation (insecure; only for testing).
								</Field.Description>
							</Field.Content>
						</Field.Field>

						<Field.Field class="gap-1">
							<Field.Label class="font-medium">Certificate expiry alerts</Field.Label>
							<Field.Content class="flex flex-row items-center gap-2">
								<Switch bind:checked={certificateExpiryNotification} />
								<input
									type="checkbox"
									name="config.certificateExpiryNotification"
									class="hidden"
									bind:checked={certificateExpiryNotification}
								/>
								<Field.Description>
									Notify when the TLS certificate is close to expiring.
								</Field.Description>
							</Field.Content>
						</Field.Field>
					</Accordion.Content>
				</Accordion.Item>
			</Accordion.Root>
		</Field.Set>
	</Card.Content>
</Card.Root>
