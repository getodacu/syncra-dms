package billing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"ai.ro/syncra/internal/auth"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	signupBonusDays    = 0
	syncraInvoiceSerie = "SYN"
)

const (
	billingProfileNameMaxLength               = 255
	billingProfileEmailMaxLength              = 320
	billingProfileAddressLineMaxLength        = 255
	billingProfileCityMaxLength               = 160
	billingProfileRegionMaxLength             = 160
	billingProfilePostalCodeMaxLength         = 40
	billingProfileFiscalCodeMaxLength         = 80
	billingProfileRegistrationNumberMaxLength = 120
	invoiceSerieMaxLength                     = 40
	invoiceLineNameMaxLength                  = 255
)

const billingProfileCountryCodes = "AD AE AF AG AI AL AM AO AQ AR AS AT AU AW AX AZ BA BB BD BE BF BG BH BI BJ BL BM BN BO BQ BR BS BT BV BW BY BZ CA CC CD CF CG CH CI CK CL CM CN CO CR CU CV CW CX CY CZ DE DJ DK DM DO DZ EC EE EG EH ER ES ET FI FJ FK FM FO FR GA GB GD GE GF GG GH GI GL GM GN GP GQ GR GS GT GU GW GY HK HM HN HR HT HU ID IE IL IM IN IO IQ IR IS IT JE JM JO JP KE KG KH KI KM KN KP KR KW KY KZ LA LB LC LI LK LR LS LT LU LV LY MA MC MD ME MF MG MH MK ML MM MN MO MP MQ MR MS MT MU MV MW MX MY MZ NA NC NE NF NG NI NL NO NP NR NU NZ OM PA PE PF PG PH PK PL PM PN PR PS PT PW PY QA RE RO RS RU RW SA SB SC SD SE SG SH SI SJ SK SL SM SN SO SR SS ST SV SX SY SZ TC TD TF TG TH TJ TK TL TM TN TO TR TT TV TW TZ UA UG UM US UY UZ VA VC VE VG VI VN VU WF WS YE YT ZA ZM ZW"

var (
	ErrProviderMetadataConflict = errors.New("billing: provider metadata belongs to another order")
	ErrBillingInvoiceExists     = errors.New("billing: invoice already exists for order")
	ErrBillingOrderNotPaid      = errors.New("billing: order is not paid")
	ErrBillingProfileRequired   = errors.New("billing: billing profile is required")
	ErrInsufficientCredits      = errors.New("billing: insufficient credits")
)

type CreateCreditOrderInput struct {
	UserID  string
	Credits int
}

type UpsertBillingProfileInput struct {
	UserID             string
	EntityType         BillingEntityType
	BillingName        string
	BillingEmail       string
	CountryCode        string
	AddressLine1       string
	AddressLine2       *string
	City               string
	Region             *string
	PostalCode         string
	FiscalCode         *string
	RegistrationNumber *string
}

type CreateBillingInvoiceInput struct {
	UserID       string
	OrderID      *uuid.UUID
	InvoiceSerie string
	InvoiceDate  time.Time
	Lines        []CreateBillingInvoiceLineInput
}

type CreateBillingInvoiceLineInput struct {
	Name          string
	Quantity      int
	UnitPrice     string
	VATPercentage string
}

type MarkCreditOrderPaidInput struct {
	OrderID                   uuid.UUID
	ProviderCheckoutSessionID *string
	ProviderPaymentIntentID   *string
	PaidAt                    time.Time
}

type DebitCreditsInput struct {
	UserID         string
	RelatedJobID   uuid.UUID
	Credits        int
	IdempotencyKey string
	Now            time.Time
}

type AdjustCreditsInput struct {
	UserID         string
	Delta          int
	Now            time.Time
	IdempotencyKey string
}

type GrantSignupBonusInput struct {
	UserID  string
	Credits int
	Now     time.Time
}

type ListCreditLedgerTransactionsInput struct {
	UserID      string
	EntryType   *CreditLedgerEntryType
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	Cursor      *CreditLedgerTransactionCursor
	Size        int
	Sort        string
}

type ListBillingOrdersInput struct {
	UserID      string
	Status      *OrderStatus
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	Cursor      *BillingOrderCursor
	Size        int
	Sort        string
}

type ListAdminBillingOrdersInput struct {
	UserID         *string
	Status         *OrderStatus
	CreatedFrom    *time.Time
	CreatedTo      *time.Time
	Cursor         *BillingOrderCursor
	Size           int
	Sort           string
	WithoutInvoice bool
}

type ListAdminBillingInvoicesInput struct {
	Search      string
	UserID      *string
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	Cursor      *BillingInvoiceCursor
	Size        int
	Sort        string
}

type CreditLedgerTransactionCursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
	Sort      string
}

type BillingOrderCursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
	Sort      string
}

type BillingInvoiceCursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
	Sort      string
}

type CreditLedgerTransactionPage struct {
	Entries    []CreditLedgerEntry
	NextCursor *CreditLedgerTransactionCursor
}

type BillingOrderPage struct {
	Orders     []BillingOrder
	NextCursor *BillingOrderCursor
}

type BillingInvoicePage struct {
	Invoices   []BillingInvoice
	NextCursor *BillingInvoiceCursor
}

type CreditBalance struct {
	Available int
}

type normalizedBillingInvoiceInput struct {
	UserID       string
	OrderID      *uuid.UUID
	InvoiceSerie string
	InvoiceDate  time.Time
	Lines        []BillingInvoiceLine
	NetAmount    decimal.Decimal
	VATAmount    decimal.Decimal
	TotalAmount  decimal.Decimal
}

type billingInvoiceBuyer struct {
	UserID                 string
	BillingProfileID       *uuid.UUID
	BillingName            string
	BillingEmail           string
	BillingFiscalCode      *string
	BillingProfileSnapshot datatypes.JSON
}

type billingUserSnapshot struct {
	Source       string `json:"source"`
	UserID       string `json:"user_id"`
	BillingName  string `json:"billing_name"`
	BillingEmail string `json:"billing_email"`
}

func GetBillingProfile(ctx context.Context, db *gorm.DB, userID string) (BillingProfile, error) {
	if db == nil {
		return BillingProfile{}, errors.New("billing: nil db")
	}
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return BillingProfile{}, errors.New("billing: user id is required")
	}

	var profile BillingProfile
	if err := db.WithContext(ctx).First(&profile, "user_id = ?", userID).Error; err != nil {
		return BillingProfile{}, err
	}
	return profile, nil
}

func UpsertBillingProfile(ctx context.Context, db *gorm.DB, input UpsertBillingProfileInput) (BillingProfile, error) {
	if db == nil {
		return BillingProfile{}, errors.New("billing: nil db")
	}
	profile, err := normalizedBillingProfile(input)
	if err != nil {
		return BillingProfile{}, err
	}

	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "user_id"}},
				DoUpdates: clause.Assignments(billingProfileUpsertAssignments(profile)),
			},
			clause.Returning{},
		).Create(&profile).Error
	})
	if err != nil {
		return BillingProfile{}, err
	}
	return profile, nil
}

func CreateBillingInvoice(ctx context.Context, db *gorm.DB, input CreateBillingInvoiceInput) (BillingInvoice, error) {
	if db == nil {
		return BillingInvoice{}, errors.New("billing: nil db")
	}
	normalized, err := normalizedBillingInvoice(input)
	if err != nil {
		return BillingInvoice{}, err
	}

	var invoice BillingInvoice
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		created, err := createBillingInvoiceTx(tx, normalized)
		if err != nil {
			return err
		}
		invoice = created
		return nil
	})
	if err != nil {
		return BillingInvoice{}, err
	}
	return invoice, nil
}

func CreateBillingInvoiceForPaidOrder(ctx context.Context, db *gorm.DB, orderID uuid.UUID, invoiceDate time.Time) (BillingInvoice, error) {
	if db == nil {
		return BillingInvoice{}, errors.New("billing: nil db")
	}
	if orderID == uuid.Nil {
		return BillingInvoice{}, errors.New("billing: order id is required")
	}

	var invoice BillingInvoice
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order BillingOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, "id = ?", orderID).Error; err != nil {
			return err
		}
		if order.Status != OrderStatusPaid {
			return ErrBillingOrderNotPaid
		}

		var existing BillingInvoice
		if err := tx.Where("order_id = ?", order.ID).First(&existing).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		} else {
			return ErrBillingInvoiceExists
		}

		orderID := order.ID
		normalized, err := normalizedBillingInvoice(CreateBillingInvoiceInput{
			UserID:       order.UserID,
			OrderID:      &orderID,
			InvoiceSerie: syncraInvoiceSerie,
			InvoiceDate:  invoiceDate,
			Lines: []CreateBillingInvoiceLineInput{
				{
					Name:          fmt.Sprintf("SYNCRA SaaS %d credits", order.Credits),
					Quantity:      1,
					UnitPrice:     orderAmountUnitPrice(order.AmountCents),
					VATPercentage: "0.00",
				},
			},
		})
		if err != nil {
			return err
		}
		created, err := createBillingInvoiceTx(tx, normalized)
		if err != nil {
			return err
		}
		invoice = created
		return nil
	})
	if err != nil {
		return BillingInvoice{}, err
	}
	return invoice, nil
}

func createBillingInvoiceTx(tx *gorm.DB, normalized normalizedBillingInvoiceInput) (BillingInvoice, error) {
	buyer, err := billingInvoiceBuyerForUser(tx, normalized.UserID)
	if err != nil {
		return BillingInvoice{}, err
	}

	if err := ensureBillingInvoiceCounter(tx, normalized.InvoiceSerie); err != nil {
		return BillingInvoice{}, err
	}
	nextNumber, err := incrementBillingInvoiceCounter(tx, normalized.InvoiceSerie)
	if err != nil {
		return BillingInvoice{}, err
	}
	if nextNumber <= 0 {
		return BillingInvoice{}, errors.New("billing: failed to allocate invoice number")
	}

	lines, err := billingInvoiceLinesJSON(normalized.Lines)
	if err != nil {
		return BillingInvoice{}, err
	}

	userID := buyer.UserID
	invoice := BillingInvoice{
		UserID:                 &userID,
		OrderID:                normalized.OrderID,
		BillingProfileID:       buyer.BillingProfileID,
		BillingName:            buyer.BillingName,
		BillingEmail:           buyer.BillingEmail,
		BillingFiscalCode:      buyer.BillingFiscalCode,
		BillingProfileSnapshot: buyer.BillingProfileSnapshot,
		Lines:                  lines,
		NetAmount:              normalized.NetAmount,
		VATAmount:              normalized.VATAmount,
		TotalAmount:            normalized.TotalAmount,
		InvoiceDate:            normalized.InvoiceDate,
		InvoiceSerie:           normalized.InvoiceSerie,
		InvoiceNumber:          nextNumber,
	}
	if err := tx.Create(&invoice).Error; err != nil {
		return BillingInvoice{}, err
	}
	return invoice, nil
}

func billingInvoiceBuyerForUser(tx *gorm.DB, userID string) (billingInvoiceBuyer, error) {
	var profile BillingProfile
	if err := tx.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return billingInvoiceBuyer{}, err
		}
		var user auth.User
		if err := tx.First(&user, "id = ?", userID).Error; err != nil {
			return billingInvoiceBuyer{}, err
		}
		snapshot, err := billingUserSnapshotJSON(user)
		if err != nil {
			return billingInvoiceBuyer{}, err
		}
		return billingInvoiceBuyer{
			UserID:                 user.ID,
			BillingName:            user.Name,
			BillingEmail:           user.Email,
			BillingProfileSnapshot: snapshot,
		}, nil
	}

	snapshot, err := billingProfileSnapshotJSON(profile)
	if err != nil {
		return billingInvoiceBuyer{}, err
	}
	profileID := profile.ID
	return billingInvoiceBuyer{
		UserID:                 profile.UserID,
		BillingProfileID:       &profileID,
		BillingName:            profile.BillingName,
		BillingEmail:           profile.BillingEmail,
		BillingFiscalCode:      profile.FiscalCode,
		BillingProfileSnapshot: snapshot,
	}, nil
}

func GrantSignupBonus(ctx context.Context, db *gorm.DB, input GrantSignupBonusInput) (CreditBucket, error) {
	if db == nil {
		return CreditBucket{}, errors.New("billing: nil db")
	}
	if input.UserID == "" {
		return CreditBucket{}, errors.New("billing: user id is required")
	}
	if input.Credits <= 0 {
		return CreditBucket{}, errors.New("billing: credits must be positive")
	}
	if input.Now.IsZero() {
		return CreditBucket{}, errors.New("billing: grant time is required")
	}

	var expiresAt *time.Time
	if signupBonusDays > 0 {
		expiresAt = timePtr(input.Now.AddDate(0, 0, signupBonusDays))
	}

	grantInput := grantCreditsInput{
		UserID:         input.UserID,
		SourceType:     CreditSourceSignupBonus,
		Credits:        input.Credits,
		ValidFrom:      input.Now,
		ExpiresAt:      expiresAt,
		EntryType:      CreditLedgerEntryGrant,
		IdempotencyKey: "signup_bonus:" + input.UserID,
	}
	return grantCredits(ctx, db, grantInput)
}

func CreateCreditOrder(ctx context.Context, db *gorm.DB, input CreateCreditOrderInput) (BillingOrder, error) {
	if db == nil {
		return BillingOrder{}, errors.New("billing: nil db")
	}
	if input.UserID == "" {
		return BillingOrder{}, errors.New("billing: user id is required")
	}
	quote, err := QuoteCreditPurchase(input.Credits)
	if err != nil {
		return BillingOrder{}, err
	}

	order := BillingOrder{
		UserID:          input.UserID,
		OrderType:       OrderTypeCreditTopup,
		Status:          OrderStatusPending,
		Provider:        BillingProviderStripe,
		PricingTier:     quote.Tier,
		UnitAmountCents: quote.UnitAmountCents,
		Credits:         quote.Credits,
		AmountCents:     quote.AmountCents,
		Currency:        quote.Currency,
	}
	if err := db.WithContext(ctx).Create(&order).Error; err != nil {
		return BillingOrder{}, err
	}
	return order, nil
}

func AttachCreditOrderCheckoutSession(ctx context.Context, db *gorm.DB, orderID uuid.UUID, checkoutSessionID string) error {
	if db == nil {
		return errors.New("billing: nil db")
	}
	if orderID == uuid.Nil {
		return errors.New("billing: order id is required")
	}
	if checkoutSessionID == "" {
		return errors.New("billing: checkout session id is required")
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockIdempotencyKey(tx, "provider_checkout_session:"+checkoutSessionID); err != nil {
			return err
		}
		var order BillingOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", orderID).
			First(&order).Error; err != nil {
			return err
		}
		if order.OrderType != OrderTypeCreditTopup {
			return fmt.Errorf("billing: unsupported order type %q", order.OrderType)
		}
		if order.ProviderCheckoutSessionID != nil {
			if *order.ProviderCheckoutSessionID == checkoutSessionID {
				return nil
			}
			return fmt.Errorf("billing: order already has checkout session %q", *order.ProviderCheckoutSessionID)
		}
		if err := ensureProviderCheckoutSessionAvailable(tx, order.ID, checkoutSessionID); err != nil {
			return err
		}
		return tx.Model(&BillingOrder{}).
			Where("id = ?", order.ID).
			Update("provider_checkout_session_id", checkoutSessionID).Error
	})
}

func MarkCreditOrderPaidAndGrantCredits(ctx context.Context, db *gorm.DB, input MarkCreditOrderPaidInput) (CreditBucket, error) {
	if db == nil {
		return CreditBucket{}, errors.New("billing: nil db")
	}
	if input.OrderID == uuid.Nil {
		return CreditBucket{}, errors.New("billing: order id is required")
	}
	if input.PaidAt.IsZero() {
		return CreditBucket{}, errors.New("billing: paid time is required")
	}

	var bucket CreditBucket
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if providerPaymentIntentID := nonEmptyString(input.ProviderPaymentIntentID); providerPaymentIntentID != "" {
			if err := lockIdempotencyKey(tx, "provider_payment_intent:"+providerPaymentIntentID); err != nil {
				return err
			}
			if err := ensureProviderPaymentIntentAvailable(tx, input.OrderID, providerPaymentIntentID); err != nil {
				return err
			}
			existing, err := bucketForProviderPaymentIntent(tx, providerPaymentIntentID)
			if err != nil {
				return err
			}
			if existing != nil {
				bucket = *existing
				return nil
			}
		}
		if checkoutSessionID := nonEmptyString(input.ProviderCheckoutSessionID); checkoutSessionID != "" {
			if err := lockIdempotencyKey(tx, "provider_checkout_session:"+checkoutSessionID); err != nil {
				return err
			}
			if err := ensureProviderCheckoutSessionAvailable(tx, input.OrderID, checkoutSessionID); err != nil {
				return err
			}
			existing, err := bucketForProviderCheckoutSession(tx, checkoutSessionID)
			if err != nil {
				return err
			}
			if existing != nil {
				bucket = *existing
				return nil
			}
		}

		ledgerKey := "topup_paid:" + input.OrderID.String()
		if err := lockIdempotencyKey(tx, ledgerKey); err != nil {
			return err
		}
		existing, err := bucketForLedgerIdempotencyKey(tx, ledgerKey)
		if err != nil {
			return err
		}
		if existing != nil {
			bucket = *existing
			return nil
		}

		var order BillingOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", input.OrderID).
			First(&order).Error; err != nil {
			return err
		}
		if order.OrderType != OrderTypeCreditTopup {
			return fmt.Errorf("billing: unsupported order type %q", order.OrderType)
		}
		if order.Status != OrderStatusPending && order.Status != OrderStatusPaid {
			return fmt.Errorf("billing: order status %q cannot be paid", order.Status)
		}
		if checkoutSessionID := nonEmptyString(input.ProviderCheckoutSessionID); checkoutSessionID != "" {
			if order.ProviderCheckoutSessionID != nil && *order.ProviderCheckoutSessionID != checkoutSessionID {
				return fmt.Errorf("billing: checkout session %q does not match order", checkoutSessionID)
			}
		}

		updates := map[string]any{
			"status":  OrderStatusPaid,
			"paid_at": input.PaidAt,
		}
		if providerPaymentIntentID := nonEmptyString(input.ProviderPaymentIntentID); providerPaymentIntentID != "" {
			updates["provider_payment_intent_id"] = providerPaymentIntentID
		}
		if checkoutSessionID := nonEmptyString(input.ProviderCheckoutSessionID); checkoutSessionID != "" {
			updates["provider_checkout_session_id"] = checkoutSessionID
		}
		if err := tx.Model(&BillingOrder{}).
			Where("id = ?", order.ID).
			Updates(updates).Error; err != nil {
			return err
		}

		granted, err := grantCredits(ctx, tx, grantCreditsInput{
			UserID:         order.UserID,
			SourceType:     CreditSourceTopupPurchase,
			OrderID:        &order.ID,
			Credits:        order.Credits,
			ValidFrom:      input.PaidAt,
			EntryType:      CreditLedgerEntryPurchase,
			IdempotencyKey: ledgerKey,
		})
		if err != nil {
			return err
		}
		bucket = granted
		return nil
	})
	return bucket, err
}

func MarkCreditOrderFailed(ctx context.Context, db *gorm.DB, orderID uuid.UUID, failedAt time.Time) error {
	if db == nil {
		return errors.New("billing: nil db")
	}
	if orderID == uuid.Nil {
		return errors.New("billing: order id is required")
	}
	if failedAt.IsZero() {
		return errors.New("billing: failed time is required")
	}
	return markCreditOrderTerminal(ctx, db, orderID, OrderStatusFailed, "failed_at", failedAt)
}

func MarkCreditOrderCanceled(ctx context.Context, db *gorm.DB, orderID uuid.UUID, canceledAt time.Time) error {
	if db == nil {
		return errors.New("billing: nil db")
	}
	if orderID == uuid.Nil {
		return errors.New("billing: order id is required")
	}
	if canceledAt.IsZero() {
		return errors.New("billing: canceled time is required")
	}
	return markCreditOrderTerminal(ctx, db, orderID, OrderStatusCanceled, "canceled_at", canceledAt)
}

func AvailableCredits(ctx context.Context, db *gorm.DB, userID string, now time.Time) (CreditBalance, error) {
	if db == nil {
		return CreditBalance{}, errors.New("billing: nil db")
	}
	if userID == "" {
		return CreditBalance{}, errors.New("billing: user id is required")
	}
	if now.IsZero() {
		return CreditBalance{}, errors.New("billing: time is required")
	}
	var total int
	if err := db.WithContext(ctx).
		Model(&CreditBucket{}).
		Select("COALESCE(SUM(credits_remaining), 0)").
		Where("user_id = ? AND credits_remaining > 0 AND valid_from <= ? AND voided_at IS NULL AND (expires_at IS NULL OR expires_at > ?)", userID, now, now).
		Scan(&total).Error; err != nil {
		return CreditBalance{}, err
	}
	return CreditBalance{Available: total}, nil
}

func AdjustCredits(ctx context.Context, db *gorm.DB, input AdjustCreditsInput) (CreditBalance, error) {
	if db == nil {
		return CreditBalance{}, errors.New("billing: nil db")
	}
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return CreditBalance{}, errors.New("billing: user id is required")
	}
	if input.Delta == 0 {
		return CreditBalance{}, errors.New("billing: credits delta must be non-zero")
	}
	if input.Now.IsZero() {
		return CreditBalance{}, errors.New("billing: time is required")
	}
	idempotencyKey := strings.TrimSpace(input.IdempotencyKey)
	if idempotencyKey == "" {
		return CreditBalance{}, errors.New("billing: idempotency key is required")
	}

	var balance CreditBalance
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if input.Delta > 0 {
			if _, err := grantCredits(ctx, tx, grantCreditsInput{
				UserID:         userID,
				SourceType:     CreditSourceAdjustment,
				Credits:        input.Delta,
				ValidFrom:      input.Now,
				EntryType:      CreditLedgerEntryAdjustment,
				IdempotencyKey: idempotencyKey,
			}); err != nil {
				return err
			}
			next, err := AvailableCredits(ctx, tx, userID, input.Now)
			if err != nil {
				return err
			}
			balance = next
			return nil
		}

		if err := lockIdempotencyKey(tx, idempotencyKey); err != nil {
			return err
		}
		existing, err := adjustmentLedgerEntriesExist(tx, idempotencyKey)
		if err != nil {
			return err
		}
		if existing {
			next, err := AvailableCredits(ctx, tx, userID, input.Now)
			if err != nil {
				return err
			}
			balance = next
			return nil
		}

		amount := -input.Delta
		buckets, err := lockEligibleCreditBuckets(tx, userID, input.Now)
		if err != nil {
			return err
		}
		total := 0
		for _, bucket := range buckets {
			total += bucket.CreditsRemaining
		}
		if total < amount {
			return ErrInsufficientCredits
		}

		remaining := amount
		for i, bucket := range buckets {
			if remaining == 0 {
				break
			}
			consume := bucket.CreditsRemaining
			if consume > remaining {
				consume = remaining
			}
			if consume <= 0 {
				continue
			}
			if err := tx.Model(&CreditBucket{}).
				Where("id = ?", bucket.ID).
				Update("credits_remaining", bucket.CreditsRemaining-consume).Error; err != nil {
				return err
			}
			bucketID := bucket.ID
			entryIDempotencyKey := idempotencyKey
			if i > 0 {
				entryIDempotencyKey = fmt.Sprintf("%s:%s", idempotencyKey, bucket.ID.String())
			}
			entry := CreditLedgerEntry{
				UserID:         userID,
				BucketID:       &bucketID,
				EntryType:      CreditLedgerEntryAdjustment,
				CreditsDelta:   -consume,
				IdempotencyKey: entryIDempotencyKey,
				Metadata:       adjustmentMetadata(idempotencyKey),
			}
			if err := tx.Create(&entry).Error; err != nil {
				return err
			}
			remaining -= consume
		}

		next, err := AvailableCredits(ctx, tx, userID, input.Now)
		if err != nil {
			return err
		}
		balance = next
		return nil
	})
	return balance, err
}

func ListCreditLedgerTransactions(ctx context.Context, db *gorm.DB, input ListCreditLedgerTransactionsInput) (CreditLedgerTransactionPage, error) {
	if db == nil {
		return CreditLedgerTransactionPage{}, errors.New("billing: nil db")
	}
	if input.UserID == "" {
		return CreditLedgerTransactionPage{}, errors.New("billing: user id is required")
	}
	if input.Size < 1 || input.Size > 100 {
		return CreditLedgerTransactionPage{}, errors.New("billing: size must be between 1 and 100")
	}
	if input.Sort != "asc" && input.Sort != "desc" {
		return CreditLedgerTransactionPage{}, errors.New("billing: sort must be asc or desc")
	}
	if input.EntryType != nil && *input.EntryType != CreditLedgerEntryPurchase && *input.EntryType != CreditLedgerEntryDebit {
		return CreditLedgerTransactionPage{}, errors.New("billing: transaction type must be purchase or debit")
	}
	if input.CreatedFrom != nil && input.CreatedTo != nil && input.CreatedFrom.After(*input.CreatedTo) {
		return CreditLedgerTransactionPage{}, errors.New("billing: created_from must be before or equal to created_to")
	}
	if input.Cursor != nil {
		if input.Cursor.ID == uuid.Nil || input.Cursor.CreatedAt.IsZero() || input.Cursor.Sort != input.Sort {
			return CreditLedgerTransactionPage{}, errors.New("billing: invalid cursor")
		}
	}

	query := db.WithContext(ctx).
		Model(&CreditLedgerEntry{}).
		Where("user_id = ?", input.UserID).
		Where("entry_type IN ?", []CreditLedgerEntryType{CreditLedgerEntryPurchase, CreditLedgerEntryDebit})
	if input.EntryType != nil {
		query = query.Where("entry_type = ?", *input.EntryType)
	}
	if input.CreatedFrom != nil {
		query = query.Where("created_at >= ?", *input.CreatedFrom)
	}
	if input.CreatedTo != nil {
		query = query.Where("created_at <= ?", *input.CreatedTo)
	}
	if input.Cursor != nil {
		operator := "<"
		if input.Sort == "asc" {
			operator = ">"
		}
		query = query.Where("(created_at, id) "+operator+" (?, ?)", input.Cursor.CreatedAt, input.Cursor.ID)
	}

	order := "created_at desc, id desc"
	if input.Sort == "asc" {
		order = "created_at asc, id asc"
	}
	var entries []CreditLedgerEntry
	if err := query.Order(order).Limit(input.Size + 1).Find(&entries).Error; err != nil {
		return CreditLedgerTransactionPage{}, err
	}

	var nextCursor *CreditLedgerTransactionCursor
	if len(entries) > input.Size {
		entries = entries[:input.Size]
		if len(entries) > 0 {
			last := entries[len(entries)-1]
			nextCursor = &CreditLedgerTransactionCursor{
				CreatedAt: last.CreatedAt.UTC(),
				ID:        last.ID,
				Sort:      input.Sort,
			}
		}
	}
	return CreditLedgerTransactionPage{Entries: entries, NextCursor: nextCursor}, nil
}

func ListBillingOrders(ctx context.Context, db *gorm.DB, input ListBillingOrdersInput) (BillingOrderPage, error) {
	if db == nil {
		return BillingOrderPage{}, errors.New("billing: nil db")
	}
	if input.UserID == "" {
		return BillingOrderPage{}, errors.New("billing: user id is required")
	}
	return listBillingOrders(ctx, db, listBillingOrdersInput{
		UserID:         &input.UserID,
		Status:         input.Status,
		CreatedFrom:    input.CreatedFrom,
		CreatedTo:      input.CreatedTo,
		Cursor:         input.Cursor,
		Size:           input.Size,
		Sort:           input.Sort,
		IncludeInvoice: true,
	})
}

func ListAdminBillingOrders(ctx context.Context, db *gorm.DB, input ListAdminBillingOrdersInput) (BillingOrderPage, error) {
	if db == nil {
		return BillingOrderPage{}, errors.New("billing: nil db")
	}
	if input.UserID != nil {
		userID := strings.TrimSpace(*input.UserID)
		if userID == "" {
			return BillingOrderPage{}, errors.New("billing: user id is required")
		}
		input.UserID = &userID
	}
	return listBillingOrders(ctx, db, listBillingOrdersInput{
		UserID:         input.UserID,
		Status:         input.Status,
		CreatedFrom:    input.CreatedFrom,
		CreatedTo:      input.CreatedTo,
		Cursor:         input.Cursor,
		Size:           input.Size,
		Sort:           input.Sort,
		IncludeUser:    true,
		IncludeInvoice: true,
		WithoutInvoice: input.WithoutInvoice,
	})
}

func ListAdminBillingInvoices(ctx context.Context, db *gorm.DB, input ListAdminBillingInvoicesInput) (BillingInvoicePage, error) {
	if db == nil {
		return BillingInvoicePage{}, errors.New("billing: nil db")
	}
	input.Search = strings.TrimSpace(input.Search)
	if input.UserID != nil {
		userID := strings.TrimSpace(*input.UserID)
		if userID == "" {
			return BillingInvoicePage{}, errors.New("billing: user id is required")
		}
		input.UserID = &userID
	}
	if input.Size < 1 || input.Size > 100 {
		return BillingInvoicePage{}, errors.New("billing: size must be between 1 and 100")
	}
	if input.Sort != "asc" && input.Sort != "desc" {
		return BillingInvoicePage{}, errors.New("billing: sort must be asc or desc")
	}
	if input.CreatedFrom != nil && input.CreatedTo != nil && input.CreatedFrom.After(*input.CreatedTo) {
		return BillingInvoicePage{}, errors.New("billing: created_from must be before or equal to created_to")
	}
	if input.Cursor != nil {
		if input.Cursor.ID == uuid.Nil || input.Cursor.CreatedAt.IsZero() || input.Cursor.Sort != input.Sort {
			return BillingInvoicePage{}, errors.New("billing: invalid cursor")
		}
	}

	query := db.WithContext(ctx).Model(&BillingInvoice{})
	if input.Search != "" {
		term := "%" + escapeSQLLike(strings.ToLower(input.Search)) + "%"
		query = query.Where(
			`(
				lower(billing_name) LIKE ? ESCAPE '\' OR
				lower(billing_email) LIKE ? ESCAPE '\' OR
				lower(invoice_serie || '-' || lpad(invoice_number::text, 5, '0')) LIKE ? ESCAPE '\' OR
				lower(invoice_serie || '-' || invoice_number::text) LIKE ? ESCAPE '\' OR
				invoice_number::text LIKE ? ESCAPE '\'
			)`,
			term,
			term,
			term,
			term,
			term,
		)
	}
	if input.UserID != nil {
		query = query.Where("user_id = ?", *input.UserID)
	}
	if input.CreatedFrom != nil {
		query = query.Where("created_at >= ?", *input.CreatedFrom)
	}
	if input.CreatedTo != nil {
		query = query.Where("created_at <= ?", *input.CreatedTo)
	}
	if input.Cursor != nil {
		operator := "<"
		if input.Sort == "asc" {
			operator = ">"
		}
		query = query.Where("(created_at, id) "+operator+" (?, ?)", input.Cursor.CreatedAt, input.Cursor.ID)
	}

	order := "created_at desc, id desc"
	if input.Sort == "asc" {
		order = "created_at asc, id asc"
	}
	var invoices []BillingInvoice
	if err := query.Order(order).Limit(input.Size + 1).Find(&invoices).Error; err != nil {
		return BillingInvoicePage{}, err
	}

	var nextCursor *BillingInvoiceCursor
	if len(invoices) > input.Size {
		invoices = invoices[:input.Size]
		if len(invoices) > 0 {
			last := invoices[len(invoices)-1]
			nextCursor = &BillingInvoiceCursor{
				CreatedAt: last.CreatedAt.UTC(),
				ID:        last.ID,
				Sort:      input.Sort,
			}
		}
	}
	return BillingInvoicePage{Invoices: invoices, NextCursor: nextCursor}, nil
}

func escapeSQLLike(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(value)
}

type listBillingOrdersInput struct {
	UserID         *string
	Status         *OrderStatus
	CreatedFrom    *time.Time
	CreatedTo      *time.Time
	Cursor         *BillingOrderCursor
	Size           int
	Sort           string
	IncludeUser    bool
	IncludeInvoice bool
	WithoutInvoice bool
}

func listBillingOrders(ctx context.Context, db *gorm.DB, input listBillingOrdersInput) (BillingOrderPage, error) {
	if input.Size < 1 || input.Size > 100 {
		return BillingOrderPage{}, errors.New("billing: size must be between 1 and 100")
	}
	if input.Sort != "asc" && input.Sort != "desc" {
		return BillingOrderPage{}, errors.New("billing: sort must be asc or desc")
	}
	if input.Status != nil {
		if err := validateOrderStatus(*input.Status); err != nil {
			return BillingOrderPage{}, err
		}
	}
	if input.CreatedFrom != nil && input.CreatedTo != nil && input.CreatedFrom.After(*input.CreatedTo) {
		return BillingOrderPage{}, errors.New("billing: created_from must be before or equal to created_to")
	}
	if input.Cursor != nil {
		if input.Cursor.ID == uuid.Nil || input.Cursor.CreatedAt.IsZero() || input.Cursor.Sort != input.Sort {
			return BillingOrderPage{}, errors.New("billing: invalid cursor")
		}
	}

	query := db.WithContext(ctx).
		Model(&BillingOrder{})
	if input.IncludeUser {
		query = query.Preload("User")
	}
	if input.IncludeInvoice {
		query = query.Preload("Invoice")
	}
	if input.UserID != nil {
		query = query.Where("user_id = ?", *input.UserID)
	}
	if input.Status != nil {
		query = query.Where("status = ?", *input.Status)
	}
	if input.WithoutInvoice {
		query = query.Where("NOT EXISTS (SELECT 1 FROM billing_invoices WHERE billing_invoices.order_id = billing_orders.id)")
	}
	if input.CreatedFrom != nil {
		query = query.Where("created_at >= ?", *input.CreatedFrom)
	}
	if input.CreatedTo != nil {
		query = query.Where("created_at <= ?", *input.CreatedTo)
	}
	if input.Cursor != nil {
		operator := "<"
		if input.Sort == "asc" {
			operator = ">"
		}
		query = query.Where("(created_at, id) "+operator+" (?, ?)", input.Cursor.CreatedAt, input.Cursor.ID)
	}

	order := "created_at desc, id desc"
	if input.Sort == "asc" {
		order = "created_at asc, id asc"
	}
	var orders []BillingOrder
	if err := query.Order(order).Limit(input.Size + 1).Find(&orders).Error; err != nil {
		return BillingOrderPage{}, err
	}

	var nextCursor *BillingOrderCursor
	if len(orders) > input.Size {
		orders = orders[:input.Size]
		if len(orders) > 0 {
			last := orders[len(orders)-1]
			nextCursor = &BillingOrderCursor{
				CreatedAt: last.CreatedAt.UTC(),
				ID:        last.ID,
				Sort:      input.Sort,
			}
		}
	}
	return BillingOrderPage{Orders: orders, NextCursor: nextCursor}, nil
}

func DebitCreditsForJob(ctx context.Context, db *gorm.DB, input DebitCreditsInput) error {
	if db == nil {
		return errors.New("billing: nil db")
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return DebitCreditsForJobTx(ctx, tx, input)
	})
}

func DebitCreditsForJobTx(ctx context.Context, tx *gorm.DB, input DebitCreditsInput) error {
	if tx == nil {
		return errors.New("billing: nil db")
	}
	if input.UserID == "" {
		return errors.New("billing: user id is required")
	}
	if input.RelatedJobID == uuid.Nil {
		return errors.New("billing: related job id is required")
	}
	if input.Credits <= 0 {
		return errors.New("billing: credits must be positive")
	}
	if input.IdempotencyKey == "" {
		return errors.New("billing: idempotency key is required")
	}
	if input.Now.IsZero() {
		return errors.New("billing: time is required")
	}

	tx = tx.WithContext(ctx)
	if err := lockIdempotencyKey(tx, input.IdempotencyKey); err != nil {
		return err
	}
	var existing int64
	if err := tx.Model(&CreditLedgerEntry{}).
		Where("related_job_id = ? AND metadata @> ?::jsonb", input.RelatedJobID, debitMetadataQuery(input.IdempotencyKey)).
		Count(&existing).Error; err != nil {
		return err
	}
	if existing > 0 {
		return nil
	}

	buckets, err := lockEligibleCreditBuckets(tx, input.UserID, input.Now)
	if err != nil {
		return err
	}
	total := 0
	for _, bucket := range buckets {
		total += bucket.CreditsRemaining
	}
	if total < input.Credits {
		return ErrInsufficientCredits
	}

	remaining := input.Credits
	for _, bucket := range buckets {
		if remaining == 0 {
			break
		}
		consume := bucket.CreditsRemaining
		if consume > remaining {
			consume = remaining
		}
		if consume <= 0 {
			continue
		}
		if err := tx.Model(&CreditBucket{}).
			Where("id = ?", bucket.ID).
			Update("credits_remaining", bucket.CreditsRemaining-consume).Error; err != nil {
			return err
		}
		bucketID := bucket.ID
		entry := CreditLedgerEntry{
			UserID:         input.UserID,
			BucketID:       &bucketID,
			EntryType:      CreditLedgerEntryDebit,
			CreditsDelta:   -consume,
			RelatedJobID:   &input.RelatedJobID,
			IdempotencyKey: fmt.Sprintf("%s:%s", input.IdempotencyKey, bucket.ID.String()),
			Metadata:       debitMetadata(input.IdempotencyKey),
		}
		if err := tx.Create(&entry).Error; err != nil {
			return err
		}
		remaining -= consume
	}
	return nil
}

func RefundCreditsForJob(ctx context.Context, db *gorm.DB, userID string, jobID uuid.UUID, now time.Time) error {
	if db == nil {
		return errors.New("billing: nil db")
	}
	if userID == "" {
		return errors.New("billing: user id is required")
	}
	if jobID == uuid.Nil {
		return errors.New("billing: related job id is required")
	}
	if now.IsZero() {
		return errors.New("billing: time is required")
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockIdempotencyKey(tx, "refund_job:"+jobID.String()); err != nil {
			return err
		}
		var debits []CreditLedgerEntry
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND related_job_id = ? AND entry_type = ?", userID, jobID, CreditLedgerEntryDebit).
			Order("created_at asc, id asc").
			Find(&debits).Error; err != nil {
			return err
		}
		for _, debit := range debits {
			refundKey := "refund:" + debit.ID.String()
			existing, err := bucketForLedgerIdempotencyKey(tx, refundKey)
			if err != nil {
				return err
			}
			if existing != nil {
				continue
			}
			if debit.BucketID == nil {
				return errors.New("billing: debit ledger entry has no bucket")
			}
			credits := -debit.CreditsDelta
			if credits <= 0 {
				continue
			}
			var original CreditBucket
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&original, "id = ?", *debit.BucketID).Error; err != nil {
				return err
			}
			targetBucket := original
			if bucketIsRefundableInPlace(original, now) {
				if err := tx.Model(&CreditBucket{}).
					Where("id = ?", original.ID).
					Update("credits_remaining", gorm.Expr("credits_remaining + ?", credits)).Error; err != nil {
					return err
				}
			} else {
				var err error
				targetBucket, err = createRefundBucket(tx, original, userID, credits, now)
				if err != nil {
					return err
				}
			}
			bucketID := targetBucket.ID
			entry := CreditLedgerEntry{
				UserID:         userID,
				BucketID:       &bucketID,
				EntryType:      CreditLedgerEntryRefund,
				CreditsDelta:   credits,
				RelatedJobID:   &jobID,
				IdempotencyKey: refundKey,
			}
			if err := tx.Create(&entry).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

type grantCreditsInput struct {
	UserID         string
	SourceType     CreditSourceType
	OrderID        *uuid.UUID
	Credits        int
	ValidFrom      time.Time
	ExpiresAt      *time.Time
	EntryType      CreditLedgerEntryType
	IdempotencyKey string
}

func grantCredits(ctx context.Context, db *gorm.DB, input grantCreditsInput) (CreditBucket, error) {
	var bucket CreditBucket
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockIdempotencyKey(tx, input.IdempotencyKey); err != nil {
			return err
		}
		existing, err := bucketForLedgerIdempotencyKey(tx, input.IdempotencyKey)
		if err != nil {
			return err
		}
		if existing != nil {
			bucket = *existing
			return nil
		}

		bucket = CreditBucket{
			UserID:           input.UserID,
			SourceType:       input.SourceType,
			OrderID:          input.OrderID,
			CreditsGranted:   input.Credits,
			CreditsRemaining: input.Credits,
			ValidFrom:        input.ValidFrom,
			ExpiresAt:        input.ExpiresAt,
		}
		if err := tx.Create(&bucket).Error; err != nil {
			return err
		}

		entry := CreditLedgerEntry{
			UserID:         input.UserID,
			BucketID:       &bucket.ID,
			EntryType:      input.EntryType,
			CreditsDelta:   input.Credits,
			IdempotencyKey: input.IdempotencyKey,
		}
		return tx.Create(&entry).Error
	})
	return bucket, err
}

func BucketForLedgerIdempotencyKey(ctx context.Context, db *gorm.DB, idempotencyKey string) (*CreditBucket, error) {
	if db == nil {
		return nil, errors.New("billing: nil db")
	}
	if idempotencyKey == "" {
		return nil, errors.New("billing: idempotency key is required")
	}
	return bucketForLedgerIdempotencyKey(db.WithContext(ctx), idempotencyKey)
}

func bucketForLedgerIdempotencyKey(db *gorm.DB, idempotencyKey string) (*CreditBucket, error) {
	var entry CreditLedgerEntry
	err := db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("idempotency_key = ?", idempotencyKey).
		First(&entry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if entry.BucketID == nil {
		return nil, errors.New("billing: idempotent ledger entry has no bucket")
	}
	var bucket CreditBucket
	if err := db.First(&bucket, "id = ?", *entry.BucketID).Error; err != nil {
		return nil, err
	}
	return &bucket, nil
}

func bucketForProviderPaymentIntent(db *gorm.DB, providerPaymentIntentID string) (*CreditBucket, error) {
	var order BillingOrder
	err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("provider_payment_intent_id = ? AND status = ?", providerPaymentIntentID, OrderStatusPaid).
		First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return bucketForLedgerIdempotencyKey(db, "topup_paid:"+order.ID.String())
}

func ensureProviderPaymentIntentAvailable(db *gorm.DB, orderID uuid.UUID, providerPaymentIntentID string) error {
	return ensureProviderMetadataAvailable(db, orderID, "provider_payment_intent_id", providerPaymentIntentID)
}

func bucketForProviderCheckoutSession(db *gorm.DB, providerCheckoutSessionID string) (*CreditBucket, error) {
	var order BillingOrder
	err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("provider_checkout_session_id = ? AND status = ?", providerCheckoutSessionID, OrderStatusPaid).
		First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return bucketForLedgerIdempotencyKey(db, "topup_paid:"+order.ID.String())
}

func ensureProviderCheckoutSessionAvailable(db *gorm.DB, orderID uuid.UUID, providerCheckoutSessionID string) error {
	return ensureProviderMetadataAvailable(db, orderID, "provider_checkout_session_id", providerCheckoutSessionID)
}

func ensureProviderMetadataAvailable(db *gorm.DB, orderID uuid.UUID, column string, value string) error {
	if value == "" {
		return nil
	}
	var order BillingOrder
	err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(column+" = ?", value).
		First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if order.ID != orderID {
		return ErrProviderMetadataConflict
	}
	return nil
}

func markCreditOrderTerminal(ctx context.Context, db *gorm.DB, orderID uuid.UUID, status OrderStatus, timestampColumn string, timestamp time.Time) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order BillingOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", orderID).
			First(&order).Error; err != nil {
			return err
		}
		if order.OrderType != OrderTypeCreditTopup {
			return fmt.Errorf("billing: unsupported order type %q", order.OrderType)
		}
		if order.Status == OrderStatusPaid || order.Status == OrderStatusRefunded {
			return nil
		}
		if status == OrderStatusFailed && order.Status != OrderStatusPending && order.Status != OrderStatusFailed {
			return nil
		}
		if status == OrderStatusCanceled && order.Status != OrderStatusPending && order.Status != OrderStatusCanceled {
			return nil
		}
		return tx.Model(&BillingOrder{}).
			Where("id = ?", order.ID).
			Updates(map[string]any{
				"status":        status,
				timestampColumn: timestamp,
			}).Error
	})
}

func lockEligibleCreditBuckets(tx *gorm.DB, userID string, now time.Time) ([]CreditBucket, error) {
	var buckets []CreditBucket
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND credits_remaining > 0 AND valid_from <= ? AND voided_at IS NULL AND (expires_at IS NULL OR expires_at > ?)", userID, now, now).
		Order(`
CASE source_type
  WHEN 'signup_bonus' THEN 1
  WHEN 'topup_purchase' THEN 2
  WHEN 'refund' THEN 3
  ELSE 4
END,
expires_at ASC NULLS LAST,
created_at ASC,
id ASC`).
		Find(&buckets).Error
	return buckets, err
}

func bucketIsRefundableInPlace(bucket CreditBucket, now time.Time) bool {
	if bucket.VoidedAt != nil {
		return false
	}
	return bucket.ExpiresAt == nil || bucket.ExpiresAt.After(now)
}

func createRefundBucket(tx *gorm.DB, original CreditBucket, userID string, credits int, now time.Time) (CreditBucket, error) {
	refund := CreditBucket{
		UserID:           userID,
		SourceType:       CreditSourceRefund,
		CreditsGranted:   credits,
		CreditsRemaining: credits,
		ValidFrom:        now,
	}
	if original.ExpiresAt != nil {
		expiresAt := now.AddDate(0, 0, 7)
		if original.ExpiresAt.After(expiresAt) {
			expiresAt = *original.ExpiresAt
		}
		refund.ExpiresAt = &expiresAt
	}
	if err := tx.Create(&refund).Error; err != nil {
		return CreditBucket{}, err
	}
	return refund, nil
}

func debitMetadata(idempotencyKey string) datatypes.JSON {
	raw, _ := json.Marshal(map[string]string{"debit_idempotency_key": idempotencyKey})
	return datatypes.JSON(raw)
}

func debitMetadataQuery(idempotencyKey string) string {
	raw, _ := json.Marshal(map[string]string{"debit_idempotency_key": idempotencyKey})
	return string(raw)
}

func adjustmentLedgerEntriesExist(db *gorm.DB, idempotencyKey string) (bool, error) {
	var count int64
	if err := db.Model(&CreditLedgerEntry{}).
		Where("metadata @> ?::jsonb", adjustmentMetadataQuery(idempotencyKey)).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func adjustmentMetadata(idempotencyKey string) datatypes.JSON {
	raw, _ := json.Marshal(map[string]string{"adjustment_idempotency_key": idempotencyKey})
	return datatypes.JSON(raw)
}

func adjustmentMetadataQuery(idempotencyKey string) string {
	raw, _ := json.Marshal(map[string]string{"adjustment_idempotency_key": idempotencyKey})
	return string(raw)
}

func lockIdempotencyKey(db *gorm.DB, idempotencyKey string) error {
	if idempotencyKey == "" {
		return errors.New("billing: idempotency key is required")
	}
	return db.Exec("SELECT pg_advisory_xact_lock(hashtext(?))", idempotencyKey).Error
}

func nonEmptyString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func orderAmountUnitPrice(amountCents int) string {
	return decimal.NewFromInt(int64(amountCents)).Div(decimal.NewFromInt(100)).StringFixed(2)
}

func normalizedBillingInvoice(input CreateBillingInvoiceInput) (normalizedBillingInvoiceInput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return normalizedBillingInvoiceInput{}, errors.New("billing: user id is required")
	}
	var orderID *uuid.UUID
	if input.OrderID != nil {
		if *input.OrderID == uuid.Nil {
			return normalizedBillingInvoiceInput{}, errors.New("billing: order id is required")
		}
		value := *input.OrderID
		orderID = &value
	}
	invoiceSerie, err := normalizedInvoiceSerie(input.InvoiceSerie)
	if err != nil {
		return normalizedBillingInvoiceInput{}, err
	}
	invoiceDate := invoiceDateOnly(input.InvoiceDate)
	if len(input.Lines) == 0 {
		return normalizedBillingInvoiceInput{}, errors.New("billing: invoice lines are required")
	}

	lines := make([]BillingInvoiceLine, 0, len(input.Lines))
	netAmount := decimal.Zero
	vatAmount := decimal.Zero
	for index, inputLine := range input.Lines {
		line, lineNetAmount, lineVATAmount, err := normalizedBillingInvoiceLine(inputLine)
		if err != nil {
			return normalizedBillingInvoiceInput{}, fmt.Errorf("billing: invoice line %d: %w", index+1, err)
		}
		lines = append(lines, line)
		netAmount = netAmount.Add(lineNetAmount)
		vatAmount = vatAmount.Add(lineVATAmount)
	}
	netAmount = netAmount.Round(2)
	vatAmount = vatAmount.Round(2)
	totalAmount := netAmount.Add(vatAmount).Round(2)

	return normalizedBillingInvoiceInput{
		UserID:       userID,
		OrderID:      orderID,
		InvoiceSerie: invoiceSerie,
		InvoiceDate:  invoiceDate,
		Lines:        lines,
		NetAmount:    netAmount,
		VATAmount:    vatAmount,
		TotalAmount:  totalAmount,
	}, nil
}

func normalizedBillingInvoiceLine(input CreateBillingInvoiceLineInput) (BillingInvoiceLine, decimal.Decimal, decimal.Decimal, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return BillingInvoiceLine{}, decimal.Zero, decimal.Zero, errors.New("line name is required")
	}
	if err := validateMaxLength("line name", name, invoiceLineNameMaxLength); err != nil {
		return BillingInvoiceLine{}, decimal.Zero, decimal.Zero, err
	}
	if input.Quantity <= 0 {
		return BillingInvoiceLine{}, decimal.Zero, decimal.Zero, errors.New("quantity must be positive")
	}
	unitPrice, err := parseNonNegativeInvoiceDecimal("unit price", input.UnitPrice)
	if err != nil {
		return BillingInvoiceLine{}, decimal.Zero, decimal.Zero, err
	}
	vatPercentage, err := parseNonNegativeInvoiceDecimal("vat percentage", input.VATPercentage)
	if err != nil {
		return BillingInvoiceLine{}, decimal.Zero, decimal.Zero, err
	}
	if vatPercentage.GreaterThan(decimal.NewFromInt(100)) {
		return BillingInvoiceLine{}, decimal.Zero, decimal.Zero, errors.New("vat percentage must be between 0 and 100")
	}

	unitPrice = unitPrice.Round(2)
	vatPercentage = vatPercentage.Round(2)
	lineNetAmount := unitPrice.Mul(decimal.NewFromInt(int64(input.Quantity))).Round(2)
	lineVATAmount := lineNetAmount.Mul(vatPercentage).Div(decimal.NewFromInt(100)).Round(2)
	lineTotalAmount := lineNetAmount.Add(lineVATAmount).Round(2)

	return BillingInvoiceLine{
		Name:           name,
		Quantity:       input.Quantity,
		UnitPrice:      unitPrice.StringFixed(2),
		VATPercentage:  vatPercentage.StringFixed(2),
		TotalVATAmount: lineVATAmount.StringFixed(2),
		TotalAmount:    lineTotalAmount.StringFixed(2),
	}, lineNetAmount, lineVATAmount, nil
}

func parseNonNegativeInvoiceDecimal(field string, value string) (decimal.Decimal, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return decimal.Zero, fmt.Errorf("%s is required", field)
	}
	parsed, err := decimal.NewFromString(trimmed)
	if err != nil {
		return decimal.Zero, fmt.Errorf("%s is invalid", field)
	}
	if parsed.IsNegative() {
		return decimal.Zero, fmt.Errorf("%s must be non-negative", field)
	}
	return parsed, nil
}

func normalizedInvoiceSerie(value string) (string, error) {
	serie := strings.ToUpper(strings.TrimSpace(value))
	if serie == "" {
		return "", errors.New("billing: invoice serie is required")
	}
	if len(serie) > invoiceSerieMaxLength {
		return "", fmt.Errorf("billing: invoice serie must be %d characters or fewer", invoiceSerieMaxLength)
	}
	for _, r := range serie {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			continue
		}
		return "", errors.New("billing: invoice serie contains invalid characters")
	}
	return serie, nil
}

func invoiceDateOnly(value time.Time) time.Time {
	if value.IsZero() {
		value = time.Now().UTC()
	} else {
		value = value.UTC()
	}
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func ensureBillingInvoiceCounter(tx *gorm.DB, invoiceSerie string) error {
	counter := BillingInvoiceCounter{InvoiceSerie: invoiceSerie}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "invoice_serie"}},
		DoNothing: true,
	}).Create(&counter).Error
}

func incrementBillingInvoiceCounter(tx *gorm.DB, invoiceSerie string) (int64, error) {
	var nextNumber int64
	err := tx.Raw(`
UPDATE "billing_invoice_counters"
SET "last_number" = "last_number" + 1,
    "updated_at" = NOW()
WHERE "invoice_serie" = ?
RETURNING "last_number"`, invoiceSerie).Scan(&nextNumber).Error
	if err != nil {
		return 0, err
	}
	return nextNumber, nil
}

func billingProfileSnapshotJSON(profile BillingProfile) (datatypes.JSON, error) {
	raw, err := json.Marshal(BillingProfileSnapshot{
		ID:                 profile.ID,
		UserID:             profile.UserID,
		EntityType:         profile.EntityType,
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
	})
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(raw), nil
}

func billingUserSnapshotJSON(user auth.User) (datatypes.JSON, error) {
	raw, err := json.Marshal(billingUserSnapshot{
		Source:       "user",
		UserID:       user.ID,
		BillingName:  user.Name,
		BillingEmail: user.Email,
	})
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(raw), nil
}

func billingInvoiceLinesJSON(lines []BillingInvoiceLine) (datatypes.JSON, error) {
	raw, err := json.Marshal(lines)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(raw), nil
}

func normalizedBillingProfile(input UpsertBillingProfileInput) (BillingProfile, error) {
	profile := BillingProfile{
		UserID:             strings.TrimSpace(input.UserID),
		EntityType:         BillingEntityType(strings.TrimSpace(string(input.EntityType))),
		BillingName:        strings.TrimSpace(input.BillingName),
		BillingEmail:       strings.TrimSpace(input.BillingEmail),
		CountryCode:        strings.ToUpper(strings.TrimSpace(input.CountryCode)),
		AddressLine1:       strings.TrimSpace(input.AddressLine1),
		AddressLine2:       trimmedOptionalString(input.AddressLine2),
		City:               strings.TrimSpace(input.City),
		Region:             trimmedOptionalString(input.Region),
		PostalCode:         strings.TrimSpace(input.PostalCode),
		FiscalCode:         trimmedOptionalString(input.FiscalCode),
		RegistrationNumber: trimmedOptionalString(input.RegistrationNumber),
	}

	if profile.UserID == "" {
		return BillingProfile{}, errors.New("billing: user id is required")
	}
	if err := validateBillingEntityType(profile.EntityType); err != nil {
		return BillingProfile{}, err
	}
	if profile.BillingName == "" {
		return BillingProfile{}, errors.New("billing name is required")
	}
	if err := validateMaxLength("billing name", profile.BillingName, billingProfileNameMaxLength); err != nil {
		return BillingProfile{}, err
	}
	if err := validateMaxLength("billing email", profile.BillingEmail, billingProfileEmailMaxLength); err != nil {
		return BillingProfile{}, err
	}
	if _, err := mail.ParseAddress(profile.BillingEmail); err != nil {
		return BillingProfile{}, errors.New("billing email is invalid")
	}
	if profile.CountryCode == "" {
		return BillingProfile{}, errors.New("country is required")
	}
	if !validCountryCode(profile.CountryCode) {
		return BillingProfile{}, errors.New("country code is invalid")
	}
	if profile.AddressLine1 == "" {
		return BillingProfile{}, errors.New("address line 1 is required")
	}
	if err := validateMaxLength("address line 1", profile.AddressLine1, billingProfileAddressLineMaxLength); err != nil {
		return BillingProfile{}, err
	}
	if err := validateOptionalMaxLength("address line 2", profile.AddressLine2, billingProfileAddressLineMaxLength); err != nil {
		return BillingProfile{}, err
	}
	if profile.City == "" {
		return BillingProfile{}, errors.New("city is required")
	}
	if err := validateMaxLength("city", profile.City, billingProfileCityMaxLength); err != nil {
		return BillingProfile{}, err
	}
	if err := validateOptionalMaxLength("region", profile.Region, billingProfileRegionMaxLength); err != nil {
		return BillingProfile{}, err
	}
	if profile.PostalCode == "" {
		return BillingProfile{}, errors.New("postal code is required")
	}
	if err := validateMaxLength("postal code", profile.PostalCode, billingProfilePostalCodeMaxLength); err != nil {
		return BillingProfile{}, err
	}
	if profile.EntityType == BillingEntityCompany && profile.CountryCode == "RO" && profile.FiscalCode == nil {
		return BillingProfile{}, errors.New("fiscal code is required for Romanian companies")
	}
	if err := validateOptionalMaxLength("fiscal code", profile.FiscalCode, billingProfileFiscalCodeMaxLength); err != nil {
		return BillingProfile{}, err
	}
	if err := validateOptionalMaxLength("registration number", profile.RegistrationNumber, billingProfileRegistrationNumberMaxLength); err != nil {
		return BillingProfile{}, err
	}
	return profile, nil
}

func billingProfileUpsertAssignments(profile BillingProfile) map[string]any {
	return map[string]any{
		"entity_type":         profile.EntityType,
		"billing_name":        profile.BillingName,
		"billing_email":       profile.BillingEmail,
		"country_code":        profile.CountryCode,
		"address_line1":       profile.AddressLine1,
		"address_line2":       optionalStringUpdateValue(profile.AddressLine2),
		"city":                profile.City,
		"region":              optionalStringUpdateValue(profile.Region),
		"postal_code":         profile.PostalCode,
		"fiscal_code":         optionalStringUpdateValue(profile.FiscalCode),
		"registration_number": optionalStringUpdateValue(profile.RegistrationNumber),
		"updated_at":          gorm.Expr("NOW()"),
	}
}

func validCountryCode(value string) bool {
	if len(value) != 2 {
		return false
	}
	return strings.Contains(" "+billingProfileCountryCodes+" ", " "+value+" ")
}

func validateMaxLength(field string, value string, maxLength int) error {
	if utf8.RuneCountInString(value) > maxLength {
		return fmt.Errorf("%s must be %d characters or fewer", field, maxLength)
	}
	return nil
}

func validateOptionalMaxLength(field string, value *string, maxLength int) error {
	if value == nil {
		return nil
	}
	return validateMaxLength(field, *value, maxLength)
}

func trimmedOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func optionalStringUpdateValue(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func timePtr(value time.Time) *time.Time {
	return &value
}
