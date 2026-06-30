# Syncra DMS Go Server

Lean Go API scaffold for Syncra DMS.

## Setup

```sh
cp .env.example .env
go test ./...
go run ./cmd/api
```

`DSN` must target `syncra_dms`. `DSN_DEV` must target `syncra_dms_dev`. `ATLAS_DATABASE_URL` must target `syncra_dms`. `ATLAS_DEV_DATABASE_URL` must target a separate empty scratch database.

## Endpoints

- `GET /healthz`
- `GET /readyz`
- `GET /version`

## Atlas

```sh
go run ./cmd/atlas-loader
atlas migrate validate --dir file://migrations
```

The initial scaffold has no domain models. Add models to `internal/database.ApplicationModels()` as feature plans introduce persistent entities.
