# Frontend Agent Conventions

These rules apply inside `frontend/` in addition to the root conventions.

## Stack

- SvelteKit with Svelte 5 and TypeScript.
- Package manager: `pnpm`.
- Styling: Tailwind CSS with local shadcn-svelte-style components in `src/lib/components/ui`.
- UI primitives: Bits UI and local wrappers.
- Icons: prefer existing lucide-svelte or Tabler icon usage based on nearby components.
- Data/client helpers live under `src/lib/client` when shared across routes.

## Environment

- Copy `.env.example` to `.env` for local development.
- In worktrees, copy `frontend/.env` from the source checkout before running frontend commands.
- `SYNCRA_API_BASE_URL` defaults to the local API. Keep server-only values private unless they intentionally use SvelteKit's public prefix.
- Do not expose server secrets to browser code. Use `.server.ts` files, server routes, hooks, or `$lib/server` patterns for private logic.

## Commands

```sh
pnpm install
pnpm dev
pnpm check
pnpm test
pnpm build
pnpm preview
```

Run commands from `frontend/`.

## Svelte and SvelteKit Conventions

- Use Svelte 5 patterns for new code: `$state`, `$derived`, `$props`, snippets, and direct event attributes such as `onclick`.
- Prefer `$derived` for computed values. Use `$effect` only for side effects or integration with external systems.
- Treat props as changing values; derive dependent state rather than freezing initial props accidentally.
- Use keyed `{#each}` blocks with stable object identifiers, not array indexes.
- Keep route-specific components near their route; put reusable components/utilities in `src/lib`.
- Use `$lib` imports for shared frontend code.
- Keep private server logic out of browser-reachable modules.
- Sanitize or safely render rich OCR/Markdown/HTML content; do not introduce raw `{@html}` without confirming sanitization.

## UI Conventions

- Reuse existing `src/lib/components/ui` components before adding new primitives.
- Match nearby component patterns for spacing, typography, states, and icon family.
- Use icons for compact actions where an established icon exists.
- Keep operational app screens dense, scannable, and task-focused rather than marketing-style.
- Ensure loading, empty, error, disabled, and destructive-action states are handled for user-facing flows.
- Keep text within controls responsive and non-overlapping on mobile and desktop.

## Tests

- Add or update Vitest tests for client helpers, route server handlers, state utilities, and behavior that can be isolated from the DOM.
- Run `pnpm check` after TypeScript/Svelte changes.
- Run `pnpm test` for logic or route-handler changes.
- Run `pnpm build` for routing, environment, SvelteKit config, or integration-sensitive changes.
- When writing or changing Svelte components, verify with the official Svelte MCP/autofixer when available.
