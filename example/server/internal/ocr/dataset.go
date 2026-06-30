package ocr

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
)

type Dataset struct {
	ID             uuid.UUID        `gorm:"type:uuid;primaryKey;index:idx_datasets_user_created_id,priority:3" json:"id"`
	UserID         string           `gorm:"column:user_id;type:uuid;not null;index;index:idx_datasets_user_created_id,priority:1" json:"user_id"`
	User           *auth.User       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	SchemaID       uuid.UUID        `gorm:"type:uuid;not null;index" json:"schema_id"`
	Schema         ExtractionSchema `gorm:"foreignKey:SchemaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	CreatedAt      time.Time        `gorm:"index:idx_datasets_user_created_id,priority:2" json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	Name           string           `gorm:"not null;size:160" json:"name"`
	SelectedFields datatypes.JSON   `gorm:"type:jsonb;not null" json:"selected_fields"`
}

func (dataset *Dataset) BeforeCreate(_ *gorm.DB) error {
	if dataset.ID == uuid.Nil {
		dataset.ID = uuid.New()
	}
	return nil
}

func (Dataset) TableName() string {
	return "datasets"
}
