package dashboard

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"ai.ro/syncra/internal/billing"
	"ai.ro/syncra/internal/ocr"
)

type RangeKey string

const (
	Range7D  RangeKey = "7d"
	Range30D RangeKey = "30d"
	Range90D RangeKey = "90d"
)

var ErrInvalidRange = errors.New("invalid dashboard range")

type RangeWindow struct {
	Key     RangeKey
	StartAt time.Time
	EndAt   time.Time
	Bucket  string
}

type Metrics struct {
	DocumentsProcessed int
	PagesProcessed     int
	JobsCompleted      int
	JobsFailed         int
	JobsProcessing     int
	CompletionRate     float64
	CreditsSpent       int
	DatasetCount       int
	SchemaCount        int
}

type DocumentBucket struct {
	Date               string
	DocumentsProcessed int
}

type RecentDocument struct {
	ID               string
	OriginalFilename string
	SchemaID         *string
	SchemaName       *string
	PageCount        int
	CreatedAt        time.Time
}

type SchemaThroughput struct {
	SchemaID           *string
	SchemaName         string
	DocumentsProcessed int
}

type RecentDataset struct {
	ID         string
	Name       string
	SchemaName string
	FieldCount int
	CreatedAt  time.Time
}

type DatasetSummary struct {
	TotalCount int
	Recent     []RecentDataset
}

type CreditSummary struct {
	AvailableCredits int
	CreditsSpent     int
	LowCredit        bool
}

type OnboardingSummary struct {
	HasSchema            bool
	HasCompletedDocument bool
	HasDataset           bool
	HasAPIKey            bool
	HasWebhook           bool
	ShowOnboarding       bool
}

type WarningSection string

const (
	WarningRecentDocuments  WarningSection = "recent_documents"
	WarningSchemaThroughput WarningSection = "schema_throughput"
	WarningDatasetSummary   WarningSection = "dataset_summary"
	WarningCreditSummary    WarningSection = "credit_summary"
)

type Warning struct {
	Section WarningSection
	Message string
}

type Summary struct {
	Range            RangeWindow
	Metrics          Metrics
	DocumentBuckets  []DocumentBucket
	RecentDocuments  []RecentDocument
	SchemaThroughput []SchemaThroughput
	DatasetSummary   DatasetSummary
	CreditSummary    CreditSummary
	Onboarding       OnboardingSummary
	Warnings         []Warning
}

func ParseRange(raw string) (RangeKey, error) {
	switch RangeKey(strings.TrimSpace(raw)) {
	case "":
		return Range30D, nil
	case Range7D:
		return Range7D, nil
	case Range30D:
		return Range30D, nil
	case Range90D:
		return Range90D, nil
	default:
		return "", ErrInvalidRange
	}
}

func WindowForRange(key RangeKey, now time.Time) (RangeWindow, error) {
	days, err := daysForRange(key)
	if err != nil {
		return RangeWindow{}, err
	}
	endAt := now.UTC()
	startAt := time.Date(endAt.Year(), endAt.Month(), endAt.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -(days - 1))
	return RangeWindow{Key: key, StartAt: startAt, EndAt: endAt, Bucket: "day"}, nil
}

func LoadSummary(ctx context.Context, db *gorm.DB, userID string, key RangeKey, now time.Time) (Summary, error) {
	if db == nil {
		return Summary{}, errors.New("dashboard: nil db")
	}
	window, err := WindowForRange(key, now)
	if err != nil {
		return Summary{}, err
	}

	summary := Summary{
		Range:            window,
		DocumentBuckets:  dailyBuckets(window, nil),
		RecentDocuments:  []RecentDocument{},
		SchemaThroughput: []SchemaThroughput{},
		DatasetSummary:   DatasetSummary{Recent: []RecentDataset{}},
		Warnings:         []Warning{},
	}

	var documentMetrics struct {
		DocumentsProcessed int
		PagesProcessed     int
	}
	if err := db.WithContext(ctx).
		Table("ocr_documents").
		Select("COUNT(*) AS documents_processed, COALESCE(SUM(page_count), 0) AS pages_processed").
		Where("user_id = ? AND created_at >= ? AND created_at <= ? AND deleted_at IS NULL", userID, window.StartAt, window.EndAt).
		Scan(&documentMetrics).Error; err != nil {
		return Summary{}, err
	}
	summary.Metrics.DocumentsProcessed = documentMetrics.DocumentsProcessed
	summary.Metrics.PagesProcessed = documentMetrics.PagesProcessed

	var bucketRows []struct {
		Date               string
		DocumentsProcessed int
	}
	if err := db.WithContext(ctx).
		Table("ocr_documents").
		Select("TO_CHAR(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD') AS date, COUNT(*) AS documents_processed").
		Where("user_id = ? AND created_at >= ? AND created_at <= ? AND deleted_at IS NULL", userID, window.StartAt, window.EndAt).
		Group("TO_CHAR(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD')").
		Scan(&bucketRows).Error; err != nil {
		return Summary{}, err
	}
	bucketCounts := make(map[string]int, len(bucketRows))
	for _, row := range bucketRows {
		bucketCounts[row.Date] = row.DocumentsProcessed
	}
	summary.DocumentBuckets = dailyBuckets(window, bucketCounts)

	var statusRows []struct {
		Status string
		Count  int
	}
	if err := db.WithContext(ctx).
		Table("ocr_jobs").
		Select("status, COUNT(*) AS count").
		Where("user_id = ? AND created_at >= ? AND created_at <= ? AND deleted_at IS NULL", userID, window.StartAt, window.EndAt).
		Group("status").
		Scan(&statusRows).Error; err != nil {
		return Summary{}, err
	}
	for _, row := range statusRows {
		switch ocr.OCRJobStatus(row.Status) {
		case ocr.OCRJobStatusCompleted:
			summary.Metrics.JobsCompleted = row.Count
		case ocr.OCRJobStatusFailed:
			summary.Metrics.JobsFailed = row.Count
		case ocr.OCRJobStatusProcessing, ocr.OCRJobStatusQueued:
			summary.Metrics.JobsProcessing += row.Count
		}
	}
	completedAndFailed := summary.Metrics.JobsCompleted + summary.Metrics.JobsFailed
	if completedAndFailed > 0 {
		summary.Metrics.CompletionRate = float64(summary.Metrics.JobsCompleted) / float64(completedAndFailed)
	}

	var schemaCount int64
	if err := db.WithContext(ctx).
		Table("extraction_schemas").
		Where("user_id = ?", userID).
		Count(&schemaCount).Error; err != nil {
		return Summary{}, err
	}
	summary.Metrics.SchemaCount = int(schemaCount)

	var datasetCount int64
	if err := db.WithContext(ctx).
		Table("datasets").
		Where("user_id = ?", userID).
		Count(&datasetCount).Error; err != nil {
		return Summary{}, err
	}
	summary.Metrics.DatasetCount = int(datasetCount)
	summary.DatasetSummary.TotalCount = int(datasetCount)

	loadCreditSummary(ctx, db, userID, window, now, &summary)
	loadRecentDocuments(ctx, db, userID, &summary)
	loadSchemaThroughput(ctx, db, userID, window, &summary)
	loadDatasetSummary(ctx, db, userID, &summary)
	if err := loadOnboarding(ctx, db, userID, &summary); err != nil {
		return Summary{}, err
	}

	return summary, nil
}

func daysForRange(key RangeKey) (int, error) {
	switch key {
	case Range7D:
		return 7, nil
	case Range30D:
		return 30, nil
	case Range90D:
		return 90, nil
	default:
		return 0, ErrInvalidRange
	}
}

func dailyBuckets(window RangeWindow, counts map[string]int) []DocumentBucket {
	days, err := daysForRange(window.Key)
	if err != nil {
		return []DocumentBucket{}
	}
	buckets := make([]DocumentBucket, 0, days)
	day := window.StartAt
	for i := 0; i < days; i++ {
		key := day.Format(time.DateOnly)
		buckets = append(buckets, DocumentBucket{Date: key, DocumentsProcessed: counts[key]})
		day = day.AddDate(0, 0, 1)
	}
	return buckets
}

func loadCreditSummary(ctx context.Context, db *gorm.DB, userID string, window RangeWindow, now time.Time, summary *Summary) {
	var creditsSpent int
	var availableCredits int
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Table("credit_ledger_entries").
			Select("COALESCE(SUM(ABS(credits_delta)), 0)").
			Where("user_id = ? AND entry_type = ? AND created_at >= ? AND created_at <= ?", userID, billing.CreditLedgerEntryDebit, window.StartAt, window.EndAt).
			Scan(&creditsSpent).Error; err != nil {
			return err
		}
		balance, err := billing.AvailableCredits(ctx, tx, userID, now.UTC())
		if err != nil {
			return err
		}
		availableCredits = balance.Available
		return nil
	})
	if err != nil {
		addWarning(summary, WarningCreditSummary, "Credit summary is temporarily unavailable.")
		return
	}
	summary.Metrics.CreditsSpent = creditsSpent
	summary.CreditSummary.CreditsSpent = creditsSpent
	summary.CreditSummary.AvailableCredits = availableCredits
	summary.CreditSummary.LowCredit = availableCredits <= 10
}

func loadRecentDocuments(ctx context.Context, db *gorm.DB, userID string, summary *Summary) {
	type recentDocumentRow struct {
		ID               string
		OriginalFilename string
		SchemaID         sql.NullString
		SchemaName       sql.NullString
		PageCount        int
		CreatedAt        time.Time
	}

	var rows []recentDocumentRow
	if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.
			Table("ocr_documents AS d").
			Select("d.id::text AS id, d.original_filename, CASE WHEN s.id IS NOT NULL THEN d.schema_id::text ELSE NULL END AS schema_id, s.name AS schema_name, d.page_count, d.created_at").
			Joins("LEFT JOIN extraction_schemas AS s ON s.id = d.schema_id AND s.user_id = ?", userID).
			Where("d.user_id = ? AND d.deleted_at IS NULL", userID).
			Order("d.created_at DESC, d.id DESC").
			Limit(5).
			Scan(&rows).Error
	}); err != nil {
		addWarning(summary, WarningRecentDocuments, "Recent documents are temporarily unavailable.")
		return
	}

	summary.RecentDocuments = make([]RecentDocument, 0, len(rows))
	for _, row := range rows {
		summary.RecentDocuments = append(summary.RecentDocuments, RecentDocument{
			ID:               row.ID,
			OriginalFilename: row.OriginalFilename,
			SchemaID:         nullableStringPtr(row.SchemaID),
			SchemaName:       nullableStringPtr(row.SchemaName),
			PageCount:        row.PageCount,
			CreatedAt:        row.CreatedAt,
		})
	}
}

func loadSchemaThroughput(ctx context.Context, db *gorm.DB, userID string, window RangeWindow, summary *Summary) {
	type schemaThroughputRow struct {
		SchemaID           sql.NullString
		SchemaName         string
		DocumentsProcessed int
	}
	schemaIDExpression := "CASE WHEN s.id IS NOT NULL THEN d.schema_id::text ELSE NULL END"
	schemaNameExpression := `CASE
		WHEN s.id IS NOT NULL THEN s.name
		WHEN d.inline_schema_json IS NOT NULL AND d.inline_schema_json <> 'null'::jsonb THEN 'Inline schema'
		ELSE 'No schema'
	END`

	var rows []schemaThroughputRow
	if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.
			Table("ocr_documents AS d").
			Select(`
			CASE WHEN s.id IS NOT NULL THEN d.schema_id::text ELSE NULL END AS schema_id,
			CASE
				WHEN s.id IS NOT NULL THEN s.name
				WHEN d.inline_schema_json IS NOT NULL AND d.inline_schema_json <> 'null'::jsonb THEN 'Inline schema'
				ELSE 'No schema'
			END AS schema_name,
			COUNT(*) AS documents_processed`).
			Joins("LEFT JOIN extraction_schemas AS s ON s.id = d.schema_id AND s.user_id = ?", userID).
			Where("d.user_id = ? AND d.created_at >= ? AND d.created_at <= ? AND d.deleted_at IS NULL", userID, window.StartAt, window.EndAt).
			Group(schemaIDExpression).
			Group(schemaNameExpression).
			Order("documents_processed DESC, MAX(d.created_at) DESC, schema_name ASC").
			Limit(5).
			Scan(&rows).Error
	}); err != nil {
		addWarning(summary, WarningSchemaThroughput, "Schema throughput is temporarily unavailable.")
		return
	}

	summary.SchemaThroughput = make([]SchemaThroughput, 0, len(rows))
	for _, row := range rows {
		summary.SchemaThroughput = append(summary.SchemaThroughput, SchemaThroughput{
			SchemaID:           nullableStringPtr(row.SchemaID),
			SchemaName:         row.SchemaName,
			DocumentsProcessed: row.DocumentsProcessed,
		})
	}
}

func loadDatasetSummary(ctx context.Context, db *gorm.DB, userID string, summary *Summary) {
	type recentDatasetRow struct {
		ID         string
		Name       string
		SchemaName string
		FieldCount int
		CreatedAt  time.Time
	}

	var rows []recentDatasetRow
	if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.
			Table("datasets AS d").
			Select("d.id::text AS id, d.name, COALESCE(s.name, '') AS schema_name, COALESCE(jsonb_array_length(d.selected_fields), 0) AS field_count, d.created_at").
			Joins("LEFT JOIN extraction_schemas AS s ON s.id = d.schema_id AND s.user_id = d.user_id").
			Where("d.user_id = ?", userID).
			Order("d.created_at DESC, d.id DESC").
			Limit(5).
			Scan(&rows).Error
	}); err != nil {
		addWarning(summary, WarningDatasetSummary, "Dataset summary is temporarily unavailable.")
		return
	}

	summary.DatasetSummary.Recent = make([]RecentDataset, 0, len(rows))
	for _, row := range rows {
		summary.DatasetSummary.Recent = append(summary.DatasetSummary.Recent, RecentDataset(row))
	}
}

func loadOnboarding(ctx context.Context, db *gorm.DB, userID string, summary *Summary) error {
	summary.Onboarding.HasSchema = summary.Metrics.SchemaCount > 0
	summary.Onboarding.HasDataset = summary.Metrics.DatasetCount > 0

	hasCompletedDocument, err := existsForUser(ctx, db, "ocr_documents", userID, "deleted_at IS NULL")
	if err != nil {
		return err
	}
	summary.Onboarding.HasCompletedDocument = hasCompletedDocument

	hasAPIKey, err := existsForUser(ctx, db, "api_keys", userID, "")
	if err != nil {
		return err
	}
	summary.Onboarding.HasAPIKey = hasAPIKey

	hasWebhook, err := existsForUser(ctx, db, "webhooks", userID, "")
	if err != nil {
		return err
	}
	summary.Onboarding.HasWebhook = hasWebhook
	summary.Onboarding.ShowOnboarding = !summary.Onboarding.HasCompletedDocument || !summary.Onboarding.HasDataset
	return nil
}

func existsForUser(ctx context.Context, db *gorm.DB, table, userID, extraCondition string) (bool, error) {
	query := db.WithContext(ctx).Table(table).Where("user_id = ?", userID)
	if extraCondition != "" {
		query = query.Where(extraCondition)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func nullableStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func addWarning(summary *Summary, section WarningSection, message string) {
	summary.Warnings = append(summary.Warnings, Warning{Section: section, Message: message})
}
