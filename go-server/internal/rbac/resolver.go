package rbac

import (
	"context"
	"errors"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"gorm.io/gorm"
)

type Resolver struct {
	db *gorm.DB
}

type Check struct {
	UserID             string
	Permission         string
	OrganizationUnitID *string
}

type Grant struct {
	PermissionCode     string
	ScopeType          ScopeType
	OrganizationUnitID *string
	Source             string
}

const (
	GrantSourceUserRole             = "user_role"
	GrantSourceGroupRole            = "group_role"
	GrantSourceOrganizationUnitRole = "organization_unit_role"
)

func NewResolver(db *gorm.DB) *Resolver {
	return &Resolver{db: db}
}

func (r *Resolver) Can(ctx context.Context, check Check) (bool, error) {
	grants, err := r.EffectiveGrants(ctx, check.UserID)
	if err != nil {
		return false, err
	}
	for _, grant := range grants {
		if grant.PermissionCode != check.Permission {
			continue
		}
		matches, err := r.scopeMatches(ctx, grant, check.OrganizationUnitID)
		if err != nil {
			return false, err
		}
		if matches {
			return true, nil
		}
	}
	return false, nil
}

func (r *Resolver) EffectiveGrants(ctx context.Context, userID string) ([]Grant, error) {
	var user auth.User
	if err := r.db.WithContext(ctx).
		Select("id", "status", "primary_organization_unit_id").
		First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if user.Status != string(UserStatusActive) {
		return nil, nil
	}

	grants := make([]Grant, 0)
	userGrants, err := r.userRoleGrants(ctx, userID)
	if err != nil {
		return nil, err
	}
	grants = append(grants, userGrants...)

	groupGrants, err := r.groupRoleGrants(ctx, userID)
	if err != nil {
		return nil, err
	}
	grants = append(grants, groupGrants...)

	if user.PrimaryOrganizationUnitID != nil {
		organizationUnitGrants, err := r.organizationUnitRoleGrants(ctx, *user.PrimaryOrganizationUnitID)
		if err != nil {
			return nil, err
		}
		grants = append(grants, organizationUnitGrants...)
	}

	return grants, nil
}

func (r *Resolver) userRoleGrants(ctx context.Context, userID string) ([]Grant, error) {
	var grants []Grant
	if err := r.db.WithContext(ctx).
		Table("user_roles").
		Select("? AS source, permissions.code AS permission_code, user_roles.scope_type, user_roles.organization_unit_id", GrantSourceUserRole).
		Joins("JOIN roles ON roles.id = user_roles.role_id AND roles.is_active = ?", true).
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("user_roles.user_id = ?", userID).
		Scan(&grants).Error; err != nil {
		return nil, err
	}
	return grants, nil
}

func (r *Resolver) groupRoleGrants(ctx context.Context, userID string) ([]Grant, error) {
	var grants []Grant
	if err := r.db.WithContext(ctx).
		Table("group_users").
		Select("? AS source, permissions.code AS permission_code, group_roles.scope_type, group_roles.organization_unit_id", GrantSourceGroupRole).
		Joins("JOIN groups ON groups.id = group_users.group_id AND groups.is_active = ?", true).
		Joins("JOIN group_roles ON group_roles.group_id = groups.id").
		Joins("JOIN roles ON roles.id = group_roles.role_id AND roles.is_active = ?", true).
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("group_users.user_id = ?", userID).
		Scan(&grants).Error; err != nil {
		return nil, err
	}
	return grants, nil
}

func (r *Resolver) organizationUnitRoleGrants(ctx context.Context, organizationUnitID string) ([]Grant, error) {
	var grants []Grant
	if err := r.db.WithContext(ctx).
		Table("organization_unit_roles").
		Select("? AS source, permissions.code AS permission_code, organization_unit_roles.scope_type, organization_unit_roles.organization_unit_id", GrantSourceOrganizationUnitRole).
		Joins("JOIN roles ON roles.id = organization_unit_roles.role_id AND roles.is_active = ?", true).
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("organization_unit_roles.organization_unit_id = ?", organizationUnitID).
		Scan(&grants).Error; err != nil {
		return nil, err
	}
	return grants, nil
}

func (r *Resolver) scopeMatches(ctx context.Context, grant Grant, requested *string) (bool, error) {
	switch grant.ScopeType {
	case ScopeGlobal:
		return true, nil
	case ScopeOrganizationUnit:
		return grant.OrganizationUnitID != nil && requested != nil && *grant.OrganizationUnitID == *requested, nil
	case ScopeOrganizationUnitAndChildren:
		if grant.OrganizationUnitID == nil || requested == nil {
			return false, nil
		}
		if *grant.OrganizationUnitID == *requested {
			return true, nil
		}
		var units []orgunits.Unit
		if err := r.db.WithContext(ctx).Where("archived_at IS NULL").Find(&units).Error; err != nil {
			return false, err
		}
		return orgunits.DescendantIDs(*grant.OrganizationUnitID, units)[*requested], nil
	default:
		return false, nil
	}
}
