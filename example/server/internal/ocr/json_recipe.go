package ocr

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type JSONRecipe struct {
	ID          uuid.UUID           `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Title       string              `gorm:"not null;size:160" json:"title"`
	Description string              `gorm:"type:text" json:"description"`
	JSON        datatypes.JSON      `gorm:"column:json;type:jsonb;not null" json:"json"`
	Counter     int64               `gorm:"not null;default:0;check:chk_json_recipes_counter,counter >= 0" json:"counter"`
	CategoryID  *uuid.UUID          `gorm:"type:uuid;index" json:"category_id"`
	Category    *JSONRecipeCategory `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"category"`
}

func (recipe *JSONRecipe) BeforeCreate(_ *gorm.DB) error {
	if recipe.ID == uuid.Nil {
		recipe.ID = uuid.New()
	}
	return nil
}

func (JSONRecipe) TableName() string {
	return "json_recipes"
}

type JSONRecipeCategory struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TitleEn   string    `gorm:"column:title_en;not null;size:160" json:"title_en"`
	TitleRo   string    `gorm:"column:title_ro;not null;size:160" json:"title_ro"`
}

func (category *JSONRecipeCategory) BeforeCreate(_ *gorm.DB) error {
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	return nil
}

func (JSONRecipeCategory) TableName() string {
	return "json_recipe_categories"
}
