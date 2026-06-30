<script lang="ts">
	import ChevronsUpDownIcon from "@lucide/svelte/icons/chevrons-up-down";

	import { buttonVariants } from "$lib/components/ui/button/index.js";
	import * as Command from "$lib/components/ui/command/index.js";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { cn } from "$lib/utils.js";

	import {
		countryOptions,
		getCountryName,
		getCountryOption,
		type CountryLocale,
		type CountryOption
	} from "./country-options.js";

	type Props = {
		value?: string;
		locale?: CountryLocale;
		id?: string;
		disabled?: boolean;
		placeholder?: string;
		searchPlaceholder?: string;
		emptyMessage?: string;
		class?: string;
		"aria-describedby"?: string;
		"aria-invalid"?: boolean | "false" | "true" | "grammar" | "spelling";
	};

	let {
		value = $bindable(""),
		locale = "en",
		id,
		disabled = false,
		placeholder = "Select country",
		searchPlaceholder = "Search countries...",
		emptyMessage = "No countries found.",
		class: className,
		"aria-describedby": ariaDescribedby,
		"aria-invalid": ariaInvalid
	}: Props = $props();

	let open = $state(false);

	const selectedCountry = $derived(getCountryOption(value));
	const selectedCountryName = $derived(
		selectedCountry ? getCountryName(selectedCountry, locale) : ""
	);

	$effect(() => {
		if (disabled && open) {
			open = false;
		}
	});

	function selectCountry(country: CountryOption) {
		value = country.code;
		open = false;
	}
</script>

<Popover.Root bind:open>
	<Popover.Trigger
		{id}
		type="button"
		role="combobox"
		aria-expanded={open}
		aria-describedby={ariaDescribedby}
		aria-invalid={ariaInvalid}
		{disabled}
		class={cn(
			buttonVariants({ variant: "outline" }),
			"h-9 w-full justify-between rounded-lg px-2.5 text-left text-sm focus-visible:ring-primary/20",
			!selectedCountry && "text-muted-foreground",
			className
		)}
	>
		<span class="flex min-w-0 items-center gap-2">
			{#if selectedCountry}
				<img
					src={selectedCountry.flagSrc}
					alt=""
					width="24"
					height="24"
					loading="lazy"
					class="size-6 shrink-0 rounded-sm object-cover ring-1 ring-border/40"
				/>
				<span class="min-w-0 truncate">{selectedCountryName}</span>
				<span class="shrink-0 text-xs font-semibold text-muted-foreground">
					{selectedCountry.code}
				</span>
			{:else}
				<span class="truncate">{placeholder}</span>
			{/if}
		</span>
		<ChevronsUpDownIcon class="size-4 shrink-0 text-muted-foreground" />
	</Popover.Trigger>
	<Popover.Content class="w-[min(calc(100vw-2rem),24rem)] p-0" align="start">
		<Command.Root>
			<Command.Input placeholder={searchPlaceholder} />
			<Command.List>
				<Command.Empty class="text-muted-foreground">{emptyMessage}</Command.Empty>
				{#each countryOptions as country (country.code)}
					<Command.Item
						value={country.code}
						keywords={country.keywords}
						data-checked={selectedCountry?.code === country.code}
						onSelect={() => selectCountry(country)}
					>
						<img
							src={country.flagSrc}
							alt=""
							width="24"
							height="24"
							loading="lazy"
							class="size-6 shrink-0 rounded-sm object-cover ring-1 ring-border/40"
						/>
						<span class="min-w-0 flex-1 truncate">{getCountryName(country, locale)}</span>
						<span class="shrink-0 text-xs font-semibold text-muted-foreground">
							{country.code}
						</span>
					</Command.Item>
				{/each}
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
