# Syncra DMS Frontend

Lean SvelteKit scaffold for Syncra DMS.

## Setup

```sh
cp .env.example .env
pnpm install
pnpm dev
```

`SYNCRA_API_BASE_URL` points server-side load functions at the Go API. Keep it private.

`BODY_SIZE_LIMIT` configures adapter-node request body handling. Keep it at least `26M` so document uploads can clear the Go API's 25 MiB limit plus multipart overhead.

## Commands

```sh
pnpm check
pnpm test
pnpm build
pnpm preview
```

## Structure

- `src/routes/+page.server.ts` loads backend operational status.
- `src/routes/+page.svelte` renders the environment status screen.
- `src/lib/server/` contains server-only API clients.
- `src/lib/components/ui/` contains local UI primitives.
