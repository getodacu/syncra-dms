<script lang="ts">
	import LoaderIcon from "@lucide/svelte/icons/loader-circle";
	import RefreshCwIcon from "@lucide/svelte/icons/refresh-cw";
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { createMutation, createQuery, useQueryClient } from "@tanstack/svelte-query";
	import { toast } from "svelte-sonner";

	import {
		PERSONAL_SCHEMA_OPTIONS_QUERY_KEY,
		upsertPersonalSchemaOption,
		type PersonalSchemaOption,
	} from "$lib/client/schemas";
	import * as Alert from "$lib/components/ui/alert/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import SchemaEditor from "../../schema-editor.svelte";
	import {
		cloneSchema,
		getSchema,
		isSchemaNotFoundError,
		shouldRetrySchemaQuery,
		updateSchema,
		type SchemaEditorSubmitInput,
		type SchemaResponse,
	} from "../../api";

	type SchemaMutationVariables = {
		input: SchemaEditorSubmitInput;
		feedbackVersion: number;
		schemaId: string;
	};

	const queryClient = useQueryClient();
	const schemaId = $derived(page.params.id ?? "");
	let saved = $state<SchemaResponse | null>(null);
	let displayedServerError = $state<string | null>(null);
	let dirty = $state(false);
	let pendingSchemaId = $state<string | null>(null);
	let activeEditorKey = $state("");
	let feedbackVersion = 0;

	const schemaQuery = createQuery<SchemaResponse, Error>(() => ({
		queryKey: ["schema", schemaId],
		queryFn: () => getSchema(fetch, schemaId),
		enabled: Boolean(schemaId),
		retry: shouldRetrySchemaQuery,
	}));
	const mutation = createMutation<SchemaResponse, Error, SchemaMutationVariables>(() => ({
		mutationFn: ({ input, schemaId: submittedSchemaId }) =>
			updateSchema(fetch, submittedSchemaId, input),
		onSuccess: (result, variables) => {
			queryClient.setQueryData<SchemaResponse>(["schema", variables.schemaId], result);
			queryClient.setQueryData<PersonalSchemaOption[]>(
				PERSONAL_SCHEMA_OPTIONS_QUERY_KEY,
				(current) => upsertPersonalSchemaOption(current, result)
			);
			void queryClient.invalidateQueries({ queryKey: ["schemas"] });
			void queryClient.invalidateQueries({ queryKey: ["schema", variables.schemaId] });
			if (variables.feedbackVersion !== feedbackVersion || variables.schemaId !== schemaId) return;

			saved = result;
			dirty = false;
			toast.success(m.schemas_saved_success({ name: result.name }));
		},
		onError: (error, variables) => {
			if (variables.feedbackVersion !== feedbackVersion || variables.schemaId !== schemaId) return;

			displayedServerError = error.message;
		},
		onSettled: (_result, _error, variables) => {
			if (variables.schemaId === pendingSchemaId) pendingSchemaId = null;
		},
	}));
	const cloneMutation = createMutation<SchemaResponse, Error, { schema: SchemaResponse }>(() => ({
		mutationKey: ["schemas", "clone"],
		mutationFn: ({ schema }) => cloneSchema(fetch, schema),
		onSuccess: (result) => {
			queryClient.setQueryData<SchemaResponse>(["schema", result.id], result);
			queryClient.setQueryData<PersonalSchemaOption[]>(
				PERSONAL_SCHEMA_OPTIONS_QUERY_KEY,
				(current) => upsertPersonalSchemaOption(current, result)
			);
			void queryClient.invalidateQueries({ queryKey: ["schemas"] });
			void goto(`/app/schemas/edit/${result.id}`);
		},
		onError: (error) => {
			displayedServerError = error.message;
		},
	}));
	const serverError = $derived(displayedServerError);
	const successMessage = $derived(
		saved ? m.schemas_saved_feedback({ name: saved.name, id: saved.id }) : null
	);
	const isCurrentSchemaSavePending = $derived(mutation.isPending && pendingSchemaId === schemaId);
	const schemaNotFound = $derived(isSchemaNotFoundError(schemaQuery.error));
	const cleanEditorKey = $derived(
		schemaQuery.data ? `${schemaQuery.data.id}:${schemaQuery.data.updated_at}` : ""
	);
	const editorKey = $derived(activeEditorKey || cleanEditorKey);

	function clearFeedback() {
		feedbackVersion += 1;
		saved = null;
		displayedServerError = null;
		return feedbackVersion;
	}

	function markDirty() {
		dirty = true;
		clearFeedback();
	}

	function submit(input: SchemaEditorSubmitInput) {
		const submissionFeedbackVersion = clearFeedback();
		pendingSchemaId = schemaId;
		mutation.mutate({ input, feedbackVersion: submissionFeedbackVersion, schemaId });
	}

	function cloneCurrentSchema() {
		if (!schemaQuery.data || cloneMutation.isPending) return;

		clearFeedback();
		cloneMutation.mutate({ schema: schemaQuery.data });
	}

	$effect(() => {
		schemaId;
		dirty = false;
		pendingSchemaId = null;
		activeEditorKey = "";
		clearFeedback();
	});

	$effect(() => {
		if (!dirty) activeEditorKey = cleanEditorKey;
	});
</script>

{#if schemaQuery.isLoading}
	<div class="flex flex-1 items-center justify-center gap-2 px-4 py-10 text-sm text-muted-foreground">
		<LoaderIcon class="size-4 animate-spin" />
		{m.schemas_loading_schema()}
	</div>
{:else if schemaQuery.isError}
	<div class="flex flex-1 items-center justify-center px-4 py-10">
		{#if schemaNotFound}
			<Alert.Root class="max-w-xl">
				<Alert.Title>{m.schemas_not_found_title()}</Alert.Title>
				<Alert.Description>
					{m.schemas_not_found_body()}
					<div class="mt-4">
						<Button href="/app/schemas" variant="outline" size="sm">{m.schemas_view_schemas()}</Button>
					</div>
				</Alert.Description>
			</Alert.Root>
		{:else}
			<Alert.Root variant="destructive" class="max-w-xl">
				<Alert.Title>{m.schemas_could_not_load()}</Alert.Title>
				<Alert.Description>{schemaQuery.error.message}</Alert.Description>
				<Alert.Action>
					<Button type="button" variant="ghost" size="icon" onclick={() => schemaQuery.refetch()}>
						<RefreshCwIcon class="size-4" />
						<span class="sr-only">{m.common_retry()}</span>
					</Button>
				</Alert.Action>
			</Alert.Root>
		{/if}
	</div>
{:else if schemaQuery.data}
	{#key editorKey}
		<SchemaEditor
			title={m.schemas_edit_title()}
			descriptionText={m.schemas_edit_description()}
			submitLabel={m.schemas_save_changes()}
			pending={isCurrentSchemaSavePending}
			initial={schemaQuery.data}
			schemaId={schemaQuery.data.id}
			{serverError}
			{successMessage}
			clonePending={cloneMutation.isPending}
			onSubmit={submit}
			onClone={cloneCurrentSchema}
			onDirty={markDirty}
		/>
	{/key}
{/if}
