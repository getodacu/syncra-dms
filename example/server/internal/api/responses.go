package api

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"ai.ro/syncra/internal/billing"
	"ai.ro/syncra/internal/ocr"
	"ai.ro/syncra/internal/webhooks"
)

// ErrorResponse is returned when an API request fails.
//
// swagger:model errorResponse
type ErrorResponse struct {
	// Error message.
	// required: true
	Error string `json:"error"`
}

// PaymentRequiredResponse is returned when a request needs more credits.
//
// swagger:model paymentRequiredResponse
type PaymentRequiredResponse struct {
	// Error message.
	//
	// Required: true
	Error string `json:"error"`
	// Credits required to queue the requested operation.
	//
	// Required: true
	RequiredCredits int `json:"required_credits"`
	// Credits currently available to the API key owner.
	//
	// Required: true
	AvailableCredits int `json:"available_credits"`
}

// userIDString is a UUID string.
//
// swagger:strfmt uuid
type userIDString string

// UserID is a UUID string.
//
// swagger:strfmt uuid
type UserID string

// CreditBalanceResponse describes a user's available credit balance.
//
// swagger:model creditBalanceResponse
type CreditBalanceResponse struct {
	// User id owned by the supplied public API key.
	//
	// Required: true
	UserID userIDString `json:"user_id"`
	// Available OCR credits for the API key owner.
	//
	// Required: true
	AvailableCredits int `json:"available_credits"`
}

// BillingProfileResponse describes a user's billing profile.
//
// swagger:model billingProfileResponse
type BillingProfileResponse struct {
	ID                 uuid.UUID `json:"id"`
	UserID             UserID    `json:"user_id"`
	EntityType         string    `json:"entity_type"`
	BillingName        string    `json:"billing_name"`
	BillingEmail       string    `json:"billing_email"`
	CountryCode        string    `json:"country_code"`
	AddressLine1       string    `json:"address_line1"`
	AddressLine2       *string   `json:"address_line2,omitempty"`
	City               string    `json:"city"`
	Region             *string   `json:"region,omitempty"`
	PostalCode         string    `json:"postal_code"`
	FiscalCode         *string   `json:"fiscal_code,omitempty"`
	RegistrationNumber *string   `json:"registration_number,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// BillingProfileEnvelopeResponse describes an optional billing profile.
type BillingProfileEnvelopeResponse struct {
	Profile *BillingProfileResponse `json:"profile"`
}

func billingProfileResponse(profile billing.BillingProfile) BillingProfileResponse {
	return BillingProfileResponse{
		ID:                 profile.ID,
		UserID:             UserID(profile.UserID),
		EntityType:         string(profile.EntityType),
		BillingName:        profile.BillingName,
		BillingEmail:       profile.BillingEmail,
		CountryCode:        profile.CountryCode,
		AddressLine1:       profile.AddressLine1,
		AddressLine2:       profile.AddressLine2,
		City:               profile.City,
		Region:             profile.Region,
		PostalCode:         profile.PostalCode,
		FiscalCode:         profile.FiscalCode,
		RegistrationNumber: profile.RegistrationNumber,
		CreatedAt:          profile.CreatedAt,
		UpdatedAt:          profile.UpdatedAt,
	}
}

// BillingInvoiceResponse describes a generated billing invoice.
//
// swagger:model billingInvoiceResponse
type BillingInvoiceResponse struct {
	ID                     uuid.UUID                    `json:"id"`
	UserID                 *userIDString                `json:"user_id,omitempty"`
	OrderID                *uuid.UUID                   `json:"order_id,omitempty"`
	BillingProfileID       *uuid.UUID                   `json:"billing_profile_id,omitempty"`
	BillingName            string                       `json:"billing_name"`
	BillingEmail           string                       `json:"billing_email"`
	BillingFiscalCode      *string                      `json:"billing_fiscal_code,omitempty"`
	BillingProfileSnapshot json.RawMessage              `json:"billing_profile_snapshot"`
	Lines                  []BillingInvoiceLineResponse `json:"lines"`
	NetAmount              string                       `json:"net_amount"`
	VATAmount              string                       `json:"vat_amount"`
	TotalAmount            string                       `json:"total_amount"`
	InvoiceDate            string                       `json:"invoice_date"`
	InvoiceSerie           string                       `json:"invoice_serie"`
	InvoiceNumber          int64                        `json:"invoice_number"`
	PDFPath                *string                      `json:"pdf_path,omitempty"`
	EmailDeliveryClaimedAt *time.Time                   `json:"email_delivery_claimed_at,omitempty"`
	EmailSentAt            *time.Time                   `json:"email_sent_at,omitempty"`
	CreatedAt              time.Time                    `json:"created_at"`
	UpdatedAt              time.Time                    `json:"updated_at"`
}

// BillingInvoiceLineResponse describes one generated invoice line.
//
// swagger:model billingInvoiceLineResponse
type BillingInvoiceLineResponse struct {
	Name           string `json:"name"`
	Quantity       int    `json:"quantity"`
	UnitPrice      string `json:"unit_price"`
	VATPercentage  string `json:"vat_percentage"`
	TotalVATAmount string `json:"total_vat_amount"`
	TotalAmount    string `json:"total_amount"`
}

func billingInvoiceResponse(invoice billing.BillingInvoice) BillingInvoiceResponse {
	var userID *userIDString
	if invoice.UserID != nil {
		value := userIDString(*invoice.UserID)
		userID = &value
	}
	var invoiceLines []billing.BillingInvoiceLine
	_ = json.Unmarshal(invoice.Lines, &invoiceLines)
	lines := make([]BillingInvoiceLineResponse, 0, len(invoiceLines))
	for _, line := range invoiceLines {
		lines = append(lines, BillingInvoiceLineResponse{
			Name:           line.Name,
			Quantity:       line.Quantity,
			UnitPrice:      line.UnitPrice,
			VATPercentage:  line.VATPercentage,
			TotalVATAmount: line.TotalVATAmount,
			TotalAmount:    line.TotalAmount,
		})
	}
	return BillingInvoiceResponse{
		ID:                     invoice.ID,
		UserID:                 userID,
		OrderID:                invoice.OrderID,
		BillingProfileID:       invoice.BillingProfileID,
		BillingName:            invoice.BillingName,
		BillingEmail:           invoice.BillingEmail,
		BillingFiscalCode:      invoice.BillingFiscalCode,
		BillingProfileSnapshot: json.RawMessage(invoice.BillingProfileSnapshot),
		Lines:                  lines,
		NetAmount:              invoice.NetAmount.StringFixed(2),
		VATAmount:              invoice.VATAmount.StringFixed(2),
		TotalAmount:            invoice.TotalAmount.StringFixed(2),
		InvoiceDate:            invoice.InvoiceDate.Format("2006-01-02"),
		InvoiceSerie:           invoice.InvoiceSerie,
		InvoiceNumber:          invoice.InvoiceNumber,
		PDFPath:                invoice.PDFPath,
		EmailDeliveryClaimedAt: invoice.EmailDeliveryClaimedAt,
		EmailSentAt:            invoice.EmailSentAt,
		CreatedAt:              invoice.CreatedAt,
		UpdatedAt:              invoice.UpdatedAt,
	}
}

// BillingInvoiceEmailDeliveryClaimResponse describes an invoice email delivery claim.
//
// swagger:model billingInvoiceEmailDeliveryClaimResponse
type BillingInvoiceEmailDeliveryClaimResponse struct {
	Status  string                  `json:"status"`
	Invoice *BillingInvoiceResponse `json:"invoice,omitempty"`
}

// BillingInvoiceEmailDeliverySentResponse describes a recorded invoice email delivery.
//
// swagger:model billingInvoiceEmailDeliverySentResponse
type BillingInvoiceEmailDeliverySentResponse struct {
	Status  string                 `json:"status"`
	Invoice BillingInvoiceResponse `json:"invoice"`
}

// CreditUsageHistoryEntryResponse describes one credit usage history entry.
//
// swagger:model creditUsageHistoryEntryResponse
type CreditUsageHistoryEntryResponse struct {
	ID             uuid.UUID  `json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	EntryType      string     `json:"entry_type"`
	CreditsDelta   int        `json:"credits_delta"`
	RelatedOrderID *uuid.UUID `json:"related_order_id,omitempty"`
	RelatedJobID   *uuid.UUID `json:"related_job_id,omitempty"`
}

// CreditUsageHistoryListResponse describes a cursor-paginated credit usage history list.
//
// swagger:model creditUsageHistoryListResponse
type CreditUsageHistoryListResponse struct {
	CreditUsageHistory []CreditUsageHistoryEntryResponse `json:"credit_usage_history"`
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	NextCursor *string `json:"next_cursor"`
}

func creditUsageHistoryEntryResponse(entry billing.CreditLedgerEntry) CreditUsageHistoryEntryResponse {
	return CreditUsageHistoryEntryResponse{
		ID:             entry.ID,
		CreatedAt:      entry.CreatedAt,
		EntryType:      string(entry.EntryType),
		CreditsDelta:   entry.CreditsDelta,
		RelatedOrderID: entry.RelatedOrderID,
		RelatedJobID:   entry.RelatedJobID,
	}
}

// BillingOrderResponse describes a credit purchase order.
//
// swagger:model billingOrderResponse
type BillingOrderResponse struct {
	ID                        uuid.UUID                    `json:"id"`
	UserID                    userIDString                 `json:"user_id"`
	Invoice                   *BillingOrderInvoiceResponse `json:"invoice,omitempty"`
	OrderType                 string                       `json:"order_type"`
	Status                    string                       `json:"status"`
	Provider                  string                       `json:"provider"`
	PricingTier               string                       `json:"pricing_tier"`
	UnitAmountCents           int                          `json:"unit_amount_cents"`
	Credits                   int                          `json:"credits"`
	AmountCents               int                          `json:"amount_cents"`
	Currency                  string                       `json:"currency"`
	ProviderCheckoutSessionID *string                      `json:"provider_checkout_session_id,omitempty"`
	ProviderPaymentIntentID   *string                      `json:"provider_payment_intent_id,omitempty"`
	CreatedAt                 time.Time                    `json:"created_at"`
	UpdatedAt                 time.Time                    `json:"updated_at"`
	PaidAt                    *time.Time                   `json:"paid_at,omitempty"`
	FailedAt                  *time.Time                   `json:"failed_at,omitempty"`
	RefundedAt                *time.Time                   `json:"refunded_at,omitempty"`
	CanceledAt                *time.Time                   `json:"canceled_at,omitempty"`
}

// BillingOrderInvoiceResponse describes invoice metadata attached to a billing order.
//
// swagger:model billingOrderInvoiceResponse
type BillingOrderInvoiceResponse struct {
	ID            uuid.UUID `json:"id"`
	InvoiceSerie  string    `json:"invoice_serie"`
	InvoiceNumber int64     `json:"invoice_number"`
	InvoiceDate   string    `json:"invoice_date"`
	PDFPath       *string   `json:"pdf_path,omitempty"`
}

// BillingOrderListResponse describes a cursor-paginated billing order list.
//
// swagger:model billingOrderListResponse
type BillingOrderListResponse struct {
	Orders []BillingOrderResponse `json:"orders"`
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	NextCursor *string `json:"next_cursor"`
}

func billingOrderResponse(order billing.BillingOrder) BillingOrderResponse {
	return BillingOrderResponse{
		ID:                        order.ID,
		UserID:                    userIDString(order.UserID),
		Invoice:                   billingOrderInvoiceResponse(order.Invoice),
		OrderType:                 string(order.OrderType),
		Status:                    string(order.Status),
		Provider:                  string(order.Provider),
		PricingTier:               string(order.PricingTier),
		UnitAmountCents:           order.UnitAmountCents,
		Credits:                   order.Credits,
		AmountCents:               order.AmountCents,
		Currency:                  order.Currency,
		ProviderCheckoutSessionID: order.ProviderCheckoutSessionID,
		ProviderPaymentIntentID:   order.ProviderPaymentIntentID,
		CreatedAt:                 order.CreatedAt,
		UpdatedAt:                 order.UpdatedAt,
		PaidAt:                    order.PaidAt,
		FailedAt:                  order.FailedAt,
		RefundedAt:                order.RefundedAt,
		CanceledAt:                order.CanceledAt,
	}
}

func billingOrderInvoiceResponse(invoice *billing.BillingInvoice) *BillingOrderInvoiceResponse {
	if invoice == nil {
		return nil
	}
	return &BillingOrderInvoiceResponse{
		ID:            invoice.ID,
		InvoiceSerie:  invoice.InvoiceSerie,
		InvoiceNumber: invoice.InvoiceNumber,
		InvoiceDate:   invoice.InvoiceDate.Format("2006-01-02"),
		PDFPath:       invoice.PDFPath,
	}
}

// APIKeyResponse describes a user-owned API key. The api_key field is only returned immediately after creation.
//
// swagger:model apiKeyResponse
type APIKeyResponse struct {
	ID        uuid.UUID    `json:"id"`
	UserID    userIDString `json:"user_id"`
	Name      string       `json:"name"`
	KeyPrefix string       `json:"key_prefix"`
	APIKey    string       `json:"api_key,omitempty"`
	ExpiresAt *time.Time   `json:"expires_at,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// APIKeyListResponse describes a user's API keys.
//
// swagger:model apiKeyListResponse
type APIKeyListResponse struct {
	APIKeys []APIKeyResponse `json:"api_keys"`
}

// DeleteAPIKeyResponse confirms deletion of one API key.
//
// swagger:model deleteAPIKeyResponse
type DeleteAPIKeyResponse struct {
	DeletedID    uuid.UUID `json:"deleted_id"`
	DeletedCount int       `json:"deleted_count"`
}

// WebhookResponse describes a user-owned webhook endpoint. The secret_key field is only returned immediately after creation or regeneration.
//
// swagger:model webhookResponse
type WebhookResponse struct {
	ID           uuid.UUID    `json:"id"`
	UserID       userIDString `json:"user_id"`
	URL          string       `json:"url"`
	EventsActive []string     `json:"events_active"`
	HasSecret    bool         `json:"has_secret"`
	SecretKey    string       `json:"secret_key,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// WebhookEnvelopeResponse describes an optional user webhook.
//
// swagger:model webhookEnvelopeResponse
type WebhookEnvelopeResponse struct {
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	Webhook *WebhookResponse `json:"webhook"`
}

// DeleteWebhookResponse confirms deletion of one user webhook.
//
// swagger:model deleteWebhookResponse
type DeleteWebhookResponse struct {
	DeletedID    uuid.UUID `json:"deleted_id"`
	DeletedCount int       `json:"deleted_count"`
}

func webhookResponse(hook webhooks.Webhook, secret string) WebhookResponse {
	events := webhooks.DecodeEvents(hook.EventsActive)
	outEvents := make([]string, 0, len(events))
	for _, event := range events {
		outEvents = append(outEvents, string(event))
	}
	return WebhookResponse{
		ID:           hook.ID,
		UserID:       userIDString(hook.UserID),
		URL:          hook.URL,
		EventsActive: outEvents,
		HasSecret:    hook.SecretKey != "",
		SecretKey:    secret,
		CreatedAt:    hook.CreatedAt,
		UpdatedAt:    hook.UpdatedAt,
	}
}

// SchemaResponse describes a saved extraction schema.
//
// swagger:model schemaResponse
type SchemaResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Owner user id. Missing means the entity is system-wide.
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	UserID      *optionalUserID `json:"user_id,omitempty" format:"uuid"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	// Extraction JSON Schema object.
	// required: true
	// type: object
	Schema json.RawMessage `json:"schema" swaggertype:"object"`
	Strict bool            `json:"strict"`
}

// SchemaListResponse describes a cursor-paginated extraction schema list.
//
// swagger:model schemaListResponse
type SchemaListResponse struct {
	Schemas []SchemaResponse `json:"schemas"`
	// Cursor for the next page. Null when there is no next page.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	NextCursor *string `json:"next_cursor"`
}

// JSONRecipeResponse describes an admin-managed JSON recipe.
//
// swagger:model jsonRecipeResponse
type JSONRecipeResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	// Recipe JSON Schema object.
	// required: true
	// type: object
	JSON json.RawMessage `json:"json" swaggertype:"object"`
	// Number of successful deploys.
	// required: true
	Counter int64 `json:"counter"`
	// Optional category id. Null means the recipe is shown under Others.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	CategoryID *uuid.UUID `json:"category_id"`
	// Optional category details. Null means the recipe is shown under Others.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	Category *JSONRecipeCategoryResponse `json:"category"`
}

// JSONRecipeListResponse describes a cursor-paginated JSON recipe list.
//
// swagger:model jsonRecipeListResponse
type JSONRecipeListResponse struct {
	Recipes []JSONRecipeResponse `json:"recipes"`
	// Cursor for the next page. Null when there is no next page.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	NextCursor *string `json:"next_cursor"`
}

// JSONRecipeDeployResponse describes a deployed recipe and created schema.
//
// swagger:model jsonRecipeDeployResponse
type JSONRecipeDeployResponse struct {
	Recipe JSONRecipeResponse `json:"recipe"`
	Schema SchemaResponse     `json:"schema"`
}

// JSONRecipeCategoryTitleResponse describes localized category titles.
//
// swagger:model jsonRecipeCategoryTitleResponse
type JSONRecipeCategoryTitleResponse struct {
	En string `json:"en"`
	Ro string `json:"ro"`
}

// JSONRecipeCategoryResponse describes an admin-managed JSON recipe category.
//
// swagger:model jsonRecipeCategoryResponse
type JSONRecipeCategoryResponse struct {
	ID        uuid.UUID                       `json:"id"`
	CreatedAt time.Time                       `json:"created_at"`
	UpdatedAt time.Time                       `json:"updated_at"`
	Title     JSONRecipeCategoryTitleResponse `json:"title"`
}

// JSONRecipeCategoryListResponse describes the JSON recipe categories list.
//
// swagger:model jsonRecipeCategoryListResponse
type JSONRecipeCategoryListResponse struct {
	Categories []JSONRecipeCategoryResponse `json:"categories"`
}

func jsonRecipeResponse(recipe ocr.JSONRecipe) JSONRecipeResponse {
	var category *JSONRecipeCategoryResponse
	if recipe.Category != nil {
		out := jsonRecipeCategoryResponse(*recipe.Category)
		category = &out
	}
	return JSONRecipeResponse{
		ID:          recipe.ID,
		CreatedAt:   recipe.CreatedAt,
		UpdatedAt:   recipe.UpdatedAt,
		Title:       recipe.Title,
		Description: recipe.Description,
		JSON:        json.RawMessage(recipe.JSON),
		Counter:     recipe.Counter,
		CategoryID:  recipe.CategoryID,
		Category:    category,
	}
}

func jsonRecipeCategoryResponse(category ocr.JSONRecipeCategory) JSONRecipeCategoryResponse {
	return JSONRecipeCategoryResponse{
		ID:        category.ID,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
		Title: JSONRecipeCategoryTitleResponse{
			En: category.TitleEn,
			Ro: category.TitleRo,
		},
	}
}

// collectionUserIDString is a UUID string in collection responses.
//
// swagger:strfmt uuid
type collectionUserIDString string

// CollectionResponse describes a user-owned OCR document collection.
//
// swagger:model collectionResponse
type CollectionResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Owner user id.
	UserID        collectionUserIDString `json:"user_id"`
	Name          string                 `json:"name"`
	SchemaIDs     []uuid.UUID            `json:"schema_ids"`
	SchemaCount   int                    `json:"schema_count"`
	DocumentCount int64                  `json:"document_count"`
}

// CollectionListResponse describes a cursor-paginated collection list.
//
// swagger:model collectionListResponse
type CollectionListResponse struct {
	Collections []CollectionResponse `json:"collections"`
	// Cursor for the next page. Null when there is no next page.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	NextCursor *string `json:"next_cursor"`
}

// DatasetFieldResponse describes one selected dataset field.
//
// swagger:model datasetFieldResponse
type DatasetFieldResponse struct {
	Path  string `json:"path"`
	Key   string `json:"key"`
	Label string `json:"label"`
}

// DatasetResponse describes a user-owned OCR dataset projection.
//
// swagger:model datasetResponse
type DatasetResponse struct {
	ID             uuid.UUID              `json:"id"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	UserID         collectionUserIDString `json:"user_id"`
	Name           string                 `json:"name"`
	SchemaID       uuid.UUID              `json:"schema_id"`
	SchemaName     string                 `json:"schema_name"`
	SelectedFields []DatasetFieldResponse `json:"selected_fields"`
	FieldCount     int                    `json:"field_count"`
}

// DatasetListResponse describes a cursor-paginated dataset list.
//
// swagger:model datasetListResponse
type DatasetListResponse struct {
	Datasets []DatasetResponse `json:"datasets"`
	// Cursor for the next page. Null when there is no next page.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	NextCursor *string `json:"next_cursor"`
}

// DatasetColumnResponse describes one projected dataset table column.
//
// swagger:model datasetColumnResponse
type DatasetColumnResponse struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Path  string `json:"path"`
}

// DatasetRowResponse describes one OCR document projected into a dataset row.
//
// swagger:model datasetRowResponse
type DatasetRowResponse struct {
	DocumentID uuid.UUID      `json:"document_id"`
	Filename   string         `json:"filename"`
	CreatedAt  time.Time      `json:"created_at"`
	Values     map[string]any `json:"values"`
}

// DatasetRowsResponse describes a cursor-paginated dataset row list.
//
// swagger:model datasetRowsResponse
type DatasetRowsResponse struct {
	Dataset DatasetResponse         `json:"dataset"`
	Columns []DatasetColumnResponse `json:"columns"`
	Rows    []DatasetRowResponse    `json:"rows"`
	// Cursor for the next page. Null when there is no next page.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	NextCursor *string `json:"next_cursor"`
}

// schemaIDRequestString is a UUID string in request bodies.
//
// swagger:strfmt uuid
type schemaIDRequestString string

// DeleteSchemasRequest describes extraction schema ids to delete.
//
// swagger:model deleteSchemasRequest
type DeleteSchemasRequest struct {
	// Schema ids to delete.
	// required: true
	// min items: 1
	// max items: 100
	// items:
	//   type: string
	//   format: uuid
	IDs []schemaIDRequestString `json:"ids"`
}

// DeleteSchemasResponse describes deleted extraction schema ids.
//
// swagger:model deleteSchemasResponse
type DeleteSchemasResponse struct {
	DeletedIDs   []uuid.UUID `json:"deleted_ids"`
	DeletedCount int         `json:"deleted_count"`
}

// OCRDocumentResponse describes the OCR result for an uploaded document.
//
// swagger:model ocrDocumentResponse
type OCRDocumentResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Owner user id. Missing means the entity is system-wide.
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	UserID           *optionalUserID `json:"user_id,omitempty" format:"uuid"`
	OriginalFilename string          `json:"original_filename"`
	MimeType         string          `json:"mime_type"`
	FileSize         int64           `json:"file_size"`
	PageCount        int             `json:"page_count"`
	DocumentHash     string          `json:"document_hash"`
	SchemaID         *uuid.UUID      `json:"schema_id,omitempty"`
	HasInlineSchema  bool            `json:"has_inline_schema"`
	Markdown         string          `json:"markdown"`
	// Document annotation JSON returned by Mistral.
	// type: object
	AnnotationJSON json.RawMessage `json:"annotation_json,omitempty" swaggertype:"object"`
	Cached         bool            `json:"cached"`
}

// OCRDocumentListResponse describes a cursor-paginated OCR document list.
//
// swagger:model ocrDocumentListResponse
type OCRDocumentListResponse struct {
	Documents []OCRDocumentListItemResponse `json:"documents"`
	// Cursor for the next page. Null when there is no next page.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	NextCursor *string `json:"next_cursor"`
}

// UpdateOCRDocumentRequest describes mutable OCR document metadata.
//
// swagger:model updateOCRDocumentRequest
type UpdateOCRDocumentRequest struct {
	// New document filename/title.
	// required: true
	OriginalFilename string `json:"original_filename"`
}

// DeleteOCRDocumentsRequest describes OCR document ids to delete.
//
// swagger:model deleteOCRDocumentsRequest
type DeleteOCRDocumentsRequest struct {
	IDs []string `json:"ids"`
}

// DeleteOCRDocumentsResponse describes deleted OCR document ids.
//
// swagger:model deleteOCRDocumentsResponse
type DeleteOCRDocumentsResponse struct {
	DeletedIDs   []uuid.UUID `json:"deleted_ids"`
	DeletedCount int         `json:"deleted_count"`
}

// DeleteOCRJobsRequest describes OCR job ids to delete.
//
// swagger:model deleteOCRJobsRequest
type DeleteOCRJobsRequest struct {
	IDs []string `json:"ids"`
}

// DeleteOCRJobsResponse describes deleted OCR job ids.
//
// swagger:model deleteOCRJobsResponse
type DeleteOCRJobsResponse struct {
	DeletedIDs   []uuid.UUID `json:"deleted_ids"`
	DeletedCount int         `json:"deleted_count"`
}

// MoveOCRDocumentsToCollectionsRequest describes OCR document ids and target collection ids.
//
// swagger:model moveOCRDocumentsToCollectionsRequest
type MoveOCRDocumentsToCollectionsRequest struct {
	IDs           []string `json:"ids"`
	CollectionIDs []string `json:"collection_ids,omitempty"`
}

// MoveOCRDocumentsToCollectionsResponse describes moved OCR document ids and target collection ids.
//
// swagger:model moveOCRDocumentsToCollectionsResponse
type MoveOCRDocumentsToCollectionsResponse struct {
	MovedIDs      []uuid.UUID `json:"moved_ids"`
	MovedCount    int         `json:"moved_count"`
	CollectionIDs []uuid.UUID `json:"collection_ids"`
}

// OCRDocumentCollectionSummaryResponse describes one collection attached to an OCR document.
//
// swagger:model ocrDocumentCollectionSummaryResponse
type OCRDocumentCollectionSummaryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// OCRDocumentListItemResponse describes one OCR document in a list response.
//
// swagger:model ocrDocumentListItemResponse
type OCRDocumentListItemResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Owner user id. Missing means the entity is system-wide.
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	UserID           *optionalUserID                        `json:"user_id,omitempty" format:"uuid"`
	OriginalFilename string                                 `json:"original_filename"`
	MimeType         string                                 `json:"mime_type"`
	FileSize         int64                                  `json:"file_size"`
	PageCount        int                                    `json:"page_count"`
	DocumentHash     string                                 `json:"document_hash"`
	SchemaID         *uuid.UUID                             `json:"schema_id,omitempty"`
	HasInlineSchema  bool                                   `json:"has_inline_schema"`
	Collections      []OCRDocumentCollectionSummaryResponse `json:"collections"`
}

// OCRJobResponse describes an async OCR queue job.
//
// swagger:model ocrJobResponse
type OCRJobResponse struct {
	// OCR job id.
	//
	// Required: true
	ID uuid.UUID `json:"id"`
	// Time when the OCR job was created.
	//
	// Required: true
	CreatedAt time.Time `json:"created_at"`
	// Original filename submitted with the OCR job.
	//
	// Required: true
	OriginalFilename string `json:"original_filename"`
	// Detected MIME type for the uploaded file.
	//
	// Required: true
	MimeType string `json:"mime_type"`
	// Current job status. Values are queued, processing, completed, or failed.
	//
	// Required: true
	// Enum: ["queued", "processing", "completed", "failed"]
	Status string `json:"status"`
	// Uploaded file size in bytes.
	//
	// Required: true
	FileSize int64 `json:"file_size"`
	// Number of pages known for the job. Queued jobs may report 0 until processing inspects the document.
	//
	// Required: true
	PageCount int `json:"page_count"`
	// Saved schema id used for structured extraction, when the job was created with schema_id.
	SchemaID *uuid.UUID `json:"schema_id,omitempty"`
	// Saved schema name used for structured extraction, when the job was created with schema_id.
	SchemaName *string `json:"schema_name,omitempty"`
	// Whether the job was submitted with an inline extraction schema.
	//
	// Required: true
	HasInlineSchema bool `json:"has_inline_schema"`
	// Inline JSON Schema used for structured extraction, when supplied during job creation.
	InlineSchema json.RawMessage `json:"inline_schema,omitempty" swaggertype:"object"`
	// Failure detail for failed jobs.
	ErrorMessage string `json:"error_message,omitempty"`
	// Generated document id. Null until the job completes.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	DocumentID *uuid.UUID `json:"document_id"`
}

// PublicOCRJobResponse describes a public async OCR queue job.
//
// swagger:model publicOCRJobResponse
type PublicOCRJobResponse struct {
	// OCR job id returned by POST /v1/ocr/jobs.
	//
	// Required: true
	ID uuid.UUID `json:"id"`
	// Time when the OCR job was created.
	//
	// Required: true
	CreatedAt time.Time `json:"created_at"`
	// Original filename submitted with the OCR job.
	//
	// Required: true
	OriginalFilename string `json:"original_filename"`
	// Current job status. Values are queued, processing, completed, or failed.
	//
	// Required: true
	// Enum: ["queued", "processing", "completed", "failed"]
	Status string `json:"status"`
	// Whether the job was submitted with an inline extraction schema.
	//
	// Required: true
	HasInlineSchema bool `json:"has_inline_schema"`
	// OCR document data. Omitted until the job has a loadable linked document.
	Document *PublicOCRJobDocumentResponse `json:"document,omitempty"`
}

// PublicOCRJobDocumentResponse describes the OCR document included in public job status responses.
//
// swagger:model publicOCRJobDocumentResponse
type PublicOCRJobDocumentResponse struct {
	// OCR document id linked to the completed job.
	//
	// Required: true
	DocumentID uuid.UUID `json:"document_id"`
	// Uploaded file size in bytes.
	//
	// Required: true
	FileSize int64 `json:"file_size"`
	// Number of pages reported for the OCR document.
	//
	// Required: true
	PageCount int `json:"page_count"`
	// Raw pages array from the OCR provider response. Missing provider pages are returned as an empty array.
	//
	// Required: true
	Pages []json.RawMessage `json:"pages"`
	// Raw document_annotation value from the OCR provider response. Missing provider annotations are returned as null.
	//
	// Required: true
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	DocumentAnnotation json.RawMessage `json:"document_annotation" swaggertype:"object"`
}

// OCRJobListResponse describes a cursor-paginated OCR job list.
//
// swagger:model ocrJobListResponse
type OCRJobListResponse struct {
	Jobs []OCRJobListItemResponse `json:"jobs"`
	// Cursor for the next page. Null when there is no next page.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	NextCursor *string `json:"next_cursor"`
}

// OCRJobListItemResponse describes one OCR job in a list response.
//
// swagger:model ocrJobListItemResponse
type OCRJobListItemResponse struct {
	ID               uuid.UUID  `json:"id"`
	CreatedAt        time.Time  `json:"created_at"`
	OriginalFilename string     `json:"original_filename"`
	MimeType         string     `json:"mime_type"`
	Status           string     `json:"status"`
	FileSize         int64      `json:"file_size"`
	PageCount        int        `json:"page_count"`
	SchemaID         *uuid.UUID `json:"schema_id,omitempty"`
	SchemaName       *string    `json:"schema_name,omitempty"`
	HasInlineSchema  bool       `json:"has_inline_schema"`
	// Failure detail for failed jobs.
	ErrorMessage string `json:"error_message,omitempty"`
	// Generated document id. Null until the job completes.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	DocumentID *uuid.UUID `json:"document_id"`
}

type OCRValidationResponse struct {
	OriginalFilename string     `json:"original_filename"`
	MimeType         string     `json:"mime_type"`
	FileSize         int64      `json:"file_size"`
	DocumentHash     string     `json:"document_hash"`
	HasInlineSchema  bool       `json:"has_inline_schema"`
	SchemaID         *uuid.UUID `json:"schema_id"`
}

func ocrDocumentResponse(doc ocr.OCRDocument, cached bool) OCRDocumentResponse {
	return OCRDocumentResponse{
		ID:               doc.ID,
		CreatedAt:        doc.CreatedAt,
		UpdatedAt:        doc.UpdatedAt,
		UserID:           optionalUserIDResponse(doc.UserID),
		OriginalFilename: doc.OriginalFilename,
		MimeType:         doc.MimeType,
		FileSize:         doc.FileSize,
		PageCount:        doc.PageCount,
		DocumentHash:     doc.DocumentHash,
		SchemaID:         doc.SchemaID,
		HasInlineSchema:  len(doc.InlineSchemaJSON) > 0,
		Markdown:         doc.Markdown,
		AnnotationJSON:   json.RawMessage(doc.AnnotationJSON),
		Cached:           cached,
	}
}

func ocrDocumentListItemResponse(doc ocr.OCRDocument, collections []OCRDocumentCollectionSummaryResponse) OCRDocumentListItemResponse {
	if collections == nil {
		collections = []OCRDocumentCollectionSummaryResponse{}
	}
	return OCRDocumentListItemResponse{
		ID:               doc.ID,
		CreatedAt:        doc.CreatedAt,
		UpdatedAt:        doc.UpdatedAt,
		UserID:           optionalUserIDResponse(doc.UserID),
		OriginalFilename: doc.OriginalFilename,
		MimeType:         doc.MimeType,
		FileSize:         doc.FileSize,
		PageCount:        doc.PageCount,
		DocumentHash:     doc.DocumentHash,
		SchemaID:         doc.SchemaID,
		HasInlineSchema:  len(doc.InlineSchemaJSON) > 0,
		Collections:      collections,
	}
}

func ocrJobResponse(job ocr.OCRJob) OCRJobResponse {
	return OCRJobResponse{
		ID:               job.ID,
		CreatedAt:        job.CreatedAt,
		OriginalFilename: job.OriginalFilename,
		MimeType:         job.MimeType,
		Status:           string(job.Status),
		FileSize:         job.FileSize,
		PageCount:        job.PageCount,
		SchemaID:         job.SchemaID,
		SchemaName:       ocrJobSchemaName(job),
		HasInlineSchema:  len(job.InlineSchemaJSON) > 0,
		InlineSchema:     json.RawMessage(job.InlineSchemaJSON),
		ErrorMessage:     job.ErrorMessage,
		DocumentID:       job.DocumentID,
	}
}

func publicOCRJobResponse(job ocr.OCRJob, doc *ocr.OCRDocument) PublicOCRJobResponse {
	response := PublicOCRJobResponse{
		ID:               job.ID,
		CreatedAt:        job.CreatedAt,
		OriginalFilename: job.OriginalFilename,
		Status:           string(job.Status),
		HasInlineSchema:  len(job.InlineSchemaJSON) > 0,
	}
	if doc != nil {
		pages, documentAnnotation := publicOCRJobDocumentData(doc.RawResponseJSON)
		response.Document = &PublicOCRJobDocumentResponse{
			DocumentID:         doc.ID,
			FileSize:           doc.FileSize,
			PageCount:          doc.PageCount,
			Pages:              pages,
			DocumentAnnotation: documentAnnotation,
		}
	}
	return response
}

func publicOCRJobDocumentData(raw []byte) ([]json.RawMessage, json.RawMessage) {
	data := struct {
		Pages              json.RawMessage `json:"pages"`
		DocumentAnnotation json.RawMessage `json:"document_annotation"`
	}{}
	_ = json.Unmarshal(raw, &data)

	pages := []json.RawMessage{}
	if len(data.Pages) > 0 && string(data.Pages) != "null" {
		_ = json.Unmarshal(data.Pages, &pages)
		if pages == nil {
			pages = []json.RawMessage{}
		}
	}

	documentAnnotation := json.RawMessage(`null`)
	if len(data.DocumentAnnotation) > 0 {
		documentAnnotation = data.DocumentAnnotation
	}

	return pages, documentAnnotation
}

func ocrJobListItemResponse(job ocr.OCRJob) OCRJobListItemResponse {
	return OCRJobListItemResponse{
		ID:               job.ID,
		CreatedAt:        job.CreatedAt,
		OriginalFilename: job.OriginalFilename,
		MimeType:         job.MimeType,
		Status:           string(job.Status),
		FileSize:         job.FileSize,
		PageCount:        job.PageCount,
		SchemaID:         job.SchemaID,
		SchemaName:       ocrJobSchemaName(job),
		HasInlineSchema:  len(job.InlineSchemaJSON) > 0,
		ErrorMessage:     job.ErrorMessage,
		DocumentID:       job.DocumentID,
	}
}

func ocrJobSchemaName(job ocr.OCRJob) *string {
	if job.Schema == nil {
		return nil
	}
	return &job.Schema.Name
}

// DashboardRangeResponse describes the selected dashboard time window.
//
// swagger:model dashboardRangeResponse
type DashboardRangeResponse struct {
	Key     string    `json:"key"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Bucket  string    `json:"bucket"`
}

// DashboardMetricsResponse describes top-level dashboard metrics.
//
// swagger:model dashboardMetricsResponse
type DashboardMetricsResponse struct {
	DocumentsProcessed int     `json:"documents_processed"`
	PagesProcessed     int     `json:"pages_processed"`
	JobsCompleted      int     `json:"jobs_completed"`
	JobsFailed         int     `json:"jobs_failed"`
	JobsProcessing     int     `json:"jobs_processing"`
	CompletionRate     float64 `json:"completion_rate"`
	CreditsSpent       int     `json:"credits_spent"`
	DatasetCount       int     `json:"dataset_count"`
	SchemaCount        int     `json:"schema_count"`
}

// DashboardDocumentBucketResponse describes one daily throughput bucket.
//
// swagger:model dashboardDocumentBucketResponse
type DashboardDocumentBucketResponse struct {
	Date               string `json:"date"`
	DocumentsProcessed int    `json:"documents_processed"`
}

// dashboardUUIDString is a UUID string in dashboard responses.
//
// swagger:strfmt uuid
type dashboardUUIDString string

// DashboardRecentDocumentResponse describes a recent completed OCR document.
//
// swagger:model dashboardRecentDocumentResponse
type DashboardRecentDocumentResponse struct {
	ID               dashboardUUIDString `json:"id"`
	OriginalFilename string              `json:"original_filename"`
	// Saved extraction schema id. Null when no user-owned schema is attached.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	SchemaID *dashboardUUIDString `json:"schema_id"`
	// Saved extraction schema name. Null when no user-owned schema is attached.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	SchemaName *string   `json:"schema_name"`
	PageCount  int       `json:"page_count"`
	CreatedAt  time.Time `json:"created_at"`
}

// DashboardSchemaThroughputResponse describes completed document count by schema.
//
// swagger:model dashboardSchemaThroughputResponse
type DashboardSchemaThroughputResponse struct {
	// Saved extraction schema id. Null for inline or no-schema documents.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	SchemaID           *dashboardUUIDString `json:"schema_id"`
	SchemaName         string               `json:"schema_name"`
	DocumentsProcessed int                  `json:"documents_processed"`
}

// DashboardRecentDatasetResponse describes a recent dataset projection.
//
// swagger:model dashboardRecentDatasetResponse
type DashboardRecentDatasetResponse struct {
	ID         dashboardUUIDString `json:"id"`
	Name       string              `json:"name"`
	SchemaName string              `json:"schema_name"`
	FieldCount int                 `json:"field_count"`
	CreatedAt  time.Time           `json:"created_at"`
}

// DashboardDatasetSummaryResponse describes dataset totals and recent datasets.
//
// swagger:model dashboardDatasetSummaryResponse
type DashboardDatasetSummaryResponse struct {
	TotalCount int                              `json:"total_count"`
	Recent     []DashboardRecentDatasetResponse `json:"recent"`
}

// DashboardCreditSummaryResponse describes current and range-bound credit usage.
//
// swagger:model dashboardCreditSummaryResponse
type DashboardCreditSummaryResponse struct {
	AvailableCredits int  `json:"available_credits"`
	CreditsSpent     int  `json:"credits_spent"`
	LowCredit        bool `json:"low_credit"`
}

// DashboardOnboardingResponse describes low-activity dashboard state.
//
// swagger:model dashboardOnboardingResponse
type DashboardOnboardingResponse struct {
	HasSchema            bool `json:"has_schema"`
	HasCompletedDocument bool `json:"has_completed_document"`
	HasDataset           bool `json:"has_dataset"`
	HasAPIKey            bool `json:"has_api_key"`
	HasWebhook           bool `json:"has_webhook"`
	ShowOnboarding       bool `json:"show_onboarding"`
}

// DashboardWarningResponse describes a public-safe partial data warning.
//
// swagger:model dashboardWarningResponse
type DashboardWarningResponse struct {
	Section string `json:"section"`
	Message string `json:"message"`
}

// DashboardSummaryResponse describes the authenticated app dashboard.
//
// swagger:model dashboardSummaryResponse
type DashboardSummaryResponse struct {
	Range            DashboardRangeResponse              `json:"range"`
	Metrics          DashboardMetricsResponse            `json:"metrics"`
	DocumentBuckets  []DashboardDocumentBucketResponse   `json:"document_buckets"`
	RecentDocuments  []DashboardRecentDocumentResponse   `json:"recent_documents"`
	SchemaThroughput []DashboardSchemaThroughputResponse `json:"schema_throughput"`
	DatasetSummary   DashboardDatasetSummaryResponse     `json:"dataset_summary"`
	CreditSummary    DashboardCreditSummaryResponse      `json:"credit_summary"`
	Onboarding       DashboardOnboardingResponse         `json:"onboarding"`
	Warnings         []DashboardWarningResponse          `json:"warnings"`
}
