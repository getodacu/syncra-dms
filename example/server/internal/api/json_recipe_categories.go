package api

import (
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	ocrsvc "ai.ro/syncra/internal/ocr"
)

// upsertJSONRecipeCategoryTitleRequest contains localized category titles.
//
// swagger:model upsertJSONRecipeCategoryTitleRequest
type upsertJSONRecipeCategoryTitleRequest struct {
	// English category title.
	// required: true
	En string `json:"en"`
	// Romanian category title.
	// required: true
	Ro string `json:"ro"`
}

// upsertJSONRecipeCategoryRequest contains fields used to create or update a JSON recipe category.
//
// swagger:model upsertJSONRecipeCategoryRequest
type upsertJSONRecipeCategoryRequest struct {
	// Localized category titles.
	// required: true
	Title upsertJSONRecipeCategoryTitleRequest `json:"title"`
}

type validatedJSONRecipeCategoryFields struct {
	TitleEn string
	TitleRo string
}

func validateJSONRecipeCategoryFields(req upsertJSONRecipeCategoryRequest) (validatedJSONRecipeCategoryFields, error) {
	titleEn := strings.TrimSpace(req.Title.En)
	titleRo := strings.TrimSpace(req.Title.Ro)
	if titleEn == "" {
		return validatedJSONRecipeCategoryFields{}, errors.New("english title is required")
	}
	if titleRo == "" {
		return validatedJSONRecipeCategoryFields{}, errors.New("romanian title is required")
	}
	if utf8.RuneCountInString(titleEn) > maxSchemaNameCharacters || utf8.RuneCountInString(titleRo) > maxSchemaNameCharacters {
		return validatedJSONRecipeCategoryFields{}, errors.New("titles must be at most 160 characters")
	}
	return validatedJSONRecipeCategoryFields{TitleEn: titleEn, TitleRo: titleRo}, nil
}

func parseJSONRecipeCategoryID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, errors.New("invalid category id")
	}
	return id, nil
}

func (h *Handler) validateJSONRecipeCategoryExists(c *gin.Context, categoryID *uuid.UUID) error {
	if categoryID == nil {
		return nil
	}
	var count int64
	if err := h.DB.WithContext(c.Request.Context()).Model(&ocrsvc.JSONRecipeCategory{}).Where("id = ?", *categoryID).Count(&count).Error; err != nil {
		return errors.New("failed to validate category")
	}
	if count == 0 {
		return errors.New("category_id is invalid")
	}
	return nil
}

func (h *Handler) ListJSONRecipeCategories(c *gin.Context) {
	var categories []ocrsvc.JSONRecipeCategory
	if err := h.DB.WithContext(c.Request.Context()).
		Order("title_en asc, title_ro asc, id asc").
		Find(&categories).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list JSON recipe categories")
		return
	}

	out := make([]JSONRecipeCategoryResponse, 0, len(categories))
	for _, category := range categories {
		out = append(out, jsonRecipeCategoryResponse(category))
	}
	c.JSON(http.StatusOK, JSONRecipeCategoryListResponse{Categories: out})
}

func (h *Handler) CreateJSONRecipeCategory(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAdminRequestBytes)
	var req upsertJSONRecipeCategoryRequest
	if !decodeStrictAdminJSON(c, &req, "invalid JSON recipe category payload") {
		return
	}
	fields, err := validateJSONRecipeCategoryFields(req)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	category := ocrsvc.JSONRecipeCategory{
		TitleEn: fields.TitleEn,
		TitleRo: fields.TitleRo,
	}
	if err := h.DB.WithContext(c.Request.Context()).Create(&category).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to save JSON recipe category")
		return
	}

	c.JSON(http.StatusCreated, jsonRecipeCategoryResponse(category))
}

func (h *Handler) GetJSONRecipeCategory(c *gin.Context) {
	id, err := parseJSONRecipeCategoryID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var category ocrsvc.JSONRecipeCategory
	if err := h.DB.WithContext(c.Request.Context()).First(&category, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "JSON recipe category not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load JSON recipe category")
		return
	}
	c.JSON(http.StatusOK, jsonRecipeCategoryResponse(category))
}

func (h *Handler) UpdateJSONRecipeCategory(c *gin.Context) {
	id, err := parseJSONRecipeCategoryID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAdminRequestBytes)
	var req upsertJSONRecipeCategoryRequest
	if !decodeStrictAdminJSON(c, &req, "invalid JSON recipe category payload") {
		return
	}
	fields, err := validateJSONRecipeCategoryFields(req)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var category ocrsvc.JSONRecipeCategory
	if err := h.DB.WithContext(c.Request.Context()).First(&category, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "JSON recipe category not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load JSON recipe category")
		return
	}

	category.TitleEn = fields.TitleEn
	category.TitleRo = fields.TitleRo
	if err := h.DB.WithContext(c.Request.Context()).Save(&category).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to update JSON recipe category")
		return
	}
	c.JSON(http.StatusOK, jsonRecipeCategoryResponse(category))
}

func (h *Handler) DeleteJSONRecipeCategory(c *gin.Context) {
	id, err := parseJSONRecipeCategoryID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var category ocrsvc.JSONRecipeCategory
	if err := h.DB.WithContext(c.Request.Context()).First(&category, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "JSON recipe category not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load JSON recipe category")
		return
	}

	var assigned int64
	if err := h.DB.WithContext(c.Request.Context()).Model(&ocrsvc.JSONRecipe{}).Where("category_id = ?", id).Count(&assigned).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to inspect JSON recipe category")
		return
	}
	if assigned > 0 {
		writeError(c, http.StatusConflict, "category has assigned JSON recipes")
		return
	}

	if err := h.DB.WithContext(c.Request.Context()).Delete(&category).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to delete JSON recipe category")
		return
	}
	c.Status(http.StatusNoContent)
}
