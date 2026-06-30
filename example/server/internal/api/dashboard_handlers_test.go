package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/ocr"
)

func testDashboardRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := apiPostgresTx(t)
	fixedNow := time.Date(2026, 6, 25, 15, 45, 0, 0, time.UTC)
	return NewRouter(&Handler{
		DB:               db,
		InternalAPIToken: testInternalAPIToken,
		Now: func() time.Time {
			return fixedNow
		},
	}), db
}

func dashboardJSON(t *testing.T, router *gin.Engine, target string, internalToken string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	if internalToken != "" {
		req.Header.Set(internalAPIHeader, internalToken)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func decodeDashboardResponse(t *testing.T, w *httptest.ResponseRecorder) DashboardSummaryResponse {
	t.Helper()
	var response DashboardSummaryResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("decode dashboard response: %v body=%s", err, w.Body.String())
	}
	return response
}

func TestGetDashboardSummaryRejectsInvalidQueries(t *testing.T) {
	router, _ := testDashboardRouter(t)

	for _, tt := range []struct {
		name string
		path string
		want string
	}{
		{name: "missing user", path: "/api/dashboard/summary", want: "user_id is required"},
		{name: "invalid user", path: "/api/dashboard/summary?user_id=not-a-uuid", want: "invalid user_id"},
		{name: "invalid range", path: "/api/dashboard/summary?user_id=00000000-0000-0000-0000-000000000001&range=14d", want: "invalid range"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			w := dashboardJSON(t, router, tt.path, testInternalAPIToken)
			if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), tt.want) {
				t.Fatalf("status = %d body=%s, want bad request containing %q", w.Code, w.Body.String(), tt.want)
			}
		})
	}
}

func TestGetDashboardSummaryReturnsSummary(t *testing.T) {
	router, db := testDashboardRouter(t)
	user := createTestUser(t, db, "dashboard-handler@example.com")
	schema := ocr.ExtractionSchema{
		UserID:     &user.ID,
		Name:       "Invoices",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object"}`)),
		Strict:     true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}

	w := dashboardJSON(t, router, "/api/dashboard/summary?user_id="+user.ID+"&range=7d", testInternalAPIToken)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeDashboardResponse(t, w)
	if got.Range.Key != "7d" || got.Range.Bucket != "day" {
		t.Fatalf("range = %#v, want 7d day", got.Range)
	}
	if got.Range.StartAt.Format(time.RFC3339) != "2026-06-19T00:00:00Z" {
		t.Fatalf("start_at = %s, want 2026-06-19T00:00:00Z", got.Range.StartAt.Format(time.RFC3339))
	}
	if got.Range.EndAt.Format(time.RFC3339) != "2026-06-25T15:45:00Z" {
		t.Fatalf("end_at = %s, want 2026-06-25T15:45:00Z", got.Range.EndAt.Format(time.RFC3339))
	}
	if got.Metrics.SchemaCount != 1 {
		t.Fatalf("schema_count = %d, want 1", got.Metrics.SchemaCount)
	}
	if len(got.DocumentBuckets) != 7 {
		t.Fatalf("document bucket count = %d, want 7", len(got.DocumentBuckets))
	}
	if got.DocumentBuckets[0].Date != "2026-06-19" {
		t.Fatalf("first bucket date = %q, want 2026-06-19", got.DocumentBuckets[0].Date)
	}
	if got.DocumentBuckets[6].Date != "2026-06-25" {
		t.Fatalf("last bucket date = %q, want 2026-06-25", got.DocumentBuckets[6].Date)
	}
	if got.Warnings == nil {
		t.Fatal("warnings = nil, want empty array")
	}
}

func TestGetDashboardSummaryDefaultsToThirtyDayRange(t *testing.T) {
	router, db := testDashboardRouter(t)
	user := createTestUser(t, db, "dashboard-default-range@example.com")

	w := dashboardJSON(t, router, "/api/dashboard/summary?user_id="+user.ID, testInternalAPIToken)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	assertDashboardEmptyArrays(t, w)
	got := decodeDashboardResponse(t, w)
	if got.Range.Key != "30d" || got.Range.Bucket != "day" {
		t.Fatalf("range = %#v, want 30d day", got.Range)
	}
	if got.Range.StartAt.Format(time.RFC3339) != "2026-05-27T00:00:00Z" {
		t.Fatalf("start_at = %s, want 2026-05-27T00:00:00Z", got.Range.StartAt.Format(time.RFC3339))
	}
	if got.Range.EndAt.Format(time.RFC3339) != "2026-06-25T15:45:00Z" {
		t.Fatalf("end_at = %s, want 2026-06-25T15:45:00Z", got.Range.EndAt.Format(time.RFC3339))
	}
	if len(got.DocumentBuckets) != 30 {
		t.Fatalf("document bucket count = %d, want 30", len(got.DocumentBuckets))
	}
	if got.DocumentBuckets[0].Date != "2026-05-27" {
		t.Fatalf("first bucket date = %q, want 2026-05-27", got.DocumentBuckets[0].Date)
	}
	if got.DocumentBuckets[29].Date != "2026-06-25" {
		t.Fatalf("last bucket date = %q, want 2026-06-25", got.DocumentBuckets[29].Date)
	}
}

func assertDashboardEmptyArrays(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	var raw struct {
		RecentDocuments  []json.RawMessage `json:"recent_documents"`
		SchemaThroughput []json.RawMessage `json:"schema_throughput"`
		DatasetSummary   struct {
			Recent []json.RawMessage `json:"recent"`
		} `json:"dataset_summary"`
		Warnings []json.RawMessage `json:"warnings"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("decode dashboard raw response: %v body=%s", err, w.Body.String())
	}
	if raw.RecentDocuments == nil || len(raw.RecentDocuments) != 0 {
		t.Fatalf("recent_documents = %#v, want empty array", raw.RecentDocuments)
	}
	if raw.SchemaThroughput == nil || len(raw.SchemaThroughput) != 0 {
		t.Fatalf("schema_throughput = %#v, want empty array", raw.SchemaThroughput)
	}
	if raw.DatasetSummary.Recent == nil || len(raw.DatasetSummary.Recent) != 0 {
		t.Fatalf("dataset_summary.recent = %#v, want empty array", raw.DatasetSummary.Recent)
	}
	if raw.Warnings == nil || len(raw.Warnings) != 0 {
		t.Fatalf("warnings = %#v, want empty array", raw.Warnings)
	}
}
