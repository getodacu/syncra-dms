# Syncra DMS

Syncra DMS is a lean work environment scaffold for a document management system. It follows the technology and development methods from `example/server` and `example/frontend` while omitting example product features such as OCR, billing, email, PDF generation, and authentication.

## Layout

```text
.
├── go-server/              # Go API scaffold
├── frontend/               # SvelteKit scaffold
├── docs/superpowers/specs/ # approved design specs
└── example/                # read-only reference implementation
```

## Backend

The backend uses Go, Gin, PostgreSQL, GORM, and Atlas.

```sh
cd go-server
cp .env.example .env
go test ./...
go run ./cmd/api
```

Required local databases:

- `syncra_dev` for application development and tests.
- a separate empty Atlas scratch database, for example `syncra_atlas`.

Operational endpoints:

- `GET /healthz`
- `GET /readyz`
- `GET /version`

## Frontend

The frontend uses SvelteKit, Svelte 5, TypeScript, pnpm, Tailwind CSS, adapter-node, and Paraglide/Inlang.

```sh
cd frontend
cp .env.example .env
pnpm install
pnpm dev
```

Quality commands:

```sh
pnpm check
pnpm test
pnpm build
```

## Development Notes

- Prefix shell commands with `rtk` in this workspace.
- Do not add Docker Compose for the initial scaffold.
- Keep private environment values server-side.
- Add Organization Units in a follow-up feature plan, not in this environment setup pass.
