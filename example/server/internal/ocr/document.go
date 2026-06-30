package ocr

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
)

type ExtractionSchema struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      *string        `gorm:"column:user_id;type:uuid;index" json:"user_id,omitempty"`
	User        *auth.User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Name        string         `gorm:"not null;size:160" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	SchemaJSON  datatypes.JSON `gorm:"type:jsonb;not null" json:"schema"`
	Strict      bool           `gorm:"not null;default:true" json:"strict"`
}

func (schema *ExtractionSchema) BeforeCreate(_ *gorm.DB) error {
	if schema.ID == uuid.Nil {
		schema.ID = uuid.New()
	}
	return nil
}

func (ExtractionSchema) TableName() string {
	return "extraction_schemas"
}

type OCRDocument struct {
	ID               uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           *string           `gorm:"column:user_id;type:uuid;index" json:"user_id,omitempty"`
	User             *auth.User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	DeletedAt        gorm.DeletedAt    `gorm:"index" json:"-"`
	OriginalFilename string            `gorm:"not null;size:255" json:"original_filename"`
	MimeType         string            `gorm:"not null;size:120" json:"mime_type"`
	FileSize         int64             `gorm:"not null" json:"file_size"`
	PageCount        int               `gorm:"not null;default:0" json:"page_count"`
	DocumentHash     string            `gorm:"not null;size:64;index" json:"document_hash"`
	JobID            *uuid.UUID        `gorm:"type:uuid;index" json:"job_id,omitempty"`
	Job              *OCRJob           `gorm:"foreignKey:JobID;-:migration;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	SchemaID         *uuid.UUID        `gorm:"type:uuid;index" json:"schema_id,omitempty"`
	Schema           *ExtractionSchema `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	InlineSchemaJSON datatypes.JSON    `gorm:"type:jsonb" json:"inline_schema,omitempty"`
	Markdown         string            `gorm:"type:text;not null" json:"markdown"`
	AnnotationJSON   datatypes.JSON    `gorm:"type:jsonb" json:"annotation_json,omitempty"`
	RawResponseJSON  datatypes.JSON    `gorm:"type:jsonb;not null" json:"raw_response_json"`
}

func (doc *OCRDocument) BeforeCreate(_ *gorm.DB) error {
	if doc.ID == uuid.Nil {
		doc.ID = uuid.New()
	}
	pageCount, err := CountRawResponsePages(doc.RawResponseJSON)
	if err != nil {
		return err
	}
	doc.PageCount = pageCount
	return nil
}

func (OCRDocument) TableName() string {
	return "ocr_documents"
}
