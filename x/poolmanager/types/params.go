package types

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/osmomath"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	ZeroDec = osmomath.ZeroDec()
	OneDec  = osmomath.OneDec()
)

func validateDenomPairTakerFees(pairs []DenomPairTakerFee) error {
	if len(pairs) == 0 {
		return fmt.Errorf("Empty denom pair taker fee")
	}

	for _, record := range pairs {
		if record.TokenInDenom == record.TokenOutDenom {
			return fmt.Errorf("TokenInDenom and TokenOutDenom must be different")
		}

		if sdk.ValidateDenom(record.TokenInDenom) != nil {
			return fmt.Errorf("TokenInDenom is invalid: %s", sdk.ValidateDenom(record.TokenInDenom))
		}

		if sdk.ValidateDenom(record.TokenOutDenom) != nil {
			return fmt.Errorf("TokenOutDenom is invalid: %s", sdk.ValidateDenom(record.TokenOutDenom))
		}

		takerFee := record.TakerFee
		if takerFee.IsNegative() || takerFee.GTE(OneDec) {
			return fmt.Errorf("taker fee must be between 0 and 1: %s", takerFee.String())
		}
	}
	return nil
}
