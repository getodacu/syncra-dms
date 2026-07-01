package documents

import (
	"errors"
	"path"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const MaxFolderNameCharacters = 160
const MaxDocumentDisplayNameCharacters = 255

type Folder struct {
	ID                 string        `gorm:"type:uuid;primaryKey;index:idx_document_folders_parent_name_id,priority:4;uniqueIndex:idx_document_folders_id_organization_unit_unique,priority:1" json:"id"`
	ParentID           *string       `gorm:"column:parent_id;type:uuid;index:idx_document_folders_parent_name_id,priority:2;uniqueIndex:idx_document_folders_child_name_unique,priority:2,where:parent_id IS NOT NULL AND deleted_at IS NULL" json:"parentId,omitempty"`
	Parent             *Folder       `gorm:"foreignKey:ParentID,OrganizationUnitID;references:ID,OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	OrganizationUnitID string        `gorm:"column:organization_unit_id;type:uuid;not null;index:idx_document_folders_parent_name_id,priority:1;uniqueIndex:idx_document_folders_root_name_unique,priority:1,where:parent_id IS NULL AND deleted_at IS NULL;uniqueIndex:idx_document_folders_child_name_unique,priority:1,where:parent_id IS NOT NULL AND deleted_at IS NULL;uniqueIndex:idx_document_folders_id_organization_unit_unique,priority:2" json:"organizationUnitId"`
	OrganizationUnit   orgunits.Unit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Name               string        `gorm:"not null;size:160;index:idx_document_folders_parent_name_id,priority:3;uniqueIndex:idx_document_folders_root_name_unique,priority:2,where:parent_id IS NULL AND deleted_at IS NULL;uniqueIndex:idx_document_folders_child_name_unique,priority:3,where:parent_id IS NOT NULL AND deleted_at IS NULL" json:"name"`
	Description        *string       `gorm:"type:text" json:"description,omitempty"`
	CreatedByUserID    string        `gorm:"column:created_by_user_id;type:uuid;not null;index" json:"createdByUserId"`
	CreatedByUser      auth.User     `gorm:"foreignKey:CreatedByUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	UpdatedByUserID    *string       `gorm:"column:updated_by_user_id;type:uuid;index" json:"updatedByUserId,omitempty"`
	UpdatedByUser      *auth.User    `gorm:"foreignKey:UpdatedByUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	DeletedAt          *time.Time    `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
	CreatedAt          time.Time     `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt          time.Time     `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (f *Folder) BeforeCreate(_ *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	return nil
}

func (Folder) TableName() string { return "document_folders" }

type Document struct {
	ID                 string        `gorm:"type:uuid;primaryKey" json:"id"`
	FolderID           string        `gorm:"column:folder_id;type:uuid;not null;index;uniqueIndex:idx_documents_active_folder_hash_unique,priority:1,where:deleted_at IS NULL" json:"folderId"`
	Folder             Folder        `gorm:"foreignKey:FolderID,OrganizationUnitID;references:ID,OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	OrganizationUnitID string        `gorm:"column:organization_unit_id;type:uuid;not null;index" json:"organizationUnitId"`
	OrganizationUnit   orgunits.Unit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	OriginalFileName   string        `gorm:"column:original_file_name;not null;size:255" json:"originalFileName"`
	DisplayName        string        `gorm:"column:display_name;not null;size:255" json:"displayName"`
	MimeType           string        `gorm:"column:mime_type;not null;size:255" json:"mimeType"`
	Extension          *string       `gorm:"size:32" json:"extension,omitempty"`
	SizeBytes          int64         `gorm:"column:size_bytes;not null;check:chk_documents_size_bytes_non_negative,size_bytes >= 0" json:"sizeBytes"`
	SHA256Hash         string        `gorm:"column:sha256_hash;type:char(64);not null;check:chk_documents_sha256_hash_lower_hex,length(sha256_hash) = 64 AND replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(sha256_hash,'0',''),'1',''),'2',''),'3',''),'4',''),'5',''),'6',''),'7',''),'8',''),'9',''),'a',''),'b',''),'c',''),'d',''),'e',''),'f','') = '';uniqueIndex:idx_documents_active_folder_hash_unique,priority:2,where:deleted_at IS NULL" json:"sha256Hash"`
	StorageKey         string        `gorm:"column:storage_key;not null;type:text" json:"-"`
	CreatedByUserID    string        `gorm:"column:created_by_user_id;type:uuid;not null;index" json:"createdByUserId"`
	CreatedByUser      auth.User     `gorm:"foreignKey:CreatedByUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	UpdatedByUserID    *string       `gorm:"column:updated_by_user_id;type:uuid;index" json:"updatedByUserId,omitempty"`
	UpdatedByUser      *auth.User    `gorm:"foreignKey:UpdatedByUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	DeletedAt          *time.Time    `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
	CreatedAt          time.Time     `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt          time.Time     `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (d *Document) BeforeCreate(_ *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	return nil
}

func (Document) TableName() string { return "documents" }

func NormalizeFolderName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", errors.New("folder name is required")
	}
	if utf8.RuneCountInString(name) > MaxFolderNameCharacters {
		return "", errors.New("folder name must be 160 characters or fewer")
	}
	return name, nil
}

func NormalizeDescription(raw string) *string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil
	}
	return &value
}

func NormalizeDisplayName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", errors.New("document name is required")
	}
	if utf8.RuneCountInString(name) > MaxDocumentDisplayNameCharacters {
		return "", errors.New("document name must be 255 characters or fewer")
	}
	return name, nil
}

func SafeOriginalFileName(raw string) string {
	normalizedPath := strings.ReplaceAll(raw, "\\", "/")
	name := strings.TrimSpace(path.Base(normalizedPath))
	name = strings.Map(func(char rune) rune {
		if unicode.IsControl(char) {
			return -1
		}
		return char
	}, name)
	name = strings.TrimSpace(name)
	if name == "." || name == "/" || name == "" {
		return "upload"
	}
	if utf8.RuneCountInString(name) > 255 {
		return string([]rune(name)[:255])
	}
	return name
}
