package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/pot/types"
)

// RegisterInvariants registers all x/pot invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "compliance-jailed-valid", ComplianceJailedValidInvariant(k))
}

// ComplianceJailedValidInvariant checks that:
// 1. Every address in the compliance-jailed set corresponds to an existing validator.
// 2. Every compliance-jailed validator is actually disabled (jailed/power=0).
func ComplianceJailedValidInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		broken := false

		complianceJailed := k.GetAllComplianceJailed(ctx)
		for _, addrBytes := range complianceJailed {
			consAddr := sdk.ConsAddress(addrBytes)

			// Check that the validator exists.
			_, err := k.validatorSet.GetValidatorByConsAddr(ctx, consAddr)
			if err != nil {
				msg += fmt.Sprintf("\tcompliance-jailed validator %s does not exist in validator set\n", consAddr)
				broken = true
				continue
			}

			// Check that the validator is actually jailed/disabled.
			jailed, jErr := k.validatorSet.IsValidatorJailed(ctx, consAddr)
			if jErr != nil {
				msg += fmt.Sprintf("\tcompliance-jailed validator %s: error checking jailed status: %v\n",
					consAddr, jErr)
				broken = true
				continue
			}
			if !jailed {
				msg += fmt.Sprintf("\tcompliance-jailed validator %s is not disabled in validator set\n", consAddr)
				broken = true
			}
		}

		return sdk.FormatInvariant(types.ModuleName, "compliance-jailed-valid", msg), broken
	}
}
