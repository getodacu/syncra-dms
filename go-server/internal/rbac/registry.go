package rbac

const (
	SystemAdministratorRoleCode       = "system_administrator"
	OrganizationAdministratorRoleCode = "organization_administrator"
	UnitManagerRoleCode               = "unit_manager"
	ViewerRoleCode                    = "viewer"
)

type PermissionDefinition struct {
	Code        string
	Name        string
	Description string
	Category    string
}

type RoleDefinition struct {
	Code            string
	Name            string
	Description     string
	PermissionCodes []string
}

var PermissionRegistry = []PermissionDefinition{
	{Code: "system.admin", Name: "System administration", Category: "System"},
	{Code: "user.view", Name: "View users", Category: "User Management"},
	{Code: "user.create", Name: "Create users", Category: "User Management"},
	{Code: "user.update", Name: "Update users", Category: "User Management"},
	{Code: "user.delete", Name: "Delete users", Category: "User Management"},
	{Code: "user.activate", Name: "Activate users", Category: "User Management"},
	{Code: "user.suspend", Name: "Suspend users", Category: "User Management"},
	{Code: "user.assign_role", Name: "Assign user roles", Category: "User Management"},
	{Code: "user.assign_group", Name: "Assign user groups", Category: "User Management"},
	{Code: "user.assign_unit", Name: "Assign user units", Category: "User Management"},
	{Code: "role.view", Name: "View roles", Category: "Role Management"},
	{Code: "role.create", Name: "Create roles", Category: "Role Management"},
	{Code: "role.update", Name: "Update roles", Category: "Role Management"},
	{Code: "role.delete", Name: "Delete roles", Category: "Role Management"},
	{Code: "role.assign_permissions", Name: "Assign role permissions", Category: "Role Management"},
	{Code: "role.assign_users", Name: "Assign roles to users", Category: "Role Management"},
	{Code: "group.view", Name: "View groups", Category: "Group Management"},
	{Code: "group.create", Name: "Create groups", Category: "Group Management"},
	{Code: "group.update", Name: "Update groups", Category: "Group Management"},
	{Code: "group.delete", Name: "Delete groups", Category: "Group Management"},
	{Code: "group.manage_users", Name: "Manage group users", Category: "Group Management"},
	{Code: "group.assign_roles", Name: "Assign group roles", Category: "Group Management"},
	{Code: "organization_unit.view", Name: "View organization units", Category: "Organization Unit Management"},
	{Code: "organization_unit.create", Name: "Create organization units", Category: "Organization Unit Management"},
	{Code: "organization_unit.update", Name: "Update organization units", Category: "Organization Unit Management"},
	{Code: "organization_unit.delete", Name: "Delete organization units", Category: "Organization Unit Management"},
	{Code: "organization_unit.manage_users", Name: "Manage organization unit users", Category: "Organization Unit Management"},
	{Code: "organization_unit.manage_roles", Name: "Manage organization unit roles", Category: "Organization Unit Management"},
	{Code: "organization_unit.manage_permissions", Name: "Manage organization unit permissions", Category: "Organization Unit Management"},
	{Code: "organization_unit.manage_hierarchy", Name: "Manage organization unit hierarchy", Category: "Organization Unit Management"},
	{Code: "organization_unit.view_audit", Name: "View organization unit audit", Category: "Organization Unit Management"},
}

var organizationAdministratorPermissionCodes = []string{
	"user.view",
	"user.create",
	"user.update",
	"user.delete",
	"user.activate",
	"user.suspend",
	"user.assign_role",
	"user.assign_group",
	"user.assign_unit",
	"role.view",
	"role.create",
	"role.update",
	"role.delete",
	"role.assign_permissions",
	"role.assign_users",
	"group.view",
	"group.create",
	"group.update",
	"group.delete",
	"group.manage_users",
	"group.assign_roles",
	"organization_unit.view",
	"organization_unit.create",
	"organization_unit.update",
	"organization_unit.delete",
	"organization_unit.manage_users",
	"organization_unit.manage_roles",
	"organization_unit.manage_permissions",
	"organization_unit.manage_hierarchy",
	"organization_unit.view_audit",
}

func PermissionCodes() []string {
	codes := make([]string, 0, len(PermissionRegistry))
	for _, definition := range PermissionRegistry {
		codes = append(codes, definition.Code)
	}
	return codes
}

func PermissionByCode(code string) (PermissionDefinition, bool) {
	for _, definition := range PermissionRegistry {
		if definition.Code == code {
			return definition, true
		}
	}
	return PermissionDefinition{}, false
}

func DefaultRoles() []RoleDefinition {
	allPermissions := PermissionCodes()
	organizationAdministratorPermissions := append([]string(nil), organizationAdministratorPermissionCodes...)

	return []RoleDefinition{
		{
			Code:            SystemAdministratorRoleCode,
			Name:            "System Administrator",
			Description:     "Full system administration access.",
			PermissionCodes: allPermissions,
		},
		{
			Code:            OrganizationAdministratorRoleCode,
			Name:            "Organization Administrator",
			Description:     "Administrative access within the organization.",
			PermissionCodes: organizationAdministratorPermissions,
		},
		{
			Code:        UnitManagerRoleCode,
			Name:        "Unit Manager",
			Description: "Manage assigned organization units and users.",
			PermissionCodes: []string{
				"organization_unit.view",
				"organization_unit.update",
				"organization_unit.manage_users",
			},
		},
		{
			Code:        ViewerRoleCode,
			Name:        "Viewer",
			Description: "Read-only organization unit access.",
			PermissionCodes: []string{
				"organization_unit.view",
			},
		},
	}
}
