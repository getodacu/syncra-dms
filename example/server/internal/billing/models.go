package billing

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
)

type BillingProvider string

const (
	BillingProviderStripe BillingProvider = "stripe"
)

type BillingEntityType string

const (
	BillingEntityIndividual BillingEntityType = "individual"
	BillingEntityCompany    BillingEntityType = "company"
)

type OrderType string

const (
	OrderTypeCreditTopup OrderType = "credit_topup"
)

type OrderStatus string

const (
	OrderStatusPending  OrderStatus = "pending"
	OrderStatusPaid     OrderStatus = "paid"
	OrderStatusFailed   OrderStatus = "failed"
	OrderStatusRefunded OrderStatus = "refunded"
	OrderStatusCanceled OrderStatus = "canceled"
)

type CreditSourceType string

const (
	CreditSourceSignupBonus   CreditSourceType = "signup_bonus"
	CreditSourceTopupPurchase CreditSourceType = "topup_purchase"
	CreditSourceRefund        CreditSourceType = "refund"
	CreditSourceAdjustment    CreditSourceType = "adjustment"
)

type CreditLedgerEntryType string

const (
	CreditLedgerEntryGrant      CreditLedgerEntryType = "grant"
	CreditLedgerEntryPurchase   CreditLedgerEntryType = "purchase"
	CreditLedgerEntryDebit      CreditLedgerEntryType = "debit"
	CreditLedgerEntryRefund     CreditLedgerEntryType = "refund"
	CreditLedgerEntryExpiry     CreditLedgerEntryType = "expiry"
	CreditLedgerEntryAdjustment CreditLedgerEntryType = "adjustment"
)

type BillingProfile struct {
	ID                 uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID             string            `gorm:"column:user_id;type:uuid;not null;uniqueIndex" json:"user_id"`
	User               auth.User         `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	EntityType         BillingEntityType `gorm:"column:entity_type;not null;size:40;check:chk_billing_profiles_entity_type,entity_type IN ('individual','company')" json:"entity_type"`
	BillingName        string            `gorm:"column:billing_name;not null;size:255" json:"billing_name"`
	BillingEmail       string            `gorm:"column:billing_email;not null;size:320" json:"billing_email"`
	CountryCode        string            `gorm:"column:country_code;not null;size:2" json:"country_code"`
	AddressLine1       string            `gorm:"column:address_line1;not null;size:255" json:"address_line1"`
	AddressLine2       *string           `gorm:"column:address_line2;size:255" json:"address_line2,omitempty"`
	City               string            `gorm:"not null;size:160" json:"city"`
	Region             *string           `gorm:"size:160" json:"region,omitempty"`
	PostalCode         string            `gorm:"column:postal_code;not null;size:40" json:"postal_code"`
	FiscalCode         *string           `gorm:"column:fiscal_code;size:80" json:"fiscal_code,omitempty"`
	RegistrationNumber *string           `gorm:"column:registration_number;size:120" json:"registration_number,omitempty"`
	CreatedAt          time.Time         `gorm:"not null" json:"created_at"`
	UpdatedAt          time.Time         `gorm:"not null" json:"updated_at"`
}

func (profile *BillingProfile) BeforeCreate(_ *gorm.DB) error {
	if profile.ID == uuid.Nil {
		profile.ID = uuid.New()
	}
	return validateBillingEntityType(profile.EntityType)
}

func (profile *BillingProfile) BeforeUpdate(tx *gorm.DB) error {
	if !updateWillWrite(tx, "EntityType", "entity_type") {
		return nil
	}
	if entityType, ok := billingEntityTypeFromUpdate(tx, profile.EntityType, "EntityType", "entity_type"); ok {
		return validateBillingEntityType(entityType)
	}
	return nil
}

func (BillingProfile) TableName() string {
	return "billing_profiles"
}

type BillingInvoice struct {
	ID                     uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	UserID                 *string         `gorm:"column:user_id;type:uuid;index" json:"user_id,omitempty"`
	User                   *auth.User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	OrderID                *uuid.UUID      `gorm:"column:order_id;type:uuid;uniqueIndex:idx_billing_invoices_order_id,where:order_id IS NOT NULL" json:"order_id,omitempty"`
	Order                  *BillingOrder   `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	BillingProfileID       *uuid.UUID      `gorm:"column:billing_profile_id;type:uuid;index" json:"billing_profile_id,omitempty"`
	BillingProfile         *BillingProfile `gorm:"foreignKey:BillingProfileID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	BillingName            string          `gorm:"column:billing_name;not null;size:255" json:"billing_name"`
	BillingEmail           string          `gorm:"column:billing_email;not null;size:320" json:"billing_email"`
	BillingFiscalCode      *string         `gorm:"column:billing_fiscal_code;size:80;index:idx_billing_invoices_billing_fiscal_code,where:billing_fiscal_code IS NOT NULL" json:"billing_fiscal_code,omitempty"`
	BillingProfileSnapshot datatypes.JSON  `gorm:"column:billing_profile_snapshot;type:jsonb;not null" json:"billing_profile_snapshot"`
	Lines                  datatypes.JSON  `gorm:"type:jsonb;not null" json:"lines"`
	NetAmount              decimal.Decimal `gorm:"column:net_amount;type:numeric(20,2);not null;check:chk_billing_invoices_net_amount,net_amount >= 0" json:"net_amount"`
	VATAmount              decimal.Decimal `gorm:"column:vat_amount;type:numeric(20,2);not null;check:chk_billing_invoices_vat_amount,vat_amount >= 0" json:"vat_amount"`
	TotalAmount            decimal.Decimal `gorm:"column:total_amount;type:numeric(20,2);not null;check:chk_billing_invoices_total_amount,total_amount >= 0 AND total_amount = net_amount + vat_amount" json:"total_amount"`
	InvoiceDate            time.Time       `gorm:"column:invoice_date;type:date;not null;index" json:"invoice_date"`
	InvoiceSerie           string          `gorm:"column:invoice_serie;not null;size:40;uniqueIndex:idx_billing_invoices_serie_number,priority:1" json:"invoice_serie"`
	InvoiceNumber          int64           `gorm:"column:invoice_number;not null;uniqueIndex:idx_billing_invoices_serie_number,priority:2;check:chk_billing_invoices_invoice_number,invoice_number > 0" json:"invoice_number"`
	PDFPath                *string         `gorm:"column:pdf_path;type:text" json:"pdf_path,omitempty"`
	EmailDeliveryClaimedAt *time.Time      `gorm:"column:email_delivery_claimed_at" json:"email_delivery_claimed_at,omitempty"`
	EmailSentAt            *time.Time      `gorm:"column:email_sent_at" json:"email_sent_at,omitempty"`
	CreatedAt              time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt              time.Time       `gorm:"not null" json:"updated_at"`
}

func (invoice *BillingInvoice) BeforeCreate(_ *gorm.DB) error {
	if invoice.ID == uuid.Nil {
		invoice.ID = uuid.New()
	}
	return nil
}

func (BillingInvoice) TableName() string {
	return "billing_invoices"
}

type BillingInvoiceCounter struct {
	InvoiceSerie string    `gorm:"column:invoice_serie;primaryKey;size:40" json:"invoice_serie"`
	LastNumber   int64     `gorm:"column:last_number;not null;default:0;check:chk_billing_invoice_counters_last_number,last_number >= 0" json:"last_number"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at"`
}

func (BillingInvoiceCounter) TableName() string {
	return "billing_invoice_counters"
}

type BillingInvoiceLine struct {
	Name           string `json:"name"`
	Quantity       int    `json:"quantity"`
	UnitPrice      string `json:"unit_price"`
	VATPercentage  string `json:"vat_percentage"`
	TotalVATAmount string `json:"total_vat_amount"`
	TotalAmount    string `json:"total_amount"`
}

type BillingProfileSnapshot struct {
	ID                 uuid.UUID         `json:"id"`
	UserID             string            `json:"user_id"`
	EntityType         BillingEntityType `json:"entity_type"`
	BillingName        string            `json:"billing_name"`
	BillingEmail       string            `json:"billing_email"`
	CountryCode        string            `json:"country_code"`
	AddressLine1       string            `json:"address_line1"`
	AddressLine2       *string           `json:"address_line2,omitempty"`
	City               string            `json:"city"`
	Region             *string           `json:"region,omitempty"`
	PostalCode         string            `json:"postal_code"`
	FiscalCode         *string           `json:"fiscal_code,omitempty"`
	RegistrationNumber *string           `json:"registration_number,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

type BillingOrder struct {
	ID                        uuid.UUID          `gorm:"type:uuid;primaryKey" json:"id"`
	UserID                    string             `gorm:"column:user_id;type:uuid;not null;index" json:"user_id"`
	User                      auth.User          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	OrderType                 OrderType          `gorm:"column:order_type;not null;size:40;index;check:chk_billing_orders_order_type,order_type IN ('credit_topup')" json:"order_type"`
	Status                    OrderStatus        `gorm:"not null;size:40;index;check:chk_billing_orders_status,status IN ('pending','paid','failed','refunded','canceled')" json:"status"`
	Provider                  BillingProvider    `gorm:"not null;size:40;check:chk_billing_orders_provider,provider IN ('stripe')" json:"provider"`
	PricingTier               CreditPurchaseTier `gorm:"column:pricing_tier;not null;size:40;check:chk_billing_orders_pricing_tier,pricing_tier IN ('tier_1','tier_2','tier_3','tier_4')" json:"pricing_tier"`
	UnitAmountCents           int                `gorm:"column:unit_amount_cents;not null;check:chk_billing_orders_unit_amount_cents,unit_amount_cents > 0" json:"unit_amount_cents"`
	Credits                   int                `gorm:"not null;default:0;check:chk_billing_orders_credits,credits > 0" json:"credits"`
	AmountCents               int                `gorm:"not null;default:0;check:chk_billing_orders_amount_cents,amount_cents >= 0" json:"amount_cents"`
	Currency                  string             `gorm:"not null;size:3" json:"currency"`
	ProviderCheckoutSessionID *string            `gorm:"column:provider_checkout_session_id;size:255;uniqueIndex:idx_billing_orders_provider_checkout_session_id,where:provider_checkout_session_id IS NOT NULL" json:"provider_checkout_session_id,omitempty"`
	ProviderPaymentIntentID   *string            `gorm:"column:provider_payment_intent_id;size:255;uniqueIndex:idx_billing_orders_provider_payment_intent_id,where:provider_payment_intent_id IS NOT NULL" json:"provider_payment_intent_id,omitempty"`
	CreatedAt                 time.Time          `json:"created_at"`
	UpdatedAt                 time.Time          `json:"updated_at"`
	PaidAt                    *time.Time         `gorm:"column:paid_at" json:"paid_at,omitempty"`
	FailedAt                  *time.Time         `gorm:"column:failed_at" json:"failed_at,omitempty"`
	RefundedAt                *time.Time         `gorm:"column:refunded_at" json:"refunded_at,omitempty"`
	CanceledAt                *time.Time         `gorm:"column:canceled_at" json:"canceled_at,omitempty"`
	Invoice                   *BillingInvoice    `gorm:"foreignKey:OrderID" json:"-"`
}

func (order *BillingOrder) BeforeCreate(_ *gorm.DB) error {
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	return order.validate()
}

func (order *BillingOrder) BeforeUpdate(tx *gorm.DB) error {
	if !updateWillWrite(tx, "OrderType", "order_type") &&
		!updateWillWrite(tx, "Status", "status") &&
		!updateWillWrite(tx, "Provider", "provider") &&
		!updateWillWrite(tx, "PricingTier", "pricing_tier") {
		return nil
	}
	return order.validateUpdate(tx)
}

func (order *BillingOrder) validate() error {
	if err := validateOrderType(order.OrderType); err != nil {
		return err
	}
	if err := validateOrderStatus(order.Status); err != nil {
		return err
	}
	if err := validateBillingProvider(order.Provider); err != nil {
		return err
	}
	return validateCreditPurchaseTier(order.PricingTier)
}

func (order *BillingOrder) validateUpdate(tx *gorm.DB) error {
	if orderType, ok := orderTypeFromUpdate(tx, order.OrderType, "OrderType", "order_type"); ok {
		if err := validateOrderType(orderType); err != nil {
			return err
		}
	}
	if status, ok := orderStatusFromUpdate(tx, order.Status, "Status", "status"); ok {
		if err := validateOrderStatus(status); err != nil {
			return err
		}
	}
	if provider, ok := billingProviderFromUpdate(tx, order.Provider, "Provider", "provider"); ok {
		if err := validateBillingProvider(provider); err != nil {
			return err
		}
	}
	if pricingTier, ok := creditPurchaseTierFromUpdate(tx, order.PricingTier, "PricingTier", "pricing_tier"); ok {
		return validateCreditPurchaseTier(pricingTier)
	}
	return nil
}

func (BillingOrder) TableName() string {
	return "billing_orders"
}

type CreditBucket struct {
	ID               uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           string           `gorm:"column:user_id;type:uuid;not null;index;index:idx_credit_buckets_user_available,priority:1,where:credits_remaining > 0 AND voided_at IS NULL;uniqueIndex:idx_credit_buckets_one_signup_bonus,where:source_type = 'signup_bonus'" json:"user_id"`
	User             auth.User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	SourceType       CreditSourceType `gorm:"column:source_type;not null;size:40;index;check:chk_credit_buckets_source_type,source_type IN ('signup_bonus','topup_purchase','refund','adjustment')" json:"source_type"`
	OrderID          *uuid.UUID       `gorm:"column:order_id;type:uuid;index" json:"order_id,omitempty"`
	Order            *BillingOrder    `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	CreditsGranted   int              `gorm:"not null;check:chk_credit_buckets_credits_granted,credits_granted > 0" json:"credits_granted"`
	CreditsRemaining int              `gorm:"not null;check:chk_credit_buckets_credits_remaining,credits_remaining >= 0 AND credits_remaining <= credits_granted" json:"credits_remaining"`
	ValidFrom        time.Time        `gorm:"column:valid_from;not null;index" json:"valid_from"`
	ExpiresAt        *time.Time       `gorm:"column:expires_at;index;index:idx_credit_buckets_user_available,priority:2,where:credits_remaining > 0 AND voided_at IS NULL" json:"expires_at,omitempty"`
	VoidedAt         *time.Time       `gorm:"column:voided_at;index" json:"voided_at,omitempty"`
	CreatedAt        time.Time        `gorm:"index:idx_credit_buckets_user_available,priority:3,where:credits_remaining > 0 AND voided_at IS NULL" json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

func (bucket *CreditBucket) BeforeCreate(_ *gorm.DB) error {
	if bucket.ID == uuid.Nil {
		bucket.ID = uuid.New()
	}
	return validateCreditSourceType(bucket.SourceType)
}

func (bucket *CreditBucket) BeforeUpdate(tx *gorm.DB) error {
	if sourceType, ok := creditSourceTypeFromUpdate(tx, bucket.SourceType, "SourceType", "source_type"); ok {
		return validateCreditSourceType(sourceType)
	}
	return nil
}

func (CreditBucket) TableName() string {
	return "credit_buckets"
}

type CreditLedgerEntry struct {
	ID             uuid.UUID             `gorm:"type:uuid;primaryKey;index:idx_credit_ledger_entries_transactions,priority:3,where:entry_type = 'purchase' OR entry_type = 'debit'" json:"id"`
	UserID         string                `gorm:"column:user_id;type:uuid;not null;index;index:idx_credit_ledger_entries_transactions,priority:1,where:entry_type = 'purchase' OR entry_type = 'debit'" json:"user_id"`
	User           auth.User             `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	BucketID       *uuid.UUID            `gorm:"column:bucket_id;type:uuid;index" json:"bucket_id,omitempty"`
	Bucket         *CreditBucket         `gorm:"foreignKey:BucketID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	EntryType      CreditLedgerEntryType `gorm:"column:entry_type;not null;size:40;index;check:chk_credit_ledger_entries_entry_type,entry_type IN ('grant','purchase','debit','refund','expiry','adjustment')" json:"entry_type"`
	CreditsDelta   int                   `gorm:"not null;check:chk_credit_ledger_entries_credits_delta,credits_delta <> 0" json:"credits_delta"`
	RelatedJobID   *uuid.UUID            `gorm:"column:related_job_id;type:uuid;index" json:"related_job_id,omitempty"`
	RelatedOrderID *uuid.UUID            `gorm:"column:related_order_id;type:uuid;index" json:"related_order_id,omitempty"`
	RelatedOrder   *BillingOrder         `gorm:"foreignKey:RelatedOrderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	IdempotencyKey string                `gorm:"column:idempotency_key;not null;size:255;uniqueIndex:idx_credit_ledger_entries_idempotency_key" json:"idempotency_key"`
	Metadata       datatypes.JSON        `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt      time.Time             `gorm:"index:idx_credit_ledger_entries_transactions,priority:2,where:entry_type = 'purchase' OR entry_type = 'debit'" json:"created_at"`
}

func (entry *CreditLedgerEntry) BeforeCreate(_ *gorm.DB) error {
	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	return validateCreditLedgerEntryType(entry.EntryType)
}

func (entry *CreditLedgerEntry) BeforeUpdate(_ *gorm.DB) error {
	return errors.New("credit ledger entries are append-only")
}

func (CreditLedgerEntry) TableName() string {
	return "credit_ledger_entries"
}

func validateBillingEntityType(entityType BillingEntityType) error {
	switch entityType {
	case BillingEntityIndividual, BillingEntityCompany:
		return nil
	default:
		return fmt.Errorf("invalid billing entity type %q", entityType)
	}
}

func validateBillingProvider(provider BillingProvider) error {
	switch provider {
	case BillingProviderStripe:
		return nil
	default:
		return fmt.Errorf("invalid billing provider %q", provider)
	}
}

func validateOrderType(orderType OrderType) error {
	switch orderType {
	case OrderTypeCreditTopup:
		return nil
	default:
		return fmt.Errorf("invalid billing order type %q", orderType)
	}
}

func validateOrderStatus(status OrderStatus) error {
	switch status {
	case OrderStatusPending, OrderStatusPaid, OrderStatusFailed, OrderStatusRefunded, OrderStatusCanceled:
		return nil
	default:
		return fmt.Errorf("invalid billing order status %q", status)
	}
}

func validateCreditPurchaseTier(tier CreditPurchaseTier) error {
	switch tier {
	case CreditPurchaseTier1, CreditPurchaseTier2, CreditPurchaseTier3, CreditPurchaseTier4:
		return nil
	default:
		return fmt.Errorf("invalid credit purchase tier %q", tier)
	}
}

func validateCreditSourceType(sourceType CreditSourceType) error {
	switch sourceType {
	case CreditSourceSignupBonus, CreditSourceTopupPurchase, CreditSourceRefund, CreditSourceAdjustment:
		return nil
	default:
		return fmt.Errorf("invalid credit source type %q", sourceType)
	}
}

func validateCreditLedgerEntryType(entryType CreditLedgerEntryType) error {
	switch entryType {
	case CreditLedgerEntryGrant, CreditLedgerEntryPurchase, CreditLedgerEntryDebit, CreditLedgerEntryRefund, CreditLedgerEntryExpiry, CreditLedgerEntryAdjustment:
		return nil
	default:
		return fmt.Errorf("invalid credit ledger entry type %q", entryType)
	}
}

func updateWillWrite(tx *gorm.DB, fieldName string, columnName string) bool {
	if tx == nil || tx.Statement == nil {
		return false
	}
	if tx.Statement.Changed(fieldName) || tx.Statement.Changed(columnName) {
		return true
	}
	selectColumns, restricted := tx.Statement.SelectAndOmitColumns(false, true)
	if selected, ok := selectColumns[columnName]; ok {
		return selected
	}
	if selected, ok := selectColumns[fieldName]; ok {
		return selected
	}
	if restricted {
		return false
	}
	return false
}

func billingEntityTypeFromUpdate(tx *gorm.DB, fallback BillingEntityType, fieldName string, columnName string) (BillingEntityType, bool) {
	value, ok := updateValue(tx, fieldName, columnName)
	if !ok {
		if updateWillWrite(tx, fieldName, columnName) {
			return fallback, true
		}
		return "", false
	}
	return BillingEntityType(fmt.Sprint(value)), true
}

func billingProviderFromUpdate(tx *gorm.DB, fallback BillingProvider, fieldName string, columnName string) (BillingProvider, bool) {
	value, ok := updateValue(tx, fieldName, columnName)
	if !ok {
		if updateWillWrite(tx, fieldName, columnName) {
			return fallback, true
		}
		return "", false
	}
	return BillingProvider(fmt.Sprint(value)), true
}

func orderTypeFromUpdate(tx *gorm.DB, fallback OrderType, fieldName string, columnName string) (OrderType, bool) {
	value, ok := updateValue(tx, fieldName, columnName)
	if !ok {
		if updateWillWrite(tx, fieldName, columnName) {
			return fallback, true
		}
		return "", false
	}
	return OrderType(fmt.Sprint(value)), true
}

func orderStatusFromUpdate(tx *gorm.DB, fallback OrderStatus, fieldName string, columnName string) (OrderStatus, bool) {
	value, ok := updateValue(tx, fieldName, columnName)
	if !ok {
		if updateWillWrite(tx, fieldName, columnName) {
			return fallback, true
		}
		return "", false
	}
	return OrderStatus(fmt.Sprint(value)), true
}

func creditPurchaseTierFromUpdate(tx *gorm.DB, fallback CreditPurchaseTier, fieldName string, columnName string) (CreditPurchaseTier, bool) {
	value, ok := updateValue(tx, fieldName, columnName)
	if !ok {
		if updateWillWrite(tx, fieldName, columnName) {
			return fallback, true
		}
		return "", false
	}
	return CreditPurchaseTier(fmt.Sprint(value)), true
}

func creditSourceTypeFromUpdate(tx *gorm.DB, fallback CreditSourceType, fieldName string, columnName string) (CreditSourceType, bool) {
	value, ok := updateValue(tx, fieldName, columnName)
	if !ok {
		if updateWillWrite(tx, fieldName, columnName) {
			return fallback, true
		}
		return "", false
	}
	return CreditSourceType(fmt.Sprint(value)), true
}

func updateValue(tx *gorm.DB, fieldName string, columnName string) (any, bool) {
	if tx == nil || tx.Statement == nil || tx.Statement.Dest == nil {
		return nil, false
	}
	return updateValueFromReflect(reflect.ValueOf(tx.Statement.Dest), fieldName, columnName)
}

func updateValueFromReflect(value reflect.Value, fieldName string, columnName string) (any, bool) {
	if !value.IsValid() {
		return nil, false
	}
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, false
		}
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Map:
		return updateValueFromMap(value, fieldName, columnName)
	case reflect.Struct:
		return updateValueFromStruct(value, fieldName)
	default:
		return nil, false
	}
}

func updateValueFromMap(updateMap reflect.Value, fieldName string, columnName string) (any, bool) {
	for _, key := range updateMap.MapKeys() {
		if key.Kind() != reflect.String {
			continue
		}
		switch key.String() {
		case fieldName, columnName:
			value := updateMap.MapIndex(key)
			if !value.IsValid() {
				return nil, true
			}
			if value.Kind() == reflect.Interface && value.IsNil() {
				return nil, true
			}
			return value.Interface(), true
		}
	}
	return nil, false
}

func updateValueFromStruct(updateStruct reflect.Value, fieldName string) (any, bool) {
	field := updateStruct.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, false
	}
	if !field.CanInterface() {
		return nil, false
	}
	return field.Interface(), true
}
