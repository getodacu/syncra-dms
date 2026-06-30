package rbac

import (
	"fmt"
	"time"

	"ai.ro/syncra/dms/internal/auth"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func BootstrapLegacyAdmins(db *gorm.DB) error {
	var adminRole Role
	if err := db.First(&adminRole, "code = ?", SystemAdministratorRoleCode).Error; err != nil {
		return fmt.Errorf("load %s role: %w", SystemAdministratorRoleCode, err)
	}

	var users []auth.User
	if err := db.
		Where("role = ? AND status = ? AND deleted_at IS NULL", auth.UserRoleAdmin, string(UserStatusActive)).
		Find(&users).Error; err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, user := range users {
		link := UserRole{
			UserID:    user.ID,
			RoleID:    adminRole.ID,
			ScopeType: ScopeGlobal,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&link).Error; err != nil {
			return err
		}
	}
	return nil
}
