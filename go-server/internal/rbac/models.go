package rbac

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScopeType string

const (
	ScopeGlobal                      ScopeType = "global"
	ScopeOrganizationUnit            ScopeType = "organization_unit"
	ScopeOrganizationUnitAndChildren ScopeType = "organization_unit_and_children"
)

func (s ScopeType) Valid() bool {
	return s == ScopeGlobal || s == ScopeOrganizationUnit || s == ScopeOrganizationUnitAndChildren
}

type UserStatus string

const (
	UserStatusInvited   UserStatus = "invited"
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

func (s UserStatus) Valid() bool {
	return s == UserStatusInvited || s == UserStatusActive || s == UserStatusInactive || s == UserStatusSuspended || s == UserStatusDeleted
}

type Role struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"not null;size:160" json:"name"`
	Code        string    `gorm:"not null;size:80;uniqueIndex" json:"code"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	IsSystem    bool      `gorm:"column:is_system;not null;default:false" json:"isSystem"`
	IsActive    bool      `gorm:"column:is_active;not null;default:true;index" json:"isActive"`
	CreatedAt   time.Time `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (r *Role) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}

func (Role) TableName() string { return "roles" }

type Permission struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	Code        string    `gorm:"not null;size:120;uniqueIndex" json:"code"`
	Name        string    `gorm:"not null;size:160" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	Category    string    `gorm:"not null;size:80;index" json:"category"`
	IsSystem    bool      `gorm:"column:is_system;not null;default:true" json:"isSystem"`
	CreatedAt   time.Time `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (p *Permission) BeforeCreate(_ *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return nil
}

func (Permission) TableName() string { return "permissions" }

type RolePermission struct {
	ID           string     `gorm:"type:uuid;primaryKey" json:"id"`
	RoleID       string     `gorm:"column:role_id;type:uuid;not null;uniqueIndex:idx_role_permission_unique,priority:1" json:"roleId"`
	Role         Role       `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	PermissionID string     `gorm:"column:permission_id;type:uuid;not null;uniqueIndex:idx_role_permission_unique,priority:2" json:"permissionId"`
	Permission   Permission `gorm:"foreignKey:PermissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null" json:"createdAt"`
}

func (r *RolePermission) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}

func (RolePermission) TableName() string { return "role_permissions" }

type UserRole struct {
	ID                 string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID             string         `gorm:"column:user_id;type:uuid;not null;index;uniqueIndex:idx_user_role_scope_unique,priority:1" json:"userId"`
	User               auth.User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	RoleID             string         `gorm:"column:role_id;type:uuid;not null;index;uniqueIndex:idx_user_role_scope_unique,priority:2" json:"roleId"`
	Role               Role           `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	ScopeType          ScopeType      `gorm:"column:scope_type;not null;size:40;uniqueIndex:idx_user_role_scope_unique,priority:3" json:"scopeType"`
	OrganizationUnitID *string        `gorm:"column:organization_unit_id;type:uuid;uniqueIndex:idx_user_role_scope_unique,priority:4" json:"organizationUnitId,omitempty"`
	OrganizationUnit   *orgunits.Unit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	CreatedAt          time.Time      `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (r *UserRole) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}

func (UserRole) TableName() string { return "user_roles" }

type Group struct {
	ID                 string         `gorm:"type:uuid;primaryKey" json:"id"`
	Name               string         `gorm:"not null;size:160" json:"name"`
	Code               string         `gorm:"not null;size:80;uniqueIndex" json:"code"`
	Description        *string        `gorm:"type:text" json:"description,omitempty"`
	OrganizationUnitID *string        `gorm:"column:organization_unit_id;type:uuid;index" json:"organizationUnitId,omitempty"`
	OrganizationUnit   *orgunits.Unit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	IsActive           bool           `gorm:"column:is_active;not null;default:true;index" json:"isActive"`
	CreatedAt          time.Time      `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (g *Group) BeforeCreate(_ *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.NewString()
	}
	return nil
}

func (Group) TableName() string { return "groups" }

type GroupUser struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	GroupID   string    `gorm:"column:group_id;type:uuid;not null;uniqueIndex:idx_group_user_unique,priority:1" json:"groupId"`
	Group     Group     `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	UserID    string    `gorm:"column:user_id;type:uuid;not null;uniqueIndex:idx_group_user_unique,priority:2" json:"userId"`
	User      auth.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"createdAt"`
}

func (g *GroupUser) BeforeCreate(_ *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.NewString()
	}
	return nil
}

func (GroupUser) TableName() string { return "group_users" }

type GroupRole struct {
	ID                 string         `gorm:"type:uuid;primaryKey" json:"id"`
	GroupID            string         `gorm:"column:group_id;type:uuid;not null;uniqueIndex:idx_group_role_scope_unique,priority:1" json:"groupId"`
	Group              Group          `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	RoleID             string         `gorm:"column:role_id;type:uuid;not null;uniqueIndex:idx_group_role_scope_unique,priority:2" json:"roleId"`
	Role               Role           `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	ScopeType          ScopeType      `gorm:"column:scope_type;not null;size:40;uniqueIndex:idx_group_role_scope_unique,priority:3" json:"scopeType"`
	OrganizationUnitID *string        `gorm:"column:organization_unit_id;type:uuid;uniqueIndex:idx_group_role_scope_unique,priority:4" json:"organizationUnitId,omitempty"`
	OrganizationUnit   *orgunits.Unit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	CreatedAt          time.Time      `gorm:"column:created_at;not null" json:"createdAt"`
}

func (g *GroupRole) BeforeCreate(_ *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.NewString()
	}
	return nil
}

func (GroupRole) TableName() string { return "group_roles" }

type OrganizationUnitRole struct {
	ID                 string        `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationUnitID string        `gorm:"column:organization_unit_id;type:uuid;not null;uniqueIndex:idx_organization_unit_role_unique,priority:1" json:"organizationUnitId"`
	OrganizationUnit   orgunits.Unit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	RoleID             string        `gorm:"column:role_id;type:uuid;not null;uniqueIndex:idx_organization_unit_role_unique,priority:2" json:"roleId"`
	Role               Role          `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	ScopeType          ScopeType     `gorm:"column:scope_type;not null;size:40;uniqueIndex:idx_organization_unit_role_unique,priority:3" json:"scopeType"`
	CreatedAt          time.Time     `gorm:"column:created_at;not null" json:"createdAt"`
}

func (o *OrganizationUnitRole) BeforeCreate(_ *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.NewString()
	}
	return nil
}

func (OrganizationUnitRole) TableName() string { return "organization_unit_roles" }

var codeNonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func NormalizeCode(raw string) (string, error) {
	code := strings.TrimSpace(strings.ToLower(raw))
	code = codeNonAlnum.ReplaceAllString(code, "_")
	code = strings.Trim(code, "_")
	if code == "" {
		return "", errors.New("code is required")
	}
	if utf8.RuneCountInString(code) > 80 {
		return "", errors.New("code must be at most 80 characters")
	}
	return code, nil
}
