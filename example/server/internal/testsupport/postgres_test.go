package testsupport

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type testWidget struct {
	ID   int64 `gorm:"primaryKey"`
	Name string
}

func TestOpenPostgresTxHidesExistingCommittedRows(t *testing.T) {
	dsn := PostgresTestDSN(t)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	if err := db.AutoMigrate(&testWidget{}); err != nil {
		t.Fatalf("auto migrate postgres: %v", err)
	}
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	})

	name := fmt.Sprintf("test-committed-schema-%d-%s", time.Now().UnixNano(), strings.ReplaceAll(t.Name(), "/", "-"))
	committed := testWidget{Name: name}
	if err := db.Create(&committed).Error; err != nil {
		t.Fatalf("create committed widget: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Unscoped().Delete(&testWidget{}, committed.ID).Error; err != nil {
			t.Fatalf("delete committed widget: %v", err)
		}
	})

	tx := OpenPostgresTx(t, &testWidget{})

	var count int64
	if err := tx.Model(&testWidget{}).Where("id = ?", committed.ID).Count(&count).Error; err != nil {
		t.Fatalf("count committed widget in tx: %v", err)
	}
	if count != 0 {
		t.Fatalf("committed widget visible inside test transaction count = %d, want 0", count)
	}

	if err := tx.Rollback().Error; err != nil {
		t.Fatalf("rollback test transaction: %v", err)
	}
	if err := db.Model(&testWidget{}).Where("id = ?", committed.ID).Count(&count).Error; err != nil {
		t.Fatalf("count committed widget after rollback: %v", err)
	}
	if count != 1 {
		t.Fatalf("committed widget after rollback count = %d, want 1", count)
	}
}

func TestOpenPostgresDBMigratesAndCleansTables(t *testing.T) {
	if !t.Run("commits row during test", func(t *testing.T) {
		db := OpenPostgresDB(t, &testWidget{})

		var firstBackendPID int
		if err := db.Raw("SELECT pg_backend_pid()").Scan(&firstBackendPID).Error; err != nil {
			t.Fatalf("select first backend pid: %v", err)
		}

		widget := testWidget{Name: "committed"}
		if err := db.Create(&widget).Error; err != nil {
			t.Fatalf("create widget: %v", err)
		}

		var count int64
		if err := db.Model(&testWidget{}).Count(&count).Error; err != nil {
			t.Fatalf("count widgets: %v", err)
		}
		if count != 1 {
			t.Fatalf("widget count = %d, want 1", count)
		}

		var secondBackendPID int
		if err := db.Raw("SELECT pg_backend_pid()").Scan(&secondBackendPID).Error; err != nil {
			t.Fatalf("select second backend pid: %v", err)
		}
		if secondBackendPID != firstBackendPID {
			t.Fatalf("backend pid changed from %d to %d", firstBackendPID, secondBackendPID)
		}
	}) {
		return
	}

	dsn := PostgresTestDSN(t)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	})

	var count int64
	if err := db.Model(&testWidget{}).Count(&count).Error; err != nil {
		t.Fatalf("count widgets after cleanup: %v", err)
	}
	if count != 0 {
		t.Fatalf("widget count after cleanup = %d, want 0", count)
	}
}
