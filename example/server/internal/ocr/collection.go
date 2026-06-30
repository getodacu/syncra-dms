package ocr

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ai.ro/syncra/internal/auth"
)

type Collection struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;index:idx_collections_user_created_id,priority:3" json:"id"`
	UserID    string     `gorm:"column:user_id;type:uuid;not null;index;index:idx_collections_user_created_id,priority:1" json:"user_id"`
	User      *auth.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	CreatedAt time.Time  `gorm:"index:idx_collections_user_created_id,priority:2" json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Name      string     `gorm:"not null;size:160" json:"name"`
}

func (collection *Collection) BeforeCreate(_ *gorm.DB) error {
	if collection.ID == uuid.Nil {
		collection.ID = uuid.New()
	}
	return nil
}

func (Collection) TableName() string {
	return "collections"
}

type CollectionSchema struct {
	CollectionID uuid.UUID        `gorm:"type:uuid;primaryKey;index" json:"collection_id"`
	Collection   Collection       `gorm:"foreignKey:CollectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	SchemaID     uuid.UUID        `gorm:"type:uuid;primaryKey;index" json:"schema_id"`
	Schema       ExtractionSchema `gorm:"foreignKey:SchemaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	CreatedAt    time.Time        `json:"created_at"`
}

func (CollectionSchema) TableName() string {
	return "collection_schemas"
}

type CollectionDocument struct {
	CollectionID uuid.UUID   `gorm:"type:uuid;primaryKey;index" json:"collection_id"`
	Collection   Collection  `gorm:"foreignKey:CollectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	DocumentID   uuid.UUID   `gorm:"type:uuid;primaryKey;index" json:"document_id"`
	Document     OCRDocument `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	CreatedAt    time.Time   `json:"created_at"`
}

func (CollectionDocument) TableName() string {
	return "collection_documents"
}

func LinkDocumentToMatchingCollections(ctx context.Context, db *gorm.DB, doc OCRDocument) error {
	if db == nil {
		return errors.New("database is required")
	}
	return LinkDocumentToMatchingCollectionsForSchema(ctx, db, doc.ID, doc.UserID, doc.SchemaID)
}

func LinkDocumentToMatchingCollectionsForSchema(ctx context.Context, db *gorm.DB, documentID uuid.UUID, userID *string, schemaID *uuid.UUID) error {
	if db == nil {
		return errors.New("database is required")
	}
	if documentID == uuid.Nil || schemaID == nil || userID == nil {
		return nil
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var doc OCRDocument
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id").
			Where("id = ? AND user_id = ?", documentID, *userID).
			First(&doc).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		return tx.Exec(`
INSERT INTO collection_documents (collection_id, document_id, created_at)
SELECT cs.collection_id, ?, now()
FROM collection_schemas cs
JOIN collections c ON c.id = cs.collection_id
WHERE cs.schema_id = ? AND c.user_id = ?
ON CONFLICT DO NOTHING
`, documentID, *schemaID, *userID).Error
	})
}
