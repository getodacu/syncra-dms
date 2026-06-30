package rbac

import (
	"fmt"
	"time"

	"ai.ro/syncra/dms/internal/auth"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func BootstrapLegacyAdmins(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var marker BootstrapMarker
		result := tx.Where("name = ?", LegacyAdminsBootstrapMarkerName).Limit(1).Find(&marker)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			return nil
		}

		var adminRole Role
		if err := tx.First(&adminRole, "code = ?", SystemAdministratorRoleCode).Error; err != nil {
			return fmt.Errorf("load %s role: %w", SystemAdministratorRoleCode, err)
		}

		var users []auth.User
		if err := tx.
			Where("role = ? AND status IN ? AND deleted_at IS NULL", auth.UserRoleAdmin, []string{string(UserStatusActive), ""}).
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
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&link).Error; err != nil {
				return err
			}
		}
		marker = BootstrapMarker{
			Name:      LegacyAdminsBootstrapMarkerName,
			CreatedAt: now,
		}
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&marker).Error
	})
}
