// Package api Syncra API
//
// Syncra document OCR API.
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 0.1
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package api

// swagger:model createCollectionRequest
type createCollectionRequestModel struct {
	// Collection name.
	// required: true
	Name string `json:"name"`
	// Owner user id.
	// required: true
	UserID schemaIDRequestString `json:"user_id"`
	// Extraction schema ids to attach to the collection.
	// max items: 100
	SchemaIDs []schemaIDRequestString `json:"schema_ids"`
}

// swagger:model updateCollectionRequest
type updateCollectionRequestModel struct {
	// Collection name.
	// required: true
	Name string `json:"name"`
	// Extraction schema ids to attach to the collection.
	// max items: 100
	SchemaIDs []schemaIDRequestString `json:"schema_ids"`
}

// swagger:model datasetFieldRequest
type datasetFieldRequestModel struct {
	// JSON pointer path into the extraction schema.
	// required: true
	Path string `json:"path"`
	// Stable column key for exported dataset values.
	// required: true
	Key string `json:"key"`
	// Display label for the selected field.
	// required: true
	Label string `json:"label"`
}

// swagger:model createDatasetRequest
type createDatasetRequestModel struct {
	// Dataset name.
	// required: true
	Name string `json:"name"`
	// Owner user id.
	// required: true
	UserID schemaIDRequestString `json:"user_id"`
	// User-owned extraction schema id.
	// required: true
	SchemaID schemaIDRequestString `json:"schema_id"`
	// Selected fields projected from OCR annotations.
	// required: true
	// min items: 1
	// max items: 100
	SelectedFields []datasetFieldRequestModel `json:"selected_fields"`
}

// swagger:model updateDatasetRequest
type updateDatasetRequestModel struct {
	// Dataset name.
	// required: true
	Name string `json:"name"`
	// User-owned extraction schema id.
	// required: true
	SchemaID schemaIDRequestString `json:"schema_id"`
	// Selected fields projected from OCR annotations.
	// required: true
	// min items: 1
	// max items: 100
	SelectedFields []datasetFieldRequestModel `json:"selected_fields"`
}

// swagger:model createBillingOrderRequest
type createBillingOrderRequestModel struct {
	// Owner user id.
	// required: true
	UserID schemaIDRequestString `json:"user_id"`
	// Number of credits to buy. Must be a positive multiple of 1000.
	// required: true
	Credits int `json:"credits"`
}

// swagger:model createBillingInvoiceLineRequest
type createBillingInvoiceLineRequestModel struct {
	// Invoice line item name.
	// required: true
	Name string `json:"name"`
	// Positive integer quantity.
	// required: true
	Quantity int `json:"quantity"`
	// Unit price as a decimal string.
	// required: true
	UnitPrice string `json:"unit_price"`
	// VAT percentage as a decimal string.
	// required: true
	VATPercentage string `json:"vat_percentage"`
}

// swagger:model createBillingInvoiceRequest
type createBillingInvoiceRequestModel struct {
	// Owner user id.
	// required: true
	UserID schemaIDRequestString `json:"user_id"`
	// Invoice series. Trimmed, uppercased, max 40 characters, using letters, numbers, underscore, or hyphen.
	// required: true
	InvoiceSerie string `json:"invoice_serie"`
	// Invoice date in YYYY-MM-DD format. Defaults to current UTC date when omitted.
	InvoiceDate *string `json:"invoice_date"`
	// Invoice lines used to compute totals.
	// required: true
	// min items: 1
	Lines []createBillingInvoiceLineRequestModel `json:"lines"`
}

// swagger:model generateInvoicePDFRequest
type generateInvoicePDFRequestModel struct {
	// Existing billing invoice id.
	// required: true
	InvoiceID schemaIDRequestString `json:"invoice_id"`
}

// swagger:model billingInvoiceEmailDeliveryRequest
type billingInvoiceEmailDeliveryRequestModel struct {
	// Owner user id. The invoice must belong to this user.
	// required: true
	UserID schemaIDRequestString `json:"user_id"`
}

// swagger:model upsertBillingProfileRequest
type upsertBillingProfileRequestModel struct {
	// Owner user id.
	// required: true
	UserID schemaIDRequestString `json:"user_id"`
	// Billing entity type.
	// required: true
	EntityType string `json:"entity_type"`
	// Legal billing name.
	// required: true
	BillingName string `json:"billing_name"`
	// Billing email address.
	// required: true
	BillingEmail string `json:"billing_email"`
	// ISO 3166-1 alpha-2 country code.
	// required: true
	CountryCode string `json:"country_code"`
	// Street address line 1.
	// required: true
	AddressLine1 string `json:"address_line1"`
	// Street address line 2.
	AddressLine2 *string `json:"address_line2"`
	// City.
	// required: true
	City string `json:"city"`
	// Region, state, or county.
	Region *string `json:"region"`
	// Postal code.
	// required: true
	PostalCode string `json:"postal_code"`
	// Tax or fiscal code.
	FiscalCode *string `json:"fiscal_code"`
	// Company registration number.
	RegistrationNumber *string `json:"registration_number"`
}

// swagger:model adminUserResponse
type adminUserResponseModel struct {
	// User id.
	// required: true
	ID schemaIDRequestString `json:"id"`
	// Display name.
	// required: true
	Name string `json:"name"`
	// Email address.
	// required: true
	Email string `json:"email"`
	// Whether the email address has been verified.
	// required: true
	EmailVerified bool `json:"email_verified"`
	// User role. Values are user or admin. Role editing is not exposed by admin user management v1.
	// required: true
	Role string `json:"role"`
	// Optional profile image URL.
	Image *string `json:"image"`
	// Creation timestamp.
	// required: true
	CreatedAt dateTimeString `json:"created_at"`
	// Last update timestamp.
	// required: true
	UpdatedAt dateTimeString `json:"updated_at"`
	// Last successful email/password login timestamp. Null when the user has not logged in.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	LastLoginAt *dateTimeString `json:"last_login_at"`
}

// swagger:model adminUserDetailResponse
type adminUserDetailResponseModel struct {
	adminUserResponseModel
	// Current available credit balance.
	// required: true
	AvailableCredits int `json:"available_credits"`
	// Billing profile. Null when no billing profile has been saved.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	BillingProfile *BillingProfileResponse `json:"billing_profile"`
}

// swagger:model adminUserListResponse
type adminUserListResponseModel struct {
	// Users on this page.
	// required: true
	Users []adminUserResponseModel `json:"users"`
	// Cursor for the next page. Null when no additional page is available.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	NextCursor *string `json:"next_cursor"`
}

// swagger:model adminBillingOrderUserResponse
type adminBillingOrderUserResponseModel struct {
	// User id.
	// required: true
	ID schemaIDRequestString `json:"id"`
	// Display name.
	// required: true
	Name string `json:"name"`
	// Email address.
	// required: true
	Email string `json:"email"`
}

// swagger:model adminBillingOrderResponse
type adminBillingOrderResponseModel struct {
	// Order id.
	// required: true
	ID schemaIDRequestString `json:"id"`
	// Owner user id.
	// required: true
	UserID schemaIDRequestString `json:"user_id"`
	// Owner user summary.
	// required: true
	User adminBillingOrderUserResponseModel `json:"user"`
	// Order type.
	// required: true
	OrderType string `json:"order_type"`
	// Order status. Values are pending, paid, failed, refunded, or canceled.
	// required: true
	Status string `json:"status"`
	// Billing provider.
	// required: true
	Provider string `json:"provider"`
	// Pricing tier id.
	// required: true
	PricingTier string `json:"pricing_tier"`
	// Unit amount in cents.
	// required: true
	UnitAmountCents int `json:"unit_amount_cents"`
	// Credits purchased.
	// required: true
	Credits int `json:"credits"`
	// Total amount in cents.
	// required: true
	AmountCents int `json:"amount_cents"`
	// Currency code.
	// required: true
	Currency string `json:"currency"`
	// Provider checkout session id.
	ProviderCheckoutSessionID *string `json:"provider_checkout_session_id"`
	// Provider payment intent id.
	ProviderPaymentIntentID *string `json:"provider_payment_intent_id"`
	// Creation timestamp.
	// required: true
	CreatedAt dateTimeString `json:"created_at"`
	// Last update timestamp.
	// required: true
	UpdatedAt dateTimeString `json:"updated_at"`
	// Payment timestamp.
	PaidAt *dateTimeString `json:"paid_at"`
	// Failure timestamp.
	FailedAt *dateTimeString `json:"failed_at"`
	// Refund timestamp.
	RefundedAt *dateTimeString `json:"refunded_at"`
	// Cancellation timestamp.
	CanceledAt *dateTimeString `json:"canceled_at"`
}

// swagger:model adminBillingOrderListResponse
type adminBillingOrderListResponseModel struct {
	// Orders on this page.
	// required: true
	Orders []adminBillingOrderResponseModel `json:"orders"`
	// Cursor for the next page. Null when no additional page is available.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	NextCursor *string `json:"next_cursor"`
}

// swagger:model adminBillingInvoiceListResponse
type adminBillingInvoiceListResponseModel struct {
	// Invoices on this page.
	// required: true
	Invoices []BillingInvoiceResponse `json:"invoices"`
	// Cursor for the next page. Null when no additional page is available.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	NextCursor *string `json:"next_cursor"`
}

// swagger:model patchAdminUserRequest
type patchAdminUserRequestModel struct {
	// Display name to set. Omit to leave unchanged.
	Name *string `json:"name"`
	// Email address to set. Omit to leave unchanged. Email verification state is preserved.
	Email *string `json:"email"`
}

// swagger:model adminSetPasswordRequest
type adminSetPasswordRequestModel struct {
	// New password from 8 to 128 characters.
	// required: true
	Password string `json:"password"`
}

// swagger:model adminSetPasswordResponse
type adminSetPasswordResponseModel struct {
	// True when the password was reset and existing target-user sessions were invalidated.
	// required: true
	OK bool `json:"ok"`
}

// swagger:model adminAdjustUserBalanceRequest
type adminAdjustUserBalanceRequestModel struct {
	// Positive to add credits, negative to subtract credits. Zero is rejected.
	// required: true
	CreditsDelta int `json:"credits_delta"`
}

// swagger:model adminUpsertBillingProfileRequest
type adminUpsertBillingProfileRequestModel struct {
	// Billing entity type.
	// required: true
	EntityType string `json:"entity_type"`
	// Legal billing name.
	// required: true
	BillingName string `json:"billing_name"`
	// Billing email address.
	// required: true
	BillingEmail string `json:"billing_email"`
	// ISO 3166-1 alpha-2 country code.
	// required: true
	CountryCode string `json:"country_code"`
	// Street address line 1.
	// required: true
	AddressLine1 string `json:"address_line1"`
	// Street address line 2.
	AddressLine2 *string `json:"address_line2"`
	// City.
	// required: true
	City string `json:"city"`
	// Region, state, or county.
	Region *string `json:"region"`
	// Postal code.
	// required: true
	PostalCode string `json:"postal_code"`
	// Tax or fiscal code.
	FiscalCode *string `json:"fiscal_code"`
	// Company registration number.
	RegistrationNumber *string `json:"registration_number"`
}

// swagger:model billingProfileEnvelopeResponse
type billingProfileEnvelopeResponseModel struct {
	// Billing profile. Null when no billing profile has been saved.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	Profile *struct {
		// Billing profile id.
		ID schemaIDRequestString `json:"id"`
		// Owner user id.
		UserID schemaIDRequestString `json:"user_id"`
		// Billing entity type.
		EntityType string `json:"entity_type"`
		// Legal billing name.
		BillingName string `json:"billing_name"`
		// Billing email address.
		BillingEmail string `json:"billing_email"`
		// ISO 3166-1 alpha-2 country code.
		CountryCode string `json:"country_code"`
		// Street address line 1.
		AddressLine1 string `json:"address_line1"`
		// Street address line 2.
		AddressLine2 *string `json:"address_line2"`
		// City.
		City string `json:"city"`
		// Region, state, or county.
		Region *string `json:"region"`
		// Postal code.
		PostalCode string `json:"postal_code"`
		// Tax or fiscal code.
		FiscalCode *string `json:"fiscal_code"`
		// Company registration number.
		RegistrationNumber *string `json:"registration_number"`
		// Creation timestamp.
		CreatedAt dateTimeString `json:"created_at"`
		// Last update timestamp.
		UpdatedAt dateTimeString `json:"updated_at"`
	} `json:"profile"`
}

// swagger:model attachCheckoutSessionRequest
type attachCheckoutSessionRequestModel struct {
	// Stripe Checkout Session id.
	// required: true
	CheckoutSessionID string `json:"checkout_session_id"`
}

// dateTimeString is an RFC3339 timestamp string.
//
// swagger:strfmt date-time
type dateTimeString string

// swagger:model markBillingOrderPaidRequest
type markBillingOrderPaidRequestModel struct {
	// Stripe Checkout Session id.
	CheckoutSessionID *string `json:"checkout_session_id"`
	// Stripe Payment Intent id.
	PaymentIntentID *string `json:"payment_intent_id"`
	// Payment completion timestamp.
	// required: true
	PaidAt dateTimeString `json:"paid_at"`
}

// swagger:model createAPIKeyRequest
type createAPIKeyRequestModel struct {
	// Owner user id.
	// required: true
	UserID schemaIDRequestString `json:"user_id"`
	// Human-readable API key name.
	// required: true
	Name string `json:"name"`
	// Optional expiration timestamp. Null or omitted means the key does not expire.
	ExpiresAt *dateTimeString `json:"expires_at"`
}

// swagger:model upsertWebhookRequest
type upsertWebhookRequestModel struct {
	// Owner user id.
	// required: true
	UserID schemaIDRequestString `json:"user_id"`
	// Absolute http or https webhook endpoint URL.
	// required: true
	URL string `json:"url"`
	// Active webhook events. Supported values are job.started, job.failed, and job.succeeded. Omit or pass an empty list for no active events.
	EventsActive []string `json:"events_active"`
}

func swaggerOperations() {
	// swagger:operation GET /api/dashboard/summary dashboard getDashboardSummary
	//
	// Get dashboard summary.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Owner user id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: range
	//   in: query
	//   description: Dashboard range. Defaults to 30d.
	//   required: false
	//   type: string
	//   enum:
	//   - 7d
	//   - 30d
	//   - 90d
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/dashboardSummaryResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getDashboardSummary"

	// swagger:operation POST /api/ocr/schemas schemas createSchema
	//
	// Create extraction schema.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: schema
	//   in: body
	//   description: Extraction schema with a valid JSON Schema object.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/createSchemaRequest"
	// responses:
	//   '201':
	//     description: Created.
	//     schema:
	//       "$ref": "#/definitions/schemaResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "createSchema"

	// swagger:operation GET /api/ocr/schemas schemas listSchemas
	//
	// List extraction schemas.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Optional owner user id. Omit or leave empty to list system-wide schemas.
	//   required: false
	//   type: string
	//   format: uuid
	//   x-nullable: true
	// - name: cursor
	//   in: query
	//   description: Cursor returned by a previous page.
	//   required: false
	//   type: string
	// - name: size
	//   in: query
	//   description: Page size from 1 to 100. Defaults to 20.
	//   required: false
	//   type: integer
	// - name: sort
	//   in: query
	//   description: Sort direction by creation time. Defaults to desc.
	//   required: false
	//   type: string
	//   enum:
	//   - asc
	//   - desc
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/schemaListResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listSchemas"

	// swagger:operation GET /api/ocr/schemas/{id} schemas getSchema
	//
	// Get extraction schema.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Schema id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Optional owner user id used to scope schema lookup. Omit to use id-only lookup.
	//   required: false
	//   type: string
	//   format: uuid
	//   x-nullable: true
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/schemaResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getSchema"

	// swagger:operation PUT /api/ocr/schemas/{id} schemas updateSchema
	//
	// Update extraction schema.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Schema id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Optional owner user id used to scope schema update. Omit or leave empty to update system-wide schemas.
	//   required: false
	//   type: string
	//   format: uuid
	//   x-nullable: true
	// - name: schema
	//   in: body
	//   description: Updated extraction schema with a valid JSON Schema object.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/updateSchemaRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/schemaResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "updateSchema"

	// swagger:operation DELETE /api/ocr/schemas schemas deleteSchemas
	//
	// Delete extraction schemas.
	//
	// Deletes scoped extraction schemas by id. Unknown ids and ids outside the requested scope are omitted from the response.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Optional owner user id. Omit or leave empty to delete system-wide schemas.
	//   required: false
	//   type: string
	//   format: uuid
	//   x-nullable: true
	// - name: body
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/deleteSchemasRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/deleteSchemasResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "deleteSchemas"

	// swagger:operation DELETE /api/ocr/schemas/{id} schemas deleteSchema
	//
	// Delete extraction schema.
	//
	// Deletes a scoped extraction schema by id.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Schema id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Optional owner user id used to scope schema deletion. Omit or leave empty to delete system-wide schemas.
	//   required: false
	//   type: string
	//   format: uuid
	//   x-nullable: true
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/deleteSchemasResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "deleteSchema"

	// swagger:operation POST /api/ocr ocr createOCRDocument
	//
	// OCR a document.
	//
	// Upload a PDF, PNG, or JPEG document for OCR. Provide either an inline schema or a saved schema id, but not both. When a schema is provided, Mistral must return valid document annotation JSON. Upload size is limited by MAX_UPLOAD_BYTES.
	//
	// ---
	// consumes:
	// - multipart/form-data
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: file
	//   in: formData
	//   description: PDF, PNG, or JPEG document. Filename must be at most 255 characters.
	//   required: true
	//   type: file
	// - name: schema
	//   in: formData
	//   description: Inline JSON Schema. Mutually exclusive with schema_id.
	//   required: false
	//   type: string
	// - name: schema_id
	//   in: formData
	//   description: Saved schema id. Mutually exclusive with schema.
	//   required: false
	//   type: string
	// - name: user_id
	//   in: formData
	//   description: Optional owner user id. Omit or leave empty to create/use system-wide OCR documents.
	//   required: false
	//   type: string
	//   format: uuid
	//   x-nullable: true
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/ocrDocumentResponse"
	//   '201':
	//     description: Created.
	//     schema:
	//       "$ref": "#/definitions/ocrDocumentResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '502':
	//     description: Bad gateway.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "createOCRDocument"

	// swagger:operation GET /api/ocr/documents documents listOCRDocuments
	//
	// List OCR documents.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Optional owner user id. Omit or leave empty to list system-wide OCR documents.
	//   required: false
	//   type: string
	//   format: uuid
	//   x-nullable: true
	// - name: collection_id
	//   in: query
	//   description: Optional collection id filter. Requires user_id.
	//   required: false
	//   type: string
	//   format: uuid
	// - name: filename
	//   in: query
	//   description: Optional case-insensitive partial original filename filter.
	//   required: false
	//   type: string
	// - name: created_from
	//   in: query
	//   description: Inclusive lower created_at bound as an RFC3339 timestamp.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: created_to
	//   in: query
	//   description: Inclusive upper created_at bound as an RFC3339 timestamp.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: sort
	//   in: query
	//   description: Created-at sort direction. Defaults to desc.
	//   required: false
	//   type: string
	// - name: cursor
	//   in: query
	//   description: Opaque cursor returned by the previous response.
	//   required: false
	//   type: string
	// - name: size
	//   in: query
	//   description: Page size from 1 to 100. Defaults to 20.
	//   required: false
	//   type: integer
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/ocrDocumentListResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listOCRDocuments"

	// swagger:operation DELETE /api/ocr/documents documents deleteOCRDocuments
	//
	// Delete OCR documents.
	//
	// Soft-deletes scoped OCR documents by id. Unknown or already-deleted ids are omitted from the response.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Optional owner user id. Omit or leave empty to delete system-wide OCR documents.
	//   required: false
	//   type: string
	//   format: uuid
	//   x-nullable: true
	// - name: body
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/deleteOCRDocumentsRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/deleteOCRDocumentsResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "deleteOCRDocuments"

	// swagger:operation PUT /api/ocr/documents/collections documents moveOCRDocumentsToCollections
	//
	// Move OCR documents to collections.
	//
	// Replaces all collection associations for scoped OCR documents. Empty or omitted collection_ids removes documents from all collections.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Owner user id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: body
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/moveOCRDocumentsToCollectionsRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/moveOCRDocumentsToCollectionsResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Collection not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "moveOCRDocumentsToCollections"

	// swagger:operation PATCH /api/ocr/documents/{id} documents updateOCRDocument
	//
	// Update an OCR document.
	//
	// Updates mutable OCR document metadata, currently the displayed filename/title.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: OCR document id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Optional owner user id used to scope document update. Omit or leave empty to update system-wide OCR documents.
	//   required: false
	//   type: string
	//   format: uuid
	//   x-nullable: true
	// - name: body
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/updateOCRDocumentRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/ocrDocumentResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "updateOCRDocument"

	// swagger:operation DELETE /api/ocr/documents/{id} documents deleteOCRDocument
	//
	// Delete an OCR document.
	//
	// Soft-deletes an OCR document by id. The document is hidden from normal reads and list responses after deletion.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: OCR document id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Optional owner user id used to scope document deletion. Omit to use the legacy unscoped delete behavior.
	//   required: false
	//   type: string
	//   format: uuid
	//   x-nullable: true
	// responses:
	//   '204':
	//     description: No content.
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "deleteOCRDocument"

	// swagger:operation GET /api/ocr/document/{id} documents getOCRDocument
	//
	// Get an OCR document.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: OCR document id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Optional owner user id used to scope document lookup.
	//   required: false
	//   type: string
	//   format: uuid
	//   x-nullable: true
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/ocrDocumentResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getOCRDocument"

	// swagger:operation POST /api/ocr/jobs jobs createOCRJob
	//
	// Queue a document for async OCR.
	//
	// Upload a PDF, PNG, or JPEG document for async OCR. Provide either an inline schema or a saved schema id, but not both. This endpoint only registers queued work and stores the uploaded file; it does not execute OCR. Uploads without a schema may contain at most 1000 pages; uploads with an inline or saved schema may contain at most 150 pages.
	//
	// ---
	// consumes:
	// - multipart/form-data
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: file
	//   in: formData
	//   description: PDF, PNG, or JPEG document. Filename must be at most 255 characters.
	//   required: true
	//   type: file
	// - name: schema
	//   in: formData
	//   description: Inline JSON Schema. Mutually exclusive with schema_id.
	//   required: false
	//   type: string
	// - name: schema_id
	//   in: formData
	//   description: Saved schema id. Mutually exclusive with schema.
	//   required: false
	//   type: string
	// - name: user_id
	//   in: formData
	//   description: Owner user id used for credit checks.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '202':
	//     description: Accepted.
	//     schema:
	//       "$ref": "#/definitions/ocrJobResponse"
	//   '400':
	//     description: Bad request. Returned for invalid upload input or when the document exceeds the page limit.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '402':
	//     description: Payment required.
	//     schema:
	//       "$ref": "#/definitions/paymentRequiredResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "createOCRJob"

	// swagger:operation GET /api/ocr/jobs jobs listOCRJobs
	//
	// List async OCR jobs.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Optional owner user id. Omit or leave empty to list system-wide OCR jobs.
	//   required: false
	//   type: string
	//   format: uuid
	//   x-nullable: true
	// - name: status
	//   in: query
	//   description: Optional job status filter. One of queued, processing, completed, or failed.
	//   required: false
	//   type: string
	// - name: created_from
	//   in: query
	//   description: Inclusive lower created_at bound as an RFC3339 timestamp.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: created_to
	//   in: query
	//   description: Inclusive upper created_at bound as an RFC3339 timestamp.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: sort
	//   in: query
	//   description: Created-at sort direction. Defaults to desc.
	//   required: false
	//   type: string
	// - name: cursor
	//   in: query
	//   description: Opaque cursor returned by the previous response.
	//   required: false
	//   type: string
	// - name: size
	//   in: query
	//   description: Page size from 1 to 100. Defaults to 20.
	//   required: false
	//   type: integer
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/ocrJobListResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listOCRJobs"

	// swagger:operation DELETE /api/ocr/jobs jobs deleteOCRJobs
	//
	// Delete OCR jobs.
	//
	// Soft-deletes scoped OCR jobs by id. Unknown or already-deleted ids are omitted from the response.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Optional owner user id. Omit or leave empty to delete system-wide OCR jobs.
	//   required: false
	//   type: string
	//   format: uuid
	//   x-nullable: true
	// - name: request
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/deleteOCRJobsRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/deleteOCRJobsResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "deleteOCRJobs"

	// swagger:operation GET /api/ocr/jobs/{id} jobs getOCRJob
	//
	// Get async OCR job status.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: OCR job id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Optional owner user id used to scope polling.
	//   required: false
	//   type: string
	//   format: uuid
	//   x-nullable: true
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/ocrJobResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getOCRJob"

	// swagger:operation POST /v1/ocr/jobs public-ocr createPublicOCRJob
	//
	// Queue a document for async OCR using an API key.
	//
	// ---
	// summary: Queue public async OCR job
	// description: |
	//   Upload a PDF, PNG, or JPEG document as `multipart/form-data` and queue it for asynchronous OCR processing.
	//
	//   Authentication uses the public API key in the `Authorization` header. The header accepts either the raw API key value or `Bearer <api_key>`.
	//
	//   Ownership is derived from the API key. The request must not include `user_id`; user_id is not accepted in either the query string or the form body.
	//
	//   Extraction schema options:
	//   - Omit both `schema` and `schema_id` to run OCR without structured extraction.
	//   - Send `schema` with an inline JSON Schema object when the extraction shape is request-specific.
	//   - Send `schema_id` with a saved schema UUID owned by the same API key owner.
	//   - schema and schema_id are mutually exclusive.
	//
	//   Page limits:
	//   - Without a schema, uploads may contain at most 1000 pages.
	//   - With an inline or saved schema, uploads may contain at most 150 pages.
	//
	//   The endpoint returns `202 Accepted` after the file and schema options are validated, one credit is reserved, and the job is stored with status `queued`. OCR processing happens asynchronously. Poll GET /v1/ocr/jobs/{id} with the returned `id` to track status and retrieve the completed public document response.
	//
	//   Accepted job response example:
	//
	//   ```json
	//   {
	//     "id": "f2415ed5-5d0b-4c8f-8da4-46fd5cf8f7cb",
	//     "created_at": "2026-06-10T09:24:42Z",
	//     "status": "queued",
	//     "original_filename": "invoice.pdf",
	//     "mime_type": "application/pdf",
	//     "file_size": 123456,
	//     "page_count": 0,
	//     "has_inline_schema": true,
	//     "inline_schema": {
	//       "type": "object",
	//       "properties": {
	//         "invoice_number": {
	//           "type": "string"
	//         }
	//       }
	//     },
	//     "document_id": null
	//   }
	//   ```
	// consumes:
	// - multipart/form-data
	// produces:
	// - application/json
	// parameters:
	// - name: Authorization
	//   in: header
	//   description: Public API key for the OCR job owner. Use either the raw API key or `Bearer <api_key>`.
	//   required: true
	//   type: string
	// - name: file
	//   in: formData
	//   description: PDF, PNG, or JPEG document to process. Filename must be at most 255 characters.
	//   required: true
	//   type: file
	// - name: schema
	//   in: formData
	//   description: Inline JSON Schema object for structured extraction; mutually exclusive with schema_id.
	//   required: false
	//   type: string
	// - name: schema_id
	//   in: formData
	//   description: saved schema UUID owned by the same API key owner; mutually exclusive with schema.
	//   required: false
	//   type: string
	//   format: uuid
	// responses:
	//   '202':
	//     description: Accepted. The OCR job is stored with status queued. Poll GET /v1/ocr/jobs/{id} with the returned id for status and results.
	//     schema:
	//       "$ref": "#/definitions/ocrJobResponse"
	//     examples:
	//       application/json:
	//         id: f2415ed5-5d0b-4c8f-8da4-46fd5cf8f7cb
	//         created_at: "2026-06-10T09:24:42Z"
	//         status: queued
	//         original_filename: invoice.pdf
	//         mime_type: application/pdf
	//         file_size: 123456
	//         page_count: 0
	//         has_inline_schema: true
	//         inline_schema:
	//           type: object
	//           properties:
	//             invoice_number:
	//               type: string
	//   '400':
	//     description: Bad request. Returned for a missing file, unsupported file type, filename longer than 255 characters, invalid inline schema JSON, both schema and schema_id being supplied, forbidden user_id input, or an upload exceeding the page limit.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized. Returned when the Authorization header is missing or the API key is invalid.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '402':
	//     description: Payment required. Returned when the API key owner does not have enough available credits to queue the job.
	//     schema:
	//       "$ref": "#/definitions/paymentRequiredResponse"
	//   '404':
	//     description: Not found. Returned when schema_id references a saved schema that does not exist or is not owned by the API key owner.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error. Returned for an unexpected server-side failure while storing the OCR job.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "createPublicOCRJob"

	// swagger:operation GET /v1/ocr/jobs/{id} public-ocr getPublicOCRJob
	//
	// Get public async OCR job status.
	//
	// ---
	// summary: Get public async OCR job status
	// description: |
	//   Poll this endpoint with the job id returned by `POST /v1/ocr/jobs` to track an asynchronous OCR request.
	//
	//   Authentication uses the public API key in the `Authorization` header. The header accepts either the raw API key value or `Bearer <api_key>`.
	//
	//   Jobs are scoped to the API key owner. A valid API key can only retrieve jobs created by the same owner; another owner's job id is returned as `404`.
	//
	//   Status values:
	//   - `queued`: the job was accepted and is waiting for OCR processing.
	//   - `processing`: OCR processing has started and the result document is not available yet.
	//   - `completed`: OCR processing finished successfully. When the linked document row is available, the response includes `document`.
	//   - `failed`: OCR processing did not complete successfully. Public responses do not expose internal failure details.
	//
	//   The top-level response always contains `id`, `created_at`, `status`, `original_filename`, and `has_inline_schema`. The `document` object is omitted until a completed job has a loadable linked document. If the linked document was deleted or cannot be loaded, the request still returns `200` and omits `document`.
	//
	//   Completed job response example:
	//
	//   ```json
	//   {
	//     "id": "f2415ed5-5d0b-4c8f-8da4-46fd5cf8f7cb",
	//     "created_at": "2026-06-10T09:24:42Z",
	//     "status": "completed",
	//     "original_filename": "invoice.pdf",
	//     "has_inline_schema": true,
	//     "document": {
	//       "document_id": "7db2e1e4-9552-4616-9e89-7b452b1f7793",
	//       "file_size": 123456,
	//       "page_count": 1,
	//       "pages": [
	//         {
	//           "page_number": 1,
	//           "markdown": "# Invoice\n\nTotal: 100.00"
	//         }
	//       ],
	//       "document_annotation": {
	//         "invoice_number": "INV-1001",
	//         "total": 100
	//       }
	//     }
	//   }
	//   ```
	//
	//   Queued job response example:
	//
	//   ```json
	//   {
	//     "id": "f2415ed5-5d0b-4c8f-8da4-46fd5cf8f7cb",
	//     "created_at": "2026-06-10T09:24:42Z",
	//     "status": "queued",
	//     "original_filename": "invoice.pdf",
	//     "has_inline_schema": true
	//   }
	//   ```
	// produces:
	// - application/json
	// parameters:
	// - name: Authorization
	//   in: header
	//   description: Public API key for the owner that created the OCR job. Use either the raw API key or `Bearer <api_key>`.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: OCR job id returned by `POST /v1/ocr/jobs`.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: Job status for a job owned by the API key. The document is included only when the job has a loadable linked OCR document.
	//     schema:
	//       "$ref": "#/definitions/publicOCRJobResponse"
	//     examples:
	//       application/json:
	//         id: f2415ed5-5d0b-4c8f-8da4-46fd5cf8f7cb
	//         created_at: "2026-06-10T09:24:42Z"
	//         status: completed
	//         original_filename: invoice.pdf
	//         has_inline_schema: true
	//         document:
	//           document_id: 7db2e1e4-9552-4616-9e89-7b452b1f7793
	//           file_size: 123456
	//           page_count: 1
	//           pages:
	//           - page_number: 1
	//             markdown: "# Invoice\n\nTotal: 100.00"
	//           document_annotation:
	//             invoice_number: INV-1001
	//             total: 100
	//   '400':
	//     description: Bad request. Returned when the job id is missing or invalid, including ids that are not valid UUIDs.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized. Returned when the Authorization header is missing or the API key is invalid.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found. Returned when the job does not exist or is not owned by the API key owner.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error. Returned for an unexpected server-side failure while loading job status.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getPublicOCRJob"

	// swagger:operation POST /api/collection collections createCollection
	//
	// Create a document collection.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: collection
	//   in: body
	//   description: Collection name, owner user id, and optional extraction schema ids.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/createCollectionRequest"
	// responses:
	//   '201':
	//     description: Created.
	//     schema:
	//       "$ref": "#/definitions/collectionResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "createCollection"

	// swagger:operation GET /api/collections collections listCollections
	//
	// List document collections.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Required owner user id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: cursor
	//   in: query
	//   description: Cursor returned by a previous page.
	//   required: false
	//   type: string
	// - name: size
	//   in: query
	//   description: Page size from 1 to 100. Defaults to 20.
	//   required: false
	//   type: integer
	// - name: sort
	//   in: query
	//   description: Sort direction by creation time. Defaults to desc.
	//   required: false
	//   type: string
	//   enum:
	//   - asc
	//   - desc
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/collectionListResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listCollections"

	// swagger:operation GET /api/collections/{id} collections getCollection
	//
	// Get a document collection.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Collection id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Required owner user id used to scope collection lookup.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/collectionResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getCollection"

	// swagger:operation PUT /api/collection/{id} collections updateCollection
	//
	// Update a document collection.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Collection id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Required owner user id used to scope collection update.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: collection
	//   in: body
	//   description: Updated collection name and extraction schema ids.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/updateCollectionRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/collectionResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "updateCollection"

	// swagger:operation DELETE /api/collection/{id} collections deleteCollection
	//
	// Delete a document collection.
	//
	// Deletes a collection and its collection links for the requested owner.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Collection id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Required owner user id used to scope collection deletion.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '204':
	//     description: No content.
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "deleteCollection"

	// swagger:operation POST /api/datasets datasets createDataset
	//
	// Create a dataset.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: dataset
	//   in: body
	//   description: Dataset name, owner user id, schema id, and selected fields.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/createDatasetRequest"
	// responses:
	//   '201':
	//     description: Created.
	//     schema:
	//       "$ref": "#/definitions/datasetResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "createDataset"

	// swagger:operation GET /api/datasets datasets listDatasets
	//
	// List datasets.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Required owner user id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: cursor
	//   in: query
	//   description: Cursor returned by a previous page.
	//   required: false
	//   type: string
	// - name: size
	//   in: query
	//   description: Page size from 1 to 100. Defaults to 20.
	//   required: false
	//   type: integer
	// - name: sort
	//   in: query
	//   description: Sort direction by creation time. Defaults to desc.
	//   required: false
	//   type: string
	//   enum:
	//   - asc
	//   - desc
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/datasetListResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listDatasets"

	// swagger:operation GET /api/datasets/{id}/rows datasets getDatasetRows
	//
	// List dataset rows.
	//
	// Projects matching OCR documents into the dataset's selected field columns.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Dataset id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Required owner user id used to scope dataset lookup.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: created_from
	//   in: query
	//   description: Include source documents created at or after this RFC3339 timestamp.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: created_to
	//   in: query
	//   description: Include source documents created at or before this RFC3339 timestamp.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: cursor
	//   in: query
	//   description: Cursor returned by a previous page.
	//   required: false
	//   type: string
	// - name: size
	//   in: query
	//   description: Page size from 1 to 100. Defaults to 20.
	//   required: false
	//   type: integer
	// - name: sort
	//   in: query
	//   description: Sort direction by source document creation time. Defaults to desc.
	//   required: false
	//   type: string
	//   enum:
	//   - asc
	//   - desc
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/datasetRowsResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getDatasetRows"

	// swagger:operation GET /api/datasets/{id}/export datasets exportDataset
	//
	// Export a dataset.
	//
	// Exports matching OCR documents as CSV or XLSX.
	//
	// ---
	// produces:
	// - text/csv
	// - application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Dataset id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Required owner user id used to scope dataset lookup.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: format
	//   in: query
	//   description: Export format. Defaults to csv.
	//   required: false
	//   type: string
	//   enum:
	//   - csv
	//   - xlsx
	// - name: created_from
	//   in: query
	//   description: Include source documents created at or after this RFC3339 timestamp.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: created_to
	//   in: query
	//   description: Include source documents created at or before this RFC3339 timestamp.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: sort
	//   in: query
	//   description: Sort direction by source document creation time. Defaults to desc.
	//   required: false
	//   type: string
	//   enum:
	//   - asc
	//   - desc
	// responses:
	//   '200':
	//     description: Export file.
	//     schema:
	//       type: file
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "exportDataset"

	// swagger:operation GET /api/datasets/{id} datasets getDataset
	//
	// Get a dataset.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Dataset id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Required owner user id used to scope dataset lookup.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/datasetResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getDataset"

	// swagger:operation PUT /api/datasets/{id} datasets updateDataset
	//
	// Update a dataset.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Dataset id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Required owner user id used to scope dataset update.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: dataset
	//   in: body
	//   description: Updated dataset name, schema id, and selected fields.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/updateDatasetRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/datasetResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "updateDataset"

	// swagger:operation DELETE /api/datasets/{id} datasets deleteDataset
	//
	// Delete a dataset.
	//
	// Deletes only the dataset row for the requested owner.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Dataset id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Required owner user id used to scope dataset deletion.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '204':
	//     description: No content.
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "deleteDataset"

	// swagger:operation POST /api/auth/sign-up/email auth signUpEmail
	//
	// Sign up with email credentials.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: credentials
	//   in: body
	//   description: Email sign-up credentials.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/signUpEmailRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/signUpEmailResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "signUpEmail"

	// swagger:operation POST /api/auth/sign-in/email auth signInEmail
	//
	// Sign in with email credentials.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: credentials
	//   in: body
	//   description: Email sign-in credentials.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/signInEmailRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/signInEmailResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Forbidden.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "signInEmail"

	// swagger:operation POST /api/auth/oauth/google/start auth startGoogleOAuth
	//
	// Start Google OAuth sign-in.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: request
	//   in: body
	//   description: Google OAuth start request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/googleOAuthStartRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/googleOAuthStartResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '503':
	//     description: Google OAuth is not configured.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "startGoogleOAuth"

	// swagger:operation POST /api/auth/oauth/google/callback auth signInGoogleOAuth
	//
	// Complete Google OAuth sign-in.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: request
	//   in: body
	//   description: Google OAuth callback request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/googleOAuthCallbackRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/signInEmailResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Invalid Google identity.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '502':
	//     description: Google token exchange failed.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '503':
	//     description: Google OAuth is not configured.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "signInGoogleOAuth"

	// swagger:operation POST /api/auth/oauth/github/start auth startGitHubOAuth
	//
	// Start GitHub OAuth sign-in.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: request
	//   in: body
	//   description: GitHub OAuth start request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/githubOAuthStartRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/githubOAuthStartResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '503':
	//     description: GitHub OAuth is not configured.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "startGitHubOAuth"

	// swagger:operation POST /api/auth/oauth/github/callback auth signInGitHubOAuth
	//
	// Complete GitHub OAuth sign-in.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: request
	//   in: body
	//   description: GitHub OAuth callback request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/githubOAuthCallbackRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/signInEmailResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Invalid GitHub identity.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '502':
	//     description: GitHub token exchange failed.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '503':
	//     description: GitHub OAuth is not configured.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "signInGitHubOAuth"

	// swagger:operation GET /api/auth/sessions auth listAuthSessions
	//
	// List active sessions for the current user.
	//
	// Returns only unexpired sessions owned by the authenticated user. Session tokens are never returned.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for the current user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/authSessionListResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listAuthSessions"

	// swagger:operation DELETE /api/auth/sessions/{id} auth revokeAuthSession
	//
	// Revoke one non-current session owned by the current user.
	//
	// The current session cannot be revoked through this endpoint.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for the current user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Session id to revoke.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/deleteAuthSessionResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Session not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "revokeAuthSession"

	// swagger:operation GET /api/auth/accounts auth listAuthAccounts
	//
	// List linked sign-in methods for the current user.
	//
	// Returns provider metadata only. Provider subject ids, OAuth tokens, and password hashes are never returned.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for the current user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/authAccountListResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listAuthAccounts"

	// swagger:operation DELETE /api/auth/accounts/{provider_id} auth unlinkAuthAccount
	//
	// Unlink one OAuth sign-in provider from the current user.
	//
	// Only google and github are supported. The endpoint rejects unlinking the last remaining sign-in method.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for the current user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: provider_id
	//   in: path
	//   description: OAuth provider id to unlink.
	//   required: true
	//   type: string
	//   enum:
	//   - google
	//   - github
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/deleteAuthAccountResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Linked account not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "unlinkAuthAccount"

	// swagger:operation POST /api/auth/accounts/google/start auth startGoogleAccountLink
	//
	// Start Google OAuth account linking for the current user.
	//
	// Uses link-specific OAuth state and never creates a login session.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for the current user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: request
	//   in: body
	//   description: Google OAuth link start request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/googleOAuthStartRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/googleOAuthStartResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '503':
	//     description: Google OAuth is not configured.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "startGoogleAccountLink"

	// swagger:operation POST /api/auth/accounts/google/callback auth linkGoogleAccount
	//
	// Complete Google OAuth account linking for the current user.
	//
	// Attaches the Google identity to the authenticated user only. Existing sessions are not created or replaced.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for the current user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: request
	//   in: body
	//   description: Google OAuth link callback request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/googleOAuthCallbackRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/authAccountListItemResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized or invalid Google identity.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '409':
	//     description: Google identity is already linked to another user.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '502':
	//     description: Google token exchange failed.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '503':
	//     description: Google OAuth is not configured.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "linkGoogleAccount"

	// swagger:operation POST /api/auth/accounts/github/start auth startGitHubAccountLink
	//
	// Start GitHub OAuth account linking for the current user.
	//
	// Uses link-specific OAuth state and never creates a login session.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for the current user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: request
	//   in: body
	//   description: GitHub OAuth link start request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/githubOAuthStartRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/githubOAuthStartResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '503':
	//     description: GitHub OAuth is not configured.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "startGitHubAccountLink"

	// swagger:operation POST /api/auth/accounts/github/callback auth linkGitHubAccount
	//
	// Complete GitHub OAuth account linking for the current user.
	//
	// Attaches the GitHub identity to the authenticated user only. Existing sessions are not created or replaced.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for the current user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: request
	//   in: body
	//   description: GitHub OAuth link callback request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/githubOAuthCallbackRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/authAccountListItemResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized or invalid GitHub identity.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '409':
	//     description: GitHub identity is already linked to another user.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '502':
	//     description: GitHub token exchange or profile fetch failed.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '503':
	//     description: GitHub OAuth is not configured.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "linkGitHubAccount"

	// swagger:operation POST /api/auth/sign-out auth signOut
	//
	// Sign out the current email session.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/signOutResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "signOut"

	// swagger:operation GET /api/auth/get-session auth getSession
	//
	// Get the current email session.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: OK. Returns null when no active session exists.
	//     schema:
	//       "$ref": "#/definitions/getSessionResponse"
	//       x-nullable: true
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getSession"

	// swagger:operation GET /api/admin/users admin listAdminUsers
	//
	// List users for the admin portal.
	//
	// Requires both the internal API token and an authenticated admin session cookie.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for an admin user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: search
	//   in: query
	//   description: Case-insensitive search across name and email.
	//   required: false
	//   type: string
	// - name: sort
	//   in: query
	//   description: Sort field. Defaults to created_at.
	//   required: false
	//   type: string
	//   enum:
	//   - created_at
	//   - last_login_at
	// - name: direction
	//   in: query
	//   description: Sort direction. Defaults to desc.
	//   required: false
	//   type: string
	//   enum:
	//   - asc
	//   - desc
	// - name: cursor
	//   in: query
	//   description: Cursor returned by a previous page.
	//   required: false
	//   type: string
	// - name: size
	//   in: query
	//   description: Page size from 1 to 100. Defaults to 20.
	//   required: false
	//   type: integer
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/adminUserListResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized. Missing or invalid session.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Forbidden. The authenticated user is not an admin.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listAdminUsers"

	// swagger:operation GET /api/admin/billing/orders admin listAdminBillingOrders
	//
	// List billing orders across all users for the admin portal.
	//
	// Requires both the internal API token and an authenticated admin session cookie.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for an admin user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Optional user id to filter orders by owner.
	//   required: false
	//   type: string
	//   format: uuid
	// - name: status
	//   in: query
	//   description: Optional order status filter.
	//   required: false
	//   type: string
	//   enum: [pending, paid, failed, refunded, canceled]
	// - name: without_invoice
	//   in: query
	//   description: When true, include only orders with no linked invoice.
	//   required: false
	//   type: boolean
	// - name: created_from
	//   in: query
	//   description: Include orders created at or after this RFC3339 timestamp.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: created_to
	//   in: query
	//   description: Include orders created at or before this RFC3339 timestamp.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: cursor
	//   in: query
	//   description: Opaque pagination cursor from a previous response.
	//   required: false
	//   type: string
	// - name: size
	//   in: query
	//   description: Page size from 1 to 100.
	//   required: false
	//   type: integer
	// - name: sort
	//   in: query
	//   description: Sort direction by created_at.
	//   required: false
	//   type: string
	//   enum: [asc, desc]
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/adminBillingOrderListResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized. Missing or invalid session.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Forbidden. The authenticated user is not an admin.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listAdminBillingOrders"

	// swagger:operation GET /api/admin/billing/invoices admin listAdminBillingInvoices
	//
	// List billing invoices across all users for the admin portal.
	//
	// Requires both the internal API token and an authenticated admin session cookie.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for an admin user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Optional user id to filter invoices by owner.
	//   required: false
	//   type: string
	//   format: uuid
	// - name: created_from
	//   in: query
	//   description: Include invoices created at or after this RFC3339 timestamp.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: created_to
	//   in: query
	//   description: Include invoices created at or before this RFC3339 timestamp.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: cursor
	//   in: query
	//   description: Opaque pagination cursor from a previous response.
	//   required: false
	//   type: string
	// - name: size
	//   in: query
	//   description: Page size from 1 to 100.
	//   required: false
	//   type: integer
	// - name: sort
	//   in: query
	//   description: Sort direction by created_at.
	//   required: false
	//   type: string
	//   enum: [asc, desc]
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/adminBillingInvoiceListResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized. Missing or invalid session.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Forbidden. The authenticated user is not an admin.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listAdminBillingInvoices"

	// swagger:operation POST /api/admin/billing/orders/{id}/invoice admin createAdminBillingOrderInvoice
	//
	// Generate an invoice for a paid billing order from the admin portal.
	//
	// Requires both the internal API token and an authenticated admin session cookie.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for an admin user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Billing order id.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '201':
	//     description: Created.
	//     schema:
	//       "$ref": "#/definitions/billingInvoiceResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized. Missing or invalid session.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Forbidden. The authenticated user is not an admin.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Billing order not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '409':
	//     description: Billing order already has an invoice.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "createAdminBillingOrderInvoice"

	// swagger:operation GET /api/admin/users/{id} admin getAdminUser
	//
	// Get one user for the admin portal.
	//
	// Requires both the internal API token and an authenticated admin session cookie.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for an admin user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Target user id.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/adminUserDetailResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized. Missing or invalid session.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Forbidden. The authenticated user is not an admin.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getAdminUser"

	// swagger:operation PATCH /api/admin/users/{id} admin patchAdminUser
	//
	// Update editable user profile fields from the admin portal.
	//
	// Role is intentionally not accepted by this endpoint. Email verification state is preserved when email changes.
	//
	// Requires both the internal API token and an authenticated admin session cookie.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for an admin user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Target user id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user
	//   in: body
	//   description: Editable user fields.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/patchAdminUserRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/adminUserResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized. Missing or invalid session.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Forbidden. The authenticated user is not an admin.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '409':
	//     description: Conflict. Email already exists.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "patchAdminUser"

	// swagger:operation POST /api/admin/users/{id}/impersonation admin startAdminUserImpersonation
	//
	// Start an admin impersonation session for a normal user.
	//
	// The existing admin session cookie is reused. The session owner remains the admin, while user-space routes resolve the selected normal user as the effective user.
	//
	// Requires both the internal API token and an authenticated admin session cookie.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for an admin user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Normal user id to impersonate.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/getSessionResponse"
	//   '400':
	//     description: Bad request. The target is not a normal user or is the admin themself.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized. Missing or invalid session.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Forbidden. The authenticated user is not an admin.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Target user not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '409':
	//     description: The session is already impersonating another user.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "startAdminUserImpersonation"

	// swagger:operation POST /api/admin/impersonation/stop admin stopAdminImpersonation
	//
	// Stop the current admin impersonation session.
	//
	// The admin session remains signed in. When no impersonation is active, this endpoint is idempotent and returns the current admin session.
	//
	// Requires both the internal API token and an authenticated admin session cookie.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for an admin user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/getSessionResponse"
	//   '401':
	//     description: Unauthorized. Missing or invalid session.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Forbidden. The authenticated user is not an admin.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "stopAdminImpersonation"

	// swagger:operation POST /api/admin/users/{id}/password admin setAdminUserPassword
	//
	// Set a new password for a user from the admin portal.
	//
	// Existing sessions for the target user are invalidated after the password is reset.
	//
	// Requires both the internal API token and an authenticated admin session cookie.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for an admin user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Target user id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: password
	//   in: body
	//   description: New password.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/adminSetPasswordRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/adminSetPasswordResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized. Missing or invalid session.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Forbidden. The authenticated user is not an admin.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "setAdminUserPassword"

	// swagger:operation POST /api/admin/users/{id}/balance-adjustment admin adjustAdminUserBalance
	//
	// Add or subtract credits from a user balance.
	//
	// Positive deltas create adjustment credits. Negative deltas consume available credits and are rejected when they would make the balance negative.
	//
	// Requires both the internal API token and an authenticated admin session cookie.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for an admin user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Target user id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: adjustment
	//   in: body
	//   description: Credit balance adjustment delta.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/adminAdjustUserBalanceRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/adminUserDetailResponse"
	//   '400':
	//     description: Bad request. Includes insufficient credits for negative adjustments.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized. Missing or invalid session.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Forbidden. The authenticated user is not an admin.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "adjustAdminUserBalance"

	// swagger:operation PUT /api/admin/users/{id}/billing-profile admin upsertAdminUserBillingProfile
	//
	// Create or update a user's billing profile from the admin portal.
	//
	// Uses the same billing validation as the user billing profile endpoint. The target user is taken from the path.
	//
	// Requires both the internal API token and an authenticated admin session cookie.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: Cookie
	//   in: header
	//   description: Auth session cookie for an admin user. The internal token alone is not sufficient.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Target user id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: profile
	//   in: body
	//   description: Billing profile fields for the target user.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/adminUpsertBillingProfileRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/billingProfileResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized. Missing or invalid session.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Forbidden. The authenticated user is not an admin.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "upsertAdminUserBillingProfile"

	// swagger:operation PATCH /api/auth/user auth patchAuthUser
	//
	// Update the current authenticated user's profile, password, and preferred language.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - in: body
	//   name: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/patchAuthUserRequest"
	// responses:
	//   "200":
	//     description: Updated user
	//     schema:
	//       "$ref": "#/definitions/authUserResponse"
	//   "400":
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   "401":
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   "409":
	//     description: Conflict.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   "500":
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "patchAuthUser"

	// swagger:operation POST /api/auth/apikeys auth createAPIKey
	//
	// Create a user API key.
	//
	// The raw api_key secret is returned only in this creation response.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: api_key
	//   in: body
	//   description: API key details.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/createAPIKeyRequest"
	// responses:
	//   '201':
	//     description: Created.
	//     schema:
	//       "$ref": "#/definitions/apiKeyResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "createAPIKey"

	// swagger:operation GET /api/auth/apikeys/{user_id} auth listAPIKeys
	//
	// List a user's API keys.
	//
	// Returned keys include metadata and key_prefix only; raw api_key secrets are not returned.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: path
	//   description: Owner user id.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/apiKeyListResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listAPIKeys"

	// swagger:operation DELETE /api/auth/apikeys auth deleteAPIKey
	//
	// Delete one user API key.
	//
	// Deletes only the key matching both user_id and api_key_id.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Owner user id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: api_key_id
	//   in: query
	//   description: API key id to delete.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/deleteAPIKeyResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "deleteAPIKey"

	// swagger:operation GET /api/auth/webhook/{user_id} auth getWebhook
	//
	// Get a user's webhook.
	//
	// Existing webhooks include metadata and has_secret only; plaintext and encrypted secrets are not returned.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: path
	//   description: Owner user id.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/webhookEnvelopeResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getWebhook"

	// swagger:operation POST /api/auth/webhook auth upsertWebhook
	//
	// Create or update a user webhook.
	//
	// The plaintext secret_key is returned only when a webhook is created. Updating an existing webhook keeps the existing encrypted secret and omits secret_key.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: webhook
	//   in: body
	//   description: Webhook settings.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/upsertWebhookRequest"
	// responses:
	//   '200':
	//     description: Updated.
	//     schema:
	//       "$ref": "#/definitions/webhookResponse"
	//   '201':
	//     description: Created.
	//     schema:
	//       "$ref": "#/definitions/webhookResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "upsertWebhook"

	// swagger:operation PATCH /api/auth/webhook/{user_id}/secret auth regenerateWebhookSecret
	//
	// Regenerate a user's webhook secret.
	//
	// The plaintext secret_key is returned only in this regeneration response.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: path
	//   description: Owner user id.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/webhookResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "regenerateWebhookSecret"

	// swagger:operation DELETE /api/auth/webhook auth deleteWebhook
	//
	// Delete a user webhook.
	//
	// Deletes only the webhook matching user_id.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Owner user id.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/deleteWebhookResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "deleteWebhook"

	// swagger:operation POST /api/auth/email-otp/send-verification-otp auth sendVerificationOTP
	//
	// Send an email verification OTP.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: request
	//   in: body
	//   description: Email verification request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/sendVerificationOTPRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/sendVerificationOTPResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "sendVerificationOTP"

	// swagger:operation POST /api/auth/email-otp/verify-email auth verifyEmailOTP
	//
	// Verify an email OTP.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: request
	//   in: body
	//   description: Email verification OTP.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/verifyEmailOTPRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/verifyEmailOTPResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "verifyEmailOTP"

	// swagger:operation POST /api/auth/password-reset/request auth requestPasswordReset
	//
	// Request a password reset link.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: request
	//   in: body
	//   description: Password reset request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/requestPasswordResetRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/requestPasswordResetResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "requestPasswordReset"

	// swagger:operation POST /api/auth/password-reset/confirm auth confirmPasswordReset
	//
	// Confirm a password reset token and set a new password.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: request
	//   in: body
	//   description: Password reset confirmation.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/confirmPasswordResetRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/confirmPasswordResetResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "confirmPasswordReset"

	// swagger:operation GET /api/billing/balance billing getCreditBalance
	//
	// Get credit balance.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Owner user id.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/creditBalanceResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getCreditBalance"

	// swagger:operation GET /v1/get-balance public-billing getPublicCreditBalance
	//
	// Get public credit balance.
	//
	// ---
	// summary: Get public credit balance
	// description: |
	//   Return the available OCR credits for the owner of the supplied public API key.
	//
	//   Authentication uses the public API key in the `Authorization` header. The header accepts either the raw API key value or `Bearer <api_key>`.
	//
	//   Ownership is derived from the API key. The request must not include `user_id`; user_id is not accepted as a query parameter. If the owner has no credit ledger entries, `available_credits` is returned as `0`.
	//
	//   Balance response example:
	//
	//   ```json
	//   {
	//     "user_id": "550e8400-e29b-41d4-a716-446655440000",
	//     "available_credits": 750
	//   }
	//   ```
	// produces:
	// - application/json
	// parameters:
	// - name: Authorization
	//   in: header
	//   description: Public API key for the credit owner. Use either the raw API key or `Bearer <api_key>`.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: Current available credits for the API key owner.
	//     schema:
	//       "$ref": "#/definitions/creditBalanceResponse"
	//     examples:
	//       application/json:
	//         user_id: 550e8400-e29b-41d4-a716-446655440000
	//         available_credits: 750
	//   '400':
	//     description: Bad request. Returned when a forbidden user_id query parameter is supplied.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized. Returned when the Authorization header is missing or the API key is invalid.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error. Returned for an unexpected server-side failure while loading available credits.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getPublicCreditBalance"

	// swagger:operation GET /api/billing/profile billing getBillingProfile
	//
	// Get billing profile.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted server-to-server billing reads.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Owner user id.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: OK. Profile is null when missing.
	//     schema:
	//       "$ref": "#/definitions/billingProfileEnvelopeResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getBillingProfile"

	// swagger:operation PUT /api/billing/profile billing upsertBillingProfile
	//
	// Upsert billing profile.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted server-to-server billing mutations.
	//   required: true
	//   type: string
	// - name: profile
	//   in: body
	//   description: Billing profile request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/upsertBillingProfileRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/billingProfileResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "upsertBillingProfile"

	// swagger:operation POST /api/billing/invoices billing createBillingInvoice
	//
	// Generate a billing invoice.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted server-to-server billing mutations.
	//   required: true
	//   type: string
	// - name: invoice
	//   in: body
	//   description: Invoice generation request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/createBillingInvoiceRequest"
	// responses:
	//   '201':
	//     description: Created. The response includes the profile snapshot, computed line totals, and allocated invoice number.
	//     schema:
	//       "$ref": "#/definitions/billingInvoiceResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "createBillingInvoice"

	// swagger:operation POST /api/billing/generate-invoice-pdf billing generateBillingInvoicePDF
	//
	// Generate and store the PDF rendering for an existing billing invoice.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted server-to-server billing mutations.
	//   required: true
	//   type: string
	// - name: invoice
	//   in: body
	//   description: Invoice PDF generation request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/generateInvoicePDFRequest"
	// responses:
	//   '200':
	//     description: OK. The response includes the updated pdf_path.
	//     schema:
	//       "$ref": "#/definitions/billingInvoiceResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Billing invoice not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '502':
	//     description: Gotenberg conversion failed.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "generateBillingInvoicePDF"

	// swagger:operation GET /api/billing/invoices/{id}/pdf billing serveUserBillingInvoicePDF
	//
	// Serve a generated billing invoice PDF for the specified owner.
	//
	// ---
	// produces:
	// - application/pdf
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted server-to-server billing reads.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Billing invoice id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: user_id
	//   in: query
	//   description: Owner user id. The invoice must belong to this user.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: download
	//   in: query
	//   description: Set to 1 to serve the PDF as an attachment.
	//   required: false
	//   type: string
	// responses:
	//   '200':
	//     description: PDF bytes.
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Billing invoice PDF not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "serveUserBillingInvoicePDF"

	// swagger:operation POST /api/billing/invoices/{id}/email-delivery/claim billing claimBillingInvoiceEmailDelivery
	//
	// Claim invoice email delivery for a generated invoice PDF.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted server-to-server billing mutations.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Billing invoice id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: delivery
	//   in: body
	//   description: Invoice email delivery claim request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/billingInvoiceEmailDeliveryRequest"
	// responses:
	//   '200':
	//     description: OK. Status is claimed, claim_active, already_sent, or not_ready.
	//     schema:
	//       "$ref": "#/definitions/billingInvoiceEmailDeliveryClaimResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Billing invoice not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "claimBillingInvoiceEmailDelivery"

	// swagger:operation POST /api/billing/invoices/{id}/email-delivery/sent billing markBillingInvoiceEmailSent
	//
	// Mark invoice email delivery as sent.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted server-to-server billing mutations.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Billing invoice id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: delivery
	//   in: body
	//   description: Invoice email delivery sent request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/billingInvoiceEmailDeliveryRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/billingInvoiceEmailDeliverySentResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Billing invoice not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "markBillingInvoiceEmailSent"

	// swagger:operation GET /static/invoice/{invoice_pdf} billing serveBillingInvoicePDF
	//
	// Serve a generated billing invoice PDF.
	//
	// ---
	// produces:
	// - application/pdf
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted server-to-server billing reads.
	//   required: true
	//   type: string
	// - name: invoice_pdf
	//   in: path
	//   description: Invoice PDF filename in the form {invoice_id}.pdf.
	//   required: true
	//   type: string
	// - name: download
	//   in: query
	//   description: Set to 1 to serve the PDF as an attachment.
	//   required: false
	//   type: string
	// responses:
	//   '200':
	//     description: PDF bytes.
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Billing invoice PDF not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "serveBillingInvoicePDF"

	// swagger:operation GET /api/billing/credit-usage-history billing listCreditUsageHistory
	//
	// List credit usage history.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted server-to-server billing reads.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: Owner user id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: type
	//   in: query
	//   description: Optional entry type filter.
	//   required: false
	//   type: string
	//   enum: [purchase, debit]
	// - name: created_from
	//   in: query
	//   description: Inclusive created_at lower bound.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: created_to
	//   in: query
	//   description: Inclusive created_at upper bound.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: cursor
	//   in: query
	//   description: Opaque cursor returned by the previous response.
	//   required: false
	//   type: string
	// - name: size
	//   in: query
	//   description: Page size from 1 to 100.
	//   required: false
	//   type: integer
	// - name: sort
	//   in: query
	//   description: Sort direction by created_at.
	//   required: false
	//   type: string
	//   enum: [asc, desc]
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/creditUsageHistoryListResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listCreditUsageHistory"

	// swagger:operation GET /api/billing/orders billing listBillingOrders
	//
	// List billing orders for a user.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted server-to-server billing reads.
	//   required: true
	//   type: string
	// - name: user_id
	//   in: query
	//   description: User id to list orders for.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: status
	//   in: query
	//   description: Optional order status filter.
	//   required: false
	//   type: string
	//   enum: [pending, paid, failed, refunded, canceled]
	// - name: created_from
	//   in: query
	//   description: Include orders created at or after this RFC3339 timestamp.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: created_to
	//   in: query
	//   description: Include orders created at or before this RFC3339 timestamp.
	//   required: false
	//   type: string
	//   format: date-time
	// - name: cursor
	//   in: query
	//   description: Opaque pagination cursor from a previous response.
	//   required: false
	//   type: string
	// - name: size
	//   in: query
	//   description: Page size from 1 to 100.
	//   required: false
	//   type: integer
	// - name: sort
	//   in: query
	//   description: Sort direction by created_at.
	//   required: false
	//   type: string
	//   enum: [asc, desc]
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/billingOrderListResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listBillingOrders"

	// swagger:operation POST /api/billing/orders billing createBillingOrder
	//
	// Create a credit purchase order.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted server-to-server billing mutations.
	//   required: true
	//   type: string
	// - name: order
	//   in: body
	//   description: Credit purchase order request.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/createBillingOrderRequest"
	// responses:
	//   '201':
	//     description: Created.
	//     schema:
	//       "$ref": "#/definitions/billingOrderResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "createBillingOrder"

	// swagger:operation POST /api/billing/orders/{id}/checkout-session billing attachBillingOrderCheckoutSession
	//
	// Attach a Stripe Checkout Session to a billing order.
	//
	// ---
	// consumes:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted server-to-server billing mutations.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Billing order id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: checkout_session
	//   in: body
	//   description: Stripe Checkout Session id.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/attachCheckoutSessionRequest"
	// responses:
	//   '204':
	//     description: No content.
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "attachBillingOrderCheckoutSession"

	// swagger:operation POST /api/billing/orders/{id}/paid billing markBillingOrderPaid
	//
	// Mark a billing order paid and grant credits.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted server-to-server billing mutations.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Billing order id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: payment
	//   in: body
	//   description: Completed Stripe payment metadata.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/markBillingOrderPaidRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/billingOrderResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '409':
	//     description: Conflict.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "markBillingOrderPaid"

	// swagger:operation POST /api/billing/orders/{id}/failed billing markBillingOrderFailed
	//
	// Mark a billing order failed without granting credits.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted server-to-server billing mutations.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Billing order id.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/billingOrderResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "markBillingOrderFailed"

	// swagger:operation GET /api/json-recipes json-recipes listPublicJSONRecipes
	//
	// List public system JSON recipes.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: cursor
	//   in: query
	//   description: Cursor returned by a previous page.
	//   required: false
	//   type: string
	// - name: size
	//   in: query
	//   description: Page size from 1 to 100. Defaults to 20.
	//   required: false
	//   type: integer
	// - name: sort
	//   in: query
	//   description: Sort direction by creation time. Defaults to desc.
	//   required: false
	//   type: string
	//   enum:
	//   - asc
	//   - desc
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/jsonRecipeListResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Unauthorized.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listPublicJSONRecipes"

	// swagger:operation GET /api/admin/json-recipe-categories json-recipes listJSONRecipeCategories
	//
	// List admin-managed JSON recipe categories.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/jsonRecipeCategoryListResponse"
	//   '401':
	//     description: Authentication required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Admin access required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listJSONRecipeCategories"

	// swagger:operation POST /api/admin/json-recipe-categories json-recipes createJSONRecipeCategory
	//
	// Create an admin-managed JSON recipe category.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: category
	//   in: body
	//   description: English and Romanian category titles.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/upsertJSONRecipeCategoryRequest"
	// responses:
	//   '201':
	//     description: Created.
	//     schema:
	//       "$ref": "#/definitions/jsonRecipeCategoryResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Authentication required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Admin access required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "createJSONRecipeCategory"

	// swagger:operation GET /api/admin/json-recipe-categories/{id} json-recipes getJSONRecipeCategory
	//
	// Get an admin-managed JSON recipe category.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Category id.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/jsonRecipeCategoryResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Authentication required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Admin access required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getJSONRecipeCategory"

	// swagger:operation PUT /api/admin/json-recipe-categories/{id} json-recipes updateJSONRecipeCategory
	//
	// Update an admin-managed JSON recipe category.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Category id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: category
	//   in: body
	//   description: English and Romanian category titles.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/upsertJSONRecipeCategoryRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/jsonRecipeCategoryResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Authentication required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Admin access required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "updateJSONRecipeCategory"

	// swagger:operation DELETE /api/admin/json-recipe-categories/{id} json-recipes deleteJSONRecipeCategory
	//
	// Delete an unused admin-managed JSON recipe category.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Category id.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '204':
	//     description: Deleted.
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Authentication required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Admin access required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '409':
	//     description: Category has assigned recipes.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "deleteJSONRecipeCategory"

	// swagger:operation GET /api/admin/json-recipes json-recipes listJSONRecipes
	//
	// List admin-managed JSON recipes.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: cursor
	//   in: query
	//   description: Cursor returned by a previous page.
	//   required: false
	//   type: string
	// - name: size
	//   in: query
	//   description: Page size from 1 to 100. Defaults to 20.
	//   required: false
	//   type: integer
	// - name: sort
	//   in: query
	//   description: Sort direction by creation time. Defaults to desc.
	//   required: false
	//   type: string
	//   enum:
	//   - asc
	//   - desc
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/jsonRecipeListResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Authentication required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Admin access required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "listJSONRecipes"

	// swagger:operation POST /api/admin/json-recipes json-recipes createJSONRecipe
	//
	// Create an admin-managed JSON recipe.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: recipe
	//   in: body
	//   description: Recipe title, description, and valid JSON Schema object.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/upsertJSONRecipeRequest"
	// responses:
	//   '201':
	//     description: Created.
	//     schema:
	//       "$ref": "#/definitions/jsonRecipeResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Authentication required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Admin access required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "createJSONRecipe"

	// swagger:operation GET /api/admin/json-recipes/{id} json-recipes getJSONRecipe
	//
	// Get an admin-managed JSON recipe.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Recipe id.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/jsonRecipeResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Authentication required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Admin access required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "getJSONRecipe"

	// swagger:operation PUT /api/admin/json-recipes/{id} json-recipes updateJSONRecipe
	//
	// Update an admin-managed JSON recipe.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Recipe id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: recipe
	//   in: body
	//   description: Recipe title, description, and valid JSON Schema object.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/upsertJSONRecipeRequest"
	// responses:
	//   '200':
	//     description: OK.
	//     schema:
	//       "$ref": "#/definitions/jsonRecipeResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Authentication required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Admin access required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "updateJSONRecipe"

	// swagger:operation DELETE /api/admin/json-recipes/{id} json-recipes deleteJSONRecipe
	//
	// Delete an admin-managed JSON recipe.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Recipe id.
	//   required: true
	//   type: string
	//   format: uuid
	// responses:
	//   '204':
	//     description: Deleted.
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '401':
	//     description: Authentication required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '403':
	//     description: Admin access required.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "deleteJSONRecipe"

	// swagger:operation POST /api/json-recipes/{id}/deploy json-recipes deployJSONRecipe
	//
	// Clone a JSON recipe into a user-owned extraction schema.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: X-Syncra-Internal-Token
	//   in: header
	//   description: Shared internal token for trusted frontend-server to Go API requests.
	//   required: true
	//   type: string
	// - name: id
	//   in: path
	//   description: Recipe id.
	//   required: true
	//   type: string
	//   format: uuid
	// - name: deploy
	//   in: body
	//   description: User id that will own the cloned extraction schema.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/deployJSONRecipeRequest"
	// responses:
	//   '201':
	//     description: Created.
	//     schema:
	//       "$ref": "#/definitions/jsonRecipeDeployResponse"
	//   '400':
	//     description: Bad request.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '404':
	//     description: Not found.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	//   '500':
	//     description: Internal server error.
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	_ = "deployJSONRecipe"
}
