# Frontend Agent Conventions

These rules apply inside `frontend/` in addition to the repository conventions.

## Stack

- SvelteKit with Svelte 5 and TypeScript.
- Package manager: `pnpm`.
- Styling: Tailwind CSS with local shadcn-svelte-style primitives under `src/lib/components/ui`.
- Localization: Paraglide/Inlang with source messages in `messages/`.
- Deployment adapter: `@sveltejs/adapter-node`.

## Environment

- Copy `.env.example` to `.env` for local development.
- `SYNCRA_API_BASE_URL` defaults to the local Go API at `http://localhost:8080`.
- Private values must stay in server-only modules, server routes, server load functions, or hooks.
- Browser code must not import private environment modules.

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

## Svelte Conventions

- Use Svelte 5 runes and direct event attributes for new code.
- Prefer `$derived` for computed values.
- Keep route-specific components near their route and shared primitives under `src/lib`.
- Reuse existing UI primitives before adding new ones.
- Keep operational screens dense, scannable, and work-focused.
