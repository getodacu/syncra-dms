# Syncra DMS Go Server

Lean Go API scaffold for Syncra DMS.

## Setup

```sh
cp .env.example .env
rtk go test ./...
rtk go run ./cmd/syncra api
rtk go run ./cmd/syncra api --port 8090
```

`DSN` must target `syncra_dms`. `DSN_DEV` must target `syncra_dms_dev`. `ATLAS_DATABASE_URL` must target `syncra_dms`. `ATLAS_DEV_DATABASE_URL` must target a separate empty scratch database.

## Endpoints

- `GET /healthz`
- `GET /readyz`
- `GET /version`
- `GET /swagger/index.html`
- `GET /swagger/doc.json`
- `/api/auth/*` for trusted SvelteKit-to-Go authentication requests

## Atlas

```sh
rtk go run ./cmd/syncra migrate
rtk atlas migrate validate --dir file://migrations
```

Auth models are included in `internal/database.ApplicationModels()` for Atlas migration output. Add future domain models there as feature plans introduce persistent entities.

## Swagger

```sh
rtk go run ./cmd/syncra swagger
```

This generates and validates `docs/swagger.json` using the installed `swagger` binary.
The API serves the embedded spec at `/swagger/doc.json` and mounts Swagger UI at `/swagger/index.html`.
