# Syncra DMS Work Environment Scaffold

## Summary

Set up a lean Syncra DMS workspace using the examples' technology and work methods, without copying OCR, billing, authentication, or other product feature code. The backend lives in `go-server/`, the frontend lives in `frontend/`, and local PostgreSQL is expected to exist outside this repository.

## Architecture

The backend is a Go service using the same infrastructure style as the example server: Gin for HTTP routing, environment-driven configuration, structured application wiring, PostgreSQL through GORM/pgx, Atlas migrations, and package-local tests. The Go module identity is `ai.ro/syncra/dms`.

The frontend is a SvelteKit application using the same tooling style as the example frontend: Svelte 5, TypeScript, pnpm, Tailwind CSS, adapter-node, Paraglide/Inlang, and local shadcn-style UI primitives.

Organization Units are intentionally not implemented in this scaffold. The scaffold prepares the backend and frontend boundaries for Organization Units to be added as the first domain feature in a later implementation plan.

## Backend Interface

The initial backend exposes operational endpoints only:

- `GET /healthz` returns process liveness and does not depend on the database.
- `GET /readyz` reports whether required runtime configuration and PostgreSQL connectivity are valid.
- `GET /version` returns application metadata including app name, module path, and version.

Backend configuration follows the example conventions:

- `SERVER_HOST_PORT`
- `DEBUG`
- `DSN`
- `DSN_DEV`
- `ATLAS_DATABASE_URL`
- `ATLAS_DEV_DATABASE_URL`

`DSN` and `DSN_DEV` must target the local `syncra_dev` database. `ATLAS_DEV_DATABASE_URL` must target a separate scratch database.

## Frontend Interface

The initial frontend exposes a minimal DMS shell and a server-loaded status surface that calls the Go backend through private server-side configuration. Browser code must not import private environment values.

Frontend configuration follows the example conventions:

- `SYNCRA_API_BASE_URL`
- `SYNCRA_APP_ORIGIN`
- `SYNCRA_INTERNAL_API_TOKEN`

`SYNCRA_API_BASE_URL` defaults to `http://localhost:8080`.

## Developer Workflow

Commands should be run from their service directories and prefixed with `rtk` in this workspace.

Backend commands:

- `go test ./...`
- `go run ./cmd/api`
- `go run ./cmd/atlas-loader`
- `atlas migrate validate --dir file://migrations`

Frontend commands:

- `pnpm install`
- `pnpm dev`
- `pnpm check`
- `pnpm test`
- `pnpm build`

The repository does not add Docker Compose. Developers provide PostgreSQL locally and configure DSNs in `.env` files copied from `.env.example`.

## Verification

The scaffold is complete when:

- backend tests pass with `go test ./...`;
- frontend checks/tests/build pass with `pnpm check`, `pnpm test`, and `pnpm build`;
- the backend can serve liveness, readiness, and version endpoints;
- the frontend server-side status route can read the backend version/readiness endpoints;
- docs describe setup, environment, and command conventions.
