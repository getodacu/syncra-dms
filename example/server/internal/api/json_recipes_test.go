package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/ocr"
)

const validRecipeSchema = `{"type":"object","properties":{"number":{"type":"string"}}}`

type jsonRecipeCategoryTitleTestResponse struct {
	En string `json:"en"`
	Ro string `json:"ro"`
}

type jsonRecipeCategoryTestResponse struct {
	ID        uuid.UUID                           `json:"id"`
	Title     jsonRecipeCategoryTitleTestResponse `json:"title"`
	CreatedAt time.Time                           `json:"created_at"`
	UpdatedAt time.Time                           `json:"updated_at"`
}

type jsonRecipeCategoryListTestResponse struct {
	Categories []jsonRecipeCategoryTestResponse `json:"categories"`
}

type jsonRecipeTestResponse struct {
	ID          uuid.UUID                       `json:"id"`
	Title       string                          `json:"title"`
	Description string                          `json:"description"`
	JSON        json.RawMessage                 `json:"json"`
	Counter     int64                           `json:"counter"`
	CategoryID  *uuid.UUID                      `json:"category_id"`
	Category    *jsonRecipeCategoryTestResponse `json:"category"`
	CreatedAt   time.Time                       `json:"created_at"`
	UpdatedAt   time.Time                       `json:"updated_at"`
}

type jsonRecipeListTestResponse struct {
	Recipes    []jsonRecipeTestResponse `json:"recipes"`
	NextCursor *string                  `json:"next_cursor"`
}

type jsonRecipeDeployTestResponse struct {
	Recipe jsonRecipeTestResponse `json:"recipe"`
	Schema SchemaResponse         `json:"schema"`
}

func createTestJSONRecipe(t *testing.T, db *gorm.DB, recipe ocr.JSONRecipe) ocr.JSONRecipe {
	t.Helper()
	if recipe.Title == "" {
		recipe.Title = "Invoice"
	}
	if len(recipe.JSON) == 0 {
		recipe.JSON = datatypes.JSON([]byte(validRecipeSchema))
	}
	if err := db.Create(&recipe).Error; err != nil {
		t.Fatalf("create JSON recipe: %v", err)
	}
	return recipe
}

func createTestJSONRecipeCategory(t *testing.T, db *gorm.DB, category ocr.JSONRecipeCategory) ocr.JSONRecipeCategory {
	t.Helper()
	if category.TitleEn == "" {
		category.TitleEn = "Invoices"
	}
	if category.TitleRo == "" {
		category.TitleRo = "Facturi"
	}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("create JSON recipe category: %v", err)
	}
	return category
}

func TestAdminJSONRecipesRequireAdminSession(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-json-recipes-authz@example.com", "admin")
	normal := createAdminTestUser(t, db, "normal-json-recipes-authz@example.com", "user")
	adminCookie := createAdminTestSession(t, db, admin, "admin-json-recipes-authz-session")
	normalCookie := createAdminTestSession(t, db, normal, "normal-json-recipes-authz-session")

	noSession := adminJSON(t, router, http.MethodGet, "/api/admin/json-recipes", "", nil)
	if noSession.Code != http.StatusUnauthorized {
		t.Fatalf("no-session status = %d body=%s", noSession.Code, noSession.Body.String())
	}

	nonAdmin := adminJSON(t, router, http.MethodGet, "/api/admin/json-recipes", "", normalCookie)
	if nonAdmin.Code != http.StatusForbidden {
		t.Fatalf("non-admin status = %d body=%s", nonAdmin.Code, nonAdmin.Body.String())
	}

	ok := adminJSON(t, router, http.MethodGet, "/api/admin/json-recipes", "", adminCookie)
	if ok.Code != http.StatusOK {
		t.Fatalf("admin status = %d body=%s", ok.Code, ok.Body.String())
	}

	categoryNoSession := adminJSON(t, router, http.MethodGet, "/api/admin/json-recipe-categories", "", nil)
	if categoryNoSession.Code != http.StatusUnauthorized {
		t.Fatalf("category no-session status = %d body=%s", categoryNoSession.Code, categoryNoSession.Body.String())
	}

	categoryNonAdmin := adminJSON(t, router, http.MethodGet, "/api/admin/json-recipe-categories", "", normalCookie)
	if categoryNonAdmin.Code != http.StatusForbidden {
		t.Fatalf("category non-admin status = %d body=%s", categoryNonAdmin.Code, categoryNonAdmin.Body.String())
	}

	categoryOK := adminJSON(t, router, http.MethodGet, "/api/admin/json-recipe-categories", "", adminCookie)
	if categoryOK.Code != http.StatusOK {
		t.Fatalf("category admin status = %d body=%s", categoryOK.Code, categoryOK.Body.String())
	}
}

func TestAdminJSONRecipeCategoryCRUDValidationAndDeleteConflict(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-json-recipe-categories@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-json-recipe-categories-session")

	missingTitle := adminJSON(t, router, http.MethodPost, "/api/admin/json-recipe-categories", `{"title":{"en":"Invoices"}}`, cookie)
	if missingTitle.Code != http.StatusBadRequest {
		t.Fatalf("missing-title status = %d body=%s", missingTitle.Code, missingTitle.Body.String())
	}

	unknownField := adminJSON(t, router, http.MethodPost, "/api/admin/json-recipe-categories", `{"title":{"en":"Invoices","ro":"Facturi"},"slug":"invoice"}`, cookie)
	if unknownField.Code != http.StatusBadRequest {
		t.Fatalf("unknown-field status = %d body=%s", unknownField.Code, unknownField.Body.String())
	}

	create := adminJSON(t, router, http.MethodPost, "/api/admin/json-recipe-categories", `{"title":{"en":" Invoices ","ro":" Facturi "}}`, cookie)
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d body=%s", create.Code, create.Body.String())
	}
	created := decodeAdminResponse[jsonRecipeCategoryTestResponse](t, create)
	if created.ID == uuid.Nil || created.Title.En != "Invoices" || created.Title.Ro != "Facturi" {
		t.Fatalf("unexpected category: %#v", created)
	}

	get := adminJSON(t, router, http.MethodGet, "/api/admin/json-recipe-categories/"+created.ID.String(), "", cookie)
	if get.Code != http.StatusOK {
		t.Fatalf("get status = %d body=%s", get.Code, get.Body.String())
	}
	got := decodeAdminResponse[jsonRecipeCategoryTestResponse](t, get)
	if got.ID != created.ID || got.Title.En != "Invoices" || got.Title.Ro != "Facturi" {
		t.Fatalf("unexpected get response: %#v", got)
	}

	update := adminJSON(t, router, http.MethodPut, "/api/admin/json-recipe-categories/"+created.ID.String(), `{"title":{"en":" Receipts ","ro":" Bonuri "}}`, cookie)
	if update.Code != http.StatusOK {
		t.Fatalf("update status = %d body=%s", update.Code, update.Body.String())
	}
	updated := decodeAdminResponse[jsonRecipeCategoryTestResponse](t, update)
	if updated.Title.En != "Receipts" || updated.Title.Ro != "Bonuri" {
		t.Fatalf("unexpected update response: %#v", updated)
	}

	list := adminJSON(t, router, http.MethodGet, "/api/admin/json-recipe-categories", "", cookie)
	if list.Code != http.StatusOK {
		t.Fatalf("list status = %d body=%s", list.Code, list.Body.String())
	}
	categories := decodeAdminResponse[jsonRecipeCategoryListTestResponse](t, list)
	if len(categories.Categories) != 1 || categories.Categories[0].ID != created.ID {
		t.Fatalf("unexpected category list: %#v", categories)
	}

	categoryID := created.ID
	createTestJSONRecipe(t, db, ocr.JSONRecipe{Title: "Invoice recipe", CategoryID: &categoryID})
	conflict := adminJSON(t, router, http.MethodDelete, "/api/admin/json-recipe-categories/"+created.ID.String(), "", cookie)
	if conflict.Code != http.StatusConflict {
		t.Fatalf("conflict status = %d body=%s", conflict.Code, conflict.Body.String())
	}

	if err := db.Model(&ocr.JSONRecipe{}).Where("category_id = ?", categoryID).Update("category_id", nil).Error; err != nil {
		t.Fatalf("clear recipe category: %v", err)
	}
	deleteResponse := adminJSON(t, router, http.MethodDelete, "/api/admin/json-recipe-categories/"+created.ID.String(), "", cookie)
	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d body=%s", deleteResponse.Code, deleteResponse.Body.String())
	}
}

func TestAdminJSONRecipeCRUDValidationAndCounterOwnership(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-json-recipes-crud@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-json-recipes-crud-session")
	category := createTestJSONRecipeCategory(t, db, ocr.JSONRecipeCategory{TitleEn: "Invoices", TitleRo: "Facturi"})

	emptyTitle := adminJSON(t, router, http.MethodPost, "/api/admin/json-recipes", `{"title":"   ","description":"","json":{"type":"object"}}`, cookie)
	if emptyTitle.Code != http.StatusBadRequest {
		t.Fatalf("empty-title status = %d body=%s", emptyTitle.Code, emptyTitle.Body.String())
	}

	invalidSchema := adminJSON(t, router, http.MethodPost, "/api/admin/json-recipes", `{"title":"Invoice","description":"","json":{"type":5}}`, cookie)
	if invalidSchema.Code != http.StatusBadRequest {
		t.Fatalf("invalid-schema status = %d body=%s", invalidSchema.Code, invalidSchema.Body.String())
	}

	counterSet := adminJSON(t, router, http.MethodPost, "/api/admin/json-recipes", `{"title":"Invoice","description":"","json":{"type":"object"},"counter":9}`, cookie)
	if counterSet.Code != http.StatusBadRequest {
		t.Fatalf("counter-set status = %d body=%s", counterSet.Code, counterSet.Body.String())
	}

	invalidCategory := adminJSON(t, router, http.MethodPost, "/api/admin/json-recipes", `{"title":"Invoice","description":"","json":{"type":"object"},"category_id":"`+uuid.NewString()+`"}`, cookie)
	if invalidCategory.Code != http.StatusBadRequest {
		t.Fatalf("invalid-category status = %d body=%s", invalidCategory.Code, invalidCategory.Body.String())
	}

	create := adminJSON(t, router, http.MethodPost, "/api/admin/json-recipes", `{"title":" Invoice ","description":"Invoice fields","json":`+validRecipeSchema+`,"category_id":"`+category.ID.String()+`"}`, cookie)
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d body=%s", create.Code, create.Body.String())
	}
	created := decodeAdminResponse[jsonRecipeTestResponse](t, create)
	if created.ID == uuid.Nil || created.Title != "Invoice" || created.Description != "Invoice fields" || created.Counter != 0 {
		t.Fatalf("unexpected created recipe: %#v", created)
	}
	if created.CategoryID == nil || *created.CategoryID != category.ID || created.Category == nil || created.Category.Title.En != "Invoices" || created.Category.Title.Ro != "Facturi" {
		t.Fatalf("unexpected created category response: %#v", created)
	}
	assertJSONEqual(t, created.JSON, validRecipeSchema)

	get := adminJSON(t, router, http.MethodGet, "/api/admin/json-recipes/"+created.ID.String(), "", cookie)
	if get.Code != http.StatusOK {
		t.Fatalf("get status = %d body=%s", get.Code, get.Body.String())
	}
	got := decodeAdminResponse[jsonRecipeTestResponse](t, get)
	if got.ID != created.ID || got.Title != "Invoice" {
		t.Fatalf("unexpected get response: %#v", got)
	}
	if got.CategoryID == nil || *got.CategoryID != category.ID || got.Category == nil {
		t.Fatalf("unexpected get category response: %#v", got)
	}

	updateCounter := adminJSON(t, router, http.MethodPut, "/api/admin/json-recipes/"+created.ID.String(), `{"title":"Receipt","description":"","json":{"type":"object"},"counter":5}`, cookie)
	if updateCounter.Code != http.StatusBadRequest {
		t.Fatalf("update-counter status = %d body=%s", updateCounter.Code, updateCounter.Body.String())
	}

	update := adminJSON(t, router, http.MethodPut, "/api/admin/json-recipes/"+created.ID.String(), `{"title":" Receipt ","description":"Receipt fields","json":{"type":"object"},"category_id":null}`, cookie)
	if update.Code != http.StatusOK {
		t.Fatalf("update status = %d body=%s", update.Code, update.Body.String())
	}
	updated := decodeAdminResponse[jsonRecipeTestResponse](t, update)
	if updated.Title != "Receipt" || updated.Description != "Receipt fields" || updated.Counter != 0 {
		t.Fatalf("unexpected updated recipe: %#v", updated)
	}
	if updated.CategoryID != nil || updated.Category != nil {
		t.Fatalf("updated category = %v/%v, want nil", updated.CategoryID, updated.Category)
	}

	deleteResponse := adminJSON(t, router, http.MethodDelete, "/api/admin/json-recipes/"+created.ID.String(), "", cookie)
	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d body=%s", deleteResponse.Code, deleteResponse.Body.String())
	}

	missing := adminJSON(t, router, http.MethodGet, "/api/admin/json-recipes/"+created.ID.String(), "", cookie)
	if missing.Code != http.StatusNotFound {
		t.Fatalf("missing status = %d body=%s", missing.Code, missing.Body.String())
	}
}

func TestAdminJSONRecipesListPaginationAndSort(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-json-recipes-list@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-json-recipes-list-session")

	base := time.Date(2026, 6, 22, 10, 0, 0, 0, time.UTC)
	first := createTestJSONRecipe(t, db, ocr.JSONRecipe{Title: "First", CreatedAt: base, UpdatedAt: base})
	second := createTestJSONRecipe(t, db, ocr.JSONRecipe{Title: "Second", CreatedAt: base.Add(time.Minute), UpdatedAt: base.Add(time.Minute)})
	third := createTestJSONRecipe(t, db, ocr.JSONRecipe{Title: "Third", CreatedAt: base.Add(2 * time.Minute), UpdatedAt: base.Add(2 * time.Minute)})

	pageOne := adminJSON(t, router, http.MethodGet, "/api/admin/json-recipes?sort=asc&size=2", "", cookie)
	if pageOne.Code != http.StatusOK {
		t.Fatalf("page one status = %d body=%s", pageOne.Code, pageOne.Body.String())
	}
	got := decodeAdminResponse[jsonRecipeListTestResponse](t, pageOne)
	assertJSONRecipeListIDs(t, got.Recipes, first.ID, second.ID)
	if got.NextCursor == nil {
		t.Fatal("next_cursor = nil, want cursor")
	}

	pageTwo := adminJSON(t, router, http.MethodGet, "/api/admin/json-recipes?sort=asc&size=2&cursor="+*got.NextCursor, "", cookie)
	if pageTwo.Code != http.StatusOK {
		t.Fatalf("page two status = %d body=%s", pageTwo.Code, pageTwo.Body.String())
	}
	gotNext := decodeAdminResponse[jsonRecipeListTestResponse](t, pageTwo)
	assertJSONRecipeListIDs(t, gotNext.Recipes, third.ID)
	if gotNext.NextCursor != nil {
		t.Fatalf("next_cursor = %q, want nil", *gotNext.NextCursor)
	}
}

func TestInternalJSONRecipesListPaginationSortAndValidation(t *testing.T) {
	router, db := testRouter(t)

	base := time.Date(2026, 6, 22, 10, 0, 0, 0, time.UTC)
	category := createTestJSONRecipeCategory(t, db, ocr.JSONRecipeCategory{TitleEn: "Invoices", TitleRo: "Facturi"})
	first := createTestJSONRecipe(t, db, ocr.JSONRecipe{Title: "First", CreatedAt: base, UpdatedAt: base, CategoryID: &category.ID})
	second := createTestJSONRecipe(t, db, ocr.JSONRecipe{Title: "Second", CreatedAt: base.Add(time.Minute), UpdatedAt: base.Add(time.Minute)})
	third := createTestJSONRecipe(t, db, ocr.JSONRecipe{Title: "Third", CreatedAt: base.Add(2 * time.Minute), UpdatedAt: base.Add(2 * time.Minute)})

	pageOne := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/api/json-recipes?sort=asc&size=2", nil)
	router.ServeHTTP(pageOne, req)
	if pageOne.Code != http.StatusOK {
		t.Fatalf("page one status = %d body=%s", pageOne.Code, pageOne.Body.String())
	}
	got := decodeAdminResponse[jsonRecipeListTestResponse](t, pageOne)
	assertJSONRecipeListIDs(t, got.Recipes, first.ID, second.ID)
	if got.Recipes[0].CategoryID == nil || *got.Recipes[0].CategoryID != category.ID || got.Recipes[0].Category == nil {
		t.Fatalf("first category response = %#v, want category", got.Recipes[0])
	}
	if got.Recipes[1].CategoryID != nil || got.Recipes[1].Category != nil {
		t.Fatalf("second category response = %#v, want nil category", got.Recipes[1])
	}
	if got.NextCursor == nil {
		t.Fatal("next_cursor = nil, want cursor")
	}

	pageTwo := httptest.NewRecorder()
	req = newTestRequest(http.MethodGet, "/api/json-recipes?sort=asc&size=2&cursor="+*got.NextCursor, nil)
	router.ServeHTTP(pageTwo, req)
	if pageTwo.Code != http.StatusOK {
		t.Fatalf("page two status = %d body=%s", pageTwo.Code, pageTwo.Body.String())
	}
	gotNext := decodeAdminResponse[jsonRecipeListTestResponse](t, pageTwo)
	assertJSONRecipeListIDs(t, gotNext.Recipes, third.ID)
	if gotNext.NextCursor != nil {
		t.Fatalf("next_cursor = %q, want nil", *gotNext.NextCursor)
	}

	for _, path := range []string{
		"/api/json-recipes?sort=sideways",
		"/api/json-recipes?size=0",
		"/api/json-recipes?cursor=not-a-cursor",
		"/api/json-recipes?sort=desc&cursor=" + *got.NextCursor,
	} {
		w := httptest.NewRecorder()
		req := newTestRequest(http.MethodGet, path, nil)
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("%s status = %d body=%s, want 400", path, w.Code, w.Body.String())
		}
	}
}

func TestDeployJSONRecipeCreatesSchemaAndIncrementsCounterOnSuccess(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "json-recipe-deploy@example.com")
	category := createTestJSONRecipeCategory(t, db, ocr.JSONRecipeCategory{TitleEn: "Invoices", TitleRo: "Facturi"})
	recipe := createTestJSONRecipe(t, db, ocr.JSONRecipe{
		Title:       "Invoice recipe",
		Description: "Invoice fields",
		Counter:     4,
		CategoryID:  &category.ID,
	})

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/json-recipes/"+recipe.ID.String()+"/deploy", bytes.NewBufferString(`{"user_id":"`+user.ID+`"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("deploy status = %d body=%s", w.Code, w.Body.String())
	}
	var got jsonRecipeDeployTestResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode deploy response: %v body=%s", err, w.Body.String())
	}
	if got.Recipe.ID != recipe.ID || got.Recipe.Counter != 5 {
		t.Fatalf("unexpected recipe response: %#v", got.Recipe)
	}
	if got.Recipe.CategoryID == nil || *got.Recipe.CategoryID != category.ID || got.Recipe.Category == nil {
		t.Fatalf("unexpected deploy category response: %#v", got.Recipe)
	}
	if got.Schema.ID == uuid.Nil || got.Schema.Name != recipe.Title || got.Schema.Description != recipe.Description || !got.Schema.Strict {
		t.Fatalf("unexpected schema response: %#v", got.Schema)
	}
	if got.Schema.UserID == nil || string(*got.Schema.UserID) != user.ID {
		t.Fatalf("schema user_id = %v, want %s", got.Schema.UserID, user.ID)
	}
	assertJSONEqual(t, got.Schema.Schema, validRecipeSchema)

	var storedRecipe ocr.JSONRecipe
	if err := db.First(&storedRecipe, "id = ?", recipe.ID).Error; err != nil {
		t.Fatalf("load stored recipe: %v", err)
	}
	if storedRecipe.Counter != 5 {
		t.Fatalf("stored counter = %d, want 5", storedRecipe.Counter)
	}

	var storedSchema ocr.ExtractionSchema
	if err := db.First(&storedSchema, "id = ?", got.Schema.ID).Error; err != nil {
		t.Fatalf("load stored schema: %v", err)
	}
	if storedSchema.UserID == nil || *storedSchema.UserID != user.ID || !storedSchema.Strict {
		t.Fatalf("unexpected stored schema: %#v", storedSchema)
	}
}

func TestDeployJSONRecipeRejectsInvalidUserWithoutIncrement(t *testing.T) {
	router, db := testRouter(t)
	recipe := createTestJSONRecipe(t, db, ocr.JSONRecipe{Title: "Invoice recipe", Counter: 2})

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/json-recipes/"+recipe.ID.String()+"/deploy", bytes.NewBufferString(`{"user_id":"`+uuid.NewString()+`"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("deploy status = %d body=%s", w.Code, w.Body.String())
	}
	var storedRecipe ocr.JSONRecipe
	if err := db.First(&storedRecipe, "id = ?", recipe.ID).Error; err != nil {
		t.Fatalf("load stored recipe: %v", err)
	}
	if storedRecipe.Counter != 2 {
		t.Fatalf("stored counter = %d, want 2", storedRecipe.Counter)
	}
	var schemas int64
	if err := db.Model(&ocr.ExtractionSchema{}).Where("name = ?", recipe.Title).Count(&schemas).Error; err != nil {
		t.Fatalf("count schemas: %v", err)
	}
	if schemas != 0 {
		t.Fatalf("schemas created = %d, want 0", schemas)
	}
}

func TestDeleteJSONRecipeLeavesDeployedSchemas(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-json-recipes-delete@example.com", "admin")
	user := createAdminTestUser(t, db, "json-recipes-delete-owner@example.com", "user")
	cookie := createAdminTestSession(t, db, admin, "admin-json-recipes-delete-session")
	userID := user.ID
	recipe := createTestJSONRecipe(t, db, ocr.JSONRecipe{Title: "Invoice recipe"})
	schema := ocr.ExtractionSchema{
		UserID:      &userID,
		Name:        recipe.Title,
		Description: recipe.Description,
		SchemaJSON:  datatypes.JSON(recipe.JSON),
		Strict:      true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create deployed schema: %v", err)
	}

	w := adminJSON(t, router, http.MethodDelete, "/api/admin/json-recipes/"+recipe.ID.String(), "", cookie)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d body=%s", w.Code, w.Body.String())
	}

	var gotSchema ocr.ExtractionSchema
	if err := db.First(&gotSchema, "id = ?", schema.ID).Error; err != nil {
		t.Fatalf("load deployed schema after recipe delete: %v", err)
	}
}

func assertJSONRecipeListIDs(t *testing.T, recipes []jsonRecipeTestResponse, want ...uuid.UUID) {
	t.Helper()
	if len(recipes) != len(want) {
		t.Fatalf("recipe count = %d, want %d: %#v", len(recipes), len(want), recipes)
	}
	for i, wantID := range want {
		if recipes[i].ID != wantID {
			t.Fatalf("recipe[%d].id = %s, want %s", i, recipes[i].ID, wantID)
		}
	}
}
