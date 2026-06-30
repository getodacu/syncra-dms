package billing

import "testing"

func TestCreditPurchaseQuoteBoundaries(t *testing.T) {
	tests := []struct {
		name           string
		credits        int
		wantTier       CreditPurchaseTier
		wantUnitAmount int
		wantAmount     int
	}{
		{name: "1000 credits use tier 1", credits: 1000, wantTier: CreditPurchaseTier1, wantUnitAmount: 1000, wantAmount: 1000},
		{name: "4000 credits use tier 1", credits: 4000, wantTier: CreditPurchaseTier1, wantUnitAmount: 1000, wantAmount: 4000},
		{name: "5000 credits use tier 2", credits: 5000, wantTier: CreditPurchaseTier2, wantUnitAmount: 950, wantAmount: 4750},
		{name: "9000 credits use tier 2", credits: 9000, wantTier: CreditPurchaseTier2, wantUnitAmount: 950, wantAmount: 8550},
		{name: "10000 credits use tier 3", credits: 10000, wantTier: CreditPurchaseTier3, wantUnitAmount: 900, wantAmount: 9000},
		{name: "19000 credits use tier 3", credits: 19000, wantTier: CreditPurchaseTier3, wantUnitAmount: 900, wantAmount: 17100},
		{name: "20000 credits use tier 4", credits: 20000, wantTier: CreditPurchaseTier4, wantUnitAmount: 850, wantAmount: 17000},
		{name: "50000 credits use tier 4", credits: 50000, wantTier: CreditPurchaseTier4, wantUnitAmount: 850, wantAmount: 42500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quote, err := QuoteCreditPurchase(tt.credits)
			if err != nil {
				t.Fatalf("QuoteCreditPurchase(%d) error = %v", tt.credits, err)
			}

			if quote.Credits != tt.credits {
				t.Fatalf("Credits = %d, want %d", quote.Credits, tt.credits)
			}
			if quote.CreditBlocks != tt.credits/CreditPurchaseBlockSize {
				t.Fatalf("CreditBlocks = %d, want %d", quote.CreditBlocks, tt.credits/CreditPurchaseBlockSize)
			}
			if quote.Tier != tt.wantTier {
				t.Fatalf("Tier = %q, want %q", quote.Tier, tt.wantTier)
			}
			if quote.UnitAmountCents != tt.wantUnitAmount {
				t.Fatalf("UnitAmountCents = %d, want %d", quote.UnitAmountCents, tt.wantUnitAmount)
			}
			if quote.AmountCents != tt.wantAmount {
				t.Fatalf("AmountCents = %d, want %d", quote.AmountCents, tt.wantAmount)
			}
			if quote.Currency != "EUR" {
				t.Fatalf("Currency = %q, want %q", quote.Currency, "EUR")
			}
		})
	}
}

func TestCreditPurchaseQuoteRejectsInvalidCreditAmounts(t *testing.T) {
	for _, credits := range []int{0, -1000, 1, 999, 1001, 1500} {
		t.Run("invalid credits", func(t *testing.T) {
			if _, err := QuoteCreditPurchase(credits); err == nil {
				t.Fatalf("QuoteCreditPurchase(%d) error = nil, want error", credits)
			}
		})
	}
}
