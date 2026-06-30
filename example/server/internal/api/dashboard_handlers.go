package api

import (
	"net/http"
	"time"

	"ai.ro/syncra/internal/dashboard"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetDashboardSummary(c *gin.Context) {
	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	rangeKey, err := dashboard.ParseRange(c.Query("range"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid range")
		return
	}

	now := time.Now().UTC()
	if h.Now != nil {
		now = h.Now().UTC()
	}
	summary, err := dashboard.LoadSummary(c.Request.Context(), h.DB, userID, rangeKey, now)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to load dashboard summary")
		return
	}
	c.JSON(http.StatusOK, dashboardSummaryResponse(summary))
}

func dashboardSummaryResponse(summary dashboard.Summary) DashboardSummaryResponse {
	return DashboardSummaryResponse{
		Range: DashboardRangeResponse{
			Key:     string(summary.Range.Key),
			StartAt: summary.Range.StartAt,
			EndAt:   summary.Range.EndAt,
			Bucket:  summary.Range.Bucket,
		},
		Metrics: DashboardMetricsResponse{
			DocumentsProcessed: summary.Metrics.DocumentsProcessed,
			PagesProcessed:     summary.Metrics.PagesProcessed,
			JobsCompleted:      summary.Metrics.JobsCompleted,
			JobsFailed:         summary.Metrics.JobsFailed,
			JobsProcessing:     summary.Metrics.JobsProcessing,
			CompletionRate:     summary.Metrics.CompletionRate,
			CreditsSpent:       summary.Metrics.CreditsSpent,
			DatasetCount:       summary.Metrics.DatasetCount,
			SchemaCount:        summary.Metrics.SchemaCount,
		},
		DocumentBuckets:  dashboardDocumentBucketResponses(summary.DocumentBuckets),
		RecentDocuments:  dashboardRecentDocumentResponses(summary.RecentDocuments),
		SchemaThroughput: dashboardSchemaThroughputResponses(summary.SchemaThroughput),
		DatasetSummary: DashboardDatasetSummaryResponse{
			TotalCount: summary.DatasetSummary.TotalCount,
			Recent:     dashboardRecentDatasetResponses(summary.DatasetSummary.Recent),
		},
		CreditSummary: DashboardCreditSummaryResponse{
			AvailableCredits: summary.CreditSummary.AvailableCredits,
			CreditsSpent:     summary.CreditSummary.CreditsSpent,
			LowCredit:        summary.CreditSummary.LowCredit,
		},
		Onboarding: DashboardOnboardingResponse{
			HasSchema:            summary.Onboarding.HasSchema,
			HasCompletedDocument: summary.Onboarding.HasCompletedDocument,
			HasDataset:           summary.Onboarding.HasDataset,
			HasAPIKey:            summary.Onboarding.HasAPIKey,
			HasWebhook:           summary.Onboarding.HasWebhook,
			ShowOnboarding:       summary.Onboarding.ShowOnboarding,
		},
		Warnings: dashboardWarningResponses(summary.Warnings),
	}
}

func dashboardDocumentBucketResponses(buckets []dashboard.DocumentBucket) []DashboardDocumentBucketResponse {
	out := make([]DashboardDocumentBucketResponse, 0, len(buckets))
	for _, bucket := range buckets {
		out = append(out, DashboardDocumentBucketResponse(bucket))
	}
	return out
}

func dashboardRecentDocumentResponses(documents []dashboard.RecentDocument) []DashboardRecentDocumentResponse {
	out := make([]DashboardRecentDocumentResponse, 0, len(documents))
	for _, document := range documents {
		out = append(out, DashboardRecentDocumentResponse{
			ID:               dashboardUUIDString(document.ID),
			OriginalFilename: document.OriginalFilename,
			SchemaID:         dashboardSchemaIDResponse(document.SchemaID),
			SchemaName:       document.SchemaName,
			PageCount:        document.PageCount,
			CreatedAt:        document.CreatedAt,
		})
	}
	return out
}

func dashboardSchemaThroughputResponses(items []dashboard.SchemaThroughput) []DashboardSchemaThroughputResponse {
	out := make([]DashboardSchemaThroughputResponse, 0, len(items))
	for _, item := range items {
		out = append(out, DashboardSchemaThroughputResponse{
			SchemaID:           dashboardSchemaIDResponse(item.SchemaID),
			SchemaName:         item.SchemaName,
			DocumentsProcessed: item.DocumentsProcessed,
		})
	}
	return out
}

func dashboardSchemaIDResponse(schemaID *string) *dashboardUUIDString {
	if schemaID == nil {
		return nil
	}
	value := dashboardUUIDString(*schemaID)
	return &value
}

func dashboardRecentDatasetResponses(datasets []dashboard.RecentDataset) []DashboardRecentDatasetResponse {
	out := make([]DashboardRecentDatasetResponse, 0, len(datasets))
	for _, dataset := range datasets {
		out = append(out, DashboardRecentDatasetResponse{
			ID:         dashboardUUIDString(dataset.ID),
			Name:       dataset.Name,
			SchemaName: dataset.SchemaName,
			FieldCount: dataset.FieldCount,
			CreatedAt:  dataset.CreatedAt,
		})
	}
	return out
}

func dashboardWarningResponses(warnings []dashboard.Warning) []DashboardWarningResponse {
	out := make([]DashboardWarningResponse, 0, len(warnings))
	for _, warning := range warnings {
		out = append(out, DashboardWarningResponse{
			Section: string(warning.Section),
			Message: warning.Message,
		})
	}
	return out
}
