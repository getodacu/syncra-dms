package ocr

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
)

type OCRJobStatus string

const (
	OCRJobStatusQueued     OCRJobStatus = "queued"
	OCRJobStatusProcessing OCRJobStatus = "processing"
	OCRJobStatusCompleted  OCRJobStatus = "completed"
	OCRJobStatusFailed     OCRJobStatus = "failed"
)

type OCRJob struct {
	ID               uuid.UUID         `gorm:"type:uuid;primaryKey;index:idx_ocr_jobs_user_created_id,priority:3;index:idx_ocr_jobs_user_status_created_id,priority:4" json:"id"`
	UserID           *string           `gorm:"column:user_id;type:uuid;index;index:idx_ocr_jobs_user_created_id,priority:1;index:idx_ocr_jobs_user_status_created_id,priority:1" json:"user_id,omitempty"`
	User             *auth.User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	CreatedAt        time.Time         `gorm:"index:idx_ocr_jobs_user_created_id,priority:2;index:idx_ocr_jobs_user_status_created_id,priority:3" json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	DeletedAt        gorm.DeletedAt    `gorm:"index" json:"-"`
	OriginalFilename string            `gorm:"not null;size:255" json:"original_filename"`
	MimeType         string            `gorm:"not null;size:120" json:"mime_type"`
	FileSize         int64             `gorm:"not null" json:"file_size"`
	PageCount        int               `gorm:"not null;default:0" json:"page_count"`
	DocumentHash     string            `gorm:"not null;size:64;index" json:"document_hash"`
	FilePath         string            `gorm:"not null;type:text" json:"-"`
	SchemaID         *uuid.UUID        `gorm:"type:uuid;index" json:"schema_id,omitempty"`
	Schema           *ExtractionSchema `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	InlineSchemaJSON datatypes.JSON    `gorm:"type:jsonb" json:"inline_schema,omitempty"`
	DocumentID       *uuid.UUID        `gorm:"type:uuid;index" json:"document_id,omitempty"`
	Document         *OCRDocument      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Status           OCRJobStatus      `gorm:"not null;size:40;index;index:idx_ocr_jobs_user_status_created_id,priority:2;default:queued" json:"status"`
	ErrorMessage     string            `gorm:"type:text" json:"error_message,omitempty"`
}

func (job *OCRJob) BeforeCreate(_ *gorm.DB) error {
	if job.ID == uuid.Nil {
		job.ID = uuid.New()
	}
	return validateOCRJobStatus(job.Status)
}

func (job *OCRJob) BeforeUpdate(tx *gorm.DB) error {
	status, ok := ocrJobStatusFromUpdate(tx)
	if !ok {
		return nil
	}
	return validateOCRJobStatus(status)
}

func validateOCRJobStatus(status OCRJobStatus) error {
	switch status {
	case OCRJobStatusQueued, OCRJobStatusProcessing, OCRJobStatusCompleted, OCRJobStatusFailed:
		return nil
	default:
		return fmt.Errorf("invalid OCR job status %q", status)
	}
}

func ocrJobStatusFromUpdate(tx *gorm.DB) (OCRJobStatus, bool) {
	if tx == nil || tx.Statement == nil {
		return "", false
	}

	dest := reflect.ValueOf(tx.Statement.Dest)
	for dest.IsValid() && dest.Kind() == reflect.Ptr {
		if dest.IsNil() {
			return "", false
		}
		dest = dest.Elem()
	}
	if !dest.IsValid() {
		return "", false
	}

	switch dest.Kind() {
	case reflect.Map:
		return ocrJobStatusFromUpdateMap(dest)
	case reflect.Struct:
		statusField := dest.FieldByName("Status")
		if !ocrJobStatusStructUpdateWillWrite(tx, statusField) &&
			!tx.Statement.Changed("Status") &&
			!tx.Statement.Changed("status") {
			return "", false
		}
		return ocrJobStatusFromValue(statusField)
	default:
		return "", false
	}
}

func ocrJobStatusStructUpdateWillWrite(tx *gorm.DB, statusField reflect.Value) bool {
	if tx == nil || tx.Statement == nil {
		return false
	}

	selectColumns, restricted := tx.Statement.SelectAndOmitColumns(false, true)
	if selected, ok := selectColumns["status"]; ok {
		return selected
	}
	if restricted {
		return false
	}
	return !ocrJobStatusValueIsZero(statusField)
}

func ocrJobStatusFromUpdateMap(updateMap reflect.Value) (OCRJobStatus, bool) {
	for _, key := range updateMap.MapKeys() {
		if key.Kind() != reflect.String {
			continue
		}
		switch key.String() {
		case "Status", "status":
			return ocrJobStatusFromValue(updateMap.MapIndex(key))
		}
	}
	return "", false
}

func ocrJobStatusFromValue(value reflect.Value) (OCRJobStatus, bool) {
	if !value.IsValid() {
		return "", false
	}
	if value.Kind() == reflect.Interface {
		if value.IsNil() {
			return "", true
		}
		value = value.Elem()
	}
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return "", true
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.String {
		return OCRJobStatus(fmt.Sprint(value.Interface())), true
	}
	return OCRJobStatus(value.String()), true
}

func ocrJobStatusValueIsZero(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}
	if value.Kind() == reflect.Interface {
		if value.IsNil() {
			return true
		}
		value = value.Elem()
	}
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return true
		}
		value = value.Elem()
	}
	return value.IsZero()
}

func (OCRJob) TableName() string {
	return "ocr_jobs"
}
