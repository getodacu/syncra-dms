package api

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
	"ai.ro/syncra/internal/ocr"
	"ai.ro/syncra/internal/testsupport"
	"ai.ro/syncra/internal/webhooks"
)

var apiPostgresGroup *testsupport.PostgresGroup

func TestMain(m *testing.M) {
	group, cleanup, err := testsupport.OpenPostgresGroupForMain(apiPostgresModels()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open API postgres test group: %v\n", err)
		os.Exit(1)
	}
	apiPostgresGroup = group

	code := m.Run()
	if err := cleanup(); err != nil {
		fmt.Fprintf(os.Stderr, "cleanup API postgres test group: %v\n", err)
		if code == 0 {
			code = 1
		}
	}
	os.Exit(code)
}

func apiPostgresTx(t *testing.T) *gorm.DB {
	t.Helper()
	if apiPostgresGroup != nil {
		return apiPostgresGroup.Tx(t)
	}
	return testsupport.OpenPostgresTx(t, apiPostgresModels()...)
}

func apiPostgresModels() []any {
	return []any{
		&auth.User{},
		&auth.AuthAccount{},
		&auth.Session{},
		&auth.Verification{},
		&auth.APIKey{},
		&billing.BillingProfile{},
		&billing.BillingInvoiceCounter{},
		&billing.BillingInvoice{},
		&billing.BillingOrder{},
		&billing.CreditBucket{},
		&billing.CreditLedgerEntry{},
		&ocr.ExtractionSchema{},
		&ocr.JSONRecipeCategory{},
		&ocr.JSONRecipe{},
		&ocr.OCRDocument{},
		&ocr.OCRJob{},
		&ocr.Collection{},
		&ocr.CollectionSchema{},
		&ocr.CollectionDocument{},
		&ocr.Dataset{},
		&webhooks.Webhook{},
		&auth.AdminImpersonationEvent{},
	}
}
