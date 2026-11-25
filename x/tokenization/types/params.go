package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		AllowedDenoms: []string{
			NativeDenom,
			// These defaults only apply to chains starting from genesis. Mainnet
			// keeps its params in state — the canonical USDC denom reaches it via
			// a governance proposal carrying tokenization.MsgUpdateParams (which
			// replaces the WHOLE params object: submit the full current params
			// with the denom appended), not via an upgrade migration.
			//
			// Canonical USDC, routed via Injective. Preferred for anything new.
			USDCDenom,
			// Legacy Noble-direct USDC. Still allowed: collections with a backed
			// path against it cannot be repointed without moving their escrow,
			// since the escrow address derives from the denom string.
			USDCNobleDenom,
			ATOMDenom,
			OSMODenom,
		},
		AffiliatePercentage: sdkmath.NewUint(500),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if !p.AffiliatePercentage.IsNil() && p.AffiliatePercentage.GT(sdkmath.NewUint(10000)) {
		return fmt.Errorf("affiliatePercentage cannot exceed 10000, got %s", p.AffiliatePercentage.String())
	}
	for _, denom := range p.AllowedDenoms {
		if err := sdk.ValidateDenom(denom); err != nil {
			return fmt.Errorf("invalid allowed denom %q: %w", denom, err)
		}
	}
	return nil
}
