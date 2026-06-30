package dashboard

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
	"ai.ro/syncra/internal/ocr"
	"ai.ro/syncra/internal/testsupport"
	"ai.ro/syncra/internal/webhooks"
)

func dashboardTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	return testsupport.OpenPostgresTx(t,
		&auth.User{},
		&auth.APIKey{},
		&billing.CreditBucket{},
		&billing.CreditLedgerEntry{},
		&ocr.ExtractionSchema{},
		&ocr.OCRJob{},
		&ocr.OCRDocument{},
		&ocr.Dataset{},
		&webhooks.Webhook{},
	)
}

func createDashboardUser(t *testing.T, db *gorm.DB, email string) auth.User {
	t.Helper()
	user := auth.User{Name: "Dashboard User", Email: email}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create dashboard user: %v", err)
	}
	return user
}

func ptrString(value string) *string {
	return &value
}

func createDashboardSchema(t *testing.T, db *gorm.DB, userID, name string) ocr.ExtractionSchema {
	t.Helper()
	schema := ocr.ExtractionSchema{
		UserID:     ptrString(userID),
		Name:       name,
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","properties":{"total":{"type":"number"}}}`)),
		Strict:     true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create dashboard schema: %v", err)
	}
	return schema
}

func createDashboardJob(t *testing.T, db *gorm.DB, userID string, createdAt time.Time, status ocr.OCRJobStatus, schemaID *uuid.UUID, pageCount int) ocr.OCRJob {
	t.Helper()
	createdAt = createdAt.UTC()
	job := ocr.OCRJob{
		UserID:           ptrString(userID),
		CreatedAt:        createdAt,
		UpdatedAt:        createdAt,
		OriginalFilename: "job-" + uuid.NewString() + ".pdf",
		MimeType:         "application/pdf",
		FileSize:         int64(pageCount * 100),
		PageCount:        pageCount,
		DocumentHash:     uuid.NewString(),
		FilePath:         "/tmp/" + uuid.NewString(),
		SchemaID:         schemaID,
		Status:           status,
	}
	if err := db.Create(&job).Error; err != nil {
		t.Fatalf("create dashboard job: %v", err)
	}
	return job
}

func createDashboardDocument(t *testing.T, db *gorm.DB, userID string, createdAt time.Time, filename string, schemaID *uuid.UUID, inline bool, pageCount int) ocr.OCRDocument {
	t.Helper()
	createdAt = createdAt.UTC()
	doc := ocr.OCRDocument{
		UserID:           ptrString(userID),
		CreatedAt:        createdAt,
		UpdatedAt:        createdAt,
		OriginalFilename: filename,
		MimeType:         "application/pdf",
		FileSize:         int64(pageCount * 100),
		DocumentHash:     uuid.NewString(),
		SchemaID:         schemaID,
		Markdown:         "# OCR",
		AnnotationJSON:   datatypes.JSON([]byte(`{"total":42}`)),
		RawResponseJSON:  dashboardRawPages(pageCount),
	}
	if inline {
		doc.InlineSchemaJSON = datatypes.JSON([]byte(`{"type":"object","properties":{"total":{"type":"number"}}}`))
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create dashboard document: %v", err)
	}
	return doc
}

func dashboardRawPages(pageCount int) datatypes.JSON {
	raw := `{"pages":[`
	for i := 0; i < pageCount; i++ {
		if i > 0 {
			raw += ","
		}
		raw += `{"index":` + strconv.Itoa(i) + `}`
	}
	raw += `]}`
	return datatypes.JSON([]byte(raw))
}

func createDashboardCreditBucket(t *testing.T, db *gorm.DB, userID string, credits int, createdAt time.Time) billing.CreditBucket {
	t.Helper()
	createdAt = createdAt.UTC()
	bucket := billing.CreditBucket{
		UserID:           userID,
		SourceType:       billing.CreditSourceAdjustment,
		CreditsGranted:   credits,
		CreditsRemaining: credits,
		ValidFrom:        createdAt,
		CreatedAt:        createdAt,
		UpdatedAt:        createdAt,
	}
	if err := db.Create(&bucket).Error; err != nil {
		t.Fatalf("create dashboard credit bucket: %v", err)
	}
	return bucket
}

func createDashboardLedgerEntry(t *testing.T, db *gorm.DB, userID string, bucketID uuid.UUID, entryType billing.CreditLedgerEntryType, creditsDelta int, createdAt time.Time) billing.CreditLedgerEntry {
	t.Helper()
	entry := billing.CreditLedgerEntry{
		UserID:         userID,
		BucketID:       &bucketID,
		EntryType:      entryType,
		CreditsDelta:   creditsDelta,
		IdempotencyKey: "dashboard:" + uuid.NewString(),
		Metadata:       datatypes.JSON([]byte(`{}`)),
		CreatedAt:      createdAt.UTC(),
	}
	if err := db.Create(&entry).Error; err != nil {
		t.Fatalf("create dashboard ledger entry: %v", err)
	}
	return entry
}

func createDashboardDataset(t *testing.T, db *gorm.DB, userID string, schemaID uuid.UUID, name string, selectedFields datatypes.JSON, createdAt time.Time) ocr.Dataset {
	t.Helper()
	createdAt = createdAt.UTC()
	dataset := ocr.Dataset{
		UserID:         userID,
		SchemaID:       schemaID,
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt,
		Name:           name,
		SelectedFields: selectedFields,
	}
	if err := db.Create(&dataset).Error; err != nil {
		t.Fatalf("create dashboard dataset: %v", err)
	}
	return dataset
}

func TestParseRangeAndWindow(t *testing.T) {
	now := time.Date(2026, 6, 25, 15, 45, 0, 0, time.UTC)

	for _, tt := range []struct {
		name      string
		raw       string
		wantKey   RangeKey
		wantStart string
	}{
		{name: "default", raw: "", wantKey: Range30D, wantStart: "2026-05-27T00:00:00Z"},
		{name: "7d", raw: "7d", wantKey: Range7D, wantStart: "2026-06-19T00:00:00Z"},
		{name: "30d", raw: "30d", wantKey: Range30D, wantStart: "2026-05-27T00:00:00Z"},
		{name: "90d", raw: "90d", wantKey: Range90D, wantStart: "2026-03-28T00:00:00Z"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			key, err := ParseRange(tt.raw)
			if err != nil {
				t.Fatalf("ParseRange(%q): %v", tt.raw, err)
			}
			if key != tt.wantKey {
				t.Fatalf("key = %q, want %q", key, tt.wantKey)
			}
			window, err := WindowForRange(key, now)
			if err != nil {
				t.Fatalf("WindowForRange(%q): %v", key, err)
			}
			if window.StartAt.Format(time.RFC3339) != tt.wantStart {
				t.Fatalf("start = %s, want %s", window.StartAt.Format(time.RFC3339), tt.wantStart)
			}
			if !window.EndAt.Equal(now) {
				t.Fatalf("end = %s, want %s", window.EndAt.Format(time.RFC3339), now.Format(time.RFC3339))
			}
			if window.Bucket != "day" {
				t.Fatalf("bucket = %q, want day", window.Bucket)
			}
		})
	}

	if _, err := ParseRange("14d"); !errors.Is(err, ErrInvalidRange) {
		t.Fatalf("ParseRange(\"14d\") error = %v, want ErrInvalidRange", err)
	}
}

func TestLoadSummaryEmptyAccount(t *testing.T) {
	db := dashboardTestDB(t)
	user := createDashboardUser(t, db, "dashboard-empty@example.com")
	now := time.Date(2026, 6, 25, 15, 45, 0, 0, time.UTC)

	summary, err := LoadSummary(t.Context(), db, user.ID, Range30D, now)
	if err != nil {
		t.Fatalf("LoadSummary: %v", err)
	}

	if summary.Range.Key != Range30D || summary.Range.Bucket != "day" {
		t.Fatalf("range = %#v, want 30d day range", summary.Range)
	}
	if summary.Metrics != (Metrics{}) {
		t.Fatalf("metrics = %#v, want zero metrics", summary.Metrics)
	}
	if len(summary.DocumentBuckets) != 30 {
		t.Fatalf("document bucket count = %d, want 30", len(summary.DocumentBuckets))
	}
	if summary.DocumentBuckets[0] != (DocumentBucket{Date: "2026-05-27", DocumentsProcessed: 0}) {
		t.Fatalf("first bucket = %#v", summary.DocumentBuckets[0])
	}
	if summary.DocumentBuckets[29] != (DocumentBucket{Date: "2026-06-25", DocumentsProcessed: 0}) {
		t.Fatalf("last bucket = %#v", summary.DocumentBuckets[29])
	}
	if !summary.Onboarding.ShowOnboarding {
		t.Fatalf("show_onboarding = false, want true for empty account")
	}
	if summary.Warnings == nil {
		t.Fatal("warnings = nil, want empty slice")
	}
}

func TestLoadSummaryAggregatesUserThroughput(t *testing.T) {
	db := dashboardTestDB(t)
	user := createDashboardUser(t, db, "dashboard-throughput@example.com")
	other := createDashboardUser(t, db, "dashboard-throughput-other@example.com")
	now := time.Date(2026, 6, 25, 15, 45, 0, 0, time.UTC)
	schema := createDashboardSchema(t, db, user.ID, "Invoice")
	otherSchema := createDashboardSchema(t, db, other.ID, "Other")

	createDashboardJob(t, db, user.ID, now.Add(-5*time.Hour), ocr.OCRJobStatusCompleted, &schema.ID, 3)
	createDashboardJob(t, db, user.ID, now.Add(-4*time.Hour), ocr.OCRJobStatusFailed, &schema.ID, 1)
	createDashboardJob(t, db, user.ID, now.Add(-3*time.Hour), ocr.OCRJobStatusProcessing, &schema.ID, 2)
	createDashboardJob(t, db, other.ID, now.Add(-2*time.Hour), ocr.OCRJobStatusCompleted, &otherSchema.ID, 4)

	createDashboardDocument(t, db, user.ID, now.Add(-2*time.Hour), "invoice-a.pdf", &schema.ID, false, 3)
	createDashboardDocument(t, db, user.ID, now.Add(-1*time.Hour), "invoice-b.pdf", nil, true, 2)
	createDashboardDocument(t, db, user.ID, now.AddDate(0, 0, -31), "invoice-archive.pdf", &schema.ID, false, 4)
	createDashboardDocument(t, db, other.ID, now.Add(-30*time.Minute), "invoice-other.pdf", &otherSchema.ID, false, 4)

	datasetFields := datatypes.JSON([]byte(`[{"path":"/total","key":"total","label":"Total"}]`))
	createDashboardDataset(t, db, user.ID, schema.ID, "Invoices", datasetFields, now.Add(-45*time.Minute))
	createDashboardDataset(t, db, other.ID, otherSchema.ID, "Other invoices", datasetFields, now.Add(-30*time.Minute))

	bucket := createDashboardCreditBucket(t, db, user.ID, 50, now.Add(-6*time.Hour))
	createDashboardLedgerEntry(t, db, user.ID, bucket.ID, billing.CreditLedgerEntryDebit, -7, now.Add(-90*time.Minute))
	createDashboardLedgerEntry(t, db, user.ID, bucket.ID, billing.CreditLedgerEntryPurchase, 100, now.Add(-80*time.Minute))

	summary, err := LoadSummary(t.Context(), db, user.ID, Range30D, now)
	if err != nil {
		t.Fatalf("LoadSummary: %v", err)
	}

	if summary.Metrics.DocumentsProcessed != 2 {
		t.Fatalf("documents processed = %d, want 2", summary.Metrics.DocumentsProcessed)
	}
	if summary.Metrics.PagesProcessed != 5 {
		t.Fatalf("pages processed = %d, want 5", summary.Metrics.PagesProcessed)
	}
	if summary.Metrics.JobsCompleted != 1 || summary.Metrics.JobsFailed != 1 || summary.Metrics.JobsProcessing != 1 {
		t.Fatalf("job metrics = completed %d failed %d processing %d, want 1/1/1", summary.Metrics.JobsCompleted, summary.Metrics.JobsFailed, summary.Metrics.JobsProcessing)
	}
	if summary.Metrics.CompletionRate != 0.5 {
		t.Fatalf("completion rate = %f, want 0.5", summary.Metrics.CompletionRate)
	}
	if summary.Metrics.CreditsSpent != 7 || summary.CreditSummary.CreditsSpent != 7 {
		t.Fatalf("credits spent = metrics %d summary %d, want 7", summary.Metrics.CreditsSpent, summary.CreditSummary.CreditsSpent)
	}
	if summary.CreditSummary.AvailableCredits != 50 || summary.CreditSummary.LowCredit {
		t.Fatalf("credit summary = %#v, want 50 available and not low", summary.CreditSummary)
	}
	if summary.Metrics.SchemaCount != 1 {
		t.Fatalf("schema count = %d, want 1", summary.Metrics.SchemaCount)
	}
	if summary.Metrics.DatasetCount != 1 {
		t.Fatalf("dataset count = %d, want 1", summary.Metrics.DatasetCount)
	}
	if summary.DatasetSummary.TotalCount != 1 {
		t.Fatalf("dataset summary total = %d, want 1", summary.DatasetSummary.TotalCount)
	}
	if len(summary.DatasetSummary.Recent) != 1 {
		t.Fatalf("recent dataset count = %d, want 1", len(summary.DatasetSummary.Recent))
	}
	if got := summary.DatasetSummary.Recent[0]; got.Name != "Invoices" || got.SchemaName != "Invoice" || got.FieldCount != 1 {
		t.Fatalf("recent dataset = %#v, want Invoices with Invoice schema and 1 field", got)
	}
	if len(summary.SchemaThroughput) != 2 {
		t.Fatalf("schema throughput count = %d, want 2: %#v", len(summary.SchemaThroughput), summary.SchemaThroughput)
	}
	throughput := map[string]int{}
	for _, item := range summary.SchemaThroughput {
		throughput[item.SchemaName] = item.DocumentsProcessed
	}
	if throughput["Invoice"] != 1 || throughput["Inline schema"] != 1 {
		t.Fatalf("schema throughput = %#v, want Invoice and Inline schema counts", summary.SchemaThroughput)
	}
	if len(summary.RecentDocuments) != 3 {
		t.Fatalf("recent documents count = %d, want 3", len(summary.RecentDocuments))
	}
	if summary.RecentDocuments[0].OriginalFilename != "invoice-b.pdf" || summary.RecentDocuments[1].OriginalFilename != "invoice-a.pdf" {
		t.Fatalf("recent documents = %#v, want invoice-b.pdf then invoice-a.pdf", summary.RecentDocuments)
	}
	if summary.RecentDocuments[2].OriginalFilename != "invoice-archive.pdf" {
		t.Fatalf("third recent document = %q, want invoice-archive.pdf", summary.RecentDocuments[2].OriginalFilename)
	}
	if !summary.Onboarding.HasSchema || !summary.Onboarding.HasCompletedDocument || !summary.Onboarding.HasDataset || summary.Onboarding.ShowOnboarding {
		t.Fatalf("onboarding = %#v, want schema, completed document, dataset, and hidden onboarding", summary.Onboarding)
	}
}

func TestLoadSummaryRecentDocumentsHideForeignSchemaReference(t *testing.T) {
	db := dashboardTestDB(t)
	user := createDashboardUser(t, db, "dashboard-foreign-schema@example.com")
	other := createDashboardUser(t, db, "dashboard-foreign-schema-other@example.com")
	now := time.Date(2026, 6, 25, 15, 45, 0, 0, time.UTC)
	otherSchema := createDashboardSchema(t, db, other.ID, "Other")

	createDashboardDocument(t, db, user.ID, now.Add(-time.Hour), "foreign-schema.pdf", &otherSchema.ID, false, 1)

	summary, err := LoadSummary(t.Context(), db, user.ID, Range30D, now)
	if err != nil {
		t.Fatalf("LoadSummary: %v", err)
	}

	if len(summary.RecentDocuments) != 1 {
		t.Fatalf("recent documents count = %d, want 1", len(summary.RecentDocuments))
	}
	got := summary.RecentDocuments[0]
	if got.OriginalFilename != "foreign-schema.pdf" {
		t.Fatalf("recent document filename = %q, want foreign-schema.pdf", got.OriginalFilename)
	}
	if got.SchemaID != nil {
		t.Fatalf("schema id = %q, want nil for foreign schema", *got.SchemaID)
	}
	if got.SchemaName != nil {
		t.Fatalf("schema name = %q, want nil for foreign schema", *got.SchemaName)
	}
}

func TestLoadSummaryBuildsDocumentBucketsForSupportedRanges(t *testing.T) {
	db := dashboardTestDB(t)
	user := createDashboardUser(t, db, "dashboard-buckets@example.com")
	other := createDashboardUser(t, db, "dashboard-buckets-other@example.com")
	now := time.Date(2026, 6, 25, 15, 45, 0, 0, time.UTC)
	createDashboardDocument(t, db, user.ID, time.Date(2026, 6, 18, 23, 59, 0, 0, time.UTC), "before-start.pdf", nil, false, 1)
	createDashboardDocument(t, db, user.ID, time.Date(2026, 6, 19, 12, 0, 0, 0, time.UTC), "start-day.pdf", nil, false, 1)
	createDashboardDocument(t, db, user.ID, time.Date(2026, 6, 25, 12, 0, 0, 0, time.UTC), "end-day.pdf", nil, false, 1)
	createDashboardDocument(t, db, user.ID, now.Add(time.Minute), "after-now.pdf", nil, false, 1)
	createDashboardDocument(t, db, other.ID, time.Date(2026, 6, 25, 12, 30, 0, 0, time.UTC), "other-user.pdf", nil, false, 1)

	for _, tt := range []struct {
		key       RangeKey
		wantCount int
		wantFirst string
		wantLast  string
	}{
		{key: Range7D, wantCount: 7, wantFirst: "2026-06-19", wantLast: "2026-06-25"},
		{key: Range30D, wantCount: 30, wantFirst: "2026-05-27", wantLast: "2026-06-25"},
		{key: Range90D, wantCount: 90, wantFirst: "2026-03-28", wantLast: "2026-06-25"},
	} {
		t.Run(string(tt.key), func(t *testing.T) {
			summary, err := LoadSummary(t.Context(), db, user.ID, tt.key, now)
			if err != nil {
				t.Fatalf("LoadSummary: %v", err)
			}
			if len(summary.DocumentBuckets) != tt.wantCount {
				t.Fatalf("bucket count = %d, want %d", len(summary.DocumentBuckets), tt.wantCount)
			}
			if summary.DocumentBuckets[0].Date != tt.wantFirst {
				t.Fatalf("first bucket = %q, want %q", summary.DocumentBuckets[0].Date, tt.wantFirst)
			}
			if summary.DocumentBuckets[len(summary.DocumentBuckets)-1].Date != tt.wantLast {
				t.Fatalf("last bucket = %q, want %q", summary.DocumentBuckets[len(summary.DocumentBuckets)-1].Date, tt.wantLast)
			}
			if tt.key == Range7D {
				wantBuckets := []DocumentBucket{
					{Date: "2026-06-19", DocumentsProcessed: 1},
					{Date: "2026-06-20", DocumentsProcessed: 0},
					{Date: "2026-06-21", DocumentsProcessed: 0},
					{Date: "2026-06-22", DocumentsProcessed: 0},
					{Date: "2026-06-23", DocumentsProcessed: 0},
					{Date: "2026-06-24", DocumentsProcessed: 0},
					{Date: "2026-06-25", DocumentsProcessed: 1},
				}
				for i, want := range wantBuckets {
					if summary.DocumentBuckets[i] != want {
						t.Fatalf("bucket[%d] = %#v, want %#v", i, summary.DocumentBuckets[i], want)
					}
				}
				if summary.Metrics.DocumentsProcessed != 2 {
					t.Fatalf("documents processed = %d, want 2", summary.Metrics.DocumentsProcessed)
				}
				if summary.Metrics.PagesProcessed != 2 {
					t.Fatalf("pages processed = %d, want 2", summary.Metrics.PagesProcessed)
				}
			}
		})
	}
}
