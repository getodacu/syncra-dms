package billing

import "fmt"

type CreditPurchaseTier string

const (
	CreditPurchaseTier1 CreditPurchaseTier = "tier_1"
	CreditPurchaseTier2 CreditPurchaseTier = "tier_2"
	CreditPurchaseTier3 CreditPurchaseTier = "tier_3"
	CreditPurchaseTier4 CreditPurchaseTier = "tier_4"

	CreditPurchaseBlockSize = 1000
)

type CreditPurchaseQuote struct {
	Credits         int
	CreditBlocks    int
	Tier            CreditPurchaseTier
	UnitAmountCents int
	AmountCents     int
	Currency        string
}

func QuoteCreditPurchase(credits int) (CreditPurchaseQuote, error) {
	if credits < CreditPurchaseBlockSize {
		return CreditPurchaseQuote{}, fmt.Errorf("billing: credit purchase must be at least %d credits", CreditPurchaseBlockSize)
	}
	if credits%CreditPurchaseBlockSize != 0 {
		return CreditPurchaseQuote{}, fmt.Errorf("billing: credit purchase must be in %d-credit blocks", CreditPurchaseBlockSize)
	}

	blocks := credits / CreditPurchaseBlockSize
	tier, unitAmountCents := quoteTierAndUnitAmount(credits)

	return CreditPurchaseQuote{
		Credits:         credits,
		CreditBlocks:    blocks,
		Tier:            tier,
		UnitAmountCents: unitAmountCents,
		AmountCents:     blocks * unitAmountCents,
		Currency:        "EUR",
	}, nil
}

func quoteTierAndUnitAmount(credits int) (CreditPurchaseTier, int) {
	switch {
	case credits >= 20000:
		return CreditPurchaseTier4, 850
	case credits >= 10000:
		return CreditPurchaseTier3, 900
	case credits >= 5000:
		return CreditPurchaseTier2, 950
	default:
		return CreditPurchaseTier1, 1000
	}
}
