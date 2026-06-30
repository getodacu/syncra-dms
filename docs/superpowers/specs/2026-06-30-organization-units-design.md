# Organization Units Design

## Summary

Implement Organization Units as the first Syncra DMS domain feature across the Go API and SvelteKit frontend. The first slice delivers a shared company hierarchy that admins manage and regular authenticated users can view.

The feature follows the current scaffold boundaries: Go domain behavior lives in `go-server/`, frontend behavior lives in `frontend/`, and `example/` remains read-only reference material for patterns only. This slice does not add OCR, billing, email, PDF, or document workflow wiring.

## Goals

- Provide a global Organization Unit hierarchy for the DMS company structure.
- Let admins create, edit, move, and archive units.
- Let regular authenticated users browse active units and view unit details.
- Keep unit records recoverable by archiving instead of hard deleting.
- Match the example frontend's dense, operational app style while staying within the lean current scaffold.
- Add behavior tests before production behavior.

## Non-Goals

- Per-unit membership, managers, permissions, or role assignment.
- Linking documents, workflows, files, billing records, or OCR jobs to units.
- Hard delete, restore workflows, audit timelines, or import/export.
- A separate admin portal shell.

## Architecture

The Go API owns Organization Unit persistence, validation, role enforcement, and response shaping. It exposes trusted internal endpoints under the existing SvelteKit-to-Go internal API pattern. Reads require an authenticated session. Mutations require an authenticated user with role `admin`; the backend enforces this even when the frontend hides controls.

The frontend adds one integrated `/app/organization-units` route. The route is available to authenticated users. It uses server-only helpers to call the Go API with the internal token and current session context. Admin users receive page data that enables mutation controls; regular users receive read-only page data.

The UI uses a tree and details split: the tree remains visible on the left while the selected unit's details and actions appear on the right. This matches the example app's operational layout patterns without copying product feature modules.

## Data Model

Add an `organization_units` table represented by a focused Go domain model:

- `id uuid primary key`
- `parent_id uuid null`, self-referencing `organization_units.id`
- `name varchar(160) not null`
- `code varchar(40) null`
- `description text null`
- `archived_at timestamptz null`
- `created_at timestamptz not null`
- `updated_at timestamptz not null`

The hierarchy uses an adjacency-list parent pointer. A root unit has no `parent_id`. Parent references must point to active units. A unit cannot be moved under itself or one of its descendants.

Codes are optional. When present, a code is trimmed, normalized to uppercase, and unique among active units. Archived units keep their historical code without blocking reuse by active units. The database should enforce active-code uniqueness with a partial unique index for non-null active codes.

Archive is recursive for the selected subtree. Archived units are excluded from the active tree and cannot be used as parents for new or moved units. Active tree responses are ordered by unit name ascending, then id, at each sibling level.

## Backend API

Expose focused internal API operations:

- `GET /api/organization-units/tree` returns the active hierarchy.
- `GET /api/organization-units/archived` returns archived units for admin-facing diagnostics or future restore work and requires an admin user.
- `POST /api/organization-units` creates a root or child unit.
- `PATCH /api/organization-units/:id` updates `name`, `code`, and `description`.
- `PATCH /api/organization-units/:id/parent` moves a unit to another active parent or root.
- `POST /api/organization-units/:id/archive` archives a unit and its descendants.

Request validation happens at the API boundary:

- `name` is required after trimming and limited to 160 characters.
- `code` is optional after trimming and limited to 40 characters.
- `description` is optional after trimming.
- `parent_id` must be empty or a valid active unit id.
- move requests reject cycles and archived targets.
- mutations reject non-admin users with `403`.

Responses use typed JSON shapes and public error messages consistent with existing API helpers.

## Frontend Behavior

Add `/app/organization-units` under the protected app area. The page shows:

- a compact page header;
- a left active-unit tree with selected state;
- a right details panel for the selected unit;
- loading, empty, error, and retry states;
- read-only presentation for normal users;
- admin controls for create root, create child, edit, move, and archive.

Admin create/edit flows use existing local UI primitives before adding new ones. Archive uses a confirmation dialog because it affects the selected subtree. Move uses a parent selector that excludes the selected unit, descendants, and archived units.

The app navigation adds Organization Units as a first-class item. It should not introduce a marketing page or separate admin portal.

## Error Handling

The backend returns specific public errors:

- `400` for invalid JSON, missing names, invalid parent ids, or invalid field lengths.
- `401` when the trusted session is missing or invalid.
- `403` for non-admin mutations.
- `404` for missing or archived units targeted by active operations.
- `409` for duplicate active codes or invalid move cycles.
- `500` with generic text for unexpected failures.

The frontend maps these to user-facing states:

- inline form errors for validation and conflict failures;
- page-level retry alerts for load or network failures;
- clear read-only messaging for users without admin rights;
- confirmation text that states archiving applies to descendants.

## Testing

Backend tests cover:

- route protection and admin-only mutations;
- create validation and successful create;
- update validation and successful update;
- active tree shape and stable ordering;
- move validation, including cycle prevention;
- archive cascade behavior;
- duplicate active code conflicts;
- model registration for Atlas schema generation.

Frontend tests cover:

- server-only Organization Unit API helper success and error mapping;
- page server load for admin and regular users;
- tree utility behavior for flattening, selection, descendant exclusion, and empty states;
- mutation helper request shapes.

Verification commands for implementation:

- backend: `rtk go test ./...`
- backend migration validation: `rtk atlas migrate validate --dir file://migrations`
- frontend: `rtk pnpm test`
- frontend: `rtk pnpm check`
- frontend: `rtk pnpm build`
- browser inspection of `/app/organization-units` for admin and read-only users

## Implementation Constraints

- Keep private values in `.env` files and server-only modules.
- Do not expose internal API tokens to browser code.
- Preserve existing auth route work and do not revert unrelated worktree changes.
- Use `example/` for frontend structure and interaction patterns only; do not copy OCR, billing, email, PDF, or other product feature modules.
- Update Swagger docs when the Go API surface changes.
- Keep the implementation focused on Organization Units and avoid unrelated refactors.
