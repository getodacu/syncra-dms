# RBAC Foundation Design

## Summary

Implement the Syncra DMS RBAC foundation across the Go API and SvelteKit frontend. This slice covers users, roles, permissions, groups, organization-unit membership, scoped role assignment, and backend permission enforcement for existing admin behavior.

The implementation does not add document storage, workflow routing, document visibility, or audit logging. It prepares the authorization model those later features will use.

## Goals

- Replace the current `user.role == admin` authorization checks with a permission resolver.
- Add persistent roles, permissions, groups, role-permission assignments, user-role assignments, group membership, group-role assignments, and organization-unit role assignments.
- Seed a fixed permission registry and default system roles.
- Expand user management with status and organization-unit profile fields needed for RBAC.
- Block inactive, suspended, and deleted users from signing in or continuing sessions.
- Add admin frontend screens for users, roles, permissions, and groups.
- Keep browser code behind SvelteKit server routes and server-only Go API clients.

## Non-Goals

- Document tables, document owner unit, document visibility, or document access checks.
- Workflow task routing or workflow-specific enforcement.
- Audit log implementation.
- External identity provider sync, SCIM, LDAP, or Active Directory integration.
- Direct user permissions, explicit deny rules, temporary access expiry, or field-level permissions.
- Separate admin portal shell.

## Architecture

The Go API owns RBAC persistence, validation, seeding, and authorization decisions. A new backend domain package, likely `internal/rbac`, defines role, permission, group, assignment, scope, and resolver behavior.

Existing auth remains responsible for credentials and sessions. The current `auth.User.Role` enum stays only as a migration bridge while default RBAC roles are introduced. Effective authorization moves to permission checks. Existing organization-unit mutations switch from enum checks to permissions such as `organization_unit.create`, `organization_unit.update`, `organization_unit.delete`, or `organization_unit.manage_hierarchy`.

All RBAC admin endpoints remain internal Go API endpoints protected by the trusted internal request token and an authenticated session. The frontend follows the existing call chain:

browser component -> Svelte `/api/...` route -> server-only frontend client -> Go API

This keeps private tokens and environment values out of browser code.

## Data Model

Add these backend tables:

- `roles`: `id`, `name`, `code`, `description`, `is_system`, `is_active`, `created_at`, `updated_at`.
- `permissions`: `id`, `code`, `name`, `description`, `category`, `is_system`, `created_at`, `updated_at`.
- `role_permissions`: `id`, `role_id`, `permission_id`, `created_at`.
- `user_roles`: `id`, `user_id`, `role_id`, `scope_type`, `organization_unit_id`, `created_at`, `updated_at`.
- `groups`: `id`, `name`, `code`, `description`, `organization_unit_id`, `is_active`, `created_at`, `updated_at`.
- `group_users`: `id`, `group_id`, `user_id`, `created_at`.
- `group_roles`: `id`, `group_id`, `role_id`, `scope_type`, `organization_unit_id`, `created_at`.
- `organization_unit_roles`: `id`, `organization_unit_id`, `role_id`, `scope_type`, `created_at`.

Extend the existing user table with:

- `status`
- `primary_organization_unit_id`
- `manager_user_id`
- `job_title`
- `phone`
- `deleted_at`

Supported statuses are `invited`, `active`, `inactive`, `suspended`, and `deleted`. Existing verified users migrate to `active`; unverified users migrate to `invited`.

Supported scope types for this slice are:

- `global`
- `organization_unit`
- `organization_unit_and_children`

## Backend API

Add internal API endpoints for:

- Users: list, read, create, update, activate, deactivate, suspend, soft-delete, assign primary organization unit, assign and remove scoped roles, assign and remove groups.
- Roles: list, read, create, update, activate, deactivate, delete if unused, list permissions, assign permissions, remove permissions.
- Permissions: list registry and list categories.
- Groups: list, read, create, update, delete if unused, add and remove users, assign and remove scoped roles.
- Current user authorization: `/api/me`, `/api/me/permissions`, and `/api/auth/check-permission`.

Hard delete should be avoided for users. Roles and groups may be deleted only when unused and only if they are not system records.

## Permission Resolution

The resolver denies by default. It allows only when all of these are true:

- the user exists and has status `active`;
- the requested permission code exists in the registry;
- the permission is granted through one of the supported assignment paths;
- the assignment scope matches the requested resource scope.

Supported assignment paths:

- user role assignments;
- group role assignments through group membership;
- organization-unit role assignments through the user's primary organization unit.

Scope matching:

- `global` matches every resource.
- `organization_unit` matches only the specified organization unit.
- `organization_unit_and_children` matches the specified organization unit and its descendants using the existing organization-unit tree helpers.

Direct user permissions and explicit deny rules are deferred.

## Auth Behavior

Sign-in fails for inactive, suspended, or deleted users. Existing sessions for users in those statuses are rejected when loaded. Suspending or deactivating a user should revoke active sessions.

The existing `admin` enum can be used during migration to seed or assign the default `system_administrator` role, but new route authorization should call the RBAC resolver.

## Frontend

Add admin routes under `/app/admin`:

- `/app/admin/users`: user list, create user, edit profile/status, assign primary organization unit, assign scoped roles, assign groups.
- `/app/admin/roles`: role list, create/edit/activate/deactivate roles, assign permissions through a category matrix.
- `/app/admin/permissions`: read-only permission registry grouped by category.
- `/app/admin/groups`: group list, create/edit groups, manage members, assign scoped roles.

The app shell adds an Admin navigation section only when the current user has relevant management permissions. Existing `/app/organization-units` keeps its layout but shows mutation controls based on effective permissions instead of `role === 'admin'`.

The UI should stay dense and operational: tables, inline filters, compact forms, scoped selectors, and confirmation dialogs for destructive changes.

## Error Handling

Backend errors:

- `400`: invalid JSON, invalid IDs, invalid status, invalid scope, invalid assignments.
- `401`: missing or invalid authenticated session.
- `403`: active user lacks the required permission.
- `404`: missing user, role, permission, group, or organization unit.
- `409`: duplicate codes, unsafe deletion, system record mutation, or conflicting assignment.
- `500`: unexpected backend failure with generic public text.

Frontend routes map backend errors to public-safe JSON and user-facing page states. The Go API remains authoritative even when the frontend hides unauthorized actions.

## Testing

Backend behavior tests cover:

- RBAC model validation and migration registration.
- Permission registry seeding and idempotency.
- Permission resolver paths and scope matching.
- Sign-in and session rejection for inactive, suspended, and deleted users.
- User, role, permission, and group APIs.
- Organization-unit mutation enforcement through permissions.

Frontend tests cover:

- Server-only RBAC API clients.
- Svelte API proxy auth and error behavior.
- Page load permission gating.
- Critical user, role, permission, and group UI state behavior.
- Organization-unit mutation controls switching from enum checks to effective permissions.

Verification commands:

```sh
cd go-server
rtk go test ./...
rtk atlas migrate validate --dir file://migrations

cd frontend
rtk pnpm test
rtk pnpm check
rtk pnpm build
```

## Implementation Constraints

- Keep private values in `.env` files and server-only modules.
- Do not expose internal API tokens to browser code.
- Preserve the existing auth and organization-unit flows while replacing authorization checks incrementally.
- Do not copy product feature modules from `example/`.
- Update Swagger docs when Go API routes or shapes change.
- Add behavior tests before production behavior.
