package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const CredentialProviderID = "credential"
const GoogleProviderID = "google"
const GitHubProviderID = "github"

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type User struct {
	ID                        string     `gorm:"type:uuid;primaryKey;index:idx_user_created_id,priority:2;index:idx_user_last_login_id,priority:2" json:"id"`
	Name                      string     `gorm:"not null;size:255" json:"name"`
	Email                     string     `gorm:"not null;size:320;uniqueIndex" json:"email"`
	EmailVerified             bool       `gorm:"column:email_verified;not null;default:false" json:"emailVerified"`
	Image                     *string    `gorm:"type:text" json:"image,omitempty"`
	PreferredLanguage         string     `gorm:"column:preferred_language;not null;size:5;default:en;check:chk_user_preferred_language,preferred_language IN ('en','ro')" json:"preferredLanguage"`
	Role                      UserRole   `gorm:"not null;size:40;default:user;index;check:chk_user_role,role IN ('user','admin')" json:"role"`
	Status                    string     `gorm:"not null;size:40;default:active;index;check:chk_user_status,status IN ('invited','active','inactive','suspended','deleted')" json:"status"`
	PrimaryOrganizationUnitID *string    `gorm:"column:primary_organization_unit_id;type:uuid;index" json:"primaryOrganizationUnitId,omitempty"`
	ManagerUserID             *string    `gorm:"column:manager_user_id;type:uuid;index" json:"managerUserId,omitempty"`
	JobTitle                  *string    `gorm:"column:job_title;size:160" json:"jobTitle,omitempty"`
	Phone                     *string    `gorm:"size:80" json:"phone,omitempty"`
	LastLoginAt               *time.Time `gorm:"column:last_login_at;index:idx_user_last_login_id,priority:1" json:"lastLoginAt,omitempty"`
	DeletedAt                 *time.Time `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
	CreatedAt                 time.Time  `gorm:"column:created_at;not null;index:idx_user_created_id,priority:1" json:"createdAt"`
	UpdatedAt                 time.Time  `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	if u.Role == "" {
		u.Role = UserRoleUser
	}
	if u.Status == "" {
		if u.EmailVerified {
			u.Status = "active"
		} else {
			u.Status = "invited"
		}
	}
	if u.PreferredLanguage == "" {
		u.PreferredLanguage = "en"
	}
	return nil
}

func (User) TableName() string {
	return "user"
}

type AuthAccount struct {
	ID                    string     `gorm:"type:uuid;primaryKey" json:"id"`
	AccountID             string     `gorm:"column:account_id;not null;size:255;index;uniqueIndex:idx_account_provider_account,priority:2" json:"accountId"`
	ProviderID            string     `gorm:"column:provider_id;not null;size:120;index;uniqueIndex:idx_account_provider_account,priority:1" json:"providerId"`
	UserID                string     `gorm:"column:user_id;type:uuid;not null;index" json:"userId"`
	User                  User       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	AccessToken           *string    `gorm:"column:access_token;type:text" json:"accessToken,omitempty"`
	RefreshToken          *string    `gorm:"column:refresh_token;type:text" json:"refreshToken,omitempty"`
	IDToken               *string    `gorm:"column:id_token;type:text" json:"idToken,omitempty"`
	AccessTokenExpiresAt  *time.Time `gorm:"column:access_token_expires_at" json:"accessTokenExpiresAt,omitempty"`
	RefreshTokenExpiresAt *time.Time `gorm:"column:refresh_token_expires_at" json:"refreshTokenExpiresAt,omitempty"`
	Scope                 *string    `gorm:"type:text" json:"scope,omitempty"`
	Password              string     `gorm:"type:text" json:"-"`
	CreatedAt             time.Time  `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt             time.Time  `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (a *AuthAccount) BeforeCreate(_ *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	return nil
}

func (AuthAccount) TableName() string {
	return "account"
}

type Session struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null;index" json:"expiresAt"`
	Token     string    `gorm:"not null;size:255;uniqueIndex" json:"token"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updatedAt"`
	IPAddress string    `gorm:"column:ip_address;type:text" json:"ipAddress,omitempty"`
	UserAgent string    `gorm:"column:user_agent;type:text" json:"userAgent,omitempty"`
	UserID    string    `gorm:"column:user_id;type:uuid;not null;index" json:"userId"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (s *Session) BeforeCreate(_ *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	return nil
}

func (Session) TableName() string {
	return "session"
}

type Verification struct {
	ID         string    `gorm:"type:uuid;primaryKey" json:"id"`
	Identifier string    `gorm:"not null;size:512;uniqueIndex:idx_verification_identifier_unique" json:"identifier"`
	Value      string    `gorm:"type:text;not null" json:"value"`
	ExpiresAt  time.Time `gorm:"column:expires_at;not null;index" json:"expiresAt"`
	CreatedAt  time.Time `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (v *Verification) BeforeCreate(_ *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.NewString()
	}
	return nil
}

func (Verification) TableName() string {
	return "verification"
}
