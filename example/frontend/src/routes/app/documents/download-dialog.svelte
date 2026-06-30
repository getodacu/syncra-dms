<script lang="ts">
	import BracesIcon from "@lucide/svelte/icons/braces";
	import DownloadIcon from "@lucide/svelte/icons/download";
	import FileTextIcon from "@lucide/svelte/icons/file-text";

	import * as Alert from "$lib/components/ui/alert/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import type { DownloadFormat, OCRDocumentListItemResponse } from "./api";

	let {
		open = $bindable(false),
		documents = [],
		pending = false,
		error = null,
		onDownload,
	}: {
		open: boolean;
		documents?: OCRDocumentListItemResponse[];
		pending?: boolean;
		error?: Error | string | null;
		onDownload: (format: DownloadFormat) => void;
	} = $props();

	const selectedCount = $derived(documents.length);
	const hasJSONOption = $derived(
		documents.some((document) => Boolean(document.schema_id || document.has_inline_schema))
	);
	const title = $derived(
		selectedCount === 1
			? m.documents_download_dialog_title_one()
			: m.documents_download_dialog_title_other({ count: selectedCount })
	);
	const description = $derived(
		selectedCount === 1 ? documents[0]?.original_filename : m.documents_selected_documents()
	);
	const errorMessage = $derived(typeof error === "string" ? error : error?.message);
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-sm">
		<Dialog.Header>
			<Dialog.Title>{title}</Dialog.Title>
			<Dialog.Description>{description}</Dialog.Description>
		</Dialog.Header>

		{#if errorMessage}
			<Alert.Root variant="destructive">
				<Alert.Description>{errorMessage}</Alert.Description>
			</Alert.Root>
		{/if}

		<div class="grid gap-2">
			<Button
				type="button"
				variant="outline"
				class="h-10 justify-start"
				disabled={pending}
				onclick={() => onDownload("markdown")}
			>
				<FileTextIcon class="size-4 text-muted-foreground" aria-hidden="true" />
				{m.documents_format_markdown()}
			</Button>
			<Button
				type="button"
				variant="outline"
				class="h-10 justify-start"
				disabled={pending}
				onclick={() => onDownload("html")}
			>
				<DownloadIcon class="size-4 text-muted-foreground" aria-hidden="true" />
				{m.documents_format_html()}
			</Button>
			{#if hasJSONOption}
				<Button
					type="button"
					variant="outline"
					class="h-10 justify-start"
					disabled={pending}
					onclick={() => onDownload("json")}
				>
					<BracesIcon class="size-4 text-muted-foreground" aria-hidden="true" />
					{m.documents_format_json()}
				</Button>
			{/if}
		</div>

		{#if pending}
			<div class="flex items-center gap-2 text-xs text-muted-foreground">
				<Spinner class="size-3.5" />
				<span>{m.documents_preparing_download()}</span>
			</div>
		{/if}
	</Dialog.Content>
</Dialog.Root>
