package testsupport

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/dbmigrate"
)

// Serializes transaction-local truncation when tests share a Postgres database.
const postgresTestAdvisoryLockID int64 = 7442088320615058761

type PostgresGroup struct {
	db *gorm.DB
}

func PostgresTestDSN(t testing.TB) string {
	t.Helper()
	dsn, err := PostgresTestDSNValue()
	if err != nil {
		t.Fatal(err)
	}
	return dsn
}

func PostgresTestDSNValue() (string, error) {
	if err := loadServerEnvValue(); err != nil {
		return "", err
	}
	dsn := strings.TrimSpace(os.Getenv("DSN_DEV"))
	if dsn == "" {
		return "", errors.New("DSN_DEV is required in server/.env for Postgres-backed tests")
	}
	return dsn, nil
}

func OpenPostgresTx(t testing.TB, migrate ...any) *gorm.DB {
	t.Helper()
	dsn := PostgresTestDSN(t)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		closeDB(t, db)
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", postgresTestAdvisoryLockID).Error; err != nil {
		_ = tx.Rollback().Error
		closeDB(t, db)
		t.Fatalf("lock postgres test transaction: %v", err)
	}
	if len(migrate) > 0 {
		if err := migratePostgresTestModels(tx, migrate...); err != nil {
			_ = tx.Rollback().Error
			closeDB(t, db)
			t.Fatalf("migrate postgres test models: %v", err)
		}
		if err := truncateMigratedTables(tx, migrate...); err != nil {
			_ = tx.Rollback().Error
			closeDB(t, db)
			t.Fatalf("truncate migrated postgres tables: %v", err)
		}
	}

	t.Cleanup(func() {
		_ = tx.Rollback().Error
		closeDB(t, db)
	})
	return tx
}

func OpenPostgresDB(t testing.TB, migrate ...any) *gorm.DB {
	t.Helper()
	dsn := PostgresTestDSN(t)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get postgres sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	if err := db.Exec("SELECT pg_advisory_lock(?)", postgresTestAdvisoryLockID).Error; err != nil {
		closeDB(t, db)
		t.Fatalf("lock postgres test database: %v", err)
	}

	if len(migrate) > 0 {
		if err := migratePostgresTestModels(db, migrate...); err != nil {
			unlockAndCloseDB(t, db)
			t.Fatalf("migrate postgres test models: %v", err)
		}
		if err := truncateMigratedTables(db, migrate...); err != nil {
			unlockAndCloseDB(t, db)
			t.Fatalf("truncate migrated postgres tables: %v", err)
		}
	}

	t.Cleanup(func() {
		defer unlockAndCloseDB(t, db)
		if len(migrate) > 0 {
			if err := truncateMigratedTables(db, migrate...); err != nil {
				t.Fatalf("cleanup truncate migrated postgres tables: %v", err)
			}
		}
	})
	return db
}

func OpenPostgresGroup(t testing.TB, migrate ...any) *PostgresGroup {
	t.Helper()
	group, cleanup, err := OpenPostgresGroupForMain(migrate...)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := cleanup(); err != nil {
			t.Fatal(err)
		}
	})
	return group
}

func OpenPostgresGroupForMain(migrate ...any) (*PostgresGroup, func() error, error) {
	dsn, err := PostgresTestDSNValue()
	if err != nil {
		return nil, nil, err
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("open postgres: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		closeDBValue(db)
		return nil, nil, fmt.Errorf("get postgres sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	if err := db.Exec("SELECT pg_advisory_lock(?)", postgresTestAdvisoryLockID).Error; err != nil {
		closeDBValue(db)
		return nil, nil, fmt.Errorf("lock postgres test group: %w", err)
	}

	if len(migrate) > 0 {
		if err := migratePostgresTestModels(db, migrate...); err != nil {
			_ = unlockAndCloseDBValue(db)
			return nil, nil, fmt.Errorf("migrate postgres test models: %w", err)
		}
		if err := truncateMigratedTables(db, migrate...); err != nil {
			_ = unlockAndCloseDBValue(db)
			return nil, nil, fmt.Errorf("truncate migrated postgres tables: %w", err)
		}
	}

	cleanup := func() error {
		if len(migrate) > 0 {
			if err := truncateMigratedTables(db, migrate...); err != nil {
				_ = unlockAndCloseDBValue(db)
				return fmt.Errorf("cleanup truncate migrated postgres tables: %w", err)
			}
		}
		return unlockAndCloseDBValue(db)
	}
	return &PostgresGroup{db: db}, cleanup, nil
}

func (g *PostgresGroup) Tx(t testing.TB) *gorm.DB {
	t.Helper()
	tx := g.db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin postgres test transaction: %v", tx.Error)
	}
	t.Cleanup(func() {
		_ = tx.Rollback().Error
	})
	return tx
}

func (g *PostgresGroup) DB() *gorm.DB {
	return g.db
}

func migratePostgresTestModels(db *gorm.DB, migrate ...any) error {
	if err := dbmigrate.ResetLegacyIntegerIDTables(db); err != nil {
		return fmt.Errorf("reset legacy integer-id tables: %w", err)
	}
	if err := dbmigrate.MigrateOCRDocumentHash(db); err != nil {
		return fmt.Errorf("migrate OCR document hash: %w", err)
	}
	if err := dbmigrate.MigrateCreditOnlyBilling(db); err != nil {
		return fmt.Errorf("migrate credit-only billing: %w", err)
	}
	if err := db.AutoMigrate(migrate...); err != nil {
		return fmt.Errorf("auto migrate postgres: %w", err)
	}
	if err := dbmigrate.MigrateOCRDocumentJobForeignKey(db); err != nil {
		return fmt.Errorf("migrate OCR document job foreign key: %w", err)
	}
	if err := dbmigrate.MigrateOCRDocumentLifecycle(db); err != nil {
		return fmt.Errorf("migrate OCR document lifecycle: %w", err)
	}
	if err := dbmigrate.MigrateOCRDocumentPageCount(db); err != nil {
		return fmt.Errorf("migrate OCR document page count: %w", err)
	}
	if err := dbmigrate.MigrateOCRDocumentListIndexes(db); err != nil {
		return fmt.Errorf("migrate OCR document list indexes: %w", err)
	}
	if err := dbmigrate.MigrateOCRJobStatus(db); err != nil {
		return fmt.Errorf("migrate OCR job status: %w", err)
	}
	if err := dbmigrate.MigrateOwnerForeignKeyCascades(db); err != nil {
		return fmt.Errorf("migrate owner foreign keys: %w", err)
	}
	return nil
}

func unlockAndCloseDB(t testing.TB, db *gorm.DB) {
	t.Helper()
	if err := unlockAndCloseDBValue(db); err != nil {
		t.Fatal(err)
	}
}

func unlockAndCloseDBValue(db *gorm.DB) error {
	var unlocked bool
	if err := db.Raw("SELECT pg_advisory_unlock(?)", postgresTestAdvisoryLockID).Scan(&unlocked).Error; err != nil {
		closeDBValue(db)
		return fmt.Errorf("unlock postgres test database: %w", err)
	}
	if !unlocked {
		closeDBValue(db)
		return errors.New("unlock postgres test database: lock was not held by this connection")
	}
	closeDBValue(db)
	return nil
}

func truncateMigratedTables(db *gorm.DB, models ...any) error {
	tableNames, err := migratedTableNames(db, models...)
	if err != nil {
		return err
	}
	if len(tableNames) == 0 {
		return nil
	}

	quoted := make([]string, 0, len(tableNames))
	for _, tableName := range tableNames {
		quoted = append(quoted, quoteTableName(tableName))
	}

	return db.Exec("TRUNCATE TABLE " + strings.Join(quoted, ", ") + " RESTART IDENTITY CASCADE").Error
}

func migratedTableNames(db *gorm.DB, models ...any) ([]string, error) {
	tableNames := make([]string, 0, len(models))
	seen := make(map[string]struct{}, len(models))
	for _, model := range models {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			return nil, fmt.Errorf("parse model table: %w", err)
		}
		if stmt.Schema == nil || stmt.Schema.Table == "" {
			return nil, fmt.Errorf("parse model table: missing table name for %T", model)
		}
		tableName := stmt.Schema.Table
		if _, ok := seen[tableName]; ok {
			continue
		}
		seen[tableName] = struct{}{}
		tableNames = append(tableNames, tableName)
	}
	return tableNames, nil
}

func quoteTableName(tableName string) string {
	parts := strings.Split(tableName, ".")
	for i, part := range parts {
		parts[i] = `"` + strings.ReplaceAll(part, `"`, `""`) + `"`
	}
	return strings.Join(parts, ".")
}

func closeDB(t testing.TB, db *gorm.DB) {
	t.Helper()
	closeDBValue(db)
}

func closeDBValue(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}

func loadServerEnv(t testing.TB) {
	t.Helper()
	if err := loadServerEnvValue(); err != nil {
		t.Fatal(err)
	}
}

func loadServerEnvValue() error {
	root, err := serverRootValue()
	if err != nil {
		return err
	}
	envPath := filepath.Join(root, ".env")
	if _, err := os.Stat(envPath); err != nil {
		return fmt.Errorf("server/.env is required for tests; copy it into the worktree first: %w", err)
	}
	if err := godotenv.Overload(envPath); err != nil {
		return fmt.Errorf("load server/.env: %w", err)
	}
	return nil
}

func serverRoot(t testing.TB) string {
	t.Helper()
	root, err := serverRootValue()
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func serverRootValue() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("cannot resolve testsupport path")
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		next := filepath.Dir(dir)
		if next == dir {
			return "", errors.New("cannot find server go.mod")
		}
		dir = next
	}
}
