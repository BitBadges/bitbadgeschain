package v35

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// conversionDec is ConversionFactor as a LegacyDec, for the share fields.
var conversionDec = sdkmath.LegacyNewDecFromInt(ConversionFactor)

// StakingMigrationResult reports what RescaleStaking touched.
type StakingMigrationResult struct {
	Validators            int
	Delegations           int
	UnbondingDelegations  int
	Redelegations         int
}

// RescaleStaking multiplies every staking amount by ConversionFactor so the
// module stays consistent with the redenominated bank balances.
//
// The bonded and not-bonded pools are ordinary module-account bank balances, so
// RedenominateBank already moved them. What it cannot move is the accounting
// x/staking keeps in its own store: validator token totals, delegator shares,
// and the in-flight unbonding and redelegation records.
//
// Tokens and shares are scaled by the same factor deliberately. A validator's
// exchange rate is tokens/shares, so scaling both leaves it unchanged and every
// delegator's claim converts cleanly — scaling only tokens would silently
// reprice every delegation by 10^9.
//
// MinSelfDelegation is scaled too. It is a token amount, and leaving it at the
// old scale would turn a 1-BADGE floor into a 10^-9-BADGE floor.
//
// NOTE: this does not touch consensus power. See PowerReduction in
// app/params — power is derived as tokens/PowerReduction, so the reduction has
// to move by the same factor or every validator's power jumps by 10^9.
func RescaleStaking(ctx sdk.Context, sk *stakingkeeper.Keeper) (StakingMigrationResult, error) {
	var res StakingMigrationResult

	// Validators: token totals and the share denominator.
	vals, err := sk.GetAllValidators(ctx)
	if err != nil {
		return res, fmt.Errorf("reading validators: %w", err)
	}
	for _, v := range vals {
		v.Tokens = v.Tokens.Mul(ConversionFactor)
		v.DelegatorShares = v.DelegatorShares.Mul(conversionDec)
		if !v.MinSelfDelegation.IsNil() {
			v.MinSelfDelegation = v.MinSelfDelegation.Mul(ConversionFactor)
		}
		if err := sk.SetValidator(ctx, v); err != nil {
			return res, fmt.Errorf("setting validator %s: %w", v.OperatorAddress, err)
		}
		res.Validators++
	}

	// Delegations: shares only. The token value follows from the validator's
	// exchange rate, which is unchanged.
	dels, err := sk.GetAllDelegations(ctx)
	if err != nil {
		return res, fmt.Errorf("reading delegations: %w", err)
	}
	for _, d := range dels {
		d.Shares = d.Shares.Mul(conversionDec)
		if err := sk.SetDelegation(ctx, d); err != nil {
			return res, fmt.Errorf("setting delegation %s: %w", d.DelegatorAddress, err)
		}
		res.Delegations++
	}

	// Unbonding delegations: both balances on every entry. These are token
	// amounts already detached from any validator, so nothing else scales them.
	// Collected before writing: mutating inside the iterator is not safe.
	var ubds []stakingtypes.UnbondingDelegation
	if err := sk.IterateUnbondingDelegations(ctx, func(_ int64, ubd stakingtypes.UnbondingDelegation) bool {
		ubds = append(ubds, ubd)
		return false
	}); err != nil {
		return res, fmt.Errorf("reading unbonding delegations: %w", err)
	}
	for _, u := range ubds {
		for i := range u.Entries {
			u.Entries[i].InitialBalance = u.Entries[i].InitialBalance.Mul(ConversionFactor)
			u.Entries[i].Balance = u.Entries[i].Balance.Mul(ConversionFactor)
		}
		if err := sk.SetUnbondingDelegation(ctx, u); err != nil {
			return res, fmt.Errorf("setting unbonding delegation %s: %w", u.DelegatorAddress, err)
		}
		res.UnbondingDelegations++
	}

	// Redelegations: InitialBalance is tokens, SharesDst is shares against the
	// destination validator. Both move by the same factor for the same reason
	// tokens and shares do above.
	var reds []stakingtypes.Redelegation
	if err := sk.IterateRedelegations(ctx, func(_ int64, red stakingtypes.Redelegation) bool {
		reds = append(reds, red)
		return false
	}); err != nil {
		return res, fmt.Errorf("reading redelegations: %w", err)
	}
	for _, r := range reds {
		for i := range r.Entries {
			r.Entries[i].InitialBalance = r.Entries[i].InitialBalance.Mul(ConversionFactor)
			r.Entries[i].SharesDst = r.Entries[i].SharesDst.Mul(conversionDec)
		}
		if err := sk.SetRedelegation(ctx, r); err != nil {
			return res, fmt.Errorf("setting redelegation %s: %w", r.DelegatorAddress, err)
		}
		res.Redelegations++
	}

	ctx.Logger().Info(
		"v35: rescaled staking",
		"factor", ConversionFactor.String(),
		"validators", res.Validators,
		"delegations", res.Delegations,
		"unbonding_delegations", res.UnbondingDelegations,
		"redelegations", res.Redelegations,
	)

	return res, nil
}

// assert the staking types we rely on still carry the fields this migration
// scales; a field rename upstream should fail the build, not silently skip.
var (
	_ = stakingtypes.Validator{}.Tokens
	_ = stakingtypes.Delegation{}.Shares
)
