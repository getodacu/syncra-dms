<script lang="ts">
	import CheckIcon from "@lucide/svelte/icons/check";
	import ChevronsUpDownIcon from "@lucide/svelte/icons/chevrons-up-down";
	import CircleAlertIcon from "@lucide/svelte/icons/circle-alert";
	import CircleCheckIcon from "@lucide/svelte/icons/circle-check";
	import CloudUploadIcon from "@lucide/svelte/icons/cloud-upload";
	import CreditCardIcon from "@lucide/svelte/icons/credit-card";
	import DownloadIcon from "@lucide/svelte/icons/download";
	import EyeIcon from "@lucide/svelte/icons/eye";
	import FileIcon from "@lucide/svelte/icons/file";
	import FileImageIcon from "@lucide/svelte/icons/file-image";
	import FileTextIcon from "@lucide/svelte/icons/file-text";
	import LoaderIcon from "@lucide/svelte/icons/loader-circle";
	import RocketIcon from "@lucide/svelte/icons/rocket";
	import SearchIcon from "@lucide/svelte/icons/search";
	import Trash2Icon from "@lucide/svelte/icons/trash-2";
	import XIcon from "@lucide/svelte/icons/x";
	import { page } from "$app/state";
	import { createQuery, useQueryClient } from "@tanstack/svelte-query";
	import { onDestroy } from "svelte";

	import { publicApiErrorMessage } from "$lib/client/api-errors";
	import { CREDIT_BALANCE_QUERY_KEY } from "$lib/client/billing";
	import { Button, buttonVariants } from "$lib/components/ui/button/index.js";
	import * as Command from "$lib/components/ui/command/index.js";
	import { DocumentPreviewDialog } from "$lib/components/ui/document-preview/index.js";
	import { Field, FieldDescription, FieldError } from "$lib/components/ui/field/index.js";
	import * as FileDropZone from "$lib/components/extra-ui/file-drop-zone/index.js";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import * as Tooltip from "$lib/components/ui/tooltip/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import { cn } from "$lib/utils.js";
	import {
		canSubmitJobs,
		displayJobStatus,
		getSubmittableRows,
		isInsufficientCreditsError,
		isTerminalStatus,
		jobStatusPatch,
		MAX_UNKNOWN_STATUS_POLLS,
		percentComplete,
		shouldShowSubmitButton,
		type CreateStatus,
		type JobDisplayStatus,
		type JobStatus,
	} from "./workflow";
	import { buildSchemaTree } from "./schema-tree";
	import { isSelectedSchemaIdUnavailable, schemaIdFromSearchParams } from "./schema-query";

	const ACCEPTED_FILE_TYPES = ".pdf,.png,.jpg,.jpeg,application/pdf,image/png,image/jpeg";
	const MAX_FILE_SIZE = 20 << 20;
	const MAX_FILES = 100;
	const POLL_INTERVAL_MS = 2500;
	const MAX_POLL_FAILURES = 3;

	type JsonSchemaObject = Record<string, unknown>;

	type SchemaResponse = {
		id: string;
		name: string;
		description: string;
		strict: boolean;
		schema: JsonSchemaObject;
	};

	type OCRJobResponse = {
		id: string;
		file_size: number;
		page_count: number;
		document_id: string | null;
		status: string;
		error_message?: string;
	};

	type OCRDocumentResponse = {
		id: string;
		original_filename: string;
		markdown: string;
		annotation_json?: unknown;
	};

	type JobRow = {
		localId: string;
		file: File;
		submitted: boolean;
		createStatus: CreateStatus;
		jobStatus: JobStatus;
		jobId?: string;
		documentId?: string;
		errorMessage?: string;
		pollFailures: number;
	};

	let selectedSchemaId = $state(schemaIdFromSearchParams(page.url.searchParams));
	let schemaPopoverOpen = $state(false);
	let rows = $state<JobRow[]>([]);
	let fileError = $state<string | null>(null);
	let hadFileRejection = false;
	let isSubmitting = $state(false);
	let isPolling = $state(false);
	let pollTimer = $state<number | null>(null);
	let pollingGeneration = 0;
	let previewOpen = $state(false);
	let previewDocumentId = $state<string | null>(null);
	let previewFileName = $state<string | null>(null);
	let isDestroyed = false;

	const queryClient = useQueryClient();
	const schemasQuery = createQuery<SchemaResponse[], Error>(() => ({
		queryKey: ["schemas", "mine"],
		queryFn: () => fetchMySchemas(),
	}));
	const documentQuery = createQuery<OCRDocumentResponse, Error>(() => ({
		queryKey: ["ocr-document", previewDocumentId],
		queryFn: () => {
			if (!previewDocumentId) throw new Error(m.new_job_missing_document_id());
			return fetchOCRDocument(previewDocumentId);
		},
		enabled: previewOpen && previewDocumentId !== null,
		staleTime: 1000 * 60 * 2,
	}));

	const schemas = $derived(schemasQuery.data ?? []);
	const selectedSchema = $derived(schemas.find((schema) => schema.id === selectedSchemaId));
	const selectedSchemaLabel = $derived(selectedSchema?.name ?? m.new_job_select_schema_placeholder());
	const submittableRows = $derived(getSubmittableRows(rows));
	const submittedRows = $derived(rows.filter((row) => row.submitted));
	const unsubmittedRows = $derived(rows.filter((row) => !row.submitted));
	const showSubmitButton = $derived(shouldShowSubmitButton(selectedSchemaId, rows.length));
	const canSubmit = $derived(
		canSubmitJobs({
			selectedSchemaId,
			rows,
			isSubmitting,
			isPolling,
		})
	);

	const completedCount = $derived(submittedRows.filter(r => displayJobStatus(r) === "completed").length);
	const processingCount = $derived(submittedRows.filter(r => displayJobStatus(r) === "processing").length);
	const queueingCount = $derived(submittedRows.filter(r => displayJobStatus(r) === "queueing").length);
	const failedCount = $derived(submittedRows.filter(r => displayJobStatus(r) === "failed").length);
	const totalSubmitted = $derived(submittedRows.length);
	const overallProgress = $derived(percentComplete(completedCount + failedCount, totalSubmitted));
	const schemaTree = $derived(selectedSchema ? buildSchemaTree(selectedSchema.schema) : []);
	const selectedSchemaUnavailable = $derived(
		isSelectedSchemaIdUnavailable({
			selectedSchemaId,
			schemas,
			schemasLoaded: schemasQuery.isSuccess,
		})
	);

	$effect(() => {
		if (selectedSchemaUnavailable) selectedSchemaId = "";
	});

	async function readResponseJSON(response: Response) {
		const text = await response.text();
		if (!text) return null;
		try {
			return JSON.parse(text) as unknown;
		} catch {
			return null;
		}
	}

	function createJobErrorMessage(error: unknown) {
		const message = error instanceof Error ? error.message : m.new_job_failed_create();
		if (isInsufficientCreditsError(message)) {
			return m.new_job_insufficient_credits_buy();
		}
		return message;
	}

	async function fetchOCRDocument(id: string) {
		const response = await fetch(`/api/ocr/document/${encodeURIComponent(id)}`);
		const data = await readResponseJSON(response);
		if (!response.ok) {
			throw new Error(publicApiErrorMessage(response.status, data, m.new_job_failed_load_document()));
		}
		if (!isOCRDocumentResponse(data)) throw new Error(m.new_job_invalid_document_response());
		return data;
	}

	function isOCRDocumentResponse(value: unknown): value is OCRDocumentResponse {
		return (
			isRecord(value) &&
			typeof value.id === "string" &&
			typeof value.original_filename === "string" &&
			typeof value.markdown === "string"
		);
	}

	function openDocumentPreview(row: JobRow) {
		if (!row.documentId) return;
		previewFileName = row.file.name;
		previewDocumentId = row.documentId;
		previewOpen = true;
	}

	function isRecord(value: unknown): value is Record<string, unknown> {
		return typeof value === "object" && value !== null && !Array.isArray(value);
	}

	function isSchemaList(value: unknown): value is SchemaResponse[] {
		return (
			Array.isArray(value) &&
			value.every(
				(item) =>
					isRecord(item) &&
					typeof item.id === "string" &&
					typeof item.name === "string" &&
					typeof item.description === "string" &&
					typeof item.strict === "boolean" &&
					isRecord(item.schema)
			)
		);
	}

	function isOCRJobResponse(value: unknown): value is OCRJobResponse {
		return (
			isRecord(value) &&
			typeof value.id === "string" &&
			typeof value.file_size === "number" &&
			typeof value.page_count === "number" &&
			"document_id" in value &&
			(value.document_id === null || typeof value.document_id === "string") &&
			typeof value.status === "string" &&
			!("error_message" in value && typeof value.error_message !== "string")
		);
	}

	async function fetchMySchemas(fetchFn: typeof fetch = fetch) {
		const response = await fetchFn("/api/schemas?scope=mine");
		const data = await readResponseJSON(response);
		if (!response.ok) {
			throw new Error(publicApiErrorMessage(response.status, data, m.new_job_failed_load_schemas()));
		}
		if (!isSchemaList(data)) throw new Error(m.new_job_invalid_schema_response());
		return data;
	}

	async function createJob(file: File, schemaId?: string) {
		const formData = new FormData();
		formData.set("file", file);
		if (schemaId) formData.set("schema_id", schemaId);

		const response = await fetch("/api/ocr/jobs", { method: "POST", body: formData });
		const data = await readResponseJSON(response);
		if (!response.ok) {
			throw new Error(publicApiErrorMessage(response.status, data, m.new_job_failed_create()));
		}
		if (!isOCRJobResponse(data)) throw new Error(m.new_job_invalid_job_response());
		return data;
	}

	async function fetchJob(id: string) {
		const response = await fetch(`/api/ocr/jobs/${encodeURIComponent(id)}`);
		const data = await readResponseJSON(response);
		if (!response.ok) {
			throw new Error(publicApiErrorMessage(response.status, data, m.new_job_failed_load_job()));
		}
		if (!isOCRJobResponse(data)) throw new Error(m.new_job_invalid_job_response());
		return data;
	}

	function updateRow(localId: string, patch: Partial<JobRow>) {
		rows = rows.map((row) => (row.localId === localId ? { ...row, ...patch } : row));
	}

	async function handleUpload(files: File[]) {
		if (!hadFileRejection) fileError = null;

		const remaining = MAX_FILES - rows.length;
		const accepted = files.slice(0, remaining);
		rows = [
			...rows,
			...accepted.map((file) => ({
				localId: crypto.randomUUID(),
				file,
				submitted: false,
				createStatus: "pending" as const,
				jobStatus: "unknown" as const,
				pollFailures: 0,
			})),
		];
		hadFileRejection = false;
	}

	function handleFileRejected(event: { reason: FileDropZone.FileRejectedReason; file: File }) {
		hadFileRejection = true;
		fileError = `${event.file.name}: ${event.reason}`;
	}

	function removeRow(localId: string) {
		if (isSubmitting || isPolling) return;
		rows = rows.filter((row) => row.localId !== localId);
	}

	function clearUnsubmittedRows() {
		if (isSubmitting || isPolling) return;
		rows = rows.filter((row) => row.submitted);
		fileError = null;
	}

	async function submitJobs() {
		if (submittableRows.length === 0 || isSubmitting) return;

		const schemaId = selectedSchemaId || undefined;
		const targetRows = submittableRows;

		isSubmitting = true;
		stopPolling();
		rows = rows.map((row) => ({
			...row,
			...(targetRows.some((targetRow) => targetRow.localId === row.localId)
				? {
						submitted: true,
						createStatus: "pending" as const,
						jobStatus: "unknown" as const,
						jobId: undefined,
						documentId: undefined,
						errorMessage: undefined,
						pollFailures: 0,
					}
				: {}),
		}));

		for (const row of targetRows) {
			if (isDestroyed) return;

			updateRow(row.localId, { createStatus: "creating" });
			try {
				const job = await createJob(row.file, schemaId);
				if (isDestroyed) return;

				updateRow(row.localId, {
					createStatus: "queued",
					jobId: job.id,
					...jobStatusPatch(job, 0, MAX_UNKNOWN_STATUS_POLLS),
				});
			} catch (error) {
				if (isDestroyed) return;

				updateRow(row.localId, {
					createStatus: "create_failed",
					jobStatus: "failed",
					errorMessage: createJobErrorMessage(error),
				});
			}
		}

		isSubmitting = false;
		startPolling();
	}

	function activeJobRows() {
		return rows.filter((row) => row.jobId && !isTerminalStatus(row.jobStatus));
	}

	function stopPolling() {
		pollingGeneration += 1;
		if (pollTimer) window.clearTimeout(pollTimer);
		pollTimer = null;
		isPolling = false;
	}

	function startPolling() {
		if (isDestroyed) return;

		if (activeJobRows().length === 0) {
			stopPolling();
			return;
		}

		isPolling = true;
		if (pollTimer) window.clearTimeout(pollTimer);
		const generation = pollingGeneration + 1;
		pollingGeneration = generation;
		pollTimer = window.setTimeout(() => void pollOnce(generation), POLL_INTERVAL_MS);
	}

	async function pollOnce(generation: number) {
		if (isDestroyed || generation !== pollingGeneration) return;

		const activeRows = activeJobRows();
		if (activeRows.length === 0) {
			stopPolling();
			return;
		}
		let completedJob = false;

		await Promise.all(
			activeRows.map(async (row) => {
				if (!row.jobId) return;
				try {
					const job = await fetchJob(row.jobId);
					if (isDestroyed || generation !== pollingGeneration) return;
					const patch = jobStatusPatch(job, row.pollFailures, MAX_UNKNOWN_STATUS_POLLS);
					if (row.jobStatus !== "completed" && patch.jobStatus === "completed") {
						completedJob = true;
					}

					updateRow(row.localId, patch);
				} catch (error) {
					if (isDestroyed || generation !== pollingGeneration) return;

					const failures = row.pollFailures + 1;
					updateRow(row.localId, {
						pollFailures: failures,
						errorMessage: error instanceof Error ? error.message : m.new_job_failed_poll_job(),
						jobStatus: failures >= MAX_POLL_FAILURES ? "failed" : row.jobStatus,
					});
				}
			})
		);

		if (isDestroyed || generation !== pollingGeneration) return;
		if (completedJob) {
			void queryClient.invalidateQueries({ queryKey: CREDIT_BALANCE_QUERY_KEY });
		}
		startPolling();
	}

	function selectSchema(schemaId: string) {
		if (isSubmitting || isPolling) return;

		selectedSchemaId = schemaId;
		schemaPopoverOpen = false;
	}

	function fileCountText(count: number) {
		return count === 1
			? m.new_job_file_count_one({ count })
			: m.new_job_file_count_other({ count });
	}

	function isPDFFile(file: File) {
		return file.type === "application/pdf" || file.name.toLowerCase().endsWith(".pdf");
	}

	function isImageFile(file: File) {
		return file.type.startsWith("image/") || /\.(png|jpe?g)$/i.test(file.name);
	}

	function fileIconClass(file: File) {
		if (isPDFFile(file)) return "border-blue-100 bg-blue-50 text-blue-600";
		if (isImageFile(file)) return "border-rose-100 bg-rose-50 text-rose-600";
		return "border-indigo-100 bg-indigo-50 text-indigo-600";
	}

	function statusText(status: JobDisplayStatus) {
		if (status === "queueing") return m.jobs_status_queued();
		if (status === "processing") return m.jobs_status_processing();
		if (status === "completed") return m.jobs_status_completed();
		if (status === "failed") return m.jobs_status_failed();
		return status;
	}

	function statusPillClass(status: JobDisplayStatus) {
		return cn(
			"inline-flex h-7 items-center gap-1.5 rounded-full px-3 text-[11px] font-bold uppercase tracking-wider",
			status === "queueing" && "bg-sky-500/10 text-sky-600 dark:text-sky-400 ring-1 ring-sky-500/20",
			status === "processing" && "bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 ring-1 ring-indigo-500/20",
			status === "completed" && "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 ring-1 ring-emerald-500/20",
			status === "failed" && "bg-red-500/10 text-red-600 dark:text-red-400 ring-1 ring-red-500/20"
		);
	}

	function rowCardClass(status: JobDisplayStatus) {
		return cn(
			"grid min-h-20 grid-cols-[auto_minmax(0,1fr)] items-center gap-4 rounded-xl border bg-card/60 backdrop-blur-xs p-4 text-sm shadow-2xs hover:shadow-xs transition-all duration-300 hover:border-primary/20 sm:grid-cols-[auto_minmax(0,1fr)_auto]",
			status === "queueing" && "border-l-4 border-l-sky-500 bg-sky-500/2",
			status === "processing" && "border-l-4 border-l-indigo-500 bg-indigo-500/2",
			status === "completed" && "border-l-4 border-l-emerald-500 bg-emerald-500/2",
			status === "failed" && "border-l-4 border-l-red-500 bg-red-500/2"
		);
	}

	function rowDetail(row: JobRow) {
		const status = displayJobStatus(row);
		if (status === "failed" && isInsufficientCreditsError(row.errorMessage)) {
			return m.new_job_insufficient_credits_document();
		}
		if (status === "failed") return row.errorMessage ?? m.new_job_processing_failed();
		if (status === "completed") {
			return row.documentId ? m.new_job_document_id({ id: row.documentId }) : m.new_job_processed();
		}
		if (status === "queueing") {
			return row.createStatus === "creating"
				? m.new_job_creating_job()
				: m.new_job_queued_processing();
		}
		return row.errorMessage ?? m.new_job_extracting_entities();
	}

	function canRemoveRow(row: JobRow) {
		return (
			!isSubmitting &&
			!isPolling &&
			(!row.submitted || row.createStatus === "create_failed" || row.jobStatus === "failed")
		);
	}

	onDestroy(() => {
		isDestroyed = true;
		stopPolling();
	});
</script>

<div class="@container/main flex flex-1 flex-col bg-muted/20">
	<main class="mx-auto flex w-full max-w-5xl flex-1 flex-col gap-8 px-4 py-6 sm:px-6 lg:px-8 lg:py-8">

		<!-- Interactive Timeline Stepper -->
		<div class="grid grid-cols-1 gap-4 rounded-xl border bg-card p-5 shadow-2xs md:grid-cols-3">
			<!-- Step 1 -->
			<div class="flex items-center gap-3.5 px-2 py-1.5 transition-all duration-300">
				<div class={cn(
					"flex size-10 items-center justify-center rounded-full font-bold text-sm shadow-inner transition-all duration-300 border-2 shrink-0",
					selectedSchemaId !== "" 
						? "bg-indigo-50 border-indigo-200 text-indigo-600 dark:bg-indigo-950/30 dark:border-indigo-900/50 dark:text-indigo-400" 
						: "bg-muted border-muted-foreground/10 text-muted-foreground"
				)}>
					{#if selectedSchemaId !== ""}
						<CheckIcon class="size-5" />
					{:else}
						01
					{/if}
				</div>
				<div class="min-w-0">
					<p class={cn("text-sm font-semibold truncate", selectedSchemaId !== "" ? "text-foreground" : "text-muted-foreground")}>{m.new_job_select_schema()}</p>
					<p class="text-xs text-muted-foreground truncate">{selectedSchemaId ? selectedSchemaLabel : m.new_job_configure_payload_format()}</p>
				</div>
			</div>

			<!-- Step 2 -->
			<div class="flex items-center gap-3.5 px-2 py-1.5 transition-all duration-300">
				<div class={cn(
					"flex size-10 items-center justify-center rounded-full font-bold text-sm shadow-inner transition-all duration-300 border-2 shrink-0",
					rows.length > 0 
						? "bg-indigo-50 border-indigo-200 text-indigo-600 dark:bg-indigo-950/30 dark:border-indigo-900/50 dark:text-indigo-400" 
						: selectedSchemaId !== "" 
							? "bg-background border-indigo-500/50 text-indigo-500 animate-pulse" 
							: "bg-muted border-muted-foreground/10 text-muted-foreground"
				)}>
					{#if rows.length > 0}
						<CheckIcon class="size-5" />
					{:else}
						02
					{/if}
				</div>
				<div class="min-w-0">
					<p class={cn("text-sm font-semibold truncate", rows.length > 0 ? "text-foreground" : "text-muted-foreground")}>{m.new_job_upload_documents()}</p>
					<p class="text-xs text-muted-foreground truncate">
						{#if rows.length > 0}
							{rows.length === 1
								? m.new_job_files_selected_one({ count: rows.length })
								: m.new_job_files_selected_other({ count: rows.length })}
						{:else}
							{m.new_job_drag_or_browse_files()}
						{/if}
					</p>
				</div>
			</div>

			<!-- Step 3 -->
			<div class="flex items-center gap-3.5 px-2 py-1.5 transition-all duration-300">
				<div class={cn(
					"flex size-10 items-center justify-center rounded-full font-bold text-sm shadow-inner transition-all duration-300 border-2 shrink-0",
					submittedRows.length > 0 && submittedRows.every(row => isTerminalStatus(row.jobStatus))
						? "bg-emerald-50 border-emerald-200 text-emerald-600 dark:bg-emerald-950/30 dark:border-emerald-900/50 dark:text-emerald-400"
						: (isSubmitting || isPolling)
							? "bg-background border-indigo-500/50 text-indigo-500 animate-pulse"
							: "bg-muted border-muted-foreground/10 text-muted-foreground"
				)}>
					{#if submittedRows.length > 0 && submittedRows.every(row => isTerminalStatus(row.jobStatus))}
						<CircleCheckIcon class="size-5" />
					{:else}
						03
					{/if}
				</div>
				<div class="min-w-0">
					<p class={cn("text-sm font-semibold truncate", submittedRows.length > 0 ? "text-foreground" : "text-muted-foreground")}>{m.new_job_run_monitor()}</p>
					<p class="text-xs text-muted-foreground truncate">
						{#if isPolling}
							{m.new_job_processing_batch()}
						{:else}
							{m.new_job_start_extraction_pipeline()}
						{/if}
					</p>
				</div>
			</div>
		</div>

		<!-- Step 1: Select Schema Card -->
		<section class="rounded-xl border bg-card p-6 shadow-2xs transition-all duration-300 hover:border-primary/20">
			<div class="mb-4 flex flex-col gap-1">
				<h2 class="text-lg font-bold tracking-tight text-foreground flex items-center gap-2">
					<span class="flex size-6 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary">1</span>
					{m.new_job_select_extraction_schema()}
				</h2>
				<p class="text-sm text-muted-foreground">{m.new_job_select_schema_description()}</p>
			</div>

			<Field>
				<Popover.Root bind:open={schemaPopoverOpen}>
					<Popover.Trigger
						class={cn(
							buttonVariants({ variant: "outline" }),
							"h-14 w-full justify-between rounded-lg border-input bg-background px-4 text-left text-base shadow-2xs hover:bg-accent/40 transition-colors cursor-pointer"
						)}
						aria-label={m.new_job_select_extraction_schema_aria()}
						disabled={isSubmitting || isPolling}
					>
						<span class="flex min-w-0 items-center gap-3">
							<SearchIcon class="size-5 shrink-0 text-muted-foreground" />
							<span class="min-w-0 truncate font-semibold">{selectedSchemaLabel}</span>
						</span>
						<ChevronsUpDownIcon class="size-5 shrink-0 text-muted-foreground" />
					</Popover.Trigger>
					<Popover.Content class="w-[min(calc(100vw-2rem),48rem)] p-0" align="start">
						<Command.Root>
							<Command.Input placeholder={m.new_job_search_schemas()} />
							<Command.List>
								{#if schemasQuery.isLoading}
									<div class="px-3 py-6 text-center text-sm text-muted-foreground">
										{m.new_job_loading_schemas()}
									</div>
								{:else if schemasQuery.isError}
									<div class="px-3 py-6 text-center text-sm text-destructive">
										{schemasQuery.error.message}
									</div>
								{:else}
									<Command.Empty>{m.new_job_no_schemas_found()}</Command.Empty>
									<Command.Item
										value="__no_schema__"
										keywords={["No schema", "OCR only", "plain text"]}
										disabled={isSubmitting || isPolling}
										onSelect={() => selectSchema("")}
									>
										<div class="flex min-w-0 flex-col">
											<span class="truncate">{m.new_job_no_schema_ocr_only()}</span>
											<span class="truncate text-xs text-muted-foreground">
												{m.new_job_no_schema_description()}
											</span>
										</div>
										{#if selectedSchemaId === ""}
											<CheckIcon class="ml-auto size-4" />
										{/if}
									</Command.Item>
									{#each schemas as schema (schema.id)}
										<Command.Item
											value={schema.id}
											keywords={[schema.name, schema.description, schema.id]}
											disabled={isSubmitting || isPolling}
											onSelect={() => selectSchema(schema.id)}
										>
											<div class="flex min-w-0 flex-col">
												<span class="truncate">{schema.name}</span>
												{#if schema.description}
													<span class="truncate text-xs text-muted-foreground">
														{schema.description}
													</span>
												{/if}
											</div>
											{#if selectedSchemaId === schema.id}
												<CheckIcon class="ml-auto size-4" />
											{/if}
										</Command.Item>
									{/each}
								{/if}
							</Command.List>
						</Command.Root>
					</Popover.Content>
				</Popover.Root>
				{#if schemasQuery.isError}
					<FieldError>{schemasQuery.error.message}</FieldError>
				{:else if !schemasQuery.isLoading && schemas.length === 0}
					<FieldDescription>
						{m.new_job_no_personal_schemas()}
						<a class="font-medium underline underline-offset-4" href="/app/schemas/new">
							{m.new_job_create_one()}
						</a>
					</FieldDescription>
				{:else}
					<FieldDescription>
						{#if selectedSchema}
							{m.new_job_selected_schema_help()}
						{:else}
							{m.new_job_no_schema_selected_help()}
						{/if}
					</FieldDescription>
				{/if}
			</Field>

			<!-- Dynamic Schema Visualization Details -->
			{#if selectedSchema}
				<div class="mt-4 rounded-xl bg-muted/40 p-4.5 border border-border/50 space-y-3.5 transition-all duration-300">
					<div class="flex items-center justify-between">
						<span class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">{m.new_job_target_mapped_fields({ count: schemaTree.length })}</span>
						{#if selectedSchema.strict}
							<span class="inline-flex items-center gap-1 rounded-md bg-amber-500/10 px-2 py-0.5 text-[10px] font-bold text-amber-600 dark:bg-amber-500/20 dark:text-amber-400">
								{m.schemas_strict_mode()}
							</span>
						{/if}
					</div>
					{#if schemaTree.length > 0}
						<div class="flex flex-wrap gap-2">
							{#each schemaTree as node (node.path)}
								<span class="inline-flex items-center gap-1.5 rounded-full bg-background border px-3 py-1 text-xs font-semibold shadow-3xs hover:border-indigo-500/30 transition-colors">
									<span class="h-1.5 w-1.5 rounded-full bg-indigo-500"></span>
									<span class="text-foreground">{node.name}</span>
									<span class="text-muted-foreground text-[10px] font-mono">({node.type})</span>
									{#if node.required}
										<span class="text-red-500 font-bold text-[10px]" title={m.common_required()}>*</span>
									{/if}
								</span>
							{/each}
						</div>
					{:else}
						<p class="text-xs text-muted-foreground italic">{m.new_job_no_fields_defined()}</p>
					{/if}
				</div>
			{:else}
				<div class="mt-4 rounded-xl bg-amber-500/5 p-4.5 border border-amber-500/10 flex items-start gap-3 transition-all duration-300">
					<CircleAlertIcon class="size-5 text-amber-600 shrink-0 mt-0.5 dark:text-amber-400 animate-pulse" />
					<div>
						<p class="text-sm font-bold text-amber-800 dark:text-amber-300">{m.new_job_ocr_only_mode_active()}</p>
						<p class="text-xs text-amber-700/80 dark:text-amber-400/80 leading-normal">
							{m.new_job_ocr_only_mode_body()}
						</p>
					</div>
				</div>
			{/if}
		</section>

		<!-- Step 2: Upload Documents Card -->
		<section class="rounded-xl border bg-card p-6 shadow-2xs transition-all duration-300 hover:border-primary/20">
			<div class="mb-4 flex flex-col gap-1">
				<h2 class="text-lg font-bold tracking-tight text-foreground flex items-center gap-2">
					<span class="flex size-6 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary">2</span>
					{m.new_job_upload_documents()}
				</h2>
				<p class="text-sm text-muted-foreground">{m.new_job_upload_documents_description({ count: MAX_FILES })}</p>
			</div>

			<Field>
				<FileDropZone.Root
					accept={ACCEPTED_FILE_TYPES}
					maxFiles={MAX_FILES}
					fileCount={rows.length}
					maxFileSize={MAX_FILE_SIZE}
					onUpload={handleUpload}
					onFileRejected={handleFileRejected}
					disabled={isSubmitting || isPolling}
				>
					<FileDropZone.Trigger class="block">
						<div
							class="flex min-h-60 flex-col items-center justify-center gap-4 rounded-xl border-2 border-dashed border-muted-foreground/20 bg-background/30 px-6 py-8 text-center transition-all duration-300 cursor-pointer group-aria-disabled/file-drop-zone-trigger:opacity-60 hover:border-indigo-500/50 hover:bg-indigo-500/5 group-aria-disabled/file-drop-zone-trigger:hover:cursor-not-allowed"
						>
							<div class="flex size-14 items-center justify-center rounded-xl bg-indigo-500/10 text-indigo-600 transition-transform duration-300 group-hover/file-drop-zone-trigger:scale-110">
								<CloudUploadIcon class="size-7" />
							</div>
							<div class="space-y-1">
								<p class="text-base font-bold text-foreground">
									{m.new_job_dropzone_title()}
								</p>
								<p class="text-xs text-muted-foreground max-w-sm mx-auto">
									{m.new_job_dropzone_description({ size: FileDropZone.displaySize(MAX_FILE_SIZE) })}
								</p>
							</div>
							<span
								class={cn(
									buttonVariants({ variant: "outline" }),
									"pointer-events-none h-9.5 rounded-lg px-5 text-sm font-semibold shadow-2xs hover:bg-accent/40"
								)}
							>
								{m.new_job_browse_files()}
							</span>
						</div>
					</FileDropZone.Trigger>
				</FileDropZone.Root>
				{#if fileError}
					<FieldError>{fileError}</FieldError>
				{/if}
			</Field>

			<!-- Pending Upload Queue Grid -->
			{#if unsubmittedRows.length > 0}
				<div class="mt-6 space-y-3">
					<div class="flex items-center justify-between border-b pb-2">
						<span class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
							{m.new_job_pending_upload_queue({ count: unsubmittedRows.length })}
						</span>
						<Button
							type="button"
							variant="ghost"
							size="sm"
							class="h-8 text-xs text-muted-foreground hover:text-destructive hover:bg-destructive/5 cursor-pointer rounded-md transition-colors"
							disabled={isSubmitting || isPolling}
							onclick={clearUnsubmittedRows}
						>
							<XIcon class="mr-1 size-3.5" />
							{m.new_job_clear_all()}
						</Button>
					</div>

					<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
						{#each unsubmittedRows as row (row.localId)}
							<div class="group flex items-center justify-between gap-3 rounded-lg border bg-background/50 p-3 shadow-2xs hover:border-indigo-500/20 hover:shadow-xs transition-all duration-200">
								<div class="flex items-center gap-3 min-w-0">
									<div class={cn("flex size-9.5 shrink-0 items-center justify-center rounded-md border shadow-3xs", fileIconClass(row.file))}>
										{#if isImageFile(row.file)}
											<FileImageIcon class="size-4.5" />
										{:else if isPDFFile(row.file)}
											<FileTextIcon class="size-4.5" />
										{:else}
											<FileIcon class="size-4.5" />
										{/if}
									</div>
									<div class="min-w-0 space-y-0.5">
										<p class="truncate text-sm font-bold text-foreground" title={row.file.name}>{row.file.name}</p>
										<p class="text-xs text-muted-foreground">{FileDropZone.displaySize(row.file.size)}</p>
									</div>
								</div>
								
								<Button
									type="button"
									variant="ghost"
									size="icon-xs"
									class="h-7 w-7 rounded-md opacity-0 group-hover:opacity-100 transition-opacity hover:bg-destructive/5 hover:text-destructive cursor-pointer shrink-0"
									disabled={isSubmitting || isPolling}
									onclick={() => removeRow(row.localId)}
									aria-label={m.new_job_remove_file()}
								>
									<XIcon class="size-3.5" />
								</Button>
							</div>
						{/each}
					</div>
				</div>
			{/if}
		</section>

		<!-- Step 3: Job Status Dashboard -->
		{#if showSubmitButton}
		<section class="rounded-xl border bg-card p-6 shadow-2xs transition-all duration-300 hover:border-primary/20 space-y-5">
			<div class="flex items-center justify-between gap-3 border-b border-border/80 pb-4">
				<div class="flex items-center gap-2">
					<span class="flex size-6 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary">3</span>
					<h2 class="text-lg font-bold tracking-tight text-foreground">{m.new_job_extraction_queue_results()}</h2>
				</div>
				<span class="inline-flex items-center rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold text-muted-foreground">
					{m.new_job_total({ label: fileCountText(submittedRows.length) })}
				</span>
			</div>

			{#if submittedRows.length > 0}
				<!-- Real-time Batch Progress Dashboard -->
				<div class="rounded-xl border bg-background/30 p-4.5 space-y-4 shadow-3xs">
					<div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
						<div>
							<h3 class="text-sm font-bold text-foreground">{m.new_job_active_batch_status()}</h3>
							<p class="text-xs text-muted-foreground">{m.new_job_active_batch_description()}</p>
						</div>
						<div class="flex items-center gap-2">
							<span class="text-xs font-mono font-bold px-2 py-0.5 rounded bg-muted text-muted-foreground shadow-3xs">
								{m.new_job_progress({ progress: overallProgress })}
							</span>
						</div>
					</div>

					<!-- Linear Progress Bar -->
					<div class="relative w-full h-2 rounded-full bg-muted overflow-hidden">
						<div 
							class={cn(
								"h-full rounded-full transition-all duration-500 ease-out",
								failedCount > 0 ? "bg-linear-to-r from-indigo-500 to-red-500" : "bg-linear-to-r from-indigo-500 to-emerald-500",
								(isSubmitting || isPolling) && "animate-pulse"
							)}
							style="width: {overallProgress}%"
						></div>
					</div>

					<!-- Metrics Cards Grid -->
					<div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
						<!-- Metric 1: Total -->
						<div class="rounded-lg border bg-card/40 p-3 text-center space-y-0.5 shadow-3xs">
							<p class="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">{m.new_job_total_files()}</p>
							<p class="text-xl font-black text-foreground tabular-nums">{totalSubmitted}</p>
						</div>
						<!-- Metric 2: Completed -->
						<div class="rounded-lg border bg-card/40 p-3 text-center space-y-0.5 shadow-3xs hover:border-emerald-500/20 transition-colors">
							<p class="text-[10px] font-bold uppercase tracking-wider text-muted-foreground flex items-center justify-center gap-1">
								<span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span> {m.new_job_completed()}
							</p>
							<p class="text-xl font-black text-emerald-600 dark:text-emerald-400 tabular-nums">{completedCount}</p>
						</div>
						<!-- Metric 3: Processing -->
						<div class="rounded-lg border bg-card/40 p-3 text-center space-y-0.5 shadow-3xs hover:border-indigo-500/20 transition-colors">
							<p class="text-[10px] font-bold uppercase tracking-wider text-muted-foreground flex items-center justify-center gap-1">
								<span class="h-1.5 w-1.5 rounded-full bg-indigo-500 animate-pulse"></span> {m.new_job_processing()}
							</p>
							<p class="text-xl font-black text-indigo-600 dark:text-indigo-400 tabular-nums">{processingCount + queueingCount}</p>
						</div>
						<!-- Metric 4: Failed -->
						<div class="rounded-lg border bg-card/40 p-3 text-center space-y-0.5 shadow-3xs hover:border-red-500/20 transition-colors">
							<p class="text-[10px] font-bold uppercase tracking-wider text-muted-foreground flex items-center justify-center gap-1">
								<span class="h-1.5 w-1.5 rounded-full bg-red-500"></span> {m.new_job_failed()}
							</p>
							<p class="text-xl font-black text-red-600 dark:text-red-400 tabular-nums">{failedCount}</p>
						</div>
					</div>
				</div>
			{/if}

			{#if submittedRows.length === 0}
				<div class="rounded-xl border border-dashed bg-background/20 px-6 py-10 text-center space-y-2">
					<div class="mx-auto flex size-12 items-center justify-center rounded-full bg-muted/40 text-muted-foreground">
						<RocketIcon class="size-6" />
					</div>
					<p class="text-sm font-bold text-foreground">{m.new_job_no_active_extraction_jobs()}</p>
					<p class="text-xs text-muted-foreground max-w-sm mx-auto">
						{m.new_job_no_active_extraction_jobs_body()}
					</p>
				</div>
			{:else}
				<div class="space-y-3 max-h-[30rem] overflow-y-auto pr-1">
					{#each submittedRows as row (row.localId)}
						{@const status = displayJobStatus(row)}
						<div class={rowCardClass(status)}>
							<div class={cn("flex size-11 items-center justify-center rounded-lg border shadow-3xs shrink-0 transition-transform duration-300 hover:scale-105", fileIconClass(row.file))}>
								{#if isImageFile(row.file)}
									<FileImageIcon class="size-5" />
								{:else if isPDFFile(row.file)}
									<FileTextIcon class="size-5" />
								{:else}
									<FileIcon class="size-5" />
								{/if}
							</div>

							<div class="min-w-0 space-y-1">
								<p class="truncate font-bold text-foreground text-sm sm:text-base" title={row.file.name}>{row.file.name}</p>
								<div class="flex flex-wrap items-center gap-x-2 gap-y-1 text-xs">
									<span class="text-muted-foreground font-semibold">{FileDropZone.displaySize(row.file.size)}</span>
									<span class="text-muted-foreground/60">&bull;</span>
									<span class={cn("font-medium", status === "failed" ? "text-destructive" : "text-muted-foreground")}>
										{rowDetail(row)}
									</span>
								</div>
							</div>

							<div class="col-span-2 flex items-center gap-2 justify-self-start sm:col-span-1 sm:justify-self-end w-full sm:w-auto mt-2 sm:mt-0 justify-between sm:justify-end">
								<span class={statusPillClass(status)}>
									{#if status === "completed"}
										<CircleCheckIcon class="size-3.5" />
									{:else if status === "failed"}
										<CircleAlertIcon class="size-3.5" />
									{:else}
										<LoaderIcon class="size-3.5 animate-spin" />
									{/if}
									{statusText(status)}
								</span>

								<div class="flex items-center gap-1 shrink-0">
									{#if status === "completed"}
										{#if row.documentId}
											<Button
												type="button"
												variant="ghost"
												size="icon-sm"
												class="hover:bg-indigo-500/5 hover:text-indigo-500 rounded-md cursor-pointer"
												aria-label={m.new_job_preview_document()}
												onclick={() => openDocumentPreview(row)}
											>
												<EyeIcon class="size-4.5" />
											</Button>
										{:else}
											<Tooltip.Root>
												<Tooltip.Trigger>
													{#snippet child({ props })}
														<span {...props} class="inline-flex">
															<Button
																type="button"
																variant="ghost"
																size="icon-sm"
																disabled
																class="rounded-md"
																aria-label={m.new_job_preview_document()}
															>
																<EyeIcon class="size-4.5" />
															</Button>
														</span>
													{/snippet}
												</Tooltip.Trigger>
												<Tooltip.Content>{m.new_job_preview_unavailable()}</Tooltip.Content>
											</Tooltip.Root>
										{/if}
									{:else if status === "failed"}
										{#if isInsufficientCreditsError(row.errorMessage)}
											<Button
												href="/app/billing"
												variant="outline"
												size="sm"
												class="h-8 gap-1.5 rounded-md"
											>
												<CreditCardIcon class="size-3.5" />
												{m.nav_billing()}
											</Button>
										{/if}
										{#if canRemoveRow(row)}
											<Button
												type="button"
												variant="ghost"
												size="icon-sm"
												class="hover:bg-destructive/5 hover:text-destructive rounded-md cursor-pointer"
												aria-label={m.new_job_remove_failed_job()}
												onclick={() => removeRow(row.localId)}
											>
												<Trash2Icon class="size-4.5" />
											</Button>
										{/if}
									{/if}
								</div>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</section>
		{/if}

		<DocumentPreviewDialog
			bind:open={previewOpen}
			filename={previewFileName}
			markdown={documentQuery.data?.markdown}
			annotationJson={documentQuery.data?.annotation_json}
			isLoading={documentQuery.isLoading}
			error={documentQuery.error}
			onRetry={() => documentQuery.refetch()}
		/>

		{#if showSubmitButton}
			<div class="sticky bottom-0 -mx-4 border-t bg-background/95 px-4 py-4 backdrop-blur-md sm:static sm:mx-0 sm:border-0 sm:bg-transparent sm:p-0 flex justify-center">
				<Button
					type="button"
					class="flex h-12 w-full max-w-sm rounded-xl bg-indigo-600 px-8 text-base font-bold text-white shadow-lg shadow-indigo-500/20 hover:bg-indigo-700 hover:scale-102 hover:shadow-indigo-500/30 disabled:scale-100 disabled:shadow-none transition-all duration-200 cursor-pointer"
					onclick={submitJobs}
					disabled={!canSubmit}
				>
					{#if isSubmitting}
						<LoaderIcon class="mr-2 size-5 animate-spin" />
						{m.new_job_queueing_documents()}
					{:else if isPolling}
						<LoaderIcon class="mr-2 size-5 animate-spin" />
						{m.new_job_extracting_content()}
					{:else}
						<RocketIcon class="mr-2 size-5 transition-transform" />
						{unsubmittedRows.length === 1
							? m.new_job_run_extraction_one({ count: unsubmittedRows.length })
							: m.new_job_run_extraction_other({ count: unsubmittedRows.length })}
					{/if}
				</Button>
			</div>
		{/if}
	</main>
</div>
