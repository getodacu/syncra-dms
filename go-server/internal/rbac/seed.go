package rbac

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedDefaults(db *gorm.DB) error {
	now := time.Now().UTC()
	return db.Transaction(func(tx *gorm.DB) error {
		permissionIDs := map[string]string{}
		for _, definition := range PermissionRegistry {
			permission := Permission{
				Code:      definition.Code,
				Name:      definition.Name,
				Category:  definition.Category,
				IsSystem:  true,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if definition.Description != "" {
				description := definition.Description
				permission.Description = &description
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "code"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "description", "category", "is_system", "updated_at"}),
			}).Create(&permission).Error; err != nil {
				return err
			}
			var saved Permission
			if err := tx.First(&saved, "code = ?", definition.Code).Error; err != nil {
				return err
			}
			permissionIDs[definition.Code] = saved.ID
		}

		for _, definition := range DefaultRoles() {
			role := Role{
				Code:      definition.Code,
				Name:      definition.Name,
				IsSystem:  true,
				IsActive:  true,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if definition.Description != "" {
				description := definition.Description
				role.Description = &description
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "code"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "description", "is_system", "is_active", "updated_at"}),
			}).Create(&role).Error; err != nil {
				return err
			}
			var saved Role
			if err := tx.First(&saved, "code = ?", definition.Code).Error; err != nil {
				return err
			}
			if err := syncDefaultRolePermissions(tx, saved.ID, definition, permissionIDs, now); err != nil {
				return err
			}
		}
		return nil
	})
}

func syncDefaultRolePermissions(tx *gorm.DB, roleID string, definition RoleDefinition, permissionIDs map[string]string, now time.Time) error {
	expectedPermissionIDSet := make(map[string]struct{}, len(definition.PermissionCodes))
	expectedPermissionIDs := make([]string, 0, len(definition.PermissionCodes))
	for _, permissionCode := range definition.PermissionCodes {
		permissionID := permissionIDs[permissionCode]
		if permissionID == "" {
			return fmt.Errorf("role %q references unknown permission code %q", definition.Code, permissionCode)
		}
		if _, ok := expectedPermissionIDSet[permissionID]; ok {
			continue
		}
		expectedPermissionIDSet[permissionID] = struct{}{}
		expectedPermissionIDs = append(expectedPermissionIDs, permissionID)
	}

	for _, permissionID := range expectedPermissionIDs {
		link := RolePermission{RoleID: roleID, PermissionID: permissionID, CreatedAt: now}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&link).Error; err != nil {
			return err
		}
	}

	deleteExtra := tx.Where("role_id = ?", roleID)
	if len(expectedPermissionIDs) > 0 {
		deleteExtra = deleteExtra.Where("permission_id NOT IN ?", expectedPermissionIDs)
	}
	return deleteExtra.Delete(&RolePermission{}).Error
}
