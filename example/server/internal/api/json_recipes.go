package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ai.ro/syncra/internal/auth"
	ocrsvc "ai.ro/syncra/internal/ocr"
)

// upsertJSONRecipeRequest contains the fields used to create or update a JSON recipe.
//
// swagger:model upsertJSONRecipeRequest
type upsertJSONRecipeRequest struct {
	// Recipe title.
	// required: true
	Title string `json:"title"`
	// Recipe description.
	Description string `json:"description"`
	// Valid JSON Schema object to clone into user extraction schemas.
	// required: true
	// type: object
	JSON json.RawMessage `json:"json" swaggertype:"object"`
	// Optional category id. Null or omitted places the recipe under Others.
	CategoryID *string `json:"category_id"`
}

// deployJSONRecipeRequest contains the target user id for a recipe deploy.
//
// swagger:model deployJSONRecipeRequest
type deployJSONRecipeRequest struct {
	// Owner user id for the cloned extraction schema.
	// required: true
	UserID string `json:"user_id" format:"uuid"`
}

type validatedJSONRecipeFields struct {
	Title       string
	Description string
	JSON        json.RawMessage
	CategoryID  *uuid.UUID
}

type jsonRecipeListCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
	Sort      string    `json:"sort"`
}

func validateJSONRecipeFields(title string, description string, raw json.RawMessage, rawCategoryID *string) (validatedJSONRecipeFields, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return validatedJSONRecipeFields{}, errors.New("title is required")
	}
	if utf8.RuneCountInString(title) > maxSchemaNameCharacters {
		return validatedJSONRecipeFields{}, errors.New("title must be at most 160 characters")
	}
	if !isJSONObject(raw) {
		return validatedJSONRecipeFields{}, errors.New("json must be a JSON object")
	}
	if int64(len(raw)) > maxSchemaJSONBytes {
		return validatedJSONRecipeFields{}, errors.New("json is too large")
	}
	if err := validateJSONSchema(raw); err != nil {
		return validatedJSONRecipeFields{}, errors.New("json must be a valid JSON Schema")
	}
	categoryID, err := parseOptionalJSONRecipeCategoryID(rawCategoryID)
	if err != nil {
		return validatedJSONRecipeFields{}, err
	}
	return validatedJSONRecipeFields{
		Title:       title,
		Description: description,
		JSON:        raw,
		CategoryID:  categoryID,
	}, nil
}

func parseOptionalJSONRecipeCategoryID(raw *string) (*uuid.UUID, error) {
	if raw == nil {
		return nil, nil
	}
	id, err := uuid.Parse(strings.TrimSpace(*raw))
	if err != nil || id == uuid.Nil {
		return nil, errors.New("invalid category id")
	}
	return &id, nil
}

func parseJSONRecipeID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, errors.New("invalid recipe id")
	}
	return id, nil
}

func (h *Handler) CreateJSONRecipe(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAdminRequestBytes)
	var req upsertJSONRecipeRequest
	if !decodeStrictAdminJSON(c, &req, "invalid JSON recipe payload") {
		return
	}
	fields, err := validateJSONRecipeFields(req.Title, req.Description, req.JSON, req.CategoryID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validateJSONRecipeCategoryExists(c, fields.CategoryID); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	recipe := ocrsvc.JSONRecipe{
		Title:       fields.Title,
		Description: fields.Description,
		JSON:        datatypes.JSON(fields.JSON),
		CategoryID:  fields.CategoryID,
	}
	if err := h.DB.WithContext(c.Request.Context()).Create(&recipe).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to save JSON recipe")
		return
	}
	if err := h.DB.WithContext(c.Request.Context()).Preload("Category").First(&recipe, "id = ?", recipe.ID).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to load JSON recipe")
		return
	}

	c.JSON(http.StatusCreated, jsonRecipeResponse(recipe))
}

func (h *Handler) ListJSONRecipes(c *gin.Context) {
	sortDirection, err := parseOCRJobListSort(c.Query("sort"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	size, err := parseOCRJobListSize(c.Query("size"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	cursor, err := parseJSONRecipeListCursor(c.Query("cursor"), sortDirection)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	query := h.DB.WithContext(c.Request.Context()).Model(&ocrsvc.JSONRecipe{}).Preload("Category")
	if cursor != nil {
		operator := "<"
		if sortDirection == "asc" {
			operator = ">"
		}
		query = query.Where("(created_at, id) "+operator+" (?, ?)", cursor.CreatedAt, cursor.ID)
	}

	order := "created_at desc, id desc"
	if sortDirection == "asc" {
		order = "created_at asc, id asc"
	}
	var recipes []ocrsvc.JSONRecipe
	if err := query.Order(order).Limit(size + 1).Find(&recipes).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list JSON recipes")
		return
	}

	var nextCursor *string
	if len(recipes) > size {
		recipes = recipes[:size]
		if len(recipes) > 0 {
			encoded, err := encodeJSONRecipeListCursor(recipes[len(recipes)-1], sortDirection)
			if err != nil {
				writeError(c, http.StatusInternalServerError, "failed to encode next cursor")
				return
			}
			nextCursor = &encoded
		}
	}

	out := make([]JSONRecipeResponse, 0, len(recipes))
	for _, recipe := range recipes {
		out = append(out, jsonRecipeResponse(recipe))
	}
	c.JSON(http.StatusOK, JSONRecipeListResponse{Recipes: out, NextCursor: nextCursor})
}

func parseJSONRecipeListCursor(raw string, sortDirection string) (*jsonRecipeListCursor, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.New("invalid cursor")
	}
	var cursor jsonRecipeListCursor
	if err := json.Unmarshal(decoded, &cursor); err != nil {
		return nil, errors.New("invalid cursor")
	}
	if cursor.ID == uuid.Nil || cursor.CreatedAt.IsZero() || (cursor.Sort != "asc" && cursor.Sort != "desc") {
		return nil, errors.New("invalid cursor")
	}
	if cursor.Sort != sortDirection {
		return nil, errors.New("cursor sort does not match sort")
	}
	return &cursor, nil
}

func encodeJSONRecipeListCursor(recipe ocrsvc.JSONRecipe, sortDirection string) (string, error) {
	raw, err := json.Marshal(jsonRecipeListCursor{
		CreatedAt: recipe.CreatedAt.UTC(),
		ID:        recipe.ID,
		Sort:      sortDirection,
	})
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func (h *Handler) GetJSONRecipe(c *gin.Context) {
	id, err := parseJSONRecipeID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var recipe ocrsvc.JSONRecipe
	if err := h.DB.WithContext(c.Request.Context()).Preload("Category").First(&recipe, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "JSON recipe not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load JSON recipe")
		return
	}
	c.JSON(http.StatusOK, jsonRecipeResponse(recipe))
}

func (h *Handler) UpdateJSONRecipe(c *gin.Context) {
	id, err := parseJSONRecipeID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAdminRequestBytes)
	var req upsertJSONRecipeRequest
	if !decodeStrictAdminJSON(c, &req, "invalid JSON recipe payload") {
		return
	}
	fields, err := validateJSONRecipeFields(req.Title, req.Description, req.JSON, req.CategoryID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validateJSONRecipeCategoryExists(c, fields.CategoryID); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var recipe ocrsvc.JSONRecipe
	if err := h.DB.WithContext(c.Request.Context()).First(&recipe, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "JSON recipe not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load JSON recipe")
		return
	}

	recipe.Title = fields.Title
	recipe.Description = fields.Description
	recipe.JSON = datatypes.JSON(fields.JSON)
	recipe.CategoryID = fields.CategoryID
	if err := h.DB.WithContext(c.Request.Context()).Save(&recipe).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to update JSON recipe")
		return
	}
	if err := h.DB.WithContext(c.Request.Context()).Preload("Category").First(&recipe, "id = ?", recipe.ID).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to load JSON recipe")
		return
	}

	c.JSON(http.StatusOK, jsonRecipeResponse(recipe))
}

func (h *Handler) DeleteJSONRecipe(c *gin.Context) {
	id, err := parseJSONRecipeID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	result := h.DB.WithContext(c.Request.Context()).Delete(&ocrsvc.JSONRecipe{}, "id = ?", id)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to delete JSON recipe")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "JSON recipe not found")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeployJSONRecipe(c *gin.Context) {
	id, err := parseJSONRecipeID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAdminRequestBytes)
	var req deployJSONRecipeRequest
	if !decodeStrictAdminJSON(c, &req, "invalid JSON recipe deploy payload") {
		return
	}
	userID, err := parseRequiredUserID(req.UserID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var recipe ocrsvc.JSONRecipe
	var schema ocrsvc.ExtractionSchema
	if err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Category").First(&recipe, "id = ?", id).Error; err != nil {
			return err
		}

		var count int64
		if err := tx.Model(&auth.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return errInvalidUserID
		}

		schema = ocrsvc.ExtractionSchema{
			UserID:      &userID,
			Name:        recipe.Title,
			Description: recipe.Description,
			SchemaJSON:  datatypes.JSON(recipe.JSON),
			Strict:      true,
		}
		if err := tx.Create(&schema).Error; err != nil {
			return err
		}

		recipe.Counter++
		if err := tx.Save(&recipe).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "JSON recipe not found")
			return
		}
		if errors.Is(err, errInvalidUserID) {
			writeError(c, http.StatusBadRequest, errInvalidUserID.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to deploy JSON recipe")
		return
	}

	c.JSON(http.StatusCreated, JSONRecipeDeployResponse{
		Recipe: jsonRecipeResponse(recipe),
		Schema: schemaResponse(schema),
	})
}
