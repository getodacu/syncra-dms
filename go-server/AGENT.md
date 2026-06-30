# Go Server Agent Conventions

These rules apply inside `go-server/` in addition to the repository conventions.

## Stack

- Go module: `ai.ro/syncra/dms`.
- HTTP framework: Gin.
- Database: PostgreSQL through GORM and pgx-related dependencies.
- Configuration: environment variables loaded from `go-server/.env`.
- Schema management: Atlas migrations in `go-server/migrations`.

## Database Rules

- Local runtime uses the `syncra_dms` database.
- Postgres-backed tests use the `syncra_dms_dev` database.
- `DSN` must target `syncra_dms`.
- `DSN_DEV` must target `syncra_dms_dev`.
- `ATLAS_DATABASE_URL` must target the same `syncra_dms` database in URL form.
- `ATLAS_DEV_DATABASE_URL` must target a separate empty scratch database, not `syncra_dms` or `syncra_dms_dev`.
- Do not run migration apply commands until DSN targets have been checked.

## Commands

```sh
go test ./...
go run ./cmd/api
go run ./cmd/atlas-loader
atlas migrate validate --dir file://migrations
```

Run commands from `go-server/`.

## Code Conventions

- Keep HTTP routing in `internal/api`.
- Keep application startup and runtime wiring in `internal/app`.
- Keep database opening and model-list wiring in `internal/database`.
- Add domain packages only when the corresponding feature plan is approved.
- Preserve the operational endpoints: `GET /healthz`, `GET /readyz`, and `GET /version`.
