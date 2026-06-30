<script lang="ts">
	import "highlight.js/styles/github-dark.css";

	import { Button } from "$lib/components/ui/button/index.js";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import * as Tabs from "$lib/components/ui/tabs/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import CopyIcon from "@tabler/icons-svelte/icons/copy";
	import CheckIcon from "@tabler/icons-svelte/icons/check";
	import LoaderIcon from "@lucide/svelte/icons/loader-circle";
	import PencilIcon from "@lucide/svelte/icons/pencil";

	import { renderHighlightedJSON, renderMarkdown, formatAnnotationJSON } from "./document-preview-utils";

	type Props = {
		open: boolean;
		filename?: string | null;
		markdown?: string | null;
		annotationJson?: unknown;
		isLoading?: boolean;
		error?: Error | string | null;
		onRetry?: () => void;
		onRename?: (originalFilename: string) => void | Promise<void>;
		renamePending?: boolean;
	};

	let {
		open = $bindable(false),
		filename = null,
		markdown = null,
		annotationJson = undefined,
		isLoading = false,
		error = null,
		onRetry,
		onRename,
		renamePending = false,
	}: Props = $props();

	let copiedMarkdown = $state(false);
	let copiedHTML = $state(false);
	let copiedJSON = $state(false);
	let renaming = $state(false);
	let draftFilename = $state("");
	let renameSaving = $state(false);
	let renameError = $state<string | null>(null);
	let renameInputRef = $state<HTMLInputElement | null>(null);

	const title = $derived(filename ?? m.documents_preview_fallback_title());
	const description = m.documents_preview_description();
	const errorMessage = $derived(typeof error === "string" ? error : error?.message);
	const hasPreview = $derived(markdown !== null && markdown !== undefined);
	const markdownHTML = $derived(renderMarkdown(markdown, m.documents_no_markdown_content()));
	const highlightedJSON = $derived(
		renderHighlightedJSON(annotationJson, m.documents_no_json_annotation())
	);
	const canRename = $derived(Boolean(filename && onRename));
	const renameEditorPending = $derived(renameSaving || renamePending);

	function closePreview() {
		open = false;
	}

	function startRenaming() {
		if (!canRename || renameEditorPending) return;

		draftFilename = filename ?? "";
		renameError = null;
		renaming = true;
	}

	function cancelRenaming() {
		if (renameEditorPending) return;

		draftFilename = filename ?? "";
		renameError = null;
		renaming = false;
	}

	async function commitRename() {
		const nextFilename = draftFilename.trim();
		if (!onRename || !filename || !nextFilename || nextFilename === filename) {
			cancelRenaming();
			return;
		}

		renameSaving = true;
		renameError = null;
		try {
			await onRename?.(nextFilename);
			renaming = false;
		} catch (error) {
			renameError = error instanceof Error ? error.message : m.documents_failed_rename();
		} finally {
			renameSaving = false;
		}
	}

	function handleRenameKeydown(event: KeyboardEvent) {
		if (event.key === "Enter") {
			event.preventDefault();
			void commitRename();
			return;
		}

		if (event.key === "Escape") {
			event.preventDefault();
			cancelRenaming();
		}
	}

	async function copyMarkdown() {
		try {
			await navigator.clipboard.writeText(markdown ?? "");
			copiedMarkdown = true;
			setTimeout(() => {
				copiedMarkdown = false;
			}, 2000);
		} catch (err) {
			console.error("Failed to copy markdown: ", err);
		}
	}

	async function copyHTML() {
		try {
			await navigator.clipboard.writeText(markdownHTML);
			copiedHTML = true;
			setTimeout(() => {
				copiedHTML = false;
			}, 2000);
		} catch (err) {
			console.error("Failed to copy HTML: ", err);
		}
	}

	async function copyJSON() {
		try {
			await navigator.clipboard.writeText(
				formatAnnotationJSON(annotationJson, m.documents_no_json_annotation())
			);
			copiedJSON = true;
			setTimeout(() => {
				copiedJSON = false;
			}, 2000);
		} catch (err) {
			console.error("Failed to copy JSON: ", err);
		}
	}

	$effect(() => {
		if (open) {
			copiedMarkdown = false;
			copiedHTML = false;
			copiedJSON = false;
			renaming = false;
			renameError = null;
			draftFilename = filename ?? "";
		}
	});

	$effect(() => {
		if (!renaming) {
			draftFilename = filename ?? "";
		}
	});

	$effect(() => {
		if (!renaming || !renameInputRef) return;

		renameInputRef.focus();
		renameInputRef.select();
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Content
		class="flex h-[min(92vh,80rem)] w-full flex-col overflow-hidden sm:max-w-[75vw] lg:max-w-[50vw]"
	>
		<Dialog.Header>
			{#if renaming && canRename}
				<div class="flex min-w-0 items-center gap-2">
					<Dialog.Title class="sr-only">{m.documents_rename_document_title()}</Dialog.Title>
					<Input
						bind:ref={renameInputRef}
						value={draftFilename}
						class="h-9 min-w-0 flex-1 px-2 text-xl font-medium"
						disabled={renameEditorPending}
						aria-label={m.documents_rename_file({ name: filename ?? "" })}
						aria-invalid={renameError ? "true" : undefined}
						title={renameError ?? filename ?? undefined}
						oninput={(event) => {
							draftFilename = event.currentTarget.value;
						}}
						onkeydown={handleRenameKeydown}
						onblur={cancelRenaming}
					/>
					{#if renameEditorPending}
						<LoaderIcon class="size-4 shrink-0 animate-spin text-muted-foreground" aria-hidden="true" />
					{/if}
				</div>
			{:else}
				<div class="flex min-w-0 items-center gap-2">
					<Dialog.Title class="min-w-0 truncate text-xl">{title}</Dialog.Title>
					{#if canRename}
						<Button
							type="button"
							variant="ghost"
							size="icon-sm"
							class="text-muted-foreground hover:text-foreground"
							disabled={renameEditorPending}
							aria-label={m.documents_rename_file({ name: filename ?? "" })}
							onclick={startRenaming}
						>
							<PencilIcon class="size-4" aria-hidden="true" />
						</Button>
					{/if}
				</div>
			{/if}
			{#if renameError}
				<p class="text-sm text-destructive">{renameError}</p>
			{/if}
			<Dialog.Description>{description}</Dialog.Description>
		</Dialog.Header>

		<div class="min-h-0 flex-1 overflow-y-auto pr-2">
			{#if isLoading}
				<p class="text-sm text-muted-foreground">{m.documents_loading_document()}</p>
			{:else if errorMessage}
				<div class="flex flex-wrap items-center gap-3 text-sm text-destructive">
					<span>{errorMessage}</span>
					{#if onRetry}
						<Button type="button" variant="outline" size="sm" onclick={onRetry}>{m.documents_retry()}</Button>
					{/if}
				</div>
			{:else if hasPreview}
				<Tabs.Root value="markdown" class="w-full">
					<Tabs.List>
						<Tabs.Trigger value="markdown">{m.documents_format_markdown()}</Tabs.Trigger>
						<Tabs.Trigger value="json">{m.documents_format_json()}</Tabs.Trigger>
					</Tabs.List>

					<Tabs.Content value="markdown" class="mt-3 min-h-0 relative">
						<div class="rounded-md border bg-muted/30 p-5 pt-24 sm:pt-12 relative">
							<div class="absolute top-3 right-3 left-3 flex flex-wrap justify-end gap-2 sm:left-auto">
								<Button
									type="button"
									variant="outline"
									size="sm"
									onclick={copyMarkdown}
									class="flex h-8 items-center gap-1.5 px-2.5 text-xs text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
								>
									{#if copiedMarkdown}
										<CheckIcon class="h-3.5 w-3.5 text-emerald-500" />
										<span class="font-medium text-emerald-600 dark:text-emerald-400">{m.documents_copied()}</span>
									{:else}
										<CopyIcon class="h-3.5 w-3.5" />
										<span>{m.documents_copy_markdown()}</span>
									{/if}
								</Button>
								<Button
									type="button"
									variant="outline"
									size="sm"
									onclick={copyHTML}
									class="flex h-8 items-center gap-1.5 px-2.5 text-xs text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
								>
									{#if copiedHTML}
										<CheckIcon class="h-3.5 w-3.5 text-emerald-500" />
										<span class="font-medium text-emerald-600 dark:text-emerald-400">{m.documents_copied()}</span>
									{:else}
										<CopyIcon class="h-3.5 w-3.5" />
										<span>{m.documents_copy_html()}</span>
									{/if}
								</Button>
							</div>
							<div class="markdown-body">
								{@html markdownHTML}
							</div>
						</div>
					</Tabs.Content>
					<Tabs.Content value="json" class="mt-3 min-h-0 relative">
						<div class="relative">
							<Button
								type="button"
								variant="outline"
								size="sm"
								onclick={copyJSON}
								class="absolute top-3 right-3 flex h-8 items-center gap-1.5 px-2.5 text-xs text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
							>
								{#if copiedJSON}
									<CheckIcon class="h-3.5 w-3.5 text-emerald-500" />
									<span class="font-medium text-emerald-600 dark:text-emerald-400">{m.documents_copied()}</span>
								{:else}
									<CopyIcon class="h-3.5 w-3.5" />
									<span>{m.documents_copy_json()}</span>
								{/if}
							</Button>
							<pre
								class="overflow-x-auto rounded-md border bg-muted/30 p-4 pt-12 text-xs whitespace-pre-wrap"
							><code class="hljs language-json bg-transparent p-0">{@html highlightedJSON}</code></pre>
						</div>
					</Tabs.Content>
				</Tabs.Root>
			{:else}
				<p class="text-sm text-muted-foreground">{m.documents_no_preview_available()}</p>
			{/if}
		</div>

		<Dialog.Footer class="mt-3 flex justify-end">
			<Button type="button" variant="outline" onclick={closePreview}>{m.documents_close()}</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<style>
	.markdown-body {
		font-family: var(--font-sans);
		color: var(--foreground);
		line-height: 1.625;
		font-size: 0.875rem;
	}

	.markdown-body :global(h1),
	.markdown-body :global(h2),
	.markdown-body :global(h3),
	.markdown-body :global(h4),
	.markdown-body :global(h5),
	.markdown-body :global(h6) {
		font-weight: 600;
		line-height: 1.25;
		margin-top: 1.5rem;
		margin-bottom: 0.75rem;
		color: var(--foreground);
	}

	.markdown-body :global(h1:first-child),
	.markdown-body :global(h2:first-child),
	.markdown-body :global(h3:first-child) {
		margin-top: 0;
	}

	.markdown-body :global(h1) {
		font-size: 1.5rem;
		border-bottom: 1px solid var(--border);
		padding-bottom: 0.3em;
	}
	.markdown-body :global(h2) {
		font-size: 1.25rem;
		border-bottom: 1px solid var(--border);
		padding-bottom: 0.3em;
	}
	.markdown-body :global(h3) {
		font-size: 1.125rem;
	}
	.markdown-body :global(h4) {
		font-size: 1rem;
	}

	.markdown-body :global(p) {
		margin-top: 0;
		margin-bottom: 0.875rem;
	}

	.markdown-body :global(p:last-child) {
		margin-bottom: 0;
	}

	.markdown-body :global(a) {
		color: var(--primary);
		text-decoration: underline;
		text-underline-offset: 4px;
		font-weight: 500;
		transition: opacity 150ms ease;
	}

	.markdown-body :global(a:hover) {
		opacity: 0.8;
	}

	.markdown-body :global(ul),
	.markdown-body :global(ol) {
		margin-top: 0;
		margin-bottom: 0.875rem;
		padding-left: 1.5rem;
	}

	.markdown-body :global(ul) {
		list-style-type: disc;
	}

	.markdown-body :global(ol) {
		list-style-type: decimal;
	}

	.markdown-body :global(li) {
		margin-top: 0.25rem;
		margin-bottom: 0.25rem;
	}

	.markdown-body :global(code) {
		font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono",
			"Courier New", monospace;
		font-size: 0.85em;
		background-color: var(--muted);
		padding: 0.2em 0.4em;
		border-radius: 4px;
		border: 1px solid var(--border);
	}

	.markdown-body :global(pre) {
		margin-top: 0;
		margin-bottom: 0.875rem;
		padding: 1rem;
		overflow-x: auto;
		background-color: var(--muted);
		border-radius: 6px;
		border: 1px solid var(--border);
	}

	.markdown-body :global(pre code) {
		background-color: transparent;
		padding: 0;
		font-size: 0.85em;
		border-radius: 0;
		border: none;
	}

	.markdown-body :global(blockquote) {
		margin: 0 0 0.875rem 0;
		padding: 0 1rem;
		color: var(--muted-foreground);
		border-left: 0.25rem solid var(--muted-foreground);
		font-style: italic;
	}

	.markdown-body :global(table) {
		display: block;
		width: 100%;
		overflow: auto;
		margin-top: 0;
		margin-bottom: 0.875rem;
		border-spacing: 0;
		border-collapse: collapse;
	}

	.markdown-body :global(tr) {
		background-color: var(--background);
		border-top: 1px solid var(--border);
	}

	.markdown-body :global(tr:nth-child(2n)) {
		background-color: var(--muted);
	}

	.markdown-body :global(th),
	.markdown-body :global(td) {
		padding: 6px 13px;
		border: 1px solid var(--border);
	}

	.markdown-body :global(th) {
		font-weight: 600;
	}
</style>
