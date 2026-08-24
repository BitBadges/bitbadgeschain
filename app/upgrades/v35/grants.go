package v35

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	transfertypes "github.com/cosmos/ibc-go/v11/modules/apps/transfer/types"
)

// GrantMigrationResult reports what RescaleGrants touched.
type GrantMigrationResult struct {
	FeeAllowances int
	Authorizations int
}

// RescaleGrants converts the spend limits recorded on fee grants and
// authorizations.
//
// These are permissions with a number attached, and the number is a coin
// amount. Left at the old scale a limit of "1 BADGE" becomes a limit of
// 10^-9 BADGE — the grant still exists and still looks right in a query, but
// every use of it fails for insufficient allowance.
//
// The failure is quiet and delayed: nothing breaks at upgrade time, and the
// first symptom is a user's fee grant or authz grant mysteriously refusing a
// transaction that should be within limits.
func RescaleGrants(
	ctx sdk.Context,
	fk feegrantkeeper.Keeper,
	azk authzkeeper.Keeper,
) (GrantMigrationResult, error) {
	var res GrantMigrationResult

	// --- x/feegrant ---
	type feeGrant struct {
		granter, grantee sdk.AccAddress
		allowance        feegrant.FeeAllowanceI
	}
	var feeGrants []feeGrant

	if err := fk.IterateAllFeeAllowances(ctx, func(grant feegrant.Grant) bool {
		allowance, err := grant.GetGrant()
		if err != nil {
			return false
		}
		granter, err := sdk.AccAddressFromBech32(grant.Granter)
		if err != nil {
			return false
		}
		grantee, err := sdk.AccAddressFromBech32(grant.Grantee)
		if err != nil {
			return false
		}
		feeGrants = append(feeGrants, feeGrant{granter: granter, grantee: grantee, allowance: allowance})
		return false
	}); err != nil {
		return res, fmt.Errorf("reading fee allowances: %w", err)
	}

	for _, g := range feeGrants {
		scaled, changed := scaleAllowance(g.allowance)
		if !changed {
			continue
		}
		if err := fk.UpdateAllowance(ctx, g.granter, g.grantee, scaled); err != nil {
			return res, fmt.Errorf("updating fee allowance %s->%s: %w", g.granter, g.grantee, err)
		}
		res.FeeAllowances++
	}

	// --- x/authz ---
	type grantEntry struct {
		granter, grantee sdk.AccAddress
		grant            authz.Grant
	}
	var grants []grantEntry

	azk.IterateGrants(ctx, func(granter, grantee sdk.AccAddress, grant authz.Grant) bool {
		grants = append(grants, grantEntry{granter: granter, grantee: grantee, grant: grant})
		return false
	})

	for _, g := range grants {
		auth, err := g.grant.GetAuthorization()
		if err != nil {
			// A grant whose authorization cannot be decoded is already broken;
			// do not fail the whole upgrade over it, but do not pretend it was
			// migrated either.
			ctx.Logger().Error("v35: skipping undecodable authz grant",
				"granter", g.granter.String(), "grantee", g.grantee.String(), "error", err)
			continue
		}

		var (
			scaled  authz.Authorization
			changed bool
		)
		switch a := auth.(type) {
		case *banktypes.SendAuthorization:
			if converted, did := convertCoinsChanged(a.SpendLimit); did {
				a.SpendLimit = converted
				scaled, changed = a, true
			}
		case *transfertypes.TransferAuthorization:
			for i := range a.Allocations {
				if converted, did := convertCoinsChanged(a.Allocations[i].SpendLimit); did {
					a.Allocations[i].SpendLimit = converted
					changed = true
				}
			}
			if changed {
				scaled = a
			}
		default:
			// Other authorization types carry no coin amount.
			continue
		}

		if !changed {
			continue
		}

		var expiration = g.grant.Expiration
		if err := azk.SaveGrant(ctx, g.grantee, g.granter, scaled, expiration); err != nil {
			return res, fmt.Errorf("saving authz grant %s->%s: %w", g.granter, g.grantee, err)
		}
		res.Authorizations++
	}

	ctx.Logger().Info(
		"v35: rescaled grants",
		"factor", ConversionFactor.String(),
		"fee_allowances", res.FeeAllowances,
		"authorizations", res.Authorizations,
	)
	return res, nil
}

// scaleAllowance rewrites the coin amounts inside a fee allowance, reporting
// whether anything changed. Allowance types nest, so the periodic case has to
// recurse into its basic allowance.
func scaleAllowance(a feegrant.FeeAllowanceI) (feegrant.FeeAllowanceI, bool) {
	switch al := a.(type) {
	case *feegrant.BasicAllowance:
		converted, changed := convertCoinsChanged(al.SpendLimit)
		if !changed {
			return a, false
		}
		al.SpendLimit = converted
		return al, true

	case *feegrant.PeriodicAllowance:
		changed := false
		if converted, did := convertCoinsChanged(al.Basic.SpendLimit); did {
			al.Basic.SpendLimit = converted
			changed = true
		}
		if converted, did := convertCoinsChanged(al.PeriodSpendLimit); did {
			al.PeriodSpendLimit = converted
			changed = true
		}
		if converted, did := convertCoinsChanged(al.PeriodCanSpend); did {
			al.PeriodCanSpend = converted
			changed = true
		}
		if !changed {
			return a, false
		}
		return al, true

	case *feegrant.AllowedMsgAllowance:
		inner, err := al.GetAllowance()
		if err != nil {
			return a, false
		}
		scaled, changed := scaleAllowance(inner)
		if !changed {
			return a, false
		}
		if err := al.SetAllowance(scaled); err != nil {
			return a, false
		}
		return al, true
	}
	return a, false
}

// convertCoinsChanged is convertCoins plus a flag for whether the legacy denom
// was actually present, so callers can skip writes that would be no-ops.
func convertCoinsChanged(coins sdk.Coins) (sdk.Coins, bool) {
	found := false
	for _, c := range coins {
		if c.Denom == legacyDenom() {
			found = true
			break
		}
	}
	if !found {
		return coins, false
	}
	return convertCoins(coins, legacyDenom(), newDenom()), true
}
