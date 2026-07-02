<script lang="ts">
	import UploadCloudIcon from '@lucide/svelte/icons/upload-cloud';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { UploadQueueItem } from './upload-queue';

	let {
		canCreate,
		selectedFolderId,
		queue,
		isUploading,
		onFilesSelected,
		onUpload
	}: {
		canCreate: boolean;
		selectedFolderId: string | null;
		queue: UploadQueueItem[];
		isUploading: boolean;
		onFilesSelected: (files: FileList) => void;
		onUpload: () => Promise<void>;
	} = $props();

	const hasPendingUploads = $derived(
		queue.some((item) => item.status === 'queued' || item.status === 'failed')
	);
	const canUpload = $derived(canCreate && selectedFolderId !== null && hasPendingUploads && !isUploading);

	function handleFileChange(event: Event) {
		const input = event.currentTarget as HTMLInputElement;
		if (input.files?.length) {
			onFilesSelected(input.files);
			input.value = '';
		}
	}

	async function submitUpload(event: SubmitEvent) {
		event.preventDefault();
		if (!canUpload) return;

		try {
			await onUpload();
		} catch {
		}
	}

	function statusVariant(status: UploadQueueItem['status']) {
		switch (status) {
			case 'queued':
				return 'outline';
			case 'uploading':
				return 'secondary';
			case 'uploaded':
				return 'default';
			case 'failed':
				return 'destructive';
		}
	}
</script>

<Card.Root>
	<Card.Header>
		<div class="flex items-start justify-between gap-3">
			<div>
				<Card.Title>Upload</Card.Title>
				<Card.Description>{selectedFolderId ? 'Selected folder' : 'No folder selected'}</Card.Description>
			</div>
			<UploadCloudIcon class="size-4 text-muted-foreground" />
		</div>
	</Card.Header>
	<Card.Content>
		<form class="grid gap-3" onsubmit={submitUpload}>
			<label class="grid gap-1.5 text-sm font-medium">
				Files
				<input
					class="h-9 rounded-md border bg-background px-3 py-1.5 text-sm file:mr-3 file:rounded-md file:border-0 file:bg-muted file:px-2 file:py-1 file:text-xs file:font-medium disabled:cursor-not-allowed disabled:opacity-60"
					type="file"
					multiple
					disabled={!canCreate || !selectedFolderId || isUploading}
					onchange={handleFileChange}
				/>
			</label>

			<div class="grid gap-2">
				{#each queue as item (item.id)}
					<div class="grid gap-1 rounded-md border bg-muted/20 px-3 py-2 text-sm sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center">
						<div class="min-w-0">
							<div class="truncate font-medium" title={item.file.name}>{item.file.name}</div>
							{#if item.error}
								<div class="truncate text-xs text-destructive" title={item.error}>{item.error}</div>
							{:else}
								<div class="text-xs text-muted-foreground">{item.file.size} bytes</div>
							{/if}
						</div>
						<Badge variant={statusVariant(item.status)} class="w-fit">{item.status}</Badge>
					</div>
				{:else}
					<div
						class="flex h-16 items-center justify-center rounded-md border border-dashed bg-muted/20 px-3 text-sm text-muted-foreground"
					>
						No files queued
					</div>
				{/each}
			</div>

			<Button type="submit" size="sm" class="w-fit gap-2" disabled={!canUpload}>
				<UploadCloudIcon class="size-4" />
				Upload
			</Button>
		</form>
	</Card.Content>
</Card.Root>
