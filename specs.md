Syncra DMS — User Management, RBAC, Permissions and Organization Units Specification
1. Purpose

Syncra DMS shall provide a secure user management and authorization system that controls access to documents, workflows, administration features, organization units, and audit data.

The system shall support:

Users
Roles
Permissions
Groups
Organization Units
Unit-based access control
Document-level access control
Workflow permissions
Administrative permissions
Auditability

The authorization model shall be based on RBAC, extended with Organization Unit scope.

2. Core Concepts
2.1 User

A user represents a person who can authenticate and use Syncra DMS.

A user may be:

internal employee
manager
administrator
external collaborator
auditor
signer
viewer
service account

Each user shall have:

name
email
status
primary organization unit
secondary organization units
roles
groups
permissions
manager
job title
authentication method
2.2 Role

A role is a reusable collection of permissions.

Example roles:

System Administrator
Organization Administrator
Unit Manager
Document Owner
Document Contributor
Document Viewer
Workflow Designer
Workflow Approver
Digital Signer
Auditor
External Collaborator

Roles should not be hard-coded. The system shall provide default roles, but administrators shall be able to create custom roles.

2.3 Permission

A permission is a low-level action allowed in the system.

Example:

document.view
document.create
document.update
document.delete
document.sign
workflow.start
workflow.approve
organization_unit.manage
user.create
role.assign
audit.view

Permissions are assigned to roles, groups, users, or organization units.

2.4 Group

A group is a collection of users used for easier assignment.

Examples:

Finance Approvers
Legal Reviewers
Executive Board
HR Managers
External Auditors
Signing Committee

Groups may have roles and permissions.

2.5 Organization Unit

An Organization Unit represents a company structure element, such as:

company
division
department
team
office
branch
cost center
project unit
external unit

Users, documents, permissions, workflows, and repository spaces may reference Organization Units.

3. Authorization Model

Syncra DMS shall use:

RBAC = Role-Based Access Control
+
Unit Scope = Organization Unit-based access restriction
+
Document Scope = document ownership and sharing rules

A permission decision shall depend on:

user permissions
role permissions
group permissions
organization unit membership
document owner unit
document sharing rules
workflow assignment
administrative scope

Example:

Ana has role Invoice Approver in Accounting Unit.
She can approve invoices owned by Accounting Unit.
She cannot approve contracts owned by Legal Unit.
4. Permission Scopes

Permissions shall support scope.

Recommended scopes:

global
organization_unit
organization_unit_and_children
own_documents
assigned_documents
shared_documents
workflow_assigned

Example:

document.view + organization_unit_and_children

Means the user can view documents owned by their unit and child units.

5. Functional Requirements
5.1 User Management

Administrators shall be able to:

create users
invite users
edit users
activate users
deactivate users
suspend users
delete users only if unused
assign users to organization units
assign roles to users
assign users to groups
define primary organization unit
define secondary organization units
define user manager
reset authentication credentials
view user activity

User statuses:

invited
active
inactive
suspended
deleted
5.2 Role Management

Administrators shall be able to:

create roles
edit roles
delete roles if unused
assign permissions to roles
assign roles to users
assign roles to groups
assign roles within organization unit scope
clone predefined roles
activate/deactivate roles

Default roles:

System Administrator
DMS Administrator
Organization Administrator
Unit Manager
Document Manager
Document Contributor
Document Viewer
Workflow Designer
Workflow Approver
Signer
Auditor
External Viewer
5.3 Permission Management

The system shall provide a fixed permission registry.

Administrators may assign permissions, but should not create arbitrary low-level permission codes from the UI.

Permission categories:

User Management
Role Management
Organization Unit Management
Document Repository
Document Metadata
Document Versioning
Workflow
Digital Signature
OCR and AI Processing
Audit
System Settings
Reports
5.4 Group Management

Administrators shall be able to:

create groups
edit groups
delete groups if unused
add users to groups
remove users from groups
assign roles to groups
assign permissions to groups
use groups in workflow routing

Example:

Group: Legal Reviewers
Members: legal users
Role: Document Reviewer
Scope: Legal Department
5.5 Organization Unit-Based Access

The system shall support permissions assigned to Organization Units.

Example:

Legal Department:
- document.view
- document.create
- workflow.approve_contract

Users assigned to that unit may inherit those permissions depending on configuration.

Inheritance options:

no inheritance
inherit from parent unit
inherit to child units
inherit both directions only if explicitly enabled

Recommended MVP behavior:

parent-to-child inheritance only
5.6 Document Access Control

Each document shall have:

owner organization unit
created_by user
visibility level
access list
workflow state
document category
confidentiality level

Access rules shall consider:

owner unit
user unit
role
group
explicit share
workflow assignment
document confidentiality

Document visibility levels:

private
unit
unit_and_children
organization
restricted
external_shared
6. Recommended Permissions
6.1 User Permissions
user.view
user.create
user.update
user.delete
user.activate
user.suspend
user.assign_role
user.assign_group
user.assign_unit
6.2 Role Permissions
role.view
role.create
role.update
role.delete
role.assign_permissions
role.assign_users
6.3 Group Permissions
group.view
group.create
group.update
group.delete
group.manage_users
group.assign_roles
6.4 Organization Unit Permissions
organization_unit.view
organization_unit.create
organization_unit.update
organization_unit.delete
organization_unit.manage_users
organization_unit.manage_roles
organization_unit.manage_permissions
organization_unit.manage_hierarchy
organization_unit.view_audit
6.5 Document Permissions
document.view
document.create
document.update
document.delete
document.archive
document.restore
document.upload_version
document.download
document.share
document.manage_permissions
document.view_audit
document.change_owner_unit
6.6 Workflow Permissions
workflow.view
workflow.create
workflow.update
workflow.delete
workflow.start
workflow.cancel
workflow.approve
workflow.reject
workflow.delegate
workflow.escalate
workflow.design
workflow.manage_templates
6.7 Signature Permissions
signature.request
signature.sign
signature.validate
signature.reject
signature.view_certificate
6.8 AI/OCR Permissions
ocr.start
ocr.view_result
ocr.correct_result
ocr.export_result
ai.summarize_document
ai.classify_document
ai.extract_structured_data
6.9 Audit Permissions
audit.view
audit.export
audit.view_security_events
audit.view_document_events
audit.view_user_events
7. Data Model Specification
7.1 users
id
email
first_name
last_name
display_name
password_hash
status
primary_organization_unit_id
manager_user_id
job_title
phone
auth_provider
last_login_at
created_at
updated_at
deleted_at
7.2 roles
id
name
code
description
is_system
is_active
created_at
updated_at
7.3 permissions
id
code
name
description
category
is_system
created_at
updated_at
7.4 role_permissions
id
role_id
permission_id
created_at
7.5 user_roles
id
user_id
role_id
scope_type
organization_unit_id
created_at
updated_at

Scope types:

global
organization_unit
organization_unit_and_children
own_documents
assigned_documents
7.6 groups
id
name
code
description
organization_unit_id
is_active
created_at
updated_at
7.7 group_users
id
group_id
user_id
created_at
7.8 group_roles
id
group_id
role_id
scope_type
organization_unit_id
created_at
7.9 organization_unit_roles
id
organization_unit_id
role_id
scope_type
created_at
7.10 document_permissions
id
document_id
principal_type
principal_id
permission_id
scope_type
created_at
updated_at

Principal types:

user
role
group
organization_unit
8. Permission Resolution Logic

When checking access, the backend shall evaluate:

1. Is user active?
2. Is document active and accessible?
3. Does user have global permission?
4. Does user have role permission?
5. Does user have group permission?
6. Does user have unit-based permission?
7. Does document owner unit match user scope?
8. Is user explicitly assigned to workflow step?
9. Is user explicitly shared on the document?
10. Is document confidentiality level compatible?

Access shall be denied by default.

Default = deny
Explicit allow = allow
Explicit deny = deny
Deny overrides allow
9. REST API Specification
9.1 User API
GET    /api/users
GET    /api/users/:id
POST   /api/users
PUT    /api/users/:id
DELETE /api/users/:id
POST   /api/users/:id/activate
POST   /api/users/:id/suspend
POST   /api/users/:id/roles
DELETE /api/users/:id/roles/:roleId
POST   /api/users/:id/groups
DELETE /api/users/:id/groups/:groupId
POST   /api/users/:id/organization-units
DELETE /api/users/:id/organization-units/:unitId
9.2 Role API
GET    /api/roles
GET    /api/roles/:id
POST   /api/roles
PUT    /api/roles/:id
DELETE /api/roles/:id
GET    /api/roles/:id/permissions
POST   /api/roles/:id/permissions
DELETE /api/roles/:id/permissions/:permissionId
9.3 Permission API
GET /api/permissions
GET /api/permissions/categories
9.4 Group API
GET    /api/groups
GET    /api/groups/:id
POST   /api/groups
PUT    /api/groups/:id
DELETE /api/groups/:id
POST   /api/groups/:id/users
DELETE /api/groups/:id/users/:userId
POST   /api/groups/:id/roles
DELETE /api/groups/:id/roles/:roleId
9.5 Authorization API
GET  /api/me
GET  /api/me/permissions
POST /api/auth/check-permission

Example:

{
  "user_id": "user-id",
  "permission": "document.view",
  "resource_type": "document",
  "resource_id": "document-id"
}
10. Frontend Requirements

Admin screens:

Users
User Details
Create/Edit User
Roles
Role Details
Permissions
Groups
Group Details
Organization Unit Users
Organization Unit Roles
Access Matrix
Audit Log

Useful UI components:

role selector
permission matrix
organization unit tree selector
user assignment dialog
group membership dialog
scope selector
effective permissions viewer
11. MVP Scope

For MVP, implement:

User CRUD
Role CRUD
Permission registry
Role-permission assignment
User-role assignment
Group CRUD
Group-user assignment
Organization Unit membership
Primary Organization Unit
Unit Manager
Document owner unit
Document visibility rules
Backend permission enforcement
Admin UI
Audit log

Out of scope for MVP:

external identity provider sync
Active Directory / LDAP sync
SCIM provisioning
advanced ABAC
complex deny rules
temporary access expiry
multi-tenant shared SaaS
fine-grained field-level permissions
12. Acceptance Criteria

The implementation is acceptable when:

admin can create users
admin can assign users to organization units
admin can create roles
admin can assign permissions to roles
admin can assign roles to users and groups
admin can assign scoped roles by organization unit
users can access only allowed documents
unit managers can access documents in their unit if permitted
workflow tasks can be routed to roles, groups, users, or unit managers
all authorization decisions are enforced in backend
frontend hides unauthorized actions
audit events are recorded for user, role, permission, and unit changes
inactive users cannot log in
suspended users immediately lose access
13. Recommended MVP Authorization Rule

For the first version, use this simple rule:

A user may access a document if:
1. the user has the required permission;
2. the permission scope matches the document owner unit;
3. the document visibility allows access;
4. the user is active.