package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"ai.ro/syncra/dms/internal/orgunits"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const maxOrganizationUnitRequestBytes int64 = 1 << 20

var errOrganizationUnitResponseWritten = errors.New("organization unit response already written")

type organizationUnitHandler struct {
	db   *gorm.DB
	auth *authHandler
}

type organizationUnitRequest struct {
	ParentID    *string `json:"parentId"`
	Name        string  `json:"name"`
	Code        *string `json:"code"`
	Description *string `json:"description"`
}

type moveOrganizationUnitRequest struct {
	ParentID *string
}

func (r *moveOrganizationUnitRequest) UnmarshalJSON(data []byte) error {
	type rawRequest struct {
		ParentID *string `json:"parentId"`
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if _, ok := raw["parentId"]; !ok {
		return errors.New("parentId is required")
	}
	var out rawRequest
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	r.ParentID = out.ParentID
	return nil
}

type organizationUnitResponse struct {
	ID          string                     `json:"id"`
	ParentID    *string                    `json:"parentId,omitempty"`
	Name        string                     `json:"name"`
	Code        *string                    `json:"code,omitempty"`
	Description *string                    `json:"description,omitempty"`
	ArchivedAt  *string                    `json:"archivedAt,omitempty"`
	CreatedAt   string                     `json:"createdAt"`
	UpdatedAt   string                     `json:"updatedAt"`
	Children    []organizationUnitResponse `json:"children"`
}

type organizationUnitListResponse struct {
	Units []organizationUnitResponse `json:"units"`
}

type normalizedOrganizationUnitInput struct {
	Name                string
	Code                *string
	Description         *string
	CodeProvided        bool
	DescriptionProvided bool
}

func newOrganizationUnitHandler(options RouterOptions, auth *authHandler) *organizationUnitHandler {
	return &organizationUnitHandler{db: options.DB, auth: auth}
}

func (h *organizationUnitHandler) listTree(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "organization_unit.view", nil); !ok {
		return
	}
	var units []orgunits.Unit
	if err := h.db.WithContext(c.Request.Context()).
		Where("archived_at IS NULL").
		Order("name asc, id asc").
		Find(&units).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list organization units")
		return
	}
	c.JSON(http.StatusOK, organizationUnitListResponse{Units: organizationUnitTreeResponse(orgunits.BuildTree(units))})
}

func (h *organizationUnitHandler) listArchived(c *gin.Context) {
	if _, ok := requireAnyPermission(c, h.auth, []string{"organization_unit.view_audit", "organization_unit.manage_hierarchy"}, nil); !ok {
		return
	}
	var units []orgunits.Unit
	if err := h.db.WithContext(c.Request.Context()).
		Where("archived_at IS NOT NULL").
		Order("archived_at desc, name asc, id asc").
		Find(&units).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list archived organization units")
		return
	}
	out := make([]organizationUnitResponse, 0, len(units))
	for _, unit := range units {
		out = append(out, organizationUnitResponseFromUnit(unit))
	}
	c.JSON(http.StatusOK, organizationUnitListResponse{Units: out})
}

func (h *organizationUnitHandler) create(c *gin.Context) {
	if _, ok := requireAnyPermission(c, h.auth, []string{"organization_unit.create", "organization_unit.manage_hierarchy"}, nil); !ok {
		return
	}
	req, ok := bindOrganizationUnitRequest(c)
	if !ok {
		return
	}
	input, ok := normalizeOrganizationUnitInput(c, req)
	if !ok {
		return
	}
	if ok := h.ensureActiveCodeAvailable(c, input.Code, ""); !ok {
		return
	}
	now := time.Now().UTC()
	var unit orgunits.Unit
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		parentID, ok := h.validateOptionalActiveParentWithDB(c, tx, req.ParentID)
		if !ok {
			return errOrganizationUnitResponseWritten
		}
		unit = orgunits.Unit{
			ParentID:    parentID,
			Name:        input.Name,
			Code:        input.Code,
			Description: input.Description,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		return tx.Create(&unit).Error
	})
	if err != nil {
		if errors.Is(err, errOrganizationUnitResponseWritten) {
			return
		}
		writeOrganizationUnitMutationError(c, err, "failed to create organization unit")
		return
	}
	c.JSON(http.StatusCreated, organizationUnitResponseFromUnit(unit))
}

func (h *organizationUnitHandler) update(c *gin.Context) {
	id, ok := parseOrganizationUnitID(c, c.Param("id"))
	if !ok {
		return
	}
	if _, ok := requireAnyPermission(c, h.auth, []string{"organization_unit.update", "organization_unit.manage_hierarchy"}, &id); !ok {
		return
	}
	req, ok := bindOrganizationUnitRequest(c)
	if !ok {
		return
	}
	input, ok := normalizeOrganizationUnitInput(c, req)
	if !ok {
		return
	}
	if input.CodeProvided {
		if ok := h.ensureActiveCodeAvailable(c, input.Code, id); !ok {
			return
		}
	}
	updates := map[string]any{
		"name":       input.Name,
		"updated_at": time.Now().UTC(),
	}
	if input.CodeProvided {
		updates["code"] = nullableStringValue(input.Code)
	}
	if input.DescriptionProvided {
		updates["description"] = nullableStringValue(input.Description)
	}
	var unit orgunits.Unit
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&orgunits.Unit{}).Where("id = ? AND archived_at IS NULL", id).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.First(&unit, "id = ?", id).Error
	})
	if err != nil {
		writeOrganizationUnitMutationError(c, err, "failed to update organization unit")
		return
	}
	c.JSON(http.StatusOK, organizationUnitResponseFromUnit(unit))
}

func (h *organizationUnitHandler) move(c *gin.Context) {
	id, ok := parseOrganizationUnitID(c, c.Param("id"))
	if !ok {
		return
	}
	if _, ok := requirePermission(c, h.auth, "organization_unit.manage_hierarchy", &id); !ok {
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxOrganizationUnitRequestBytes)
	var req moveOrganizationUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	var unit orgunits.Unit
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		parentID, ok := h.validateOptionalActiveParentWithDB(c, tx, req.ParentID)
		if !ok {
			return errOrganizationUnitResponseWritten
		}
		if parentID != nil && *parentID == id {
			writeError(c, http.StatusConflict, "organization unit cannot be moved under itself")
			return errOrganizationUnitResponseWritten
		}
		if parentID != nil {
			wouldCreateCycle, ok := h.parentWouldCreateCycleWithDB(c, tx, id, *parentID)
			if !ok {
				return errOrganizationUnitResponseWritten
			}
			if wouldCreateCycle {
				writeError(c, http.StatusConflict, "organization unit cannot be moved under its descendant")
				return errOrganizationUnitResponseWritten
			}
		}
		result := tx.Model(&orgunits.Unit{}).Where("id = ? AND archived_at IS NULL", id).Updates(map[string]any{
			"parent_id":  nullableStringValue(parentID),
			"updated_at": time.Now().UTC(),
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.First(&unit, "id = ?", id).Error
	})
	if err != nil {
		if errors.Is(err, errOrganizationUnitResponseWritten) {
			return
		}
		writeOrganizationUnitMutationError(c, err, "failed to move organization unit")
		return
	}
	c.JSON(http.StatusOK, organizationUnitResponseFromUnit(unit))
}

func (h *organizationUnitHandler) archive(c *gin.Context) {
	id, ok := parseOrganizationUnitID(c, c.Param("id"))
	if !ok {
		return
	}
	if _, ok := requireAnyPermission(c, h.auth, []string{"organization_unit.delete", "organization_unit.manage_hierarchy"}, &id); !ok {
		return
	}
	now := time.Now().UTC()
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var units []orgunits.Unit
		if err := tx.Where("archived_at IS NULL").Find(&units).Error; err != nil {
			return err
		}
		found := false
		for _, unit := range units {
			if unit.ID == id {
				found = true
				break
			}
		}
		if !found {
			return gorm.ErrRecordNotFound
		}
		ids := []string{id}
		for descendantID := range orgunits.DescendantIDs(id, units) {
			ids = append(ids, descendantID)
		}
		return tx.Model(&orgunits.Unit{}).Where("id IN ?", ids).Updates(map[string]any{
			"archived_at": now,
			"updated_at":  now,
		}).Error
	})
	if err != nil {
		writeOrganizationUnitMutationError(c, err, "failed to archive organization unit")
		return
	}
	c.JSON(http.StatusOK, okResponse{OK: true})
}

func bindOrganizationUnitRequest(c *gin.Context) (organizationUnitRequest, bool) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxOrganizationUnitRequestBytes)
	var req organizationUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return organizationUnitRequest{}, false
	}
	return req, true
}

func normalizeOrganizationUnitInput(c *gin.Context, req organizationUnitRequest) (normalizedOrganizationUnitInput, bool) {
	name, err := orgunits.NormalizeName(req.Name)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return normalizedOrganizationUnitInput{}, false
	}
	codeRaw := ""
	if req.Code != nil {
		codeRaw = *req.Code
	}
	code, err := orgunits.NormalizeCode(codeRaw)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return normalizedOrganizationUnitInput{}, false
	}
	descriptionRaw := ""
	if req.Description != nil {
		descriptionRaw = *req.Description
	}
	return normalizedOrganizationUnitInput{
		Name:                name,
		Code:                code,
		Description:         orgunits.NormalizeDescription(descriptionRaw),
		CodeProvided:        req.Code != nil,
		DescriptionProvided: req.Description != nil,
	}, true
}

func parseOrganizationUnitID(c *gin.Context, raw string) (string, bool) {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil || id == uuid.Nil {
		writeError(c, http.StatusBadRequest, "invalid organization unit id")
		return "", false
	}
	return id.String(), true
}

func (h *organizationUnitHandler) validateOptionalActiveParentWithDB(c *gin.Context, db *gorm.DB, raw *string) (*string, bool) {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return nil, true
	}
	parentID, ok := parseOrganizationUnitID(c, *raw)
	if !ok {
		return nil, false
	}
	var count int64
	if err := db.WithContext(c.Request.Context()).Model(&orgunits.Unit{}).Where("id = ? AND archived_at IS NULL", parentID).Count(&count).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate parent organization unit")
		return nil, false
	}
	if count == 0 {
		writeError(c, http.StatusNotFound, "parent organization unit not found")
		return nil, false
	}
	return &parentID, true
}

func (h *organizationUnitHandler) ensureActiveCodeAvailable(c *gin.Context, code *string, currentID string) bool {
	if code == nil {
		return true
	}
	query := h.db.WithContext(c.Request.Context()).Model(&orgunits.Unit{}).Where("code = ? AND archived_at IS NULL", *code)
	if currentID != "" {
		query = query.Where("id <> ?", currentID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate organization unit code")
		return false
	}
	if count > 0 {
		writeError(c, http.StatusConflict, "organization unit code already exists")
		return false
	}
	return true
}

func nullableStringValue(value *string) any {
	if value == nil {
		return gorm.Expr("NULL")
	}
	return *value
}

func (h *organizationUnitHandler) parentWouldCreateCycleWithDB(c *gin.Context, db *gorm.DB, unitID string, parentID string) (bool, bool) {
	var units []orgunits.Unit
	if err := db.WithContext(c.Request.Context()).Where("archived_at IS NULL").Find(&units).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate organization unit move")
		return false, false
	}
	descendants := orgunits.DescendantIDs(unitID, units)
	return descendants[parentID], true
}

func organizationUnitTreeResponse(nodes []orgunits.TreeNode) []organizationUnitResponse {
	out := make([]organizationUnitResponse, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, organizationUnitResponse{
			ID:          node.ID,
			ParentID:    node.ParentID,
			Name:        node.Name,
			Code:        node.Code,
			Description: node.Description,
			CreatedAt:   node.CreatedAt,
			UpdatedAt:   node.UpdatedAt,
			Children:    organizationUnitTreeResponse(node.Children),
		})
	}
	return out
}

func organizationUnitResponseFromUnit(unit orgunits.Unit) organizationUnitResponse {
	var archivedAt *string
	if unit.ArchivedAt != nil {
		value := unit.ArchivedAt.UTC().Format(time.RFC3339Nano)
		archivedAt = &value
	}
	return organizationUnitResponse{
		ID:          unit.ID,
		ParentID:    unit.ParentID,
		Name:        unit.Name,
		Code:        unit.Code,
		Description: unit.Description,
		ArchivedAt:  archivedAt,
		CreatedAt:   unit.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:   unit.UpdatedAt.UTC().Format(time.RFC3339Nano),
		Children:    []organizationUnitResponse{},
	}
}

func writeOrganizationUnitMutationError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(c, http.StatusNotFound, "organization unit not found")
		return
	}
	if isUniqueConstraintError(err) {
		writeError(c, http.StatusConflict, "organization unit code already exists")
		return
	}
	writeError(c, http.StatusInternalServerError, fallback)
}

func isUniqueConstraintError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	message := err.Error()
	return strings.Contains(message, "idx_organization_units_active_code_unique") ||
		strings.Contains(message, "UNIQUE constraint failed")
}
