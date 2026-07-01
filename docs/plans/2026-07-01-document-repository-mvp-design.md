# Document Repository MVP Design

## Summary

Implement the first document management slice as a folder-first repository MVP. The slice includes folder hierarchy, document metadata, local file upload and download, organization-unit-scoped RBAC enforcement, and soft delete/archive behavior.

This design intentionally defers categories, collections, upload sessions, versions, preview, explicit sharing, OCR, classification, signing, and workflow integration. The goal is to create the repository spine before adding higher-level document management features.

## Goals

- Let authorized users browse folders scoped to an organization unit.
- Let authorized users create, rename, move, and archive folders.
- Let authorized users upload one file per request into a selected folder.
- Store document metadata in PostgreSQL and file bytes under a configurable local storage root.
- Let authorized users view metadata and download files.
- Enforce access in the Go backend with existing RBAC and organization-unit scope.
- Keep records and file bytes recoverable by using soft delete/archive behavior.

## Non-Goals

- Categories, collections, document sharing, upload sessions, version history, preview, OCR, classification, signing, workflow, and audit timelines.
- External file storage such as S3-compatible buckets.
- Hard deletion of document records or file bytes.
- Browser exposure of private storage paths or server-only configuration.

## Architecture

The Go API owns folder hierarchy, document metadata, local file persistence, permission checks, and download streaming. The SvelteKit frontend uses the existing trusted server proxy pattern to call the Go API with the internal token and current session cookie.

Repository records are scoped by organization unit. Folders and documents store `organization_unit_id`, and backend reads/mutations require the matching `document.*` permission for that unit through the existing RBAC resolver. This aligns document access with the existing user, role, group, and organization-unit foundation without adding explicit sharing in the first slice.

Local file bytes are stored under `SYNCRA_DOCUMENT_STORAGE_ROOT`. The database stores metadata plus an internal `storage_key`; it must not expose absolute filesystem paths to browser code or public API responses. Downloads are streamed by the backend after permission checks.

Delete behavior is soft delete/archive. Active repository views exclude archived folders and documents. File bytes remain on disk for future restore, audit, and retention behavior.

## Data Model

Add `document_folders`:

- `id uuid primary key`
- `parent_id uuid null`, self-referencing `document_folders.id`
- `organization_unit_id uuid not null`, referencing `organization_units.id`
- `name varchar(160) not null`
- `description text null`
- `created_by_user_id uuid not null`, referencing `user.id`
- `updated_by_user_id uuid null`, referencing `user.id`
- `deleted_at timestamptz null`
- `created_at timestamptz not null`
- `updated_at timestamptz not null`

Add `documents`:

- `id uuid primary key`
- `folder_id uuid not null`, referencing `document_folders.id`
- `organization_unit_id uuid not null`, referencing `organization_units.id`
- `original_file_name varchar(255) not null`
- `display_name varchar(255) not null`
- `mime_type varchar(255) not null`
- `extension varchar(32) null`
- `size_bytes bigint not null`
- `sha256_hash char(64) not null`
- `storage_key text not null`
- `created_by_user_id uuid not null`, referencing `user.id`
- `updated_by_user_id uuid null`, referencing `user.id`
- `deleted_at timestamptz null`
- `created_at timestamptz not null`
- `updated_at timestamptz not null`

Active folder names are unique per parent and organization unit. Root folders use `parent_id = null`. Folder moves must reject cycles and cross-organization-unit moves.

Document display names do not need to be unique in a folder for MVP. Duplicate upload protection uses `sha256_hash` plus `folder_id` among active documents.

## Backend API

Expose internal Go API endpoints:

- `GET /api/document-folders/tree?organizationUnitId=...`
- `POST /api/document-folders`
- `PATCH /api/document-folders/:id`
- `PATCH /api/document-folders/:id/parent`
- `POST /api/document-folders/:id/archive`
- `GET /api/document-folders/:id/contents`
- `POST /api/documents/upload`
- `GET /api/documents/:id`
- `GET /api/documents/:id/download`
- `PATCH /api/documents/:id`
- `POST /api/documents/:id/archive`

Permission mapping:

- `document.view`: tree, folder contents, and document metadata.
- `document.create`: folder creation and uploads.
- `document.update`: folder/document rename and folder moves.
- `document.delete`: folder/document archive.
- `document.download`: document download.

All endpoints require the trusted internal request header and an authenticated session. The backend enforces permissions even when the frontend hides actions.

## Frontend Behavior

Add `/app/documents` as a dense repository workspace in the protected app shell. The route shows an organization-unit selector when the user can access multiple units, a folder tree on the left, and selected folder contents on the right.

Folder contents combine child folders and documents in one sortable list. The UI exposes compact actions for creating folders, renaming, moving, archiving, uploading, and downloading according to page permission flags.

Uploads follow the one-file-per-request model from the specification. The browser may allow selecting multiple files, but each file is submitted independently to a SvelteKit endpoint and forwarded as a single multipart upload to Gin. Each file receives separate progress, success, and error state. Upload sessions are deferred.

Downloads stay server-mediated. The browser calls a SvelteKit route, SvelteKit forwards to Gin with session and internal headers, and Gin streams the stored file after `document.download` authorization. Responses use safe content headers and never expose the storage root.

The app navigation shows Documents to users with any relevant `document.*` permission. Users without permission do not see the nav item and receive a forbidden state if they reach the route manually.

## Error Handling

The backend returns public JSON errors consistent with existing helpers:

- `400` for invalid folder IDs, missing names, oversized uploads, unsupported file types, invalid multipart bodies, and invalid move targets.
- `401` for missing or invalid sessions.
- `403` for missing `document.*` permission in the target organization unit.
- `404` for missing or archived folders/documents.
- `409` for duplicate active folder names, invalid move cycles, cross-unit moves, or duplicate upload hash in the same folder.
- `500` for unexpected database/storage failures without exposing absolute storage paths.

The frontend maps these into page-level load failures, inline form errors, upload-row errors, and forbidden/read-only states.

## Testing

Backend tests cover:

- folder create, tree, update, move, archive, and active-view exclusion;
- duplicate folder names and folder move cycle rejection;
- document upload metadata, SHA-256 hash, storage key, and duplicate hash handling;
- download permission checks and safe missing-file handling;
- soft delete exclusion for folders and documents;
- model registration for Atlas schema generation.

Frontend tests cover:

- SvelteKit server proxy routes and response validation;
- permission-derived page flags and navigation gating;
- folder tree and content selection behavior;
- upload queue behavior with per-file success and error state;
- action gating and empty/error states.

Verification commands:

- `rtk go test ./...`
- `rtk atlas migrate validate --dir file://migrations`
- `rtk pnpm test`
- `rtk pnpm check`
- `rtk pnpm build`

## Implementation Constraints

- Prefix shell commands with `rtk`.
- Keep private configuration in `.env` files and server-only modules.
- Do not expose `SYNCRA_DOCUMENT_STORAGE_ROOT`, internal API tokens, or absolute storage paths to browser code.
- Preserve existing worktree changes and do not revert unrelated files.
- Keep `example/` read-only and use it only for patterns.
- Update Swagger docs when the Go API surface changes.
- Keep the implementation focused on the repository MVP and avoid unrelated refactors.
