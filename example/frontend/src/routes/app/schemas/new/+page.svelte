<script lang="ts">
	import { createMutation, useQueryClient } from "@tanstack/svelte-query";
	import { goto } from "$app/navigation";
	import { toast } from "svelte-sonner";
	import {
		PERSONAL_SCHEMA_OPTIONS_QUERY_KEY,
		upsertPersonalSchemaOption,
		type PersonalSchemaOption,
	} from "$lib/client/schemas";
	import { m } from "$lib/paraglide/messages.js";
	import SchemaEditor from "../schema-editor.svelte";
	import { createSchema, type SchemaEditorSubmitInput, type SchemaResponse } from "../api";
	import { schemaHasFields } from "./schema-validation";

	type SchemaMutationVariables = {
		input: SchemaEditorSubmitInput;
		feedbackVersion: number;
	};

	let saved = $state<SchemaResponse | null>(null);
	let displayedServerError = $state<string | null>(null);
	let feedbackVersion = 0;

	const queryClient = useQueryClient();
	const mutation = createMutation<SchemaResponse, Error, SchemaMutationVariables>(() => ({
		mutationFn: ({ input }) => createSchema(fetch, input),
		onSuccess: (result, variables) => {
			queryClient.setQueryData<PersonalSchemaOption[]>(
				PERSONAL_SCHEMA_OPTIONS_QUERY_KEY,
				(current) => upsertPersonalSchemaOption(current, result)
			);
			void queryClient.invalidateQueries({ queryKey: ["schemas"] });
			if (variables.feedbackVersion !== feedbackVersion) return;

			saved = result;
			toast.success(m.schemas_saved_success_with_id({ name: result.name, id: result.id }));
			void goto("/app/schemas");
		},
		onError: (error, variables) => {
			if (variables.feedbackVersion !== feedbackVersion) return;

			displayedServerError = error.message;
		},
	}));
	const serverError = $derived(displayedServerError);

	function clearFeedback() {
		feedbackVersion += 1;
		saved = null;
		displayedServerError = null;
		return feedbackVersion;
	}

	function submit(input: SchemaEditorSubmitInput) {
		const submissionFeedbackVersion = clearFeedback();
		if (!schemaHasFields(input.schema)) {
			displayedServerError = m.schemas_empty_schema_error();
			return;
		}

		mutation.mutate({ input, feedbackVersion: submissionFeedbackVersion });
	}
</script>

<SchemaEditor
	title={m.schemas_new_title()}
	descriptionText={m.schemas_new_description()}
	submitLabel={m.schemas_save_schema()}
	pending={mutation.isPending}
	{serverError}
	onSubmit={submit}
	onDirty={clearFeedback}
/>
