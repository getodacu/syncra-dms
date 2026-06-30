package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
	"ai.ro/syncra/internal/ocr"
	"ai.ro/syncra/internal/webhooks"
)

func OpenPostgres(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func ApplicationModels() []any {
	return []any{
		&auth.User{},
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
		&auth.AuthAccount{},
		&auth.Session{},
		&auth.Verification{},
		&auth.AdminImpersonationEvent{},
	}
}
