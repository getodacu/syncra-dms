<script lang="ts">
	import { onMount } from "svelte";
	import { getLocale } from "$lib/paraglide/runtime.js";
	import { mode } from "mode-watcher";
	import { jsonJoyBuilderConfigForLocale } from "./json-schema-builder-locale";
	import "jsonjoy-builder/styles.css";

	import AlertTriangleIcon from "@lucide/svelte/icons/alert-triangle";
	import RefreshCwIcon from "@lucide/svelte/icons/refresh-cw";

	// Avoid importing jsonjoy-builder's JsonSchema type here because its declarations
	// depend on React types, which this Svelte app does not install.
	type JsonSchemaValue = boolean | Record<string, unknown>;
	type ReactModule = {
		createElement: (type: unknown, props: Record<string, unknown>) => unknown;
	};
	type ReactRoot = {
		render: (node: unknown) => void;
		unmount: () => void;
	};
	type ReactDomClientModule = {
		createRoot: (container: Element) => ReactRoot;
	};
	type JsonJoyBuilderModule = {
		SchemaBuilder: unknown;
	};
	type MonacoEditorModule = {
		editor: {
			setTheme: (themeName: string) => void;
		};
	};
	type MonacoWorkerModule = {
		default: new () => Worker;
	};
	type MonacoGlobalScope = typeof globalThis & {
		monaco?: MonacoEditorModule;
		MonacoEnvironment?: {
			getWorker: (workerId: string, label: string) => Worker;
		};
	};

	type Props = {
		value: JsonSchemaValue;
		onChange: (nextSchema: JsonSchemaValue) => void;
		class?: string;
	};

	let { value, onChange, class: className = "" }: Props = $props();

	let target: HTMLDivElement;
	let react: ReactModule | null = null;
	let root: ReactRoot | null = null;
	let SchemaBuilder: unknown = null;
	let monaco: MonacoEditorModule | null = null;
	let themeObserver: MutationObserver | null = null;
	let pendingThemeFrame: number | null = null;
	let loading = $state(true);
	let errorMessage = $state<string | null>(null);
	let destroyed = false;

	const jsonJoyTheme = $derived(mode.current === "dark" ? "dark" : "light");
	const schemaBuilderClassName = $derived(
		jsonJoyTheme === "dark" ? "h-full min-h-0 dark" : "h-full min-h-0"
	);
	const monacoThemeName = $derived(
		jsonJoyTheme === "dark" ? "appDarkTheme" : "appLightTheme"
	);
	const jsonJoyBuilderConfig = $derived(jsonJoyBuilderConfigForLocale(getLocale()));

	function applyMonacoTheme() {
		if (!monaco) return;

		try {
			monaco.editor.setTheme(monacoThemeName);
		} catch {
			// JSONJoy registers appLightTheme/appDarkTheme when Monaco mounts.
		}
	}

	function scheduleMonacoThemeSync() {
		if (pendingThemeFrame !== null) {
			cancelAnimationFrame(pendingThemeFrame);
		}

		pendingThemeFrame = requestAnimationFrame(() => {
			pendingThemeFrame = null;
			applyMonacoTheme();
		});
	}

	function prepareMonacoThemeSync(
		monacoModule: unknown,
		editorWorkerModule: unknown,
		jsonWorkerModule: unknown
	) {
		const scope = globalThis as MonacoGlobalScope;
		if (editorWorkerModule && jsonWorkerModule && !scope.MonacoEnvironment) {
			const EditorWorker = (editorWorkerModule as MonacoWorkerModule).default;
			const JsonWorker = (jsonWorkerModule as MonacoWorkerModule).default;

			scope.MonacoEnvironment = {
				getWorker: (_workerId: string, label: string) =>
					label === "json" ? new JsonWorker() : new EditorWorker(),
			};
		}

		monaco = monacoModule as MonacoEditorModule;
		scope.monaco = monaco;
		scheduleMonacoThemeSync();
	}

	function observeThemeTarget() {
		if (themeObserver || !target) return;

		themeObserver = new MutationObserver(() => scheduleMonacoThemeSync());
		themeObserver.observe(target, { childList: true, subtree: true });
	}

	function stopThemeSync() {
		themeObserver?.disconnect();
		themeObserver = null;

		if (pendingThemeFrame !== null) {
			cancelAnimationFrame(pendingThemeFrame);
			pendingThemeFrame = null;
		}
	}

	function renderBuilder() {
		if (!react || !root || !SchemaBuilder) return;

		root.render(
			react.createElement(SchemaBuilder, {
				value,
				onChange,
				readOnly: false,
				className: schemaBuilderClassName,
				locale: jsonJoyBuilderConfig.locale,
				messages: jsonJoyBuilderConfig.messages,
			})
		);
		scheduleMonacoThemeSync();
	}

	async function mountBuilder() {
		loading = true;
		errorMessage = null;
		try {
			const [
				reactModule,
				reactDomClientModule,
				jsonJoyBuilderModule,
				monacoEditorModule,
				editorWorkerModule,
				jsonWorkerModule,
			] = await Promise.all([
				// @ts-expect-error React is dynamically loaded for the browser island without @types/react.
				import("react"),
				// @ts-expect-error ReactDOM is dynamically loaded for the browser island without @types/react-dom.
				import("react-dom/client"),
				import("jsonjoy-builder"),
				import("monaco-editor").catch(() => null),
				import("monaco-editor/esm/vs/editor/editor.worker?worker").catch(() => null),
				import("monaco-editor/esm/vs/language/json/json.worker?worker").catch(() => null),
			]);

			if (destroyed) return;

			react = reactModule as ReactModule;
			const reactDomClient = reactDomClientModule as ReactDomClientModule;
			SchemaBuilder = (jsonJoyBuilderModule as JsonJoyBuilderModule).SchemaBuilder;
			if (monacoEditorModule) {
				prepareMonacoThemeSync(monacoEditorModule, editorWorkerModule, jsonWorkerModule);
			}
			
			// Clean up previous root if retrying
			if (root) {
				try {
					root.unmount();
				} catch {
					// safe fallback
				}
			}
			
			root = reactDomClient.createRoot(target);
			observeThemeTarget();
			renderBuilder();
			loading = false;
		} catch (error) {
			if (destroyed) return;

			errorMessage = error instanceof Error ? error.message : "Unable to load schema builder.";
			loading = false;
		}
	}

	onMount(() => {
		destroyed = false;
		void mountBuilder();

		return () => {
			destroyed = true;
			stopThemeSync();
			if (root) {
				try {
					root.unmount();
				} catch {
					// safe fallback
				}
			}
			root = null;
		};
	});

	$effect(() => {
		value;
		jsonJoyTheme;
		jsonJoyBuilderConfig;
		renderBuilder();
	});
</script>

<div
	class={`relative h-full min-h-[480px] w-full overflow-hidden rounded-lg border border-border/80 bg-background ${className}`}
	data-json-schema-builder
>
	<div bind:this={target} class="h-full min-h-0 w-full"></div>

	{#if loading}
		<div class="absolute inset-0 flex flex-col bg-background p-6 z-10 select-none">
			<!-- Editor Header skeleton -->
			<div class="flex items-center justify-between border-b border-border/60 pb-4 mb-6">
				<div class="flex items-center gap-3">
					<div class="w-24 h-6 rounded bg-muted animate-pulse"></div>
					<div class="w-16 h-5 rounded bg-muted/60 animate-pulse"></div>
				</div>
				<div class="w-28 h-8 rounded bg-muted animate-pulse"></div>
			</div>
			
			<!-- Editor layout skeleton representing connection trees -->
			<div class="space-y-6 flex-1 pl-4 border-l border-dashed border-border/80">
				<!-- Root object skeleton -->
				<div class="flex items-center gap-3">
					<div class="w-4 h-4 rounded-full bg-primary/20 animate-pulse shrink-0"></div>
					<div class="w-32 h-8 rounded bg-muted animate-pulse"></div>
					<div class="w-16 h-5 rounded bg-muted/60 animate-pulse"></div>
				</div>
				
				<!-- Children rows skeleton -->
				<div class="space-y-4 pl-8 border-l border-dashed border-border/80">
					<div class="flex items-center gap-4 py-2 border border-border/30 rounded-lg p-3 bg-muted/5">
						<div class="w-4 h-4 rounded bg-muted animate-pulse shrink-0"></div>
						<div class="w-36 h-8 rounded bg-muted animate-pulse"></div>
						<div class="w-24 h-8 rounded bg-muted animate-pulse"></div>
						<div class="w-16 h-5 rounded-full bg-muted/60 animate-pulse"></div>
						<div class="ml-auto w-8 h-8 rounded bg-muted/40 animate-pulse shrink-0"></div>
					</div>
					
					<div class="flex items-center gap-4 py-2 border border-border/30 rounded-lg p-3 bg-muted/5">
						<div class="w-4 h-4 rounded bg-muted shrink-0 animate-pulse"></div>
						<div class="w-40 h-8 rounded bg-muted animate-pulse"></div>
						<div class="w-24 h-8 rounded bg-muted animate-pulse"></div>
						<div class="w-16 h-5 rounded-full bg-muted/60 animate-pulse"></div>
						<div class="ml-auto w-8 h-8 rounded bg-muted/40 animate-pulse shrink-0"></div>
					</div>
					
					<div class="flex items-center gap-4 py-2 border border-border/30 rounded-lg p-3 bg-muted/5">
						<div class="w-4 h-4 rounded bg-muted shrink-0 animate-pulse"></div>
						<div class="w-28 h-8 rounded bg-muted animate-pulse"></div>
						<div class="w-24 h-8 rounded bg-muted animate-pulse"></div>
						<div class="w-16 h-5 rounded-full bg-muted/60 animate-pulse"></div>
						<div class="ml-auto w-8 h-8 rounded bg-muted/40 animate-pulse shrink-0"></div>
					</div>
				</div>
				
				<!-- Field creation action skeleton -->
				<div class="pl-8 flex items-center gap-2">
					<div class="w-5 h-5 rounded bg-muted animate-pulse"></div>
					<div class="w-24 h-5 rounded bg-muted animate-pulse"></div>
				</div>
			</div>
			
			<!-- Overlay with active feedback information -->
			<div class="absolute inset-0 flex flex-col items-center justify-center bg-background/50 backdrop-blur-[1px]">
				<div class="flex flex-col items-center gap-3 p-5 rounded-xl border border-border bg-card shadow-md max-w-[280px] text-center">
					<div class="size-8 rounded-full bg-primary/10 flex items-center justify-center text-primary">
						<span class="animate-spin size-4 border-2 border-primary border-t-transparent rounded-full"></span>
					</div>
					<div class="space-y-1">
						<p class="text-sm font-semibold text-foreground">Interactive Builder</p>
						<p class="text-xs text-muted-foreground">Initializing visual schema design canvas...</p>
					</div>
				</div>
			</div>
		</div>
	{/if}

	{#if errorMessage}
		<div class="absolute inset-0 grid min-h-[480px] place-items-center bg-destructive/5 backdrop-blur-[1px] p-6 z-20">
			<div class="max-w-md w-full rounded-xl border border-destructive/25 bg-card p-6 text-center shadow-lg animate-shake">
				<div class="mx-auto size-12 rounded-full bg-destructive/10 flex items-center justify-center text-destructive mb-4">
					<AlertTriangleIcon class="size-6" />
				</div>
				<h4 class="text-base font-bold text-foreground mb-1">Editor Failed to Load</h4>
				<p class="text-xs text-muted-foreground mb-4 leading-relaxed">
					{errorMessage}
				</p>
				<button
					type="button"
					onclick={mountBuilder}
					class="inline-flex items-center justify-center rounded-md bg-destructive px-4 py-2 text-xs font-semibold text-destructive-foreground shadow hover:bg-destructive/90 transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring cursor-pointer"
				>
					<RefreshCwIcon class="size-3.5 mr-1.5" />
					Retry Loading Environment
				</button>
			</div>
		</div>
	{/if}
</div>
