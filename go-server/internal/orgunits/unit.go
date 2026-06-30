package orgunits

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const MaxNameCharacters = 160
const MaxCodeCharacters = 40

type Unit struct {
	ID          string     `gorm:"type:uuid;primaryKey;index:idx_organization_units_parent_name_id,priority:3" json:"id"`
	ParentID    *string    `gorm:"column:parent_id;type:uuid;index:idx_organization_units_parent_name_id,priority:1" json:"parentId,omitempty"`
	Parent      *Unit      `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Name        string     `gorm:"not null;size:160;index:idx_organization_units_parent_name_id,priority:2" json:"name"`
	Code        *string    `gorm:"size:40;index;uniqueIndex:idx_organization_units_active_code_unique,where:code IS NOT NULL AND archived_at IS NULL" json:"code"`
	Description *string    `gorm:"type:text" json:"description"`
	ArchivedAt  *time.Time `gorm:"column:archived_at;index" json:"archivedAt"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (u *Unit) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return nil
}

func (Unit) TableName() string {
	return "organization_units"
}

func NormalizeName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", errors.New("name is required")
	}
	if utf8.RuneCountInString(name) > MaxNameCharacters {
		return "", errors.New("name must be 160 characters or fewer")
	}
	return name, nil
}

func NormalizeCode(raw string) (*string, error) {
	code := strings.ToUpper(strings.TrimSpace(raw))
	if code == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(code) > MaxCodeCharacters {
		return nil, errors.New("code must be 40 characters or fewer")
	}
	for _, char := range code {
		if (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-' || char == '_' {
			continue
		}
		return nil, errors.New("code may only contain letters, numbers, hyphens, and underscores")
	}
	return &code, nil
}

func NormalizeDescription(raw string) *string {
	description := strings.TrimSpace(raw)
	if description == "" {
		return nil
	}
	return &description
}
