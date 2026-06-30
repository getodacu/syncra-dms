<script lang="ts">
	import CopyPlusIcon from "@lucide/svelte/icons/copy-plus";
	import PencilIcon from "@lucide/svelte/icons/pencil";
	import RocketIcon from "@lucide/svelte/icons/rocket";
	import Trash2Icon from "@lucide/svelte/icons/trash-2";
	import { Button } from "$lib/components/ui/button/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import { buildNewJobPath } from "../new-job/schema-query";
	import type { SchemaListItemResponse } from "./api";

	let {
		schema,
		onClone,
		onDelete,
		clonePending = false,
		deletePending = false,
	}: {
		schema: SchemaListItemResponse;
		onClone: (schema: SchemaListItemResponse) => void;
		onDelete: (schema: SchemaListItemResponse) => void;
		clonePending?: boolean;
		deletePending?: boolean;
	} = $props();
</script>

<div class="flex items-center justify-end gap-1">
	<Button
		href={`/app/schemas/edit/${schema.id}`}
		variant="ghost"
		size="icon-sm"
		class="text-muted-foreground hover:text-foreground transition-colors"
		aria-label={m.schemas_edit_aria({ name: schema.name })}
	>
		<PencilIcon class="size-4" aria-hidden="true" />
	</Button>

	<Button
		href={buildNewJobPath(schema.id)}
		variant="ghost"
		size="icon-sm"
		class="text-muted-foreground hover:text-foreground transition-colors"
		aria-label={m.schemas_create_job_with({ name: schema.name })}
		title={m.schemas_create_job_with({ name: schema.name })}
	>
		<RocketIcon class="size-4" aria-hidden="true" />
	</Button>

	<Button
		type="button"
		variant="ghost"
		size="icon-sm"
		class="text-muted-foreground hover:text-foreground transition-colors"
		disabled={clonePending}
		aria-label={m.schemas_clone_aria({ name: schema.name })}
		onclick={() => onClone(schema)}
	>
		<CopyPlusIcon class="size-4" aria-hidden="true" />
	</Button>

	<Button
		type="button"
		variant="ghost"
		size="icon-sm"
		class="text-muted-foreground hover:bg-destructive/10 hover:text-destructive dark:hover:bg-destructive/20 transition-all"
		disabled={deletePending}
		aria-label={m.schemas_delete_aria({ name: schema.name })}
		onclick={() => onDelete(schema)}
	>
		<Trash2Icon class="size-4" aria-hidden="true" />
	</Button>
</div>
