# Syncra DMS Agent Conventions

## Shell

Always prefix shell commands with `rtk`.

Examples:

```sh
rtk git status
rtk go test ./...
rtk pnpm check
```

Use `rtk proxy <cmd>` only when the command needs raw output or `rtk` does not support a shell feature.

## Repository Layout

- `go-server/` contains the lean Go API scaffold.
- `frontend/` contains the SvelteKit application scaffold.
- `example/` is read-only reference material. Do not copy product feature modules from it unless a plan explicitly calls for that.
- `docs/superpowers/specs/` contains approved design specs.

## Work Style

- Keep the initial scaffold free of OCR, billing, email, PDF, and auth feature wiring.
- Organization Units are the first expected domain feature, but they are not part of the environment scaffold.
- Add behavior tests before production code for new backend or frontend behavior.
- Keep private values in `.env` files and server-only modules. Do not expose private env vars to browser code.
