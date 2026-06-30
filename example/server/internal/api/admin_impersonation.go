package api

import (
	"errors"
	"net/http"

	"ai.ro/syncra/internal/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	errAdminImpersonationAlreadyActive = errors.New("already impersonating another user")
	errAdminImpersonationTargetNotUser = errors.New("only normal users can be impersonated")
	errAdminImpersonationSelf          = errors.New("cannot impersonate yourself")
	errAdminImpersonationTargetMissing = errors.New("user not found")
)

func (h *Handler) StartAdminUserImpersonation(c *gin.Context) {
	targetUserID, ok := parseAdminUserIDParam(c)
	if !ok {
		return
	}

	session, ok := adminSessionFromContext(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "authentication required")
		return
	}
	if session.UserID == targetUserID {
		writeError(c, http.StatusBadRequest, errAdminImpersonationSelf.Error())
		return
	}

	var updated auth.Session
	if err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("User").
			Preload("ImpersonatedUser").
			First(&updated, "id = ?", session.ID).Error; err != nil {
			return err
		}
		if updated.ImpersonatedUserID != nil {
			if *updated.ImpersonatedUserID == targetUserID {
				return nil
			}
			return errAdminImpersonationAlreadyActive
		}

		var target auth.User
		if err := tx.First(&target, "id = ?", targetUserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errAdminImpersonationTargetMissing
			}
			return err
		}
		if target.Role != auth.UserRoleUser {
			return errAdminImpersonationTargetNotUser
		}

		now := h.now().UTC()
		if err := tx.Model(&auth.Session{}).
			Where("id = ?", updated.ID).
			Updates(map[string]any{
				"impersonated_user_id":     target.ID,
				"impersonation_started_at": now,
			}).Error; err != nil {
			return err
		}
		if err := tx.Create(&auth.AdminImpersonationEvent{
			EventType:       auth.AdminImpersonationEventStart,
			SessionID:       updated.ID,
			AdminUserID:     updated.UserID,
			AdminUserEmail:  updated.User.Email,
			TargetUserID:    target.ID,
			TargetUserEmail: target.Email,
			IPAddress:       c.ClientIP(),
			UserAgent:       c.Request.UserAgent(),
			CreatedAt:       now,
		}).Error; err != nil {
			return err
		}

		return tx.Preload("User").Preload("ImpersonatedUser").First(&updated, "id = ?", updated.ID).Error
	}); err != nil {
		h.writeAdminImpersonationError(c, err)
		return
	}

	c.JSON(http.StatusOK, authSessionPayloadJSON(updated))
}

func (h *Handler) StopAdminImpersonation(c *gin.Context) {
	session, ok := adminSessionFromContext(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var updated auth.Session
	if err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("User").
			Preload("ImpersonatedUser").
			First(&updated, "id = ?", session.ID).Error; err != nil {
			return err
		}
		if updated.ImpersonatedUserID == nil {
			return nil
		}

		targetUserID := *updated.ImpersonatedUserID
		targetEmail := ""
		if updated.ImpersonatedUser != nil {
			targetEmail = updated.ImpersonatedUser.Email
		}
		now := h.now().UTC()
		if err := clearSessionImpersonation(c.Request.Context(), tx, updated.ID); err != nil {
			return err
		}
		if err := tx.Create(&auth.AdminImpersonationEvent{
			EventType:       auth.AdminImpersonationEventStop,
			SessionID:       updated.ID,
			AdminUserID:     updated.UserID,
			AdminUserEmail:  updated.User.Email,
			TargetUserID:    targetUserID,
			TargetUserEmail: targetEmail,
			IPAddress:       c.ClientIP(),
			UserAgent:       c.Request.UserAgent(),
			CreatedAt:       now,
		}).Error; err != nil {
			return err
		}

		return tx.Preload("User").Preload("ImpersonatedUser").First(&updated, "id = ?", updated.ID).Error
	}); err != nil {
		h.writeAdminImpersonationError(c, err)
		return
	}

	c.JSON(http.StatusOK, authSessionPayloadJSON(updated))
}

func (h *Handler) writeAdminImpersonationError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errAdminImpersonationAlreadyActive):
		writeError(c, http.StatusConflict, errAdminImpersonationAlreadyActive.Error())
	case errors.Is(err, errAdminImpersonationSelf), errors.Is(err, errAdminImpersonationTargetNotUser):
		writeError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, errAdminImpersonationTargetMissing):
		writeError(c, http.StatusNotFound, errAdminImpersonationTargetMissing.Error())
	case errors.Is(err, gorm.ErrRecordNotFound):
		writeError(c, http.StatusUnauthorized, "authentication required")
	default:
		writeError(c, http.StatusInternalServerError, "failed to update impersonation")
	}
}
