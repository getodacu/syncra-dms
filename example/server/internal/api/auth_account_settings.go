package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	auth "ai.ro/syncra/internal/auth"
)

var (
	errLastSignInMethod      = errors.New("at least one sign-in method is required")
	errOAuthAccountLinked    = errors.New("account is already linked to another user")
	errProviderAlreadyLinked = errors.New("provider is already linked")
)

func (h *Handler) ListAuthSessions(c *gin.Context) {
	session, ok, err := h.loadAuthenticatedSession(c)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var sessions []auth.Session
	if err := h.DB.WithContext(c.Request.Context()).
		Where("user_id = ? AND expires_at > ?", session.UserID, h.now()).
		Order("created_at DESC, id DESC").
		Find(&sessions).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list sessions")
		return
	}

	items := make([]authSessionListItemResponse, 0, len(sessions))
	for _, item := range sessions {
		items = append(items, authSessionListItemJSON(item, session.ID))
	}
	c.JSON(http.StatusOK, authSessionListResponse{Sessions: items})
}

func (h *Handler) RevokeAuthSession(c *gin.Context) {
	session, ok, err := h.loadAuthenticatedSession(c)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	sessionID := strings.TrimSpace(c.Param("id"))
	if sessionID == "" {
		writeError(c, http.StatusBadRequest, "session id is required")
		return
	}
	if sessionID == session.ID {
		writeError(c, http.StatusBadRequest, "current session cannot be revoked here")
		return
	}

	var target auth.Session
	result := h.DB.WithContext(c.Request.Context()).
		Where("id = ? AND user_id = ? AND expires_at > ?", sessionID, session.UserID, h.now()).
		Limit(1).
		Find(&target)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to load session")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "session not found")
		return
	}

	deleteResult := h.DB.WithContext(c.Request.Context()).Where("id = ? AND user_id = ?", sessionID, session.UserID).Delete(&auth.Session{})
	if deleteResult.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to revoke session")
		return
	}
	c.JSON(http.StatusOK, deleteAuthSessionResponse{DeletedID: sessionID, DeletedCount: int(deleteResult.RowsAffected)})
}

func (h *Handler) ListAuthAccounts(c *gin.Context) {
	session, ok, err := h.loadAuthenticatedSession(c)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var accounts []auth.AuthAccount
	if err := h.DB.WithContext(c.Request.Context()).
		Where("user_id = ?", session.UserID).
		Order("provider_id ASC, created_at ASC").
		Find(&accounts).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list linked accounts")
		return
	}

	items := make([]authAccountListItemResponse, 0, len(accounts))
	for _, account := range accounts {
		items = append(items, authAccountListItemJSON(account))
	}
	c.JSON(http.StatusOK, authAccountListResponse{Accounts: items})
}

func (h *Handler) UnlinkAuthAccount(c *gin.Context) {
	session, ok, err := h.loadAuthenticatedSession(c)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	providerID := strings.TrimSpace(c.Param("provider_id"))
	if providerID != auth.GoogleProviderID && providerID != auth.GitHubProviderID {
		writeError(c, http.StatusBadRequest, "only google and github accounts can be unlinked")
		return
	}

	var deleted int64
	if err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var account auth.AuthAccount
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND provider_id = ?", session.UserID, providerID).
			Limit(1).
			Find(&account)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		var remaining []auth.AuthAccount
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND id <> ?", session.UserID, account.ID).
			Find(&remaining).Error; err != nil {
			return err
		}
		if signInMethodCount(remaining) == 0 {
			return errLastSignInMethod
		}

		result = tx.Where("id = ? AND user_id = ?", account.ID, session.UserID).Delete(&auth.AuthAccount{})
		if result.Error != nil {
			return result.Error
		}
		deleted = result.RowsAffected
		return nil
	}); err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			writeError(c, http.StatusNotFound, "linked account not found")
		case errors.Is(err, errLastSignInMethod):
			writeError(c, http.StatusBadRequest, errLastSignInMethod.Error())
		default:
			writeError(c, http.StatusInternalServerError, "failed to unlink linked account")
		}
		return
	}

	c.JSON(http.StatusOK, deleteAuthAccountResponse{DeletedProviderID: providerID, DeletedCount: int(deleted)})
}

func signInMethodCount(accounts []auth.AuthAccount) int {
	count := 0
	for _, account := range accounts {
		if account.ProviderID == auth.CredentialProviderID {
			if strings.TrimSpace(account.Password) != "" {
				count++
			}
			continue
		}
		count++
	}
	return count
}
