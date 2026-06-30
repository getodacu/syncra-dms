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
	ID                string     `gorm:"type:uuid;primaryKey;index:idx_user_created_id,priority:2;index:idx_user_last_login_id,priority:2" json:"id"`
	Name              string     `gorm:"not null;size:255" json:"name"`
	Email             string     `gorm:"not null;size:320;uniqueIndex" json:"email"`
	EmailVerified     bool       `gorm:"column:email_verified;not null;default:false" json:"emailVerified"`
	Image             *string    `gorm:"type:text" json:"image,omitempty"`
	PreferredLanguage string     `gorm:"column:preferred_language;not null;size:5;default:en;check:chk_user_preferred_language,preferred_language IN ('en','ro')" json:"preferredLanguage"`
	Role              UserRole   `gorm:"not null;size:40;default:user;index;check:chk_user_role,role IN ('user','admin')" json:"role"`
	LastLoginAt       *time.Time `gorm:"column:last_login_at;index:idx_user_last_login_id,priority:1" json:"lastLoginAt,omitempty"`
	CreatedAt         time.Time  `gorm:"column:created_at;not null;index:idx_user_created_id,priority:1" json:"createdAt"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	if u.Role == "" {
		u.Role = UserRoleUser
	}
	if u.PreferredLanguage == "" {
		u.PreferredLanguage = "en"
	}
	return nil
}

// Better Auth expects this exact singular table name; quote "user" in raw SQL.
func (User) TableName() string {
	return "user"
}

type APIKey struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string     `gorm:"column:user_id;type:uuid;not null;index" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Name      string     `gorm:"not null;size:255" json:"name"`
	KeyHash   string     `gorm:"column:key_hash;not null;size:64;uniqueIndex" json:"-"`
	KeyPrefix string     `gorm:"column:key_prefix;not null;size:8" json:"key_prefix"`
	ExpiresAt *time.Time `gorm:"column:expires_at;index" json:"expires_at,omitempty"`
	CreatedAt time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (k *APIKey) BeforeCreate(_ *gorm.DB) error {
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}
	return nil
}

func (APIKey) TableName() string {
	return "api_keys"
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
	ID                     string     `gorm:"type:uuid;primaryKey" json:"id"`
	ExpiresAt              time.Time  `gorm:"column:expires_at;not null;index" json:"expiresAt"`
	Token                  string     `gorm:"not null;size:255;uniqueIndex" json:"token"`
	CreatedAt              time.Time  `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt              time.Time  `gorm:"column:updated_at;not null" json:"updatedAt"`
	IPAddress              string     `gorm:"column:ip_address;type:text" json:"ipAddress,omitempty"`
	UserAgent              string     `gorm:"column:user_agent;type:text" json:"userAgent,omitempty"`
	UserID                 string     `gorm:"column:user_id;type:uuid;not null;index" json:"userId"`
	User                   User       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	ImpersonatedUserID     *string    `gorm:"column:impersonated_user_id;type:uuid;index" json:"impersonatedUserId,omitempty"`
	ImpersonatedUser       *User      `gorm:"foreignKey:ImpersonatedUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	ImpersonationStartedAt *time.Time `gorm:"column:impersonation_started_at" json:"impersonationStartedAt,omitempty"`
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

type AdminImpersonationEventType string

const (
	AdminImpersonationEventStart AdminImpersonationEventType = "start"
	AdminImpersonationEventStop  AdminImpersonationEventType = "stop"
)

type AdminImpersonationEvent struct {
	ID              uuid.UUID                   `gorm:"type:uuid;primaryKey" json:"id"`
	EventType       AdminImpersonationEventType `gorm:"column:event_type;not null;size:20;index;check:chk_admin_impersonation_event_type,event_type IN ('start','stop')" json:"event_type"`
	SessionID       string                      `gorm:"column:session_id;type:uuid;not null;index" json:"session_id"`
	AdminUserID     string                      `gorm:"column:admin_user_id;type:uuid;not null;index" json:"admin_user_id"`
	AdminUserEmail  string                      `gorm:"column:admin_user_email;not null;size:320" json:"admin_user_email"`
	TargetUserID    string                      `gorm:"column:target_user_id;type:uuid;not null;index" json:"target_user_id"`
	TargetUserEmail string                      `gorm:"column:target_user_email;not null;size:320" json:"target_user_email"`
	IPAddress       string                      `gorm:"column:ip_address;type:text" json:"ip_address,omitempty"`
	UserAgent       string                      `gorm:"column:user_agent;type:text" json:"user_agent,omitempty"`
	CreatedAt       time.Time                   `gorm:"column:created_at;not null;index" json:"created_at"`
}

func (e *AdminImpersonationEvent) BeforeCreate(_ *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now().UTC()
	}
	return nil
}

func (AdminImpersonationEvent) TableName() string {
	return "admin_impersonation_events"
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
