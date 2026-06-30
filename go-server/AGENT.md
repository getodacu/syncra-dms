# Go Server Agent Conventions

These rules apply inside `go-server/` in addition to the repository conventions.

## Stack

- Go module: `ai.ro/syncra/dms`.
- HTTP framework: Gin.
- Database: PostgreSQL through GORM and pgx-related dependencies.
- Configuration: environment variables loaded from `go-server/.env`.
- Schema management: Atlas migrations in `go-server/migrations`.

## Database Rules

- Local development uses the `syncra_dev` database.
- `DSN` and `DSN_DEV` must both target `syncra_dev`.
- `ATLAS_DATABASE_URL` must target the same `syncra_dev` database in URL form.
- `ATLAS_DEV_DATABASE_URL` must target a separate empty scratch database, not `syncra_dev`.
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
